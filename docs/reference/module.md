# 模块条目参考

> 唯一信源：`app.yaml` 中 `modules.<name>` 这一条**怎么跑起来**——启动命令、挂载、注入的环境变量、
> 路由与落点。字段清单与类型见 [`app-yaml.md`](./app-yaml.md)（生成物）；模块与 agent 的
> 接口文件（stdin/stdout/exit code/sink-tasks）见 [`features/agent-interface-files`](../../features/agent-interface-files/README.md)。
> 最后校验：2026-10-21，对照 go-scalebox@09cd5a3。

## 1. 条目结构

没有严格必填字段：`base_image` 省略时缺省为 `scalebox.net/platform/agent`。字段按用途分组：

| 组 | 字段 | 作用 |
|---|---|---|
| 启动 | `base_image`、`command`、`volumes`、`environments` | 容器从哪个镜像起、怎么起、挂什么、注入什么 |
| 落点 | `slots`、`parameters.task_dist_mode`、`arguments.slot_options` | 在哪些节点、以什么方式建 slot |
| 数据 | `arguments.code_path`、`volumes`、`slot_options` 的映射开关 | 容器内可见路径 |
| 路由 | `sink_modules`、`sink_vmodules`、`parameters.vtask_role` | 下游是谁、任务如何回流 |
| 调优 | `arguments.*`、`parameters.*` | 超时、轮询、心跳、闸门、分组 |

## 2. 容器启动

`command` 为空时由平台生成模板：

| 运行时 | 缺省命令 |
|---|---|
| docker | `docker run -d --network=host --rm {{ENVS}} {{VOLUMES}} {{IMAGE}}` |
| docker（`slot_options` 含 `reserved_on_exit`） | 去掉 `--rm`，退出后保留容器便于排查 |
| singularity（`.sif` 镜像） | `singularity exec {{ENVS}} {{VOLUMES}} {{IMAGE}} entrypoint.sh` |
| singularity（docker 镜像源） | `singularity exec {{ENVS}} {{VOLUMES}} docker://{{IMAGE}} entrypoint.sh` |

模板变量：`{{ENVS}}`（`-e K='V'` 序列）、`{{VOLUMES}}`（`-v 宿主:容器`）、`{{IMAGE}}`。

- 自定义 `command` 时，`MODULE_ID` 的注入依赖命令里含 `GRPC_SERVER`（平台把 `-e GRPC_SERVER`
  替换为 `-e MODULE_ID=N -e GRPC_SERVER`；singularity 替换 `--env GRPC_SERVER`）；
  写自定义命令必须保留该变量。
- `slot_options` 含 `tmpfs_workdir` 且非 singularity 时，命令追加 `--tmpfs /work`。
- `reserved_on_exit` 也可写在 `arguments.reserved_on_exit: yes`（与写在 `slot_options` 中等价）。
- 容器入口是 agent（`entrypoint.sh` → `goagent`）；agent 以 `/bin/bash -c "<ACTION_RUN> <body> <headers>"`
  执行任务，`ACTION_RUN` 等由 `environments` 或镜像预设。

## 3. 挂载规则

| 容器内路径 | 宿主路径 | 开关 |
|---|---|---|
| `/app/bin`（只读） | `arguments.code_path`（须为绝对路径） | — |
| `/cluster_data_root` | 集群 `data_root` | `slot_options` 含 `disable_data_mapping` 时不挂 |
| `/local_data_root` | `/` | `slot_options` 含 `disable_local_mapping` 时不挂 |
| `/tmp`、`/dev/shm` | 同名宿主路径 | `disable_tmp_mapping`；singularity 不挂 |
| `volumes` 的容器路径 | `volumes` 的宿主路径 | 相对路径会拼到 `data_root` 下 |

`volumes` 写法固定为两段 `宿主:容器`，例如 `shared/raw:/data` → `<data_root>/shared/raw`。
**不支持第三段挂载选项**（`:ro` 等；写成三段会被整条忽略）。

## 4. 环境变量注入

平台按以下顺序合成容器环境（后者覆盖前者）：

1. app 级变量（`scalebox.env` / `--env-file`）；
2. 平台注入：`PLAT_MODULE_NAME`、`CLUSTER`、`GRPC_SERVER`、`LOCAL_IP_INDEX`、`SINK_MODULE`、
   `MODULE_ID`（命令模板注入）、`APP_ID`/`SLOT_ID`（agent 注入）；
3. `arguments` 映射：`task_max_seconds` → `PLAT_TASK_MAX_SECONDS`、`task_min_seconds` →
   `PLAT_TASK_MIN_SECONDS`、`poll_interval_seconds` → `PLAT_POLL_INTERVAL_SECONDS`、
   `max_idle_polls` → `PLAT_MAX_IDLE_POLLS`、`heartbeat_seconds` → `PLAT_HEARTBEAT_SECONDS`、
   `max_sink_tasks` → `MAX_SINK_TASKS`、
   `output_text_size` → `PLAT_OUTPUT_TEXT_SIZE`、`text_trunc_mode` → `PLAT_TEXT_TRUNC_MODE`；
   `dir_quota_gb` → `PLAT_DIR_QUOTA_GB`、`free_space_gb` → `PLAT_FREE_SPACE_GB`（后两者以 hex 编码传入）；
   模块 `parameters` 含 `vtask_size` 或 `node_progress_gap` 时另注入 `PLAT_SRV_THROTTLE=yes`；
   app 参数 `default_idle_polls` 在模块未显式配置 `max_idle_polls` 时作为其缺省；
