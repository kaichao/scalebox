# recipe — debug-stuck-chain（链卡住的处置流程）

> 用于：app 状态 RUNNING，但任务不再推进。**症状细节全部在 [`../20-diagnosis.md`](../20-diagnosis.md)**，
> 本页只给处置顺序与恢复验证，不重复诊断表。

## 前置

- 知道 `app_id`；能在控制面执行 `scalebox` 与（必要时）`docker exec database psql`。

## 步骤

1. 60 秒四层快照（顺序不要变）：

```bash
scalebox app show --app-id <app_id>              # ① app 状态
scalebox module list --app-id <app_id>           # ② 模块与 module_id
scalebox slot status --module-id <module_id>     # ③ slot 汇总（逐模块）
scalebox task list --app-id <app_id>             # ④ 任务与 status_code
```

2. 分诊（每个分支的判定与修复见 20-diagnosis 同名小节）：

| 快照特征 | 分支 |
|---|---|
| 一个任务都没有 | 入口未投递 / `start_task` 未生效 |
| slot 0 行 | `slots` 匹配问题（INV-15） |
| slot STARTING 不 READY | 镜像 / SSH / 节点可达（INV-16） |
| 任务 `-1` 不动 | slot 未建、归属错、被闸门拦 |
| 任务 `-3` 很久 | 真在跑 vs 进程已死；看超时参数 |
| 有 `>0` 失败码 | 112 / 123 / 124 / 1–10 各自的处置 |

3. 按分支执行修复（常见动作）：

```bash
# 配置类（slots/参数/env）→ 必须重建 app
scalebox app delete --app-id <app_id>
scalebox run --app-file=app.yaml < /dev/null

# 重复投递 112 → 清记录后重投（带 _vtask_id 与 real_path 等业务 header）
# 闸门拦截 → 释放空间，等待自动放行（不要重启 app）

# 单任务查证据
scalebox task show <task_id>                     # stdout/stderr/extras（位置参数）
```

4. 恢复验证：同一批任务继续流动并完成，**且根因已消除**（不是靠重试绕过）。

5. 回流：如果症状与 20-diagnosis 的描述不符，把新症状补进去（这是知识回流的一部分）。

## 判据

- 修复后任务到达终态且 `status_code` 符合预期。
- 清理现场后无残留 app/容器（见 20-diagnosis「清理现场」）。

## 失败下一步

| 现象 | 下一步 |
|---|---|
| 四层都正常但任务不走 | 查 `t_task_exec` 的 headers 与 `t_slot` 明细；可能是调度侧问题 |
| 修复动作无效 | 不要反复重启；收集 `t_task_exec.stderr`/`extras` 与 controld 日志 |
| 同一症状反复出现 | 补一条 invariant 或 feature 坑记录（`_template-invariant.md`） |
