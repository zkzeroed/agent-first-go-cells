# 10. Guardrail Tooling

Guardrails are executable. The task names below are the source of truth; use `task --list` for the complete current surface.

| Guardrail | Command | Enforces |
| --- | --- | --- |
| Manifests | `task check-manifests` | strict YAML, resolved entrypoints, exact direct API dependencies, validation |
| Generated metadata | `task check-index` | fresh index and context packs |
| Structure | `task structure-test` | cell schema, API imports, size limits, no `init()` |
| Import policy | `task policy` | internal dependency direction and allow/deny rules |
| Agent guides | `task check-agents` | required guide sections |
| Full handoff | `task ready` | doctor, impact, tests, and status |

Cell-to-cell imports must use `internal/cells/<id>/api`. The import policy and structural tests both enforce this boundary. `AGENT_OVERRIDE` applies only to the documented function that carries it.

Run `task quick-check` after a local edit and `task ready` before handoff. Generated files are refreshed by `task index`, never edited directly.
