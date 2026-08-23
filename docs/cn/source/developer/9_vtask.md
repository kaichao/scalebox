# 9. VTask 虚拟任务

**VTask**（Virtual Task）是 Scalebox 在细粒度 task 之上建立的应用级粗粒度计算单元。一个 vtask 是跨模块的 task 集合，内置信号量（semaphore）和共享变量（variable），由 wait-queue（可选）、vtask-head、一个或多个级联 vtask-core、vtask-tail 组成管道，统一管理流控、资源绑定、状态追踪和全链路上下文传播。

## 9.1 为什么需要 VTask

### 9.1.1 task 模型的局限

Scalebox 的基础模型是 module → slot → task。每个 task 是一次独立的执行单元，在细粒度调度层面工作良好，但在应用层面有三个缺口：

1. **缺少跨模块状态管理**：同一 App 内多个 module 协作时，各 module 内的 task 彼此孤立，没有"这批 task 属于同一个计算任务"的概念
2. **无法回答应用级问题**："某个计算任务执行到哪了"、"一共多少子任务、完成多少"、"失败了该释放哪些资源"
3. **缺少资源隔离**：不同计算任务在同一节点上执行时无资源预留机制

### 9.1.2 VTask 的解决方案

VTask 在 task 之上叠加跨模块状态管理：

```
细粒度层（task）：  单次执行单元，slot 调度，状态码 -1/-2/-3/0/>0
粗粒度层（vtask）：  跨模块 task 集合，流控 + 状态 + 容错 + 资源绑定
```

| 能力 | 实现 | 效果 |
|------|------|------|
| 跨模块关联 | `vtask` 列 + `_vtask_id` header 串联 head→core→tail | 脚本无需手动维护 task 间引用 |
| 流控 | `vtask_size` 信号量限制并发数 | 不需要额外模块做排队逻辑 |
| 状态追踪 | `vtask get` / `vtask list` 实时展示进度 | 一行命令取代手动查 task 表 |
| 粗粒度容错 | `vtask fail` 终止管道 + 释放资源 | 整 vtask 重试，不需逐 task 清理 |
| 资源绑定 | bind/unbind 预留/归还计算节点 | 声明式管理，复用而非硬编码 |
| 上下文透传 | `_` 前缀 header 全链路自动传播 | 中间模块零代码感知上下文来源 |

## 9.2 概念模型

### 9.2.1 三种类型

| | DEFAULT | HOST-BOUND | GROUP-BOUND |
|----|---------|------------|-------------|
| 资源绑定 | 无 | 单节点 | 节点组 |
| wait-queue | 无 | 有 | 有 |
| 流控信号量 | `vtask_size:<mod>` | `host_vtask_size:<mod>:<host>` | `group_vtask_size:<mod>:<组号>` |
| vtask-head 模式 | 不设 task_dist_mode | `HOST-BOUND` | `SLOT-BOUND` |
| vtask-core 模式 | 不设 task_dist_mode | `HOST-BOUND` | `HOST-BOUND` |
| 适用场景 | 轻量批处理 | 单节点独占计算 | 多节点协同计算 |

**管道结构**（`vtask-core` 可级联多个阶段）：

```
DEFAULT:       task → vtask-head → vtask-core-1 → vtask-core-2 → ... → vtask-tail
HOST-BOUND:    task → wait-queue → vtask-head → vtask-core → vtask-tail
GROUP-BOUND:   task → wait-queue → vtask-head → vtask-core-1 → vtask-core-2 → ... → vtask-tail
```

### 9.2.2 能力矩阵

| 能力 | DEFAULT | HOST-BOUND | GROUP-BOUND |
|------|:--:|:--:|:--:|
| 流控 + 追踪 + 容错 | ✅ | ✅ | ✅ |
| header DAG 透传 | ✅ | ✅ | ✅ |
| 确定性性能（节点独占） | ❌ | ✅ | ✅ |
| 本地 I/O（NVMe） | ❌ | ✅ | ✅ |
| 故障隔离 | ❌ | ✅ | ✅ |
| vtask 多节点并行 | ❌ | ❌ | ✅ |

