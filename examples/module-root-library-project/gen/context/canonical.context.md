<!-- source-hash: sha256:cfaf0814bae8399f0231777eb275d6603ef3974898f98c9f9505017df486fe99 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# canonical

**Purpose:** Canonicalize short user-facing labels for downstream consumers

**Kind:** library-package

**Public:** true

**Conformance:**

- Basis: engineering-primitive
- Status: conformant
- Evidence: verified
- Rationale: The package demonstrates a stable module-root API backed by a private normalization helper; its text-folding behavior is intentional library policy.

**Entrypoints:**

- `api.go`

**Dependencies:**

- text-fold

**Validation:**

- `go test ./...`

## Cell Guide

# Cell: canonical

## Purpose

Provide a small, stable module-root API for canonical labels.

## Start Here

Read `cell.yaml`, `api.go`, and the `text-fold` helper.

## Invariants

- Keep exported API stable and do not expose private implementation types.
- Folded labels are lowercase and have no repeated whitespace.

## Common Tasks

Declare every private helper dependency in `cell.yaml`.

## Reliability & Concurrency

Document concurrency and error behavior for exported operations.

## Validation

- `go test ./...`
