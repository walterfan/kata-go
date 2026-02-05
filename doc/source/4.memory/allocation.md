# 内存分配

```{contents} 目录
:depth: 3
```

## Go 内存分配器

Go 使用 **TCMalloc** 风格的内存分配器，特点：

- **分级管理**：小对象、大对象分开处理
- **线程缓存**：减少锁竞争
- **Span 管理**：内存按 Span 管理

## 内存布局

```{mermaid}
graph TD
    subgraph "Go 内存分配器"
        H[Heap]
        H --> A[Arena 64MB]
        A --> S1[Span 8KB]
        A --> S2[Span 8KB]
        S1 --> O1[Object]
        S1 --> O2[Object]
    end
    
    subgraph "P 本地缓存"
        MC[mcache]
        MC --> SC1[Span Class 1]
        MC --> SC2[Span Class 2]
    end
```

## 对象大小分类

| 类别 | 大小范围 | 分配方式 |
|------|----------|----------|
| Tiny | < 16B 且无指针 | Tiny 分配器 |
| Small | 16B - 32KB | mcache → mcentral → mheap |
| Large | > 32KB | 直接从 mheap 分配 |

## 分配流程

### 小对象分配

```go
// 小对象分配流程
// 1. 从 P 的 mcache 获取对应 size class 的 span
// 2. 如果 span 中有空闲对象，直接分配
// 3. 如果 span 已满，从 mcentral 获取新 span
// 4. 如果 mcentral 也没有，从 mheap 分配
```

### 大对象分配

```go
// 大对象直接从 mheap 分配
// 可能触发 GC
bigSlice := make([]byte, 64*1024) // 64KB
```

## 内存分配函数

### new vs make

```go
// new: 分配零值，返回指针
p := new(int)       // *int, 值为 0
s := new(MyStruct)  // *MyStruct, 字段为零值

// make: 初始化 slice, map, channel
s := make([]int, 10)      // 长度 10 的 slice
m := make(map[string]int) // 空 map
c := make(chan int, 10)   // 容量 10 的 channel
```

### 栈分配 vs 堆分配

```go
// 栈分配（快，自动释放）
func stackAlloc() {
    x := 42  // 通常在栈上
    _ = x
}

// 堆分配（需要 GC）
func heapAlloc() *int {
    x := 42
    return &x  // x 逃逸到堆
}
```

## 减少内存分配

### 1. 预分配

```go
// ❌ 多次扩容
var slice []int
for i := 0; i < 1000; i++ {
    slice = append(slice, i)
}

// ✅ 预分配容量
slice := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    slice = append(slice, i)
}
```

### 2. 复用 buffer

```go
var bufPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 4096)
    },
}

func process(data []byte) []byte {
    buf := bufPool.Get().([]byte)
    buf = buf[:0]  // 复用但重置长度
    defer bufPool.Put(buf)
    
    // 使用 buf
    buf = append(buf, data...)
    result := make([]byte, len(buf))
    copy(result, buf)
    return result
}
```

### 3. 使用数组而非切片

```go
// 切片有额外的 header (24 bytes)
func withSlice(data []byte) {
    // slice header + 底层数组
}

// 数组直接存储
func withArray(data [1024]byte) {
    // 直接是数据，无额外开销
}
```

### 4. 内联小结构体

```go
// ❌ 指针增加分配
type Node struct {
    Value *Data
    Next  *Node
}

// ✅ 内联值（如果 Data 较小）
type Node struct {
    Value Data  // 直接嵌入
    Next  *Node
}
```

## 内存对齐

Go 会自动进行内存对齐，但了解对齐规则有助于优化结构体大小。

### 查看对齐

```go
import "unsafe"

type Example struct {
    a bool    // 1 byte
    b int64   // 8 bytes
    c bool    // 1 byte
}

func main() {
    fmt.Println(unsafe.Sizeof(Example{}))  // 24 (因为对齐)
    fmt.Println(unsafe.Alignof(Example{})) // 8
}
```

### 优化结构体布局

```go
// ❌ 浪费空间（24 bytes）
type Bad struct {
    a bool    // 1 + 7 padding
    b int64   // 8
    c bool    // 1 + 7 padding
}

// ✅ 优化布局（16 bytes）
type Good struct {
    b int64   // 8
    a bool    // 1
    c bool    // 1 + 6 padding
}
```

### 使用 fieldalignment 工具

```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
fieldalignment -fix ./...
```

## 内存统计

```go
import "runtime"

func printMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc = %v KB\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc = %v KB\n", m.TotalAlloc/1024)
    fmt.Printf("Sys = %v KB\n", m.Sys/1024)
    fmt.Printf("Mallocs = %v\n", m.Mallocs)
    fmt.Printf("Frees = %v\n", m.Frees)
    fmt.Printf("HeapAlloc = %v KB\n", m.HeapAlloc/1024)
    fmt.Printf("HeapSys = %v KB\n", m.HeapSys/1024)
    fmt.Printf("HeapIdle = %v KB\n", m.HeapIdle/1024)
    fmt.Printf("HeapInuse = %v KB\n", m.HeapInuse/1024)
    fmt.Printf("StackSys = %v KB\n", m.StackSys/1024)
}
```

## 最佳实践

1. **使用 `-benchmem` 测量分配**
2. **使用 `sync.Pool` 复用对象**
3. **预分配切片和 map 容量**
4. **优化结构体字段顺序**
5. **避免不必要的指针**
6. **使用 pprof 分析热点**

## 参考资源

- [A visual guide to Go Memory Allocator](https://blog.learngoprogramming.com/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132f7c6c2c3)
- [Go Memory Management](https://povilasv.me/go-memory-management/)
