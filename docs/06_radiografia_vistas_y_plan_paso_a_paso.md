# MANUAL MAESTRO DE INGENIERÍA · AUDITORÍA DEL ESQUELETO Y PLAN DETALLADO

# Guía de Oficios: Radiografía Vista por Vista del Esqueleto y Plan de Trabajo Detallado

> Documento técnico para desarrollar en equipo (2 personas) revisando y rehaciendo cada pantalla al 100%: qué tenemos hoy como esqueleto base, qué le falta exactamente al Login, Registro y a cada una de las 15 vistas de Stitch, cómo correr todo local y los 14 pasos secuenciales paso a paso.

| ESTADO REAL DEL CÓDIGO | AUDITORÍA DE PANTALLAS | PLAN DE EJECUCIÓN |
| :---: | :---: | :---: |
| **Esqueleto Base (Todo por Pulir)** | **15 Vistas Desmenuzadas al Detalle** | **14 Pasos Secuenciales con Ramas** |

---

## 1. Funcionamiento de la App y Estado Real del Esqueleto Actual

Guía de Oficios (GDO) es una plataforma vecinal hiper-local para CABA y GBA pensada con accesibilidad extrema (WCAG AAA, botones táctiles de 56dp, tipografías grandes y alto contraste) para que vecinos y adultos mayores encuentren profesionales verificados de su propio barrio (Plomeros, Electricistas, Gasistas Matriculados, Aire Acondicionado, Pintores, Albañiles, Jardineros y Cerrajeros).

> [!IMPORTANT]
> **Premisa Fundamental: Lo que hay hoy es un ESQUELETO BASE. Hay que revisar y completar TODO.**
> 1. **Ninguna pantalla está 'terminada':** Si bien el proyecto compila y tiene una navegación inicial, al revisar el código contra los diseños originales de Stitch (`guia-de-oficios-docs-ui`) y las reglas de negocio, **TODAS** las vistas empezando por el Login y el Registro tienen faltantes importantes de UI, validaciones débiles, diálogos improvisados, datos hardcodeados o flujos que se cortan.
> 2. **Por qué no debemos avanzar 'a lo genérico':** Si dejamos el Login en un `AlertDialog` que se cierra solo cuando hay error, o si el Registro no pide dirección/barrio ni redirige según el rol, todas las pantallas siguientes arrastran bugs. La meta de este documento es ir pantalla por pantalla, desmenuzando cada botón, cada campo, cada estado de error y cada endpoint antes de pasar a la siguiente.

### 1.1. Estructura de los 3 Repositorios Locales

| Repositorio | Stack Principal | Qué contiene hoy en el esqueleto y qué debemos llevar a producción |
| :--- | :--- | :--- |
| **`guia-de-oficios-backend`** | Go 1.27 + PostgreSQL 16 + PostGIS + Docker | Contiene las migraciones SQL (`001_init.sql` y `002_seed.sql`) y handlers en `internal/handlers/api.go`. Debemos reemplazar el hash provisorio por bcrypt, exigir JWT en rutas privadas, validar dueño del perfil y completar endpoints de viaje y ranking. |
| **`guia-de-oficios-android`** | Kotlin + Jetpack Compose M3 + Retrofit + Coil | Tiene todas las pantallas concentradas en `HomeScreen.kt` y en un único archivo monolítico (`StitchWorkflowScreens.kt` de 2250 líneas). Debemos modularizar cada vista, corregir el Login/Registro, quitar los fallbacks de Roberto Gómez y crear las 4 vistas faltantes. |
| **`guia-de-oficios-docs-ui`** | HTML + Tailwind (Stitch) + PNG + Markdown | Contiene las 15 carpetas de referencia visual (`code.html` + `screen.png`) y `DESIGN.md`. Es nuestra regla de medición: cada pantalla en Android debe verse y comportarse igual a su `code.html`. |

---

### 1.2. Mapa Maestro de las 15 Vistas de Stitch y su Estado Real en el Esqueleto

| # | Carpeta en `docs-ui` | Rol | Estado en el Esqueleto | Resumen de lo que falta hacer |
| :---: | :--- | :---: | :---: | :--- |
| **1** | `logo_provisional_gdo` | Ambos | **ESQUELETO PARCIAL** | Se usa un ícono genérico de llave inglesa en varios headers en vez del escudo oficial con tilde verde. |
| **2** | `bienvenida_y_selecci_n_de_rol` | Ambos | **INCOMPLETO / CON BUGS** | El Login es un AlertDialog básico que se cierra solo antes de validar; Google Auth es simulado; falta recuperar clave. |
| **3** | `registro_google_o_manual` | Ambos | **INCOMPLETO / CON BUGS** | Validaciones mínimas, no pide dirección/barrio (hardcodea Av. Rivadavia 5420), falta checkbox de Términos y redirige mal al profesional. |
| **4** | `cliente_inicio_y_categor_as` | Vecino | **ESQUELETO PARCIAL** | TopBar repite acción en ambos botones; BottomBar manda 'Mis contactos' al buscador; diagnóstico viene pre-escrito. |
| **5** | `cliente_resultados_de_profesionales` | Vecino | **INCOMPLETO** | Las tarjetas NO muestran la foto del profesional; bug else-if oculta sello de urgencia en matriculados; falta tira de seguridad. |
| **6** | `cliente_perfil_del_profesional_y_contacto` | Vecino | **INCOMPLETO / HARDCODED** | Falta foto en cabecera; falta sección 'Seguridad y Garantía Vecinal' (3 tarjetas); el portfolio NO muestra las fotos. |
| **7** | `comunidad_top_profesionales_del_mes` | Vecino | **NO EXISTE EN ANDROID** | No fue maquetada. Al tocar el banner en Home abre un listado común en vez del Podio #1 Oro, #2 Plata y #3 Bronce. |
| **8** | `gu_a_de_oficios_servicios_y_oficios_vecinales` | Ambos | **NO EXISTE EN ANDROID** | Selector inicial de rol ('Necesito un oficio' vs 'Ofrezco mis servicios') con tarjetas grandes de 96dp y radio buttons. |
| **9** | `cliente_solicitud_guiada_a_roberto` | Vecino | **INCOMPLETO** | Usa chips chicos en vez de las 2 tarjetas grandes de urgencia (72dp); no muestra miniatura de la foto subida ni botón borrar. |
| **10** | `cliente_elegir_contacto_y_confirmaci_n` | Vecino | **INCOMPLETO / HARDCODED** | Falta ícono gigante de confirmación, falta botón 'Ver profesional en camino y mapa en vivo' y tiene fallback hardcodeado. |
| **11** | `cliente_profesional_en_camino_y_mapa_en_vivo` | Vecino | **NO EXISTE EN ANDROID** | Falta toda la pantalla de seguimiento en vivo con ETA ('En aprox. 12 min'), mapa barrial, llamada rápida y checklist. |
| **12** | `cliente_trabajo_terminado_y_calificaci_n` | Vecino | **INCOMPLETO / HARDCODED** | Falta tarjeta resumen del trabajo; usa FilterChips en vez de 5 estrellas grandes de 48dp; trae comentario pre-escrito. |
| **13** | `profesional_registro_y_verificaci_n_dni` | Profesional | **INCOMPLETO / HARDCODED** | No deja elegir rubro (manda siempre 1=Plomero); trae matrícula pre-escrita y manda URLs falsas si no adjuntan DNI. |
| **14** | `profesional_panel_y_solicitudes_recibidas` | Profesional | **INCOMPLETO / HARDCODED** | Falta foto y matrícula en cabecera; las tarjetas de pedidos no tienen click para abrir el detalle del viaje. |
| **15** | `profesional_aceptar_trabajo_e_iniciar_viaje` | Profesional | **NO EXISTE EN ANDROID** | Falta toda la pantalla operativa para ver foto del desperfecto, activar GPS, aceptar trabajo e iniciar viaje al domicilio. |

---

## 2. Guía Completa para Correr la Solución Localmente

Antes de tocar la primera línea del Login, ambos integrantes del equipo deben poder levantar el Backend con PostGIS en Docker y la App en Android Studio conectada a la base de datos local.

