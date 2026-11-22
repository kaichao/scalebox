# 3. VTask 虚拟任务

> 本章讲 vtask 设计与操作语义；完整命令 flags 见 [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)。

**VTask**（Virtual Task）是 Scalebox 在细粒度 task 之上建立的应用级粗粒度计算单元。一个 vtask 是跨模块的 task 集合，内置信号量（semaphore）和共享变量（variable），由 wait-queue（可选）、vtask-head、一个或多个级联 vtask-core、vtask-tail 组成管道，统一管理流控、资源绑定、状态追踪和全链路上下文传播。

## 3.1 为什么需要 VTask

VTask 是对计算逻辑中跨模块协作部分（并行逻辑）及其状态的粗粒度封装。它不封装算法逻辑；算法仍在各模块内。它也不是独立调度器：调度由平台与主路由完成。全局变量是同一物理集群内跨应用配置，不属于 vtask 作用域；vtask 主要使用 vtask 作用域的信号量/变量和资源信号量。

概念定义见 {doc}`核心概念 §1.6 <../application/1_core_concepts>`；实现与优化见 {doc}`主路由与应用状态 <2_main_router_status>`。

### 3.1.1 task 模型的局限

Scalebox 的基础模型是 module → slot → task。每个 task 是一次独立的执行单元，在细粒度调度层面工作良好，但在应用层面有三个缺口：

1. **缺少跨模块状态管理**：同一 App 内多个 module 协作时，各 module 内的 task 彼此孤立，没有"这批 task 属于同一个计算任务"的概念
2. **无法回答应用级问题**："某个计算任务执行到哪了"、"一共多少子任务、完成多少"、"失败了该释放哪些资源"
3. **缺少资源隔离**：不同计算任务在同一节点上执行时无资源预留机制

### 3.1.2 VTask 的解决方案

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

## 3.2 概念模型

### 3.2.1 三种类型

| | DEFAULT | HOST-BOUND | GROUP-BOUND |
|----|---------|------------|-------------|
| 资源绑定 | 无 | 单节点 | 节点组 |
| wait-queue | 无 | 有 | 有 |
| vtask-head 模式 | 不设 task_dist_mode | `HOST-BOUND` | `SLOT-BOUND` |
| vtask-core 模式 | 不设 task_dist_mode | `HOST-BOUND` | `HOST-BOUND` |
| 适用场景 | 轻量批处理 | 单节点独占计算 | 多节点协同计算 |

三种类型的流控信号量命名见 §3.2.4 的标识体系。

**管道结构**（`vtask-core` 可级联多个阶段）：

```
DEFAULT:       task → vtask-head → vtask-core-1 → vtask-core-2 → ... → vtask-tail
HOST-BOUND:    task → wait-queue → vtask-head → vtask-core → vtask-tail
GROUP-BOUND:   task → wait-queue → vtask-head → vtask-core-1 → vtask-core-2 → ... → vtask-tail
```

### 3.2.2 能力矩阵

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

### 3.2.3 两层抽象

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

### 3.2.4 计算资源标识

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

**节点组的成员没有平台级定义**：组内包含哪些节点、组内如何分发，仍由应用业务逻辑决定（vtask-head 脚本路由 + vtask-core 以 HOST-BOUND 部署在组内节点上）。平台只管理"组"这个锚点维度。

组内节点可以用**独立组节点**（组 slot）承载：在 vtask-core 的 `slots` 中写 `<host>:<count>:<group_prefix>`（如 `y00-00:4:x43`），平台在 `host` 上建 `count` 个组 slot，每个 slot 领取 `to_host` 匹配 `group_prefix.+` 或未指定 `to_host` 的 task。组节点的排空与替换见 §3.3.3。

**资源占用与释放**：占用 = `vtask bind` 原子扣减 + `_vtask_head_module`/`_vtask_host`/`_vtask_group` 随 task 链传播（HOST-BOUND 根 task 带 `_vtask_host`；GROUP-BOUND 根 task 带 `_vtask_group`，`_vtask_host` 由 core 任务携带）；`_vtask_size_sema` 已降为兼容/内部（新任务不再写，由这三个属性重建）；释放由 controld 自动完成（错误回滚用 `vtask unbind`），见 §3.4.5 与 §3.3.6。每个主机/组的在制 vtask 数 ≤ `vtask_size`。

