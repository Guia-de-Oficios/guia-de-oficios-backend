# Guía de Oficios — Backend API REST (Golang + PostgreSQL/PostGIS en Docker) ⚙️

Repositorio oficial del servidor Backend para **Guía de Oficios**. Es una **API REST JSON 100% pura** (sin vistas HTML ni repositorios falsos en memoria), conectada directamente a **PostgreSQL 16 + PostGIS** para dar soporte a las **15 pantallas oficiales generadas en Stitch**.

## 🔗 Repositorios Hermanos de la Organización
* 📱 **App Android (Kotlin + Jetpack Compose):** [`Guia-de-Oficios/guia-de-oficios-android`](https://github.com/Guia-de-Oficios/guia-de-oficios-android)
* 🎨 **Diseño Stitch (15 Vistas Oficiales) y Documentación Maestra:** [`Guia-de-Oficios/guia-de-oficios-docs-ui`](https://github.com/Guia-de-Oficios/guia-de-oficios-docs-ui)

---

## 📚 Documentación y Plan de Etapas en `/docs`
1. 🚀 **[`docs/06_plan_etapas_funcionales_y_tareas.md`](docs/06_plan_etapas_funcionales_y_tareas.md)**: **[HOJA DE RUTA DEL EQUIPO]** División en **4 Etapas Funcionales**, estado actual de la **Etapa 1 (Completada)** y lista detallada de tareas Backend y Android para las **Etapas 2, 3 y 4**.
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

| Método | Ruta | Autenticación | Etapa / Funcionalidad |
| :--- | :--- | :---: | :--- |
| `GET` | `/api/v1/health` | Pública | **Infra:** Health check contra PostgreSQL + versión de PostGIS |
| `GET` | `/api/v1/metrics` | Pública | **Infra:** Snapshot de métricas clave y señales de oro (Datadog) |
| `POST` | `/api/v1/auth/register` | Pública | **Etapa 1 ✅:** Registro manual (Vecino o Profesional) + emisión de JWT |
| `POST` | `/api/v1/auth/login` | Pública | **Etapa 1 ✅:** Login por **correo electrónico o teléfono celular** + contraseña |
| `POST` | `/api/v1/auth/google` | Pública | **Etapa 1 ✅:** Login/Registro con cuenta Google real del dispositivo Android |
| `GET` | `/api/v1/auth/me` | `Bearer JWT` | **Etapa 1 ✅:** Perfil real del usuario autenticado (`401` si no hay token válido) |
| `PATCH` | `/api/v1/users/me/address` | `Bearer JWT` | **Etapa 1 ✅:** Actualización de dirección guardada y coordenadas en PostgreSQL |
| `GET` | `/api/v1/categories` | Pública | **Etapa 1 ✅:** Grilla de 8 rubros oficiales con badges (`Matriculados`, `Urgencias`) |
| `POST` | `/api/v1/diagnostic` | Pública | **Etapa 1 ✅:** Asistente *"¿No sabés qué profesional necesitás?"* |
| `GET` | `/api/v1/professionals` | Pública | **Etapa 2 🔜:** Búsqueda geoespacial PostGIS (`category_id`, `q`, `only_urgent`, `only_licensed`, `top_rated`) |
| `GET` | `/api/v1/professionals/{id}` | Pública | **Etapa 2 🔜:** Ficha completa del profesional, portfolio de fotos y reseñas |
| `POST` | `/api/v1/uploads` | `Bearer JWT` | **Etapa 3/4 ⏳:** Subida real `multipart/form-data` de fotos de problemas, DNI y matrícula |
| `POST` | `/api/v1/requests` | `Bearer JWT` | **Etapa 3 ⏳:** Creación de solicitud guiada con foto y urgencia |
| `GET` | `/api/v1/requests/my` | `Bearer JWT` | **Etapa 3 ⏳:** Historial de solicitudes del vecino (`Mis contactos`) |
| `GET` | `/api/v1/requests/{id}` | `Bearer JWT` | **Etapa 3 ⏳:** Detalle de solicitud y estado de viaje |
| `POST` | `/api/v1/requests/{id}/reviews` | `Bearer JWT` | **Etapa 3 ⏳:** Calificación transaccional 1-5 estrellas + 5 aspectos destacables |
| `POST` | `/api/v1/professionals/{id}/verification` | `Bearer JWT` | **Etapa 4 ⏳:** Carga de oficio, cobertura, DNI frente/dorso y matrícula |
| `GET` | `/api/v1/professionals/{id}/dashboard` | `Bearer JWT` | **Etapa 4 ⏳:** Panel del profesional, métricas y solicitudes recibidas |
| `PATCH` | `/api/v1/professionals/{id}/availability` | `Bearer JWT` | **Etapa 4 ⏳:** Cambio de disponibilidad en tiempo real |
