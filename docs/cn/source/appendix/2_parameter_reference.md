# 2. 参数语义手册

> 字段清单与类型见 [reference/app-yaml.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)；本章讲参数语义与 schema 未覆盖项。

## 2.1 应用配置参数表

### 2.1.1 app-parameters参数表

字段清单与类型见 [`reference/app-yaml.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md) 的 `app.parameters`。

### 2.1.2 module-arguments参数表

字段清单与类型见同一参考页的 `module.arguments`。`arguments` 允许扩展键；`slot_options` 的子选项（`always_running`、`reserved_on_exit`、`tmpfs_workdir`、`disable_local_mapping`、`disable_data_mapping`、`slot_on_head`）以该页为准。`PLAT_MODULE_NAME` 等环境变量见 §2.3。

### 2.1.3 module-parameters参数表

字段清单与类型见同一参考页的 `module.parameters`。schema 未覆盖的补充参数：

| 参数名 | 说明 |
|---|---|
| `slot_max_tasks` | slot 生命周期内累计执行 task 数上限（成败都算），达到后 slot 优雅退出并置 READY；0 表示不限；`PLAT_ALWAYS_RUNNING=yes` 的 slot 忽略此参数 |

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
| _vtask_attr_${name} | 调用方（应用）   | 应用自定义属性，vtask范围内可见：随任务进入wait-queue后由标准模块原样转发给vtask-head，并沿vtask管道向下游模块自动传播。平台只运输、不解释内容；${name}由应用定义，平台保留该子命名空间 |

其中，from_ip、from_module、from_module_last等，由系统自动生成。

- 针对HOST-BOUND的模块，需在代码中设定to_host；或通过模块的pod_id相等来设定。
- 针对SLOT-BOUND的模块，需在代码中设定to_slot

task-header可自定义参数，其中 ```_```开头的参数表示标准传递参数；```_vtask_```开头的参数为vtask参数，由系统维护；```_vtask_attr_```子命名空间保留给应用，用于声明vtask范围内可见的自定义属性（平台只运输、不解释内容）。

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

`slot get-parameter` 未设置时返回空字符串；`set-parameter` 的值以字符串存入 `t_slot.parameters`，可用于 slot 启动命令表达式 `((@p:...))`。

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
| PLAT_MODULE_NAME | 当前 module 名称（平台注入；用户程序无需访问）                                        |
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

