# VTask 应用示例

> 完整设计文档见 [go-scalebox docs/vtask-design.md](../../go-scalebox/docs/vtask-design.md)。

## 概念

**VTask**（Virtual Task）是 scalebox 在细粒度 task 之上建立的应用级粗粒度计算单元。一个 vtask 是跨模块的 task 集合，内置信号量（semaphore）和共享变量（variable），由 wait-queue（可选）、vtask-head、一个或多个级联的 vtask-core、vtask-tail 组成管道，统一管理流控、资源绑定、状态追踪和全链路上下文传播。

### 三种类型

| 类型 | 资源绑定 | 流控信号量 | 适用场景 |
|------|---------|-----------|---------|
| **DEFAULT** | 无 | `vtask_size:vtask-head` | 轻量批处理、日志汇总、数据格式转换 |
| **HOST-BOUND** | 单节点 | `host_vtask_size:vtask-head:<hostname>` | 单节点独占计算、GPU 推理、超长计算 |
| **GROUP-BOUND** | 节点组 | `group_vtask_size:vtask-head:<组号>` | 多节点协作、分布式训练、大规模并行 FFT |

### 两层抽象

| 层级 | 概念 | 取值 | 关注点 |
|------|------|------|--------|
| **vtask（应用层）** | 资源绑定形态 | HOST-BOUND / GROUP-BOUND / DEFAULT | 用户视角：vtask 跑在单节点还是节点组 |
| **module（底层）** | task 路由分发方式 | HOST-BOUND / SLOT-BOUND / 无 | 框架视角：task 如何路由到具体 slot |

**映射关系**：
- HOST-BOUND vtask：所有模块的 `task_dist_mode=HOST-BOUND`
- GROUP-BOUND vtask：vtask-head 用 `task_dist_mode=SLOT-BOUND`（head slot 为节点组锚点），vtask-core 用 `HOST-BOUND`

用户只需理解三种 vtask 类型。SLOT-BOUND 是模块级实现细节，不在 vtask 层面暴露。

**环境变量分层**（与两层抽象一一对应）：

| 变量 | 层次 | 值域 | 消费方 |
|------|------|------|--------|
| `VTASK_MODE` | vtask 应用层 | `DEFAULT` / `HOST-BOUND` / `GROUP-BOUND` | main-router 脚本、wait-queue 标准模块 |
| `TASK_DIST_MODE` | 平台实现层 | `HOST-BOUND` / `SLOT-BOUND` / 空 | 仅驱动 vtask-head 模块的 `task_dist_mode` 平台参数 |

映射关系收敛在 env 文件一处（如 `group-bound.env`）：`VTASK_MODE=GROUP-BOUND` 与 `TASK_DIST_MODE=SLOT-BOUND` 成对出现。应用作者在 vtask 层只接触 `VTASK_MODE`。

### 计算资源标识

vtask bind 分配的资源本质是**计算位置（placement）**：vtask 锚定到某台主机（HOST-BOUND）或某个节点组（GROUP-BOUND），之后 task 通过 `to_host` / `to_slot` 路由到该位置上的 slot 执行。平台不感知物理容量（CPU 核数、内存量），**容量由应用作者通过 `vtask_size` 参数配置**。

**资源池来源**：可用资源 = head 模块（`vtask_role=head`）的 slot 展开结果：

| vtask 类型 | 资源池 | 标识 | 示例 |
|-----------|--------|------|------|
| HOST-BOUND | head 模块 `slots` 正则展开的每台在线主机 | 短主机名（去 cluster 后缀） | `n0-0` |
| GROUP-BOUND | head 模块的每个 slot（slot 列表长度即组数） | 组号 = head slot 的 `seq` | `0` |

`vtask bind` 返回的 resource 即信号量名后缀；`vtask get` 显示 `host=n0-0` / `group=2`。**节点组的成员没有平台级定义**——组内哪些节点、组内如何分发，由应用业务逻辑决定（`from-vtask-head.sh` 的路由 + vtask-core 以 HOST-BOUND 部署在组内节点上）；平台只管理"组"这个锚点维度。

