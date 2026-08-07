# 1. 系统安装部署

## 1.1 环境要求

### 1.1.1 硬件要求

#### 芯片架构及操作系统

|  芯片架构  | 操作系统   |  说明              |
| -------- | --------- | ----------------- |
| x86_64   | Linux     |                   |
| arm64    | Linux     |                   |
| x86_64   | MacOS     | 用于开发环境        |
| arm64    | MacOS     | 待测试，用于开发环境 |
| x86_64   | Win64/WSL | 待测试，用于开发环境 |

生产环境中，所有节点安装64位Linux（CentOS7/8/9、Debian 12/13、Ubuntu 20/22/24/26）等。

#### 内存需求
- 头节点：≥8GB
- 计算节点：按需，推荐 ≥8GB

#### 存储空间
- 头节点：≥100GB
- 计算节点：按需，推荐 ≥100GB

### 1.1.2 软件依赖

#### 容器引擎/容器运行时

- 计算节点

|  容器引擎/容器运行时     | 版本号             |  说明     |
| --------------------- | ----------------- |  ------- |
| docker-ce             | 20.10<sup>+</sup> |          |
| singularity           | 3.8<sup>+</sup>   |          |
| podman                | 4.8<sup>+</sup>   | 待测试    |
| containerd + nerdctl  | 1.6<sup>+</sup>   | 待测试    |
| apptainer             |                   | 待测试    |
| Kata Containers       |                   | 待测试    |

- 头节点：docker-ce 20.10<sup>+</sup>，Docker Compose v2<sup>+</sup>

#### 数据库
- 头节点：postgresql 18<sup>+</sup>，以容器化部署

#### 集群存储（可选）
- 所有节点上安装集群存储客户端，构建统一的集群存储
  - glusterfs
  - NFS

#### 网络配置

- 头节点、计算节点属于内网的同一子网内
- 若需支持跨集群计算，头节点及相关传输节点需与其他集群通信的外网地址

## 1.2 安装步骤

### 1.2.1 单节点快速部署（Docker Compose）

适用于开发、测试环境，所有服务运行在单台机器上。

**获取代码**：
```bash
git clone https://github.com/kaichao/scalebox.git
cd scalebox
```

**生成密钥体系**（证书 + JWT 签名密钥）：
```bash
cd build && bash gen-secrets.sh   # → ./secrets/
```

**构建镜像**：
```bash
make -C build/                    # 构建全部服务镜像
```

**启动服务**：
```bash
docker compose -f build/compose.yaml up -d
```

启动后服务包括：
- **controld**：gRPC 控制服务（:50051）
- **actuator**：容器启动器
- **database**：PostgreSQL 数据库（:5432）
- **webui**：Web 管理界面（:8088）

**查看服务状态**：
```bash
docker compose -f build/compose.yaml ps
```

**安全层（可选）**：

启用 gRPC TLS 传输加密：
```bash
docker compose -f build/compose.yaml -f build/compose.tls.yaml up -d
```

启用数据库证书认证：
```bash
docker compose -f build/compose.yaml -f build/compose.db-cert.yaml up -d
```

启用 JWT 认证（含 Token Service）：
```bash
docker compose -f build/compose.yaml -f build/compose.security.yaml up -d
```

### 1.2.2 多节点集群部署

生产环境多节点部署，以 1 个 HEAD 节点 + 4 个 NODE 节点为例：

| 名称 | 类型 | IP 地址 |
| --- | ---- | ------- |
| h0  | HEAD | 10.0.6.100 |
| n0  | NODE | 10.0.6.101 |
| n1  | NODE | 10.0.6.102 |
| n2  | NODE | 10.0.6.103 |
| n3  | NODE | 10.0.6.104 |

**HEAD 节点配置**：
- 安装 Docker 20.10+
- 克隆代码仓库，执行 `make -C build/` 构建镜像
- 启动 controld + actuator + database + webui 容器
- 设置 `/etc/hosts`，包含所有节点的 IP 映射

**计算节点配置**：
- 安装 Docker 20.10+
- 设置 HEAD → NODE 的 SSH 免密登录（actuator 通过 SSH 在计算节点启动容器）
- 设置主机名和 `/etc/hosts`

**集群存储（可选）**：
- 安装 glusterfs 客户端，挂载共享存储卷
- 详见 :doc:`附录 7 - 基础软件安装 <../appendix/7_software_install>`

## 1.3 配置说明

### 1.3.1 环境变量

常用环境变量：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `GRPC_SERVER` | controld gRPC 地址 | `localhost:50051` |
| `DATABASE_URL` | 数据库连接串 | — |
| `LOG_LEVEL` | 日志级别（debug/info/warn/error） | `info` |
| `MY_APP_CLUSTER` | 默认集群名 | — |

数据库连接优先级：`DATABASE_URL` > `PGHOST`/`PGPORT`/`PGUSER`/`PGDB`/`PGPASS`

### 1.3.2 安全配置

| 变量 | 说明 |
|------|------|
| `SECURITY_ENABLED` | 安全总开关（默认 false） |
| `GRPC_TLS_ENABLED` | gRPC TLS 加密 |
| `AUTH_TOKEN` | JWT 认证令牌 |

安全配置详见 :doc:`安全框架文档 <../../go-scalebox/docs/security>`（英文）。

## 1.4 验证安装

### 1.4.1 健康检查

```bash
# 检查服务状态
docker compose -f build/compose.yaml ps

# 检查 controld gRPC 连通性
grpcurl -plaintext localhost:50051 list

# CLI 连接测试
scalebox cluster list
```

### 1.4.2 运行示例应用

```bash
cd examples/hello-scalebox
echo "Hello Scalebox" | scalebox run
scalebox app list
```

### 1.4.3 访问 WebUI

浏览器打开 `http://<head-node-ip>:8088`，查看 Dashboard、App 列表、集群资源等。