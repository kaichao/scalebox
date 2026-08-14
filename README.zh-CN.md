# Scalebox — 云原生的流式计算引擎

Scalebox 是一种云原生的流式计算引擎，可在分布式、异构计算集群上运行容器化的用户算法，以流水线方式组织模块层级的大规模并行处理，并支持任务级容错。与已有大数据处理、并行计算等框架相比，其技术特点特别适用于数据分布、算力分散、算法复杂、数据规模触及内存墙等应用场景。

## ✨ 核心特性

- **云原生设计**：所有算法模块容器化封装，通过边车模式嵌入数据处理流水线；控制消息与数据通道分离，前后模块以消息总线关联，实现多语言非侵入式并行编程。

- **虚拟任务（VTask / 任务组）**：在细粒度 Task 之上建立应用级粗粒度计算单元。跨模块 Task 集合内置信号量和共享变量，由 wait-queue → vtask-head → vtask-core → vtask-tail 管道统一管理流控、资源绑定和状态追踪。支持 DEFAULT / HOST-BOUND / GROUP-BOUND 三种模式，覆盖从轻量批处理到多节点协同计算的全场景。

- **跨集群计算**：归一化处理算法模块和传输模块，通过流水线统一处理集群内/跨集群的数据。t_cluster 全网格复制 + gRPC 代理转发，CLI 始终连接本地 controld，跨集群操作由服务端透明路由。

- **多层级容错**：针对硬件故障、软件 bug、网络问题、数据异常等偶发性错误，基于退出码规则实现自动重试。其中，细粒度的任务级容错可在不可靠硬件上实现可信数据处理；VTask 任务组级容错支持节点出错时的一致性任务恢复。

- **本地计算优化**：以节点本地存储为中心，通过"算子级空间展开 + 本地存储驻留 + 显式数据流编排"，将节点间通信前移到节点内部完成，显著降低对网络带宽和外部存储的依赖。

- **WebUI 管理界面**：REST 网关（:8088），76 个端点覆盖全部管理功能。SPA 前端内嵌单文件部署，提供 Dashboard 统计大盘、App 卡片网格、DAG 拓扑图、Task/VTask 管理、信号量/变量层级树。

- **VS Code 插件**：侧边栏 TreeView（App → Modules/Tasks/VTasks 三级懒加载树）、任务日志 Webview、DAG 力导向图（ECharts）、app.yaml LSP（语法校验/自动补全/悬停提示）。Docker 打包，零本地 Node.js 依赖。

- **安全框架**：JWT 认证 + RBAC 角色体系（admin/viewer/operator/automation）+ gRPC TLS 加密 + 数据库证书认证。Compose 分层可选启用，默认关闭零开销。

- **多级并行化**：模块内算法并行、模块级数据并行、模块间流水线并行。

- **多容器引擎**：Docker（默认）、Singularity/Apptainer、Podman。

## 🚀 快速开始

### 1. 环境准备

```bash
# Docker 26.1+
curl -fsSL https://get.docker.com | sh
sudo systemctl enable --now docker
docker --version
```

### 2. 启动 Scalebox

```bash
git clone https://github.com/kaichao/scalebox.git
cd scalebox

# 生成密钥（可选）
cd build && bash gen-secrets.sh && cd ..

# 构建并启动
make -C build/
docker compose -f build/compose.yaml up -d
```

浏览器打开 `http://localhost:8088` 进入 WebUI。

### 3. 运行示例

```bash
cd examples/hello-scalebox
echo "Hello Scalebox" | scalebox run
scalebox app list
```

### 4. 安装 VS Code 插件

```bash
cd vscode-scalebox && make install
```

## 📊 核心概念

| 概念 | 英文 | 中文 | 说明 |
|:---|:---|:---|:---|
| App | Application | 应用 | 完成特定计算任务的应用程序，包含多个 Module |
| Module | Module | 模块 | 容器化封装的算法组件，通过级联形成流水线 |
| Task | Task | 任务 | 基本运行单位，输入数据在 Module 上的执行过程 |
| VTask | Virtual Task | 任务组 | 跨模块 Task 集合，内置信号量和变量，统一流控和状态管理 |
| Host | Host | 节点 | 执行计算任务的服务器 |
| Slot | Slot | 插槽 | 节点上的计算资源切片，Task 调度的基础单位 |
| Cluster | Cluster | 集群 | 计算资源的逻辑分组，支持跨广域网协同 |

## 📚 文档

### 中文用户文档（[→ 完整目录](docs/cn/source/index.rst)）

