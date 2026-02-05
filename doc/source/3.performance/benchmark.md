# 基准测试 (Benchmark)

```{contents} 目录
:depth: 3
```

## 基准测试基础

Go 内置了基准测试框架，用于测量代码性能。

### 基本语法

```go
// benchmark_test.go
package mypackage

import "testing"

func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 被测试的代码
        MyFunction()
    }
}
```

### 运行基准测试

```bash
# 运行所有基准测试
go test -bench=.

# 运行特定基准测试
go test -bench=BenchmarkFunction

# 运行并显示内存分配
go test -bench=. -benchmem

# 指定运行时间
go test -bench=. -benchtime=5s

# 运行多次取平均
go test -bench=. -count=5
```

## 理解测试结果

```
BenchmarkFunction-8    1000000    1234 ns/op    256 B/op    4 allocs/op
```

| 字段 | 含义 |
|------|------|
| `BenchmarkFunction-8` | 测试名称-GOMAXPROCS |
| `1000000` | 运行次数 |
| `1234 ns/op` | 每次操作耗时 |
| `256 B/op` | 每次操作分配的字节数 |
| `4 allocs/op` | 每次操作的内存分配次数 |

## 高级基准测试技巧

### 1. 子基准测试

```go
func BenchmarkConcat(b *testing.B) {
    sizes := []int{10, 100, 1000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                concatStrings(size)
            }
        })
    }
}
```

### 2. 重置计时器

```go
func BenchmarkExpensiveSetup(b *testing.B) {
    // 昂贵的初始化
    data := loadLargeDataSet()
    
    b.ResetTimer() // 重置计时器，不计算初始化时间
    
    for i := 0; i < b.N; i++ {
        processData(data)
    }
}
```

### 3. 暂停/恢复计时

```go
func BenchmarkWithCleanup(b *testing.B) {
    for i := 0; i < b.N; i++ {
        data := generateData()
        
        b.StartTimer()
        processData(data)
        b.StopTimer()
        
        cleanupData(data)
    }
}
```

### 4. 并行基准测试

```go
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // 并行执行的代码
            myFunction()
        }
    })
}
```

### 5. 报告自定义指标

```go
func BenchmarkThroughput(b *testing.B) {
    totalBytes := int64(0)
    
    for i := 0; i < b.N; i++ {
        n := processData(data)
        totalBytes += int64(n)
    }
    
    b.SetBytes(totalBytes / int64(b.N)) // 报告吞吐量
}
```

## 实际案例

### 案例 1：字符串拼接对比

```go
func BenchmarkStringConcat(b *testing.B) {
    strs := []string{"Hello", " ", "World", "!"}
    
    b.Run("Plus", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var s string
            for _, str := range strs {
                s += str
            }
        }
    })
    
    b.Run("StringBuilder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var sb strings.Builder
            for _, str := range strs {
                sb.WriteString(str)
            }
            _ = sb.String()
        }
    })
    
    b.Run("Join", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = strings.Join(strs, "")
        }
    })
}
```

### 案例 2：Map vs Slice 查找

```go
func BenchmarkLookup(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    
    for _, size := range sizes {
        // 准备数据
        slice := make([]int, size)
        m := make(map[int]bool, size)
        for i := 0; i < size; i++ {
            slice[i] = i
            m[i] = true
        }
        target := size / 2
        
        b.Run(fmt.Sprintf("Slice-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                for _, v := range slice {
                    if v == target {
                        break
                    }
                }
            }
        })
        
        b.Run(fmt.Sprintf("Map-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                _ = m[target]
            }
        })
    }
}
```

### 案例 3：sync.Pool vs 直接分配

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func BenchmarkAllocation(b *testing.B) {
    b.Run("DirectAlloc", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            buf := make([]byte, 1024)
            _ = buf
        }
    })
    
    b.Run("SyncPool", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            buf := bufferPool.Get().([]byte)
            bufferPool.Put(buf)
        }
    })
}
```

## 使用 benchstat 比较结果

### 安装

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

### 使用

```bash
# 运行旧版本并保存结果
go test -bench=. -count=10 > old.txt

# 修改代码后运行新版本
go test -bench=. -count=10 > new.txt

# 比较结果
benchstat old.txt new.txt
```

### 输出示例

```
name          old time/op  new time/op  delta
Function-8    1.23µs ± 2%  0.98µs ± 1%  -20.33%  (p=0.001 n=10+10)

name          old alloc/op new alloc/op delta
Function-8    256B ± 0%    128B ± 0%    -50.00%  (p=0.001 n=10+10)
```

## 常见陷阱

### 陷阱 1：编译器优化导致代码被消除

```go
// ❌ 结果可能被优化掉
func BenchmarkBad(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = compute() // 编译器可能消除这行
    }
}

// ✅ 使用全局变量防止优化
var result int

func BenchmarkGood(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = compute()
    }
    result = r // 防止编译器优化
}
```

### 陷阱 2：未预热缓存

```go
func BenchmarkWithWarmup(b *testing.B) {
    // 预热
    for i := 0; i < 100; i++ {
        myFunction()
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        myFunction()
    }
}
```

### 陷阱 3：测试数据不具代表性

```go
// ❌ 只测试最好情况
func BenchmarkBad(b *testing.B) {
    data := []int{1, 2, 3} // 数据太小
    for i := 0; i < b.N; i++ {
        sort(data)
    }
}

// ✅ 测试多种情况
func BenchmarkGood(b *testing.B) {
    for _, size := range []int{10, 100, 1000, 10000} {
        b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
            data := generateRandomData(size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                dataCopy := make([]int, len(data))
                copy(dataCopy, data)
                sort(dataCopy)
            }
        })
    }
}
```

## 最佳实践

1. **使用 `-benchmem`**：总是关注内存分配
2. **多次运行**：使用 `-count=10` 确保结果稳定
3. **使用 benchstat**：科学比较优化前后的差异
4. **测试真实数据**：使用生产环境的数据分布
5. **避免编译器优化**：使用全局变量保存结果
6. **隔离测试环境**：关闭其他程序，固定 CPU 频率

## 参考资源

- [Go Testing Package](https://pkg.go.dev/testing)
- [How to Write Benchmarks in Go](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
