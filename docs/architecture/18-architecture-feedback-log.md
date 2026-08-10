# 18. architecture Feedback Log

> Durable working log for improving the Agent-First Go architecture and project structure.
>
> Use this file to capture learnings that directly affect the cell architecture: cell schema, inside-out build flow, manifests, wiring/dependencies, context packs, impact analysis, and structural guardrails.

## Purpose

This project is trying to converge on the optimal agent-first Go cell architecture. Agents should record architectural learning only when it directly affects how cells are named, scaffolded, structured, wired, validated, indexed, navigated, or safely modified.

## When to Add an Entry

Add an entry when any of the following happens:

- A cell schema, scaffold, manifest, dependency, wiring, context-pack, or impact-analysis convention makes cell work noticeably easier.
- A cell architecture rule or structural guardrail causes confusion or slows cell work down.
- A cell validation command fails because of an architecture assumption.
- You discover a repeatable fix or workaround for cell creation, modification, wiring, indexing, or validation.
- You identify a cell guardrail that should be added, removed, relaxed, or made stricter.
- You see cell naming, dependency, wiring, manifest, or context-pack drift.
- You make a cell architecture decision that future agents should understand.

Do not add entries for generic repository maintenance such as `.editorconfig`, `.github`, CI metadata, issue templates, dependency update configuration, licensing, README polish, or editor settings unless the issue exposes a direct problem with the cell architecture itself.

## Entry Format

Use this format for each entry:

```md
### YYYY-MM-DD — Short title

- **Context:** What cell architecture task or area was being worked on?
- **What worked:** What helped agents create, modify, wire, validate, or navigate cells more safely?
- **Difficulty:** What cell-architecture rule, file schema, dependency pattern, or validation step caused friction?
- **Cause:** What appears to be the root cause?
- **Fix applied:** What was done now, if anything?
- **Suggested improvement:** What should be amended, changed, removed, automated, or documented next for the cell architecture?
- **Files/commands:** Relevant cell paths, architecture docs, or commands.
```

## Log Entries

### 2026-08-10 — Keep all validation caches writable

- **Context:** Ran the complete architecture handoff in a sandbox with a read-only home directory.
- **What worked:** Taskfile defaults already redirected Go build and ccache state to a writable temporary location.
- **Difficulty:** `task lint` still emitted extensive cache-write warnings from golangci-lint despite passing.
- **Cause:** golangci-lint uses its own cache location, independent of `GOCACHE` and `CCACHE_DIR`.
- **Fix applied:** Added a temporary `GOLANGCI_LINT_CACHE` default and documented all three cache overrides in the agent workflow.
- **Suggested improvement:** When adding a validation tool, identify and redirect its persistent state before making it part of the agent handoff loop.
- **Files/commands:** `Taskfile.yml`, `AGENTS.md`, `docs/architecture/14-taskfile-targets.md`, `task lint`.

### 2026-08-10 — Keep generated context visible in Zed

- **Context:** Audited the Zed project configuration against the cell-navigation workflow.
- **What worked:** Root `AGENTS.md`, project-local skills, and Taskfile tasks are directly consumable in a trusted Zed worktree.
- **Difficulty:** Zed excluded `gen/` from scans, searches, and the project tree even though it contains the generated cell index and bounded context packs.
- **Cause:** Generated files were treated as generic build output rather than as agent-navigation inputs.
- **Fix applied:** Restored Zed's standard exclusions, removed `gen/` from exclusions, and made validation tasks save open buffers before executing.
- **Suggested improvement:** Keep editor configuration small; expose new Taskfile commands in Zed only when their required inputs can be supplied safely without duplicating workflow logic.
- **Files/commands:** `.zed/settings.json`, `.zed/tasks.json`, `AGENTS.md`, `task context`, `task scope`.

### 2026-08-10 — Make cell scope explicit and fail closed

- **Context:** Reviewed how an agent can refactor a cell without silently expanding into unrelated code.
- **What worked:** API membranes, source-backed manifests, and transitive impact already expose the dependency graph.
- **Difficulty:** They did not state which files a specific task was authorized to change, and shared integration surfaces had incomplete impact coverage.
- **Cause:** Navigation and validation were separate commands without an explicit task-local boundary.
- **Fix applied:** Added read-only scope and fail-closed scope verification commands, including target removal from `HEAD`; impact now treats contracts, platform, and wiring as full-project shared surfaces.
- **Suggested improvement:** Keep verification opt-in until real feature work establishes a safe task-identity mechanism for automatic handoff enforcement.
- **Files/commands:** `tools/agent/scope/`, `tools/agent/impact/`, `task scope`, `task verify-scope`.

