# VTask 应用示例

> 完整设计文档见 [go-scalebox docs/vtask-design.md](../../go-scalebox/docs/vtask-design.md)。

## 概念

**VTask**（Virtual Task）是 scalebox 在细粒度 task 之上建立的应用级粗粒度计算单元。一个 vtask 是跨模块的 task 集合，内置信号量（semaphore）和共享变量（variable），由 wait-queue（可选）、vtask-head、一个或多个级联的 vtask-core、vtask-tail 组成管道，统一管理流控、资源绑定、状态追踪和全链路上下文传播。

### 三种类型

| 类型 | 资源绑定 | 流控信号量 | 适用场景 |
|------|---------|-----------|---------|
| **DEFAULT** | 无 | `vtask_size:vtask-head` | 轻量批处理、日志汇总、数据格式转换 |
| **HOST-BOUND** | 单节点 | `host_vtask_size:vtask-head:<hostname>` | 单节点独占计算、GPU 推理、超长计算 |
| **GROUP-BOUND** | 节点组 | `slot_vtask_size:vtask-head:<slot_seq>` | 多节点协作、分布式训练、大规模并行 FFT |

### 两层抽象

| 层级 | 概念 | 取值 | 关注点 |
|------|------|------|--------|
| **vtask（应用层）** | 资源绑定形态 | HOST-BOUND / GROUP-BOUND / DEFAULT | 用户视角：vtask 跑在单节点还是节点组 |
| **module（底层）** | task 路由分发方式 | HOST-BOUND / SLOT-BOUND / 无 | 框架视角：task 如何路由到具体 slot |

**映射关系**：
- HOST-BOUND vtask：所有模块的 `task_dist_mode=HOST-BOUND`
- GROUP-BOUND vtask：vtask-head 用 `task_dist_mode=SLOT-BOUND`（head slot 为节点组锚点），vtask-core 用 `HOST-BOUND`

用户只需理解三种 vtask 类型。SLOT-BOUND 是模块级实现细节，不在 vtask 层面暴露。

### 双信号量机制

每个流控信号量有两个版本：

| 版本 | 名称 | 管理者 | 值范围 |
|------|------|--------|--------|
| **流控版** | `host_vtask_size:vtask-head:n0-0` | controld 自动 | `[0, n]` |
| **可编程版** | `:host_vtask_size:vtask-head:n0-0`（冒号前缀） | 脚本操作 | `[-1, n]` |

- 流控版控制并发 vtask 数（`vtask_size` 参数），controld 在拾取 task 时自动 decrement、tail 完成时自动 increment
- 可编程版由 `vtask bind` / `vtask unbind` 操作，记录资源占用情况
- `vtask_size_sema_copy: yes` 时，创建 app 自动建立两者，初值一致

## 模块结构

示例包含以下模块：

```
main-router → default.sh（入口路由）
  │
  ├─ DEFAULT: 直连 vtask-head
  └─ HOST-BOUND / GROUP-BOUND: 进入 wait-queue
                                  │
                                  ▼
                            wait-queue（串行化 + 资源分配）
                              from-wait-queue.sh:
                                ① vtask bind（原子化绑定资源）
                                ② vtask add-subtask --direct（gRPC 直连，需同步感知错误）
                                ③ 失败时 vtask unbind 回滚
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
                            vtask-tail（完结标记 + 资源释放）
                              from-vtask-tail.sh:
                                ① vtask unbind（仅 HOST-BOUND/SLOT-BOUND）
                                ② controld: doVTaskFinished（流控信号量 increment）
```

| 模块 | 部署位置 | 关键参数 | 职责 |
|------|---------|---------|------|
| main-router | head 节点 (h0) | `main_router: main-router` | 入口路由，按 `TASK_DIST_MODE` 分发 |
| wait-queue | head 节点 (h0) | `vtask_size: 1` | 串行化准入 + `vtask bind` 分配资源 |
| vtask-head | 计算节点 / head 节点 | `vtask_role: head`, `vtask_size: 2` | vtask 根标记 + 业务路由 + gate 释放 |
| vtask-core | 计算节点 | `vtask_role: core` | 核心计算（可多级联，如 vtask-core-1 → vtask-core-2） |
| vtask-tail | head 节点 (h0) | `vtask_role: tail` | 完结标记 + `vtask unbind` 归还资源 |

## 模块详解

### main-router（入口路由）

main-router 是 app 的入口点，`run.sh` 根据 `from_module` 为空时调用 `default.sh`，按 `TASK_DIST_MODE` 做路由决策：

