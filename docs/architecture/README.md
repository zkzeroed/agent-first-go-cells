# Agent-First Go Architecture

This repository is a bootstrap for Go projects maintained primarily by coding agents. The design optimizes for local context, explicit boundaries, deterministic navigation, and executable guardrails.

## Operating loop

```bash
task orient
task find-cell QUERY=<topic>
task context ID=<cell-id>
task deps ID=<cell-id>
task scope ID=<cell-id>
# edit
task changed
task verify-scope ID=<cell-id>
task ready
```

## Core model

- A **cell** is one vertical capability in `internal/cells/<id>`.
- Its public contract is `internal/cells/<id>/api`; other cells import only that package.
- The cell root contains implementation, tests, `cell.yaml`, and `AGENTS.md`.
- `cell.yaml` supplies the machine-readable graph; dependencies use exact cell IDs.
- `gen/cells.json` and bounded `gen/context/` packs accelerate orientation.
- `internal/app/wiring.go` is the explicit composition root.

## Invariants

- Files are at most 300 lines and functions at most 40 lines unless the specific function documents `AGENT_OVERRIDE`.
- No `init()` functions.
- Internal imports obey `policy/imports.yaml`; cell dependencies cross only `api` packages.
- Generated metadata is never manually edited.
- `task verify-scope` rejects task edits outside the explicitly declared cell
  and optional integration scope.

## Reference

Read only the document relevant to the task:

1. [Architecture principles](01-architecture-principles.md)
2. [Cell model](02-cell-model.md)
3. [Boundaries and guardrails](03-boundaries-and-guardrails.md)
4. [Agent operations](04-agent-operations.md)
5. [Bootstrap workflow](05-bootstrap-workflow.md)
6. [Architecture feedback log](06-architecture-feedback-log.md)

The repository files and Taskfile are authoritative where they differ from prose.
