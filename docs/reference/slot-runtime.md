# Slot 运行态参考

> 唯一信源：`t_slot.status` 的生命周期状态与 `t_slot_runtime` 的运行态字段。
> 语义、设计取舍见 go-scalebox `docs/slot-runtime-design.md`；排障步骤见
> [`docs/agents/20-diagnosis.md`](../agents/20-diagnosis.md)。
> 最后校验：2026-11-22，对照 go-scalebox@762bb55 + 空闲预算工作区改动。

## 1. 两层状态

| 层 | 数据源 | 取值 | 用途 |
|---|---|---|---|
| 生命周期 | `t_slot.status` | `READY` / `STARTING` / `ON` / `OFF` / `ERROR` / `SSHFAILED` / `TIMEOUT` / `GREEDY` | slot 是否可被启动、是否应派发任务 |
| 运行态 | `t_slot_runtime` 派生 | `RUNNING` / `ASSIGNED` / `WAITING` / `BLOCKED` / `STALE` / `UNKNOWN` | `ON` 期间具体在做什么、为什么没领到任务 |

生命周期不是 `ON` 时，运行态直接显示生命周期值（如 `READY`）。“ON 但没有任务”不等于故障；
先看运行态和等待原因。

## 2. 运行态

| 状态 | 判定 | 含义 |
|---|---|---|
| `RUNNING` | 存在未完成的 `t_task_exec.ts_started` | 正在执行 task；页面显示已运行秒数 |
| `ASSIGNED` | 无运行中 task，但存在待启动 exec（`ts_started IS NULL`） | task 已派发到 slot，agent 尚未开始执行 |
| `WAITING` | 无 task，且等待原因是 `NO_ELIGIBLE_TASK` / `FLOW_NODE_PROGRESS` / `FLOW_VTASK` / `DEDUP` | 正常等待：没任务、流控、节点进度或去重 |
| `BLOCKED` | 无 task，且等待原因是 `NOT_DISPATCHABLE` / `DISPATCH_ERROR` | 调度前置不满足（slot/module/host 状态、主机匹配、派发错误） |
| `STALE` | 无 task/exec，且心跳年龄超过阈值 | 心跳过期；可能 agent 已死或长任务期间心跳缺失 |
| `UNKNOWN` | 无运行/待启动 task，且没有有效等待原因 | 信息不足；等待原因超过 TTL 后也会回到这里 |

判定优先级：`RUNNING` > `ASSIGNED` > `STALE` > 等待原因。运行中或已有待启动 exec 时，
不用 `STALE` 覆盖运行态；心跳单独标红。

## 3. 等待原因

| `wait_reason` | 分类 | 典型含义 |
|---|---|---|
| `NO_ELIGIBLE_TASK` | WAITING | `GetTaskList` 返回空批次：当前没有符合条件的 task |
| `FLOW_VTASK` | WAITING | vtask 额度信号量不足，`GetTaskList` 主动停领 |
| `FLOW_NODE_PROGRESS` | WAITING | `node_progress` 节点进度差限制 |
| `DEDUP` | WAITING | 去重缓存导致本批任务被过滤 |
| `NOT_DISPATCHABLE` | BLOCKED | slot 非 ON、module 非 RUNNING、host 不可用或模块/host 不匹配 |
| `DISPATCH_ERROR` | BLOCKED | `GetTaskList` 查询/派发错误 |
| `UNKNOWN` | — | 未采集到原因 |

`wait_detail_json` 是原因相关明细（如信号量名、当前值、批次大小）；等待状态超过
120 秒未更新时，运行态回到 `UNKNOWN`，避免显示过期原因。

## 4. `t_slot_runtime`

`t_slot_runtime` 与 `t_slot` 1:1，随 slot 删除级联删除。过渡期 `t_slot.last_active` /
`t_slot.create_code` 仍双写，不删除。

| 字段 | 更新点 | 说明 |
|---|---|---|
| `run_started_at` | `SLOT_START` / `SetSlotInitialized` / `RegisterSlot` | 本次 agent run 起点；用于过滤上一轮 run 的孤儿 exec |
| `last_seen_at` | 每次 `SubmitSlotEvent` | 最近一次看到 slot 活动；旧列 `last_active` 双写 |
| `last_poll_at` | `GetTaskList` | 最近一次轮询；内存缓存批量落表 |
| `wait_reason` / `wait_detail` / `wait_since` | `GetTaskList` | 当前等待原因、明细、开始时间 |
| `start_error_code` / `start_error_at` | actuator / `SLOT_CREATE` | 最近一次启动命令非零退出；成功启动或新 `STARTING` 清除；旧列 `create_code` 双写 |
| `idle_polls_remaining` | agent 上报 `extras["num_idle_polls"]` | 最后一次上报的空闲轮询剩余次数；`NULL` 表示未知 |
| `updated_at` | 每次行更新 | 最近写入时间 |

