package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"guiaoficios/backend/internal/auth"
	"guiaoficios/backend/internal/models"
	"guiaoficios/backend/internal/repository"
	"guiaoficios/backend/internal/telemetry"
)

// testStore implementa repository.Store exclusivamente para pruebas unitarias de handlers HTTP (archivo *_test.go).
// El binario de producción usa exclusivamente repository.PostgresStore (PostgreSQL + PostGIS).
type testStore struct {
	users    map[string]models.User
	pwHashes map[string]string
	pros     map[string]models.ProfessionalDetail
	requests map[string]models.ServiceRequest
}

func newTestStore() *testStore {
	pwHash, _ := auth.HashPassword("Clave1234")
	s := &testStore{
		users: map[string]models.User{
			"11111111-1111-1111-1111-111111111111": {
				ID:           "11111111-1111-1111-1111-111111111111",
				FullName:     "María Elena Castro",
				Phone:        "11 5544-3322",
				Email:        "maria.elena@correo.com",
				Role:         "client",
				AuthProvider: "manual",
				SavedAddress: "Av. Rivadavia 4520, Caballito",
				SavedLat:     -34.6185,
				SavedLng:     -58.4365,
			},
		},
		pwHashes: map[string]string{
			"maria.elena@correo.com": pwHash,
		},
		pros: map[string]models.ProfessionalDetail{
			"aaaa1111-1111-1111-1111-111111111111": {
				ProfessionalSummary: models.ProfessionalSummary{
					ID:                 "aaaa1111-1111-1111-1111-111111111111",
					UserID:             "22222222-2222-2222-2222-222222222222",
					FullName:           "Roberto Carlos Gómez",
					TradeTitle:         "Plomero y Gasista Matriculado",
					PrimaryCategoryID:  1,
					CategoryName:       "Plomero",
					IsDNIVerified:      true,
					IsLicensed:         true,
					IsAvailableNow:     true,
					AttendsEmergencies: true,
					DistanceKm:         1.2,
					RatingAvg:          4.9,
					ReviewsCount:       128,
					Phone:              "5491145678901",
				},
			},
		},
		requests: make(map[string]models.ServiceRequest),
	}
	return s
}

func (s *testStore) Ping(_ context.Context) (string, error) {
	return "PostgreSQL + PostGIS 3.4", nil
}

func (s *testStore) RegisterUser(_ context.Context, input models.RegisterUserInput) (*models.User, error) {
	rawStreet := strings.TrimSpace(input.Address)
	if rawStreet == "" {
		rawStreet = strings.TrimSpace(input.SavedAddress)
	}
	cleanNeighborhood, resolvedAddr, resolvedLat, resolvedLng := repository.ResolveNeighborhoodCoordinates(
		input.Neighborhood,
		rawStreet,
		input.SavedLat,
		input.SavedLng,
	)
	u := models.User{
		ID:           "user-new-001",
		FullName:     input.FullName,
		Phone:        input.Phone,
		Email:        strings.ToLower(input.Email),
		Role:         input.Role,
		AuthProvider: "manual",
		SavedAddress: resolvedAddr,
		SavedLat:     resolvedLat,
		SavedLng:     resolvedLng,
	}
	if u.Role == "professional" {
		u.ProfessionalID = "prof-new-001"
		s.pros[u.ProfessionalID] = models.ProfessionalDetail{
			ProfessionalSummary: models.ProfessionalSummary{
				ID:           u.ProfessionalID,
				UserID:       u.ID,
				FullName:     u.FullName,
				Neighborhood: cleanNeighborhood,
			},
		}
	}
	if input.Password != "" {
		h, _ := auth.HashPassword(input.Password)
		s.pwHashes[u.Email] = h
	}
	s.users[u.ID] = u
	return &u, nil
}

