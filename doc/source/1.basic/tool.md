# Go Tools

## 🚀 **开发 & 代码质量**
### 1️⃣ 代码格式化
- **[`gofmt`](https://golang.org/cmd/gofmt/)**  
  - Go 官方格式化工具，自动整理代码风格  
  - 使用方式：
    ```sh
    gofmt -w .
    ```

- **[`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)**  
  - 自动格式化代码并整理 `import` 语句  
  - 替换 `gofmt`，更智能  

---

### 2️⃣ 代码检查 & Lint
- **[`golangci-lint`](https://golangci-lint.run/)**  
  - 集成多种 lint 工具的强大 linter  
  - 使用方式：
    ```sh
    golangci-lint run ./...
    ```

- **[`staticcheck`](https://staticcheck.dev/)**  
  - 高级静态分析工具，可检测错误、性能问题和代码风格问题  

---

## 🐞 **调试 & 性能分析**
### 3️⃣ 调试工具
- **[`delve (dlv)`](https://github.com/go-delve/delve)**  
  - Go 官方推荐的调试器，支持断点、单步执行、变量查看  
  - 使用方式：
    ```sh
    dlv debug main.go
    ```

---

### 4️⃣ 性能分析 & 追踪
- **[`pprof`](https://pkg.go.dev/net/http/pprof)**  
  - Go 内置性能分析工具，可以分析 CPU、内存等
  - 使用方式：
    ```sh
    go tool pprof http://localhost:6060/debug/pprof/profile
    ```
  - 可以结合 **[`go-torch`](https://github.com/uber/go-torch)** 生成火焰图

- **[`trace`](https://pkg.go.dev/runtime/trace)**  
  - Go 官方提供的运行时跟踪工具，用于分析 goroutine 执行情况
  - 使用方式：
    ```sh
    go run main.go
    go tool trace trace.out
    ```

---

## 📦 **依赖管理**
### 5️⃣ 依赖管理工具
- **[`Go Modules`](https://golang.org/ref/mod)** （官方推荐）  
  - 现代化的 Go 依赖管理工具  
  - 使用方式：
    ```sh
    go mod init my_project
    go mod tidy
    ```

- **[`goproxy`](https://goproxy.cn/)**  
  - 国内推荐使用 `https://goproxy.cn` 加速 Go 模块下载  
  - 配置方式：
    ```sh
    go env -w GOPROXY=https://goproxy.cn,direct
    ```

---

## 🛠 **开发辅助**
### 6️⃣ 代码生成 & API 工具
- **[`swaggo`](https://github.com/swaggo/swag)**  
  - 用于自动生成 **Swagger API 文档**，适用于 RESTful 服务  
  - 使用方式：
    ```sh
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init
    ```

- **[`ent`](https://entgo.io/)**  
  - 强大的 Go ORM 框架，适用于数据库操作  
  - 使用方式：
    ```sh
    go install entgo.io/ent/cmd/ent@latest
    ent init User
    ```

- **[`gomock`](https://github.com/golang/mock)**  
  - 单元测试 mock 工具，可用于模拟接口  

---

## ☁️ **云开发 & 部署**
### 7️⃣ DevOps & 部署
- **[`Air`](https://github.com/cosmtrek/air)**  
  - 热重载工具，支持代码改动后自动重启服务  
  - 使用方式：
    ```sh
    go install github.com/cosmtrek/air@latest
    air
    ```

- **[`goreleaser`](https://goreleaser.com/)**  
  - 自动化 Go 项目发布，生成二进制文件并发布到 GitHub Release  
  - 使用方式：
    ```sh
    go install github.com/goreleaser/goreleaser@latest
    goreleaser init
    ```

---

## 🎯 **推荐的 Go 开发环境**
### 8️⃣ IDE & 编辑器
- **[GoLand](https://www.jetbrains.com/go/)** （最强大的 Go 开发 IDE，支持调试、代码分析）
- **[VS Code](https://code.visualstudio.com/)** （轻量级，搭配 `gopls` 扩展）
  - 推荐插件：
    - `Go`（官方插件）
    - `golangci-lint`（代码检查）
    - `Go Test Explorer`（测试可视化）

---

## 🔥 **总结**
| **类别** | **工具** | **作用** |
|----------|---------|---------|
| **代码格式化** | `gofmt`, `goimports` | 代码自动格式化 |
| **代码检查** | `golangci-lint`, `staticcheck` | 代码规范与静态分析 |
| **调试工具** | `dlv` (delve) | 断点调试 |
| **性能分析** | `pprof`, `trace` | 监控 CPU、内存、goroutine |
| **依赖管理** | `Go Modules`, `goproxy` | 依赖管理 |
| **开发辅助** | `swaggo`, `ent`, `gomock` | API 文档、ORM、单元测试 |
| **部署工具** | `Air`, `goreleaser` | 热重载、自动化发布 |
| **IDE** | `GoLand`, `VS Code` | 开发环境 |
