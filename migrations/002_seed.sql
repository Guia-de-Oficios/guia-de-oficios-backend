-- =============================================================================
-- Guía de Oficios — Datos Iniciales Reales en PostgreSQL + PostGIS (GDO-9 / BE-5.1)
-- Cubre los 8 rubros oficiales con al menos 16 profesionales verificados,
-- avatares, portfolios con imágenes reales y reseñas barriales.
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
INSERT INTO users (id, full_name, email, phone, password_hash, auth_provider, avatar_url, saved_address, saved_lat, saved_lng, role) VALUES
('00000000-0000-0000-0000-000000000001', 'María Elena Castro', 'maria.castro@vecinos.ar', '11 5544-3322', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=240', 'Av. Rivadavia 5420, Caballito', -34.6186, -58.4368, 'client'),
-- 1. Plomeros
('11111111-1111-1111-1111-111111111111', 'Roberto Carlos Gómez', 'roberto.gomez@guiaoficios.ar', '+5491145890000', '$2a$10$X9j0kHu0e7UmDvcemNqE5ujNhhKppb6coqyKnnkjJKJAMJ2Y8O1D.', 'manual', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=240', 'Av. Acoyte 240, Caballito, CABA', -34.6150, -58.4330, 'professional'),
('22222222-2222-2222-2222-222222222222', 'Martín Edgardo Silva', 'martin.silva@guiaoficios.ar', '+5491155667788', '$2a$10$DYXz4hs.721eXHdkzox5je3eEtaaxRWe2B3Oq6XfMH4HXbmKr.9fO', 'manual', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=240', 'Av. Directorio 1420, Parque Chacabuco, CABA', -34.6295, -58.4410, 'professional'),
-- 2. Electricistas
('11111111-1111-1111-1111-111111111103', 'Jorge Luis Pereyra', 'jorge.pereyra@guiaoficios.ar', '+5491144112233', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=240', 'Av. Díaz Vélez 4850, Caballito, CABA', -34.6085, -58.4355, 'professional'),
('11111111-1111-1111-1111-111111111104', 'Luciano Ezequiel Ríos', 'luciano.rios@guiaoficios.ar', '+5491144223344', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=240', 'Av. Carabobo 320, Flores, CABA', -34.6265, -58.4580, 'professional'),
-- 3. Gasistas
('11111111-1111-1111-1111-111111111105', 'Horacio Daniel Ferreyra', 'horacio.ferreyra@guiaoficios.ar', '+5491144334455', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=240', 'Av. Pedro Goyena 890, Caballito, CABA', -34.6230, -58.4385, 'professional'),
('11111111-1111-1111-1111-111111111106', 'Néstor Fabián Domínguez', 'nestor.dominguez@guiaoficios.ar', '+5491144445566', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1522075469751-3a6694fb2f61?w=240', 'Av. Medrano 640, Almagro, CABA', -34.6060, -58.4215, 'professional'),
-- 4. Aire acondicionado
('11111111-1111-1111-1111-111111111107', 'Diego Sebastián Molina', 'diego.molina@guiaoficios.ar', '+5491144556677', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=240', 'Neuquén 1150, Caballito, CABA', -34.6145, -58.4450, 'professional'),
('11111111-1111-1111-1111-111111111108', 'Federico Andrés Paz', 'federico.paz@guiaoficios.ar', '+5491144667788', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?w=240', 'Av. Corrientes 4520, Villa Crespo, CABA', -34.6015, -58.4310, 'professional'),
-- 5. Pintores
('11111111-1111-1111-1111-111111111109', 'Rubén Oscar Maidana', 'ruben.maidana@guiaoficios.ar', '+5491144778899', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1501196354995-cbb51c65aaea?w=240', 'Rosario 540, Caballito, CABA', -34.6198, -58.4342, 'professional'),
('11111111-1111-1111-1111-111111111110', 'Claudio Marcelo Vera', 'claudio.vera@guiaoficios.ar', '+5491144889900', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=240', 'Av. Boedo 980, Boedo, CABA', -34.6248, -58.4162, 'professional'),
-- 6. Albañiles
('11111111-1111-1111-1111-111111111112', 'Ramón Antonio Giménez', 'ramon.gimenez@guiaoficios.ar', '+5491144990011', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1560250097-0b93528c311a?w=240', 'Av. La Plata 620, Caballito, CABA', -34.6215, -58.4265, 'professional'),
('11111111-1111-1111-1111-111111111113', 'Víctor Hugo Ledesma', 'victor.ledesma@guiaoficios.ar', '+5491145001122', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1570295999919-56ceb5ecca61?w=240', 'Av. Asamblea 920, Parque Chacabuco, CABA', -34.6335, -58.4395, 'professional'),
-- 7. Jardineros
('11111111-1111-1111-1111-111111111114', 'Miguel Ángel Carrizo', 'miguel.carrizo@guiaoficios.ar', '+5491145112233', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1527980965255-d3b416303d12?w=240', 'Av. Honorio Pueyrredón 950, Caballito, CABA', -34.6115, -58.4465, 'professional'),
('11111111-1111-1111-1111-111111111115', 'Silvia Beatriz Acosta', 'silvia.acosta@guiaoficios.ar', '+5491145223344', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=240', 'Av. Avellaneda 1850, Flores, CABA', -34.6238, -58.4590, 'professional'),
-- 8. Cerrajeros
('11111111-1111-1111-1111-111111111116', 'Sergio Darío Alarcón', 'sergio.alarcon@guiaoficios.ar', '+5491145334455', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=240', 'Av. Rivadavia 5120, Caballito, CABA', -34.6175, -58.4348, 'professional'),
('11111111-1111-1111-1111-111111111117', 'Pablo Nicolás Quiroga', 'pablo.quiroga@guiaoficios.ar', '+5491145445566', '$2a$10$ykERU2/jYPFXfm2Hi0nQq.AQT.dHyN4HuJK3zIhljOL6Xf.GSKMwS', 'manual', 'https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?w=240', 'Av. Rivadavia 4180, Almagro, CABA', -34.6132, -58.4235, 'professional')
ON CONFLICT (id) DO UPDATE
SET password_hash = CASE
    WHEN users.password_hash IS NULL OR users.password_hash NOT LIKE '$2%'
    THEN EXCLUDED.password_hash
    ELSE users.password_hash
END,
avatar_url = COALESCE(users.avatar_url, EXCLUDED.avatar_url);

-- Perfiles de Profesionales (16 profesionales en CABA cubriendo los 8 rubros)
INSERT INTO professional_profiles (
    id, user_id, primary_category_id, trade_title, specialties, bio, years_experience,
    guarantee_percentage, neighborhood, coverage_radius_km, base_location,
    is_available_now, attends_emergencies, is_dni_verified, is_licensed,
    license_number, license_issuer, verification_status, visit_price_amount, visit_price_label,
    rating_avg, reviews_count, jobs_completed_count, monthly_contacts_count
) VALUES
-- 1. Plomeros
(
    'aaaa1111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    1,
    'Plomero y Gasista Matriculado',
    '["Destapaciones", "Calefones y Termotanques", "Termofusión"]'::jsonb,
    'Hola, soy Roberto. Trabajo hace más de 25 años haciendo reparaciones en casas y departamentos de la zona. Trabajo prolijo, puntual y con garantía en cada arreglo.',
    25, 100, 'Caballito y Flores', 15,
    ST_SetSRID(ST_MakePoint(-58.4330, -34.6150), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula oficial N° 48.912 verificada', 'Metrogas / ENARGAS', 'verified', 12000, '$12.000',
    4.90, 142, 142, 18
),
(
    'bbbb2222-2222-2222-2222-222222222222',
    '22222222-2222-2222-2222-222222222222',
    1,
    'Plomero Especialista en Filtraciones',
    '["Urgencias 24hs", "Filtraciones", "Bombas presurizadoras"]'::jsonb,
    'Especialista en reparaciones rápidas en el día, pérdidas de agua y conexiones de artefactos en toda la zona de Caballito y Parque Chacabuco.',
    14, 100, 'Caballito y Parque Chacabuco', 15,
    ST_SetSRID(ST_MakePoint(-58.4410, -34.6295), 4326)::geography,
    TRUE, TRUE, TRUE, FALSE,
    NULL, NULL, 'verified', 10000, '$10.000',
    4.80, 89, 89, 12
),
-- 2. Electricistas
(
    'aaaa1111-1111-1111-1111-111111111103',
    '11111111-1111-1111-1111-111111111103',
    2,
    'Electricista Matriculado COPIME',
    '["Tableros y térmicas", "Cortocircuitos 24hs", "Certificados APSE"]'::jsonb,
    'Técnico electricista matriculado con 19 años de experiencia en Caballito. Diagnóstico de fugas eléctricas, recableados y colocación de disyuntores.',
    19, 100, 'Caballito y Almagro', 15,
    ST_SetSRID(ST_MakePoint(-58.4355, -34.6085), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula COPIME N° 19.304', 'COPIME / ENRE', 'verified', 13000, '$13.000',
    4.95, 116, 116, 21
),
(
    'aaaa1111-1111-1111-1111-111111111104',
    '11111111-1111-1111-1111-111111111104',
    2,
    'Electricista Domiciliario e Iluminación',
    '["Enchufes y llaves", "Iluminación LED", "Porteros eléctricos"]'::jsonb,
    'Soluciones eléctricas rápidas y seguras para hogares y comercios del barrio. Presupuesto sin cargo por WhatsApp y garantía escrita.',
    11, 100, 'Flores y Caballito', 15,
    ST_SetSRID(ST_MakePoint(-58.4580, -34.6265), 4326)::geography,
    TRUE, FALSE, TRUE, TRUE,
    'Matrícula APSE N° 7.812', 'APSE', 'verified', 11000, '$11.000',
    4.84, 64, 64, 10
),
-- 3. Gasistas
(
    'aaaa1111-1111-1111-1111-111111111105',
    '11111111-1111-1111-1111-111111111105',
    3,
    'Gasista Matriculado de 1° Categoría',
    '["Calefones y estufas", "Detección de fugas", "Rehabilitaciones Metrogas"]'::jsonb,
    'Gasista de primera categoría habilitado por Metrogas. Limpieza de quemadores, cambio de termocuplas, pruebas de hermeticidad y trámites.',
    23, 100, 'Caballito, Flores y Boedo', 15,
    ST_SetSRID(ST_MakePoint(-58.4385, -34.6230), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula Metrogas N° 27.610', 'Metrogas / ENARGAS', 'verified', 14000, '$14.000',
    4.94, 131, 131, 19
),
(
    'aaaa1111-1111-1111-1111-111111111106',
    '11111111-1111-1111-1111-111111111106',
    3,
    'Gasista Matriculado y Termotanques',
    '["Cocinas y hornos", "Tiro balanceado", "Termotanques"]'::jsonb,
    'Instalación y reparación de todo tipo de artefactos a gas con normas vigentes de seguridad. Atención cordial a personas mayores.',
    16, 100, 'Almagro y Caballito', 15,
    ST_SetSRID(ST_MakePoint(-58.4215, -34.6060), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula Metrogas N° 33.418', 'Metrogas / ENARGAS', 'verified', 12500, '$12.500',
    4.87, 78, 78, 13
),
-- 4. Aire acondicionado
(
    'aaaa1111-1111-1111-1111-111111111107',
    '11111111-1111-1111-1111-111111111107',
    4,
    'Técnico en Refrigeración y Split Matriculado',
    '["Carga de gas R410A", "Instalación con bomba de vacío", "Limpieza de filtros"]'::jsonb,
    'Técnico certificado por la Cámara Argentina de Calefacción y Aire Acondicionado. Instalaciones sin romper de más y service anual preventivo.',
    13, 100, 'Caballito y Villa Crespo', 15,
    ST_SetSRID(ST_MakePoint(-58.4450, -34.6145), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Matrícula CACAAV N° 14.920', 'CACAAV', 'verified', 15000, '$15.000',
    4.91, 95, 95, 16
),
(
    'aaaa1111-1111-1111-1111-111111111108',
    '11111111-1111-1111-1111-111111111108',
    4,
    'Especialista en Equipos Inverter y Climatización',
    '["Placas electrónicas", "Pérdidas de agua interior", "Desinstalación y mudanza"]'::jsonb,
    'Reparación de equipos Split e Inverter que no enfrían o gotean hacia adentro. Repuestos originales y garantía por toda la temporada.',
    9, 100, 'Villa Crespo, Almagro y Palermo', 15,
    ST_SetSRID(ST_MakePoint(-58.4310, -34.6015), 4326)::geography,
    TRUE, FALSE, TRUE, TRUE,
    'Matrícula CACAAV N° 18.203', 'CACAAV', 'verified', 14000, '$14.000',
    4.82, 54, 54, 9
),
-- 5. Pintores
(
    'aaaa1111-1111-1111-1111-111111111109',
    '11111111-1111-1111-1111-111111111109',
    5,
    'Pintor de Obra e Interiores Prolijo',
    '["Humedad y enduido", "Látex lavable y cielorrasos", "Impermeabilización de terrazas"]'::jsonb,
    'Más de 20 años pintando departamentos y casas en Caballito. Cubrimos todos los muebles y pisos antes de empezar y entregamos limpio.',
    20, 100, 'Caballito y Flores', 15,
    ST_SetSRID(ST_MakePoint(-58.4342, -34.6198), 4326)::geography,
    TRUE, FALSE, TRUE, FALSE,
    NULL, NULL, 'verified', 8000, '$8.000',
    4.93, 108, 108, 14
),
(
    'aaaa1111-1111-1111-1111-111111111110',
    '11111111-1111-1111-1111-111111111110',
    5,
    'Pintura General, Maderas y Revestimientos',
    '["Barnizado de puertas", "Revestimiento plástico Tarquini", "Reparación de grietas"]'::jsonb,
    'Trabajos de pintura interior y exterior con materiales de primera línea. Presupuesto detallado por escrito sin sorpresas.',
    15, 100, 'Boedo, Caballito y Parque Patricios', 15,
    ST_SetSRID(ST_MakePoint(-58.4162, -34.6248), 4326)::geography,
    TRUE, FALSE, TRUE, FALSE,
    NULL, NULL, 'verified', 9000, '$9.000',
    4.85, 67, 67, 10
),
-- 6. Albañiles
(
    'aaaa1111-1111-1111-1111-111111111112',
    '11111111-1111-1111-1111-111111111112',
    6,
    'Maestro Mayor de Obras y Albañilería',
    '["Colocación de cerámicos", "Revoques y humedad de cimientos", "Refacción de baños y cocinas"]'::jsonb,
    'Maestro Mayor de Obras matriculado. Arreglos chicos de albañilería, cambio de baldosas flojas, revoques finos y remodelaciones completas.',
    24, 100, 'Caballito, Boedo y Almagro', 15,
    ST_SetSRID(ST_MakePoint(-58.4265, -34.6215), 4326)::geography,
    TRUE, FALSE, TRUE, TRUE,
    'Matrícula CPIC N° 11.450', 'CPIC', 'verified', 12000, '$12.000',
    4.89, 92, 92, 13
),
(
    'aaaa1111-1111-1111-1111-111111111113',
    '11111111-1111-1111-1111-111111111113',
    6,
    'Albañil de Confianza y Techista',
    '["Membranas y filtraciones", "Durlock y tabiques", "Contrapisos y carpetas"]'::jsonb,
    'Reparación definitiva de goteras en techos, colocación de membranas asfálticas y cerramientos en seco. Referencias comprobables en el barrio.',
    17, 100, 'Parque Chacabuco y Caballito', 15,
    ST_SetSRID(ST_MakePoint(-58.4395, -34.6335), 4326)::geography,
    TRUE, TRUE, TRUE, FALSE,
    NULL, NULL, 'verified', 10000, '$10.000',
    4.81, 59, 59, 9
),
-- 7. Jardineros
(
    'aaaa1111-1111-1111-1111-111111111114',
    '11111111-1111-1111-1111-111111111114',
    7,
    'Jardinero Paisajista y Poda en Altura',
    '["Corte de césped y bordes", "Poda de árboles y enredaderas", "Control de plagas orgánico"]'::jsonb,
    'Mantenimiento mensual o puntual de patios, jardines, balcones y terrazas verdes. Llevo todas mis herramientas y retiro las bolsas de poda.',
    18, 100, 'Caballito, Villa Crespo y Flores', 15,
    ST_SetSRID(ST_MakePoint(-58.4465, -34.6115), 4326)::geography,
    TRUE, FALSE, TRUE, FALSE,
    NULL, NULL, 'verified', 9000, '$9.000',
    4.92, 83, 83, 14
),
(
    'aaaa1111-1111-1111-1111-111111111115',
    '11111111-1111-1111-1111-111111111115',
    7,
    'Especialista en Huertas, Balcones y Riego',
    '["Diseño de maceteros", "Riego automático por goteo", "Tierra fértil y plantines"]'::jsonb,
    'Técnica en jardinería urbana. Armado y cuidado de balcones floridos, patios internos y sistemas de riego automático para cuando salís de viaje.',
    10, 100, 'Flores, Caballito y Parque Chacabuco', 15,
    ST_SetSRID(ST_MakePoint(-58.4590, -34.6238), 4326)::geography,
    TRUE, FALSE, TRUE, TRUE,
    'Certificación UBA Agronomía #4.102', 'FAUBA', 'verified', 8500, '$8.500',
    4.96, 71, 71, 15
),
-- 8. Cerrajeros
(
    'aaaa1111-1111-1111-1111-111111111116',
    '11111111-1111-1111-1111-111111111116',
    8,
    'Cerrajero de Urgencias 24 Horas',
    '["Apertura sin romper", "Cambio de combinación en el acto", "Cerraduras doble paleta y digitales"]'::jsonb,
    'Vecino de Caballito con local a la calle. Si te quedaste afuera o se trabó la llave, llego en 15 a 20 minutos a cualquier hora del día o la noche.',
    21, 100, 'Caballito (Av. Rivadavia y Acoyte)', 15,
    ST_SetSRID(ST_MakePoint(-58.4348, -34.6175), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Cámara de Cerrajeros CABA N° 1.184', 'CACEBA', 'verified', 14000, '$14.000',
    4.95, 164, 164, 27
),
(
    'aaaa1111-1111-1111-1111-111111111117',
    '11111111-1111-1111-1111-111111111117',
    8,
    'Cerrajería de Seguridad y Puertas Blindadas',
    '["Cerrojos suplementarios", "Destrabe de cortinas y portones", "Copias codificadas"]'::jsonb,
    'Instalación de cerrojos adicionales de alta seguridad y reparación de puertas Pentágono o blindadas en Almagro y Caballito.',
    12, 100, 'Almagro y Caballito Norte', 15,
    ST_SetSRID(ST_MakePoint(-58.4235, -34.6132), 4326)::geography,
    TRUE, TRUE, TRUE, TRUE,
    'Cámara de Cerrajeros CABA N° 1.690', 'CACEBA', 'verified', 13000, '$13.000',
    4.86, 74, 74, 11
)
ON CONFLICT (id) DO NOTHING;

-- Trabajos realizados verificados (Portfolio)
INSERT INTO portfolio_items (id, professional_id, title, subtitle, image_url, is_verified) VALUES
('dddd1111-1111-1111-1111-111111111111', 'aaaa1111-1111-1111-1111-111111111111', 'Canilla monocomando cambiada', 'Instalación y sellado completo', 'https://images.unsplash.com/photo-1585704032915-c3400ca199e7?w=600', TRUE),
('dddd2222-2222-2222-2222-222222222222', 'aaaa1111-1111-1111-1111-111111111111', 'Reparación de cañería termofusión', 'Solución definitiva a filtración', 'https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=600', TRUE),
('dddd1111-1111-1111-1111-111111111103', 'aaaa1111-1111-1111-1111-111111111103', 'Tablero seccional normalizado', 'Térmicas y disyuntor diferencial Schneider', 'https://images.unsplash.com/photo-1621905251189-08b45d6a269e?w=600', TRUE),
('dddd1111-1111-1111-1111-111111111105', 'aaaa1111-1111-1111-1111-111111111105', 'Instalación de calefón tiro balanceado', 'Aprobado con oblea Metrogas', 'https://images.unsplash.com/photo-1504328345606-18bbc8c9d7d1?w=600', TRUE),
('dddd1111-1111-1111-1111-111111111107', 'aaaa1111-1111-1111-1111-111111111107', 'Instalación Split Inverter 3500 frigorías', 'Presurización con nitrógeno y vacío', 'https://images.unsplash.com/photo-1631545806609-35d4ae490431?w=600', TRUE),
('dddd1111-1111-1111-1111-111111111109', 'aaaa1111-1111-1111-1111-111111111109', 'Pintura integral de living comedor', 'Enduido completo y látex satinado lavable', 'https://images.unsplash.com/photo-1562259949-e8e7689d7828?w=600', TRUE),
('dddd1111-1111-1111-1111-111111111112', 'aaaa1111-1111-1111-1111-111111111112', 'Renovación de cerámicos en cocina', 'Porcelanato rectificado y pastina antihongos', 'https://images.unsplash.com/photo-1503387762-592deb58ef4e?w=600', TRUE),
('dddd1111-1111-1111-1111-111111111114', 'aaaa1111-1111-1111-1111-111111111114', 'Poda y puesta en valor de jardín', 'Nivelación de césped y cantero floral', 'https://images.unsplash.com/photo-1558904541-efa843a96f01?w=600', TRUE),
('dddd1111-1111-1111-1111-111111111116', 'aaaa1111-1111-1111-1111-111111111116', 'Colocación de cerradura multipunto', 'Blindaje y doble cerrojo de seguridad', 'https://images.unsplash.com/photo-1558002038-1055907df827?w=600', TRUE)
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
    'Av. Rivadavia 5420, Caballito',
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

-- Reseñas de vecinos verificados
INSERT INTO reviews (id, professional_id, author_name, author_zone, rating, highlights, comment) VALUES
(
    'ffff1111-1111-1111-1111-111111111111',
    'aaaa1111-1111-1111-1111-111111111111',
    'Graciela M.',
    'Vecina de Caballito (Av. Pedro Goyena)',
    5,
    '["Puntualidad", "Precio justo", "Honestidad y prolijidad", "Excelente atención", "Solución en el acto"]'::jsonb,
    'Vino puntual a la hora acordada, resolvió la pérdida del baño en el acto y dejó todo impecable. Muy respetuoso.'
),
(
    'ffff1111-1111-1111-1111-111111111103',
    'aaaa1111-1111-1111-1111-111111111103',
    'Héctor S.',
    'Vecino de Caballito (Primera Junta)',
    5,
    '["Puntualidad", "Precio justo", "Solución en el acto"]'::jsonb,
    'Se nos cortó la luz un domingo a la tarde, Jorge encontró el corto en el tablero en 15 minutos y nos cobró lo acordado.'
),
(
    'ffff1111-1111-1111-1111-111111111105',
    'aaaa1111-1111-1111-1111-111111111105',
    'Norma B.',
    'Vecina de Caballito Sur',
    5,
    '["Honestidad y prolijidad", "Excelente atención"]'::jsonb,
    'Limpió el calefón que se apagaba solo, revisó que no hubiera pérdidas en la cocina y me explicó todo con mucha paciencia.'
),
(
    'ffff1111-1111-1111-1111-111111111116',
    'aaaa1111-1111-1111-1111-111111111116',
    'Marta R.',
    'Vecina de Av. Rivadavia 5200',
    5,
    '["Puntualidad", "Solución en el acto", "Precio justo"]'::jsonb,
    'Me quedé afuera de noche, llegó en 15 minutos reloj y abrió sin dañar la cerradura. Totalmente recomendable.'
)
ON CONFLICT (id) DO NOTHING;
