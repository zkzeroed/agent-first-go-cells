# 11. Multi-Agent Collaboration

## 11.1 Ownership Model

- **Strong ownership:** Each cell is owned by one agent at a time. Two agents do not edit the same cell concurrently.
- **Weak/shared ownership:** `internal/platform/`, `internal/contracts/`, `pkg/`, `AGENTS.md` are shared. Changes require a coordination message.
- **Cross-cell dependencies:** A cell may depend on another cell's public interface only. If an agent needs to change a sibling cell's interface, the sibling cell's owner is notified.

## 11.2 Protocol: Claim → Work → Integrate → Release

```
┌─────────────────────────────────────────────────────────────────┐
│  1. CLAIM: Agent reads gen/cells.json and claims a cell through │
│     the team's coordinator or work tracker                      │
│                                                                 │
│  2. WORK: Agent works entirely within its cell directory        │
│     → Creates/modifies files following the fixed schema         │
│     → Updates cell.yaml if dependencies/invariants change       │
│     → Runs tests: go test ./internal/cells/{name}/...           │
│     → Does NOT touch wiring.go or other cells                   │
│                                                                 │
│  3. INTEGRATE: Agent updates wiring.go atomically               │
│     → Adds New() call to wiring function                        │
│     → Adds route mapping to handler map                         │
│     → Runs task index to regenerate gen/cells.json              │
│     → Runs task impact to check affected cells                  │
│     → Runs full test suite: task test                           │
│                                                                 │
│  4. RELEASE: Agent records release with the coordinator or      │
│     work tracker, then produces a handoff summary (see below)   │
│     → Notifies coordinator that cell is complete                │
└─────────────────────────────────────────────────────────────────┘
```

## 11.3 Handoff Format

```markdown
## Handoff: user-authenticate

**Status:** Complete / In-Progress / Blocked
**Files Modified:**

- internal/cells/user-authenticate/user_authenticate.go (new interface)
- internal/cells/user-authenticate/service.go (implementation)
- internal/cells/user-authenticate/cell.yaml (new manifest)
- internal/app/wiring.go (added wiring)

**Tests:** All passing (`go test ./internal/cells/user-authenticate/...`)
**Linting:** Clean (`golangci-lint run ./internal/cells/user-authenticate/...`)
**Impact:** No other cells affected (`task impact` confirmed)
**Next Steps for Reviewer:**

1. Verify interface matches requirements in docs/specs/user-authenticate-spec.md
2. Run integration tests: `task test-integration`
3. Check wiring.go for correct dependency injection
```
