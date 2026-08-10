# 15. Trade-offs & Mitigations

## For Human Developers

| Trade-off                       | Impact                               | Mitigation                                                                                  |
| ------------------------------- | ------------------------------------ | ------------------------------------------------------------------------------------------- |
| **More files per feature**      | Humans open 6–8 files instead of 1–2 | Fixed schema means muscle memory. `task new-cell` scaffolds automatically.                  |
| **No shared model packages**    | No single `models.User` everywhere   | Each cell owns its types. Cross-cell transfer via interfaces or DTOs. Prevents god-objects. |
| **Strict file/function limits** | May feel constraining                | 300 LOC file and 40 LOC function limits keep context bounded; `AGENT_OVERRIDE` documents function exceptions. |
| **Manual wiring**               | `wiring.go` grows with each cell     | < 300 LOC limit, split by concern. Cell-to-cell imports reduce wiring burden.               |
| **No `init()` functions**       | No auto-registration                 | Explicit wiring is a feature. Dependency graph is visible and testable.                     |

## For Agents

| Trade-off                     | Impact                                      | Mitigation                                                                                           |
| ----------------------------- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| **More boilerplate per cell** | Agent generates more files                  | `task new-cell ID=<id>` scaffolds all files. Fixed schema = mechanical generation.                   |
| **Wiring file bottleneck**    | Multiple agents coordinate on `wiring.go`   | Atomic updates, < 300 LOC, split by concern. Cell-to-cell imports reduce wiring.                     |
| **Manifest maintenance**      | `cell.yaml` can drift                       | `task check-manifests` in CI. Only 5 mandatory fields — minimal overhead.                            |
| **Index staleness**           | `gen/cells.json` can become stale           | Hash-based `task check-index` in CI. No git dependency.                                              |
| **Impact analysis limitations** | Declared metadata may miss undeclared coupling | Go-based and testable. Covers cell ownership plus declared dependencies; broader tests catch remaining. |
| **Contracts coordination**    | Adding shared interface touches shared file | One file per concern in `contracts/`. Cell-to-cell interfaces stay in cells.                         |
| **Over-normalization**        | Too many tiny cells                         | Sub-action guidance: use sub-actions within domain cells for related operations.                     |
