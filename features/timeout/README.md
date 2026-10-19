# timeout

单任务超时：超时后 agent 中断任务进程，任务以超时退出码结束。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `task_max_seconds` | `module.arguments` | 秒 | 0（不限） | 单任务最大运行时间；超时返回 124 |
| `task_global_timeout_scale` | `module.parameters` | 浮点 | 2.0 | 全局超时倍数：外部原因导致任务长期 `-3` 时，按 `task_max_seconds × scale` 恢复为 123 |

## 行为与不变式

- `task_max_seconds > 0` 时，agent 用 `timeout Ns` 包裹 `ACTION_RUN` / `ACTION_SETUP` / `ACTION_TEARDOWN` / `ACTION_CHECK`。
- 单任务超时退出码为 124；全局超时（任务一直处于运行态）返回 123。
- `task_max_seconds = 0` 表示不限制。
- 超时是任务失败（`status_code > 0`），不会回流下游；是否重试由 `retry_rules` 决定。

## 坑

- **`task_min_seconds` 不是超时**：它是 slot 健康阈值（连续多次低于该值时判 GREEDY 并退出），不是"补足任务时长"。
- **超时作用于四类脚本**：`ACTION_CHECK` 也会被 `timeout` 包裹；自定义检查脚本不要依赖无限等待。
- **全局超时不是杀进程**：它把长期 `-3` 的任务状态改为 123，进程可能仍留在节点上，需要人工清理。
- 超时退出码 124 与普通失败（1）在 `retry_rules` 中可分别配置。

## 最小例子

见 `feature.yaml`：

```yaml
modules:
  timeout:
    arguments:
      code_path: ${PWD}/code
      task_max_seconds: ${TIMEOUT_SECONDS}
      slot_options: slot_on_head
    parameters:
      start_task: 0
    environments:
      - SLEEP_SECONDS=${SLEEP_SECONDS}
```

`scalebox.env` 中 `TIMEOUT_SECONDS=10`、`SLEEP_SECONDS=20`，用于制造超时。

## 测试

```sh
./test.sh   # 阶段 3 提供
# 当前手工：
scalebox run --app-file=feature.yaml
# 查看 t_task_exec.status_code / stderr：
#   timeout 超过 10 秒的任务应为 124（或被 kill）
```

预期：任务在约 10 秒后被中断，返回码 124；`t_task_exec` 中可见 timeout 相关 stderr。

## 源码指针

- `go-scalebox/cmd/agent/init.go`（`PLAT_TASK_MAX_SECONDS`）—— 超时读取
- `go-scalebox/cmd/agent/run.go`、`actions.go`、`check.go`—— `timeout Ns` 包裹位置
- `go-scalebox/build/database/initdb.d/62-pgcron_task.sql`—— `task_global_timeout_scale` 与全局超时

## 相关文档

- Sphinx：`docs/cn/source/developer/7_advanced_features.md` §7.4
- 退出码：`docs/cn/source/appendix/6_exit_code_spec.md`（124/123）
- 参数表：`docs/reference/app-yaml.md`（`module.arguments` 的 `task_max_seconds`）
