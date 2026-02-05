# 并发陷阱

```{contents} 目录
:depth: 3
```

## 概述

Go 的并发模型虽然简洁，但仍然有许多容易踩坑的地方。本节总结了最常见的并发陷阱及其解决方案。

## 陷阱 1：数据竞争 (Data Race)

### 问题描述

当多个 goroutine 同时访问同一个变量，且至少有一个是写操作时，就会发生数据竞争。

```go
// ❌ 数据竞争示例
var counter int

func main() {
    for i := 0; i < 1000; i++ {
        go func() {
            counter++ // 数据竞争！
        }()
    }
    time.Sleep(time.Second)
    fmt.Println(counter) // 结果不确定
}
```

### 检测方法

```bash
go run -race main.go
go test -race ./...
```

### 解决方案

```go
// ✅ 方案 1：使用 Mutex
var (
    counter int
    mu      sync.Mutex
)

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}

// ✅ 方案 2：使用 atomic
var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}

// ✅ 方案 3：使用 channel
func main() {
    counter := make(chan int, 1)
    counter <- 0
    
    for i := 0; i < 1000; i++ {
        go func() {
            v := <-counter
            v++
            counter <- v
        }()
    }
}
```

## 陷阱 2：闭包捕获循环变量

### 问题描述

```go
// ❌ 所有 goroutine 都打印相同的值
func main() {
    for i := 0; i < 5; i++ {
        go func() {
            fmt.Println(i) // 可能全部打印 5
        }()
    }
    time.Sleep(time.Second)
}
```

### 解决方案

```go
// ✅ 方案 1：传递参数
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}

// ✅ 方案 2：创建局部变量
for i := 0; i < 5; i++ {
    i := i // 创建新变量
    go func() {
        fmt.Println(i)
    }()
}

// ✅ Go 1.22+ 自动修复（循环变量每次迭代都是新的）
```

## 陷阱 3：Goroutine 泄漏

### 问题描述

```go
// ❌ goroutine 永远阻塞
func leakyFunction() {
    ch := make(chan int)
    go func() {
        val := <-ch // 永远阻塞，因为没有发送者
        fmt.Println(val)
    }()
    // 函数返回，但 goroutine 仍然存在
}
```

### 解决方案

```go
// ✅ 使用 context 控制生命周期
func nonLeakyFunction(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return // 优雅退出
        }
    }()
}

// ✅ 使用 buffered channel
func nonLeakyFunction2() {
    ch := make(chan int, 1)
    go func() {
        ch <- 42 // 不会阻塞
    }()
    // 即使不读取，goroutine 也能完成
}
```

## 陷阱 4：死锁

### 常见死锁场景

#### 场景 1：channel 自己等自己

```go
// ❌ 死锁
func main() {
    ch := make(chan int)
    ch <- 1  // 阻塞，等待接收者
    <-ch     // 永远执行不到
}

// ✅ 修复
func main() {
    ch := make(chan int)
    go func() {
        ch <- 1
    }()
    <-ch
}
```

#### 场景 2：循环等待

```go
// ❌ 两个 goroutine 互相等待
var mu1, mu2 sync.Mutex

// goroutine 1
go func() {
    mu1.Lock()
    time.Sleep(time.Millisecond)
    mu2.Lock() // 等待 goroutine 2 释放
    mu2.Unlock()
    mu1.Unlock()
}()

// goroutine 2
go func() {
    mu2.Lock()
    time.Sleep(time.Millisecond)
    mu1.Lock() // 等待 goroutine 1 释放
    mu1.Unlock()
    mu2.Unlock()
}()

// ✅ 修复：按固定顺序获取锁
```

#### 场景 3：WaitGroup 使用错误

```go
// ❌ 死锁：Wait 在 Add 之前
var wg sync.WaitGroup

go func() {
    wg.Add(1)
    defer wg.Done()
    // ...
}()

wg.Wait() // 可能在 Add 之前执行

// ✅ 修复：Add 在启动 goroutine 之前
var wg sync.WaitGroup
wg.Add(1)

go func() {
    defer wg.Done()
    // ...
}()

wg.Wait()
```

## 陷阱 5：不正确的锁粒度

### 锁粒度太大

