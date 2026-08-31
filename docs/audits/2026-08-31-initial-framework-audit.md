# FF-Hexas 初始框架基线审计

日期：2026-08-31
状态：专项处置完成；初始发现作为历史证据保留，最终设计见 3.1 节

## 1. 审计范围

本审计用于固定 FF-Hexas 初始导入基线中的已有定制和遗留风险，不代表已经完成专项设计或生产验收。

比较范围：

- 当前版本：FF-Hexas 初始游戏化框架工作树
- 官方基线：`zeromicro/go-zero v1.10.3@925f8a2bcc159eaf3b1da0f5fc695beac26e15ff`
- 差异规模：47 个文件，约 1278 行新增、169 行删除
- 审计方式：Git 历史与 diff 对照、关键路径静态检查、根 module 与 goctl 独立 module 测试

本文第 4 节保留初始导入时的原始发现，用于解释改造动机；这些描述不再代表当前实现状态。

### 1.1 遗留问题处置原则

- 审计描述的是初始代码现状、风险和反例，不把遗留结构或行为确立为兼容契约。
- 获得专项改造授权后，每个问题都可以选择局部修复、整体重构、替换实现或删除能力，不要求沿用当前代码路径。
- 高风险项应先定义目标契约、一致性或安全边界和失败模型，再设计实现与测试；不要默认在现有错误分支上叠加补丁。
- 不要求兼容官方后续版本，也不默认兼容遗留定制。只有用户明确要求时才为旧行为增加过渡层。
- 涉及 FF-Hexas 已有调用方、存量数据、配置或部署时，重新设计仍需给出影响评估、迁移或切换方案、验收测试和回滚路径。
- 完成专项处理后更新对应审计项状态，记录最终选择、被替换或删除的能力及验证证据。

## 2. 风险等级

- **高**：可能导致数据丢失、持久化错误、服务不可用、安全边界失效或生产迁移失败。
- **中**：可能造成语义不一致、并发问题、观测缺失、资源泄漏或调用方构建/行为变化。
- **低**：主要是可维护性、命名、文档或局部测试完整性问题。

## 3. 审计结论摘要

| 编号 | 等级 | 区域 | 结论 |
| --- | --- | --- | --- |
| A-01 | 高 | Mongo 脏数据落盘 | 已关闭：删除未被业务使用且无法证明一致性的异步脏写链路 |
| A-02 | 高 | Mongo 脏数据落盘 | 已关闭：随脏写 worker 一并删除，不保留不安全生命周期契约 |
| A-03 | 高 | Redis prefix/Lua | 已关闭：补齐复合命令和 Streams 前缀；Cluster 多 key 强制共享 hash tag |
| A-04 | 中 | Mongo Database | 已关闭：Database 独占 client，collection cache 加锁并显式返回索引冲突 |
| A-05 | 高 | etcd/zRPC | 已关闭：保持 JSON-only，隔离坏值和空地址并继续构造其余节点 |
| A-06 | 高 | REST 权限 | 已关闭：Route helper 复制完整结构并保留 `Permissions` |
| A-07 | 中 | REST protobuf | 已关闭：补齐 protobuf error/write API、8 MiB 限制、content-type 和错误链 |
| A-08 | 中 | logx | 已关闭：BI 改为可选 writer 能力并进入标准日志字段、脱敏和截断路径 |
| A-09 | 中 | goctl | 已关闭：根 `go.work` 绑定本地 module，版本更新为 `1.10.3-ffhexas` |
| A-10 | 中 | 默认行为 | 已关闭：新增 FF-Hexas 默认行为契约文档并保留对应测试 |
| A-11 | 中 | 测试覆盖 | 已关闭：新增关键失败、元数据、前缀、Cluster、并发与 BI 回归测试 |
| A-12 | 高 | GitHub workflow | 已关闭：CI 对齐 `main` 和双 module，删除上游 goctl 发布/版本任务及 gomod 自动升级 |
| A-13 | 低 | 基线格式 | 已接受：CI 只检查本次变更行，不批量改写导入快照 |

### 3.1 最终设计与迁移结论

