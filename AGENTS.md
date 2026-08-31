# AGENTS.md

本文件是 FF-Hexas 的 AI 协作与默认开发流程入口。仓库是面向游戏服务场景维护的 go-zero 独立硬分支，不是 LemonGo 业务服务仓库，也不是 `zeromicro/go-zero` 的兼容镜像。

如果文档与代码、脚本、配置或 CI 行为不一致，以可执行内容为准，并在相关改动中同步修正文档。

## 1. 强制工作流

- 先列计划，得到确认后再修改；紧急小修也要先说明涉及文件和验证方式。
- 修改前读取相关实现、测试、配置、生成入口和本文列出的专项文档。
- 优先复用已有模块、接口、测试工具和脚本，不增加平行实现。
- 不确定时明确说明证据、假设和风险，不猜测。
- 中途需要改变已确认范围时，先报告并确认。
- 计划外问题只记录或报告，不顺手扩大改动。
- 未运行验证不算完成；无法验证时必须说明原因、影响和风险。
- 未经明确要求，不执行 `git commit`、`git tag` 或 `git push`。

## 2. 必读入口

开始工作时按顺序读取：

1. `AGENTS.md`：仓库级强制规则。
2. `ai/README.md`：AI 文档索引与任务路由。
3. `ai/project-overview.md`：目录职责、定制区域、编辑边界和验证入口。
4. `ai/framework-lineage.md`：基线提交、module path 和不兼容策略。
5. `docs/audits/2026-08-31-initial-framework-audit.md`：当前已知遗留问题。

修改审计中列出的高风险区域前，必须先阅读对应条目。审计条目不是自动授权的修复清单；用户只要求审计或范围未包含修复时，不得修改运行时行为。

获得专项改造授权后，审计记录只作为问题证据和设计输入，不构成对遗留实现的兼容要求。可以根据目标契约选择局部修复、整体重构、替换实现或删除能力。

## 3. 仓库身份与版本边界

- 代码基线：`github.com/zeromicro/go-zero v1.10.3`，提交为 `925f8a2bcc159eaf3b1da0f5fc695beac26e15ff`。
- 当前工作树已经包含面向游戏服务场景的定制；这些定制属于 FF-Hexas 初始基线，不按官方同版本行为推断。
- 根 module 仍为 `github.com/zeromicro/go-zero`，Go 版本为 `1.24.0`。
- `tools/goctl` 是独立 module；它在导入基线中直接依赖官方 `github.com/zeromicro/go-zero v1.10.3`。

保留 `github.com/zeromicro/go-zero` module/import path 是既有引用策略，不代表继续跟随官方发行。FF-Hexas 不承诺兼容官方后续分支、版本、API、配置或行为；不得为了“上游兼容”主动加入兼容层、双实现、自动 fallback 或同步逻辑。

后续若要合并、挑选或对照官方改动，必须作为独立任务明确目标提交、冲突处理、行为差异和验证范围。

## 4. 目录职责

- `core/`：基础能力、弹性治理、配置、日志、服务发现、存储、并发和进程生命周期。
- `rest/`：HTTP 服务、路由、中间件、客户端、请求解析和响应编码。
- `zrpc/`：gRPC 服务端、客户端、服务发现、负载均衡和拦截器。
- `gateway/`：API 网关。
- `mcp/`：MCP 服务。
- `internal/`：仅供本 module 内部使用的实现。
- `tools/goctl/`：独立 Go module 的生成工具和模板。
- `ai/`：AI 协作上下文和项目事实。
- `docs/audits/`：基线审计、遗留风险和验证记录。
- `.github/`：工作流、Issue 模板和 GitHub Copilot 入口。

## 5. 当前定制区域

相对官方 `v1.10.3`，当前导入基线主要定制了：

- `core/discov`：etcd 发布信息增加地址和服务名结构。
- `core/logx`：格式化日志附加字段、BI 原始 JSON 输出和日志级别调整。
- `core/service`：仅停止型 Service 适配。
- `core/stores/cache`：默认 TTL 调整和脏缓存写入。
- `core/stores/mon`、`core/stores/monc`：Database 封装与 Redis 到 MongoDB 的异步脏数据落盘。
- `core/stores/redis`：key prefix、Lua 原子操作、命令结果封装。
- `rest`：权限元数据、protobuf 请求/响应、HTTP client 和日志顺序/级别。
- `zrpc`：服务发现元数据、RPC 默认超时、负载均衡配置和公开别名。

