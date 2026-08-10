# Cell: profiles/profile-create

## Purpose

Creates validated profiles through the profiles domain contract.

## Start Here

Read `api/api.go`, `profile_create.go`, and `profiles/api`.

## Invariants

- Empty profile identifiers and names are rejected.

## Common Tasks

Declare `profiles` as the exact dependency and import only `profiles/api`.

## Reliability & Concurrency

The service is immutable after construction; repository concurrency is owned by the domain adapter.

## Validation

- `go test ./internal/cells/profiles/profile-create/...`
