-- =============================================================================
-- Migración 003: Semilla de Profesionales en Mendoza, San Luis y Luján (GDO-14)
-- =============================================================================

-- Contraseña por defecto para profesionales semilla: "123456" ($2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS)
INSERT INTO users (id, full_name, email, phone, password_hash, auth_provider, saved_address, saved_lat, saved_lng, role) VALUES
-- Mendoza
('33333333-3333-3333-3333-333333333331', 'Carlos Alberto Morales', 'carlos.morales.mza@guiaoficios.ar', '+54 9 261 555-1122', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'Rivadavia y 25 de Mayo, Palmira, Mendoza', -33.054634, -68.5612924, 'professional'),
('33333333-3333-3333-3333-333333333332', 'Gustavo Daniel Fernández', 'gustavo.fernandez.mza@guiaoficios.ar', '+54 9 261 444-3322', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'Bandera de los Andes 450, Guaymallén, Mendoza', -32.895000, -68.825000, 'professional'),

-- San Luis (Juana Koslay y Gran San Luis)
('44444444-4444-4444-4444-444444444441', 'Héctor Ramón Lucero', 'hector.lucero.sl@guiaoficios.ar', '+54 9 266 455-6677', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'Av. del Viento Chorrillero 1250, Juana Koslay, San Luis', -33.2662423, -66.2337363, 'professional'),
('44444444-4444-4444-4444-444444444442', 'Walter Gómez', 'walter.gomez.sl@guiaoficios.ar', '+54 9 266 477-8899', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'Los Zorzales 300, Juana Koslay, San Luis', -33.265500, -66.234000, 'professional'),

-- Buenos Aires (Luján)
('55555555-5555-5555-5555-555555555551', 'Esteban Darío Peralta', 'esteban.peralta.lujan@guiaoficios.ar', '+54 9 2323 55-4433', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'San Martín 750, Luján, Buenos Aires', -34.570300, -59.105000, 'professional'),
('55555555-5555-5555-5555-555555555552', 'Mariano Rossi', 'mariano.rossi.lujan@guiaoficios.ar', '+54 9 2323 66-7788', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'Colón 420, Luján, Buenos Aires', -34.568000, -59.112000, 'professional')
ON CONFLICT (id) DO NOTHING;

-- Perfiles Profesionales
INSERT INTO professional_profiles (
    id, user_id, primary_category_id, trade_title, specialties, bio, years_experience,
    guarantee_percentage, neighborhood, coverage_radius_km, base_location,
    is_available_now, attends_emergencies, is_dni_verified, is_licensed,
    license_number, license_issuer, verification_status, visit_price_amount, visit_price_label,
    rating_avg, reviews_count, jobs_completed_count, monthly_contacts_count
) VALUES
-- Mendoza: Carlos Morales (Plomero y Gasista)
(
    'cccc3333-3333-3333-3333-333333333331',
    '33333333-3333-3333-3333-333333333331',
    1,
    'Plomero y Gasista Matriculado',
    '["Urgencias 24hs", "Pérdidas de agua y gas", "Calefones y termotanques"]'::jsonb,
    'Oficio familiar con más de 20 años de trayectoria en Palmira, San Martín y zona este de Mendoza. Atención rápida, presupuestos transparentes y trabajos garantizados.',
    22,
    100,
    'Palmira, San Martín y Gran Mendoza',
    35,
    ST_SetSRID(ST_MakePoint(-68.5612924, -33.054634), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula Ecogas N° 34.209', 'Ecogas / Enargas', 'verified', 12000, '$12.000',
    4.92, 94, 94, 15
),
-- Mendoza: Gustavo Fernández (Electricista)
(
    'cccc3333-3333-3333-3333-333333333332',
    '33333333-3333-3333-3333-333333333332',
    2,
    'Electricista Matriculado e Iluminación',
    '["Cortocircuitos", "Tableros eléctricos", "Instalaciones monofásicas y trifásicas"]'::jsonb,
    'Técnico electricista matriculado. Instalaciones domiciliarias, comerciales y solución inmediata a cortes de luz en Gran Mendoza.',
    16,
    100,
    'Gran Mendoza, Guaymallén y Maipú',
    30,
    ST_SetSRID(ST_MakePoint(-68.825000, -32.895000), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'COPIG Matrícula N° 8.741', 'COPIG Mendoza', 'verified', 14000, '$14.000',
    4.95, 118, 118, 22
),

-- San Luis: Héctor Lucero (Electricista y Plomero)
(
    'cccc4444-4444-4444-4444-444444444441',
    '44444444-4444-4444-4444-444444444441',
    2,
    'Electricista de Urgencias y Plomero',
    '["Térmicas y disyuntores", "Instalaciones de luz", "Urgencias en el día"]'::jsonb,
    'Vecino de Juana Koslay. Reparaciones eléctricas de urgencia, iluminación y plomería básica. Puntualidad y honestidad.',
    15,
    100,
    'Juana Koslay, El Volcán y San Luis',
    25,
    ST_SetSRID(ST_MakePoint(-66.2337363, -33.2662423), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula Profesional N° 5.120', 'Colegio de Técnicos San Luis', 'verified', 13000, '$13.000',
    4.88, 73, 73, 11
),
-- San Luis: Walter Gómez (Plomero)
(
    'cccc4444-4444-4444-4444-444444444442',
    '44444444-4444-4444-4444-444444444442',
    1,
    'Plomero y Destapaciones',
    '["Destapaciones con máquina", "Cambio de griferías", "Pérdidas de agua"]'::jsonb,
    'Servicio integral de plomería y desagües en Juana Koslay y San Luis Capital. Maquinaria propia y garantía escrita.',
    12,
    100,
    'Juana Koslay y alrededores',
    20,
    ST_SetSRID(ST_MakePoint(-66.234000, -33.265500), 4326)::geography,
    TRUE, TRUE, TRUE, FALSE,
    NULL, NULL, 'verified', 11000, '$11.000',
    4.85, 52, 52, 8
),

-- Luján: Esteban Peralta (Gasista y Plomero)
(
    'cccc5555-5555-5555-5555-555555555551',
    '55555555-5555-5555-5555-555555555551',
    3,
    'Gasista Matriculado y Plomero',
    '["Habilitaciones de gas", "Pruebas de hermeticidad", "Calefones y calderas"]'::jsonb,
    'Gasista matriculado de primera categoría en Luján. Trámites ante Naturgy, colocación de artefactos y cañerías de agua por termofusión.',
    18,
    100,
    'Luján, Jáuregui, Olivera y Cortínez',
    25,
    ST_SetSRID(ST_MakePoint(-59.105000, -34.570300), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula Naturgy N° 61.428', 'Naturgy BAN', 'verified', 13500, '$13.500',
    4.93, 106, 106, 19
),
-- Luján: Mariano Rossi (Cerrajero)
(
    'cccc5555-5555-5555-5555-555555555552',
    '55555555-5555-5555-5555-555555555552',
    8,
    'Cerrajero 24hs Urgencias',
    '["Aperturas domiciliarias", "Cambio de combinaciones", "Cerraduras de seguridad"]'::jsonb,
    'Atención inmediata las 24 horas en Luján y partidos vecinos. Llegamos en 20 a 35 minutos para aperturas de puertas trabadas.',
    14,
    100,
    'Luján y alrededores',
    25,
    ST_SetSRID(ST_MakePoint(-59.112000, -34.568000), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Cámara de Cerrajeros Reg. 941', 'CAECEBA', 'verified', 15000, '$15.000',
    4.90, 87, 87, 14
)
ON CONFLICT (id) DO NOTHING;
