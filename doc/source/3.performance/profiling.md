# 性能分析 (Profiling)

```{contents} 目录
:depth: 3
```

## pprof 概述

Go 内置了强大的性能分析工具 `pprof`，支持以下类型的分析：

| 类型 | 描述 |
|------|------|
| CPU | CPU 使用情况分析 |
| Heap | 堆内存分配分析 |
| Goroutine | Goroutine 栈分析 |
| Block | 阻塞操作分析 |
| Mutex | 互斥锁竞争分析 |
| Trace | 执行追踪 |

## 开启 pprof

### 方式 1：HTTP 服务（推荐）

```go
import (
    "net/http"
    _ "net/http/pprof" // 匿名导入，自动注册 handler
)

func main() {
    // 在单独的 goroutine 中启动 pprof 服务
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // 主程序逻辑
    runMainApplication()
}
```

访问端点：
- `http://localhost:6060/debug/pprof/` - 总览页面
- `http://localhost:6060/debug/pprof/heap` - 堆内存
- `http://localhost:6060/debug/pprof/goroutine` - Goroutine
- `http://localhost:6060/debug/pprof/profile` - CPU (需要参数 seconds)

### 方式 2：代码中手动收集

```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    // CPU 分析
    cpuFile, _ := os.Create("cpu.prof")
    defer cpuFile.Close()
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
    
    // 你的程序逻辑
    doWork()
    
    // 内存分析
    memFile, _ := os.Create("mem.prof")
    defer memFile.Close()
    pprof.WriteHeapProfile(memFile)
}
```

## CPU 性能分析

### 收集 CPU Profile

```bash
# 方式 1：通过 HTTP
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 方式 2：从文件
go tool pprof cpu.prof
```

### pprof 交互式命令

```bash
(pprof) top          # 显示消耗最多 CPU 的函数
(pprof) top10        # 显示前 10 个
(pprof) list funcName # 显示函数的源码及耗时
(pprof) web          # 在浏览器中打开调用图
(pprof) svg          # 生成 SVG 图
(pprof) pdf          # 生成 PDF 报告
```

### 示例输出

```
(pprof) top10
Showing nodes accounting for 2.10s, 84.00% of 2.50s total
Showing top 10 nodes out of 50
      flat  flat%   sum%        cum   cum%
     0.70s 28.00% 28.00%      0.70s 28.00%  runtime.memmove
     0.30s 12.00% 40.00%      0.50s 20.00%  main.processData
     0.25s 10.00% 50.00%      0.25s 10.00%  runtime.mallocgc
     ...
```

### 理解指标

- **flat**: 函数自身消耗的时间
- **cum**: 函数及其调用的所有函数消耗的时间
- **flat%**: flat 时间占总时间的百分比
- **cum%**: cum 时间占总时间的百分比

## 内存分析

### 收集 Heap Profile

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

### 分析内存分配

```bash
# 查看当前堆内存使用
(pprof) top

# 查看累计分配（包括已释放的）
(pprof) top -cum

# 查看分配次数而非字节数
go tool pprof -alloc_objects http://localhost:6060/debug/pprof/heap
```

### 内存分析类型

| 参数 | 描述 |
|------|------|
| `-inuse_space` | 当前使用的内存 (默认) |
| `-inuse_objects` | 当前使用的对象数 |
| `-alloc_space` | 累计分配的内存 |
| `-alloc_objects` | 累计分配的对象数 |

## Goroutine 分析

### 查看 Goroutine 状态

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### 检测 Goroutine 泄漏

```go
func checkGoroutineLeak() {
    before := runtime.NumGoroutine()
    
    // 执行可能泄漏的操作
    doSomething()
    
    // 等待 goroutine 完成
    time.Sleep(time.Second)
    
    after := runtime.NumGoroutine()
    if after > before {
        fmt.Printf("可能存在 Goroutine 泄漏: %d -> %d\n", before, after)
    }
}
```

## 阻塞分析

需要先开启阻塞分析：

```go
import "runtime"

func init() {
    runtime.SetBlockProfileRate(1) // 开启阻塞分析
}
```

```bash
go tool pprof http://localhost:6060/debug/pprof/block
```

## Mutex 竞争分析

需要先开启 Mutex 分析：

```go
import "runtime"

func init() {
    runtime.SetMutexProfileFraction(1) // 开启 mutex 分析
}
```

```bash
go tool pprof http://localhost:6060/debug/pprof/mutex
```

## Trace 分析

Trace 提供更细粒度的执行追踪。

### 收集 Trace

```go
import (
    "os"
    "runtime/trace"
)

func main() {
    f, _ := os.Create("trace.out")
    defer f.Close()
    
    trace.Start(f)
    defer trace.Stop()
    
    // 你的程序逻辑
}
```

### 查看 Trace

```bash
go tool trace trace.out
```

Trace 可以显示：
- Goroutine 调度
- 网络 I/O
- 系统调用
- GC 活动
- 用户自定义事件

## 可视化工具

### 使用 go tool pprof Web UI

```bash
# 需要安装 graphviz
brew install graphviz

# 启动 Web UI
go tool pprof -http=:8080 cpu.prof
```

### Flame Graph (火焰图)

```bash
# 安装 go-torch（已集成到 go tool pprof）
go tool pprof -http=:8080 -flame cpu.prof
```

## 生产环境最佳实践

### 1. 安全暴露 pprof

```go
import (
    "net/http"
    "net/http/pprof"
)

func main() {
    // 只在内部网络监听
    mux := http.NewServeMux()
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
    
    // 只监听内网地址
    go http.ListenAndServe("127.0.0.1:6060", mux)
}
```

### 2. 持续性能监控

```go
import (
    "os"
    "runtime/pprof"
    "time"
)

func periodicMemProfile() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        f, err := os.Create(fmt.Sprintf("mem-%d.prof", time.Now().Unix()))
        if err != nil {
            continue
        }
        pprof.WriteHeapProfile(f)
        f.Close()
    }
}
```

### 3. 自定义 pprof Label

```go
import "runtime/pprof"

func handleRequest(ctx context.Context, userID string) {
    labels := pprof.Labels("user_id", userID)
    pprof.Do(ctx, labels, func(ctx context.Context) {
        // 处理请求
        processRequest(ctx)
    })
}
```

## 常用分析流程

```{mermaid}
flowchart TD
    A[发现性能问题] --> B{确定问题类型}
    B -->|CPU 高| C[CPU Profile]
    B -->|内存高| D[Heap Profile]
    B -->|响应慢| E[Trace/Block Profile]
    B -->|死锁| F[Goroutine/Mutex Profile]
    
    C --> G[找出热点函数]
    D --> H[找出内存分配点]
    E --> I[找出阻塞点]
    F --> J[找出锁竞争]
    
    G --> K[优化代码]
    H --> K
    I --> K
    J --> K
    
    K --> L[验证效果]
    L -->|未解决| B
    L -->|已解决| M[完成]
```

## 参考资源

- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
- [pprof Package](https://pkg.go.dev/runtime/pprof)
- [Diagnostics](https://go.dev/doc/diagnostics)
- [Go 程序崩溃分析实战](https://www.fanyamin.com/journal/2026-01-23-golang_crash_analysis.html) — Coredump、Delve 分析、预防措施