### 3.2.5 根 task 标识

`vtask_role=head` 模块的 task 为 vtask 根。controld 创建 head 模块 task 时自动建立自引用：

```text
UPDATE t_task SET vtask = $1,
    headers = jsonb_set(headers, '{_vtask_id}', to_jsonb($1::text), true)
WHERE id = $1
```

`ListVTasks` 查询条件为 `WHERE vtask = id`。

### 3.2.6 状态派生

vtask 状态由根 task 的 `status_code` + 子任务完成比例派生：

| 根 task status | 子任务 | 派生状态 |
|:---:|------|------|
| `-1` | 0 | READY |
| `-2` | 0 | QUEUED |
| `-3` | 0 | RUNNING |
| `-3` | >0，部分完成 | RUNNING（X/Y finished） |
| `0` | 全部完成 | FINISHED（Y tasks） |
| `1` | — | FAILED（exit=N） |

## 3.3 模块结构

### 3.3.1 DEFAULT（4 模块）

```yaml
modules:
  main-router:   # 入口路由：task add --sink-module=vtask-head $body
  vtask-head:    # vtask_role: head, vtask_size: 2
  vtask-core:    # vtask_role: core
  vtask-tail:    # vtask_role: tail（无 unbind）
```

### 3.3.2 HOST-BOUND（5 模块）

```yaml
modules:
  main-router:   # task add --sink-module=wait-queue $body
  wait-queue:    # 标准模块 scalebox.net/platform/wait-queue（vtask_size ≥ slot 数的并发 bind 准入 + check.sh 准入 + run.sh bind）
  vtask-head:    # vtask_role: head, task_dist_mode: HOST-BOUND, pod_id 可选（非 vtask 兜底）
  vtask-core:    # vtask_role: core, task_dist_mode: HOST-BOUND, 优先 _vtask_host
  vtask-tail:    # vtask_role: tail
```

**关键脚本操作**：

| 脚本 | 位置 | 操作 |
|------|------|------|
| `check.sh` / `run.sh` | wait-queue 标准模块镜像 | 准入闸门（ACTION_CHECK）+ `vtask bind` → `vtask add-subtask --direct`（失败时 unbind + gate increment 回滚）；task 在本 slot 终结，作者无需编写 |
| `from-vtask-head.sh` | main-router/ | `vtask add-subtask --module=vtask-core`（async；HOST-BOUND 不传 `to_ip`/`to_host`，由 `_vtask_host` 兜底，老路径/非 vtask 由 `pod_id` 推导；GROUP-BOUND 必须显式 `to_host`）。**不再 +1 放行 gate**——gate 已由 wait-queue 的 `run.sh` 在 bind+投递成功后放行（2026-11-02 L1） |
| `from-vtask-core.sh` | main-router/ | `task add --sink-module=vtask-tail $1` |
| `from-vtask-tail.sh` | main-router/ | 空操作（额度由 controld `doVTaskFinished` 自动归还） |

> **2026-11-10 锚点变更（2026-11-12 更新）**：wait-queue 写入 `_vtask_head_module` + `_vtask_host`（HOST-BOUND）或 `_vtask_group`（GROUP-BOUND）；HOST-BOUND head→core 不再依赖 `pod_id`；GROUP-BOUND 首跳必须显式 `to_host`，core→core 继承 `_vtask_host`，`pod_id` 只作兜底。详见 §3.3.6。

> **收尾契约**：`vtask-tail` 应在全部 core 任务完成后投递一次。平台不对重复投递做隐式忽略；
> 重复的 tail 身份会让整批投递回滚（112），并可能因首个 tail 提前执行导致 vtask 提前收尾。

### 3.3.3 GROUP-BOUND

