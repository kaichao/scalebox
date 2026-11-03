# vtask 规模化设计：L1 兼容改进与 L2 自扣模式

> 状态：设计（2026-11-02），未实施。
> 相关：`docs/agents/recipes/vtask.md`（操作）、`docs/cn/source/design/3_vtask.md`（机制）、
> `examples/vtask/`（可跑示例）、`docs/agents/10-invariants.md`（红线）。

## 0. 摘要

现有 vtask 链路的**分配吞吐与 host 数量无关**：`wait-queue` 只有一个 slot，agent
串行执行任务，全平台约 15–20 个 vtask/s。500 台 host、每台容量 2、vtask 寿命 5 s
时需要 200 个/s，差一个数量级；寿命 1 s 时差两个数量级。

本设计给出三层，按需取用：

| 层 | 模型 | 吞吐上限 | 延迟 | 适用 |
|---|---|---|---|---|
| **L1** 兼容改进 | wait-queue 不变 | 400–800/s | 2 × poll | ≤ 500 host、≥ 5 s 任务 |
| **L2** 自扣模式 | 去掉 wait-queue | ~500/s | 1 × poll | 低延迟、无中心单点 |
| **L3** 分布式 | 多 controld / 分片 wait-queue | 线性 | — | > 1000 host、< 1 s 任务 |

**推荐路径：L1 先行**（4 处改动、模型与用户心智不变），L2 按需引入，L3 等实际压力出现。

---

## 1. 现状机制

### 1.1 链路

```
main-router ──► wait-queue（1 slot，head 节点）
                   │  scalebox vtask bind            → 扣额度、拿 to_host
                   │  scalebox vtask add-subtask --direct
                   ▼
              vtask-head（每 host 1 slot）
                   │  业务逻辑（展开/转发）
                   ▼
              vtask-core（HOST-BOUND，pod_id 推导 to_host）
                   ▼
              vtask-tail ──► controld doVTaskFinished → 额度 +1
```

### 1.2 两把锁，不要混为一谈

| 信号量 | 谁扣 | 谁还 | 语义 |
|---|---|---|---|
| `host_vtask_size:<head>:<host>` | `scalebox vtask bind`（预扣型）或 `updateVTask`（自扣型） | controld `doVTaskFinished`（tail 完成时） | **资源额度**：在制 vtask 数 |
| `vtask_size:wait-queue` | `updateVTask`（wait-queue 领取任务时） | `from-vtask-head.sh` 的 `semaphore increment` | **串行门**：同一时刻只允许一个 vtask 进入分配 |

额度归还**不在** wait-queue、**不在** head，而在 controld
（`cmd/controld/task/vtask.go:115 doVTaskFinished`，读 tail task 的 `_vtask_size_sema`）。
`from-vtask-head.sh` 里的 `+1` 还的是 wait-queue 自己的串行门，不是额度。

### 1.3 额度归属

`vtask bind` → `BindVtaskResource` → `semagroup.Decrement`
（`internal/semagroup/decrement.go:24-45`）：

```sql
SELECT name, value FROM t_semaphore
WHERE name ~ $1 AND app = $2
ORDER BY value DESC, name
LIMIT 1
```

选**当前值最大**的那台 host 扣一。这是 O(H) 扫描 + 排序，且事务内无行锁——
并发 bind 会读到同一行、双双扣同一个 host，另一个 host 的空额被浪费。

---

## 2. 瓶颈分解

### 2.1 单 slot 串行（硬上限）

`wait-queue` 的 `slot_options: slot_on_head` 让每个 HEAD 节点只有 **1 个 slot**
（`cmd/controld/grpc/app_slots.go:58`），而 agent 的任务执行是串行的
（`cmd/agent/run.go:69` 的 `for _, task := range taskList`）。

单个 vtask 要消耗两次 CLI 进程启动 + 两次 gRPC：

- `scalebox vtask bind`：进程启动 + gRPC + 事务（O(H) 查询 + 排序 + UPDATE）
- `scalebox vtask add-subtask --direct`：进程启动 + gRPC + 插入

