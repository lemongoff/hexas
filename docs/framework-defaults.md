# Hexas 框架默认值

本文件记录 Hexas 基于 `zeromicro/go-zero v1.10.3` 固定分支形成的运行时默认契约。项目不跟随 go-zero 后续版本，也不以官方默认值推断本分支行为。

配置默认值不再来自 `json:",default=..."` 反射标签，而由各领域的 `Default...` 构造函数显式提供。配置加载和动态快照由 `github.com/lemongoff/hexas-config` 管理，详见 [配置体系](configuration.md)。

## RPC

- RPC 客户端与服务端默认超时均为 5000 毫秒；服务端显式配置 `Timeout: 0` 才表示不设超时。
- RPC 客户端 `NonBlock` 默认 `false`，因此创建连接时默认等待连接结果；需要异步建连时必须显式配置 `NonBlock: true`。
- 方法级超时仍可通过 `MethodTimeouts` 或调用选项覆盖全局值。
- 默认 `p2c_ewma` 识别 `zrpc.WithRouteTarget`。`RouteRequire` 找不到指定实例时明确失败，`RoutePrefer` 才允许回退到普通 P2C；详见 [RPC instance routing](rpc-instance-routing.md)。
- Etcd RPC 实例必须发布 `ready` 或 `draining`。普通 P2C 排除 `draining`，显式实例路由仍可用于完成受 fencing 保护的状态迁移。
- RPC Server 先监听再发布 `ready`，`BeginDrain` 原租约更新为 `draining`，`Stop` 优雅停止并同步撤销租约。RPC Client 由调用方长期持有并显式 `Close`。

## 缓存

- 正常缓存默认 TTL 为 24 小时。
- 未找到占位值默认 TTL 为 1 分钟。
- 需要其他生命周期时必须显式使用对应配置或 `WithExpiry` / `WithNotFoundExpiry`。

## 成功日志

- REST 与 RPC 成功请求按 debug 级别记录，以限制生产环境的常规成功日志量。
- 错误、慢请求和统计日志仍按各自级别输出；排查成功请求时需要临时启用 debug 日志。

## Windows 与文件资源

- `proc.ShutdownConf` 在各平台具有相同字段，默认 `WrapUpTime` 为 1 秒、`WaitTime` 为 5.5 秒；Windows 仍不执行 Unix 信号驱动的定时停服，调用方需自行管理停止流程。
- 大小轮转的默认文件名时间格式在 Linux/macOS 保持 RFC3339；Windows 使用 `2006-01-02T15-04-05Z07-00`，避免文件名包含冒号。显式设置 `FileTimeFormat` 时，调用方负责使用平台合法的格式。日志内容时间格式不变。
- Windows 旧的默认大小轮转文件名无法正常创建，无需迁移；已有自定义格式不变。日志目录、压缩和保留策略不变。回滚此修复会重新引入 Windows 默认轮转失败。
- Trace 的文件导出器拥有输出文件，停止时在导出完成后关闭文件；重复停止安全，停止后的导出不再写入。非文件导出器行为不变。
- `syncx.WithRefreshIntervalOnFailure(0)` 在每次失败后立即允许重试，不依赖系统时钟前进；成功结果仍然缓存。

## 变更要求

修改上述默认值属于行为变更，必须同步更新配置示例、受影响测试和本文件，并在发布说明中写明回滚值。
