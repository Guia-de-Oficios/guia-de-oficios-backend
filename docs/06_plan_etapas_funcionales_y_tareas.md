# Plan de Ejecución por Etapas Funcionales y Tareas del Equipo — "Guía de Oficios"

> [!IMPORTANT]
> **Filosofía de Trabajo del Equipo:**
> 1. **Cero Mocks, Cero Demos, Cero HTML en la API:** El backend en Go (`guia-de-oficios-backend`) es exclusivamente una API REST JSON conectada a **PostgreSQL 16 + PostGIS** en Docker. No sirve vistas HTML ni utiliza usuarios o UUIDs cableados por defecto.
> 2. **Fidelidad 1:1 con las 15 Vistas de Stitch:** Cada pantalla en Android (`guia-de-oficios-android`) debe respetar estrictamente el diseño gráfico (`code.html` + `screen.png` + `DESIGN.md`) alojado en `guia-de-oficios-docs-ui`.
> 3. **Avance Sólido Step-by-Step:** El proyecto se divide en **4 Etapas Funcionales**. No se avanza a la siguiente etapa hasta que la etapa actual esté 100% funcional, conectada a PostgreSQL y validada en dispositivo físico/emulador.

---

## 📦 Inventario Completo de las 15 Vistas Oficiales de Stitch

