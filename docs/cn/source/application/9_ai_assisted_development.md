# 9. AI 辅助开发

本章说明如何让 AI 编程工具（Claude Code、Codex、Cursor 等）在应用仓中开发 Scalebox 应用。
两种接法可以并用：AI 读取应用仓的 `AGENTS.md`，通过 `scalebox` CLI 执行操作；也可以让
AI 通过只读 MCP server 查询应用、任务与 slot 状态。

## 9.1 适用范围

适合：

- 编写和修改 `app.yaml`、模块脚本、主路由代码；
- 排查任务链、任务状态与模块接口；
- 为新应用生成初始目录、配置和验证命令。

不适合：

- 代替平台安装与集群部署，见[安装部署](../deployment/1_installation.md)；
- 代替人工审查生产变更、集群凭据和权限配置。

使用前需要：

- AI 工具能读取应用仓；
- AI 工具能读取 Scalebox 仓（本地路径或 GitHub）；
- `scalebox` CLI 在 `PATH` 中，控制面已启动。

## 9.2 给应用仓加 AGENTS.md

在应用仓根目录新建 `AGENTS.md`。它是跨工具入口，只写本应用的事实；平台规则不复制，
统一指向 Scalebox 仓。

最小模板：

````markdown
# AGENTS.md — <应用名>

本仓库是基于 Scalebox 的应用，不是平台仓库。平台规则见：
- 本地：`<scalebox 仓路径>/AGENTS.md`
- 线上：https://github.com/kaichao/scalebox/blob/dev/AGENTS.md

## 仓库地图

| 路径 | 内容 |
|---|---|
| `app.yaml` | 应用定义 |
| `*.env` | 站点与环境变量 |
| `modules/` | 本应用模块 |
| `docs/` | 本应用业务说明 |

## 运行

```sh
scalebox run --env-file=<env-file> --app-file=<app.yaml> < /dev/null
scalebox app list
scalebox task list --app-id <id>
scalebox task show <task-id>
```

## 本地约定

- <目录布局 / 数据格式 / 命名规则>
- <本应用特有的坑>

## 平台资料

| 问题 | 位置 |
|---|---|
| 心智模型与术语 | scalebox `AGENTS.md`、`docs/agents/00-model.md` |
| 红线与坑 | `docs/agents/10-invariants.md` |
| 卡链诊断 | `docs/agents/20-diagnosis.md` |
| 建 app / 加模块 / 流控 / vtask | `docs/agents/recipes/` |
| CLI / 字段 / header | `docs/reference/` |
| 可执行特性 | `features/` |

## 知识回流

平台级机制按 scalebox `docs/agents/contributing.md` 写回平台仓；
本应用私有约定留在本仓库。
````

写作边界：

- **只写本应用**：流水线拓扑、目录布局、环境文件、业务命名、已验证的坑。
- **平台规则只放链接**：红线、字段、命令、诊断步骤都指向 Scalebox 仓；复制一份必然过期。
- **命令要能直接执行**：写清工作目录、env 文件和入口模块，避免只写“运行应用”。
- **状态要标注日期**：实验结论、已知缺陷、待办要有日期或提交号，AI 不会把旧状态当现状。

## 9.3 接入 AI 工具

### 9.3.1 Claude Code

Claude Code 读取 `CLAUDE.md`。如果应用仓已有 `CLAUDE.md`，在文件顶部加一行导入：

```markdown
@AGENTS.md
```

如果已有详细的 `CLAUDE.md`，也可以保留它，并在其中指向 `AGENTS.md`；不要让两份内容各自维护。
`AGENTS.md` 负责跨工具入口，`CLAUDE.md` 可以继续放本仓库的详细工作说明。

### 9.3.2 其他工具

Codex、Cursor 等工具若支持 `AGENTS.md`，会直接读取；不支持时，在工具的规则文件中写一行
“先读仓库根 `AGENTS.md`”，不要另建一套规则。

### 9.3.3 Scalebox MCP server（只读）

Scalebox 自带一个只读 MCP server，向支持 MCP 的客户端提供平台状态查询：

```sh
scalebox mcp serve --read-only
```

Claude Code 的 `.mcp.json` 配置示例：

