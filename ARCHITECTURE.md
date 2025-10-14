# netdbg – Project Architecture & AI Prompt

Este documento describe la arquitectura, los patrones de diseño y las convenciones de modularidad del proyecto netdbg. Está pensado para que una IA (o cualquier desarrollador) pueda entender cómo mantener el código limpio, modular y coherente al añadir nuevas funcionalidades.

---

## 1. Filosofía de diseño

- **Modularidad:** Cada funcionalidad principal debe estar en su propio módulo bajo `internal/`, con separación clara entre lógica de negocio, helpers, opciones y tests.
- **Patrón Command/Executor:** Cada módulo sigue el patrón:
  - `command.go`: integración con Cobra y parsing de flags.
  - `executor.go`: interfaz Executor y lógica principal.
  - `options.go`: struct Options y helpers de flags.
  - Helpers adicionales (por ejemplo, `check.go`, `build.go`, `kubectl.go`, etc.) para lógica auxiliar.
- **Inyección de dependencias:** Todas las funciones que interactúan con el sistema (exec.Command, http.Get, etc.) usan variables globales inyectables para facilitar el testeo.
- **Logging centralizado:** Todo el logging debe hacerse a través de `internal/logger` para trazabilidad y consistencia.
- **Tests exhaustivos:** Cada helper y flujo principal debe tener tests unitarios y de integración, cubriendo caminos de éxito y error, usando el patrón de inyección y setupLogger().

---

## 2. Estructura de carpetas relevante

- `cmd/`
  Comandos principales de la CLI (netcat.go, revdns.go, kexec.go, root.go, etc.).
  Cada comando solo debe hacer wiring y delegar en su módulo correspondiente de `internal/`.

- `internal/netcat/`
  - Lógica y helpers para el comando netcat.
  - Tests: `*_test.go` en cada módulo.

- `internal/revdns/`
  - Lógica y helpers para el comando revdns.
  - Tests: `*_test.go` en cada módulo.

- `internal/kexec/`
  Lógica y helpers para el comando kexec (ejecución de netdbg en pods de Kubernetes).
  - `command.go`
  - `executor.go`
  - `options.go`
  - `check.go`
  - `build.go`
  - `download.go`
  - `kubectl.go`
  - Tests: `*_test.go` en cada módulo.

- `internal/logger/`
  Sistema de logging centralizado, inicializado en main.go y en cada test con setupLogger().

- `main.go`
  Inicializa el logger y ejecuta la CLI.

---

## 3. Convenciones para nuevas funcionalidades

- Crea una carpeta bajo `internal/` para cada nueva funcionalidad (por ejemplo, `internal/traceroute/`).
- Sigue el patrón command/executor/options/helpers/tests.
- Usa variables inyectables para dependencias externas.
- Añade tests exhaustivos y setupLogger() en todos los tests.
- Integra el comando en `cmd/` solo como wiring.
- Usa el logger centralizado para todos los mensajes.
- Si la funcionalidad requiere recursos externos (red, archivos, etc.), asegúrate de que los tests puedan simularlos mediante inyección.

---

## 4. Ejemplo de estructura para una nueva funcionalidad

```
internal/traceroute/
  command.go      // ExecuteCommand para integración con Cobra
  executor.go     // Executor interface y DefaultExecutor
  options.go      // Struct Options y helpers de flags
  helpers.go      // Lógica auxiliar (por ejemplo, ICMP, TCP, etc.)
  traceroute.go   // Lógica principal de traceroute
  traceroute_test.go
```

En `cmd/traceroute.go`:
```go
import "github.com/feliux/netdbg/internal/traceroute"

var tracerouteCmd = &cobra.Command{
    Use:   "traceroute",
    Short: "Network path tracing tool",
    Run: func(cmd *cobra.Command, args []string) {
        traceroute.ExecuteCommand(cmd, args)
    },
}
```

---

## 5. Pautas para la IA/desarrollador

- Antes de añadir código, revisa la estructura modular y los patrones existentes.
- Mantén la separación de responsabilidades: parsing de flags en command.go, lógica en executor.go, helpers en archivos separados.
- Usa el sistema de logging y el patrón de inyección para facilitar los tests.
- Añade tests unitarios y de integración para cada helper y flujo principal.
- Documenta cualquier helper o convención nueva en este archivo.
- Si tienes dudas, revisa los módulos existentes (netcat, revdns, kexec) como referencia.

---

## 6. Carpetas importantes para leer ficheros de código

- `cmd/`
- `internal/netcat/`
- `internal/revdns/`
- `internal/kexec/`
- `internal/logger/`
- `main.go`

---

## 7. Ejemplo de prompt para la IA

> Lee la arquitectura y patrones de diseño en ARCHITECTURE.md.
> Añade la funcionalidad X siguiendo el patrón modular (command/executor/options/helpers/tests), usando logging centralizado y variables inyectables para dependencias externas.
> Asegúrate de que los tests cubran todos los caminos y usa setupLogger() en cada test.
> Integra el comando en cmd/ solo como wiring.
