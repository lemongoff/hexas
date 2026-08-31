# zero-skills 来源与选择

本 Skill 基于以下固定快照适配：

- 仓库：`https://github.com/zeromicro/zero-skills.git`
- 提交：`943a13c5d82cbd3d8f896134ea6b34182bb4c32f`
- 引入日期：2026-08-31
- 上游声明许可证：MIT

## 引入内容

保留以下通用知识快照，内容位于 `upstream/`，用于按需参考：

- REST、RPC、数据库、弹性治理和 goctl 命令参考。
- 生产最佳实践。
- 常见业务服务问题排查。

未引入安装教程、工具对比、演示项目和通用 Agent 模板，因为它们不改变 FF-Hexas 框架开发决策，且部分安装流程会要求 `goctl@latest` 或自动跟随上游。

## 适配方式

- `SKILL.md` 已重写为 FF-Hexas 项目入口。
- `upstream/` 文件保持所选提交的原始内容，不能单独视为本项目规范。
- `ai/context/` 提供项目工作流层，`AGENTS.md` 始终拥有最高优先级。
- 本仓库不使用 submodule，也不自动更新该快照。

更新时应重新固定提交、审计通用规则与 FF-Hexas 的冲突，并运行 Skill、链接和文档验证。
