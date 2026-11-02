# Guía de Diseño UX/UI y Accesibilidad — "Guía de Oficios"

> [!IMPORTANT]
> **Objetivo central de UX:** Que cualquier persona —incluyendo adultos mayores de 65+ años o personas con poca experiencia usando celulares— pueda encontrar y contactar a un trabajador de confianza **sin miedo a tocar algo equivocado, sin leer letra chica y sin perderse entre pantallas**.

---

## 1. Paleta de Colores Armoniosa y Serena

Descartamos el uso de amarillos estridentes o fondos negros duros. Utilizamos una paleta inspirada en aplicaciones médicas, bancarias simplificadas y de bienestar: **tonos azules serenos, verdes amables y fondos cálidos de alto descanso visual**.

| Rol Visual | Nombre del Tono | Código Hex | Uso en la Aplicación |
| :--- | :--- | :--- | :--- |
| **Fondo Principal** | Blanco Alabastro / Suave | `#F8FAFC` | Fondo general de todas las pantallas (evita el encandilamiento del blanco puro `#FFFFFF`). |
| **Superficie de Tarjetas** | Blanco Puro | `#FFFFFF` | Tarjetas de oficios y fichas de profesionales, con bordes suaves `#E2E8F0`. |
| **Primario (Marca y Confianza)** | Azul Sereno Profundo | `#1E40AF` / `#2563EB` | Encabezados suaves, botones principales de navegación, íconos de rubros y enlaces. |
| **Secundario (Calidez / Cercanía)** | Azul Cielo Claro | `#EFF6FF` | Fondos de etiquetas informativas, pasos guiados y recuadros de ayuda. |
| **Acción Principal (Contacto)** | Verde Esmeralda Amable | `#15803D` | Botón principal de **"Hablar por WhatsApp"** o **"Confirmar Solicitud"**. |
| **Sellos de Confianza** | Verde Agua Suave | `#DCFCE7` (Fondo)<br>`#166534` (Texto) | Etiquetas de *"✓ DNI Verificado"* y *"✓ Matriculado"*. |
| **Destacado / Estrellas** | Ámbar Cálido Suave | `#D97706` | Estrellas de calificación (`★ 4.9`) y etiqueta suave de profesional destacado, sin saturar. |
| **Texto Principal** | Gris Pizarra Oscuro | `#0F172A` | Títulos y nombres de profesionales (máxima legibilidad). |
| **Texto Secundario** | Gris Medio Legible | `#475569` | Descripciones y distancias (nunca usamos grises demasiado claros que se pierdan). |

---

## 2. Reglas de Oro de Accesibilidad y Didáctica

### A. Grilla de Oficios de 2 Columnas con Explicación
En lugar de poner 4 íconos pequeños por fila donde apenas entra la palabra "Aire Acond.", usamos **2 tarjetas grandes por fila** que incluyen:
1. Un ícono claro sobre un círculo de color suave.
2. El nombre del oficio en negrita grande (**Plomero**).
3. Ejemplos cotidianos debajo (*"Canillas, caños, tanque, pérdidas"*), porque muchas veces una persona mayor no sabe si el calefón lo arregla el plomero o el gasista.

### B. Ayuda Contextual ("¿No sabés qué rubro elegir?")
Arriba de la lista de oficios hay un cuadro didáctico muy claro: podés escribir con tus palabras qué te pasa (*"No sale agua caliente"*) y la aplicación te muestra directamente al profesional indicado.

### C. Cero Jerga Técnica o Iconos sin Texto
* **Mal (Apps tradicionales):** Un ícono de un escudo solo o una campanita sin etiqueta.
* **Bien (Guía de Oficios):** Todos los íconos van acompañados de su palabra en castellano claro: `🏠 Inicio`, `📋 Mis contactos`, `❓ Ayuda`.

### D. Contacto Dual: WhatsApp + Llamada Telefónica
Mientras que el público joven prefiere solo chatear, la gente grande muchas veces prefiere **hablar por teléfono** o mandar un **audio de WhatsApp**. Por eso, el paso de contacto ofrece ambos botones bien grandes y diferenciados:
* 🟢 **Hablar por WhatsApp** (Abre el chat listo para mandar mensaje o audio).
* 🔵 **Llamar por Teléfono** (Abre el marcador del celular con el número cargado).

---

## 3. Anatomía de las Vistas en Jetpack Compose (Material 3)

Para trasladar este diseño a **Android Studio**, definimos en `android-app/` el archivo `Theme.kt` con:
* **Tipografía escalada:** Respetando el tamaño de fuente configurado en el teléfono Android del usuario (`sp`).
* **Esquinas redondeadas amables:** `RoundedCornerShape(16.dp)` en tarjetas y `RoundedCornerShape(12.dp)` en botones para dar una sensación orgánica y moderna.
* **Elevación sutil:** Sombras suaves de `2.dp` con borde de `1.dp` para que cada botón se distinga claramente del fondo como un elemento presionable.
