# AGENTS.md — ai-agent-hello

本目录是 Scalebox 的 AI agent 最小示例，不是平台仓库。平台规则见：

- 本地：`../../AGENTS.md`
- 线上：https://github.com/kaichao/scalebox/blob/dev/AGENTS.md

## 仓库地图

| 路径 | 内容 |
|---|---|
| `app.yaml` | 单模块应用定义 |
| `scalebox.env` | `CODE_PATH` 默认值 |
| `code/run.sh` | 模块入口；打印 task body，忽略末尾的 headers JSON |
| `.mcp.json` | 只读 MCP server 配置 |
| `README.md` | 人工运行步骤与示例提示词 |

## 运行

```sh
scalebox run --app-file=app.yaml
# 输出形如 {"app_id":1}

echo "hello agent" | scalebox task add --app-id 1 --sink-module hello
scalebox task list --app-id 1
scalebox task show <task-id>
```

平台调用模块的 argv 形式是 `run.sh <task body...> '<headers-json>'`；示例脚本会去掉最后一个 headers 参数，只打印 task body。

`code_path` 来自 `CODE_PATH`。本机在示例目录运行时，`scalebox.env` 默认设为 `${PWD}/code`；如果 CLI 连接远端控制面，必须覆盖为远端节点可见的路径，例如：

```sh
CODE_PATH=/shared/scalebox/examples/ai-agent-hello/code scalebox run --app-file=app.yaml
```

`code_path` 是挂载目录，修改 `run.sh` 后重新投递任务即可；修改 `app.yaml` 中的镜像、参数或环境变量后，必须删除并重建 app。

## 用 MCP 查状态

`.mcp.json` 启动只读 MCP server：

```sh
scalebox mcp serve --read-only
```

推荐查询顺序：

1. `list_apps` 找到 `ai-agent-hello`；
2. `get_app` 取 module id、slot 数和任务计数；
3. `list_tasks` 按应用、模块或状态查任务；
4. `get_task` 看任务详情和 exec 摘要；
5. `get_task_exec` 看单次执行的 stdout/stderr；
6. `list_slots` 查 slot 是否 READY、心跳和主机信息。

MCP 只读，不创建应用、不投递任务、不执行 shell。需要写操作时使用 `scalebox` CLI，并先向人确认。

## 修改规则

- 只改 `code/run.sh`：重新投递任务即可。
- 改 `app.yaml` 或 `scalebox.env`（含 `CODE_PATH`）：app 创建时快照环境变量、镜像与参数，必须重建 app。
- 删除 app、停止容器、修改集群配置前先确认影响范围。
- 投递时 task body 经 stdin 传给 `scalebox task add`，只读一次；模块侧按 argv 接收，末尾附带 headers JSON。

## 平台资料

| 问题 | 位置 |
|---|---|
| 平台心智模型、术语 | `../../AGENTS.md`、`../../docs/agents/00-model.md` |
| 红线与已知坑 | `../../docs/agents/10-invariants.md` |
| 卡链诊断 | `../../docs/agents/20-diagnosis.md` |
| CLI、app.yaml、header | `../../docs/reference/` |
| AI 辅助开发说明 | `../../docs/cn/source/application/9_ai_assisted_development.md` |

## 验收

- `scalebox task show <task-id>` 的 stdout 包含 `hello from ai-agent-hello: <body>`；
- 任务 `status_code=0`；
- MCP `get_task_exec` 能看到同一段 stdout。
