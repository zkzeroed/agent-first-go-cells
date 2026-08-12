# Agent-First Go Cell Architecture

> A vertical-slice modular-monolith architecture for Go, with API membranes and
> tooling that makes a codebase easy for coding agents to navigate, change, and
> validate safely.

This is a bootstrap repository, not an application. Start a project here when
you want behavior-oriented modules, explicit composition, and guardrails that
are executable rather than merely documented.

It supports both private application cells and exportable library-package
cells. Application cells live in `internal/cells`; public package cells live
at normal Go import paths and are registered in `policy/architecture.yaml`.
`internal/cells` is the fixed private-cell convention; public packages may use
the module root or any non-`internal` relative path, including `pkg/`.

## The idea

A **cell** is one business capability. It owns its vocabulary, API contract,
implementation, tests, manifest, and guide in one directory. Cells form a
modular monolith: they communicate only through explicit `api` packages.

For an executable application, `cmd/<name>/main.go` owns process lifecycle
(configuration, signals, server start and shutdown) and delegates dependency
construction to `internal/app/wiring*.go`. For a library, do not add `cmd/` or
`internal/app/wiring.go`: the registered public package is its composition root
and constructs its declared private helpers behind the exported API.

The familiar foundations are vertical slices, modular boundaries, and
ports-and-adapters. The differentiator is the operational layer for agents:
machine-readable manifests, generated context packs, dependency and impact
analysis, explicit edit scopes, scaffolding, and structural checks.

```
task → cell → api/api.go → explicit wiring → validation
         │
         ├── cell.yaml     machine-readable intent and dependencies
         ├── AGENTS.md     local guidance
         └── gen/context   bounded orientation pack
```

## Start in two minutes

```bash
# See the project, its health, and your working tree.
task orient

# Find the capability, then read only its bounded context.
task find-cell QUERY=<topic>
task context ID=<cell-id>
task deps ID=<cell-id>
task scope ID=<cell-id>

# After an edit, prove that the change stayed within its declared boundary.
task changed
task verify-scope ID=<cell-id>
task ready
```

Create a focused capability with `task new-cell ID=user-authenticate`; it
starts with the manifest, guide, API contract, implementation, and test. Use
`task new-cell-ext ID=user-authenticate` when the capability needs the fuller
application template for types, errors, store, service, and handler layers.

For an exportable library package, use `task new-package ID=field PATH=field`.
Consumers import `<module>/field`; implementation helpers stay unexported in
that package, `field/internal/`, or declared private cells. A public package
may import only its declared `internal/cells` implementations; it must not
reach into module-level app, platform, or contract internals.

Choose the layout that matches the consumer:

```text
# Application: main owns lifecycle; wiring owns composition.
cmd/myapp/main.go
internal/app/wiring.go
internal/cells/
└── user-authenticate/
    ├── api/api.go
    ├── cell.yaml
    └── service.go

# Library: stable Go import path; no cmd/ or application wiring.
field/
├── cell.yaml            # kind: library-package, public: true
├── field.go             # exported downstream API
└── internal/            # package-private helpers (optional)
internal/cells/
└── field-encode/        # declared private helper cell (optional)
```

## Research-driven agentic engineering

Public library packages can record the provenance of their behavior directly in
`cell.yaml`. A conformance record tells the next agent whether the package
implements paper-defined mathematics, fixed profile policy, or an engineering
primitive; it also records implementation status, local paper citations, and
known gaps.

```yaml
conformance:
  basis: paper-defined-math
  status: simplified
  evidence: unverified
  citations:
    - file: docs/papers/protocol.pdf
      locator:
        type: pdf-pages
        pages: [12, 13]
      symbols: [Prove, Verify]
  gaps:
    - "Uses the reference algorithm; optimized proving is not implemented."
```

`task context ID=<id>` and `task cells-json` expose this provenance without
opening every source file. The manifest validator requires citations for
paper-defined work, rationale for policy or engineering work, and explicit
gaps for every non-conformant status. Set `evidence: verified` only after the
validator can resolve every local citation and its exported symbols; use
`task conformance ID=<id>` for that focused check. This turns research reading
into a bounded, reviewable input to implementation rather than an undocumented
assumption.

## What an agent can prove before handoff

| Question | Evidence | Command |
| --- | --- | --- |
| Where does this work belong? | Cell metadata and a bounded context pack | `task find-cell`, `task context` |
| What may change? | Exact cell ownership plus explicitly declared shared surfaces | `task scope` |
| What else can it affect? | Exact dependency graph and conservative shared-surface analysis | `task deps`, `task impact` |
| What research governs an export? | Conformance basis, citations, and known gaps | `task context`, `task cells-json` |
| Did the work stay safe? | Fail-closed scope check plus architecture, test, and status checks | `task verify-scope`, `task ready` |