| 类型 | 定位 | 典型场景 |
|------|------|---------|
| DEFAULT | 敏捷模式 | 日志汇总、配置下发、批量缩略图、数据格式转换 |
| HOST-BOUND | 性能模式 | 单节点 GPU 推理、大规模数据变换 |
| GROUP-BOUND | 规模模式 | 多通道波束合成、多节点协同数据分析 |

### 9.2.3 两层抽象

| 层级 | 概念 | 关注点 |
|------|------|--------|
| **vtask**（应用层） | 资源绑定形态：HOST-BOUND / GROUP-BOUND / DEFAULT | 用户视角：跑在单节点还是节点组 |
| **module**（底层） | task 路由方式：HOST-BOUND / SLOT-BOUND / DEFAULT | 框架视角：task 如何路由到具体 slot |

`SLOT-BOUND` 是模块级实现细节，只在 GROUP-BOUND 的 head 模块上使用，不存在"纯 SLOT-BOUND vtask"。

**环境变量分层**（与两层抽象一一对应）：

| 变量 | 层次 | 值域 | 消费方 |
|------|------|------|--------|
| `VTASK_MODE` | vtask 应用层 | `DEFAULT` / `HOST-BOUND` / `GROUP-BOUND` | 应用脚本（main-router 路由、wait-queue 标准模块） |
| `TASK_DIST_MODE` | 平台实现层 | `HOST-BOUND` / `SLOT-BOUND` / 空 | 仅驱动 vtask-head 模块的 `task_dist_mode` 平台参数 |

映射关系收敛在应用 env 文件一处（如 `group-bound.env`）：

```bash
VTASK_MODE=GROUP-BOUND
TASK_DIST_MODE=SLOT-BOUND   # 仅驱动 vtask-head 的平台参数，实现层枚举
```

应用作者在 vtask 层只接触 `VTASK_MODE`，`SLOT-BOUND` 仅在实现层出现。

### 9.2.4 计算资源标识

vtask bind 分配的资源本质是**计算位置（placement）**——vtask 锚定到某台主机（HOST-BOUND）或某个节点组（GROUP-BOUND），之后 task 通过 `to_host` / `to_slot` 路由到该位置上的 slot 执行。平台不感知物理容量（CPU 核数、内存量），**容量由应用作者通过 `vtask_size` 参数配置**（该位置同时承载的在制 vtask 数）。

**资源池来源**：可用资源 = head 模块（`vtask_role=head`）的 slot 展开结果：

| vtask 类型 | 资源池 | 标识 | 示例 |
|-----------|--------|------|------|
| HOST-BOUND | head 模块 `slots` 正则展开的每台主机（`t_host` 中 `status='ON'`） | 短主机名（去 cluster 后缀） | `n0-0` |
| GROUP-BOUND | head 模块的每个 slot（slot 列表长度即组数） | 组号 = head slot 的 `seq` | `0` |

信号量在 app 创建时建立、动态扩 slot 时补建，初值 = `vtask_size`。

**标识体系**（信号量命名即标识）：

| 类型 | 信号量名 | 标识部分 |
|------|---------|---------|
| DEFAULT | `vtask_size:<module>` | 无（全局） |
| HOST-BOUND | `host_vtask_size:<module>:<hostname>` | 短主机名 |
| GROUP-BOUND | `group_vtask_size:<module>:<组号>` | head slot 的 seq |

`vtask bind` 返回的 resource 即信号量名后缀；`vtask get` 显示 `host=n0-0` / `group=2`。

**节点组的成员没有平台级定义**：组内包含哪些节点、组内如何分发，是应用业务逻辑（vtask-head 脚本路由 + vtask-core 以 HOST-BOUND 部署在组内节点上）。平台只管理"组"这个锚点维度。

