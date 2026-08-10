#!/usr/bin/env bash
# new-cell.sh — Scaffold a new flat cell.
#
# Usage: task new-cell ID=user-authenticate
#
# Creates a complete cell directory with all required files following the
# Agent-First Go architecture file schema. Each file contains TODO markers
# for the agent or developer to fill in.
#
# WHEN TO USE: When creating a new capability that has one primary behavior.
# For capabilities with multiple distinct behaviors sharing types, use
# `task scaffold-domain` instead.
#
# WHY IT HELPS AGENTS: An agent can scaffold a new cell with a single command,
# getting the correct file schema automatically. The TODO markers guide the
# agent on what to implement, reducing the chance of missing required files.
set -euo pipefail

id="${1:?usage: new-cell.sh <cell-id>}"
if ! [[ "${id}" =~ ^[a-z][a-z0-9]*(-[a-z0-9]+)*(/[a-z][a-z0-9]*(-[a-z0-9]+)*)?$ ]]; then
  echo "invalid cell ID: ${id} (use kebab-case, optionally domain/action)" >&2
  exit 1
fi
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

cat > "${path}/AGENTS.md" <<EOF
# Cell: ${id}

## Purpose

TODO: Describe what this cell does.

## Start Here

1. \`cell.yaml\` — metadata, dependencies, invariants
2. \`api/api.go\` — public interface and shared contract types
3. \`service.go\` — business logic
4. \`service_test.go\` — expected behavior

## Invariants

- TODO: Add invariants

## Common Tasks

### Add a new endpoint
1. Add method to the public interface in \`api/api.go\`
2. Implement in \`service.go\`
3. Add handler case in \`handler.go\`
4. Add table-driven test
5. Update \`cell.yaml\` if interface changed
6. Run \`go test ./internal/cells/${id}/...\`

## Reliability & Concurrency

- TODO: Document concurrency safety, context handling, error taxonomy, retry/idempotency.

## Validation

- \`go test ./internal/cells/${id}/...\`
- \`golangci-lint run ./internal/cells/${id}/...\`
EOF

echo "Created cell ${id} at ${path}"
echo "Next: fill in the TODOs, wire in internal/app/wiring.go, run 'task index && task test'"
