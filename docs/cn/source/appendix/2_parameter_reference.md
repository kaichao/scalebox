# 2. 参数配置手册

## 2.1 应用配置参数表

### 2.1.1 app-parameters参数表
| 参数名称             |   含义                            |
| ------------------- | -------------------------------- |
| initial_status      | 'RUNNING' / 'INITIAL' / 'PAUSED' |
| main_router         |  主路由模块名                      |
| default_idle_polls  | 所属module的缺省的空闲轮询次数       |
| is_cluster_admin    |                                  |
| default_sleep_count | 各模块sleep_count的默认值          |
| global_directories  | 全局目录列表，用于所属模块的读写数据量分类统计。（local/global/tmpfs）|
| slot_group          | map形式的多节点的多module-slot配置，用于新增节点的slot自动创建。例：'{"module0":n0,"module1",n1}' |

### 2.1.2 module-arguments参数表

| 参数名                  | 标准环境变量             | 含义                                                                    |
| ---------------------- | ---------------------- | ---------------------------------------------------------------------- |
| grpc_server           | GRPC_SERVER            | controld的服务端点（endpoint），```${ip_addr}:${port}```，port缺省值为50051 |
| code_path             |                        | 模块代码的目录，通过容器的数据卷Volume映射到容器内/app/bin                    |
| local_ip_index        | LOCAL_IP_INDEX         | hostname -I'返回IP地址列表，该参数指定列表中的第n个IP地址作为本机IP地址。 (agent中暂未使用？)      |
|                       | CLUSTER                | 所在的集群名 (与CLUSTER_NAME、_CLUSTER的区别？)                                                      |
|                       | PLAT_MODULE_NAME       | 当前module名称                             |
|                       | MODULE_ID              | module-id                                 |
|                       | FROM_MODULE               | module-name                                  |
|                       | FROM_IP                | from-ip                                         |
|                       | SINK_MODULE               | 缺省sink_module的名称                          |
|                       | IS_SINGULARITY         | 容器引擎为singularity或apptainer      |
| task_max_seconds      | PLAT_TASK_MAX_SECONDS       | 每个task运行中超时设置的秒数，若运行时间超过该时限，task运行中断，返回超时码124；全局超时返回123。 |
| task_min_seconds      | PLAT_TASK_MIN_SECONDS       | 每个task运行的最小秒数，若运行时间连续多次低于此值，则判定slot为GREEDY异常，并退出 |
| poll_interval_seconds | PLAT_POLL_INTERVAL_SECONDS | Slot 检查新任务的轮询间隔。slot睡眠并定期检查task可用，该参数指定以秒计的时间间隔，缺省值为6秒。 |
| max_idle_polls        | PLAT_MAX_IDLE_POLLS    | Slot 在退出前可进行的最大空闲轮询次数。slot退出前的最多睡眠次数。缺省值为100（10分钟）                              |
| dir_quota_gb          | PLAT_DIR_QUOTA_GB      | 标准流控参数，目录本身的容量配额限制。用于指定目录以GB计的最大空间。格式为： ```'{"/dir-1":10,"/dir-2":100}'```  |
| free_space_gb         | PLAT_FREE_SPACE_GB     | 标准流控参数，目录所在磁盘需要保留的最小空间。用于指定目录所在分区以GB计的最小保留空间。格式为```'{"/dir-3":10,"/dir-4":100}'``` |
| heartbeat_seconds     | PLAT_HEARTBEAT_SECONDS | 以秒计的心跳间隔，缺省值为60；若为非正整数，则禁用心跳操作 |
| output_text_size      | PLAT_OUTPUT_TEXT_SIZE  | task运行记录t_task_exec中，大文本字段（stdout/stderr/auxout）的最大字节数。缺省值为65535，最大值可以为10MB(for varchar) 或1GB(for text) |
| text_trunc_mode       | PLAT_TEXT_TRUNC_MODE   | HEAD'/'TAIL', default value is 'HEAD'，头截断，保留末尾部分  |
| timezone_mode         |                        | HOST'/'UTC'/'NONE'                                                |
| slot_options          |                        | 逗号分隔的slot选项                                                  |
|  - always_running     | PLAT_ALWAYS_RUNNING    | 设定slot一直运行，不主动退出（一般仅用于调试）    |
|  - reserved_on_exit   |                        | slot退出后，保留容器，以便排错。(docker-only，命令行去掉--rm)    |
|  - tmpfs_workdir      |   TMPFS_WORKDIR        | 用tmpfs文件系统存放工作目录/work（针对docker，解析后的命令行加上--tmpfs /work；针对singularity，解析后增加环境变量TMPFS_WORKDIR=yes）|
|  - disable_local_mapping |                     | 不生成本地根目录到容器内/local_data_root的映射  |
|  - disable_data_mapping  |                     | 不生成集群数据目录到容器内/cluster_data_root的映射 |
|  - slot_on_head          |                     | 在头节点上生成单个slot实例                         |

