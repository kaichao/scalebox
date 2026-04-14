# CLAUDE.md

该文件为 Claude Code (claude.ai/code) 在操作此仓库代码时提供指导。

## 项目概述

Scalebox 是一个云原生流计算引擎，专为分布式异构计算集群而设计。它使容器化的单机用户算法能够在流水线组织的分层大规模并行处理上运行，并提供任务级容错能力。

## 开发工作流程

### 快速开始
```bash
# 克隆并设置
git clone https://github.com/kaichao/scalebox.git
cd scalebox

# 设置开发环境
cd runtime && make prepare && make pull-all && make get-cli

# 启动控制平面服务
cd runtime && make all

# 构建并运行 hello 示例
cd examples/hello-scalebox && make build
echo "hello" | scalebox run

# 运行分布式质数计算
cd examples/app-primes && make run
```

### 开发周期
1. **模块开发**: 在 `dockerfiles/` 中创建/修改模块
2. **构建测试**: 在模块目录中使用 `make build`
3. **集成**: 在 `examples/` 中创建应用 YAML 配置
4. **部署**: 使用 `scalebox run` 运行应用
5. **监控**: 使用 `scalebox app list` 和 `scalebox app logs`

## 项目结构

### 核心目录
- `runtime/`: 控制平面服务 (controld, actuator, database)
- `dockerfiles/`: 标准可重用模块和 Docker 定义
- `examples/`: 多语言实现的应用示例
- `pkg/`: 共享 Go 包 (auth, postgres, semaphore, task 等)
- `docs/`: 英文和中文文档
- `features/`: 特定功能测试和演示

### 关键概念
- **应用 (App)**: 在 YAML 中定义的完整工作流，包含模块和参数
- **模块 (Module)**: 具有消息输入/输出的容器化处理组件
- **任务 (Task)**: 由消息触发的工作执行单元
- **槽位 (Slot)**: 集群节点上的计算资源分配
- **集群 (Cluster)**: 资源集合（内联管理或外部 Slurm/K8s）

## 关键命令

### 环境设置
```bash
# 安装依赖并准备环境
cd ~/scalebox/runtime && make prepare

# 拉取容器镜像并获取 CLI 工具
cd ~/scalebox/runtime && make pull-all
make get-cli

# 启动 scalebox 控制服务
cd ~/scalebox/runtime && make all
```

### 应用开发

#### 构建应用
```bash
# 构建 hello-scalebox 示例
cd examples/hello-scalebox && make build

# 构建 app-primes 示例  
cd examples/app-primes && make build

# 使用自定义标签构建
cd examples/app-primes && make build TAG=v2
```

#### 运行应用
```bash
# 运行 hello-scalebox (消息触发)
cd examples/hello-scalebox && echo "hello scalebox" | scalebox run

# 运行分布式质数计算
cd examples/app-primes && make run NUM_GROUPS=4 CALC_NODE=local NUM_PARALLEL=2

# 使用自定义集群运行
cd examples/app-primes && make run CLUSTER=inline NUM_GROUPS=8

# 列出正在运行的应用
scalebox app list

# 停止应用
scalebox app stop <app-name>

# 查看应用日志
scalebox app logs <app-name>
```

#### 应用配置
应用使用 YAML 文件定义 (如 `app.yaml`)：
- `name`: 带变量的应用标识符 (如 `app-primes-g${NUM_GROUPS}-${TAG}`)
- `cluster`: 目标集群名称
- `modules`: 具有容器镜像和槽位分配的处理组件
- `parameters`: 包含任务分发模式的运行时配置

### 模块开发

#### 构建模块
```bash
# 构建所有标准模块
cd dockerfiles && make build

# 构建特定模块
cd dockerfiles/dir-list && make build

# 使用自定义标签构建
cd dockerfiles/dir-list && make build TAG=v2
```

#### 模块结构
标准模块遵循以下模式：
- `Dockerfile`: 包含入口点的容器定义
- `main.go`: 具有消息处理的 Go 实现
- `go.mod`: Go 模块依赖
- `*_test.go` 文件中的测试

