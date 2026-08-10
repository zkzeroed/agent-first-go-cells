---
name: architecture-feedback
description: Record learnings, friction, fixes, and improvement proposals that directly affect the Agent-First Go cell architecture; cell schema, manifests, wiring/dependencies, context packs, impact analysis, and structural guardrails.
---

# architecture Feedback

Use this skill only when work reveals something useful about the Agent-First Go cell architecture.

The goal is to improve the cell model itself: how cells are named, scaffolded, structured, wired, validated, indexed, navigated, and safely modified.

## Feedback Log

Write durable observations to:

```text
docs/architecture/06-architecture-feedback-log.md
```

Do not create scattered notes elsewhere unless the user asks. Keep architecture-learning entries centralized in this log.

## When to Use

Use this skill when you encounter or notice any of the following:

- A cell schema, scaffold, manifest, dependency, wiring, context-pack, or impact-analysis convention works especially well.
- A cell architecture rule or structural guardrail creates friction while creating or modifying cells.
- A cell validation command fails because of an architecture assumption.
- A fix or workaround should be reused by future agents working on cells.
- A cell guardrail should be amended, removed, strengthened, automated, or moved.
- A cell naming, dependency, manifest, context-pack, or wiring convention is ambiguous.
- The architecture docs are missing cell-specific guidance needed to complete a real use case.

## What Not to Record

Do not use this skill for generic repository maintenance such as `.editorconfig`, `.github`, CI metadata, issue templates, dependency update configuration, licensing, README polish, or editor settings unless the issue exposes a direct problem with the cell architecture itself.

## What to Record

Add a new dated entry using this format:

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

## How to Decide Whether to Edit Rules Too

- If the observation is speculative or based on one instance, only add a log entry.
- If the same issue is clearly repeatable or already has an obvious low-risk fix, update the relevant docs, skill, or rule in the same change.
- If changing a rule could alter architecture behavior or constrain future agents, log the recommendation and ask the user before changing the rule.

## Keep Entries Useful

- Be specific. Name files, commands, and failure modes.
- Prefer root causes over symptoms.
- Separate what was fixed now from what should be improved later.
- Do not use the log as a generic progress journal; record cell-architecture learning only.
- Keep entries concise enough that future agents can scan them quickly.
