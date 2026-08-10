#

This file is the first thing an AI agent should read when arriving at this
codebase. It provides orientation, operating instructions, and validation
commands.

For full architecture documentation, see docs/architecture/README.md

## Purpose

This is a bootstrap codebase for the Agent-First Go architecture — a Go project
structure designed from first principles for AI coding agents as primary
maintainers. It provides the tooling, guardrails, and scaffolding to quickly
start new projects using the cell-based architecture.

## Start Here

1. **Read the architecture overview:** `docs/architecture/README.md`
2. **List existing cells:** `task cells`
3. **Check architecture health:** `task doctor`
4. **Create a new cell:** `task new-cell ID=<behavior-name>`
5. **Find a cell by concept:** `task find-cell QUERY=<text>`

## architecture Summary

- **Cells** are vertical slices — one directory per capability (`internal/cells/<name>/`)
- **Fixed file schema** per cell: `cell.yaml`, `AGENTS.md`, `doc.go`, `<name>.go`, `types.go`, `errors.go`, `service.go`, `store.go`, `handler.go`, `*_test.go`
- **API packages as membranes** — cells expose behavior through `api` packages only
- **Explicit wiring** — all DI in `internal/app/wiring.go` (< 300 LOC, no `init()`)
- **Guardrails enforced by tooling** — linters, structural tests, import policy

## Build Convention: Inside-Out

Cells are built **inside-out** — from domain vocabulary to transport layer.
This mirrors hexagonal architecture's "domain outwards to adapters" but with
cell-specific ordering:

1. **`types.go`** — domain types (the vocabulary everything else uses)
2. **`errors.go`** — error taxonomy (sentinel errors + typed errors)
3. **`api/api.go`** — public interfaces, shared contract types, and public errors (the membrane)
4. **`<name>.go`** — private implementation + `New()` constructor
5. **`store.go`** — data access behind an interface (enables testing without DB)
6. **`service.go`** — business logic (stateless, concurrent-safe, uses types + store)
7. **`handler.go`** — HTTP handler (thin: parse, delegate, respond)
8. **`*_test.go`** — table-driven tests using `t.Context()`
9. **`cell.yaml`** — manifest (metadata, exact cell deps, validation, invariants)
10. **`AGENTS.md`** — agent guide for the next agent
11. **Wire in `internal/app/wiring.go`** — the ONLY place cells are constructed

**Why this order matters:** Each layer depends only on layers above it. Types
and errors have no dependencies. The public `api` package is the only
cross-cell boundary; the root cell package stays private. The store depends on
types, the service on the API contract + store, and the handler on the service.

## Invariants

- Files ≤ 300 LOC, functions ≤ 40 LOC (enforced by structural tests)
- No `init()` functions (enforced by structural tests + `gochecknoinits` linter)
- Cells import only other cells' `api` packages, never implementations
- All DI in `internal/app/wiring.go` (< 300 LOC)
- Behavior-first kebab-case naming: `user-authenticate`, not `auth`
- Every cell has: `cell.yaml`, `doc.go`, `AGENTS.md`
- Store is always an interface; service is stateless after construction

## Common Tasks

### Create a new cell

```bash
task new-cell ID=user-authenticate
# Build inside-out: API → types/errors → store → service → handler → tests
# Fill in TODOs in the generated files
# Wire in internal/app/wiring.go
task index && task test
```

### Create a domain cell with sub-actions

```bash
task scaffold-domain ID=users
task new-cell ID=users/user-invite
task new-cell ID=users/user-register
```

### Modify a cell

```bash
task context ID=user-authenticate   # Read context pack first
task deps ID=user-authenticate      # Check blast radius
task scope ID=user-authenticate     # Declare the permitted edit boundary
# ... make changes (inside-out if touching multiple layers) ...
task changed                       # Status + impact (what did I touch?)
task verify-scope ID=user-authenticate # Reject undeclared edits
task quick-check                   # Fast post-edit validation
task test-cell ID=user-authenticate
```

`scope` permits only the selected cell (not nested child cells), its generated
context pack, and the generated index. For intentional cross-cell or shared
integration work, declare each additional boundary explicitly, for example
`task verify-scope ID=user-authenticate WITH=profiles,@wiring`. Dependencies
and dependents never expand scope implicitly.

### Find which cell handles a concept

```bash
task find-cell QUERY=payment
```

### Agent workflow loop

```bash
task orient                  # First command: cells + doctor + status
task find-cell QUERY=<topic> # Find the right cell
task context ID=<id>         # Read its context pack
task deps ID=<id>            # Check blast radius
task scope ID=<id>           # Declare permitted edit paths
# ... edit inside-out ...
task changed                 # Status + impact
task verify-scope ID=<id>    # Confirm no unrelated files changed
task ready                   # Pre-handoff: doctor + impact + test + status
```

## Reliability & Concurrency