- Mongo 持久化采用同步写入并在成功后失效缓存；删除 `ModelDirtyWorker`、`SetWithDirtyCtx`、脏标记 Lua 和 `UpdateOneWithCacheDirty`。仓库内没有这些 API 的业务调用方，因此无需数据迁移；外部调用方必须改用 `UpdateOne` 等同步路径。
- `Database` 不再借用全局共享 client：每个实例创建、持有并关闭自己的 client；collection lazy cache 使用互斥锁。`MustNewDatabase` 删除未使用的 collection 参数，未消费的 `DatabaseOptions` API 被删除，索引选项冲突直接返回错误。
- Redis 对公开 API 统一接收逻辑 key。`Keys`/`Scan` 自动限定前缀并返回去前缀的逻辑 key；复合 destination、MSET key、ZStore key 和 Stream key 均加前缀。Cluster 多 key 操作必须让所有 key 使用同一非空 `{hash-tag}`，否则返回 `ErrCrossSlot`。
- etcd 服务发现坚持 JSON-only 硬切换，不恢复旧纯地址格式兼容；单条损坏或空地址记录只被隔离，不再截断合法节点列表。
- REST Route 的元数据在 helper 组合中完整保留。Protobuf HTTP 使用 `application/pb`、8 MiB 请求上限和独立 success/error 写入入口。
- `logx.Writer` 不再强制实现 BI；实现 `BIWriter` 可选择独立路由，否则落到 Info。两条路径都使用 caller、全局字段、`channel=bi`、敏感值遮罩和内容截断。
- goctl 通过根 `go.work` 使用当前框架源码，显示版本为 `1.10.3-ffhexas`。本仓库不发布上游 goctl tag，也不自动跟随 go-zero 后续版本。
- RPC、缓存和成功日志的默认值见 [`../framework-defaults.md`](../framework-defaults.md)。

回滚时应按条目独立回退代码和测试；不要恢复已删除的异步脏写链路。若必须重新引入异步持久化，应作为新设计完成 outbox/ack、幂等、崩溃恢复和真实故障注入验证。

## 4. 详细发现

### A-01：脏数据失败后仍可能清除脏标记

等级：高

证据：`core/stores/monc/modeldirtyworker.go:67-124`。

当前 worker 先用 `RPOP COUNT` 从队列移除 key，再批量读取 Redis 并生成 MongoDB update：

- 单个 key 在 `persist` 中失败时只记录日志并 `continue`，没有重新入队。
- 某个 collection 的 `BulkWrite` 失败时会把整批 key 重新入队，但循环结束后仍无条件对整批 key 执行 `SREM dirty:set`。
- Redis/Mongo 操作的重入队、清标记和落库之间没有事务或可证明的 at-least-once 协议。
- 进程在出队后、落库前退出时，队列项已经移除。

影响：缓存中的最新状态可能未写入 MongoDB，却失去脏标记或可靠重试路径；缓存过期后可能形成持久化数据回退。

后续建议：先定义一致性模型和状态机，再设计 ack/retry/dead-letter、幂等版本和崩溃恢复测试。不要只在当前错误分支追加局部 `LPUSH`。

### A-02：脏数据 worker 的输入与生命周期不安全

等级：高

证据：`core/stores/monc/modeldirtyworker.go:59-76`、`:128-144`。

- `strings.Split` 后直接访问索引 2、3、4，短 key 会触发越界 panic；key/prefix 中额外冒号也可能改变字段语义。
- `ModelPoolPackage.Get()` 取出的对象从未调用 `Put()` 归还。
- worker 是无限循环，sleep 不监听 `ctx.Done()`，取消后仍会继续重试。
- 休眠计算使用 `time.Now().Second()`，它只返回分钟内秒数；跨分钟时差值会失真，负 duration 会立即返回并可能形成忙轮询。
- `AddModel` 写普通 map，worker 同时读取时没有同步约束。

影响：恶意或损坏 key 可终止 worker；停服过程不可控；对象池失效；并发注册可能触发 data race。

后续建议：把 key schema 变为显式编解码契约，worker 接入可取消 timer、同步注册和对象归还，并补 panic、取消、跨分钟与 race 测试。

### A-03：Redis prefix 与 Cluster 契约不完整

等级：高

证据：

- prefix 入口：`core/stores/redis/redis.go:173-186`
- 脏标记脚本：`core/stores/redis/redis.go:216-229`
- 未处理 key 的命令：`MsetCtx` 位于 `:1261-1268`，`PfmergeCtx` 位于 `:1327-1334`，`SunionstoreCtx` 位于 `:1742-1749`
- `ZunionstoreCtx` 只处理 destination，未处理 `store.Keys`：`:2741-2749`
- `KeysCtx`、`ScanCtx` 不加 prefix，也不剥离返回值：`:1083-1095`、`:1471-1484`

当前大部分单 key 命令会加 prefix，但多 key、pattern、destination 和复合参数并不一致。调用方无法从统一契约判断输入输出是逻辑 key 还是物理 key。

脏标记 Lua 同时使用 `dirty:set`、`dirty:queue` 和数据 key。在 Redis Cluster 中：

