# Scalebox 平台不变式（第一批）

> 面向 AI agent 与平台开发者。开发或修改 scalebox 应用前先通读 TL;DR；命中相关条目时读正文。
> 来源：`app-fxcorr` 38 条 + `app-mwa` 29 条实战实证（完整清单见 `docs/feature-inventory.md` 附录 A/C），按 `_template-invariant.md` 模板蒸馏。
> 置信度：**已交叉验证** = 两个应用或应用+官方文档同时出现；**单项目来源** = 仅一个应用/一次实验实证。

## TL;DR（20 条）

1. 建 app 即快照：改 env/镜像/slots/code_path/run.sh 后必须重建 app。
2. 未定义的 `${VAR}` 会保留字面量：每个变量显式定义（空值也写）。
3. 非 tty stdin 会被当成初始任务：脚本包装必须 `< /dev/null`。
4. 回流只靠 `sink-tasks.txt`：成功才写、失败不写、空文件是错误、嵌套脚本先 `rm -f`。
5. `task add` 不带 `--sink-module` 投给 `main_router`；没有自研 router 时会被无差别转发。
6. 任务身份 = `(mod_id, from_mod, to_host, vtask, body)`；重复投递退出码 112。
7. `OVERWRITE` 有条件：router 重建下游时禁用；人工重投时可用（官方边界待确认）。
8. 计数器初值必须精确：多一卡死、少一提前触发；归零只能发生一次。
9. wait-queue 标准模块已完成 bind：路由端 `fromWaitQueue` 必须空实现。
10. vtask-head 必须对 `vtask_size:wait-queue` 加 1：漏做后续 vtask 进不来。
11. vtask-tail 不二次归还：`doVTaskFinished` 自动归还资源。
12. 两套信号量 API 不可混用：`vtask.*` 与 `semaphore.*` 各自配套。
13. `_` 前缀 header 自动传播，但只对 `vtask_role≠""` 的 sink 模块放行；wait-queue 一跳需显式转发。
14. 任务级 `free_space_gb` + 组原子性 = 死锁：空间闸门挂批入口。
15. `slots` 条目是主机正则，匹配 hostname 第一段（短名，如 `h0`）；不是逗号列表，也不要写 FQDN；写错则 slot 全不建、链卡死。
16. 节点池必须 ON 且 ssh 可达：Slurm 重命名或僵尸 ON 会静默卡死。
17. 容器可见路径只有 `/tmp`、`/dev/shm`、`/cluster_data_root`、`/local_data_root`、`/app/bin`。
18. 状态码：-1 等待（含被闸门拦）、-3 执行中、0 完成、>0 失败。
19. 计时用 `t_task_exec.t0–t3`；`t_task` 的时间戳与容器时钟不一致。
20. task header 优先于同名环境变量；测试可用 `task add -h` 覆盖。

---

## 1. 应用与配置

### INV-01 建 app 即快照

- **现象**：app 创建后修改 env、镜像名、slots、`code_path` 或 `run.sh` 权限，再跑仍是旧值。
- **规则**：任何上述改动后重建 app；控制机上运行的仓库副本也要同步，否则旧副本会静默生效。
- **为什么**：app 创建时把模块参数、环境变量、slot 展开结果写入 DB；slot 容器启动时从快照注入，`run.sh` 权限在容器启动时缓存。
- **反例/后果**：旧 env 静默生效；镜像名不匹配表现为 `STARTING → SSHFAILED`、日志 `Unable to find image ... locally`。
- **源码指针**：`cmd/controld/grpc/app_create.go`（CreateApp 持久化）、`cmd/agent/init.go`（环境注入）。
- **置信度**：单项目来源（app-fxcorr 多次实证）。
- **适用范围**：所有 app。

### INV-02 未定义的 `${VAR}` 保留字面量