### 统一信号量机制

每个流控信号量只有一套（无冒号前缀），由 `vtask bind` 预扣、controld 自动加回：

| 管理者 | 操作 | 时机 |
|--------|------|------|
| wait-queue `check.sh`（slot 级准入闸门，ACTION_CHECK） | 读 semagroup max，≤ 0 时返回非 0 → agent **停领 task** 并轮询重试 | agent 每轮 poll、`GetTaskList` 之前 |
| `vtask bind`（服务端 `BindVtaskResource`） | 原子 decrement（semagroup 选值最大成员，返回组号/hostname） | vtask 分配资源时（wait-queue 串行准入） |
| controld `updateVTask` | 带 `_vtask_size_sema` header 的 task 已预扣，领取时跳过；仅无 header 裸任务兜底扣减 | head slot 拾取 task 时 |
| controld `doVTaskFinished` | increment（读 tail task 的 `_vtask_size_sema`） | vtask-tail 模块 task 完成时 |
| controld `FailVtask` | increment（读 root task 的 `_vtask_size_sema`） | vtask 失败时 |

- check 准入闸门保证 bind 发生时必有空额：在制 vtask 数 ≤ `vtask_size`，信号量值 ∈ `[0, vtask_size]`，不会为负
- `vtask bind` 的原子扣减兼具两个职责：流控（在制数 ≤ `vtask_size`）与资源分配（返回绑定的组号/hostname）
- HOST-BOUND/GROUP-BOUND 的 head 模块领取**不再受信号量值 gate**，否则 bind 预扣后值=0，队列中的已绑任务无法领取（DEFAULT 模式无 bind，保留读值 gate）
- `vtask_size_sema_copy` 参数已废弃：不再创建冒号前缀的可编程副本，bind/unbind 与 controld 流控操作同一信号量

## 模块结构

示例包含以下模块：

```
main-router → default.sh（入口路由）
  │
  ├─ DEFAULT: 直连 vtask-head
  └─ HOST-BOUND / GROUP-BOUND: 进入 wait-queue
                                  │
                                  ▼
                            wait-queue（标准模块镜像，准入控制 + 资源分配）
                              check.sh（ACTION_CHECK，领取前执行）:
                                ① semagroup max ≤ 0 → exit 1，agent 停领并轮询
                                ② max > 0 才放行领取 task
                              run.sh（镜像内置，无需作者编写）:
                                ① vtask bind（原子化绑定资源）
                                ② vtask add-subtask --direct（gRPC 直连，需同步感知错误）
                                ③ 失败时 vtask unbind 回滚
                                ④ task 在本 slot 终结，不回流 main-router
                                  │
                                  ▼
                            vtask-head（vtask 根标记 + 业务路由）
                              from-vtask-head.sh:
                                ① vtask add-subtask（异步路径，_ 前缀自动传播）
                                ② gate increment（放行下一个 vtask）
                                  │
                                  ▼
                            vtask-core（核心计算，可多级联）
                              from-vtask-core.sh
                                  │
                                  ▼
                            vtask-tail（完结标记）
                              controld: doVTaskFinished
                                ① 读 tail task 的 _vtask_size_sema
                                ② increment 流控信号量（脚本无需操作）
```

| 模块 | 部署位置 | 关键参数 | 职责 |
|------|---------|---------|------|
| main-router | head 节点 (h0) | `main_router: main-router` | 入口路由，按 `VTASK_MODE` 分发 |
| wait-queue | head 节点 (h0) | `vtask_size: 1` | 标准模块（`scalebox.net/platform/wait-queue`）：slot 级准入闸门（check.sh）+ 串行化 + `vtask bind` 分配资源 |
| vtask-head | 计算节点 / head 节点 | `vtask_role: head`, `vtask_size: 2` | vtask 根标记 + 业务路由 + gate 释放 |
| vtask-core | 计算节点 | `vtask_role: core` | 核心计算（可多级联，如 vtask-core-1 → vtask-core-2） |
| vtask-tail | head 节点 (h0) | `vtask_role: tail` | 完结标记（controld 自动归还资源） |

## 模块详解

