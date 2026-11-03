package telemetry

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ctxKey string

const (
	traceIDKey   ctxKey = "trace_id"
	ddTraceIDKey ctxKey = "dd_trace_id"
)

// TraceIDFromContext devuelve el trace_id activo en el contexto o "no-trace" si no existe.
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return "no-trace"
	}
	if v, ok := ctx.Value(traceIDKey).(string); ok && v != "" {
		return v
	}
	return "no-trace"
}

// DDTraceIDFromContext devuelve el ID numérico de 64 bits compatible con Datadog APM.
func DDTraceIDFromContext(ctx context.Context) uint64 {
	if ctx == nil {
		return 0
	}
	if v, ok := ctx.Value(ddTraceIDKey).(uint64); ok {
		return v
	}
	return 0
}

// WithTraceID inyecta el trace_id y dd_trace_id en el context.Context.
func WithTraceID(ctx context.Context, traceID string, ddTraceID uint64) context.Context {
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	ctx = context.WithValue(ctx, ddTraceIDKey, ddTraceID)
	return ctx
}

// GenerateTraceIDs genera un trace_id legible ("gdo-...") y un uint64 para correlación en Datadog.
func GenerateTraceIDs(incomingHeader string) (string, uint64) {
	var b [8]byte
	_, _ = rand.Read(b[:])
	ddID := binary.BigEndian.Uint64(b[:]) & 0x7FFFFFFFFFFFFFFF

	clean := strings.TrimSpace(incomingHeader)
	if clean != "" {
		return clean, ddID
	}
	return fmt.Sprintf("gdo-%d-%s", time.Now().UnixMilli(), hex.EncodeToString(b[:4])), ddID
}

// LogAction registra una acción de negocio o de base de datos incluyendo siempre el trace_id y atributos clave.
func LogAction(ctx context.Context, action string, details map[string]any) {
	traceID := TraceIDFromContext(ctx)
	ddTraceID := DDTraceIDFromContext(ctx)

	if strings.EqualFold(os.Getenv("LOG_FORMAT"), "json") {
		payload := map[string]any{
			"timestamp":   time.Now().Format(time.RFC3339Nano),
			"level":       "INFO",
			"service":     getEnvOrDefault("DD_SERVICE", "guia-de-oficios-api"),
			"env":         getEnvOrDefault("DD_ENV", "mvp"),
			"trace_id":    traceID,
			"dd.trace_id": strconv.FormatUint(ddTraceID, 10),
			"action":      action,
		}
		for k, v := range details {
			payload[k] = v
		}
		b, _ := json.Marshal(payload)
		log.Println(string(b))
		return
	}

	var kvPairs []string
	for k, v := range details {
		kvPairs = append(kvPairs, fmt.Sprintf("%s=%v", k, v))
	}
	log.Printf("🔹 [TRACE:%s] [ACTION:%s] %s", traceID, action, strings.Join(kvPairs, " | "))
}

// LogError registra un error de negocio o de base de datos con el trace_id correspondiente.
func LogError(ctx context.Context, action string, err error, details map[string]any) {
	traceID := TraceIDFromContext(ctx)
	var kvPairs []string
	for k, v := range details {
		kvPairs = append(kvPairs, fmt.Sprintf("%s=%v", k, v))
	}
	log.Printf("❌ [TRACE:%s] [ERROR:%s] err=%v | %s", traceID, action, err, strings.Join(kvPairs, " | "))
}

// =============================================================================
// Datadog Key Metrics Collector (DogStatsD UDP + Snapshot en memoria)
// Solo instrumenta las métricas clave del negocio y Golden Signals de la API.
// =============================================================================

type MetricsCollector struct {
	serviceName string
	env         string
	udpAddr     string
	conn        net.Conn

	// Contadores clave en memoria para inspección rápida vía /api/v1/metrics
	httpRequestsTotal       atomic.Int64
	httpErrorsTotal         atomic.Int64
	postgisSearchesTotal    atomic.Int64
	postgisZeroResultsTotal atomic.Int64
	serviceRequestsCreated  atomic.Int64
	urgentRequestsCreated   atomic.Int64
	reviewsSubmittedTotal   atomic.Int64
	authSuccessTotal        atomic.Int64
	authFailuresTotal       atomic.Int64
	availabilityToggles     atomic.Int64

	mu             sync.RWMutex
	lastSearchMs   float64
	lastHttpMs     float64
	byCategoryReqs map[int]int64
}

var defaultCollector *MetricsCollector
var once sync.Once

