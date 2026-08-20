# Agent Operations

`AGENTS.md` is the always-on guide. `.agents/skills/` provides task-specific
procedures, and `Taskfile.yml` is the portable automation contract.

## Standard workflow

```text
orient → find-cell → context + deps + scope → edit → changed → verify-scope → ready
```

```bash
task orient
task find-cell QUERY=<topic>
task context ID=<cell-id>
task deps ID=<cell-id>
task scope ID=<cell-id>
# edit
task changed
task verify-scope ID=<cell-id> WITH=<optional-scope>
task ready
```

Use `task --list` for the complete command surface. Human-readable and JSON
forms are available for navigation commands. Task defaults place Go, ccache,
and golangci-lint state in writable temporary locations; callers may override
`GOCACHE`, `CCACHE_DIR`, and `GOLANGCI_LINT_CACHE`.

## Library workflow

Use `task new-package ID=<id> PATH=<path>` for an exportable package cell.
It registers the package in `policy/architecture.yaml`; `PATH=.` requires
`PACKAGE=<go-package-name>`. Public package APIs must not leak `internal` or
private-cell types. All navigation and validation commands accept
`ROOT=<project-root>`, including `context`, `test-cell`, and `validate-cell`.

For research-driven packages, update the scaffolded `conformance` record before
implementation. Use `task context ID=<id>` to see its basis, status, citations,
and gaps; `task cells-json` exposes the same structured data to agents and
automation. Cite local source files, paper sections, PDF pages, and affected
symbols rather than placing untraceable paper claims in code comments. Use
`task conformance ID=<id>` before declaring evidence verified; it resolves the
local source and citation locator and verifies cited exported symbols.

## Editor support

Zed loads the root guide and project skills in trusted worktrees. Its committed
tasks expose the root workflow, and its settings keep generated context visible.

VS Code task templates expose the same Taskfile workflow, including prompted
cell IDs and scope expansion. Workspace settings enable root and nested
`AGENTS.md` discovery; VS Code discovers `.agents/skills/` directly. The
recommended Go extension formats with `gopls` using gofumpt rules and organizes
imports on save. The recommended YAML extension validates every `cell.yaml`
against `policy/manifest.schema.json`. `Agent: Ready` is the default build task
and `Agent: Test` is the default test task. Do not duplicate portable
instructions in editor-specific files.

## Multi-agent work

One agent owns a cell at a time. `internal/platform/`, `internal/contracts/`,
and application wiring are shared integration surfaces and require explicit
coordination. A worker stays in its declared scope, tests its cell, then hands
off the integration work with affected files, validation results, and impact.

The coordinator or work tracker manages claims and release; this repository
does not add persistent lock state. Before handoff, report status, changed
files, tests, lint, impact, and any remaining integration step.
