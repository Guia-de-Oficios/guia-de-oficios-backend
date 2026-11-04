# Guía de Trabajo Colaborativo y Reparto de Tareas por Pasos — "Guía de Oficios"

> [!NOTE]
> Este documento está diseñado para que ambos integrantes del equipo tengan claridad total sobre cómo se estructura la organización **[Guia-de-Oficios](https://github.com/orgs/Guia-de-Oficios/repositories)**, cómo levantar el entorno local sin saturar la memoria RAM (usando **Docker** para el Backend + PostGIS y **Wi-Fi ADB / Teléfono Físico** para Android Studio), y cómo avanzar paso a paso sin usar mocks siguiendo el plan de 14 pasos.

---

## 1. Arquitectura de Repositorios en la Organización

El proyecto está dividido en **3 repositorios oficiales** dentro de `https://github.com/Guia-de-Oficios`:

| Repositorio | Propósito Principal | Entorno Recomendado |
| :--- | :--- | :--- |
| **[`guia-de-oficios-backend`](https://github.com/Guia-de-Oficios/guia-de-oficios-backend)** | API REST pura en **Golang** (sin HTML ni mocks), migraciones y semillas de **PostgreSQL 16 + PostGIS**, `Dockerfile` + `docker-compose.yml` y tests unitarios. | **Docker Compose** (o GoLand solo cuando se edita código Go) |
| **[`guia-de-oficios-android`](https://github.com/Guia-de-Oficios/guia-de-oficios-android)** | Aplicación móvil nativa Android en **Kotlin + Jetpack Compose**, cliente Retrofit (`GuiaDeOficiosApi`), persistencia JWT en `DataStore` y pantallas 1:1 con Stitch. | **Android Studio (JBR 21 + Gradle 8.13)** |
| **[`guia-de-oficios-docs-ui`](https://github.com/Guia-de-Oficios/guia-de-oficios-docs-ui)** | Las **15 pantallas y activos originales exportados de Stitch** (`code.html` + `screen.png`), el sistema de diseño `DESIGN.md` y la documentación maestra de Producto, Pasos y Negocio. | **VS Code / Navegador** |

*(Cada repositorio mantiene sincronizada su carpeta `/docs`, incluyendo **[`docs/06_radiografia_vistas_y_plan_paso_a_paso.md`](06_radiografia_vistas_y_plan_paso_a_paso.md)** con el checklist de tareas de cada paso).*

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

## 3. División del Trabajo en los 14 Pasos Secuenciales

El detalle completo de tareas con sus códigos (`[BE-1.5]`, `[AND-1.1]`, etc.) se encuentra en **[`docs/06_radiografia_vistas_y_plan_paso_a_paso.md`](06_radiografia_vistas_y_plan_paso_a_paso.md)**:

* **Paso 1:** Modularización de Pantallas, Selector de Rol Inicial y Login Bien Hecho (Vistas #1, #2 y #8).
* **Paso 2:** Registro Manual Completo, Dirección Barrial y Google Auth Seguro (Vista #3).
* **Paso 3:** Redirección Inteligente por Rol, Sesión JWT y Blindaje de Endpoints.
* **Paso 4:** Barra Superior, Barra Inferior y Pantalla 'Mi Cuenta' Diferenciada por Rol.
* **Paso 5:** Home del Vecino, Diagnóstico Inteligente y Semilla de los 8 Rubros (Vista #4).
* **Paso 6:** Pantalla de Resultados PostGIS 1:1 con Stitch (Vista #5).
* **Paso 7:** Ficha Completa del Profesional y Galería de Fotos Antes/Después (Vista #6).
* **Paso 8:** Pantalla Comunidad: Top Profesionales del Mes / Podio Barrial (Vista #7).
* **Paso 9:** Solicitud Guiada (Paso 1 de 2) 1:1 con Stitch y Foto con Preview (Vista #9).
* **Paso 10:** Confirmación de Solicitud (Paso 2 de 2) y Disparo a WhatsApp / Teléfono (Vista #10).
* **Paso 11:** Seguimiento en Vivo con Mapa y Pestaña 'Mis Contactos' (Vista #11).
* **Paso 12:** Trabajo Terminado y Calificación Transaccional 1 a 5 Estrellas (Vista #12).
* **Paso 13:** Onboarding del Profesional: Rubro, Especialidades y Fotos Reales de DNI/Matrícula (Vista #13).
* **Paso 14:** Panel del Profesional, Vista 'Aceptar Trabajo e Iniciar Viaje' y APK Final (Vistas #14 y #15).

---

## 4. Flujo de Git Acordado

1. **`main`**: Rama estable que siempre compila (`./gradlew assembleDebug` y `go test ./...`).
2. **Regla de oro anti-conflictos:** Refactorizar primero `StitchWorkflowScreens.kt` en archivos individuales (`RoleSelectionScreen.kt`, `WelcomeLoginScreen.kt`, etc.) en el Paso 1 para evitar choques en el monolito de 2250 líneas.
3. **Ramas exactas por Paso:**
   * Backend: `feat/paso-XX-...` o `fix/paso-XX-...`
   * Android: `feat/paso-XX-...` o `fix/paso-XX-...`
4. Cada vez que se completa un Paso, ambos integrantes verifican que compila y pasa tests, hacen merge a `main`, hacen `git pull` y prueban juntos en el dispositivo antes de iniciar el siguiente Paso.
