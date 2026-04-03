# 2. 核心模块详解

## 2.1 基础模块（agent）

## 2.2 文件传输相关模块

- uri定义
- 通用环境变量（用表格）
- 主要概念解释

通过rsync-over-ssh，将本地目录/文件传输到远端。

- task-body: 待传输到远端的本地相对路径

- task-headers: 

|   任务头    |       缺省值     |               说明                       |
| ---------- | --------------- | --------------------------------------- |
| action     | "FILE-TRANSFER" |                                         |
| source_url |                 | "/local/dir"                            |
| target_url |                 | "user@remote-ip:remote-port/remote/dir" |
| keep_source | "yes"          | "yes" / "no"，是否保留源端（本地）文件/目录  |


### 2.2.1 file-copy模块

- 模块配置参数与环境变量

### 2.2.2 dir-copy模块

### 2.2.3 dir-list模块





## 2.3 辅助功能模块

### 2.3.1 cron模块

### 2.3.2 cluster-head模块


### 2.3.3 node-agent模块