### main-router（入口路由）

main-router 是 app 的入口点，`run.sh` 根据 `from_module` 为空时调用 `default.sh`，按 `VTASK_MODE` 做路由决策：

```
DEFAULT → 直连 vtask-head（跳过资源分配）
HOST-BOUND / GROUP-BOUND → 进入 wait-queue（等待资源分配）
```

注意：`default.sh` 是在 main-router slot 内执行的第一个脚本，此时 task 尚未绑定任何 vtask。`VTASK_MODE` 环境变量从 `app.yaml` 的 `environments` 注入。

### wait-queue（准入控制 + 资源分配，标准模块）

wait-queue 是整个 vtask 管道的准入控制点，**全部逻辑封装在标准模块镜像 `scalebox.net/platform/wait-queue` 中**，应用作者只需在 app.yaml 引用镜像并配置 `VTASK_MODE` / `vtask_size`，无需编写任何脚本。分两层：

**第一层：slot 级准入闸门（`check.sh`，ACTION_CHECK）**

agent 主循环在每次 `GetTaskList` 领取 task **之前**执行 `check.sh`（见 go-scalebox `cmd/agent/run.go`），实现"无空额不领取"：

```bash
# check.sh 按 VTASK_MODE 选信号量组前缀（HEAD_MODULE 可自定义 head 模块名，缺省 vtask-head）
# HOST-BOUND:  ^host_vtask_size:vtask-head
# GROUP-BOUND:  ^group_vtask_size:vtask-head
ret=$(scalebox semagroup max "$sema_prefix")
max_value=${ret##*:}
[ "$max_value" -le 0 ] && exit 1   # 额度耗尽 → agent 停领，sleep 一个 poll 间隔后重试
exit 0                              # 有空额 → 放行领取
```

- check 发生在领取**之前**：exit 非 0 时 task 从未被拾取，不存在失败任务，也不需要重试机制——agent 只是这一轮 poll 不调 `GetTaskList`
- 额度由 `doVTaskFinished` 在 tail 完成时自动加回，check 轮询到 max > 0 后自然放行
- 作用：保证 bind 发生时必有空额 → 在制 vtask 数 ≤ `vtask_size`，信号量值不为负
- DEFAULT 模式不经 wait-queue，check.sh 直接 exit 0

**第二层：串行化 gate（`vtask_size: 1`）**

核心参数 `vtask_size: 1` 保证了**同一时刻只有一个 vtask 在执行 bind**，这是并发安全的前提——bind 内部（`semagroup.Decrement`）当前未加 `FOR UPDATE`，依赖单线程串行化保证正确性。

镜像内置 `run.sh` 的执行流程（wait-queue slot 拾取 task 后）：

```bash
① vtask bind
   ↓ 内部：查 module 配置 → 构造 semagroup（如 group_vtask_size:vtask-head）
   → semagroup.Decrement 原子减信号量 → 提取 hostname/组号 → 返回 resource
   ↓ HOST-BOUND 返回 "n0-0"（主机名），GROUP-BOUND 返回 "0"（组号）

② 构造 _vtask_size_sema（如 host_vtask_size:vtask-head:n0-0）

③ vtask add-subtask --direct --module=vtask-head
      --header _vtask_size_sema=$sema_name
      --header to_host=$to_host $body
   ↓ 内部：doAddLocalTaskList → INSERT t_task（vtask 列暂未设置）
   ↓ vtask-head slot 拾取 → controld 建立 vtask = id 自引用

④ 失败处理：
   - vtask unbind --sema-name="${sema_name}"  ← 回滚资源
   - semaphore increment vtask_size:${PLAT_MODULE_NAME}  ← 放行下一个
   - exit 1
```

**关键细节**：
- 步骤 ① 和 ③ 之间无事务保证。如果进程在两步之间被 kill，资源泄漏。这是已知的设计权衡（`vtask_size=1` 串行化 + 进程崩溃概率极低）
- run.sh 不写 sink-tasks.txt：vtask-head 任务已由 `--direct` 直投，task 在 wait-queue slot 直接终结（不回流 main-router）
- 所有命令均不传 `--app-id`：`APP_ID` 环境变量在 agent 容器内已设定，`param` 包自动从环境变量解析
- `_vtask_size_sema` 在 bind 返回后**由 run.sh 构造**，作为 header 传给 vtask-head。后续 controld 的 `updateVTask` 对已预扣 task 跳过 decrement，`doVTaskFinished` 按此 header 定位加回

