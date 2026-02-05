# Cobra 命令行框架

```{contents} 目录
:depth: 3
```

## Cobra 概述

Cobra 是一个用于创建强大的现代 CLI 应用程序的库，被广泛使用（kubectl、docker、hugo 等）。

## 安装

```bash
go get github.com/spf13/cobra/cobra
```

## 项目结构

```
myapp/
├── cmd/
│   ├── root.go
│   ├── serve.go
│   └── version.go
├── main.go
└── go.mod
```

## 基本用法

### main.go

```go
package main

import "myapp/cmd"

func main() {
    cmd.Execute()
}
```

### cmd/root.go

```go
package cmd

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "A brief description of your application",
    Long:  `A longer description...`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    
    // 全局标志
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
    
    // 本地标志
    rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        viper.SetConfigName("config")
        viper.AddConfigPath(".")
    }
    
    viper.AutomaticEnv()
    viper.ReadInConfig()
}
```

### cmd/serve.go

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the server",
    Long:  `Start the HTTP server on the specified port.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Starting server on port %d\n", port)
        // 启动服务器
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)
    
    serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
}
```

### cmd/version.go

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    Version   = "dev"
    GitCommit = "none"
    BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Version: %s\n", Version)
        fmt.Printf("Git Commit: %s\n", GitCommit)
        fmt.Printf("Build Date: %s\n", BuildDate)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
```

## 子命令

```go
// 创建子命令层级: myapp user create
var userCmd = &cobra.Command{
    Use:   "user",
    Short: "User management commands",
}

var userCreateCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new user",
    Args:  cobra.ExactArgs(1),  // 必须有一个参数
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        email, _ := cmd.Flags().GetString("email")
        fmt.Printf("Creating user: %s (%s)\n", name, email)
    },
}

func init() {
    rootCmd.AddCommand(userCmd)
    userCmd.AddCommand(userCreateCmd)
    
    userCreateCmd.Flags().StringP("email", "e", "", "User email")
    userCreateCmd.MarkFlagRequired("email")
}
```

## 参数验证

```go
var cmd = &cobra.Command{
    Use:  "process [file]",
    Args: cobra.ExactArgs(1),  // 必须正好一个参数
    // 其他选项：
    // cobra.NoArgs        - 不接受参数
    // cobra.MinimumNArgs(n) - 至少 n 个参数
    // cobra.MaximumNArgs(n) - 最多 n 个参数
    // cobra.RangeArgs(min, max) - 参数数量在范围内
    Run: func(cmd *cobra.Command, args []string) {
        // ...
    },
}

// 自定义验证
var customCmd = &cobra.Command{
    Use:  "process",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 1 {
            return errors.New("requires at least one argument")
        }
        if !isValidFile(args[0]) {
            return fmt.Errorf("invalid file: %s", args[0])
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        // ...
    },
}
```

## 标志类型

```go
func init() {
    // 字符串
    cmd.Flags().StringP("name", "n", "", "Name")
    
    // 整数
    cmd.Flags().IntP("count", "c", 0, "Count")
    
    // 布尔
    cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
    
    // 字符串数组
    cmd.Flags().StringArrayP("tags", "t", []string{}, "Tags")
    
    // 持久标志（子命令也可用）
    cmd.PersistentFlags().String("config", "", "Config file")
    
    // 必填标志
    cmd.MarkFlagRequired("name")
}
```

## Pre/Post 钩子

```go
var cmd = &cobra.Command{
    Use: "myapp",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // 在此命令及所有子命令执行前运行
        fmt.Println("PersistentPreRun")
    },
    PreRun: func(cmd *cobra.Command, args []string) {
        // 在 Run 之前执行
        fmt.Println("PreRun")
    },
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Run")
    },
    PostRun: func(cmd *cobra.Command, args []string) {
        // 在 Run 之后执行
        fmt.Println("PostRun")
    },
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
        // 在此命令及所有子命令执行后运行
        fmt.Println("PersistentPostRun")
    },
}
```

## 自动补全

```go
// 生成 bash 补全脚本
var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish]",
    Short: "Generate completion script",
    Args:  cobra.ExactValidArgs(1),
    ValidArgs: []string{"bash", "zsh", "fish"},
    Run: func(cmd *cobra.Command, args []string) {
        switch args[0] {
        case "bash":
            rootCmd.GenBashCompletion(os.Stdout)
        case "zsh":
            rootCmd.GenZshCompletion(os.Stdout)
        case "fish":
            rootCmd.GenFishCompletion(os.Stdout, true)
        }
    },
}
```

## 最佳实践

1. **分离命令**：每个命令一个文件
2. **使用 Viper**：配合 Viper 管理配置
3. **验证输入**：使用 Args 验证参数
4. **提供帮助**：编写清晰的 Short 和 Long 描述
5. **错误处理**：使用 RunE 返回错误

```go
// 使用 RunE 处理错误
var cmd = &cobra.Command{
    Use: "process",
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := doSomething(); err != nil {
            return fmt.Errorf("processing failed: %w", err)
        }
        return nil
    },
}
```

## 参考资源

- [Cobra GitHub](https://github.com/spf13/cobra)
- [Cobra User Guide](https://github.com/spf13/cobra/blob/main/user_guide.md)
