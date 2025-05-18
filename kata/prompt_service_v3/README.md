# Prompt Management Service

## example

This Go project is a **prompt management service**, providing APIs to manage prompts used for interacting with LLMs (Large Language Models). It includes core functionalities such as CRUD operations and search, backed by a SQLite database and exposed metrics for monitoring.

---

### 🧩 主要功能模块

| 模块 | 功能说明 |
|------|----------|
| **main.go** | 程序入口，使用 `cobra` 支持命令行参数（如监听端口），集成 `zap` 日志系统。 |
| **pkg/database/sqlite.go** | 初始化 SQLite 数据库连接，并提供数据迁移和初始化样本数据的功能。 |
| **pkg/models/prompt.go** | 定义 [Prompt](pkg/models/prompt.go#L4-L14) 结构体，映射数据库表结构，包含字段如 [Name](pkg/models/prompt.go#L6-L6), [Description](pkg/models/prompt.go#L7-L7), [Tags](pkg/models/prompt.go#L10-L10), [UserPrompt](pkg/models/prompt.go#L9-L9), [SystemPrompt](pkg/models/prompt.go#L8-L8) 等。 |
| **pkg/handlers/prompt_handler.go** | 提供 RESTful API 接口：  
- `POST /api/v1/prompts`: 创建 Prompt  
- `GET /api/v1/prompts/:id`: 获取单个 Prompt  
- `PUT /api/v1/prompts/:id`: 更新 Prompt  
- `DELETE /api/v1/prompts/:id`: 删除 Prompt  
- `GET /api/v1/prompts`: 支持关键字搜索与分页 |
| **pkg/metrics/metrics.go** | 集成 Prometheus 指标监控，记录 HTTP 请求次数、耗时等信息。 |

---

### 📦 技术栈

| 技术 | 用途 |
|------|------|
| **Gin** | Web 框架，用于构建 HTTP 服务。 |
| **GORM + SQLite** | ORM 和数据库，用于持久化存储 prompts 数据。 |
| **Prometheus + Metrics Middleware** | 监控接口调用次数、延迟等运行指标。 |
| **Zap** | 高性能日志库，用于记录服务日志。 |
| **Cobra** | CLI 命令行支持，用于解析启动参数（如监听端口）。 |

---

### 📈 特性

- ✅ **RESTful API 设计**：清晰的接口设计，易于集成到前端或 AI 应用中。
- ✅ **Prometheus 指标暴露**：通过 `/metrics` 接口可接入监控系统。
- ✅ **SQLite 轻量级存储**：适合快速部署和开发测试环境。
- ✅ **CLI 参数支持**：可通过命令行配置服务监听端口。
- ✅ **日志统一管理**：使用 `zap` 提升日志性能与结构化输出能力。
- ✅ **分页搜索功能**：支持关键字模糊匹配和分页查询。

---

### 🔄 示例使用场景

- 存储常见提示语模板，供不同 LLM 使用
- 快速检索特定任务相关的 prompt（例如 Golang 编程技巧）
- 统计 prompt 使用频率和更新时间
- 监控服务健康状态和接口响应性能

## 测试命令

请先启动服务：

```bash
go run main.go --port 8080
```
---

### 🟢 创建 (Create) - `POST /api/v1/prompts`

```bash
curl -X POST http://localhost:8080/api/v1/prompts \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Explain Goroutines",
        "description": "Describe how goroutines work in Go.",
        "systemPrompt": "You are a Go language expert.",
        "userPrompt": "What is a goroutine and how does it differ from a thread?",
        "tags": "concurrency,golang"
      }'
```

---

### 🔵 查询单个 (Read) - `GET /api/v1/prompts/:id`

```bash
curl http://localhost:8080/api/v1/prompts/1
```

---

### 🟡 更新 (Update) - `PUT /api/v1/prompts/:id`

```bash
curl -X PUT http://localhost:8080/api/v1/prompts/1 \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Updated Goroutines Prompt",
        "description": "Updated description for goroutines."
      }'
```

---

### 🔴 删除 (Delete) - `DELETE /api/v1/prompts/:id`

```bash
curl -X DELETE http://localhost:8080/api/v1/prompts/1
```

---

### 🔍 查询列表 (Search with Pagination) - `GET /api/v1/prompts`

- **带关键字搜索和分页**

```bash
curl "http://localhost:8080/api/v1/prompts?q=golang&pageNum=1&pageSize=20"
```

- **不带参数查询所有**

```bash
curl http://localhost:8080/api/v1/prompts
```

---

### ✅ 示例输出字段说明

返回的 JSON 数据结构如下：

```json
{
  "id": 1,
  "name": "Explain Goroutines",
  "desc": "Describe how goroutines work in Go.",
  "systemPrompt": "You are a Go language expert.",
  "userPrompt": "What is a goroutine and how does it differ from a thread?",
  "tags": "concurrency,golang",
  "createdAt": 1715000000,
  "updatedAt": 1715000000
}
```

---

### Login

```

curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
        "username": "admin",
        "password": "defaultpassword"
      }'
```