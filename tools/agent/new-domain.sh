#!/usr/bin/env bash
# new-domain.sh — Scaffold a new domain cell with sub-action support.
#
# Usage: task scaffold-domain ID=users
#
# Creates a domain cell directory with cell.go, errors.go, model.go, doc.go,
# AGENTS.md, and cell.yaml. Does NOT create sub-actions — those are added
# individually via `task new-cell ID=<domain>/<action>`.
#
# WHEN TO USE: When a capability has multiple distinct behaviors with shared
# types and error taxonomy (e.g., `users` — invite, register, profile).
# For simple capabilities with one primary behavior, use `task new-cell` instead.
#
# WHY IT HELPS AGENTS: An agent can scaffold a domain cell with a single
# command, getting the correct domain-level file schema. Sub-actions are then
# added individually, maintaining the same fixed schema at the leaf level.
set -euo pipefail

root="."
if [[ "${1:-}" == "--root" ]]; then
  root="${2:?usage: new-domain.sh --root <project-root> <domain-id>}"
  shift 2
fi
id="${1:?usage: new-domain.sh [--root <project-root>] <domain-id>}"
if ! [[ "${id}" =~ ^[a-z][a-z0-9]*(-[a-z0-9]+)*$ ]]; then
  echo "invalid domain ID: ${id} (use kebab-case)" >&2
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
purpose: "TODO: describe what this domain cell does"
entrypoints:
  - file: api/api.go
    symbol: Cell
dependencies: []
validation:
  - go test ./internal/cells/${id}/...
EOF

cat > "${path}/doc.go" <<EOF
// Package ${name} implements the ${id} domain cell.
//
// Cell ID: ${id}
// This is a domain cell. Sub-actions are in subdirectories.
package ${name}
EOF

cat > "${path}/api/api.go" <<EOF
// Package api defines the public contract for the ${id} domain cell.
package api

// Cell is the public behavior exposed by this domain.
// TODO: Define domain-level behavior shared by sub-actions.
type Cell = any
EOF

cat > "${path}/cell.go" <<EOF
package ${name}

// TODO: Implement NewCell constructor and aggregate sub-action handlers.

type Deps struct {
	// TODO: Add sub-action dependencies.
}

type cell struct {
	// TODO: Add sub-action fields.
}
EOF

cat > "${path}/errors.go" <<EOF
package ${name}

import "errors"

// Shared error taxonomy for all sub-actions in this domain.
var (
	ErrTODO = errors.New("${id}: TODO")
)

func IsTODO(err error) bool { return errors.Is(err, ErrTODO) }
EOF

cat > "${path}/model.go" <<EOF
package ${name}

// Shared domain types for all sub-actions in this domain.
// TODO: Define shared types.
EOF

cat > "${path}/AGENTS.md" <<EOF
# Cell: ${id} (domain)

## Purpose

TODO: Describe what this domain cell does.

## Start Here

1. \`cell.yaml\` — metadata, dependencies, invariants
2. \`api/api.go\` — domain public interface and shared contract types
3. \`cell.go\` — domain implementation and constructor
4. \`model.go\` — shared domain types
5. \`errors.go\` — shared error taxonomy

## Invariants

- TODO: Add invariants

## Common Tasks

### Add a new sub-action
1. Run \`task new-cell ID=${id}/<action-name>\`
2. Implement the sub-action following the flat cell schema
3. Add public behavior to \`api/api.go\`
4. Wire sub-action in \`cell.go\` NewCell constructor
5. Wire domain cell in \`internal/app/wiring.go\`
6. Run \`go test ./internal/cells/${id}/...\`

## Reliability & Concurrency

- TODO: Document concurrency safety, context handling, error taxonomy, retry/idempotency.

## Validation

- \`go test ./internal/cells/${id}/...\`
- \`golangci-lint run ./internal/cells/${id}/...\`
EOF

echo "Created domain cell ${id} at ${path}"
echo "Next: add sub-actions with 'task new-cell ID=${id}/<action-name>'"
echo "Then: fill in TODOs, wire in internal/app/wiring.go, run 'task index && task test'"
