# 3. 模块设计与实现

模块分为算法模块、辅助模块、集成模块等。

![slot_arch](../diagrams/slot_arch.drawio.svg)

## 3.1 模块及代码规范


### 代码结构

（模块结构图）

支持事件驱动式结构。

模块连接代码：

算法代码；

- 减少中间存储量，可以生成部分结果，尽可能早生成最终结果。

### 数据目录
在数据处理量较大的计算模块中，通常需按照计算过程中的使用特性，将目录分类并支持对其灵活配置，以适配不同的计算环境中不同类型的存储形式（内存缓存、本机SSD、本机HDD、网络存储等）,进而提升数据I/O效率。

主要目录如下：
- 代码目录/配置目录：位于共享存储或容器镜像内。通常数据量较小，重点考虑可灵活配置。
- 输入数据目录：综合考虑读取频度、数据总量，可灵活配置不同数据存储形式。如配置为本机存储，则输入数据由上游模块生成，或通过上游模块从共享存储传输过来。
- 中间结果目录：通常位于本机存储
- 输出结果目录：综合考虑写出频度、数据总量，可灵活配置不同数据存储形式。如配置为本机存储，则输出数据可供下游游模块读取，或通过下游模块传输至共享存储。

### 配置参数

在运行时需调整的配置参数通常以环境变量形式传递。在外部编排系统中可以修改配置参数。(运行线程数、)

### 配置参数与task头



## 3.2 模块镜像定义


## 3.3 模块脚本

sidecar模式：
- run: 任务的单次运行
- check: 任务运行的前置条件检测
- setup: 设置环境
- teardown: 清除环境

### 3.3.1 算法运行run.sh

### 3.3.2 流控检测check.sh

### 3.3.3 初始设置setup.sh

### 3.3.4 结束退出teardown.sh

## 3.4 模块单元测试

用户程序：用任意语言写；

集成程序：一般用bash写。将用户程序的结果写回。

用户程序与集成程序间接口：
- 运行结果文件：${WORK_DIR}/task-exec.json
- 输出task文件：${WORK_DIR}/sink-tasks.txt。原始文件为：${WORK_DIR}/messages.txt
- 用户自定义时间戳：${WORK_DIR}/timestamps.txt
- 运行附加属性文件：${WORK_DIR}/extra-attributes.txt，存放于t_task_exec表中extras的extra_attributes中。
- 用户数据文件：${WORK_DIR}/custom-out.txt
- 输入文件列表：${WORK_DIR}/input-files.txt
- 输出文件列表：${WORK_DIR}/output-files.txt
- 待删除文件列表：${WORK_DIR}/removed-files.txt

- 批量task文件：${WORK_DIR}/batch-tasks.txt，单次处理多个消息，可用于消息路由中批量消息的高效处理。设定环境变量TASK_BATCH_SIZE=n，缺省值为1.

- 多个初始化消息

## 3.5 迭代优化

用户程序（run.sh）运行结束后，agent对用户程序的运行进行统计，并纪录到核心数据库中。两者之间主要通过用户程序的标准输出（stdout）、标准错误（stderr）以及以下文件来交换信息：

|  文件名                  |  文件说明                                   |
| ----------------------- |  ----------------------------------------- |
| /work/task-exec.json    | 主控制文件，以json形式，用户程序运行结果         |
| /work/sink-tasks.txt    | 后续任务列表文件，每行一个任务                  |
| /work/timestamps.txt    |  时间戳文件，用于调试程序、测试程序性能时使用     |
| /work/input-files.txt   | 输入文件列表，用于统计输入文件字节数             |
| /work/output-files.txt  | 输出文件列表，用于统计输出文件字节数             |
| /work/removed-files.txt | 待删除文件列表，一般用于统计文件字节数后再删除     |
| /work/auxout.txt        | 辅助输出文件，纪录在最终的任务执行(task_exec)表中 |

- 采集输入、输出字节数
- 增加时间戳


## 3.6 标准模块的可编程特性

- 标准模块：其功能脚本可放在/app/share/bin下，子模块的功能脚本在/app/bin下。模块识别规范是优先使用/app/bin，再搜索/app/share/bin目录下；
- main-router任务体格式定制：按标准模块的任务体格式定制，便于使用标准模块功能；

这样可充分利用标准模块的功能。

## 3.7 关键模块的task运行排序

task运行排序是scalebox应用运行的重要基础。

对于多节点协同处理的模块来说，排序可将同一时段上所有数据在不同计算资源上并行运行，以在后续模块中做协同（分组等）。

通过设置以下参数，实现排序。

### 3.7.1  排序标签
  消息头中sort_tag，是用于消息排序的专用标签，具有最高高优先级，通常由main-router设置

### 3.7.2 消息分组号
- group_regex：正则表达式，从消息体中提取相关分组字符串。
- group_index：正则表达式对应的分组编号。

### 3.7.3 任务处理顺序

若未设置前述排序方式，则缺省按消息生成的顺序进行处理

- 幂等性：节点本地计算存在节点失败的可能性，导致本地存储失效，无法保证task级的幂等性，在跨模块的vtask层级上实现幂等性。


### 模块编程的标准环境变量表

| 环境变量名        | 说明                                                            |
| --------------- | --------------------------------------------------------------- |
| ACTION_CHECK    |  自定义check脚本路径                                              |
| ACTION_RUN      |  自定义run脚本路径                                                |
| ACTION_SETUP    |  自定义setup脚本路径                                              |
| ACTION_TEARDOWN |  自定义teardown脚本路径                                           |
| APP_ID          |                                                                 |
| MODULE_ID       |                                                                 |
| SLOT_ID         |                                                                 |
| TASK_ID         |                                                                 |
| WORK_DIR        | 本地工作目录                                                      |
| LOG_LEVEL       | 'info'/'debug'/'trace'，trace原则上仅用于单模块调试、排错。           |
| LOCAL_IP        | 本机IP地址                                                        |
| FROM_IP         | headers中已包含？                                                 |
| FROM_MODULE     | headers中已包含？                                                 |
| SINK_MODULE     | 路由应用的非路由模块指向路由模块；无路由应用指向下一个模块。               |
| REMOTE_SERVER   | 跨集群应用的远端grpc_server地址                                     |
| GRPC_SERVER     | server端的grpc地址，格式：```server_name[:port]```，缺省port为50051。|
| LOCAL_SHMDIR    | 本地tmpfs下工作目录（/dev/shm）                                    |
| LOCAL_TMPDIR    | 本地/tmp下工作目录                                                 |
| TMPDIR_GROUP    | TMPDIR中分组号（数字）                                             |
| TMPFS_WORKDIR   | 本地工作目录设置使用tmpfs，以提高大文件加载性能                        |
| SLOT_ROLE       | ''/'group'，组计算模式中的组slot，可跨节点获取组内所有task             |

PLAT_ 开头的所有环境变量，用户程序无需访问。

