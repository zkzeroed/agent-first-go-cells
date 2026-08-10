<!-- source-hash: sha256:a880d838358059b8ce4fb90aa7a23234ba93c0fc39503dbd71f9627fb50e5a02 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# greeting-compose

**Purpose:** Compose a greeting from a validated name

**Entrypoints:**

- `api/api.go`

**Invariants:**

- Empty names are rejected before store access

**Validation:**

- `go test ./internal/cells/greeting-compose/...`

## Cell Guide

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
