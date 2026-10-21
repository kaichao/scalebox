# 信号量与变量参考

> 唯一信源：流控原语的事实——存储、作用域、API、命名规范、生命周期。设计动机与
> vtask 管道原理见 Sphinx `developer/9_vtask` 与 go-scalebox `docs/vtask-design.md`。
> 最后校验：2026-10-21，对照 go-scalebox@09cd5a3。

## 1. 三类存储与两个作用域

| 原语 | 表 | 作用域 | CLI |
|---|---|---|---|
| 信号量 | `t_semaphore` | app（`vtask IS NULL`）或 vtask | `semaphore` / `semagroup` / `vtask *-semaphore*` |
| 变量 | `t_variable` | app（`vtask IS NULL`）或 vtask | `variable` / `vtask get-variable|set-variable` |
| 全局变量 | `t_global` | 跨 app（无 app 列） | `global` |

**两套 API 不可混用**（INV-12）：`vtask create-semaphore` 建的信号量只能用
`vtask add-semaphore-value` / `vtask get-semaphore` 操作；`semaphore create` 建的用
`semaphore increment/decrement/add-value`。混用时表现为 `semaphore not found`。

## 2. 数据模型

- `t_semaphore(name, value, value0, app, vtask)`：`value` 当前值，`value0` 创建时的初始值；
  `vtask IS NULL` 表示 app 全局作用域。
- `t_variable(name, value, app, vtask)`：`vtask IS NULL` 表示 app 全局。
- `t_global(name, value, t, ext)`：进程无关的全局键值，不按 app 隔离。

约束：
- `t_semaphore`：`uniq_semaphore UNIQUE NULLS NOT DISTINCT (name, vtask, app)`；
- `t_variable`：`uniq_var UNIQUE NULLS NOT DISTINCT (name, vtask, app)` **与**
  `UNIQUE (name, app)` 并存——后者更严格，因此**变量名在 app 内唯一**，vtask 变量不能与
  同 app 的其他 vtask/全局变量重名；
- `t_global`：`name` 为主键。

## 3. 信号量 API

### 全局（app 作用域，`vtask IS NULL`）

| 命令 | 语义 |
|---|---|
| `semaphore create <name> <value>` | 创建；冲突行为由 **controld 进程**的环境变量 `CONFLICT_ACTION` 控制：`OVERWRITE` 覆盖、`IGNORE` 跳过、缺省报错。注意：当前 `createOn` 的 `ON CONFLICT (name, app)` 与实际唯一约束 `(name, vtask, app)` 不匹配，`OVERWRITE`/`IGNORE` 预计直接报 constraint 错（源码待修） |
| `semaphore get <name>` | 读当前值 |
| `semaphore increment <name>` / `decrement <name>` | 值 ±1 |
| `semaphore add-value <name> <delta>` | 原子加 delta，输出新值 |
| `semaphore delete <name>` | 删除 |
| `semaphore ls [prefix]` | 按前缀/正则列出 |

### semagroup（组操作）

组 = 名字匹配同一表达式（前缀或正则）的一组信号量（用于 `host_vtask_size:*`、`group_vtask_size:*`）。
表达式首字符不是正则元字符时自动加 `^`（前缀匹配）；`diff-min`/`diff-max` 使用
`(group-prefix)name` 形式，`(group-prefix)` 会展开为 `^<prefix>.*`。组内选择按
`value DESC, name` 取最大成员（`semagroup.decrement` 即如此选）。

| 命令 | 语义 |
|---|---|
| `semagroup min <expr>` / `max <expr>` | 组内最小值/最大值（返回成员名与值） |
| `semagroup diff-min <name> <expr>` | 当前值 − 组最小值 |
| `semagroup diff-max <name> <expr>` | 组最大值 − 当前值 |
| `semagroup increment/decrement <expr>` | 对组操作 |

### vtask 作用域

`vtask create-semaphore <name> <value> --vtask-id N`、`vtask get-semaphore`、
`vtask add-semaphore-value <name> <delta> --vtask-id N`、`vtask delete-semaphore`。

