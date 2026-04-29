# Scalebox - 云原生的流式计算引擎

Scalebox是一种云原生的流式计算引擎，可在分布式、异构计算集群上运行容器化的用户算法，以流水线组织模块层级上大规模并行处理，支持任务级容错。

## ✨ 核心特性

- **云原生设计**：所有算法模块容器化封装，通过边车模式嵌入到面向云环境的数据处理流水线；控制消息与数据通道分离，前后模块间以消息总线关联，实现多语言的非侵入式并行编程

- **跨集群计算**：归一化处理算法模块和传输模块，通过流水线统一处理集群内/跨集群的数据，屏蔽数据和计算的跨集群差异

- **任务级容错**：对于硬件故障、软件bug、网络问题、数据异常等原因导致的偶发性出错，基于规则实现自动容错处理

- **本地计算优化**：以节点本地存储为中心，通过"算子级空间展开 + 本地存储驻留 + 显式数据流编排"，把大量本该发生在节点间互联通信前移到节点内部完成

- **多级并行化**：
  - 模块内算法并行（多线程、GPU加速等）
  - 模块级数据并行（同一模块多个实例处理不同数据）
  - 模块间流水线并行（不同模块通过流水线方式并行执行）

- **多语言支持**：支持Python、Go、C++、Java、Shell等多种语言的算法实现

## 🚀 快速开始

### 1. 环境准备

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh
sudo systemctl enable --now docker

# 验证 Docker 安装
docker --version
docker run hello-world
```

### 2. 获取并运行 Scalebox

```bash
# 克隆仓库
git clone https://github.com/kaichao/scalebox.git
cd scalebox

# 启动控制平面服务
cd runtime && make all

# 验证服务状态
docker ps
```

### 3. 运行 Hello Scalebox 示例

```bash
# 运行第一个应用
cd examples/hello-scalebox
echo "Docker-based_Scalebox" | scalebox run

# 查看应用状态
scalebox app list

# 查看任务执行情况
scalebox task list
```

### 4. 运行分布式计算示例

```bash
# 运行质数计算应用
cd examples/app-primes
make run NUM_GROUPS=4 CALC_NODE=local NUM_PARALLEL=2

