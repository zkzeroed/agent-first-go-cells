# Cell: canonical

## Purpose

Provide a small, stable module-root API for canonical labels.

## Start Here

Read `cell.yaml`, `api.go`, and the `text-fold` helper.

## Invariants

- Keep exported API stable and do not expose private implementation types.
- Folded labels are lowercase and have no repeated whitespace.

## Common Tasks

Declare every private helper dependency in `cell.yaml`.

## Reliability & Concurrency

Document concurrency and error behavior for exported operations.

## Validation

- `go test ./...`
