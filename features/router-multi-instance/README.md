# router-multi-instance

主路由按**来源模块**分流到多个实例：同一个 app 里跑多个 main-router slot，每个来源模块
声明的 `router_index` 决定它的回流任务落到哪个实例，从而把重负载来源与轻负载来源隔离开。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `main_router` | `app.parameters` | 模块名 | — | 全 app 只能指定一个主路由模块 |
| `task_dist_mode` | 主路由模块的 `module.parameters` | `SLOT-BOUND` | `DEFAULT` | **必须写在主路由模块上**；只有它为 SLOT-BOUND 时 controld 才会给投往该模块的任务写 `to_slot` |
| `slots` | 主路由模块的 `module.slots` | `<host>:<N>` | — | N 个实例 → `t_slot.seq` = 0..N-1；**host 必须是单节点** |
| `slot_start_batch_size` | 主路由模块的 `module.parameters` | 整数 | 1 | actuator 每 5 秒一轮、按 (模块, host) 分区；多实例时取 N 可一轮起完 |
| `router_index` | **来源模块**的 `module.parameters` | 整数 | 0 | 该模块的回流投到 `t_slot.seq = router_index` 的实例 |

## 行为与不变式

- controld 在创建"投往主路由"的任务时，按来源模块的 `router_index` 查
  `t_slot WHERE mod_id=<router> AND seq=<router_index>`，命中则写 `to_slot`
  （`cmd/controld/task/cache.go` 的 `cacheMainRouterSlotID`，5 分钟缓存）。
- **主路由代码零改动**：分派仍按 `from_module` 走同一份 handler 表，多实例只是部署形态。
- **缺省回落 seq=0**：未声明 `router_index` 的来源（含入口任务、`from_module` 为空）全部进
  第 0 个实例。
- **解析失败不卡死**：查不到目标 slot 时不写 `to_slot`，而 `getSlotBoundRows` 的 SQL 带
  `to_slot IS NULL` 分支，任何实例都能领——任务退化为跨实例负载均衡，不会形成永久 `-1`
  的毒任务。代价是钉住语义丢失，2026-11-03 起 controld 会为每个这类 `from_module` 打一条
  `slot-bound router: no slot resolved ...` 的 WARN（每来源只打一次）。
- 各实例的领任务 SQL 带 `FOR UPDATE SKIP LOCKED`，多实例并发不会重复领同一条任务。
- 计数器 / 信号量在 DB 里，多实例天然共享；同一来源的 handler 之间没有内存状态需要同步。

## 坑

- **`task_dist_mode: SLOT-BOUND` 放错模块**：写在来源模块上无效。它控制的是"投给**主路由**
  的任务要不要写 `to_slot`"，所以必须落在主路由模块的 `parameters` 里。
- **主路由 `slots` 跨多节点**：`cacheMainRouterSlotID` 的查询是
  `SELECT id FROM t_slot WHERE mod_id=? AND seq=?`，多节点时每个节点各有一份 seq 0..N-1，
  查询命中多行、随机取一个——钉住变成掷骰子。多实例必须落在单节点上。
- **`router_index` 超出实例数**：例如 `router_index: 1` 但主路由只有 1 个 slot，查询 0 行 →
  静默降级为负载均衡。配置错误时只看得到性能不达标，看不到报错；靠 WARN 日志发现。
- **`slot_start_batch_size` 默认 1**：actuator 每 5 秒一轮、每轮每 (模块, host) 只起 1 个
  slot，2 个实例要 10 秒才起齐。取实例数可一轮起完。
- **并发度上限没变**：多实例只提高"回流处理"的并发，不改变下游模块的 slot 数。上游产出
  速率超过下游消费能力时，拆多少实例都没用。
- **实例数不是越多越好**：每个实例都要读同一份 workdir / 发同样的 CLI 调用，超过来源模块
  的并发产出能力后只是空转。

## 实证（app-fxcorr，2026-11-03）

2 节点、0.512 s × 120 批的 B 方案。单实例时 router 的处理时间构成：`vtask-head` 204 s（63%）、
`fxcorr-fp` 111 s（34%）、`fxcorr-xp` 5.3 s（1.6%）。`from_xp` 自身只花 11 ms，却要排在
`from_fp`（3840 次、p90 26.67 s）后面平均等 5.67 s。

给 `fxcorr-fp` 声明 `router_index: 1`、主路由开 2 个 slot（`slots: ${IO_NODES}:2`）之后：

| 指标 | 单实例 | 双实例 |
|---|---|---|
| `from_xp` 回流排队 avg / p90 | 5.67 / 9.60 s | **0.84 / 1.73 s** |
| `from_head` 回流排队 avg / p50 | 5.69 / 5.05 s | 2.57 / 0.98 s |
| xp slot 利用率 | 83.4% | **91.2%** |
| 全链 `wall_core` | 811.3 s | **754.9 s（−6.9%）** |

`to_slot` 校验：4563 个主路由任务全部带上 `to_slot`，fp 的 3840 个全部指向 `seq=1`，
controld 无 `no slot resolved` 告警。

## 最小例子

见 `feature.yaml`。关键片段：

```yaml
parameters:
  main_router: router

modules:
  router:
    parameters:
      task_dist_mode: SLOT-BOUND   # 写在主路由模块上
      slot_start_batch_size: 2
    slots:
      - h0:2                       # seq 0 / 1，同一节点

  src-a:
    parameters:
      router_index: 1              # 该模块的回流进 seq=1
```

## 相关

- 任务头 `to_slot` / `to_slot_index` 的语义：`docs/reference/task-headers.md`
- `task_dist_mode` 的三种取值：`docs/reference/app-yaml.md`
- 主路由本身的编程约束（同一 from_module 只能单实例）：`docs/cn/source/design/1_programming_model.md` §1.7.5