| # | Carpeta en `ui-stitch/` | Etapa | Rol | Propósito Funcional |
| :---: | :--- | :---: | :---: | :--- |
| **1** | `logo_provisional_gdo` | **Etapa 1** | Transversal | Escudo oficial GDO (`#1E40AF`) + tilde verde (`#15803D`) + identidad de marca. |
| **2** | `bienvenida_y_selecci_n_de_rol` | **Etapa 1** | Ambos | Puerta de entrada: Login con Google, Login con Correo/Teléfono y acceso a Registro con elección de rol. |
| **3** | `registro_google_o_manual` | **Etapa 1** | Ambos | Creación de cuenta con Google o formulario manual (Rol Vecino/Profesional, Nombre, Teléfono, Email, Contraseña con visibilidad). |
| **4** | `cliente_inicio_y_categor_as` | **Etapa 1** | Cliente | Home del vecino: Saludo real, cambio de dirección, buscador, asistente diagnóstico, grilla de 8 oficios y tarjeta Top del Mes. |
| **5** | `cliente_resultados_de_profesionales` | **Etapa 2** | Cliente | Listado geoespacial (PostGIS) con chips de filtro (`Todos`, `Urgencias`, `Matriculados`, `Mejor puntuados`) y tarjeta de ayuda vecinal. |
| **6** | `cliente_perfil_del_profesional_y_contacto` | **Etapa 2** | Cliente | Ficha detallada del profesional: foto, matrícula verificada, años de experiencia, galería de trabajos (Coil), reseñas y barra fija de contacto. |
| **7** | `comunidad_top_profesionales_del_mes` | **Etapa 2** | Cliente | Cuadro de honor y podio barrial (#1, #2, #3) con medallas, puntajes y menciones votadas por los vecinos del barrio. |
| **8** | `gu_a_de_oficios_servicios_y_oficios_vecinales` | **Etapa 2** | Cliente | Catálogo extendido de servicios barriales y garantías de seguridad de la red vecinal. |
| **9** | `cliente_solicitud_guiada_a_roberto` | **Etapa 3** | Cliente | Paso 1 de 2: Selección de urgencia (`Lo antes posible` / `Coordinar visita`), descripción del problema y subida real de foto (`Multipart`). |
| **10** | `cliente_elegir_contacto_y_confirmaci_n` | **Etapa 3** | Cliente | Paso 2 de 2: Resumen del pedido registrado en PostgreSQL y disparo de Intent nativo a WhatsApp o Llamada Telefónica. |
| **11** | `cliente_profesional_en_camino_y_mapa_en_vivo` | **Etapa 3** | Cliente | Seguimiento en tiempo real cuando el profesional acepta el trabajo e inicia el viaje hacia el domicilio del vecino (con ETA y mapa). |
| **12** | `cliente_trabajo_terminado_y_calificaci_n` | **Etapa 3** | Cliente | Cierre del servicio: calificación de 1 a 5 estrellas, selección de los 5 aspectos destacables y comentario verificado. |
| **13** | `profesional_registro_y_verificaci_n_dni` | **Etapa 4** | Profesional | Alta de oficio, especialidades, radio de cobertura en km, carga de fotos de DNI (Frente y Dorso) y Matrícula oficial. |
| **14** | `profesional_panel_y_solicitudes_recibidas` | **Etapa 4** | Profesional | Panel de control del profesional: switch de disponibilidad en vivo, Bento Grid de métricas y bandeja de solicitudes de vecinos. |
| **15** | `profesional_aceptar_trabajo_e_iniciar_viaje` | **Etapa 4** | Profesional | Vista operativa del profesional para aceptar una solicitud urgente/programada, confirmar horario e iniciar viaje hacia la casa del vecino. |

---

## ✅ ETAPA 1: Identidad, Autenticación Real y Home de Categorías (COMPLETADA)

### Resumen de lo implementado y consolidado
1. **Limpieza del Backend Go (`guia-de-oficios-backend`):**
   * Eliminada la ruta `/ui/` y cualquier archivo HTML de demostración en `cmd/api/main.go`.
   * Eliminados los UUIDs por defecto (`11111111-...` / `22222222-...`) en `GET /api/v1/auth/me` y `PATCH /api/v1/users/me/address`; ahora requieren JWT válido (`401 Unauthorized` si no hay sesión).
   * `POST /api/v1/auth/login` soporta ingreso por **correo electrónico o número de teléfono celular** + contraseña.
   * `POST /api/v1/auth/google` eliminó la cuenta ficticia `vecino.google@guiaoficios.ar` y persiste la cuenta real de Google del usuario en PostgreSQL.
   * Dockerizado completo (`guia_oficios_db` + `guia_oficios_api`) para no consumir RAM innecesaria con múltiples IDEs abiertos.
2. **Consolidación en Android (`guia-de-oficios-android`):**
   * Eliminado el usuario hardcodeado `"María Elena Castro"` al iniciar `MainViewModel.kt` y `SessionManager.kt`.
   * Creado `StitchBrandComponents.kt` con el logo vectorial oficial `GdoBrandLogo`, el ícono de 4 colores `GoogleColorIcon`, el `StitchTopAppBar` y el `StitchBottomNavBar` (`Inicio`, `Mis contactos`, `Ayuda`, `Mi cuenta`).
   * `WelcomeRoleScreen`: Diseño 1:1 con Stitch, selector de cuenta Google nativo de Android (`AccountManager`) + diálogo de respaldo, diálogo de login por correo/teléfono con toggle de visibilidad de contraseña, y tarjeta inferior de acceso a registro.
   * `RegisterAccountScreen`: Diseño 1:1 con Stitch, tarjetas de rol de 2 columnas con radio button e ícono (`Busco oficios` vs `Soy profesional`), 4 campos con encabezados e íconos, ojo para mostrar/ocultar contraseña y validaciones reales.
   * `HomeScreen` + `UserAccountScreen`: Saludo con el nombre real del usuario logueado, cambio de dirección persistido en PostgreSQL, buscador, asistente diagnóstico, grilla de 8 rubros con íconos circulares `#EAEDFF` y badges (`Matriculados`, `Urgencias`), tarjeta de *Top Profesionales del Mes*, y pantalla **Mi Cuenta** con botón de **Cerrar sesión**.

---

## 🔜 ETAPA 2: Búsqueda Geoespacial PostGIS, Perfil del Profesional y Podio del Mes

> **Objetivo para el equipo:** Que el vecino pueda tocar cualquiera de las 8 categorías, buscar por texto o entrar al *Top Profesionales del Mes*, viendo resultados reales ordenados por cercanía en PostGIS y accediendo a la ficha completa del profesional tal cual las vistas de Stitch.

### Tareas de Backend (Golang + PostgreSQL/PostGIS)
- [ ] **`[ETAPA2-BE-01]` Enriquecer Semilla y Consultas de Profesionales por Rubro (`migrations/002_seed.sql` y `store.go`):**
  * Asegurar que los 8 rubros (`Plomero`, `Electricista`, `Gasista`, `Aire acondicionado`, `Pintor`, `Albañil`, `Jardinero`, `Cerrajero`) tengan profesionales reales cargados en distintas coordenadas de CABA/GBA para probar filtros y distancias reales con `ST_Distance`.
- [ ] **`[ETAPA2-BE-02]` Endpoint de Podio Barrial (`GET /api/v1/professionals/top-month` o filtro `top_rated=true`):**
  * Devolver el ranking mensual ordenado por `rating_avg DESC, reviews_count DESC, jobs_completed_count DESC`, incluyendo medalla (`#1 Oro`, `#2 Plata`, `#3 Bronce`) y etiquetas destacadas para alimentar `comunidad_top_profesionales_del_mes`.
- [ ] **`[ETAPA2-BE-03]` Validación de Portfolio y Reseñas en `GET /api/v1/professionals/{id}`:**
  * Verificar que la ficha del profesional devuelva correctamente `portfolio`, `reviews`, `specialties`, `license_number`, `visit_price_label`, `years_experience` y `guarantee_percentage`.

### Tareas de Android (Kotlin + Jetpack Compose)
- [ ] **`[ETAPA2-AND-01]` Pulir `SearchResultsScreen` (`cliente_resultados_de_profesionales`):**
  * Alinear 1:1 con `cliente_resultados_de_profesionales/code.html` y `screen.png`:
    * Header con botón volver, título del rubro elegido y subtítulo con el barrio del usuario (`Caballito, CABA`).
    * Fila de chips horizontales (`Todos`, `Urgencias`, `Matriculados`, `Mejor puntuados`) que disparan la consulta real a `GET /api/v1/professionals`.
    * Tarjetas de profesionales con avatar, badge verde `DISPONIBLE AHORA`, insignias `DNI Verificado` y `Matrícula oficial`, distancia en cuadras/km calculada por PostGIS, precio de visita y botón `"Ver perfil completo y teléfonos →"`.
- [ ] **`[ETAPA2-AND-02]` Pulir `ProfessionalProfileScreen` (`cliente_perfil_del_profesional_y_contacto`):**
  * Alinear 1:1 con `cliente_perfil_del_profesional_y_contacto/code.html`:
    * Eliminar cualquier dato fallback en `MainActivity.kt` al abrir el perfil (mostrar estado de carga mientras responde `GET /api/v1/professionals/{id}`).
    * Renderizar fotos reales de trabajos realizados (`AsyncImage` de Coil), desglose de precios orientativos, opiniones verificadas de vecinos y barra inferior fija con los botones **"Contactar por WhatsApp"** y **"Llamar por teléfono"**.
- [ ] **`[ETAPA2-AND-03]` Implementar `TopProfessionalsMonthScreen` (`comunidad_top_profesionales_del_mes`):**
  * Crear el Composable dedicado para la pantalla `comunidad_top_profesionales_del_mes/code.html` (accesible desde la tarjeta *"🏆 Top Profesionales del Mes"* del Home).
  * Mostrar el podio barrial (#1, #2, #3), menciones de honor votadas por los vecinos y acceso directo al perfil de cada profesional.

---

## ⏳ ETAPA 3: Solicitud Guiada con Fotos, Mapa en Vivo y Calificación del Vecino

> **Objetivo para el equipo:** Que el vecino pueda crear una solicitud real adjuntando una foto de su problema, contactar al profesional por WhatsApp/Teléfono, seguir en vivo cuando el profesional está en camino, y calificar el trabajo finalizado actualizando el promedio en PostgreSQL.

### Tareas de Backend (Golang + PostgreSQL/PostGIS)
- [ ] **`[ETAPA3-BE-01]` Estados de Solicitud y Tracking en Vivo (`internal/handlers/api.go` y `store.go`):**
  * Soportar la transición de estados de una solicitud en la tabla `service_requests`:
    * `pending` (creada por el vecino)
    * `accepted` / `en_route` (aceptada por el profesional e iniciado el viaje hacia el domicilio)
    * `completed` (trabajo terminado y listo para calificar)
    * `cancelled`
  * Exponer `PATCH /api/v1/requests/{id}/status` para actualizar el estado y retornar las coordenadas/ETA para la pantalla `cliente_profesional_en_camino_y_mapa_en_vivo`.
- [ ] **`[ETAPA3-BE-02]` Historial de Solicitudes del Vecino (`GET /api/v1/requests/my`):**
  * Devolver todas las solicitudes del vecino autenticado (extraído del JWT) ordenadas por fecha para alimentar la pestaña **"Mis contactos"** del `StitchBottomNavBar`.
- [ ] **`[ETAPA3-BE-03]` Transacción de Reseña (`POST /api/v1/requests/{id}/reviews`):**
  * Asegurar que al calificar una solicitud se valide que no haya sido calificada antes, se inserte en `reviews` con los `highlights` elegidos y se recalcule `rating_avg`, `reviews_count` y `jobs_completed_count` del profesional.

### Tareas de Android (Kotlin + Jetpack Compose)
- [ ] **`[ETAPA3-AND-01]` Pulir `GuidedRequestScreen` (`cliente_solicitud_guiada_a_roberto`):**
  * Alinear 1:1 con `cliente_solicitud_guiada_a_roberto/code.html`:
    * Indicador visual `Paso 1 de 2`, tarjetas de urgencia (`Lo antes posible (Urgencia)` vs `Coordinar visita en la semana`), campo de texto descriptivo con sugerencias claras, y selector real de cámara/galería (`PickVisualMedia`) que sube la imagen a `POST /api/v1/uploads` y muestra miniatura de previsualización.
- [ ] **`[ETAPA3-AND-02]` Pulir `RequestConfirmationScreen` (`cliente_elegir_contacto_y_confirmaci_n`):**
  * Alinear 1:1 con `cliente_elegir_contacto_y_confirmaci_n/code.html`:
    * Mostrar el resumen real de la solicitud creada en PostgreSQL, botones de acción directa a **WhatsApp** (con texto pre-armado) y **Llamada Telefónica**, y botón para ver el estado del viaje o calificar el trabajo.
- [ ] **`[ETAPA3-AND-03]` Implementar `ProfessionalEnRouteMapScreen` (`cliente_profesional_en_camino_y_mapa_en_vivo`):**
  * Crear el Composable correspondiente a `cliente_profesional_en_camino_y_mapa_en_vivo/code.html`:
    * Vista de mapa barrial con el trayecto entre el profesional y el domicilio del vecino, tarjeta flotante con foto del profesional, matrícula verificada, tiempo estimado de llegada (`Llega en ~12 min`), y botones de contacto rápido / marcar trabajo como terminado.
- [ ] **`[ETAPA3-AND-04]` Pulir `JobCompletedReviewScreen` (`cliente_trabajo_terminado_y_calificaci_n`) y Pestaña `Mis contactos`:**
  * Conectar la pestaña **"Mis contactos"** del `StitchBottomNavBar` a `GET /api/v1/requests/my` para que el vecino vea sus pedidos activos/históricos.
  * Alinear 1:1 la pantalla de calificación de 1 a 5 estrellas + chips de los 5 aspectos destacables (`Puntualidad`, `Precio justo`, `Honestidad y prolijidad`, `Excelente atención`, `Solución en el acto`).

---

## ⏳ ETAPA 4: Flujo Completo del Profesional (Verificación DNI, Panel y Aceptar Viaje)

> **Objetivo para el equipo:** Que un usuario registrado con rol `professional` pueda completar su legajo de seguridad (DNI Frente/Dorso + Matrícula), gestionar su disponibilidad en tiempo real, ver los pedidos que le envían los vecinos y aceptar un trabajo iniciando el viaje.

### Tareas de Backend (Golang + PostgreSQL/PostGIS)
- [ ] **`[ETAPA4-BE-01]` Vinculación Automática de Perfil Profesional en JWT y `/api/v1/auth/me`:**
  * Verificar que todo usuario registrado con `role = 'professional'` tenga siempre su `professional_id` asociado y que `POST /api/v1/professionals/{id}/verification` actualice `is_dni_verified = TRUE` y `is_licensed = TRUE` al adjuntar los documentos reales.
- [ ] **`[ETAPA4-BE-02]` Endpoint de Aceptación de Trabajo e Inicio de Viaje (`PATCH /api/v1/requests/{id}/accept`):**
  * Permitir que el profesional acepte una solicitud recibida en su panel, actualizando el estado a `en_route`, registrando la nota de horario estimado y habilitando la vista de viaje en curso.

### Tareas de Android (Kotlin + Jetpack Compose)
- [ ] **`[ETAPA4-AND-01]` Pulir `ProfessionalVerificationScreen` (`profesional_registro_y_verificaci_n_dni`):**
  * Alinear 1:1 con `profesional_registro_y_verificaci_n_dni/code.html` y `screen.png`:
    * Selector de oficio principal, chips de especialidades, radio de cobertura en kilómetros, subida real de fotos de **DNI Frente**, **DNI Dorso** y **Credencial/Matrícula** mediante `POST /api/v1/uploads` y guardado en PostgreSQL.
- [ ] **`[ETAPA4-AND-02]` Pulir `ProfessionalDashboardScreen` (`profesional_panel_y_solicitudes_recibidas`):**
  * Alinear 1:1 con `profesional_panel_y_solicitudes_recibidas/code.html` y `screen.png`:
    * Switch grande de disponibilidad en tiempo real (`Estás DISPONIBLE para recibir trabajos` / `Pausado`), Bento Grid con métricas reales (`Nuevas hoy`, `Urgencias`, `Contactos del mes`), y tarjetas de solicitudes de vecinos sin datos hardcodeados.
- [ ] **`[ETAPA4-AND-03]` Implementar `ProfessionalAcceptJobScreen` (`profesional_aceptar_trabajo_e_iniciar_viaje`):**
  * Crear el Composable correspondiente a `profesional_aceptar_trabajo_e_iniciar_viaje/code.html` y `screen.png`:
    * Vista donde el profesional revisa el pedido urgente del vecino (con la foto del desperfecto adjuntada, dirección exacta y distancia), confirma la aceptación del trabajo e inicia el viaje hacia el domicilio o abre WhatsApp/Llamada con el vecino.

---

## 🛠️ Guía Rápida de Entorno para el Compañero de Equipo

1. **Levantar Base de Datos (PostgreSQL 16 + PostGIS) y API Go en Docker (sin abrir GoLand):**
   ```powershell
   cd backend
   docker compose up -d --build
   ```
   * La API queda escuchando en `http://localhost:8080/api/v1/health`.
   * Consume ~15 MB de RAM en lugar de abrir un segundo IDE pesado.

2. **Probar desde un Celular Android Físico (Recomendado para ahorrar RAM) o Emulador:**
   * En el celular: activar *Opciones de desarrollador* -> *Depuración por USB* (o *Depuración inalámbrica*) + *Instalar vía USB*.
   * Si se conecta por cable o Wi-Fi ADB, ejecutar:
     ```powershell
     adb reverse tcp:8080 tcp:8080
     ```
   * Además, `RetrofitClient.kt` ya incluye **failover automático de IP** (`127.0.0.1:8080` -> IP LAN Wi-Fi -> `10.0.2.2:8080`), por lo que la app se conecta automáticamente tanto en teléfono físico como en emulador.