实测每次合计 40–80 ms → **全局上限 12–25 vtask/s，与 host 数无关**。

### 2.2 2 × poll 的节拍

`vtask_size:wait-queue = 1` 时，一个 vtask 要穿过两个轮询边界才能让下一个进入：

1. wait-queue 领取 → 执行 → 投 head
2. head 领取 → 执行 → 回流 main-router → `from-vtask-head.sh` 执行 → `+1`
3. wait-queue 下一次轮询才能领下一个

实测（p419，C 2 节点，poll=6 s）：wait-queue 相邻放行间隔 **11.60 s ≈ 2 × 6 s**。

全模块 `poll_interval_seconds=1` 后节拍降到 2.77 s（examples/vtask 三轮验证：
1232 → 770 → 344 s）。

### 2.3 需求侧对照

| 场景 | 需求 (vtask/s) | 现状供给 | 差 |
|---|---|---|---|
| 500 host × 2 容量，T=5 s | 200 | 15–20 | 10× |
| 500 host × 4 容量，T=5 s | 400 | 15–20 | 20× |
| 500 host × 2 容量，T=1 s | 1000 | 15–20 | 50× |
| 50 group × 4 容量，T=5 s | 40 | 15–20 | 2× |

### 2.4 一个判断

现有模型清晰的根本原因，是**把所有分配权集中在一处**。集中才能简单，简单才能清晰。
500 台 host 时，这个结构天然撞上单点上限——清晰与可扩展在这里是直接冲突的。

---

## 3. L1：兼容改进（保留 wait-queue 模型）

### 3.1 目标

**不改用户心智模型**：仍然是「进 wait-queue → 拿到 to_host → 进 head」。
只提升 wait-queue 本身的吞吐并去掉不必要的串行。

### 3.2 四项改动

#### 3.2.1 提前归还串行门

`dockerfiles/wait-queue/run.sh`：在 `add-subtask` 成功路径末尾各加一行。

```bash
# bind 与投递已完成，串行门到此结束——下一个 vtask 可以进入
scalebox semaphore increment "vtask_size:${module_name}"
```

HOST-BOUND / GROUP-BOUND 两个分支都要加。

同时**删掉** `examples/vtask/main-router/from-vtask-head.sh` 末尾的整段：

```bash
code=$?
if [ $code -eq 0 ]; then
    wait_queue_module=${WAIT_QUEUE_MODULE:-wait-queue}
    scalebox semaphore increment "vtask_size:${wait_queue_module}"
    code=$?
fi
exit $code
```

改为 `exit $?`。**不删会双重归还**：`vtask_size:wait-queue` 从 1 涨到 2 后不再回落。

gRPC 落点：应用自己的 router 若也做同样动作（如 app-fxcorr 的 `router/from_head.go`
两处 `semaphore.AddValue(fxtask.GateWaitQueue, 1, appID)`），同样删掉。

效果：串行门的持有时间从「head 全流程（跨 2 个轮询边界）」缩到
「bind + add-subtask（约 60 ms）」，节拍 `2×poll + head执行` → `1×poll`。

#### 3.2.2 扩 slot

`slot_on_head` 固定每 HEAD 节点 1 个 slot，改成显式配置：

```yaml
wait-queue:
  arguments:
    slot_options:          # 留空（不再用 slot_on_head）
  slots:
    - ${WQ_SLOTS}          # 例：h0:16  → h0 上 16 个 wait-queue slot
```

slot 数 M 决定并发度。M 个 agent 各自领取、各自 bind，互相独立
（`BindVtaskResource` 是独立 gRPC handler，不走 `GetTaskList` 的 per-module mutex）。

#### 3.2.3 `Decrement` 加行锁

`internal/semagroup/decrement.go`：把「SELECT 最大 → UPDATE」换成单语句：

