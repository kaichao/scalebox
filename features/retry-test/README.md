# retry-test

按退出码自动重试：任务失败（`status_code > 0`）后，由平台定时扫描并按规则重置为待执行状态。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `retry_rules` | `module.parameters` | 字符串数组，如 `"['1:3','2:2','*']"` | 空 | `退出码:重试次数`；`*` 为兜底；次数缺省 1 |

YAML 中写成字符串（单引号数组）；创建 app 时，controld 会把单引号转成双引号并保存为 JSON 数组。

## 行为与不变式

- 触发：pgcron 每 5 分钟扫描 `status_code > 0` 的任务。
- 条件：该 task 的 `t_task_exec` 记录数 ≤ 配置的重试次数（`num`），且退出码匹配规则或命中 `*`。
- 动作：`status_code` 重置为 `-1`，并写入 `headers.repeatable = "yes"`（允许重复分发）。
- 次数语义：`num` 是**重试次数**；`num=1` 表示首次失败后再执行一次。
- 没有配置 `retry_rules` 的模块不会重试。

## 坑

- **不是即时重试**：由 5 分钟一次的 pgcron 触发，最短等待接近 0、最长约 5 分钟。
- **重试任务带 `repeatable=yes`**：这会影响任务去重；排查"任务为什么又跑了一次"时先看该 header。
- **与人工重投的区别**：自动重试是平台行为；人工重投需要自己清记录/设 `repeatable`（参见 feature `task-identity`，待建）。`CONFLICT_ACTION=OVERWRITE` 的适用条件在应用间存在差异（见 `docs/feature-inventory.md` 待确认项），不要与自动重试混用。
- 退出码为 0 的任务不会重试；超时（124）等退出码同样按规则匹配。

## 最小例子

见 `feature.yaml`：

```yaml
modules:
  retry_test:
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
    parameters:
      start_task: 0
      retry_rules: "['1:3','2:2','*']"
    sink_modules:
      - retry_test
```

`code/run.sh` 按预设返回码退出，用于观察重试次数与 `t_task_exec` 记录数。

## 测试

```sh
./test.sh   # 慢测：等 pgcron（每 5 分钟）重置任务，通常 2–5 分钟
# 当前手工：
scalebox run --app-file=feature.yaml
# 等待 pgcron 周期后查看：
#   t_task.status_code 是否回到 -1
#   headers.repeatable 是否为 yes
#   t_task_exec 记录数是否达到 num+1
```

## 源码指针

- `go-scalebox/cmd/controld/grpc/app_create.go`（`saveModule`）—— 把 YAML 字符串转为 JSON 数组
- `go-scalebox/build/database/initdb.d/62-pgcron_task.sql`（`get_retryable_tasks` / `recover_repeatable_tasks`）—— 扫描与重置
- `go-scalebox/cmd/agent/`—— `repeatable` header 的去重语义

## 相关文档

- Sphinx：`docs/cn/source/design/7_advanced_features.md` §7.1（摘要）
- 参数表：`docs/reference/app-yaml.md`（`module.parameters` 的 `retry_rules`）
- 退出码：`docs/cn/source/reference/4_exit_code_spec.md`
