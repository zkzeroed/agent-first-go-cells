# Cell: greeting-render

## Purpose

Renders greeting messages using the private text normalization contract.

## Start Here

Read `api/api.go`, `greeting_render.go`, then `text-normalize/api`.

## Invariants

- Empty normalized names are rejected.

## Common Tasks

Import other private cells only through their `api` package.

## Reliability & Concurrency

The renderer is immutable after construction and safe for concurrent use.

## Validation

- `go test ./internal/cells/greeting-render/...`
