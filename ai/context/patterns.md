# Hexas 框架模式

本文件只列会改变框架维护决策的项目模式。业务服务示例模式按需读取 `ai/skills/zero-skills`。

## 公共契约

- module path 与官方相同不代表行为相同；修改公开类型、接口、默认值或配置时检查仓库内全部构造、实现、mock 和生成模板。
- 不为官方后续版本自动增加兼容层。Hexas 调用方迁移仍需明确说明。
- adapter、resolver 和 handler 保持薄层，但框架内部不强制套用业务项目的 Handler/Logic/Model 目录结构。

## Context 与生命周期

- 外部 I/O、阻塞和后台 worker 应真实响应 `context.Context` 取消。
- 明确 goroutine、timer、client、连接池和全局注册项的创建者与关闭者。
- 并发 map、回调和状态切换需要同步策略以及 race 测试。

## Redis、缓存与持久化

- key prefix 必须覆盖单 key、多 key、destination、pattern、pipeline 和 Lua；区分逻辑 key 与物理 key。
- Redis Cluster 多 key 操作必须定义 hash slot 策略，不能只用单节点或 miniredis 证明正确。
- 异步落盘先定义一致性、ack、重试、幂等、崩溃恢复和脏标记生命周期。
- Mongo client、Database 和 Collection 必须有明确所有权及并发模型。

## 服务发现与 RPC

- 注册值格式变化需要部署顺序、混合数据、坏节点隔离和空地址测试。
- 单个坏值不能无说明地截断其他健康节点。
- 超时、阻塞连接和负载均衡默认值属于公共运行时契约。

## REST 与日志

- 重建或复制 `Route` 时保持权限等全部元数据，并测试 helper 组合顺序。
- 请求体限制、content type、成功与错误响应必须形成完整契约。
- 日志扩展要覆盖 Writer 实现、结构化字段、脱敏、截断、失败行为和调用方迁移。

## goctl

- 生成器源码、模板、fixture、黄金文件和生成后编译共同构成契约。
- 优先在临时目录验证可重复生成和工作树清洁性。
- goctl 独立 module 固定依赖 Hexas 正式版本；仓库内由 `go.work` 使用当前源码，`GOWORK=off` 则验证发布依赖。两种模式都必须可构建。
- 已有服务只更新 RPC 适配器时使用 `--skip-scaffold`，不得生成或覆盖服务入口与 bootstrap YAML。