与 HOST-BOUND 结构相同，差异：
- `vtask-head.task_dist_mode: SLOT-BOUND`（平台实现值），bind 返回组号（即 head slot 的 seq）
- 流控信号量格式 `group_vtask_size`，head slot 作为节点组锚点
- 每个 core task 同时属于**组**和**主机**：`_vtask_group` 是组号，`_vtask_host` 是该 core 任务实际运行的 host
- `from-vtask-head.sh` 按 `VTASK_MODE=GROUP-BOUND` 分支从 `_vtask_group`（老 wait-queue 回退 `_vtask_size_sema`）计算 `to_host`（业务逻辑）；`vtask add-subtask` 必须显式 `to_host`，平台不猜组内节点
- head→core 使用脚本显式 `to_host`；head 是 SLOT-BOUND，`pod_id` 推导不触发
- core→core：显式 `to_host`/`to_ip` > `_vtask_host`（上一跳 core 的实际 host）> `pod_id`（兜底）；不再要求所有 core 模块配置相同的 `pod_id`
- controld 在 core 入库时把 `_vtask_host` 更新为显式 `to_host`；最终无法确定 `to_host` 的 GROUP-BOUND core 任务拒绝入库
- 组号优先取 `_vtask_group`（老 wait-queue 回退 `_vtask_size_sema`）；配额 `group_vtask_size` 只用组号，`_vtask_host` 不参与 sema 重建

**独立组节点**：vtask-core（及其他 HOST-BOUND 模块）的 `slots` 支持组 slot 语法：

```
<host>:<count>:<group_prefix>
```

| 配置 | 含义 |
|---|---|
| `y00-00:4:x43` | 在组锚点节点 y00-00 上建 4 个组 slot，领取 `to_host` 匹配 `x43.+` 的 task |
| `z00-00:4:x` | 第三段为单个字母：前缀为 `x`，匹配所有 `x` 开头的 host |

组 slot 在组锚点节点本机执行 task；`to_host` 未指定的 task 也可被组 slot 领取。组 slot 的 `group_prefix` 只影响领取条件，不改变组的资源标识（组号仍由 head slot 的 seq 决定）。

**组节点生命周期**：

```
ON ──drain──> DRAINING ──排空完成──> OFF ──换 IP / 主机名──> ON
```

- `DRAINING` 期间 host 仍参与 `to_host` 位置解析，避免本地链下游 task 退化为空位置；
- 配置 `parameters.skip_draining_host: 'yes'` 的模块在该 host 上停止领取，agent 收到 `slot_action=stop` 后退出，slot 回到 `READY`，且 `GetSlotList` 不再启动它；
- 未配置该参数的模块继续领取，消费组节点本地已有的数据；
- 排空完成后才置 `OFF` / 换 IP；直接 `OFF` 会切断领取资格，本地链会断掉；
- host 置 `OFF` 时不再等空闲预算：`SubmitSlotEvent` 返回 `slot_action=stop`，agent 处理完当前已拉取批次后退出（预分配 slot → READY，动态 slot → OFF）；controld 同时级联复位该 host 上 `ON/STARTING/GREEDY` 的 slot，作为 agent 不可达时的兜底。`DRAINING` 不是 `OFF`，未配 `skip_draining_host` 的模块仍继续领取；
- 运维命令：`scalebox host drain <hostname>` / `scalebox host undrain <hostname>`（或 `host set-status`）。

### 3.3.4 vtask-core 级联

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
- `_` 前缀 header（`_vtask_id`、`_vtask_head_module`、`_vtask_host`/`_vtask_group`；`_vtask_size_sema` 仅老任务/自定义）沿 core 链自动传播
- 每个 core 模块独立配置 slot 数量、`task_dist_mode` 和镜像
- HOST-BOUND vtask 的 core 链优先用 `_vtask_host` 保持同节点；GROUP-BOUND 的 core→core 也优先继承 `_vtask_host`（上一跳实际 host），`pod_id` 仍是兜底

### 3.3.5 位置路由：to_host 的来源与优先级

HOST-BOUND 任务落到哪台主机，按以下优先级确定（`cmd/controld/task/add_local.go`）：