// InitMetrics inicializa el emisor DogStatsD para Datadog (ej. localhost:8125).
func GetMetrics() *MetricsCollector {
	once.Do(func() {
		host := getEnvOrDefault("DD_AGENT_HOST", "127.0.0.1")
		port := getEnvOrDefault("DD_DOGSTATSD_PORT", "8125")
		addr := net.JoinHostPort(host, port)
		conn, _ := net.DialTimeout("udp", addr, 500*time.Millisecond)

		defaultCollector = &MetricsCollector{
			serviceName:    getEnvOrDefault("DD_SERVICE", "guia-de-oficios-api"),
			env:            getEnvOrDefault("DD_ENV", "mvp"),
			udpAddr:        addr,
			conn:           conn,
			byCategoryReqs: make(map[int]int64),
		}
	})
	return defaultCollector
}

// sendDogStatsD envía una métrica en formato DogStatsD (`metric.name:value|type|#tags`) por UDP sin bloquear.
func (m *MetricsCollector) sendDogStatsD(name string, value string, metricType string, tags ...string) {
	if m.conn == nil {
		return
	}
	baseTags := []string{
		"service:" + m.serviceName,
		"env:" + m.env,
	}
	allTags := append(baseTags, tags...)
	payload := fmt.Sprintf("%s:%s|%s|#%s", name, value, metricType, strings.Join(allTags, ","))
	_, _ = m.conn.Write([]byte(payload))
}

// 1. Golden Signal HTTP: Latencia y Errores por endpoint
func (m *MetricsCollector) RecordHTTPRequest(method, route string, statusCode int, duration time.Duration) {
	ms := float64(duration.Microseconds()) / 1000.0
	m.httpRequestsTotal.Add(1)
	m.mu.Lock()
	m.lastHttpMs = ms
	m.mu.Unlock()

	statusClass := fmt.Sprintf("%dxx", statusCode/100)
	tags := []string{
		"method:" + method,
		"route:" + route,
		"status_class:" + statusClass,
		fmt.Sprintf("status_code:%d", statusCode),
	}
	m.sendDogStatsD("guia_oficios.http.request.duration_ms", fmt.Sprintf("%.2f", ms), "ms", tags...)
	m.sendDogStatsD("guia_oficios.http.request.hits", "1", "c", tags...)

	if statusCode >= 400 {
		m.httpErrorsTotal.Add(1)
		m.sendDogStatsD("guia_oficios.http.request.errors", "1", "c", tags...)
	}
}

// 2. Métrica Clave PostGIS: Latencia de búsqueda geoespacial y búsquedas sin resultados (demanda insatisfecha)
func (m *MetricsCollector) RecordPostGISSearch(categoryID int, onlyUrgent, onlyLicensed bool, resultsCount int, duration time.Duration) {
	ms := float64(duration.Microseconds()) / 1000.0
	m.postgisSearchesTotal.Add(1)
	m.mu.Lock()
	m.lastSearchMs = ms
	m.mu.Unlock()

	tags := []string{
		fmt.Sprintf("category_id:%d", categoryID),
		fmt.Sprintf("only_urgent:%t", onlyUrgent),
		fmt.Sprintf("only_licensed:%t", onlyLicensed),
	}
	m.sendDogStatsD("guia_oficios.search.postgis.duration_ms", fmt.Sprintf("%.2f", ms), "ms", tags...)
	m.sendDogStatsD("guia_oficios.search.results_count", fmt.Sprintf("%d", resultsCount), "h", tags...)

	if resultsCount == 0 {
		m.postgisZeroResultsTotal.Add(1)
		m.sendDogStatsD("guia_oficios.search.zero_results", "1", "c", tags...)
	}
}

// 3. KPI de Conversión: Solicitudes de servicio creadas (por rubro, urgencia y canal)
func (m *MetricsCollector) RecordServiceRequestCreated(categoryID int, urgency, contactChannel string) {
	m.serviceRequestsCreated.Add(1)
	if urgency == "urgente" {
		m.urgentRequestsCreated.Add(1)
	}
	m.mu.Lock()
	m.byCategoryReqs[categoryID]++
	m.mu.Unlock()

	tags := []string{
		fmt.Sprintf("category_id:%d", categoryID),
		"urgency:" + urgency,
		"channel:" + contactChannel,
	}
	m.sendDogStatsD("guia_oficios.funnel.service_requests", "1", "c", tags...)
}

// 4. KPI de Calidad Vecinal: Calificaciones publicadas y puntaje otorgado
func (m *MetricsCollector) RecordReviewSubmitted(professionalID string, rating int) {
	m.reviewsSubmittedTotal.Add(1)
	tags := []string{
		"professional_id:" + professionalID,
		fmt.Sprintf("rating:%d", rating),
	}
	m.sendDogStatsD("guia_oficios.funnel.reviews_submitted", "1", "c", tags...)
	m.sendDogStatsD("guia_oficios.funnel.review_rating", fmt.Sprintf("%d", rating), "h", tags...)
}