#### 模块类型
1. **文件操作**: `dir-list`, `file-copy`, `rsync-copy`, `ftp-copy`, `rsyncd`
2. **数据处理**: `data-grouping-2d` 用于 2D 数据集分组操作
3. **工具**: `cron` (定时消息), `actuator` (密钥生成)

#### Go 模块开发
```bash
# 运行模块测试
cd dockerfiles/dir-list && go test ./...

# 运行覆盖率测试
cd dockerfiles/data-grouping-2d && go test -coverprofile=coverage.out ./...

# 基准测试
cd pkg/semagroup && go test -bench=. ./...
```

#### 共享包
- `pkg/auth/`: 认证和授权 (JWT, basic auth)
- `pkg/postgres/`: PostgreSQL 数据库操作
- `pkg/semaphore/`: 基于信号量的流控制
- `pkg/task/`: 任务管理实用程序
- `pkg/common/`: 通用网络和 JSON 实用程序

### 文档
```bash
# 构建文档
cd docs && make html
cd docs/cn && make html
```

## 架构概述

### 核心组件
1. **控制平面** (runtime/ 目录):
   - `actuator`: 启动器服务，用于在计算节点上启动槽位
   - `controld`: 基于 gRPC 的控制服务，用于管理执行器和计算节点
   - `database`: PostgreSQL 数据库，存储应用/模块/任务/槽位数据

2. **标准模块** (dockerfiles/ 目录):
   - 文件操作: `dir-list`, `file-copy`, `rsync-copy`, `ftp-copy`, `rsyncd`
   - 数据处理: `data-grouping-2d`
   - 工具: `cron` (定时消息), `actuator` (密钥生成)

3. **应用示例** (examples/ 目录):
   - `hello-scalebox`: 基础的消息触发应用
   - `app-primes`: 分布式质数计算
   - `app-copy`: 跨集群数据传输演示

4. **集群管理** (runtime/ 目录):
   - 安装和配置指南
   - 单节点和多节点集群设置
   - 环境准备和服务管理

### 关键概念
- **应用 (App)**: 完整的应用工作流
- **模块 (Module)**: 单个处理组件
- **任务 (Task)**: 工作执行单元
- **槽位 (Slot)**: 计算资源分配
- **集群 (Cluster)**: 计算资源的集合（内置或外部）

### 开发模式
1. **模块开发**: 创建具有特定功能的 Docker 容器
2. **应用组合**: 使用 YAML 配置将模块链接在一起
3. **多语言支持**: 支持 Python, Go, C, Java, Julia, R 等多种语言示例
4. **跨集群操作**: 内置支持跨异构集群的分布式计算

### 文件结构
- `examples/`: 多语言实现的应用示例
- `dockerfiles/`: 标准可重用模块
- `tests/`: 功能测试（容错、超时、任务视角等）
- `docs/`: 英文和中文文档
- `runtime/`: 控制平面服务和集群管理

### 测试和验证

#### 运行测试
```bash
# 运行 Go 模块单元测试
cd dockerfiles/dir-list && go test ./...
cd dockerfiles/data-grouping-2d && go test ./...

# 运行功能测试
cd features/retry_test && make test
cd features/timeout && make test
cd features/check_test && make test
cd features/task-perspective && make test
cd features/cross-cluster-primes && make test
cd features/singularity && make test
```

#### 测试类别
- **容错测试** (`features/retry_test/`): 任务重试和恢复机制
- **超时处理** (`features/timeout/`): 任务超时配置和管理
- **流控管理** (`features/check_test/`): 流控制和背压
- **任务透视** (`features/task-perspective/`): 任务监控和调试
- **跨集群计算** (`features/cross-cluster-primes/`): 多集群分布式处理
- **Singularity 支持** (`features/singularity/`): 替代容器运行时支持

#### 验证
- 多语言质数计算示例（Python, Go, C, Java, Julia, R）
- 跨集群功能验证
- PostgreSQL 数据库集成测试
- 组件间 gRPC 通信测试

该架构支持云原生设计，具有容器化算法、消息驱动处理和任务级容错能力，适用于大规模分布式计算应用。