```
显式 to_host  >  显式 to_ip  >  _vtask_host（core）  >  pod_id 推导
```

- **显式 `to_host` / `to_ip` 优先**：调用方给出的路由是明确意图，`_vtask_host`/`pod_id` 不会覆盖。
- **`_vtask_host` 优先于 `pod_id`**：对 `HOST-BOUND` + `vtask_role=core` 的任务生效。HOST-BOUND vtask 用不可变 bind 锚点；GROUP-BOUND vtask 用上一跳 core 的实际 host。GROUP-BOUND 的 head→core 首跳没有 `_vtask_host`，必须显式 `to_host`（详见 §3.3.6）。
- **`pod_id` 推导条件**（全部满足才生效）：
  1. 目标模块 `task_dist_mode=HOST-BOUND`；
  2. 来源模块 `task_dist_mode=HOST-BOUND`（SLOT-BOUND / DEFAULT 来源不触发）；
  3. 两侧 `pod_id` 非空且相同；
  4. 调用方未指定 `to_host` 或 `to_ip`。
- 经 main-router 中转时，agent 写 `pre_router_from_module` / `pre_router_from_ip`（过渡期双写旧名 `last_from_module` / `last_from_ip`），controld 解包主路由后按真实前序模块判断。
- 推导出的 `to_host` 会补 `.<cluster>` 后缀（`t_slot.host` 的存储格式）。
- 解析失败（来源 IP 不在 `t_host` 中）不会写入 `to_host`，controld 打 warn；此时任务会停在 READY，需要修正来源主机注册或显式指定路由。

vtask 各模式的配置：

| 模式 | head→core | core→core | 配置要点 |
|------|-----------|-----------|---------|
| HOST-BOUND | `_vtask_host`（wait-queue 写入；老路径 `pod_id`） | `_vtask_host`/`pod_id` | head/core 同锚点；`pod_id` 仅 core→core/非 vtask |
| GROUP-BOUND | 脚本显式 `to_host`（由 `_vtask_group` 计算） | 显式 > `_vtask_host`（上一跳实际 host）> `pod_id` | head 是 SLOT-BOUND；首跳必须显式 |
| DEFAULT | 无绑定 | 无绑定 | 不需要 `pod_id` |

### 3.3.6 vtask 锚点属性（2026-11-10，2026-11-12 更新）

HOST-BOUND vtask 的语义是 core 段回到 head 锚点；GROUP-BOUND vtask 的语义是 core 段在组内按
host 分发并本地调度。2026-11-10 起平台用三个保留属性表达锚点；2026-11-12 起 GROUP-BOUND
的 `_vtask_host` 定义为该 core 任务实际运行的 host：

| 属性 | 作用域 | 含义 |
|---|---|---|
| `_vtask_head_module` | 所有 vtask | vtask-head 模块名（sema 重建用） |
| `_vtask_host` | HOST-BOUND / GROUP-BOUND | HOST-BOUND：绑定锚点完整名（如 `n0-0.inline`），生命周期内不可变。GROUP-BOUND：该 core 任务实际运行的 host；随显式 `to_host` 更新、沿 core 链传播 |
| `_vtask_group` | GROUP-BOUND | 组号（head slot 的 seq），**不是主机** |
| `_vtask_size_sema` | 内部/兼容 | 老任务或自定义 `--sema-name`；新任务不再写，由上面三个重建 |

**写入**：wait-queue `run.sh` 在 `vtask bind` 后把 `_vtask_head_module` + `_vtask_host`
（HOST-BOUND）或 `_vtask_group`（GROUP-BOUND）写入 head task；不再写 `_vtask_size_sema`。
GROUP-BOUND 不写 `_vtask_host`；它由 controld 在首个 core 任务入库时写入实际 host。
`_vtask_*` 随 agent 的 `_` 前缀规则沿 core 链传播；router 是 `SLOT-BOUND`，自身不消费锚点，
只透传，不会被拉到计算节点。

**消费**：controld 在任务入库时：

