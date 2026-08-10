# 6. Generated Artifacts

`task index` generates the following files from every `cell.yaml` and adjacent
`AGENTS.md`; never edit them by hand.

## `gen/cells.json`

The index has `schemaVersion`, `hash`, and `cells`. Each cell record contains
its ID, path, Go package name, purpose, entrypoint files, exact cell-ID
dependencies, and validation commands. Its hash covers the manifest and guide
for every cell, so `task check-index` detects stale metadata.

## `gen/context/<id>.context.md`

Each bounded context pack contains the manifest's purpose, entrypoints,
dependencies, invariants, validation commands, and the cell guide. Its
`source-hash` covers that cell's `cell.yaml` and `AGENTS.md`.

`task index`, `task check-index`, and `task index-json` accept
`ROOT=<project-path>` for a compatible project such as
`examples/reference-project`.
