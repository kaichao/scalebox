# 11. Security Framework

The Scalebox security framework provides authentication and authorization capabilities, securing the gRPC channel through JWT + RBAC + TLS. Security is off by default (`SECURITY_ENABLED=false`); when enabled, there is no extra runtime overhead beyond the performance cost.

## 11.1 Communication Security Matrix

| | controld | database | agent | actuator | CLI |
|---|:---:|:---:|:---:|:---:|:---:|
| **controld** | — | PG client | gRPC server | gRPC server | gRPC server |
| **database** | PG server | — | — | — | — |
| **agent** | gRPC client | — | — | — | — |
| **actuator** | gRPC client | — | — | — | — |
| **CLI** | gRPC client | PG client | — | — | — |

Security credentials carried by each module:

| Module | Role | Needs |
|------|------|------|
| controld | gRPC server + PG client | server TLS cert, JWT public key (signature verification), DB client cert |
| agent | gRPC client | CA cert (server verification), JWT token |
| actuator | gRPC client + PG client | CA cert, JWT token, DB client cert |
| CLI | gRPC client + PG client | CA cert, JWT token, DB client cert |
| database | PG server | DB server cert (PostgreSQL SSL) |

All certificates + JWT keys are generated uniformly by `gen-secrets.sh`, signed by the same internal CA.

## 11.2 Security Boundaries

- **Actuator SSH**: Ed25519 passwordless SSH login to compute nodes, belonging to infrastructure-layer security
- **Direct database connections**: operations that bypass gRPC and access the database directly are not subject to RBAC control
- **RBAC covers only the gRPC channel**: authentication and authorization take effect at the gRPC interceptor layer

## 11.3 Enabling

### 11.3.1 Environment Variables

| Variable | Description | Default |
|------|------|--------|
| `SECURITY_ENABLED` | security master switch | `false` |
| `GRPC_TLS_ENABLED` | gRPC TLS transport encryption | `false` |
| `GRPC_TLS_CA_FILE` | CA certificate path | — |
| `GRPC_TLS_CERT_FILE` | server certificate path | — |
| `GRPC_TLS_KEY_FILE` | server private key path | — |
| `AUTH_TOKEN` | JWT authentication token | — |
| `SCALEBOX_CERTS_DIR` | certificate directory (enables DB certificate authentication when set) | — |

### 11.3.2 Compose Layers

| File | Function |
|------|------|
| `compose.yaml` | base services (no security) |
| `compose.tls.yaml` | gRPC transport encryption (TLS) |
| `compose.db-cert.yaml` | database certificate authentication |
| `compose.security.yaml` | JWT authentication (Token Service) |
| `compose.security-local-key.yaml` | JWT authentication (local private key, no Token Service) |

```bash
# Enable full security
docker compose \
  -f compose.yaml \
  -f compose.tls.yaml \
  -f compose.db-cert.yaml \
  -f compose.security.yaml \
  up -d
```

## 11.4 Certificate Generation

```bash
cd build
bash gen-secrets.sh
# → ./secrets/
#   ├── ca.key           CA private key (move it away for safekeeping before deployment)
#   ├── private.pem      JWT signing private key
#   ├── db/              PostgreSQL server certificates
#   ├── controld/        gRPC server certificates
#   └── shared/          client certificates + JWT public key
```

**Production SAN configuration**:

```bash
CONTROLD_SAN="DNS:controld,IP:10.255.129.20" \
DB_SAN="DNS:scalebox-db,IP:10.255.129.20" \
bash gen-secrets.sh
```

**Enabling certificate authentication**: clients must explicitly set `SCALEBOX_CERTS_DIR` to use certificate connections to the database; when unset, they automatically fall back to password authentication.

```bash
export SCALEBOX_CERTS_DIR=./secrets/shared
scalebox cluster list
```

## 11.5 Credentials Embedded in the Agent Image

Copied to the agent image via `make copy-secrets`:

| File | Path in image | Purpose |
|------|-----------|------|
| `ca.crt` | `/usr/local/etc/certs/ca.crt` | TLS/DB CA root certificate |
| `client.crt` | `/usr/local/etc/certs/client.crt` | DB client certificate |
| `client.key` | `/usr/local/etc/certs/client.key` | DB client private key |
| `public.pem` | `/usr/local/etc/jwt/public.pem` | JWT verification public key |

