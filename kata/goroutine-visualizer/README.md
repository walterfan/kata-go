# Goroutine 调度器可视化工具

这是一个使用 Go + WebSocket + Web 页面实现的 goroutine 调度器可视化工具，能够实时监控和展示 Go 程序中 goroutine 的运行状态。

## 功能特性

- 📊 **实时监控**: 实时显示 goroutine 数量、状态分布和系统信息
- 📈 **趋势图表**: 展示 goroutine 数量和内存使用的变化趋势
- 📋 **详细列表**: 列出所有 goroutine 的详细信息
- 🔌 **WebSocket 连接**: 使用 WebSocket 实现低延迟的实时数据传输
- 🎨 **现代化 UI**: 美观的响应式界面设计
- 📱 **移动端适配**: 支持移动设备访问

## 项目结构

```
goroutine-visualizer/
├── go.mod                    # Go 模块文件
├── main.go                   # 主程序入口
├── internal/
│   ├── monitor/
│   │   └── goroutine.go      # Goroutine 监控器
│   └── ws/
│       └── websocket.go      # WebSocket 处理器
├── web/
│   ├── index.html            # 前端页面
│   ├── style.css             # 样式文件
│   └── script.js             # JavaScript 脚本
├── Makefile                  # 构建脚本
└── README.md                 # 项目说明
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行程序

```bash
# 直接运行
go run main.go

# 或者使用 Makefile
make run
```

### 3. 访问界面

打开浏览器访问: http://localhost:8080

## 使用说明

### 系统概览
- **Goroutine 总数**: 当前系统中 goroutine 的总数量
- **CPU 核心数**: 系统 CPU 核心数
- **GOMAXPROCS**: Go 程序可以同时使用的 CPU 核心数
- **内存使用**: 当前内存使用情况

### 图表说明
- **Goroutine 数量趋势**: 显示 goroutine 数量随时间的变化
- **内存使用趋势**: 显示内存使用量随时间的变化
- **状态分布**: 饼图显示不同状态的 goroutine 数量分布

### Goroutine 状态
- **running**: 正在运行的 goroutine
- **runnable**: 可运行但等待调度的 goroutine
- **waiting**: 等待资源的 goroutine
- **blocked**: 被阻塞的 goroutine
- **dead**: 已结束的 goroutine

### 控制选项
- **自动滚动**: 自动滚动到最新的 goroutine 列表
- **最大显示数量**: 限制显示的 goroutine 数量以提高性能

## 技术栈

- **后端**: Go 1.21+
- **WebSocket**: gorilla/websocket
- **前端**: HTML5 + CSS3 + JavaScript
- **图表库**: Chart.js
- **监控**: Go runtime 包

## 模拟任务

程序内置了几种模拟任务来展示 goroutine 的行为：

1. **周期性任务**: 每 2 秒创建计算密集型任务
2. **批量处理**: 每 3 秒创建一批并发任务
3. **网络模拟**: 每秒模拟网络请求

## 开发说明

### 添加新的监控指标

1. 在 `internal/monitor/goroutine.go` 中的 `SystemInfo` 结构体添加新字段
2. 在 `collectSystemInfo` 方法中收集新数据
3. 在前端 `script.js` 中添加处理逻辑
4. 在 `index.html` 和 `style.css` 中添加展示元素

### 自定义模拟任务

在 `main.go` 的 `startSimulatedTasks` 函数中添加你的模拟任务逻辑。

## 性能优化

- WebSocket 连接使用缓冲通道避免阻塞
- 图表数据限制在 50 个数据点以内
- Goroutine 列表支持限制显示数量
- 使用 requestAnimationFrame 优化动画性能

## 故障排除

### WebSocket 连接失败
- 检查防火墙设置
- 确认端口 8080 未被占用
- 查看浏览器控制台错误信息

### 数据不更新
- 检查 WebSocket 连接状态
- 确认后端服务正常运行
- 刷新页面重新连接

### 性能问题
- 降低最大显示数量
- 关闭自动滚动
- 检查浏览器性能

## 贡献指南

1. Fork 此项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

此项目使用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 作者

- 您的名字 - 初始版本

## 致谢

- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket 支持
- [Chart.js](https://www.chartjs.org/) - 图表库
- [Go Runtime](https://golang.org/pkg/runtime/) - 运行时信息收集 