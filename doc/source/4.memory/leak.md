# 内存泄漏排查

```{contents} 目录
:depth: 3
```

## Go 中的内存泄漏

Go 有垃圾回收，但仍可能发生"逻辑泄漏"——对象仍被引用但不再使用。

## 常见泄漏场景

### 场景 1：Goroutine 泄漏

```go
// ❌ Goroutine 永远阻塞
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // 永远阻塞
        fmt.Println(val)
    }()
    // ch 永远没有发送者
}

// ✅ 使用 context 控制
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return
        }
    }()
}
```

### 场景 2：切片引用

```go
// ❌ 小切片引用大数组
var global []byte

func leak() {
    big := make([]byte, 1<<20)  // 1MB
    // ... 填充数据
    global = big[:10]  // 只需要 10 字节，但 1MB 无法释放
}

// ✅ 复制需要的数据
func noLeak() {
    big := make([]byte, 1<<20)
    // ... 填充数据
    global = make([]byte, 10)
    copy(global, big[:10])
}
```

### 场景 3：time.Ticker 未停止

```go
// ❌ Ticker 泄漏
func leak() {
    ticker := time.NewTicker(time.Second)
    // 忘记调用 ticker.Stop()
}

// ✅ 确保停止
func noLeak() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    // 使用 ticker
}
```

### 场景 4：资源未关闭

```go
// ❌ 文件未关闭
func leak() {
    f, _ := os.Open("file.txt")
    // 忘记 f.Close()
}

// ✅ 使用 defer 关闭
func noLeak() {
    f, err := os.Open("file.txt")
    if err != nil {
        return
    }
    defer f.Close()
    
    // 使用文件
}
```

### 场景 5：全局 Map 无限增长

```go
// ❌ Map 无限增长
var cache = make(map[string]Data)

func process(key string, data Data) {
    cache[key] = data  // 只增不减
}

// ✅ 使用 LRU 或定期清理
import "github.com/hashicorp/golang-lru"

var cache, _ = lru.New(1000)  // 最多 1000 项

func process(key string, data Data) {
    cache.Add(key, data)  // 自动淘汰旧数据
}
```

### 场景 6：闭包捕获

```go
// ❌ 闭包捕获大对象
func leak() []func() {
    var funcs []func()
    for i := 0; i < 100; i++ {
        bigData := make([]byte, 1<<20)  // 1MB
        funcs = append(funcs, func() {
            _ = bigData  // 捕获 bigData
        })
    }
    return funcs  // 100 个 1MB 对象
}

// ✅ 只捕获需要的数据
func noLeak() []func() {
    var funcs []func()
    for i := 0; i < 100; i++ {
        bigData := make([]byte, 1<<20)
        summary := computeSummary(bigData)
        funcs = append(funcs, func() {
            _ = summary  // 只捕获摘要
        })
    }
    return funcs
}
```

## 检测内存泄漏

### 使用 pprof

```bash
# 获取堆分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 比较两个时间点的堆
go tool pprof -base heap1.prof heap2.prof
```

### 使用 runtime.MemStats

```go
func monitorMemory() {
    var m runtime.MemStats
    for {
        runtime.ReadMemStats(&m)
        fmt.Printf("HeapAlloc = %v MB\n", m.HeapAlloc/1024/1024)
        fmt.Printf("NumGoroutine = %v\n", runtime.NumGoroutine())
        time.Sleep(10 * time.Second)
    }
}
```

### 使用 goleak（测试中）

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

## 内存泄漏排查流程

```{mermaid}
flowchart TD
    A[发现内存增长] --> B{确定泄漏类型}
    B -->|Goroutine 数量增长| C[pprof goroutine]
    B -->|堆内存增长| D[pprof heap]
    
    C --> E[找出阻塞的 goroutine]
    D --> F[找出分配来源]
    
    E --> G[检查 channel/锁/sleep]
    F --> H[检查是否有引用未释放]
    
    G --> I[修复代码]
    H --> I
    
    I --> J[验证修复]
```

## 最佳实践

1. **使用 context 控制 goroutine 生命周期**
2. **使用 defer 释放资源**
3. **避免全局容器无限增长**
4. **定期监控 goroutine 数量和内存**
5. **在测试中使用 goleak**
6. **使用有界的缓存（LRU）**

## 参考资源

- [Finding Memory Leaks in Go Programs](https://go.dev/blog/pprof)
- [goleak - Goroutine Leak Detector](https://github.com/uber-go/goleak)
- [警惕！你的 Go 程序正在偷偷"泄漏"](https://www.fanyamin.com/journal/2025-12-13-go-goroutine-leak-jing-ti-ni-de-cheng-xu-zheng-zai-tou-tou-x.html) — Goroutine Leak 实战案例与 pprof 排查
