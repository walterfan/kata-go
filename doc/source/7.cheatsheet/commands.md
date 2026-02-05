# Go 命令速查表

```{contents} 目录
:depth: 2
```

```{seealso}
延伸阅读：`Debug Build 的两种哲学：C++ 宏 vs Go 链接器注入 <https://www.fanyamin.com/journal/2026-02-03-debug_build_cpp_vs_go.html>`_ — `-ldflags -X` 与 Build Tags 详解。
```

## 基本命令

```bash
# 运行
go run main.go
go run .

# 编译
go build
go build -o myapp
go build -ldflags "-X main.Version=1.0.0"

# 安装
go install

# 获取依赖
go get github.com/pkg/errors
go get -u github.com/pkg/errors  # 更新

# 清理
go clean
go clean -cache      # 清理构建缓存
go clean -testcache  # 清理测试缓存
```

## 模块管理

```bash
# 初始化模块
go mod init github.com/user/project

# 整理依赖
go mod tidy

# 下载依赖
go mod download

# 验证依赖
go mod verify

# 查看依赖
go list -m all
go list -m -versions github.com/pkg/errors

# 依赖图
go mod graph

# 编辑 go.mod
go mod edit -require github.com/pkg/errors@v0.9.1
go mod edit -replace old=new
```

## 测试

```bash
# 运行测试
go test
go test ./...           # 所有包
go test -v              # 详细输出
go test -run TestName   # 运行特定测试
go test -count=1        # 禁用缓存

# 覆盖率
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# 基准测试
go test -bench=.
go test -bench=. -benchmem
go test -bench=. -benchtime=5s
go test -bench=. -count=5

# 竞态检测
go test -race
```

## 性能分析

```bash
# CPU 分析
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# 内存分析
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# HTTP pprof
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof http://localhost:6060/debug/pprof/heap

# pprof 交互
(pprof) top
(pprof) web
(pprof) list funcName

# trace
go test -trace=trace.out
go tool trace trace.out
```

## 代码质量

```bash
# 格式化
go fmt ./...
gofmt -w .

# 静态检查
go vet ./...

# 文档
go doc fmt
go doc fmt.Println

# 生成
go generate ./...
```

## 交叉编译

```bash
# Linux
GOOS=linux GOARCH=amd64 go build

# Windows
GOOS=windows GOARCH=amd64 go build

# macOS ARM
GOOS=darwin GOARCH=arm64 go build

# 查看支持的平台
go tool dist list
```

## 环境变量

```bash
# 查看环境
go env
go env GOPATH
go env GOPROXY

# 设置环境
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on

# 常用变量
GOPATH      # 工作空间
GOROOT      # Go 安装目录
GOPROXY     # 代理
GONOPROXY   # 不使用代理
GOPRIVATE   # 私有模块
GOFLAGS     # 默认标志
```

## 常用工具

```bash
# 安装工具
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install go.uber.org/mock/mockgen@latest

# goimports
goimports -w .

# golangci-lint
golangci-lint run

# mockgen
mockgen -source=interface.go -destination=mock.go
```

## 调试

```bash
# 安装 delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试
dlv debug
dlv debug -- arg1 arg2
dlv test
dlv attach <pid>

# delve 命令
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) step
(dlv) print varName
(dlv) goroutines
```
