.. Go 实战开发指南 documentation master file

Golang 实战开发指南
=======================================

.. image:: https://go.dev/images/gophers/ladder.svg
   :alt: Go Gopher
   :width: 200px
   :align: center

.. include:: links.ref
.. include:: tags.ref
.. include:: abbrs.ref

============= ==============================================================
**Abstract**  Golang 实战开发指南 - 专注于易错和易忽略的知识点
**Authors**   Walter Fan
**Category**  learning note
**Status**    WIP
**Updated**   |date|
**License**   |CC-BY-NC-ND|
============= ==============================================================

欢迎阅读 **Golang 实战开发指南**！本文档专注于 Go 语言开发中程序员 **易错** 和 **易忽略** 的知识点，
以及 Go 相比其他语言 **特殊** 的地方。

目标读者
--------

- 有一定编程基础，想深入学习 Go 的开发者
- 从 C++/Java/Python 等语言转向 Go 的程序员
- 希望避免常见陷阱，写出高质量 Go 代码的工程师

文档特点
--------

- **实用导向**：每个知识点都有实际应用场景
- **避坑指南**：重点讲解易错、易混淆的概念
- **最佳实践**：提供业界推荐的编码方式
- **代码示例**：可直接运行的完整代码

.. toctree::
   :maxdepth: 2
   :caption: 目录
   :numbered:

   1.basic/index
   2.concurrency/index
   3.performance/index
   4.memory/index
   5.network/index
   6.library/index
   7.cheatsheet/index
   references

快速导航
--------

.. grid:: 2
   :gutter: 3

   .. grid-item-card:: 🚀 基础语法与陷阱
      :link: 1.basic/index
      :link-type: doc

      Go 语言核心概念，常见陷阱与解决方案，从其他语言迁移的注意事项。

   .. grid-item-card:: ⚡ 并发编程
      :link: 2.concurrency/index
      :link-type: doc

      Goroutine、Channel、sync 包详解，并发模式与最佳实践。

   .. grid-item-card:: 📊 性能调优
      :link: 3.performance/index
      :link-type: doc

      pprof 性能分析、基准测试、逃逸分析、编译优化。

   .. grid-item-card:: 🧠 内存管理
      :link: 4.memory/index
      :link-type: doc

      GC 原理、内存分配、内存泄漏排查、对象池使用。

   .. grid-item-card:: 🌐 网络编程
      :link: 5.network/index
      :link-type: doc

      net/http、gRPC、WebSocket、TCP/UDP 编程。

   .. grid-item-card:: 📦 常用库
      :link: 6.library/index
      :link-type: doc

      Viper、Cobra、Zap、GORM、Gin 等常用库详解。


Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
