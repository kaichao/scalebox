# dir-copy

本应用功能将数据根目录下的某个数据目录，拷贝到另一个集群。该应用包括3个标准模块+定制的message-router。


## message-router

主程序，起到各个标准模块粘结剂的作用，用shell编写。

## dir-list

列出源目录下的所有文件，以消息形式传递给cluster-file-copy。

## rsync-copy

将单个文件从源集群拷贝到目标集群。

## 模块运行

```sh
DIR_NAME=/etc/postfix~scalebox@10.255.128.1/tmp REGEX_FILTER=^.*cf\$ scalebox app create

DIR_NAME=scalebox@10.255.128.1/etc/postfix~/tmp/postfix scalebox app create

DIR_NAME=ftp://ftp.ncbi.nlm.nih.gov/1000genomes/ftp/release/2008_12~/tmp/2008_12 scalebox app create

```


```sh
SOURCE_URL=scalebox@10.255.128.1/etc/postfix TARGET_URL=/tmp/postfix DIR_NAME=.  scalebox app create
```