4. `environments` 列表（`K=V`，最后覆盖）。

**空值不注入**：合成后值为空串的变量会被删除。需要"变量存在但为空"时不能依赖它。
完整变量表见 [`env-vars.md`](./env-vars.md)。

## 5. 路由与下游

- `sink_modules`：声明的下游模块名列表（DAG 边来源）；容器内 `SINK_MODULE` 取**第一个**。
- `main_router` 模块（`app.parameters.main_router`）未声明 `sink_modules` 时，`SINK_MODULE`
  回退为 app.yaml 中第一个非 router 模块。
- 其他模块未声明 `sink_modules` 时，`SINK_MODULE` 回退为 **main_router 模块名**；
  没有 main_router 时为空串。
- 模块名支持一个 `[...]` 展开（如 `calc-[12]` → `calc-1`、`calc-2` 两个模块条目）。
  方括号内按字符逐一展开，**`[a-z]` 会展开为 `a`、`-`、`z`**，不是区间；不支持 `{}` 语法。
- 算法模块成功执行后写 `${WORK_DIR}/sink-tasks.txt` 回流；路由模块默认转发规则见
  `docs/agents/10-invariants.md` INV-05。
- `parameters.vtask_role`（`head` / `core` / `tail`）把模块接入 vtask 管道，并决定
  `_vtask_*` header 的继承门控（见 `reference/task-headers.md` §4）。

## 6. 落点与分发

- `slots`：主机正则列表，匹配 hostname 第一段（短名），只取本 cluster 内 `status='ON'` 的主机；
  写错则 slot 全不建（INV-15）。host 非 ON 时 actuator 不启动其 slot，运行中变为非 ON 后
  `GetTaskList` 也不再派发任务（INV-16）。条目语法：
  - `host_regex`：每台匹配主机建 1 个 slot（seq=0）；
  - `host_regex:N`：每台主机建 N 个 slot（seq=0…N-1）；
  - `host_regex:N:group_regex`：GROUP slots，按组名分配（组名为单个字母时直接作为组前缀）；
  - 正则不含元字符时自动按前缀匹配（追加 `.*`）。
  - `slot_options` 含 `slot_on_head` 时忽略 `slots`，只在集群 HEAD 节点建 1 个 slot（seq=0）。
- `parameters.task_dist_mode`：`DEFAULT` / `HOST-BOUND` / `SLOT-BOUND`（保存时空值归一化为
  `DEFAULT`），决定任务是否要求 `to_host` / `to_slot`。
- `slot_options`：`always_running`、`reserved_on_exit`、`tmpfs_workdir`、`disable_local_mapping`、
  `disable_data_mapping`、`disable_tmp_mapping`、`slot_on_head`（调试用）。
- `slot_max_tasks`（`parameters`，schema 未覆盖）：slot 生命周期内累计执行任务数上限，达到后
  slot 优雅退出并置 READY；0 不限；`PLAT_ALWAYS_RUNNING=yes` 忽略。

## 7. 运行时参数语义

| 参数（`arguments`） | 语义 |
|---|---|
| `task_max_seconds` | 单任务超时，缺省 0（不超时）；超时中断，任务返回 124。全局兜底：status=-3 且超过 `task_max_seconds × task_global_timeout_scale`（模块参数，缺省 2.0）的任务由 pgcron（每 15 分钟）置为 123 |
| `task_min_seconds` | 单任务最小合理耗时，缺省 0；连续低于阈值判定 slot 异常（GREEDY）并退出 |
| `poll_interval_seconds` | slot 轮询新任务的间隔，缺省 6 秒 |
| `max_idle_polls` | 空闲轮询上限，超过后 slot 退出并置 READY；缺省 100（约 10 分钟）；app 参数 `default_idle_polls` 可提供缺省。host 非 ON 时 agent 收空批次，依赖此机制回 READY（`always_running` 例外，INV-16） |
| `max_sink_tasks` | 单次 `AddTaskList` 投递的最大行数；0 = 不分批（整份 `sink-tasks.txt` 一次投递）。缺省：router 0，其他模块 100；显式分批时跨批不保证事务原子性 |
| `heartbeat_seconds` | 心跳间隔，缺省 60；非正整数禁用 |
| `output_text_size` | 执行记录里 stdout/stderr/auxout 的最大字节数，缺省 65535，最小 255 |
| `text_trunc_mode` | `TAIL` = 保留文本开头（截掉尾部）；其他或空 = 保留文本末尾（截掉开头），即缺省等价于 `HEAD` |
| `dir_quota_gb` | 目录配额闸门，格式 `'{"/dir":10}'`；挂批入口，不挂任务级（INV-14） |
| `free_space_gb` | 磁盘保留空间闸门，格式同上；同上 |

## 8. 模块类型

| 类型 | 识别 | 接口要点 |
|---|---|---|
| 算法模块 | 普通 `base_image` + `ACTION_RUN` | 从 `argv[1]` 读 body、`argv[2]` 读 headers JSON；成功写 `sink-tasks.txt` |
| 路由模块 | `app.parameters.main_router` 指定 | 按 `from_module` / header 分派，不写业务回流 |
| wait-queue | 引用标准 wait-queue 镜像；无 `vtask_role` | `ACTION_CHECK` 做准入、`run.sh` 做 bind 并直投 vtask-head；转发 `_vtask_attr_*` |
| 传输模块 | `file-copy` / `rsync-copy` 等 | 按 body 描述收发文件；见 `features/` |
