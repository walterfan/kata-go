# 性能优化技巧

```{contents} 目录
:depth: 3
```

## 字符串优化

### 使用 strings.Builder

```go
// ❌ 每次拼接都分配新内存
func concat(strs []string) string {
    var result string
    for _, s := range strs {
        result += s
    }
    return result
}

// ✅ 使用 Builder
func concat(strs []string) string {
    var sb strings.Builder
    for _, s := range strs {
        sb.WriteString(s)
    }
    return sb.String()
}

// ✅ 预分配容量
func concat(strs []string) string {
    var sb strings.Builder
    size := 0
    for _, s := range strs {
        size += len(s)
    }
    sb.Grow(size)
    for _, s := range strs {
        sb.WriteString(s)
    }
    return sb.String()
}
```

### 字符串与字节切片转换

```go
// ❌ 产生拷贝
s := string(bytes)
b := []byte(str)

// ✅ 零拷贝转换（unsafe，谨慎使用）
import "unsafe"

func stringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}

func bytesToString(b []byte) string {
    return unsafe.String(unsafe.SliceData(b), len(b))
}
```

## 切片优化

### 预分配容量

```go
// ❌ 多次扩容
func createSlice(n int) []int {
    var s []int
    for i := 0; i < n; i++ {
        s = append(s, i)
    }
    return s
}

// ✅ 预分配
func createSlice(n int) []int {
    s := make([]int, 0, n)
    for i := 0; i < n; i++ {
        s = append(s, i)
    }
    return s
}
```

### 避免切片引用导致的内存泄漏

```go
// ❌ 小切片引用大数组，导致大数组无法被 GC
func getFirstElement(data []int) []int {
    return data[:1]  // 仍然引用原始大数组
}

// ✅ 复制数据
func getFirstElement(data []int) []int {
    result := make([]int, 1)
    copy(result, data[:1])
    return result
}
```

### 复用切片

```go
// ✅ 使用 slice[:0] 复用底层数组
func processInPlace(data []int) []int {
    result := data[:0]
    for _, v := range data {
        if v > 0 {
            result = append(result, v)
        }
    }
    return result
}
```

## Map 优化

### 预分配容量

```go
// ❌ 动态扩容
m := make(map[string]int)

// ✅ 预分配
m := make(map[string]int, expectedSize)
```

### 使用结构体 key 而非字符串

```go
// ❌ 字符串 key 需要 hash 和比较
type cacheKey string
cache := make(map[cacheKey]Value)

// ✅ 固定大小的 key 更高效
type cacheKey struct {
    userID  int64
    itemID  int64
}
cache := make(map[cacheKey]Value)
```

## 并发优化

### 减少锁竞争

```go
// ❌ 全局锁
type Counter struct {
    mu    sync.Mutex
    count int
}

// ✅ 分片锁减少竞争
type ShardedCounter struct {
    shards [256]struct {
        mu    sync.Mutex
        count int
    }
}

func (c *ShardedCounter) Inc(key string) {
    shard := &c.shards[hash(key)%256]
    shard.mu.Lock()
    shard.count++
    shard.mu.Unlock()
}
```

### 使用 atomic 代替 Mutex

```go
// ❌ 使用 Mutex
type Counter struct {
    mu    sync.Mutex
    count int64
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}

// ✅ 使用 atomic
type Counter struct {
    count atomic.Int64
}

func (c *Counter) Inc() {
    c.count.Add(1)
}
```

## 内存分配优化

### 使用 sync.Pool

```go
var bufPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process(data []byte) []byte {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)
    }()
    
    // 使用 buf 处理数据
    buf.Write(data)
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result
}
```

### 避免不必要的指针

```go
// ❌ 指针增加 GC 压力
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

## I/O 优化

### 使用 bufio

```go
// ❌ 每次读取都是系统调用
func readLines(filename string) ([]string, error) {
    f, _ := os.Open(filename)
    defer f.Close()
    
    var lines []string
    buf := make([]byte, 1)
    for {
        _, err := f.Read(buf)
        // ...
    }
    return lines, nil
}

// ✅ 使用 bufio
func readLines(filename string) ([]string, error) {
    f, _ := os.Open(filename)
    defer f.Close()
    
    var lines []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}
```

### 使用 io.Copy

```go
// ❌ 手动复制
func copyFile(src, dst string) error {
    data, _ := os.ReadFile(src)  // 全部读入内存
    return os.WriteFile(dst, data, 0644)
}

// ✅ 流式复制
func copyFile(src, dst string) error {
    srcFile, _ := os.Open(src)
    defer srcFile.Close()
    
    dstFile, _ := os.Create(dst)
    defer dstFile.Close()
    
    _, err := io.Copy(dstFile, srcFile)
    return err
}
```

## JSON 优化

### 使用 json.Decoder/Encoder

```go
// ❌ 中间缓冲区
data, _ := json.Marshal(obj)
w.Write(data)

// ✅ 直接写入
json.NewEncoder(w).Encode(obj)
```

### 使用 jsoniter 或 sonic

```go
import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func encode(v interface{}) ([]byte, error) {
    return json.Marshal(v)
}
```

## 编译优化

### 内联优化

```go
// 小函数会被自动内联
//go:noinline  // 禁用内联
func add(a, b int) int {
    return a + b
}
```

### 边界检查消除

```go
// ❌ 每次访问都检查边界
func sum(s []int) int {
    total := 0
    for i := 0; i < len(s); i++ {
        total += s[i]
    }
    return total
}

// ✅ 使用 range 避免边界检查
func sum(s []int) int {
    total := 0
    for _, v := range s {
        total += v
    }
    return total
}

// ✅ 手动消除边界检查
func sum(s []int) int {
    total := 0
    _ = s[len(s)-1]  // 边界检查提前
    for i := 0; i < len(s); i++ {
        total += s[i]  // 不再检查
    }
    return total
}
```

## 性能优化检查清单

```{admonition} 检查清单
:class: tip

1. [ ] 使用 `-race` 检查数据竞争
2. [ ] 使用 pprof 找出热点
3. [ ] 使用 `-benchmem` 检查内存分配
4. [ ] 预分配切片和 map 容量
5. [ ] 使用 strings.Builder 拼接字符串
6. [ ] 使用 sync.Pool 复用对象
7. [ ] 减少锁的粒度和持有时间
8. [ ] 使用 bufio 进行 I/O
9. [ ] 检查逃逸分析结果
10. [ ] 使用基准测试验证优化效果
```

## 参考资源

- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
- [Go Performance Tuning](https://go.dev/doc/diagnostics)
- [Effective Go](https://go.dev/doc/effective_go)
