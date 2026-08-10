# 19. Agent Workflow Automation

The portable automation contract is the Taskfile. It supports both human-readable and JSON output without depending on a particular IDE.

```text
orient → find-cell → context + deps + scope → edit → changed → verify-scope → ready
```

## Automation rules

- Use `task quick-check` for frequent structural feedback.
- Use `task scope ID=<id>` before editing a cell. It is the explicit change
  contract; `task verify-scope` fails closed on edits outside it.
- Expand scope only with exact cell IDs or `@contracts`, `@platform`, and
  `@wiring`; dependencies and dependents are never authorized implicitly.
- Use `task ready` before handoff; it runs doctor, impact analysis, tests, and status.
- `task test` covers the bootstrap tooling and reference project as well as application tests.
- Use `task fuzz` for a bounded parser-hardening pass when changing manifest parsing; it is intentionally opt-in.
- Use `task index` after changing a cell manifest or guide. `check-index` rejects stale, missing, and orphan context packs.
- Consume `*-json` tasks from external tooling when structured output is needed.

IDE hooks may invoke these tasks, but they must not be the only enforcement layer. Git hooks and CI should run `task doctor` and `task test`.
