# ai-agent-hello

给 AI agent 使用的最小 Scalebox 应用示例：一个模块、一个任务，附带 `AGENTS.md`、`CLAUDE.md` 和只读 MCP 配置。

## 目录

```text
ai-agent-hello/
├── AGENTS.md       # 跨工具入口：运行、MCP 查询、修改边界
├── CLAUDE.md       # Claude Code 导入 AGENTS.md
├── .mcp.json       # 只读 MCP server 配置
├── app.yaml        # 单模块应用
├── scalebox.env    # CODE_PATH 默认值
└── code/run.sh     # 打印 task body，忽略 headers
```

## 前置条件

- Scalebox 控制面已启动；
- `scalebox` CLI 已安装，且支持 `scalebox mcp serve`；
- 当前目录的 `cluster: local` 与本地控制面匹配；其他集群请修改 `app.yaml`；
- `scalebox.env` 默认把 `CODE_PATH` 指向当前目录的 `code/`；连接远端控制面时需覆盖，见下文。

## 手动运行

```sh
cd examples/ai-agent-hello
echo "hello agent" | scalebox run --app-file=app.yaml
# 输出形如 {"app_id":1}

scalebox task list --app-id 1
scalebox task show <task-id>
```

预期 stdout：

```text
hello from ai-agent-hello: hello agent
```

平台调用模块时会在参数末尾附加 headers JSON；`run.sh` 会去掉最后一个参数，只打印 task body。

### CODE_PATH

`app.yaml` 使用 `code_path: ${CODE_PATH}`，默认值在 `scalebox.env`：

```sh
CODE_PATH=${PWD}/code
```

如果 `scalebox` CLI 连接的是远端控制面，`CODE_PATH` 必须是远端计算节点可见的路径。例如代码已同步到 `/shared/scalebox/examples/ai-agent-hello/` 时：

```sh
export CODE_PATH=/shared/scalebox/examples/ai-agent-hello/code
echo "hello agent" | scalebox run --app-file=app.yaml
```

shell 中显式设置的 `CODE_PATH` 优先于 `scalebox.env`，不需要改文件。

## 用 Claude Code

1. 在 `examples/ai-agent-hello/` 目录启动 Claude Code。
2. 首次加载 `.mcp.json` 时，确认信任这个只读 MCP server。
3. 给出类似提示：

```text
先读 AGENTS.md。用 MCP 找到 ai-agent-hello 的 app、module 和 slot 状态，
再告诉我运行和验证这个应用需要哪些命令。不要直接删除 app。
```

Claude Code 会读取 `CLAUDE.md`，其中导入 `AGENTS.md`；MCP 用于查询应用、任务、exec 和 slot 的动态状态。

## 修改与清理

- 只改 `code/run.sh`：重新投递任务即可，`code_path` 是挂载目录。
- 改 `app.yaml` 或 `scalebox.env`（含 `CODE_PATH`）：删除并重建 app，创建时的配置不会热更新。
- 清理：

```sh
scalebox app delete --app-id <app-id>
```

MCP server 只读；创建、投递和删除仍通过 `scalebox` CLI 完成。