```sql
WITH picked AS (
  SELECT name FROM t_semaphore
  WHERE name ~ $1 AND app = $2 AND value > 0
  ORDER BY value DESC, name
  LIMIT 1
  FOR UPDATE SKIP LOCKED
)
UPDATE t_semaphore s SET value = s.value - 1
FROM picked WHERE s.name = picked.name AND s.app = $2
RETURNING s.name, s.value
```

- `SKIP LOCKED`：并发 bind 各取不同 host，不互相阻塞
- `value > 0`：堵住扣成负数（现在只能靠 `check.sh` 的 gate 防）

#### 3.2.4 批量领取（可选）

`getVtaskBatchSize` 对非 head 模块返回 `min(semaValue, task_batch_size)`。
把 `task_batch_size` 从 1 调到 4/8，一次领取多个，摊销「进程启动 + gRPC」的固定开销。

代价：`updateVTask` 的扣减是 asyncbatch（10 ms 批处理），一次领多个时 gate 读值与
实际扣减之间的窗口放大，需要同步把批大小计入 gate（`n = min(value, TaskBatchSize)`
已经是这样，安全）。

> **这一项不是必须的**。前三项做完，M=16 已能到 400/s。批量领取只在 M 受限于
> 容器数时才有意义。

### 3.3 容量估算

单 slot 约 25 vtask/s（40 ms/个，保守）：

| M | 吞吐 | 覆盖 |
|---|---|---|
| 4 | 100/s | 500 host × 2 容量，T=10 s |
| 8 | 200/s | 500 host × 2 容量，T=5 s（临界） |
| 16 | 400/s | 500 host × 2 容量，T=5 s（2 倍余量） |
| 32 | 800/s | 500 host × 4 容量，T=5 s |

**你举的场景（500 host、5 s）M=16 即可，不需要引入 L2。**

### 3.4 L1 的边界

- **M 个 slot 全在 head 节点上**：物理上仍是单点。head 节点故障 → 分配停摆
- **O(H) 扫描仍在**：500 行排序约 0.1 ms，可忽略；10 万行时会成为问题
- **延迟仍是 2 × poll**：poll=1 时 2 s，对 1 s 的任务是 200% 开销
- **模型不变**：用户看到的还是 `wait-queue → head → core → tail`

---

## 4. L2：自扣模式（高性能）

### 4.1 机制

去掉 wait-queue 这一跳。任务以 `to_host` 为空的状态进入 **vtask-head 的队列**，
各 host 上的 head slot 竞争领取；领取时 controld 按 slot 所在 host 扣该 host 的额度。

```
main-router ──► [vtask-head 队列，to_host 为空]
                        │ 各 host 的 head slot 竞争领取（FOR UPDATE SKIP LOCKED）
                        ▼
                  host A 的 head slot
                        │ 领取时：扣 host_vtask_size:vtask-head:A
                        │         自动写入 _vtask_size_sema header
                        ▼
                  head 执行业务 → 投 core（to_host 由 pod_id/from_ip 推导 = A）
                        ▼
                  vtask-tail ──► doVTaskFinished → A 的额度 +1
```

关键点：**扣减、header 写入、归还全部由 controld 完成，应用脚本里一行都没有**。

### 4.2 为什么能并行

- 平台领取任务的 SQL 已经是 `FOR UPDATE SKIP LOCKED`（`cmd/controld/task/get.go:370`）
- `to_host IS NULL` 的任务任何 slot 都能领（`get.go:186`）
- 每台 host 的额度是一个独立信号量行，`SKIP LOCKED` 下互不干扰

分配从「中心挑一台」变成「谁空谁领」，500 台天然并行。

### 4.3 组件与配置

**唯一组件 = 现有的 vtask-head 模块**。

```yaml
vtask-head:
  parameters:
    task_dist_mode: HOST-BOUND
    vtask_role: head
    vtask_size: ${K}            # 每 host 在制上限
    vtask_assign_mode: self     # 新参数：启用自扣
    pod_id: mypod
    task_batch_size: 1          # 必须为 1（见 4.5）
  slots:
    - ${NODES}                  # 每 host 1 个
```