**资源占用与释放**：占用 = `vtask bind` 原子扣减信号量 + `_vtask_size_sema` header 随 root task 传播；释放 = `doVTaskFinished` / `FailVtask` 按 header increment 加回（错误回滚用 `vtask unbind`）。每个主机/组的在制 vtask 数 ≤ `vtask_size`。

### 9.2.4 根 task 标识

`vtask_role=head` 模块的 task 为 vtask 根。controld 创建 head 模块 task 时自动建立自引用：

```text
UPDATE t_task SET vtask = $1,
    headers = jsonb_set(headers, '{_vtask_id}', to_jsonb($1::text), true)
WHERE id = $1
```

`ListVTasks` 查询条件为 `WHERE vtask = id`。

### 9.2.5 状态派生

vtask 状态由根 task 的 `status_code` + 子任务完成比例派生：

| 根 task status | 子任务 | 派生状态 |
|:---:|------|------|
| `-1` | 0 | READY |
| `-2` | 0 | QUEUED |
| `-3` | 0 | RUNNING |
| `-3` | >0，部分完成 | RUNNING（X/Y finished） |
| `0` | 全部完成 | FINISHED（Y tasks） |
| `1` | — | FAILED（exit=N） |

## 9.3 模块结构

### 9.3.1 DEFAULT（4 模块）

```yaml
modules:
  main-router:   # 入口路由：task add --sink-module=vtask-head $body
  vtask-head:    # vtask_role: head, vtask_size: 2
  vtask-core:    # vtask_role: core
  vtask-tail:    # vtask_role: tail（无 unbind）
```

### 9.3.2 HOST-BOUND（5 模块）

```yaml
modules:
  main-router:   # task add --sink-module=wait-queue $body
  wait-queue:    # 标准模块 scalebox.net/platform/wait-queue（vtask_size: 1 串行化门控 + check.sh 准入 + run.sh bind）
  vtask-head:    # vtask_role: head, task_dist_mode: HOST-BOUND
  vtask-core:    # vtask_role: core, task_dist_mode: HOST-BOUND
  vtask-tail:    # vtask_role: tail
```

**关键脚本操作**：

| 脚本 | 位置 | 操作 |
|------|------|------|
| `check.sh` / `run.sh` | wait-queue 标准模块镜像 | 准入闸门（ACTION_CHECK）+ `vtask bind` → `vtask add-subtask --direct`（失败时 unbind + gate increment 回滚）；task 在本 slot 终结，作者无需编写 |
| `from-vtask-head.sh` | main-router/ | `vtask add-subtask --module=vtask-core`（async），最后 `semaphore increment vtask_size:wait-queue` |
| `from-vtask-core.sh` | main-router/ | `task add --sink-module=vtask-tail $1` |
| `from-vtask-tail.sh` | main-router/ | 空操作（额度由 controld `doVTaskFinished` 自动归还） |

### 9.3.3 GROUP-BOUND

与 HOST-BOUND 结构相同，差异：
- `vtask-head.task_dist_mode: SLOT-BOUND`（平台实现值），bind 返回组号（即 head slot 的 seq）
- 流控信号量格式 `group_vtask_size`，head slot 作为节点组锚点
- `from-vtask-head.sh` 按 `VTASK_MODE=GROUP-BOUND` 分支计算 `to_host`（业务逻辑）

### 9.3.4 vtask-core 级联

`vtask-core` 可包含多个计算阶段：

```yaml
modules:
  vtask-core-beamforming:     # 数据分发 / 预处理
    vtask_role: core
    task_dist_mode: HOST-BOUND
  vtask-core-cross-correlate: # 核心计算
    vtask_role: core
    task_dist_mode: HOST-BOUND
  vtask-core-merge:           # 结果合并
    vtask_role: core
    task_dist_mode: DEFAULT
```