The controld server certificate is mounted via volume (see `compose.tls.yaml`), not packaged into the image.

## 11.6 Automation Token

agent / actuator automatically obtain JWT tokens through `t_global.automation_token`:

```
actuator startup → env AUTH_TOKEN or db t_global → gRPC stub
    │
    └─ startSlot() → docker run -e AUTH_TOKEN=xxx ...
                           │
                    agent container startup → env AUTH_TOKEN → gRPC stub
```

The agent first reads from the `AUTH_TOKEN` environment variable; if empty, it falls back to a database query.

### Token Rotation (rolling update, zero downtime)

```bash
# 1. Issue a new token
NEW_TOKEN=$(scalebox-sec jwt sign --sub user-2 --name automation --role automation \
  --ttl 8760h --key /etc/scalebox/jwt/private.pem)

# 2. Write to the database (the old token remains valid until expiry)
docker exec database psql -U scalebox -d scalebox \
  -c "UPDATE t_global SET value = '$NEW_TOKEN' WHERE name = 'automation_token'"

# 3. Rolling-restart the actuator
docker compose -f compose.yaml up -d --force-recreate actuator

# 4. After old agent containers exit naturally, all nodes switch to the new token
```

### Emergency Token Revocation

```bash
docker exec database psql -U scalebox -d scalebox \
  -c "INSERT INTO t_token_blacklist (jti, user_id, expires_at)
      VALUES ('<jti>', <user_id>, now() + interval '30 days')"
```

## 11.7 RBAC Role System

### 11.7.1 System-Level Roles

| Role | Permissions | Responsibility |
|------|------|------|
| `admin` | all resources, all operations | platform administrator |
| `viewer` | read on all resources | global read-only (audit, reporting, monitoring) |
| `operator` | cluster read, host read/update | infrastructure operations |
| `automation` | read/update/create/execute on cluster/slot/task/app/event | agent/actuator internal |

### 11.7.2 App-Level Permissions (DAC)

- **Creators automatically gain management rights**: the user who creates the app (`t_app.user_id`) has all operation permissions under that app
- **Member authorization**: read-only permissions granted through `t_app_member(app_id, user_id, role)` (`role='viewer'`)

### 11.7.3 Fine-Grained Scope

`t_role_binding.scope` (JSONB array) limits the role's scope. When NULL, it applies globally:

```json
[{"type":"cluster","id":"cluster-a"}]
[{"type":"cluster","id":"cluster-a"},{"type":"cluster","id":"cluster-b"}]
```

### 11.7.4 Resources and Operations

| Resource | Covers sub-resources | Operations |
|------|-----------|------|
| `cluster` | host, slot | read, create, update, delete |
| `app` | module, task | read, create, update, delete, start, stop |
| `user` | — | read, create, update, delete |

## 11.8 CLI Management Commands

Security management is done through the `scalebox-sec` CLI (standalone project [scalebox-security](https://github.com/kaichao/scalebox-security)):

```bash
# User management
scalebox-sec user create --name alice --password <pwd>
scalebox-sec user list
scalebox-sec user set-status --id 1 --status DISABLED

# Role binding
scalebox-sec role bind --user-id 1 --role admin
scalebox-sec role bind --user-id 3 --role operator \
  --scope '[{"type":"cluster","id":"cluster-a"}]'
scalebox-sec role unbind --id 5

# API Keys (CI/CD)
scalebox-sec apikey create --user-id 1 --name "Jenkins CI" --expires 720h
scalebox-sec apikey revoke --id 3

# JWT management
scalebox-sec jwt genkey -o /etc/scalebox/jwt/
scalebox-sec jwt sign --sub 1 --role admin --ttl 720h
scalebox-sec jwt verify --token <token>
```

> See [scalebox-security](https://github.com/kaichao/scalebox-security/blob/main/docs/README.md) for the complete documentation, and :doc:`Security Management (in English) <../../go-scalebox/docs/security>`.
