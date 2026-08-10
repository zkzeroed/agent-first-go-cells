# 3. Directory Tree

```text
AGENTS.md                 root operating guide
Taskfile.yml              executable agent workflow
.agents/skills/           project-local task guidance
.vscode/                  optional VS Code tasks, settings, and recommendations
internal/
  app/wiring.go           explicit composition root
  cells/<id>/
    api/api.go            public cell contract
    cell.yaml             validated metadata
    AGENTS.md             cell-local guidance
    types.go, errors.go, service.go, store.go, handler.go
  contracts/              cross-cutting contracts only
  platform/               concrete infrastructure
policy/imports.yaml       enforced internal import policy
tools/agent/              agent navigation, guardrail, and scaffold tools
gen/cells.json            generated v2 cell index
gen/context/              generated bounded context packs
docs/architecture/        design and operating references
```

Flat cells live directly below `internal/cells/`. Domain cells may contain one level of sub-actions, such as `users/user-invite`; each manifest-backed directory has its own `api/api.go` contract package.

The cell root is implementation. Other cells import only its `api` package. All construction remains in `internal/app/wiring.go`.
