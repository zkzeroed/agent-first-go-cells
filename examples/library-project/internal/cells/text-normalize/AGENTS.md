# Cell: text-normalize

## Purpose

Owns name normalization used by the library's private composition cells.

## Start Here

Read `api/api.go`, then `text_normalize.go` and tests.

## Invariants

- Normalized names are trimmed and title cased by words.

## Common Tasks

Keep cross-cell behavior in `api/api.go`.

## Reliability & Concurrency

The normalizer is immutable and safe for concurrent use.

## Validation

- `go test ./internal/cells/text-normalize/...`
