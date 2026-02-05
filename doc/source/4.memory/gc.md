# 垃圾回收 (GC)

```{contents} 目录
:depth: 3
```

## Go GC 概述

Go 使用**并发、三色标记、写屏障**的垃圾回收器。

| 特性 | 描述 |
|------|------|
| 算法 | 三色标记清除 |
| 类型 | 并发、非分代 |
| STW | 极短（通常 < 1ms） |
| 目标 | 低延迟 |

## GC 触发条件

1. **堆内存达到阈值**：当堆大小达到上次 GC 后大小的一定比例（由 GOGC 控制）
2. **定时触发**：2 分钟强制执行一次 GC
3. **手动触发**：调用 `runtime.GC()`

## GOGC 环境变量

```bash
# 默认值 100：当堆增长 100% 时触发 GC
GOGC=100

# 更激进的 GC（更低内存，更高 CPU）
GOGC=50

# 更宽松的 GC（更高内存，更低 CPU）
GOGC=200

# 禁用 GC（危险！）
GOGC=off
```

### 运行时设置

```go
import "runtime/debug"

func main() {
    // 设置 GOGC
    debug.SetGCPercent(50)
    
    // 设置内存限制（Go 1.19+）
    debug.SetMemoryLimit(1 << 30) // 1GB
}
```

## 三色标记算法

```{mermaid}
graph LR
    subgraph "三色标记"
        W[白色<br/>未访问] --> G[灰色<br/>待处理]
        G --> B[黑色<br/>已访问]
    end
```

1. **白色**：未被访问的对象，GC 结束后会被回收
2. **灰色**：已被访问但其引用的对象未全部访问
3. **黑色**：已被访问且其引用的对象已全部访问

### 标记过程

1. STW，启用写屏障
2. 从根对象（栈、全局变量）开始标记
3. 并发标记（与用户程序并行）
4. STW，完成标记
5. 并发清除

## GC 监控

### 查看 GC 日志

```bash
GODEBUG=gctrace=1 ./myprogram
```

输出示例：
```
gc 1 @0.012s 2%: 0.026+0.44+0.003 ms clock, 0.10+0.32/0.27/0+0.012 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
```

| 字段 | 含义 |
|------|------|
| `gc 1` | 第 1 次 GC |
| `@0.012s` | 程序启动后 0.012 秒 |
| `2%` | GC 占用 CPU 时间比例 |
| `0.026+0.44+0.003 ms` | STW 时间 + 并发标记时间 + STW 时间 |
| `4->4->0 MB` | GC 前堆大小 -> GC 时堆大小 -> GC 后存活大小 |
| `5 MB goal` | 下次 GC 触发的目标堆大小 |

### 程序内监控

```go
import "runtime"

func printGCStats() {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)
    
    fmt.Printf("Alloc = %v MB\n", stats.Alloc/1024/1024)
    fmt.Printf("TotalAlloc = %v MB\n", stats.TotalAlloc/1024/1024)
    fmt.Printf("Sys = %v MB\n", stats.Sys/1024/1024)
    fmt.Printf("NumGC = %v\n", stats.NumGC)
    fmt.Printf("PauseTotalNs = %v ms\n", stats.PauseTotalNs/1e6)
}
```

## GC 调优

### 1. 减少分配

```go
// ❌ 频繁分配
func process(items []Item) {
    for _, item := range items {
        result := &Result{}  // 每次循环都分配
        result.Process(item)
    }
}

// ✅ 复用对象
func process(items []Item) {
    result := &Result{}
    for _, item := range items {
        result.Reset()
        result.Process(item)
    }
}
```

### 2. 使用 sync.Pool

```go
var resultPool = sync.Pool{
    New: func() interface{} {
        return &Result{}
    },
}

func process(item Item) {
    result := resultPool.Get().(*Result)
    defer resultPool.Put(result)
    
    result.Reset()
    result.Process(item)
}
```

### 3. 预分配内存

```go
// ❌ 动态增长
var data []byte
for i := 0; i < 1000; i++ {
    data = append(data, byte(i))
}

// ✅ 预分配
data := make([]byte, 0, 1000)
for i := 0; i < 1000; i++ {
    data = append(data, byte(i))
}
```

### 4. 使用值类型

```go
// ❌ 指针增加 GC 扫描压力
type Item struct {
    Name  *string
    Value *int
}

// ✅ 值类型
type Item struct {
    Name  string
    Value int
}
```

### 5. 设置内存限制（Go 1.19+）

```go
import "runtime/debug"

func init() {
    // 设置软内存限制
    debug.SetMemoryLimit(512 << 20) // 512MB
}
```

## Go 1.19+ GOMEMLIMIT

Go 1.19 引入了 `GOMEMLIMIT`，提供更精确的内存控制。

```bash
# 设置内存限制
GOMEMLIMIT=512MiB ./myprogram

# 单位支持: B, KiB, MiB, GiB, TiB
```

### GOMEMLIMIT vs GOGC

| 特性 | GOGC | GOMEMLIMIT |
|------|------|------------|
| 控制方式 | 相对增长率 | 绝对内存限制 |
| 适用场景 | 通用 | 容器环境 |
| OOM 风险 | 可能 OOM | 更可控 |

## GC Pacer

Go 1.18+ 改进了 GC Pacer 算法，提供更平滑的 GC 行为。

```go
// 查看 Pacer 状态
GODEBUG=gcpacertrace=1 ./myprogram
```

## 常见问题

### 问题 1：GC 频繁

**症状**：GC 日志显示频繁 GC，CPU 使用高

**解决**：
- 增大 GOGC
- 使用 sync.Pool
- 减少分配

### 问题 2：大堆内存

**症状**：内存持续增长

**解决**：
- 检查内存泄漏
- 使用 GOMEMLIMIT
- 降低 GOGC

### 问题 3：STW 时间长

**症状**：GC pause 时间长（> 10ms）

**解决**：
- 减少指针数量
- 减少全局变量
- 使用更小的堆

## 参考资源

- [Go GC Guide](https://tip.golang.org/doc/gc-guide)
- [Getting to Go: The Journey of Go's Garbage Collector](https://go.dev/blog/ismmkeynote)
- [Go 1.19 Memory Limit](https://go.dev/doc/go1.19#runtime)
