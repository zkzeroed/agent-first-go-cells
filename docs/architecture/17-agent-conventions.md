# 17. Agent Conventions

`AGENTS.md` is the always-on project guide. Project-local skills in `.agents/skills/` provide deeper workflows when relevant. The Taskfile is the portable automation interface; IDE-specific configuration is optional.

For Zed, commit only project-scoped editing settings and task templates. This
starter's `.zed/tasks.json` exposes the root workflow and the selected-project
reference workflow; validation tasks save buffers before running. `AGENTS.md`
and `.agents/skills/` load in trusted worktrees. Keep `gen/` visible in Zed:
the generated index and context packs are part of the agent navigation
contract, not disposable build output.

For VS Code, `.vscode/tasks.json` exposes the portable Taskfile workflow;
input-backed tasks collect exact cell IDs and optional scope expansion rather
than inferring them. `.vscode/settings.json` enables root and nested
`AGENTS.md` discovery. VS Code discovers the project-local `.agents/skills/`
directly, so do not duplicate those workflows in editor-specific instruction
files.

## Build order

Build cells inside-out:

1. `api/api.go` (public interfaces, shared types, public errors), then private types and errors
2. store interface and service
3. handler and tests
4. manifest and cell guide
5. `internal/app/wiring.go`

## Operating loop

Before editing a cell, run `task context ID=<id>`, `task deps ID=<id>`, and
`task scope ID=<id>`. After editing, run `task changed` and
`task verify-scope ID=<id>`; before handoff, run `task ready`. Cross-cell or
shared integration work requires explicit `WITH=` scope expansion.

Use `task tools:check` to verify the local toolchain.