Agents may extend scope only with exact cell IDs or the explicit shared-surface
tokens `@contracts`, `@platform`, and `@wiring`; dependencies and dependents do
not authorize edits implicitly.

## A cell at a glance

```text
internal/cells/user-authenticate/
├── api/api.go             public API contract; the only cross-cell membrane
├── types.go               domain vocabulary
├── errors.go              error taxonomy
├── user_authenticate.go   private implementation and constructor
├── store.go               data access interface and adapter
├── service.go             stateless business logic
├── handler.go             thin transport adapter
├── cell.yaml              purpose, exact dependencies, validation
├── AGENTS.md              guide for the next agent
└── *_test.go              behavior tests
```

## Non-negotiable boundaries

- Cells import another cell's `api` package, never its implementation.
- In an application, `cmd/<name>/main.go` starts the process and
  `internal/app/wiring*.go` is the sole private-cell composition root; neither
  uses `init()`.
- In a library, the registered public package is the composition root; it has
  no `cmd/` or application wiring.
- Every cell has `cell.yaml`, `doc.go`, and `AGENTS.md`.
- Files stay within 300 lines and functions within 40 lines unless a specific
  exception documents `AGENT_OVERRIDE`.
- Manifests declare exact existing cell IDs; generated metadata is never edited
  by hand.

## Guardrails that earn their place

| Concern | Command | Purpose |
| --- | --- | --- |
| Orientation | `task orient` | Cells, architecture health, and repository state |
| Fast check | `task quick-check` | Structure, policy, manifests, and generated index |
| Change scope | `task scope ID=<id>` | Explicit permitted edit boundary |
| Scope check | `task verify-scope ID=<id>` | Reject changes outside that boundary |
| Change impact | `task changed` | Git state and affected cells |
| Handoff | `task ready` | Doctor, impact analysis, tests, and state |
| Lint | `task lint` | Bootstrap tooling and reference module |
| Secret scan | `task secrets` | Working-tree scan before each commit |
| Parser hardening | `task fuzz` | Bounded fuzzing for manifest input |

Install the provided Git hooks with `task install-hooks`.
`task test` covers the bootstrap tooling, application tests, and reference
project; `task fuzz` is intentionally opt-in rather than part of every handoff.
`task secrets-history` scans all reachable Git history and is appropriate for
an initial audit or an incident response.
Run `task --list` for the complete command surface; [agent operations](docs/architecture/04-agent-operations.md)
explains the navigation and machine-readable variants.

## Agent integration

`AGENTS.md`, `.agents/skills/`, `Taskfile.yml`, manifests, and generated
metadata are the portable agent interface. Project-local skills cover
navigation, building, modification, modern Go, and token-efficient command
output.

Zed users can run the root and reference workflows from `.zed/tasks.json`; its
validation tasks save open buffers first. Zed Agent loads this repository's
`AGENTS.md` and project-local skills in a trusted worktree. The generated
`gen/` index and context packs intentionally remain visible for navigation.

VS Code users can run the same workflow through `.vscode/tasks.json`, including
input-backed context, dependency, scope, and scope-verification tasks. Its
workspace settings enable root and cell-local `AGENTS.md` discovery and reuse
the existing `.agents/skills/` directly.

## Modern Go

This repository targets **Go 1.26.5**. Prefer current standard-library and
language idioms: `any`, `errors.Is`/`errors.AsType`, `slices`, `cmp.Or`,
integer `range`, `t.Context()`, `wg.Go`, intentional `omitzero` tags, and
`b.Loop()`.

## Learn by running it

The [reference project](examples/reference-project/README.md) contains six
cells, API dependencies, a nested domain action, explicit wiring, and a
cryptographic token flow. Its complete workflow is available from the root:

```bash
task cells ROOT=examples/reference-project
task ready ROOT=examples/reference-project
```

The [library project](examples/library-project/README.md) demonstrates an
exportable `greeting` package composing private helper cells:

```bash
task cells ROOT=examples/library-project
task deps ID=greeting ROOT=examples/library-project
task ready ROOT=examples/library-project
```

## Architecture reference

The [architecture guide](docs/architecture/README.md) is organized into focused
operational references: read only the document relevant to the task. The
Taskfile and repository files are authoritative where prose and behavior differ.

## License

MIT
