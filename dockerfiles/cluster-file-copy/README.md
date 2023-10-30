## 一、模块介绍
基于rsync-over-ssh，实现不同集群间的单文件同步。模块消息的文件路径为从CLUSTER_DATA_ROOT开始的相对路径。

## 二、环境变量

| 环境变量   | 描述  |
|  ----  | ----  |
| SOURCE_CLUSTER | 固定源端集群的名称（可在消息中指定） |
| TARGET_CLUSTER | 固定的目标端集群名称（可在消息中指定）|
| JUMP_SERVER_OPTION | 取值为：'source'/'target'/'source/target'，是否启用source/target端的跳板参数|
| ZSTD_CLEVEL | 若为非空的整数，则为启用zstd的实时压缩/解压的压缩等级 |
| ENABLE_LOCAL_RELAY| 若为'yes'，则通过容器内工作目录中转文件 |

## 三、输入消息格式

| 条件   | 描述  |
|  ----  | ----  |
| <source_cluster>~<relative-path>#<relative-path-file>~ | 从远端拷贝文件到本地集群 |
| ~<relative-path>#<relative-path-file>~<target_cluster>  | 拷贝本地文件到远端集群  |
|  <relative-path>#<relative-path-file> | 由环境变量SOURCE_CLUSTER、TARGET_CLUSTER确定拷贝方向，两个环境变量一个为空，一个非空 |


## 四、应用退出码

| 代码   | 描述  |
|  ----  | ----  |
| 0 | OK | 
| 1 | source file not exist |
| 2 | target dir not allowed |

## 五、输出消息格式
- 若退出码为0，则输出与输入消息相同的消息。
- 退出码非0，则不输出消息
