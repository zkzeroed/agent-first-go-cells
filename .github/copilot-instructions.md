# Copilot Instructions

This repository uses the Agent-First Go architecture.

## Start Here

- Read `AGENTS.md` before making changes.
- Read `docs/architecture/README.md` for the architecture overview.
- Use behavior-first cell names such as `user-authenticate`, not generic names like `auth`.

## Cell Rules

Every capability is a cell under `internal/cells/<id>/` with a fixed schema:

- `types.go`
- `errors.go`
- `<name>.go`
- `store.go`
- `service.go`
- `handler.go`
- `*_test.go`
- `cell.yaml`
- `doc.go`
- `AGENTS.md`

Build cells inside-out: public API, types and errors, store, service, handler, tests, manifest, agent guide, wiring.

## Constraints

- Files must stay ≤ 300 LOC.
- Functions must stay ≤ 40 LOC.
- Do not add `init()` functions.
- Keep dependency injection in `internal/app/wiring.go`.
- Use modern Go idioms from `AGENTS.md`.
- Tests should use `t.Context()`.

## Validation

Run the most specific relevant checks, then broader checks:

```bash
task orient                  # First command: cells + doctor + status
task test-cell ID=<cell-id>  # Test the cell you changed
task changed                 # Status + impact (what did I touch?)
task verify-scope ID=<cell-id> # Reject undeclared edits
task quick-check             # Fast post-edit validation (no tests)
task ready                   # Pre-handoff: doctor + impact + test + status
```

## Cell architecture Feedback

Update `docs/architecture/06-architecture-feedback-log.md` only when work reveals a learning that directly affects the cell architecture: cell schema, inside-out build flow, manifests, wiring/dependencies, context packs, impact analysis, or structural guardrails. Do not log generic repository maintenance.