### 2026-08-10 — Make impact follow generated guide context

- **Context:** Revalidated impact output against the generated context-pack inputs before extending agent feedback.
- **What worked:** Manifests expose a deterministic dependency graph suitable for breadth-first downstream analysis.
- **Difficulty:** Cell-local `AGENTS.md` changes changed generated context but were absent from impact, and only direct dependents were reported.
- **Cause:** File filtering covered only Go/YAML and dependency traversal stopped after one hop.
- **Fix applied:** Impact now maps cell-local guides, uses directory-boundary ownership, and reports sorted transitive downstream cells.
- **Suggested improvement:** Keep the analyzer limited to declared cell relationships; dynamic runtime coupling belongs in broader integration validation.
- **Files/commands:** `tools/agent/impact/`, `docs/architecture/07-impact-analysis.md`, `task impact`.

### 2026-08-10 — Verify metadata against the source graph

- **Context:** Revalidated the agent navigation and validation loop before extending its guardrails.
- **What worked:** Manifests, API packages, and generated context packs already provide a compact model agents can navigate.
- **Difficulty:** Entrypoint symbols and dependency declarations could drift from source while metadata checks still passed.
- **Cause:** Parsing retained only entrypoint paths, and import policy validated allowed paths without comparing them to the manifest graph.
- **Fix applied:** Scaffolders now fail safely on existing paths and emit valid API entrypoints; metadata validation resolves entrypoint files/symbols and exact direct API dependencies.
- **Suggested improvement:** Keep this AST-level check focused on declared cell boundaries; add type-checking only if evidence shows syntax-level validation is insufficient.
- **Files/commands:** `tools/agent/new-cell.sh`, `tools/agent/new-domain.sh`, `tools/agent/manifest/`, `task check-manifests`.

### 2026-08-10 — Domain example completes the schema exercise

- **Context:** Built a nested `profiles/profile-create` domain action after a cold-start documentation and skill walkthrough.
- **What worked:** Exact parent dependencies, `api` contracts, root-aware indexes, and reverse dependency output made the domain relationship explicit.
- **Difficulty:** The walkthrough exposed a few remaining pre-migration code fragments in the build skill and wiring chapter.
- **Cause:** Earlier conformance checks covered terminology but not code identifiers and list numbering.
- **Fix applied:** Corrected the snippets and added the domain example with generated context packs.
- **Suggested improvement:** Extend doc-contract tests to parse or compile code snippets only if documentation volume grows again.
- **Files/commands:** `examples/reference-project/internal/cells/profiles/`, `.agents/skills/build-cell/SKILL.md`, `task deps ROOT=examples/reference-project`.

### 2026-08-10 — Selected-project impact ownership

- **Context:** Ran the normal root `task ready` after adding token cells to the nested reference project.
- **What worked:** Untracked example source files appeared in impact output.
- **Difficulty:** They had no owning or affected cells because impact mapped against only the root project's manifests.
- **Cause:** Navigation accepted `ROOT`, but impact did not share that project-selection input.
- **Fix applied:** Added `ROOT`/`--root` to impact, changed, and ready; impact now filters repository changes into selected-project-relative paths before mapping ownership.
- **Suggested improvement:** Keep Git discovery repository-scoped; selected-project support should remain a path filter rather than introduce nested Git semantics.
- **Files/commands:** `task ready`, `tools/agent/impact/`, `examples/reference-project/`.

### 2026-08-10 — Cold-start skill drift surfaced by a real cell build

- **Context:** Simulated a new agent reading the architecture guide and adding a cryptographic token flow to the reference project.
- **What worked:** The root guide, cell schema, root-aware index, dependency query, and generated context packs led directly to a valid four-cell graph.
- **Difficulty:** The project-local navigation/build skills and manifest example still described the retired root-package interface membrane and non-cell dependency syntax.
- **Cause:** The API-boundary migration updated primary docs but not all task-specific guidance.
- **Fix applied:** Updated those skills and the manifest example to use `api/api.go` and exact cell dependencies.
- **Suggested improvement:** Keep architecture migrations paired with a targeted scan of project-local skills and examples.
- **Files/commands:** `.agents/skills/`, `docs/architecture/05-manifest-schema.md`, `task index ROOT=examples/reference-project`.

