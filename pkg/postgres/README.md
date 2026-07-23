# 数据库连接配置

go-scalebox 系列项目的统一数据库连接管理，支持密码和证书两种认证方式。

## 认证模式

| 模式 | 场景 | 客户端凭证 | pg_hba.conf |
|------|------|-----------|-------------|
| **密码** | 开发/测试 | `PGPASS` 环境变量 | `scram-sha-256` |
| **证书** | 生产环境 | `client.crt` + `client.key` | `cert`（CN 匹配用户名） |

模式由客户端自动检测：`$PG_CERT_DIR/client.crt` 存在则走证书，否则走密码。

## 客户端环境变量

### 连接目标

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `DATABASE_URL` | — | 完整连接串，最高优先级 |
| `PGURL` | — | 同上，存量兼容 |
| `PGHOST` | fallback 链¹ | 主机地址 |
| `PGPORT` | `5432` | 端口 |
| `PGUSER` | `scalebox` | 用户名 |
| `PGDB` | `scalebox` | 数据库名 |

### 认证

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `PGPASS` | — | 密码（密码模式） |
| `PG_CERT_DIR` | `./certs` | 证书目录（证书模式） |

### 连接池

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `PG_MAX_CONNS` | `20` | pgxpool 最大连接（GetPgxPool） |
| `PG_MIN_CONNS` | `5` | pgxpool 最小连接（GetPgxPool） |
| `PG_MAX_IDLE_CONNS` | `1` | sql.DB 最大空闲连接（GetDB） |
| `PG_MAX_OPEN_CONNS` | `4` | sql.DB 最大打开连接（GetDB） |
| `PG_MAX_CONN_LIFETIME_MIN` | `30` | pgxpool 连接最大存活（分钟） |
| `PG_MAX_CONN_IDLE_TIME_MIN` | `5` | pgxpool 连接最大空闲（分钟） |

> ¹ PGHOST fallback：`PGHOST` → `GRPC_SERVER` 的 IP → `LOCAL_ADDR` → 本机 IP

## 连接串构建优先级

```
1. DATABASE_URL          完整控制权，有值直接返回
2. PGURL                 存量兼容
3. 证书文件存在           检测 $PG_CERT_DIR/client.crt + client.key
   └─ sslmode=verify-full（CA 存在）/ require（CA 不存在）
4. PGPASS                密码认证
```

## 证书体系

同一内部 CA 签发所有服务间通信证书，由 `gen-certs.sh` 一次性生成：

```
./certs/
├── ca.crt              CA 公钥（所有容器挂载）
├── ca.key              CA 私钥（仅签发用，不部署）
│
├── db/                 PostgreSQL 服务端
│   ├── server.crt        CN=scalebox-db
│   └── server.key
│
├── controld/           controld gRPC 服务端（compose.tls.yaml 使用）
│   ├── server.crt        CN=controld
│   └── server.key
│
└── client/             所有数据库客户端共用（PG_CERT_DIR 指向此目录）
    ├── ca.crt            验证服务端
    ├── client.crt        CN=scalebox
    └── client.key
```

```bash
cd build
bash gen-certs.sh              # 生成所有证书到 ./certs/
```

### 各容器挂载

| 容器 | 挂载 | 用途 |
|------|------|------|
| database | `./certs/ca.crt` + `./certs/db/` | PostgreSQL SSL |
| controld | `./certs/controld/` + `./certs/client/` | gRPC TLS + DB 客户端 |
| actuator | `./certs/client/` | DB 客户端 |
| agent | `scp certs/client/ agent:/etc/scalebox/certs/` | DB 客户端 |
| CLI | `export PG_CERT_DIR=./certs/client` | DB 客户端 |

## 连接变化检测

`GetDB()` 和 `GetPgxPool()` 每次调用时比较当前连接串与上次是否一致。任何环境变量变化（`DATABASE_URL`、`PGHOST`、`PGPASS`、`PG_CERT_DIR` 等）都会触发旧连接关闭、新连接重建。

## 数据库服务端

### 环境变量

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `AUTH_MODE` | `password` | `password` 或 `cert` |
| `POSTGRES_PASSWORD` | —（password 模式必填） | 数据库密码 |
| `PASSWORD_FALLBACK_CIDR` | — | cert 模式下允许密码降级的 IP 范围 |

## 部署

go-scalebox 的 `build/` 目录提供分层 compose 文件，各层独立开关，按需组合：

| 文件 | 功能 |
|------|------|
| `compose.yaml` | 基础服务（controld、actuator、database、redis） |
| `compose.tls.yaml` | gRPC 传输加密（controld↔agent↔CLI） |
| `compose.db-cert.yaml` | 数据库证书认证 |
| `compose.security.yaml` | JWT 安全认证（Token Service） |
| `compose.security-local-key.yaml` | JWT 安全认证（本地私钥，无需 Token Service） |

### 用法

```bash
cd build

# 首次部署前生成证书
bash gen-certs.sh

# 开发环境：仅密码模式
docker compose -f compose.yaml up -d

# gRPC TLS
docker compose -f compose.yaml -f compose.tls.yaml up -d

# gRPC TLS + 数据库证书
docker compose -f compose.yaml -f compose.tls.yaml -f compose.db-cert.yaml up -d

# gRPC TLS + 数据库证书 + JWT 认证
docker compose -f compose.yaml -f compose.tls.yaml -f compose.db-cert.yaml -f compose.security.yaml up -d

### 本地开发（直连数据库）

```bash
# 证书模式（需先运行 gen-certs.sh）
export PGHOST=localhost
export PG_CERT_DIR=./build/certs/client

# 密码模式
export PGHOST=localhost
export PGPASS=dev123
```

## 项目覆盖

| 项目 | 获取连接方式 | 改动 |
|------|-------------|------|
| go-scalebox | `postgres.GetDB()` / `postgres.GetPgxPool()` | 自动受益 |
| workspace/scalebox | `pkg/postgres`（源头） | 修改 |
| exastore | `DATABASE_URL` → `pgxpool.NewWithConfig()` | 部署时 DATABASE_URL 加证书参数 |
| scalebox-security | `DATABASE_URL` → `pgxpool.New()` | 同上 |
| gopkg | 接收外部 `*pgxpool.Pool`（DI） | 无需改动 |
