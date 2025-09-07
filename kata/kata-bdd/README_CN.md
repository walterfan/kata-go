# 编码器 BDD 演示（Go + Godog）

本项目使用 `godog` 展示如何在 Go 语言中进行行为驱动开发（BDD）。
我们实现了一个简单的编码/解码器，支持 Base64、Hex 和 URL 编码。

## 什么是 BDD？

BDD（行为驱动开发）是一种通过业务可读的示例来定义软件行为的实践。
这些示例以用例（场景）的形式表达，并自动化为验收测试，指导开发。

要点：
- 规格以可执行示例（Scenario）表达
- 常见结构：`Given`（前置条件）、`When`（动作）、`Then`（结果）
- 促进开发、测试与业务的共同理解

## 什么是 Godog？

Godog 是 Go 语言的 Cucumber 实现。它允许你用 Gherkin（`.feature` 文件）描述行为，并将每个步骤绑定到 Go 函数。Godog 执行这些场景并输出结果，使之成为可执行规格与验收测试。

常见使用方式：
- 作为独立 CLI 工具（直接运行特性）
- 集成到 `go test`（在测试套件内运行）

### 在本项目中的用法

- 功能文件位于 `bdd/*.feature`
- 步骤定义与测试运行器位于 `bdd/encoder_steps_test.go`
- 通过 `go test` 集成：构造 `godog.TestSuite` 并指向特性路径

可选：安装 Godog CLI 直接运行

```
go install github.com/cucumber/godog/cmd/godog@latest
# 在项目根目录下，运行 bdd/ 下的特性
godog bdd
```

## 项目结构

```
.
├── bdd
│   ├── encoder.feature          # 功能与场景（Gherkin）
│   └── encoder_steps_test.go    # 步骤定义 + 测试运行器
├── cmd
│   └── encoder
│       └── main.go              # 手动使用的命令行工具
├── pkg
│   └── encoder
│       └── encoder.go           # 编码/解码实现
└── go.mod
```

## 功能文件（Gherkin）

行为在 `bdd/encoder.feature` 中描述：

```
Feature: Encoding and decoding text
  As a user of the encoder tool
  I want to encode and decode strings using base64, hex, and url schemes
  So that I can transform text reliably

  Scenario Outline: Encode text
    Given I have the text "<plain>"
    When I encode using "<type>"
    Then the result should be "<encoded>"

    Examples:
      | type   | plain        | encoded                         |
      | base64 | hello world  | aGVsbG8gd29ybGQ=                |
      | hex    | hello        | 68656c6c6f                      |
      | url    | hello world! | hello+world%21                  |

  Scenario Outline: Decode text
    Given I have the text "<encoded>"
    When I decode using "<type>"
    Then the result should be "<plain>"

    Examples:
      | type   | encoded                         | plain        |
      | base64 | aGVsbG8gd29ybGQ=                | hello world  |
      | hex    | 68656c6c6f                      | hello        |
      | url    | hello+world%21                  | hello world! |
```

## 运行 BDD 测试

先决条件：
- Go 1.20+

安装依赖并运行测试：

```
go test ./...
```

你会看到 Godog 套件执行并通过。

只运行 BDD 包：

```
go test ./bdd -run TestGodog
```

## 命令行工具用法

构建并运行编码器 CLI：

```
go run ./cmd/encoder --mode encode --type base64 --text "hello world"
# 输出：aGVsbG8gd29ybGQ=

go run ./cmd/encoder --mode decode --type url --text "hello+world%21"
# 输出：hello world!
```

## BDD 如何指导实现

1. 编写功能文件（可执行规范）
2. 实现步骤定义，将 Gherkin 步骤绑定到 Go 代码
3. 实现最小代码以使场景通过
4. 在保持场景为绿色的同时进行重构

该示例展示了紧密的反馈循环：特性文件阐明期望行为；步骤断言期望；代码演进直至满足行为。

## 参考
- Godog: https://github.com/cucumber/godog
- Cucumber BDD: https://cucumber.io/docs/bdd/
