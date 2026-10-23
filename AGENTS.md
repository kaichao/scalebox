# AGENTS.md — Scalebox 平台入口（面向 AI agent）

本页是 agent 的**第一入口**：先读这里，再按导航表深入。平台深度规则在
[`docs/agents/10-invariants.md`](./docs/agents/10-invariants.md)，可执行样例在
[`features/`](./features/)，生成的结构参考在 [`docs/reference/`](./docs/reference/)。

## 心智模型

- **App**：一次应用实例；`app.yaml` 定义模块与参数。创建时快照环境变量、镜像与参数，改配置要重建 app。
- **Module**：容器化算法单元；模块之间用 `sink_modules` / `sink-tasks.txt` 串联。
- **Task**：模块的一次执行；`body` 是输入，headers 携带上下文，`status_code` 是结果（0 成功）。
- **Slot**：节点上的执行资源；模块按 `slots` 正则建 slot，agent 在 slot 里领任务。
- **VTask**：跨模块的任务集合；用于批/组级流控（wait-queue → vtask-head → vtask-tail）。
- **控制面/数据面**：controld/actuator/database 管调度；数据走 `/cluster_data_root` 等挂载。

## 导航表

| 要做的事 | 先读 |
|---|---|
| **先建立全局图** | [`docs/agents/00-model.md`](./docs/agents/00-model.md)（心智模型）、[`glossary.md`](./docs/agents/glossary.md)（术语） |
| **按步骤做一件事** | [`docs/agents/recipes/`](./docs/agents/recipes/README.md)：建 app / 加模块 / 流控 / vtask / 多节点 / 排障 / 验收 |
| 从零写一个 app.yaml | [`docs/reference/app-yaml.md`](./docs/reference/app-yaml.md)、[`examples/hello-scalebox`](./examples/hello-scalebox/) |
| 查 CLI 命令与 flags | [`docs/reference/cli.md`](./docs/reference/cli.md) |
| 模块怎么跑起来（挂载/注入/路由） | [`docs/reference/module.md`](./docs/reference/module.md) |
| task header / 信号量 / 变量 / 环境变量 / RPC | [`docs/reference/`](./docs/reference/)（`task-headers.md`、`semaphore-variable.md`、`env-vars.md`、`rpc.md`） |
| 写模块脚本、接口文件（sink/task-exec） | [`features/agent-interface-files/README.md`](./features/agent-interface-files/README.md) |
| 任务回流 / `_` header 传播 | [`features/headers-propagation/README.md`](./features/headers-propagation/README.md) |
| 任务身份 / 重复投递 112 | [`features/task-identity/README.md`](./features/task-identity/README.md) |
| 超时 / 退出码 | [`features/timeout/README.md`](./features/timeout/README.md)、[`docs/reference/exit-codes.md`](./docs/reference/exit-codes.md) |
| 空间/准入流控 | [`features/check-test/README.md`](./features/check-test/README.md) |
| slot 落点 / HOST-BOUND / pod | [`features/pod/README.md`](./features/pod/README.md) |
| vtask 流控 | [`docs/agents/10-invariants.md`](./docs/agents/10-invariants.md) INV-08…INV-14 |
| **链卡住/任务不动** | [`docs/agents/20-diagnosis.md`](./docs/agents/20-diagnosis.md) |
| 跑/加 feature 集成测试 | [`features/README.md`](./features/README.md)（`./run-tests.sh`） |
| 生成物来源与再生成 | [`docs/reference/README.md`](./docs/reference/README.md) |

## 铁律摘要（详见 invariants）

1. **app 创建时快照**环境变量、镜像、参数；改完要重建 app（INV-01）。
2. 未定义变量**保留字面量** `${VAR}`，不会自动为空（INV-02）。
3. 非 tty stdin 的任务列表**只读一次**；不要指望循环读 stdin（INV-03）。
4. 非路由模块写 `${WORK_DIR}/sink-tasks.txt` 回流；`_` 前缀 header 才自动继承（INV-04/05）。
5. **`slots` 用短名/正则**（`h0`、`n0-[01]`），不要写完整 FQDN（`h0.local`），不要逗号串（INV-15）。
6. 相同 `(mod, from, host, vtask, body)` 重复投递触发 **uniq_task / 112**；重跑先清记录或显式处理冲突（INV-06/07）。
7. **任务级空间闸门会死锁**；`free_space_gb` 挂批入口（vtask-head）而不是每个任务（INV-12）。
8. vtask-head 负责释放 gate；vtask-tail 不要二次归还（INV-13/14）。
9. **两套信号量 API 不要混用**（INV-11）。
10. 模块容器优先用 `code_path` 挂到 `/app/bin`；容器路径有白名单（INV-17）。

## 常用命令

```sh
scalebox run --app-file=app.yaml            # 建 app，输出 {"app_id":N}
echo "body" | scalebox task add --app-id N --sink-module producer -H _demo=hello
scalebox app list
scalebox module list --app-id N              # 取 module-id
scalebox slot list --module-id M            # 查 slot
scalebox task list --app-id N               # 查任务
scalebox task show T                        # 查单任务（位置参数，含 stdout/stderr）
scalebox app delete --app-id N              # 删 app（不自动停本地容器）
```

测试部署通常没有默认基础镜像；feature 统一用
`base_image: scalebox.net/platform/agent:latest`。`features/` 的 `test.sh` 封装了
建 app → 等任务 → 断言 → 清理的完整生命周期，可当作最小范例读。

## MCP server

`scalebox mcp serve --read-only` 向 MCP 客户端提供只读平台查询：
`list_apps`、`get_app`、`list_tasks`、`get_task`、`get_task_exec`、`list_slots`。
它只读、不执行 shell；`AGENTS.md` 仍是静态规则的入口。客户端配置见
`docs/cn/source/application/9_ai_assisted_development.md`。

## 生成物

`docs/reference/` 分两类：

- **生成物**（`cli.md`、`app-yaml.md`、`exit-codes.md`、`rpc.md`、`env-vars.md`）：由 go-scalebox 的
  `scripts/gen-reference/` 生成，**不要手改**；源码变化后重新生成并提交 diff。
- **手写参考**（`task-headers.md`、`module.md`、`semaphore-variable.md`）：结构以源码为准，修改时对照源码。

语义说明写在本页、`docs/agents/` 或 `features/`，不要写进生成器。
