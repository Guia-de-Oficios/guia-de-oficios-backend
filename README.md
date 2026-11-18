# Guía de Oficios — Backend API REST (Golang + PostgreSQL/PostGIS en Docker) ⚙️

Repositorio oficial del servidor Backend para **Guía de Oficios**. Es una **API REST JSON 100% pura** (sin vistas HTML ni repositorios falsos en memoria), conectada directamente a **PostgreSQL 16 + PostGIS** para dar soporte a las **15 pantallas oficiales generadas en Stitch**.

## 🔗 Repositorios Hermanos de la Organización
* 📱 **App Android (Kotlin + Jetpack Compose):** [`Guia-de-Oficios/guia-de-oficios-android`](https://github.com/Guia-de-Oficios/guia-de-oficios-android)
* 🎨 **Diseño Stitch (15 Vistas Oficiales) y Documentación Maestra:** [`Guia-de-Oficios/guia-de-oficios-docs-ui`](https://github.com/Guia-de-Oficios/guia-de-oficios-docs-ui)

---

## 📚 Documentación y Plan de Pasos en `/docs`
1. 🚀 **[`docs/06_radiografia_vistas_y_plan_paso_a_paso.md`](docs/06_radiografia_vistas_y_plan_paso_a_paso.md)**: **[HOJA DE RUTA DEL EQUIPO]** Radiografía completa de las 15 vistas, auditoría del esqueleto base, guía de ejecución local y checklist de los **14 Pasos Secuenciales** (ramas `feat/paso-XX-...` y `fix/paso-XX-...`).
2. [`docs/01_especificacion_mvp_y_evolutivos.md`](docs/01_especificacion_mvp_y_evolutivos.md): Especificación funcional del MVP y contratos REST.
3. [`docs/02_guia_ux_ui_accesible.md`](docs/02_guia_ux_ui_accesible.md): Lineamientos de accesibilidad y sistema de diseño.
4. [`docs/03_trabajo_colaborativo_y_onboarding.md`](docs/03_trabajo_colaborativo_y_onboarding.md): Guía de trabajo en equipo y conexión con dispositivos Android.
5. [`docs/04_vision_negocio_y_estrategia.md`](docs/04_vision_negocio_y_estrategia.md): Estrategia de producto y fases de monetización.

---

## 🗺️ Estado de los 14 Pasos en Backend

| Paso | Rama Backend | Alcance en Servidor Go + PostgreSQL/PostGIS | Estado |
| :---: | :--- | :--- | :---: |
| **Paso 1** | `feat/paso-01-auth-bcrypt` | Hash `bcrypt` real en `internal/auth/jwt.go` y actualización de `migrations/002_seed.sql` | ✅ **Completado** |
| **Paso 2** | `feat/paso-02-registro-barrio` | Registro con barrio/dirección, geocodificación nacional, autocompletado en vivo (`/geocode/suggest`) y validación criptográfica de `id_token` en Google Auth | ✅ **Completado** |
| **Paso 3** | `fix/paso-03-middleware-jwt` | Middleware estricto `withRequiredAuth` (`401`), cierre de fuga de datos en `GET /api/v1/requests/my` (`BUG-03`), validación de dueño en endpoints de profesional (`BUG-07`) y flag `is_dni_verified` en `User` | ✅ **Completado** |
| **Paso 4** | `feat/paso-04-perfil-usuario` | Edición integral de perfil (`PUT /api/v1/users/me`) e historial real de solicitudes del vecino (`GET /api/v1/requests/my`) | ✅ **Completado** |
| **Paso 5** | `feat/paso-05-seed-y-coordenadas` | Semilla de al menos 16 profesionales verificados cubriendo los 8 rubros oficiales y geocodificador por barrio | 🔜 **Próximo paso** |
| **Pasos 6–14** | `feat/paso-06...14` | Filtros PostGIS, distancia dinámica en ficha, ranking barrial, solicitudes reales, transición de estados `en_route`/`completed` y calificación transaccional | ⏳ Pendiente |

---

## 🐳 Cómo Levantar todo el Backend con Docker (Sin consumir RAM de GoLand)

El archivo `docker-compose.yml` levanta tanto la base de datos **PostgreSQL 16 + PostGIS** (`guia_oficios_db` en el puerto `5432`) como el servidor **API Go** (`guia_oficios_api` en el puerto `8080`), ejecutando automáticamente las migraciones `migrations/001_init.sql`, `migrations/002_seed.sql` y `migrations/003_seed_regions.sql`:

```powershell
docker compose up -d --build
```

Para ejecutar los tests unitarios de la API localmente:
```powershell
go test -v ./...
```

---

## 📡 Endpoints REST Disponibles (`http://localhost:8080`)

| Método | Ruta | Autenticación | Paso del Plan / Funcionalidad |
| :--- | :--- | :---: | :--- |
| `GET` | `/api/v1/health` | Pública | **Infra:** Health check contra PostgreSQL + versión de PostGIS |
| `GET` | `/api/v1/metrics` | Pública | **Infra:** Snapshot de métricas clave y señales de oro (Datadog) |
| `POST` | `/api/v1/auth/login` | Pública | **Paso 1:** Login por correo o celular + contraseña con `bcrypt` |
| `POST` | `/api/v1/auth/register` | Pública | **Paso 2:** Registro manual con barrio/dirección, coordenadas reales y alta de perfil profesional |
| `POST` | `/api/v1/auth/google` | Pública | **Paso 2:** Login/Registro con token Google verificado criptográficamente |
| `GET` | `/api/v1/geocode` | Pública | **Paso 2:** Validación y geocodificación de dirección y localidad en Argentina |
| `GET` | `/api/v1/geocode/suggest` | Pública | **Paso 2:** Autocompletado en vivo de localidades y calles normalizadas (Georef AR + OSM) |
| `GET` | `/api/v1/auth/me` | `Bearer JWT` (estricto) | **Paso 3:** Perfil del usuario autenticado (`401` si falta o expira) |
| `PATCH` | `/api/v1/users/me/address` | `Bearer JWT` (estricto) | **Paso 3 / 5:** Actualización de dirección y geocodificación de barrio |
| `PUT` | `/api/v1/users/me` | `Bearer JWT` (estricto) | **Paso 4:** Edición completa de perfil de vecino o profesional |
| `GET` | `/api/v1/categories` | Pública | **Paso 5:** Grilla de 8 rubros oficiales con badges (`Matriculados`, `Urgencias`) |
| `POST` | `/api/v1/diagnostic` | Pública | **Paso 5:** Asistente diagnóstico interactivo |
| `GET` | `/api/v1/professionals` | Pública | **Paso 6 / 8:** Búsqueda geoespacial PostGIS (filtros combinados y podio mensual) |
| `GET` | `/api/v1/professionals/{id}` | Pública | **Paso 7:** Ficha completa, cálculo ST_Distance real, portfolio y reseñas |
| `POST` | `/api/v1/uploads` | `Bearer JWT` | **Paso 9 / 13:** Subida real `multipart/form-data` de fotos de problemas, DNI y matrícula |
| `POST` | `/api/v1/requests` | `Bearer JWT` (estricto) | **Paso 3 / 9:** Creación de solicitud guiada con datos reales del JWT |
| `GET` | `/api/v1/requests/my` | `Bearer JWT` (estricto) | **Paso 3 / 4:** Historial de solicitudes del vecino sin fuga de datos (`BUG-03`) |
| `GET` | `/api/v1/requests/{id}` | `Bearer JWT` | **Paso 11:** Detalle de solicitud y seguimiento de viaje |
| `POST` | `/api/v1/requests/{id}/reviews` | `Bearer JWT` (estricto) | **Paso 3 / 12:** Calificación transaccional 1-5 estrellas + 5 aspectos destacables |
| `POST` | `/api/v1/professionals/{id}/verification` | `Bearer JWT` + Dueño | **Paso 3 / 13:** Carga y validación de oficio, cobertura, DNI y matrícula (`BUG-07`) |
| `GET` | `/api/v1/professionals/{id}/dashboard` | Pública / `Bearer JWT` | **Paso 14:** Panel del profesional, métricas y solicitudes recibidas |
| `PATCH` | `/api/v1/professionals/{id}/availability` | `Bearer JWT` + Dueño | **Paso 3 / 14:** Cambio de disponibilidad en tiempo real validando propietario (`BUG-07`) |
