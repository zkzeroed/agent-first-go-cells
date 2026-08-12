# Architecture Principles

This bootstrap optimizes a Go modular monolith for coding agents while keeping
the result straightforward for human developers.

## Design goals

- **Locality:** Everything needed to understand a capability lives in one cell
  directory, so an agent does not need to reconstruct a feature from layers.
- **Predictability:** Every cell follows the same schema and behavior-first
  naming, making search and scaffolding deterministic.
- **Explicitness:** Public APIs, dependencies, wiring, invariants, and
  validation are declared rather than inferred from convention.
- **Evidence:** Agent guidance is backed by manifests, generated context,
  import policy, structural tests, impact analysis, and scope verification.

## Core rules

1. A cell is one business capability, named in behavior-first kebab-case.
2. Files stay below 300 lines and functions below 40 lines unless the affected
   function documents `AGENT_OVERRIDE`; AST-based structural tests enforce
   both and `cyclop` has a maximum 16 complexity limit.
3. A cell exposes only `api/api.go`; other cells never import its implementation.
4. In an application, `cmd/<name>/main.go` owns process lifecycle and
   `internal/app/wiring*.go` is the sole explicit private-cell composition root.
   In a library, the registered public package is the composition root; it has
   no `cmd/` or application wiring. No `init()`, reflection DI, or service
   locators.
5. `cell.yaml` declares exact dependencies and validation. Generated metadata
   is derived from it and is never edited manually.

## Trade-offs

The fixed schema creates more files and explicit manifests, while centralized
wiring requires deliberate integration. In exchange, agents and developers get
bounded context, visible dependency direction, predictable scaffolding, and
deterministic checks. Use domain sub-actions for closely related behaviors; do
not create tiny cells merely to maximize separation.
