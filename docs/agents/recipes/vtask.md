# recipe — vtask（建任务组管道）

> 用于：一批任务需要整体流控、绑定到单节点/节点组、或需要 `vtask get` 看整体进度。
> 原理见 Sphinx `design/3_vtask`；信号量与命令事实见 [`reference/semaphore-variable.md`](../../reference/semaphore-variable.md)。
> 可直接复制的完整示例：`examples/vtask/`。

## 前置

- app 已有入口模块，了解三种模式：`DEFAULT`（无绑定）/ `HOST-BOUND`（单节点）/ `GROUP-BOUND`（节点组）。
- 节点池可用（HOST-BOUND/GROUP-BOUND 需要 wait-queue 的 `slot_on_head` 与 head 的 `slots`）。

## 步骤

1. 选模式并准备 env 文件（示例已备好三份）：

```sh
# DEFAULT
VTASK_MODE=DEFAULT   TASK_DIST_MODE=      HEAD_SLOTS=h0    CORE_MODE=      CORE_SLOTS=h0
# HOST-BOUND
VTASK_MODE=HOST-BOUND TASK_DIST_MODE=HOST-BOUND HEAD_SLOTS=n0-[01] CORE_MODE=HOST-BOUND CORE_SLOTS=n0-[01]
# GROUP-BOUND
VTASK_MODE=GROUP-BOUND TASK_DIST_MODE=SLOT-BOUND HEAD_SLOTS=<head 节点> CORE_MODE=HOST-BOUND CORE_SLOTS=<组内节点>
```

`VTASK_MODE` 与 `TASK_DIST_MODE` 成对出现，收敛在 env 文件一处；应用脚本只读 `VTASK_MODE`。

2. 核对 `app.yaml` 的模块骨架：

| 模块 | 关键参数 |
|---|---|
| `main-router` | `parameters.main_router: main-router` |
| `wait-queue` | 标准镜像；`parameters.vtask_size: 1`（串行 bind） |
| `vtask-head` | `parameters.vtask_role: head`、`vtask_size: <在制上限>` |
| `vtask-core` | `parameters.vtask_role: core` |
| `vtask-tail` | `parameters.vtask_role: tail`；完成时平台自动归还额度 |

3. 建 app 并投递一批：

```bash
scalebox run --env-file=<mode>.env --app-file=app.yaml < /dev/null
cat <mode>-tasks.txt | scalebox task add --app-id <app_id> --sink-module main-router
```

4. 观察：

```bash
scalebox vtask list --app-id <app_id>
scalebox vtask get --vtask-id <vtask_id>          # 状态、进度、绑定资源
scalebox semaphore ls --vtask-size-            # 或 semaphore get vtask_size:vtask-head
```

## 判据

- 每个 body 对应一个 vtask；`vtask get` 显示 `RUNNING` 后转为完成，子任务数符合预期。
- `vtask_size:*` 信号量在每批完成后回到初值（无泄漏）。
- 在制 vtask 数 ≤ `vtask_size`；`wait-queue` 不长时间停领。

## 失败下一步

| 现象 | 下一步 |
|---|---|
| 第一个 vtask 过后 wait-queue 停领 | head 未对 `vtask_size:wait-queue` 加 1（INV-10） |
| 在制数超过 `vtask_size` | 脚本重复 unBind 或重复 bind（INV-11）；查 `semaphore get` 的值 |
| `vtask bind` 报找不到信号量 | `host_vtask_size`/`group_vtask_size` 名不匹配；两套 API 混用（INV-12） |
| `_vtask_attr_*` 在 head 丢失 | wait-queue 一跳需显式转发（INV-13） |
| 链卡在 wait-queue | 先看 `check.sh` 读的 semagroup max 是否 ≤ 0 → [`debug-stuck-chain`](./debug-stuck-chain.md) |
