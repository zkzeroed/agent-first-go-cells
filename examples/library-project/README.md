# Library Project

This runnable library-mode example exposes `example.com/agent-first-library/greeting`.
Its public package composes the private `greeting-render` cell, which depends on
the nested private `text-normalize` cell. Consumers can import only `greeting`;
the implementation packages remain under `internal/`.

```bash
task cells ROOT=examples/library-project
task deps ID=greeting ROOT=examples/library-project
task deps ID=greeting-render ROOT=examples/library-project
task context ID=greeting ROOT=examples/library-project
task scope ID=greeting-render ROOT=examples/library-project
task ready ROOT=examples/library-project
```

Use `go test ./...` from this directory to validate the consumer-facing package.
