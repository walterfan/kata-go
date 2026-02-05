# kata-go

Kata should be called "routine" in Chinese. The secret of practicing martial arts is to master various routines.

First, you must learn from the strengths of hundreds of schools and be familiar with the routines of various schools in the world. Only then can you integrate them and achieve the state of no moves being better than moves.

Dave Thomas - the author of "The Pragmatic Programmer", proposed the idea of ​​Code Kata. Dave also collected a small practice project on his website (http://codekata.com/).

As a professional programmer, I hope to practice some routines that can be often used in work, such as some small routines for file modification, image cutting, and network sending and receiving, so I will organize and collect some routines here.


## Kata

1. [cron-service.go](./kata/cron): demostration of how to implement cron job
1. [encoding_tool](./kata/encoding_tool): demostrate how to implement encoding command line tool by cobra library
1. [ds_web_console](./kata/ds_web_console/): demostrate how to call deep seek api and provide a web console by gin library
1. [ata-auth](./kata/kata-auth/): demostrate how to implement a simple auth system by casbin
1. [list_files.go](./kata/kata-files/list_files.go): demostrate how to list files in a directory
1. [links.go](./kata/kata-http/links.go): demostrate how to use http client to get links from a web page
1. [llm-agent-go](./kata/llm-agent-go/): demostrate how to implement a simple llm agent by gin and vue.js
1. [prompt_service](./kata/prompt_service): demostrate how to implement a simple prompt management service by gin and gorm with sqlite db
2. [prompt_service_v2](./kata/prompt_service_v2): add login and metrics endpoint for prompt management service
3. [service-monitor](./kata/service_monitor): demostrate how to monitor a backend service with prometheus exportor
4. [simple-ai-agent](./kata/simple-ai-agent): demostrate how to use function calling and tools of openai
5. [unix_socket](./kata/unix_socket): demostrate how to use unix socket to communicate with backend service


## Golang 实战开发指南

本项目包含一份完整的 Go 语言开发文档，专注于**易错**和**易忽略**的知识点。

### 1. 基础语法与陷阱
- [Go 概述](doc/source/1.basic/overview.md) - Go 语言核心概念
- [常见陷阱](doc/source/1.basic/trap.md) - Go 编程中的常见错误
- [开发工具](doc/source/1.basic/tool.md) - Go 开发工具链
- [最佳实践](doc/source/1.basic/best_practice.md) - Go 编程最佳实践

### 2. 并发编程
- [Goroutine 深入理解](doc/source/2.concurrency/goroutine.md) - GMP 模型、生命周期、泄漏检测
- [Channel 详解](doc/source/2.concurrency/channel.md) - Channel 类型、select、常用模式
- [sync 包详解](doc/source/2.concurrency/sync.md) - Mutex、RWMutex、WaitGroup、Pool
- [并发模式](doc/source/2.concurrency/patterns.md) - Worker Pool、Rate Limiter、Circuit Breaker
- [并发陷阱](doc/source/2.concurrency/pitfalls.md) - 数据竞争、死锁、泄漏

### 3. 性能调优
- [性能分析 (Profiling)](doc/source/3.performance/profiling.md) - pprof 使用指南
- [基准测试 (Benchmark)](doc/source/3.performance/benchmark.md) - 编写和分析基准测试
- [逃逸分析](doc/source/3.performance/escape_analysis.md) - 栈 vs 堆分配
- [性能优化技巧](doc/source/3.performance/optimization.md) - 字符串、切片、Map、并发优化

### 4. 内存管理
- [垃圾回收 (GC)](doc/source/4.memory/gc.md) - GC 原理、GOGC、GOMEMLIMIT
- [内存分配](doc/source/4.memory/allocation.md) - 分配器原理、内存对齐
- [内存泄漏排查](doc/source/4.memory/leak.md) - 常见泄漏场景与检测方法

### 5. 网络编程
- [HTTP 编程](doc/source/5.network/http.md) - net/http 服务端与客户端
- [gRPC](doc/source/5.network/grpc.md) - Protocol Buffers、服务定义、拦截器
- [TCP/UDP 编程](doc/source/5.network/tcp_udp.md) - 底层网络编程

### 6. 常用库
- [Viper 配置管理](doc/source/6.library/viper.md) - 配置文件、环境变量、热重载
- [Cobra 命令行框架](doc/source/6.library/cobra.md) - CLI 应用开发
- [Zap 日志库](doc/source/6.library/zap.md) - 高性能结构化日志
- [GORM ORM 框架](doc/source/6.library/gorm.md) - 数据库操作
- [Gin Web 框架](doc/source/6.library/gin.md) - HTTP 服务开发

### 7. 速查表
- [Go 语法速查表](doc/source/7.cheatsheet/syntax.md) - 语法快速参考
- [Go 命令速查表](doc/source/7.cheatsheet/commands.md) - 常用命令
- [Map 操作](doc/source/7.cheatsheet/map.md) - Map 使用技巧

### 延伸阅读（博客文章）

以下博客文章与文档主题相关，可作为补充：

| 主题 | 文章 |
|------|------|
| 常见陷阱 | [Go 语言的常见陷阱](https://www.fanyamin.com/journal/2025-03-25-go-yu-yan-de-chang-jian-xian-jing.html) |
| 并发哲学 | [通过通信来共享内存](https://www.fanyamin.com/journal/2025-03-26-tong-guo-tong-xin-lai-gong-xiang-nei-cun-er-bu-shi-tong-guo.html) |
| 访问控制 | [Go Casbin 实践指南](https://www.fanyamin.com/journal/2025-07-13-go-casbin-wei-fu-wu-fang-wen-kong-zhi-zhi-shi-jian-zhi-nan.html) |
| 代码结构 | [SoC Code Structure](https://www.fanyamin.com/journal/2025-08-25-soc-code-structure-in-golang.html)、[Go 代码组织](https://www.fanyamin.com/journal/2025-08-29-go-ying-yong-cheng-xu-de-dai-ma-zu-zhi.html) |
| Context | [Context in Go](https://www.fanyamin.com/journal/2025-08-28-context-in-go.html) |
| Goroutine 泄漏 | [Goroutine Leak 详解](https://www.fanyamin.com/journal/2025-12-13-go-goroutine-leak-jing-ti-ni-de-cheng-xu-zheng-zai-tou-tou-x.html) |
| 崩溃分析 | [Go 崩溃分析实战](https://www.fanyamin.com/journal/2026-01-23-golang_crash_analysis.html) |
| Debug Build | [C++ 宏 vs Go 链接器注入](https://www.fanyamin.com/journal/2026-02-03-debug_build_cpp_vs_go.html) |

### 构建文档

```bash
cd doc
pip install -r requirements.txt
make html
make serve  # 访问 http://localhost:8000
```

## Cheat sheet of golang

* [Go cheatsheet 1](go-cheat-sheet.md)
* [Go cheatsheet 2 ](https://devhints.io/go)
* [Go cheatsheet 3](https://quickref.me/go.html)

## Go tutorial

* [Go tutorial](https://tour.golang.org/welcome/1)
* [Go by example](https://gobyexample.com/)
* [Go Style Guide](https://google.github.io/styleguide/go/guide)
* [Go Style Decisions](https://google.github.io/styleguide/go/decisions)
* [Go Style Best Practices](https://google.github.io/styleguide/go/best-practices)

## Go Tools
### build and run
```shell
go build xxx.go
go run xxx.go
```
### check dependency
go list, go get, go mod xxx

```
go mod init my_project
go mod tidy
```
### format code
go fmt, gofmt

### Debug
```bash
dlv debug main.go

```
### documentation
go doc, godoc

```shell
go install github.com/swaggo/swag/cmd/swag@latest
swag init

```
### unit test

```
# Run all tests with verbose output for entire project
go test -v ./...

# Run tests in specific package with verbose output
go test -v ./internal/cmd

# Run specific test function
go test -v ./internal/cmd -run Test_Metrics

# Run tests matching a pattern
go test -v ./internal/cmd -run "Test.*eks.*"

# Run tests with regex pattern
go test -v ./internal/cmd -run "^TestMonitor.*"

# Skip specific tests
go test -v ./internal/cmd -skip "TestFileSync"

# Skip multiple test patterns
go test -v ./internal/cmd -skip "TestFileSync|TestMetrics"

# Run tests but skip long-running ones
go test -v ./internal/cmd -short
```
### static analysis
go vet, golangci-lint

e.g.
  ```bash
  brew install golangci-lint
  golangci-lint run ./...
  ```
### performance profile
go tool pprof, go tool trace

```bash
go tool pprof http://localhost:6060/debug/pprof/profile

go run main.go
go tool trace trace.out
```
### upgrade
go tool fix
### bug report
go bug



## go with vscode
* install go extensions and [delve](https://github.com/go-delve/delve/blob/master/Documentation/installation/osx/install.md)
* configuration of go debug in vscode

```json
 {
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug file",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${file}"
        },
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}"
        }
    ]
}
```

## Reference
* [Go Official Site](https://go.dev/)
* [Go Documentation](https://pkg.go.dev/)
* [Go Playground](https://go.dev/play/)
* [Go Blog](https://blog.golang.org/)
* [Go Wiki](https://github.com/golang/go/wiki)
* [Effective Go](https://golang.org/doc/effective_go)
* [Go Memory Model](https://golang.org/ref/mem)
* [Go Standards](https://google.github.io/styleguide/go/guide)
* [Awesome go projects](https://github.com/avelino/awesome-go)
