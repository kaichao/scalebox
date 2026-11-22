# 7. WebUI 使用指南

WebUI 是 Scalebox 的统一 Web 管理界面，作为 REST 网关（:8088）连接浏览器和 VS Code 插件，通过 gRPC 与 controld 通信。

## 7.1 架构

```
Browser / VS Code ──HTTP──▶ WebUI (:8088) ──gRPC──▶ controld (:50051)
```

WebUI 进程的角色：
- **REST 网关**：将 HTTP JSON 请求转为 gRPC 调用；109 个 RPC 端点覆盖管理功能，另有 `/auth/me`
- **SPA 宿主**：内嵌 React 前端（`embed.FS`），单文件部署
- **只读 Data View 入口**：执行预定义 SQL 模板，供 App 的 Data tab / Overview 定制视图使用
- **Queue 查看入口**：Coordination → Queues 查看跨集群队列积压并 peek/purge（不替代 CLI 的 push/pull）

## 7.2 启动访问

### 7.2.1 Docker Compose（推荐）

```bash
docker compose -f build/compose.yaml up -d webui
```

浏览器打开 `http://<host-ip>:8088`。

### 7.2.2 本地开发

```bash
# 后端
go run -tags dev ./cmd/webui/

# 前端（热更新）
cd cmd/webui/frontend && npm run dev
# 访问 http://localhost:5173
```

### 7.2.3 配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `GRPC_SERVER` | controld 地址 | `controld:50051` |
| `SERVER_PORT` | WebUI 监听端口 | `8088` |

## 7.3 导航结构

```
/login                              → 登录（安全关闭时跳过）
/                                   → Dashboard（统计卡片 + 趋势图）
/resources/clusters                 → 集群列表
/resources/clusters/{name}          → 集群详情（Overview / Hosts / Slots）
/resources/hosts                    → 主机全局列表
/resources/slots                    → Slot 色标矩阵
/apps                               → App 卡片网格
/apps/new                           → 可视化建 App（Graph / YAML 双 tab）
/apps/{id}                          → App Overview（DAG + 统计 + 定制视图）
/apps/{id}/modules                  → 模块列表
/apps/{id}/modules/{name}           → 模块统计（吞吐/错误/利用率/耗时）
/apps/{id}/tasks                    → 任务列表（分页 + 筛选 + header 操作）
/apps/{id}/vtasks                   → VTask 列表（可展开子任务）
/apps/{id}/vtasks/{vtaskId}         → VTask 详情
/apps/{id}/slots                    → App Slots（生命周期 + 运行态矩阵）
/apps/{id}/slots/{slotId}/execs     → 单个 slot 的 TaskExec 历史
/apps/{id}/coordination             → App 信号量/变量层级树
/apps/{id}/data                     → App Data View（只读 SQL 视图）
/apps/{id}/settings                 → App 设置（Basic Info / 状态 / 成员 / 删除）
/data-views                         → Data View 定义管理（admin 可编辑/预览）
/coordination                       → Semaphore / Variable / Global / Queue
/users                              → 用户管理（仅 admin）
/profile                            → 个人中心
```

旧路径 `/apps/{id}/dag` 会回到 Overview（DAG 已内嵌）；`/apps/{id}/semaphores`、
`/apps/{id}/variables` 兼容到 `/apps/{id}/coordination`。

### 侧栏可见性

| 菜单 | admin | operator | viewer | 开发者 |
|------|:---:|:---:|:---:|:---:|
| 📊 Dashboard | ✅ | ✅ | ✅ | ✅（仅自有） |
| 📦 Apps | ✅ | ❌ | ✅ | ✅（仅自有） |
| 🖥 Resources | ✅ | ✅ | ✅ | ✅ |
| ⚙ Coordination | ✅ | ❌ | ✅ | ✅（自 App） |
| 📋 Data Views | ✅ | ❌ | ✅ | ✅ |
| 👥 Users | ✅ | ❌ | ❌ | ❌ |

## 7.4 主要功能

### 7.4.1 Dashboard

系统总览页，展示：
- 集群/主机/Slot 计数
- 运行中 App 数和活跃 Task 数
- App 级趋势图和失败列表

### 7.4.2 App 管理

- **卡片网格**：所有 App 以卡片展示，按状态颜色区分（RUNNING 绿 / PAUSED 黄 / OFF 灰）
- **快速操作**：Stop / Set Finished / Delete
- **新建 App**：`/apps/new` 支持 **Graph / YAML 双 tab**：YAML 是唯一事实源，图编辑回写 YAML；
  提交前走 `ValidateApp`，创建时快照环境变量、镜像、参数与 slot 展开结果
- **重建提醒**：修改 `app.yaml`、`scalebox.env`、`code_path` 或 `slots` 后必须重建 App

### 7.4.3 App Overview

Overview 是运行时视图，顶部是 DAG，中间是紧凑统计条，底部是 Module / Router / Errors 三个 tab：

- **DAG**：X6 + dagre 分层 + `manhattan` 正交避障路由；支持 physical / logical 两种模式，
  节点位置按模式保存到 `t_app.dag_layout`；admin 可拖动并保存，非 admin 只读
- **统计条**：Running Tasks、Errors(24h)、Recovered(24h)、Active/Total Slots
- **Module Stats**：全量 `v_module_stat`（含 FINISHED/零任务模块）
- **Router Stats**：`v_router_stat`，按主路由模块 + 来源模块聚合排队/执行时间
- **Errors**：当前 task 错误 + `t_task_exec` 曾失败但 task 已恢复/重试中的记录
- **定制视图**：绑定到 App 且 `layout.overview=true` 的 Data View 以卡片/Tab 挂在 Overview；
  支持表格、三维散点、二维 heatmap 伴随视图

