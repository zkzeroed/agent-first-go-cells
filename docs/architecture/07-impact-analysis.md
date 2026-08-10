# 7. Impact Analysis

## 7.1 `tools/agent/impact/`

Go-based impact analyzer. More maintainable and testable than shell scripts.

**Algorithm:**

1. Read tracked changes plus untracked files to get changed files. In addition to
   Go and YAML files, include a changed `AGENTS.md` when it is local to a cell.
2. Map each file to its most-specific owning cell using directory-segment
   boundaries, so a change under `orders/create/` is owned by that nested cell,
   not its `orders/` parent.
3. Read `cell.yaml` manifests to find all direct and transitive downstream
   dependents of the changed cell IDs. Affected cell IDs are sorted for stable
   output.
4. Classify changes in `internal/contracts/`, `internal/platform/`, and
   `internal/app/wiring*.go` as shared surfaces. A shared-surface change affects
   every cell and requires full-project validation.
5. When a changed `cell.yaml` has been removed, read that manifest from `HEAD`
   so its ownership and dependents remain visible.
6. Print: changed files → owning cells → shared surfaces → downstream affected
   cells → required validation commands. JSON output includes additive
   `sharedSurfaces` and `fullProjectValidation` fields.

**Output example:**

```
=== Impact Analysis ===

Changed files:
  internal/cells/user-authenticate/service.go
  internal/contracts/clock.go

Owning cells:
  user-authenticate
  contracts (shared)

Shared surfaces:
  contracts

Affected cells (depend on changed cells):
  user-authenticate → user-invite (declared in cell.yaml)
  user-authenticate → user-register (declared in cell.yaml)

Validation commands to run:
  go test ./...
```

## 7.2 Usage

```bash
task impact  # Show blast-radius for changes from HEAD plus untracked files
```

Impact is Git-repository scoped and considers changed `.go`, `.yaml`, and
`.yml` files plus cell-local `AGENTS.md`. Every changed file under the shared
contracts, platform, and wiring surfaces is also included. Use
`ROOT=<project-path>` to scope analysis to a compatible project under the same
Git repository.
