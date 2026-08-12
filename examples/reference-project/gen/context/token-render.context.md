<!-- source-hash: sha256:c58f59437db4fe6fe295514577b98ee098ad5dffd5f2eb4630ad690c8800b51a -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# token-render

**Purpose:** Render newly issued opaque tokens over HTTP

**Kind:** private cell

**Public:** false

**Entrypoints:**

- `api/api.go`

**Dependencies:**

- token-issue

**Invariants:**

- Empty subjects return HTTP 400

**Validation:**

- `go test ./internal/cells/token-render/...`

## Cell Guide

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
