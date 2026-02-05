# Go 语法速查表

```{contents} 目录
:depth: 2
```

## 变量声明

```go
// 完整声明
var name string = "Go"

// 类型推断
var name = "Go"

// 短声明（只能在函数内）
name := "Go"

// 多变量声明
var a, b, c int
x, y := 1, 2

// 常量
const Pi = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)

// iota 枚举
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
)
```

## 基本类型

```go
// 数值类型
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
float32, float64
complex64, complex128

// 其他
bool
string
byte    // uint8 别名
rune    // int32 别名，表示 Unicode 码点

// 零值
数值: 0
布尔: false
字符串: ""
指针/切片/map/channel/函数/接口: nil
```

## 控制结构

```go
// if
if x > 0 {
    // ...
} else if x < 0 {
    // ...
} else {
    // ...
}

// if 带初始化
if v := compute(); v > 0 {
    // v 的作用域仅在 if 块内
}

// for
for i := 0; i < 10; i++ { }

// while 风格
for condition { }

// 无限循环
for { }

// range
for i, v := range slice { }
for k, v := range map { }
for i, c := range "string" { }  // c 是 rune

// switch
switch x {
case 1:
    // 自动 break
case 2, 3:
    // 多值匹配
default:
    // 默认
}

// switch 无表达式
switch {
case x > 0:
    // ...
case x < 0:
    // ...
}

// type switch
switch v := i.(type) {
case int:
    // v 是 int
case string:
    // v 是 string
}
```

## 函数

```go
// 基本函数
func add(a, b int) int {
    return a + b
}

// 多返回值
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 命名返回值
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return
}

// 可变参数
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 闭包
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// defer
func readFile(name string) {
    f, _ := os.Open(name)
    defer f.Close()  // 函数返回时执行
    // 使用 f
}
```

## 结构体

```go
// 定义
type Person struct {
    Name string
    Age  int
}

// 创建
p1 := Person{Name: "Alice", Age: 30}
p2 := Person{"Bob", 25}
p3 := new(Person)  // *Person

// 匿名结构体
point := struct {
    X, Y int
}{10, 20}

// 嵌入
type Employee struct {
    Person  // 嵌入 Person
    Salary float64
}

// 方法
func (p Person) Greet() string {
    return "Hello, " + p.Name
}

func (p *Person) SetAge(age int) {
    p.Age = age
}
```

## 接口

```go
// 定义
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}

// 空接口
var any interface{}
any = 42
any = "hello"

// 类型断言
s, ok := any.(string)
if ok {
    fmt.Println(s)
}
```

## 切片

```go
// 创建
s1 := []int{1, 2, 3}
s2 := make([]int, 5)      // len=5, cap=5
s3 := make([]int, 0, 10)  // len=0, cap=10

// 追加
s1 = append(s1, 4, 5)
s1 = append(s1, s2...)  // 追加另一个切片

// 切片操作
s := []int{0, 1, 2, 3, 4}
s[1:3]   // [1, 2]
s[:3]    // [0, 1, 2]
s[2:]    // [2, 3, 4]
s[:]     // [0, 1, 2, 3, 4]

// 复制
dst := make([]int, len(src))
copy(dst, src)

// 删除元素
s = append(s[:i], s[i+1:]...)
```

## Map

```go
// 创建
m1 := map[string]int{"a": 1, "b": 2}
m2 := make(map[string]int)

// 操作
m["key"] = value      // 设置
v := m["key"]         // 获取
v, ok := m["key"]     // 检查是否存在
delete(m, "key")      // 删除
len(m)                // 长度

// 遍历
for k, v := range m {
    fmt.Println(k, v)
}
```

## Channel

```go
// 创建
ch := make(chan int)        // 无缓冲
ch := make(chan int, 10)    // 有缓冲

// 操作
ch <- value  // 发送
v := <-ch    // 接收
v, ok := <-ch  // 检查是否关闭
close(ch)    // 关闭

// select
select {
case v := <-ch1:
    // ...
case ch2 <- x:
    // ...
case <-time.After(time.Second):
    // 超时
default:
    // 非阻塞
}
```

## 错误处理

```go
// 检查错误
if err != nil {
    return err
}

// 创建错误
err := errors.New("something went wrong")
err := fmt.Errorf("failed to process %s: %w", name, originalErr)

// 错误包装 (Go 1.13+)
if errors.Is(err, os.ErrNotExist) { }
var pathErr *os.PathError
if errors.As(err, &pathErr) { }

// panic/recover
func safeCall() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    panic("oops")
}
```

## 并发

```go
// Goroutine
go func() {
    // 并发执行
}()

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // ...
}()
wg.Wait()

// Mutex
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

// Once
var once sync.Once
once.Do(func() {
    // 只执行一次
})

// Context
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
```

## 常用包

```go
import (
    "fmt"      // 格式化 I/O
    "os"       // 操作系统功能
    "io"       // I/O 原语
    "strings"  // 字符串操作
    "strconv"  // 类型转换
    "time"     // 时间
    "encoding/json"  // JSON
    "net/http"       // HTTP
    "context"        // 上下文
    "sync"           // 同步原语
)
```
