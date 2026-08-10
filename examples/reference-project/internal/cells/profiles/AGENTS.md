# Cell: profiles

## Purpose

Owns shared profile vocabulary and repository contracts for profile actions.

## Start Here

Read `api/api.go`, then `cell.go`.

## Invariants

- Repository access is concurrent-safe.

## Common Tasks

Add behavior as a sub-action that depends on `profiles` and imports only `profiles/api`.

## Reliability & Concurrency

The in-memory repository uses a read/write mutex; production adapters retain the same contract.

## Validation

- `go test ./internal/cells/profiles/...`
