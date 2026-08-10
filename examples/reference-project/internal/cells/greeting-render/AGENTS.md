# Cell: greeting-render

## Purpose

Exposes greeting composition over HTTP.

## Start Here

Read `api/api.go`, `handler.go`, and the greeting-compose API contract.

## Invariants

- Invalid names return HTTP 400.

## Common Tasks

Import only `greeting-compose/api` for cross-cell behavior.

## Reliability & Concurrency

The handler is immutable after construction and delegates all business logic.

## Validation

- `go test ./internal/cells/greeting-render/...`
