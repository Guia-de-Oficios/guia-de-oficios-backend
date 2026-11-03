# Backlog Técnico Completo del MVP (100% Real, Sin Mocks + APK Instalable) — "Guía de Oficios"

> [!IMPORTANT]
> **Objetivo del Backlog:** Llegar al **MVP 100% funcional** en producción, eliminando cualquier repositorio en memoria (`InMemoryStore`) o dato harcodeado en la app, operando exclusivamente contra **PostgreSQL 16 + PostGIS**, con subida real de fotos desde el teléfono, autenticación real (JWT + Google), e generando el archivo **`guia-de-oficios-mvp-demo.apk`** listo para instalar en dispositivos Android físicos.

> [!TIP]
> **Importación directa en Jira:**  
> En esta misma carpeta se encuentra el archivo **[`05_backlog_jira_mvp.csv`](05_backlog_jira_mvp.csv)** formateado con las columnas oficiales de Jira (`Issue Type`, `Issue Key`, `Issue ID`, `Parent ID`, `Summary`, `Description`, `Priority`, `Story Points`, `Labels`, `Component`).  
> Para importarlo en Jira: **Configuración de Proyecto / Sistema -> Importar desde CSV (External System Import -> CSV)** y seleccionar `05_backlog_jira_mvp.csv`.

---

## 📊 Resumen de Épicas y Estimación (Total: 94 Story Points)

| Key | Épica | Historias Vinculadas | Story Points | Foco Principal |
| :--- | :--- | :--- | :---: | :--- |
| **`GDO-1`** | **Infraestructura Base de Datos Real y Almacenamiento (PostgreSQL + PostGIS + Storage)** | `GDO-8`, `GDO-9`, `GDO-10` | **11 SP** | Eliminar `InMemoryStore`, driver `pgx/v5`, migraciones SQL y subida real de fotos (`Multipart`). |
| **`GDO-2`** | **Autenticación Real y Seguridad (JWT + Bcrypt + Google OAuth2)** | `GDO-11`, `GDO-12`, `GDO-13` | **13 SP** | Registro manual con hash `bcrypt`, login, Google Sign-In real y middleware JWT. |
| **`GDO-3`** | **API Backend Real en Golang (100% SQL/PostGIS sin Mocks)** | `GDO-14`, `GDO-15`, `GDO-16`, `GDO-17`, `GDO-18` | **23 SP** | Consultas espaciales `ST_Distance`/`ST_DWithin`, solicitudes reales, reseñas transaccionales y panel profesional. |
| **`GDO-4`** | **Arquitectura Base Android y Networking Real (Gradle + Hilt + Retrofit + DataStore)** | `GDO-19`, `GDO-20` | **8 SP** | Proyecto Gradle compilable, Hilt, OkHttp + `AuthInterceptor` JWT y `DataStore` de sesión. |
| **`GDO-5`** | **Flujos Android de Onboarding y Cliente Conectados al Backend (Pantallas Stitch 1 a 8)** | `GDO-21`, `GDO-22`, `GDO-23`, `GDO-24`, `GDO-25` | **26 SP** | Pantallas Stitch 1 a 8 en Compose conectadas a la API, cámara nativa y apertura real de WhatsApp y Teléfono. |
| **`GDO-6`** | **Flujos Android del Profesional Conectados al Backend (Pantallas Stitch 9 y 10)** | `GDO-26` | **8 SP** | Subida real de DNI Frente/Dorso y Matrícula + Panel de Solicitudes y switch de disponibilidad en vivo. |
| **`GDO-7`** | **Despliegue Cloud del Backend y Generación de APK Demo Instalable** | `GDO-27` | **5 SP** | Deploy con HTTPS del backend + PostGIS y generación del `.apk` firmado para instalar en celulares Android. |

---

## 📋 Detalle de Todas las Historias de Usuario y Criterios de Aceptación (DoD)

### 🗄️ Épica 1 (`GDO-1`): Infraestructura Base de Datos Real y Almacenamiento