// 5. KPI de Adquisición y Seguridad: Registro/Login/Google Auth exitosos o fallidos
func (m *MetricsCollector) RecordAuthEvent(action, role string, success bool) {
	outcome := "success"
	if success {
		m.authSuccessTotal.Add(1)
	} else {
		outcome = "failure"
		m.authFailuresTotal.Add(1)
	}
	tags := []string{
		"action:" + action,
		"role:" + role,
		"outcome:" + outcome,
	}
	m.sendDogStatsD("guia_oficios.auth.events", "1", "c", tags...)
}

// 6. KPI de Oferta Activa: Cambio de disponibilidad de profesionales
func (m *MetricsCollector) RecordAvailabilityToggle(isAvailable bool) {
	m.availabilityToggles.Add(1)
	tags := []string{
		fmt.Sprintf("available:%t", isAvailable),
	}
	m.sendDogStatsD("guia_oficios.professionals.availability_toggles", "1", "c", tags...)
}

// Snapshot devuelve un resumen JSON de las métricas clave recolectadas.
func (m *MetricsCollector) Snapshot() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catCopy := make(map[string]int64, len(m.byCategoryReqs))
	for k, v := range m.byCategoryReqs {
		catCopy[fmt.Sprintf("category_%d", k)] = v
	}

	return map[string]any{
		"datadog_dogstatsd_target": m.udpAddr,
		"service":                  m.serviceName,
		"env":                      m.env,
		"golden_signals": map[string]any{
			"http_requests_total":     m.httpRequestsTotal.Load(),
			"http_errors_total":       m.httpErrorsTotal.Load(),
			"last_http_duration_ms":   m.lastHttpMs,
			"postgis_searches_total":  m.postgisSearchesTotal.Load(),
			"postgis_zero_results":    m.postgisZeroResultsTotal.Load(),
			"last_postgis_search_ms":  m.lastSearchMs,
		},
		"business_kpis": map[string]any{
			"service_requests_created":   m.serviceRequestsCreated.Load(),
			"urgent_requests_created":    m.urgentRequestsCreated.Load(),
			"requests_by_category":       catCopy,
			"reviews_submitted_total":    m.reviewsSubmittedTotal.Load(),
			"auth_success_total":         m.authSuccessTotal.Load(),
			"auth_failures_total":        m.authFailuresTotal.Load(),
			"availability_toggles_total": m.availabilityToggles.Load(),
		},
	}
}

// =============================================================================
// Middleware HTTP de Trazabilidad (trace_id) + Logging de Entrada/Salida + Datadog
// =============================================================================

type statusRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}

func ObservabilityMiddleware(next http.Handler) http.Handler {
	metrics := GetMetrics()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Extraer o generar Trace ID y correlacionarlo con Datadog
		incomingTrace := r.Header.Get("X-Trace-ID")
		if incomingTrace == "" {
			incomingTrace = r.Header.Get("X-Request-ID")
		}
		if incomingTrace == "" {
			incomingTrace = r.Header.Get("x-datadog-trace-id")
		}
		traceID, ddTraceID := GenerateTraceIDs(incomingTrace)

		ctx := WithTraceID(r.Context(), traceID, ddTraceID)
		r = r.WithContext(ctx)

		// Devolver siempre los headers de trazabilidad al cliente (Android / Web)
		w.Header().Set("X-Trace-ID", traceID)
		w.Header().Set("X-Datadog-Trace-ID", strconv.FormatUint(ddTraceID, 10))

		query := ""
		if r.URL.RawQuery != "" {
			query = "?" + r.URL.RawQuery
		}

		// Loguear el llamado entrante
		log.Printf("➡️  [TRACE:%s] --> %s %s%s (remote=%s)", traceID, r.Method, r.URL.Path, query, r.RemoteAddr)

		rec := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		ms := float64(duration.Microseconds()) / 1000.0

		// Loguear el resultado del llamado
		log.Printf("⬅️  [TRACE:%s] <-- %s %s | status=%d | duration=%.2fms | bytes=%d",
			traceID, r.Method, r.URL.Path, rec.statusCode, ms, rec.bytesWritten)

		// Emitir métrica a Datadog si es una ruta de la API
		if strings.HasPrefix(r.URL.Path, "/api/") {
			routeLabel := normalizeRoutePattern(r.URL.Path)
			metrics.RecordHTTPRequest(r.Method, routeLabel, rec.statusCode, duration)
		}
	})
}

func normalizeRoutePattern(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 4 && parts[0] == "api" && parts[1] == "v1" {
		switch parts[2] {
		case "professionals":
			if len(parts) == 4 {
				return "/api/v1/professionals/:id"
			}
			if len(parts) == 5 {
				return "/api/v1/professionals/:id/" + parts[4]
			}
		case "requests":
			if len(parts) == 4 && parts[3] != "my" {
				return "/api/v1/requests/:id"
			}
			if len(parts) == 5 {
				return "/api/v1/requests/:id/" + parts[4]
			}
		}
	}
	return path
}

func getEnvOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