### AddValue 语义细节

- 原子执行 `UPDATE ... SET value = value + delta RETURNING value`，**返回新值**；
  信号量已存在时，返回 0 的那一次被视为"归零者"（INV-08 的触发点）。
- 目标不存在时：`autoCreate` 为真则按 `delta` 创建（新值 = delta）并**返回 0**；
  否则报 `semaphore not found`。autoCreate 由调用方设置 `SEMAPHORE_AUTO_CREATE=yes`
  （CLI 透传为 gRPC metadata `semaphore-auto-create: yes`；controld 进程环境变量同样生效）。
- 注意 autoCreate 的"返回 0"与归零者的"返回 0"不可区分——依赖归零判断的计数器应先 `create`。

## 4. 流控信号量命名

| 分发模式 | 格式 | 示例 |
|---|---|---|
| DEFAULT | `vtask_size:<module>` | `vtask_size:vtask-head` |
| HOST-BOUND | `host_vtask_size:<module>:<host>` | `host_vtask_size:vtask-head:n0-0` |
| SLOT-BOUND / GROUP-BOUND | `group_vtask_size:<module>:<group>` | `group_vtask_size:vtask-head:2` |

wait-queue 的串行化 gate 用 `vtask_size:wait-queue`（值通常为 1）。

## 5. 生命周期：谁预扣、谁归还

| 步骤 | 执行者 | 操作 |
|---|---|---|
| 准入 | wait-queue `check.sh` | 读 semagroup max，≤ 0 时非 0 退出，agent 停领任务（任务保持 READY） |
| 预扣 | `vtask bind`（服务端 `BindVtaskResource`） | 对组原子 decrement，返回 hostname/组号，写入 `_vtask_size_sema` |
| 领取 | controld `updateVTask` | 已预扣的 task 跳过扣减；DEFAULT 裸任务兜底扣减 |
| 归还 | controld `doVTaskFinished` | vtask-tail 完成时按 `_vtask_size_sema` increment |
| 失败归还 | `vtask fail` / `FailVtask` | 对 root task 的 `_vtask_size_sema` increment，**并**释放所有 gate 信号量（app 内 `vtask_size>0` 且无 `vtask_role` 的模块 → `vtask_size:<name>`）；随后级联把未完成子任务置失败 |

约定：
- **脚本不要重复归还**（`vtask unbind` 只用于错误回滚，INV-11）。
- **vtask-head 必须对 `vtask_size:wait-queue` 加 1**（INV-10），否则后续 vtask 进不来。
- 正常稳态下信号量值在 `[0, vtask_size]` 内；为负说明有脚本重复扣减，或 bind 绕过了
  wait-queue 的准入（`semagroup.decrement` 本身不设下限）。

## 6. 变量

| 命令 | 语义 |
|---|---|
| `variable get <name>` / `set <name> <value>` / `delete` / `ls` | app 作用域（`vtask IS NULL`），读写 `t_variable` |
| `vtask get-variable <name> --vtask-id N` / `vtask set-variable ...` | vtask 作用域 |
| `global get/set/delete/ls` | `t_global`，跨 app |

变量与 header 的分工：**header 随任务传播**，适合一次性参数；**变量存在 DB 中**，适合跨任务、
跨模块共享的状态（如进度、累计值）。模块脚本通过 `scalebox::*` 函数或 CLI 读写。

## 7. 常见错误

| 现象 | 原因 | 处理 |
|---|---|---|
| `semaphore not found: sql: no rows` | 两套 API 混用，或作用域不对（vtask vs 全局） | 按 §1 对齐创建与读取的 API |
| 重复 `create` 报 duplicate key | `t_semaphore` 唯一约束 | 用 `CONFLICT_ACTION=OVERWRITE`/`IGNORE`，或先 delete |
| 链永不完成 / 提前完成 | 计数器初值与实际条数不一致 | 初值必须精确（INV-08） |
| 在制 vtask 超限或为负 | 重复 bind/unbind、head 未释放 gate | 检查 INV-09/10/11 |