### vtask-head（vtask 根标记 + 业务路由 + gate 释放）

vtask-head 是整个 vtask 树的根节点，承担最多的职责。

**controld 侧（自动）**：

当 slot 拾取 task 时，`GetTaskList` → `updateVTask` → `addLocalTaskList` 自动完成：

1. **流控信号量 decrement**（`updateVTask`）：
   - 带 `_vtask_size_sema` header 的 task 已由 `vtask bind` 原子预扣，领取时跳过扣减
   - 无 header 的裸任务（DEFAULT 模式）按 slot 计算值兜底扣减，保持流控计数

2. **vtask 自引用建立**（`addLocalTaskList`）：
   ```sql
   UPDATE t_task SET vtask = $1,
       headers = jsonb_set(headers, '{_vtask_id}', to_jsonb($1::text), true)
   WHERE id = $1
   ```
   根 task 的 `vtask = id`（自引用），`ListVTasks` 查询 `WHERE vtask = id`

**脚本侧**（`from-vtask-head.sh`）：

```bash
① vtask add-subtask --module=vtask-core
      [--header to_host=$to_host | --header to_ip=$from_ip] $body
   ↓ _vtask_id、_vtask_size_sema 由 agent 自动传播（异步路径，无需显式传）
   ↓ HOST-BOUND：传 to_ip，controld 自动转 to_host（同节点路由）
   ↓ GROUP-BOUND：按 body 哈希计算 to_host（业务逻辑），显式传 --header to_host

② 成功：semaphore increment vtask_size:${WAIT_QUEUE_MODULE:-wait-queue}
   ↓ 释放 gate，允许下一个 vtask 进入 wait-queue
   ↓ 此时当前 vtask 的 core 模块可能还在执行——这就是流水线化
```

**gate 释放时机分析**：gate 在 add-subtask 到 vtask-core **成功**后才释放，不是在 bind 后就释放。如果 add-subtask 失败，gate 不释放 → wait-queue 永久阻塞。这是有意为之：阻止坏 vtask 继续消耗资源。生产环境需外部监控检测 wait-queue 阻塞。

**HOST-BOUND 和 GROUP-BOUND 在 vtask-head 的差异**：

| | HOST-BOUND | GROUP-BOUND |
|--|-----------|-------------|
| `task_dist_mode`（平台参数） | HOST-BOUND | SLOT-BOUND |
| to_host 来源 | `from_ip`（agent 自动）→ controld 转 `to_host` | 脚本按 body 哈希计算 |
| 节点分配 | 单节点，同 node 路由 | 节点组，按业务逻辑分发到不同 node |
| _vtask_size_sema | `host_vtask_size:vtask-head:<hostname>` | `group_vtask_size:vtask-head:<组号>` |

### vtask-core（核心计算）

`from-vtask-core.sh` 是实际计算发生的地方。示例中只有一行：

```bash
scalebox task add --sink-module=vtask-tail "$1"
```

**为什么用 `task add` 而非 `vtask add-subtask`**：

- `vtask add-subtask`：agent 内默认走异步路径，`_vtask_id` 由 `addSinkTasks` 自动传播（同 `task add`）；`--direct` 或 agent 外需显式传 `--header _vtask_id=X`
- `task add` 创建的子 task，其 `_vtask_id` header 同样由 agent 的 `_` 前缀规则自动传播
- 两者在 agent 内行为统一：`_` 前缀 header 均自动传播，调用方无需关心
- `vtask add-subtask` 额外设置 `vtask` 列建立父子关系，`task add` 不设

**多级联扩展**（GROUP-BOUND 典型场景）：

```
vtask-head → vtask-core-1（数据预处理）→ vtask-core-2（并行计算）→ vtask-tail
```

