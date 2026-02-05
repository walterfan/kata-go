# 延伸阅读

以下文章来自 [Walter Fan 的博客](https://www.fanyamin.com)，与本文档主题相关，可作为深入学习的补充材料。

## 基础与陷阱

- [**Go 语言的常见陷阱**](https://www.fanyamin.com/journal/2025-03-25-go-yu-yan-de-chang-jian-xian-jing.html) — 15 个常见陷阱详解：短变量声明、指针、nil、for range、切片、字符串、switch、goroutine、channel、方法接收者、break、闭包、错误处理、并发安全、包导入
- [**通过通信来共享内存，而不是通过共享内存来通信**](https://www.fanyamin.com/journal/2025-03-26-tong-guo-tong-xin-lai-gong-xiang-nei-cun-er-bu-shi-tong-guo.html) — Go 并发哲学，C++/Java/Go 三种语言实现事件循环对比

## 架构与代码组织

- [**SoC Code Structure in Golang**](https://www.fanyamin.com/journal/2025-08-25-soc-code-structure-in-golang.html) — 关注点分离，Go 项目目录结构最佳实践
- [**Go 应用程序的代码组织**](https://www.fanyamin.com/journal/2025-08-29-go-ying-yong-cheng-xu-de-dai-ma-zu-zhi.html) — MVC 模式、依赖注入、控制反转在 Go 中的应用

## 并发与 Context

- [**Context in Go**](https://www.fanyamin.com/journal/2025-08-28-context-in-go.html) — Context 详解：取消机制、超时、请求范围值、最佳实践与常见错误
- [**警惕！你的 Go 程序正在偷偷"泄漏"**](https://www.fanyamin.com/journal/2025-12-13-go-goroutine-leak-jing-ti-ni-de-cheng-xu-zheng-zai-tou-tou-x.html) — Goroutine Leak 详解：排查工具（pprof、goleak）、修复方案

## 安全与访问控制

- [**Go 微服务访问控制之 Casbin 实践指南**](https://www.fanyamin.com/journal/2025-07-13-go-casbin-wei-fu-wu-fang-wen-kong-zhi-zhi-shi-jian-zhi-nan.html) — Casbin + JWT + Gin 实现 RBAC 权限控制

## 调试与崩溃分析

- [**Go 程序崩溃分析实战**](https://www.fanyamin.com/journal/2026-01-23-golang_crash_analysis.html) — Coredump 生成、Delve 分析、预防措施
- [**Debug Build 的两种哲学：C++ 宏 vs Go 链接器注入**](https://www.fanyamin.com/journal/2026-02-03-debug_build_cpp_vs_go.html) — `-ldflags -X` 与 Build Tags 详解