### 2.1.3 module-parameters参数表

| 参数名                  | 说明                                                                        |
| -------------------- | ---------------------------------------------------------------------------- |
| priority             | 优先级(暂未使用)                                            |
| key_group_regex      | 从消息中提取分组的正则表达式  （改名为group_regex?）            |
| key_group_index      | 分组排序的编号               (改名为group_index?)            |
| task_dist_mode       | task分发模式，'HOST-BOUND'/'SLOT-BOUND'/'DEFAULT'   |
| start_task        | 给定初始消息，若为'FILE:{filename}'，则将文件中每一行作为一个初始消息  |
| initial_task_status  | task的初始状态，'READY'/'INITIAL'                               |
| initial_slot_status  | slot的初始状态，'READY'/'OFF'                                   |
| retry_rules          | 基于退出码的重试规则<br>```"['exit_code_1:num_retries',...,'exit_code_n:num_retries']"```。num_retries缺省值为1，退出码通配符为'*'。拟修改格式为```{"exit_code_1":num_retries',...,"exit_code_n":num_retries}``` |
| slot_recoverable     | 'yes'，支持将出错后已退出的slot从'ERROR'设置为'READY'，以支持slot级重试 (以slot_max_retries替换？,TIMED-OUT/GREEDY分别处理)  |
| slot_max_retries     | slot状态从'TIMED-OUT'设置为'READY'的重试次数(?)          |
| slot_timeout_minutes | 若slot未正常启动，则一直处于'STARTING'状态。设置以分钟计的timeout，到期后将状态转换为'TIMEOUT'。缺省值为15分钟。对于不允许重复启动的slot实例（用GPU等），可设置较大值。 |
| slot_max_tasks | slot生命周期策略参数：设置slot退出前累计执行的task数上限（成败都算），达到后slot优雅退出，退出后slot置READY，新task到来时由actuator重启新容器；缺省值0表示不限；PLAT_ALWAYS_RUNNING=yes的slot忽略此参数 |
| task_global_timeout_scale | 若外部原因（slot异常退出等）导致task一直处于运行状态（状态码-3）。通过全局超时设置，恢复task状态码为123。该值为相对task_max_seconds的倍数，缺省值为2.0。全局退出设定返回码123。 拟改为task_timeout_scale_factor ？|
| router_index | 多主路由实例的应用中，指定当前module发给第n个主路由。缺省值为0，通常设置值>0，以指定特定main-router实例  |
| pod_id               | 标识本module属于pod管理，若消息来源的pod也有相同的pod_id，则所有task标识为采用本地计算（task_dist_mode为HOST_BOUND）  |
| task_dedup_cache_ttl_minutes | 任务去重缓存的生存时间，在高负载时需设置。设定重复task-id检测的cache过期时间（分钟数），缺省值为30分钟，清除时间为n+1分钟。避免出现同一task的多次分发。通常情况下，其时间需大于```task_max_seconds```的值。 |
| task_batch_size       | 批量获取任务的数量。设置slot单批次读取的最大消息数，缺省值为1。针对运行时长在5秒以内的任务，可设置批量读取消息，避免读取频繁而导致server端过载、数据不一致。 |
| visible              | 在流水线逻辑图中是否可见。缺省值为'yes'                                          |
| task_id_in_headers   | 'no'/'yes'，该模块中task的headers中，包含task_id值。|
| node_progress_gap    | 标准流控参数，针对指定module同一组内node间运行同步，最快node与最慢node间的运行的task最大差值，其值为整数。在对应slot生成时，自动创建对应信号量，其名称为```node_progress:${mod_name}:${hostname}```，初值为0。该参数格式示例为```{"prefix1":4,"node_prefix2":6}```。该参数拟替换为node_progress_max_gap？ |
| vtask_role   | 'head' / 'core' / 'tail'，vtask处理中当前的角色，仅针对非路由模块有效。head是vtask处理的起始模块；tail是结束模块；core是算法模块。|
| vtask_size   | 标准流控参数，定义可同时处于就绪/运行状态的vtask 数量上限，在app解析时，创建对应信号量及初值。<br/>用于全局vtask流控的信号量名为：```vtask_size:${mod_name}```；<br/>用于SLOT-BOUND（GROUP-BOUND vtask）的vtask流控信号量名为：```group_vtask_size:${mod_name}:${group_seq}```；<br/>用于HOST-BOUND的vtask流控信号量名为：```host_vtask_size:${mod_name}:${hostname}```|
| vtask_size_sema_copy | 'no' / 'yes'。为vtask_size信号量生成一个拷贝，用于编程控制。拷贝信号量名称为原始信号量前加冒号```:```，初值与原始信号量相同。|