# 查看应用日志
scalebox app logs <app-name>
```

## 📊 架构概述

### 核心概念

| 概念 | 英文名称 | 中文名称 | 说明 |
| :--- | :--- | :--- | :--- |
| 应用程序 | **App** | **应用** | 完成特定计算任务的应用程序，包含多个模块 |
| 软件单元 | **Module** | **模块** | 构成应用的基本可编程单元，容器化封装的程序代码集合 |
| 执行实例 | **Task** | **任务** | 基本运行单位，输入数据在具体模块上的执行过程 |
| 计算节点 | **Host** | **节点** | 执行计算任务的服务器节点 |
| 资源单元 | **Slot** | **插槽** | 节点上与模块对应的计算资源切片 |

### 系统架构

```
┌─────────────────────────────────────────┐
│           应用层 (Application)           │
│   ┌─────────┐  ┌─────────┐  ┌─────────┐ │
│   │   App   │  │   App   │  │   App   │ │
└───┴─────────┴──┴─────────┴──┴─────────┴─┘
┌─────────────────────────────────────────┐
│           模块层 (Module)                │
│   ┌─────────┐  ┌─────────┐  ┌─────────┐ │
│   │ Module  │  │ Module  │  │ Module  │ │
└───┴─────────┴──┴─────────┴──┴─────────┴─┘
┌─────────────────────────────────────────┐
│         运行时层 (Runtime)                │
│  controld │ actuator │ database │ ...   │
└─────────────────────────────────────────┘
```

### 核心组件

- **controld**：gRPC-based 控制服务，管理执行器和计算节点
- **actuator**：启动服务，通过 SSH 或外部调度器在计算节点上启动插槽
- **database**：PostgreSQL 数据库，存储应用/模块/任务/插槽元数据

## 🎯 应用场景

### 大规模数据的复杂处理
- **天文计算**：大型天文望远镜观测数据的处理
- **基因组学和生物信息学**：基因组测序、组装、比对等
- **高能物理**：粒子对撞机实验数据分析

### 跨集群算力应用
- 大规模数据传输与处理
- 分布式计算资源统一调度
- 跨广域网异构算力集群的统一管理

### 基于容器化的跨平台嵌入式仿真
- 软件模块标准化封装
- 多平台、多模块集成测试
- 全系统仿真

### 大模型训练
- 原生支持数据并行、流水线并行
- 支持模型数据按层划分的张量并行

## 📚 详细文档

### 入门教程
- [Scalebox 简介](docs/cn/source/started/1_introduction.md) - 技术起源、核心概念、主要特性
- [快速上手](docs/cn/source/started/2_quick_start.md) - 环境准备、部署运行、验证安装
- [下一步](docs/cn/source/started/3_next_steps.md) - 学习路径和进阶指南

### 使用指南
- [安装部署](docs/cn/source/user/1_installation.md) - 系统要求、Docker 安装、集群配置
- [核心概念](docs/cn/source/user/2_core_concepts.md) - App、Module、Task、Slot 详解
- [运行应用](docs/cn/source/user/3_running_apps.md) - 应用创建、运行、监控和管理
- [示例应用](docs/cn/source/user/4_example_apps.md) - 各种示例应用详解
- [标准模块](docs/cn/source/user/6_standard_modules.md) - 文件操作、数据处理、工具模块

### 编程指南
- [编程模型](docs/cn/source/developer/1_programming_model.md) - 两级编程模型、消息驱动机制
- [模块开发](docs/cn/source/developer/2_module_development.md) - 模块设计、实现、测试
- [应用设计](docs/cn/source/developer/3_app_design.md) - 应用架构设计、YAML 配置
- [性能优化](docs/cn/source/developer/6_performance_optimization.md) - 并行优化、I/O 优化
- [高级特性](docs/cn/source/developer/7_advanced_features.md) - 容错机制、准入控制

### 附录
- [技术规格](docs/cn/source/appendix/1_tech_specifications.md) - 应用规范、模块规范
- [命令行工具](docs/cn/source/appendix/3_commandline_tools.md) - 所有命令行工具详解
- [最佳实践](docs/cn/source/appendix/5_best_practices.md) - 开发、部署、运维最佳实践
- [状态码](docs/cn/source/appendix/6_status_codes.md) - 状态码说明和错误处理

## 🧪 示例应用

本仓库包含多个应用示例：

- **[hello-scalebox](examples/hello-scalebox/)** - Scalebox 的第一个入门应用
- **[app-primes](examples/app-primes/)** - 计算区间内质数总数量
- **[remote-primes](examples/remote-primes/)** - 跨集群质数计算演示
- **[app-copy](examples/app-copy/)** - 跨集群数据拷贝示例
- **[cluster-dir-copy](examples/cluster-dir-copy/)** - 集群目录拷贝示例

## 🔧 特性测试

- **[retry_test](tests/retry_test/)** - 容错支持测试
- **[timeout-gen](tests/timeout-gen/)** - 超时设置测试
- **[check_test](tests/check_test/)** - 流控管理测试
- **[task-perspective](tests/task-perspective/)** - 任务透视测试
- **[cross-cluster-primes](tests/cross-cluster-primes/)** - 跨集群计算测试

## 🏗️ 标准模块

Scalebox 提供多种标准模块：

### 文件操作模块
- **dir-list** - 目录列表功能
- **file-copy** - 文件拷贝功能
- **rsync-copy** - 基于 rsync 的文件传输
- **ftp-copy** - 基于 FTP 的文件传输
- **rsyncd** - rsync 服务器模块

### 数据处理模块
- **data-grouping-2d** - 2D 数据集分组操作

### 工具模块
- **cron** - 定时消息发送
- **actuator** - 密钥生成器

## 🤝 贡献指南

非常欢迎你的加入！

1. [提交 Issue](https://github.com/kaichao/scalebox/issues/new) 报告问题或建议新功能
2. Fork 项目并提交 Pull Request
3. 参与文档改进、测试用例编写
4. 分享使用经验和应用案例

## 📄 许可证

[Apache License 2.0](LICENSE) © Kaichao Wu

## 🔗 相关软件

- [Docker](https://www.docker.com/) - 容器化平台
- [PostgreSQL](https://github.com/postgres/postgres) - Scalebox 后台数据库
- [gRPC](https://github.com/grpc/grpc) - 不同软件模块间的高效通信协议
- [Go](https://github.com/golang/go) - 云原生应用的程序语言

---

*想要了解更多？请访问 [详细技术文档](docs/cn/source/index.rst) 获取完整信息。*