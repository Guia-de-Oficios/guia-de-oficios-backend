-- =============================================================================
-- Guía de Oficios — Datos Iniciales Reales en PostgreSQL + PostGIS (GDO-9)
-- =============================================================================

INSERT INTO categories (id, slug, name, helper_examples, icon_name, badge_text, badge_type, display_order) VALUES
(1, 'plomero', 'Plomero', 'Canillas, cañerías y pérdidas', 'water_drop', NULL, NULL, 1),
(2, 'electricista', 'Electricista', 'Enchufes, cables e iluminación', 'bolt', NULL, NULL, 2),
(3, 'gasista', 'Gasista', 'Calefones, estufas y conexiones', 'local_fire_department', 'Matriculados', 'verified', 3),
(4, 'aire-acondicionado', 'Aire acondicionado', 'Equipos y mantenimiento', 'mode_fan', NULL, NULL, 4),
(5, 'pintor', 'Pintor', 'Paredes y terminaciones', 'format_paint', NULL, NULL, 5),
(6, 'albanil', 'Albañil', 'Reparaciones y construcción', 'handyman', NULL, NULL, 6),
(7, 'jardinero', 'Jardinero', 'Plantas y mantenimiento', 'yard', NULL, NULL, 7),
(8, 'cerrajero', 'Cerrajero', 'Cerraduras y puertas', 'key', 'Urgencias', 'urgent', 8)
ON CONFLICT (id) DO NOTHING;

-- Usuarios iniciales (contraseña por defecto: "123456" con hash bcrypt costo 10)
INSERT INTO users (id, full_name, email, phone, password_hash, auth_provider, saved_address, saved_lat, saved_lng, role) VALUES
('00000000-0000-0000-0000-000000000001', 'María Elena Castro', 'maria.castro@vecinos.ar', '11 5544-3322', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'Av. Rivadavia 5420, Caballito', -34.6186, -58.4368, 'client'),
('11111111-1111-1111-1111-111111111111', 'Roberto Carlos Gómez', 'roberto.gomez@guiaoficios.ar', '+5491145890000', '$2a$10$X9j0kHu0e7UmDvcemNqE5ujNhhKppb6coqyKnnkjJKJAMJ2Y8O1D.', 'manual', 'Caballito, CABA', -34.6150, -58.4330, 'professional'),
('22222222-2222-2222-2222-222222222222', 'Martín Edgardo Silva', 'martin.silva@guiaoficios.ar', '+5491155667788', '$2a$10$DYXz4hs.721eXHdkzox5je3eEtaaxRWe2B3Oq6XfMH4HXbmKr.9fO', 'manual', 'Caballito, CABA', -34.6110, -58.4420, 'professional')
ON CONFLICT (id) DO UPDATE
SET password_hash = CASE
    WHEN users.password_hash IS NULL OR users.password_hash NOT LIKE '$2%'
    THEN EXCLUDED.password_hash
    ELSE users.password_hash
END;

