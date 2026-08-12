#!/usr/bin/env bash
# new-cell.sh — Scaffold a new flat cell.
#
# Usage: task new-cell ID=user-authenticate
#
# Creates the lean cell directory required by the Agent-First Go architecture.
# Pass --extended to add application-oriented types, errors, service, store,
# and handler files. Each file contains TODO markers for the next developer.
#
# WHEN TO USE: When creating a new capability that has one primary behavior.
# For capabilities with multiple distinct behaviors sharing types, use
# `task scaffold-domain` instead.
#
# WHY IT HELPS AGENTS: An agent can start with the smallest correct ownership
# boundary, then opt into the broader application template when needed.
set -euo pipefail

root="."
extended=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --extended) extended=true; shift ;;
    --root) root="${2:?usage: new-cell.sh [--extended] --root <project-root> <cell-id>}"; shift 2 ;;
    *) break ;;
  esac
done
id="${1:?usage: new-cell.sh [--root <project-root>] <cell-id>}"
if ! [[ "${id}" =~ ^[a-z][a-z0-9]*(-[a-z0-9]+)*(/[a-z][a-z0-9]*(-[a-z0-9]+)*)?$ ]]; then
  echo "invalid cell ID: ${id} (use kebab-case, optionally domain/action)" >&2
  exit 1
fi
cd "$root"
path="internal/cells/${id}"
name=$(basename "${id}" | tr -d '-')
if [[ -e "${path}" ]]; then
  echo "cell path already exists: ${path}" >&2
  exit 1
fi
mkdir -p "$(dirname "${path}")"
mkdir "${path}"
mkdir -p "${path}/api"

cat > "${path}/cell.yaml" <<EOF
id: ${id}
purpose: "TODO: describe what this cell does"
entrypoints:
  - file: api/api.go
    symbol: Service
dependencies: []
validation:
  - go test ./internal/cells/${id}/...
EOF

cat > "${path}/doc.go" <<EOF
// Package ${name} implements the ${id} cell.
//
// Cell ID: ${id}
package ${name}
EOF

cat > "${path}/api/api.go" <<EOF
// Package api defines the public contract for the ${id} cell.
package api

// Service is the behavior other cells may depend on.
// TODO: Define the public methods and types for this cell.
type Service = any
EOF

cat > "${path}/${name}.go" <<EOF
package ${name}

// TODO: Implement the cell. Keep public contracts in api/api.go.
EOF

cat > "${path}/${name}_test.go" <<EOF
package ${name}

import "testing"

func TestTODO(t *testing.T) {
    // TODO: Add table-driven tests.
}
EOF

if [[ "$extended" == true ]]; then
rm "${path}/${name}_test.go"

cat > "${path}/types.go" <<EOF
package ${name}

// TODO: Define domain types.
EOF

cat > "${path}/errors.go" <<EOF
package ${name}

import "errors"

// TODO: Define sentinel errors.
var ErrTODO = errors.New("${id}: TODO")
EOF

cat > "${path}/service.go" <<EOF
package ${name}

// TODO: Implement business logic.
EOF

cat > "${path}/store.go" <<EOF
package ${name}

// TODO: Define Store interface and implementation.
EOF

cat > "${path}/handler.go" <<EOF
package ${name}

// TODO: Implement HTTP handler.
EOF

cat > "${path}/service_test.go" <<EOF
package ${name}

import "testing"

func TestTODO(t *testing.T) {
	// TODO: Add table-driven tests.
}
EOF
fi

cat > "${path}/AGENTS.md" <<EOF
# Cell: ${id}

## Purpose

TODO: Describe what this cell does.

## Start Here

1. \`cell.yaml\` — metadata, dependencies, invariants
2. \`api/api.go\` — public interface and shared contract types
3. \`${name}.go\` — implementation
4. \`${name}_test.go\` — expected behavior

## Invariants

- TODO: Add invariants

## Common Tasks

### Extend this cell
1. Add method to the public interface in \`api/api.go\`
2. Implement in \`${name}.go\`
3. Add a table-driven test
4. Update \`cell.yaml\` if the interface changed
5. Run \`go test ./internal/cells/${id}/...\`

## Reliability & Concurrency

- TODO: Document concurrency safety, context handling, error taxonomy, retry/idempotency.

## Validation

- \`go test ./internal/cells/${id}/...\`
- \`golangci-lint run ./internal/cells/${id}/...\`
EOF

echo "Created $(if [[ "$extended" == true ]]; then echo "extended "; fi)cell ${id} at ${path}"
echo "Next: fill in the TODOs, wire in internal/app/wiring.go, run 'task index && task test'"