### 7.4.4 App Slots（生命周期 + 运行态）

App → Slots 用矩阵展示每个 host 上的 slot，每个单元格 4 项：

```text
[状态 + HB]  [Details]  [TaskExec]  [SlotLog]
```

- **状态**：生命周期 `READY/STARTING/ON/OFF/ERROR/SSHFAILED/TIMEOUT/GREEDY` 用颜色区分；
  `ON` 期间叠加运行态 `RUNNING/ASSIGNED/WAITING/BLOCKED/STALE`；`RUN` 显示当前 task 已运行秒数，
  `HB` 显示最近心跳年龄
- **Details**：slot 生命周期、运行态、等待原因/明细、当前 task/exec、心跳、启动错误、
  **Idle Budget**（`remaining/max × interval ≈ 剩余时间`）；状态切换按钮也在这里
- **TaskExec**：该 slot 的 exec 历史，从新到旧；可继续打开 exec/task 详情
- **SlotLog**：`t_slot_log` 生命周期/心跳事件，从新到旧，分页 100
- **快照语义**：页面不做本地秒数递增；`RUN`/`HB`/空闲预算都冻结在本次取数时刻，表头上方显示
  `Updated YYYY/MM/DD HH:MM:SS`；只有点击 Refresh 或重新进入页面才重新取数
- **STALE 不等于 OFF**：心跳过期但生命周期仍是 ON 时显示 `STALE`；host OFF 的退出/级联见
  [slot 运行态参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/slot-runtime.md)

### 7.4.5 Data View

Data View 是**只读 SQL 查看组件**：

- `/data-views` 管理定义：名称、标题、SQL 模板、参数、列布局、`max_rows`、状态；admin 可
  新建/编辑/删除/预览，viewer 只读
- App → Data tab 展示该 App 的绑定；同一视图可用不同参数绑定多次，用 binding id 寻址
- SQL 只允许单条 `SELECT` / `WITH`，占位符 `{{name}}` 编译为参数绑定；只读事务 + 超时 +
  行数/字节截断 + 敏感表黑名单
- `layout.renderer` 支持 `table` / `scatter3d`；3D 可配 `axisTypes`、`axisRange`、
  `colorMap`、`companion`（二维 heatmap，`avg/count/min`）
- 设计、参数绑定、REST 端点与安全边界见 [Data View 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/data-view.md)

### 7.4.6 Task 管理

- 分页表格 + 状态/模块筛选
- 点击行展开抽屉：stdout / stderr 文本区 + Headers 编辑表
- 状态色标：等待=蓝，完成=绿，失败=红，运行中=橙
- 支持按 `task_id` 深链打开详情；`from_task` header 可用于逐跳反查任务链（当前需按 header 查询）

### 7.4.7 VTask 管理

- VTask 列表（可展开行查看子任务）
- Fail 按钮终止未完成的 VTask

### 7.4.8 Coordination

层级树展示 Semaphore / Variable / Global / Queue：

- Semaphore / Variable / Global 支持前缀过滤和叶子节点过滤
- **Queues**：列出队列名、积压数、最早消息时间；Drawer 内按 `priority DESC, id ASC` 查看
  前 N 条消息（不删除）；admin 可 purge。push/pull 仍用 `scalebox queue`
- Queue 是跨物理集群的管理面消息队列，不是 wait-queue

### 7.4.9 用户管理（admin）

用户 CRUD + 角色绑定（admin / operator / viewer）。

## 7.5 构建部署

### 7.5.1 单二进制构建

```bash
cd cmd/webui/frontend && npm run build    # 前端 → dist/
cd ../.. && go build ./cmd/webui/         # Go 嵌入 dist/
./webui                                    # 启动 :8088
```

### 7.5.2 Docker 构建

```bash
docker build -f build/webui/Dockerfile -t scalebox/webui:latest .
```

Dockerfile 三阶段：
1. `node:22-slim` — 前端 `npm ci + npm run build`
2. `golang:1.27.2` — 后端 `go build`（embed 前端产出）
3. `debian:13-slim` — 运行镜像（CGO_ENABLED=0，~27MB）

## 7.6 REST API

全部 **109 个 RPC REST 端点**（前缀 `/api/v1`），另有 `/auth/me`；结构事实以
`docs/reference/rpc.md` 与 go-scalebox `docs/rest-api-reference.md` 为准。

端点分组：

| 分组 | 端点数 | 说明 |
|------|:-----:|------|
| Cluster | 11 | 集群 CRUD + allocate/release/renew |
| Host | 7 | 主机 CRUD + 状态/续期 |
| Slot | 6 | 插槽 CRUD + 状态 + SlotLog |
| App / Task / Module / VTask | 43 | App CRUD、模块、任务、VTask、shared-body、DAG 布局、错误任务等 |
| Coordination | 16 | Semaphore + Variable + Global + Queue（不含 Semagroup 与 push/pull） |
| Dashboard | 5 | 系统/App/模块/路由统计 |
| Data View | 11 | 定义 CRUD + Preview + App 绑定 + 查询 |
| User | 9 | 用户 CRUD + 角色 + Login/Logout |
| Control | 1 | 仅 AddTaskList（`POST /tasks`） |
| **合计** | **109** | 另有 `/auth/me` |

> 端点数量与分组会随 RPC 变化；`scalebox/docs/reference/rpc.md` 是生成物，源码改动后重新生成。
