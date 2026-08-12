# Cell: greeting

## Purpose

Provides the stable, importable greeting package for downstream consumers.

## Start Here

Read `greeting.go`, then the `greeting-render` private cell contract.

## Invariants

- Exported API never exposes private implementation types.
- Whitespace-only names return `ErrInvalidName`.

## Common Tasks

Keep private-cell dependencies declared in `cell.yaml` and translate their
errors at this public boundary when necessary.

## Reliability & Concurrency

The greeter is immutable after construction and safe for concurrent use.

## Validation

- `go test ./greeting/...`
