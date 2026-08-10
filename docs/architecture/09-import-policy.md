# 9. Import Policy as Code

## 9.1 `policy/imports.yaml` (Standard+ tier)

```yaml
rules:
  - from: internal/cells/**
    allow:
      - pkg/**
      - internal/contracts/**
      - internal/platform/**
      - internal/cells/**/api # explicit cell contract packages only
    deny:
      - internal/app/** # cells must not import wiring or app internals

  - from: pkg/**
    deny:
      - internal/** # pkg must never import internal

  - from: cmd/**
    allow:
      - internal/app/**
    deny:
      - internal/cells/**
```

## 9.2 Validator

`scripts/imports/` parses `policy/imports.yaml`, normalizes module-local imports,
and walks the import graph using `go/ast`. For internal imports, matching rules
enforce both deny and allow lists. Run via `task policy`.