### 2.1.4 task-headers标准参数表

| 参数名称          | 写入方           | 含义                                            |
| ---------------- | --------------- | ----------------------------------------------- |
| to_ip            | 调用方          | 当前task的待处理主机ip                             |
| to_host          | 调用方          | 当前task的待处理主机名(t_host主键)                  |
| from_ip          | 系统自动        | 生成task的主机ip                                 |
| from_ip_last     | 系统自动        | 若from_module为主路由，主机ip                     |
| from_host        | 系统自动        | 生成task的主机名(t_host主键)                      |
| from_module      | 系统自动        | 生成task的模块名                                 |
| from_module_last | 系统自动        | 若from_module为主路由，主路由之前的from_module     |
| from_cluster     | 系统自动        | 跨集群应用的上一集群名                             |
| to_slot          | 调用方          | SLOT-BOUND模块中当前task的待处理slot_id           |
| to_slot_index    | 调用方          | 以本机seq表示的slot_id                           |
| repeatable       | 调用方          | 缺省task在指定时间内不可重复分发，缺省值可通过module的task_cache_expired_minutes参数定制；在retry操作、特定场景下，需支持消息的重复分发，则设为该参数'yes'|
| action           | 调用方          | 模块action路由（tape-mgr/tape-io类模块用）          |
| conflict_action  | 调用方          | 任务冲突处理策略（task add的--conflict-action参数）  |
| slot_last_task   | 调用方          | 值为'yes'时声明本task是当前slot的最后一个task：agent执行完该task后退出（处理完当前已拉取批次），不处理后续task。退出后slot置READY，新task到来时由actuator重启新容器。PLAT_ALWAYS_RUNNING=yes的slot忽略此标记 |

其中，from_ip、from_module、from_module_last等，由系统自动生成。

- 针对HOST-BOUND的模块，需在代码中设定to_host；或通过模块的pod_id相等来设定。
- 针对SLOT-BOUND的模块，需在代码中设定to_slot

task-header可自定义参数，其中 ```_```开头的参数表示标准传递参数；```_vtask_```开头的参数为vtask参数，由系统维护。

## 2.2 资源配置参数表

