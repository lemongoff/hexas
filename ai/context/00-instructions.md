# FF-Hexas AI 工作流上下文

本目录基于 `zeromicro/ai-context` 的两层上下文思路整理，但服务对象是 FF-Hexas 框架源码，不是使用 go-zero 搭建的业务服务。

## 指令优先级

发生冲突时按以下顺序执行：

1. 仓库根目录 `AGENTS.md`。
2. `ai/project-overview.md`、`ai/framework-lineage.md` 和对应审计条目。
3. 本目录的工作流、工具和模式说明。
4. `ai/skills/zero-skills/SKILL.md` 及其按需引用的通用资料。
5. 官方在线文档。

代码、测试、配置和脚本的实际行为高于文档描述；发现漂移时在当前改动范围内同步修正文档。

## 先判断任务对象

- 修改 `core/`、`rest/`、`zrpc/`、`gateway/`、`mcp/` 或 `internal/`：按框架源码任务处理。
- 修改 `tools/goctl/`：按独立 Go module 和生成器任务处理。
- 创建示例或验证 goctl 产物：只有用户明确要求时，才按业务服务模式使用 `.api`、`.proto` 和 Handler/Logic/Model 参考。
- 审计或设计遗留能力：旧实现只作为证据，可以修复、重构、替换或删除；未获得实现授权时只记录结论。

## 核心规则

- 修改前读取实现、调用点、配置、测试和相关审计，先列计划并等待确认。
- 不把通用业务服务的三层目录规则强加给框架内部包。
- 不自动安装或调用 `goctl@latest`；检查生成器时使用 `tools/goctl` 当前源码并单独验证。
- 不按官方最新行为推断 FF-Hexas；本仓库只固定 `v1.10.3` 基线并独立演进。
- 高风险改动先定义目标契约、失败模型、迁移和回滚，再决定实现方式。
- 生成文件优先从源定义和本地生成器重建；运行生成命令前检查覆盖范围和工作树。
- 只执行与用户范围一致的写操作，不自动同步上游、更新依赖、发布、提交或推送。

## 决策入口

```text
任务
├─ 框架缺陷或功能 → workflows.md 的框架变更流程
├─ Redis/Mongo/缓存 → patterns.md + 审计 A-01 至 A-04
├─ 服务发现/RPC → patterns.md + 审计 A-05、A-10
├─ REST/权限/protobuf → patterns.md + 审计 A-06、A-07
├─ 日志 → patterns.md + 审计 A-08
├─ goctl/模板 → tools.md + 审计 A-09
├─ 业务服务示例 → zero-skills 中对应上游参考
└─ 文档/审计 → AGENTS.md 的文档验证矩阵
```

来源和适配记录见 [UPSTREAM.md](UPSTREAM.md)。