- 两个 `KEYS` 没有共同 hash tag，可能直接产生 `CROSSSLOT`。
- 数据 key 作为 `ARGV` 传入但脚本内部访问，仍缺少与执行节点同 slot 的保证。

影响：开启 prefix 后部分命令读写到不同命名空间；Cluster 模式下脏写链路可能不可用。

后续建议：先列出全部 Redis API 的 key 参数模型，定义 prefix 和物理 key 逃生口；Cluster Lua 使用明确 hash tag/单 slot 设计，配真实 Cluster 集成测试。现有 miniredis 测试不足以证明该行为。

### A-04：Mongo Database 并发与 client 所有权不明确

等级：中

证据：`core/stores/mon/database.go:25-81`、`core/stores/mon/clientmanager.go`。

- `Database.Collection` 对普通 map 做无锁 lazy read/write，多 goroutine 首次访问 collection 时可能产生 data race 或并发 map 写。
- `getClient` 通过全局 resource manager 按 URI 共享 client；`Database.Close` 却直接断开底层 client，可能影响同 URI 的其他 Model/Database。
- `MustNewDatabase` 的 `collection` 参数未使用。
- `CreateIndex` 通过错误字符串包含 `with different options` 判断并返回成功，实际索引配置不一致仍会被隐藏。
- `DatabaseOptions`、`WithIndex` 等结构当前没有被执行路径消费。

影响：并发启动和停服存在不确定行为；索引声明可能与数据库实际状态漂移。

后续建议：明确 client 的引用计数/进程级所有权，collection cache 使用同步结构；索引差异应返回结构化结果或执行显式迁移。

### A-05：etcd 注册格式是硬切换，坏值会截断节点列表

等级：高

证据：

- 新 JSON 类型：`core/discov/publisher.go`
- 服务端发布：`zrpc/internal/rpcpubserver.go`
- resolver 解码：`zrpc/resolver/internal/discovbuilder.go:24-43`

RPC server 不再发布纯 `host:port`，而是发布 `{"Addr":...,"ServerName":...}`。resolver 对每个值强制 JSON 解码；遇到第一个旧格式或损坏值时执行 `break`，不会继续处理后面的合法节点，随后仍用部分地址调用 `UpdateState`。

影响：滚动部署、etcd 残留旧值或单个异常节点可能让客户端看到空集或不完整节点集，造成流量集中或不可用。

这不是要求恢复官方兼容；即使采用一次性硬切换，也需要部署前清理、版本闸门和坏节点隔离策略。

后续建议：明确原子切换流程；单节点解码失败至少应隔离该节点并继续处理其他值，同时增加格式、空地址和混合集合测试。

### A-06：REST Route helper 会丢失权限元数据

等级：高

证据：`rest/types.go` 为 `Route` 增加 `Permissions`；`rest/server.go:221-250` 和 `:268-280` 重建 Route 时未复制该字段。

- `WithMiddleware` 返回的新 Route 只保留 Method、Path、Handler。
- `WithPermissionsMiddleware` 应用权限后也丢弃字段。
- `WithPrefix` 同样丢弃权限。

影响：helper 组合顺序会改变权限中间件收到的数据。若先调用普通 middleware 或 prefix 再应用权限 middleware，权限列表可能变为空，形成授权绕过风险。

后续建议：定义 Route 元数据复制方式，避免多个 helper 手写字段；测试不同 helper 顺序和空/多权限场景。

### A-07：protobuf HTTP 支持未闭环

等级：中

证据：

- `rest/httpx/responses.go:19-25`、`:99-115`
- `rest/httpx/requests.go:128-139`

`SetErrorPbHandlerCtx` 会写入全局 `errorPbHandler`，但当前没有任何响应函数读取该变量，因此配置不会生效。`ParsePbBody` 使用 `io.ReadAll`，没有 JSON 解析路径的 8 MiB 限制；错误消息写成 `json to proto failed`，且未使用 `%w` 保留底层原因。

影响：调用方可能误以为 protobuf error handler 已启用；大请求可造成不受控内存占用；排障信息不足。

后续建议：先确定完整的 protobuf success/error API 和 content-type 契约，再实现统一 body limit、错误包装和测试。

### A-08：日志接口和 BI 输出形成新的独立契约

等级：中

证据：`core/logx/writer.go:20-44`、`:279-309`。

- `Writer` 接口新增 `Bi(any)`，所有外部自定义 Writer 都必须实现该方法才能重新编译。
- concrete `Bi` 直接 `writeJson` 到 info 输出，不经过普通 `output` 路径，因此不自动附加 timestamp、level、caller、global fields，也不执行普通字段的 sensitive masking 和 content 截断。
- 格式化日志只识别 variadic 参数末尾的 `*LogFields`，这是额外的隐式调用约定。

