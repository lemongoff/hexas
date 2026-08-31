# 框架来源与版本策略

## 1. 固定基线

本仓库在 2026-08-31 固定以下代码基线：

| 项目 | 值 |
| --- | --- |
| 项目 | `github.com/zeromicro/go-zero` |
| 版本 | `v1.10.3` |
| 基线提交 | `925f8a2bcc159eaf3b1da0f5fc695beac26e15ff` |

FF-Hexas 保留自己的 Git 仓库和远端，不复制其他仓库的 `.git` 历史。

## 2. 与官方 go-zero 的关系

官方 `github.com/zeromicro/go-zero v1.10.3` 对应提交
`925f8a2bcc159eaf3b1da0f5fc695beac26e15ff`。当前工作树在该基础上保留日志、Redis、MongoDB、服务生命周期、REST 和 RPC 等已有定制，因此不能把 FF-Hexas 视为官方标签的无差异副本。

## 3. module path 说明

根 `go.mod` 仍声明：

```go
module github.com/zeromicro/go-zero
```

这是有意保留的引用路径，避免在导入阶段重写全部 Go import。它不代表：

- FF-Hexas 由 zeromicro 官方发布或维护；
- FF-Hexas 与官方同版本号具有相同行为；
- 官方 issue、文档或升级指南可直接套用；
- 本仓库需要持续兼容官方新版本。

消费方应通过明确的 `replace`、工作区或内部发布版本引用 FF-Hexas，并把替换关系视为应用构建配置的一部分。

## 4. 不兼容策略

FF-Hexas 从本次导入起独立演进：

- 不承诺兼容官方后续分支、tag、API、配置、默认值、生成结果或运行时行为。
- 不设置自动跟随官方 `master` 的流程。
- 不为了官方兼容主动保留旧行为、增加双实现或静默 fallback。
- 如确需引入官方修复或能力，按单独任务固定源提交，评估内部定制冲突，并为目标行为编写测试。
- FF-Hexas 自身的稳定性和迁移策略应由内部版本、变更记录和调用方验证定义，不能借用官方版本号推断。

“不做后续兼容”并不意味着可以无说明地破坏 FF-Hexas 自己的调用方。公共契约变化仍需明确影响、迁移方式、验证和回滚方案。

## 5. goctl 的独立边界

`tools/goctl` 有独立 `go.mod`：

- module：`github.com/zeromicro/go-zero/tools/goctl`
- 直接依赖：官方 `github.com/zeromicro/go-zero v1.10.3`
- 当前 `BuildVersion`：`1.10.2`

因此根目录测试不会覆盖 goctl，goctl 构建默认也不会使用本地 FF-Hexas 根 module。若未来希望生成器与本地框架定制联动，需要单独设计 workspace/replace、版本命名和发布流程；当前不做调整。

## 6. 维护本文件

只有发生以下事件时更新代码基线：

- 明确重新建立官方版本或指定提交基线；
- 明确执行上游合并或 cherry-pick；
- 修改根 module path；
- 建立 FF-Hexas 自身版本与发布策略；
- 改变 goctl 与根 module 的依赖关系。

更新时必须记录精确 commit，不使用“最新”“当前版本”等可漂移描述。
