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

## Editor support

Zed loads the root guide and project skills in trusted worktrees. Its committed
tasks expose the root and reference workflows, and its settings keep generated
context visible.

VS Code task templates expose the same Taskfile workflow, including prompted
cell IDs and scope expansion. Workspace settings enable root and nested
`AGENTS.md` discovery; VS Code discovers `.agents/skills/` directly. Do not
duplicate portable instructions in editor-specific files.

## Multi-agent work

One agent owns a cell at a time. `internal/platform/`, `internal/contracts/`,
`pkg/`, and wiring are shared integration surfaces and require explicit
coordination. A worker stays in its declared scope, tests its cell, then hands
off the integration work with affected files, validation results, and impact.

The coordinator or work tracker manages claims and release; this repository
does not add persistent lock state. Before handoff, report status, changed
files, tests, lint, impact, and any remaining integration step.