每个 core 模块设 `task_dist_mode=HOST-BOUND` + `vtask_role: core`，级联通过 `task add --sink-module=vtask-core-N` 实现。节点内路由由 `pod_id` 参数自动匹配，节点间路由由显式 `--header to_host` 控制。

### vtask-tail（完结标记）

tail 模块的完成触发 controld 的 `doVTaskFinished`，这是 vtask 生命周期的终点：

**controld 侧（自动）**：

```go
func doVTaskFinished(te *TaskExecMessage) error {
    // 1. 查 tail task 的 _vtask_size_sema
    // 2. semaphore.AddValue(semaName, 1) → 流控信号量 increment
    // 3. UPDATE t_task SET headers.vtask_tail_id = taskID WHERE id = vtaskID
}
```

**脚本侧**（`from-vtask-tail.sh`）：无需任何操作。统一信号量后，increment 完全由 controld 负责：

```
vtask bind → 原子 decrement（分配时）
      ↓ vtask 执行……
tail task 完成 → controld doVTaskFinished → increment（同一信号量）
```

decrement 与 increment 通过 `_vtask_size_sema` header 操作同一个信号量，脚本不再参与资源归还。

### 管道数据流总览

```
[entry task]
  │ body, headers
  ▼
main-router (default.sh)
  │ body
  ▼
wait-queue (标准模块：check.sh + run.sh)
  │ check.sh: semagroup max ≤ 0 → agent 停领（准入闸门）
  │ run.sh: vtask bind 原子 decrement（预扣额度）+ add-subtask --direct
  │ task 在本 slot 终结，vtask-head 任务直投
  │ body + _vtask_size_sema + to_host
  ▼
vtask-head (from-vtask-head.sh)
  │ controld: vtask=id, headers._vtask_id=id
  │ controld: updateVTask（已预扣 task 跳过扣减）
  │ script: add-subtask + gate increment
  │ body + _vtask_id + _vtask_size_sema + to_host/to_ip
  ▼
vtask-core (from-vtask-core.sh)
  │ agent: _ 前缀 header 自动传播
  │ body + _vtask_id + _vtask_size_sema
  ▼
vtask-tail (from-vtask-tail.sh)
  │ controld: doVTaskFinished increment 流控信号量
  ▼
[vtask 完成]
```

整个管道中，`_vtask_id` 和 `_vtask_size_sema` 是两个贯穿全链路的关键 header。`_vtask_id` 维护 vtask 树结构，`_vtask_size_sema` 确保 decrement/increment 操作同一个信号量。

## vtask 高级命令

| 命令 | 功能 | 典型调用方 |
|------|------|-----------|
| `vtask bind` | 原子化绑定计算资源，返回 hostname（HOST-BOUND）或组号（GROUP-BOUND） | wait-queue 标准模块 run.sh |
| `vtask unbind` | 释放计算资源（increment 信号量） | 错误回滚补偿（正常完成由 controld 自动加回） |
| `vtask add-subtask` | 向 vtask 添加子任务（agent 内默认异步路径自动传播 `_` 前缀，`--direct` 强制 gRPC 直连） | wait-queue 标准模块 run.sh（`--direct`）, from-vtask-head.sh（异步） |
| `vtask get` | 查询 vtask 详情（派生状态、资源绑定、子任务完成比例） | 运维 / 调试 |
| `vtask fail` | 强制终止 vtask（释放 gate 与 `_vtask_size_sema` + 标记 status_code=1） | 运维 / 脚本错误处理 |
| `vtask list` | 列出 app 下所有 vtask | 运维 |
| `vtask list-subtasks` | 列出 vtask 的所有子任务（含 from_module） | 调试 |

### 用法示例

