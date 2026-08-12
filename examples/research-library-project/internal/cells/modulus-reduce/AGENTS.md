# Cell: modulus-reduce

## Purpose

TODO: Describe what this cell does.

## Start Here

1. `cell.yaml` — metadata, dependencies, invariants
2. `api/api.go` — public interface and shared contract types
3. `service.go` — business logic
4. `service_test.go` — expected behavior

## Invariants

- TODO: Add invariants

## Common Tasks

### Add a new endpoint
1. Add method to the public interface in `api/api.go`
2. Implement in `service.go`
3. Add handler case in `handler.go`
4. Add table-driven test
5. Update `cell.yaml` if interface changed
6. Run `go test ./internal/cells/modulus-reduce/...`

## Reliability & Concurrency

- TODO: Document concurrency safety, context handling, error taxonomy, retry/idempotency.

## Validation

- `go test ./internal/cells/modulus-reduce/...`
- `golangci-lint run ./internal/cells/modulus-reduce/...`
