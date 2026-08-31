# Contributing to FF-Hexas

感谢你参与 FF-Hexas。本文说明提交 Issue、修改代码和发起 Pull Request 的基本流程。

FF-Hexas 是基于 `github.com/JellyGoFF/FF-Hexas v1.10.3` 维护的游戏服务框架独立硬分支，不是 go-zero 官方发行版或兼容镜像。仓库保留 `github.com/JellyGoFF/FF-Hexas` module/import path，但不默认合并、升级或兼容官方后续版本。

参与贡献即表示你同意遵守 [Code of Conduct](code-of-conduct.md)。安全问题请按 [SECURITY.md](SECURITY.md) 私下报告，不要在公开 Issue 中披露漏洞细节。

## 开始之前

请先阅读：

- [README.md](README.md)：项目定位、固定基线和主要入口。
- [AGENTS.md](AGENTS.md)：仓库级开发、验证和交付规则。
- [项目概览](ai/project-overview.md)：目录职责、定制区域和验证建议。
- [框架来源关系](ai/framework-lineage.md)：module path、固定版本和不兼容策略。
- [框架默认值](docs/framework-defaults.md)：RPC、缓存和日志的默认行为。
- [初始框架审计](docs/audits/2026-08-31-initial-framework-audit.md)：历史问题、最终设计和迁移边界。

代码、配置、测试和 CI 的实际行为优先于文档。发现不一致时，应在同一改动中修正文档。

## 提交 Issue

提交缺陷或改进建议时，请提供足够的信息让维护者能够复现和判断范围：

- 使用的提交、Go 版本、操作系统和相关配置。
- 预期行为、实际行为及最小复现步骤。
- 完整错误信息、必要日志和受影响模块。
- 是否涉及公共 API、配置、持久化数据、Redis key、服务发现格式或生成代码。
- 已尝试的排查方式和可行的回滚路径。

功能建议应先说明目标契约和使用场景。不要仅以“与最新 go-zero 保持一致”作为升级理由；上游改动必须明确目标提交、行为差异、迁移影响和验证范围。

## 开发流程

1. 从 `main` 创建范围明确的主题分支。
2. 修改前阅读相关实现、测试、配置、生成入口和专项文档。
3. 优先复用现有模块和测试工具，避免增加职责重复的平行实现。
4. 只处理当前 Issue 或已确认计划内的内容；发现额外问题时单独记录。
5. 提交前运行最接近改动面的测试、静态检查和格式检查。
6. 发起 PR，并完整说明行为变化、验证结果、风险、迁移和回滚方式。

分支名应简短表达目的，例如 `fix/redis-prefix`、`feat/rest-protobuf` 或 `docs/contributing`。提交信息使用英文、采用祈使语气并保持单一职责，例如 `fix redis stream key prefix`。

不要在提交中包含密钥、账号、私有地址、环境专用绝对路径、日志数据或无关格式化改动。

## 代码与设计要求

- Go 代码必须通过 `gofmt`；导出符号、错误语义和生命周期注释应与实现一致。
- 包装错误时使用 `%w` 保留原因；不能用空返回、吞错或静默 fallback 掩盖未定义契约。
- 阻塞操作和外部 I/O 应接收并响应 `context.Context`。
- 并发访问的 map、slice、全局回调和生命周期状态必须有明确同步策略。
- 新增行为应覆盖正常、反例、失败路径；并发或生命周期改动应增加 race 测试。
- 生成文件优先从源定义重新生成，不直接维护第二份生成结果。
- 修改公共 API、配置或默认值时，同步更新调用点、mock、示例、测试和文档。

本项目不为官方后续 go-zero 版本自动增加兼容层。需要兼容旧调用方或存量数据时，必须在 PR 中明确迁移期限、切换方法和移除条件。

## 高风险改动

以下改动需要在 PR 描述中单独列出目标契约、失败模型和回滚方式：

- Redis：prefix、pattern、pipeline、Lua、Streams、Cluster 和多 key 操作。Cluster 多 key 必须使用共享 `{hash-tag}`。
- MongoDB 与缓存：client 所有权、索引、缓存失效、持久化顺序、幂等和数据丢失窗口。
- 服务发现和 RPC：注册值格式、混合部署、异常节点隔离、超时与阻塞连接默认值。
- REST：Route 权限元数据、中间件组合、请求大小限制和 Protobuf content type。
- 日志：字段 schema、敏感值遮罩、截断、日志级别和外部 Writer。
- `tools/goctl`：模板、生成结果、protobuf 工具版本和根 module 的联动。

真实 Redis Cluster、MongoDB 或故障注入环境不可用时，应运行可执行的单元测试，并明确未覆盖的集成验证及上线风险。

## 验证

所有改动至少运行：

```bash
git diff --check
```

根 module 的标准验证：

```bash
go mod verify
go vet ./...
go test ./...
```

并发或高风险包应增加：

```bash
go test -race <受影响包>
```

`tools/goctl` 是独立 module，相关改动必须在其目录单独验证：

```bash
cd tools/goctl
go mod verify
go vet ./...
go test ./...
```

goctl 的完整生成器测试需要 `protoc`、`protoc-gen-go v1.36.11` 和 `protoc-gen-go-grpc v1.5.1`。不得以根目录测试代替 goctl module 测试。

只修改文档时不要求运行 Go 测试，但必须检查路径、链接、命令和 `git diff --check`。无法执行某项验证时，应在 PR 中写明原因和风险，不能标记为已通过。

## Pull Request

PR 应保持范围集中，并包含：

- 问题和目标契约。
- 主要实现选择及拒绝其他方案的原因。
- 用户可见行为、API、配置、数据或部署影响。
- 实际执行的验证命令和结果。
- 未验证项及剩余风险。
- 迁移、发布顺序和回滚方式。
- 对应 Issue、审计项或设计文档。

不要混入无关重构、依赖升级或整仓格式化。大改动应拆成可独立审查、可独立回滚的逻辑提交，但不能用多个局部补丁规避对同一业务语义的整体检查。

Review 反馈应覆盖同一语义涉及的实现、测试、配置、文档和生成入口。CI 通过是合并的必要条件，但不能替代对数据一致性、安全边界和兼容影响的审查。

## 许可证

提交到本仓库的内容按项目 [MIT License](LICENSE) 分发。提交贡献即表示你确认有权提交相关代码、文档和资源，并同意按该 MIT License 授权你的贡献。贡献中包含第三方材料时，必须保留必要的版权与许可证声明，并更新 [Third-Party Notices](THIRD_PARTY_NOTICES.md)。
