# 8. VS Code 插件

Scalebox VS Code 插件提供侧边栏 TreeView、任务日志查看、DAG 拓扑图和 app.yaml 语法支持，通过 HTTP 调用 WebUI REST API。

## 8.1 安装

### 8.1.1 从 .vsix 文件安装

```bash
# 构建 .vsix（需要 Docker）
cd vscode-scalebox && make package

# 安装
code --install-extension vscode-scalebox-0.1.0.vsix
```

或通过 VS Code 界面：Extensions → `...` → Install from VSIX。

### 8.1.2 配置

安装后在 VS Code 设置中配置：

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `scalebox.serverUrl` | WebUI 地址 | `http://localhost:8088` |
| `scalebox.clusterName` | 默认集群名 | 空（自动选择） |
| `scalebox.authToken` | JWT 认证令牌 | 空 |

安全环境需设置 `scalebox.authToken`，令牌通过 WebUI 登录获取。

## 8.2 功能

### 8.2.1 侧边栏 TreeView

打开 Scalebox 视图（Activity Bar → Scalebox 图标），显示三级懒加载树：

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

状态图标：
- App：🟢 RUNNING / 🟡 PAUSED / 🔴 FAILED / ⚪ OFF
- Task：⏳ READY（-1）/ 🔵 QUEUED（-2）/ 🟠 RUNNING（-3）/ ✅ 完成（0）/ ❌ 失败（>0）
- VTask：⏳ / 🔵 / ✅ / ❌（派生状态 + 子任务进度）

右键菜单操作：Stop App / Set App Status / New Task / View Task Logs / View DAG / Copy App ID / Open in Browser。

### 8.2.2 命令面板

`Ctrl+Shift+P` → 输入 `Scalebox`：

| 命令 | 说明 |
|------|------|
| Scalebox: Refresh Apps | 刷新 TreeView |
| Scalebox: Open in Browser | 在浏览器打开 WebUI |
| Scalebox: View App Detail | 查看 App 概览 |
| Scalebox: Copy App ID | 复制 App ID 到剪贴板 |
| Scalebox: New Task | 向 App 添加新 Task |
| Scalebox: View Task Logs | 查看 Task 日志 |
| Scalebox: Stop App | 停止 App |
| Scalebox: Set App Status | 修改 App 状态 |
| Scalebox: View DAG | 打开 DAG 拓扑图 |
| Scalebox: Select Cluster | 切换集群 |

### 8.2.3 状态栏

| 状态栏项 | 说明 |
|---------|------|
| 🖥 Cluster | 当前集群名（点击切换） |
| 📦 Running: N | 运行中 App 数量 |
| ❌ Failed: N | 失败 Task 数量 |

30 秒自动刷新，也可手动点击刷新。

### 8.2.4 Task 日志查看

点击 Task 节点或右键 → View Task Logs，打开 Webview Panel：
- **stdout**：绿色背景代码块
- **stderr**：红色背景代码块
- 刷新 / 复制 / WebUI 跳转按钮

### 8.2.5 DAG 拓扑图

右键 App → View DAG，打开 Webview + ECharts 力导向图：
- 节点 = 模块，按 vtask_role 着色（head 蓝 / core 绿 / tail 橙）
- 边 = 数据流方向
- controld 未实现 `GetAppDAG` 时自动降级为模块拓扑

### 8.2.6 app.yaml LSP

编辑 `app*.yaml` 文件时自动激活：
- **Diagnostics**：必填字段检查、值校验（vtask_role 枚举、slot 正则）、YAML 解析错误
- **Completions**：输入 `:` 触发，顶层 key + 模块级 key + vtask_role 值上下文补全
- **Hover**：悬停显示字段说明、类型、是否必填

Schema 对齐 `cmd/controld/grpc/app_types.go`，14 个模块字段完整覆盖。

## 8.3 构建

### 8.3.1 本地开发

```bash
cd vscode-scalebox
npm install
npx tsc -p ./
```

按 F5 启动 Extension Development Host。

### 8.3.2 Docker 打包

```bash
cd vscode-scalebox
make package      # 生成 .vsix
make install      # 打包 + 安装
make clean        # 清理
```

Makefile 使用 `node:22-slim` 镜像，无需本地 Node.js。

## 8.4 目录结构

```
vscode-scalebox/
├── package.json          # 扩展清单（activationEvents, contributes）
├── tsconfig.json         # TypeScript 配置
├── Makefile              # Docker 构建 + .vsix 打包
├── src/
│   ├── extension.ts      # activate() 入口
│   ├── api/
│   │   ├── client.ts     # HTTP 客户端（fetch + JWT）
│   │   └── types.ts      # TypeScript 类型定义
│   ├── tree/
│   │   └── appTree.ts    # TreeDataProvider 实现
│   ├── commands/
│   │   └── index.ts      # 10 个命令实现
│   ├── statusBar.ts      # 状态栏项
│   └── lsp/
│       └── appYaml.ts    # app.yaml 校验/补全/悬停
```

> 完整设计文档见 :doc:`WebUI / VS Code 设计方案 <../../go-scalebox/docs/webui-vscode-plan>`（英文）。