### Checklist de Instalación en Windows
- **[OK] Go 1.27.2 windows/amd64:** Instalado en `C:\Program Files\Go\bin\go.exe`.
- **[OK] Docker Desktop:** Instalado y listo para levantar `postgis/postgis:16-3.4-alpine` + API Go.
- **Android Studio:** Descargar de `developer.android.com/studio`. Durante el wizard marcar `Android SDK Platform 34` y `Android SDK Platform-Tools (adb)`. El propio Android Studio instala su JDK 21 embebido.

### 2.1. Levantar Base de Datos PostgreSQL/PostGIS y Backend Go

```powershell
# 1. Entrar al repositorio del backend:
cd backend

# 2. Construir y levantar en segundo plano (ejecuta 001_init.sql y 002_seed.sql automáticamente):
docker compose up -d --build

# 3. Verificar estado de salud y conexión PostGIS:
curl http://localhost:8080/api/v1/health

# 4. Si modificás migrations/001_init.sql o 002_seed.sql y querés recrear la base desde cero:
docker compose down -v && docker compose up -d --build
```

### 2.2. Conectar Android Studio (Celular Físico por Wi-Fi o Cable USB)

- **Opción A — Celular físico con ADB Reverse (Recomendada, consume 0 MB de RAM de emulador):**
  Conectá el celular por USB o por Depuración Inalámbrica (`Device Manager -> Pair using Wi-Fi`).
  Ejecutá `adb reverse tcp:8080 tcp:8080` para que el puerto 8080 de tu PC sea accesible directo desde el teléfono.
- **Opción B — Emulador Android:** Creá un Virtual Device (Pixel 7 API 34). El emulador accede automáticamente a tu localhost:8080 a través de la IP especial `http://10.0.2.2:8080/` configurada en `RetrofitClient.kt`.

### 2.3. Usuarios Semilla Actuales para Pruebas (`migrations/002_seed.sql`)

| Rol en DB | Nombre | Correo Electrónico | Teléfono | Clave Semilla |
| :--- | :--- | :--- | :--- | :--- |
| `client` | María Elena Castro | `maria.castro@vecinos.ar` | `11 5544-3322` | `123456` (password_hash es NULL) |
| `professional` | Roberto Carlos Gómez | `roberto.gomez@guiaoficios.ar` | `+5491145890000` | `123456` (password_hash es NULL) |
| `professional` | Martín Edgardo Silva | `martin.silva@guiaoficios.ar` | `+5491155667788` | `123456` (password_hash es NULL) |

---

## 3. Radiografía Vista por Vista: Qué Tiene el Esqueleto y Qué Falta

A continuación desmenuzamos las 15 vistas del sistema comparando el diseño original de Stitch (`code.html`) contra lo que hay programado hoy en el esqueleto de Android y Backend. Usen esta radiografía como lista de control para no pasar por alto ningún detalle.

---

### VISTA 1 y 8 · Selector Inicial de Rol y Escudo de Marca [NO EXISTE / MEZCLADA]
- **Stitch:** `gu_a_de_oficios_servicios_y_oficios_vecinales/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 96-250)
- **Diseño Stitch:** Pantalla limpia con pill superior *'Buenos Aires y Gran Buenos Aires'*, logo oficial h-16, pregunta *'¿Cómo querés usar Guía de Oficios hoy?'*, 2 tarjetas Radio de 96dp (*'Necesito un oficio para mi casa'* vs *'Ofrezco mis servicios de oficio'*), botón *'Comenzar ->'* de 56dp, link *'¿Ya tenés una cuenta? Iniciar sesión'* y sello *'Red vecinal protegida y verificada'*.
- **Faltantes / Bugs:**
  - **Falta pantalla de selección de rol inicial:** El esqueleto salta directo a los botones de login sin preguntarle primero al usuario nuevo cómo quiere usar la guía ni explicarle la propuesta de valor.
  - **Logo inconsistente:** En vez de renderizar el escudo oficial (`logo_provisional_gdo`), varias pantallas dibujan solo un ícono Build genérico sin el tilde verde de verificación.

---

### VISTA 2 · Bienvenida e Inicio de Sesión (Login Manual y Google) [CRÍTICO · MUY INCOMPLETO]
- **Stitch:** `bienvenida_y_selecci_n_de_rol/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 96-375) + `MainActivity.kt` (L. 142-182) + `jwt.go`
- **Diseño Stitch:** Botones táctiles de 56dp (*'Entrar con Google'*, *'Entrar con correo o teléfono'*), flujo de credenciales accesible con letra grande, manejo de errores en línea sin perder lo escrito, recuperación de contraseña y redirección según rol.
- **Faltantes / Bugs:**
  1. **El Login está escondido en un AlertDialog precario** (`LoginCredentialsDialog`, L. 253-325): En vez de una vista o hoja dedicada accesible para adultos mayores, abre un popup básico con 2 `OutlinedTextField` sin íconos, sin `KeyboardOptions` (email/phone, `ImeAction.Next/Done`) y sin botón de *'¿Olvidaste tu contraseña?'*.
  2. **BUG GRAVE al fallar el Login** (L. 244-247): Al tocar *'Ingresar'*, el código hace `showLoginDialog = false` ANTES de esperar la respuesta del backend. Si el usuario puso mal un dígito de su contraseña, el diálogo se cierra igual, pierde todo lo que escribió y solo sale un Snackbar fugaz abajo.
  3. **Sin estado de carga en el botón de Login:** El botón *'Ingresar'* del diálogo no muestra `CircularProgressIndicator` ni se deshabilita mientras viaja el request HTTP, permitiendo múltiples clicks repetidos.
  4. **Login con Google es FALSO / SIMULADO** (`MainActivity.kt` L. 148-159 y `MainViewModel.kt` L. 224-237): Usa `AccountManager.newChooseAccountIntent` (que solo devuelve el string del email sin autenticar nada) y arma un token inventado `'google-android-account-171...'`. Peor aún: si falla el picker, abre `GoogleAccountInputDialog` (L. 330-375) donde cualquiera escribe cualquier email ajeno y entra sin contraseña.
  5. **Backend no verifica firma de Google** (`store.go` L. 200-265): `AuthenticateGoogle` decodifica el payload base64 del token sin validar la firma criptográfica de Google y, si el email ya existe en la tabla `users`, le entrega un JWT válido sin comprobar nada.
  6. **Contraseñas sin bcrypt y puerta trasera en semilla** (`jwt.go` L. 93-125): El backend usa SHA-256 casero en lugar de bcrypt (costo >= 10) y tiene un `if storedHash == "" { return password == "123456" }` porque en `002_seed.sql` los usuarios tienen `password_hash = NULL`.
  7. **Redirección post-login rota** (`MainActivity.kt` L. 59, 165, 176): Cuando Roberto Gómez (`role='professional'`) inicia sesión, `MainActivity` lo manda siempre a `ClientHomeCategories` (la pantalla de buscar plomeros para vecinos) en vez de llevarlo a su Panel de Profesional.

---

