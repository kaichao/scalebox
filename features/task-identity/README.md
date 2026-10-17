# task-identity

任务身份与幂等：平台用唯一约束定义"同一个任务"，重复投递会被拒绝；重跑需要显式处理。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `conflict_action` | `task add --conflict-action` / task header | `OVERWRITE` 等 | 空 | 冲突时的处理策略 |
| `repeatable` | task header | `yes` | 不设 | 允许同一任务重复分发（重试/重跑） |
| `task_dedup_cache_ttl_minutes` | `module.parameters` | 分钟 | 30 | 去重缓存窗口；应大于任务最长运行时间 |

任务身份（`uniq_task` 唯一约束）= `(mod_id, from_mod, to_host, vtask, body)`。

## 行为与不变式

- 相同身份的任务重复插入 → duplicate key，任务退出码 112。
- `repeatable=yes` 时允许重复分发；自动重试（`retry_rules`）会由平台设置该 header。
- 去重窗口由 `task_dedup_cache_ttl_minutes` 控制；窗口内相同 body 不会被再次分发。
- 人工重跑应清理该任务的 `t_task_exec` / `t_task` 记录（带 app 过滤），再重新投递。
- `OVERWRITE` 的适用条件存在场景差异：人工重投可用；router 重建下游任务时禁用。边界待平台确认。

## 坑

- **112 不是业务错误**：是幂等约束；排查时先看是否重复投递。
- **不要用 `OVERWRITE` 绕过 router 冲突**：`router` 重建下游任务会撞 `uniq_task`。
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
./test.sh   # 阶段 3 提供
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
- `go-scalebox/cmd/controld/task/add_local.go`（`conflict_action`）—— 冲突处理
- `go-scalebox/cmd/cli/task/task.go`—— `--conflict-action` / `--header`
- `go-scalebox/cmd/controld/task/cache_attrset.go`—— `task_dedup_cache_ttl_minutes`

## 相关文档

- 参数表：`docs/cn/source/appendix/2_parameter_reference.md` §2.1.3–2.1.4
- Sphinx：`docs/cn/source/developer/7_advanced_features.md` §7.1（重试）
- 实战：`docs/agents/10-invariants.md` INV-06 / INV-07
