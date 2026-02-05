# sync 包详解

```{contents} 目录
:depth: 3
```

## sync 包概述

`sync` 包提供了基本的同步原语，用于低级别的内存访问同步。

```{warning}
sync 包中的类型在使用后不能被复制！
```

## sync.Mutex (互斥锁)

### 基本用法

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

### ⚠️ 常见陷阱

#### 陷阱 1：忘记解锁

```go
// ❌ 如果中间 return，锁不会释放
func (c *SafeCounter) Inc() {
    c.mu.Lock()
    if someCondition {
        return // 锁没有释放！
    }
    c.count++
    c.mu.Unlock()
}

// ✅ 使用 defer
func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if someCondition {
        return // defer 会确保解锁
    }
    c.count++
}
```

#### 陷阱 2：锁的复制

```go
type Counter struct {
    mu    sync.Mutex
    count int
}

// ❌ 值传递会复制锁
func (c Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++ // 修改的是副本
}

// ✅ 使用指针接收者
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

#### 陷阱 3：死锁

```go
// ❌ 同一个 goroutine 重复加锁
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.addOne() // 如果 addOne 也加锁，死锁！
}

func (c *Counter) addOne() {
    c.mu.Lock() // 死锁！
    defer c.mu.Unlock()
    c.count++
}

// ✅ 内部方法不加锁，或使用读写锁
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.addOneInternal()
}

func (c *Counter) addOneInternal() {
    c.count++ // 假设调用者已持有锁
}
```

## sync.RWMutex (读写锁)

适用于读多写少的场景。

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

// 读操作使用 RLock
func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

// 写操作使用 Lock
func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}
```

### 读写锁规则

| 操作 | 条件 |
|------|------|
| RLock | 可以多个 goroutine 同时持有读锁 |
| Lock | 必须等待所有读锁和写锁释放 |
| RUnlock | 释放一个读锁 |
| Unlock | 释放写锁 |

## sync.WaitGroup

等待一组 Goroutine 完成。

```go
func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }
    
    wg.Wait() // 等待所有 goroutine 完成
    fmt.Println("All workers done")
}
```

### ⚠️ WaitGroup 陷阱

#### 陷阱 1：Add 在 goroutine 外调用

```go
// ❌ 可能出现 race condition
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    go func(id int) {
        wg.Add(1) // 太晚了！Wait 可能已经返回
        defer wg.Done()
        // ...
    }(i)
}
wg.Wait()

// ✅ Add 在启动 goroutine 前调用
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // ...
    }(i)
}
wg.Wait()
```

#### 陷阱 2：Done 调用次数不匹配

```go
// ❌ Done 调用太多次会 panic
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    defer wg.Done() // panic: negative WaitGroup counter
}()
```

## sync.Once

确保某个操作只执行一次，常用于单例模式。

```go
var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        instance = &Database{}
        instance.Connect()
    })
    return instance
}
```

### ⚠️ Once 陷阱

```go
// ❌ 如果 Do 中的函数 panic，Once 仍然被标记为已执行
var once sync.Once
once.Do(func() {
    panic("oops") // panic 后，once 被标记为已执行
})
once.Do(func() {
    fmt.Println("这行不会执行")
})
```

## sync.Cond (条件变量)

用于等待或通知条件状态变化。

```go
type Queue struct {
    items []int
    cond  *sync.Cond
}

func NewQueue() *Queue {
    q := &Queue{}
    q.cond = sync.NewCond(&sync.Mutex{})
    return q
}

func (q *Queue) Push(item int) {
    q.cond.L.Lock()
    defer q.cond.L.Unlock()
    
    q.items = append(q.items, item)
    q.cond.Signal() // 通知一个等待者
}

func (q *Queue) Pop() int {
    q.cond.L.Lock()
    defer q.cond.L.Unlock()
    
    // 等待队列非空
    for len(q.items) == 0 {
        q.cond.Wait() // 释放锁并等待信号
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item
}
```

### Signal vs Broadcast

- `Signal()`: 唤醒一个等待的 goroutine
- `Broadcast()`: 唤醒所有等待的 goroutine

## sync.Map

并发安全的 map，适用于以下场景：

1. key 只会写入一次，但读取多次
2. 多个 goroutine 读写不相交的 key 集合

```go
var m sync.Map

// 存储
m.Store("key", "value")

// 读取
value, ok := m.Load("key")

// 读取或存储
actual, loaded := m.LoadOrStore("key", "default")

// 删除
m.Delete("key")

// 遍历
m.Range(func(key, value interface{}) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true // 返回 false 停止遍历
})
```

### ⚠️ sync.Map 的适用场景

```go
// ❌ 不适合：频繁写入的场景
// sync.Map 写入性能比 map + RWMutex 差

// ✅ 适合：读多写少，或 key 分布不均匀
// 例如：缓存、配置、连接池
```

## sync.Pool

对象池，用于复用临时对象，减少内存分配。

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func GetBuffer() *bytes.Buffer {
    return bufferPool.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
    buf.Reset()
    bufferPool.Put(buf)
}

// 使用示例
func ProcessData(data []byte) string {
    buf := GetBuffer()
    defer PutBuffer(buf)
    
    buf.Write(data)
    // 处理数据...
    return buf.String()
}
```

### ⚠️ sync.Pool 注意事项

1. **Pool 中的对象可能随时被 GC 回收**
2. **不能用于连接池**（连接有状态）
3. **放回前要重置对象状态**

```go
// ❌ 危险：放回的 buffer 还有旧数据
bufferPool.Put(buf)

// ✅ 正确：重置后再放回
buf.Reset()
bufferPool.Put(buf)
```

## atomic 包

提供原子操作，比锁更轻量。

```go
import "sync/atomic"

var counter int64

// 原子增加
atomic.AddInt64(&counter, 1)

// 原子读取
value := atomic.LoadInt64(&counter)

// 原子存储
atomic.StoreInt64(&counter, 100)

// CAS (Compare And Swap)
swapped := atomic.CompareAndSwapInt64(&counter, 100, 200)
```

### Go 1.19+ 新增的 atomic 类型

```go
var counter atomic.Int64

counter.Add(1)
counter.Load()
counter.Store(100)
counter.CompareAndSwap(100, 200)
```

## 选择合适的同步原语

| 场景 | 推荐方案 |
|------|----------|
| 简单计数器 | atomic |
| 读多写少 | sync.RWMutex |
| 临时对象复用 | sync.Pool |
| 单例初始化 | sync.Once |
| 等待多个 goroutine | sync.WaitGroup |
| 并发安全 map | sync.Map 或 map + RWMutex |
| goroutine 间通信 | channel (首选) |

## 参考资源

- [Go sync package](https://pkg.go.dev/sync)
- [Go atomic package](https://pkg.go.dev/sync/atomic)
- [The Go Memory Model](https://go.dev/ref/mem)
