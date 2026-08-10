# 7. Impact Analysis

## 7.1 `scripts/impact/`

Go-based impact analyzer. More maintainable and testable than shell scripts.

**Algorithm:**

1. Read tracked changes plus untracked files to get changed files.
2. Map each file to its owning cell (by directory path).
3. Read `cell.yaml` manifests to find cells with an exact dependency on the changed cell ID.
4. Print: changed files → owning cells → affected cells → required validation commands.

**Output example:**

```
=== Impact Analysis ===

Changed files:
  internal/cells/user-authenticate/service.go
  internal/contracts/clock.go

Owning cells:
  user-authenticate
  contracts (shared)

Affected cells (depend on changed cells):
  user-authenticate → user-invite (declared in cell.yaml)
  user-authenticate → user-register (declared in cell.yaml)

Validation commands to run:
  go test ./internal/cells/user-authenticate/...
  go test ./internal/cells/users/user-invite/...
  go test ./internal/cells/users/user-register/...
  go test ./internal/contracts/...
```

## 7.2 Usage

```bash
task impact  # Show blast-radius for changes from HEAD plus untracked files
```

Impact is Git-repository scoped and considers changed `.go`, `.yaml`, and
`.yml` files. Use `ROOT=<project-path>` to scope analysis to a compatible
project under the same Git repository.