```json
{
  "mcpServers": {
    "scalebox": {
      "command": "scalebox",
      "args": ["mcp", "serve", "--read-only"]
    }
  }
}
```

当前提供 6 个工具：

| 工具 | 用途 |
|---|---|
| `list_apps` | 查找应用 |
| `get_app` | 应用、模块、slot 数与任务计数 |
| `list_tasks` | 按应用、模块、状态、时间查任务 |
| `get_task` | 任务详情与执行记录摘要 |
| `get_task_exec` | 单次执行的 stdout/stderr 等日志 |
| `list_slots` | slot 状态、心跳与错误定位 |

MCP server 只读，不创建应用、不投递任务、不执行 shell。它提供动态平台状态；
`AGENTS.md` 仍负责静态项目约定，两者互补。MCP server 使用当前 CLI 的集群配置和环境文件；
在应用目录启动时会读取该目录的 `scalebox.env`，也可用全局 `--env-file` 指定。

### 9.3.4 可选权限配置

在应用仓的 `.claude/settings.local.json` 中放行常用命令，可以减少逐条确认：

```json
{
  "permissions": {
    "allow": [
      "Bash(scalebox *)",
      "Read(//path/to/scalebox/**)"
    ]
  }
}
```

`Read` 的路径按本机实际 Scalebox 仓填写。不要把令牌、密码或个人凭据写进
`AGENTS.md`、`CLAUDE.md` 或权限文件；环境文件如有敏感值，保持在版本控制之外。

## 9.4 平台资料入口

AI 在应用仓工作时，按任务查以下资料：

| 要做的事 | 先读 |
|---|---|
| 建立平台整体模型 | [AGENTS.md](https://github.com/kaichao/scalebox/blob/dev/AGENTS.md)、[00-model.md](https://github.com/kaichao/scalebox/blob/dev/docs/agents/00-model.md) |
| 确认平台红线与已知坑 | [10-invariants.md](https://github.com/kaichao/scalebox/blob/dev/docs/agents/10-invariants.md) |
| 任务卡住或失败 | [20-diagnosis.md](https://github.com/kaichao/scalebox/blob/dev/docs/agents/20-diagnosis.md) |
| 建 app、加模块、配置流控、VTask、多节点 | [recipes/](https://github.com/kaichao/scalebox/blob/dev/docs/agents/recipes/README.md) |
| 查 CLI、字段、header、信号量 | [docs/reference/](https://github.com/kaichao/scalebox/blob/dev/docs/reference/) |
| 找可运行的最小特性与测试 | [features/](https://github.com/kaichao/scalebox/blob/dev/features/README.md) |
| 把新发现写回平台 | [contributing.md](https://github.com/kaichao/scalebox/blob/dev/docs/agents/contributing.md) |

站内先读[核心概念](1_core_concepts.md)、[运行与管理应用](2_running_apps.md)、
[模块开发](4_module_development.md)；完整示例见[示例应用](5_example_apps.md)。

## 9.5 一次典型工作流

1. 在应用仓启动 AI 工具，给出目标和边界。
2. AI 先读 `AGENTS.md`、平台不变式与现有 `app.yaml`，再提出改动清单。
3. 人工确认后修改配置或代码；涉及镜像、环境变量或参数时重建 app。
4. 用 CLI 验证：`scalebox run` → `scalebox task add` → `scalebox task list` →
   `scalebox task show`。
5. 平台级发现按 `docs/agents/contributing.md` 回写；应用私有约定留在应用仓。

给 AI 的任务描述可以写成：

```text
先读 AGENTS.md、docs/agents/10-invariants.md 和现有 app.yaml。
为当前应用增加一个模块，先给出改动清单、影响范围和验证命令；
不要直接执行 scalebox run，等我确认。
```

## 9.6 边界与验收

AI 可以生成配置、代码和命令，但不能替代以下判断：

- 生产环境的删除、停机和权限变更；
- 集群凭据、密钥和数据访问边界；
- 平台不变式的最终确认。

验收标准：

- AI 能说出某个操作该查哪份平台文档，而不是凭记忆编命令；
- 修改后能用 `scalebox task show` 看到预期 stdout 与 `status_code=0`；
- 新发现的平台级问题进入 `docs/agents/`、`features/` 或 `docs/reference/`，不是只留在应用仓。
