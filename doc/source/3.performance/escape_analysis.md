# 逃逸分析

```{contents} 目录
:depth: 3
```

## 什么是逃逸分析

逃逸分析是 Go 编译器的一项优化技术，用于决定变量应该分配在栈上还是堆上。

| 分配位置 | 特点 |
|----------|------|
| **栈** | 快速分配/释放，函数返回自动清理 |
| **堆** | 需要 GC 管理，分配较慢 |

## 查看逃逸分析结果

```bash
go build -gcflags='-m' main.go
go build -gcflags='-m -m' main.go  # 更详细
```

## 常见逃逸场景

### 场景 1：返回局部变量的指针

```go
// 发生逃逸：返回局部变量的指针
func createUser() *User {
    u := User{Name: "Alice"}  // u 逃逸到堆
    return &u
}

// 输出：moved to heap: u
```

### 场景 2：接口类型

```go
// 发生逃逸：赋值给接口类型
func printAny(v interface{}) {
    fmt.Println(v)
}

func main() {
    x := 42
    printAny(x)  // x 逃逸到堆（需要装箱）
}
```

### 场景 3：闭包引用

```go
// 发生逃逸：闭包引用外部变量
func closure() func() int {
    x := 0  // x 逃逸到堆
    return func() int {
        x++
        return x
    }
}
```

### 场景 4：切片扩容

```go
// 可能逃逸：切片容量不确定
func appendSlice(s []int) []int {
    return append(s, 1, 2, 3)  // 如果需要扩容，新切片可能在堆上
}
```

### 场景 5：大对象

```go
// 发生逃逸：对象太大
func createLargeArray() [1024 * 1024]byte {
    var arr [1024 * 1024]byte  // 太大，逃逸到堆
    return arr
}
```

### 场景 6：动态类型

```go
// 发生逃逸：类型在编译期未知
func createValue(typ string) interface{} {
    switch typ {
    case "int":
        return 42
    case "string":
        return "hello"
    }
    return nil
}
```

## 避免逃逸的技巧

### 技巧 1：使用值类型而非指针

```go
// ❌ 逃逸
func getUser() *User {
    return &User{Name: "Alice"}
}

// ✅ 不逃逸
func getUser() User {
    return User{Name: "Alice"}
}
```

### 技巧 2：预分配切片容量

```go
// ❌ 可能逃逸
func process(n int) []int {
    var result []int
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}

// ✅ 减少逃逸风险
func process(n int) []int {
    result := make([]int, 0, n)  // 预分配
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}
```

### 技巧 3：避免在循环中使用接口

```go
// ❌ 每次迭代都逃逸
func sum(nums []int) int {
    var total interface{} = 0
    for _, n := range nums {
        total = total.(int) + n
    }
    return total.(int)
}

// ✅ 使用具体类型
func sum(nums []int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
```

### 技巧 4：使用 sync.Pool 复用对象

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process(data []byte) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    buf.Write(data)
    return buf.String()
}
```

### 技巧 5：在函数内使用固定大小数组

```go
// ❌ 切片逃逸
func hash(data []byte) []byte {
    result := make([]byte, 32)
    // ...
    return result
}

// ✅ 数组不逃逸（如果不返回指针）
func hash(data []byte, out *[32]byte) {
    // 写入 out
}
```

## 逃逸分析实战

### 示例：JSON 编码优化

```go
// ❌ 多次逃逸
func toJSON(v interface{}) ([]byte, error) {
    return json.Marshal(v)
}

// ✅ 使用 Encoder 减少分配
func toJSON(w io.Writer, v interface{}) error {
    enc := json.NewEncoder(w)
    return enc.Encode(v)
}

// ✅ 使用 sync.Pool
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func toJSONBytes(v interface{}) ([]byte, error) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    if err := json.NewEncoder(buf).Encode(v); err != nil {
        return nil, err
    }
    
    // 复制结果，因为 buf 会被复用
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result, nil
}
```

### 示例：HTTP Handler 优化

```go
// ❌ 每次请求都分配
func handler(w http.ResponseWriter, r *http.Request) {
    response := &Response{
        Code:    200,
        Message: "OK",
        Data:    getData(),
    }
    json.NewEncoder(w).Encode(response)
}

// ✅ 复用 Response 结构
var responsePool = sync.Pool{
    New: func() interface{} {
        return new(Response)
    },
}

func handler(w http.ResponseWriter, r *http.Request) {
    response := responsePool.Get().(*Response)
    defer responsePool.Put(response)
    
    response.Code = 200
    response.Message = "OK"
    response.Data = getData()
    
    json.NewEncoder(w).Encode(response)
}
```

## 何时关注逃逸

| 场景 | 是否需要关注 |
|------|-------------|
| 高性能服务的热点路径 | ✅ 是 |
| 批处理任务 | ❌ 通常不需要 |
| 每秒处理数万请求 | ✅ 是 |
| 一次性脚本 | ❌ 不需要 |
| 内存受限环境 | ✅ 是 |

## 参考资源

- [Go Escape Analysis](https://go.dev/wiki/CompilerOptimizations#escape-analysis)
- [Allocation Efficiency in High-Performance Go Services](https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/)
