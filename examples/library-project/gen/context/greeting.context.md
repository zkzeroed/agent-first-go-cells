<!-- source-hash: sha256:93d78bfc60c7baa9cda48cb1a586c9eac701007c617fbdd058b08aeb04ab01ef -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# greeting

**Purpose:** Expose stable greeting messages to downstream Go consumers

**Kind:** library-package

**Public:** true

**Conformance:**

- Basis: engineering-primitive
- Status: conformant
- Evidence: verified
- Rationale: This example demonstrates a stable exported package composed from private cells; it does not implement a paper-defined protocol.

**Entrypoints:**

- `greeting.go`

**Dependencies:**

- greeting-render
- text-normalize

**Invariants:**

- Exported API never exposes private cell types
- Whitespace-only names are rejected

**Validation:**

- `go test ./greeting/...`

## Cell Guide

# Cell: greeting

## Purpose

Provides the stable, importable greeting package for downstream consumers.

## Start Here

Read `greeting.go`, then the `greeting-render` private cell contract.

## Invariants

- Exported API never exposes private implementation types.
- Whitespace-only names return `ErrInvalidName`.

## Common Tasks

Keep private-cell dependencies declared in `cell.yaml` and translate their
errors at this public boundary when necessary.

## Reliability & Concurrency

The greeter is immutable after construction and safe for concurrent use.

## Validation

- `go test ./greeting/...`
