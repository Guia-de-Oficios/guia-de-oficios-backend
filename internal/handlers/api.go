package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"guiaoficios/backend/internal/auth"
	"guiaoficios/backend/internal/models"
	"guiaoficios/backend/internal/repository"
	"guiaoficios/backend/internal/telemetry"
)

type contextKey string

const claimsContextKey contextKey = "jwt_claims"

type APIHandler struct {
	store     repository.Store
	mode      string
	uploadDir string
	metrics   *telemetry.MetricsCollector
}

func NewAPIHandler(store repository.Store, mode string, uploadDir string) *APIHandler {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	_ = os.MkdirAll(uploadDir, 0o755)
	return &APIHandler{
		store:     store,
		mode:      mode,
		uploadDir: uploadDir,
		metrics:   telemetry.GetMetrics(),
	}
}

func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	// Health check real contra PostgreSQL + PostGIS (GDO-8) y Métricas Clave Datadog
	mux.HandleFunc("GET /api/v1/health", h.handleHealth)
	mux.HandleFunc("GET /api/v1/metrics", h.handleMetrics)

	// Archivos estáticos de fotos de problemas, DNI y matrículas (GDO-10)
	fileServer := http.FileServer(http.Dir(h.uploadDir))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", fileServer))
	mux.HandleFunc("POST /api/v1/uploads", h.withOptionalAuth(h.handleUploadFile))

	// Pantallas 1 & 2: Bienvenida, Selección de Rol, Registro, Login y Google Auth (GDO-11, GDO-12, GDO-13)
	mux.HandleFunc("POST /api/v1/auth/register", h.handleRegisterUser)
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLoginUser)
	mux.HandleFunc("POST /api/v1/auth/google", h.handleGoogleAuth)
	mux.HandleFunc("GET /api/v1/auth/me", h.withOptionalAuth(h.handleGetCurrentUser))
	mux.HandleFunc("PATCH /api/v1/users/me/address", h.withOptionalAuth(h.handleUpdateUserAddress))

	// Pantalla 3: Cliente Inicio, Categorías y Asistente "¿No sabés qué profesional necesitás?" (GDO-14)
	mux.HandleFunc("GET /api/v1/categories", h.handleListCategories)
	mux.HandleFunc("POST /api/v1/diagnostic", h.handleDiagnosticSuggestion)

	// Pantallas 4 & 5: Resultados de Profesionales con PostGIS y Perfil Detallado (GDO-15, GDO-16)
	mux.HandleFunc("GET /api/v1/professionals", h.handleListProfessionals)
	mux.HandleFunc("GET /api/v1/professionals/{id}", h.handleGetProfessional)

	// Pantallas 6, 7 & 8: Solicitud Guiada, Confirmación de Contacto y Calificación (GDO-16, GDO-17)
	mux.HandleFunc("POST /api/v1/requests", h.withOptionalAuth(h.handleCreateRequest))
	mux.HandleFunc("GET /api/v1/requests", h.withOptionalAuth(h.handleListRequests))
	mux.HandleFunc("GET /api/v1/requests/my", h.withOptionalAuth(h.handleListMyClientRequests))
	mux.HandleFunc("GET /api/v1/requests/{id}", h.withOptionalAuth(h.handleGetRequest))
	mux.HandleFunc("POST /api/v1/requests/{id}/reviews", h.withOptionalAuth(h.handleCreateReview))

	// Pantallas 9 & 10: Profesional Registro/Verificación DNI + Matrícula y Panel de Solicitudes (GDO-18)
	mux.HandleFunc("POST /api/v1/professionals/{id}/verification", h.withOptionalAuth(h.handleUpdateVerification))
	mux.HandleFunc("PATCH /api/v1/professionals/{id}/availability", h.withOptionalAuth(h.handleUpdateAvailability))
	mux.HandleFunc("GET /api/v1/professionals/{id}/dashboard", h.withOptionalAuth(h.handleGetProfessionalDashboard))
}

