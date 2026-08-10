# 19. Agent Workflow Automation

The portable automation contract is the Taskfile. It supports both human-readable and JSON output without depending on a particular IDE.

```text
orient → find-cell → context + deps → edit → changed → ready
```

## Automation rules

- Use `task quick-check` for frequent structural feedback.
- Use `task ready` before handoff; it runs doctor, impact analysis, tests, and status.
- `task test` covers the bootstrap tooling and reference project as well as application tests.
- Use `task fuzz` for a bounded parser-hardening pass when changing manifest parsing; it is intentionally opt-in.
- Use `task index` after changing a cell manifest or guide. `check-index` rejects stale, missing, and orphan context packs.
- Consume `*-json` tasks from external tooling when structured output is needed.

IDE hooks may invoke these tasks, but they must not be the only enforcement layer. Git hooks and CI should run `task doctor` and `task test`.