-- Perfiles de Profesionales
INSERT INTO professional_profiles (
    id, user_id, primary_category_id, trade_title, specialties, bio, years_experience,
    guarantee_percentage, neighborhood, coverage_radius_km, base_location,
    is_available_now, attends_emergencies, is_dni_verified, is_licensed,
    license_number, license_issuer, verification_status, visit_price_amount, visit_price_label,
    rating_avg, reviews_count, jobs_completed_count, monthly_contacts_count
) VALUES
(
    'aaaa1111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    1,
    'Plomero y Gasista Matriculado',
    '["Destapaciones", "Calefones y Termotanques"]'::jsonb,
    'Hola, soy Roberto. Trabajo hace más de 25 años haciendo reparaciones en casas y departamentos de la zona. Trabajo prolijo, puntual y con garantía en cada arreglo.',
    25,
    100,
    'Caballito y Flores',
    5,
    ST_SetSRID(ST_MakePoint(-58.4330, -34.6150), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula oficial N° 48.912 verificada', 'Metrogas / ENARGAS', 'verified', 12000, '$12.000',
    4.90, 142, 142, 18
),
(
    'bbbb2222-2222-2222-2222-222222222222',
    '22222222-2222-2222-2222-222222222222',
    1,
    'Plomero y Gasista',
    '["Urgencias 24hs", "Filtraciones"]'::jsonb,
    'Especialista en reparaciones rápidas en el día, pérdidas de agua y conexiones de artefactos en toda la zona de Caballito.',
    14,
    100,
    'Caballito y Parque Chacabuco',
    5,
    ST_SetSRID(ST_MakePoint(-58.4420, -34.6110), 4326)::geography,
    TRUE, TRUE, TRUE, FALSE,
    NULL, NULL, 'verified', 10000, '$10.000',
    4.80, 89, 89, 12
)
ON CONFLICT (id) DO NOTHING;

-- Trabajos realizados de Roberto Carlos Gómez
INSERT INTO portfolio_items (id, professional_id, title, subtitle, image_url, is_verified) VALUES
('dddd1111-1111-1111-1111-111111111111', 'aaaa1111-1111-1111-1111-111111111111', 'Canilla monocomando cambiada', 'Instalación y sellado completo', 'https://lh3.googleusercontent.com/aida-public/AB6AXuCh10RiTc50jP69tmrjG8REQ7uQZftm52vnIAc3RMhe3SukIRmt1qsbmlKpf6XUc585FmJriOr-k8Z0QvUfu8A7nNCkOGW2jJx3avdmiKF-Nkrta-clOO_eVsmwIDHmRz30tgPahuYwFtN8PQ1fyh9hPzU65q3-IQsxo-p-0lExk7Sal_lolRdoZ-44c8UC0ZDRI9Y7OqNDb2WZlr5Ly30L2CAZmaoIXjO55wdLSbRG', TRUE),
('dddd2222-2222-2222-2222-222222222222', 'aaaa1111-1111-1111-1111-111111111111', 'Reparación de cañería termofusión', 'Solución definitiva a filtración', 'https://lh3.googleusercontent.com/aida-public/AB6AXuCTODIWHLBehXm6SVzYW9pLhJ4C1zuw4YW31S21vmoWs31kV-fwgGBcwNkWLWT8eQAC0le8E2mTVOCpr78K1LSgHxq1Wz_TIbADpF91xbBHUTJEVtBkUGHSa_saUCNhM56oSA0exrcwQueuH6FVaBDeEcWP8z0qj_SvfAmVejvoezsktrrA52SZpDaRBLiNFN1lS3RBFovpPsOpuNUk6SDaVRHc2VbiBDJKWl6kekvw', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Solicitudes recibidas en el Panel Profesional de Roberto
INSERT INTO service_requests (
    id, client_id, client_name, client_phone, client_address, client_distance_label,
    professional_id, category_id, urgency, problem_description, photo_url, status, scheduled_note
) VALUES
(
    'eeee1111-1111-1111-1111-111111111111',
    '00000000-0000-0000-0000-000000000001',
    'María Elena Castro',
    '11 5544-3322',
    'Av. Rivadavia 5420',
    'a 600m',
    'aaaa1111-1111-1111-1111-111111111111',
    1,
    'urgente',
    'Pierde agua la canilla de la cocina debajo de la mesada. Gotea constante y cerré la llave de paso provisoria.',
    'Foto_caño_cocina.jpg',
    'pending',
    NULL
),
(
    'eeee2222-2222-2222-2222-222222222222',
    NULL,
    'Jorge Benítez',
    '11 6677-8899',
    'Caballito Norte',
    'a 1.4 km',
    'aaaa1111-1111-1111-1111-111111111111',
    1,
    'normal',
    'Cambio de fuelle y flexibles en baño principal.',
    NULL,
    'scheduled',
    'Trabajo coordinado para el sábado'
)
ON CONFLICT (id) DO NOTHING;

-- Reseñas de vecinos en el perfil de Roberto
INSERT INTO reviews (id, professional_id, author_name, author_zone, rating, highlights, comment) VALUES
(
    'ffff1111-1111-1111-1111-111111111111',
    'aaaa1111-1111-1111-1111-111111111111',
    'Graciela M.',
    'Vecina de Caballito (Av. Pedro Goyena)',
    5,
    '["Puntualidad", "Precio justo", "Honestidad y prolijidad", "Excelente atención", "Solución en el acto"]'::jsonb,
    'Vino puntual a la hora acordada, resolvió la pérdida del baño en el acto y dejó todo impecable. Muy respetuoso.'
)
ON CONFLICT (id) DO NOTHING;
