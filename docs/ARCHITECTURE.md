# netdbg – Architecture & AI Prompt

This document describes the architecture, design patterns, and modularity conventions of the netdbg project. It is intended to help an AI (or any developer) understand how to keep the codebase clean, modular, and consistent when adding new features.

---

## 1. Design Philosophy

- **Modularity:** Each core feature must live in its own module under `internal/`, with a clear separation between business logic, helpers, options, and tests.
- **Command/Executor Pattern:** Every module follows this pattern:
  - `command.go`: Cobra integration and flag parsing.
  - `executor.go`: Executor interface and core logic.
  - `options.go`: Options struct and flag helpers.
  - Additional helpers (e.g., `check.go`, `build.go`, `kubectl.go`, etc.) for supporting logic.
- **Dependency Injection:** All functions interacting with the system (`exec.Command`, `http.Get`, etc.) must use injectable global variables to simplify testing.
- **Centralized Logging:** All logging must go through `internal/logger` for consistency and traceability.
- **Comprehensive Testing:** Every helper and main workflow must include unit and integration tests covering both success and failure paths, using the injection pattern and `setupLogger()`.

---

## 2. Relevant Folder Structure

- `cmd/`
  Main CLI commands (`netcat.go`, `revdns.go`, `kexec.go`, `root.go`, etc.).
  Each command should only handle wiring and delegate execution to its corresponding module in `internal/`.

- `internal/netcat/`
  - Logic and helpers for the netcat command.
  - Tests: `*_test.go` files for each module.

- `internal/revdns/`
  - Logic and helpers for the revdns command.
  - Tests: `*_test.go` files for each module.

- `internal/kexec/`
  Logic and helpers for the kexec command (executing netdbg inside Kubernetes pods).
  - `command.go`
  - `executor.go`
  - `options.go`
  - `check.go`
  - `build.go`
  - `download.go`
  - `kubectl.go`
  - Tests: `*_test.go` files for each module.

- `internal/logger/`
  Centralized logging system, initialized in `main.go` and in every test through `setupLogger()`.

- `main.go`
  Initializes the logger and runs the CLI.

---

## 3. Conventions for New Features

- Create a folder under `internal/` for each new feature (e.g., `internal/traceroute/`).
- Follow the command/executor/options/helpers/tests pattern.
- Use injectable variables for external dependencies.
- Add comprehensive tests and use `setupLogger()` in all tests.
- Integrate the command into `cmd/` only as wiring.
- Use the centralized logger for all messages.
- If the feature requires external resources (network, files, etc.), ensure tests can mock them through dependency injection.

---

## 4. Example Structure for a New Feature

```text
internal/traceroute/
  command.go      // ExecuteCommand for Cobra integration
  executor.go     // Executor interface and DefaultExecutor
  options.go      // Options struct and flag helpers
  helpers.go      // Supporting logic (e.g., ICMP, TCP, etc.)
  traceroute.go   // Core traceroute logic
  traceroute_test.go
````

---

## 5. Guidelines for the AI/Developer

* Before adding code, review the modular structure and existing patterns.
* Maintain separation of concerns: flag parsing in `command.go`, logic in `executor.go`, helpers in separate files.
* Use the logging system and dependency injection pattern to simplify testing.
* Add unit and integration tests for every helper and main workflow.
* Document any new helper or convention in this file.
* If in doubt, use the existing modules (`netcat`, `revdns`, `kexec`) as references.

---

## 6. Important Folders to Read Source Code

* `cmd/`
* `internal/netcat/`
* `internal/revdns/`
* `internal/kexec/`
* `internal/logger/`
* `main.go`

---

## 7. Example AI Prompt

> Read the architecture and design patterns described in ARCHITECTURE.md.
> Add feature X following the modular pattern (command/executor/options/helpers/tests), using centralized logging and injectable variables for external dependencies.
> Ensure tests cover all execution paths and use `setupLogger()` in every test.
> Integrate the command into `cmd/` only as wiring.
```