### VISTA 3 · Registro de Cuenta Nueva (Google o Manual) [CRÍTICO · MUY INCOMPLETO]
- **Stitch:** `registro_google_o_manual/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 380-615) + `store.go` (L. 82-155)
- **Diseño Stitch:** Selector de rol (*'Soy Vecino / Cliente'* vs *'Soy Profesional'*) con tilde verde, botón *'Registrarse rápido con Google'*, separador *'O completá tus datos manualmente'*, Nombre, Teléfono con ayudante de WhatsApp, Email opcional, Contraseña con ojito, botón *'Crear mi cuenta gratis'* y aceptación de Términos de Convivencia Vecinal.
- **Faltantes / Bugs:**
  1. **No pide Dirección ni Barrio al registrarse:** En el formulario no hay campo de Barrio ni Calle/Altura. Como consecuencia, el backend en `CreateUser` (`store.go` L. 106) le clava a TODOS los usuarios nuevos la dirección fija `'Av. Rivadavia 5420, Caballito'` y coordenadas `(-34.6186, -58.4368)`.
  2. **Validaciones de campos súper débiles** (L. 566-573): Solo mira `phone.length < 6` (permite letras o teléfonos inválidos en vez de normalizar formato celular argentino 11 / +54 9), el email solo mira `.contains("@")` (acepta '@' solo), y no valida fortaleza mínima ni confirmación de contraseña.
  3. **Falta bloque legal de Términos de Convivencia Vecinal:** En `registro_google_o_manual/code.html` (L. 258-262) hay un texto interactivo obligatorio con enlaces a *'Términos de Convivencia Vecinal'* y *'Política de Privacidad'* que no está en Android.
  4. **Registro de Profesional no crea fila inicial ni redirige al Onboarding DNI** (`MainActivity.kt` L. 206, 218): Si alguien elige *'Soy Profesional'* y se registra, la app lo manda al Home de Vecino en vez de llevarlo al Paso 2 de su alta (`profesional_registro_y_verificaci_n_dni`).
  5. **Falta feedback visual por campo (`isError` + `supportingText`):** Si un campo está mal, hoy tira un texto genérico arriba en vez de marcar en rojo el `OutlinedTextField` exacto y hacer scroll hacia él.

---

### VISTA 4 · Inicio del Vecino, Buscador, Diagnóstico y 8 Categorías [ESQUELETO INCOMPLETO]
- **Stitch:** `cliente_inicio_y_categor_as` (`HomeScreen.kt`) | **Actual:** `HomeScreen.kt` (L. 62-475) + `StitchNavigation.kt` + `MainActivity.kt` (L. 242-279)
- **Diseño Stitch:** TopBar con logo y botón Mi Cuenta, saludo con nombre real, tarjeta de ubicación con cambio de dirección, buscador + botón *'¿No sabés qué profesional necesitás?'*, grilla 2x4 de los 8 rubros, banner *'Top Profesionales del Mes'*, garantías y `BottomNavBar` de 3 pestañas.
- **Faltantes / Bugs:**
  1. **TopAppBar con dos botones que hacen exactamente lo mismo** (`MainActivity.kt` L. 249-250): Tanto el botón izquierdo de menú (`onMenuClick`) como el botón derecho 'Mi Cuenta' (`onActionClick`) abren `UserAccount`.
  2. **BottomNavBar con pestañas 'Mis contactos' y 'Ayuda' rotas** (`MainActivity.kt` L. 116-125): Tocar 'Mis contactos' llama a `searchProfessionals(0)` (abre el buscador de rubros en vez de las solicitudes del vecino). Tocar 'Ayuda' solo dispara un Snackbar en vez de abrir una pantalla de Ayuda/Catálogo.
  3. **Diálogo de Diagnóstico viene pre-rellenado** (`HomeScreen.kt` L. 80): La variable `problemDraft` arranca hardcodeada con *'No enciende el calefón de la cocina'* en vez de estar vacía con sugerencias tocables.
  4. **Cambio de dirección no actualiza coordenadas reales** (`MainViewModel.kt` L. 258-265 y `store.go` L. 267): Aunque el vecino escriba una calle de Flores, Belgrano o Palermo, la app envía siempre `lat=-34.6186` y `lng=-58.4368` (Caballito), por lo que el cálculo de cuadras PostGIS da mal.
  5. **Banner 'Top Profesionales del Mes' mal conectado** (`MainActivity.kt` L. 269): Al tocar *'Ver ranking completo'*, navega a `ClientSearchResults(0)` en vez de abrir la vista dedicada `comunidad_top_profesionales_del_mes`.

---

### VISTA 5 · Resultados de Profesionales por Cercanía (PostGIS) [INCOMPLETO / CON BUGS]
- **Stitch:** `cliente_resultados_de_profesionales/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 1295-1535) + `store.go` (L. 399-473)
- **Diseño Stitch:** Header sticky con *'<- Volver'* y marca, título *'[Rubro] cerca de vos'*, indicador de barrio, chips horizontales con scroll (*'Todos (8)'*, *'Atienden urgencias'*, *'Matriculados'*, *'Mejor valorados'*), aviso azul para mayores, tarjetas completas con foto de 64dp y distancia en cuadras, y caja 0800 al pie.
- **Faltantes / Bugs:**
  1. **¡Las tarjetas NO muestran la foto del profesional!** (`StitchWorkflowScreens.kt` L. 1420-1424): En Stitch cada tarjeta tiene la foto de 64x64dp con borde y badge en la esquina inferior derecha. En el esqueleto actual solo pusieron `Text(pro.full_name)` e ignoraron por completo `pro.avatar_url`.
  2. **Bug lógico 'else if' oculta el sello de Urgencias** (L. 1441-1452): El código hace `if (pro.is_licensed) { ... } else if (pro.attends_emergencies) { ... }`. Por culpa de ese `else if`, si un plomero o gasista es Matriculado Y ADEMÁS atiende urgencias (como Roberto Gómez), el vecino NUNCA ve el sello naranja *'[!] Atiende urgencias hoy'*. Deben poder mostrarse ambos sellos.
  3. **Falta la tira azul 'Profesionales de tu barrio' para adultos mayores** (`code.html` L. 181-189): Card informativa con ícono de escudo que explica que todos tienen identidad presencial validada.
  4. **Fila de chips sin scroll horizontal ni contador** (L. 1344): Usa un `Row` fijo con 3 chips que se cortan en pantallas chicas; debe ser un `LazyRow` con los 4 filtros de Stitch incluyendo el contador *'Todos (N)'* y *'Mejor valorados'*.
  5. **Semilla vacía en 7 de los 8 rubros** (`002_seed.sql`): Solo hay 2 plomeros cargados. Si tocás Electricista, Gasista, Pintor, etc., la lista queda vacía y ni siquiera hay un componente de Estado Vacío (*'No encontramos profesionales con estos filtros — Limpiar filtros'*).

---

### VISTA 6 · Perfil Completo del Profesional y Contacto [INCOMPLETO / HARDCODED]
- **Stitch:** `cliente_perfil_del_profesional_y_contacto/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 1545-1715) + `MainActivity.kt` (L. 337-357)
- **Diseño Stitch:** Barra *'<- Volver a resultados'* + barrio, cabecera con foto de 80dp con borde verde, nombre, oficio, rating, 2 cajas (Experiencia y Distancia), bloque *'SEGURIDAD Y GARANTÍA VECINAL'* (3 tarjetas: DNI, Matrícula, Urgencias), *'Sobre mi trabajo'*, Precio de visita, Galería de fotos h-48 (192dp), Opiniones con estrellas y barra fija inferior con WhatsApp y Llamar.
- **Faltantes / Bugs:**
  1. **Fallback hardcodeado de Roberto Gómez en `MainActivity.kt`** (L. 337-344): Si `selectedProfessional` tarda en cargar, instancia un `ProfessionalSummaryDto` falso de Roberto Gómez en memoria.
  2. **Falta la foto de perfil de 80dp con tilde verde en la cabecera** (L. 1588-1600): Solo muestra textos sin la imagen del trabajador.
  3. **Falta TODA la sección 'SEGURIDAD Y GARANTÍA VECINAL'** (`code.html` L. 173-215): Son 3 tarjetas horizontales clave con íconos grandes (1. Documento de identidad verificado, 2. Matrícula oficial verificada con ente regulador, 3. Disponible hoy para urgencias). En el esqueleto las reemplazaron por un texto plano.
  4. **La galería 'Fotos de trabajos realizados' NO muestra imágenes** (BUG-05, L. 1659-1670): Itera sobre `professional.portfolio` pero solo imprime `Text(item.title)` y `Text(item.subtitle)`. Falta el `AsyncImage(model = item.image_url)` de 192dp de alto con el badge 'Verificado' arriba a la derecha.
  5. **Reseñas de vecinos incompletas** (L. 1672-1688): No dibuja las 5 estrellas doradas, ni los chips de aspectos destacados (highlights), ni el pie verde *'[OK] Trabajo contratado por Guía de Oficios'*.
  6. **Distancia fija en el backend** (`store.go` L. 498): `GetProfessionalByID` devuelve el string clavado `'A 8 cuadras (650 metros)'` en vez de calcular `ST_Distance` con la ubicación del vecino.

---

### VISTA 7 · Comunidad: Top Profesionales del Mes (Podio Barrial) [NO EXISTE EN ANDROID]
- **Stitch:** `comunidad_top_profesionales_del_mes/code.html` | **Actual:** Falta crear `TopProfessionalsMonthScreen.kt`
- **Diseño Stitch:** Header con selector de barrio, Hero Banner azul *'Top Profesionales del Mes — Elegidos por los vecinos de tu barrio'*, Podio Bento Grid (#1 Puesto Oro a ancho completo con corona y comentario destacado, #2 Plata y #3 Bronce en 2 columnas), chips de filtro por oficio, listado de puestos #4 al #10 y tarjeta *'¿Cómo se calcula este ranking?'*.
- **Faltantes / Bugs:**
  1. **Pantalla 100% faltante en Android:** Hay que maquetarla completa en Jetpack Compose siguiendo `comunidad_top_profesionales_del_mes/code.html`.
  2. **Lógica de Podio y Filtros:** Debe permitir filtrar el ranking por rubro (Todos los oficios, Plomeros, Electricistas, Gasistas, Cerrajeros) y al tocar *'Ver perfil completo'* o *'Contactar'* abrir la ficha real del profesional.

---

### VISTA 9 · Solicitud Guiada al Profesional (Paso 1 de 2) [INCOMPLETO]
- **Stitch:** `cliente_solicitud_guiada_a_roberto/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 1720-1875)
- **Diseño Stitch:** Header *'Solicitud directa · Paso 1 de 2'*, tarjeta resumen del profesional con foto, Paso 1 *'¿Con qué urgencia lo necesitás?'* con 2 tarjetas Radio grandes de 72dp (*'Normal · Sin apuro'* vs *'Urgente hoy · Atención rápida'* con punto rojo pulsante), Paso 2 *'Contale brevemente qué pasa'* con chips rápidos, Paso 3 *'Sumá una foto'* con preview real y botón borrar, aviso *'Tu tranquilidad es prioridad'* y modal de confirmación.
- **Faltantes / Bugs:**
  1. **Falta la tarjeta de contexto del profesional arriba** (`code.html` L. 156-179): Con su foto, nombre, oficio y barrio para que el vecino sepa a quién le está escribiendo.
  2. **Selector de urgencia usa FilterChips diminutos en vez de las 2 Tarjetas Radio de 72dp** (L. 1765-1778): Para un adulto mayor, dos chips chiquitos son difíciles de leer y tocar; Stitch exige 2 tarjetas anchas con título, subtítulo explicativo e indicador circular.
  3. **Carga de foto sin miniatura ni opción de quitar** (L. 1808-1834): Al subir una foto con `PickVisualMedia` solo muestra un texto *'Foto cargada y verificada en servidor'*, pero NO muestra la imagen con `AsyncImage` ni tiene el botón *'Eliminar / Cambiar foto'* por si el vecino se equivocó de imagen.
  4. **Backend usa nombre de cliente hardcodeado si no lee la DB** (`store.go` L. 564): En `CreateServiceRequest` el fallback es `'María Elena Castro'` en vez de exigir el usuario autenticado.