```go
// ❌ 锁住整个操作，包括 I/O
func (s *Server) HandleRequest(req Request) Response {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // 网络请求被锁保护，导致所有请求串行
    result := s.callExternalService(req)
    return result
}

// ✅ 只锁住需要保护的部分
func (s *Server) HandleRequest(req Request) Response {
    // 网络请求不需要锁
    result := s.callExternalService(req)
    
    s.mu.Lock()
    s.cache[req.ID] = result
    s.mu.Unlock()
    
    return result
}
```

### 锁粒度太小

```go
// ❌ 多次加锁解锁，开销大
func (c *Counter) AddAndGet() int {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
    
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count // 可能已经被其他 goroutine 修改
}

// ✅ 原子操作
func (c *Counter) AddAndGet() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
    return c.count
}
```

## 陷阱 6：错误的 channel 关闭

```go
// ❌ 多次关闭 channel
func bad() {
    ch := make(chan int)
    close(ch)
    close(ch) // panic!
}

// ❌ 向已关闭的 channel 发送
func bad2() {
    ch := make(chan int)
    close(ch)
    ch <- 1 // panic!
}

// ✅ 使用 sync.Once 确保只关闭一次
type SafeChannel struct {
    ch   chan int
    once sync.Once
}

func (sc *SafeChannel) Close() {
    sc.once.Do(func() {
        close(sc.ch)
    })
}
```

## 陷阱 7：select 中的优先级问题

```go
// ❌ 可能永远不会处理 done
func process(done chan struct{}, data chan int) {
    for {
        select {
        case <-done:
            return
        case v := <-data:
            process(v)
        }
    }
}

// 如果 data 一直有数据，done 可能永远不被选中

// ✅ 每次循环都检查 done
func process(done chan struct{}, data chan int) {
    for {
        select {
        case <-done:
            return
        default:
        }
        
        select {
        case <-done:
            return
        case v := <-data:
            process(v)
        }
    }
}
```

## 陷阱 8：time.After 在循环中使用

```go
// ❌ 每次循环都创建新的 timer，内存泄漏
for {
    select {
    case <-ch:
        // ...
    case <-time.After(time.Second): // 每次都分配新的 timer
        // ...
    }
}

// ✅ 复用 timer
timer := time.NewTimer(time.Second)
defer timer.Stop()

for {
    select {
    case <-ch:
        if !timer.Stop() {
            <-timer.C
        }
        timer.Reset(time.Second)
    case <-timer.C:
        timer.Reset(time.Second)
    }
}
```

## 陷阱 9：map 并发读写

```go
// ❌ panic: concurrent map read and map write
var m = make(map[string]int)

go func() {
    for {
        m["key"] = 1
    }
}()

go func() {
    for {
        _ = m["key"]
    }
}()

// ✅ 使用 sync.RWMutex
var (
    m  = make(map[string]int)
    mu sync.RWMutex
)

func set(k string, v int) {
    mu.Lock()
    m[k] = v
    mu.Unlock()
}

func get(k string) int {
    mu.RLock()
    defer mu.RUnlock()
    return m[k]
}

// ✅ 或使用 sync.Map
var m sync.Map

m.Store("key", 1)
v, _ := m.Load("key")
```

## 陷阱 10：忽略 context 取消

```go
// ❌ 忽略 context，无法取消
func longOperation(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        heavyComputation()
    }
    return nil
}

// ✅ 定期检查 context
func longOperation(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        if ctx.Err() != nil {
            return ctx.Err()
        }
        heavyComputation()
    }
    return nil
}

// ✅ 在长操作中使用 select
func longOperation(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            heavyComputation()
        }
    }
    return nil
}
```

## 总结

| 陷阱 | 检测方法 | 解决方案 |
|------|----------|----------|
| 数据竞争 | `-race` 标志 | mutex/atomic/channel |
| 闭包捕获 | 代码审查 | 传参/局部变量 |
| Goroutine 泄漏 | 监控 NumGoroutine | context/buffered channel |
| 死锁 | `-race`/死锁检测器 | 固定锁顺序/超时 |
| 锁粒度问题 | 性能分析 | 合理划分临界区 |
| channel 关闭 | 代码审查 | sync.Once |
| select 优先级 | 代码审查 | 双重 select |
| time.After 泄漏 | 内存分析 | 复用 timer |
| map 并发 | `-race` 标志 | RWMutex/sync.Map |
| 忽略 context | 代码审查 | 定期检查 ctx.Done() |

## 参考资源

- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Uber Go Style Guide - Concurrency](https://github.com/uber-go/guide/blob/master/style.md#dont-fire-and-forget-goroutines)
