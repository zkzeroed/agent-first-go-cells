<!-- source-hash: sha256:ad0c1a5a09b2e6e595eb853608ad79298d51a3df067999c61141e9541fb68527 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# greeting-render

**Purpose:** Render greetings over HTTP

**Kind:** private cell

**Public:** false

**Entrypoints:**

- `api/api.go`

**Dependencies:**

- greeting-compose

**Invariants:**

- Invalid names return HTTP 400

**Validation:**

- `go test ./internal/cells/greeting-render/...`

## Cell Guide

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
