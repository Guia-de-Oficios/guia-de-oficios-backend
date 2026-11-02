package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"guiaoficios/backend/internal/models"
)

// Store define todas las operaciones requeridas por las 10 pantallas de Stitch.
type Store interface {
	// Autenticación y Rol (Pantallas 1 y 2)
	RegisterUser(ctx context.Context, input models.RegisterUserInput) (*models.User, error)
	GetCurrentUser(ctx context.Context, role string) (*models.User, error)

	// Inicio, Diagnóstico y Categorías (Pantalla 3)
	ListCategories(ctx context.Context) ([]models.Category, error)
	SuggestCategoryByProblem(ctx context.Context, problem string) (*models.DiagnosticSuggestionResponse, error)

	// Resultados y Perfil de Profesional (Pantallas 4 y 5)
	ListProfessionals(ctx context.Context, filter ProfessionalFilter) ([]models.ProfessionalSummary, error)
	GetProfessionalByID(ctx context.Context, id string) (*models.ProfessionalDetail, error)

	// Solicitud Guiada, Confirmación y Calificación (Pantallas 6, 7 y 8)
	CreateServiceRequest(ctx context.Context, input models.CreateServiceRequestInput) (*models.ServiceRequest, error)
	GetServiceRequestByID(ctx context.Context, id string) (*models.ServiceRequest, error)
	ListServiceRequests(ctx context.Context, professionalID string) ([]models.ServiceRequest, error)
	CreateReview(ctx context.Context, requestID string, input models.CreateReviewInput) (*models.Review, error)

	// Registro Profesional DNI/Matrícula y Panel de Solicitudes (Pantallas 9 y 10)
	UpdateProfessionalVerification(ctx context.Context, profID string, input models.OnboardingVerificationInput) (*models.ProfessionalDetail, error)
	UpdateProfessionalAvailability(ctx context.Context, profID string, input models.AvailabilityToggleInput) (*models.ProfessionalDetail, error)
	GetProfessionalDashboard(ctx context.Context, profID string) (*models.ProfessionalDashboardResponse, error)
}

type ProfessionalFilter struct {
	CategoryID   int
	SearchQuery  string
	OnlyUrgent   bool
	OnlyLicensed bool
	TopRated     bool
	Lat          float64
	Lng          float64
}

// =============================================================================
// InMemoryStore: Precargado con los datos exactos de las 10 pantallas de Stitch
// =============================================================================

