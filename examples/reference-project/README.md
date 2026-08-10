# Reference Project

This runnable example demonstrates six cells, including a nested profiles domain,
explicit `api` boundaries,
exact manifest dependencies, tests, and composition in `cmd/reference/main.go`. The token flow uses
`crypto/rand` for entropy and `crypto/sha256` for digest-only storage.

```bash
cd examples/reference-project
go test ./...
go run ./cmd/reference
curl http://localhost:8080/greeting/Ada
curl http://localhost:8080/token/ada
```

The greeting response is `Hello, Ada!`; the token response is JSON containing
an opaque token and its SHA-256 digest.

`greeting-render` imports only `greeting-compose/api`; its manifest declares
the exact `greeting-compose` dependency. `token-render` follows the same
pattern with `token-issue/api`; `greeting-compose` and `token-issue` have no
cell dependencies. `profiles/profile-create` depends on the parent `profiles`
cell and imports only `profiles/api`. This directory is reference material, not
starter application code to delete.
