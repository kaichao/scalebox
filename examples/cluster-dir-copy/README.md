# cluster-dir-copy

本应用功能将数据根目录下的某个数据目录，拷贝到另一个集群。该应用包括3个标准模块+定制的message-router。


## message-router

主程序，起到各个标准模块粘结剂的作用，用shell编写。

## cluster-dir-list

列出源目录下的所有文件，以消息形式传递给cluster-file-copy。

## cluster-file-copy

将单个文件从源集群拷贝到目标集群。

## data-grouping-2d

可选模块，主要功能是对数据集的分组，了解文件拷贝的进度。
本应用拷贝的数据目录为2D数据集，可用其分组功能了解文件拷贝的进度，并可在文件拷贝完成后，设置应用状态为FINISHED。

## 模块运行

- PUSH

```sh

DIR_NAME=~fits-fz#Dec+6007_09_03~qiu REGEX_FILTER=^.+-M[0-9]{2}_000[1-2].+\$ scalebox app create

DIR_NAME=~fits-fz#Dec+6007_09_03~qiu REGEX_FILTER=^.+-M01_000[1-3].+\$ scalebox app create

```

- PULL
```sh
CLUSTER=qiu DIR_NAME=main~fits-fz#Dec+6007_09_03~ REGEX_FILTER=^.+-M03_000[1-3].+\$ scalebox app create
```