这些实现已经作为基线代码保留，但不等于已完成设计或验证，也不要求后续原样保留。已知问题统一见初始审计，不得用静默降级掩盖未定义契约。

## 6. 编辑边界

- 优先修改现有文件；只有职责明确且无合适入口时才新增文件。
- 处理遗留问题时先定义目标契约，再选择修复、重构、替换或删除；不得因旧代码已经存在就默认延续其结构和行为。
- 重新设计不默认兼容遗留实现。涉及 FF-Hexas 已有调用方、存量数据、配置或部署时，必须说明影响、迁移、切换、验证和回滚；只有用户明确要求时才增加兼容层。
- 公共 API、配置字段、默认值、存储语义、服务发现编码、日志格式和生成模板都属于高风险变更。
- 修改接口时检查所有实现、mock、调用点和独立 module；不能只补当前编译错误。
- 修改配置语义时同步检查默认值、解析标签、示例、文档和启动行为。
- 修改 Redis key 处理时覆盖单节点、Cluster、多 key、Lua、pipeline 和 prefix 组合。
- 修改缓存或异步落盘时必须说明一致性模型、重试、幂等、进程退出、数据丢失窗口和监控方式。
- 修改服务发现值格式时必须说明部署顺序、混合版本行为和异常节点处理。
- 修改 REST `Route` 时检查所有复制/重建 Route 的 helper，确保元数据不丢失。
- 修改 `tools/goctl` 时在其目录内按独立 module 验证；根目录的 `go test ./...` 不覆盖该 module。
- 生成文件优先由现有生成命令重建，不直接手改，除非项目现状明确要求补丁。
- 不写死密钥、账号、私有地址、绝对路径、端口或环境差异逻辑。

## 7. 代码规范

- 遵循 Go 官方风格，使用 `gofmt`，保持 package 和导出符号注释准确。
- 阻塞或外部 I/O API 优先接收 `context.Context`，并真实响应取消。
- 错误必须被处理或有意返回；包装错误时保留 `%w`。
- 并发访问的 map、slice、全局回调和生命周期状态必须有明确同步策略。
- handler、resolver、adapter 保持薄层；可复用业务语义放在对应领域包。
- 新增行为必须配套正向、反例、失败路径及必要的并发测试。
- 不用 broad fallback、空返回、吞错或仅记录日志来伪装成功。

## 8. 文档规则

- `README.md` 只放项目定位、基线、重要入口和快速验证。
- 仓库事实、目录和工作流放在 `ai/`；历史来源放在 `ai/framework-lineage.md`。
- 已发现但本轮不处理的问题放在 `docs/audits/`，注明证据、影响、验证状态、重设计方向和建议顺序。
- 用户可见 API、配置、部署、生成方式或兼容边界变化时同步更新文档。
- 上游 `readme-cn.md`、`readme-ko.md` 和 `docs/upstream-readme.md` 只是导入快照，不作为本仓库开发规则。

## 9. 验证矩阵

文档和 AI 资产：

```bash
git diff --check
```

根 module：

```bash
go test ./...
go vet ./...
```

窄范围示例：

```bash
go test ./core/stores/redis/...
go test ./core/stores/mon/... ./core/stores/monc/...
go test ./core/discov/... ./zrpc/...
go test ./rest/...
go test -race <受影响包>
```

`goctl` 独立 module：

```bash
cd tools/goctl
go test ./...
go vet ./...
```

部分测试会监听本机临时端口，并可能需要 `protoc`、`protoc-gen-go` 和 `protoc-gen-go-grpc`。环境不满足时运行可执行的窄范围验证并明确缺口。

## 10. Review 与交付

Review 优先级依次为：数据一致性与安全、并发与生命周期、API/配置语义、回归风险、测试缺口、可维护性。

默认交付说明包含：

- 改了什么以及为什么。
- 验证命令和实际结果。
- 未验证项及原因。
- 已知风险和影响范围。
- 回滚方式。

未执行的测试不得描述为“通过”。未经用户明确要求，不提交或推送代码。