func (h *APIHandler) withOptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			writeJSON(w, http.StatusNoContent, nil)
			return
		}
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token := strings.TrimSpace(authHeader[7:])
			if claims, err := auth.ParseAndVerifyToken(token); err == nil {
				ctx := context.WithValue(r.Context(), claimsContextKey, claims)
				r = r.WithContext(ctx)
			}
		}
		next(w, r)
	}
}

func getClaims(r *http.Request) *auth.Claims {
	if val := r.Context().Value(claimsContextKey); val != nil {
		if claims, ok := val.(*auth.Claims); ok {
			return claims
		}
	}
	return nil
}

func (h *APIHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	dbInfo, err := h.store.Ping(r.Context())
	dbStatus := dbInfo
	httpStatus := http.StatusOK
	if err != nil {
		dbStatus = "error: " + err.Error()
		httpStatus = http.StatusServiceUnavailable
		telemetry.LogError(r.Context(), "db.ping", err, nil)
	} else {
		telemetry.LogAction(r.Context(), "db.ping", map[string]any{"engine": dbInfo})
	}
	writeJSON(w, httpStatus, map[string]any{
		"status":   "ok",
		"service":  "Guía de Oficios API (PostgreSQL + PostGIS)",
		"storage":  h.mode,
		"database": dbStatus,
		"trace_id": telemetry.TraceIDFromContext(r.Context()),
	})
}

func (h *APIHandler) handleMetrics(w http.ResponseWriter, r *http.Request) {
	telemetry.LogAction(r.Context(), "metrics.snapshot", map[string]any{"provider": "datadog_dogstatsd"})
	writeJSON(w, http.StatusOK, h.metrics.Snapshot())
}

func (h *APIHandler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Datos de registro inválidos")
		return
	}
	if strings.TrimSpace(input.FullName) == "" || strings.TrimSpace(input.Email) == "" {
		writeError(w, http.StatusBadRequest, "Nombre completo y correo electrónico son obligatorios")
		return
	}
	if input.AuthProvider != "google" && len(strings.TrimSpace(input.Password)) < 6 {
		writeError(w, http.StatusBadRequest, "La contraseña debe tener al menos 6 números o letras")
		return
	}
	user, err := h.store.RegisterUser(r.Context(), input)
	if err != nil {
		h.metrics.RecordAuthEvent("register", input.Role, false)
		telemetry.LogError(r.Context(), "auth.register", err, map[string]any{"email": input.Email, "role": input.Role})
		writeError(w, http.StatusBadRequest, "No se pudo completar el registro: "+err.Error())
		return
	}
	h.metrics.RecordAuthEvent("register", user.Role, true)
	telemetry.LogAction(r.Context(), "auth.register", map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})
	token, _ := auth.GenerateToken(user.ID, user.Email, user.Role, user.ProfessionalID)
	writeJSON(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *APIHandler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var input models.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Formato de inicio de sesión inválido")
		return
	}
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		writeError(w, http.StatusBadRequest, "Ingresá tu correo electrónico o teléfono y tu contraseña")
		return
	}
	user, err := h.store.AuthenticateUser(r.Context(), input.Email, input.Password)
	if err != nil {
		h.metrics.RecordAuthEvent("login", "unknown", false)
		telemetry.LogError(r.Context(), "auth.login", err, map[string]any{"identifier": input.Email})
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	h.metrics.RecordAuthEvent("login", user.Role, true)
	telemetry.LogAction(r.Context(), "auth.login", map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})
	token, _ := auth.GenerateToken(user.ID, user.Email, user.Role, user.ProfessionalID)
	writeJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *APIHandler) handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	var input models.GoogleAuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Datos de Google inválidos")
		return
	}
	user, err := h.store.AuthenticateGoogle(r.Context(), input)
	if err != nil {
		h.metrics.RecordAuthEvent("google", input.Role, false)
		telemetry.LogError(r.Context(), "auth.google", err, map[string]any{"email": input.Email})
		writeError(w, http.StatusBadRequest, "Error en autenticación con Google: "+err.Error())
		return
	}
	h.metrics.RecordAuthEvent("google", user.Role, true)
	telemetry.LogAction(r.Context(), "auth.google", map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})
	token, _ := auth.GenerateToken(user.ID, user.Email, user.Role, user.ProfessionalID)
	writeJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *APIHandler) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := getClaims(r)
	if claims == nil || claims.UserID == "" {
		writeError(w, http.StatusUnauthorized, "Sesión no iniciada o token inválido")
		return
	}
	user, err := h.store.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	telemetry.LogAction(r.Context(), "users.get_me", map[string]any{"user_id": user.ID, "source": "jwt"})
	writeJSON(w, http.StatusOK, user)
}

