# 2. 核心模块详解

-  核心模块的可编程特性

- 其功能脚本可放在/app/share/bin下，子模块的功能脚本在/app/bin下。模块识别规范是优先使用/app/bin，再搜索/app/share/bin目录下；
- main-router任务体格式定制：按标准模块的任务体格式定制，便于使用标准模块功能；

这样可充分利用标准模块的功能。

- 算法模块：
- 辅助模块：
  - 数据传输模块
    - 基于ssh的数据同步拷贝
    - 基于其他协议的数据传输（wdt/rcloud/gridftp/...）
  - 其他辅助模块：目录列表
- 路由模块

## 2.1 基础模块（agent）

- 提供了模块底座功能；
  - 应用模块的中任务执行的生命周期管理
  - 与server端runtime的controld交互，纪录任务执行状态

### 2.1.5 非agent模块的构建


## 2.2 文件传输相关模块

- uri定义
- 通用环境变量（用表格）
- 主要概念解释

通过rsync-over-ssh，将本地目录/文件传输到远端。

- task-body: 待传输到远端的本地相对路径

- task-headers: 

|   任务头    |       缺省值     |               说明                       |
| ---------- | --------------- | --------------------------------------- |
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