```
DEFAULT → 直连 vtask-head（跳过资源分配）
HOST-BOUND / GROUP-BOUND → 进入 wait-queue（等待资源分配）
```

注意：`default.sh` 是在 main-router slot 内执行的第一个脚本，此时 task 尚未绑定任何 vtask。`TASK_DIST_MODE` 环境变量从 `app.yaml` 的 `environments` 注入。

### wait-queue（串行化 + 资源分配）

wait-queue 是整个 vtask 管道的准入控制点。核心参数 `vtask_size: 1` 保证了**同一时刻只有一个 vtask 在执行 bind**，这是并发安全的前提——bind 内部（`semagroup.Decrement`）当前未加 `FOR UPDATE`，依赖单线程串行化保证正确性。

`from-wait-queue.sh` 的执行流程：

```bash
① vtask bind
   ↓ 内部：查 module 配置 → 构造 semagroup（如 :host_vtask_size:vtask-head）
   → semagroup.Decrement 原子减信号量 → 提取 hostname/slot_seq → 返回 resource
   ↓ HOST-BOUND 返回 "n0-0"，GROUP-BOUND 返回 "0"

② 构造 _vtask_size_sema（如 host_vtask_size:vtask-head:n0-0）

③ vtask add-subtask --direct --module=vtask-head
      --header _vtask_size_sema=$sema_name
      --header to_host=$to_host $body
   ↓ 内部：doAddLocalTaskList → INSERT t_task（vtask 列暂未设置）
   ↓ vtask-head slot 拾取 → controld 建立 vtask = id 自引用

④ 失败处理：
   - vtask unbind --sema-name=":${sema_name}"  ← 回滚资源
   - semaphore increment vtask_size:wait-queue  ← 放行下一个
   - exit 1
```

**关键细节**：
- 步骤 ① 和 ③ 之间无事务保证。如果进程在两步之间被 kill，资源泄漏。这是已知的设计权衡（`vtask_size=1` 串行化 + 进程崩溃概率极低）
- 所有命令均不传 `--app-id`：`APP_ID` 环境变量在 agent 容器内已设定，`param` 包自动从环境变量解析
- `_vtask_size_sema` 在 bind 返回后**由脚本构造**，作为 header 传给 vtask-head。后续 controld 的 `updateVTask` 按此值分组聚合 decrement

### vtask-head（vtask 根标记 + 业务路由 + gate 释放）

vtask-head 是整个 vtask 树的根节点，承担最多的职责。

**controld 侧（自动）**：

当 slot 拾取 task 时，`GetTaskList` → `updateVTask` → `addLocalTaskList` 自动完成：

1. **流控信号量 decrement**（`updateVTask`）：
   - 按每个 task 的 `_vtask_size_sema` header 分组聚合
   - 例：2 个 task 分别来自 `host_vtask_size:vtask-head:n0-0` 和 `host_vtask_size:vtask-head:n0-1`，则分别减 1
   - 这是经历 bug 修复后的方案——多 slot 场景下 decrement 和 increment 必须操作同一个信号量名

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

② 成功：semaphore increment vtask_size:wait-queue
   ↓ 释放 gate，允许下一个 vtask 进入 wait-queue
   ↓ 此时当前 vtask 的 core 模块可能还在执行——这就是流水线化
```

**gate 释放时机分析**：gate 在 add-subtask 到 vtask-core **成功**后才释放，不是在 bind 后就释放。如果 add-subtask 失败，gate 不释放 → wait-queue 永久阻塞。这是有意为之：阻止坏 vtask 继续消耗资源。生产环境需外部监控检测 wait-queue 阻塞。

**HOST-BOUND 和 GROUP-BOUND 在 vtask-head 的差异**：

| | HOST-BOUND | GROUP-BOUND |
|--|-----------|-------------|
| `task_dist_mode` | HOST-BOUND | SLOT-BOUND |
| to_host 来源 | `from_ip`（agent 自动）→ controld 转 `to_host` | 脚本按 body 哈希计算 |
| 节点分配 | 单节点，同 node 路由 | 节点组，按业务逻辑分发到不同 node |
| _vtask_size_sema | `host_vtask_size:vtask-head:<hostname>` | `slot_vtask_size:vtask-head:<slot_seq>` |

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

### vtask-tail（完结标记 + 资源释放）

tail 模块的完成触发 controld 的 `doVTaskFinished`，这是 vtask 生命周期的终点：

**controld 侧（自动）**：

```go
func doVTaskFinished(te *TaskExecMessage) error {
    // 1. 查 tail task 的 _vtask_size_sema
    // 2. semaphore.AddValue(semaName, 1) → 流控版信号量 increment
    // 3. UPDATE t_task SET headers.vtask_tail_id = taskID WHERE id = vtaskID
}
```

**脚本侧**（`from-vtask-tail.sh`）：

```bash
sema=$(scalebox::task_header "$2" "_vtask_size_sema")
case "$sema" in
    host_vtask_size:*|slot_vtask_size:*)
        scalebox vtask unbind --sema-name=":${sema}"
        ;;
