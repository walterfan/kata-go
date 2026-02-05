# HTTP 编程

```{contents} 目录
:depth: 3
```

## net/http 包

Go 标准库的 `net/http` 包功能强大，可用于生产环境。

## HTTP 服务端

### 基本用法

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### 使用自定义 ServeMux

```go
func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/api/users", usersHandler)
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    server.ListenAndServe()
}
```

### ⚠️ 常见陷阱

#### 陷阱 1：忽略 http.Server 配置

```go
// ❌ 无超时配置，可能被慢客户端攻击
http.ListenAndServe(":8080", nil)

// ✅ 配置超时
server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
server.ListenAndServe()
```

#### 陷阱 2：未检查 ResponseWriter.Write 错误

```go
// ❌ 忽略错误
func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("data"))  // 忽略了返回的错误
}

// ✅ 检查错误
func handler(w http.ResponseWriter, r *http.Request) {
    _, err := w.Write([]byte("data"))
    if err != nil {
        log.Printf("write error: %v", err)
        return
    }
}
```

#### 陷阱 3：在写入 Header 后修改

```go
// ❌ Header 已发送后无法修改
func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("data"))
    w.Header().Set("X-Custom", "value")  // 无效！
}

// ✅ 先设置 Header
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("X-Custom", "value")
    w.Write([]byte("data"))
}
```

## HTTP 客户端

### 基本用法

```go
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
```

### ⚠️ 客户端陷阱

#### 陷阱 1：使用默认 Client

```go
// ❌ 默认 Client 无超时
resp, _ := http.Get(url)

// ✅ 配置超时
client := &http.Client{
    Timeout: 10 * time.Second,
}
resp, _ := client.Get(url)
```

#### 陷阱 2：忘记关闭 Response Body

```go
// ❌ 泄漏连接
resp, _ := http.Get(url)
// 忘记 resp.Body.Close()

// ✅ 使用 defer 关闭
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
```

#### 陷阱 3：不读取 Response Body

```go
// ❌ 连接无法复用
resp, _ := http.Get(url)
defer resp.Body.Close()
// 未读取 body，连接无法复用

// ✅ 读取并丢弃 body
resp, _ := http.Get(url)
defer resp.Body.Close()
io.Copy(io.Discard, resp.Body)  // 确保读完
```

### 自定义 Transport

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
}

client := &http.Client{
    Transport: transport,
    Timeout:   30 * time.Second,
}
```

## 中间件模式

```go
type Middleware func(http.Handler) http.Handler

// 日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// 链式使用
func main() {
    handler := LoggingMiddleware(AuthMiddleware(http.HandlerFunc(myHandler)))
    http.ListenAndServe(":8080", handler)
}
```

## 优雅关闭

```go
func main() {
    server := &http.Server{Addr: ":8080", Handler: mux}
    
    // 启动服务器
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Server stopped")
}
```

## HTTP/2 支持

```go
// Go 默认支持 HTTP/2（HTTPS 时）
server := &http.Server{
    Addr:    ":443",
    Handler: mux,
}
server.ListenAndServeTLS("cert.pem", "key.pem")
```

## 参考资源

- [net/http Package](https://pkg.go.dev/net/http)
- [Writing Web Applications](https://go.dev/doc/articles/wiki/)
