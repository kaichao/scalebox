# 10. 跨集群架构

Scalebox 支持跨广域网异构算力集群的统一调度和协同计算。通过全网格数据复制和 gRPC 代理转发，CLI 始终只连本地 controld，跨集群操作由服务端透明路由。

## 10.1 核心原则

```
CLI ──(始终)──▶ 本地 controld ──(gRPC 代理)──▶ 远端 controld
```

- CLI **永远只连接本地集群**的 controld
- 跨集群操作由本地 controld 解析 YAML 中的 `cluster` 字段，比对 `grpc_server` 后代理转发
- `t_cluster` 表在所有集群数据库间全网格复制（`name='local'` 标识本地集群），保证所有 controld 都能查到任意集群的 `grpc_server`

## 10.2 数据复制模型

| 表 | 复制策略 | 说明 |
|----|:---:|------|
| `t_cluster` | 全网格 | 每个 controld 持有全部集群的定义 |
| `t_app` | 不复制的 | 各集群独立管理自己的 App |
| `t_app_rlink` | 级联创建 | 跨集群 App 引用（原 `t_app_remote_link`） |
| `t_module_rvlink` | 级联创建 | 跨集群模块虚拟链接（原 `t_module_remote_vlink`） |

## 10.3 跨集群 App 创建

### 10.3.1 流程

```
CLI: scalebox run（含 app.yaml，cluster 字段指向远端）
  │
  ▼
本地 controld: CreateApp()
  ├── yaml.Unmarshal → app.Cluster
  ├── isRemoteCluster(ctx, app.Cluster)
  │     → 比较目标集群 grpc_server 与 'local' 的 grpc_server
  │
  ├── 远端 → proxyCreateAppToRemote()
  │           │ dialAppClient(remoteGrpc)
  │           └── 远端 controld.CreateApp(req)
  │                 └── 远端 registerCrossCluster()
  │                       ├── 在远端建立 t_app_rlink（主 App ↔ 新建 App）
  │                       └── 在所有集群级联创建 t_app_rlink / t_module_rvlink
  │
  └── 本地 → 本地创建（原逻辑）
```

### 10.3.2 CreateApp 关键参数

```go
type CreateAppRequest struct {
    YamlContent    string            // app.yaml 内容
    EnvVars        map[string]string // 环境变量（模板替换）
    MainAppId      int32             // 主 App ID（跨集群联动时传入）
    MainGrpcServer string            // 主集群 gRPC 地址
}
```

### 10.3.3 RegisterRemoteAppLink

`CreateApp` 完成后，controld 自动调用 `RegisterRemoteAppLink` 在所有集群间建立引用条目。引用条目包含：
- 本地 App ID ↔ 远端 App ID 映射
- 跨集群 sink-module 路由信息（`t_module_rvlink`）

## 10.4 地址解析

### 10.4.1 两步解析

| 步骤 | 函数 | 用途 | 查询字段 |
|------|------|------|---------|
| 身份比较 | `resolveGrpcServer` | 判断目标是否远端（`isRemoteCluster`） | `grpc_server` |
| 实际拨号 | `resolveReachableAddr` | 建立 gRPC 连接 | `grpc_server` + `remote_grpc_server` |

### 10.4.2 拨号决策

```
目标集群 grpc_server ?= 'local' 行 grpc_server
  ├─ 相同 → 本地集群，返回 grpc_server
  └─ 不同 → 远端集群
              ├─ remote_grpc_server 存在 → 返回 remote_grpc_server
              └─ remote_grpc_server 不存在 → 返回 grpc_server（同网络可达）
```

关键点：
- 不依赖 CIDR 子网判断（两个集群可能同大内网但网络隔离）
- `remote_grpc_server` 的语义：专门给其他集群连接本集群使用的地址（公网 IP 或跨网段 IP）
- `registerCrossCluster` 中的 `MainGrpcServer` 会被反查集群名后重解析（因为该函数运行在新 app 所在集群，需要从自身网络视角连主集群）

## 10.5 跨集群 Task 路由

### 10.5.1 task add

```
CLI → 本地 controld → doAddRemoteTaskList()
  │
  ├── 1. 查 t_app_rlink，本地 app-id → 远端 app-id
  ├── 2. 忽略本地 moduleID（两端 t_module ID 各自独立）
  ├── 3. 调用远端 AddTaskList(moduleID=0, appID=远端, sinkModule=原样)
  └── 4. 远端 getValidModuleID 自行解析：
        ├── sinkModule != "" → 按模块名查 appID
        └── sinkModule == "" → app 第一个模块
```

关键点：
- **不传本地 moduleID**：模块解析完全由远端负责
- **sinkModule 原样传递**：不做任何填充或替换
- 远端通过 `grpc_server` 比较判断本地/远端，递归代理

### 10.5.2 run 的 stdin start tasks

CLI 始终连本地 controld，`clusterName` 作为 `AddTaskList` 路由参数传入。

## 10.6 各命令连接目标

| 命令域 | 连接目标 |
|--------|---------|
| cluster / host / slot 全部 | 本地 controld |
| app list / get-router / set-status / add-slots | 本地 controld |
| module list | 本地 controld |
| vtask 全部 | 本地 controld |
| task add | 本地 controld（`--remote-cluster` 作路由参数） |
| task list / get / header | 本地 controld |
| run | 本地 controld（controld 内部代理到远端） |
| semaphore / semagroup / variable / global | 本地 controld |

## 10.7 GRPC_SERVER 解析优先级

```
1. 环境变量 GRPC_SERVER（用户显式设置，最高优先级）
2. t_cluster 查询 name='local' 的 parameters->>'grpc_server'
3. 回退：127.0.0.1:50051
```

## 10.8 t_cluster.parameters 关键字段

| 字段 | 说明 |
|------|------|
| `grpc_server` | 集群 gRPC 地址 |
| `remote_grpc_server` | 远端可访问的 gRPC 地址（供跨集群拨号使用） |
| `data_root` | 集群数据根目录 |

> 详细设计见 :doc:`CLI gRPC 连接架构 <../../go-scalebox/docs/cli-grpc-architecture>`（英文）。
