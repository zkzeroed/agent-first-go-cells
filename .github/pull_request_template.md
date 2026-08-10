## Summary

-

## Type of change

- [ ] Cell implementation or modification
- [ ] architecture/tooling change
- [ ] Documentation change
- [ ] Bug fix
- [ ] Refactor

## architecture checklist

- [ ] I followed the cell file schema and inside-out build order, if applicable.
- [ ] Cell dependencies are declared in `cell.yaml`, if applicable.
- [ ] New or changed cells include/maintain `AGENTS.md`, `doc.go`, and `cell.yaml`.
- [ ] Files stay within 300 LOC and functions stay within 40 LOC, or an `AGENT_OVERRIDE` is documented.
- [ ] No `init()` functions were added.
- [ ] All DI remains in `internal/app/wiring.go`, if applicable.

## Validation

- [ ] `task doctor`
- [ ] `task test`
- [ ] `task impact`, if cell behavior or dependencies changed

## Cell architecture feedback

If this change revealed friction or an improvement opportunity that directly affects the cell architecture — cell schema, inside-out build flow, manifests, wiring/dependencies, context packs, impact analysis, or structural guardrails — add an entry to `docs/architecture/06-architecture-feedback-log.md`. Do not log generic repository maintenance.
