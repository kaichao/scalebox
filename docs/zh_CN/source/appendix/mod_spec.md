# 3. 模块定义技术规范

为解耦模块算法，以非侵入式调用，主算法及相关辅助程序作为独立应用程序，封装在容器镜像中。
- 模块定义规范：通过Dockerfile定义容器化规范，以此为基础构建容器镜像。


## 3.1 模块介绍

模块（Job）是Scalebox应用程序的基本单元，可部署在单个计算集群的多个计算节点（Host）上，每节点上可运行一个或多个实例（Slot）；每个Slot可处理多个计算任务（Task）。

从Slot视角看，每个Task为一条该模块对应的待处理消息。其处理过程为读取消息，进行计算处理，纪录计算结果，并发送消息到后续模块（sink-job），做进一步处理。

从实现角度看，模块为容器化封装的算法代码。



附图 模块结构示意图



- 模块内主要脚本

| 环境变量 | 脚本名 | 脚本说明 |
| --------------- | ----------- | ------------ |
| ACTION_RUN      | run.sh      | 主算法脚本      |
| ACTION_CHECK    | check.sh    | 定制流控脚本。返回值为0，OK；非0，流控限制 |
| ACTION_SETUP    | setup.sh    | 初始设置脚本。返回值为0，OK；非0，初始化失败，设置对应slot为错误。 |
| ACTION_TEARDOWN | teardown.sh | 结束退出脚本。返回值为0，OK；非0，退出失败，设置对应slot为错误。 |


- 脚本返回码



## 3.2 消息体body格式

task含有body、headers。body是task的唯一标识，headers用于task标识辅助信息。

headers中包含运行调度信息（to_host、to_slot）

```

消息体body ::= <不含逗号的字符串>|<json字符串>。
扩展消息体 ::= <消息体body>,<json扩展消息头>

```

在task-add命令、消息文件中，支持用扩展消息体。

## 3.3 消息头格式

- 消息头名一般为大写字母、数字、下划线组成。

### 标准消息头

### 2.3.2 标准task-headers参数表

| 参数名称      | 含义 |
| --------------- | ----------------------------------------------- |
| to_ip           | 当前task的待处理主机ip                             |
| to_host         | 当前task的待处理主机名(t_host主键)                  |
| from_ip         | 生成task消息的主机ip                               |
| from_host       | 生成task消息的主机名(t_host主键)                    |
| from_job        | 生成task消息的job名                               |
| from_job_last   | 若from_job为消息路由，消息路由之前的from_job         |
| to_slot         | 当前task的待处理slot_id                           |
| repeatable      | 缺省task对应消息在指定时间内不可重复分发，缺省值可通过job的task_cache_expired_minutes参数定制；在retry操作、特定场景下，需支持消息的重复分发，则设为该参数'yes'|
| slot_broadcast  | 仅用于cli的命令行参数，针对所有slot，按广播形式生成一组消息 |
| host_broadcast  | 仅用于cli的命令行参数，针对所有host，按广播形式生成一组消息 |

其中，from_ip、from_job、from_job_last等，由系统自动生成。

- 针对HOST-BOUND的job，需在代码中设定to_host；或通过job的pod_id相等来设定。
- 针对SLOT-BOUND的job，需在代码中设定to_slot


### 自定义消息头

- 自定义header标识task的辅助信息。
- 通过scalebox task add 中，--header header_name=header_value

- 比如标准模块file-copy中，定义了自定义消息头source_url、target_url，用于标识源端、目标端的URL。



## 3.4 环境变量

- 环境变量可设置模块运行的初始参数，还可以设定消息头的缺省值。

- 环境变量名一般为大写字母、数字、下划线组成。

- 标准环境变量表

| 环境变量名 |     环境变量说明 |
| --------------- | ------------ |
| ACTION_RUN      |   主算法脚本      |
| ACTION_CHECK    | 流控脚本。返回值为0，OK；非0，流控限制 |
| ACTION_SETUP    | 初始设置脚本。返回值为0，OK；非0，初始化失败，设置对应slot为错误。 |
| ACTION_TEARDOWN | 结束退出脚本。返回值为0，OK；非0，退出失败，设置对应slot为错误。 |
| JOB_ID          | 当前JOB_ID |
| JOB_NAME          | 当前JOB_NAME |
| SINK_JOB          | SINK_JOB |
| CLUSTER          | CLUSTER |
| GRPC_SERVER          | GRPC_SERVER |
| TASK_TIMEOUT_SECONDS          | TASK_TIMEOUT_SECONDS |
| OUTPUT_TEXT_SIZE          | OUTPUT_TEXT_SIZE |
| TEXT_TRANC_MODE          | TEXT_TRANC_MODE |
| LOCAL_IP_INDEX          | LOCAL_IP_INDEX |
| LOCAL_IP          | LOCAL_IP |
| HOST_NAME          | hostname in scalebox, used in progress-counter_* |
| INTERVAL_SECONDS          | INTERVAL_SECONDS |
| PGHOST          | PGHOST, query job-attr-set from db  |
| HEART_BEAT_SECONDS          | HEART_BEAT_SECONDS |
| DIR_LIMIT_GB          | DIR_LIMIT_GB, directory capacity limit, sep sign is '~' |
| DIR_FREE_GB          | DIR_FREE_GB |
| MAX_SLEEP_COUNT          | MAX_SLEEP_COUNT |
| SLEEP_INTERVAL_SECONDS          | SLEEP_INTERVAL_SECONDS |
| IS_SINGULARITY          | IS_SINGULARITY |
| WORK_DIR          | WORK_DIR, 工作目录，缺省为/work |


## 3.5 模块内目录

### 数据目录映射
- 临时目录 ```/dev/shm```、```/tmp```，自动映射到算法容器中
- 集群数据目录映射到算法容器中的```/cluster_data_root```，同时设置环境```CLUSTER_DATA_ROOT```指向实际的集群数据目录。
- 计算节点的```/```映射到容器中的```/local```

如需额外目录映射，通过```paths```再做定制映射。

### 代码目录

按以下顺序：
- 环境变量：ACTION_RUN、ACTION_CHECK、ACTION_SETUP、ACTION_TEARDOWN
- /app/bin/{run.sh,check.sh,setup.sh,teardown.sh}
- /app/share/bin/{run.sh,check.sh,setup.sh,teardown.sh}

argument:code_path，映射为容器中的/app/bin。


## 3.6 封装脚本

支持用多种语言实现，推荐使用shell。

- run.sh
- setup.sh
- teardown.sh
- check.sh

### 数据目录路径

- Cluster数据目录（base_data_dir），可在容器中用/data引用。
- 本机目录：在目录前加上/local，可访问


## 3.7 返回码

- 通用返回码
  - 0：OK
  - >0：错误

- 主算法返回码

## 3.8 模块的构建规范


### 3.8.1 Dockerfile示例

```Dockerfile
FROM my-image


```

