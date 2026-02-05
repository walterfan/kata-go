# Best Practices

```{seealso}
延伸阅读：`SoC Code Structure in Golang <https://www.fanyamin.com/journal/2025-08-25-soc-code-structure-in-golang.html>`_、`Go 应用程序的代码组织 <https://www.fanyamin.com/journal/2025-08-29-go-ying-yong-cheng-xu-de-dai-ma-zu-zhi.html>`_ — 项目结构与 MVC、依赖注入最佳实践。
```

## Best Practices for Go Backend Development
### 1. **Project Structure**
Organize your project in a way that makes it easy to navigate and maintain. A common structure for a Go backend service is:

```
/my-service
  ├── /cmd
  │   └── /my-service
  │       └── main.go
  ├── /internal
  │   ├── /handlers
  │   ├── /models
  │   ├── /services
  │   └── /repositories
  ├── /pkg
  │   └── /utils
  ├── /configs
  ├── /migrations
  ├── /api
  ├── /scripts
  └── go.mod
```

- **`cmd/my-service/main.go`**: Entry point of the application.
- **`internal/`**: Contains the core application logic. This directory is private to your module.
- **`pkg/`**: Reusable utility functions or libraries.
- **`configs/`**: Configuration files (e.g., YAML, JSON).
- **`migrations/`**: Database migration scripts.
- **`api/`**: API definitions (e.g., OpenAPI/Swagger specs).
- **`scripts/`**: Helper scripts for deployment, testing, etc.

---

### 2. **Error Handling**
Go encourages explicit error handling. Always check for errors and handle them gracefully.

```go
func someFunction() error {
    result, err := doSomething()
    if err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    // Use result
    return nil
}
```

- Use `fmt.Errorf` with `%w` to wrap errors for better context.
- Avoid panics in production code; handle errors instead.

---

### 3. **Concurrency**
Go's goroutines and channels make concurrency easy, but misuse can lead to bugs or resource leaks.

- Use `sync.WaitGroup` to wait for goroutines to finish:
  ```go
  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
      defer wg.Done()
      // Do work
  }()
  wg.Wait()
  ```

- Use channels for communication between goroutines:
  ```go
  ch := make(chan int)
  go func() {
      ch <- 42
  }()
  value := <-ch
  ```

- Avoid leaking goroutines by ensuring they always exit.

---

### 4. **Network Communication**
For network communication, Go's `net/http` package is commonly used. Here are some best practices:

#### a. **HTTP Server**
- Use `http.HandlerFunc` or `http.Handler` for routing:
  ```go
  http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusOK)
      w.Write([]byte("OK"))
  })
  ```

- Use middleware for cross-cutting concerns (e.g., logging, authentication):
  ```go
  func loggingMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          log.Println(r.Method, r.URL.Path)
          next.ServeHTTP(w, r)
      })
  }
  ```

- Use `context` for request-scoped values and cancellation:
  ```go
  func handler(w http.ResponseWriter, r *http.Request) {
      ctx := r.Context()
      // Pass ctx to downstream functions
  }
  ```

#### b. **HTTP Client**
- Use `http.Client` for making HTTP requests:
  ```go
  client := &http.Client{}
  req, err := http.NewRequest("GET", "https://example.com", nil)
  if err != nil {
      log.Fatal(err)
  }
  resp, err := client.Do(req)
  if err != nil {
      log.Fatal(err)
  }
  defer resp.Body.Close()
  ```

- Always close the response body to avoid resource leaks.

#### c. **gRPC**
- Use `grpc-go` for high-performance RPC communication:
  ```go
  conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
  if err != nil {
      log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  client := pb.NewMyServiceClient(conn)
  ```

---

### 5. **Configuration Management**
Use environment variables or configuration files to manage settings.

- Use `os.Getenv` for environment variables:
  ```go
  port := os.Getenv("PORT")
  if port == "" {
      port = "8080"
  }
  ```

- Use libraries like `viper` for advanced configuration management:
  ```go
  viper.SetConfigFile("config.yaml")
  err := viper.ReadInConfig()
  if err != nil {
      log.Fatalf("failed to read config: %v", err)
  }
  port := viper.GetString("port")
  ```

---

### 6. **Logging**
Use structured logging for better observability.

- Use `log` or `logrus` for structured logging:
  ```go
  log.WithFields(log.Fields{
      "event": "request",
      "method": r.Method,
      "path":   r.URL.Path,
  }).Info("Handling request")
  ```

---

### 7. **Testing**
Write unit tests and integration tests for your code.

- Use `testing` package for unit tests:
  ```go
  func TestAdd(t *testing.T) {
      result := Add(2, 3)
      if result != 5 {
          t.Errorf("Expected 5, got %d", result)
      }
  }
  ```

- Use `httptest` for testing HTTP handlers:
  ```go
  req := httptest.NewRequest("GET", "/health", nil)
  w := httptest.NewRecorder()
  handler(w, req)
  if w.Code != http.StatusOK {
      t.Errorf("Expected status 200, got %d", w.Code)
  }
  ```

---

### 8. **Dependency Management**
Use Go modules for dependency management.

- Initialize a module:
  ```bash
  go mod init my-service
  ```

- Add dependencies:
  ```bash
  go get github.com/some/package
  ```

- Use `go mod tidy` to clean up unused dependencies.

---

### 9. **Security**
- Validate and sanitize all inputs to prevent injection attacks.
- Use HTTPS for secure communication.
- Avoid hardcoding sensitive information (e.g., API keys, passwords).

---

### 10. **Documentation**
- Use Go's built-in documentation tool (`godoc`) to document your code.
- Write clear and concise comments for exported functions and types.
