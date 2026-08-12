# Architecture Feedback Log

> Durable working log for improving the Agent-First Go architecture and
> project structure.

Use this file to capture learnings that directly affect the cell architecture:
cell schema, inside-out build flow, manifests, wiring/dependencies, context
packs, impact analysis, and structural guardrails.

## When to Add an Entry

Add an entry when any of the following happens:

- A cell schema, scaffold, manifest, dependency, wiring, context-pack, or
  impact-analysis convention makes cell work noticeably easier.
- A cell architecture rule or structural guardrail causes confusion or slows
  cell work down.
- A cell validation command fails because of an architecture assumption.
- You discover a repeatable fix or workaround for cell creation, modification,
  wiring, indexing, or validation.
- You identify a cell guardrail that should be added, removed, relaxed, or made
  stricter.
- You see cell naming, dependency, wiring, manifest, or context-pack drift.
- You make a cell architecture decision that future agents should understand.

Do not add entries for generic repository maintenance such as `.editorconfig`,
`.github`, CI metadata, issue templates, dependency update configuration,
licensing, README polish, or editor settings unless the issue exposes a direct
problem with the cell architecture itself.

## Entry Format

Use this format for each entry:

```md
### YYYY-MM-DD — Short title

- **Context:** What cell architecture task or area was being worked on?
- **What worked:** What helped agents create, modify, wire, validate, or
  navigate cells more safely?
- **Difficulty:** What cell-architecture rule, file schema, dependency pattern,
  or validation step caused friction?
- **Cause:** What appears to be the root cause?
- **Fix applied:** What was done now, if anything?
- **Suggested improvement:** What should be amended, changed, removed,
  automated, or documented next for the cell architecture?
- **Files/commands:** Relevant cell paths, architecture docs, or commands.
```

## Log Entries

<!-- Add new entries below this line, newest first. -->

### 2026-08-12 — Index freshness must include its schema contract

- **Context:** Orienting a new repository regenerated an otherwise
  hash-current empty `gen/cells.json` solely because the checked-in schema was
  `agent-first/v2` while the generator emitted `agent-first/v5`.
- **What worked:** The manifest hash continued to detect source and guide
  changes, and regeneration made the schema mismatch immediately visible in
  the working tree.
- **Difficulty:** `task check-index` reported the v2 index fresh, so `task
  orient` was unexpectedly mutating in a clean checkout.
- **Cause:** Freshness compared only the manifest hash; `schemaVersion` was
  emitted but not part of the freshness contract.
- **Fix applied:** Established `agent-first/v1` as the completed baseline and
  made `check-index` reject a mismatched schema version before checking hash
  and context packs.
- **Suggested improvement:** Treat every future incompatible index shape as a
  deliberate schema-version migration and regenerate all tracked indexes in
  the same change.
- **Files/commands:** `tools/agent/index`; `task index`; `task check-index`;
  `task orient`.

### 2026-08-12 — Registry-driven public-package boundaries

- **Context:** Audited whether public packages at the module root, `pkg/`, and
  named top-level directories received equivalent import protection.
- **What worked:** The library-package registry already identifies the exact
  directories entitled to compose declared private cells, and manifest/source
  validation already matches those direct imports to exact dependencies.
- **Difficulty:** A `pkg/**` policy rule protected only that naming convention;
  registered public packages at `.` or another path could otherwise import
  module-level internal surfaces.
- **Cause:** The import policy mixed static source-path globs with dynamic
  library registration, while `cellsRoot` was configurable in discovery but
  not consistently in scaffolding and structural tooling.
- **Fix applied:** Made private cells a fixed `internal/cells` convention,
  removed the partial `cellsRoot` configuration, and made import validation
  registry-driven. Every registered library package may directly import only
  private-cell implementations; `internal/app`, `internal/platform`, and
  `internal/contracts` are rejected regardless of its public path.
- **Suggested improvement:** Keep any future private-root customization as an
  all-tooling migration, rather than reintroducing a partial configuration
  field.
- **Files/commands:** `tools/agent/imports`, `tools/agent/projectconfig`,
  `policy/architecture.yaml`, `task policy`, `task ready` across all examples.

### 2026-08-12 — First research-library package walkthrough

- **Context:** Created `examples/research-library-project` solely through the
  public-package and private-cell scaffolds, then implemented a cited modular
  reduction package and navigated it through index, context, scope, dependency,
  targeted validation, and readiness checks.
- **What worked:** The library scaffold registered the importable package; its
  conformance placeholder was immediately visible in `task scope`. Replacing
  it with a paper-defined record appeared in `task context` and `task
  cells-json`. The manifest/source entrypoint check stopped an incorrect
  private-cell symbol before generated metadata could be trusted. `task ready`
  validated both the public package and its private helper.
- **Difficulty:** A new project cannot run `task orient` until `task index` has
  been run once. The private-cell scaffold is broad for a narrow helper and
  creates several unrelated TODO files. The conformance validator accepts a
  nonexistent citation path and accepts `pdfPages` for a non-PDF source, so a
  record can look more authoritative than its local evidence warrants.
- **Cause:** `orient` delegates to the index-reading `cells` command; the
  generic application-cell template is reused for library helpers; citation
  syntax is validated but citations are not resolved or type-checked.
- **Fix applied:** Added the runnable research-library example, corrected a
  manifest entrypoint mismatch found by `task index`, made `orient` refresh the
  index, changed `new-cell` to a lean default with `new-cell-ext` retaining the
  extended application template, and added typed, locally resolved citation
  locators plus verified-evidence enforcement.
- **Suggested improvement:** Consider validating PDF page bounds in a future
  PDF-aware evidence checker; current validation confirms local file type and
  declared positive page locators, but does not parse PDF pagination.
- **Files/commands:** `examples/research-library-project`; `task index`,
  `task orient`, `task context`, `task scope`, `task validate-cell`, `task
  ready`.