func (s *testStore) AuthenticateUser(_ context.Context, email, password string) (*models.User, error) {
	id := strings.ToLower(strings.TrimSpace(email))
	hash, ok := s.pwHashes[id]
	if !ok || !auth.VerifyPassword(password, hash) {
		return nil, errors.New("credenciales inválidas")
	}
	for _, u := range s.users {
		if u.Email == id {
			uCopy := u
			return &uCopy, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (s *testStore) AuthenticateGoogle(_ context.Context, input models.GoogleAuthInput) (*models.User, error) {
	claims, err := auth.VerifyGoogleIDToken(input.IDToken)
	if err != nil {
		return nil, err
	}
	u := models.User{
		ID:           "user-google-001",
		FullName:     claims.Name,
		Email:        claims.Email,
		Role:         input.Role,
		AuthProvider: "google",
	}
	s.users[u.ID] = u
	return &u, nil
}

func (s *testStore) GetUserByID(_ context.Context, userID string) (*models.User, error) {
	if u, ok := s.users[userID]; ok {
		return &u, nil
	}
	return nil, errors.New("not found")
}

func (s *testStore) UpdateUserAddress(_ context.Context, userID string, input models.UpdateAddressInput) (*models.User, error) {
	u, ok := s.users[userID]
	if !ok {
		return nil, errors.New("not found")
	}
	u.SavedAddress = input.SavedAddress
	u.SavedLat = input.SavedLat
	u.SavedLng = input.SavedLng
	s.users[userID] = u
	return &u, nil
}

func (s *testStore) ListCategories(_ context.Context) ([]models.Category, error) {
	return []models.Category{
		{ID: 1, Name: "Plomero", Slug: "plomero"},
		{ID: 2, Name: "Gasista", Slug: "gasista"},
	}, nil
}

func (s *testStore) SuggestCategoryByProblem(_ context.Context, problem string) (*models.DiagnosticSuggestionResponse, error) {
	return &models.DiagnosticSuggestionResponse{
		ProblemText: problem,
		SuggestedCategory: models.Category{
			ID:   2,
			Name: "Gasista",
			Slug: "gasista",
		},
		Explanation: "Te sugerimos un Gasista Matriculado.",
	}, nil
}

func (s *testStore) ListProfessionals(_ context.Context, _ repository.ProfessionalFilter) ([]models.ProfessionalSummary, error) {
	var list []models.ProfessionalSummary
	for _, p := range s.pros {
		list = append(list, p.ProfessionalSummary)
	}
	return list, nil
}

func (s *testStore) GetProfessionalByID(_ context.Context, id string) (*models.ProfessionalDetail, error) {
	if p, ok := s.pros[id]; ok {
		return &p, nil
	}
	return nil, errors.New("not found")
}

func (s *testStore) CreateServiceRequest(_ context.Context, input models.CreateServiceRequestInput) (*models.ServiceRequest, error) {
	sr := models.ServiceRequest{
		ID:                 "req-test-001",
		ClientID:           input.ClientID,
		ClientName:         input.ClientName,
		ClientPhone:        input.ClientPhone,
		ProfessionalID:     input.ProfessionalID,
		CategoryID:         input.CategoryID,
		Urgency:            input.Urgency,
		ProblemDescription: input.ProblemDescription,
		PhotoURL:           input.PhotoURL,
		ContactChannel:     input.ContactChannel,
		Status:             "contacted",
	}
	s.requests[sr.ID] = sr
	return &sr, nil
}

func (s *testStore) GetServiceRequestByID(_ context.Context, id string) (*models.ServiceRequest, error) {
	if r, ok := s.requests[id]; ok {
		return &r, nil
	}
	return nil, errors.New("not found")
}

func (s *testStore) ListServiceRequests(_ context.Context, _ string) ([]models.ServiceRequest, error) {
	var list []models.ServiceRequest
	for _, r := range s.requests {
		list = append(list, r)
	}
	return list, nil
}

func (s *testStore) ListClientServiceRequests(_ context.Context, _ string) ([]models.ServiceRequest, error) {
	return s.ListServiceRequests(context.Background(), "")
}

func (s *testStore) CreateReview(_ context.Context, _ string, input models.CreateReviewInput) (*models.Review, error) {
	return &models.Review{
		ID:         "rev-001",
		AuthorName: input.AuthorName,
		Rating:     input.Rating,
		Highlights: input.Highlights,
		Comment:    input.Comment,
	}, nil
}

func (s *testStore) UpdateProfessionalVerification(_ context.Context, profID string, input models.OnboardingVerificationInput) (*models.ProfessionalDetail, error) {
	p, ok := s.pros[profID]
	if !ok {
		return nil, errors.New("not found")
	}
	p.DNIFrontURL = input.DNIFrontURL
	p.DNIBackURL = input.DNIBackURL
	p.LicenseDocURL = input.LicenseDocURL
	s.pros[profID] = p
	return &p, nil
}

func (s *testStore) UpdateProfessionalAvailability(_ context.Context, profID string, input models.AvailabilityToggleInput) (*models.ProfessionalDetail, error) {
	p, ok := s.pros[profID]
	if !ok {
		return nil, errors.New("not found")
	}
	p.IsAvailableNow = input.IsAvailableNow
	s.pros[profID] = p
	return &p, nil
}

func (s *testStore) GetProfessionalDashboard(_ context.Context, profID string) (*models.ProfessionalDashboardResponse, error) {
	p, ok := s.pros[profID]
	if !ok {
		return nil, errors.New("not found")
	}
	return &models.ProfessionalDashboardResponse{
		Professional: p,
	}, nil
}

func TestStitchFullWorkflowAPI(t *testing.T) {
	store := newTestStore()
	tmpUploadDir := t.TempDir()
	h := NewAPIHandler(store, "postgresql+postgis", tmpUploadDir)
	rawMux := http.NewServeMux()
	h.RegisterRoutes(rawMux)
	mux := telemetry.ObservabilityMiddleware(rawMux)

	// 1. Pantalla 2: Registro de usuario con barrio/dirección real (BE-2.1) y obtención de JWT + Propagación de X-Trace-ID
	regPayload := models.RegisterUserInput{
		FullName:     "Carlos Alberto Gómez",
		Phone:        "+54 9 11 5544-3322",
		Email:        "carlos.nuevo@correo.com",
		Password:     "MiClaveSegura2026",
		Role:         "client",
		AuthProvider: "manual",
		Neighborhood: "Palermo, CABA",
		Address:      "Av. Santa Fe 3420",
	}
	regBody, _ := json.Marshal(regPayload)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(regBody))
	reqReg.Header.Set("X-Trace-ID", "and-test-trace-001")
	recReg := httptest.NewRecorder()
	mux.ServeHTTP(recReg, reqReg)
	if recReg.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 en registro, obtuvo %d: %s", recReg.Code, recReg.Body.String())
	}
	if gotTrace := recReg.Header().Get("X-Trace-ID"); gotTrace != "and-test-trace-001" {
		t.Fatalf("esperaba X-Trace-ID=and-test-trace-001, obtuvo %q", gotTrace)
	}
	var authRes models.AuthResponse
	_ = json.Unmarshal(recReg.Body.Bytes(), &authRes)
	if authRes.Token == "" {
		t.Fatalf("esperaba token JWT en registro")
	}
	if authRes.User.SavedAddress != "Av. Santa Fe 3420, Palermo, CABA" || authRes.User.SavedLat != -34.5889 {
		t.Fatalf("esperaba dirección y coordenadas reales de Palermo, obtuvo addr=%q lat=%v", authRes.User.SavedAddress, authRes.User.SavedLat)
	}

	// 1b. Google Auth debe rechazar email en texto plano sin id_token firmado (BE-2.2)
	unsignedGoogleBody := []byte(`{"id_token":"google-android-account-123","email":"intruso@gmail.com","role":"client"}`)
	reqBadGoogle := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google", bytes.NewReader(unsignedGoogleBody))
	recBadGoogle := httptest.NewRecorder()
	mux.ServeHTTP(recBadGoogle, reqBadGoogle)
	if recBadGoogle.Code != http.StatusUnauthorized {
		t.Fatalf("esperaba 401 al intentar Google Auth sin id_token firmado, obtuvo %d: %s", recBadGoogle.Code, recBadGoogle.Body.String())
	}

	// 1c. Google Auth con id_token firmado criptográficamente debe autenticar (BE-2.2)
	signedGToken, _ := auth.GenerateGoogleIDToken("vecina.belgrano@gmail.com", "Marta López", "g-sub-999")
	validGooglePayload := models.GoogleAuthInput{
		IDToken: signedGToken,
		Email:   "vecina.belgrano@gmail.com",
		Role:    "client",
	}
	validGoogleBytes, _ := json.Marshal(validGooglePayload)
	reqGoodGoogle := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google", bytes.NewReader(validGoogleBytes))
	recGoodGoogle := httptest.NewRecorder()
	mux.ServeHTTP(recGoodGoogle, reqGoodGoogle)
	if recGoodGoogle.Code != http.StatusOK {
		t.Fatalf("esperaba 200 en Google Auth con id_token firmado, obtuvo %d: %s", recGoodGoogle.Code, recGoodGoogle.Body.String())
	}

	// 2. Login con credenciales válidas
	loginPayload := models.LoginInput{
		Email:    "carlos.nuevo@correo.com",
		Password: "MiClaveSegura2026",
	}
	loginBody, _ := json.Marshal(loginPayload)
	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	recLogin := httptest.NewRecorder()
	mux.ServeHTTP(recLogin, reqLogin)
	if recLogin.Code != http.StatusOK {
		t.Fatalf("esperaba 200 en login, obtuvo %d: %s", recLogin.Code, recLogin.Body.String())
	}

	// 2b. Login con contraseña incorrecta debe devolver 401 Unauthorized
	badLoginPayload := models.LoginInput{
		Email:    "carlos.nuevo@correo.com",
		Password: "ClaveIncorrecta999",
	}
	badLoginBody, _ := json.Marshal(badLoginPayload)
	reqBadLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(badLoginBody))
	recBadLogin := httptest.NewRecorder()
	mux.ServeHTTP(recBadLogin, reqBadLogin)
	if recBadLogin.Code != http.StatusUnauthorized {
		t.Fatalf("esperaba 401 en login con clave inválida, obtuvo %d: %s", recBadLogin.Code, recBadLogin.Body.String())
	}

	// 3. Upload real multipart/form-data (GDO-10)
	var bodyBuf bytes.Buffer
	mw := multipart.NewWriter(&bodyBuf)
	_ = mw.WriteField("category", "problem_photo")
	fw, _ := mw.CreateFormFile("file", "canilla_rota.jpg")
	_, _ = fw.Write([]byte("fake-jpeg-bytes-for-upload-test"))
	_ = mw.Close()

	reqUp := httptest.NewRequest(http.MethodPost, "/api/v1/uploads", &bodyBuf)
	reqUp.Header.Set("Content-Type", mw.FormDataContentType())
	reqUp.Header.Set("Authorization", "Bearer "+authRes.Token)
	recUp := httptest.NewRecorder()
	mux.ServeHTTP(recUp, reqUp)
	if recUp.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 en upload multipart, obtuvo %d: %s", recUp.Code, recUp.Body.String())
	}
	var upRes models.UploadResponse
	_ = json.Unmarshal(recUp.Body.Bytes(), &upRes)
	if upRes.Filename == "" || !strings.Contains(upRes.URL, "/uploads/") {
		t.Fatalf("respuesta de upload inválida: %+v", upRes)
	}

	// 4. Pantalla 3: Diagnóstico "¿Qué está pasando?"
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

	// 5. Pantalla 4: Listado de profesionales con filtro Matriculados
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

	// 6. Pantalla 6 y 7: Solicitud Guiada a Roberto y Confirmación
	createReq := models.CreateServiceRequestInput{
		ClientName:         "María Elena Castro",
		ClientPhone:        "11 5544-3322",
		ProfessionalID:     pros[0].ID,
		CategoryID:         1,
		Urgency:            "urgente",
		ProblemDescription: "Pérdida en canilla de cocina",
		PhotoURL:           upRes.URL,
		ContactChannel:     "whatsapp",
	}
	reqBytes, _ := json.Marshal(createReq)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/requests", bytes.NewReader(reqBytes))
	httpReq.Header.Set("Authorization", "Bearer "+authRes.Token)
	httpRec := httptest.NewRecorder()
	mux.ServeHTTP(httpRec, httpReq)
	if httpRec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 al crear solicitud, obtuvo %d", httpRec.Code)
	}
	var createdSR models.ServiceRequest
	_ = json.Unmarshal(httpRec.Body.Bytes(), &createdSR)

	// 7. Pantalla 8: Calificación de trabajo finalizado con los 5 aspectos destacados
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

	// 8. Pantalla 9 & 10: Panel del Profesional y cambio de disponibilidad
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

	// 9. Verificar Snapshot de Métricas Clave de Datadog (/api/v1/metrics)
	reqMetrics := httptest.NewRequest(http.MethodGet, "/api/v1/metrics", nil)
	recMetrics := httptest.NewRecorder()
	mux.ServeHTTP(recMetrics, reqMetrics)
	if recMetrics.Code != http.StatusOK {
		t.Fatalf("esperaba 200 en /api/v1/metrics, obtuvo %d", recMetrics.Code)
	}
	t.Logf("Snapshot Métricas Datadog: %s", recMetrics.Body.String())
}
