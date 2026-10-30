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
| `slot_max_tasks` | slot 生命周期内累计执行 task 数上限（成败都算），达到后 slot 优雅退出并置 READY；0 表示不限；`slot_options` 含 `always_running` 的 slot 忽略此参数 |

### 2.1.4 task-headers标准参数表

task header 的名称、写入方、传播规则与转义以
[`reference/task-headers.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/task-headers.md) 为唯一信源；本节只保留两条落点约定：

- 针对 HOST-BOUND 的模块，需在代码中设定 `to_host`；或通过模块的 `pod_id` 相等来设定。
- 针对 SLOT-BOUND 的模块，需在代码中设定 `to_slot`。

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

