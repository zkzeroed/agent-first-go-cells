# 17. Agent Conventions

`AGENTS.md` is the always-on project guide. Project-local skills in `.agents/skills/` provide deeper workflows when relevant. The Taskfile is the portable automation interface; IDE-specific configuration is optional.

For Zed, commit only project-scoped editing settings and task templates. This
starter's `.zed/tasks.json` exposes the root workflow and the selected-project
reference workflow; `AGENTS.md` and `.agents/skills/` load in trusted worktrees.

## Build order

Build cells inside-out:

1. `api/api.go` (public interfaces, shared types, public errors), then private types and errors
2. store interface and service
3. handler and tests
4. manifest and cell guide
5. `internal/app/wiring.go`

## Operating loop

Before editing a cell, run `task context ID=<id>` and `task deps ID=<id>`. After editing, run `task changed`; before handoff, run `task ready`.

Use `task tools:check` to verify the local toolchain.
