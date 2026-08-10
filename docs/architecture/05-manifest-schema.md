# 5. Manifest Schema

## 5.1 `cell.yaml` — 5 Mandatory Fields + 1 Optional

```yaml
# cell.yaml — Machine-readable manifest for user-authenticate cell
id: user-authenticate
purpose: "Authenticate users via login, manage sessions, validate tokens"
entrypoints:
  - file: api/api.go
    symbol: Authenticator
dependencies: []
validation:
  - go test ./internal/cells/user-authenticate/...
  - golangci-lint run ./internal/cells/user-authenticate/...

# Optional (Standard/Full tier):
invariants:
  - "Tokens are generated exactly once"
  - "Expired tokens are rejected before any side effect"
  - "Session IDs are cryptographically random"
```

## 5.2 Field Rules

- **`id`** (required): Behavior-first kebab-case. Must match directory name.
- **`purpose`** (required): One-sentence description of what the cell does.
- **`entrypoints`** (required): List of `{file, symbol}` pairs used as
  orientation metadata. Each must name a Go file within its cell and a
  top-level symbol declared in that file.
- **`dependencies`** (required): List of exact IDs of other cells. Every value must resolve to one existing cell exactly once and match a direct import of that cell's `api` package; packages, interfaces, and types do not belong here.
- **`validation`** (required): List of commands to validate the cell. First command should be fast (`go test`). `task validate-cell ID=<id>` reads and runs these.
- **`invariants`** (optional, Standard+): List of properties that must always hold. High-value for agents — tells them what must not be broken.

## 5.3 Validation

`task check-manifests` validates:

1. YAML has exactly one document, no unknown fields, and all 5 mandatory fields.
2. `id` uses the supported kebab-case form and matches its directory name.
3. Every entrypoint names a contained Go file and a symbol declared in it, and the manifest has at least one validation command.
4. Cell-to-cell dependency IDs resolve exactly once, have no self-dependencies or duplicates, and match direct cell API imports.
