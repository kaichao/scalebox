# 11. 安全框架

Scalebox 安全框架提供认证（Authentication）和授权（Authorization）能力，通过 JWT + RBAC + TLS 实现 gRPC 通道安全。安全默认关闭（`SECURITY_ENABLED=false`），启用后无性能损耗外的额外运行时开销。

## 11.1 通信安全矩阵

| | controld | database | agent | actuator | CLI |
|---|:---:|:---:|:---:|:---:|:---:|
| **controld** | — | PG 客户端 | gRPC 服务端 | gRPC 服务端 | gRPC 服务端 |
| **database** | PG 服务端 | — | — | — | — |
| **agent** | gRPC 客户端 | — | — | — | — |
| **actuator** | gRPC 客户端 | — | — | — | — |
| **CLI** | gRPC 客户端 | PG 客户端 | — | — | — |

各模块携带的安全凭证：

| 模块 | 角色 | 需要 |
|------|------|------|
| controld | gRPC 服务端 + PG 客户端 | server TLS cert、JWT 公钥（验签）、DB client cert |
| agent | gRPC 客户端 | CA cert（验证服务端）、JWT token |
| actuator | gRPC 客户端 + PG 客户端 | CA cert、JWT token、DB client cert |
| CLI | gRPC 客户端 + PG 客户端 | CA cert、JWT token、DB client cert |
| database | PG 服务端 | DB server cert（PostgreSQL SSL） |

所有证书 + JWT 密钥由 `gen-secrets.sh` 统一生成，同一内部 CA 签发。

## 11.2 安全边界

- **Actuator SSH**：Ed25519 免密 SSH 登录计算节点，属基础设施层安全
- **数据库直连**：绕过 gRPC 直接访问数据库的操作不受 RBAC 控制
- **RBAC 仅覆盖 gRPC 通道**：认证与授权在 gRPC 拦截器层生效

## 11.3 启用方式

### 11.3.1 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SECURITY_ENABLED` | 安全总开关 | `false` |
| `GRPC_TLS_ENABLED` | gRPC TLS 传输加密 | `false` |
| `GRPC_TLS_CA_FILE` | CA 证书路径 | — |
| `GRPC_TLS_CERT_FILE` | 服务端证书路径 | — |
| `GRPC_TLS_KEY_FILE` | 服务端私钥路径 | — |
| `AUTH_TOKEN` | JWT 认证令牌 | — |
| `SCALEBOX_CERTS_DIR` | 证书目录（设置后启用 DB 证书认证） | — |

### 11.3.2 Compose 分层

| 文件 | 功能 |
|------|------|
| `compose.yaml` | 基础服务（无安全） |
| `compose.tls.yaml` | gRPC 传输加密（TLS） |
| `compose.db-cert.yaml` | 数据库证书认证 |
| `compose.security.yaml` | JWT 认证（Token Service） |
| `compose.security-local-key.yaml` | JWT 认证（本地私钥，无 Token Service） |

```bash
# 启用全部安全
docker compose \
  -f compose.yaml \
  -f compose.tls.yaml \
  -f compose.db-cert.yaml \
  -f compose.security.yaml \
  up -d
```

## 11.4 证书生成

```bash
cd build
bash gen-secrets.sh
# → ./secrets/
#   ├── ca.key           CA 私钥（部署前移走保管）
#   ├── private.pem      JWT 签发私钥
#   ├── db/              PostgreSQL 服务端证书
#   ├── controld/        gRPC 服务端证书
#   └── shared/          客户端证书 + JWT 公钥
```

**生产环境 SAN 配置**：

```bash
CONTROLD_SAN="DNS:controld,IP:10.255.129.20" \
DB_SAN="DNS:scalebox-db,IP:10.255.129.20" \
bash gen-secrets.sh
```

**证书认证启用**：客户端需显式设置 `SCALEBOX_CERTS_DIR` 才会使用证书连接数据库，未设置时自动回退到密码认证。

```bash
export SCALEBOX_CERTS_DIR=./secrets/shared
scalebox cluster list
```

## 11.5 Agent 镜像内置凭证

通过 `make copy-secrets` 拷贝到 agent 镜像：

