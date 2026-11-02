package models

import "time"

// Category representa cada rubro de la pantalla "cliente_inicio_y_categor_as".
type Category struct {
	ID             int    `json:"id"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	HelperExamples string `json:"helper_examples"`
	IconName       string `json:"icon_name"`
	BadgeText      string `json:"badge_text,omitempty"` // Ej: "Matriculados" o "Urgencias"
	BadgeType      string `json:"badge_type,omitempty"` // "verified" | "urgent"
	DisplayOrder   int    `json:"display_order"`
}

// RegisterUserInput representa el formulario de "registro_google_o_manual".
type RegisterUserInput struct {
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	Role         string `json:"role"`          // "client" ("Busco oficios") | "professional" ("Soy profesional")
	AuthProvider string `json:"auth_provider"` // "manual" | "google"
	SavedAddress string `json:"saved_address,omitempty"`
}

// User representa un vecino o profesional autenticado.
type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`
	AuthProvider string    `json:"auth_provider"`
	SavedAddress string    `json:"saved_address"`
	CreatedAt    time.Time `json:"created_at"`
}

// PortfolioItem representa una foto real en "cliente_perfil_del_profesional_y_contacto".
type PortfolioItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	ImageURL   string `json:"image_url"`
	IsVerified bool   `json:"is_verified"`
}

