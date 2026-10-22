<!-- AUTO-GENERATED from `cmd/webui/frontend/src/schema/app-schema.json` (JSON Schema) by `make docs-gen` (go-scalebox); DO NOT EDIT -->
<!-- source: go-scalebox@233a3c4 (2026-10-21) -->


# App YAML 参考

> 字段结构由 webui 的 `app-schema.json` 生成；参数语义（何时用、副作用）见 `features/*/README.md` 与 Sphinx 参数表。
> schema 与 controld 的 AppDef 对齐，但 `arguments` / `parameters` 允许扩展键，以下只列 schema 已知项。

## App 顶层字段

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| `name` | `string` | 是 | — | App 名称（字母/数字/下划线/连字符/点） |
| `label` | `string` | 否 | — | App 显示标签 |
| `version` | `string` | 否 | — | App 版本号 |
| `cluster` | `string` | 是 | — | 目标集群名 |
| `comment` | `string` | 否 | — | 备注说明 |
| `parameters` | `appParameters` | 否 | — |  |
| `modules` | object<string, `module`> | 是 | — | 模块定义（至少一个） |

## app.parameters

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| `initial_status` | `string` | 否 | — | App 初始状态 取值：`RUNNING`, `INITIAL`, `PAUSED` |
| `main_router` | `string` | 否 | — | 主路由模块名 |
| `default_idle_polls` | string / number | 否 | — | 所属 module 的缺省的空闲轮询次数 |
| `is_cluster_admin` | `string` | 否 | — | 集群管理员标记 |
| `default_sleep_count` | string / number | 否 | — | 各模块 sleep_count 的默认值 |
| `global_directories` | `string` | 否 | — | 全局目录列表，用于所属模块的读写数据量分类统计。（local/global/tmpfs） |
| `slot_group` | `string` | 否 | — | map 形式的多节点的多 module-slot 配置，用于新增节点的 slot 自动创建。例：'{"module0":n0,"module1",n1}' |

## modules.<name> 字段

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| `label` | `string` | 否 | — | 模块显示标签 |
| `base_image` | `string` | 是 | — | Docker 镜像地址 |
| `cluster` | `string` | 否 | — | 模块绑定集群（默认继承 App 集群） |
| `command` | `string` | 否 | — | 容器启动命令 |
| `comment` | `string` | 否 | — | 备注说明 |
| `volumes` | array<`string`> | 否 | — | 挂载卷列表 |
| `environments` | array<`string`> | 否 | — | 环境变量列表 |
| `slots` | array<`string`> | 否 | — | Slot 列表（支持正则；匹配 hostname 第一段，如 h0、n0-[01]） |
| `sink_modules` | array<`string`> | 否 | — | 下游模块集合（VTask 管道路由，DAG 边来源） |
| `sink_vmodules` | array<`string`> | 否 | — | 下游虚拟模块集合 |
| `arguments` | object（任意键值） | 否 | — | 容器启动参数（module-arguments 参数表） |
| `parameters` | object（任意键值） | 否 | — | 模块参数（module-parameters 参数表） |

### module.arguments 已知键

> arguments 允许扩展键（additionalProperties: true）；完整参数表见 Sphinx reference/2_parameter_reference。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| `grpc_server` | `string` | 否 | — | controld 的服务端点，${ip_addr}:${port}，port 缺省值为 50051 |
| `code_path` | `string` | 否 | — | 模块代码的目录，通过数据卷映射到容器内 /app/bin |
| `local_ip_index` | `string` | 否 | — | hostname -I 返回 IP 地址列表，该参数指定列表中的第 n 个 IP 作为本机 IP |
| `task_max_seconds` | string / number | 否 | — | 每个 task 运行中超时设置的秒数，超时中断返回码 124；全局超时返回 123 |
| `task_min_seconds` | string / number | 否 | — | 每个 task 运行的最小秒数，连续多次低于此值则判定 slot 为 GREEDY 异常并退出 |
| `poll_interval_seconds` | string / number | 否 | — | Slot 检查新任务的轮询间隔（秒），缺省 6 秒 |
| `max_idle_polls` | string / number | 否 | — | Slot 在退出前可进行的最大空闲轮询次数，缺省 100（10 分钟） |
| `dir_quota_gb` | `string` | 否 | — | 目录容量配额限制（GB）。格式：'{"/dir-1":10,"/dir-2":100}' |
| `free_space_gb` | `string` | 否 | — | 目录所在磁盘需要保留的最小空间（GB）。格式：'{"/dir-3":10,"/dir-4":100}' |
| `heartbeat_seconds` | string / number | 否 | — | 心跳间隔（秒），缺省 60；非正整数则禁用心跳 |
| `output_text_size` | string / number | 否 | — | task 执行记录中大文本字段（stdout/stderr/auxout）的最大字节数，缺省 65535 |
| `text_trunc_mode` | `string` | 否 | — | 文本截断模式：HEAD（头截断，保留末尾）/ TAIL，缺省 HEAD 取值：`HEAD`, `TAIL` |
| `timezone_mode` | `string` | 否 | — | 时区模式 取值：`HOST`, `UTC`, `NONE` |
| `slot_options` | `string` | 否 | — | 逗号分隔的 slot 选项：always_running（slot 一直运行，仅调试）/ reserved_on_exit（slot 退出后保留容器）/ tmpfs_workdir（tmpfs 存放工作目录 /work）/ disable_local_mapping / disable_data_mapping / slot_on_head（头节点生成单个 slot 实例） |

