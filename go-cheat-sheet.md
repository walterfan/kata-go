- [无废话 Go 语言手册](#无废话-go-语言手册)
  - [1. 基本类型](#1-基本类型)
  - [2. 变量与常量](#2-变量与常量)
  - [3. 表达式](#3-表达式)
  - [4. 流程控制](#4-流程控制)
    - [分支语句](#分支语句)
    - [循环语句](#循环语句)
  - [5. 复合类型](#5-复合类型)
    - [数组](#数组)
    - [切片](#切片)
    - [字符串](#字符串)
    - [映射](#映射)
    - [结构体](#结构体)
  - [6. 函数](#6-函数)
  - [7. 指针](#7-指针)
  - [8. 接口](#8-接口)
  - [9. 错误处理](#9-错误处理)
  - [10. 模块与包](#10-模块与包)
  - [11. 并发编程](#11-并发编程)
  - [12. 标准库](#12-标准库)
  - [13. 常用库](#13-常用库)
  - [14. 常用框架](#14-常用框架)
  - [15. unsafe](#15-unsafe)
  - [16. 命令行应用](#16-命令行应用)
  - [17. 数据库操作](#17-数据库操作)
  - [18. 文件读写](#18-文件读写)
  - [19. grpc](#19-grpc)
  - [20. web 应用](#20-web-应用)
  - [21. 设计模式](#21-设计模式)


# 无废话 Go 语言手册

## 1. 基本类型
- **整数**: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`
- **浮点数**: `float32`, `float64`
- **复数**: `complex64`, `complex128`
- **布尔值**: `bool`
- **字符串**: `string`
- **字节**: `byte` (alias for `uint8`)
- **空接口**: `interface{}`

## 2. 变量与常量
- **变量声明**:
  ```go
  var x int
  var y = 10
  x, y := 5, 6
  ```
- **常量声明**:
  ```go
  const Pi = 3.14
  const (
    A = 1
    B = 2
  )
  ```

## 3. 表达式
- **算术运算符**: `+`, `-`, `*`, `/`, `%`
- **关系运算符**: `==`, `!=`, `<`, `>`, `<=`, `>=`
- **逻辑运算符**: `&&`, `||`, `!`
- **位运算符**: `&`, `|`, `^`, `<<`, `>>`
- **类型转换**: `float64(x)`, `int(y)`

## 4. 流程控制

### 分支语句
- **if .. else**:
  ```go
  if x > 0 {
    fmt.Println("Positive")
  } else {
    fmt.Println("Non-positive")
  }
  ```
- **switch .. case**:
  ```go
  switch day {
    case 1:
      fmt.Println("Monday")
    case 2:
      fmt.Println("Tuesday")
    default:
      fmt.Println("Other day")
  }
  ```

### 循环语句
- **for**:
  ```go
  for i := 0; i < 10; i++ {
    fmt.Println(i)
  }
  ```
  - Infinite loop:
    ```go
    for {
      // infinite loop
    }
    ```

## 5. 复合类型

### 数组
- **声明与初始化**:
  ```go
  var arr [5]int
  arr := [5]int{1, 2, 3, 4, 5}
  ```

### 切片
- **声明与初始化**:
  ```go
  var slice []int
  slice := []int{1, 2, 3}
  ```
- **切片操作**:
  ```go
  slice = append(slice, 4)
  slice = slice[1:3]
  ```

### 字符串
- **基本操作**:
  ```go
  str := "Hello, World!"
  fmt.Println(len(str))
  fmt.Println(str[0]) // byte value
  ```

### 映射
- **声明与初始化**:
  ```go
  var m map[string]int
  m := map[string]int{"a": 1, "b": 2}
  ```
- **访问与删除**:
  ```go
  value := m["a"]
  delete(m, "b")
  ```

### 结构体
- **声明与初始化**:
  ```go
  type Person struct {
    Name string
    Age  int
  }
  
  p := Person{"Alice", 30}
  ```

## 6. 函数
```go
func name(parameter-list)(result-list) {
  body
}
```
- **返回多个值**:
  ```go
  func swap(a, b int) (int, int) {
    return b, a
  }
  ```

## 7. 指针
- **声明与使用**:
  ```go
  var ptr *int
  x := 58
  ptr = &x
  fmt.Println(*ptr) // Dereference
  ```

## 8. 接口
- **定义与实现**:
  ```go
  type Speaker interface {
    Speak() string
  }

  type Person struct {
    Name string
  }

  func (p Person) Speak() string {
    return "Hello, " + p.Name
  }
  ```

## 9. 错误处理
- **基本错误处理**:
  ```go
  if err != nil {
    fmt.Println("Error:", err)
  }
  ```

## 10. 模块与包
- **导入包**:
  ```go
  import "fmt"
  ```
- **自定义包**:
  - Create a file `mypackage.go`:
    ```go
    package mypackage
    var Var = "Hello"
    ```

## 11. 并发编程
- **goroutine**:
  ```go
  go func() {
    fmt.Println("Running concurrently")
  }()
  ```
- **Channel**:
  ```go
  ch := make(chan int)
  go func() {
    ch <- 42
  }()
  value := <-ch
  fmt.Println(value)
  ```

## 12. 标准库
- **fmt**: 格式化输出
- **math**: 数学函数
- **strings**: 字符串处理
- **os**: 操作系统相关操作
- **net/http**: HTTP 请求与响应

## 13. 常用库
- **Gorilla Mux**: 高效的 HTTP 路由器
- **Go-Redis**: Redis 客户端
- **Gin**: Web 框架

## 14. 常用框架
- **Gin**: 高效 Web 框架
- **Echo**: Web 框架
- **Beego**: Web 框架
- **Revel**: MVC Web 框架

## 15. unsafe
- **使用 unsafe**:
  ```go
  import "unsafe"
  ptr := unsafe.Pointer(&x)
  ```

## 16. 命令行应用
- **使用 `os.Args` 获取命令行参数**:
  ```go
  import "os"
  fmt.Println(os.Args)
  ```

## 17. 数据库操作
- **MySQL 操作**:
  ```go
  import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
  )

  db, err := sql.Open("mysql", "user:password@/dbname")
  ```

## 18. 文件读写
- **读取文件**:
  ```go
  data, err := ioutil.ReadFile("file.txt")
  ```
- **写入文件**:
  ```go
  err := ioutil.WriteFile("file.txt", []byte("Hello, World!"), 0644)
  ```

## 19. grpc
- **定义 proto 文件**:
  ```proto
  service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
  }
  ```
- **Go 代码**:
  ```go
  import "google.golang.org/grpc"
  ```

## 20. web 应用
- **使用 net/http**:
  ```go
  import "net/http"

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
  })
  http.ListenAndServe(":8080", nil)
  ```

## 21. 设计模式
- **单例模式**:
  ```go
  package singleton
  import "sync"

  var instance *Singleton
  var once sync.Once

  type Singleton struct{}

  func GetInstance() *Singleton {
    once.Do(func() {
      instance = &Singleton{}
    })
    return instance
  }
  ```

