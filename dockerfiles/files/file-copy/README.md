# file-copy

## 一、模块介绍

支持以SSH/RSYNC_OVER_SSH/RSYNC等多种方式，完成跨节点、跨集群的文件传输。

## 二、 参数/环境变量

| 参数名         | 环境变量名       | 说明                                                     |
| ------------  | -------------- | -------------------------------------------------------- |
|  source_url   | SOURCE_URL     | 源端URL前缀值                                              | 
|  target_url   | TARGET_URL     | 目标端URL前缀值                                            | 
|               | SOURCE_MODE    | 源端URL模式                                               | 
|               | TARGET_MODE    | 目源端URL模式                                             | 
| source_jump   | SOURCE_JUMP    | 源端跳板                                                  | 
| target_jump   | TARGET_JUMP    | 目标端跳板                                                | 
|               | KEEP_SOURCE    | 若为'no'，则不保留源端文件                                  | 
| checksum_algo | CHECKSUM_ALGO  | 若不为空，则后续消息头中产生当前文件的校验和                    | 
| bwlimit       | BWLIMIT        | 若为RSYNC_OVER_SSH/RSYNC，则限制传输的最大带宽               | 
|               | RSYNCD_MODULE  | 以原生rsync传输时，server端模块名。                          | 
|               | RSYNC_PASSWORD |                                                          | 
|               | ZSTD_CLEVEL    | 若不为空，则选择传输过程中启用zstd实时压缩，压缩等级为ZSTD_CLEVEL | 

## 三、主要参数格式说明

### 3.1 URL格式

SOURCE_URL/TARGET_URL

| MODE           |  实例                             | 描述                        |
| -------------- | -------------------------------- | ----------------------------|
| LOCAL          | /path/to                         | 本地目录路径                  |
| SSH            | user@ssh-host[:ssh-port]/path/to | 基于ssh的目录路径             | 
| RSYNC_OVER_SSH | user@ssh-host[:ssh-port]/path/to | 基于ssh的目录路径，用rsync传输 | 
| RSYNC          | rsync://user@rsync-host[:rsync-port]/path/to | rsync的路径      | 

- SOURCE_MODE/TARGET_MODE的设置
  - 取值为：LOCAL/SSH/RSYNC_OVER_SSH/RSYNC
  - SOURCE_MODE/TARGET_MODE一般可通过对应的URL确定。
  - RSYNC_OVER_SSH/SSH的URL格式相同，缺省为RSYNC_OVER_SSH，若需用SSH传输，则需显式设定SOURCE_MODE/TARGET_MODE。

### 3.2 JUMP格式

- SOURCE_JUMP/TARGET_JUMP

基于ssh的跳板设置，以','分隔多项，单项格式为user@host:port

### 3.3 CHECKSUM_ALGO的设置

取值为：'md5'/'sha1'/'sha256'/'sha512'，最终结果为16进制文本。


## 四、模块用法

### 4.1 ssh-server to local

```sh

echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.101/etc \
TARGET_URL=/tmp \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.101/etc \
TARGET_URL=/tmp/etc \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'master.cf' | \
SOURCE_URL=root@10.0.6.101/etc/postfix \
TARGET_URL=/tmp \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'master.cf' | \
SOURCE_MODE=SSH \
SOURCE_URL=root@10.0.6.101/etc/postfix \
TARGET_URL=/tmp \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```

```sh
echo 'master.cf' | \
SOURCE_URL=root@10.0.6.101/tmp \
TARGET_URL=/tmp \
KEEP_SOURCE=no \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.101/tmp \
TARGET_URL=/tmp \
KEEP_SOURCE=no \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```

```sh
echo '1115381072/1115381072_1115382688_combined.tar' | \
SOURCE_URL=scalebox@10.255.128.1:10022/raid0/mwa \
TARGET_URL=/tmp \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo '1115381072/1115381072_1115382688_combined.tar' | \
SOURCE_URL=scalebox@10.255.128.1:10022/raid0/mwa \
TARGET_URL=/tmp SOURCE_MODE=SSH \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```

### 4.2 local to ssh-server
```sh
echo 'postfix/master.cf' | \
SOURCE_URL=/etc \
TARGET_URL=root@10.0.6.101/tmp \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'postfix/master.cf' | \
SOURCE_URL=/etc \
TARGET_URL=root@10.0.6.101/tmp \
TARGET_MODE=SSH \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

```

### 4.3 ssh-server to ssh-server
```sh
echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.102/etc \
TARGET_URL=root@10.0.6.101/tmp/etc \
scalebox run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```

## 五、输出文件

模块运行后会在 `WORK_DIR` 下生成以下文件追踪记录：

### 5.1 input-files.txt

记录从**本地磁盘读取**的源文件，每行一个文件路径。以下传输模式会写入该文件：

| 传输模式 | 说明 |
| -------- | ---- |
| LOCAL → SSH | 本地文件经 SSH 发往远程 |
| LOCAL → RSYNC_OVER_SSH | 本地文件经 rsync over SSH 发往远程 |
| LOCAL → RSYNC | 本地文件经原生 rsync 发往远程 |

### 5.2 output-files.txt

记录**写入本地磁盘**的目标文件，每行一个文件路径。以下传输模式会写入该文件：

| 传输模式 | 说明 |
| -------- | ---- |
| SSH → LOCAL | 从远程经 SSH 拉取文件到本地 |
| RSYNC_OVER_SSH → LOCAL | 从远程经 rsync over SSH 拉取文件到本地 |
| RSYNC → LOCAL | 从远程经原生 rsync 拉取文件到本地 |

### 5.3 network-files.txt

记录**经网络传输**的文件，每行格式为 `direction,remote-ip,filename`，其中 `direction` 取值为 `from`（从网络读取）或 `to`（写出到网络）。所有非 LOCAL→LOCAL 的传输模式均会写入该文件：

| 传输模式 | 记录内容 |
| -------- | -------- |
| LOCAL → SSH | `to,<target_host>,<file>` |
| LOCAL → RSYNC_OVER_SSH | `to,<target_host>,<file>` |
| LOCAL → RSYNC | `to,<target_host>,<file>` |
| SSH → LOCAL | `from,<source_host>,<file>` |
| SSH → SSH | `from,<source_host>,<file>` + `to,<target_host>,<file>` |
| RSYNC_OVER_SSH → LOCAL | `from,<source_host>,<file>` |
| RSYNC → LOCAL | `from,<source_host>,<file>` |

> `remote-ip` 为远端主机名（不含用户名），由 `get_host_from_url` 从 `source_url` / `target_url` 中提取。
