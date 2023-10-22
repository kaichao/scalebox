## 模块介绍
基于rsync-over-ssh，实现不同集群间的单文件同步

## 环境变量

| 环境变量   | 描述  |
|  ----  | ----  |
| SOURCE_CLUSTER | 源端集群名称 |
| TARGET_CLUSTER | 目标端集群名称|
| JUMP_SERVER_OPTION | 取值为：'source'/'target'/'source/target'，是否启用source/target端的跳板参数|
| ZSTD_CLEVEL | 若为非空的整数，则为启用zstd的实时压缩/解压的压缩等级 |
| ENABLE_LOCAL_RELAY| 若为'yes'，则通过容器内工作目录中转文件 |

## 输入消息格式：

| 条件   | 描述  |
|  ----  | ----  |
|  若SOURCE_CLUSTER非空、TARGET_CLUSTER为空  | 从远端拷贝文件到本地，输入消息为相对于CLUSTER_DATA_ROOT的文件路径  |
|  若SOURCE_CLUSTER为空、TARGET_CLUSTER非空  | 从本地拷贝文件到远端，输入消息为相对于CLUSTER_DATA_ROOT的文件路径  |
|  若SOURCE_CLUSTER非空、TARGET_CLUSTER非空  | 执行双远端文件拷贝，从SOURCE_CLUSTER端拷贝到本地，再从本地拷贝到TARGET_CLUSTER端，输入消息为相对于CLUSTER_DATA_ROOT的文件路径  |
|  若SOURCE_CLUSTER为空、TARGET_CLUSTER为空  | 以消息体格式决定拷贝模式。  <br/> <SOURCE_CLUSTER>~<relative_path> : 远端拷贝到本地；<br/>~<relative_path>~<TARGET_CLUSTER>:本地拷贝到远端 |

## 用户应用的退出码
- 0 : OK 
- 1 : source file not exist
- 2 : target dir not allowed
- 

## 输出消息格式
- 若退出码为0，则输出与输入消息相同的消息。
- 退出码非0，则不输出消息
