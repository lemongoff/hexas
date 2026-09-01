# Hexas AI 文档入口

本目录只维护 Hexas 自身的项目事实和协作约束。官方 go-zero 的通用资料可作为参考，但不能覆盖本仓库的硬分支边界、已有定制和审计结论。

## 阅读顺序

1. `../AGENTS.md`：强制工作流、编辑边界和验证规则。
2. `project-overview.md`：代码结构、定制区域和任务路由。
3. `framework-lineage.md`：基线提交、module path 和版本策略。
4. `context/00-instructions.md`：适配本仓库的 AI 工作流和决策入口。
5. `skills/zero-skills/SKILL.md`：按任务加载的 go-zero 知识路由。
6. `../docs/audits/2026-08-31-initial-framework-audit.md`：初始基线的遗留问题和重设计输入。

## 按任务选择入口

| 任务 | 必读内容 | 最小验证 |
| --- | --- | --- |
| Redis、缓存、异步落盘 | 初始审计 A-01 至 A-04 | `go test ./core/stores/cache/... ./core/stores/redis/... ./core/stores/mon/... ./core/stores/monc/...` |
| etcd、RPC、服务发现 | 初始审计 A-05、A-10 | `go test ./core/discov/... ./zrpc/...` |
| REST、权限、protobuf | 初始审计 A-06、A-07 | `go test ./rest/...` |
| 日志 | 初始审计 A-08 | `go test ./core/logx/...` |
| goctl、模板、生成器 | `context/tools.md`、初始审计 A-09、zero-skills goctl 参考 | 在 `tools/goctl` 内运行 `go test ./...` |
| 业务服务示例、API/RPC 用法 | `skills/zero-skills/SKILL.md` 及其对应上游参考 | 使用本地 goctl，在隔离目录构建和测试 |
| 文档或 AI 入口 | 本文件、`AGENTS.md` | `git diff --check` |
| 上游同步或版本升级 | `framework-lineage.md` 全文 | 必须单独制定差异审计和全量验证计划 |

## 两层 AI 资产

- `context/` 是工作流层，回答“当前仓库应该怎样做”。其内容基于 `zeromicro/ai-context@bc525eedc924fe53b5d26f27e41595cfdb347477` 重写。
- `skills/zero-skills/` 是知识层，回答具体 REST、RPC、数据库、弹性治理和 goctl 问题。通用参考固定自 `zeromicro/zero-skills@943a13c5d82cbd3d8f896134ea6b34182bb4c32f`。
- 两者都以普通文件快照纳入仓库，不使用 submodule，不自动更新。
- 上游参考面向使用 go-zero 的业务服务，只有项目入口明确路由时才按需读取，不能覆盖 Hexas 规则。

## 维护原则

- 只记录已由代码、Git 元数据或实际命令确认的事实。
- 推测必须标为“待验证”，不能写成既定契约。
- 审计文档记录问题不等于授权修复。
- 获得专项改造授权后，审计中的旧行为只作为证据，不是必须保留的契约；可以选择修复、重构、替换或删除。
- 重新设计前先明确目标契约，并评估 Hexas 调用方、存量数据、配置、部署切换和回滚；只有用户明确要求时才保留旧行为兼容。
- 运行时契约变化必须同步更新项目概览、来源说明或对应审计状态。
- 不将官方 go-zero 的“兼容性”默认带入 Hexas；本仓库按独立框架演进。
- 更新 `context/` 或 `skills/zero-skills/` 的上游内容时固定新 commit，记录选择范围并重新审计项目覆盖规则。
