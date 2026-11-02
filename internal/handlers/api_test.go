package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"guiaoficios/backend/internal/models"
	"guiaoficios/backend/internal/repository"
)

func TestStitchFullWorkflowAPI(t *testing.T) {
	store := repository.NewInMemoryStore()
	h := NewAPIHandler(store, "in-memory-test")
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// 1. Pantalla 2: Registro de usuario
	regPayload := models.RegisterUserInput{
		FullName:     "Carlos Alberto Gómez",
		Phone:        "11 5544-3322",
		Email:        "carlos.nuevo@correo.com",
		Role:         "client",
		AuthProvider: "manual",
	}
	regBody, _ := json.Marshal(regPayload)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(regBody))
	recReg := httptest.NewRecorder()
	mux.ServeHTTP(recReg, reqReg)
	if recReg.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 en registro, obtuvo %d: %s", recReg.Code, recReg.Body.String())
	}

	// 2. Pantalla 3: Diagnóstico "¿Qué está pasando?"
	diagBody := []byte(`{"problem":"No enciende el calefón de la cocina"}`)
	reqDiag := httptest.NewRequest(http.MethodPost, "/api/v1/diagnostic", bytes.NewReader(diagBody))
	recDiag := httptest.NewRecorder()
	mux.ServeHTTP(recDiag, reqDiag)
	if recDiag.Code != http.StatusOK {
		t.Fatalf("esperaba 200 en diagnostico, obtuvo %d", recDiag.Code)
	}
	var diagRes models.DiagnosticSuggestionResponse
	_ = json.Unmarshal(recDiag.Body.Bytes(), &diagRes)
	if diagRes.SuggestedCategory.Slug != "gasista" {
		t.Fatalf("esperaba sugerencia gasista, obtuvo %s", diagRes.SuggestedCategory.Slug)
	}

	// 3. Pantalla 4: Listado de profesionales con filtro Matriculados
	reqPros := httptest.NewRequest(http.MethodGet, "/api/v1/professionals?only_licensed=true", nil)
	recPros := httptest.NewRecorder()
	mux.ServeHTTP(recPros, reqPros)
	if recPros.Code != http.StatusOK {
		t.Fatalf("esperaba 200 en profesionales, obtuvo %d", recPros.Code)
	}
	var pros []models.ProfessionalSummary
	_ = json.Unmarshal(recPros.Body.Bytes(), &pros)
	if len(pros) == 0 || pros[0].FullName != "Roberto Carlos Gómez" {
		t.Fatalf("esperaba encontrar a Roberto Carlos Gómez matriculado, obtuvo %+v", pros)
	}

	// 4. Pantalla 6 y 7: Solicitud Guiada a Roberto y Confirmación
	createReq := models.CreateServiceRequestInput{
		ClientName:         "María Elena Castro",
		ClientPhone:        "11 5544-3322",
		ProfessionalID:     pros[0].ID,
		CategoryID:         1,
		Urgency:            "urgente",
		ProblemDescription: "Pérdida en canilla de cocina",
		PhotoURL:           "Foto_caño_cocina.jpg",
		ContactChannel:     "whatsapp",
	}
	reqBytes, _ := json.Marshal(createReq)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/requests", bytes.NewReader(reqBytes))
	httpRec := httptest.NewRecorder()
	mux.ServeHTTP(httpRec, httpReq)
	if httpRec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 al crear solicitud, obtuvo %d", httpRec.Code)
	}
	var createdSR models.ServiceRequest
	_ = json.Unmarshal(httpRec.Body.Bytes(), &createdSR)

	// 5. Pantalla 8: Calificación de trabajo finalizado con los 5 aspectos destacados
	revInput := models.CreateReviewInput{
		AuthorName: "María Elena Castro",
		Rating:     5,
		Highlights: []string{"Puntualidad", "Precio justo", "Honestidad y prolijidad", "Excelente atención", "Solución en el acto"},
		Comment:    "Impecable trabajo. Roberto vino en el horario acordado.",
	}
	revBytes, _ := json.Marshal(revInput)
	reqRev := httptest.NewRequest(http.MethodPost, "/api/v1/requests/"+createdSR.ID+"/reviews", bytes.NewReader(revBytes))
	recRev := httptest.NewRecorder()
	mux.ServeHTTP(recRev, reqRev)
	if recRev.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 al calificar, obtuvo %d: %s", recRev.Code, recRev.Body.String())
	}

	// 6. Pantalla 9 & 10: Panel del Profesional y cambio de disponibilidad
	availBody := []byte(`{"is_available_now":false}`)
	reqAvail := httptest.NewRequest(http.MethodPatch, "/api/v1/professionals/"+pros[0].ID+"/availability", bytes.NewReader(availBody))
	recAvail := httptest.NewRecorder()
	mux.ServeHTTP(recAvail, reqAvail)
	if recAvail.Code != http.StatusOK {
		t.Fatalf("esperaba 200 al cambiar disponibilidad, obtuvo %d", recAvail.Code)
	}

	reqDash := httptest.NewRequest(http.MethodGet, "/api/v1/professionals/"+pros[0].ID+"/dashboard", nil)
	recDash := httptest.NewRecorder()
	mux.ServeHTTP(recDash, reqDash)
	if recDash.Code != http.StatusOK {
		t.Fatalf("esperaba 200 en dashboard profesional, obtuvo %d", recDash.Code)
	}
}