- **现象**：env 文件里漏写一个变量，app.yaml 或模块环境里出现字面量 `${VAR}`。
- **规则**：每个被引用的变量都必须显式定义；即使值应为空，也要写成 `VAR=`。
- **为什么**：CLI 的模板展开对未定义变量保留原样，不会替换为空串（`internal/common/yaml_template.go`）。
- **反例/后果**：模块收到字面量 `${VAR}`；路径、阈值、镜像名错误但无报错，链静默跑错。
- **源码指针**：`internal/common/yaml_template.go`（`RenderYAMLTemplate`）、`docs/env-loading.md`。
- **置信度**：单项目来源 + 设计文档（app-fxcorr）。
- **适用范围**：所有使用 `--env-file` / env 文件的应用。

### INV-03 非 tty stdin 是任务列表

- **现象**：脚本里调用 `scalebox run`，脚本剩余文本被创建成任务。
- **规则**：脚本/管道包装中给 `scalebox run` 加 `< /dev/null`（需要投递初始任务时除外）。
- **为什么**：`scalebox run` 在 stdin 非 tty 时把 stdin 内容当起始任务列表。
- **反例/后果**：意外创建任务；链行为与预期不符。
- **源码指针**：`cmd/cli/run/stdin_start_tasks.go`。
- **置信度**：单项目来源（app-fxcorr）。
- **适用范围**：所有脚本化的 `scalebox run`。

## 2. 任务流转与身份

### INV-04 回流只靠 `sink-tasks.txt`

- **现象**：算法模块正常退出，但下游没有任何任务。
- **规则**：成功时把下游任务写入 `${WORK_DIR}/sink-tasks.txt`；失败时不要写；文件为空会被判为后处理错误；嵌套调用子脚本时先 `rm -f` 再写自己的。
- **为什么**：agent 在任务成功后读取该文件并创建下游任务；没有文件就没有回流。
- **反例/后果**：无回流；空文件 → `null sink-tasks.txt` 后处理错误；嵌套脚本覆盖 → 重复投递或投错。
- **源码指针**：`cmd/agent/sink_tasks.go`（`addSinkTasks`）。
- **置信度**：已交叉验证（app-fxcorr、app-mwa）。
- **适用范围**：所有算法模块。

### INV-05 入口默认路由与缺省 router

- **现象**：不带 `--sink-module` 投递任务，任务被投到意外的模块。
- **规则**：入口任务默认投给 `main_router` 指定的模块；未自研 router 时，回流任务会被缺省 router 无差别转发回算法模块。
- **为什么**：`main_router` 决定入口；缺省 router 不按 `from_module` 分派。
- **反例/后果**：单模块测试出现 112 噪音；多阶段应用任务回环。
- **源码指针**：`cmd/controld/task/add_local.go`、`cmd/agent/sink_tasks.go`。
- **置信度**：已交叉验证（app-fxcorr、app-mwa）。
- **适用范围**：所有带 router 的应用。

### INV-06 任务身份与 112

- **现象**：重复投递报 duplicate key，路由/重跑失败，退出码 112。
- **规则**：任务身份 = `(mod_id, from_mod, to_host, vtask, body)`；不要重复投递同一身份的任务；重跑前清 `t_task_exec` / `t_task` 并带 app 过滤。
- **为什么**：数据库 `uniq_task` 唯一约束。
- **反例/后果**：112；链停在中间；`OVERWRITE` 与 router 重建下游任务冲突时也会 112。
- **源码指针**：`build/database/initdb.d/25-updated-schema.sql`（`uniq_task`）、`cmd/controld/task/add_local.go`（`conflict_action`）。
- **置信度**：已交叉验证（app-fxcorr 明确 112；app-mwa 提到去重窗口）。
- **适用范围**：所有任务投递。

### INV-07 `OVERWRITE` 的适用条件（待官方确认）

