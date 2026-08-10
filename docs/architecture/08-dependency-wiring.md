# 8. Dependency & Wiring Rules

## 8.1 Cell-to-Cell Imports

Cells may import only another cell's explicit **`api` package**. This is a real
Go package boundary, so implementation symbols are not importable by accident:

```go
// ALLOWED: importing the public cell contract
import "myproject/internal/cells/user-authenticate/api"

// FORBIDDEN: importing a sub-package of internal files
import "myproject/internal/cells/user-authenticate"  // NO
```

Rules:

1. Cell A may import Cell B's `api` package.
2. Cell A must not import Cell B's implementation package or subpackages.
3. The dependency must be listed in Cell A's `cell.yaml` `dependencies` field.
4. Structural test (`TestCellImportsUseAPIPackages`) enforces this via AST analysis.

## 8.2 Wiring File (`internal/app/wiring.go`)

The single place where all cells are constructed and connected:

```go
// internal/app/wiring.go
package app

import (
    "myproject/internal/cells/user-authenticate"
    "myproject/internal/cells/users"
    "myproject/internal/cells/users/user-invite"
    "myproject/internal/platform/config"
    "myproject/internal/platform/logging"
    "myproject/internal/platform/server"
)

type AppDeps struct {
    Config   config.Config
    Logger   *slog.Logger
    DB       *sql.DB
}

func wireCapabilities(deps AppDeps) (map[string]http.Handler, error) {
    authCap := userauthenticate.New(userauthenticate.Deps{
        Store:    userauthenticate.NewSQLStore(deps.DB),
        Logger:   deps.Logger.With("cell", "user-authenticate"),
        TokenTTL: deps.Config.AuthTokenTTL,
    })

    inviteCap := userinvite.New(userinvite.Deps{
        Store:  userinvite.NewSQLStore(deps.DB),
        Sender: emailSender,
        Logger: deps.Logger.With("cell", "user-invite"),
    })

    return map[string]http.Handler{
        "/auth/":   userauthenticate.NewHandler(authCap, deps.Logger),
        "/invite/": userinvite.NewHandler(inviteCap, deps.Logger),
    }, nil
}
```

Rules:

- `wiring.go` must stay under 300 LOC.
- If it grows, split by concern: `wiring_platform.go`, `wiring_cells.go`, `wiring_routes.go`.
- No `init()` functions. All wiring is explicit.
- Cell-to-cell dependencies are wired here (e.g., passing `authCap` to a cell that needs it).

## 8.3 Contracts

`internal/contracts/` holds cross-cutting interfaces with no natural owner. **One file per concern:**

```
internal/contracts/
├── clock.go       // Clock interface
├── eventbus.go    // EventBus interface
├── logger.go      // Logger interface (if abstracting slog)
└── doc.go
```

Cell-to-cell interfaces stay in the cell (not in `contracts/`). Only true cross-cutting concerns go here.
