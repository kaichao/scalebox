# idle-exit

验证 agent 的空闲预算：`max_idle_polls` 是**最后一次非空 task 批次之后**的轮询预算，不是 slot 总寿命；
预算耗尽后 slot 退出并回 READY，`always_running` 豁免，非正值回落缺省值。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `max_idle_polls` | `module.arguments` | 正整数 | 100 | 最后一次非空 task 批次后的空闲轮询上限；`<=0` 或非法值回落 100 |
| `poll_interval_seconds` | `module.arguments` | 正整数 | 6 | 每次轮询间隔；`<=0` 回落 6 |
| `slot_options` | `module.arguments` | 逗号分隔 | 空 | 含 `always_running` 时豁免空闲退出 |
| `ACTION_CHECK` | `module.environments` | 脚本路径 | 空 | 本测试用自定义 check 让 agent 每个 poll 上报 `idle_polls_remaining` |

空闲预算时长 = `max_idle_polls × poll_interval_seconds`。例如 `20 × 1s = 20s`。

## 行为与不变式

- agent 领到并处理 task 后，`numIdlePolls` 重置为 `max_idle_polls`；空轮询、`GetTaskList` 失败、
  `check` 失败各递减一次。
- 预算耗尽后，非 `always_running` 的 slot 优雅退出：`SLOT_COMPLETE` → `SetSlotTerminated`，
  slot 回 READY（或模块 `initial_slot_status`）。
- 退出事件写入：

  ```text
  slot exit=idle timeout
  slot exit:num_idle_polls=<remaining>
  slot exit:max_idle_polls=<max>
  slot exit:poll_interval_seconds=<interval>
  ```

- `max_idle_polls <= 0` 或非法值回落 100；`poll_interval_seconds <= 0` 回落 6。
  需要 one-shot 时用 `slot_max_tasks: 1` 或 `task header slot_last_task=yes`，不要依赖 `max_idle_polls=0`。
- host OFF 不走空闲预算：controld 下发 `slot_action=stop` 并级联复位活动 slot；只有 agent 不可达时
  才轮到空闲退出兜底。

## 坑

- **把 `max_idle_polls` 当 slot 总寿命**：有 task 持续到达时不会退出；它是“闲下来之后还能等多久”。
- **用 0 表示立即退出**：当前 agent 把 0 当非法值回落 100，约 10 分钟（按 poll=6s）。
- **没有自定义 check 时，运行态/剩余预算更新慢**：agent 只在心跳/任务事件上报；本测试用
  `ACTION_CHECK` 把上报频率提到每个 poll。
- **`always_running` 只豁免空闲退出**：host OFF 的 `slot_action=stop` 仍然生效。

## 最小例子

见 `feature.yaml`。关键片段：

```yaml
modules:
  idle-reset:
    arguments:
      code_path: ${PWD}/code
      max_idle_polls: 20
      poll_interval_seconds: 1
      slot_options: slot_on_head
    parameters:
      start_task: reset-probe
    environments:
      - ACTION_CHECK=/app/bin/check.sh
```

## 测试

```sh
./test.sh
```

断言点：

1. `idle-reset` 的 start_task 完成后，`t_slot_runtime.idle_polls_remaining` 回到 20 附近（>=15）；
2. `idle-exit` 的 slot 在预算耗尽后变为 READY，`SLOT_COMPLETE` 的 `slot exit=idle timeout`；
3. `idle-exit` 的退出明细包含 `max_idle_polls=5`、`poll_interval_seconds=1`；
4. `idle-always`（`always_running`）在预算窗口后仍为 ON；
5. `idle-clamp`（`max_idle_polls=0`）不立即退出，仍为 ON。

## 源码指针

- `go-scalebox/cmd/agent/run.go`：`idlePollBudget` / `idleExitDetails` / 空闲退出判断
- `go-scalebox/cmd/agent/module_args.go`：`normalizeIdlePolls` / `normalizePollInterval`
- `go-scalebox/cmd/agent/check.go`：`num_idle_polls` 上报
- `go-scalebox/cmd/controld/grpc/slot.go`：`num_idle_polls` 落表
- `go-scalebox/internal/slot/store_detail.go`：`idle_polls_remaining` 查询

## 相关文档

- [`docs/reference/slot-runtime.md`](../../docs/reference/slot-runtime.md)：运行态、心跳、空闲预算
- [`docs/reference/module.md`](../../docs/reference/module.md)：`max_idle_polls` / `poll_interval_seconds`
- [`docs/agents/10-invariants.md`](../../docs/agents/10-invariants.md) INV-25