- 目标模块 `task_dist_mode=HOST-BOUND` + `vtask_role=core` + 无 `_vtask_group` + 无显式
  `to_host`/`to_ip`：`to_host = _vtask_host`；
- GROUP-BOUND core：有显式 `to_host`/`to_ip` 时优先，并把 `_vtask_host` 更新为该 host；
  无显式路由时 `to_host = _vtask_host`（上一跳 core 的实际 host）；最后用 `pod_id` 兜底；
  仍无法确定则拒绝入库。`_vtask_group` 只用于 sema 重建/观测；
- 显式 `to_host`/`to_ip` 优先；`pod_id` 保留为 HOST-BOUND core→core/非 vtask 的兜底。

**优先级**：

```
显式 to_host  >  显式 to_ip  >  _vtask_host（core）  >  pod_id 推导
```

HOST-BOUND 与 GROUP-BOUND 的 core→core 都包含 `_vtask_host` 这一项；GROUP-BOUND 的
head→core 首跳 `_vtask_host` 为空，必须显式 `to_host`，不存在"组号直接当主机"的路径。

**sema 重建**（`VTaskSizeSema`）：

- HOST-BOUND：`host_vtask_size:<_vtask_head_module>:<short(_vtask_host)>`
- GROUP-BOUND：`group_vtask_size:<_vtask_head_module>:<_vtask_group>`；`_vtask_host` 不参与
- `updateVTask` 的“已预扣”判断、`doVTaskFinished`、`FailVtask`、`vtask unbind`、
  `vtask get` 都先看 `_vtask_size_sema`（老任务/自定义），没有再按上述重建；
- 缺字段时显式告警，不静默跳过释放；tail 的完成标记不被阻断。

**约束**：

- `_vtask_host` 用完整主机名（与 `t_slot.host` 同格式）；sema 的短名由平台截取。
- GROUP-BOUND 的 `_vtask_host` 是该 core 任务的实际 host，不是 head/入口节点；首跳必须显式。
- GROUP-BOUND core 显式覆盖 `to_host` 是受支持的：controld 更新 `_vtask_host`，下游默认跟随。
- HOST-BOUND vtask 的 core 数据面任务不建议显式覆盖 `to_host`；如覆盖，下游必须显式跟随
  实际来源，否则 `/dev/shm` 等本地数据读不到。
- router 的 `SLOT-BOUND` 不因 `_vtask_host` 改变；它不会被分发到计算节点。
- `vtask add-subtask` 在 GROUP-BOUND 的 head→core 首跳必须带 `to_host`；`--direct` 时还需
  显式转发 `_vtask_id`/`_vtask_group`/`_vtask_head_module`（CLI 不自动读取当前任务 headers）。

## 3.4 信号量机制

### 3.4.1 命名规范

