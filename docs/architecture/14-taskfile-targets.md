# 14. Task Targets

`Taskfile.yml` is the executable source of truth. Run `task --list` for the current complete list; do not copy a Taskfile from this documentation.

## Agent workflow

```bash
task orient
task find-cell QUERY=<topic>
task context ID=<cell-id>
task deps ID=<cell-id>
# edit
task changed
task ready
```

## High-signal targets

| Target | Use |
| --- | --- |
| `task orient` | cells, architecture health, and working-tree state |
| `task new-cell ID=<id>` | scaffold a flat cell |
| `task scaffold-domain ID=<id>` | scaffold a domain cell |
| `task context ID=<id>` | read a generated, bounded cell context pack |
| `task find-cell QUERY=<text>` | search cell metadata and guides |
| `task deps ID=<id>` | inspect exact declared dependencies |
| `task quick-check` | fast structural validation after an edit |
| `task changed` | changed files and affected cells |
| `task ready` | pre-handoff doctor, impact, tests, and status |
| `task lint` | lint the bootstrap tooling and reference module |
| `task tools:check` | verify the local Go development toolchain |
| `task fuzz` | opt-in bounded fuzzing for manifest parsing |

Machine-readable forms are available as `cells-json`, `index-json`, `deps-json`, `impact-json`, and `context-json`.

`index`, `check-index`, `index-json`, `cells`, `cells-json`, `deps`,
`deps-json`, `find-cell`, `impact`, `impact-json`, `changed`, and `ready`
accept `ROOT=<project-path>` for a compatible project, such as the reference
example. `context` operates only on the repository root.
