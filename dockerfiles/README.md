# dockerfiles

## Introduction
scalebox应用中常见的公用模块。

## file-related

### module list
|  模块名 | 模块描述 |
|  ----  | ---- | 
| cluster-dir-list  | 集群内目录的文件列表   | 
| cluster-file-copy  | 基于rsync-over-ssh的跨集群文件复制 |
| dir-list  | 本地或远程目录的文件列表，支持本地目录、远端rsync目录、远端rsync-over-ssh目录、远端ftp目录等| 
| rsync-copy  | 基于rsync-over-ssh、rsync的远端文件复制 |
| ftp-copy  | 基于ftp的远端文件复制 |
| rsyncd  | rsync的服务端 |

- Cluster related configuration
| 配置参数   | 描述  |
|  ----  | ----  |
| base_data_dir | The cluster's data base directory |
| storage_endpoint | <user>@<ip-addr>:[<port>],default port is 22 |
| relay_endpoint | <user>@<ip-addr>:[<port>]<user>@<ip-addr>:[<port>] |


## data-grouping-2d
基于2维数据集的数据分组，是数据处理中的常见模式。将数据集中数据实体按id组织为2维数据集，并支持x、y方向上对数据进行分组。
## cron
定时消息生成模块，以启动后续模块。消息体一般用当前时间戳来表示。

## actuator
支持在标准actuator模块中自动生成自定义的公私钥。

## node-manager
单个节点上的actuator，用于在当前节点上启动slot。
