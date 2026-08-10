# 4. Cell File Schema

## 4.1 Flat Cell (Simple Capability)

```
internal/cells/{name}/
├── api/api.go          # Public interfaces and shared contract types
├── cell.yaml              # Machine-readable manifest
├── AGENTS.md              # Agent guide (purpose, start-here, tasks, validation, reliability)
├── doc.go                 # Package doc with cell ID anchor
├── {name}.go              # Implementation + constructor
├── types.go               # Domain types (tiny structs, no logic)
├── errors.go              # Sentinel errors + typed error struct
├── service.go             # Business logic
├── store.go               # Data access interface + implementation
├── handler.go             # HTTP/gRPC handler
└── *_test.go              # Collocated table-driven tests
```

## 4.2 Domain Cell (Complex, with Sub-Actions)

```
internal/cells/{domain}/
├── api/api.go          # Domain public contract
├── cell.yaml              # Domain-level manifest
├── AGENTS.md              # Domain-level agent guide
├── doc.go                 # Package doc with domain anchor
├── cell.go                # Domain implementation and handler aggregation
├── errors.go              # Shared error taxonomy
├── model.go               # Shared domain types
│
├── {action}/              # Sub-action (behavior-first leaf)
│   ├── cell.yaml          # Sub-action manifest
│   ├── AGENTS.md          # Sub-action guide
│   ├── doc.go
│   ├── api/api.go         # Sub-action public contract
│   ├── {action}.go        # Sub-action implementation + constructor
│   ├── types.go, errors.go
│   ├── service.go
│   ├── store.go
│   ├── handler.go
│   └── *_test.go
│
└── {action}/
    └── ...
```

## 4.3 When to Use Flat vs Domain Cell

- **Flat cell:** The capability has one primary behavior (e.g., `user-authenticate` — login, logout, validate are all part of the same session lifecycle).
- **Domain cell:** The capability has multiple distinct behaviors with shared types/errors (e.g., `users` — invite, register, profile are distinct actions sharing `User` type and error taxonomy).
- **Guidance:** Only create a new cell when the capability has its own domain types, error taxonomy, and lifecycle. Use sub-actions within a domain cell for closely related operations.
