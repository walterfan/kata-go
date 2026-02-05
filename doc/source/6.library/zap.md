# Zap 日志库

```{contents} 目录
:depth: 3
```

## Zap 概述

Zap 是 Uber 开发的高性能、结构化日志库，特点：

- **极快**：零内存分配的日志记录
- **结构化**：支持结构化日志
- **灵活**：可定制的编码器和输出

## 安装

```bash
go get go.uber.org/zap
```

## 基本用法

### 快速开始

```go
package main

import "go.uber.org/zap"

func main() {
    // 开发模式：可读性好，性能略低
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()
    
    logger.Info("This is an info message",
        zap.String("key", "value"),
        zap.Int("count", 42),
    )
    
    // 生产模式：JSON 格式，性能最优
    logger, _ = zap.NewProduction()
    defer logger.Sync()
    
    logger.Info("Production log",
        zap.String("service", "myapp"),
    )
}
```

### SugaredLogger（更方便）

```go
logger, _ := zap.NewProduction()
sugar := logger.Sugar()
defer sugar.Sync()

// printf 风格
sugar.Infof("Processing request %s", requestID)

// 键值对风格
sugar.Infow("User logged in",
    "username", "alice",
    "ip", "192.168.1.1",
)
```

## 日志级别

```go
logger.Debug("Debug message")   // -1
logger.Info("Info message")     //  0
logger.Warn("Warning message")  //  1
logger.Error("Error message")   //  2
logger.DPanic("DPanic message") //  3 (开发模式会 panic)
logger.Panic("Panic message")   //  4 (会 panic)
logger.Fatal("Fatal message")   //  5 (会调用 os.Exit(1))
```

## 自定义配置

```go
cfg := zap.Config{
    Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
    Development: false,
    Encoding:    "json",  // 或 "console"
    EncoderConfig: zapcore.EncoderConfig{
        TimeKey:        "timestamp",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        FunctionKey:    zapcore.OmitKey,
        MessageKey:     "message",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    },
    OutputPaths:      []string{"stdout", "/var/log/myapp.log"},
    ErrorOutputPaths: []string{"stderr"},
}

logger, _ := cfg.Build()
```

## 结构化字段

```go
// 强类型字段（性能最好）
logger.Info("Request processed",
    zap.String("method", "GET"),
    zap.String("path", "/api/users"),
    zap.Int("status", 200),
    zap.Duration("latency", time.Millisecond*150),
    zap.Time("timestamp", time.Now()),
    zap.Any("headers", headers),
    zap.Error(err),
)

// 对象字段
type User struct {
    ID   int
    Name string
}

func (u User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddInt("id", u.ID)
    enc.AddString("name", u.Name)
    return nil
}

logger.Info("User created", zap.Object("user", User{ID: 1, Name: "Alice"}))
```

## 日志轮转

Zap 本身不支持日志轮转，需要配合 lumberjack：

```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

func getLogWriter() zapcore.WriteSyncer {
    lumberJackLogger := &lumberjack.Logger{
        Filename:   "/var/log/myapp.log",
        MaxSize:    100,    // MB
        MaxBackups: 5,
        MaxAge:     30,     // days
        Compress:   true,
    }
    return zapcore.AddSync(lumberJackLogger)
}

func main() {
    writeSyncer := getLogWriter()
    encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
    
    core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
    logger := zap.New(core, zap.AddCaller())
    
    logger.Info("Using lumberjack for log rotation")
}
```

## 全局 Logger

```go
// 设置全局 logger
logger, _ := zap.NewProduction()
zap.ReplaceGlobals(logger)

// 使用全局 logger
zap.L().Info("Using global logger")
zap.S().Infof("Using global sugar logger: %s", message)
```

## 添加上下文

```go
// 创建带有固定字段的 logger
logger := zap.NewExample()
contextLogger := logger.With(
    zap.String("service", "user-service"),
    zap.String("version", "1.0.0"),
)

contextLogger.Info("Request received") // 自动包含 service 和 version
```

## ⚠️ 常见陷阱

### 陷阱 1：忘记 Sync

```go
// ❌ 日志可能丢失
func main() {
    logger, _ := zap.NewProduction()
    logger.Info("Hello")
    // 程序退出，缓冲区未刷新
}

// ✅ 确保 Sync
func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()  // 刷新缓冲区
    logger.Info("Hello")
}
```

### 陷阱 2：性能敏感场景使用 Sugar

```go
// ❌ Sugar 稍慢（但仍然很快）
sugar.Infof("User %s logged in", username)

// ✅ 极致性能使用强类型
logger.Info("User logged in", zap.String("username", username))
```

### 陷阱 3：错误处理

```go
// ❌ Error 字段传 nil
logger.Error("Failed", zap.Error(nil))  // 输出 error: null

// ✅ 检查 error
if err != nil {
    logger.Error("Failed", zap.Error(err))
}
```

## 与 Gin 集成

```go
import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func GinZapLogger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        logger.Info("Request",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", time.Since(start)),
            zap.String("client_ip", c.ClientIP()),
        )
    }
}

func main() {
    logger, _ := zap.NewProduction()
    
    r := gin.New()
    r.Use(GinZapLogger(logger))
    // ...
}
```

## 参考资源

- [Zap GitHub](https://github.com/uber-go/zap)
- [Zap Documentation](https://pkg.go.dev/go.uber.org/zap)