应用侧删除 wait-queue 模块；`main-router` 的入口分支直接把任务投给 vtask-head
（不再带 `to_host`）。

### 4.4 controld 改动清单

| # | 文件 | 改动 | 优先级 |
|---|---|---|---|
| 1 | `cmd/controld/task/cache_attrset.go` | `ModuleAttrSet` 加 `VtaskAssignMode`，SQL 读 `parameters->>'vtask_assign_mode'` | P0 |
| 2 | `cmd/controld/task/vtask.go:23` | `getVtaskBatchSize` 对自扣型 head 也走信号量 gate | P0 |
| 3 | `cmd/controld/grpc/app_service.go:736` | 参数白名单加 `vtask_assign_mode` | P0 |
| 4 | `cmd/controld/task/get.go:362` | 领取 SQL 原子写入 `_vtask_size_sema` | P0（见 4.5 风险 1） |
| 5 | `cmd/controld/grpc/vtask_ops.go:276` | `FailVtask` 的 gate 释放加条件 | P0（见 4.5 风险 2） |
| 6 | `cmd/controld/grpc/app_service.go` | 创建校验（见 4.5 的五条防呆） | P1 |
| 7 | `cmd/controld/task/get.go:172,190` | 日志降级（Info → Debug） | P2 |

`getVtaskBatchSize` 的改法：

```go
selfAssign := module.VtaskAssignMode == "self"
if module.VtaskSize <= 0 || (!selfAssign && module.VtaskRole == "head" &&
    (module.TaskDistMode == "HOST-BOUND" || module.TaskDistMode == "SLOT-BOUND")) {
    return module.TaskBatchSize
}
```

往下已有分支不动——按 slot 的 host 拼出 `host_vtask_size:<module>:<host>`，读值，
`n = min(value, TaskBatchSize)`，`n <= 0` 时 `getHostBoundRows` 返回 nil，slot 领不到。

**不需要改的**（已确认现成）：

| 机制 | 位置 |
|---|---|
| 空 `to_host` 任何 slot 可领 | `get.go:186` |
| 领取时按 slot 的 host 扣额度 | `vtask.go:87-100` |
| 自动写 `_vtask_size_sema` | `vtask.go:103-110` |
| head → core 的 to_host 推导 | `shouldDeriveToHostByPodID` |
| tail 完成归还 | `vtask.go:115` |
| 并发领取不重复 | `get.go:370` |
| `t_task(mod_id, status_code)` 索引 | `build/database/initdb.d/25-updated-schema.sql:48` |

### 4.5 设计审查：必须处理的风险

#### P0-1 `FailVtask` 在自扣模式下失效（额度永久泄漏）

`updateVTask` 写的 `_vtask_size_sema` **只在内存里**（`vtask.go:103-110`），不落库。
而 `FailVtask` 读的是数据库（`vtask_ops.go:298`）：

```sql
SELECT COALESCE(headers->>'_vtask_size_sema', '') FROM t_task WHERE id = $1
```

预扣模式下 header 由 wait-queue 通过 `--header` 写入、库里有；自扣模式下库里是空的
→ **head 任务失败时这一格额度永远不还**。

修复：把写入挪进领取 SQL，和 `status_code = -2` 同一个事务：

```sql
UPDATE t_task
SET status_code = -2,
    headers = CASE
      WHEN COALESCE(t_task.headers->>'_vtask_size_sema','') <> '' THEN t_task.headers
      ELSE jsonb_set(COALESCE(t_task.headers,'{}'::jsonb),
                     '{_vtask_size_sema}', to_jsonb($4::text), true)
    END
FROM cte WHERE t_task.id = cte.id
```

预扣模式下 header 已存在、`CASE` 不覆盖，行为不变。

#### P0-2 `FailVtask` 会误放别人的 gate

