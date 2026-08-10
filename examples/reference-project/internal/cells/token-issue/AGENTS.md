# Cell: token-issue

## Purpose

Issues cryptographically random opaque tokens and stores only their SHA-256 digests.

## Start Here

Read `api/api.go`, then `service.go` and `service_test.go`.

## Invariants

- Raw tokens are never written to the store.
- Empty subjects are rejected before random bytes are generated.

## Common Tasks

Keep public token vocabulary in `api/api.go`; preserve cryptographic randomness and digest-only storage.

## Reliability & Concurrency

The service is immutable after construction. The store must be safe for concurrent use in production.

## Validation

- `go test ./internal/cells/token-issue/...`