三种分发模式的信号量命名、作用域与生命周期见
[`reference/semaphore-variable.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/semaphore-variable.md) §4–§5；
本章不重复事实表。vtask 相关的关键点是：资源容量信号量与 gate 锁是两类不同生命周期的信号量（§3.4.5）。

### 3.4.2 统一信号量

每个流控信号量只有一套（无冒号前缀），由 `vtask bind` 预扣、controld 自动加回，脚本不再参与资源归还；
各环节的扣减/归还者见 reference §5。

> 历史版本的双信号量（流控版 + `:` 前缀可编程版）已废弃，对应的参数不再生效。

### 3.4.3 流控流程

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

### 3.4.4 等待队列门控

wait-queue 模块的 `vtask_size` 保持 1（串行化）。gate 由 wait-queue `run.sh` 在 bind+投递成功后立即 increment（2026-11-02 L1）；`from-vtask-head.sh` 不再 +1。

### 3.4.5 vtask 生命周期中的信号量变化

> 本小节面向平台贡献者/高级开发者；普通应用开发者可先读 §3.4.3 流控流程。

管道涉及两类信号量，生命周期与归还责任不同：

| 类型 | 信号量 | 作用 | 归还责任 |
|------|--------|------|---------|
| 资源容量 | `host_vtask_size:<head>:<host>` / `group_vtask_size:<head>:<组>` | 该 host/组在制 vtask 数 ≤ `vtask_size` | **controld 自动**（tail 完成 / fail）|
| gate 锁 | `vtask_size:<wait-queue>` | wait-queue 段串行化（同一时刻只有一个 vtask 在执行 bind）| **wait-queue run.sh 解锁（L1）**；失败回滚在 run.sh；`vtask fail` 兜底 |

生命周期各阶段的变化：

| 阶段 | 触发点 | 资源容量 | gate 锁 |
|------|--------|---------|---------|
| 领取 | controld 领取 wait-queue 的裸任务（无 `_vtask_size_sema`/新锚点）| — | **−1**（`updateVTask` 兜底扣）|
| 准入 | wait-queue `check.sh`（ACTION_CHECK）| semagroup max ≤ 0 → 停领 | —（不检查 gate）|
| bind | wait-queue run.sh `vtask bind` | **−1**（原子预扣，返回 host/组号）| — |
| 建链回滚 | run.sh `vtask add-subtask --direct` 失败 | **+1**（`vtask unbind`）| **+1**（run.sh increment，放行）|
| vtask-head 运行 | 主路由 `from-vtask-head` | —（已预扣，head 领取时跳过）| —（gate 已在 wait-queue run.sh 放行）|
| vtask 管道运行 | core 各阶段 | — | — |
| 完成 | controld `doVTaskFinished`（tail 完成）| **+1**（按 tail 的 `_vtask_size_sema` 或由 `_vtask_head_module`+host/group 重建后加回）| — |
| 失败 | `FailVtask` / `vtask fail` | **+1**（释放资源）| **+1**（释放门控）|

要点：

- gate 锁的实际 gate 点是 controld `getVtaskBatchSize`：无 `vtask_role` 的模块（如
  wait-queue，DEFAULT 分发）每轮领取数 = min(`vtask_size:<mod>` 当前值, TaskBatchSize)，
  值为 0 时停领——与 check.sh（查 host/group 容量）是两个独立闸门；
- gate 的**扣减**由平台做（controld 领取裸任务时），**解锁**由 wait-queue
  `run.sh` 做（`scalebox semaphore increment vtask_size:<wait-queue>`，2026-11-02 L1；
  bind + 投递成功后立即放行下一个 vtask；失败回滚也在 run.sh）；
- 应用作者无需编写 wait-queue 相关脚本（标准模块承担）；gate 放行已在 wait-queue
  `run.sh`（2026-11-02 L1）。主路由/`from-vtask-head` **不要**再 +1，否则双重归还；
  异常兜底用 `vtask fail` + 人工 `semaphore increment` 补偿。

## 3.5 管道流程

### 3.5.1 header 传播

header 的名称、写入方与 `_` 前缀传播规则见
[`reference/task-headers.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/task-headers.md) §4；
下图只标 vtask 链路上的关键节点。

```
main-router task add
  → agent addSinkTasks: _vtask_* 门控（仅 vtask_role≠"" 的 sink 模块通过）
  → wait-queue task: 不传入 _vtask_id / 锚点属性；run.sh 写 _vtask_head_module + _vtask_host/_vtask_group

wait-queue task 执行标准模块 run.sh
  → --direct gRPC 创建 vtask-head task
  → vtask_role=head → self-ref vtask=id
  → _vtask_head_module + _vtask_host（HOST-BOUND）或 _vtask_group（GROUP-BOUND）由 run.sh 写入；
    不再写 _vtask_size_sema；GROUP-BOUND 的 _vtask_host 由首个 core 入库时写入
  → task 在 wait-queue slot 终结（不写 sink-tasks.txt，无回流）

vtask-head task 执行 from-vtask-head.sh
  → sink-tasks.txt async 创建 vtask-core task
  → agent addSinkTasks: _ 前缀 header 自动传播

vtask-core → vtask-tail（同上路径）

vtask-tail 完成
  → doVTaskFinished: 读取 _vtask_size_sema 或由新属性重建 → increment
  → from-vtask-tail.sh: 同上（老路径 sema / 新属性重建）
```

