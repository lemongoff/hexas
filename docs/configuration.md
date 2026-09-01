# Hexas 配置体系

Hexas 使用独立 module [`github.com/lemongoff/hexas-config`](https://github.com/lemongoff/hexas-config) 提供类型化配置加载、校验、快照和动态来源。本仓库不保留原 `core/conf`、`core/configcenter` 或 go-zero 配置标签兼容层。

## 配置边界

配置分成两个独立 Manager：

- Bootstrap Config：服务监听、数据库、Redis、日志、TLS、服务发现和 tracing。来源是代码默认值、YAML 文件、环境 overlay、环境变量和命令行显式覆盖。
- Runtime Config：功能开关、活动参数、区服策略和临时运维值。可以通过独立的 etcd Source 加载和监听。

Runtime Manager 不得包含或覆盖 Bootstrap 字段，尤其不能修改监听端口、数据库凭据、TLS 或身份材料。

## 默认值

框架默认值由各领域的显式构造函数提供：

- `service.DefaultServiceConf`
- `rest.DefaultRestConf`
- `zrpc.DefaultRpcClientConf`
- `zrpc.DefaultRpcServerConf`
- `gateway.DefaultGatewayConf`
- `mcp.DefaultMcpConf`
- Redis、SQL、日志、Tracing 等包中的对应 `Default...` 函数

配置 struct 不再使用 `json:",default=..."`、`range` 或 `options` 标签。非零或非空配置值通过显式 `Validate` 检查；程序化调用可以使用零值表示未配置，配置文件加载则总是从完整默认实例开始。

## 目录规范

goctl 新生成项目使用：

```text
config/
├── base.yaml
└── env/
    ├── dev.yaml
    ├── test.yaml
    └── prod.yaml
```

当前生成器只创建 `config/base.yaml`。环境 overlay 由项目按部署环境增加；生产秘密不能写入配置文件或提交到仓库。

## 加载示例

```go
manager, err := hexasconfig.NewManager(
    appconfig.DefaultConfig(),
    hexasconfig.YAMLFile("config/base.yaml"),
    hexasconfig.OptionalYAMLFile("config/env/dev.yaml"),
    hexasconfig.Environment("HEXAS_"),
)
if err != nil {
    return err
}
if err := manager.Load(ctx); err != nil {
    return err
}
snapshot, ok := manager.Current()
if !ok {
    return errors.New("configuration was not published")
}
configuration := snapshot.Value()
```

来源按声明顺序合并，后声明的来源优先。未知字段、类型错误和业务校验失败都会阻止发布，并保留上一版有效快照。

## etcd Runtime Config

`github.com/lemongoff/hexas-config/source/etcd` 直接使用 etcd client，监听一个精确 key。调用方负责 client 生命周期，并必须配置有界 timeout。

建议 key：

```text
/hexas/{environment}/{game}/{service}/{cluster}/runtime
```

etcd 值统一为 YAML。Snapshot metadata 记录来源 revision 和 checksum；配置内容不得原样写入日志。

## 不兼容变化

- 删除 `core/conf` 和 `core/configcenter`。
- 删除 `MustLoad`、`FillDefault`、`UseEnv`、`LoadFrom*Bytes` 和 Properties API。
- 删除 JSON5、TOML 和 properties 配置入口。
- 删除默认值、范围和枚举反射标签。
- goctl 生成配置从 `etc/<service>.yaml` 改为 `config/base.yaml`。
- 配置错误由应用入口处理，配置库不会调用 `log.Fatal` 或退出进程。

不存在旧 API wrapper、路径别名、双实现或自动 fallback。

## 本地协同与发布

两仓尚未创建正式发布 tag 时，Hexas 使用明确的本地 `replace` 联调。发布前必须：

1. 为 `hexas-config` 创建正式版本；
2. 把 Hexas 和 goctl 的依赖更新为该版本；
3. 删除指向相邻目录的本地 `replace`；
4. 在不使用本地 workspace 的环境重新构建和测试生成项目。
