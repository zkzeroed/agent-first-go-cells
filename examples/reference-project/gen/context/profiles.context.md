<!-- source-hash: sha256:76fc702873be4a6030a1e72715f0c5f40cc53b94ca128502e9738a62a7da34a3 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# profiles

**Purpose:** Own shared profile vocabulary and repository contracts

**Entrypoints:**

- `api/api.go`

**Validation:**

- `go test ./internal/cells/profiles/...`

## Cell Guide

# Cell: profiles

## Purpose

Owns shared profile vocabulary and repository contracts for profile actions.

## Start Here

Read `api/api.go`, then `cell.go`.

## Invariants

- Repository access is concurrent-safe.

## Common Tasks

Add behavior as a sub-action that depends on `profiles` and imports only `profiles/api`.

## Reliability & Concurrency

The in-memory repository uses a read/write mutex; production adapters retain the same contract.

## Validation

- `go test ./internal/cells/profiles/...`
