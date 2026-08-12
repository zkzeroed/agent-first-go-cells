---
name: new-package
description: Create an exportable Go library-package cell in an Agent-First Go project. Use when adding a public module-root or subdirectory package, registering it in library mode, or composing declared private cells behind a stable downstream API.
---

# New Public Package

Create a public package cell rather than a private application cell when
downstream Go modules must import it.

## Workflow

1. Inspect the target project: `task orient ROOT=<project-root>` and
   `task cells ROOT=<project-root>`.
2. Choose a public import path and ownership ID. Keep package-specific helpers
   as unexported code beside the API; use `internal/` for helpers consumers
   must not access.
3. Scaffold and register the package:

   ```bash
   task new-package ID=field PATH=field ROOT=<project-root>
   ```

   For the module root, specify the Go package name:

   ```bash
   task new-package ID=widget PATH=. PACKAGE=widget ROOT=<project-root>
   ```

4. Replace the scaffolded `conformance` placeholder before treating the package
   as complete. Use `paper-defined-math` with local citations (file, typed
   locator, symbols) for direct research implementations. Use
   `fixed-profile-policy` or `engineering-primitive` with a rationale for
   deliberate engineering choices. Non-conformant statuses require explicit
   gaps. A `conformant` package requires `evidence: verified`; use
   `task conformance ID=<id> ROOT=<project-root>` to resolve its evidence.
5. Define the stable exported API and its tests. Do not expose `internal` or
   private-cell types in exported signatures.
6. If the package composes private cells, import their implementation packages
   only from this configured library package and list every direct dependency in
   `cell.yaml`. Do not import `internal/app`, `internal/platform`, or
   `internal/contracts`. Private cells still import each other only through
   `api`.
7. Run `task index ROOT=<project-root>`, `task test-cell ID=<id>
   ROOT=<project-root>`, and `task ready ROOT=<project-root>`.

## Boundaries

- `kind: library-package` and `public: true` identify an exportable package.
- Every library package records validated research or engineering conformance.
- `policy/architecture.yaml` is the authoritative ID-to-path registry.
- `internal/cells/` remains for private application or helper cells.
- Public packages may use the module root or any non-`internal` relative path,
  including `pkg/`; that path is a deliberate downstream import-path choice.
