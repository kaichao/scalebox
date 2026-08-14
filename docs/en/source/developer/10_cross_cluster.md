# 10. Cross-Cluster Architecture

Scalebox supports unified scheduling and collaborative computing across heterogeneous compute clusters over WANs. Through full-mesh data replication and gRPC proxy forwarding, the CLI always connects only to the local controld, and cross-cluster operations are transparently routed by the server.

## 10.1 Core Principles

```
CLI ──(always)──▶ local controld ──(gRPC proxy)──▶ remote controld
```

- The CLI **always connects only to the local cluster's** controld
- Cross-cluster operations are handled by the local controld, which parses the `cluster` field in the YAML, compares `grpc_server`, and proxies the request
- The `t_cluster` table is replicated in full mesh across all cluster databases (`name='local'` identifies the local cluster), ensuring that every controld can look up any cluster's `grpc_server`

## 10.2 Data Replication Model

| Table | Replication strategy | Description |
|----|:---:|------|
| `t_cluster` | full mesh | every controld holds the definitions of all clusters |
| `t_app` | not replicated | each cluster independently manages its own Apps |
| `t_app_rlink` | cascaded creation | cross-cluster App references (formerly `t_app_remote_link`) |
| `t_module_rvlink` | cascaded creation | cross-cluster module virtual links (formerly `t_module_remote_vlink`) |

## 10.3 Cross-Cluster App Creation

### 10.3.1 Flow

```
CLI: scalebox run (with app.yaml, cluster field pointing to remote)
  │
  ▼
local controld: CreateApp()
  ├── yaml.Unmarshal → app.Cluster
  ├── isRemoteCluster(ctx, app.Cluster)
  │     → compare target cluster grpc_server with 'local' grpc_server
  │
  ├── remote → proxyCreateAppToRemote()
  │           │ dialAppClient(remoteGrpc)
  │           └── remote controld.CreateApp(req)
  │                 └── remote registerCrossCluster()
  │                       ├── establish t_app_rlink on the remote side (main App ↔ new App)
  │                       └── cascade-create t_app_rlink / t_module_rvlink in all clusters
  │
  └── local → create locally (original logic)
```

### 10.3.2 CreateApp Key Parameters

```go
type CreateAppRequest struct {
    YamlContent    string            // app.yaml content
    EnvVars        map[string]string // environment variables (template substitution)
    MainAppId      int32             // main App ID (passed for cross-cluster linkage)
    MainGrpcServer string            // main cluster gRPC address
}
```

### 10.3.3 RegisterRemoteAppLink

After `CreateApp` completes, controld automatically calls `RegisterRemoteAppLink` to establish reference entries across all clusters. Reference entries include:
- Local App ID ↔ remote App ID mappings
- Cross-cluster sink-module routing information (`t_module_rvlink`)

## 10.4 Address Resolution

### 10.4.1 Two-Step Resolution

| Step | Function | Purpose | Query field |
|------|------|------|---------|
| Identity comparison | `resolveGrpcServer` | determine whether the target is remote (`isRemoteCluster`) | `grpc_server` |
| Actual dialing | `resolveReachableAddr` | establish the gRPC connection | `grpc_server` + `remote_grpc_server` |

### 10.4.2 Dialing Decisions

```
Target cluster grpc_server ?= 'local' row grpc_server
  ├─ same → local cluster, return grpc_server
  └─ different → remote cluster
              ├─ remote_grpc_server exists → return remote_grpc_server
              └─ remote_grpc_server absent → return grpc_server (reachable on the same network)
```

Key points:
- Does not rely on CIDR subnet judgment (two clusters may be on the same large intranet but network-isolated)
- Semantics of `remote_grpc_server`: the address specifically provided for other clusters to connect to this cluster (public IP or cross-segment IP)
- `MainGrpcServer` in `registerCrossCluster` is reverse-resolved to a cluster name and re-resolved (because that function runs in the cluster where the new app resides and needs to connect to the main cluster from its own network perspective)

## 10.5 Cross-Cluster Task Routing

### 10.5.1 task add

```
CLI → local controld → doAddRemoteTaskList()
  │
  ├── 1. query t_app_rlink, local app-id → remote app-id
  ├── 2. ignore the local moduleID (t_module IDs are independent on both sides)
  ├── 3. call remote AddTaskList(moduleID=0, appID=remote, sinkModule=as-is)
  └── 4. remote getValidModuleID resolves on its own:
        ├── sinkModule != "" → look up appID by module name
        └── sinkModule == "" → the app's first module
```

Key points:
- **Local moduleID is not passed**: module resolution is entirely the remote side's responsibility
- **sinkModule passed as-is**: no filling or replacement
- The remote side determines local/remote by comparing `grpc_server`, recursively proxying

### 10.5.2 run's stdin start tasks

The CLI always connects to the local controld; `clusterName` is passed as an `AddTaskList` routing parameter.

## 10.6 Connection Targets by Command

| Command domain | Connection target |
|--------|---------|
| all cluster / host / slot | local controld |
| app list / get-router / set-status / add-slots | local controld |
| module list | local controld |
| all vtask | local controld |
| task add | local controld (`--remote-cluster` as routing parameter) |
| task list / get / header | local controld |
| run | local controld (controld proxies to remote internally) |
| semaphore / semagroup / variable / global | local controld |

## 10.7 GRPC_SERVER Resolution Priority

```
1. environment variable GRPC_SERVER (explicitly set by the user, highest priority)
2. t_cluster query for name='local' parameters->>'grpc_server'
3. fallback: 127.0.0.1:50051
```

## 10.8 t_cluster.parameters Key Fields

| Field | Description |
|------|------|
| `grpc_server` | the cluster's gRPC address |
| `remote_grpc_server` | the remotely accessible gRPC address (for cross-cluster dialing) |
| `data_root` | the cluster's data root directory |

> For the detailed design, see :doc:`CLI gRPC connection architecture <../../go-scalebox/docs/cli-grpc-architecture>` (in English).
