# 并发模式

```{contents} 目录
:depth: 3
```

## 概述

Go 语言中有许多经典的并发模式，掌握这些模式可以帮助你编写高效、可维护的并发代码。

## 模式 1：Worker Pool (工作池)

适用于需要限制并发数量的场景。

```go
type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID int
    Value interface{}
    Err   error
}

func WorkerPool(numWorkers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                result := processJob(job)
                results <- result
            }
        }(i)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    return results
}

func processJob(job Job) Result {
    // 处理任务
    return Result{JobID: job.ID, Value: "done"}
}

// 使用示例
func main() {
    jobs := make(chan Job, 100)
    results := WorkerPool(5, jobs)
    
    // 提交任务
    go func() {
        for i := 0; i < 100; i++ {
            jobs <- Job{ID: i, Data: i}
        }
        close(jobs)
    }()
    
    // 收集结果
    for result := range results {
        fmt.Printf("Job %d: %v\n", result.JobID, result.Value)
    }
}
```

## 模式 2：Rate Limiter (速率限制器)

### 使用 time.Ticker

```go
func rateLimiter(requests <-chan Request, rate time.Duration) <-chan Request {
    limited := make(chan Request)
    
    go func() {
        ticker := time.NewTicker(rate)
        defer ticker.Stop()
        defer close(limited)
        
        for req := range requests {
            <-ticker.C // 等待下一个 tick
            limited <- req
        }
    }()
    
    return limited
}
```

### 使用 golang.org/x/time/rate

```go
import "golang.org/x/time/rate"

func main() {
    // 每秒 10 个请求，最多突发 5 个
    limiter := rate.NewLimiter(10, 5)
    
    for i := 0; i < 20; i++ {
        if err := limiter.Wait(context.Background()); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Request %d at %v\n", i, time.Now())
    }
}
```

## 模式 3：Circuit Breaker (熔断器)

防止级联故障的模式。

```go
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreaker struct {
    mu          sync.Mutex
    state       State
    failures    int
    threshold   int
    timeout     time.Duration
    lastFailure time.Time
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:     StateClosed,
        threshold: threshold,
        timeout:   timeout,
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    
    switch cb.state {
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
        } else {
            cb.mu.Unlock()
            return errors.New("circuit breaker is open")
        }
    }
    cb.mu.Unlock()
    
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.threshold {
            cb.state = StateOpen
        }
        return err
    }
    
    cb.failures = 0
    cb.state = StateClosed
    return nil
}
```

## 模式 4：Pub/Sub (发布订阅)

```go
type PubSub struct {
    mu     sync.RWMutex
    subs   map[string][]chan interface{}
    closed bool
}

func NewPubSub() *PubSub {
    return &PubSub{
        subs: make(map[string][]chan interface{}),
    }
}

func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    ch := make(chan interface{}, 1)
    ps.subs[topic] = append(ps.subs[topic], ch)
    return ch
}

func (ps *PubSub) Publish(topic string, msg interface{}) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()
    
    if ps.closed {
        return
    }
    
    for _, ch := range ps.subs[topic] {
        select {
        case ch <- msg:
        default:
            // 如果订阅者处理太慢，跳过
        }
    }
}

func (ps *PubSub) Close() {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    if ps.closed {
        return
    }
    
    ps.closed = true
    for _, subs := range ps.subs {
        for _, ch := range subs {
            close(ch)
        }
    }
}
```

## 模式 5：Context Cancellation (上下文取消)

```go
func longRunningTask(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 执行一小步工作
            if done := doOneStep(); done {
                return nil
            }
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := longRunningTask(ctx); err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            fmt.Println("Task timed out")
        } else if errors.Is(err, context.Canceled) {
            fmt.Println("Task was canceled")
        }
    }
}
```

## 模式 6：Graceful Shutdown (优雅关闭)

```go
func main() {
    server := &http.Server{Addr: ":8080"}
    
    // 启动服务器
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("HTTP server error: %v", err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown error: %v", err)
    }
    
    fmt.Println("Server gracefully stopped")
}
```

## 模式 7：Barrier (屏障)

等待多个 goroutine 到达同一点后再继续。

```go
type Barrier struct {
    total int
    count int
    mu    sync.Mutex
    cond  *sync.Cond
}

func NewBarrier(n int) *Barrier {
    b := &Barrier{total: n}
    b.cond = sync.NewCond(&b.mu)
    return b
}

func (b *Barrier) Wait() {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    b.count++
    if b.count == b.total {
        b.count = 0 // 重置
        b.cond.Broadcast()
    } else {
        b.cond.Wait()
    }
}

// 使用示例
func main() {
    barrier := NewBarrier(3)
    
    for i := 0; i < 3; i++ {
        go func(id int) {
            fmt.Printf("Worker %d preparing\n", id)
            time.Sleep(time.Duration(id) * time.Second)
            
            barrier.Wait() // 等待所有 worker 准备好
            
            fmt.Printf("Worker %d running\n", id)
        }(i)
    }
    
    time.Sleep(5 * time.Second)
}
```

## 模式 8：Tee Channel (T 形分流)

将一个 channel 的数据分发到多个 channel。

```go
func tee(done <-chan struct{}, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
    out1 := make(chan interface{})
    out2 := make(chan interface{})
    
    go func() {
        defer close(out1)
        defer close(out2)
        
        for val := range orDone(done, in) {
            var o1, o2 = out1, out2
            for i := 0; i < 2; i++ {
                select {
                case <-done:
                    return
                case o1 <- val:
                    o1 = nil // 已发送，置为 nil
                case o2 <- val:
                    o2 = nil
                }
            }
        }
    }()
    
    return out1, out2
}
```

## 模式 9：Retry with Backoff (指数退避重试)

```go
type BackoffConfig struct {
    InitialInterval time.Duration
    MaxInterval     time.Duration
    MaxRetries      int
    Multiplier      float64
}

func RetryWithBackoff(ctx context.Context, cfg BackoffConfig, fn func() error) error {
    interval := cfg.InitialInterval
    
    for i := 0; i < cfg.MaxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if i == cfg.MaxRetries-1 {
            return fmt.Errorf("max retries exceeded: %w", err)
        }
        
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(interval):
            // 计算下一次重试间隔
            interval = time.Duration(float64(interval) * cfg.Multiplier)
            if interval > cfg.MaxInterval {
                interval = cfg.MaxInterval
            }
        }
    }
    
    return nil
}

// 使用示例
func main() {
    cfg := BackoffConfig{
        InitialInterval: 100 * time.Millisecond,
        MaxInterval:     5 * time.Second,
        MaxRetries:      5,
        Multiplier:      2.0,
    }
    
    err := RetryWithBackoff(context.Background(), cfg, func() error {
        // 可能失败的操作
        return callExternalService()
    })
    
    if err != nil {
        log.Fatal(err)
    }
}
```

## 模式选择指南

| 场景 | 推荐模式 |
|------|----------|
| 限制并发请求数 | Worker Pool |
| 限制请求频率 | Rate Limiter |
| 防止服务雪崩 | Circuit Breaker |
| 事件驱动解耦 | Pub/Sub |
| 任务超时控制 | Context Cancellation |
| 服务停止 | Graceful Shutdown |
| 等待多个任务到达 | Barrier |
| 重试失败操作 | Retry with Backoff |

## 参考资源

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Advanced Go Concurrency Patterns](https://go.dev/blog/io2013-talk-concurrency)
- [Context Package](https://go.dev/blog/context)