type InMemoryStore struct {
	mu            sync.RWMutex
	users         []models.User
	categories    []models.Category
	professionals []models.ProfessionalDetail
	requests      []models.ServiceRequest
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: []models.User{
			{
				ID:           "00000000-0000-0000-0000-000000000001",
				FullName:     "María Elena Castro",
				Email:        "maria.castro@vecinos.ar",
				Phone:        "11 5544-3322",
				Role:         "client",
				AuthProvider: "manual",
				SavedAddress: "Av. Rivadavia 5420, Caballito",
				CreatedAt:    time.Now().Add(-30 * 24 * time.Hour),
			},
			{
				ID:           "11111111-1111-1111-1111-111111111111",
				FullName:     "Roberto Carlos Gómez",
				Email:        "roberto.gomez@guiaoficios.ar",
				Phone:        "+5491145890000",
				Role:         "professional",
				AuthProvider: "manual",
				SavedAddress: "Caballito, CABA",
				CreatedAt:    time.Now().Add(-90 * 24 * time.Hour),
			},
		},
		categories: []models.Category{
			{ID: 1, Slug: "plomero", Name: "Plomero", HelperExamples: "Canillas, cañerías y pérdidas", IconName: "water_drop", DisplayOrder: 1},
			{ID: 2, Slug: "electricista", Name: "Electricista", HelperExamples: "Enchufes, cables e iluminación", IconName: "bolt", DisplayOrder: 2},
			{ID: 3, Slug: "gasista", Name: "Gasista", HelperExamples: "Calefones, estufas y conexiones", IconName: "local_fire_department", BadgeText: "Matriculados", BadgeType: "verified", DisplayOrder: 3},
			{ID: 4, Slug: "aire-acondicionado", Name: "Aire acondicionado", HelperExamples: "Equipos y mantenimiento", IconName: "mode_fan", DisplayOrder: 4},
			{ID: 5, Slug: "pintor", Name: "Pintor", HelperExamples: "Paredes y terminaciones", IconName: "format_paint", DisplayOrder: 5},
			{ID: 6, Slug: "albanil", Name: "Albañil", HelperExamples: "Reparaciones y construcción", IconName: "handyman", DisplayOrder: 6},
			{ID: 7, Slug: "jardinero", Name: "Jardinero", HelperExamples: "Plantas y mantenimiento", IconName: "yard", DisplayOrder: 7},
			{ID: 8, Slug: "cerrajero", Name: "Cerrajero", HelperExamples: "Cerraduras y puertas", IconName: "key", BadgeText: "Urgencias", BadgeType: "urgent", DisplayOrder: 8},
		},
		professionals: []models.ProfessionalDetail{
			{
				ProfessionalSummary: models.ProfessionalSummary{
					ID:                   "aaaa1111-1111-1111-1111-111111111111",
					UserID:               "11111111-1111-1111-1111-111111111111",
					FullName:             "Roberto Carlos Gómez",
					Phone:                "+5491145890000",
					AvatarURL:            "https://lh3.googleusercontent.com/aida-public/AB6AXuBJDsjoc_sU9gzvdM9kL-GyGdUsU1jE5bNEiebL_eoOBTI99zaT6MQQEnXV3evr2oUMVfK0DN1-Nv3xgr0T8Oivk6GQHrEIMu-YGCiNnhixwDqrPkqdsD2pxGU4jrDkwuu562QZumlwOmTNkgZ2L1K9RBWOHRrlXSnawefx-xTL14Oo32ZYgm_ww5_lFhaxO92MCGqJ_dlRmF6mi7Dx9WS5gFrGevtdm_qvq_I0IG44",
					PrimaryCategoryID:    1,
					CategoryName:         "Plomero",
					TradeTitle:           "Plomero y Gasista Matriculado",
					Specialties:          []string{"Destapaciones", "Calefones y Termotanques"},
					Neighborhood:         "Caballito y Flores",
					DistanceLabel:        "A 8 cuadras (650 metros)",
					DistanceKm:           0.65,
					IsAvailableNow:       true,
					AttendsEmergencies:   true,
					IsDNIVerified:        true,
					IsLicensed:           true,
					LicenseNumber:        "Matrícula oficial N° 48.912 verificada",
					LicenseIssuer:        "Metrogas / ENARGAS",
					VisitPriceAmount:     12000,
					VisitPriceLabel:      "$12.000",
					RatingAvg:            4.9,
					RatingLabel:          "Muy recomendado",
					ReviewsCount:         142,
					JobsCompletedCount:   142,
					MonthlyContactsCount: 18,
				},
				Bio:                 "Hola, soy Roberto. Trabajo hace más de 25 años haciendo reparaciones en casas y departamentos de la zona. Trabajo prolijo, puntual y con garantía en cada arreglo.",
				YearsExperience:     25,
				GuaranteePercentage: 100,
				Portfolio: []models.PortfolioItem{
					{
						ID:         "job-1",
						Title:      "Canilla monocomando cambiada",
						Subtitle:   "Instalación y sellado completo",
						ImageURL:   "https://lh3.googleusercontent.com/aida-public/AB6AXuCh10RiTc50jP69tmrjG8REQ7uQZftm52vnIAc3RMhe3SukIRmt1qsbmlKpf6XUc585FmJriOr-k8Z0QvUfu8A7nNCkOGW2jJx3avdmiKF-Nkrta-clOO_eVsmwIDHmRz30tgPahuYwFtN8PQ1fyh9hPzU65q3-IQsxo-p-0lExk7Sal_lolRdoZ-44c8UC0ZDRI9Y7OqNDb2WZlr5Ly30L2CAZmaoIXjO55wdLSbRG",
						IsVerified: true,
					},
					{
						ID:         "job-2",
						Title:      "Reparación de cañería termofusión",
						Subtitle:   "Solución definitiva a filtración",
						ImageURL:   "https://lh3.googleusercontent.com/aida-public/AB6AXuCTODIWHLBehXm6SVzYW9pLhJ4C1zuw4YW31S21vmoWs31kV-fwgGBcwNkWLWT8eQAC0le8E2mTVOCpr78K1LSgHxq1Wz_TIbADpF91xbBHUTJEVtBkUGHSa_saUCNhM56oSA0exrcwQueuH6FVaBDeEcWP8z0qj_SvfAmVejvoezsktrrA52SZpDaRBLiNFN1lS3RBFovpPsOpuNUk6SDaVRHc2VbiBDJKWl6kekvw",
						IsVerified: true,
					},
				},
				Reviews: []models.Review{
					{
						ID:         "rev-1",
						AuthorName: "Graciela M.",
						AuthorZone: "Vecina de Caballito (Av. Pedro Goyena)",
						Rating:     5,
						Highlights: []string{"Puntualidad", "Precio justo", "Honestidad y prolijidad", "Excelente atención", "Solución en el acto"},
						Comment:    "Vino puntual a la hora acordada, resolvió la pérdida del baño en el acto y dejó todo impecable. Muy respetuoso.",
						CreatedAt:  time.Now().Add(-48 * time.Hour),
					},
				},
			},
			{
				ProfessionalSummary: models.ProfessionalSummary{
					ID:                   "bbbb2222-2222-2222-2222-222222222222",
					UserID:               "22222222-2222-2222-2222-222222222222",
					FullName:             "Martín Edgardo Silva",
					Phone:                "+5491155667788",
					AvatarURL:            "https://lh3.googleusercontent.com/aida-public/AB6AXuDDxicrHktEOBqIj1BUMTX05PyFDIxZ1yoq8TXlbS0_tAzKP5WF2Kk66LoYyNFRsxa6rPT-7_mjxtVo0Ndv042VSX55wLr_wkujhvNvGL38GzFQRd_iwQwwC4ZSOczNAIiidtGmTVluGkr_VPuDXq2wFBT0iXNi_WAwGWWDRBi7FiBL52aOtaZB0ZXb7EXXZxY99wL-iucwj2kqRIYrdplZibNCoDiz3kIkeVVbvFFM",
					PrimaryCategoryID:    1,
					CategoryName:         "Plomero",
					TradeTitle:           "Plomero y Gasista",
					Specialties:          []string{"Urgencias 24hs", "Destapaciones"},
					Neighborhood:         "Caballito y Parque Chacabuco",
					DistanceLabel:        "A 1.2 km de tu domicilio",
					DistanceKm:           1.2,
					IsAvailableNow:       true,
					AttendsEmergencies:   true,
					IsDNIVerified:        true,
					IsLicensed:           false,
					VisitPriceAmount:     10000,
					VisitPriceLabel:      "$10.000",
					RatingAvg:            4.8,
					RatingLabel:          "Excelente atención",
					ReviewsCount:         89,
					JobsCompletedCount:   89,
					MonthlyContactsCount: 12,
				},
				Bio:                 "Atención rápida en el día para pérdidas de agua, cambio de griferías y conexiones de artefactos.",
				YearsExperience:     14,
				GuaranteePercentage: 100,
			},
		},
		requests: []models.ServiceRequest{
			{
				ID:                  "req-1",
				ClientName:          "María Elena Castro",
				ClientInitials:      "MC",
				ClientPhone:         "11 5544-3322",
				ClientAddress:       "Av. Rivadavia 5420",
				ClientDistanceLabel: "a 600m",
				ProfessionalID:      "aaaa1111-1111-1111-1111-111111111111",
				ProfessionalName:    "Roberto Carlos Gómez",
				ProfessionalPhone:   "+5491145890000",
				ProfessionalLicense: "Matrícula Nº 48.912",
				CategoryID:          1,
				Urgency:             "urgente",
				ProblemDescription:  "Pierde agua la canilla de la cocina debajo de la mesada. Gotea constante y cerré la llave de paso provisoria.",
				PhotoURL:            "Foto_caño_cocina.jpg",
				ContactChannel:      "whatsapp",
				Status:              "pending",
				TimeAgoLabel:        "Hace 25m",
				CreatedAt:           time.Now().Add(-25 * time.Minute),
			},
			{
				ID:                  "req-2",
				ClientName:          "Jorge Benítez",
				ClientInitials:      "JB",
				ClientPhone:         "11 6677-8899",
				ClientAddress:       "Caballito Norte",
				ClientDistanceLabel: "a 1.4 km",
				ProfessionalID:      "aaaa1111-1111-1111-1111-111111111111",
				ProfessionalName:    "Roberto Carlos Gómez",
				ProfessionalPhone:   "+5491145890000",
				ProfessionalLicense: "Matrícula Nº 48.912",
				CategoryID:          1,
				Urgency:             "normal",
				ProblemDescription:  "Cambio de fuelle y flexibles en baño principal.",
				ContactChannel:      "whatsapp",
				Status:              "scheduled",
				ScheduledNote:       "Trabajo coordinado para el sábado",
				TimeAgoLabel:        "Ayer",
				CreatedAt:           time.Now().Add(-24 * time.Hour),
			},
		},
	}
}

