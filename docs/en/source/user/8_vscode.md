# 8. VS Code Extension

The Scalebox VS Code extension provides a sidebar TreeView, task log viewing, DAG topology, and app.yaml syntax support, calling the WebUI REST API over HTTP.

## 8.1 Installation

### 8.1.1 Installing from a .vsix File

```bash
# Build the .vsix (requires Docker)
cd vscode-scalebox && make package

# Install
code --install-extension vscode-scalebox-0.1.0.vsix
```

Or through the VS Code UI: Extensions → `...` → Install from VSIX.

### 8.1.2 Configuration

After installation, configure in VS Code settings:

| Setting | Description | Default |
|--------|------|--------|
| `scalebox.serverUrl` | WebUI address | `http://localhost:8088` |
| `scalebox.clusterName` | Default cluster name | empty (auto-selected) |
| `scalebox.authToken` | JWT authentication token | empty |

In secure environments, set `scalebox.authToken`; the token is obtained by logging in to the WebUI.

## 8.2 Features

### 8.2.1 Sidebar TreeView

Open the Scalebox view (Activity Bar → Scalebox icon) to display a three-level lazy-loading tree:

```
📦 App 1 (RUNNING) 🟢
  ├── 📁 Modules (3)
  │   ├── 🔧 main-router
  │   ├── 🔧 vtask-head
  │   └── 🔧 vtask-core
  ├── 📁 Tasks (12)
  │   ├── ⏳ task-001
  │   ├── ✅ task-002
  │   └── ❌ task-003
  └── 📁 VTasks (2)
      ├── ⏳ vtask-42 (2/3)
      └── ✅ vtask-41
```

Status icons:
- App: 🟢 RUNNING / 🟡 PAUSED / 🔴 FAILED / ⚪ OFF
- Task: ⏳ READY (-1) / 🔵 QUEUED (-2) / 🟠 RUNNING (-3) / ✅ completed (0) / ❌ failed (>0)
- VTask: ⏳ / 🔵 / ✅ / ❌ (derived status + sub-task progress)

Right-click menu actions: Stop App / Set App Status / New Task / View Task Logs / View DAG / Copy App ID / Open in Browser.

### 8.2.2 Command Palette

`Ctrl+Shift+P` → type `Scalebox`:

| Command | Description |
|------|------|
| Scalebox: Refresh Apps | Refresh the TreeView |
| Scalebox: Open in Browser | Open the WebUI in a browser |
| Scalebox: View App Detail | View the App overview |
| Scalebox: Copy App ID | Copy the App ID to the clipboard |
| Scalebox: New Task | Add a new Task to the App |
| Scalebox: View Task Logs | View Task logs |
| Scalebox: Stop App | Stop the App |
| Scalebox: Set App Status | Change the App status |
| Scalebox: View DAG | Open the DAG topology |
| Scalebox: Select Cluster | Switch clusters |

### 8.2.3 Status Bar

| Status bar item | Description |
|---------|------|
| 🖥 Cluster | Current cluster name (click to switch) |
| 📦 Running: N | Number of running Apps |
| ❌ Failed: N | Number of failed Tasks |

Auto-refreshes every 30 seconds; can also be refreshed manually by clicking.

### 8.2.4 Task Log Viewing

Click a Task node or right-click → View Task Logs to open a Webview Panel:
- **stdout**: code block with a green background
- **stderr**: code block with a red background
- Refresh / Copy / WebUI jump buttons

### 8.2.5 DAG Topology

Right-click an App → View DAG to open a Webview + ECharts force-directed graph:
- Nodes = modules, colored by vtask_role (head blue / core green / tail orange)
- Edges = data flow direction
- Automatically degrades to module topology when controld does not implement `GetAppDAG`

### 8.2.6 app.yaml LSP

Automatically activated when editing `app*.yaml` files:
- **Diagnostics**: required field checks, value validation (vtask_role enums, slot regex), YAML parse errors
- **Completions**: triggered by typing `:`, top-level keys + module-level keys + vtask_role value contextual completion
- **Hover**: hover to show field descriptions, types, and whether required

Schema aligned with `cmd/controld/grpc/app_types.go`, covering all 14 module fields.

## 8.3 Building

### 8.3.1 Local Development

```bash
cd vscode-scalebox
npm install
npx tsc -p ./
```

Press F5 to start the Extension Development Host.

### 8.3.2 Docker Packaging

```bash
cd vscode-scalebox
make package      # generate .vsix
make install      # package + install
make clean        # clean up
```

The Makefile uses the `node:22-slim` image; no local Node.js required.

## 8.4 Directory Structure

```
vscode-scalebox/
├── package.json          # extension manifest (activationEvents, contributes)
├── tsconfig.json         # TypeScript configuration
├── Makefile              # Docker build + .vsix packaging
├── src/
│   ├── extension.ts      # activate() entry point
│   ├── api/
│   │   ├── client.ts     # HTTP client (fetch + JWT)
│   │   └── types.ts      # TypeScript type definitions
│   ├── tree/
│   │   └── appTree.ts    # TreeDataProvider implementation
│   ├── commands/
│   │   └── index.ts      # 10 command implementations
│   ├── statusBar.ts      # status bar items
│   └── lsp/
│       └── appYaml.ts    # app.yaml validation/completion/hover
```

> See :doc:`WebUI / VS Code design plan <../../go-scalebox/docs/webui-vscode-plan>` for the complete design document (in English).
