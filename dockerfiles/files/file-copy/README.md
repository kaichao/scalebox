# file-copy

## 一、模块介绍

支持以SSH/RSYNC_OVER_SSH/RSYNC等多种方式，完成跨节点、跨集群的文件传输。

## 二、 参数/环境变量

| 参数名         | 环境变量名         | 说明                                                    |
| ------------  | ---------------- | ------------------------------------------------------- |
|  source_url   | SOURCE_URL       | 源端URL前缀值                                            | 
|  target_url   | TARGET_URL       | 目标端URL前缀值                                          | 
|               | SOURCE_MODE      | 源端URL模式                                             | 
|               | TARGET_MODE      | 目源端URL模式                                            | 
| source_jump   | SOURCE_JUMP      | 源端跳板                                                | 
| target_jump   | TARGET_JUMP      | 目标端跳板                                               | 
|               | KEEP_SOURCE_FILE | 若为'no'，则不保留源端文件                                 | 
| checksum_algo | CHECKSUM_ALGO    | 若不为空，则后续消息头中产生当前文件的校验和                   | 
| bwlimit       | BWLIMIT          | 若为RSYNC_OVER_SSH/RSYNC，则限制传输的最大带宽              | 
|               | RSYNCD_MODULE    | 以原生rsync传输时，server端模块名。                         | 
|               | RSYNC_PASSWORD   |                                                         | 
|               | ZSTD_CLEVEL      | 若不为空，则选择传输过程中启用zstd实时压缩，压缩等级为ZSTD_CLEVEL | 

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
  - SOURCE_MODE/TARGET_MODE一般可通过对应的URL确定。RSYNC_OVER_SSH/SSH的URL格式相同，缺省为RSYNC_OVER_SSH，若需用SSH传输，则需显式设定SOURCE_MODE/TARGET_MODE。

### 3.2 JUMP格式

- SOURCE_JUMP/TARGET_JUMP

基于ssh的跳板设置，以','分隔多项，单项格式为user@host:port

### 3.3 CHECKSUM_ALGO的设置

取值为：'md5'/'sha1'/'sha256'/'sha512'，最终结果为16进制文本。


## 四、模块用法

### 4.1 ssh-server to local

```sh

echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.101/etc TARGET_URL=/tmp \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.101/etc TARGET_URL=/tmp/etc \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'master.cf' | \
SOURCE_URL=root@10.0.6.101/etc/postfix TARGET_URL=/tmp \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'master.cf' | \
SOURCE_MODE=SSH SOURCE_URL=root@10.0.6.101/etc/postfix TARGET_URL=/tmp \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```


```sh
echo 'master.cf' | \
SOURCE_URL=root@10.0.6.101/tmp TARGET_URL=/tmp KEEP_SOURCE_FILE=no \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.101/tmp TARGET_URL=/tmp KEEP_SOURCE_FILE=no \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```

```sh
echo '1115381072/1115381072_1115382688_combined.tar' | \
SOURCE_URL=scalebox@10.255.128.1:10022/raid0/mwa TARGET_URL=/tmp \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo '1115381072/1115381072_1115382688_combined.tar' | \
SOURCE_URL=scalebox@10.255.128.1:10022/raid0/mwa TARGET_URL=/tmp SOURCE_MODE=SSH \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```

### 4.2 local to ssh-server
```sh
echo 'postfix/master.cf' | \
SOURCE_URL=/etc TARGET_URL=root@10.0.6.101/tmp \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

echo 'postfix/master.cf' | \
SOURCE_URL=/etc TARGET_URL=root@10.0.6.101/tmp TARGET_MODE=SSH \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest

```

### 4.3 ssh-server to ssh-server
```sh
echo 'postfix/master.cf' | \
SOURCE_URL=root@10.0.6.102/etc TARGET_URL=root@10.0.6.101/tmp/etc \
scalebox app run --image-name hub.cstcloud.cn/scalebox/file-copy:latest
```