// Review representa una opinión verificada y las 5 categorías destacables de "cliente_trabajo_terminado_y_calificaci_n".
type Review struct {
	ID         string    `json:"id"`
	AuthorName string    `json:"author_name"`
	AuthorZone string    `json:"author_zone"`
	Rating     int       `json:"rating"`
	Highlights []string  `json:"highlights"` // ["Puntualidad", "Precio justo", "Honestidad y prolijidad", "Excelente atención", "Solución en el acto"]
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProfessionalSummary es la tarjeta que se muestra en "cliente_resultados_de_profesionales".
type ProfessionalSummary struct {
	ID                   string   `json:"id"`
	UserID               string   `json:"user_id"`
	FullName             string   `json:"full_name"`
	Phone                string   `json:"phone"`
	AvatarURL            string   `json:"avatar_url"`
	PrimaryCategoryID    int      `json:"primary_category_id"`
	CategoryName         string   `json:"category_name"`
	TradeTitle           string   `json:"trade_title"`
	Specialties          []string `json:"specialties"`
	Neighborhood         string   `json:"neighborhood"`
	DistanceLabel        string   `json:"distance_label"` // Ej: "A 8 cuadras (650 metros)"
	DistanceKm           float64  `json:"distance_km"`
	IsAvailableNow       bool     `json:"is_available_now"`
	AttendsEmergencies   bool     `json:"attends_emergencies"`
	IsDNIVerified        bool     `json:"is_dni_verified"`
	IsLicensed           bool     `json:"is_licensed"`
	LicenseNumber        string   `json:"license_number,omitempty"`
	LicenseIssuer        string   `json:"license_issuer,omitempty"`
	VisitPriceAmount     int      `json:"visit_price_amount"`
	VisitPriceLabel      string   `json:"visit_price_label"`
	RatingAvg            float64  `json:"rating_avg"`
	RatingLabel          string   `json:"rating_label"` // Ej: "Muy recomendado" / "Excelente atención"
	ReviewsCount         int      `json:"reviews_count"`
	JobsCompletedCount   int      `json:"jobs_completed_count"`
	MonthlyContactsCount int      `json:"monthly_contacts_count"`
}

// ProfessionalDetail es la ficha completa de "cliente_perfil_del_profesional_y_contacto".
type ProfessionalDetail struct {
	ProfessionalSummary
	Bio                 string          `json:"bio"`
	YearsExperience     int             `json:"years_experience"`
	GuaranteePercentage int             `json:"guarantee_percentage"`
	Portfolio           []PortfolioItem `json:"portfolio"`
	Reviews             []Review        `json:"reviews"`
}

// OnboardingVerificationInput representa los datos de "profesional_registro_y_verificaci_n_dni".
type OnboardingVerificationInput struct {
	PrimaryCategoryID int      `json:"primary_category_id"`
	TradeTitle        string   `json:"trade_title"`
	Specialties       []string `json:"specialties"`
	Neighborhood      string   `json:"neighborhood"`
	CoverageRadiusKm  int      `json:"coverage_radius_km"`
	DNIFrontURL       string   `json:"dni_front_url"`
	DNIBackURL        string   `json:"dni_back_url"`
	LicenseNumber     string   `json:"license_number,omitempty"`
	LicenseIssuer     string   `json:"license_issuer,omitempty"`
	LicenseDocURL     string   `json:"license_doc_url,omitempty"`
}

// AvailabilityToggleInput cambia el estado de disponibilidad en "profesional_panel_y_solicitudes_recibidas".
type AvailabilityToggleInput struct {
	IsAvailableNow bool   `json:"is_available_now"`
	Neighborhood   string `json:"neighborhood,omitempty"`
}

// CreateServiceRequestInput representa "cliente_solicitud_guiada_a_roberto".
type CreateServiceRequestInput struct {
	ClientName         string `json:"client_name"`
	ClientPhone        string `json:"client_phone"`
	ClientAddress      string `json:"client_address"`
	ProfessionalID     string `json:"professional_id"`
	CategoryID         int    `json:"category_id"`
	Urgency            string `json:"urgency"` // "normal" | "urgente"
	ProblemDescription string `json:"problem_description"`
	PhotoURL           string `json:"photo_url,omitempty"`
	ContactChannel     string `json:"contact_channel"` // "whatsapp" | "phone_call"
}

// ServiceRequest representa un pedido en "cliente_elegir_contacto_y_confirmaci_n" y "profesional_panel_y_solicitudes_recibidas".
type ServiceRequest struct {
	ID                  string    `json:"id"`
	ClientName          string    `json:"client_name"`
	ClientInitials      string    `json:"client_initials"`
	ClientPhone         string    `json:"client_phone"`
	ClientAddress       string    `json:"client_address"`
	ClientDistanceLabel string    `json:"client_distance_label"`
	ProfessionalID      string    `json:"professional_id"`
	ProfessionalName    string    `json:"professional_name"`
	ProfessionalPhone   string    `json:"professional_phone"`
	ProfessionalLicense string    `json:"professional_license"`
	CategoryID          int       `json:"category_id"`
	Urgency             string    `json:"urgency"` // "normal" | "urgente"
	ProblemDescription  string    `json:"problem_description"`
	PhotoURL            string    `json:"photo_url,omitempty"`
	ContactChannel      string    `json:"contact_channel"`
	Status              string    `json:"status"` // "pending" | "scheduled" | "completed"
	ScheduledNote       string    `json:"scheduled_note,omitempty"`
	TimeAgoLabel        string    `json:"time_ago_label"`
	CreatedAt           time.Time `json:"created_at"`
}

// ProfessionalDashboardResponse alimenta la pantalla "profesional_panel_y_solicitudes_recibidas".
type ProfessionalDashboardResponse struct {
	Professional   ProfessionalDetail `json:"professional"`
	NewTodayCount  int                `json:"new_today_count"`
	UrgentCount    int                `json:"urgent_count"`
	MonthContacts  int                `json:"month_contacts"`
	ActiveRequests []ServiceRequest   `json:"active_requests"`
}

// CreateReviewInput representa el formulario de "cliente_trabajo_terminado_y_calificaci_n".
type CreateReviewInput struct {
	AuthorName string   `json:"author_name"`
	AuthorZone string   `json:"author_zone"`
	Rating     int      `json:"rating"`
	Highlights []string `json:"highlights"`
	Comment    string   `json:"comment"`
}

// DiagnosticSuggestionResponse responde al asistente "¿No sabés qué profesional necesitás?" del Home.
type DiagnosticSuggestionResponse struct {
	ProblemText       string   `json:"problem_text"`
	SuggestedCategory Category `json:"suggested_category"`
	Explanation       string   `json:"explanation"`
}