### 2026-08-10 — Reference projects need portable tooling support

- **Context:** Built `examples/reference-project` with two cells and an explicit API dependency.
- **What worked:** The API package, manifests, and inside-out file layout made the cross-cell boundary clear.
- **Difficulty:** Root navigation tools were fixed to `internal/cells`, and impact analysis ignored newly created files.
- **Cause:** Tool paths and change discovery were repository-root assumptions rather than configurable project roots and untracked-aware inputs.
- **Fix applied:** Added `ROOT`/`--root` support for metadata navigation and included non-ignored untracked Go/YAML files in impact analysis.
- **Suggested improvement:** Keep impact Git-repository scoped; add fixture-specific validation only if reference projects require broader checks.
- **Files/commands:** `examples/reference-project/`, `go test ./...` from the example directory.

### 2026-08-10 — Keep prose subordinate to executable sources

- **Context:** Audited every architecture chapter, the root README, and agent guide after the API-boundary migration.
- **What worked:** The Taskfile and scaffold scripts provide compact, executable sources of truth.
- **Difficulty:** Long copied Taskfile, scaffolder, IDE, and code examples drifted after the repository changed.
- **Cause:** Static duplication had no conformance check.
- **Fix applied:** Replaced duplicated material with concise operational references and added a test for retired architecture terms.
- **Suggested improvement:** Add examples only when an executable scaffold cannot communicate the behavior.
- **Files/commands:** `docs/architecture/`, `README.md`, `AGENTS.md`, `go test ./...`.

### 2026-08-10 — Enforce the Go-level cell membrane

- **Context:** Follow-up architecture pass after hardening manifests, context packs, and impact analysis.
- **What worked:** A dedicated `api` package gives agents an import path that exactly matches the intended public boundary.
- **Difficulty:** The former interface-file rule was not enforceable because Go imports packages, not individual files.
- **Cause:** The public contract and implementation shared one importable package.
- **Fix applied:** Scaffolded `api/api.go`, restricted cell imports to `api` paths, enforced internal allow lists, and corrected the function-override guardrail.
- **Suggested improvement:** Add a complete fixture project once the first production cells exist, to exercise the full multi-cell workflow under the final package layout.
- **Files/commands:** `tools/agent/new-cell.sh`, `tools/agent/new-domain.sh`, `tools/agent/structure/structure_test.go`, `tools/agent/imports/`, `task doctor`.

### 2026-08-09 — Make agent metadata fail closed

- **Context:** Audited the bootstrap's navigation, manifest, context-pack, impact, and validation workflow.
- **What worked:** Fixed cell locality and a small task surface make the intended agent loop easy to discover.
- **Difficulty:** Several commands reported success while accepting malformed manifests, ignored test failures, stale agent guidance, or ambiguous dependency edges.
- **Cause:** Documentation expressed stronger contracts than the custom parsers and shell wrappers enforced.
- **Fix applied:** Hardened manifest parsing/validation, context freshness, dependency matching, test failure propagation, and scaffold validation; added focused regression coverage.
- **Suggested improvement:** Make cell membranes explicit `api` packages before introducing cross-cell dependencies, because Go imports packages rather than individual interface files.
- **Files/commands:** `tools/agent/manifest/`, `tools/agent/index/`, `tools/agent/impact/`, `Taskfile.yml`, `task doctor`, `task ready`.

### 2026-07-04 — Created architecture feedback workflow

- **Context:** Added a durable feedback loop for evolving the Agent-First Go architecture while real use cases are implemented.
- **What worked:** Existing project already has clear architecture docs, root `AGENTS.md`, project rules, and project-local skills, making `.agents/skills/architecture-feedback/` and `docs/architecture/18-architecture-feedback-log.md` natural homes for this workflow.
- **Difficulty:** The architecture previously described how to build and validate cells, but did not explicitly instruct agents to record architecture/tooling friction as they encounter it.
- **Cause:** Feedback capture was implicit rather than part of the agent operating protocol.
- **Fix applied:** Created this log and added an `architecture-feedback` skill plus reminders in `AGENTS.md`, `.rules`, and `.windsurfrules`.
- **Suggested improvement:** Consider adding a `task architecture-log` or `task feedback-log` target later if entries become frequent, and consider validating that feedback entries follow the template.
- **Files/commands:** `docs/architecture/18-architecture-feedback-log.md`, `.agents/skills/architecture-feedback/SKILL.md`, `AGENTS.md`, `.rules`, `.windsurfrules`.
