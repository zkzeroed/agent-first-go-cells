# Cell: token-render

## Purpose

Exposes token issuance over HTTP.

## Start Here

Read `api/api.go`, `handler.go`, and `token-issue/api`.

## Invariants

- Empty subjects return HTTP 400.

## Common Tasks

Import only `token-issue/api` for token issuance behavior.

## Reliability & Concurrency

The handler is immutable after construction and delegates token behavior to its dependency.

## Validation

- `go test ./internal/cells/token-render/...`
