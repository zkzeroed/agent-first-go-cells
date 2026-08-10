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
- **Files/commands:** `task ready`, `scripts/impact/`, `examples/reference-project/`.

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
- **Files/commands:** `scripts/new-cell.sh`, `scripts/new-domain.sh`, `scripts/structure/structure_test.go`, `scripts/imports/`, `task doctor`.

### 2026-08-09 — Make agent metadata fail closed

- **Context:** Audited the bootstrap's navigation, manifest, context-pack, impact, and validation workflow.
- **What worked:** Fixed cell locality and a small task surface make the intended agent loop easy to discover.
- **Difficulty:** Several commands reported success while accepting malformed manifests, ignored test failures, stale agent guidance, or ambiguous dependency edges.
- **Cause:** Documentation expressed stronger contracts than the custom parsers and shell wrappers enforced.
- **Fix applied:** Hardened manifest parsing/validation, context freshness, dependency matching, test failure propagation, and scaffold validation; added focused regression coverage.
- **Suggested improvement:** Make cell membranes explicit `api` packages before introducing cross-cell dependencies, because Go imports packages rather than individual interface files.
- **Files/commands:** `scripts/manifest/`, `scripts/index/`, `scripts/impact/`, `Taskfile.yml`, `task doctor`, `task ready`.

### 2026-07-04 — Created architecture feedback workflow

- **Context:** Added a durable feedback loop for evolving the Agent-First Go architecture while real use cases are implemented.
- **What worked:** Existing project already has clear architecture docs, root `AGENTS.md`, project rules, and project-local skills, making `.agents/skills/architecture-feedback/` and `docs/architecture/18-architecture-feedback-log.md` natural homes for this workflow.
- **Difficulty:** The architecture previously described how to build and validate cells, but did not explicitly instruct agents to record architecture/tooling friction as they encounter it.
- **Cause:** Feedback capture was implicit rather than part of the agent operating protocol.
- **Fix applied:** Created this log and added an `architecture-feedback` skill plus reminders in `AGENTS.md`, `.rules`, and `.windsurfrules`.
- **Suggested improvement:** Consider adding a `task architecture-log` or `task feedback-log` target later if entries become frequent, and consider validating that feedback entries follow the template.
- **Files/commands:** `docs/architecture/18-architecture-feedback-log.md`, `.agents/skills/architecture-feedback/SKILL.md`, `AGENTS.md`, `.rules`, `.windsurfrules`.
