# Cell: text-fold

## Purpose

TODO: Describe what this cell does.

## Start Here

1. `cell.yaml` — metadata, dependencies, invariants
2. `api/api.go` — public interface and shared contract types
3. `textfold.go` — implementation
4. `textfold_test.go` — expected behavior

## Invariants

- TODO: Add invariants

## Common Tasks

### Extend this cell
1. Add method to the public interface in `api/api.go`
2. Implement in `textfold.go`
3. Add a table-driven test
4. Update `cell.yaml` if the interface changed
5. Run `go test ./internal/cells/text-fold/...`

## Reliability & Concurrency

- TODO: Document concurrency safety, context handling, error taxonomy, retry/idempotency.

## Validation

- `go test ./internal/cells/text-fold/...`
- `golangci-lint run ./internal/cells/text-fold/...`
