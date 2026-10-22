# 2. 主路由与应用状态

> 本章讲实现与优化；概念定义见 {doc}`核心概念 <../user/2_core_concepts>` §1.5。
> 读者：应用开发者、平台贡献者；普通用户可先读核心概念。

## 2.1 定位与实现形态

### 2.1.1 主路由是应用控制面

主路由模块（main-router）是应用中的控制面模块：执行并行逻辑、读写应用状态、分发任务、组织模块协作。它运行形态上是应用中的一个模块（容器/脚本），功能上是应用控制面。

- 应用参数 `main_router` 指定主路由模块名。
- 模块未声明 `sink_modules` 时，回流任务默认发给主路由；没有 `main_router` 时为空。
- 主路由通常不声明 `sink_modules`，而是按 `from_module` / header 分派。
- 主路由不执行算法逻辑，不管理容器/节点/slot；这些由平台负责。

职责：

1. **执行并行逻辑**：根据回流消息与当前状态，决定下一步生成哪些任务、发往哪个模块。
2. **读写应用状态**：集中维护进度、计数、阶段、配额等控制状态。
3. **分发任务**：把任务发往算法模块或基础模块，并携带必要的任务头。
4. **组织协作**：把多个模块串成完整计算过程，直到应用级目标完成。

概念定义、边界与常见误解见 {doc}`核心概念 <../user/2_core_concepts>` §1.5。

### 2.1.2 原则上单实例与多实例

**默认：单实例**。主路由逻辑串行，避免并发写应用状态；串行协调只处理粗粒度控制消息，数据搬运与计算在模块内并行完成，因此通常不会成为性能瓶颈。

**高负载例外：多实例**。按模块分设独立主路由实例：

- 主路由使用 `task_dist_mode: SLOT-BOUND` 部署多个 slot。
- 源模块通过 `router_index` 指定发给第 n 个主路由，缺省 0。
- 任务头也可用 `to_slot` / `to_slot_index` 指定目标 slot；`to_slot` 未给时可由 `router_index` 推导。
- 把共享同一组状态的消息集中到同一个实例；不同实例负责不同的模块/状态分区。

**并发正确性由程序员保证**：采用多实例时，应用作者必须确保同一信号量/共享变量不会被多个主路由实例并发写；平台不提供跨实例的自动互斥。若两个实例写同一状态，仍会出现并发问题。

示例（app-mwa 的 beam-form 管道）：

- `router-beam-form`：`task_dist_mode: SLOT-BOUND`，`slots: h0:3`，即 3 个主路由实例。
- `pull-unpack`：`router_index: 2`；`beam-make`：`router_index: 1`；其余模块使用默认实例（索引 0）。
- 这样把 task 数量较多的模块分流到独立实例，同时每个状态变量仍只由一个实例写。

字段定义见 [app-yaml.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md) 的 `module.parameters`。

### 2.1.3 实现语言与消息入口

主路由可用 bash、Go、Python 等语言实现。常见形态：

- 按来源模块拆分脚本，例如 `from-*.sh`；每个脚本处理一类回流消息。
- 读取 task body 与 headers，按规则分支、循环、生成后续任务。
- 生成任务可写入 `sink-tasks.txt`（异步，由 agent 后处理），或在脚本中调用 `scalebox task add` / gRPC（同步，立即感知错误）。
- `task_batch_size` 控制单批次读取的最大消息数，缺省 1；主路由可一次读取多条任务并批量处理。

