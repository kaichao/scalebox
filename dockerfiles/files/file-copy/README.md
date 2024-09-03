# file-copy

## 一、模块介绍
支持以SSH/RSYNC_OVER_SSH/RSYNC等多种方式，完成跨节点、跨集群的文件传输。

## 二、 参数/环境变量

| 参数名         | 环境变量名          | 说明                                 |
| ------------  | ----------------- | ------------------------------------------------------- |
|  source_url   | SOURCE_URL   | 源端URL前缀值                   | 
|  target_url   | TARGET_URL   | 目标端URL前缀值                   | 
|               | SOURCE_MODE  | 源端URL的模式，取值为：LOCAL/SSH/RSYNC_OVER_SSH/RSYNC，通常仅针对RSYNC_OVER_SSH设置  | 
|               | TARGET_MODE  | 目源端URL的模式，取值为：LOCAL/SSH/RSYNC_OVER_SSH/RSYNC，通常仅针对RSYNC_OVER_SSH设置 | 
| source_jump_servers | SOURCE_JUMP_SERVERS | 源端跳板机，以','分隔多项，单项格式为user@host:port   | 
| target_jump_servers | TARGET_JUMP_SERVERS | 目标端跳板机，以','分隔多项，单项格式为user@host:port | 
|               | KEEP_SOURCE_FILE | 若为'no'，则不保留源端文件 | 
| checksum_algo | CHECKSUM_ALGO | 若不为空，则检验目标端校验和 | 
| rsync_bwlimit | RSYNC_BWLIMIT | 若为RSYNC_OVER_SSH/RSYNC，则限制传输的最大带宽 | 
|               | RSYNCD_MODULE    | 以原生rsync传输时，server端模块名。若server端以docker方式启动，选择'/local'；否则为缺省值'/root' | 
|               | RSYNC_PASSWORD   |           | 
|               | ZSTD_CLEVEL      | 若不为空，则选择传输过程中启用zstd实时压缩，压缩等级为ZSTD_CLEVEL | 

## 三、URL格式

SOURCE_URL/TARGET_URL

- LOCAL：
- SSH：
- RSYNC_OVER_SSH(RSYNC1)：
- RSYNC(RSYNC2)：

SOURCE_MODE/TARGET_MODE的取值

## 3. JUMP_SERVERS格式

SOURCE_JUMP_SERVERS/TARGET_JUMP_SERVERS

## 4. KEEP_SOURCE_FILE的设置

|        | LOCAL | SSH | RSYNC1 | RSYNC2 |
| ------ | ----- | --- | ------ | ------ |
| LOCAL  |   X   | yes |  yes   |  yes   |
|  SSH   |  yes  | yes |   X    |   X    |
| RSYNC1 |  yes  |  X  |   no   |   X    |
| RSYNC2 |  yes  |  X  |   X    |   no   |


上表中，左侧列表示源端，首行表示目的端。
- 值为'yes'，表示可设置KEEP_SOURCE_FILE
- 值为'no'，表示不可设置KEEP_SOURCE_FILE
- 值为'X'，表示当前模块不支持该操作
