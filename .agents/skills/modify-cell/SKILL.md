---
name: modify-cell
description: Modify an existing cell in the Agent-First Go architecture. Use when changing a cell's interface, logic, or dependencies. Covers blast-radius analysis, safe modification order, and validation.
---

# Modify a Cell (Safely)

## Before You Start

```bash
task context ID=<cell-id>      # read the cell's context pack
task deps ID=<cell-id>         # check who depends on this cell
task scope ID=<cell-id>        # declare the allowed edit boundary
```

The `task deps` output shows:

- **Dependencies**: what this cell needs (forward deps)
- **Dependents**: cells that depend on THIS cell (reverse deps)

If cells appear in "Dependents", changing the interface may break them.

## Modification Order

### If changing implementation only (no interface change):

1. **Read** `service.go`, `store.go` — understand current logic
2. **Modify** the implementation
3. **Update tests** — add/modify test cases
4. **Validate**: `task test-cell ID=<id> && task impact && task doctor`

### If changing the interface:

1. **Check dependents**: `task deps ID=<id>` — note all cells in "Dependents"
2. **Update the public contract** in `api/api.go`
3. **Update service** implementation in `service.go`
4. **Update all dependent cells** — they may need signature changes
5. **Update `cell.yaml`** if entrypoints changed
6. **Update tests** in this cell and dependent cells
7. **Validate everything**: `task index && task test && task impact && task doctor`

### If adding a new dependency:

1. **Import the target cell's `api` package** where the dependency is used
2. **Add the target cell's exact ID** to `cell.yaml` `dependencies`
3. **Add to `Deps` struct** in `<name>.go`
4. **Wire in `internal/app/wiring.go`**
5. **Validate**: `task index && task test-cell ID=<id> && task doctor`

### If adding a new method to the interface:

1. **Add method** to the interface in `api/api.go`
2. **Implement** in `service.go`
3. **Add handler** if HTTP-exposed (in `handler.go`)
4. **Add test** in `*_test.go`
5. **Update `cell.yaml`** if entrypoints changed
6. **Wire in `wiring.go`** if new deps needed
7. **Validate**: `task index && task test-cell ID=<id> && task impact && task doctor`

## After All Changes

```bash
task index          # regenerate index (hash changes)
task impact         # see full blast radius
task verify-scope ID=<cell-id> # reject undeclared edits
task test           # full test suite
task doctor         # architecture health check
```

## Rules

- **Never modify another cell's internal files.** Only use its interface.
- **Never modify `wiring.go` from within a cell.** Wiring is centralized.
- **Always run `task impact` after changes.** It catches under-tested blast radius.
- **Expand scope deliberately.** For a sibling cell or shared surface, pass
  only the needed exact IDs or `@contracts`, `@platform`, and `@wiring` through
  `WITH=`; dependencies and dependents are never implicit authorization.
- **Always run `task doctor` before committing.** It catches architecture violations.
