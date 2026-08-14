# 7. WebUI User Guide

WebUI is Scalebox's unified web management interface. As a REST gateway (:8088), it connects browsers and the VS Code extension, communicating with controld via gRPC.

## 7.1 Architecture

```
Browser / VS Code ──HTTP──▶ WebUI (:8088) ──gRPC──▶ controld (:50051)
```

Roles of the WebUI process:
- **REST gateway**: converts HTTP JSON requests into gRPC calls; 76 endpoints cover all management functions
- **SPA host**: embeds the React frontend (`embed.FS`), single-file deployment
- **WebSocket bridge**: task log streaming, system event notifications

## 7.2 Starting and Accessing

### 7.2.1 Docker Compose (recommended)

```bash
docker compose -f build/compose.yaml up -d webui
```

Open `http://<host-ip>:8088` in your browser.

### 7.2.2 Local Development

```bash
# Backend
go run -tags dev ./cmd/webui/

# Frontend (hot reload)
cd cmd/webui/frontend && npm run dev
# Access http://localhost:5173
```

### 7.2.3 Configuration

| Variable | Description | Default |
|------|------|--------|
| `GRPC_SERVER` | controld address | `controld:50051` |
| `SERVER_PORT` | WebUI listening port | `8088` |

## 7.3 Navigation Structure

```
/login                              → Login (skipped when security is off)
/                                   → Dashboard (stat cards + trend charts)
/resources/clusters                 → Cluster list
/resources/clusters/{id}            → Cluster details (Overview / Hosts / Slots)
/resources/hosts                    → Global host list
/resources/slots                    → Slot color-coded matrix
/apps                               → App card grid
/apps/{id}                          → App overview (statistics + actions + failure list)
/apps/{id}/dag                      → Module topology DAG
/apps/{id}/modules                  → Module list
/apps/{id}/modules/{name}           → Module statistics (throughput/errors/utilization/duration)
/apps/{id}/tasks                    → Task list (pagination + filtering + header actions)
/apps/{id}/vtasks                   → VTask list (expandable sub-tasks)
/apps/{id}/semaphores               → App semaphore hierarchy tree
/apps/{id}/variables                → App variable hierarchy tree
/apps/{id}/settings                 → App settings (member management + deletion)
/coordination                       → Global Semaphore / Variable / Global
/users                              → User management (admin only)
/profile                            → Personal center
```

### Sidebar Visibility

| Menu | admin | operator | viewer | Developer |
|------|:---:|:---:|:---:|:---:|
| 📊 Dashboard | ✅ | ✅ | ✅ | ✅ (own only) |
| 📦 Apps | ✅ | ❌ | ✅ | ✅ (own only) |
| 🖥 Resources | ✅ | ✅ | ✅ | ✅ |
| ⚙ Coordination | ✅ | ❌ | ✅ | ✅ (own Apps) |
| 👥 Users | ✅ | ❌ | ❌ | ❌ |

## 7.4 Main Features

### 7.4.1 Dashboard

The system overview page, showing:
- Cluster/host/Slot counts
- Number of running Apps and active Tasks
- App-level trend charts and failure lists

### 7.4.2 App Management

- **Card grid**: all Apps are displayed as cards, distinguished by status color (RUNNING green / PAUSED yellow / OFF gray)
- **Quick actions**: Stop / Set Finished / Delete
- **New App**: upload app.yaml + environment variables, automatic template substitution

### 7.4.3 Task Management

- Paginated table + status/module filtering
- Click a row to expand a drawer: stdout / stderr text areas + Headers editing table
- Status colors: waiting = blue, completed = green, failed = red, running = orange

### 7.4.4 VTask Management

- VTask list (expandable rows to view sub-tasks)
- Fail button terminates unfinished VTasks

### 7.4.5 DAG Topology

A directed graph of module relationships (ECharts force-directed graph), with nodes colored by vtask_role:
- head = blue, core = green, tail = orange

### 7.4.6 Coordination

Hierarchy trees displaying Semaphore / Variable / Global, supporting prefix filtering and leaf-node filtering.

### 7.4.7 User Management (admin)

User CRUD + role binding (admin / operator / viewer).

## 7.5 Build and Deployment

### 7.5.1 Single-Binary Build

```bash
cd cmd/webui/frontend && npm run build    # frontend → dist/
cd ../.. && go build ./cmd/webui/         # Go embeds dist/
./webui                                    # start on :8088
```

### 7.5.2 Docker Build

```bash
docker build -f build/webui/Dockerfile -t scalebox/webui:latest .
```

Three Dockerfile stages:
1. `node:22-slim` — frontend `npm ci + npm run build`
2. `golang:1.25` — backend `go build` (embedding the frontend output)
3. `debian:13-slim` — runtime image (CGO_ENABLED=0, ~27MB)

## 7.6 REST API

All 76 REST endpoints (prefix `/api/v1`); see :doc:`REST API Reference <../../go-scalebox/docs/rest-api-reference>` (in English).

Endpoint groups:

| Group | Endpoints | Description |
|------|:-----:|------|
| Cluster | 11 | Cluster CRUD + allocate/release/renew |
| Host | 7 | Host CRUD + status/renewal |
| Slot | 5 | Slot CRUD + status |
| App | 12 | App CRUD + modules/remote links/dynamic slot expansion |
| Task | 5 | Task queries + header CRUD |
| VTask | 13 | VTask queries + fail + bind/unbind + scoped variables/semaphores |
| Coordination | 14 | Semaphore(5) + Variable(4) + Global(5) |
| Dashboard | 3 | System/App/module statistics |
| Users | 7 | User CRUD + role binding |