- **Concurrency safety:** Each cell's `service` struct should be safe for concurrent use (stateless except `deps`, which is read-only after construction).
- **Context handling:** All I/O methods take `context.Context` as first parameter.
- **Error taxonomy:** Sentinel errors + typed error structs per cell. Use `errors.Is` for sentinel checks, `errors.As` for typed data.
- **No init() functions:** All construction is explicit in `wiring.go`.

## Validation

Run these commands before committing:

```bash
task ready            # Pre-handoff: doctor + impact + test + status
task doctor           # Fast architecture health check
task test             # Full test suite with race detector
task secrets          # Scan the working tree for hardcoded secrets
task fuzz             # Optional bounded fuzzing for manifest parsing
```

Fast post-edit validation:

```bash
task quick-check      # structure + policy + manifests + index (no tests)
```

### Task environment

Task commands place Go, ccache, and golangci-lint state in a writable temporary
directory by default. This keeps the standard validation loop reliable in agent
sandboxes with read-only home directories. Set `GOCACHE`, `CCACHE_DIR`, or
`GOLANGCI_LINT_CACHE` explicitly when your environment needs a different cache
location.

## Git Hooks

This repo includes git hooks in `.githooks/`:

- **pre-commit:** Runs `task doctor` and `task secrets` (architecture invariants and hardcoded-secret detection)
- **pre-push:** Runs `task test` (full test suite)

Install with: `task install-hooks`

## Modern Go

This project uses Go 1.26.5. Always use modern Go idioms:

- `any` not `interface{}`
- `errors.Is`/`errors.AsType[T]` not `==` comparison
- `slices.Contains` not manual loop
- `for i := range n` not `for i := 0; i < n; i++`
- `cmp.Or` not if-else for defaults
- `min`/`max` built-ins not if-else
- `t.Context()` not `context.WithCancel(context.Background())` in tests
- `wg.Go(func() { ... })` not `wg.Add(1)` + `go func() { defer wg.Done(); ... }()`
- `new(42)` not `x := 42; &x`
- `omitzero` when the JSON contract omits zero values; preserve `omitempty` when empty slices or maps must be omitted
- `b.Loop()` not `for i := 0; i < b.N; i++` in benchmarks

See the `use-modern-go` skill for the full reference.

## Agent Integration

`AGENTS.md`, `.agents/skills/`, and `Taskfile.yml` are the portable agent
integration surface. Add IDE-specific configuration only when the project
adopts that IDE.

### Zed

In a trusted Zed worktree, this repository's `AGENTS.md` and `.agents/skills/`
are available to the Zed Agent. Use `.zed/tasks.json` for the root workflow or
the `Reference:` tasks for `examples/reference-project`; validation tasks save
open buffers first. `.zed/settings.json` keeps `gen/` visible because its cell
index and context packs are part of the navigation contract.

### Visual Studio Code

In VS Code, use `.vscode/tasks.json` through **Terminal: Run Task**. Its
input-backed cell commands preserve the Taskfile's exact `ID` and optional
`WITH` scope contract. `.vscode/settings.json` enables root and cell-local
`AGENTS.md` discovery; VS Code also discovers the existing `.agents/skills/`
directly. The committed configuration recommends Go and YAML support only.

### Available Skills

| Skill                   | When to Use                            |
| ----------------------- | -------------------------------------- |
| `architecture-feedback` | Recording cell-architecture learnings  |
| `build-cell`            | Creating a new cell from scratch       |
| `modify-cell`           | Modifying an existing cell             |
| `navigate-cells`        | Finding and understanding cells        |
| `use-modern-go`         | Writing or modifying Go code           |

### Before starting work

```bash
task cells                    # Orient: what cells exist?
task find-cell QUERY=<topic>  # Find the right cell
task context ID=<id>          # Read its context pack
task scope ID=<id>            # Confirm the edit boundary
```

### After making changes

```bash
task test-cell ID=<id>        # Test the cell you changed
task impact                   # Check what else is affected
task verify-scope ID=<id>     # Confirm scope before handoff
task doctor                   # Verify architecture invariants
task secrets                  # Scan the working tree for hardcoded secrets
```

### Before committing

```bash
task doctor                   # architecture health
task test                     # Full tests
task secrets                  # Hardcoded-secret scan
```

### architecture feedback

Use the `architecture-feedback` skill only for learnings that directly affect the cell architecture: cell schema, inside-out build flow, manifests, wiring/dependencies, context packs, impact analysis, and structural guardrails. Do not log generic repository maintenance such as `.editorconfig`, `.github`, CI metadata, issue templates, dependency update configuration, licensing, README polish, or editor settings unless it exposes a direct cell-architecture problem.

## Do Not

- Do not import `internal/app/wiring.go` from a cell.
- Do not import another cell's implementation — use its `api` package only.
- Do not add `init()` functions.
- Do not create files over 300 LOC.
- Do not create functions over 40 LOC (use `AGENT_OVERRIDE` for documented exceptions).
- Do not skip `cell.yaml`, `doc.go`, or `AGENTS.md` when creating a cell.
- Do not build out-of-order (e.g., handler before public API). Always build inside-out.
- Do not use outdated Go patterns when modern alternatives exist (see `use-modern-go` skill).
