<!-- source-hash: sha256:bec2693e33f1ff94b79c5338262e49ad130afdfd0751027a0fbd68be4f64394b -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# text-normalize

**Purpose:** Normalize display names for private library composition

**Kind:** private cell

**Public:** false

**Entrypoints:**

- `api/api.go`

**Invariants:**

- Normalized names are trimmed and title cased by words

**Validation:**

- `go test ./internal/cells/text-normalize/...`

## Cell Guide

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