- **现象**：同样使用 `CONFLICT_ACTION=OVERWRITE`，在一个人工重投场景有效，在 router 重建下游任务时却撞 112。
- **规则**：人工重投同一任务可用 `OVERWRITE` + `repeatable=yes`；**router 重建下游任务时不要用 OVERWRITE**。两者边界待平台开发者给出权威定义。
- **为什么**：`OVERWRITE` 的冲突处理与 `uniq_task` 的组合行为依赖调用场景（是否重建下游、是否已有 exec 记录）。
- **反例/后果**：router 退出码 112；任务"没跑"或重跑不生效。
- **源码指针**：`cmd/controld/task/add_local.go`（`conflict_action`）、`cmd/cli/task/task.go`。
- **置信度**：**条件性冲突**（app-fxcorr 与 app-mwa 结论不同）。
- **适用范围**：人工重投、router 重建。

### INV-08 计数器初值必须精确

- **现象**：计数器多 1 导致链永不完成；少 1 导致下游提前触发。
- **规则**：计数器初值 = 本批预期条数；每完成一条 decrement；只有返回值为 0 的那一次触发下游。
- **为什么**：`AddValue` 返回原子新值，唯一归零者触发；多/少 1 都会改变终态。
- **反例/后果**：链卡死或提前完成；重复创建同名计数器 = duplicate key。
- **源码指针**：`internal/semaphore/add.go`（`AddValue`）、`internal/semaphore/create.go`。
- **置信度**：已交叉验证（app-fxcorr、app-mwa）。
- **适用范围**：所有倒计数聚合。

## 3. VTask 与流控

### INV-09 wait-queue 不二次 bind

- **现象**：路由端收到 wait-queue 回流后再次 bind，后续 vtask 进不来。
- **规则**：标准 wait-queue 镜像的 `run.sh` 已完成准入、`vtask bind`、`add-subtask --direct`；路由端的 `fromWaitQueue` 必须空实现，不能再 bind 或写 sink-tasks。
- **为什么**：bind 是原子预扣；重复 bind 会多扣额度，且任务已在 wait-queue slot 终结、不回流。
- **反例/后果**：额度泄漏；wait-queue 停领；链卡死。
- **源码指针**：`dockerfiles/wait-queue/run.sh`、`cmd/controld/grpc/vtask_ops.go`（`BindVtaskResource`）。
- **置信度**：单项目来源（app-mwa）+ 官方示例（`examples/vtask`）。
- **适用范围**：所有经 wait-queue 的 vtask。

### INV-10 vtask-head 必须释放 gate

- **现象**：第一个 vtask 正常，后续 vtask 永远进不了 wait-queue。
- **规则**：vtask-head 在 `add-subtask` 成功后必须对 `vtask_size:wait-queue` 加 1；失败时不要释放。
- **为什么**：wait-queue 的 gate 在 bind 时被预扣，只在 head 成功投递后归还；这是流水线化的前提。
- **反例/后果**：wait-queue 停领；外部监控可见 gate 长期为 0。
- **源码指针**：`dockerfiles/wait-queue/check.sh`、`cmd/controld/task/vtask.go`。
- **置信度**：已交叉验证（app-fxcorr、app-mwa）。
- **适用范围**：所有 HOST-BOUND / GROUP-BOUND vtask。

### INV-11 vtask-tail 不二次归还

- **现象**：tail 脚本再次调用 `unbind`，资源计数被多归还。
- **规则**：tail 完成时 controld 的 `doVTaskFinished` 自动 increment 流控信号量；脚本不要再 `unbind`（`vtask unbind` 仅用于错误回滚）。
- **为什么**：归还由 `_vtask_size_sema` header 定位信号量，自动完成。
- **反例/后果**：额度虚增；在制 vtask 数超过 `vtask_size`。
- **源码指针**：`cmd/controld/task/vtask.go`（`doVTaskFinished`）。
- **置信度**：已交叉验证（app-fxcorr、app-mwa）。
- **适用范围**：所有 vtask。

### INV-12 两套信号量 API 不可混用

- **现象**：调用返回 `semaphore not found: sql: no rows in result set`。
- **规则**：`vtask.CreateSemaphore` 建的信号量用 `vtask.AddSemaphoreValue`；`semaphore.Create` 建的用 `semaphore.AddValue`。
- **为什么**：两套 API 的作用域与查找路径不同（vtask 作用域 vs 全局）。
- **反例/后果**：读写不到同一信号量；流控失真。
- **源码指针**：`internal/semaphore/add.go`、`cmd/cli/vtask/vtask.go`。
- **置信度**：单项目来源（app-fxcorr）。
- **适用范围**：所有脚本中的信号量操作。

