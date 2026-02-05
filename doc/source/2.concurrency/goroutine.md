# Goroutine 深入理解

```{contents} 目录
:depth: 3
```

## 什么是 Goroutine

Goroutine 是 Go 语言中的轻量级线程，由 Go 运行时管理。与操作系统线程相比，Goroutine 具有以下特点：

| 特性 | OS 线程 | Goroutine |
|------|---------|-----------|
| 内存占用 | ~1MB | ~2KB |
| 创建成本 | 高 | 低 |
| 切换成本 | 高（内核态） | 低（用户态） |
| 调度 | OS 调度器 | Go 运行时调度器 |

## Goroutine 调度模型 (GMP)

Go 使用 **GMP 模型** 进行 Goroutine 调度：

- **G (Goroutine)**：Goroutine，包含栈、指令指针等信息
- **M (Machine)**：操作系统线程，执行 Goroutine 的载体
- **P (Processor)**：逻辑处理器，维护本地运行队列

```{mermaid}
graph TD
    subgraph "Go Runtime"
        GQ[Global Queue]
        subgraph "P1"
            LQ1[Local Queue]
            G1[G1]
            G2[G2]
        end
        subgraph "P2"
            LQ2[Local Queue]
            G3[G3]
            G4[G4]
        end
        M1[M1 - OS Thread]
        M2[M2 - OS Thread]
    end
    
    P1 --> M1
    P2 --> M2
    GQ --> P1
    GQ --> P2
```

### 设置 P 的数量

```go
import "runtime"

func main() {
    // 获取当前 GOMAXPROCS
    n := runtime.GOMAXPROCS(0)
    fmt.Println("GOMAXPROCS:", n)
    
    // 设置 GOMAXPROCS（通常不需要手动设置）
    runtime.GOMAXPROCS(4)
}
```

## Goroutine 的创建与生命周期

### 基本用法

```go
func main() {
    go func() {
        fmt.Println("Hello from goroutine")
    }()
    
    // 等待 goroutine 执行完毕
    time.Sleep(time.Second)
}
```

### ⚠️ 常见陷阱：主 Goroutine 退出

```go
// ❌ 错误示例：主 goroutine 退出，子 goroutine 被强制终止
func main() {
    go func() {
        time.Sleep(time.Second)
        fmt.Println("这行代码可能永远不会执行")
    }()
    // main 函数立即返回，程序退出
}

// ✅ 正确做法：使用 WaitGroup 等待
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    
    go func() {
        defer wg.Done()
        time.Sleep(time.Second)
        fmt.Println("正常执行")
    }()
    
    wg.Wait() // 等待所有 goroutine 完成
}
```

## Goroutine 泄漏

Goroutine 泄漏是 Go 程序中最常见的问题之一。

### 泄漏场景 1：Channel 阻塞

```go
// ❌ 泄漏示例：channel 永远不会被读取
func leak1() {
    ch := make(chan int)
    go func() {
        ch <- 42 // 永远阻塞，goroutine 泄漏
    }()
    // 没有接收者
}

// ✅ 修复：使用 buffered channel 或确保有接收者
func fixed1() {
    ch := make(chan int, 1) // buffered channel
    go func() {
        ch <- 42
    }()
    // 或者确保读取
}
```

### 泄漏场景 2：无限循环没有退出条件

```go
// ❌ 泄漏示例：无法停止的 goroutine
func leak2() {
    go func() {
        for {
            doSomething()
            time.Sleep(time.Second)
        }
    }()
}

// ✅ 修复：使用 context 或 done channel
func fixed2(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return // 优雅退出
            default:
                doSomething()
                time.Sleep(time.Second)
            }
        }
    }()
}
```

### 泄漏场景 3：select 中没有 default 或超时

```go
// ❌ 泄漏示例：可能永远阻塞
func leak3(ch1, ch2 chan int) {
    go func() {
        select {
        case v := <-ch1:
            process(v)
        case v := <-ch2:
            process(v)
        // 如果两个 channel 都没有数据，永远阻塞
        }
    }()
}

// ✅ 修复：添加超时或 context
func fixed3(ctx context.Context, ch1, ch2 chan int) {
    go func() {
        select {
        case v := <-ch1:
            process(v)
        case v := <-ch2:
            process(v)
        case <-ctx.Done():
            return
        case <-time.After(5 * time.Second):
            return // 超时退出
        }
    }()
}
```

## 检测 Goroutine 泄漏

### 使用 runtime 包

```go
func monitorGoroutines() {
    ticker := time.NewTicker(time.Second)
    for range ticker.C {
        fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    }
}
```

### 使用 goleak 库（推荐）

```go
import "go.uber.org/goleak"

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}

func TestNoLeak(t *testing.T) {
    defer goleak.VerifyNone(t)
    // 测试代码
}
```

## Goroutine 的栈增长

Go 的 Goroutine 使用动态栈，初始大小为 2KB，可以按需增长。

```go
// 递归深度测试
func recursiveFunc(depth int) {
    if depth == 0 {
        return
    }
    var arr [1024]byte // 1KB 局部变量
    _ = arr
    recursiveFunc(depth - 1)
}

func main() {
    // 这会触发栈增长
    recursiveFunc(1000)
}
```

## 最佳实践

### 1. 总是考虑 Goroutine 如何退出

```go
// 使用 context 控制生命周期
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("worker exiting")
            return
        default:
            // 工作逻辑
        }
    }
}
```

### 2. 限制并发数量

```go
// 使用 semaphore 限制并发
func limitedConcurrency(tasks []Task, limit int) {
    sem := make(chan struct{}, limit)
    var wg sync.WaitGroup
    
    for _, task := range tasks {
        wg.Add(1)
        sem <- struct{}{} // 获取信号量
        
        go func(t Task) {
            defer wg.Done()
            defer func() { <-sem }() // 释放信号量
            
            t.Execute()
        }(task)
    }
    
    wg.Wait()
}
```

### 3. 使用 errgroup 管理一组 Goroutine

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) error {
    g, ctx := errgroup.WithContext(ctx)
    
    for _, url := range urls {
        url := url // 避免闭包陷阱
        g.Go(func() error {
            return fetch(ctx, url)
        })
    }
    
    return g.Wait() // 等待所有完成，返回第一个错误
}
```

## 参考资源

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Go Memory Model](https://go.dev/ref/mem)
- [Uber Go Style Guide - Goroutines](https://github.com/uber-go/guide/blob/master/style.md#goroutines)
- [警惕！你的 Go 程序正在偷偷"泄漏"](https://www.fanyamin.com/journal/2025-12-13-go-goroutine-leak-jing-ti-ni-de-cheng-xu-zheng-zai-tou-tou-x.html) — Goroutine Leak 实战排查与修复
