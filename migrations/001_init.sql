-- =============================================================================
-- Guía de Oficios — Esquema Definitivo PostgreSQL 16 + PostGIS (GDO-9)
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- 1. Categorías / Oficios (Pantalla 3: cliente_inicio_y_categor_as)
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

-- 2. Usuarios (Pantallas 1 y 2: bienvenida_y_selecci_n_de_rol, registro_google_o_manual)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(120) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(40) NOT NULL,
    password_hash TEXT,
    google_sub VARCHAR(120) UNIQUE,
    auth_provider VARCHAR(20) NOT NULL DEFAULT 'manual', -- 'manual' | 'google'
    avatar_url TEXT,
    saved_address VARCHAR(200) NOT NULL DEFAULT 'Av. Rivadavia 5420, Caballito',
    saved_lat DOUBLE PRECISION NOT NULL DEFAULT -34.6186,
    saved_lng DOUBLE PRECISION NOT NULL DEFAULT -58.4368,
    role VARCHAR(20) NOT NULL DEFAULT 'client',          -- 'client' | 'professional' | 'admin'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Perfiles de Profesionales (Pantallas 4, 5, 9 y 10)
CREATE TABLE IF NOT EXISTS professional_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    primary_category_id INT NOT NULL REFERENCES categories(id),
    trade_title VARCHAR(120) NOT NULL,
    specialties JSONB NOT NULL DEFAULT '[]'::jsonb,
    bio TEXT NOT NULL,
    years_experience INT NOT NULL DEFAULT 10,
    guarantee_percentage INT NOT NULL DEFAULT 100,
    neighborhood VARCHAR(120) NOT NULL,
    coverage_radius_km INT NOT NULL DEFAULT 5,
    base_location GEOGRAPHY(Point, 4326),
    is_available_now BOOLEAN NOT NULL DEFAULT TRUE,
    attends_emergencies BOOLEAN NOT NULL DEFAULT FALSE,
    is_dni_verified BOOLEAN NOT NULL DEFAULT FALSE,
    dni_front_url TEXT,
    dni_back_url TEXT,
    is_licensed BOOLEAN NOT NULL DEFAULT FALSE,
    license_number VARCHAR(80),
    license_issuer VARCHAR(100),
    license_doc_url TEXT,
    verification_status VARCHAR(30) NOT NULL DEFAULT 'verified',
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

-- 4. Trabajos Realizados (Galería en Pantalla 5)
CREATE TABLE IF NOT EXISTS portfolio_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    professional_id UUID NOT NULL REFERENCES professional_profiles(id) ON DELETE CASCADE,
    title VARCHAR(140) NOT NULL,
    subtitle VARCHAR(140) NOT NULL,
    image_url TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Solicitudes de Servicio (Pantallas 6, 7 y 10)
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
    contact_channel VARCHAR(30) NOT NULL DEFAULT 'whatsapp',
    status VARCHAR(30) NOT NULL DEFAULT 'pending',        -- 'pending' | 'scheduled' | 'completed'
    scheduled_note VARCHAR(140),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_requests_prof ON service_requests(professional_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_requests_client ON service_requests(client_id, created_at DESC);

-- 6. Reseñas y Calificaciones (Pantalla 8: cliente_trabajo_terminado_y_calificaci_n)
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID UNIQUE REFERENCES service_requests(id) ON DELETE SET NULL,
    professional_id UUID NOT NULL REFERENCES professional_profiles(id) ON DELETE CASCADE,
    author_name VARCHAR(100) NOT NULL,
    author_zone VARCHAR(100) NOT NULL DEFAULT 'Vecina Verificada de Caballito',
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    highlights JSONB NOT NULL DEFAULT '[]'::jsonb,
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_prof ON reviews(professional_id, created_at DESC);
