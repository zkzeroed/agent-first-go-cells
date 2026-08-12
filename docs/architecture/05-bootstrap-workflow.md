# Bootstrap Workflow

Use this repository as a starter kit. The executable templates are
`Taskfile.yml`, `tools/agent/new-cell.sh`, `tools/agent/new-domain.sh`, and
`tools/agent/new-package/`;
they are authoritative over copied documentation examples.

```bash
task setup
task new-cell ID=user-authenticate
# implement the public API and then the cell inside-out
task index
task ready
```

`task new-cell` validates IDs, refuses to overwrite an existing path, and
creates a lean compilable flat or nested sub-action skeleton. Use
`task new-cell-ext` when the cell needs the former extended application
template with types, errors, service, store, and handler files. Use
`task scaffold-domain ID=users` before adding related actions such as
`task new-cell ID=users/user-invite`.
For an exportable Go module package, use `task new-package ID=field PATH=field`.
The command creates a `kind: library-package` manifest and registers the
package in `policy/architecture.yaml`. `PATH=.` creates a module-root package
when `PACKAGE=<go-package-name>` is supplied. Public packages may compose
declared private helper cells but must keep those types out of exported APIs.
They may not import module-level `internal/app`, `internal/platform`, or
`internal/contracts`. A library does not need `cmd/` or
`internal/app/wiring.go`: its registered public package is the composition
root.