### INV-13 `_` 前缀 header 传播规则

- **现象**：应用自定义的 `_vtask_attr_*` 在某个模块丢失，或 `_vtask_id` 被错误覆盖。
- **规则**：`_` 前缀 header 由 agent 自动向下游传播；`_vtask_id` / `_vtask_size_sema` 由平台注入，应用不要手写；应用自定义属性用 `_vtask_attr_<name>`；经 wait-queue 一跳时需显式转发（wait-queue 无 `vtask_role`，自动继承被门控拦下）。
- **为什么**：`addSinkTasks` 只对 `vtask_role≠""` 的 sink 模块自动继承 `_` 头。
- **反例/后果**：属性在管道中丢失；vtask 树断开。
- **源码指针**：`cmd/agent/sink_tasks.go`、`docs/vtask-design.md`。
- **置信度**：已交叉验证（app-fxcorr、app-mwa）+ 官方示例。
- **适用范围**：所有 vtask 管道。

### INV-14 任务级空间闸门会死锁

- **现象**：盘上数据不再增长、`fail=0`、任务一直 `-1`，链停住。
- **规则**：`free_space_gb` / `dir_quota_gb` 不要挂在任务级模块上形成"整组原子性"闸门；空间闸门应挂批入口（如 vtask-head，每节点 1 slot）；`VTASK_SIZE=1` 可消除首轮盲区。
- **为什么**：组内任务必须整组进入，任务级闸门等待空间而组内其他任务占着资源，形成循环等待（数学上无安全阈值窗口）。
- **反例/后果**：死锁；需人工释放空间或重启链。
- **源码指针**：`cmd/agent/check.go`（`CheckFreeSpace` / `CheckDirQuota`）。
- **置信度**：单项目来源（app-fxcorr 三轮实验）。
- **适用范围**：带空间流控的批处理应用。

## 4. 节点、数据与容器

### INV-15 `slots` 是主机正则（匹配短名）

- **现象**：`NODES` 写成 `host1,host2`，或把完整 hostname `h0.local` 写进 `slots`，计算侧 slot 全不建，链卡在 wait-queue/READY。
- **规则**：`slots` 条目按主机正则解析（如 `n0-[01]`），不是逗号列表；匹配对象是 **hostname 第一段（短名）**，如 `h0.local` 要写 `h0`，不要写 FQDN。需要多节点时用正则或字符类。
- **为什么**：`expandRegexHostList` 用 `split_part(hostname,'.',1)` 取短名后做正则匹配；逗号串或 FQDN 匹配不到任何 host。
- **反例/后果**：`host_vtask_size` 缺失；`t_slot` 0 行、任务停在 READY；wait-queue 卡住；controld 反复报 `GetSemagroupMax: no rows`。
- **源码指针**：`cmd/controld/grpc/app_slots.go`（`expandRegexHostList`）。
- **置信度**：源码 + 阶段 5 实测（`slots: h0.local` 不建 slot，改 `h0` 后正常）。
- **适用范围**：所有 `slots` 配置。

### INV-16 节点池必须 ON 且可达

- **现象**：`t_host` 里节点状态为 ON，但物理机已死或已被 Slurm 重新分配；链静默定格。
- **规则**：每轮启动前核对 `t_host` 中 ON 且 ssh 可达的节点；节点名变化时同步更新 `NODES` / `slots`。
- **为什么**：平台按 DB 状态建 slot，不主动探测物理可达性；旧名可能仍标 ON。
- **反例/后果**：slot 建不出；任务静默等待；`STARTING → SSHFAILED`。
- **源码指针**：`build/database/initdb.d/`（`t_host`）、`cmd/controld/slot/`。
- **置信度**：单项目来源（app-fxcorr / p419）。
- **适用范围**：所有多节点部署。

### INV-17 容器可见路径白名单

