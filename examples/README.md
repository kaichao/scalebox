# examples

## 1. hello-scalebox
scalebox的入门应用，展示基于消息触发的基本应用。具体操作步骤如下：

- 本目录没有独立镜像：模块脚本 `code/run.sh` 由缺省的
  `scalebox.net/platform/agent:latest` 镜像执行（需要 singularity 镜像时见目录内
  `singularity.md`）。
- 创建并运行应用（本目录没有 `app.yaml`，`scalebox run` 按命令行参数生成应用模板；
  任务 body 经 stdin 传给初始模块）
```sh
cd examples/hello-scalebox
echo "hello" | scalebox run
```
- 查看应用状态：
```bash
scalebox app list
```


## 2. app-primes
计算区间内质数数量的分布式应用。首先将区间分解为若干个子区间，再分别计算各个子区间的质数数量，最后加和得到区间内结果。

操作步骤如下：
```bash
cd examples/app-primes
make build
```
- 创建并运行应用（目录内有 `app.yaml`，`scalebox run` 自动加载；输入为区间上限）
```sh
echo 1000 | scalebox run
```
- 查看应用状态：
```bash
scalebox app list
```


## 3. app-copy
文件拷贝是scalebox的基本功能，实现集群内、跨集群的数据传输。

scalebox的基本模块支持基于rsync、rsync-over-ssh、ftp等的文件拷贝。

可运行示例见 `dockerfiles/files/app-dir-copy/`（`app.yaml` + `scalebox run`），
各传输模块的用法见 `dockerfiles/files/*/README.md`。


## 4. ai-agent-hello
给 AI agent 使用的最小应用示例：一个模块、一个任务，附带 `AGENTS.md`、
`CLAUDE.md` 和只读 `.mcp.json`。适合验证 AI 读取应用规则、用 MCP 查询
应用/任务/slot 状态，再通过 CLI 运行应用的完整闭环。

操作步骤如下：
```bash
cd examples/ai-agent-hello
scalebox run --app-file=app.yaml
echo "hello agent" | scalebox task add --app-id <id> --sink-module hello
scalebox task list --app-id <id>
scalebox task show <task-id>
```
Claude Code 的接入步骤和示例提示词见目录内 `README.md`。连接远端控制面时，
`CODE_PATH` 必须指向远端计算节点可见的路径，见目录内 `README.md`。


