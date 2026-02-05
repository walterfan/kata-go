# Viper 配置管理

```{contents} 目录
:depth: 3
```

## Viper 概述

Viper 是 Go 应用程序的完整配置解决方案，支持：

- JSON、YAML、TOML、HCL 等配置文件
- 环境变量
- 命令行标志
- 远程配置系统（etcd、Consul）
- 配置热重载

## 安装

```bash
go get github.com/spf13/viper
```

## 基本用法

### 读取配置文件

```go
package main

import (
    "fmt"
    "github.com/spf13/viper"
)

func main() {
    viper.SetConfigName("config")     // 配置文件名（不带扩展名）
    viper.SetConfigType("yaml")       // 配置文件类型
    viper.AddConfigPath(".")          // 查找路径
    viper.AddConfigPath("./config")   // 可以添加多个路径
    
    if err := viper.ReadInConfig(); err != nil {
        panic(fmt.Errorf("fatal error config file: %w", err))
    }
    
    // 读取配置
    fmt.Println(viper.GetString("database.host"))
    fmt.Println(viper.GetInt("database.port"))
}
```

### 配置文件示例

```yaml
# config.yaml
app:
  name: myapp
  port: 8080
  debug: true

database:
  host: localhost
  port: 3306
  user: root
  password: secret
  dbname: mydb

redis:
  host: localhost
  port: 6379
```

## 配置优先级

Viper 按以下顺序（从高到低）读取配置：

1. 显式调用 `Set`
2. 命令行标志
3. 环境变量
4. 配置文件
5. 远程配置
6. 默认值

```go
// 设置默认值
viper.SetDefault("database.port", 3306)

// 显式设置（最高优先级）
viper.Set("database.host", "production-db")
```

## 环境变量

### 自动绑定环境变量

```go
// 自动绑定所有环境变量
viper.AutomaticEnv()

// 设置环境变量前缀
viper.SetEnvPrefix("MYAPP")  // 将读取 MYAPP_* 环境变量

// 将 . 替换为 _
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 现在可以通过 MYAPP_DATABASE_HOST 设置 database.host
```

### 手动绑定

```go
viper.BindEnv("database.host", "DB_HOST")
```

## 与 Cobra 集成

```go
import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "My application",
}

func init() {
    cobra.OnInitialize(initConfig)
    
    rootCmd.PersistentFlags().String("config", "", "config file")
    rootCmd.PersistentFlags().Int("port", 8080, "server port")
    
    viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
}

func initConfig() {
    if cfgFile := viper.GetString("config"); cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        viper.SetConfigName("config")
        viper.AddConfigPath(".")
    }
    
    viper.AutomaticEnv()
    viper.ReadInConfig()
}
```

## 反序列化到结构体

```go
type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Database DatabaseConfig `mapstructure:"database"`
}

type AppConfig struct {
    Name  string `mapstructure:"name"`
    Port  int    `mapstructure:"port"`
    Debug bool   `mapstructure:"debug"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
}

func loadConfig() (*Config, error) {
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    return &config, nil
}
```

## 配置热重载

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("Config file changed:", e.Name)
    // 重新加载配置
    reloadConfig()
})
```

## ⚠️ 常见陷阱

### 陷阱 1：并发访问不安全

```go
// ❌ 并发读写 viper 不安全
go func() {
    for {
        viper.Set("key", value)
    }
}()
go func() {
    for {
        _ = viper.GetString("key")
    }
}()

// ✅ 使用全局配置结构体
var (
    config     *Config
    configLock sync.RWMutex
)

func GetConfig() *Config {
    configLock.RLock()
    defer configLock.RUnlock()
    return config
}

func reloadConfig() {
    newConfig := &Config{}
    viper.Unmarshal(newConfig)
    
    configLock.Lock()
    config = newConfig
    configLock.Unlock()
}
```

### 陷阱 2：大小写敏感

```go
// Viper 键名不区分大小写
viper.Set("myKey", "value")
fmt.Println(viper.GetString("MYKEY"))  // value
fmt.Println(viper.GetString("mykey"))  // value
```

### 陷阱 3：环境变量命名

```go
// 配置: database.host
// 环境变量需要使用前缀和替换规则

viper.SetEnvPrefix("APP")
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
viper.AutomaticEnv()

// 设置: APP_DATABASE_HOST=localhost
```

## 最佳实践

1. **使用结构体**：将配置反序列化到结构体，类型安全
2. **设置默认值**：确保配置缺失时有合理默认值
3. **验证配置**：在启动时验证必要配置
4. **分离敏感配置**：密码等敏感信息使用环境变量

## 参考资源

- [Viper GitHub](https://github.com/spf13/viper)
- [Viper Documentation](https://pkg.go.dev/github.com/spf13/viper)
