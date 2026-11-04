# Guía de Oficios — Backend API REST (Golang + PostgreSQL/PostGIS en Docker) ⚙️

Repositorio oficial del servidor Backend para **Guía de Oficios**. Es una **API REST JSON 100% pura** (sin vistas HTML ni repositorios falsos en memoria), conectada directamente a **PostgreSQL 16 + PostGIS** para dar soporte a las **15 pantallas oficiales generadas en Stitch**.

## 🔗 Repositorios Hermanos de la Organización
* 📱 **App Android (Kotlin + Jetpack Compose):** [`Guia-de-Oficios/guia-de-oficios-android`](https://github.com/Guia-de-Oficios/guia-de-oficios-android)
* 🎨 **Diseño Stitch (15 Vistas Oficiales) y Documentación Maestra:** [`Guia-de-Oficios/guia-de-oficios-docs-ui`](https://github.com/Guia-de-Oficios/guia-de-oficios-docs-ui)

---

## 📚 Documentación y Plan de Pasos en `/docs`
1. 🚀 **[`docs/06_radiografia_vistas_y_plan_paso_a_paso.md`](docs/06_radiografia_vistas_y_plan_paso_a_paso.md)**: **[HOJA DE RUTA DEL EQUIPO]** Radiografía completa de las 15 vistas, auditoría del esqueleto base, guía de ejecución local y checklist de los **14 Pasos Secuenciales** (ramas `feat/paso-XX-...`).
2. [`docs/01_especificacion_mvp_y_evolutivos.md`](docs/01_especificacion_mvp_y_evolutivos.md): Especificación funcional del MVP y contratos REST.
3. [`docs/02_guia_ux_ui_accesible.md`](docs/02_guia_ux_ui_accesible.md): Lineamientos de accesibilidad y sistema de diseño.
4. [`docs/03_trabajo_colaborativo_y_onboarding.md`](docs/03_trabajo_colaborativo_y_onboarding.md): Guía de trabajo en equipo y conexión con dispositivos Android.
5. [`docs/04_vision_negocio_y_estrategia.md`](docs/04_vision_negocio_y_estrategia.md): Estrategia de producto y fases de monetización.

---

## 🐳 Cómo Levantar todo el Backend con Docker (Sin consumir RAM de GoLand)

El archivo `docker-compose.yml` levanta tanto la base de datos **PostgreSQL 16 + PostGIS** (`guia_oficios_db` en el puerto `5432`) como el servidor **API Go** (`guia_oficios_api` en el puerto `8080`), ejecutando automáticamente las migraciones `migrations/001_init.sql` y `migrations/002_seed.sql`:

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
| `POST` | `/api/v1/auth/login` | Pública | **Paso 1:** Login por correo o celular + contraseña (migrando a bcrypt) |
| `POST` | `/api/v1/auth/register` | Pública | **Paso 2:** Registro manual (con barrio/dirección y coordenadas reales) |
| `POST` | `/api/v1/auth/google` | Pública | **Paso 2:** Login/Registro con token Google verificado |
| `GET` | `/api/v1/auth/me` | `Bearer JWT` (estricto) | **Paso 3:** Perfil del usuario autenticado (`401` si falta o expira) |
| `PATCH` | `/api/v1/users/me/address` | `Bearer JWT` (estricto) | **Paso 3 / 5:** Actualización de dirección y geocodificación de barrio |
| `GET` | `/api/v1/categories` | Pública | **Paso 5:** Grilla de 8 rubros oficiales con badges (`Matriculados`, `Urgencias`) |
| `POST` | `/api/v1/diagnostic` | Pública | **Paso 5:** Asistente diagnóstico interactivo |
| `GET` | `/api/v1/professionals` | Pública | **Paso 6 / 8:** Búsqueda geoespacial PostGIS (filtros combinados y podio mensual) |
| `GET` | `/api/v1/professionals/{id}` | Pública | **Paso 7:** Ficha completa, cálculo ST_Distance real, portfolio y reseñas |
| `POST` | `/api/v1/uploads` | `Bearer JWT` | **Paso 9 / 13:** Subida real `multipart/form-data` de fotos de problemas, DNI y matrícula |
| `POST` | `/api/v1/requests` | `Bearer JWT` (estricto) | **Paso 9:** Creación de solicitud guiada con datos reales del JWT |
| `GET` | `/api/v1/requests/my` | `Bearer JWT` (estricto) | **Paso 11:** Historial de solicitudes del vecino sin fuga de datos |
| `PATCH` | `/api/v1/requests/{id}/status` | `Bearer JWT` (estricto) | **Paso 11 / 14:** Transición de estados (`pending` -> `en_route` con ETA -> `completed`) |
| `GET` | `/api/v1/requests/{id}` | `Bearer JWT` | **Paso 11:** Detalle de solicitud y seguimiento de viaje |
| `POST` | `/api/v1/requests/{id}/reviews` | `Bearer JWT` (estricto) | **Paso 12:** Calificación transaccional 1-5 estrellas + 5 aspectos destacables |
| `POST` | `/api/v1/professionals/{id}/verification` | `Bearer JWT` (estricto) | **Paso 13:** Carga y validación de oficio, cobertura, DNI y matrícula |
| `GET` | `/api/v1/professionals/{id}/dashboard` | `Bearer JWT` (estricto) | **Paso 14:** Panel del profesional, métricas y solicitudes recibidas |
| `PATCH` | `/api/v1/professionals/{id}/availability` | `Bearer JWT` (estricto) | **Paso 14:** Cambio de disponibilidad en tiempo real |
