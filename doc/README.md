# Golang 实战开发指南

这是一份使用 Sphinx 构建的 Go 语言开发文档，专注于 Go 程序员**易错**和**易忽略**的知识点，以及 Go 相比其他语言**特殊**的地方。

## 文档内容

1. **基础语法与陷阱** - Go 核心概念，常见陷阱与解决方案
2. **并发编程** - Goroutine、Channel、sync 包详解
3. **性能调优** - pprof、基准测试、逃逸分析
4. **内存管理** - GC 原理、内存分配、内存泄漏排查
5. **网络编程** - HTTP、gRPC、TCP/UDP
6. **常用库** - Viper、Cobra、Zap、GORM、Gin
7. **速查表** - 语法和命令速查

## 快速开始

### 安装依赖

```bash
cd doc
pip install -r requirements.txt
```

### 构建文档

```bash
# 构建 HTML
make html

# 或者使用实时预览
make livehtml
```

### 查看文档

```bash
# 启动本地服务器
make serve

# 访问 http://localhost:8000
```

## 开发

### 文档结构

```
doc/
├── Makefile
├── README.md
├── requirements.txt
└── source/
    ├── conf.py
    ├── index.rst
    ├── 1.basic/
    ├── 2.concurrency/
    ├── 3.performance/
    ├── 4.memory/
    ├── 5.network/
    ├── 6.library/
    └── 7.cheatsheet/
```

### 添加新内容

1. 在对应目录下创建 `.md` 或 `.rst` 文件
2. 在该目录的 `index.rst` 中添加引用
3. 运行 `make html` 构建

### 使用的扩展

- **myst-parser** - 支持 Markdown
- **sphinx-design** - 美观的卡片和网格布局
- **sphinxcontrib-mermaid** - 支持 Mermaid 图表
- **sphinx-togglebutton** - 可折叠内容
- **sphinx-copybutton** - 代码复制按钮

## 参考资源

- [Go 官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Wiki](https://go.dev/wiki/)