```bash
# bind：自动查找 vtask_role=head 的模块（APP_ID 从环境变量获取）
to_host=$(scalebox vtask bind) || exit 1
# HOST-BOUND → "n0-0"（hostname）
# GROUP-BOUND → "0"（组号，即 head slot 的 seq）

# add-subtask（agent 内异步路径）：_ 前缀 header 由 agent 自动传播
scalebox vtask add-subtask --module=vtask-core \
    --header to_host=$to_host $body

# unbind：通过信号量名释放资源（正常完成由 controld 自动加回，仅错误回滚时使用）
scalebox vtask unbind --sema-name="host_vtask_size:vtask-head:${to_host}"

# 或者通过 vtask-id（服务端自动查 _vtask_size_sema）
scalebox vtask unbind --vtask-id=42

# get：查看 vtask 详情
scalebox vtask get --vtask-id=42
# 输出示例：
# ID:        42
# App ID:    5
# Body:      000
# Status:    RUNNING (2/3 subtasks finished)
# Module:    vtask-core
# Resource:  host=n0-0
# Sema:      host_vtask_size:vtask-head:n0-0
# （GROUP-BOUND 时: Resource: group=2，Sema: group_vtask_size:vtask-head:2）
# Subtasks:  3 total, 2 finished

# fail：强制终止
scalebox vtask fail --vtask-id=42

# list-subtasks：查看子任务（含 from_module 列）
scalebox vtask list-subtasks --vtask-id=42
```

### 与 `task add` 的区别

- `vtask add-subtask`：设置 `vtask` 列为父 task ID，建立父子关系链。agent 内 `_vtask_id` 由 `addSinkTasks` 自动传播；agent 外需手动传 `--header _vtask_id=X`
- `task add`：不设置 `vtask` 列，子 task 的 `_vtask_id` header 由 agent `_` 前缀规则自动传播

vtask-head → vtask-core 用 `vtask add-subtask`（建立 vtask 树），vtask-core → vtask-tail 可用 `task add`（agent 自动传播上下文）。

## task_dist_mode 配置

| 模块 | DEFAULT | HOST-BOUND | GROUP-BOUND |
|------|---------|------------|-------------|
| main-router | 无 | 无 | 无 |
| wait-queue | 无（不经过） | 无 | 无 |
| vtask-head | 无 | HOST-BOUND | **SLOT-BOUND**（平台实现值） |
| vtask-core | 无 | HOST-BOUND | HOST-BOUND |
| vtask-tail | 无 | 无 | 无 |

GROUP-BOUND 的 vtask-head 使用 SLOT-BOUND，head slot 作为节点组的锚点。用户层面呈现为 GROUP-BOUND，不暴露 SLOT-BOUND 概念。

## 信号量命名规范

| 类型 | 格式 | 示例 |
|------|------|------|
| DEFAULT（流控） | `vtask_size:<module>` | `vtask_size:vtask-head` |
| HOST-BOUND（流控） | `host_vtask_size:<module>:<hostname>` | `host_vtask_size:vtask-head:n0-0` |
| GROUP-BOUND（流控） | `group_vtask_size:<module>:<组号>` | `group_vtask_size:vtask-head:0` |

## 使用步骤

### 1. DEFAULT 模式

```bash
cd scalebox/examples/vtask

# 创建 app + 投递 task
cat default-tasks.txt | scalebox run -e scalebox.env

# 查看信号量
scalebox semaphore ls --leaf-only
# Name                       Value  Value0
# vtask_size:vtask-head      3      3

# 查看 vtask
scalebox vtask list --app-id=<app_id>

# 查看 vtask 详情
scalebox vtask get --vtask-id=<vtask_id>

# 查看子任务（含 from_module）
scalebox vtask list-subtasks --vtask-id=<vtask_id>
```

DEFAULT 模式无资源绑定，不经过 wait-queue，信号量 `vtask_size:vtask-head` 控制全局并发数。

### 2. HOST-BOUND 模式

```bash
cd scalebox/examples/vtask

# 创建 app
app_id=$(cat host-tasks.txt | scalebox run -e host-bound.env | grep -o '"app_id":[0-9]*' | cut -d: -f2)
echo "app_id=$app_id"

# 查看模块状态
scalebox module list --app-id=$app_id

# 向 wait-queue 投递 task，触发完整管道：
# wait-queue（check 准入）→ bind → vtask-head → vtask-core → vtask-tail（controld 自动归还）
echo "test-body" | scalebox task add --app-id=$app_id --sink-module=wait-queue

# 等待处理后查看信号量
sleep 5
scalebox semaphore ls --leaf-only
# Name                                    Value  Value0
# host_vtask_size:vtask-head:n0-0         2      2

# 查看 vtask
scalebox vtask list --app-id=$app_id
```

