---
name: navigate-cells
description: Navigate and orient in an Agent-First Go architecture codebase. Use when arriving at a project, finding which cell handles a concept, or understanding the cell dependency graph.
---

# Navigate Cells

## First Commands (Orientation)

```bash
task cells                        # list all cells with status
task find-cell QUERY=payment      # find cell by concept
task context ID=user-authenticate # read a cell's context pack
task deps ID=user-authenticate    # see dependency graph
task scope ID=user-authenticate   # print the explicit edit boundary
```

For the reference project, root-capable tools accept `ROOT=examples/reference-project`:

```bash
task cells ROOT=examples/reference-project
task deps ID=greeting-render ROOT=examples/reference-project
task find-cell QUERY=greeting ROOT=examples/reference-project
task impact ROOT=examples/reference-project
```

`context` is repository-root only.

## Navigation Workflow

### Arriving at a codebase

1. `task cells` — see what cells exist and if index is fresh
2. `task doctor` — check architecture health
3. Read `AGENTS.md` (root) — project purpose and conventions
4. `task find-cell QUERY=<your-task>` — find the relevant cell
5. `task context ID=<cell-id>` — read the cell's context pack (compact summary)

### Finding which cell handles a concept

```bash
task find-cell QUERY=payment
task find-cell QUERY=auth
task find-cell QUERY=email
```

This searches cell IDs, purposes, and AGENTS.md content. More precise than grep.

### Understanding blast radius before a change

```bash
task deps ID=user-authenticate    # who depends on this cell?
task impact                       # direct owners and transitive downstream cells
task scope ID=user-authenticate   # owned files and explicit integration scope
```

Before handoff, run `task verify-scope ID=<id>`. Add only explicit sibling
cells or `@contracts`, `@platform`, and `@wiring` through `WITH=` when the task
crosses those boundaries.

### Reading a cell's full context

```bash
task context ID=user-authenticate
```

This prints the generated context pack: purpose, entrypoints, dependencies, invariants, validation commands.

For full detail, read the cell's files in order:

1. `cell.yaml` — metadata
2. `AGENTS.md` — agent guide
3. `api/api.go` — public contract (the membrane)
4. `types.go` — domain types
5. `errors.go` — error taxonomy
6. `service.go` — business logic
7. `store.go` — data access
8. `handler.go` — HTTP handler

## Directory Layout

```
internal/cells/<cell-id>/
├── cell.yaml          # manifest
├── AGENTS.md          # agent guide
├── doc.go             # package doc
├── api/api.go         # public contract
├── <name>.go          # implementation + New()
├── types.go           # domain types
├── errors.go          # error taxonomy
├── service.go         # business logic
├── store.go           # data access
├── handler.go         # HTTP handler
└── *_test.go          # tests
```

Domain cells have sub-actions:

```
internal/cells/users/
├── cell.yaml
├── cell.go            # aggregates sub-actions
├── api/api.go         # domain public contract
├── model.go           # shared types
├── errors.go          # shared errors
├── user-invite/       # sub-action (same schema as flat cell)
├── user-register/     # sub-action
└── user-profile/      # sub-action
```