---

### VISTA 10 · Elegir Contacto y Confirmación de Solicitud (Paso 2 de 2) [INCOMPLETO / HARDCODED]
- **Stitch:** `cliente_elegir_contacto_y_confirmaci_n/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 1880-1955) + `MainActivity.kt` (L. 391-435)
- **Diseño Stitch:** Header *'Paso 2 de 2'*, tarjeta central con tilde verde gigante de 64dp, *'¡Listo! Ya preparamos tu pedido para [Nombre]'*, caja gris con resumen del problema/urgencia/barrio, pregunta *'¿Cómo preferís hablar con [Nombre]?'*, botón verde *'Escribir por WhatsApp'* (64dp con subtítulo), botón azul *'Llamar por teléfono'* (64dp con número visible), y botón *'Ver profesional en camino y mapa en vivo'*.
- **Faltantes / Bugs:**
  1. **Fallback hardcodeado en `MainActivity.kt`** (L. 391-400): Si `latestServiceRequest` es null, inventa en memoria el pedido `'req-demo-1001'` de Roberto Gómez con *'Pérdida de agua bajo mesada'*.
  2. **Botones de contacto sin subtítulo ni altura de 64dp** (L. 1922-1938): En Stitch, el botón de WhatsApp explica abajo *'Se envía tu problema escrito automáticamente'* y el de Llamar muestra el número formateado *'Hablar directo al 11 ...'*.
  3. **Falta el botón de navegación al Seguimiento en Vivo:** No hay forma de pasar desde esta pantalla (ni desde Mis Contactos) hacia la Vista 11 (`cliente_profesional_en_camino_y_mapa_en_vivo`).

---

### VISTA 11 · Profesional en Camino y Mapa en Vivo [NO EXISTE EN ANDROID]
- **Stitch:** `cliente_profesional_en_camino_y_mapa_en_vivo/code.html` | **Actual:** Falta crear `ProfessionalEnRouteMapScreen.kt` + endpoint PATCH status
- **Diseño Stitch:** Header *'Servicio en curso · Seguimiento en vivo'*, tarjeta verde *'¡[Nombre] aceptó tu solicitud y está en camino!'* con caja azul de ETA (*'Llegada estimada: En aprox. 12 minutos'*), Mapa barrial ilustrado con ruta punteada azul, pin pulsante del profesional y pin de *'Tu casa'*, tarjeta del profesional con botones *'Llamar ahora'* y *'WhatsApp'*, checklist de seguridad y botón para pasar a Trabajo Terminado o cancelar.
- **Faltantes / Bugs:**
  1. **Pantalla 100% faltante en Android:** Debe construirse en Compose respetando todos los bloques de `cliente_profesional_en_camino_y_mapa_en_vivo/code.html`.
  2. **Falta endpoint en Backend para actualizar estado del pedido:** Necesitamos `PATCH /api/v1/requests/{id}/status` para que cuando el profesional inicie el viaje pase a `'en_route'` (con `eta_minutes`) y cuando termine pase a `'completed'`.

---

### VISTA 12 · Trabajo Terminado y Calificación del Vecino [INCOMPLETO / HARDCODED]
- **Stitch:** `cliente_trabajo_terminado_y_calificaci_n/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 1960-2060)
- **Diseño Stitch:** Badge verde *'SERVICIO COMPLETADO'*, título *'¡Trabajo terminado!'*, tarjeta resumen del profesional + trabajo realizado + fecha *'Finalizado hoy'*, selector de 5 estrellas grandes de 48dp con etiqueta dinámica (*'5.0 · ¡Excelente trabajo!'*), 5 tarjetas seleccionables de aspectos destacables con ícono y subtítulo, campo de comentario limpio, aviso de Transparencia Vecinal y barra fija inferior.
- **Faltantes / Bugs:**
  1. **Falta la tarjeta resumen del profesional y del trabajo arriba** (`code.html` L. 163-193): La pantalla actual arranca directo con texto sin mostrar la foto del profesional ni qué trabajo se está calificando.
  2. **Estrellas hechas con FilterChips chiquitos** (*'★ 1'*, *'★ 2'*...) en vez de 5 íconos Star de 48dp (L. 1986-1998): Hay que reemplazarlos por 5 botones grandes de estrella dorada con la píldora inferior que cambia el texto según el puntaje (1='Malo', 3='Cumplió', 5='¡Excelente trabajo!').
  3. **Tarjetas de aspectos destacables incompletas** (L. 2002-2035): En Stitch cada una de las 5 opciones tiene ícono propio, título en negrita y subtítulo explicativo (ej. *'Precio justo'* -> *'Cobró lo presupuestado, sin sorpresas'*). En el esqueleto solo muestran el título.
  4. **Comentario pre-escrito hardcodeado** (L. 1967): La variable `comment` arranca con *'Excelente atención de Roberto, llegó en 20 minutos...'* escrito de fábrica. Debe arrancar vacía (`""`).

---