| 分类 | 文档 | 说明 |
|------|------|------|
| 入门 | [Scalebox 简介](docs/cn/source/started/1_introduction.md) | 产品定位、主要特性、核心价值 |
| 入门 | [快速上手](docs/cn/source/started/2_quick_start.md) | 环境准备、部署运行、验证安装 |
| 使用 | [安装部署](docs/cn/source/user/1_installation.md) | Docker Compose 单节点 / 多节点集群 |
| 使用 | [核心概念](docs/cn/source/user/2_core_concepts.md) | App/Module/Task/VTask/Host/Slot/Cluster 详解 |
| 使用 | [运行应用](docs/cn/source/user/3_running_apps.md) | 创建、管理、监控应用和任务 |
| 使用 | [示例应用](docs/cn/source/user/4_example_apps.md) | 入门到高级示例 |
| 使用 | [运维管理](docs/cn/source/user/5_operations.md) | 监控、排错、备份、集群管理 |
| 使用 | [标准模块](docs/cn/source/user/6_standard_modules.md) | 文件传输、定时任务、目录列表等 |
| 使用 | [WebUI 指南](docs/cn/source/user/7_webui.md) | REST 网关、页面导航、构建部署 |
| 使用 | [VS Code 插件](docs/cn/source/user/8_vscode.md) | 安装配置、TreeView、DAG、LSP |
| 开发 | [编程模型](docs/cn/source/developer/1_programming_model.md) | 两级编程模型、事件驱动架构 |
| 开发 | [模块开发](docs/cn/source/developer/2_module_development.md) | 模块设计、sidecar 模式、单元测试 |
| 开发 | [应用设计](docs/cn/source/developer/3_app_design.md) | 设计原则、流水线模式、最佳实践 |
| 开发 | [主路由与状态](docs/cn/source/developer/4_main_router_status.md) | 信号量/变量/全局变量、VTask 概览 |
| 开发 | [节点本地计算](docs/cn/source/developer/5_node_local_compute.md) | 存算一体、三阶段模型、性能分析 |
| 开发 | [高级特性](docs/cn/source/developer/7_advanced_features.md) | 容错、准入控制、超时、slot 自动扩缩 |
| 开发 | [VTask 设计](docs/cn/source/developer/9_vtask.md) | 概念模型、模块结构、信号量、管道流程 |
| 开发 | [跨集群架构](docs/cn/source/developer/10_cross_cluster.md) | 数据复制、地址解析、gRPC 代理 |
| 开发 | [安全框架](docs/cn/source/developer/11_security.md) | JWT + RBAC + TLS、证书管理 |
| 参考 | [技术规范](docs/cn/source/appendix/1_tech_specifications.md) | app.yaml / module / cluster 定义规范 |
| 参考 | [参数手册](docs/cn/source/appendix/2_parameter_reference.md) | 全部参数和环境变量 |
| 参考 | [CLI 命令](docs/cn/source/appendix/3_commandline_tools.md) | 24 个子命令完整参考 |
| 参考 | [Shell 编程](docs/cn/source/appendix/4_shell_programming.md) | 内置函数、容器目录、文件交换接口 |
| 参考 | [退出码规范](docs/cn/source/appendix/6_exit_code_spec.md) | 应用退出码约定和调度策略 |
| 参考 | [FAQ](docs/cn/source/faq.md) | 部署、任务、信号量、WebUI 常见问题 |

### 英文设计文档

| 文档 | 说明 |
|------|------|
| [VTask 设计](docs/vtask-design.md) | VTask 完整设计文档 |
| [WebUI / VS Code](docs/webui-vscode-plan.md) | WebUI + VS Code 插件架构方案 |
| [REST API](docs/rest-api-reference.md) | 76 端点完整参考 |
| [gRPC API](docs/grpc-api.md) | 95 RPC 调用示例 |
| [安全框架](docs/security.md) | JWT + RBAC + TLS 配置 |

## 🧪 示例应用

- **[hello-scalebox](examples/hello-scalebox/)** — 第一个入门应用
- **[app-primes](examples/app-primes/)** — 质数计算，展示数据并行
- **[app-copy](examples/app-copy/)** — 跨节点数据传输
- **[remote-primes](examples/remote-primes/)** — 跨集群质数计算
- **[vtask](examples/vtask/)** — VTask 管道示例（DEFAULT / HOST-BOUND）

## 🏗️ 标准模块

- **dir-list** — 目录列表，生成待处理文件清单
- **file-copy** — 基于 rsync-over-ssh 的单文件拷贝
- **dir-copy** — 目录级拷贝
- **rsync-copy** — 高效 rsync 传输
- **ftp-copy** — FTP 协议传输
- **cron** — 定时消息触发
- **cluster-head** / **node-agent** — 集群节点管理

## 🔗 相关软件

- [PostgreSQL](https://github.com/postgres/postgres) — 元数据存储
- [gRPC](https://github.com/grpc/grpc) — 跨组件高效通信协议
- [Go](https://github.com/golang/go) — 云原生应用开发语言

## 🤝 贡献

欢迎提交 [Issue](https://github.com/kaichao/scalebox/issues/new) 或 Pull Request。

## 📄 许可证

[Apache License 2.0](LICENSE) © Kaichao Wu