`vtask_ops.go:276-286` 无条件释放「同 app 下 `vtask_size > 0` 且 `vtask_role = ''`
的所有模块」的 `vtask_size:<name>`。自扣模式没有 wait-queue，但应用里若还有别的
DEFAULT 模块配了 `vtask_size`，会被莫名 +1，信号量越跑越高。

修复：只在应用存在预扣模块时执行，或随 wait-queue 一起删掉这段。

#### P0-3 超卖的三个入口

| 入口 | 后果 | 防呆 |
|---|---|---|
| 同 host 配了 2 个以上 head slot | gate 读值 → 10 ms 后 asyncbatch 才扣，两个 slot 都读到 1、都领走 | 校验 head slot 数 = 1 |
| `task_batch_size > 1` | 一次领 N 个，全扣同一个 host，gate 只看了「> 0」 | 校验自扣 head 的 `task_batch_size = 1` |
| `vtask_size` 没配或 ≤ 0 | 走 `VtaskSize <= 0` 第一分支，**完全不 gate** | 校验自扣 head 必须 `vtask_size > 0` |

超卖的后果不是崩溃，是 `host_vtask_size` 变负、该 host 要等更多 tail 完成才恢复，
表现为「某些节点莫名变慢」，很难定位。

#### P0-4 没有 tail 就是永久泄漏

`doVTaskFinished` 只对 `vtask_role = tail` 生效（`vtask.go:122`）。应用若把链设计成
head 直接收尾，额度扣了永远不还。修复：创建时校验「有 `vtask_role=head` 就必须有
`vtask_role=tail`」。

#### P1-5 500 个 slot 的空轮询会打满 controld 单核

**L2 没有消灭中心瓶颈，只是换了位置。**

每个 slot 每 poll 一次 `GetTaskList`，路径里串行经过 `moduleObj.mutex`
（`get.go:58`），锁内包含 `getVtaskBatchSize`（一次 `semaphore.GetValue` 查库）
+ 领取 SQL，约 1–3 ms。500 个 slot × poll=1 s：

**500 次/s × 2 ms = 1.0 s/s，单核饱和。**

`poll_interval_seconds` 是 int（`agent/init.go:234`），最小 1 秒，无法再降。

缓解（按代价排序）：
1. poll 调到 2–3 s（延迟换吞吐）
2. mutex 缩到只保护 `moduleObj` 状态，SQL 移出去
3. controld 多实例 + 按模块分片

#### P1-6 `isSlotDispatchable` 每次领取都查库

`get.go:31` 一次 `EXISTS` 查询，500 次/s 额外打库。可缓存或与 gate 合并。

#### P2-7 运维边界

- **节点下线 → 额度泄漏**：host 宕机，`host_vtask_size:<host>` 停在低值，无自动回收
- **信号量缺失 → 静默不领**：`GetValue` 出错时 `semaValue = 0`（fail-closed），
  新增 host 未走 `add-slots` 就会「这台机器什么都不干」
- **与预扣模式共存会乱**：两套 gate + P0-2 的误放。应创建时校验互斥

### 4.6 L2 的边界

- **吞吐 ~500/s**，受 controld mutex 封顶；**不解决 1000/s 以上的需求**
- **每 host 1 个 head slot**（多 slot 会超卖，见 P0-3）
- **head 必须短**：它跑在分配点上，慢了这台 host 就停
- **任务必须与 host 无关**：有数据亲和需求时不适用
- **轮询成本随 host 数线性增长**：500 slot = 500 次/s

---

## 5. 用户模型对比：清晰度的代价

### 5.1 现有模型为什么清晰

- **角色少**：wait-queue（分配）、vtask-head（入口）、vtask-tail（出口）
- **职责单一**：分配和执行分开，谁都能看懂
- **动作可见**：`vtask bind`、`--header _vtask_size_sema`、`semaphore increment`
  都在用户能打开、能改的脚本里
- **出错可见**：串行锁卡住时表现为「不动了」，容易发现

### 5.2 L2 模型的变化