### VISTA 13 · Onboarding Profesional: Registro de Oficio y Verificación DNI [INCOMPLETO / HARDCODED]
- **Stitch:** `profesional_registro_y_verificaci_n_dni/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 2065-2240) + `store.go` (L. 694-769)
- **Diseño Stitch:** Barra de progreso *'Paso 2 de 3 · Verificación de oficio e identidad · 66%'*, Paso 1 *'¿Cuál es tu oficio principal?'* (selector de los 8 rubros + chips de especialidades secundarias + botón '+ Agregar otra' + radio de cobertura en km), Paso 2 *'Fotos de tu DNI (Frente y Dorso)'* en grilla de 2 columnas con estado subido/pendiente, Paso 3 *'Matrícula o Certificado (Opcional)'* y botón *'Enviar para verificación ->'*.
- **Faltantes / Bugs:**
  1. **¡No tiene selector de los 8 rubros y manda `primary_category_id = 1` clavado!** (L. 2203): Un electricista o gasista no puede elegir su rubro; el código envía siempre 1 (Plomero).
  2. **Campos pre-rellenados con datos de Roberto Gómez** (L. 2073-2076): Arranca con *'Plomero y Gasista Matriculado'*, *'Caballito'*, matrícula *'18.492'* y ente *'ENARGAS'* hardcodeados.
  3. **Manda URLs falsas si el usuario no sube fotos del DNI** (L. 2208-2212): Si `dniFrontUrl` está vacío, envía `'https://guiaoficios.ar/uploads/dni_frente_verificado.jpg'` y el backend lo aprueba como verificado. Debe exigir que suba las 2 fotos reales antes de habilitar el envío.
  4. **Falta barra de progreso del paso (66%)**, selector de especialidades secundarias interactivo y vista previa (`AsyncImage`) de los documentos cargados.

---

### VISTA 14 · Panel del Profesional y Solicitudes Recibidas [INCOMPLETO / HARDCODED]
- **Stitch:** `profesional_panel_y_solicitudes_recibidas/code.html` | **Actual:** `StitchWorkflowScreens.kt` (L. 1030-1215) + `MainActivity.kt` (L. 475-525)
- **Diseño Stitch:** Header *'Panel Profesional'*, tarjeta del trabajador con foto, sello verificado y oficio, Switch grande de disponibilidad (*'Disponible para trabajos hoy'* con punto verde pulsante), Bento Grid de métricas (Calificación media, Reseñas, Trabajos cerrados, Vistas al perfil), lista de *'Solicitudes de vecinos'* con avatar de iniciales, etiqueta Urgente/Normal, hace cuánto llegó y botones de acción, más botón *'Editar mi zona de trabajo y precios'*.
- **Faltantes / Bugs:**
  1. **Fallback hardcodeado en `MainActivity.kt`** (L. 475-489): Mientras carga la API muestra datos falsos de Roberto Gómez.
  2. **Falta la tarjeta de identidad del profesional con foto y oficio** arriba del Switch (`code.html` L. 175-192).
  3. **Las tarjetas de solicitudes no abren la Vista 15** (Aceptar trabajo e iniciar viaje): En el esqueleto solo tienen botones de WhatsApp/Llamar, pero no permiten tocar la solicitud para ver la foto que subió el vecino, aceptar el trabajo y cambiar el estado a `'en_route'`.
  4. **Falta el botón inferior 'Editar mi zona de trabajo y precios'** para que el profesional pueda actualizar su precio de visita, radio en km o rubro.

---

### VISTA 15 · Profesional: Aceptar Trabajo e Iniciar Viaje [NO EXISTE EN ANDROID]
- **Stitch:** `profesional_aceptar_trabajo_e_iniciar_viaje/code.html` | **Actual:** Falta crear `ProfessionalAcceptJobScreen.kt` + endpoint PATCH status
- **Diseño Stitch:** Header *'Solicitud #... · Confirmar e Iniciar Viaje'* con badge 'Urgente', tarjeta del vecino (foto/iniciales, nombre, 'A 6 cuadras · 550m', caja de Dirección Exacta con botón 'Mapa'), bloque 'Detalle del problema' con texto y miniatura de la foto adjunta ('Ver foto' en grande), caja de 'Visita técnica acordada ($12.000)', Switch de compartir ubicación GPS en vivo, botón verde 'Iniciar viaje hacia el domicilio ->', botón 'Abrir en Google Maps / Waze', 'Avisar demora por WhatsApp' y 'Rechazar solicitud'.
- **Faltantes / Bugs:**
  1. **Pantalla 100% faltante en Android:** Es el eslabón que une al Profesional con el Mapa en Vivo del Vecino (Vista 11). Hay que maquetarla completa en Compose.
  2. **Visor de foto del desperfecto y disparo de Google Maps:** Debe mostrar con `AsyncImage` la foto que el vecino subió en el Paso 1 (`problem_photo_url`), permitir ampliarla y abrir Google Maps/Waze con `'geo:0,0?q='` + dirección del cliente.

---

## 4. Cómo Trabajar Ordenados de a 2 con Ramas Git (Branches)

Con este desglose tan detallado, la clave para no pisarse código entre ambos es nunca programar directo sobre `main` y dividir el trabajo por Paso (o uno hace las tareas `[BE]` de Backend y el otro las tareas `[AND]` de Android dentro del mismo Paso).

### Reglas de Oro para Dividir el Trabajo sin Conflictos
1. **Refactorizar primero `StitchWorkflowScreens.kt`:** Hoy ese archivo tiene 2.250 líneas con casi todas las pantallas juntas. Si los dos editan `StitchWorkflowScreens.kt` a la vez, van a tener conflictos de Git todo el tiempo. En el Paso 1 separamos cada pantalla en su propio archivo dentro de `ui/screens/` (ej. `LoginScreen.kt`, `RegisterScreen.kt`, `SearchResultsScreen.kt`, etc.).
2. **Una rama por Paso:** Usen los nombres exactos de rama indicados en la Sección 5. Antes de pasar al Paso N+1, ambos hacen merge a main, hacen `git pull` y prueban juntos en el celular.

```powershell
# Flujo exacto para empezar cada Paso en equipo:
git checkout main
git fetch --all && git pull origin main
git checkout -b feat/paso-01-login-completo

# Verificar que compila y pasa tests antes de subir:
# En backend: go test -v ./...
# En android: .\gradlew.bat assembleDebug