#### `GDO-8` — [BE-DB] Eliminar `InMemoryStore` y conectar driver `pgxpool` nativo a PostgreSQL 16 + PostGIS *(3 SP · Highest)*
* **Descripción:** Como desarrollador backend, quiero eliminar por completo el almacenamiento en memoria (`InMemoryStore`) e integrar `jackc/pgx/v5/pgxpool` para que el servidor opere exclusivamente contra PostgreSQL 16 + PostGIS real.
* **Criterios de Aceptación (DoD):**
  1. Se agrega `github.com/jackc/pgx/v5` a `go.mod`.
  2. Se elimina `InMemoryStore` de `internal/repository/store.go`; si PostgreSQL no está disponible al iniciar, el servidor falla explícitamente (`log.Fatal`) en lugar de usar datos falsos en memoria.
  3. El endpoint `GET /api/v1/health` ejecuta `db.Ping(ctx)` real y devuelve el estado de PostGIS (`SELECT PostGIS_Version()`).

#### `GDO-9` — [BE-DB] Completar migraciones SQL definitivas e índices espaciales GIST *(3 SP · Highest)*
* **Descripción:** Consolidar el esquema DDL en PostgreSQL para soportar usuarios con hash de contraseña, coordenadas `GEOGRAPHY(Point, 4326)`, URLs reales de documentos y estados de solicitudes.
* **Criterios de Aceptación (DoD):**
  1. Tabla `users` incluye `password_hash`, `google_sub`, `auth_provider`, `saved_address`, `saved_lat`, `saved_lng` y `role`.
  2. Tabla `professional_profiles` tiene índice espacial `GIST(base_location)` y columnas para `dni_front_url`, `dni_back_url`, `license_doc_url`, `is_available_now`, `visit_price_amount` y métricas calculadas.
  3. Migración automática al arrancar el servidor o comando CLI.

#### `GDO-10` — [BE-STORAGE] Implementar endpoint Multipart para subida real de imágenes (DNI / Problemas / Portfolio) *(5 SP · High)*
* **Descripción:** Subir fotos reales tomadas con la cámara o elegidas de la galería del celular para que se almacenen de forma persistente y devuelvan una URL pública accesible por la app.
* **Criterios de Aceptación (DoD):**
  1. Endpoint `POST /api/v1/uploads` que recibe `multipart/form-data` con el campo `file` y `category` (`dni_front`, `dni_back`, `license`, `request_photo`, `portfolio`).
  2. Validación de tipo MIME (`image/jpeg`, `image/png`, `image/webp` o `application/pdf`) y tamaño máximo (10 MB).
  3. Almacenamiento persistente en volumen/disco servido por HTTP o bucket S3/Cloudflare R2, retornando `{"url": "https://.../uploads/..."}`.

---

### 🔐 Épica 2 (`GDO-2`): Autenticación Real y Seguridad

#### `GDO-11` — [BE-AUTH] Implementar Registro Manual y Login con Bcrypt y emisión de JWT real *(5 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. `POST /api/v1/auth/register` hashea la contraseña con `golang.org/x/crypto/bcrypt`, inserta el registro en `users` (y crea la fila inicial en `professional_profiles` si `role == 'professional'`).
  2. `POST /api/v1/auth/login` verifica email y contraseña contra PostgreSQL y devuelve `{"token": "...", "user": {...}}`.
  3. Validación de email único (HTTP 409 Conflict si ya existe).

#### `GDO-12` — [BE-AUTH] Implementar autenticación real con Google Sign-In (`POST /api/v1/auth/google`) *(5 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Endpoint `POST /api/v1/auth/google` que recibe `{"id_token": "...", "role": "client|professional"}`.
  2. Valida el `id_token` de Google, crea o recupera el usuario en PostgreSQL y retorna el JWT propio de Guía de Oficios.

