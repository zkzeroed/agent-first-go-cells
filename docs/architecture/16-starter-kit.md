# 16. Starter Kit

Use this repository as the starter kit. The executable templates are `Taskfile.yml`, `scripts/new-cell.sh`, and `scripts/new-domain.sh`; they are authoritative and intentionally not duplicated here.

## Minimum bootstrap

```text
AGENTS.md                 root operating guide
Taskfile.yml              executable workflow
.agents/skills/           task-specific instructions
internal/cells/           capability cells
policy/imports.yaml       internal import policy
scripts/                  architecture tooling
gen/cells.json            generated v2 cell index
```

## First capability

```bash
task setup
task new-cell ID=user-authenticate
# implement the public contract in api/api.go, then the cell inside-out
task index
task ready
```

The scaffolder validates IDs and creates a compilable flat or domain skeleton. For a domain sub-action, use `task new-cell ID=users/user-invite`.
