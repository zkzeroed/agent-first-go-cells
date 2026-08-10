# 2. Core Design Principles

### Principle 1: Radical Locality

Everything needed to understand and modify a feature lives in one cell directory. An agent opens `internal/cells/user-authenticate/` and never needs to look elsewhere for that capability.

### Principle 2: Aggressive Smallness

Files and functions must be small, enforced by tooling:

- **Files:** 300 LOC maximum.
- **Functions:** 40 LOC maximum, unless the function documents `AGENT_OVERRIDE`.
- AST-based structural tests enforce both limits. `golangci-lint` enforces
  cyclomatic complexity with `cyclop` (maximum 16).

### Principle 3: Predictable Isomorphism

Every cell follows the same file schema. `api/api.go` is the public contract,
`{name}.go` is private implementation, `service.go` is business logic, and
`errors.go` is the cell-local taxonomy.

### Principle 4: Interfaces as Membranes

Every cell exposes behavior through its `api` package. Public interfaces,
shared contract types, and public errors live there; other cells may import
only `internal/cells/<id>/api`.

### Principle 5: Grep-Friendliness & Redundant Naming

Behavior-first kebab-case cell IDs (`user-authenticate`, `payment-capture`). Cell ID repeated in directory name, file names, interface names, and log keys. Split files embed the cell ID (`user_authenticate_login.go`). Go package name drops hyphens (`package userauthenticate`).

### Principle 6: Explicit Wiring, No Magic

All cell construction and connection happens in `internal/app/wiring.go`. No `init()` functions. No reflection DI. No service locators. The entire dependency graph is visible in one file.

### Principle 7: Self-Documenting for Machines

`AGENTS.md` at root and per cell. `doc.go` per cell with cell ID anchor. `AGENTS.md` includes task-oriented guidance: purpose, start-here files, common tasks, validation commands, reliability & concurrency notes.

### Principle 8: Machine-Readable Metadata

Lightweight `cell.yaml` manifest per cell (5 mandatory fields + optional `invariants`). Generated `gen/cells.json` index for project-wide orientation. Generated `gen/context/{id}.context.md` packs for per-cell quick orientation.

### Principle 9: Impact Visibility

`task impact` maps changed files to owning cells and affected downstream cells. Reduces under-testing and accidental cross-cell breakage.

### Principle 10: Guardrails over Guidelines

Strict `.golangci.yml` with merged linter set. AST-based structural tests enforce architecture invariants. Import policy as code (`policy/imports.yaml`). Manifest validation in CI. Index staleness checks via deterministic hash.

### Principle 11: Pragmatic Guardrails

Keep the executable guardrails that provide evidence for a change. Remove a
guardrail only with its documentation and validation together; unsupported
configuration tiers create ambiguity for agents.
