# Boundaries and Guardrails

## Dependency direction

Cells import only another cell's `api` package and declare the corresponding
exact ID in `cell.yaml`. `internal/app/wiring.go` constructs and connects cells;
it stays below 300 lines and may be split by concern when needed.
`internal/contracts/` holds only genuine cross-cutting interfaces with no
natural cell owner.

`policy/imports.yaml` codifies allowed and forbidden directions. `task policy`
checks it with Go AST parsing; structural tests independently enforce the
API-only cell import rule.

## Change boundaries

Before a cell edit, run `task deps ID=<id>` and `task scope ID=<id>`. Scope
permits the exact target cell, its context pack, and the generated index. A
nested child cell is never implicitly included.

After editing, `task verify-scope ID=<id>` fails on changed tracked or
non-ignored untracked files outside that boundary. Intentional expansion uses
only exact cell IDs or `@contracts`, `@platform`, and `@wiring` through `WITH=`.
Dependencies and dependents do not authorize edits implicitly. Deleted targets
are resolved from `HEAD` so removal cannot bypass ownership checks.

## Impact analysis

`task impact` maps changed Go, YAML, and cell-local guide files to their
most-specific owner, then reports transitive downstream dependents. Changes to
`internal/contracts/`, `internal/platform/`, or `internal/app/wiring*.go` are
shared surfaces: every cell is affected and full-project validation is required.
It accepts `ROOT=<project-path>` for compatible nested projects.

## Executable checks

| Concern | Command |
| --- | --- |
| Manifest source and dependency graph | `task check-manifests` |
| Generated index and context | `task check-index` |
| Schema, imports, size, and `init()` | `task structure-test` |
| Import direction | `task policy` |
| Required agent-guide sections | `task check-agents` |
| Fast post-edit feedback | `task quick-check` |
| Full handoff | `task ready` |

Run `task index` after a manifest or cell-guide change. Use `task fuzz` for a
bounded manifest-parser hardening pass when changing manifest parsing.

