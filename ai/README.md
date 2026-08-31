# FF-Hexas AI 文档入口

本目录只维护 FF-Hexas 自身的项目事实和协作约束。官方 go-zero 的通用资料可作为参考，但不能覆盖本仓库的硬分支边界、已有定制和审计结论。

## 阅读顺序

1. `../AGENTS.md`：强制工作流、编辑边界和验证规则。
2. `project-overview.md`：代码结构、定制区域和任务路由。
3. `framework-lineage.md`：来源分支、精确提交和版本策略。
4. `../docs/audits/2026-08-31-initial-framework-audit.md`：初始基线的遗留问题和重设计输入。

## 按任务选择入口

| 任务 | 必读内容 | 最小验证 |
| --- | --- | --- |
| Redis、缓存、异步落盘 | 初始审计 A-01 至 A-04 | `go test ./core/stores/cache/... ./core/stores/redis/... ./core/stores/mon/... ./core/stores/monc/...` |
| etcd、RPC、服务发现 | 初始审计 A-05、A-10 | `go test ./core/discov/... ./zrpc/...` |
| REST、权限、protobuf | 初始审计 A-06、A-07 | `go test ./rest/...` |
| 日志 | 初始审计 A-08 | `go test ./core/logx/...` |
| goctl、模板、生成器 | `project-overview.md` 的独立 module 说明、初始审计 A-09 | 在 `tools/goctl` 内运行 `go test ./...` |
| 文档或 AI 入口 | 本文件、`AGENTS.md` | `git diff --check` |
| 上游同步或版本升级 | `framework-lineage.md` 全文 | 必须单独制定差异审计和全量验证计划 |

## 维护原则

- 只记录已由代码、Git 元数据或实际命令确认的事实。
- 推测必须标为“待验证”，不能写成既定契约。
- 审计文档记录问题不等于授权修复。
- 获得专项改造授权后，审计中的旧行为只作为证据，不是必须保留的契约；可以选择修复、重构、替换或删除。
- 重新设计前先明确目标契约，并评估 FF-Hexas 调用方、存量数据、配置、部署切换和回滚；只有用户明确要求时才保留旧行为兼容。
- 运行时契约变化必须同步更新项目概览、来源说明或对应审计状态。
- 不将官方 go-zero 的“兼容性”默认带入 FF-Hexas；本仓库按独立框架演进。
