# 1. First-Principles Analysis

> **Target Language:** Go 1.26+
> **Primary Audience:** AI coding agents (single-agent and multi-agent)
> **Secondary Audience:** Human developers

---

## 1.1 How Coding Agents Actually Work

Coding agents interact with codebases through fundamentally different mechanisms than humans:

- **Context locality wins.** An agent's working memory is its context window — a finite token budget. If a feature's code is spread across 5 directories, the agent burns tokens on navigation. Optimal ergonomics require everything for a feature in one directory, readable in 2–3 file reads totaling < 500 lines.
- **Search is the primary navigation mode.** Agents locate code by grep, symbol lookup, and filename. Predictable, 1:1 mapping between concept and location eliminates search ambiguity. Behavior-first naming (`user-authenticate/`) is more grep-deterministic than generic terms (`auth/`, `users/`).
- **Modification is local and repetitive.** Agents generate code in bounded chunks. They excel at "add a new file matching the existing pattern" and struggle with "refactor six packages simultaneously."
- **Hallucinations multiply with implicit dependencies.** Every global helper, every cross-package mutation, every unwritten convention is a place where an agent can guess wrong. The antidote is explicit contracts enforced by the type system and directory structure.
- **Multi-agent work requires ownership boundaries.** Agents can work in parallel only when ownership is clear. The natural unit of ownership is a vertical slice, not a horizontal layer.

## 1.2 Why Traditional architectures Are Suboptimal for Agents

Layered, hexagonal, onion, and DDD architectures are designed for human mental models. For agents, they create:

- **Horizontal scattering:** A single feature touches `handler/`, `service/`, `repository/`, `domain/`, `dto/` — a navigation nightmare.
- **Search ambiguity:** A repository called `UserRepository` is less searchable than `user_authenticate.go` because the latter embeds the behavior.
- **Dense import graphs:** An agent cannot edit one layer without understanding imports into every other layer.
- **Naming drift toward generic terms:** `models.User` could be anything; `userauthenticate.Session` is unambiguous.
