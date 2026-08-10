# Examples

`reference-project/` is a small runnable application that demonstrates the
current Agent-First cell architecture. Keep it outside a project's application
tree; use it as an executable pattern when creating the first real cells.

The repository navigation tools can inspect it in place without copying root
tooling:

```bash
task index ROOT=examples/reference-project
task cells ROOT=examples/reference-project
task deps ID=greeting-render ROOT=examples/reference-project
task find-cell QUERY=greeting ROOT=examples/reference-project
task check-index ROOT=examples/reference-project
task index-json ROOT=examples/reference-project
task cells-json ROOT=examples/reference-project
task deps-json ID=greeting-render ROOT=examples/reference-project
task impact ROOT=examples/reference-project
```

`context` is repository-root only; inspect the reference project's committed
`gen/context/` packs directly.