func (h *APIHandler) handleUpdateUserAddress(w http.ResponseWriter, r *http.Request) {
	claims := getClaims(r)
	if claims == nil || claims.UserID == "" {
		writeError(w, http.StatusUnauthorized, "Debés iniciar sesión para actualizar tu dirección")
		return
	}
	var input models.UpdateAddressInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Datos de dirección inválidos")
		return
	}
	if strings.TrimSpace(input.SavedAddress) == "" {
		writeError(w, http.StatusBadRequest, "La dirección no puede estar vacía")
		return
	}
	userID := claims.UserID
	user, err := h.store.UpdateUserAddress(r.Context(), userID, input)
	if err != nil {
		telemetry.LogError(r.Context(), "users.update_address", err, map[string]any{"user_id": userID})
		writeError(w, http.StatusInternalServerError, "No se pudo actualizar la dirección: "+err.Error())
		return
	}
	telemetry.LogAction(r.Context(), "users.update_address", map[string]any{
		"user_id": userID,
		"address": input.SavedAddress,
		"lat":     input.SavedLat,
		"lng":     input.SavedLng,
	})
	writeJSON(w, http.StatusOK, user)
}

func (h *APIHandler) handleUploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "El archivo supera los 10 MB o el formato multipart es inválido")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Campo 'file' requerido en multipart/form-data")
		return
	}
	defer file.Close()

	category := strings.TrimSpace(r.FormValue("category"))
	if category == "" {
		category = strings.TrimSpace(r.FormValue("kind"))
	}
	if category == "" {
		category = "attachment"
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".pdf":  true,
	}
	if !allowedExts[ext] {
		ext = ".jpg"
	}

	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	safeFilename := fmt.Sprintf("%s_%d_%s%s", category, time.Now().Unix(), hex.EncodeToString(randBytes), ext)
	destPath := filepath.Join(h.uploadDir, safeFilename)

	out, err := os.Create(destPath)
	if err != nil {
		telemetry.LogError(r.Context(), "uploads.save_file", err, map[string]any{"filename": safeFilename})
		writeError(w, http.StatusInternalServerError, "Error al almacenar archivo en el servidor")
		return
	}
	defer out.Close()

	written, err := io.Copy(out, file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al escribir archivo subido")
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}
	publicURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, safeFilename)

	telemetry.LogAction(r.Context(), "uploads.save_file", map[string]any{
		"category":   category,
		"filename":   safeFilename,
		"size_bytes": written,
	})

	writeJSON(w, http.StatusCreated, models.UploadResponse{
		URL:      publicURL,
		Filename: safeFilename,
		Category: category,
		Size:     written,
	})
}

func (h *APIHandler) handleListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.store.ListCategories(r.Context())
	if err != nil {
		telemetry.LogError(r.Context(), "categories.list", err, nil)
		writeError(w, http.StatusInternalServerError, "Error al obtener rubros: "+err.Error())
		return
	}
	telemetry.LogAction(r.Context(), "categories.list", map[string]any{"count": len(cats)})
	writeJSON(w, http.StatusOK, cats)
}

