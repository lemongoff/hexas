# 项目概览与开发入口

## 1. 项目定位

FF-Hexas 是以 go-zero `v1.10.3` 为代码基线的游戏服务框架。当前阶段的重点是固定一份可审计的初始状态，保留已有定制，再按明确任务逐步完成游戏化改造。

根 module：

```text
github.com/zeromicro/go-zero
```

Go 版本：`1.24.0`。

module path 只是现有 import path，不表示本仓库继续作为官方 go-zero 的兼容实现。来源和版本边界见 `framework-lineage.md`。

## 2. 代码结构

| 目录 | 职责 | 常见高风险 |
| --- | --- | --- |
| `core/` | 配置、日志、服务发现、弹性治理、存储、并发、生命周期 | 数据一致性、并发、默认值、全局状态 |
| `rest/` | HTTP 服务、路由、中间件、客户端和编解码 | 权限元数据、中间件顺序、body 限制、响应契约 |
| `zrpc/` | gRPC 服务和客户端、etcd resolver、负载均衡 | 注册格式、混合部署、超时和连接方式 |
| `gateway/` | HTTP/gRPC 网关 | 协议适配、路由和流量保护 |
| `mcp/` | MCP 服务 | HTTP 生命周期和协议行为 |
| `internal/` | 根 module 内部实现 | 只允许 module 内引用 |
| `tools/goctl/` | 独立 module 的代码生成工具 | 模板、生成结果、外部 protoc 工具链 |

## 3. 导入基线中的定制

相对官方 `v1.10.3`，FF-Hexas 初始基线差异为 47 个文件，约 1278 行新增、169 行删除。定制主要包括：

- etcd 发布值由纯地址扩展为带 `ServerName` 的 JSON。
- 日志格式化 API 支持末尾 `LogFields`，增加 BI JSON 输出，并降低部分 HTTP/RPC 成功日志级别。
- Redis 增加 key prefix、多种 Lua 原子命令和 `CmdResult`。
- MongoDB 增加 Database 封装，Mongo cache 增加 Redis 脏队列异步落盘。
- Cache 默认 TTL 从七天调整为一天，并增加带脏标记写入。
- REST 增加权限字段、protobuf 请求/响应和简化的 GET/POST client。
- RPC 默认超时调整为五秒，并扩展 resolver 地址元数据和公开配置别名。

这些是“已存在的实现”，不是“已验证完成的设计”，也不是后续必须保留的兼容契约。修改前先读初始审计；获得专项授权后，可以按目标能力重新设计、替换或删除。

## 4. 编辑边界

- 当前代码快照只是后续改造起点；遗留定制不要求原样保留，也不要为了看起来接近官方而机械回退。
- 处理遗留问题时先定义目标契约，再决定局部修复、整体重构、替换实现或删除能力；不要默认围绕旧实现追加补丁。
- 重新设计不要求兼容遗留实现，但必须核对 FF-Hexas 已有调用点、存量数据、配置和部署方式，并给出必要的迁移、切换、验证和回滚方案。
- 不主动实现官方未来版本兼容，也不默认合并官方 `master`。
- 任何上游同步都必须固定来源 commit，单独审计差异和冲突。
- 数据写入、缓存、Redis Lua、服务发现编码、权限和配置默认值属于高风险区域。
- `Route`、配置 struct、Writer 等公共类型发生字段或方法变化时，检查所有构造、复制、mock 和外部实现。
- `tools/goctl` 不在根 module 的 `./...` 范围内，必须独立验证。

## 5. 验证入口

根 module 全量验证：

```bash
go test ./...
go vet ./...
```

高风险改动按影响面补充 race 测试：

```bash
go test -race ./core/stores/redis/...
go test -race ./core/stores/mon/... ./core/stores/monc/...
go test -race ./core/discov/... ./zrpc/...
go test -race ./rest/...
```

`goctl` 独立 module：

```bash
cd tools/goctl
go test ./...
go vet ./...
```

`goctl/rpc/generator` 测试需要 `protoc`、`protoc-gen-go` 和 `protoc-gen-go-grpc`。多个根模块测试会监听回环端口，用于 miniredis、etcd mock、HTTP 和 gRPC 测试服务。

文档与 AI 资产：

```bash
git diff --check
```

## 6. 交付要求

交付时区分：

- 本轮实际修改。
- 导入基线中原本就存在的行为。
- 已验证通过的命令。
- 因环境未执行的验证。
- 仅记录、尚未处理的审计项，以及已确定的重设计决策。

没有运行的验证不能描述为通过。审计项的建议不能描述为已实现。