| | 现有 | L2 |
|---|---|---|
| 模块数 | wait-queue + head + tail + core | head + tail + core |
| 分配在哪 | wait-queue 的 `run.sh` | controld 内部 |
| 分配动作 | `scalebox vtask bind` 显式 | 无（领取时顺手扣） |
| 额度从哪来 | `--header _vtask_size_sema` 显式传 | controld 自动注入 |
| 归还 | 脚本里的 `semaphore increment` | controld 自动，脚本零代码 |
| 新知识 | — | 新参数 + 5 条硬约束 + 3 个隐式机制 |

**净效果：用户要写的代码少了，要理解的东西没少。**

### 5.3 为什么不清晰

1. **`vtask-head` 的名字骗人**：它叫「head」（入口），实际还兼任「资源分配点」。
   用户配 `vtask_size: 2` 会以为是 head 的并发数，实际是「每 host 在制 vtask 上限」。
2. **分配动作不可见**：现有模型能逐行追到一个 vtask 如何落到某台机器；L2 全在
   controld 里，排查「为什么我的任务没被领走」只能看 `Filtered-by: host_vtask_size`。
3. **`vtask_role: head` 分叉成两种语义**：预扣型（不 gate）和自扣型（gate），
   靠一个参数切换。这是隐式分支点。
4. **出错是静默的**：少配 slot、忘配 tail、`task_batch_size` 没归 1，都表现为
   「某些节点变慢」或「某台机器后来不领活了」。

### 5.4 缓解措施（进 L2 时必须一并做）

1. 参数名自解释：`vtask_assign_mode: self` 不如叫 `vtask_capacity: local`
2. 五条硬约束做成**创建时强校验**，把隐式规则变成显式报错
3. `task_batch_size=1`、`vtask_size>0`、有 head 必有 tail、head slot 数 = 1、
   与 wait-queue 互斥——五条校验放在 `saveApp` 里
4. 文档里明确：L2 是**面向规模的模式**，不是默认模式；个位数 host 用现有模型

### 5.5 一个中间选项（已评估，不推荐）

保留 wait-queue 名字与两阶段结构，把它的 slot 分布到各 host，每个只 bind 本机
（需要新的 `vtask bind --local` 语义）。

- 优点：用户心智完全不变，分配动作仍显式
- 缺点：每个 vtask 多一跳 RPC；**轮询成本翻倍**（500 个 wait-queue slot + 500 个
  head slot，各 poll 一次）；需要改 `BindVtaskResource` 的 host 选择逻辑
- 结论：不划算。要分布化就干脆合并（L2），否则就集中（L1）

---

## 6. 测试方案

### 6.1 L1（examples/vtask，4 节点 bio 集群）

| 组 | 配置 | 判据 |
|---|---|---|
| 基线 | 现状（1 slot，poll=6） | 参考点 |
| L1-a | 提前归还（仍 1 slot） | 节拍 2×poll → 1×poll |
| L1-b | 16 slot + 提前归还 | 200 个 vtask 总时长 ≥ 3× 快于基线 |
| L1-c | 32 slot + 并发 bind | `vtask_size:wait-queue` 不越界、所有 `host_vtask_size` 跑完回初值 |
| L1-d | 500 个 vtask 压力 | 无重复领取、无任务丢失 |

### 6.2 L2（examples/vtask）

改造 `app.yaml`：删 wait-queue，head 加 `vtask_assign_mode: self`、`vtask_size: 2`，
每 host 1 slot；`main-router/default.sh` 的 HOST-BOUND 分支改成直接投 vtask-head。

| 判据 | 方法 |
|---|---|
| 无超卖 | 跑完所有 `host_vtask_size:*` 回到初值 |
| 无重复领取 | 每个 task 在 `t_task_exec` 里只有一条记录 |
| 无额度泄漏 | 强制 kill 一个 vtask-head 任务后，额度能被 `vtask fail` 或 tail 归位 |
| 节拍 | head 相邻执行间隔 < 1.5 × poll |
| 失败路径 | head 失败时 `FailVtask` 能释放（**验证 P0-1 的修复**） |

