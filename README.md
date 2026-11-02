# Guía de Oficios — Backend API (Golang + PostgreSQL/PostGIS) ⚙️

Repositorio oficial del servidor Backend para **Guía de Oficios**, diseñado para dar soporte completo a las **10 pantallas generadas en Stitch** y conectarse con la aplicación móvil Android.

## 🔗 Repositorios Hermanos de la Organización
* 📱 **App Android (Kotlin + Jetpack Compose):** [`Guia-de-Oficios/guia-de-oficios-android`](https://github.com/Guia-de-Oficios/guia-de-oficios-android)
* 🎨 **Diseño Stitch (10 Pantallas Intactas) y Documentación General:** [`Guia-de-Oficios/guia-de-oficios-docs-ui`](https://github.com/Guia-de-Oficios/guia-de-oficios-docs-ui)

---

## 📚 Documentación Incluida en `/docs`
Para que todo el equipo tenga el contexto completo del producto sin salir del repositorio:
1. [`docs/01_especificacion_mvp_y_evolutivos.md`](docs/01_especificacion_mvp_y_evolutivos.md): Mapa de las 10 pantallas de Stitch, funcionalidades y endpoints REST.
2. [`docs/02_guia_ux_ui_accesible.md`](docs/02_guia_ux_ui_accesible.md): Lineamientos de accesibilidad para adultos mayores y tokens de diseño.
3. [`docs/03_trabajo_colaborativo_y_onboarding.md`](docs/03_trabajo_colaborativo_y_onboarding.md): Guía de trabajo en equipo, división de tareas Backend/Android y flujo de ramas Git.
4. [`docs/04_vision_negocio_y_estrategia.md`](docs/04_vision_negocio_y_estrategia.md): Análisis estratégico del MVP, mitigación de fuga de plataforma y fases de monetización (Profesional PRO, Pagos y Membresía Hogar).

---

## 🛠️ Requisitos Previos
* **Go 1.22+** (Recomendado abrir el proyecto en **JetBrains GoLand**)
* **Docker & Docker Compose** (Opcional para desarrollo rápido gracias al modo *In-Memory* incluido, requerido para pruebas sobre **PostgreSQL 16 + PostGIS**).

---

## 🚀 Cómo Levantar el Proyecto en 1 Minuto

### Opción A: Ejecución Instantánea en GoLand (Modo In-Memory con Datos Stitch)
No necesitás tener Docker corriendo para probar la API por primera vez. El servidor arranca con un repositorio en memoria precargado con todos los datos exactos de las pantallas de Stitch (*María Elena Castro*, *Roberto Carlos Gómez*, *Martín Edgardo Silva*, *Graciela M.* y *Jorge Benítez*):

```powershell
go test ./... -v
go run ./cmd/api/main.go
```

### Opción B: Levantar PostgreSQL 16 + PostGIS con Docker
En la raíz de este repositorio se incluye `docker-compose.yml`, que levanta la base de datos en el puerto `5432` y ejecuta automáticamente las migraciones de `migrations/001_init.sql` y `migrations/002_seed.sql`:

```powershell
docker compose up -d
```

---

## 📡 Endpoints Disponibles (`http://localhost:8080`)

| Método | Ruta | Pantalla Stitch que alimenta |
| :--- | :--- | :--- |
| `GET` | `/api/v1/health` | Verificación de estado del servidor y modo de almacenamiento |
| `POST` | `/api/v1/auth/register` | **Pantalla 2:** Registro de usuario (Google o Manual) |
| `GET` | `/api/v1/auth/me` | **Pantalla 1 & 3:** Perfil actual del vecino o profesional |
| `GET` | `/api/v1/categories` | **Pantalla 3:** Grilla de 8 oficios con badges (`Matriculados`, `Urgencias`) |
| `POST` | `/api/v1/diagnostic` | **Pantalla 3:** Asistente *"¿No sabés qué profesional necesitás?"* |
| `GET` | `/api/v1/professionals` | **Pantalla 4:** Listado con filtros (`category_id`, `only_urgent`, `only_licensed`, `top_rated`) |
| `GET` | `/api/v1/professionals/{id}` | **Pantalla 5:** Ficha completa de Roberto Gómez, fotos de trabajos, precios y reseñas |
| `POST` | `/api/v1/requests` | **Pantalla 6:** Solicitud guiada (Urgencia, problema, foto opcional) |
| `GET` | `/api/v1/requests/{id}` | **Pantalla 7:** Confirmación de solicitud y datos de contacto directo |
| `POST` | `/api/v1/requests/{id}/reviews` | **Pantalla 8:** Calificación 1-5 estrellas + 5 aspectos destacables + comentario |
| `POST` | `/api/v1/professionals/{id}/verification` | **Pantalla 9:** Registro profesional, especialidades, frente/dorso DNI y matrícula |
| `GET` | `/api/v1/professionals/{id}/dashboard` | **Pantalla 10:** Panel del profesional, métricas mensuales y solicitudes de vecinos |
| `PATCH` | `/api/v1/professionals/{id}/availability` | **Pantalla 10:** Botón en vivo *"Estás DISPONIBLE"* / *"Pausar disponibilidad"* |

---

## 📂 Estructura de Paquetes
```text
.
├── cmd/api/main.go                 # Punto de entrada del servidor HTTP
├── internal/
│   ├── handlers/                   # Controladores HTTP REST y tests unitarios (api_test.go)
│   ├── models/                     # Estructuras de datos y DTOs JSON
│   └── repository/                 # Implementación PostgresStore + InMemoryStore
├── migrations/                     # Esquema SQL + PostGIS (001_init.sql) y Seed (002_seed.sql)
├── docs/                           # Documentación completa del proyecto
└── docker-compose.yml              # Contenedor PostgreSQL 16 + PostGIS
```
