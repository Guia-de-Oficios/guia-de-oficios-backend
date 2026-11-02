-- =============================================================================
-- Guía de Oficios — Esquema Completo de Base de Datos (Adaptado a las 10 Pantallas Stitch)
-- Motor: PostgreSQL 16 + PostGIS
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- 1. Categorías / Oficios (Pantalla: cliente_inicio_y_categor_as)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(80) NOT NULL,
    helper_examples VARCHAR(160) NOT NULL,
    icon_name VARCHAR(50) NOT NULL,
    badge_text VARCHAR(40),                 -- Ej: "Matriculados" o "Urgencias"
    badge_type VARCHAR(20),                 -- 'verified' | 'urgent' | NULL
    display_order INT NOT NULL DEFAULT 0
);

-- 2. Usuarios (Pantallas: bienvenida_y_selecci_n_de_rol, registro_google_o_manual)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(120) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(40) NOT NULL,
    password_hash TEXT,
    auth_provider VARCHAR(20) NOT NULL DEFAULT 'manual', -- 'manual' | 'google'
    avatar_url TEXT,
    saved_address VARCHAR(200) DEFAULT 'Av. Rivadavia 5420, Caballito',
    role VARCHAR(20) NOT NULL DEFAULT 'client',          -- 'client' | 'professional' | 'admin'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Perfiles de Profesionales (Pantallas: cliente_resultados_de_profesionales, cliente_perfil_del_profesional_y_contacto, profesional_registro_y_verificaci_n_dni, profesional_panel_y_solicitudes_recibidas)
CREATE TABLE IF NOT EXISTS professional_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    primary_category_id INT NOT NULL REFERENCES categories(id),
    trade_title VARCHAR(120) NOT NULL,
    specialties JSONB NOT NULL DEFAULT '[]'::jsonb,       -- Ej: ["Destapaciones", "Calefones y Termotanques"]
    bio TEXT NOT NULL,
    years_experience INT NOT NULL DEFAULT 10,
    guarantee_percentage INT NOT NULL DEFAULT 100,
    neighborhood VARCHAR(120) NOT NULL,                   -- Ej: "Caballito y Flores"
    coverage_radius_km INT NOT NULL DEFAULT 5,
    base_location GEOGRAPHY(Point, 4326),
    is_available_now BOOLEAN NOT NULL DEFAULT TRUE,       -- Toggle "Estás DISPONIBLE" / "Pausar disponibilidad"
    attends_emergencies BOOLEAN NOT NULL DEFAULT FALSE,   -- "Atiende urgencias hoy"
    is_dni_verified BOOLEAN NOT NULL DEFAULT FALSE,
    dni_front_url TEXT,
    dni_back_url TEXT,
    is_licensed BOOLEAN NOT NULL DEFAULT FALSE,
    license_number VARCHAR(60),                           -- Ej: "Matrícula N° 48.912"
    license_issuer VARCHAR(100),                          -- Ej: "Metrogas / COPIME"
    license_doc_url TEXT,
    verification_status VARCHAR(30) NOT NULL DEFAULT 'verified', -- 'pending_review' | 'verified'
    visit_price_amount INT NOT NULL DEFAULT 12000,
    visit_price_label VARCHAR(120) NOT NULL DEFAULT '$12.000',
    rating_avg NUMERIC(3,2) NOT NULL DEFAULT 5.00,
    reviews_count INT NOT NULL DEFAULT 0,
    jobs_completed_count INT NOT NULL DEFAULT 0,
    monthly_contacts_count INT NOT NULL DEFAULT 18,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prof_category ON professional_profiles(primary_category_id);
CREATE INDEX IF NOT EXISTS idx_prof_location ON professional_profiles USING GIST(base_location);

-- 4. Trabajos Realizados (Galería en cliente_perfil_del_profesional_y_contacto)
CREATE TABLE IF NOT EXISTS portfolio_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    professional_id UUID NOT NULL REFERENCES professional_profiles(id) ON DELETE CASCADE,
    title VARCHAR(140) NOT NULL,
    subtitle VARCHAR(140) NOT NULL,
    image_url TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Solicitudes de Servicio (Pantallas: cliente_solicitud_guiada_a_roberto, cliente_elegir_contacto_y_confirmaci_n, profesional_panel_y_solicitudes_recibidas)
CREATE TABLE IF NOT EXISTS service_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID REFERENCES users(id) ON DELETE SET NULL,
    client_name VARCHAR(120) NOT NULL,
    client_phone VARCHAR(40) NOT NULL,
    client_address VARCHAR(180) NOT NULL DEFAULT 'Av. Rivadavia 5420, Caballito',
    client_distance_label VARCHAR(60) NOT NULL DEFAULT 'a 600m',
    professional_id UUID NOT NULL REFERENCES professional_profiles(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id),
    urgency VARCHAR(30) NOT NULL DEFAULT 'normal',        -- 'normal' | 'urgente'
    problem_description TEXT NOT NULL,
    photo_url TEXT,
    contact_channel VARCHAR(30) NOT NULL DEFAULT 'whatsapp', -- 'whatsapp' | 'phone_call'
    status VARCHAR(30) NOT NULL DEFAULT 'pending',        -- 'pending' | 'scheduled' | 'completed'
    scheduled_note VARCHAR(140),                          -- Ej: "Trabajo coordinado para el sábado"
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 6. Reseñas y Calificaciones (Pantalla: cliente_trabajo_terminado_y_calificaci_n)
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID UNIQUE REFERENCES service_requests(id) ON DELETE SET NULL,
    professional_id UUID NOT NULL REFERENCES professional_profiles(id) ON DELETE CASCADE,
    author_name VARCHAR(100) NOT NULL,
    author_Zone VARCHAR(100) NOT NULL DEFAULT 'Vecina de Caballito',
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    highlights JSONB NOT NULL DEFAULT '[]'::jsonb,        -- ["Puntualidad", "Precio justo", "Honestidad y prolijidad", "Excelente atención", "Solución en el acto"]
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