### 3.5.2 生命周期

```
task 创建 → ON CONFLICT 防重
  ↓
task 调度 → FOR UPDATE SKIP LOCKED（单 slot 独占）
  ↓
task 执行 → agent → status_code: RUNNING → OK/ERROR
  ↓  _vtask_id / _vtask_head_module / _vtask_host / _vtask_group 通过 _ 前缀自动传播（_vtask_size_sema 兼容）
  ↓
tail 完成 → doVTaskFinished → sema 重建/兼容 → increment
  ↓  headers.vtask_tail_id = tail task ID
  ↓
from-vtask-tail.sh → unbind（仅 HOST-BOUND / SLOT-BOUND）
```

### 3.5.3 add-subtask 双路径

`vtask add-subtask` 在 agent 容器内运行时，默认走 sink-tasks.txt 异步路径，加 `--direct` 切换到 gRPC 同步路径：

| | 异步（sink-tasks.txt） | 同步（gRPC `--direct`） |
|--|:---:|:---:|
| 执行时机 | agent 后处理（task 执行完毕后批量发送） | 脚本内同步调用 |
| 错误感知 | 脚本无感知（写入成功 → exit 0） | 脚本立即获知（exit code ≠ 0） |
| `_` 前缀 header | agent 自动传播（`_vtask_id`、`_vtask_head_module`、`_vtask_host`、`_vtask_group`） | 调用方显式传 `--header`；GROUP-BOUND 首跳至少转发 `_vtask_id`/`_vtask_group`/`_vtask_head_module` |
| `from_module`/`from_ip` | agent 注入（更准确） | CLI 从 env 注入 |

**选择原则**：

| 脚本 | 路径 | 原因 |
|------|------|------|
| `from-vtask-head.sh` | 异步（默认） | 无资源绑定前置操作，享受 header 自动传播；GROUP-BOUND 必须显式 `to_host` |
| wait-queue 标准模块 run.sh | 同步（`--direct`） | bind 已扣除资源，须同步感知成败，失败时回滚 |
| 运维手工 CLI | 同步（gRPC） | agent 外无 sink-tasks.txt 机制 |
| GROUP-BOUND 首跳需早失败 | 可选 `--direct` | 同步感知缺 `to_host` 的拒绝；须显式转发 `_vtask_*`，失败时自行 unbind |

## 3.6 命令参考

> 完整命令、usage 与 flags 见 [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)。本节只保留 vtask 特有的操作语义。

### 3.6.1 三层操作体系

```
Layer 1: 基础资源操作（全局作用域）
  semaphore create / get / increment / decrement / ls / delete
  semagroup min / max / diff-min / diff-max / increment / decrement
  variable get / set / delete / ls
  global get / set / delete / ls

Layer 2: vtask 作用域资源操作
  vtask get-variable / set-variable
  vtask create-semaphore / get-semaphore / add-semaphore-value / delete-semaphore

Layer 3: vtask 高级操作
  vtask bind / unbind / add-subtask / get / fail / list / list-subtasks
```

### 3.6.2 vtask list / list-subtasks

`vtask list` 输出派生状态与子任务计数；`vtask list-subtasks` 列出指定 vtask 的子任务。

### 3.6.3 vtask get

输出 `ID`、`App ID`、`Status`（含完成子任务数）、`Sema`（绑定信号量名）、`Subtasks`（总数/完成数）。

### 3.6.4 vtask fail

`scalebox vtask fail --vtask-id=<id>` 强制终止 vtask：

- 原子置失败（仅 `-1/-2/-3` 可 fail；已结束的 vtask 返回 user_msg）；
- 释放 gate 信号量与资源信号量；
- 级联把未完成子任务置失败。

