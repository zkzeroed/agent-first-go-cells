<!-- source-hash: sha256:6c8cca358162103048210fee11773282a604364a60e306fc6e4cea6c723fd7a0 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# text-fold

**Purpose:** Fold whitespace and lowercase labels for module-root package composition

**Kind:** private cell

**Public:** false

**Entrypoints:**

- `api/api.go`

**Validation:**

- `go test ./internal/cells/text-fold/...`

## Cell Guide

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
