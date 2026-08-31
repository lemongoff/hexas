# ai-context 来源与适配

本目录基于以下固定快照的结构和主题重新整理：

- 仓库：`https://github.com/zeromicro/ai-context.git`
- 提交：`bc525eedc924fe53b5d26f27e41595cfdb347477`
- 引入日期：2026-08-31
- 上游声明许可证：MIT

统一的来源与许可说明见 [`../../THIRD_PARTY_NOTICES.md`](../../THIRD_PARTY_NOTICES.md) 和 [`../../LICENSES/MIT-zeromicro.txt`](../../LICENSES/MIT-zeromicro.txt)。固定上游提交没有独立 `LICENSE` 文件，本仓库不据此虚构额外版权声明。

采用了上游 `00-instructions.md`、`workflows.md`、`tools.md` 和 `patterns.md` 的“工作流层”划分，但内容已针对 FF-Hexas 重写。没有建立 submodule，也不自动跟随上游。

主要适配：

- 从业务服务生成流程改为框架源码维护流程。
- 用本地 `tools/goctl` 替代 `goctl@latest`。
- 加入审计、重新设计、独立版本和高风险验证边界。
- 删除“所有任务都必须生成服务 README/API.md”等不适用于框架仓库的通用要求。

重新引入上游变化时必须固定新提交、审计差异，并保持 `AGENTS.md` 的最高优先级。
