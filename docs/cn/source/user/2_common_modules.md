# 2. 公共模块

## 2.1 文件传输相关模块

- uri定义

- 通用环境变量（用表格）

- 主要概念解释

### 2.1.1 file-copy模块

### 2.1.2 dir-copy模块

### 2.1.3 dir-list模块

## 2.2 非文件传输模块

### 2.2.1 cron模块

### 2.2.2 cluster-head模块


### 2.2.3 node-agent模块

#### 本地命令执行

- task-body: base64编码的待执行命令
- task-headers:
|   任务头    |       缺省值     |               说明                       |
| ---------- | --------------- | --------------------------------------- |
| action     | "REMOTE-EXEC"   |                                         |
| need_response | "no"         |     "yes" / "no"，缺省值为"no"，是否需要返回结果 |

#### 文件传输/目录传输

通过rsync-over-ssh，将本地目录/文件传输到远端。

- task-body: 待传输到远端的本地相对路径

- task-headers: 

|   任务头    |       缺省值     |               说明                       |
| ---------- | --------------- | --------------------------------------- |
| action     | "FILE-TRANSFER" |                                         |
| source_url |                 | "/local/dir"                            |
| target_url |                 | "user@remote-ip:remote-port/remote/dir" |
| keep_source | "yes"          | "yes" / "no"，是否保留源端（本地）文件/目录  |

