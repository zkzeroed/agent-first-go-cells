<!-- source-hash: sha256:562466478781f49e4007d6aa243f17d2d39b36d8b257314bb5ba8319520b9153 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# profiles/profile-create

**Purpose:** Create validated profiles using the shared repository contract

**Entrypoints:**

- `api/api.go`

**Dependencies:**

- profiles

**Invariants:**

- Empty profile identifiers and names are rejected

**Validation:**

- `go test ./internal/cells/profiles/profile-create/...`

## Cell Guide

# Cell: profiles/profile-create

## Purpose

Creates validated profiles through the profiles domain contract.

## Start Here

Read `api/api.go`, `profile_create.go`, and `profiles/api`.

## Invariants

- Empty profile identifiers and names are rejected.

## Common Tasks

Declare `profiles` as the exact dependency and import only `profiles/api`.

## Reliability & Concurrency

The service is immutable after construction; repository concurrency is owned by the domain adapter.

## Validation

- `go test ./internal/cells/profiles/profile-create/...`