git add .
git commit -m "feat(paso-01): rehacer pantalla de login completa con validaciones y bcrypt"
git push -u origin feat/paso-01-login-completo
```

---

## 5. Plan de Trabajo Paso a Paso (14 Pasos de Principio a Fin)

A continuación están los 14 Pasos ordenados desde el Login hasta el APK final. Cada paso detalla exactamente qué hacer en Backend `[BE]` y en Android `[AND]` para que no quede ningún detalle suelto.

---

### PASO 1 · Modularización de Pantallas, Selector de Rol Inicial y Login Bien Hecho (Vistas #1, #2 y #8)
- **Rama Backend:** `feat/paso-01-auth-bcrypt` | **Rama Android:** `feat/paso-01-login-completo`
- **Objetivo de Calidad:** Separar el archivo monolítico de Compose, crear el selector de rol inicial y rehacer el Login completo (sin el AlertDialog que se cierra solo, con estados de carga, errores por campo y bcrypt real en backend).
- [ ] **`[AND-1.1]`** Separar `StitchWorkflowScreens.kt` en archivos individuales por pantalla dentro de `ui/screens/` (`RoleSelectionScreen.kt`, `WelcomeLoginScreen.kt`, `RegisterScreen.kt`, `SearchResultsScreen.kt`, `ProfessionalProfileScreen.kt`, etc.) para trabajar de a 2 sin conflictos de Git.
- [ ] **`[AND-1.2]`** Crear el componente oficial `GdoLogoBadge` (Vista #1) respetando `logo_provisional_gdo/code.html` (escudo azul `#1E40AF`, llave inglesa blanca y círculo esmeralda `#15803D` con tilde en la esquina inferior derecha).
- [ ] **`[AND-1.3]`** Crear `RoleSelectionScreen.kt` siguiendo `gu_a_de_oficios_servicios_y_oficios_vecinales/code.html` (Vista #8): tarjetas Radio de 96dp (*'Necesito un oficio para mi casa'* vs *'Ofrezco mis servicios de oficio'*), botón *'Comenzar ->'* de 56dp que pre-selecciona el rol para el registro, y enlace *'¿Ya tenés una cuenta? Iniciar sesión'*.
- [ ] **`[AND-1.4]`** Eliminar `LoginCredentialsDialog` (el AlertDialog que hacía `showLoginDialog = false` antes de recibir respuesta de la API) y crear una vista/BottomSheet dedicada de Login con:
  - (a) campo Identificador (Email o Celular) de 56dp con KeyboardOptions,
  - (b) campo Contraseña con ojito mostrar/ocultar,
  - (c) mensaje de error inline en rojo dentro del formulario si las credenciales fallan (sin borrar lo que el usuario escribió),
  - (d) botón 'Ingresar' con spinner `CircularProgressIndicator` mientras `isLoading == true`, y
  - (e) flujo '¿Olvidaste tu contraseña?'.
- [ ] **`[BE-1.5]`** En `internal/auth/jwt.go`, reemplazar el hash SHA-256 casero por `golang.org/x/crypto/bcrypt` (`GenerateFromPassword` y `CompareHashAndPassword`) y actualizar `migrations/002_seed.sql` para que los 3 usuarios semilla tengan su hash bcrypt real (eliminando el backdoor `if storedHash == ""`).

---

### PASO 2 · Registro Manual Completo, Dirección Barrial y Google Auth Seguro (Vista #3)
- **Rama Backend:** `feat/paso-02-registro-barrio` | **Rama Android:** `feat/paso-02-registro-ui`
- **Objetivo de Calidad:** Hacer que el Registro valide cada campo correctamente, guarde el barrio/dirección real del vecino o profesional, incluya los términos legales y elimine el diálogo de Google falso.
- [ ] **`[BE-2.1]`** En `POST /api/v1/auth/register` (`api.go` y `store.go`), aceptar los campos `'neighborhood'` y `'address'` en el registro (en vez de clavar 'Av. Rivadavia 5420, Caballito' para todos) y asignar las coordenadas lat/lng reales del barrio elegido. Si `role == 'professional'`, crear también su registro inicial en la tabla `professionals`.
- [ ] **`[BE-2.2]`** En `POST /api/v1/auth/google` (`store.go`), validar el `id_token` contra las claves públicas de Google (o el endpoint tokeninfo de Google) y no permitir el login solo por mandar un email en texto plano sin firma.
- [ ] **`[AND-2.3]`** En `RegisterScreen.kt` (`registro_google_o_manual/code.html`), implementar:
  - (a) selector de Rol (*'Soy Vecino / Cliente'* vs *'Soy Profesional'*) con borde azul, fondo `#EFF6FF` y tilde verde,
  - (b) validación individual con `isError` + `supportingText` en Nombre (>= 3 letras), Teléfono celular argentino (normalizando `+54 9 11`), Email válido (`Patterns.EMAIL_ADDRESS`), Barrio (selector desplegable de barrios de CABA/GBA) y Dirección, y Contraseña (>= 6 caracteres).
- [ ] **`[AND-2.4]`** Agregar al pie del formulario el bloque de aceptación de *'Términos de Convivencia Vecinal y Política de Privacidad'* con diálogo de lectura accesible, y eliminar `GoogleAccountInputDialog` (el popup que permitía tipear cualquier correo sin contraseña).

---

### PASO 3 · Redirección Inteligente por Rol, Sesión JWT y Blindaje de Endpoints
- **Rama Backend:** `fix/paso-03-middleware-jwt` | **Rama Android:** `fix/paso-03-navegacion-roles`
- **Objetivo de Calidad:** Que cada usuario entre a la pantalla que corresponde a su rol (Vecino vs Profesional nuevo vs Profesional verificado) y que ningún endpoint privado funcione sin JWT válido.
- [ ] **`[BE-3.1]`** En `internal/handlers/api.go`, cambiar `withOptionalAuth` por un middleware estricto `withRequiredAuth` (devuelve 401 Unauthorized si falta o expiró el Bearer token) en: `GET /api/v1/requests/my`, `POST /api/v1/requests`, `POST /api/v1/reviews`, `PATCH /api/v1/users/me/address`, `POST /api/v1/professionals/{id}/verification` y `PATCH /api/v1/professionals/{id}/availability`.
- [ ] **`[BE-3.2]`** Corregir la fuga de datos en `GET /api/v1/requests/my` (BUG-03) eliminando la condición SQL (`$1::text = ''`) y validar en los endpoints de profesional que el usuario autenticado sea dueño de ese `professional_id` (BUG-07).
- [ ] **`[AND-3.3]`** En `MainActivity.kt` (corrige BUG-01), unificar la función `navigateAfterAuth(user)`:
  - (a) si `user.role == 'client'` -> ir a `ClientHomeCategories`;
  - (b) si `user.role == 'professional'` recién registrado o sin DNI cargado -> ir a `ProfessionalVerificationOnboarding` (Vista #13);
  - (c) si `user.role == 'professional'` ya dado de alta -> ir directo a `ProfessionalDashboard` (Vista #14). Hacer lo mismo al restaurar sesión desde `SessionDataStore` al abrir la app.

---

### PASO 4 · Barra Superior, Barra Inferior y Pantalla 'Mi Cuenta' Diferenciada por Rol
- **Rama Backend:** `feat/paso-04-perfil-usuario` | **Rama Android:** `fix/paso-04-barras-y-cuenta`
- **Objetivo de Calidad:** Corregir los botones repetidos del TopAppBar, arreglar las pestañas del BottomNavBar y separar qué ve un Vecino vs un Profesional en 'Mi Cuenta'.
- [ ] **`[AND-4.1]`** En `StitchTopAppBar` (`StitchNavigation.kt` y `HomeScreen.kt`), diferenciar los botones: el escudo/logo a la izquierda lleva siempre al Inicio del rol actual, y el botón derecho 'Mi Cuenta' abre `UserAccountScreen`.
- [ ] **`[AND-4.2]`** En `StitchBottomNavBar` (`MainActivity.kt` L. 108-127):
  - (a) Pestaña 'Inicio' -> lleva a `ClientHomeCategories` (si es vecino) o `ProfessionalDashboard` (si es profesional);
  - (b) Pestaña 'Mis contactos' -> abre la nueva pantalla `MyRequestsHistoryScreen` conectada a `GET /api/v1/requests/my` (corrige BUG-02);
  - (c) Pestaña 'Ayuda' -> abre la pantalla de Ayuda Vecinal y Garantías con botón de llamado directo al 0800-444-OFICIO.
- [ ] **`[AND-4.3]`** En `UserAccountScreen`, mostrar los datos reales del usuario logueado y botones acordes a su rol: si es Vecino, cambiar dirección/barrio y ver su historial de pedidos; si es Profesional, ir a su Panel de Solicitudes o actualizar su legajo de verificación DNI/Matrícula; más el botón rojo de Cerrar Sesión que limpia DataStore y vuelve al inicio.

---

### PASO 5 · Home del Vecino, Diagnóstico Inteligente y Semilla de los 8 Rubros (Vista #4)
- **Rama Backend:** `feat/paso-05-seed-y-coordenadas` | **Rama Android:** `feat/paso-05-home-vecino`
- **Objetivo de Calidad:** Dejar el Home del Vecino 100% funcional, con coordenadas reales por barrio, diagnóstico sin textos pre-escritos y los 8 rubros poblados en PostgreSQL.
- [ ] **`[BE-5.1]`** En `migrations/002_seed.sql`, cargar al menos 16 profesionales verificados distribuidos en los 8 rubros (Plomero, Electricista, Gasista, Aire acondicionado, Pintor, Albañil, Jardinero y Cerrajero) con fotos de perfil, portfolios de trabajos con imágenes reales, matrículas y coordenadas PostGIS en Caballito, Flores, Almagro y Parque Chacabuco.
- [ ] **`[BE-5.2]`** En `UpdateUserAddress` (`store.go`), agregar un diccionario geocodificador de barrios de CABA/GBA para que al cambiar de dirección actualice `saved_lat` y `saved_lng` reales del barrio correspondiente.
- [ ] **`[AND-5.3]`** En `HomeScreen.kt`:
  - (a) vaciar el texto hardcodeado 'No enciende el calefón de la cocina' del diálogo de diagnóstico (L. 80) y agregar chips de ejemplos rápidos tocables ('Pierde agua el flexible', 'Saltó la térmica', 'Olor a gas', 'Se trabó la cerradura');
  - (b) cuando el backend responde la categoría sugerida, mostrar el botón 'Ver [Rubro]s cerca de mi casa ->' que navega directo al listado filtrado.

---

### PASO 6 · Pantalla de Resultados PostGIS 1:1 con Stitch (Vista #5)
- **Rama Backend:** `feat/paso-06-filtros-postgis` | **Rama Android:** `feat/paso-06-resultados-stitch`
- **Objetivo de Calidad:** Maquetar la pantalla `cliente_resultados_de_profesionales` idéntica al HTML de Stitch, con fotos de perfil, ambos sellos simultáneos, scroll de filtros y estado vacío.
- [ ] **`[BE-6.1]`** En `ListProfessionals` (`store.go`), devolver junto con la lista los contadores o soportar los 4 filtros combinables: Todos, Atienden urgencias (`only_urgent=true`), Matriculados (`only_licensed=true`) y Mejor valorados (`top_rated=true`).
- [ ] **`[AND-6.2]`** En `SearchResultsScreen.kt` (comparando con `cliente_resultados_de_profesionales/code.html`):
  - (a) agregar barra superior sticky con botón '<- Volver' y sello 'Guía de Oficios';
  - (b) convertir la fila de chips en `LazyRow` con los 4 filtros de Stitch incluyendo el contador 'Todos (N)';
  - (c) agregar la tira azul de seguridad 'Profesionales de tu barrio' para adultos mayores.
- [ ] **`[AND-6.3]`** En cada tarjeta de profesional:
  - (a) renderizar la foto de 64dp con `AsyncImage` (`pro.avatar_url`) y su insignia verde/naranja en la esquina;
  - (b) corregir el bug del 'else if' (L. 1445) usando dos 'if' independientes en un `FlowRow` para que un profesional matriculado que también atiende urgencias muestre AMBOS sellos ('[OK] Matrícula Verificada' y '[!] Atiende urgencias hoy');
  - (c) agregar componente de Estado Vacío con botón 'Limpiar filtros' cuando no hay coincidencias.

---

### PASO 7 · Ficha Completa del Profesional y Galería de Fotos Antes/Después (Vista #6)
- **Rama Backend:** `fix/paso-07-distancia-perfil` | **Rama Android:** `feat/paso-07-perfil-profesional`
- **Objetivo de Calidad:** Completar cada bloque de `cliente_perfil_del_profesional_y_contacto` sin datos hardcodeados y renderizando las imágenes reales del portfolio con Coil.
- [ ] **`[BE-7.1]`** En `GET /api/v1/professionals/{id}` (`api.go` y `store.go`), recibir query params opcionales `lat` y `lng` del vecino para calcular con `ST_Distance` la distancia exacta en cuadras y metros en lugar del texto fijo 'A 8 cuadras (650 metros)' (corrige BUG-08).
- [ ] **`[AND-7.2]`** En `MainActivity.kt` (L. 337-344), eliminar por completo el DTO fallback hardcodeado de Roberto Gómez y mostrar un indicador de carga limpio mientras responde la API.
- [ ] **`[AND-7.3]`** En `ProfessionalProfileScreen.kt`:
  - (a) agregar barra superior '<- Volver a resultados' con pill del barrio;
  - (b) mostrar foto de 80dp con borde verde y tilde de verificado en la tarjeta principal;
  - (c) construir el bloque 'SEGURIDAD Y GARANTÍA VECINAL' con sus 3 tarjetas completas (DNI verificado, Matrícula oficial + ente, Disponible hoy para urgencias).
- [ ] **`[AND-7.4]`** En la sección 'Fotos de trabajos realizados' (corrige BUG-05): renderizar cada ítem de `professional.portfolio` con `AsyncImage(model = item.image_url)` de 192dp de alto, badge verde 'Verificado' sobre la foto, título y subtítulo. En 'Opiniones de vecinos', dibujar las 5 estrellas doradas, los chips de highlights y el pie '[OK] Trabajo contratado por Guía de Oficios'.
- [ ] **`[AND-7.5]`** En la barra fija inferior, al tocar 'Escribir por WhatsApp' o 'Llamar por teléfono', pasar el canal preseleccionado ('whatsapp' o 'phone') a la pantalla de Solicitud Guiada.

---

### PASO 8 · Pantalla Comunidad: Top Profesionales del Mes / Podio Barrial (Vista #7)
- **Rama Backend:** `feat/paso-08-endpoint-ranking` | **Rama Android:** `feat/paso-08-podio-barrial`
- **Objetivo de Calidad:** Construir desde cero la Vista #7 de Stitch (`comunidad_top_profesionales_del_mes`) y conectarla al banner del Home.
- [ ] **`[BE-8.1]`** Verificar que `GET /api/v1/professionals?top_rated=true&category_id=...` ordene por `rating_avg DESC, reviews_count DESC, jobs_completed_count DESC` y devuelva también el comentario destacado más reciente de un vecino para mostrar en el podio.
- [ ] **`[AND-8.2]`** Crear `TopProfessionalsMonthScreen.kt` siguiendo `comunidad_top_profesionales_del_mes/code.html`:
  - (a) Header con selector de barrio;
  - (b) Hero Banner azul 'Top Profesionales del Mes — Reconocimiento a los trabajadores mejor calificados';
  - (c) Chips horizontales de filtro por oficio;
  - (d) Bento Podium con la tarjeta #1 Puesto Oro a ancho completo (borde ámbar, medalla #1, foto 80dp, cita de vecino, precio de visita y botón 'Ver perfil completo y contactar') y las tarjetas #2 Plata y #3 Bronce en grilla de 2 columnas;
  - (e) Lista de puestos #4 al #10 y tarjeta explicativa '¿Cómo se calcula este ranking?'.
- [ ] **`[AND-8.3]`** En `MainActivity.kt` (L. 269), conectar el botón 'Ver ranking completo' del Home para que navegue a `StitchRoute.CommunityTopProfessionals` en vez de `ClientSearchResults(0)`.

---

### PASO 9 · Solicitud Guiada (Paso 1 de 2) 1:1 con Stitch y Foto con Preview (Vista #9)
- **Rama Backend:** `feat/paso-09-solicitud-backend` | **Rama Android:** `feat/paso-09-solicitud-guiada`
- **Objetivo de Calidad:** Rehacer la pantalla `cliente_solicitud_guiada_a_roberto` con las tarjetas grandes de urgencia, previsualización real de foto y modal de confirmación.
- [ ] **`[BE-9.1]`** En `CreateServiceRequest` (`store.go` L. 564), eliminar el fallback 'María Elena Castro' y obtener siempre `full_name`, `phone` y `address` reales del cliente autenticado en el JWT.
- [ ] **`[AND-9.2]`** En `GuidedRequestScreen.kt` (`cliente_solicitud_guiada_a_roberto/code.html`):
  - (a) agregar barra superior 'Solicitud directa · Paso 1 de 2';
  - (b) agregar tarjeta de contexto del profesional arriba con su foto, nombre, oficio y barrio;
  - (c) reemplazar los FilterChips de urgencia por las 2 tarjetas Radio grandes de 72dp ('Normal · Sin apuro' vs 'Urgente hoy · Atención rápida' con punto rojo pulsante).
- [ ] **`[AND-9.3]`** En el Paso 3 ('Sumá una foto si querés'): cuando el vecino selecciona una foto con `PickVisualMedia` y se sube a `POST /api/v1/uploads`, mostrar la miniatura real con `AsyncImage`, el nombre/estado de la imagen y un botón 'Quitar / Cambiar foto'.
- [ ] **`[AND-9.4]`** Agregar el bloque verde 'Tu tranquilidad es prioridad' y el botón fijo inferior 'Continuar para elegir contacto ->' con validación en línea si la descripción del problema está vacía.

---

### PASO 10 · Confirmación de Solicitud (Paso 2 de 2) y Disparo a WhatsApp / Teléfono (Vista #10)
- **Rama Backend:** `feat/paso-10-confirmacion-pedido` | **Rama Android:** `feat/paso-10-confirmacion-ui`
- **Objetivo de Calidad:** Pulir `cliente_elegir_contacto_y_confirmaci_n` sin datos hardcodeados y conectarla tanto a WhatsApp/Llamada como al Mapa en Vivo.
- [ ] **`[AND-10.1]`** En `MainActivity.kt` (L. 391-400), eliminar el fallback hardcodeado `'req-demo-1001'` de Roberto Gómez.
- [ ] **`[AND-10.2]`** En `RequestConfirmationScreen.kt` (`cliente_elegir_contacto_y_confirmaci_n/code.html`):
  - (a) agregar círculo verde de 64dp con tilde blanco y título '¡Listo! Ya preparamos tu pedido para [Nombre]';
  - (b) caja de resumen con el problema real, etiqueta roja/azul de urgencia y dirección del vecino;
  - (c) botón verde gigante (64dp) 'Escribir por WhatsApp — Se envía tu problema escrito automáticamente' y botón azul (64dp) 'Llamar por teléfono — Hablar directo al [Número]'.
- [ ] **`[AND-10.3]`** Agregar el botón 'Seguir estado del pedido / Ver mapa en vivo ->' que lleva a la Vista #11, y el pie de garantía vecinal con el 0800-444-OFICIO.

---

### PASO 11 · Seguimiento en Vivo con Mapa y Pestaña 'Mis Contactos' (Vista #11)
- **Rama Backend:** `feat/paso-11-status-en-camino` | **Rama Android:** `feat/paso-11-mapa-en-vivo`
- **Objetivo de Calidad:** Crear el endpoint de transición de estados, la pantalla de historial 'Mis Contactos' y la Vista #11 de Profesional en Camino con Mapa.
- [ ] **`[BE-11.1]`** Crear endpoint `PATCH /api/v1/requests/{id}/status` en `api.go` y `store.go` que permita cambiar el estado de una solicitud ('pending' -> 'en_route' -> 'completed' o 'cancelled') guardando `eta_minutes` y validando que quien lo cambia sea el profesional asignado o el cliente dueño.
- [ ] **`[AND-11.2]`** Crear `MyRequestsHistoryScreen.kt` para la pestaña 'Mis contactos': lista todas las solicitudes del vecino desde `GET /api/v1/requests/my` mostrando estado ('Pendiente', 'En camino', 'Finalizado') y permitiendo abrir el Mapa en Vivo (si está en curso) o Calificar (si terminó).
- [ ] **`[AND-11.3]`** Crear `ProfessionalEnRouteMapScreen.kt` siguiendo `cliente_profesional_en_camino_y_mapa_en_vivo/code.html` (Vista #11):
  - (a) tarjeta verde de estado '¡[Nombre] aceptó tu solicitud y está en camino!' con caja azul de ETA ('En aprox. 12 minutos');
  - (b) mapa barrial con trazado de ruta, pin pulsante del profesional y pin del domicilio del vecino;
  - (c) tarjeta del profesional con botones directos 'Llamar ahora' y 'WhatsApp';
  - (d) checklist 'Para tu tranquilidad' y botón 'El trabajo ya terminó -> Calificar atención' que navega a la Vista #12.

---

### PASO 12 · Trabajo Terminado y Calificación Transaccional 1 a 5 Estrellas (Vista #12)
- **Rama Backend:** `feat/paso-12-reviews-tx` | **Rama Android:** `feat/paso-12-calificacion-stitch`
- **Objetivo de Calidad:** Rehacer `cliente_trabajo_terminado_y_calificaci_n` con estrellas grandes de 48dp, tarjetas completas de aspectos destacables y comentario limpio.
- [ ] **`[BE-12.1]`** En `CreateReview` (`store.go`), verificar que al insertar la reseña en la transacción SQL también marque el `service_request` como `'completed'` y recalcule `rating_avg`, `reviews_count` y `jobs_completed_count` del profesional.
- [ ] **`[AND-12.2]`** En `JobCompletedReviewScreen.kt` (`cliente_trabajo_terminado_y_calificaci_n/code.html`):
  - (a) agregar arriba la tarjeta resumen del profesional con su foto, nombre, rubro, descripción del trabajo realizado y fecha 'Finalizado hoy';
  - (b) reemplazar los FilterChips por 5 botones de estrella Star de 48dp con píldora dinámica ('5.0 · ¡Excelente trabajo!');
  - (c) maquetar las 5 tarjetas de aspectos destacables con su ícono, título y subtítulo explicativo.
- [ ] **`[AND-12.3]`** Eliminar el comentario pre-escrito sobre Roberto Gómez (L. 1967, dejar `comment = ""` con placeholder claro), agregar la tarjeta 'Transparencia Vecinal' y al enviar la reseña navegar de regreso mostrando el nuevo promedio actualizado.

---

### PASO 13 · Onboarding del Profesional: Rubro, Especialidades y Fotos Reales de DNI/Matrícula (Vista #13)
- **Rama Backend:** `feat/paso-13-verificacion-pro` | **Rama Android:** `feat/paso-13-onboarding-dni`
- **Objetivo de Calidad:** Rehacer `profesional_registro_y_verificaci_n_dni` sin datos pre-escritos, permitiendo elegir cualquiera de los 8 rubros y exigiendo fotos reales del DNI.
- [ ] **`[AND-13.1]`** En `ProfessionalVerificationScreen.kt` (`profesional_registro_y_verificaci_n_dni/code.html`):
  - (a) agregar barra de progreso superior 'Paso 2 de 3 · Verificación de oficio e identidad · 66%';
  - (b) agregar selector desplegable de los 8 rubros oficiales conectado a `categories` para enviar el `primary_category_id` real elegido (en vez de mandar 1 fijo);
  - (c) agregar chips interactivos de especialidades secundarias con botón '+ Agregar otra' y selector de radio de cobertura en km.
- [ ] **`[AND-13.2]`** En el Paso 2 ('Fotos de tu DNI'): maquetar las 2 tarjetas en grilla ('Frente del DNI' y 'Dorso del DNI') mostrando la miniatura `AsyncImage` cuando el archivo se sube a `POST /api/v1/uploads`. Eliminar por completo las URLs falsas hardcodeadas (L. 2208-2212) y validar que no se pueda enviar el formulario sin haber subido Frente y Dorso reales.
- [ ] **`[AND-13.3]`** Vaciar los textos pre-cargados ('18.492', 'ENARGAS') en el Paso 3 de Matrícula opcional y, al guardar con éxito, redirigir al profesional a su Panel de Solicitudes (Vista #14).

---

### PASO 14 · Panel del Profesional, Vista 'Aceptar Trabajo e Iniciar Viaje' y APK Final (Vistas #14 y #15)
- **Rama Backend:** `feat/paso-14-panel-y-viaje` | **Rama Android:** `feat/paso-14-aceptar-trabajo`
- **Objetivo de Calidad:** Completar el Panel Profesional, construir desde cero la Vista #15 para aceptar pedidos e iniciar viaje, y realizar la prueba cruzada completa.
- [ ] **`[AND-14.1]`** En `ProfessionalDashboardScreen.kt` (Vista #14):
  - (a) quitar el fallback hardcodeado de `MainActivity.kt` (L. 475-489);
  - (b) agregar la tarjeta superior con foto del profesional, sello verificado y oficio;
  - (c) en cada tarjeta de solicitud recibida, mostrar avatar con iniciales del vecino, etiqueta Urgente/Normal y hacer que al tocar la tarjeta o el botón principal se abra `ProfessionalAcceptJobScreen` (Vista #15);
  - (d) conectar el botón inferior 'Editar mi zona de trabajo y precios'.
- [ ] **`[AND-14.2]`** Crear `ProfessionalAcceptJobScreen.kt` siguiendo `profesional_aceptar_trabajo_e_iniciar_viaje/code.html` (Vista #15):
  - (a) Header con número de solicitud y badge 'Urgente';
  - (b) tarjeta del vecino con distancia y caja de Dirección Exacta;
  - (c) bloque 'Detalle del problema' con la descripción y la foto real subida por el vecino (`AsyncImage` con modal de zoom al tocar 'Ver foto');
  - (d) caja de 'Visita técnica acordada';
  - (e) Switch de ubicación GPS en tiempo real;
  - (f) botón verde 'Iniciar viaje hacia el domicilio ->' que invoca `PATCH /api/v1/requests/{id}/status` (`estado='en_route'`), botón 'Abrir dirección en Google Maps / Waze', botón 'Avisar demora por WhatsApp' y botón 'Rechazar solicitud'.
- [ ] **`[QA-14.3]`** **Prueba Cruzada E2E entre los 2 integrantes y armado de APK:**
  - Usuario 1 (Vecino) se registra, busca un rubro, filtra, ve el perfil con fotos reales, crea solicitud con foto ->
  - Usuario 2 (Profesional) la recibe en su Panel, abre la Vista #15, ve la foto, toca 'Iniciar viaje' ->
  - Usuario 1 ve en la Vista #11 que el profesional está en camino con ETA ->
  - al finalizar califica con 5 estrellas en la Vista #12 y el promedio sube en PostgreSQL.
  - Finalmente correr `.\gradlew.bat assembleDebug` para generar `app-debug.apk`.