#### `GDO-13` — [BE-AUTH] Implementar Middleware JWT y endpoints de perfil (`/api/v1/auth/me` y cambio de dirección) *(3 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Middleware HTTP que valida `Authorization: Bearer <token>` e inyecta `user_id` y `role` en el contexto.
  2. `GET /api/v1/auth/me` devuelve los datos reales desde PostgreSQL.
  3. `PATCH /api/v1/users/me/address` actualiza `saved_address`, `saved_lat` y `saved_lng`.

---

### ⚙️ Épica 3 (`GDO-3`): API Backend Real en Golang (100% SQL/PostGIS sin Mocks)

#### `GDO-14` — [BE-API] Implementar búsqueda geoespacial real en PostGIS (`GET /api/v1/professionals`) *(5 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. Consulta 100% SQL con PostGIS usando `ST_Distance` para calcular los metros reales y `ST_DWithin` para filtrar dentro del radio de cobertura.
  2. Soporta filtros combinables: `category_id`, `q` (`ILIKE`), `only_urgent=true`, `only_licensed=true`, `top_rated=true`.
  3. Solo devuelve profesionales con `is_available_now = TRUE`.
  4. Formatea `distance_label` dinámicamente según los metros reales (ej. `"A 6 cuadras (600 metros)"`).

#### `GDO-15` — [BE-API] Implementar consulta SQL completa de Ficha de Profesional (`GET /api/v1/professionals/{id}`) *(3 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Consulta real a `professional_profiles` + `users` + `categories` por UUID.
  2. Carga desde `portfolio_items` todas las fotos reales de trabajos verificados y desde `reviews` las opiniones con sus `highlights`.

#### `GDO-16` — [BE-API] Implementar persistencia real de Solicitud Guiada y Confirmación (`POST` y `GET /api/v1/requests`) *(5 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. `POST /api/v1/requests` inserta el pedido en PostgreSQL con la URL de la foto subida y suma `+1` a `monthly_contacts_count`.
  2. `GET /api/v1/requests/{id}` devuelve el pedido para la Pantalla 7 de confirmación.
  3. `GET /api/v1/requests/my` devuelve el historial de contactos del vecino.

#### `GDO-17` — [BE-API] Implementar transacción SQL de Calificación y recálculo de reputación (`POST /api/v1/requests/{id}/reviews`) *(5 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Transacción SQL que cambia `status` a `'completed'`, inserta la reseña con los 5 `highlights` (`Puntualidad`, `Precio justo`, `Honestidad y prolijidad`, `Excelente atención`, `Solución en el acto`) y recalcula con `AVG(rating)` y `COUNT(*)` los valores de `rating_avg`, `reviews_count` y `jobs_completed_count` en `professional_profiles`.

#### `GDO-18` — [BE-API] Implementar Onboarding Profesional DNI/Matrícula y Panel de Solicitudes en SQL (Pantallas 9 y 10) *(5 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. `POST /api/v1/professionals/{id}/verification` persiste oficio, especialidades, zona, URLs de DNI frente/dorso y matrícula en PostgreSQL.
  2. `PATCH /api/v1/professionals/{id}/availability` actualiza `is_available_now` en tiempo real.
  3. `GET /api/v1/professionals/{id}/dashboard` calcula desde SQL las métricas del Bento Grid y lista las solicitudes recibidas.

---

### 📱 Épica 4 (`GDO-4`): Arquitectura Base Android y Networking Real

#### `GDO-19` — [AND-ARCH] Configurar proyecto Gradle completo en Android Studio (Kotlin + Compose + Hilt + Retrofit + Coil) *(5 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. Archivos `settings.gradle.kts`, `build.gradle.kts` (root y `:app`), `gradle.properties`, `gradlew` y `AndroidManifest.xml` listos para ejecutar `./gradlew assembleDebug`.
  2. Dependencias: Compose Material 3, Navigation Compose, Hilt, Retrofit2 + OkHttp3, Coil (`AsyncImage`) y DataStore.
  3. Permisos de Internet, ubicación, cámara/galería y queries de paquetes para WhatsApp y marcador telefónico.