级联规则：
- 所有 core 模块共享 `vtask_role: core`，controld 按统一逻辑处理
- `_` 前缀 header（`_vtask_id`、`_vtask_size_sema`）沿 core 链自动传播
- 每个 core 模块独立配置 slot 数量、`task_dist_mode` 和镜像

## 9.4 信号量机制

### 9.4.1 命名规范

| 分发模式 | 格式 | 示例 |
|---------|------|------|
| DEFAULT | `vtask_size:<module>` | `vtask_size:vtask-head` |
| HOST-BOUND | `host_vtask_size:<module>:<host>` | `host_vtask_size:vtask-head:n0-0` |
| SLOT-BOUND（GROUP-BOUND vtask） | `group_vtask_size:<module>:<组号>` | `group_vtask_size:vtask-head:2` |

### 9.4.2 统一信号量

每个流控信号量只有一套（无冒号前缀），由 `vtask bind` 预扣、controld 自动加回，脚本不再参与资源归还：

| 管理者 | 操作 | 时机 |
|--------|------|------|
| wait-queue `check.sh`（ACTION_CHECK） | 读 semagroup max，≤ 0 时返回非 0 → agent 停领 task 并轮询重试 | agent 每轮 poll、`GetTaskList` 之前 |
| `vtask bind`（`BindVtaskResource`） | 原子 decrement（semagroup 选值最大成员，返回组号/hostname） | wait-queue 串行准入时 |
| controld `updateVTask` | 带 `_vtask_size_sema` 的 task 已预扣，领取时跳过；仅无 header 裸任务兜底扣减 | head slot 拾取 task 时 |
| controld `doVTaskFinished` / `FailVtask` | increment（按 `_vtask_size_sema` 定位） | vtask-tail 完成 / vtask 失败时 |

> 历史版本的双信号量（流控版 + `:` 前缀可编程版，`vtask_size_sema_copy: yes`）已废弃，`vtask_size_sema_copy` 参数不再生效。

### 9.4.3 流控流程

```
wait-queue slot 每轮 poll
  → check.sh: semagroup max ≤ 0 → exit 1，agent 停领（task 从未被拾取）
  → max > 0 → 领取 task → run.sh（标准模块镜像内）
  → vtask bind: semagroup.Decrement 原子预扣（返回 hostname/组号）
  → vtask-head task 拾取：updateVTask 跳过已预扣 task（DEFAULT 裸任务兜底扣）
  → vtask 管道执行
  → vtask-tail 完成 → doVTaskFinished: increment（恢复配额）
  → check.sh 下一轮轮询放行
```

### 9.4.4 等待队列门控

wait-queue 模块的 `vtask_size` 保持 1（串行化）。`from-vtask-head.sh` 末尾手动 increment 放行下一个 vtask。

## 9.5 管道流程

### 9.5.1 header 传播

```
main-router task add
  → agent addSinkTasks: _vtask_* 门控（仅 vtask_role≠"" 的 sink 模块通过）
  → wait-queue task: 不传入 _vtask_id / _vtask_size_sema

wait-queue task 执行标准模块 run.sh
  → --direct gRPC 创建 vtask-head task
  → vtask_role=head → self-ref vtask=id
  → _vtask_size_sema 由 run.sh 显式设置
  → task 在 wait-queue slot 终结（不写 sink-tasks.txt，无回流）

vtask-head task 执行 from-vtask-head.sh
  → sink-tasks.txt async 创建 vtask-core task
  → agent addSinkTasks: _ 前缀 header 自动传播

vtask-core → vtask-tail（同上路径）

vtask-tail 完成
  → doVTaskFinished: 读取 _vtask_size_sema → increment 流控版
  → from-vtask-tail.sh: 读取 _vtask_size_sema → unbind 可编程版
```

### 9.5.2 生命周期

