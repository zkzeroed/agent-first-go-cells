<!-- source-hash: sha256:3f3b915a2f2c4a4a0a6771a5e1b7e98d7f69039ce1fc8cf0a05fd67c473f1cb3 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# token-issue

**Purpose:** Issue random opaque tokens and retain only their SHA-256 digests

**Entrypoints:**

- `api/api.go`

**Invariants:**

- Raw tokens are never written to the store
- Empty subjects are rejected before random bytes are generated

**Validation:**

- `go test ./internal/cells/token-issue/...`

## Cell Guide

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
