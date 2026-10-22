# Task Headers 参考

> 唯一信源：task header 的名称、写入方、传播规则、转义都以本页为准。
> 语义说明（何时设、坑）在 `features/` 与 `docs/agents/`；字段类型见 `reference/app-yaml.md`。
> 最后校验：2026-10-21，对照 go-scalebox@09cd5a3。

## 1. 定义

Task header 是任务的键值元数据，随任务在 `CLI/agent → controld → DB → agent → 模块进程` 全链路传输。

- **类型**：`map[string]string`（DB 中为 `jsonb`）。
- **值只能是字符串**：不支持对象/数组；传复合值会被 JSON 解析拒绝。
- **值内容**：任意字符串，不要求是 JSON。

### 到达模块的形式

agent 以 `/bin/bash -c` 启动模块：

```text
<ACTION_RUN> <body> <headers-json>
```

- `argv[1]`：task body。body 中的 `"` 总是转义为 `\"`；body 匹配 JSON 形态
  （`{".+":.+,...}`）时，`,`、`(`、`)` 另作转义。经 bash 词法还原后到达模块。
- `argv[2]`：紧凑 JSON（空格已去），整体用单引号包裹；空 headers 时为 `{}`。
- 同时注入环境变量：`TASK_ID`（总是）；`FROM_MODULE`、`FROM_IP`（仅当 headers 含
  非空 `from_module` 时）；以及平台环境变量（见 `reference/env-vars.md`）。

### 取值

bash 模块先 `source /usr/local/lib/scalebox/functions.sh`，再用 `scalebox::task_header <name>`：
**先查 `argv[2]` 的 JSON，再回退同名环境变量**（名字转大写后查找，INV-20）。

### 值约定（列表）

需要列表时用逗号分隔字符串，由模块自行 split：

```sh
scalebox task add --headers='{"shard_names":"shard-0001.bin,shard-0002.bin"}'
```

不要嵌 JSON 数组字符串（值含转义引号，可读性差，历史上触发过转义 bug）。元素本身含逗号时，
由调用方另选分隔符（如 `|`）。

## 2. 调用方可写的 header

| 名称 | 含义 | 备注 |
|---|---|---|
| `to_host` | 目标主机名（`t_host` 主键） | HOST-BOUND 模块用；也可由 `pod_id` 相等自动推导 |
| `to_ip` | 目标主机 IPv4 | controld 转换为 `to_host`（自动补 `.<cluster>` 后缀） |
| `to_slot` | 目标 slot 的 id | SLOT-BOUND 模块用；未给时可由 main-router 的 router_index 推导 |
| `to_slot_index` | 目标 slot 的 seq | 仅头节点有效；controld 转换为 `to_slot` |
| `repeatable` | `yes` = 允许重复分发 | 绕过调度去重窗口；重跑、重放场景用 |
| `conflict_action` | 任务冲突处理策略 | `OVERWRITE` 转换为 `repeatable=yes` 后删除本 header；其他值无效果。`task add --conflict-action`；agent 也会从 `CONFLICT_ACTION` 环境变量映射 |
| `slot_broadcast` | `yes` = 向模块的每个 slot 各投一条任务 | controld 展开为多条 `to_slot` 任务 |
| `host_broadcast` | `yes` = 向模块的每个 host 各投一条任务 | controld 展开为多条 `to_host` 任务 |
| `action` | 应用自定义动作名 | 平台只运输、不解释；语义由模块（如 tape 类）自行约定 |
| `slot_last_task` | `yes` = 本任务是本 slot 的最后一个任务 | agent 执行完当前批次后退出；`PLAT_ALWAYS_RUNNING=yes` 的 slot 忽略 |
| `_vtask_attr_<name>` | 应用自定义 vtask 属性 | 平台只运输、不解释；`<name>` 由应用定义 |
| 其他任意名字 | 应用自定义 | 不以 `_` 开头；是否**自动**传播见 §4 |

## 3. 系统注入的 header

| 名称 | 写入方 | 含义 |
|---|---|---|
| `from_ip` | CLI `task add` / agent | 生成该任务的主机 IPv4 |
| `from_module` | agent | 生成该任务的模块名 |
| `last_from_module` | agent（router 中转时） | 主路由之前的真实上游模块 |
| `last_from_ip` | agent（router 中转时） | 主路由之前的真实上游 IP |
| `from_cluster` | CLI（跨集群投递时） | 上一集群名 |
| `task_id` | controld | 本任务自身的 id；仅当模块参数 `task_id_in_headers: yes`（缺省 `no`）时写入 |
| `_vtask_id` | controld | 任务所属 vtask id（= vtask-head 任务的 id）；`vtask_role: head` 且非 main-router 的任务创建时写入；`0`/空表示不属于 vtask |
| `_vtask_size_sema` | controld | vtask 额度信号量名；只对 `parameters.vtask_role` 非空的模块注入 |

`from_cluster` 仅在跨集群投递（`--remote-cluster`）且本机有 `CLUSTER` 环境变量时写入。

## 4. `_` 前缀与传播

- `_` 前缀是**平台保留命名空间**：`_vtask_id`、`_vtask_size_sema` 由平台维护；
  REST 入口直接拒绝外部设置 `_` 开头的 header；应用自定义属性只能放在 `_vtask_attr_*` 下。
- **自动传播规则**（agent 投递下游任务时）：

| 场景 | 哪些 header 自动带到下游 |
|---|---|
| 普通模块（非 router） | 所有 `_` 开头 header |
| 模块是 main-router | 只有 `_vtask_*`，且逐条按下游 sink 模块的 `vtask_role` 是否非空门控 |

- **wait-queue 一跳**：wait-queue 没有 `vtask_role`，`_vtask_attr_*` 不会自动继承；
  需要由 wait-queue 标准模块显式转发（标准镜像已实现）。
- 普通（不带 `_`）header **不会**自动传播；下游需要时由模块自己写进 `sink-tasks.txt` 的第三列。

## 5. 转义链路与历史问题

```text
json.Marshal(headers) → 单引号整体包裹（' 用 '\''）→ bash 词法还原 → 模块 argv[2]
```

- 2026-08-29 修复：旧实现按「逗号→引号」顺序逐字符转义且不处理反斜杠，值含 `"` 时
  marshal 出的 `\"` 被 bash 吃掉，模块收到损坏 JSON（`invalid character 's' in string escape code`）。
  现实现单引号包裹，任意字符串值无损传输。
- 回归测试：go-scalebox `cmd/agent/escape_test.go`（round-trip 穷举）。
- body 的转义是另一条路径（`doExternalRun`）：`"` 无条件转义为 `\"`；匹配 JSON 形态时
  另转义 `,`、`(`、`)`。body 含空格/单引号时仍可能被 bash 分词拆散，避免在 body 里放这类字符。

## 6. 与旧文档的关系

本页替代 Sphinx `reference/2_parameter_reference §2.1.4` 的 task-headers 表（该节已改为指针）。旧表里的
`from_host`、`from_module_last`、`from_ip_last` 在当前源码中不存在对应写入点，未收录。