| 文件 | 镜像内路径 | 用途 |
|------|-----------|------|
| `ca.crt` | `/usr/local/etc/certs/ca.crt` | TLS/DB CA 根证书 |
| `client.crt` | `/usr/local/etc/certs/client.crt` | DB 客户端证书 |
| `client.key` | `/usr/local/etc/certs/client.key` | DB 客户端私钥 |
| `public.pem` | `/usr/local/etc/jwt/public.pem` | JWT 验签公钥 |

controld 服务端证书通过 volume 挂载（见 `compose.tls.yaml`），不打包进镜像。

## 11.6 Automation Token

agent / actuator 通过 `t_global.automation_token` 自动获取 JWT token：

```
actuator 启动 → env AUTH_TOKEN 或 db t_global → gRPC stub
    │
    └─ startSlot() → docker run -e AUTH_TOKEN=xxx ...
                           │
                    agent 容器启动 → env AUTH_TOKEN → gRPC stub
```

agent 优先从环境变量 `AUTH_TOKEN` 读取，若为空则回退到数据库查询。

### Token 轮换（滚动更新，零中断）

```bash
# 1. 签发新 token
NEW_TOKEN=$(scalebox-sec jwt sign --sub user-2 --name automation --role automation \
  --ttl 8760h --key /etc/scalebox/jwt/private.pem)

# 2. 写入数据库（旧 token 仍有效直到过期）
docker exec database psql -U scalebox -d scalebox \
  -c "UPDATE t_global SET value = '$NEW_TOKEN' WHERE name = 'automation_token'"

# 3. 滚动重启 actuator
docker compose -f compose.yaml up -d --force-recreate actuator

# 4. 旧 agent 容器自然退出后，所有节点切换为新 token
```

### Token 紧急注销

```bash
docker exec database psql -U scalebox -d scalebox \
  -c "INSERT INTO t_token_blacklist (jti, user_id, expires_at)
      VALUES ('<jti>', <user_id>, now() + interval '30 days')"
```

## 11.7 RBAC 角色体系

### 11.7.1 系统层角色

| 角色 | 权限 | 职责 |
|------|------|------|
| `admin` | 全资源、全操作 | 平台管理员 |
| `viewer` | 全资源 read | 全局只读（审计、报表、监控） |
| `operator` | cluster read、host read/update | 基础设施运维 |
| `automation` | cluster/slot/task/app/event 的 read/update/create/execute | agent/actuator 内部 |

### 11.7.2 应用层权限（DAC）

- **创建者自动获得管理权**：创建 app 的用户（`t_app.user_id`）拥有该 app 下全部操作权限
- **成员授权**：通过 `t_app_member(app_id, user_id, role)` 授予只读权限（`role='viewer'`）

### 11.7.3 细粒度 Scope

`t_role_binding.scope`（JSONB 数组）限定角色作用域。为 NULL 时全局生效：

```json
[{"type":"cluster","id":"cluster-a"}]
[{"type":"cluster","id":"cluster-a"},{"type":"cluster","id":"cluster-b"}]
```

### 11.7.4 资源与操作

| 资源 | 覆盖子资源 | 操作 |
|------|-----------|------|
| `cluster` | host、slot | read、create、update、delete |
| `app` | module、task | read、create、update、delete、start、stop |
| `user` | — | read、create、update、delete |

## 11.8 CLI 管理命令

安全管理通过 `scalebox-sec` CLI（独立项目 [scalebox-security](https://github.com/kaichao/scalebox-security)）：

```bash
# 用户管理
scalebox-sec user create --name alice --password <pwd>
scalebox-sec user list
scalebox-sec user set-status --id 1 --status DISABLED

# 角色绑定
scalebox-sec role bind --user-id 1 --role admin
scalebox-sec role bind --user-id 3 --role operator \
  --scope '[{"type":"cluster","id":"cluster-a"}]'
scalebox-sec role unbind --id 5

# API Key（CI/CD）
scalebox-sec apikey create --user-id 1 --name "Jenkins CI" --expires 720h
scalebox-sec apikey revoke --id 3

# JWT 管理
scalebox-sec jwt genkey -o /etc/scalebox/jwt/
scalebox-sec jwt sign --sub 1 --role admin --ttl 720h
scalebox-sec jwt verify --token <token>
```

> 完整文档见 [scalebox-security](https://github.com/kaichao/scalebox-security/blob/main/docs/README.md) 和 :doc:`安全管理（英文） <../../go-scalebox/docs/security>`。
