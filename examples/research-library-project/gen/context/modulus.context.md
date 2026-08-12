<!-- source-hash: sha256:42beedc1b38e05aa0868b7a968bf0ae94ac1ccea6098f451b18521e6cc0ed7c1 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# modulus

**Purpose:** Reduce signed integers modulo the fixed prime 97

**Kind:** library-package

**Public:** true

**Conformance:**

- Basis: paper-defined-math
- Status: conformant
- Evidence: verified
- Citation: `docs/research/modular-reduction.md`, heading "Euclidean remainder"; symbols: Reduce

**Entrypoints:**

- `api.go`

**Dependencies:**

- modulus-reduce

**Validation:**

- `go test ./modulus/...`

## Cell Guide

# Cell: modulus

## Purpose

Expose canonical residues modulo the fixed prime 97.

## Start Here

Read `cell.yaml`, `docs/research/modular-reduction.md`, `api.go`, and tests.

## Invariants

- Keep exported API stable and do not expose private implementation types.
- Preserve the canonical range `0 <= Reduce(value) < Prime`.

## Common Tasks

Keep the cited Euclidean-remainder behavior and the private helper dependency
declared in `cell.yaml`.

## Reliability & Concurrency

Document concurrency and error behavior for exported operations.

## Validation

- `go test ./modulus/...`
