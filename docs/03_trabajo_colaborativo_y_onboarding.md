# Guía de Trabajo Colaborativo y Reparto de Tareas — "Guía de Oficios"

> [!NOTE]
> Este documento está pensado para que ambos integrantes del equipo tengan claridad total sobre cómo se estructura la organización **[Guia-de-Oficios](https://github.com/orgs/Guia-de-Oficios/repositories)**, cómo levantar el entorno local en **GoLand**, **Android Studio** y **Docker**, y cómo coordinar el desarrollo sin pisarse código.

---

## 1. Arquitectura de Repositorios en la Organización

El proyecto está dividido en **3 repositorios oficiales** dentro de `https://github.com/Guia-de-Oficios`:

| Repositorio | Propósito Principal | IDE Recomendado |
| :--- | :--- | :--- |
| **[`guia-de-oficios-backend`](https://github.com/Guia-de-Oficios/guia-de-oficios-backend)** | API REST en **Golang**, esquema y semillas de **PostgreSQL 16 + PostGIS**, `docker-compose.yml` y tests unitarios. | **JetBrains GoLand** |
| **[`guia-de-oficios-android`](https://github.com/Guia-de-Oficios/guia-de-oficios-android)** | Aplicación móvil nativa Android en **Kotlin + Jetpack Compose**, cliente Retrofit (`GuiaDeOficiosApi`), DTOs y pantallas basadas en Stitch. | **Android Studio** |
| **[`guia-de-oficios-docs-ui`](https://github.com/Guia-de-Oficios/guia-de-oficios-docs-ui)** | Las **10 pantallas originales exportadas de Stitch** (`code.html` + `screen.png`), el sistema de diseño `DESIGN.md`, el navegador interactivo `app.html` y la documentación maestra de Producto y Negocio. | **VS Code / Navegador** |

*(Además, cada repositorio incluye una carpeta `/docs` sincronizada para que ambos tengan la documentación siempre a mano sin cambiar de ventana).*

---

## 2. Propuesta de División de Trabajo (Para avanzar en paralelo)

Como el contrato de la API REST (`/api/v1/...`) y las 10 pantallas de Stitch ya están definidos y documentados, pueden trabajar **100% en paralelo** desde el día 1:

### Rol A: Especialista Backend & Base de Datos (GoLand + Docker)
1. **Persistencia completa en PostgreSQL + PostGIS:** Extender las consultas espaciales (`ST_DWithin`, `ST_Distance`) en `internal/repository/store.go` para todos los endpoints que hoy tienen fallback en memoria.
2. **Autenticación y Seguridad:** Implementar tokens JWT reales en `POST /api/v1/auth/register` y `POST /api/v1/auth/login`, y validación del `id_token` de Google Sign-In.
3. **Subida de Imágenes (DNI, Matrícula y Fotos de Problemas/Portfolio):** Generar *Presigned URLs* (S3 / Cloudflare R2 / Supabase Storage) para que la app Android suba las fotos de `Foto_caño_cocina.jpg`, `DNI_frente.jpg` y `DNI_dorso.jpg`.
4. **Notificaciones Push (FCM):** Disparar alerta al celular del profesional cuando entra una solicitud marcada como `urgente`.

### Rol B: Especialista Mobile Android (Android Studio + Jetpack Compose)
1. **Navegación y ViewModels (`StateFlow`):** Conectar el grafo de navegación de las 10 pantallas (`StitchNavigation.kt`) con sus `ViewModels` consumiendo `GuiaDeOficiosApi.kt`.
2. **Fidelidad Visual Stitch:** Ajustar márgenes, sombras, íconos Material Symbols y tipografía `Inter` tomando como referencia exacta las capturas `screen.png` de cada carpeta en `ui-stitch/`.
3. **Integración con Cámara/Galería y Intents Nativos:**
   * Selector de foto (`ActivityResultContracts.PickVisualMedia`) en la pantalla 6 (*Solicitud Guiada*) y pantalla 9 (*Verificación DNI*).
   * Intent nativo a **WhatsApp** (`https://wa.me/...`) y al **Marcador Telefónico** (`tel:...`) en las pantallas 5, 7 y 10.

---

## 3. Flujo de Git Acordado (Git Flow Simplificado)

Para evitar conflictos al trabajar juntos:
1. **`main`**: Rama estable que siempre compila y pasa los tests (`go test ./...`).
2. **Ramas por funcionalidad (`feature/...` o `fix/...`)**:
   * Ejemplo Backend: `git checkout -b feature/jwt-auth`
   * Ejemplo Android: `git checkout -b feature/screen-perfil-roberto`
3. **Pull Requests (PRs)**: Cada vez que uno termina una pantalla o un endpoint, abre un PR hacia `main` para que el compañero vea qué se agregó y estén siempre sincronizados.

---

## 4. Cómo Conectar el Emulador de Android Studio con el Backend Local

Cuando uno de ustedes esté corriendo el servidor Go en su PC (`http://localhost:8080`) y quiera probar desde el **Emulador de Android Studio**:
* Recordar que dentro del emulador de Android, `localhost` apunta al teléfono mismo.
* Para apuntar al servidor Go que corre en Windows, usar la IP especial del emulador:
  ```
  http://10.0.2.2:8080/
  ```
* Si prueban con un **teléfono Android físico** conectado al mismo Wi-Fi, usar la IP local de la PC (ej. `http://192.168.1.XX:8080/`).