接口文件格式见 [agent-interface-files](https://github.com/kaichao/scalebox/blob/dev/features/agent-interface-files/README.md)；命令见 [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)。

### 2.1.4 为什么不是工作流引擎

如果主路由只做最简单的任务分发，Scalebox 就退化为一个常见的工作流引擎。工作流引擎通常按静态 DAG 调度任务，不负责应用级状态、跨模块协作、准入控制与节点本地计算。

主路由的实现重点不是 DAG 调度器，而是：

- **状态机**：读取应用状态，按并行逻辑决定下一步。
- **路由规则**：按 task body / headers / `from_module` 分派任务。
- **批量操作**：批量生成任务、批量更新信号量，降低数据库往返。

这三点把“任务分发”提升为“计算过程的持续管理”。

## 2.2 并行逻辑的实现

### 2.2.1 app.yaml 拓扑

`app.yaml` 定义并行逻辑的静态骨架：

- `parameters.main_router`：主路由模块名。
- 模块的 `parameters.task_dist_mode`：`DEFAULT` / `HOST-BOUND` / `SLOT-BOUND`。
- 模块的 `parameters.group_regex`、`pod_id`、`vtask_role`、`router_index`、`task_batch_size` 等。
- 模块的 `slots`：实例数量与节点分布。

静态拓扑只描述“有哪些模块、怎么连接”；动态规则由主路由脚本实现。例如，一个 `dir-list → file-copy → calc → merge` 的流水线，主路由根据当前阶段和回流消息决定下一步生成哪类任务。

字段定义见 {doc}`技术规范 <../appendix/1_tech_specifications>` 与 [app-yaml.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)；应用设计方法见 {doc}`应用设计 <3_app_design>`。

### 2.2.2 主路由脚本

主路由脚本实现动态并行逻辑：

- 读取 task body 与 headers，判断当前阶段、数据分组、前后模块。
- 按规则生成后续任务，指定目标模块与任务头。
- 读写应用状态（信号量、变量、全局变量）以推进计算。
- 处理异常与重试分支；必要时调用 vtask 接口做粗粒度容错。

常用操作：

```sh
# 生成后续任务
scalebox task add --sink-module=<module> --header key=value "$body"

# 读写状态
scalebox semaphore increment my_counter
scalebox variable set my_stage compute
scalebox global set my_config value
```

批量处理是主路由效率的关键：一次读取多条任务、批量生成后续任务、批量更新信号量，减少数据库往返。实现细节见 [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md) 与 {doc}`VTask 虚拟任务 <9_vtask>`。

### 2.2.3 任务头与任务体

- **task body**：任务输入数据，通常是待处理文件名、URI 或控制消息。
- **task headers**：任务元数据与控制字段；主路由用它携带路由信息与状态。
- **`_` 前缀 header**：标记为待传递任务头，随任务在模块间自动传播。

主路由常用的控制字段：

- `to_slot` / `to_slot_index`：目标 slot（SLOT-BOUND）。
- `to_host`：目标主机（HOST-BOUND）。
- `_vtask_id` / `_vtask_size_sema`：vtask 关联与资源信号量。
- 自定义 `_` 前缀 header：跨模块透传的上下文。

**header 与 variable 的分工**：header 随任务传播，适合一次性参数；variable 存在数据库中，适合跨任务、跨模块共享的状态（如阶段、累计值）。

字段与传播规则见 [task-headers.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/task-headers.md)。

## 2.3 应用状态的实现

### 2.3.1 信号量

信号量是应用或 vtask 作用域的计数器，用于流控、准入、进度同步。常见命名如 `vtask_size:<module>`、`node_progress:<module>:<host>`。

常用命令：

```sh
scalebox semaphore create <name> <value>
scalebox semaphore get <name>
scalebox semaphore increment <name>
scalebox semaphore decrement <name>
scalebox semaphore add-value <name> <delta>
scalebox semaphore ls [prefix]
scalebox semaphore delete <name>
```

信号量组（semagroup）用于对一组信号量做 min/max/diff 等操作，适合多节点/多通道同步：

```sh
scalebox semagroup min <expr>
scalebox semagroup max <expr>
scalebox semagroup diff-min <name> <expr>
scalebox semagroup diff-max <name> <expr>
scalebox semagroup increment <expr>
scalebox semagroup decrement <expr>
```

作用域、存储与生命周期见 [semaphore-variable.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/semaphore-variable.md)；vtask 中的用法见 {doc}`VTask §3.4 <9_vtask>`。

### 2.3.2 共享变量

共享变量是应用或 vtask 作用域的字符串状态，适合保存阶段、游标、配置等跨任务/跨模块信息。

```sh
scalebox variable get <name>
scalebox variable set <name> <value>
scalebox variable delete <name>
scalebox variable ls
```

vtask 作用域：

```sh
scalebox vtask get-variable <name> --vtask-id N
scalebox vtask set-variable <name> <value> --vtask-id N
```

**header 与 variable 的分工**：header 随任务传播，适合一次性参数；variable 存在数据库中，适合跨任务、跨模块共享的状态。

### 2.3.3 全局变量

全局变量（global）用于**同一物理集群内多个应用之间**共享信息/配置。它存放在 `t_global`，不按 app 隔离。

```sh
scalebox global get <name>
scalebox global set <name> <value>
scalebox global delete <name>
scalebox global ls
```

边界：

- 不属于 vtask 作用域；vtask 使用自己的信号量/变量。
- 不跨物理集群；跨物理集群使用 `queue`（见 §2.3.5）。
- 适合跨应用共享配置、指向信息等；不适合业务数据。

### 2.3.4 读写者与批量操作

信号量/变量的读写者分三层：

- **应用层**：主路由集中写，避免应用层并发冲突。
- **平台层**：controld 也会读写信号量，例如 vtask 的 `bind` 预扣、gate 扣减、tail 完成归还；这些操作由平台层统一处理，避免应用层可能出现的并发访问问题。
- **算法模块**：通常只读，用于准入控制或分支判断。

批量操作：

- 主路由一次读取多条任务（`task_batch_size`），批量生成后续任务，批量增减信号量。
- vtask 的信号量写入使用 asyncbatch（10ms/50ms）降低逐条竞争（见 {doc}`VTask §3.7 <9_vtask>`）。
- 批量大小要平衡延迟与吞吐；不要因流控联动影响并行效率。

**多实例并发**：默认单实例；多实例时由程序员保证同一状态不被多个实例并发写，平台不提供跨实例的自动互斥（见 §2.1.2）。

### 2.3.5 跨集群通信：queue

`queue` 是跨物理集群的管理面消息队列，用于跨集群应用通信：

- 消息存放在目标集群的数据库中。
- `--cluster` 选择远端集群；本地 controld 代理转发。
- 出队顺序：`priority DESC, id ASC`（高优先级先出，同优先级 FIFO）。
- 命令：`ls`、`peek`、`pull`、`purge`、`push`；`pull` 删除消息，`peek` 不删除。

边界：

- 不是 wait-queue，也不是平台内部任务队列。
- 不是应用内状态存储；跨集群信息共享用它，应用内状态用 semaphore/variable。
- `global` 只在同一物理集群内跨应用；跨物理集群用 `queue`。

命令与连接目标见 {doc}`跨集群架构 §5.6 <10_cross_cluster>` 和 [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md) 的 `scalebox queue` 章节。

## 2.4 状态迁移

### 2.4.1 任务头状态传递

主路由通过任务头把状态传给算法模块，模块再把 `_` 前缀 header 带回主路由：

```text
主路由 → 算法模块 → 主路由
```

- 以下划线 `_` 开头的 header 标记为待传递任务头，由 agent 自动传播。
- 单次传递；复杂状态需要主路由人工维护。
- 例：`_vtask_id`、`_vtask_size_sema`、自定义 `_demo=hello`。

传播规则见 [task-headers.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/task-headers.md)。

### 2.4.2 vtask 状态保持

vtask 是对并行逻辑及其状态的粗粒度封装，把跨模块 task 串成一条计算链。它跨越输入加载、计算、结果回写的整个过程，基于计数做流控，分为全局计数、分组计数、节点计数等类型。

- `_vtask_` 前缀 header 标识 vtask 属性；vtask 属性横跨 vtask 的各个模块，自动维护。
- vtask 作用域的信号量/变量由平台维护，支持流控、资源绑定与粗粒度容错。
- 结构：wait-queue（可选）→ vtask-head → 一个或多个 vtask-core → vtask-tail。

| 机制 | 作用域 | 生命周期 | 适用 |
|---|---|---|---|
| 任务头 | 单条任务链 | 单次传递，人工维护 | 一次性参数、简单上下文 |
| vtask | 跨模块计算链 | 平台自动维护 | 流控、资源绑定、容错、粗粒度状态 |

完整设计见 {doc}`VTask 虚拟任务 <9_vtask>`。

## 2.5 性能与优化

### 2.5.1 减少消息数量

主路由是全系统集中模块；路由相关的输入任务越少，批量处理越有效：

- 使用粗粒度控制消息，避免每个文件都产生一条路由消息。
- 把同一阶段的多条数据合并成一条控制消息，由算法模块内部分发。
- 减少全系统的信号量数量，优先使用粗粒度信号量。

### 2.5.2 粗粒度信号量

信号量越细，管理数据规模越大、数据库往返越多。建议：

- 按阶段/资源组设计粗粒度计数器，而不是每个 task 一个信号量。
- 相关的计数器用 `semagroup` 统一操作。
- 只在需要精确准入/进度时引入细粒度计数。

### 2.5.3 批量处理

- 主路由一次读取多条任务（`task_batch_size`，缺省 1），批量生成后续任务。
- 批量增减信号量/变量，减少数据库操作次数。
- 批量大小要平衡延迟与吞吐；过大的批次会增加单次处理时间与内存占用。
- 避免因与流控联动影响并行效率；把流控检查放在批量处理的边界。

## 2.6 状态变量命名规范

信号量、共享变量、全局变量使用同一命名规则：

- 字符集：`[A-Za-z0-9:_-]`。
- 首字符：字母或下划线。

vtask 信号量命名模式：

- 全局：`vtask_size:${mod_name}`
- 节点级：`host_vtask_size:${mod_name}:${hostname}`
- 组级（GROUP-BOUND）：`group_vtask_size:${mod_name}:${group_seq}`

命名与作用域的完整规则见 [semaphore-variable.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/semaphore-variable.md)。

## 2.7 主路由实现最佳实践

- **应用分解**：主路由是集中模块，减少路由消息数量，有利于批量处理信号量与消息发送。
- **状态设计**：header 用于一次性参数；variable 用于跨任务状态；global 用于同集群跨应用；queue 用于跨集群通信。
- **多实例**：按模块/状态分区；同一状态只由一个实例写；并发正确性由程序员保证。
- **批量处理**：一次读取多条任务、批量生成、批量更新信号量，减少数据库往返。
- **流控设计**：信号量尽量粗粒度；准入检查放在批量边界，避免影响并行效率。
- **容错**：task 级重试用 `retry_rules`；跨模块用 vtask 级容错；模块尽量幂等。
- **可观测性**：用 `task show`、`vtask get`/`vtask list`、日志定位状态与卡点。

延伸阅读：

- {doc}`核心概念 <../user/2_core_concepts>`：计算逻辑、应用状态、主路由、VTask。
- {doc}`VTask 虚拟任务 <9_vtask>`：管道、信号量与命令。
- {doc}`节点本地计算 <5_node_local_compute>`：本地计算的准入与容错。
- {doc}`高级编程特性 <7_advanced_features>`：流控、排序、超时与扩缩容。
- {doc}`应用设计 <3_app_design>`：设计模式与最佳实践。