func (h *APIHandler) handleDiagnosticSuggestion(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Problem string `json:"problem"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	res, err := h.store.SuggestCategoryByProblem(r.Context(), body.Problem)
	if err != nil {
		telemetry.LogError(r.Context(), "diagnostic.suggest", err, map[string]any{"problem": body.Problem})
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	telemetry.LogAction(r.Context(), "diagnostic.suggest", map[string]any{
		"problem":            body.Problem,
		"suggested_category": res.SuggestedCategory.Slug,
	})
	writeJSON(w, http.StatusOK, res)
}

func (h *APIHandler) handleListProfessionals(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	q := r.URL.Query()
	catID, _ := strconv.Atoi(q.Get("category_id"))
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(q.Get("lng"), 64)

	filter := repository.ProfessionalFilter{
		CategoryID:   catID,
		SearchQuery:  q.Get("q"),
		OnlyUrgent:   q.Get("only_urgent") == "true",
		OnlyLicensed: q.Get("only_licensed") == "true",
		TopRated:     q.Get("top_rated") == "true",
		Lat:          lat,
		Lng:          lng,
	}

	pros, err := h.store.ListProfessionals(r.Context(), filter)
	if err != nil {
		telemetry.LogError(r.Context(), "professionals.search_postgis", err, map[string]any{
			"category_id": filter.CategoryID,
			"query":       filter.SearchQuery,
		})
		writeError(w, http.StatusInternalServerError, "Error al buscar profesionales: "+err.Error())
		return
	}

	elapsed := time.Since(start)
	h.metrics.RecordPostGISSearch(filter.CategoryID, filter.OnlyUrgent, filter.OnlyLicensed, len(pros), elapsed)
	telemetry.LogAction(r.Context(), "professionals.search_postgis", map[string]any{
		"category_id":   filter.CategoryID,
		"query":         filter.SearchQuery,
		"only_urgent":   filter.OnlyUrgent,
		"only_licensed": filter.OnlyLicensed,
		"results_count": len(pros),
		"duration_ms":   fmt.Sprintf("%.2f", float64(elapsed.Microseconds())/1000.0),
	})

	writeJSON(w, http.StatusOK, pros)
}

func (h *APIHandler) handleGetProfessional(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	prof, err := h.store.GetProfessionalByID(r.Context(), id)
	if err != nil {
		telemetry.LogError(r.Context(), "professionals.get_detail", err, map[string]any{"professional_id": id})
		writeError(w, http.StatusNotFound, "Profesional no encontrado")
		return
	}
	telemetry.LogAction(r.Context(), "professionals.get_detail", map[string]any{
		"professional_id": prof.ID,
		"full_name":       prof.FullName,
		"rating_avg":      prof.RatingAvg,
	})
	writeJSON(w, http.StatusOK, prof)
}

func (h *APIHandler) handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var input models.CreateServiceRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Formato de solicitud inválido")
		return
	}
	if claims := getClaims(r); claims != nil && input.ClientID == "" {
		input.ClientID = claims.UserID
	}
	if input.ProfessionalID == "" || input.ProblemDescription == "" {
		writeError(w, http.StatusBadRequest, "Faltan datos obligatorios (professional_id, problem_description)")
		return
	}
	req, err := h.store.CreateServiceRequest(r.Context(), input)
	if err != nil {
		telemetry.LogError(r.Context(), "requests.create", err, map[string]any{"professional_id": input.ProfessionalID})
		writeError(w, http.StatusInternalServerError, "No se pudo registrar la solicitud: "+err.Error())
		return
	}
	h.metrics.RecordServiceRequestCreated(req.CategoryID, req.Urgency, req.ContactChannel)
	telemetry.LogAction(r.Context(), "requests.create", map[string]any{
		"request_id":      req.ID,
		"professional_id": req.ProfessionalID,
		"client_name":     req.ClientName,
		"urgency":         req.Urgency,
		"channel":         req.ContactChannel,
	})
	writeJSON(w, http.StatusCreated, req)
}

func (h *APIHandler) handleListRequests(w http.ResponseWriter, r *http.Request) {
	profID := r.URL.Query().Get("professional_id")
	list, err := h.store.ListServiceRequests(r.Context(), profID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al listar solicitudes: "+err.Error())
		return
	}
	telemetry.LogAction(r.Context(), "requests.list", map[string]any{"professional_id": profID, "count": len(list)})
	writeJSON(w, http.StatusOK, list)
}

func (h *APIHandler) handleListMyClientRequests(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if claims := getClaims(r); claims != nil && claims.UserID != "" {
		clientID = claims.UserID
	}
	list, err := h.store.ListClientServiceRequests(r.Context(), clientID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al listar tus solicitudes: "+err.Error())
		return
	}
	telemetry.LogAction(r.Context(), "requests.list_client", map[string]any{"client_id": clientID, "count": len(list)})
	writeJSON(w, http.StatusOK, list)
}

func (h *APIHandler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req, err := h.store.GetServiceRequestByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Solicitud no encontrada")
		return
	}
	telemetry.LogAction(r.Context(), "requests.get", map[string]any{"request_id": req.ID, "status": req.Status})
	writeJSON(w, http.StatusOK, req)
}

func (h *APIHandler) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	reqID := r.PathValue("id")
	var input models.CreateReviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Formato de calificación inválido")
		return
	}
	if input.Rating < 1 || input.Rating > 5 {
		writeError(w, http.StatusBadRequest, "La calificación debe ser de 1 a 5 estrellas")
		return
	}
	rev, err := h.store.CreateReview(r.Context(), reqID, input)
	if err != nil {
		telemetry.LogError(r.Context(), "reviews.create_transactional", err, map[string]any{"request_id": reqID})
		writeError(w, http.StatusInternalServerError, "No se pudo guardar la calificación: "+err.Error())
		return
	}
	h.metrics.RecordReviewSubmitted(reqID, input.Rating)
	telemetry.LogAction(r.Context(), "reviews.create_transactional", map[string]any{
		"review_id":  rev.ID,
		"request_id": reqID,
		"rating":     rev.Rating,
		"author":     rev.AuthorName,
	})
	writeJSON(w, http.StatusCreated, rev)
}

func (h *APIHandler) handleUpdateVerification(w http.ResponseWriter, r *http.Request) {
	profID := r.PathValue("id")
	var input models.OnboardingVerificationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Datos de verificación inválidos")
		return
	}
	prof, err := h.store.UpdateProfessionalVerification(r.Context(), profID, input)
	if err != nil {
		telemetry.LogError(r.Context(), "professionals.verify_dni_license", err, map[string]any{"professional_id": profID})
		writeError(w, http.StatusNotFound, "Profesional no encontrado: "+err.Error())
		return
	}
	telemetry.LogAction(r.Context(), "professionals.verify_dni_license", map[string]any{
		"professional_id": profID,
		"is_dni_verified": prof.IsDNIVerified,
		"is_licensed":     prof.IsLicensed,
	})
	writeJSON(w, http.StatusOK, prof)
}

func (h *APIHandler) handleUpdateAvailability(w http.ResponseWriter, r *http.Request) {
	profID := r.PathValue("id")
	var input models.AvailabilityToggleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Datos de disponibilidad inválidos")
		return
	}
	prof, err := h.store.UpdateProfessionalAvailability(r.Context(), profID, input)
	if err != nil {
		telemetry.LogError(r.Context(), "professionals.toggle_availability", err, map[string]any{"professional_id": profID})
		writeError(w, http.StatusNotFound, "Profesional no encontrado: "+err.Error())
		return
	}
	h.metrics.RecordAvailabilityToggle(input.IsAvailableNow)
	telemetry.LogAction(r.Context(), "professionals.toggle_availability", map[string]any{
		"professional_id":  profID,
		"is_available_now": prof.IsAvailableNow,
	})
	writeJSON(w, http.StatusOK, prof)
}

func (h *APIHandler) handleGetProfessionalDashboard(w http.ResponseWriter, r *http.Request) {
	profID := r.PathValue("id")
	dash, err := h.store.GetProfessionalDashboard(r.Context(), profID)
	if err != nil {
		telemetry.LogError(r.Context(), "professionals.get_dashboard", err, map[string]any{"professional_id": profID})
		writeError(w, http.StatusNotFound, "Panel profesional no encontrado: "+err.Error())
		return
	}
	telemetry.LogAction(r.Context(), "professionals.get_dashboard", map[string]any{
		"professional_id": profID,
		"active_requests": len(dash.ActiveRequests),
		"new_today":       dash.NewTodayCount,
	})
	writeJSON(w, http.StatusOK, dash)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-ID, X-Request-ID, x-datadog-trace-id")
	w.Header().Set("Access-Control-Expose-Headers", "X-Trace-ID, X-Datadog-Trace-ID")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
