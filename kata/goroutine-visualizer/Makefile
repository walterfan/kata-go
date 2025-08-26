# Makefile for goroutine-visualizer

# 变量定义
APP_NAME := goroutine-visualizer
GO_VERSION := 1.21
PORT := 8080

# 默认目标
.PHONY: help
help:
	@echo "可用的命令:"
	@echo "  make run      - 运行应用程序"
	@echo "  make build    - 构建应用程序"
	@echo "  make clean    - 清理构建文件"
	@echo "  make deps     - 安装依赖"
	@echo "  make test     - 运行测试"
	@echo "  make fmt      - 格式化代码"
	@echo "  make lint     - 代码检查"
	@echo "  make dev      - 开发模式运行"
	@echo "  make docker   - 构建Docker镜像"
	@echo "  make help     - 显示帮助信息"

# 安装依赖
.PHONY: deps
deps:
	@echo "正在安装依赖..."
	go mod tidy
	go mod download

# 运行应用程序
.PHONY: run
run: deps
	@echo "正在启动 $(APP_NAME)..."
	@echo "访问地址: http://localhost:$(PORT)"
	go run main.go

# 构建应用程序
.PHONY: build
build: deps
	@echo "正在构建 $(APP_NAME)..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/$(APP_NAME) main.go
	@echo "构建完成: bin/$(APP_NAME)"

# 交叉编译
.PHONY: build-all
build-all: deps
	@echo "正在进行交叉编译..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(APP_NAME)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(APP_NAME)-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(APP_NAME)-windows-amd64.exe main.go
	@echo "交叉编译完成"

# 清理构建文件
.PHONY: clean
clean:
	@echo "正在清理构建文件..."
	rm -rf bin/
	go clean
	@echo "清理完成"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "正在格式化代码..."
	go fmt ./...
	@echo "格式化完成"

# 代码检查
.PHONY: lint
lint:
	@echo "正在进行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，使用 go vet 代替"; \
		go vet ./...; \
	fi

# 运行测试
.PHONY: test
test:
	@echo "正在运行测试..."
	go test -v ./...

# 开发模式运行（带实时重载）
.PHONY: dev
dev:
	@echo "开发模式启动..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air 未安装，使用普通模式运行"; \
		$(MAKE) run; \
	fi

# 安装开发工具
.PHONY: install-tools
install-tools:
	@echo "正在安装开发工具..."
	go install github.com/cosmtrek/air@latest
	@if command -v curl >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin; \
	else \
		echo "请手动安装 golangci-lint"; \
	fi

# 构建Docker镜像
.PHONY: docker
docker:
	@echo "正在构建Docker镜像..."
	docker build -t $(APP_NAME):latest .
	@echo "Docker镜像构建完成: $(APP_NAME):latest"

# 运行Docker容器
.PHONY: docker-run
docker-run: docker
	@echo "正在运行Docker容器..."
	docker run -p $(PORT):$(PORT) --name $(APP_NAME) $(APP_NAME):latest

# 停止Docker容器
.PHONY: docker-stop
docker-stop:
	@echo "正在停止Docker容器..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# 查看项目信息
.PHONY: info
info:
	@echo "项目信息:"
	@echo "  应用名称: $(APP_NAME)"
	@echo "  Go版本: $(GO_VERSION)"
	@echo "  端口: $(PORT)"
	@echo "  当前目录: $(shell pwd)"
	@echo "  Go版本: $(shell go version)"

# 生成版本信息
.PHONY: version
version:
	@echo "版本信息:"
	@echo "  Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "  Git branch: $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')"
	@echo "  Build time: $(shell date)"

# 检查依赖更新
.PHONY: check-updates
check-updates:
	@echo "检查依赖更新..."
	go list -u -m all

# 更新依赖
.PHONY: update-deps
update-deps:
	@echo "更新依赖..."
	go get -u ./...
	go mod tidy

# 性能测试
.PHONY: bench
bench:
	@echo "运行性能测试..."
	go test -bench=. -benchmem ./...

# 生成测试覆盖率报告
.PHONY: coverage
coverage:
	@echo "生成测试覆盖率报告..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告生成: coverage.html"

# 快速开始
.PHONY: quick-start
quick-start: deps
	@echo "快速开始..."
	@echo "1. 安装依赖完成"
	@echo "2. 正在启动服务器..."
	@echo "3. 请访问: http://localhost:$(PORT)"
	@echo "4. 按 Ctrl+C 停止服务器"
	@echo ""
	$(MAKE) run 