### 3. GROUP-BOUND 模式

```bash
cd scalebox/examples/vtask

# 创建 app（vtask-head 平台参数用 SLOT-BOUND，h0:2 表示 2 个 head slot = 2 个节点组）
app_id=$(cat group-tasks.txt | scalebox run -e group-bound.env | grep -o '"app_id":[0-9]*' | cut -d: -f2)
echo "app_id=$app_id"

# 查看信号量（2 个节点组）
scalebox semaphore ls --leaf-only
# Name                                    Value  Value0
# group_vtask_size:vtask-head:0           2      2
# group_vtask_size:vtask-head:1           2      2

# 按组投递 task
for i in $(seq 0 3); do
  echo "0${i}0" | scalebox task add --app-id=$app_id --header to_slot_index=$i
  echo "0${i}1" | scalebox task add --app-id=$app_id --header to_slot_index=$((i+1))
done
```

`to_slot_index` 路由到对应 head slot（节点组锚点），controld 自动转换为 `to_slot`。每个节点组维护独立的流控信号量（`group_vtask_size:vtask-head:<组号>`）。

## 常见问题

| 现象 | 原因 | 解决 |
|------|------|------|
| `app_id=` 为空 | `grep` 与输出格式不匹配 | 先看原始输出：`cat host-tasks.txt \| scalebox run -e host-bound.env` |
| `vtask list` 返回空 | task 尚未被 slot 拾取，vtask 自引用未建立 | `sleep 5` 等待 slot 处理周期 |
| wait-queue 不领取任务 | check.sh 检测到 semagroup max ≤ 0（额度耗尽，属正常流控） | 等待 vtask 完成释放额度；若长时间停领，查下游是否有卡死 |
| bind 返回空 | 信号量已耗尽（check 闸门失效时才会发生） | `scalebox semaphore ls --leaf-only` 查看当前值 |
| 信号量值不收敛 | decrement 和 increment 使用了不同的信号量名 | 检查 `_vtask_size_sema` header 传递是否正确 |
| `vtask fail` 后信号量未恢复 | fail handler 执行顺序问题 | 手动 `scalebox semaphore increment` 恢复 |

## 脚本清单

| 脚本 | 位置 | 职责 |
|------|------|------|
| `check.sh` / `run.sh` | 标准模块镜像 `scalebox.net/platform/wait-queue` | 准入闸门（ACTION_CHECK）+ bind 资源 + add-subtask `--direct`（失败时 unbind 回滚），作者无需编写 |
| `run.sh` | main-router/ | 总入口，按 `from_module` 分发到对应脚本（wait-queue 回流空操作，未知 from_module 丢弃） |
| `default.sh` | main-router/ | 入口路由（仅 from_module 为空的入口 task）：DEFAULT→vtask-head，其他→wait-queue |
| `from-vtask-head.sh` | main-router/ | add-subtask（异步路径，`_` 前缀自动传播）+ gate 释放 |
| `from-vtask-core.sh` | main-router/ | 实际计算（示例中仅转发到 vtask-tail） |
| `from-vtask-tail.sh` | main-router/ | 空操作（资源由 controld `doVTaskFinished` 自动归还） |

相对于旧版的主要变更：wait-queue 封装为标准模块（`check.sh` 准入闸门 + `run.sh` bind 逻辑进镜像，task 在 wait-queue slot 终结不回流）、统一信号量（`semagroup max/decrement` → `vtask bind`、`semaphore increment` → `vtask unbind` 仅错误回滚）、tail 不再 unbind（controld `doVTaskFinished` 自动归还）、取消 `vtask_size_sema_copy` 副本信号量、GROUP-BOUND 术语统一（`group_vtask_size` 信号量前缀，平台层 `task_dist_mode=SLOT-BOUND` 保持不变）。