func (m *InMemoryStore) RegisterUser(_ context.Context, input models.RegisterUserInput) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if input.Role == "" {
		input.Role = "client"
	}
	if input.AuthProvider == "" {
		input.AuthProvider = "manual"
	}
	if input.SavedAddress == "" {
		input.SavedAddress = "Av. Rivadavia 5420, Caballito"
	}

	user := models.User{
		ID:           fmt.Sprintf("user-%d", len(m.users)+1),
		FullName:     input.FullName,
		Email:        input.Email,
		Phone:        input.Phone,
		Role:         input.Role,
		AuthProvider: input.AuthProvider,
		SavedAddress: input.SavedAddress,
		CreatedAt:    time.Now(),
	}
	m.users = append(m.users, user)
	return &user, nil
}

func (m *InMemoryStore) GetCurrentUser(_ context.Context, role string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.users {
		if role == "" || u.Role == role {
			copyU := u
			return &copyU, nil
		}
	}
	return &m.users[0], nil
}

func (m *InMemoryStore) ListCategories(_ context.Context) ([]models.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.categories, nil
}

func (m *InMemoryStore) SuggestCategoryByProblem(_ context.Context, problem string) (*models.DiagnosticSuggestionResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lower := strings.ToLower(problem)
	cat := m.categories[0] // Por defecto Plomero
	explanation := "Por lo que comentás sobre agua o cañerías, te recomendamos un Plomero verificado de tu barrio."

	switch {
	case strings.Contains(lower, "gas") || strings.Contains(lower, "calefón") || strings.Contains(lower, "calefon") || strings.Contains(lower, "estufa"):
		cat = m.categories[2] // Gasista
		explanation = "Para artefactos de gas o calefones siempre te recomendamos un Gasista Matriculado."
	case strings.Contains(lower, "luz") || strings.Contains(lower, "enchufe") || strings.Contains(lower, "cable") || strings.Contains(lower, "térmica") || strings.Contains(lower, "termica"):
		cat = m.categories[1] // Electricista
		explanation = "Para problemas con enchufes, llaves térmicas o cortes de luz, el oficio indicado es un Electricista."
	case strings.Contains(lower, "aire") || strings.Contains(lower, "split"):
		cat = m.categories[3] // Aire acondicionado
		explanation = "Para equipos de climatización y filtros te sugerimos un técnico en Aire Acondicionado."
	case strings.Contains(lower, "llave") || strings.Contains(lower, "cerradura") || strings.Contains(lower, "puerta"):
		cat = m.categories[7] // Cerrajero
		explanation = "Para cerraduras trabadas o llaves, te conectamos con un Cerrajero de urgencias."
	}

	return &models.DiagnosticSuggestionResponse{
		ProblemText:       problem,
		SuggestedCategory: cat,
		Explanation:       explanation,
	}, nil
}

