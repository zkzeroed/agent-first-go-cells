# Cell: modulus

## Purpose

Expose canonical residues modulo the fixed prime 97.

## Start Here

Read `cell.yaml`, `docs/research/modular-reduction.md`, `api.go`, and tests.

## Invariants

- Keep exported API stable and do not expose private implementation types.
- Preserve the canonical range `0 <= Reduce(value) < Prime`.

## Common Tasks

Keep the cited Euclidean-remainder behavior and the private helper dependency
declared in `cell.yaml`.

## Reliability & Concurrency

Document concurrency and error behavior for exported operations.

## Validation

- `go test ./modulus/...`
