# Agent Workflow

Use this guide to navigate, change, and validate the repository safely. Treat
`Taskfile.yml` and repository tooling as the executable source of truth. Read
[the architecture map](docs/architecture/README.md) only when its focused
reference material is relevant to the task.

## Start Here

```bash
task orient
task --list
```

For cell work, find the target and read its bounded context before editing:

```bash
task find-cell QUERY=<topic>
task context ID=<cell-id>
task deps ID=<cell-id>
task scope ID=<cell-id>
```

Use the project-local skill that matches the task: `navigate-cells` for
orientation, `build-cell` for new capabilities, `modify-cell` for existing
cells, and `use-modern-go` whenever editing Go.

## Build Convention

Build each cell inside-out: domain types and errors, its `api` contract,
implementation and constructor, store, service, handler, tests, manifest, local
guide, then explicit wiring.

Use `task new-cell ID=<behavior-name>` for a single capability. Use
`task scaffold-domain ID=<domain>` before adding related sub-actions such as
`<domain>/<action>`. Generated `gen/` indexes and context packs are derived
artifacts: regenerate them with `task index`; never edit them directly.

## Invariants

- Name cells by behavior in kebab-case; each cell owns `cell.yaml`, `doc.go`,
  `AGENTS.md`, and `api/api.go`.
- Cross-cell imports target only another cell's `api` package. Declare every
  direct dependency by exact cell ID in `cell.yaml`.
- Construct dependencies only in `internal/app/wiring.go`; do not use `init()`,
  reflection DI, or service locators.
- Keep non-test cell files within 300 lines and functions within 40 lines unless
  the function documents `AGENT_OVERRIDE`.

## Common Tasks

Create a cell:

```bash
task new-cell ID=user-authenticate
# Implement inside-out, wire it, then:
task index
task test-cell ID=user-authenticate
```

Modify a cell:

```bash
task context ID=<cell-id>
task deps ID=<cell-id>
task scope ID=<cell-id>
# Edit only the declared boundary; use WITH= for intentional shared work.
task changed
task verify-scope ID=<cell-id> WITH=<optional-scope>
task quick-check
task test-cell ID=<cell-id>
```

## Reliability & Concurrency

- Services are immutable after construction and safe for concurrent use.
- I/O methods take `context.Context` first. Keep storage behind interfaces.
- Expose sentinel and typed errors deliberately; callers use `errors.Is` and
  `errors.As`.
- Target the Go version in `go.mod` and use the current idioms documented by the
  `use-modern-go` skill.

## Validation

Run the narrowest relevant checks while working, then use the handoff commands:

```bash
task quick-check              # Structural, policy, manifest, and index checks
task ready                    # Doctor, impact analysis, and race-enabled tests
task secrets                  # Working-tree hardcoded-secret scan
```

Run `task secrets-history` for an initial audit or suspected credential exposure.
Run `task fuzz` when changing manifest parsing. `task tools:check` verifies the
required local toolchain; `task setup` installs the Git hooks.

## Git Hooks

- **pre-commit:** `task doctor` and `task secrets`.
- **pre-push:** `task test`.

Install them with `task install-hooks`.

## Do Not

- Do not import `internal/app` from a cell or another cell's implementation.
- Do not add `init()` functions or construct dependencies outside wiring.
- Do not manually edit generated indexes or context packs.
- Do not expand a cell edit boundary implicitly through dependencies or
  dependents; declare each additional cell or shared surface with `WITH=`.
- Do not add generic code to `internal/contracts/` or `internal/platform/`; use
  them only for genuine shared surfaces.
- Do not record generic maintenance in the architecture feedback log. Use it
  only for learnings about cell structure, manifests, wiring, or guardrails.