```
task 创建 → ON CONFLICT 防重
  ↓
task 调度 → FOR UPDATE SKIP LOCKED（单 slot 独占）
  ↓
task 执行 → agent → status_code: RUNNING → OK/ERROR
  ↓  _vtask_id / _vtask_size_sema 通过 _ 前缀自动传播
  ↓
tail 完成 → doVTaskFinished → semaphore increment
  ↓  headers.vtask_tail_id = tail task ID
  ↓
from-vtask-tail.sh → unbind（仅 HOST-BOUND / SLOT-BOUND）
```

### 9.5.3 add-subtask 双路径

`vtask add-subtask` 在 agent 容器内运行时，默认走 sink-tasks.txt 异步路径，加 `--direct` 切换到 gRPC 同步路径：

| | 异步（sink-tasks.txt） | 同步（gRPC `--direct`） |
|--|:---:|:---:|
| 执行时机 | agent 后处理（task 执行完毕后批量发送） | 脚本内同步调用 |
| 错误感知 | 脚本无感知（写入成功 → exit 0） | 脚本立即获知（exit code ≠ 0） |
| `_` 前缀 header | agent 自动传播 | 调用方显式传 `--header` |
| `from_module`/`from_ip` | agent 注入（更准确） | CLI 从 env 注入 |

**选择原则**：

| 脚本 | 路径 | 原因 |
|------|------|------|
| `from-vtask-head.sh` | 异步（默认） | 无资源绑定前置操作，享受 header 自动传播 |
| wait-queue 标准模块 run.sh | 同步（`--direct`） | bind 已扣除资源，须同步感知成败，失败时回滚 |
| 运维手工 CLI | 同步（gRPC） | agent 外无 sink-tasks.txt 机制 |

## 9.6 命令参考

### 9.6.1 三层操作体系

```
Layer 1: 基础资源操作（全局作用域）
  semaphore create / get / increment / decrement / list / delete
  semagroup min / max / diffmin / diffmax / increment / decrement
  variable get / set / delete / list
  global get / set / delete / list

Layer 2: vtask 作用域资源操作
  vtask get-variable / set-variable
  vtask create-semaphore / get-semaphore / add-semaphore-value / delete-semaphore

Layer 3: vtask 高级操作
  vtask bind / unbind / add-subtask / get / fail / list / list-subtasks
```

### 9.6.2 vtask list / list-subtasks

```bash
# 列出应用所有 vtask（含派生状态和子任务计数）
scalebox vtask list --app-id=77

# 列出 vtask 的子任务
scalebox vtask list-subtasks --vtask-id=1964
```

### 9.6.3 vtask get

```bash
scalebox vtask get --vtask-id=42
# ID:        42
# App ID:    5
# Status:    RUNNING (2/3 finished)
# Sema:      host_vtask_size:vtask-head:n0-0
# Subtasks:  3 total, 2 finished
```

### 9.6.4 vtask fail

```bash
scalebox vtask fail --vtask-id=42
```

执行逻辑：
1. 原子化 `UPDATE status_code='1' WHERE id=$1 AND status_code IN ('-1','-2','-3')`
2. `RowsAffected==0` → 已结束，返回 user_msg
3. 释放 wait-queue 门控信号量
4. 释放资源信号量（`semaphore increment :<sema_name>`）
5. 级联标记未完成子任务

### 9.6.5 vtask bind / unbind

```bash
# 绑定（自动查找 head 模块，返回 hostname 或组号）
to_host=$(scalebox vtask bind) || exit 1

# 解绑（错误回滚时传信号量名；正常完成由 controld 自动加回）
scalebox vtask unbind --sema-name="host_vtask_size:vtask-head:n0-0"

# 解绑（运维传 vtask-id，服务端自动查 _vtask_size_sema）
scalebox vtask unbind --vtask-id=42
```

| vtask 类型 | bind 返回 |
|-----------|----------|
| HOST-BOUND | hostname（如 `n0-0`） |
| GROUP-BOUND | 组号（head slot 的 seq，如 `0`） |
| DEFAULT | 空（no-op） |

