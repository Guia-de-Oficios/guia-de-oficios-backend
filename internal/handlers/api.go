package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"guiaoficios/backend/internal/models"
	"guiaoficios/backend/internal/repository"
)

type APIHandler struct {
	store repository.Store
	mode  string
}

func NewAPIHandler(store repository.Store, mode string) *APIHandler {
	return &APIHandler{store: store, mode: mode}
}

func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	// Health
	mux.HandleFunc("GET /api/v1/health", h.handleHealth)

	// Pantallas 1 & 2: Bienvenida, Selección de Rol y Registro (Google o Manual)
	mux.HandleFunc("POST /api/v1/auth/register", h.handleRegisterUser)
	mux.HandleFunc("GET /api/v1/auth/me", h.handleGetCurrentUser)

	// Pantalla 3: Cliente Inicio, Categorías y Asistente "¿No sabés qué profesional necesitás?"
	mux.HandleFunc("GET /api/v1/categories", h.handleListCategories)
	mux.HandleFunc("POST /api/v1/diagnostic", h.handleDiagnosticSuggestion)

	// Pantallas 4 & 5: Resultados de Profesionales (con filtros Todos/Urgencias/Matriculados/Mejor Valorados) y Perfil
	mux.HandleFunc("GET /api/v1/professionals", h.handleListProfessionals)
	mux.HandleFunc("GET /api/v1/professionals/{id}", h.handleGetProfessional)

	// Pantallas 6, 7 & 8: Solicitud Guiada a Roberto, Confirmación de Contacto y Calificación de Trabajo Terminado
	mux.HandleFunc("POST /api/v1/requests", h.handleCreateRequest)
	mux.HandleFunc("GET /api/v1/requests", h.handleListRequests)
	mux.HandleFunc("GET /api/v1/requests/{id}", h.handleGetRequest)
	mux.HandleFunc("POST /api/v1/requests/{id}/reviews", h.handleCreateReview)

	// Pantallas 9 & 10: Profesional Registro/Verificación DNI + Matrícula y Panel de Solicitudes/Disponibilidad
	mux.HandleFunc("POST /api/v1/professionals/{id}/verification", h.handleUpdateVerification)
	mux.HandleFunc("PATCH /api/v1/professionals/{id}/availability", h.handleUpdateAvailability)
	mux.HandleFunc("GET /api/v1/professionals/{id}/dashboard", h.handleGetProfessionalDashboard)
}

func (h *APIHandler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "Guía de Oficios API (Stitch 10-Screen Complete Edition)",
		"storage": h.mode,
	})
}

func (h *APIHandler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Datos de registro inválidos")
		return
	}
	if input.FullName == "" || input.Email == "" {
		writeError(w, http.StatusBadRequest, "Nombre completo y correo electrónico son obligatorios")
		return
	}
	user, err := h.store.RegisterUser(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al registrar usuario: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *APIHandler) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	user, err := h.store.GetCurrentUser(r.Context(), role)
	if err != nil {
		writeError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *APIHandler) handleListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.store.ListCategories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al obtener rubros: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

func (h *APIHandler) handleDiagnosticSuggestion(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Problem string `json:"problem"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	res, err := h.store.SuggestCategoryByProblem(r.Context(), body.Problem)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *APIHandler) handleListProfessionals(w http.ResponseWriter, r *http.Request) {
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
		writeError(w, http.StatusInternalServerError, "Error al buscar profesionales: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pros)
}

func (h *APIHandler) handleGetProfessional(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	prof, err := h.store.GetProfessionalByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Profesional no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, prof)
}

func (h *APIHandler) handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var input models.CreateServiceRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Formato de solicitud inválido")
		return
	}
	if input.ProfessionalID == "" || input.ProblemDescription == "" {
		writeError(w, http.StatusBadRequest, "Faltan datos obligatorios (professional_id, problem_description)")
		return
	}
	req, err := h.store.CreateServiceRequest(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "No se pudo registrar la solicitud: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, req)
}

func (h *APIHandler) handleListRequests(w http.ResponseWriter, r *http.Request) {
	profID := r.URL.Query().Get("professional_id")
	list, err := h.store.ListServiceRequests(r.Context(), profID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al listar solicitudes: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *APIHandler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req, err := h.store.GetServiceRequestByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Solicitud no encontrada")
		return
	}
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
		writeError(w, http.StatusInternalServerError, "No se pudo guardar la calificación: "+err.Error())
		return
	}
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
		writeError(w, http.StatusNotFound, "Profesional no encontrado")
		return
	}
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
		writeError(w, http.StatusNotFound, "Profesional no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, prof)
}

func (h *APIHandler) handleGetProfessionalDashboard(w http.ResponseWriter, r *http.Request) {
	profID := r.PathValue("id")
	dash, err := h.store.GetProfessionalDashboard(r.Context(), profID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Panel profesional no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, dash)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