- **现象**：模块里访问宿主机路径失败，或 workdir 不在共享盘导致下游读不到。
- **规则**：容器内只有 `/tmp`、`/dev/shm`、`/cluster_data_root`、`/local_data_root`（以及 `code_path` 的 `/app/bin`）；工作目录必须放在 `data_root` 之下，平台不会自动加 volume。
- **为什么**：slot 启动命令行只挂载这些目录。
- **反例/后果**：`No such file or directory`；数据写到容器内层，任务结束即丢。
- **源码指针**：`cmd/controld/grpc/app_env_vol.go`（volume 生成）。
- **置信度**：单项目来源（app-fxcorr）。
- **适用范围**：所有模块。

## 5. 可观测

### INV-18 状态码语义

- **现象**：误把 `-1` 当失败，或把 `-3` 当完成。
- **规则**：`t_task_exec.status_code`：`-1` 等待/未开始（含被闸门拦）、`-3` 执行中、`0` 完成、`>0` 失败。
- **为什么**：controld 的任务状态机定义。
- **反例/后果**：误判链状态；把闸门等待当故障。
- **源码指针**：`internal/task/status_code.go`。
- **置信度**：已交叉验证（app-fxcorr、app-mwa）。
- **适用范围**：所有运维查询。

### INV-19 计时用 t0–t3

- **现象**：用 `t_task.create_time` / `begin_time` / `end_time` 算任务耗时，结果偏差几十秒。
- **规则**：任务耗时用 `t_task_exec.t0–t3`（与容器时钟一致）；`t_task` 的时间戳与 agent 时钟不一致。
- **为什么**：两套时间分别由容器和 agent 写入，时钟未对齐。
- **反例/后果**：性能数据失真。
- **源码指针**：`build/database/initdb.d/20-scalebox-erd.sql`（`t_task_exec`）。
- **置信度**：单项目来源（app-fxcorr）。
- **适用范围**：所有性能统计。

### INV-20 header 优先于同名环境变量

- **现象**：同一名字的 header 与环境变量取值不一致，模块用了预期外的那个。
- **规则**：`scalebox::task_header` 取值时 task header 优先；测试可用 `task add -h name=value` 覆盖，无需改镜像。
- **为什么**：bash 库的取值顺序是 header → 环境变量。
- **反例/后果**：模块行为与配置不符；排错时误看环境变量。
- **源码指针**：`build/agent/rootfs/usr/local/lib/scalebox/README.md`（`task_header`）。
- **置信度**：单项目来源（app-mwa）。
- **适用范围**：所有使用 `scalebox::task_header` 的脚本。

---

## 附录：未纳入本批的条目

以下条目来自实战，但属于单点特性或需要更多验证，暂留在对应 feature/盘点文档：

| 条目 | 原因 | 位置 |
|---|---|---|
| vtask-tail 不自动创建 | 单项目来源，需确认是否平台设计 | `docs/feature-inventory.md` 附录 A #6 |
| 人工重跑契约（抄回 `_vtask_id` / 业务头） | 依赖具体场景 | 附录 A #9 |
| `task_min_seconds` 是 slot 健康阈值 | 单项目来源，易误解 | 附录 C #16；`features/timeout` |
| HOST-BOUND 与数据面 | 单项目来源 | 附录 A #24 |
| 容器不净化宿主环境 | 单项目来源（MPI 场景） | 附录 A #32 |
| `node_progress_gap` 语义 | 源码 + 单项目来源 | `features/node-progress-gap` |
| 集群/主机/安全状态机 | 平台内部机制，非应用开发高频 | `docs/feature-inventory.md` §2.11 |

## 维护说明

- 新增条目：先查 `docs/feature-inventory.md` 是否已有；按 `_template-invariant.md` 填写；尽量找源码指针；标注置信度；同步更新 TL;DR。
- 条目变更影响应用侧约定时，在应用仓库的 `AGENTS.md` 中更新入口说明。
- 本文件是**规则清单**，不是教程；过程性内容放 `features/` 与 `recipes/`。
