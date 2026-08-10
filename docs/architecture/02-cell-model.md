# Cell Model

## Project layout

```text
AGENTS.md                 root operating guide
Taskfile.yml              executable agent workflow
.agents/skills/           project-local task guidance
.vscode/                  optional VS Code tasks, settings, and recommendations
internal/
  app/wiring.go           explicit composition root
  cells/<id>/             capability cells
  contracts/              cross-cutting contracts only
  platform/               concrete infrastructure
policy/imports.yaml       enforced internal import policy
tools/agent/              agent navigation, guardrail, and scaffold tools
gen/cells.json            generated cell index
gen/context/              generated bounded context packs
```

## Flat cell

```text
internal/cells/<id>/
├── api/api.go            public interfaces, types, and errors
├── cell.yaml             machine-readable manifest
├── AGENTS.md             local agent guide
├── doc.go                package documentation
├── <name>.go             private implementation and constructor
├── types.go, errors.go   domain vocabulary and taxonomy
├── service.go, store.go  business logic and data access boundary
├── handler.go            transport adapter
└── *_test.go             collocated behavior tests
```

Use a flat cell for one primary behavior. Use a domain cell when several
distinct sub-actions share a domain vocabulary and error taxonomy. A domain
cell has `cell.go`, `model.go`, and `errors.go`, with each sub-action following
the flat-cell schema beneath it.

## Manifest and generated context

Every `cell.yaml` requires:

```yaml
id: user-authenticate
purpose: "Authenticate users and manage sessions"
entrypoints:
  - file: api/api.go
    symbol: Authenticator
dependencies: []
validation:
  - go test ./internal/cells/user-authenticate/...
```

`id` must match the behavior-first directory name. Entrypoints must name a
contained Go file and declared top-level symbol. Dependencies are exact existing
cell IDs and must match direct imports of those cells' `api` packages.
`invariants` is optional and records properties an agent must preserve.

`task index` derives `gen/cells.json` and `gen/context/<id>.context.md` from
the manifest and local guide. `task check-index` rejects stale, missing, or
orphaned generated metadata. Use `task context ID=<id>` for the compact pack;
never edit generated files by hand.

