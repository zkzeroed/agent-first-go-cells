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

`context` accepts `ROOT`, so it can print a reference-project context pack:

```bash
task context ID=greeting-render ROOT=examples/reference-project
```

`library-project/` demonstrates library mode: its `greeting/` package is
externally importable and composes two private cells under `internal/cells/`.
Use the same commands with `ROOT=examples/library-project`.

`research-library-project/` demonstrates a paper-defined conformance record
on an exported `modulus/` package, with a private reduction helper:

```bash
task context ID=modulus ROOT=examples/research-library-project
task ready ROOT=examples/research-library-project
```

`module-root-library-project/` demonstrates `PATH=.` with an importable module
root package and a lean private helper cell:

```bash
task context ID=canonical ROOT=examples/module-root-library-project
task ready ROOT=examples/module-root-library-project
```
