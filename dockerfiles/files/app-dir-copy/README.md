# dir-copy

本应用功能将数据根目录下的某个数据目录，拷贝到另一个集群。

## dir-list

列出源目录下的所有文件，以消息形式传递给file-copy。

## file-copy

将单个文件从源集群拷贝到目标集群。

## 模块运行实例

- local to ssh-server
```sh
SOURCE_URL=/ TARGET_URL=scalebox@10.255.128.1/tmp DIR_NAME=etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create

SOURCE_URL=/etc TARGET_URL=scalebox@10.255.128.1/tmp DIR_NAME=postfix REGEX_FILTER=^.*cf\$ scalebox app create

SOURCE_URL=/etc/postfix TARGET_URL=scalebox@10.255.128.1/tmp DIR_NAME=. REGEX_FILTER=^.*cf\$ scalebox app create
```

- ssh-server to local
```sh
SOURCE_URL=scalebox@10.255.128.1/ TARGET_URL=/tmp DIR_NAME=etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create

SOURCE_URL=scalebox@10.255.128.1/etc TARGET_URL=/tmp DIR_NAME=postfix REGEX_FILTER=^.*cf\$ scalebox app create

SOURCE_URL=scalebox@10.255.128.1/etc/postfix TARGET_URL=/tmp DIR_NAME=. REGEX_FILTER=^.*cf\$ scalebox app create
```



```sh

SOURCE_URL=scalebox@159.226.237.136:10022/raid0/tmp/mwa/tar1301240224 TARGET_URL=/data1/tmp DIR_NAME=1301240224 REGEX_FILTER= scalebox app create

```