# 7. 高级编程特性

> 本章讲机制；具体操作步骤见 [features/](https://github.com/kaichao/scalebox/tree/dev/features)。

- §7.1 容错机制（任务级/平台级/跨模块容错 + 分级策略）
- §7.2 准入控制（存储容量/计数/同步/自定义流控）
- §7.3 任务运行排序（sort_tag / group_regex / group_index）
- §7.4 超时设置（task_max_seconds + timeout 封装建议）
- §7.5 时间戳格式（RFC3339Nano / bash / golang 示例）
- §7.6 跨集群应用
- §7.7 Slot 启动命令表达式（`((@seq))` / `((@v:var_name))` / `((@p:param_name))` 动态求解）
- §7.8 gRPC Metadata 传参（semaphore-auto-create）
- §7.9 Slot 自动扩缩容

## 7.1 容错机制

> 容错状态由应用状态与 vtask 承载；概念见 {doc}`核心概念 <../user/2_core_concepts>` §1.4，vtask 容错见 {doc}`VTask 虚拟任务 <9_vtask>`。

模块运行错误通常分为两类：
- 逻辑错误：代码逻辑导致的结果错误；
- 非逻辑错误：外部软硬件环境异常导致的结果错误；通常在外部环境恢复后，重复多次运行，可解决；

非逻辑错发生有一定的概率，大规模数据处理中，非逻辑错的自动处理，将大大提升数据处理的效率及可靠性。

### 常见的外部运行环境异常
- 不稳定的硬件和基础软件
  - 网络设备链路的不稳定
  - 环境温度升高，导致计算节点暂时停止服务
  - 磁盘盘阵中介质不稳定，导致RAID暂时失效
  - 集群存储软件的临时性、可修复的bug
- 运行中的随机异常
  - 高负载导致核心存储响应延时增加
    - 可变延时（数秒~几分钟）
    - 暂时性不可写入
  - 高负载导致低概率的内存超分，应用内存不足
  - 系统高负载，不合适的调度策略，导致算法运行超时
- 数据异常
  - 异常输入数据
  - 前导处理软件的随机错误
  - 各类运行异常导致的数据不完整

### 非逻辑错的容错处理
- 系统异常的判断与定位
  - 读取/执行/回写的时间
  - 任务返回码
  - 系统监控的异常信息
- 异常任务的容错处理
  - 自动/人工


### 7.1.1 模块级容错

按算法对非逻辑错的分类设置退出码，在 module 的 `retry_rules` 中按退出码配置重试次数。格式、示例与测试见
[retry-test](https://github.com/kaichao/scalebox/blob/dev/features/retry-test/README.md)。

### 7.1.2 平台级故障及恢复

- slot 启动后运行 task 很快出错（任务耗时异常短），单位时间内可完成多个失败的 task；
- module 的 `task_min_seconds` 参数：**连续两个**失败任务耗时低于该阈值时，slot 被置为
  `GREEDY` 并退出，以免干扰其他 task 的正常等待；
- 再结合 `retry_rules`，即可在大量 slot 运行错的情况下停止有问题的 slot；
- 最后通过手工恢复（`scalebox slot set-status` 等）。

### 7.1.3 跨模块容错

跨模块的容错以 vtask 为单位：head/core/tail 通过 `_vtask_id` 关联，失败时由
`vtask fail` 级联终止并释放资源。重试优先在 task 级完成；无法在 task 级恢复的
（如节点退出）只能整 vtask 重跑，靠 vtask 层的幂等性保证结果一致。

容错分级：

- **task 级**：基于退出码的重试；vtask 内的 task 重试需同步 vtask 状态。尽可能在
  task 级恢复。
- **slot 级**：`task_min_seconds`（GREEDY）与 slot 相关参数剔除故障 slot。
- **vtask 级**：节点本地计算模式下保证最终幂等性，vtask 层保持强一致性；部分出错
  （如节点退出）只能在 vtask 级恢复。


## 7.2 准入控制机制

> 准入控制读写应用状态（信号量/变量）；读写者与并发纪律见 {doc}`核心概念 <../user/2_core_concepts>` §1.4 与 {doc}`主路由与应用状态 <4_main_router_status>` §2.3.4。

准入控制(Admission Control)

本地计算本地资源受限，需引入准入控制，避免本地节点的关键资源（内存容量、磁盘空间大小）超载。
对于复杂应用程序的本地计算，高级准入控制可大大简化业务逻辑，避免运行异常。

算法模块的准入控制，就是控制流式数据处理的进度，以免数据处理所需不可共享的核心计算资源（内存缓存、GPU显存、本地磁盘空间）超过本地资源的物理限额，导致计算异常中断，包括运行异常退出、程序死锁等。

流量控制的实现机制是在消息处理前，通过模块的标准流控规则、自定义流控规则检查是否符合外部条件，若不符合，则在代码中显式调用sleep()，等待一段时间后再做检查。经过多次检查未果，则退出模块运行。

流控规则通过引入等待机制，降低了局部计算资源的利用率，但从整体上保证应用程序的顺利执行。
避免内存空间的峰值需求，以增加等待时间为代价，在低内存配置的机器上运行本地计算的应用。

流控分为slot级流控、node间并行同步流控。
流控支持模块运行的同步，以免系统资源受限而导致运行异常。

精准控制数据输出，若数据空间不满足，等待其他并行任务完成后释放资源，清理出空间，待空间容量满足要求后再继续运行。

### 7.2.1 准入控制分类
#### 7.2.1.1 task addmission
- 规则型准入控制：Rule-based Admission Control
  - 存储容量准入控制（Storage Capacity Admission Control）
  - 计数型准入控制（Counter-based Admission Control）
- 自定义准入控制：Custom Admission Control

#### 7.2.1.2 slot addmission
可以启动多少个slot？

### 7.2.2  存储容量准入控制

（Storage Capacity Admission Control）
#### 7.2.2.1 目录最大占用空间流控

- 标准属性名：`dir_quota_gb`，通过限制目录的最大占用存储空间来实现流控；取值格式与挂载位置见
  [`reference/module.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/module.md) §7。

#### 7.2.2.2 目录分区最少空余空间

- 标准流控属性：`free_space_gb`，通过限制目录所在存储分区的最小空余空间实现流控；取值格式同上。

### 7.2.3 vtask 计数准入控制

标准流控属性 `vtask_size`：vtask-head 启动时计数器减一，vtask-tail（或主路由）完成后加一；计数器 <= 0 时新 vtask 不能运行。按 `task_dist_mode` 生成三种信号量（`vtask_size:<mod>` / `host_vtask_size:<mod>:<host>` / `group_vtask_size:<mod>:<组号>`）；命名、作用域与生命周期见
[`reference/semaphore-variable.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/semaphore-variable.md) §4–§5，vtask 管道中的用法见 [VTask §3.4](9_vtask.md)。

### 7.2.4 节点组并行同步准入控制

标准流控属性 `node_progress_gap`（格式 `{"节点前缀":最大差值,...}`）：平台为每个
目标节点维护 `node_progress:<module>:<host>` 信号量，每完成一个 task 加一。slot 在
领取任务前用 `semagroup` 的 `DiffMin`（当前节点进度 − 组内最慢节点进度）与配置差值
比较，超出时减少本批领取数——快节点自动等慢节点。

主路由/模块按任务处理进度更新信号量，业务模块基于计数值实现同步流控。示例见
`features/node-progress-gap/README.md`。



### 7.2.5 task级自定义准入控制

- 代码定制流控（可基于信号量）
  - ACTION_CHECK/check.sh：自定义流控逻辑


## 7.3 关键模块的任务运行排序

> 排序是并行逻辑的一部分，由主路由设置 `sort_tag` 等控制信息；概念见 {doc}`核心概念 <../user/2_core_concepts>` §1.3.2。

task运行排序是scalebox应用运行的重要基础。

对于多节点协同处理的模块来说，排序可将同一时段上所有数据在不同计算资源上并行运行，以在后续模块中做协同（分组等）。

通过设置以下参数，实现排序。

### 7.3.1  排序标签
  消息头中sort_tag，是用于消息排序的专用标签，具有最高高优先级，通常由main-router设置

### 7.3.2 任务分组号

- `group_regex`：正则表达式，从任务体中提取分组字符串（旧名 `key_group_regex` 兼容）。
- `group_index`：正则表达式对应的分组编号（旧名 `key_group_index` 兼容）。

字段定义见 [`reference/app-yaml.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md) 的 `module.parameters`。

### 7.3.3 任务处理顺序

若未设置前述排序方式，则缺省按任务生成的顺序进行处理

- 幂等性：节点本地计算存在节点失败的可能性，导致本地存储失效，无法保证task级的幂等性，在跨模块的vtask层级上实现幂等性。

## 7.4 超时设置（timeout）

由于代码的bug、数据错误等原因，可能导致算法模块进入死循环，一直占用slot运行，无法正常退出。

超时设置可以解决以上问题。通过设置 module 的 `task_max_seconds` 参数，在超出该时限后，该任务（task）会退出运行，返回错误码 124。参数语义与全局兜底（`task_global_timeout_scale` → 123）见
[`reference/module.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/module.md) §7。

在封装脚本调用算法代码时，建议也在封装脚本中加入timeout的相关设置。若无此设置，外层脚本timeout，内层脚本还会运行结束，导致不确定的结果。
若算法代码名为run.py，则封装代码中对应部分可写成：
```sh
timeout ${TASK_TIMEOUT_SECONDS}s run.py $*
```


## 7.5 时间戳格式

scalebox 标准时间格式为 `2023-10-24T18:00:00.123456+0800`；`timestamps.txt` 支持的完整格式、注释与 label 规则见
[agent-interface-files](https://github.com/kaichao/scalebox/blob/dev/features/agent-interface-files/README.md)。

生成示例：bash `date +"%Y-%m-%dT%H:%M:%S.%6N%z"`；Go `time.Now().Format("2006-01-02T15:04:05.999999")`。

## 7.6 跨集群应用

见 [跨集群](10_cross_cluster.md)。

## 7.7 Slot 启动命令表达式

支持 `(( ))` 定界符的表达式；slot 启动时按外部环境动态求解。适用于单节点多 GPU、按 slot 序号设置准入参数等场景。

### 7.7.1 变量类型

表达式中的变量名以 `@` 开头：

- 数值型：`@seq` 为 slot 序号（从 0 计数）；`@i` 是向后兼容的等价写法。
- 全局字符串：`@v:var_name`，值来自 `t_global`，跨 slot 共享。
- slot 参数：`@p:param_name`，值来自 `t_slot.parameters`（jsonb），每个 slot 独立；通过 `scalebox slot set-parameter` / `get-parameter` 操作。

类型判定：

1. 表达式含 `@seq`/`@i` → 数值求值；`@v:`/`@p:` 替换后的值必须是数字（白名单校验），否则报错并保留原表达式。
2. 只含 `@v:`/`@p:` → 纯文本替换，不做数值运算。
3. 不含变量 → 数值求值。

示例：`((@seq + @v:batch_size))`，seq=5、batch_size="100" → `105`；`((@v:name))`，name="john" → `john`。

### 7.7.2 slot 序号示例

```yaml
beam-make:
  command: ${ROCM_COMMAND}
  arguments:
    free_space_gb: '{"${LOCAL_SHMDIR}":${BEAM_MAKE_FREE_GB}}'
```

```bash
LOCAL_SHMDIR=/dev/shm/scalebox/mydata
BEAM_MAKE_FREE_GB='((@seq*5+11))'
ROCM_COMMAND='singularity exec --rocm --env ROCR_VISIBLE_DEVICES=((@seq)) {{ENVS}} {{VOLUMES}} {{IMAGE}} goagent'
```

docker 形式把启动命令换成 `docker run ... --device ...`；双 GPU 按 `@seq%2` 绑定设备的示例见
[slot-expr](https://github.com/kaichao/scalebox/blob/dev/features/slot-expr/README.md)。

### 7.7.3 global 变量示例

```yaml
my-module:
  base_image: my-image:((@v:pipeline_version))
```

```sh
scalebox global set pipeline_version 20260411
```

### 7.7.4 slot 参数示例

slot 参数存于 `t_slot.parameters`（jsonb），每个 slot 独立：

```yaml
mod-slot-expr:
  command: 'docker run -d --network=host -e MY_PARAM=((@p:my_param)) -e MY_VAR=((@v:my_global)) --rm {{ENVS}} {{VOLUMES}} {{IMAGE}}'
```

```sh
scalebox slot set-parameter my_param PARAM_VALUE --slot-id 3
scalebox slot get-parameter my_param --slot-id 3
```

未设置的参数（或 jsonb 中非字符串值）按空值处理。

## 7.8 gRPC Metadata 传参

部分标志通过 gRPC metadata 传递，避免环境变量跨进程污染。

### 7.8.1 semaphore-auto-create

| metadata key | 值 | 用途 |
|-------------|------|------|
| `semaphore-auto-create` | `yes` | 信号量不存在时自动创建 |

CLI 读取环境变量 `SEMAPHORE_AUTO_CREATE=yes` 后，通过
`metadata.AppendToOutgoingContext` 透传给 controld；controld 同时接受该环境变量与
metadata。autoCreate 的创建规则与返回值语义见
[`reference/semaphore-variable.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/semaphore-variable.md) §3。

### 7.8.2 使用方式

```bash
# 设置环境变量后，CLI 自动通过 metadata 传给 controld
SEMAPHORE_AUTO_CREATE=yes scalebox semaphore get my_sema
```

## 7.9 Slot 自动扩缩容

### 7.9.1 概述

Slot 自动扩缩容根据资源使用情况自动调整任务执行槽数量，通过监控主机资源使用率智能决策扩容或缩容，优化资源利用率和任务执行效率。

### 7.9.2 功能特性

- **自动监控**：定期检查主机资源使用情况
- **智能决策**：基于使用率阈值自动决定扩容/缩容
- **双重策略**：支持资源紧张扩容和容量不足扩容
- **安全保护**：防止过度扩缩容，内置冷却机制确保系统稳定

### 7.9.3 配置参数

在模块参数（`parameters`）中设置；未显式配置的项取默认值：

| 参数 | 默认值 | 说明 |
|---|---:|---|
| `autoscaling.enabled` | `false` | 是否启用（`true` / `yes`） |
| `autoscaling.min_slots` | 1 | 最小 slot 数 |
| `autoscaling.max_slots` | 10 | 最大 slot 数 |
| `autoscaling.scale_up_threshold` | 0.8 | 使用率高于此值扩容 |
| `autoscaling.scale_down_threshold` | 0.3 | 使用率低于此值缩容 |
| `autoscaling.scale_up_step` | 2 | 单次扩容步长 |
| `autoscaling.scale_down_step` | 1 | 单次缩容步长 |
| `autoscaling.capacity_threshold` | 0.8 | 当前 slot 数低于最大可用容量的此比例时扩容 |
| `autoscaling.scale_up_cooldown_seconds` | 300 | 两次扩容之间的冷却时间 |
| `autoscaling.scale_down_cooldown_seconds` | 600 | 两次缩容之间的冷却时间 |
| `autoscaling.check_interval_seconds` | 60 | 资源检查间隔 |

### 7.9.4 扩缩容策略

#### 扩容触发条件
1. **资源紧张扩容**：当资源使用率超过扩容阈值（默认80%）时
2. **容量不足扩容**：当当前slot数量低于最大可用容量的阈值比例时

#### 缩容触发条件
1. **资源闲置缩容**：当资源使用率低于缩容阈值（默认30%）时
2. **超量运行缩容**：当slot数量超过物理容量时强制缩容

### 7.9.5 资源指标要求

系统需要三类资源指标配合工作：
- **主机总资源**（total.*）：主机配置中的总资源定义
- **主机已使用资源**（used.*）：外部监控系统定时更新的使用数据
- **Slot资源需求**（slot_req.*）：模块配置中每个slot的资源需求

三类指标需要保持相同的key名称和单位一致性。

### 7.9.6 使用示例

```
# 正常扩容场景
初始状态：slot=3，最大可用=10，使用率=90%
结果：触发扩容2个slot（使用率>80%且当前<最大容量×80%）

# 正常缩容场景  
初始状态：slot=8，最大可用=10，使用率=20%
结果：触发缩容1个slot（使用率<30%）

# 超量运行缩容
初始状态：slot=5，最大可用=1
结果：强制缩容到物理容量内
```

## 7.10 相关文档

以下高级特性有独立详述：

| 主题 | 文档 | 说明 |
|------|------|------|
| 核心概念 | {doc}`核心概念 <../user/2_core_concepts>` | 计算逻辑、应用状态、主路由、VTask |
| 主路由实现 | {doc}`主路由与应用状态 <4_main_router_status>` | 并行逻辑、状态读写、多实例、批量处理 |
| VTask 任务组 | {doc}`9_vtask` | 完整的概念模型、模块结构、信号量机制、管道流程、CLI 命令参考 |
| 跨集群架构 | {doc}`10_cross_cluster` | 数据复制模型、地址解析、gRPC 代理转发、Task 跨集群路由 |
| 安全框架 | {doc}`11_security` | JWT + RBAC + TLS、证书管理、Automation Token 轮换 |
| WebUI 使用 | {doc}`WebUI 指南 <../user/7_webui>` | REST 网关、导航结构、构建部署 |
| VS Code 插件 | {doc}`VS Code 插件 <../user/8_vscode>` | TreeView、DAG、LSP、安装配置 |
