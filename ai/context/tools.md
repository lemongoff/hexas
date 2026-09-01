# Hexas 工具入口

## 根 module

```bash
go test ./...
go vet ./...
```

按影响面缩小验证范围，涉及并发、共享 map、worker 或生命周期时补充 `go test -race`。

## 本地 goctl

`tools/goctl` 是独立 Go module。检查版本或构建本地生成器时从该目录运行：

```bash
cd tools/goctl
go run . --version
go test ./...
go vet ./...
GOWORK=off go build .
```

需要在临时示例中调用可执行文件时，先把当前源码构建到任务专用临时目录。不要执行：

```bash
go install github.com/lemongoff/hexas/tools/goctl@latest
```

`@latest` 可能引入与 Hexas 基线不同的模板、命令和运行时假设。

在已有服务中只更新 RPC 适配器和客户端时使用 `rpc protoc --skip-scaffold`，避免生成新的服务入口和 bootstrap YAML。新建完整服务时不使用该开关。

## 生成前检查

- 确认输入源是 `.api`、`.proto`、DDL 还是模板。
- 确认生成命令来自当前 `tools/goctl`，并记录参数和 style。
- 在临时目录试生成，检查将新增、覆盖或删除的文件。
- 不默认认为重复生成不会覆盖定制；以当前生成器实现和测试为准。
- 根 module 的 `go test ./...` 不覆盖 `tools/goctl`。

## 常用窄范围验证

```bash
go test ./core/stores/cache/... ./core/stores/redis/...
go test ./core/stores/mon/... ./core/stores/monc/...
go test ./core/discov/... ./zrpc/...
go test ./rest/...
go test ./core/logx/...
```

完整矩阵见 `ai/project-overview.md`。上游 goctl 命令参考只在需要对应功能时读取：
`../skills/zero-skills/upstream/references/goctl-commands.md`。