影响：日志消费者、脱敏和 schema 需要按 BI 独立评估；外部 Writer 是有意的硬破坏点。

后续建议：把 BI schema、必填字段、脱敏责任、写入失败和 Writer 扩展方式定义成显式契约。由于本项目不承诺官方兼容，不需要为官方 Writer API 增加兼容层，但需要管理 FF-Hexas 自身调用方迁移。

### A-09：goctl 与根框架不在同一依赖轨道

等级：中

证据：

- `tools/goctl/go.mod:1-23` 直接依赖官方 `github.com/zeromicro/go-zero v1.10.3`
- `tools/goctl/internal/version/version.go:8-9` 的 `BuildVersion` 为 `1.10.2`

`tools/goctl` 是独立 module。根目录 `go test ./...` 不包含它；在 goctl 目录构建时，默认下载官方 v1.10.3，而不是使用本地 FF-Hexas 根 module。

影响：生成器测试通过不能证明其模板与 FF-Hexas 定制一致；版本显示也不能代表当前导入基线。

后续建议：在需要游戏化生成能力时单独设计本地 workspace/replace、模板来源和 FF-Hexas 版本命名。本轮保持原样。

### A-10：已有默认行为变更缺少迁移契约

等级：中

证据：

- `zrpc/config.go:23-54`：client/server timeout 为 5000ms；`NonBlock` 从官方默认 true 变为 optional false。
- `core/stores/cache/cacheopt.go:5-8`：默认缓存 TTL 为一天，而官方基线为七天。
- HTTP/RPC 成功日志由 info 调整为 debug，REST log middleware 顺序发生变化。

影响：未显式配置的调用方会改变启动阻塞、调用超时、缓存回源频率和观测量。它们可能是预期定制，但目前缺少 FF-Hexas 自身的配置和迁移说明。

后续建议：建立 FF-Hexas 配置参考和变更记录；不要以官方默认值推断本仓库行为。

### A-11：现有测试通过但未覆盖关键定制反例

等级：中

根 module 和 goctl 独立 module 的现有测试均能通过，但静态检查确认以下缺口：

- 没有 `modeldirtyworker_test.go`，未验证失败重试、崩溃窗口、取消和 key schema。
- Redis prefix 没有覆盖全部命令，也没有真实 Cluster 多 key Lua 测试。
- 新 Database 缺少并发 Collection、共享 client Close 和索引冲突测试。
- resolver 没有覆盖旧格式、损坏值位于列表中间和空地址。
- REST helper 没有验证 `Permissions` 在 middleware/prefix 组合后保持。
- protobuf error handler 未被测试为实际响应路径。
- BI 没有 schema、脱敏和自定义 Writer 迁移测试。
- goctl module 的 `go vet ./...` 不能通过：存在无缓冲 signal channel 用法和 2 处不可达代码。
- `tools/goctl/api/gogen` 的测试在包目录复用并删除已跟踪的 `jwt.api`，测试结束会污染工作树；本次验证后已按源快照恢复该文件。

影响：全量绿测只能证明既有测试集通过，不能关闭 A-01 至 A-10；goctl 也尚未达到静态检查全绿，且测试不具备完整的工作树隔离性。

### A-12：复制的 GitHub workflow 仍假设上游仓库

等级：高

证据：

- `.github/workflows/go.yml` 和 `codeql-analysis.yml` 只监听 `master`，而 FF-Hexas 当前主分支为 `main`。
- `.github/dependabot.yml` 每日更新根 module 和 goctl module，可能持续引入官方依赖变化，与“不默认跟进上游”冲突。
- `.github/workflows/release.yaml` 在 `tools/goctl/*` tag 上使用 `zeromicro/go-zero-release-action@master` 发布 goctl，并固定下载 Go 1.21.13；当前两个 go.mod 均要求 Go 1.24.0。
- `.github/workflows/version-check.yml` 同样使用 Go 1.21，且当前 goctl `BuildVersion` 仍为 1.10.2。

影响：`main` 的 push/PR 可能没有预期的测试和 CodeQL 闸门；创建 goctl tag 可能触发不符合 FF-Hexas 发布边界的自动发布或因 Go 版本不匹配失败；Dependabot 会制造未经专项审计的版本漂移。

后续建议：在单独 CI/发布治理任务中决定保留、禁用或重写这些 workflow。完成前不要创建 `tools/goctl/*` tag，并把本地验证作为必需闸门。本轮按“遗留问题只审计”要求保持原文件不变。

### A-13：导入快照含既有空白格式告警

等级：低