#### `GDO-20` — [AND-ARCH] Implementar gestión de sesión real (`DataStore` + `AuthInterceptor` JWT) y repositorio sin mocks *(3 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. `SessionManager` guarda JWT, rol, usuario y dirección en Jetpack DataStore.
  2. `AuthInterceptor` inyecta el header `Authorization: Bearer <token>` en cada llamada de Retrofit.

---

### 🎨 Épica 5 (`GDO-5`): Flujos Android de Onboarding y Cliente (Pantallas Stitch 1 a 8)

#### `GDO-21` — [AND-UI] Conectar Pantallas Stitch 1 y 2 (Bienvenida, Selección de Rol, Registro Manual y Google Sign-In) *(5 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Pantallas 1 y 2 idénticas a Stitch conectadas a `POST /api/v1/auth/register`, `POST /api/v1/auth/login` y `POST /api/v1/auth/google`.

#### `GDO-22` — [AND-UI] Conectar Pantallas Stitch 3 y 4 (Inicio Cliente, Diagnóstico, Cambio de Dirección y Resultados con Filtros) *(5 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Grilla de 8 oficios alimentada por `GET /api/v1/categories`, asistente por `POST /api/v1/diagnostic`, y resultados filtrados en vivo por `GET /api/v1/professionals`.

#### `GDO-23` — [AND-UI] Conectar Pantalla Stitch 5 (Perfil del Profesional, Galería con Coil y Barra Fija de Contacto) *(3 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Ficha completa alimentada por `GET /api/v1/professionals/{id}` con imágenes reales cargadas por Coil y barra fija inferior de contacto.

#### `GDO-24` — [AND-UI] Conectar Pantallas Stitch 6 y 7 (Solicitud Guiada con Cámara/Galería Real e Intents a WhatsApp y Teléfono) *(8 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. Captura de foto real desde cámara/galería -> subida por `POST /api/v1/uploads` -> creación del pedido con `POST /api/v1/requests`.
  2. En la Pantalla 7, apertura real de **WhatsApp** (`https://wa.me/...`) con el mensaje armado y del **Marcador Telefónico** (`tel:...`).

#### `GDO-25` — [AND-UI] Conectar Pantalla Stitch 8 y Pestaña 'Mis contactos' (Listado de Solicitudes y Calificación Real) *(5 SP · High)*
* **Criterios de Aceptación (DoD):**
  1. Historial de profesionales contactados en la pestaña 'Mis contactos' y pantalla de calificación de 5 estrellas + 5 aspectos destacables conectada a `POST /api/v1/requests/{id}/reviews`.

---

### 👷 Épica 6 (`GDO-6`): Flujos Android del Profesional (Pantallas Stitch 9 y 10)

#### `GDO-26` — [AND-UI] Conectar Pantallas Stitch 9 y 10 (Onboarding Profesional con Fotos de DNI Reales y Panel Profesional en Vivo) *(8 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. Subida real de fotos de frente y dorso del DNI y credencial de matrícula en Pantalla 9.
  2. Panel Profesional (Pantalla 10) conectado a `GET /api/v1/professionals/{id}/dashboard`, switch de disponibilidad en vivo (`PATCH /api/v1/professionals/{id}/availability`) y botones reales para responder por WhatsApp o llamar al vecino.

---

### 🚀 Épica 7 (`GDO-7`): Despliegue Cloud del Backend y Generación de APK Demo Instalable

#### `GDO-27` — [DEVOPS-APK] Desplegar Backend + PostgreSQL/PostGIS en Cloud y generar APK Demo instalable para Android *(5 SP · Highest)*
* **Criterios de Aceptación (DoD):**
  1. Despliegue del contenedor Go + PostgreSQL/PostGIS con URL HTTPS pública para que el teléfono físico se conecte desde cualquier red Wi-Fi o 4G/5G.
  2. Compilación y firma de **`guia-de-oficios-mvp-demo.apk`**, validado en un dispositivo Android físico real recorriendo las 10 pantallas de Stitch de principio a fin.
