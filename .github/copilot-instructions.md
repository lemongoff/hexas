# FF-Hexas Copilot Instructions

This repository is an independently maintained, game-oriented hard fork of go-zero. It is not an upstream compatibility mirror.

Before proposing or editing code, read these files in order:

1. `AGENTS.md`
2. `ai/README.md`
3. `ai/project-overview.md`
4. `ai/framework-lineage.md`
5. `docs/audits/2026-08-31-initial-framework-audit.md`

Repository-specific rules take precedence over generic go-zero guidance.

Key boundaries:

- The code baseline is `zeromicro/go-zero v1.10.3@925f8a2bcc159eaf3b1da0f5fc695beac26e15ff`.
- The history contains upstream `zeromicro/go-zero v1.10.3@925f8a2bcc159eaf3b1da0f5fc695beac26e15ff` plus internal customizations.
- The module path remains `github.com/zeromicro/go-zero` only to preserve existing imports.
- Do not add forward-compatibility code, automatic fallbacks, or upstream synchronization unless a task explicitly requires it.
- Audit findings are not authorization to fix runtime behavior.
- Plan first and wait for approval before editing.
- Read implementations and tests before changing them; keep scope minimal.
- Treat persistence, Redis prefix/Lua, service discovery encoding, permissions, public interfaces, config defaults, and code generation as high-risk.
- Run root-module tests from the repository root. Test `tools/goctl` separately because it is an independent Go module.
- Report changes, validation, unverified items, risks, and rollback steps. Do not commit or push unless explicitly requested.
