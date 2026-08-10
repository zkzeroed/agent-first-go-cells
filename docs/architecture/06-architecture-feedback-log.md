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
