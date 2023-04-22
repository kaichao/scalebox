# Scalebox - 一种云原生的流式计算框架

Scalebox是一种流式计算引擎，以容器化统一传输及算法模块，通过分布式消息触发，实现跨集群模块间的流水线并行。软件平台支持任务级容错；以位置感知的消息分发，支持大文件的纯本地加载，提升横向扩展能力。特别适用于数据分布、算力资源分布的应用场景。

本仓库包含以下内容：

1. 基于docker-compose的scalebox服务端环境（[服务环境](https://github.com/kaichao/docker-scalebox/server/README.md)）
2. scalebox标准模块的Dockerfile定义 ([标准模块](https://github.com/kaichao/docker-scalebox/dockerfiles/README.md))
3. scalebox的应用示例（[应用示例](https://github.com/kaichao/docker-scalebox/examples/README.md)）
4. scalebox的主要特性的测试（[特性测试](https://github.com/kaichao/docker-scalebox/tests/README.md)）

## 内容列表

- [Scalebox - 一种云原生的流式计算框架](#scalebox---一种云原生的流式计算框架)
  - [内容列表](#内容列表)
  - [研究背景](#研究背景)
  - [环境安装](#环境安装)
  - [使用说明](#使用说明)
    - [单节点集群](#单节点集群)
    - [多节点集群](#多节点集群)
  - [应用示例](#应用示例)
  - [特性测试](#特性测试)
  - [相关软件](#相关软件)
  - [维护者](#维护者)
  - [如何贡献](#如何贡献)
  - [使用许可](#使用许可)

## 研究背景

大规模数据处理的计算框架主要分为两类：
- 大数据处理框架
  - 离线数据处理：Hadoop/Spark
  - 流式数据处理：Storm/Spark Streaming/FLink
- 基于高性能计算的MPI（Message Passing Interface）

scalebox提供了一种构建分布式数据处理的高效方法。用户仅需要研发单机版的算法模块，通过容器化打包后推到镜像库。基于系统标准模块、用户定义算法模块，定义流水线应用。

scalebox具有以下特性：
- 算网融合：统一处理跨网络的文件加载、算法处理，通过流水线并行实现高效处理
- 跨集群分布式消息：跨集群分布式消息传递，用跨集群的分布式算力支持单个应用
- 位置感知消息分发：基于位置的消息分发，支持纯本地存储的大文件加载，大幅提升I/O能力
- 横向扩展：大文件的存本地加载，消除I/O瓶颈，有效支持横向扩展


## 环境安装

- 操作系统
  - CentOS 7+（其他Linux版本待测试）
  - macos 10.15+(amd64)（ARM版待测试）

macos主要用于单节点集群的开发测试。

详细安装参见：[服务环境](server/README.md)

## 使用说明

### 单节点集群
单节点集群常用于测试、开发。

- 安装单机版系统环境
- 测试应用示例（hello-scalebox、app-primes）
- 基于标准模块，参照标准应用实例，构建自己的应用实例
- 定制自己的应用模块，并构建应用案例

### 多节点集群
多节点集群常用于生产环境。

- 在单节点集群基础上，参照```bio-down```集群定义，定义自己的内联集群
- 在应用定义文件，引用多节点集群的资源


## 应用示例

scalebox应用的初级示例，包括：
- [hello-scalebox](examples/hello-scalebox/)：scalebox的第一个入门应用
- [app-primes](examples/app-primes/)：计算区间内质数总数量
- [app-copy](examples/app-copy/)：演示最常用的跨集群数据拷贝


## 特性测试

- 容错支持：[retry_test](tests/retry_test/)
- 超时设置：[timeout-gen](tests/timeout-gen/)
- 流控管理：[check_test](tests/check_test/)
- 应用交互：[task-exec-files](tests/task-exec-files/)


## 相关软件

- [PostgreSQL Database Management System](https://github.com/postgres/postgres) — scalebox后台数据库
- [gRPC – An RPC library and framework](https://github.com/grpc/grpc) — 不同软件模块间的高效通信协议
- [The Go Programming Language](https://github.com/golang/go) — 云原生应用的程序语言
- [Pony ORM ER Diagram Editor](https://editor.ponyorm.com/) - 神奇的ER图工具

## 维护者

[@kaichao](https://github.com/kaichao)

## 如何贡献

非常欢迎你的加入！[提一个 Issue](https://github.com/kaichao/docker-scalebox/issues/new) 或者提交一个 Pull Request。


## 使用许可

[Apache License](LICENSE) © Kaichao Wu
