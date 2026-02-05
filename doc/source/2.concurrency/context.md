# Context 详解

```{contents} 目录
:depth: 3
```

## 什么是 Context

Context 是 Go 中用于在 API 边界和进程之间传递截止时间、取消信号和请求范围值的接口。它被设计为在调用链中传递，并可在任何时候被取消。

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

## 核心概念

- **取消机制**：发出应停止工作的信号
- **截止时间**：超过时间后自动取消
- **请求范围的值**：在调用链中携带请求特定数据

## 创建 Context

```go
// 背景 Context
ctx := context.Background()

// 带取消
ctx, cancel := context.WithCancel(ctx)
defer cancel()

// 带超时
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// 带值
ctx = context.WithValue(ctx, "userID", "123")
```

## 最佳实践

1. **始终检查 Context 取消** — 在长操作中定期检查 `ctx.Done()`
2. **使用 defer cancel()** — 防止 Context 泄漏
3. **不要将 Context 存储在结构体中** — 作为参数传递
4. **Context 作为第一个参数** — 遵循 Go 社区约定
5. **谨慎使用 Context 值** — 只存储请求范围数据，不存业务数据

## 常见错误

- 忽略 Context 取消信号
- 在循环中不检查 Context
- Context 泄漏（未调用 cancel）
- 在 Context 中存储可变数据或敏感信息

```{seealso}
延伸阅读：`Context in Go <https://www.fanyamin.com/journal/2025-08-28-context-in-go.html>`_ — 详细的最佳实践、最差实践与 contextcheck 工具说明。
```

## 参考资源

- [Go context package](https://pkg.go.dev/context)
- [Go Blog: Context](https://go.dev/blog/context)