时间字段都是 `TIMESTAMP WITHOUT TIME ZONE`；controld 用
`internal/postgres.ReinterpretTimestamp` 按进程时区解释后返回，避免 API 时间偏移。

## 5. 心跳与 STALE

- 心跳间隔 `heartbeat_seconds`（模块 `arguments`，缺省 60 秒）。
- 心跳年龄 = `now - last_seen_at`；超过 `3 × heartbeat_seconds` 标为 stale。
- 长任务期间 agent 主循环不轮询，但心跳 goroutine 仍应上报；若 `last_seen_at` 长时间不前进，
  先检查 agent 容器/网络，而不是只看 `t_slot.status=ON`。
- `STALE` 是“心跳过期”，不是生命周期状态；slot 仍可能是 `ON`，需要人工或
  `host OFF` 级联复位。

## 6. 空闲预算

agent 用两个参数控制空闲退出：

| 参数 | 位置 | 语义 |
|---|---|---|
| `max_idle_polls` | `module.arguments` | 最后一次非空 task 批次后的空闲轮询上限；缺省 100；非正整数/非法值回落 100 |
| `poll_interval_seconds` | `module.arguments` | 每次轮询间隔；缺省 6 秒；非正/非法值回落 6 |

- **非空 task 批次后重置**：只要 agent 领到并处理了 task，`numIdlePolls` 重置为
  `max_idle_polls`。`max_idle_polls` 不是 slot 总寿命，而是“最后一批 task 之后还能空转多久”。
- 空轮询、`GetTaskList` 失败、`check` 失败各递减一次。
- `slot_options` 含 `always_running` 时豁免空闲退出；此时 `idle_polls_remaining`
  可能为 0 或负数，页面显示 `Idle exit: disabled (always_running)`。
- 空闲预算 = `max_idle_polls × poll_interval_seconds`。例如
  `3600 × 1s ≈ 60 min`。
- 空闲退出时 `SLOT_COMPLETE` 的 `details` 写入：

```text
slot exit=idle timeout
slot exit:num_idle_polls=<remaining>
slot exit:max_idle_polls=<max>
slot exit:poll_interval_seconds=<interval>
```

## 7. 观测入口

| 入口 | 能看到的运行态 |
|---|---|
| WebUI `App → Slots` | 4 项单元格：状态 + HB、Details、TaskExec、SlotLog；状态颜色；运行时详情、空闲预算、等待原因、心跳年龄；手动 Refresh 快照 |
| `GET /api/v1/slots?include_runtime=true` | `ListSlots` 返回 `runtime`（状态、等待原因、心跳、当前 task、idle 字段） |
| `GET /api/v1/slots/{id}` | 单个 slot 的 runtime |
| `GET /api/v1/slots/{id}/logs?page=1&page_size=100` | SlotLog：生命周期/心跳事件，按时间倒序 |
| CLI `scalebox slot list` / `status` | 仅生命周期状态；`slot show` 当前也不打印 runtime 字段 |
| MCP `list_slots` | 生命周期、心跳时间、主机信息；不返回 `runtime_state` / `wait_reason` |

需要等待原因、空闲预算或 SlotLog 时，用 WebUI 或 REST；CLI/MCP 只覆盖生命周期层。

## 8. 源码指针

- `build/database/initdb.d/88-slot-runtime.sql`、`89-slot-log-retention.sql`、
  `90-slot-runtime-idle.sql`：表、索引、保留期、空闲预算列。
- `cmd/controld/slot/runtime/`：派生、内存等待缓存、SQLStore。
- `cmd/controld/task/get.go`：等待原因采集。
- `cmd/controld/grpc/slot_runtime.go`、`slot_log.go`：runtime 映射与 SlotLog RPC。
- `internal/slot/store_detail.go`：runtime 查询。
- `cmd/agent/run.go`、`module_args.go`：空闲预算、重置与退出原因。
- `cmd/webui/frontend/src/pages/apps/AppSlots.tsx`：App Slots UI。

## 9. 相关文档

- go-scalebox `docs/slot-runtime-design.md`：完整设计、字段语义、迁移与风险。
- [`docs/agents/20-diagnosis.md`](../agents/20-diagnosis.md)：ON 但任务不动的排查顺序。
- [`docs/reference/module.md`](./module.md)：`max_idle_polls` / `poll_interval_seconds` 的模块参数上下文。
