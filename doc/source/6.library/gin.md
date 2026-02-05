# Gin Web 框架

```{contents} 目录
:depth: 3
```

## Gin 概述

Gin 是 Go 语言中最流行的 Web 框架，以高性能和简洁 API 著称。

## 安装

```bash
go get github.com/gin-gonic/gin
```

## 基本用法

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()  // 包含 Logger 和 Recovery 中间件
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    r.Run(":8080")
}
```

## 路由

### HTTP 方法

```go
r.GET("/users", getUsers)
r.POST("/users", createUser)
r.PUT("/users/:id", updateUser)
r.DELETE("/users/:id", deleteUser)
r.PATCH("/users/:id", patchUser)
r.HEAD("/users", headUsers)
r.OPTIONS("/users", optionsUsers)

// 匹配所有方法
r.Any("/test", handleAny)
```

### 路由参数

```go
// 路径参数
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})

// 通配符
r.GET("/files/*filepath", func(c *gin.Context) {
    path := c.Param("filepath")
    c.JSON(200, gin.H{"path": path})
})

// 查询参数
r.GET("/search", func(c *gin.Context) {
    query := c.Query("q")           // ?q=xxx
    page := c.DefaultQuery("page", "1")  // 默认值
    c.JSON(200, gin.H{"query": query, "page": page})
})
```

### 路由分组

```go
v1 := r.Group("/api/v1")
{
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)
}

v2 := r.Group("/api/v2")
{
    v2.GET("/users", getUsersV2)
}

// 带中间件的分组
authorized := r.Group("/admin")
authorized.Use(AuthMiddleware())
{
    authorized.GET("/dashboard", dashboard)
}
```

## 请求处理

### 绑定请求数据

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}

func createUser(c *gin.Context) {
    var req CreateUserRequest
    
    // 绑定 JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理请求...
    c.JSON(201, gin.H{"user": req})
}

// 其他绑定方法
c.ShouldBindQuery(&req)       // 查询参数
c.ShouldBindUri(&req)         // 路径参数
c.ShouldBind(&req)            // 自动推断
c.ShouldBindHeader(&req)      // Header
```

### 获取请求数据

```go
// Header
auth := c.GetHeader("Authorization")

// Cookie
cookie, _ := c.Cookie("session")

// 表单数据
name := c.PostForm("name")
name := c.DefaultPostForm("name", "default")

// 文件上传
file, _ := c.FormFile("file")
c.SaveUploadedFile(file, "/tmp/"+file.Filename)

// 原始 Body
body, _ := c.GetRawData()
```

## 响应

```go
// JSON
c.JSON(200, gin.H{"status": "ok"})
c.JSON(200, user)  // 结构体

// XML
c.XML(200, gin.H{"status": "ok"})

// YAML
c.YAML(200, gin.H{"status": "ok"})

// String
c.String(200, "Hello %s", name)

// HTML
c.HTML(200, "index.html", gin.H{"title": "Home"})

// 重定向
c.Redirect(301, "/new-path")

// 文件
c.File("/path/to/file")
c.FileAttachment("/path/to/file", "filename.txt")

// 流式响应
c.Stream(func(w io.Writer) bool {
    w.Write([]byte("data"))
    return true  // 继续发送
})
```

## 中间件

### 内置中间件

```go
// Logger - 日志
r.Use(gin.Logger())

// Recovery - panic 恢复
r.Use(gin.Recovery())

// 静态文件
r.Static("/static", "./public")
r.StaticFile("/favicon.ico", "./public/favicon.ico")
```

### 自定义中间件

```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()  // 处理请求
        
        latency := time.Since(start)
        status := c.Writer.Status()
        
        log.Printf("%s %s %d %v", c.Request.Method, path, status, latency)
    }
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        // 验证 token...
        user := validateToken(token)
        c.Set("user", user)  // 设置上下文
        
        c.Next()
    }
}

// 使用
r.Use(LoggerMiddleware())
r.GET("/protected", AuthMiddleware(), protectedHandler)
```

### 从上下文获取数据

```go
func protectedHandler(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(500, gin.H{"error": "user not found"})
        return
    }
    c.JSON(200, gin.H{"user": user})
}
```

## ⚠️ 常见陷阱

### 陷阱 1：goroutine 中使用 Context

```go
// ❌ 在 goroutine 中直接使用 c
func handler(c *gin.Context) {
    go func() {
        time.Sleep(time.Second)
        c.JSON(200, gin.H{})  // 危险！c 可能已失效
    }()
}

// ✅ 复制需要的数据
func handler(c *gin.Context) {
    cCopy := c.Copy()  // 或者只复制需要的值
    go func() {
        time.Sleep(time.Second)
        processAsync(cCopy)
    }()
    c.JSON(200, gin.H{"status": "processing"})
}
```

### 陷阱 2：多次写入响应

```go
// ❌ 多次写入
func handler(c *gin.Context) {
    c.JSON(200, gin.H{"first": true})
    c.JSON(200, gin.H{"second": true})  // 无效！
}

// ✅ 只写入一次
func handler(c *gin.Context) {
    c.JSON(200, gin.H{"result": "done"})
    return
}
```

### 陷阱 3：Next 之后的代码

```go
// 中间件中 Next 之后的代码在响应后执行
func middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求前
        c.Next()
        // 响应后（可以用于记录日志、清理等）
    }
}
```

## 错误处理

```go
// 全局错误处理
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            c.JSON(-1, gin.H{
                "errors": c.Errors.Errors(),
            })
        }
    }
}

// 在 handler 中记录错误
func handler(c *gin.Context) {
    if err := doSomething(); err != nil {
        c.Error(err)
        c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
        return
    }
}
```

## 优雅关闭

```go
func main() {
    r := gin.Default()
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }
    
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## 参考资源

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [Go 微服务访问控制之 Casbin 实践指南](https://www.fanyamin.com/journal/2025-07-13-go-casbin-wei-fu-wu-fang-wen-kong-zhi-zhi-shi-jian-zhi-nan.html) — Gin + Casbin + JWT 实现 RBAC
