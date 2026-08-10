# Cell: greeting-compose

## Purpose

Creates greetings from a name and a prefix store.

## Start Here

Read `api/api.go`, then `service.go` and `service_test.go`.

## Invariants

- Empty names are rejected before store access.

## Common Tasks

Keep cross-cell behavior in `api/api.go`; keep storage behind `Store`.

## Reliability & Concurrency

The service is immutable after construction. Store calls accept context.

## Validation

- `go test ./internal/cells/greeting-compose/...`