释放的信号量集合与顺序见
[`reference/semaphore-variable.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/semaphore-variable.md) §5。

### 3.6.5 vtask bind / unbind

- `vtask bind` 自动查找 head 模块，返回计算资源标识。
- `vtask unbind --sema-name=...`：错误回滚时传信号量名；正常完成由 controld 自动加回。
- `vtask unbind --vtask-id=...`：运维手工；服务端查 `_vtask_size_sema`，没有则由 `_vtask_head_module` + host/group 重建。

| vtask 类型 | bind 返回 |
|-----------|----------|
| HOST-BOUND | hostname（如 `n0-0`） |
| GROUP-BOUND | 组号（head slot 的 seq，如 `0`） |
| DEFAULT | 空（no-op） |

GROUP-BOUND 的返回值是组号，不是主机；head→首个 core 的落点由应用显式 `to_host` 决定，
core→core 可由 `_vtask_host` 继承。配额仍按组维护。

### 3.6.6 vtask add-subtask

```bash
# HOST-BOUND：异步路径，_ 前缀 header 自动传播；不传路由 header，
# controld 用 _vtask_host 兜底（老路径/非 vtask 可用相同 pod_id 推导）
scalebox vtask add-subtask --module=vtask-core $body

# GROUP-BOUND 首跳（head→首个 core）：必须显式 to_host（head 是 SLOT-BOUND，
# 平台无法从组号推出组内节点）；异步路径由 agent 自动传播 _vtask_*
scalebox vtask add-subtask --module=vtask-core \
    --header to_host=$to_host $body

# GROUP-BOUND 首跳需同步感知错误：加 --direct
# CLI 不会自动读取当前任务 headers，必须显式转发 _vtask_*
scalebox vtask add-subtask --direct --module=vtask-core \
    --header to_host=$to_host \
    --header _vtask_id=$(scalebox::task_header "$2" _vtask_id) \
    --header _vtask_group=$(scalebox::task_header "$2" _vtask_group) \
    --header _vtask_head_module=$(scalebox::task_header "$2" _vtask_head_module) $body || {
    scalebox vtask unbind --sema-name="group_vtask_size:${head_module}:${group}"
    exit 1
}

# GROUP-BOUND core→core：由 sink 自动产生时继承上游 _vtask_host，
# 无需重复指定 to_host/pod_id；要偏离上一跳位置时再显式传 to_host

# Agent 外运维手工：须显式 --app-id；GROUP-BOUND core 还须显式 _vtask_group/to_host
scalebox vtask add-subtask --app-id=10 --module=vtask-core \
    --header _vtask_id=42 --header _vtask_group=0 --header to_host=n0-0 $body
```

自动注入：`from_module`（`PLAT_MODULE_NAME` env）、`from_ip`（本机 IPv4）。调用方显式传入时不覆盖。

两条路径最终都经过 controld 的 `DoAddLocalTaskList` 入库：GROUP-BOUND core 若最终没有可用
`to_host`（显式、`_vtask_host`、`pod_id` 都得不到），入库失败，不会随机落点。wait-queue
标准模块投递 `vtask-head` 是 SLOT-BOUND（`to_slot_index`），不需要 `to_host`；gate 放行由
run.sh 在 bind + 投递成功后立即执行。

### 3.6.7 vtask 作用域变量/信号量

作用域资源以 `vtask-id` 为作用域：

- 变量：`vtask get-variable` / `set-variable`
- 信号量：`vtask create-semaphore` / `get-semaphore` / `add-semaphore-value` / `delete-semaphore`

## 3.7 并发安全

| 机制 | 防护场景 |
|------|---------|
| `FOR UPDATE SKIP LOCKED` | slot 不会拿同一个 task |
| `ON CONFLICT` | 同 body 的 task 不重复创建 |
| 死锁重试 ×3 | `40P01` deadlock 自愈 |
| asyncbatch（10ms/50ms） | 信号量批量写入不逐条竞争 |
| `vtask_size ≥ slot 数`（wait-queue） | 并发 bind 准入（经验值 = 全部空额；不够会超卖，见 INV-23） |

## 3.8 运维场景速查

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

## 3.9 实现索引

> 本节是平台实现索引，面向平台贡献者；普通应用开发者可跳过。

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

> 测试详情和已知限制见 go-scalebox 仓库 `docs/vtask-review.md`（英文）。
