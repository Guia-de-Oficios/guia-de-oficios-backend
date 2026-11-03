# Guía de Trabajo Colaborativo y Reparto de Tareas por Etapas — "Guía de Oficios"

> [!NOTE]
> Este documento está diseñado para que ambos integrantes del equipo tengan claridad total sobre cómo se estructura la organización **[Guia-de-Oficios](https://github.com/orgs/Guia-de-Oficios/repositories)**, cómo levantar el entorno local sin saturar la memoria RAM (usando **Docker** para el Backend + PostGIS y **Wi-Fi ADB / Teléfono Físico** para Android Studio), y cómo avanzar etapa por etapa sin usar mocks.

---

## 1. Arquitectura de Repositorios en la Organización

El proyecto está dividido en **3 repositorios oficiales** dentro de `https://github.com/Guia-de-Oficios`:

| Repositorio | Propósito Principal | Entorno Recomendado |
| :--- | :--- | :--- |
| **[`guia-de-oficios-backend`](https://github.com/Guia-de-Oficios/guia-de-oficios-backend)** | API REST pura en **Golang** (sin HTML ni mocks), migraciones y semillas de **PostgreSQL 16 + PostGIS**, `Dockerfile` + `docker-compose.yml` y tests unitarios. | **Docker Compose** (o GoLand solo cuando se edita código Go) |
| **[`guia-de-oficios-android`](https://github.com/Guia-de-Oficios/guia-de-oficios-android)** | Aplicación móvil nativa Android en **Kotlin + Jetpack Compose**, cliente Retrofit (`GuiaDeOficiosApi`), persistencia JWT en `DataStore` y pantallas 1:1 con Stitch. | **Android Studio (JBR 21 + Gradle 8.13)** |
| **[`guia-de-oficios-docs-ui`](https://github.com/Guia-de-Oficios/guia-de-oficios-docs-ui)** | Las **15 pantallas y activos originales exportados de Stitch** (`code.html` + `screen.png`), el sistema de diseño `DESIGN.md` y la documentación maestra de Producto, Etapas y Negocio. | **VS Code / Navegador** |

*(Cada repositorio mantiene sincronizada su carpeta `/docs`, incluyendo **[`docs/06_plan_etapas_funcionales_y_tareas.md`](06_plan_etapas_funcionales_y_tareas.md)** con el checklist de tareas de cada etapa).*

---

## 2. Cómo Levantar el Entorno Consumiendo el Mínimo de Memoria RAM

Para evitar que Windows se quede sin memoria virtual al abrir dos IDEs pesados + el emulador de Android:

### Paso A: Levantar PostgreSQL/PostGIS + API Go en Docker (~35 MB de RAM total)
No hace falta abrir GoLand para usar o probar la aplicación desde Android Studio. Basta con ejecutar en `backend/`:
```powershell
docker compose up -d --build
```
Esto levanta dos contenedores en segundo plano:
* `guia_oficios_db` (`postgis/postgis:16-3.4-alpine` en el puerto `5432`)
* `guia_oficios_api` (`backend-api` en el puerto `8080`)

### Paso B: Probar la App en un Celular Android Físico (por Wi-Fi ADB o Cable USB)
En lugar de levantar un emulador pesado de 3 GB de RAM:
1. En el celular Android ir a **Ajustes -> Opciones de desarrollador** y activar **Depuración por USB** (o **Depuración inalámbrica**) e **Instalar vía USB**.
2. Emparejar el teléfono desde **Android Studio -> Device Manager -> Pair Devices Using Wi-Fi** (escaneando el QR).
3. El cliente HTTP [`RetrofitClient.kt`](../app/src/main/java/ar/guiaoficios/app/data/api/RetrofitClient.kt) tiene failover automático entre `127.0.0.1:8080`, la IP LAN de la PC y `10.0.2.2:8080`. Si se usa túnel ADB:
   ```powershell
   adb reverse tcp:8080 tcp:8080
   ```

---

## 3. División del Trabajo en las 4 Etapas Funcionales

El detalle completo de tareas con sus códigos (`[ETAPA2-BE-01]`, `[ETAPA2-AND-01]`, etc.) se encuentra en **[`docs/06_plan_etapas_funcionales_y_tareas.md`](06_plan_etapas_funcionales_y_tareas.md)**:

* ✅ **Etapa 1 (Completada):** Identidad (`logo_provisional_gdo`), Bienvenida y Login con Google/Correo/Teléfono (`bienvenida_y_selecci_n_de_rol`), Registro Manual y Google (`registro_google_o_manual`), Inicio del Cliente y 8 Categorías (`cliente_inicio_y_categor_as`) + Pantalla `Mi Cuenta` con cierre de sesión real.
* 🔜 **Etapa 2 (En curso para continuar):** Resultados de búsqueda geoespacial con PostGIS (`cliente_resultados_de_profesionales`), Perfil completo del Profesional con galería Coil (`cliente_perfil_del_profesional_y_contacto`) y Cuadro de Honor mensual (`comunidad_top_profesionales_del_mes`).
* ⏳ **Etapa 3:** Solicitud Guiada con foto real (`cliente_solicitud_guiada_a_roberto`), Confirmación y contacto WhatsApp/Llamada (`cliente_elegir_contacto_y_confirmaci_n`), Seguimiento en vivo en mapa (`cliente_profesional_en_camino_y_mapa_en_vivo`) y Calificación transaccional (`cliente_trabajo_terminado_y_calificaci_n`).
* ⏳ **Etapa 4:** Registro y Verificación DNI + Matrícula del Profesional (`profesional_registro_y_verificaci_n_dni`), Panel de Solicitudes Recibidas (`profesional_panel_y_solicitudes_recibidas`) y Aceptar Trabajo / Iniciar Viaje (`profesional_aceptar_trabajo_e_iniciar_viaje`).

---

## 4. Flujo de Git Acordado

1. **`main`**: Rama estable que siempre compila (`./gradlew assembleDebug` y `go test ./...`).
2. **Ramas por tarea de etapa (`feature/etapa2-...`)**:
   * Ejemplo Backend: `git checkout -b feature/etapa2-podio-barrial`
   * Ejemplo Android: `git checkout -b feature/etapa2-perfil-profesional`
3. Cada vez que se completa una tarea de la tabla en `06_plan_etapas_funcionales_y_tareas.md`, se marca como `[x]` y se abre Pull Request hacia `main`.