### module.parameters 已知键

> parameters 允许扩展键（additionalProperties: true）；完整参数表见 Sphinx reference/2_parameter_reference。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| `priority` | `string` | 否 | — | 优先级（暂未使用）（deprecated） |
| `group_regex` | `string` | 否 | — | 从消息中提取分组的正则表达式 |
| `group_index` | `string` | 否 | — | 分组排序的编号 |
| `key_group_regex` | `string` | 否 | — | 已改名：group_regex（兼容保留）（deprecated） |
| `key_group_index` | `string` | 否 | — | 已改名：group_index（兼容保留）（deprecated） |
| `task_dist_mode` | `string` | 否 | — | task 分发模式 取值：`HOST-BOUND`, `SLOT-BOUND`, `DEFAULT` |
| `start_task` | `string` | 否 | — | 给定初始消息；若为 'FILE:{filename}'，则将文件中每一行作为一个初始消息 |
| `initial_task_status` | `string` | 否 | — | task 的初始状态 取值：`READY`, `INITIAL` |
| `initial_slot_status` | `string` | 否 | — | slot 的初始状态 取值：`READY`, `OFF` |
| `retry_rules` | `string` | 否 | — | 基于退出码的重试规则："['exit_code_1:num_retries',...,'exit_code_n:num_retries']"。num_retries 缺省 1，退出码通配符为 '*' |
| `slot_recoverable` | `string` | 否 | — | 'yes'：支持将出错后已退出的 slot 从 'ERROR' 设置为 'READY'，以支持 slot 级重试 |
| `slot_max_retries` | string / number | 否 | — | slot 状态从 'TIMED-OUT' 设置为 'READY' 的重试次数 |
| `slot_timeout_minutes` | string / number | 否 | — | slot 未正常启动一直处于 'STARTING' 状态的超时（分钟），到期转为 'TIMEOUT'，缺省 15 分钟 |
| `task_global_timeout_scale` | string / number | 否 | — | task 全局超时恢复：相对 task_max_seconds 的倍数，缺省 2.0；恢复后状态码 123 |
| `router_index` | string / number | 否 | — | 多主路由实例的应用中，指定当前 module 发给第 n 个主路由。缺省 0 |
| `pod_id` | `string` | 否 | — | 标识本 module 属于 pod 管理，若消息来源的 pod 也有相同的 pod_id，则所有 task 标识为采用本地计算 |
| `task_dedup_cache_ttl_minutes` | string / number | 否 | — | 任务去重缓存的生存时间（分钟），缺省 30 分钟；通常需大于 task_max_seconds |
| `task_batch_size` | string / number | 否 | — | 批量获取任务的数量，slot 单批次读取的最大消息数，缺省 1 |
| `visible` | `string` | 否 | — | 在流水线逻辑图中是否可见，缺省 'yes' 取值：`yes`, `no` |
| `task_id_in_headers` | `string` | 否 | — | 该模块中 task 的 headers 中是否包含 task_id 值 取值：`no`, `yes` |
| `node_progress_gap` | `string` | 否 | — | 针对指定 module 同一组内 node 间运行同步，最快与最慢 node 间运行的 task 最大差值。格式示例：{"prefix1":4,"node_prefix2":6} |
| `vtask_role` | `string` | 否 | — | vtask 处理中的当前角色：head 是起始模块，tail 是结束模块，core 是算法模块 取值：`head`, `core`, `tail` |
| `vtask_size` | string / number | 否 | — | 可同时处于就绪/运行状态的 vtask 数量上限，app 解析时创建对应信号量及初值 |
| `vtask_size_sema_copy` | `string` | 否 | — | 为 vtask_size 信号量生成一个拷贝（名称前缀 ':'）用于编程控制 取值：`no`, `yes` |
