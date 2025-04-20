```{contents} Table of Contents
:depth: 3
```
# Go Trap

在 Go 语言编程中，尽管其设计简洁，但仍然有一些常见错误和陷阱，容易让开发者踩坑。以下是一些典型的 Go 语言编程错误及如何避免它们的建议：

---

## **1. 忘记 `defer` 的执行顺序**
Go 的 `defer` 语句按照**后进先出（LIFO）**的顺序执行，很多开发者会误解其执行顺序。

### **错误示例**
```go
func main() {
    fmt.Println("Start")
    defer fmt.Println("First defer")
    defer fmt.Println("Second defer")
    fmt.Println("End")
}
```
### **执行结果**
```
Start
End
Second defer
First defer
```
**✅ 解决方案**：牢记 `defer` 语句的执行顺序，并合理利用它来清理资源。

---

## **2. `defer` 闭包捕获变量的坑**
`defer` 语句中的匿名函数会捕获变量，而不是捕获变量的值。

### **错误示例**
```go
func main() {
    for i := 0; i < 3; i++ {
        defer func() {
            fmt.Println(i) // 可能不是预期的 2, 1, 0
        }()
    }
}
```
### **实际输出**
```
3
3
3
```
**✅ 解决方案**：
在 `defer` 中传递参数，而不是捕获外部变量：
```go
func main() {
    for i := 0; i < 3; i++ {
        defer func(n int) {
            fmt.Println(n)
        }(i)
    }
}
```
输出：
```
2
1
0
```

---

## **3. `nil` 接收者方法可能导致 `panic`**
如果方法的接收者是 `nil`，调用它时可能会导致 `panic`。

### **错误示例**
```go
type Person struct {
    Name string
}

func (p *Person) Greet() {
    fmt.Println("Hello,", p.Name)
}

func main() {
    var p *Person
    p.Greet() // panic: runtime error: invalid memory address or nil pointer dereference
}
```
**✅ 解决方案**：
在方法内部检查 `nil`：
```go
func (p *Person) Greet() {
    if p == nil {
        fmt.Println("No person to greet!")
        return
    }
    fmt.Println("Hello,", p.Name)
}
```

---

## **4. `range` 循环变量的作用域**
在 `range` 循环中，循环变量是**复用的**，而不是每次都创建新的变量。这在处理 Goroutine 时可能会引发问题。

### **错误示例**
```go
func main() {
    words := []string{"Go", "is", "awesome"}
    for _, word := range words {
        go func() {
            fmt.Println(word)
        }()
    }
    time.Sleep(time.Second)
}
```
### **可能的错误输出**
```
awesome
awesome
awesome
```
**✅ 解决方案**：
显式传递变量：
```go
func main() {
    words := []string{"Go", "is", "awesome"}
    for _, word := range words {
        go func(w string) {
            fmt.Println(w)
        }(word)
    }
    time.Sleep(time.Second)
}
```
输出：
```
Go
is
awesome
```

---

## **5. `append` 切片的容量增长问题**
Go 的 `append` 可能会创建新的底层数组，导致原切片的 `append` 操作不会影响其他变量。

### **错误示例**
```go
func main() {
    nums := []int{1, 2, 3}
    newNums := append(nums, 4)
    nums[0] = 100
    fmt.Println(newNums) // 可能不是预期的 [1 2 3 4]
}
```
**✅ 解决方案**：
要确保 `append` 操作不会影响原切片，可使用 `copy`：
```go
newNums := make([]int, len(nums))
copy(newNums, nums)
newNums = append(newNums, 4)
```

---

## **6. `map` 并发读写导致 `fatal error: concurrent map read and map write`**
Go 的 `map` 不是线程安全的，直接并发读写 `map` 会导致 `panic`。

### **错误示例**
```go
var m = make(map[int]int)

func main() {
    go func() {
        for i := 0; i < 1000; i++ {
            m[i] = i
        }
    }()
    
    go func() {
        for i := 0; i < 1000; i++ {
            _ = m[i]
        }
    }()

    time.Sleep(time.Second)
}
```
### **可能的错误**
```
fatal error: concurrent map read and map write
```
**✅ 解决方案**：
使用 `sync.Mutex` 或 `sync.Map` 来保护并发访问：
```go
var mu sync.Mutex
var m = make(map[int]int)

func main() {
    go func() {
        for i := 0; i < 1000; i++ {
            mu.Lock()
            m[i] = i
            mu.Unlock()
        }
    }()

    go func() {
        for i := 0; i < 1000; i++ {
            mu.Lock()
            _ = m[i]
            mu.Unlock()
        }
    }()

    time.Sleep(time.Second)
}
```
或者使用 `sync.Map`：
```go
var sm sync.Map

func main() {
    go func() {
        for i := 0; i < 1000; i++ {
            sm.Store(i, i)
        }
    }()

    go func() {
        for i := 0; i < 1000; i++ {
            sm.Load(i)
        }
    }()

    time.Sleep(time.Second)
}
```

---

## **7. `interface{}` 断言失败**
Go 允许使用 `interface{}` 作为通用类型，但如果类型断言失败，会导致 `panic`。

### **错误示例**
```go
func main() {
    var x interface{} = "hello"
    y := x.(int) // panic: interface conversion: string is not int
    fmt.Println(y)
}
```
**✅ 解决方案**：
使用**类型判断**：
```go
if y, ok := x.(int); ok {
    fmt.Println(y)
} else {
    fmt.Println("x is not an int")
}
```

---

## **总结**
Go 语言的设计旨在简洁和高效，但仍然有一些常见的坑点需要注意：
1. `defer` 执行顺序
2. `defer` 闭包捕获变量
3. `nil` 接收者方法导致 `panic`
4. `range` 循环变量作用域问题
5. `append` 可能导致切片数据不一致
6. `map` 并发读写导致 `panic`
7. `interface{}` 断言失败

**掌握这些坑的解决方案，可以让你的 Go 代码更加健壮和高效！** 🚀