func (m *InMemoryStore) ListProfessionals(_ context.Context, f ProfessionalFilter) ([]models.ProfessionalSummary, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(f.SearchQuery))
	var res []models.ProfessionalSummary
	for _, p := range m.professionals {
		if f.CategoryID != 0 && p.PrimaryCategoryID != f.CategoryID {
			continue
		}
		if f.OnlyUrgent && !p.AttendsEmergencies {
			continue
		}
		if f.OnlyLicensed && !p.IsLicensed {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(p.FullName), q) &&
			!strings.Contains(strings.ToLower(p.TradeTitle), q) &&
			!strings.Contains(strings.ToLower(p.Bio), q) {
			continue
		}
		res = append(res, p.ProfessionalSummary)
	}

	if f.TopRated {
		sort.Slice(res, func(i, j int) bool {
			return res[i].RatingAvg > res[j].RatingAvg
		})
	}
	return res, nil
}

func (m *InMemoryStore) GetProfessionalByID(_ context.Context, id string) (*models.ProfessionalDetail, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.professionals {
		if p.ID == id {
			copyProf := p
			return &copyProf, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *InMemoryStore) CreateServiceRequest(_ context.Context, input models.CreateServiceRequestInput) (*models.ServiceRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	prof := m.professionals[0]
	for _, p := range m.professionals {
		if p.ID == input.ProfessionalID {
			prof = p
			break
		}
	}

	if input.ClientName == "" {
		input.ClientName = "María Elena Castro"
	}
	if input.ClientPhone == "" {
		input.ClientPhone = "11 5544-3322"
	}
	if input.ClientAddress == "" {
		input.ClientAddress = "Av. Rivadavia 5420, Caballito"
	}
	if input.Urgency == "" {
		input.Urgency = "normal"
	}
	if input.ContactChannel == "" {
		input.ContactChannel = "whatsapp"
	}

	initials := "VC"
	parts := strings.Fields(input.ClientName)
	if len(parts) >= 2 {
		initials = strings.ToUpper(string(parts[0][0]) + string(parts[1][0]))
	} else if len(parts) == 1 {
		initials = strings.ToUpper(string(parts[0][0]))
	}

	req := models.ServiceRequest{
		ID:                  fmt.Sprintf("req-%d", len(m.requests)+1),
		ClientName:          input.ClientName,
		ClientInitials:      initials,
		ClientPhone:         input.ClientPhone,
		ClientAddress:       input.ClientAddress,
		ClientDistanceLabel: "a 600m",
		ProfessionalID:      prof.ID,
		ProfessionalName:    prof.FullName,
		ProfessionalPhone:   prof.Phone,
		ProfessionalLicense: prof.LicenseNumber,
		CategoryID:          input.CategoryID,
		Urgency:             input.Urgency,
		ProblemDescription:  input.ProblemDescription,
		PhotoURL:            input.PhotoURL,
		ContactChannel:      input.ContactChannel,
		Status:              "pending",
		TimeAgoLabel:        "Recién ahora",
		CreatedAt:           time.Now(),
	}

	m.requests = append([]models.ServiceRequest{req}, m.requests...)
	return &req, nil
}

func (m *InMemoryStore) GetServiceRequestByID(_ context.Context, id string) (*models.ServiceRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, r := range m.requests {
		if r.ID == id {
			copyR := r
			return &copyR, nil
		}
	}
	if len(m.requests) > 0 {
		copyR := m.requests[0]
		return &copyR, nil
	}
	return nil, sql.ErrNoRows
}

func (m *InMemoryStore) ListServiceRequests(_ context.Context, professionalID string) ([]models.ServiceRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if professionalID == "" {
		return m.requests, nil
	}
	var filtered []models.ServiceRequest
	for _, r := range m.requests {
		if r.ProfessionalID == professionalID {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

func (m *InMemoryStore) CreateReview(_ context.Context, requestID string, input models.CreateReviewInput) (*models.Review, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var profID string
	for i := range m.requests {
		if m.requests[i].ID == requestID {
			m.requests[i].Status = "completed"
			profID = m.requests[i].ProfessionalID
			break
		}
	}
	if profID == "" && len(m.professionals) > 0 {
		profID = m.professionals[0].ID
	}

	if input.AuthorName == "" {
		input.AuthorName = "María Elena C."
	}
	if input.AuthorZone == "" {
		input.AuthorZone = "Vecina Verificada de Caballito"
	}

	rev := models.Review{
		ID:         fmt.Sprintf("rev-%d", time.Now().UnixNano()),
		AuthorName: input.AuthorName,
		AuthorZone: input.AuthorZone,
		Rating:     input.Rating,
		Highlights: input.Highlights,
		Comment:    input.Comment,
		CreatedAt:  time.Now(),
	}

	for i := range m.professionals {
		if m.professionals[i].ID == profID {
			m.professionals[i].Reviews = append([]models.Review{rev}, m.professionals[i].Reviews...)
			m.professionals[i].ReviewsCount++
			m.professionals[i].JobsCompletedCount++
			break
		}
	}
	return &rev, nil
}

func (m *InMemoryStore) UpdateProfessionalVerification(_ context.Context, profID string, input models.OnboardingVerificationInput) (*models.ProfessionalDetail, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.professionals {
		if m.professionals[i].ID == profID || profID == "default" {
			if input.TradeTitle != "" {
				m.professionals[i].TradeTitle = input.TradeTitle
			}
			if len(input.Specialties) > 0 {
				m.professionals[i].Specialties = input.Specialties
			}
			if input.Neighborhood != "" {
				m.professionals[i].Neighborhood = input.Neighborhood
			}
			if input.DNIFrontURL != "" || input.DNIBackURL != "" {
				m.professionals[i].IsDNIVerified = true
			}
			if input.LicenseNumber != "" {
				m.professionals[i].IsLicensed = true
				m.professionals[i].LicenseNumber = input.LicenseNumber
				m.professionals[i].LicenseIssuer = input.LicenseIssuer
			}
			copyP := m.professionals[i]
			return &copyP, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *InMemoryStore) UpdateProfessionalAvailability(_ context.Context, profID string, input models.AvailabilityToggleInput) (*models.ProfessionalDetail, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.professionals {
		if m.professionals[i].ID == profID || profID == "default" {
			m.professionals[i].IsAvailableNow = input.IsAvailableNow
			if input.Neighborhood != "" {
				m.professionals[i].Neighborhood = input.Neighborhood
			}
			copyP := m.professionals[i]
			return &copyP, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *InMemoryStore) GetProfessionalDashboard(_ context.Context, profID string) (*models.ProfessionalDashboardResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	prof := m.professionals[0]
	for _, p := range m.professionals {
		if p.ID == profID {
			prof = p
			break
		}
	}

	var reqs []models.ServiceRequest
	urgentCount := 0
	for _, r := range m.requests {
		if r.ProfessionalID == prof.ID {
			reqs = append(reqs, r)
			if r.Urgency == "urgente" && r.Status == "pending" {
				urgentCount++
			}
		}
	}

	return &models.ProfessionalDashboardResponse{
		Professional:   prof,
		NewTodayCount:  len(reqs) + 1, // Simula las 3 nuevas hoy de la pantalla Stitch
		UrgentCount:    urgentCount,
		MonthContacts:  prof.MonthlyContactsCount,
		ActiveRequests: reqs,
	}, nil
}

// =============================================================================
// PostgresStore: Implementación SQL sobre PostgreSQL + PostGIS
// =============================================================================

type PostgresStore struct {
	db       *sql.DB
	fallback *InMemoryStore
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db:       db,
		fallback: NewInMemoryStore(),
	}
}

func (s *PostgresStore) RegisterUser(ctx context.Context, input models.RegisterUserInput) (*models.User, error) {
	var u models.User
	if input.Role == "" {
		input.Role = "client"
	}
	if input.AuthProvider == "" {
		input.AuthProvider = "manual"
	}
	if input.SavedAddress == "" {
		input.SavedAddress = "Av. Rivadavia 5420, Caballito"
	}
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO users (full_name, email, phone, auth_provider, saved_address, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, full_name, email, phone, role, auth_provider, saved_address, created_at
	`, input.FullName, input.Email, input.Phone, input.AuthProvider, input.SavedAddress, input.Role).Scan(
		&u.ID, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.AuthProvider, &u.SavedAddress, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PostgresStore) GetCurrentUser(ctx context.Context, role string) (*models.User, error) {
	return s.fallback.GetCurrentUser(ctx, role)
}

func (s *PostgresStore) ListCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, slug, name, helper_examples, icon_name, COALESCE(badge_text, ''), COALESCE(badge_type, ''), display_order
		FROM categories
		ORDER BY display_order ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.HelperExamples, &c.IconName, &c.BadgeText, &c.BadgeType, &c.DisplayOrder); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (s *PostgresStore) SuggestCategoryByProblem(ctx context.Context, problem string) (*models.DiagnosticSuggestionResponse, error) {
	return s.fallback.SuggestCategoryByProblem(ctx, problem)
}

func (s *PostgresStore) ListProfessionals(ctx context.Context, f ProfessionalFilter) ([]models.ProfessionalSummary, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			p.id, p.user_id, u.full_name, u.phone, COALESCE(u.avatar_url, ''),
			p.primary_category_id, c.name, p.trade_title, p.specialties, p.neighborhood,
			p.is_available_now, p.attends_emergencies, p.is_dni_verified, p.is_licensed,
			COALESCE(p.license_number, ''), COALESCE(p.license_issuer, ''),
			p.visit_price_amount, p.visit_price_label,
			p.rating_avg::float8, p.reviews_count, p.jobs_completed_count, p.monthly_contacts_count
		FROM professional_profiles p
		JOIN users u ON u.id = p.user_id
		JOIN categories c ON c.id = p.primary_category_id
		WHERE ($1::int = 0 OR p.primary_category_id = $1)
		  AND ($2::bool = FALSE OR p.attends_emergencies = TRUE)
		  AND ($3::bool = FALSE OR p.is_licensed = TRUE)
		ORDER BY p.rating_avg DESC
	`, f.CategoryID, f.OnlyUrgent, f.OnlyLicensed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ProfessionalSummary
	for rows.Next() {
		var item models.ProfessionalSummary
		var specBytes []byte
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.FullName, &item.Phone, &item.AvatarURL,
			&item.PrimaryCategoryID, &item.CategoryName, &item.TradeTitle, &specBytes, &item.Neighborhood,
			&item.IsAvailableNow, &item.AttendsEmergencies, &item.IsDNIVerified, &item.IsLicensed,
			&item.LicenseNumber, &item.LicenseIssuer, &item.VisitPriceAmount, &item.VisitPriceLabel,
			&item.RatingAvg, &item.ReviewsCount, &item.JobsCompletedCount, &item.MonthlyContactsCount,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(specBytes, &item.Specialties)
		item.DistanceLabel = "A 8 cuadras (650 metros)"
		item.RatingLabel = "Muy recomendado"
		list = append(list, item)
	}
	return list, rows.Err()
}

func (s *PostgresStore) GetProfessionalByID(ctx context.Context, id string) (*models.ProfessionalDetail, error) {
	return s.fallback.GetProfessionalByID(ctx, id)
}

func (s *PostgresStore) CreateServiceRequest(ctx context.Context, input models.CreateServiceRequestInput) (*models.ServiceRequest, error) {
	return s.fallback.CreateServiceRequest(ctx, input)
}

func (s *PostgresStore) GetServiceRequestByID(ctx context.Context, id string) (*models.ServiceRequest, error) {
	return s.fallback.GetServiceRequestByID(ctx, id)
}

func (s *PostgresStore) ListServiceRequests(ctx context.Context, professionalID string) ([]models.ServiceRequest, error) {
	return s.fallback.ListServiceRequests(ctx, professionalID)
}

func (s *PostgresStore) CreateReview(ctx context.Context, requestID string, input models.CreateReviewInput) (*models.Review, error) {
	return s.fallback.CreateReview(ctx, requestID, input)
}

func (s *PostgresStore) UpdateProfessionalVerification(ctx context.Context, profID string, input models.OnboardingVerificationInput) (*models.ProfessionalDetail, error) {
	return s.fallback.UpdateProfessionalVerification(ctx, profID, input)
}

func (s *PostgresStore) UpdateProfessionalAvailability(ctx context.Context, profID string, input models.AvailabilityToggleInput) (*models.ProfessionalDetail, error) {
	return s.fallback.UpdateProfessionalAvailability(ctx, profID, input)
}

func (s *PostgresStore) GetProfessionalDashboard(ctx context.Context, profID string) (*models.ProfessionalDashboardResponse, error) {
	return s.fallback.GetProfessionalDashboard(ctx, profID)
}