esac
```

**为什么用 case 匹配前缀**：

DEFAULT 模式的 `_vtask_size_sema = vtask_size:vtask-head`，但 DEFAULT 没有 bind 操作。如果无条件调 `vtask unbind`，会错误 increment 可编程版信号量。case 通过前缀匹配确保**有 bind 才有 unbind**。

**流控版 vs 可编程版的 increment 时序**：

```
controld doVTaskFinished → increment 流控版（vtask_size:...）
      ↓ 同一 tail 完成事件
脚本 from-vtask-tail.sh → vtask unbind → increment 可编程版（:vtask_size:...）
```

两者独立操作，通过 `_vtask_size_sema` header 保证操作同一个基础信号量名。

### 管道数据流总览

```
[entry task]
  │ body, headers
  ▼
main-router (default.sh)
  │ body
  ▼
wait-queue (from-wait-queue.sh)
  │ body + _vtask_size_sema + to_host
  ▼
vtask-head (from-vtask-head.sh)
  │ controld: vtask=id, headers._vtask_id=id
  │ controld: updateVTask decrement 流控信号量
  │ script: add-subtask + gate increment
  │ body + _vtask_id + _vtask_size_sema + to_host/to_ip
  ▼
vtask-core (from-vtask-core.sh)
  │ agent: _ 前缀 header 自动传播
  │ body + _vtask_id + _vtask_size_sema
  ▼
vtask-tail (from-vtask-tail.sh)
  │ controld: doVTaskFinished increment 流控信号量
  │ script: vtask unbind increment 可编程信号量
  ▼
[vtask 完成]
```

整个管道中，`_vtask_id` 和 `_vtask_size_sema` 是两个贯穿全链路的关键 header。`_vtask_id` 维护 vtask 树结构，`_vtask_size_sema` 确保 decrement/increment 操作同一个信号量。

## vtask 高级命令

| 命令 | 功能 | 典型调用方 |
|------|------|-----------|
| `vtask bind` | 原子化绑定计算资源，返回 hostname（HOST-BOUND）或 slot_seq（GROUP-BOUND） | from-wait-queue.sh |
| `vtask unbind` | 释放计算资源（increment 可编程版信号量） | from-vtask-tail.sh, 错误回滚 |
| `vtask add-subtask` | 向 vtask 添加子任务（agent 内默认异步路径自动传播 `_` 前缀，`--direct` 强制 gRPC 直连） | from-wait-queue.sh（`--direct`）, from-vtask-head.sh（异步） |
| `vtask get` | 查询 vtask 详情（派生状态、资源绑定、子任务完成比例） | 运维 / 调试 |
| `vtask fail` | 强制终止 vtask（unbind + 标记 status_code=1） | 运维 / 脚本错误处理 |
| `vtask list` | 列出 app 下所有 vtask | 运维 |
| `vtask list-subtasks` | 列出 vtask 的所有子任务（含 from_module） | 调试 |

### 用法示例

```bash
# bind：自动查找 vtask_role=head 的模块（APP_ID 从环境变量获取）
to_host=$(scalebox vtask bind) || exit 1
# HOST-BOUND → "n0-0"（hostname）
# GROUP-BOUND → "0"（slot_seq）

# add-subtask（agent 内异步路径）：_ 前缀 header 由 agent 自动传播
scalebox vtask add-subtask --module=vtask-core \
    --header to_host=$to_host $body

# unbind：通过信号量名释放资源
scalebox vtask unbind --sema-name=":host_vtask_size:vtask-head:${to_host}"

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
| vtask-head | 无 | HOST-BOUND | **SLOT-BOUND** |
| vtask-core | 无 | HOST-BOUND | HOST-BOUND |
| vtask-tail | 无 | 无 | 无 |

GROUP-BOUND 的 vtask-head 使用 SLOT-BOUND，head slot 作为节点组的锚点。用户层面呈现为 GROUP-BOUND，不暴露 SLOT-BOUND 概念。

## 信号量命名规范

