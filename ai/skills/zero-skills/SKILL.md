---
name: zero-skills
description: Work on Hexas go-zero framework internals, REST/RPC/storage/resilience behavior, local goctl generation, or go-zero consumer examples. Use project rules for framework changes and load pinned upstream patterns only when the task needs them.
license: MIT
---

# Hexas zero-skills

This project adaptation combines Hexas framework rules with a pinned selection from `zeromicro/zero-skills`.

## Authority

Before acting, read the relevant higher-priority project sources:

1. [`../../../AGENTS.md`](../../../AGENTS.md)
2. [`../../project-overview.md`](../../project-overview.md) and [`../../framework-lineage.md`](../../framework-lineage.md)
3. [`../../context/00-instructions.md`](../../context/00-instructions.md)
4. The matching entries in [`../../../docs/audits/2026-08-31-initial-framework-audit.md`](../../../docs/audits/2026-08-31-initial-framework-audit.md)

The vendored upstream references are general consumer-service guidance. They never override Hexas code, tests, audits, defaults, or approved design decisions.

## Workflow

1. Decide whether the task changes framework source, the independent goctl module, or a consumer example.
2. Read only the project context and upstream reference needed for that task.
3. Inspect current implementation, call sites, configuration, tests, and generators before planning changes.
4. Follow the repository approval boundary. An audit finding does not authorize runtime edits.
5. For authorized legacy work, define the target contract and choose repair, redesign, replacement, or removal.
6. Use the smallest meaningful validation surface, adding race, integration, generation, or migration tests when risk requires it.

## Project routing

- Framework workflow and approvals: [`../../context/workflows.md`](../../context/workflows.md)
- Local commands and goctl boundary: [`../../context/tools.md`](../../context/tools.md)
- Hexas framework invariants: [`../../context/patterns.md`](../../context/patterns.md)
- Source snapshot and adaptation record: [`UPSTREAM.md`](UPSTREAM.md)

## Pinned upstream references

Load these only when their subject is relevant:

- REST consumer services: [`upstream/references/rest-api-patterns.md`](upstream/references/rest-api-patterns.md)
- RPC consumer services: [`upstream/references/rpc-patterns.md`](upstream/references/rpc-patterns.md)
- SQL, MongoDB, Redis, and cache examples: [`upstream/references/database-patterns.md`](upstream/references/database-patterns.md)
- Breakers, load shedding, rate limiting, timeouts, and retries: [`upstream/references/resilience-patterns.md`](upstream/references/resilience-patterns.md)
- goctl command syntax and templates: [`upstream/references/goctl-commands.md`](upstream/references/goctl-commands.md)
- General production review: [`upstream/best-practices/overview.md`](upstream/best-practices/overview.md)
- Consumer-service troubleshooting: [`upstream/troubleshooting/common-issues.md`](upstream/troubleshooting/common-issues.md)

Do not load all references by default.

## Hexas overrides

- This repository develops the framework itself. Handler/Logic/Model layering applies to generated consumer services, not every framework package.
- Never install or select `goctl@latest` automatically. Build and test `tools/goctl` from this repository when generator behavior matters.
- Do not assume upstream version ranges, defaults, generated output, configuration, or troubleshooting fixes apply to Hexas.
- Do not update dependencies, sync upstream, generate into the working tree, or modify runtime code without task authorization.
- Existing customizations are redesign inputs, not compatibility requirements. Preserve old behavior only when the user explicitly requires it.
- Generated code boundaries must be established from the current generator and templates; upstream claims that regeneration is always safe are advisory only.
- When a general reference conflicts with local behavior, follow the implementation and record or correct the documentation drift.
