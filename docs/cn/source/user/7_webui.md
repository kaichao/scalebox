# 7. WebUI 使用指南

WebUI 是 Scalebox 的统一 Web 管理界面，作为 REST 网关（:8088）连接浏览器和 VS Code 插件，通过 gRPC 与 controld 通信。

## 7.1 架构

```
Browser / VS Code ──HTTP──▶ WebUI (:8088) ──gRPC──▶ controld (:50051)
```

WebUI 进程的角色：
- **REST 网关**：将 HTTP JSON 请求转为 gRPC 调用，76 个端点覆盖全部管理功能
- **SPA 宿主**：内嵌 React 前端（`embed.FS`），单文件部署
- **WebSocket 桥接**：task 日志推流、系统事件通知

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
/resources/clusters/{id}            → 集群详情（Overview / Hosts / Slots）
/resources/hosts                    → 主机全局列表
/resources/slots                    → Slot 色标矩阵
/apps                               → App 卡片网格
/apps/{id}                          → App 概览（统计 + 操作 + 失败列表）
/apps/{id}/dag                      → 模块拓扑 DAG
/apps/{id}/modules                  → 模块列表
/apps/{id}/modules/{name}           → 模块统计（吞吐/错误/利用率/耗时）
/apps/{id}/tasks                    → 任务列表（分页 + 筛选 + header 操作）
/apps/{id}/vtasks                   → VTask 列表（可展开子任务）
/apps/{id}/semaphores               → App 信号量层级树
/apps/{id}/variables                → App 变量层级树
/apps/{id}/settings                 → App 设置（成员管理 + 删除）
/coordination                       → 全局 Semaphore / Variable / Global
/users                              → 用户管理（仅 admin）
/profile                            → 个人中心
```

### 侧栏可见性

| 菜单 | admin | operator | viewer | 开发者 |
|------|:---:|:---:|:---:|:---:|
| 📊 Dashboard | ✅ | ✅ | ✅ | ✅（仅自有） |
| 📦 Apps | ✅ | ❌ | ✅ | ✅（仅自有） |
| 🖥 Resources | ✅ | ✅ | ✅ | ✅ |
| ⚙ Coordination | ✅ | ❌ | ✅ | ✅（自 App） |
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
- **新建 App**：上传 app.yaml + 环境变量，自动模板替换

### 7.4.3 Task 管理

- 分页表格 + 状态/模块筛选
- 点击行展开抽屉：stdout / stderr 文本区 + Headers 编辑表
- 状态色标：等待=蓝，完成=绿，失败=红，运行中=橙

### 7.4.4 VTask 管理

- VTask 列表（可展开行查看子任务）
- Fail 按钮终止未完成的 VTask

### 7.4.5 DAG 拓扑

模块关系有向图（ECharts 力导向图），节点按 vtask_role 着色：
- head = 蓝，core = 绿，tail = 橙

### 7.4.6 Coordination

层级树展示 Semaphore / Variable / Global，支持前缀过滤和叶子节点过滤。

### 7.4.7 用户管理（admin）

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
2. `golang:1.25` — 后端 `go build`（embed 前端产出）
3. `debian:13-slim` — 运行镜像（CGO_ENABLED=0，~27MB）

## 7.6 REST API

全部 76 个 REST 端点（前缀 `/api/v1`），详见 :doc:`REST API 参考 <../../go-scalebox/docs/rest-api-reference>`（英文）。

端点分组：

| 分组 | 端点数 | 说明 |
|------|:-----:|------|
| Cluster | 11 | 集群 CRUD + allocate/release/renew |
| Host | 7 | 主机 CRUD + 状态/续期 |
| Slot | 5 | 插槽 CRUD + 状态 |
| App | 12 | 应用 CRUD + 模块/远程链接/动态扩 slot |
| Task | 5 | 任务查询 + header CRUD |
| VTask | 13 | VTask 查询 + fail + bind/unbind + 作用域变量/信号量 |
| Coordination | 14 | Semaphore(5) + Variable(4) + Global(5) |
| Dashboard | 3 | 系统/App/模块统计 |
| Users | 7 | 用户 CRUD + 角色绑定 |
