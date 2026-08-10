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

Read only the chapter relevant to the task:

1. [Core principles](02-core-principles.md)
2. [Directory tree](03-directory-tree.md)
3. [Cell schema](04-cell-file-schema.md)
4. [Manifest schema](05-manifest-schema.md)
5. [Generated artifacts](06-generated-artifacts.md)
6. [Dependency wiring](08-dependency-wiring.md)
7. [Guardrails](10-guardrail-tooling.md)
8. [Task targets](14-taskfile-targets.md)
9. [Agent conventions](17-agent-conventions.md)

The repository files and Taskfile are authoritative where they differ from prose.