### 9.6.6 vtask add-subtask

```bash
# Agent 内（常见）：异步路径，_ 前缀 header 自动传播
scalebox vtask add-subtask --module=vtask-core --header to_ip=$from_ip $body

# Agent 内需同步感知错误（wait-queue 标准模块 run.sh）：加 --direct
scalebox vtask add-subtask --direct --module=vtask-head \
    --header to_host=$to_host --header _vtask_size_sema=$sema_name $body || {
    scalebox vtask unbind --sema-name="${sema_name}"
    scalebox semaphore increment vtask_size:wait-queue
    exit 1
}

# Agent 外运维手工：须显式传 --app-id 和 _vtask_id
scalebox vtask add-subtask --app-id=10 --module=vtask-core \
    --header _vtask_id=42 $body
```

自动注入：`from_module`（`PLAT_MODULE_NAME` env）、`from_ip`（本机 IPv4）。调用方显式传入时不覆盖。

### 9.6.7 vtask 作用域变量/信号量

```bash
# vtask 作用域变量
scalebox vtask get-variable --vtask-id <vtask-id> <var-name>
scalebox vtask set-variable --vtask-id <vtask-id> <var-name> <value>

# vtask 作用域信号量
scalebox vtask create-semaphore --vtask-id <vtask-id> <sema-name> <initial-value>
scalebox vtask get-semaphore --vtask-id <vtask-id> <sema-name>
scalebox vtask add-semaphore-value --vtask-id <vtask-id> <sema-name> <delta>
scalebox vtask delete-semaphore --vtask-id <vtask-id> <sema-name>
```

## 9.7 并发安全

| 机制 | 防护场景 |
|------|---------|
| `FOR UPDATE SKIP LOCKED` | slot 不会拿同一个 task |
| `ON CONFLICT` | 同 body 的 task 不重复创建 |
| 死锁重试 ×3 | `40P01` deadlock 自愈 |
| asyncbatch（10ms/50ms） | 信号量批量写入不逐条竞争 |
| `vtask_size=1`（wait-queue） | bind 操作串行化 |

## 9.8 运维场景速查

| 场景 | 命令 |
|------|------|
| 查看全局 vtask 进度 | `vtask list --app-id=N` |
| 查看单个 vtask 详情 | `vtask get --vtask-id=N` |
| 查看 vtask 子任务分布 | `vtask list-subtasks --vtask-id=N` |
| 终止单个 vtask | `vtask fail --vtask-id=N` |
| 批量终止 | `vtask list \| awk ... \| xargs vtask fail` |
| 暂停全局流控 | `semaphore decrement vtask_size:<mod> N` |
| 手工释放泄漏资源 | `vtask unbind --sema-name="host_vtask_size:..."` |
| 手工归还（不知信号量名） | `vtask unbind --vtask-id=N` |

## 9.9 实现索引

| 层 | 文件 | 职责 |
|----|------|------|
| Proto | `proto/app.proto` | 5 个 vtask RPC + 8 个 message |
| 服务端 | `cmd/controld/grpc/vtask_ops.go` | bind / unbind / add-subtask / get / fail handler |
| 服务端 | `cmd/controld/grpc/app_service.go` | ListVTasks / ListVTaskSubtasks handler |
| 服务端 | `cmd/controld/task/vtask.go` | `getVtaskBatchSize` / `updateVTask` / `doVTaskFinished` |
| CLI | `cmd/cli/vtask/vtask.go` | 14 个子命令 + add-subtask 异步双路径 |
| Agent | `cmd/agent/sink_tasks.go` | `addSinkTasks`: `_` 前缀 header 传播 + vtask_role 门控 |
| SDK | `pkg/vtask/` | Go 封装 |
| 示例 | `examples/vtask/main-router/` | 5 个脚本 |

> 测试详情和已知限制见 :doc:`VTask 设计审查 <../../go-scalebox/docs/vtask-review>`（英文）。
