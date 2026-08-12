# Cell Model

## Project layout

```text
AGENTS.md                 root operating guide
Taskfile.yml              executable agent workflow
.agents/skills/           project-local task guidance
.vscode/                  optional VS Code tasks, settings, and recommendations
internal/
  app/wiring.go           application composition root
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

Public library packages also require `conformance` metadata. It makes the
implementation basis, current status, local research citations, and known gaps
available through generated context. Use `paper-defined-math` with citations
for direct research implementations; use `fixed-profile-policy` or
`engineering-primitive` with a rationale for deliberate non-paper decisions.
Any status other than `conformant` must list its remaining gaps.

Each record also declares `evidence: verified` or `unverified`. A conformant
record must be verified. Verified citations are resolved inside the project:
PDF citations use `locator.type: pdf-pages`; Markdown citations use
`locator.type: markdown-heading`. The cited symbols must be exported by the
library package. Use `task conformance ID=<id>` to validate and print only this
research contract.

`task index` derives `gen/cells.json` and `gen/context/<id>.context.md` from
the manifest and local guide. `task check-index` rejects stale, missing, or
orphaned generated metadata. Use `task context ID=<id>` for the compact pack;
never edit generated files by hand.

## Library packages

An exportable package is a `kind: library-package` manifest registered by ID
and directory in `policy/architecture.yaml`. Its directory is a normal Go
package path (including `.` for the module root), not `internal/cells`. It may
compose a declared private cell, but exported declarations must not expose
private-cell or `internal` types. Use `task new-package ID=field PATH=field`
to scaffold and register one.

A library is not an executable application: do not add `cmd/` or
`internal/app/wiring.go` merely to construct it. The registered public package
is the composition root and returns its stable exported API to downstream
consumers.

A registered public package may directly import only its declared private-cell
implementations. It may not import `internal/app`, `internal/platform`, or
`internal/contracts`; package-local helpers belong beside the API or below that
package's own `internal/` directory.

Private cells always live under `internal/cells`; this is a fixed convention,
not a configurable project setting. Public packages may use any non-`internal`
relative directory—including `pkg/`—but that directory is part of the public
Go import path, not a special export mechanism.

Every library package must include a conformance record. `task new-package`
creates a clearly marked placeholder so a package cannot be mistaken for
research-backed production behavior before its provenance is documented.