证据：对相对导入前提交的完整工作树执行 `git diff HEAD --check` 时，`.github/ISSUE_TEMPLATE/`、`core/`、`gateway/` 和 `tools/goctl/` 等上游文件报告尾随空白、空白缩进或文件末尾多余空行。

影响：若 CI 直接对整个首次导入 diff 执行该检查会失败；这些格式问题主要影响维护质量，批量清理则会制造与来源快照无关的大量差异。

后续建议：首次导入保留源快照；后续修改文件时局部清理触及行，或在独立格式治理任务中统一处理。本轮不批量改写。

## 5. 已有定制与待评估能力清单

初始基线包含以下主要定制。它们用于定位影响面，不表示后续必须保留；专项设计应逐项决定保留并修复、重构、替换或删除：

- `core/discov`：发布地址与 `ServerName` JSON 编解码。
- `core/logx`：`LogFields`、格式化日志字段、BI 输出、HTTP/RPC 日志级别。
- `core/service`：`WithStop`。
- `core/stores/cache`：一天默认 TTL、`SetWithDirtyCtx`。
- `core/stores/mon`：Database/Index 封装、Mongo 日志调整。
- `core/stores/monc`：dirty worker 和 cache dirty 更新入口。
- `core/stores/redis`：prefix、GetOrSet/CAS/CompareAndDel/ZCompareHigher、dirty Lua、`CmdResult`。
- `rest`：权限字段、protobuf 请求响应、GET/POST helper、中间件顺序和日志级别。
- `zrpc`：JSON 服务发现值、ServerName、五秒默认超时、Balancer 别名和公开 ClientOptions/Acceptable。

## 6. 重新设计优先级

本节仅给出未来专项排序，不在基线导入阶段实施。每个专项开始时先输出目标契约和“保留修复 / 重构 / 替换 / 删除”决策：

1. 重新定义异步持久化一致性模型与 Redis key/Cluster 边界，处理 A-01、A-02、A-03；现有 dirty worker 和 Lua 可以整体替换。
2. 重新定义授权元数据与服务发现异常隔离契约，处理 A-05、A-06；不以保持现有 helper 或编码结构为目标。
3. 明确 Mongo client 所有权、collection 生命周期和并发模型，处理 A-04；必要时替换当前 Database 封装。
4. 治理 CI、Dependabot 和发布触发条件，处理 A-12，避免主分支无闸门或误发布。
5. 定义 protobuf、日志和配置目标契约，处理 A-07、A-08、A-10；无实际消费方的能力可以删除。
6. 决定 goctl 是否进入 FF-Hexas 发布轨道，处理 A-09；不需要时可从框架发布范围移除。
7. 每个专项都以 A-11 中的反例和目标契约测试作为完成条件，并同步关闭、替换或拆分对应审计项。

## 7. 本次验证记录

环境：macOS arm64，Go `1.24.0`。

已执行：

```bash
go test ./...
```

结果：根 module 全部通过。测试需要本机回环监听权限。

已执行：

```bash
cd tools/goctl
go test ./...
```

结果：在隔离的临时 GOPATH/GOBIN 中补齐 `protoc-gen-go` 和 `protoc-gen-go-grpc` 后全部通过。首轮因测试尝试向只读 Go toolchain bin 安装插件而失败，属于环境路径问题；重跑已验证通过。

静态检查与文档检查：

```bash
go vet ./...
cd tools/goctl && go vet ./...
git diff --check
git diff HEAD --check
```

结果：

- 根 module `go vet ./...` 通过。
- `tools/goctl` 的 `go vet ./...` 未通过：`migrate/cancel.go` 把无缓冲 `os.Signal` channel 传给 `signal.Notify`；`test/test.go` 有 2 处不可达代码。这些属于本次复制基线的遗留静态检查问题，按本轮边界仅记录、不修改。
- 本轮新增或重写的项目文档未发现尾随空白，当前未暂存文档差异的 `git diff --check` 通过。
- 完整导入基线的 `git diff HEAD --check` 未通过，均为来源快照已有的空白格式告警，详见 A-13；本轮未批量修改。
- 逐文件比对导入源的 1387 个 tracked files：除有意重写的 `readme.md`（在大小写不敏感文件系统上对应本仓库 `README.md`）和 `.github/copilot-instructions.md` 外，内容差异为 0；原英文 README 已保存为 `docs/upstream-readme.md`。

## 8. 回滚说明

本次没有修改上述运行时代码。若需要回滚本次导入，应整体撤销初始框架工作树并恢复 FF-Hexas 导入前的初始提交；不要逐项回滚审计中列出的既有定制，因为它们属于初始基线本身。
