# task-identity

任务身份与幂等：平台用唯一约束定义"同一个任务"，重复投递会被拒绝；重跑需要显式处理。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `conflict_action` | `task add --conflict-action` / task header / `sink-tasks.txt` 行级 headers / `CONFLICT_ACTION` 环境变量 | `OVERWRITE`、`IGNORE`、空 | 空 | 冲突时的处理策略；优先级：sink 行级 > task header > 环境变量 |
| `repeatable` | task header | `yes` | 不设 | 允许同一任务重复分发（重试/重跑） |
| `task_dedup_cache_ttl_minutes` | `module.parameters` | 分钟 | 30 | 去重缓存窗口；应大于任务最长运行时间 |
| `max_sink_tasks` | `module.arguments` → `MAX_SINK_TASKS` | 非负整数 | router 0（不分批），其他 100 | 单次 `AddTaskList` 最大行数；0 = 整份 `sink-tasks.txt` 一次投递 |

任务身份（`uniq_task` 唯一约束）= `(mod_id, from_mod, to_host, vtask, body)`。

## 行为与不变式

- 相同身份的任务重复插入 → duplicate key，任务退出码 112。
- 一次 `AddTaskList` 的所有行（含不同 `conflict_action` 策略分组）在**同一数据库事务**内落库：任意一行冲突或失败，整批回滚，不再出现部分成功。
- `sink-tasks.txt` 行级 headers 中的 `conflict_action` 由 controld 逐行解析，解析后从该行 headers 中移除，不写入 `t_task.headers`。
- `repeatable=yes` 时允许重复分发；自动重试（`retry_rules`）会由平台设置该 header。
- 去重窗口由 `task_dedup_cache_ttl_minutes` 控制；窗口内相同 body 不会被再次分发。
- 人工重跑应清理该任务的 `t_task_exec` / `t_task` 记录（带 app 过滤），再重新投递。
- 人工重投同一任务可用 `OVERWRITE` + `repeatable=yes`。
- 平台不对重复投递做隐式忽略：`IGNORE` 只是显式的行级策略，不是 vtask-tail 等收尾场景的默认行为。应用必须保证同一身份只投一次（例如全部 core 完成后只投一次 tail）；重复身份会让整批投递回滚并返回 112。

## 坑

- **112 不是业务错误**：是幂等约束；排查时先看是否重复投递。
- **不要用 `OVERWRITE` 绕过 router 冲突**：`router` 重建下游任务会撞 `uniq_task`；一次投递内任意一行冲突都会整批回滚，不能靠部分成功继续。
- **显式分批会牺牲原子性**：`max_sink_tasks > 0` 时按批调用、每批各自事务，跨批不保证"全部或不做"；router 缺省 0，整份文件一次投递。
- **重跑清记录必须带 app 过滤**：否则会误删其他应用的任务。
- **去重窗口要覆盖任务时长**：窗口小于任务时长会导致"任务还没跑完就被允许重投"。

## 最小例子

`feature.yaml`：单模块 + sink 自环，便于重复投递观察 112。

```yaml
modules:
  identity:
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
    parameters:
      start_task: 0
      sink_modules:
        - identity
```

## 测试

```sh
./test.sh
# 当前手工：
APP_ID=$(scalebox run --app-file=feature.yaml | ...)
# 第一次：成功
echo "body-1" | scalebox task add --app-id=$APP_ID --sink-module=identity
# 第二次：相同身份，预期退出码 112
echo "body-1" | scalebox task add --app-id=$APP_ID --sink-module=identity
```

预期：第二次投递报 duplicate key（112）；加 `--header repeatable=yes` 后可重复分发。

## 源码指针

- `go-scalebox/build/database/initdb.d/25-updated-schema.sql`（`uniq_task`）—— 身份定义
- `go-scalebox/cmd/controld/task/add_local.go`（`conflict_action` / 单事务投递）—— 冲突处理
- `go-scalebox/cmd/agent/sink_tasks.go`—— 行级策略随行传递；router 默认不分批
- `go-scalebox/cmd/controld/grpc/app_env_vol.go`—— `max_sink_tasks` → `MAX_SINK_TASKS` 映射
- `go-scalebox/cmd/cli/task/task.go`—— `--conflict-action` / `--header`
- `go-scalebox/cmd/controld/task/cache_attrset.go`—— `task_dedup_cache_ttl_minutes`

## 相关文档

- 参数表：`docs/reference/app-yaml.md`（`module.parameters` 的 `task_dedup_cache_ttl_minutes`）；`docs/cn/source/reference/2_parameter_reference.md` §2.1.4（task headers）
- Sphinx：`docs/cn/source/design/7_advanced_features.md` §7.1（重试）
- 实战：`docs/agents/10-invariants.md` INV-06 / INV-07
