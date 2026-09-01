# Hexas

Hexas 是面向游戏服务场景维护的 go-zero 内部分支。仓库以
[`github.com/zeromicro/go-zero`](https://github.com/zeromicro/go-zero) `v1.10.3`
为代码基线，保留当前工作树中已有的日志、Redis、MongoDB、服务发现、REST 和 RPC 定制，作为后续游戏化框架改造的代码基础。

这些遗留定制用于说明现状，不构成必须兼容或原样保留的约束。后续专项可以根据 Hexas 的目标契约重新设计、替换或删除。

## 基线

| 项目 | 固定值 |
| --- | --- |
| 代码基线 | `github.com/zeromicro/go-zero v1.10.3` |
| 基线提交 | `925f8a2bcc159eaf3b1da0f5fc695beac26e15ff` |
| Go module | `github.com/lemongoff/hexas` |
| 配置 module | `github.com/lemongoff/hexas-config` |
| Go 版本 | `1.24.0` |

项目 module 和 import path 已切换为 `github.com/lemongoff/hexas`。`github.com/zeromicro/go-zero` 仅表示固定的官方代码基线；Hexas 从当前基线起按独立硬分支维护，不承诺继续兼容官方后续版本、API、配置语义或运行时行为，也不默认继续合并上游。已有调用方需要同步更新 import path 和 `go.mod` 依赖，不能依赖自动兼容或 fallback。

完整来源关系和版本策略见 [ai/framework-lineage.md](ai/framework-lineage.md)。

## 文档入口

- [AGENTS.md](AGENTS.md)：AI 与开发协作的强制规则。
- [ai/README.md](ai/README.md)：AI 文档索引和任务路由。
- [ai/context/00-instructions.md](ai/context/00-instructions.md)：适配 Hexas 的日常工作流层。
- [ai/skills/zero-skills/SKILL.md](ai/skills/zero-skills/SKILL.md)：按需加载的 go-zero 知识层。
- [ai/project-overview.md](ai/project-overview.md)：目录职责、定制区域和验证命令。
- [配置体系](docs/configuration.md)：类型化默认值、Bootstrap/Runtime 边界、goctl 目录和迁移要求。
- [框架默认值](docs/framework-defaults.md)：RPC、缓存和日志默认契约。
- [初始基线审计](docs/audits/2026-08-31-initial-framework-audit.md)：已知遗留问题、重设计输入和处理优先级；基线导入阶段仅记录，未修改运行时代码。
- [上游英文说明快照](docs/upstream-readme.md)：仅作原始框架参考。

## 主要目录

- `core/`：基础设施、弹性治理、日志、服务发现、存储和并发工具；进程配置由独立 `hexas-config` module 提供。
- `rest/`：HTTP 服务、路由、中间件、客户端和编解码。
- `zrpc/`：gRPC 服务、客户端、服务发现和负载均衡。
- `gateway/`：HTTP/gRPC 网关。
- `mcp/`：MCP 服务实现。
- `tools/goctl/`：独立 Go module 的代码生成工具。
- `ai/context/`：基于 `zeromicro/ai-context` 固定快照重写的 Hexas 工作流层。
- `ai/skills/zero-skills/`：基于 `zeromicro/zero-skills` 固定快照适配的按需知识层。
- `docs/audits/`：现状审计、遗留风险、重设计输入和处理状态。

## 验证

根模块：

```bash
go test ./...
```

`goctl` 独立模块：

```bash
cd tools/goctl
go test ./...
```

部分测试会启动本机临时 Redis、etcd、HTTP 或 gRPC 监听端口。完整开发流程和窄范围验证建议见 [ai/project-overview.md](ai/project-overview.md)。

## 许可证

框架代码沿用 [MIT License](LICENSE)。[社区行为准则](code-of-conduct.md) 是基于 Contributor Covenant 2.1 改编的 `CC-BY-4.0` 材料。上游来源、固定提交、修改范围和许可文本见 [Third-Party Notices](THIRD_PARTY_NOTICES.md) 与 [`LICENSES/`](LICENSES/)。
