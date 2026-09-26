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

安装已发布版本时必须显式固定版本：

```bash
GOWORK=off go install github.com/lemongoff/hexas/tools/goctl@v1.10.4-hexas
```

在已有服务中只更新 RPC 适配器和客户端时使用 `rpc protoc --skip-scaffold`，避免生成新的服务入口和 bootstrap YAML。新建完整服务时不使用该开关。

### 项目模块识别

goctl 会先规范化生成目录，再识别其所属 module；在包含嵌套 module 的工作区内，选择最近的所属 module。`GOWORK=off` 时不读取或改写 `go.work`；启用工作区时，仅在所属 module 尚未登记时添加其根目录，识别已登记的 module 不重写工作区文件。

仅当目录不属于现有 module 或 GOPATH 项目时才执行 `go mod init`。已有 `go.mod`、`go.work` 或 Go 命令执行失败会直接返回错误，不通过创建子 module 掩盖问题。使用相对目录 `.` 创建新 module 时，默认名称取规范化后的目录名。

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
