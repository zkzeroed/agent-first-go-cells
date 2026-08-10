# Bootstrap and Reference Project

Use this repository as a starter kit. The executable templates are
`Taskfile.yml`, `tools/agent/new-cell.sh`, and `tools/agent/new-domain.sh`;
they are authoritative over copied documentation examples.

```bash
task setup
task new-cell ID=user-authenticate
# implement the public API and then the cell inside-out
task index
task ready
```

`task new-cell` validates IDs, refuses to overwrite an existing path, and
creates a compilable flat or nested sub-action skeleton. Use
`task scaffold-domain ID=users` before adding related actions such as
`task new-cell ID=users/user-invite`.

The runnable [`examples/reference-project`](../../examples/reference-project/)
contains a multi-cell graph, API dependencies, a nested domain action, explicit
wiring, and a cryptographic token flow. Use it to exercise the tooling without
inventing static examples:

```bash
task cells ROOT=examples/reference-project
task ready ROOT=examples/reference-project
```