| 类型 | 格式 | 示例 |
|------|------|------|
| DEFAULT（流控） | `vtask_size:<module>` | `vtask_size:vtask-head` |
| HOST-BOUND（流控） | `host_vtask_size:<module>:<hostname>` | `host_vtask_size:vtask-head:n0-0` |
| GROUP-BOUND（流控） | `slot_vtask_size:<module>:<slot_seq>` | `slot_vtask_size:vtask-head:0` |
| 可编程版 | 流控版前面加 `:` | `:host_vtask_size:vtask-head:n0-0` |

## 使用步骤

### 1. DEFAULT 模式

```bash
cd scalebox/examples/vtask

# 创建 app + 投递 task
cat default-tasks.txt | scalebox run -e scalebox.env

# 查看信号量
scalebox semaphore ls --leaf-only
# Name                       Value  Value0
# :vtask_size:vtask-head     3      3
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
# wait-queue → bind → vtask-head → vtask-core → vtask-tail → unbind
echo "test-body" | scalebox task add --app-id=$app_id --sink-module=wait-queue

# 等待处理后查看信号量
sleep 5
scalebox semaphore ls --leaf-only
# Name                                    Value  Value0
# :host_vtask_size:vtask-head:n0-0        2      2
# host_vtask_size:vtask-head:n0-0         2      2

# 查看 vtask
scalebox vtask list --app-id=$app_id
```

### 3. GROUP-BOUND 模式

```bash
cd scalebox/examples/vtask

# 创建 app（vtask-head 用 SLOT-BOUND，h0:2 表示 2 个 head slot = 2 个节点组）
app_id=$(cat group-tasks.txt | scalebox run -e group-bound.env | grep -o '"app_id":[0-9]*' | cut -d: -f2)
echo "app_id=$app_id"

# 查看信号量（2 个 slot）
scalebox semaphore ls --leaf-only
# Name                                    Value  Value0
# :slot_vtask_size:vtask-head:0           2      2
# :slot_vtask_size:vtask-head:1           2      2
# slot_vtask_size:vtask-head:0            2      2
# slot_vtask_size:vtask-head:1            2      2

# 按 slot 分组投递 task
for i in $(seq 0 3); do
  echo "0${i}0" | scalebox task add --app-id=$app_id --header to_slot_index=$i
  echo "0${i}1" | scalebox task add --app-id=$app_id --header to_slot_index=$((i+1))
done
```

`to_slot_index` 路由到对应 head slot，controld 自动转换为 `to_slot`。每个 head slot 维护独立的流控信号量。

## 常见问题

| 现象 | 原因 | 解决 |
|------|------|------|
| `app_id=` 为空 | `grep` 与输出格式不匹配 | 先看原始输出：`cat host-tasks.txt \| scalebox run -e host-bound.env` |
| `vtask list` 返回空 | task 尚未被 slot 拾取，vtask 自引用未建立 | `sleep 5` 等待 slot 处理周期 |
| bind 返回空 | 信号量已耗尽 | `scalebox semaphore ls --leaf-only` 查看当前值 |
| 信号量值不收敛 | decrement 和 increment 使用了不同的信号量名 | 检查 `_vtask_size_sema` header 传递是否正确 |
| `vtask fail` 后信号量未恢复 | fail handler 执行顺序问题 | 手动 `scalebox semaphore increment` 恢复 |

## 脚本清单

全部脚本在 `main-router/` 目录下：

| 脚本 | 行数 | 职责 |
|------|------|------|
| `run.sh` | 36 | 总入口，按 `from_module` 分发到对应脚本 |
| `default.sh` | 14 | main-router 路由：DEFAULT→vtask-head，其他→wait-queue |
| `from-wait-queue.sh` | 31 | bind 资源 + add-subtask `--direct`（失败时 unbind 回滚） |
| `from-vtask-head.sh` | 39 | add-subtask（异步路径，`_` 前缀自动传播）+ gate 释放 |
| `from-vtask-core.sh` | 8 | 实际计算（示例中仅转发到 vtask-tail） |
| `from-vtask-tail.sh` | 16 | unbind 归还资源（仅 HOST-BOUND / SLOT-BOUND） |

相对于旧版的主要变更：`check.sh` 删除（`bind` 原子化替代）、`semagroup max/decrement` → `vtask bind`、`semaphore increment` → `vtask unbind`、tail 加 case 前缀匹配、`from-wait-queue.sh` 使用 `--direct` 保护回滚链、`from-vtask-head.sh` 使用异步路径享受 `_` 前缀自动传播。
