---
name: build-cell
description: Build a new cell in the Agent-First Go architecture. Use when creating a new capability, scaffolding a cell, or implementing a feature from scratch. Covers the inside-out build order, file schema, and validation steps.
---

# Build a Cell (Inside-Out)

## Build Order

Cells are built **inside-out** — from domain vocabulary to transport layer. This mirrors hexagonal architecture's "domain outwards to adapters" but with cell-specific ordering.

### Step 1: Types (`types.go`)

Define domain types first. These are the vocabulary everything else uses.

```go
package userauthenticate

type Request struct {
    Email    string
    Password string
}

type Result struct {
    UserID  string
    Token   string
}
```

### Step 2: Errors (`errors.go`)

Define the error taxonomy. Sentinel errors + typed errors.

```go
package userauthenticate

import "errors"

var (
    ErrInvalidCredentials = errors.New("user-authenticate: invalid credentials")
    ErrUserNotFound       = errors.New("user-authenticate: user not found")
)
```

### Step 3: Public API (`api/api.go`)

Define the contract — the cell's membrane. This is what other cells see.

```go
package api

import "context"

type Request struct{ Email string }
type Result struct{ UserID string }

type Authenticator interface {
    Authenticate(ctx context.Context, req Request) (Result, error)
}

// Keep implementation constructors in the parent cell package.
```

### Step 4: Store (`store.go`)

Define the Store interface and implementation. Data access behind an interface.

```go
package userauthenticate

import "context"

type Store interface {
    FindByEmail(ctx context.Context, email string) (User, error)
}

type User struct {
    ID           string
    Email        string
    PasswordHash string
}
```

### Step 5: Service (`service.go`)

Implement business logic. Uses types + store interface. Must be safe for concurrent use.

```go
package userauthenticate

import "context"

type service struct {
    deps Deps
}

func (s *service) Authenticate(ctx context.Context, req Request) (Result, error) {
    // business logic here
    return Result{}, nil
}
```

### Step 6: Handler (`handler.go`)

Wire HTTP/transport to the service. Keep it thin — parse, delegate, respond.

```go
package userauthenticate

import (
    "encoding/json"
    "net/http"
)

type Handler struct{ svc api.Authenticator }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // parse request, call service, write response
}
```

### Step 7: Tests (`*_test.go`)

Write table-driven tests. Use `t.Context()` (Go 1.24+). Test the service, not the handler.

### Step 8: Manifest (`cell.yaml`)

Declare metadata, dependencies, and validation commands.

Each entrypoint must name a Go file within this cell and a top-level symbol it
declares. Every listed dependency must correspond to a direct import of that
cell's `api` package; complete the source and manifest together before running
`task index`.

```yaml
id: user-authenticate
purpose: "Authenticate users by email and password"
entrypoints:
  - file: api/api.go
    symbol: Authenticator
dependencies: []
validation:
  - go test ./internal/cells/user-authenticate/...
invariants:
  - "Authenticate never returns a token for invalid credentials"
```

### Step 9: Agent Guide (`AGENTS.md`)

Document the cell for the next agent. Use the template from `task new-cell`.

### Step 10: Wiring (`internal/app/wiring.go`)

Connect the cell to the application. This is the ONLY place cells are constructed.

```go
// in wiring.go
userAuthStore := userauthenticate.NewStore(db)
userAuth := userauthenticate.New(userauthenticate.Deps{Store: userAuthStore})
```

### Step 11: Validate

```bash
task index          # regenerate cell index
task test-cell ID=user-authenticate
task impact         # check blast radius
task doctor         # architecture health check
```

## Key Rules

- **Never skip steps 1-3.** Public API, types, and errors define the cell's identity.
- **Store is always an interface.** This enables testing without a database.
- **Service is stateless** except for `deps` (read-only after construction).
- **Handler is thin.** No business logic in handlers.
- **Wiring is centralized.** Never construct cells outside `wiring.go`.
- **Files ≤ 300 LOC, functions ≤ 40 LOC.** Enforced by structural tests.