### 2.2.1 cluster-parameters参数表
| 参数名称            |   含义                                                  |
| ------------------ | ------------------------------------------------------ |
| data_root          |  容器内/cluster_data_root目录（改为data_root?）           |
| code_dir           |  未使用                                                 |
| uname              | ssh登录用户名                                            |
| port               | ssh登录用户名                                            |
| grpc_server        | controld连接信息，${controld_ip}[:${port}]，缺省端口号50051 |
| remote_grpc_server | 供跨集群应用远端路由连接的controld连接信息                    |
| remote_check_mode  | 'ip':精确IP比对，测试集群、同网段多集群；'subnet'：子网匹配，生产集群、跨网段部署。|
| initial_status     | 'ON'/'OFF' |

### 2.2.2 host-parameters参数表

| 参数名称         |   含义                         |
| --------------- | ----------------------------- |
| uname           | ssh登录用户名                   |
| port            | ssh的端口号                     |
| node_slot       | 该节点node-agent的slot号        |
| group_id        | 该节点所属节点组的编号            |
| slurm_node      | 在slurm调度系统中对应的节点编号    |
| reg_time        | 在scalebox中注册时间             |
| slot_module_id  | 在slurm调度系统中，node-agent的slurm module id |
| default_runtime | 缺省容器引擎,'docker'/'singularity'/'podman' |
| use_home_tmp    | 针对不允许/tmp下运行程序的节点(I/O节点)，以$HOME/tmp代替 /tmp作为远程启动slot的临时目录。|
| ssh_cepher      | 'aes128-gcm@openssh.com' / 'aes256-gcm@openssh.com' / 'chacha20-poly1305@openssh.com' |

### 2.2.3 slot-parameters参数表

| 参数名称      | 含义                                     |
| ------------ | --------------------------------------- |
| reg_time     | 注册时间                                 |
| last_access  | 最后访问时间                              |
| group_prefix | slot的分组前缀，在HOST-BOUND模块中的分组slot |

## 2.3 环境变量参考

- 模块编程的标准环境变量表

| 环境变量名        | 说明                                                                              |
| --------------- | --------------------------------------------------------------------------------- |
| ACTION_CHECK    |  自定义准入控制脚本（check）的路径。若返回值为0，OK；非0，准入受限。                        |
| ACTION_RUN      |  自定义主算法脚本（run）的路径                                                         |
| ACTION_SETUP    |  自定义初始设置脚本（setup）的路径。若返回值为0，OK；非0，初始化失败，设置当前slot状态为ERROR。|
| ACTION_TEARDOWN |  自定义结束退出脚本（teardown）的路径。过返回值为0，OK；非0，退出失败，设置当前slot为ERROR。  |
| APP_ID          |                                                                                   |
| MODULE_ID       |                                                                                   |
| SLOT_ID         |                                                                                   |
| CLUSTER         | CLUSTER                                                                           |
| WORK_DIR        | 本地工作目录，缺省为/work                                                            |
| LOG_LEVEL       | 'info'/'debug'/'trace'，trace原则上仅用于单模块调试、排错。                            |
| LOCAL_IP        | 本机IP地址                                                                         |
| FROM_IP         | headers中已包含？                                                                  |
| FROM_MODULE     | headers中已包含？                                                                  |
| SINK_MODULE     | 路由应用的非路由模块指向路由模块；无路由应用指向下一个模块。                               |
| GRPC_SERVER     | server端的grpc地址，格式：```server_name[:port]```，缺省port为50051。                |
| PGHOST          | PGHOST, query job-attr-set from db                                               |
| IS_SINGULARITY  | IS_SINGULARITY                                                                  |
| LOCAL_SHMDIR    | 本地tmpfs下工作目录（/dev/shm）                                                    |
| LOCAL_TMPDIR    | 本地/tmp下工作目录                                                                |
| TMPDIR_GROUP    | TMPDIR中分组号（数字）                                                            |
| TMPFS_WORKDIR   | 本地工作目录设置使用tmpfs，以提高大文件加载性能                                        |
| SLOT_ROLE       | ''/'group'，组计算模式中的组slot，可跨节点获取组内所有task                            |

PLAT_ 开头的所有环境变量，用户程序无需访问。