### 6.3 验证不了的

bio 只有 4 个计算节点，**验证不了 500 路的 controld 轮询成本**（P1-5）。
这一项只能靠推算，或用 100+ 节点实测。

---

## 7. 迁移路径

```
现状（wait-queue 1 slot，2×poll）
   │
   ├─► L1（4 处改动，模型不变）
   │      · 提前归还
   │      · 扩 slot 到 16
   │      · Decrement 加 SKIP LOCKED
   │      · （可选）批量领取
   │      覆盖：≤ 500 host、≥ 5 s 任务
   │
   └─► L2（8 处改动，模型改变）
          · 自扣模式
          · 覆盖：低延迟、无中心单点
          · 上限：~500/s（controld mutex）
          
   未来 L3（分布式）
          · 多 controld / 分片 wait-queue
          · 覆盖：> 1000 host、< 1 s 任务
```

**判断依据**：

| 你的场景 | 用哪层 |
|---|---|
| ≤ 10 host，任意任务 | 现状 |
| ≤ 500 host，≥ 5 s 任务 | **L1** |
| ≤ 500 host，1–5 s 任务 + 要低延迟 | L2 |
| > 1000 host 或 < 1 s 任务 | L3（需另设计） |

---

## 8. 附录：其他方案评估

| 方案 | 描述 | 结论 |
|---|---|---|
| ③ 去掉 wait-queue 一跳 | 把 bind 逻辑搬进 main-router | 不划算：router 是 DEFAULT 模块，`getVtaskBatchSize` 在 `VtaskSize<=0` 时返回 `TaskBatchSize`，一次领光、全部并发 bind、超卖；还要自己实现空额 gate。收益只有 1 次 RPC |
| 分片 wait-queue | M 个 wait-queue，按 host hash 分片 | 应用层要写路由、要知道 host 列表、热点不均；只在 L1 撑不住时才考虑 |
| 批量 bind | `BindVtaskResource` 一次返回 N 个 host | 粒度粗，任务流式到达时容易空占；L1 的 3.2.4 是它的轻量替代 |
| controld 内分配器 | controld 主动给 ready 任务写 `to_host` | 中心点没消失，只是从应用层搬进 controld；除非做成一条批量 SQL，否则和 wait-queue 一样串行。且要新增「待分配」状态码 |

---

## 9. 代码位置索引

| 机制 | 文件:行 |
|---|---|
| agent 轮询间隔（默认 6 s） | `cmd/agent/init.go:234` |
| agent 串行执行任务 | `cmd/agent/run.go:69` |
| `slot_on_head` 展开为每 HEAD 1 slot | `cmd/controld/grpc/app_slots.go:58` |
| 领取批量计算（gate） | `cmd/controld/task/vtask.go:23 getVtaskBatchSize` |
| 领取时扣减 + 写 header | `cmd/controld/task/vtask.go:62 updateVTask` |
| tail 完成归还 | `cmd/controld/task/vtask.go:115 doVTaskFinished` |
| 领取 SQL（SKIP LOCKED） | `cmd/controld/task/get.go:362 sqlFmt` |
| HOST-BOUND 的 `to_host IS NULL` 分支 | `cmd/controld/task/get.go:186` |
| GROUP-BOUND 的组前缀过滤 | `cmd/controld/task/get.go:222` |
| per-module mutex | `cmd/controld/task/get.go:58` |
| `Decrement`（无行锁） | `internal/semagroup/decrement.go:24` |
| `BindVtaskResource` | `cmd/controld/grpc/vtask_ops.go:20` |
| `FailVtask` 的 gate 释放 | `cmd/controld/grpc/vtask_ops.go:276` |
| `t_task(mod_id,status_code)` 索引 | `build/database/initdb.d/25-updated-schema.sql:48` |
| wait-queue 标准模块 | `dockerfiles/wait-queue/run.sh`、`check.sh` |
| 示例链 | `examples/vtask/main-router/*.sh` |
