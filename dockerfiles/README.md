# dockerfiles

## Introduction
scalebox应用中常见的公用模块。

## file-related

### module list

|  模块名 | 模块描述 |
|  ----  | ---- | 
| dir-list  | 本地或远程目录的文件列表，支持本地目录、远端rsync目录、远端rsync-over-ssh目录等| 
| file-copy | 本地目录、远端rsync目录、远端rsync-over-ssh目录等 |
| rsync-copy  | 基于rsync-over-ssh、rsync、ssh的远端文件复制 |
| ftp-copy  | 基于ftp的远端文件复制 |
| rsyncd  | rsync的服务端 |


## data-grouping-2d
基于2维数据集的数据分组，是数据处理中的常见模式。将数据集中数据实体按id组织为2维数据集，并支持x、y方向上对数据进行分组。

## cron
定时消息生成模块，以启动后续模块。消息体一般用当前时间戳来表示。

## actuator
支持在标准actuator模块中自动生成自定义的公私钥。

## node-agent
计算节点上的模块，执行计算节点相关管理。

## cluster-head
