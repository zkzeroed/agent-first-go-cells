<!-- source-hash: sha256:e028d5f052982f26c7e904dcdc29eec6f7ae06f8ad7bf0d7f81f129382f5f0b7 -->
<!-- generated-from: cell.yaml, AGENTS.md -->

# greeting-render

**Purpose:** Render normalized greeting messages for the public library package

**Kind:** private cell

**Public:** false

**Entrypoints:**

- `api/api.go`

**Dependencies:**

- text-normalize

**Invariants:**

- Empty normalized names are rejected

**Validation:**

- `go test ./internal/cells/greeting-render/...`

## Cell Guide

# Cell: greeting-render

## Purpose

Renders greeting messages using the private text normalization contract.

## Start Here

Read `api/api.go`, `greeting_render.go`, then `text-normalize/api`.

## Invariants

- Empty normalized names are rejected.

## Common Tasks

Import other private cells only through their `api` package.

## Reliability & Concurrency

The renderer is immutable after construction and safe for concurrent use.

## Validation

- `go test ./internal/cells/greeting-render/...`
