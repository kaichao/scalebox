# 1. 命令行工具用法

命令行工具scalebox

## 1.1 命令行选项

| 选项               | 缺省值          | 描述                              |
| ----------------- | -------------- | --------------------------------- |
| -e / --env-file   | scalebox.env   | 环境变量文件，设置命令运行的环境变量。 |
| --debug           | 'no'           | 设置调试标志位，输出更多调试、排错的信息 |

环境变量是Scalebox应用程序中参数传递的重要方法。应用中环境变量定义可来自于多个环境变量定义文件、系统级环境变量，若在不同定义文件、系统级变量中存在重复的变量名，则按照以下顺序加载（若文件不存在，则忽略）：

- 系统级环境变量
- 用户自定义名的env文件
- 当前目录下scalebox.env文件
- ${HOME}/.scalebox/environments
- /etc/scalebox/environments

其中，用户自定义名的env文件，可按文件名，执行级联加载。

示例如下：
用户自定义env文件名为：p419_48nodes_1266932744.env，则按优先级从高到低，依次加载文件：
- p419_48nodes_1266932744.env
- p419_48nodes.env
- p419.env


## 1.2 子命令图

```{mermaid}

graph LR
  scalebox --> app[<a href="#app">app</a>]
  app --> app-create[<a href="#app-create">create</a>]
  app --> get-message-router[<a href="#app-get-message-router">get-message-router</a>]
  app --> app-list[<a href="#app-list">list</a>]
  app --> app-add-remote[<a href="#app-add-remote">add-remote</a>]
  app --> app-set-finished[<a href="#app-set-finished">set-finished</a>]

  scalebox --> job[<a href="#job">job</a>]
  job --> job-list[<a href="#job-list">list</a>]
  job --> job-info[info]

  scalebox --> task[<a href="#task">task</a>]
  task --> task-add[<a href="#task-add">add</a>]
  task --> task-get-header[<a href="#task-get-header">get-header</a>]
  task --> task-set-header[<a href="#task-set-header">set-header</a>]
  task --> task-remove-header[<a href="#task-remove-header">remove-header</a>]

  scalebox --> semaphore[<a href="#semaphore">semaphore</a>]
  semaphore --> sema-create[<a href="#semaphore-create">create</a>]
  semaphore --> semaphore-get[<a href="#semaphore-get">get</a>]
  semaphore --> increment[<a href="#semaphore-increment">increment</a>]
  semaphore --> decrement[<a href="#semaphore-decrement">decrement</a>]
  semaphore --> increment-n[<a href="#semaphore-increment">increment-n</a>]
  semaphore --> semaphore-group-dist[<a href="#semaphore-group-dist">group-dist</a>]

  scalebox --> variable[<a href="#variable">variable</a>]
  variable --> variable-create[<a href="#variable-create">create</a>]
  variable --> variable-get[<a href="#variable-get">get</a>]
  variable --> variable-set[<a href="#variable-set">set</a>]

  scalebox --> event[<a href="#event">event</a>]
  event --> event-task-add[<a href="#event-task-add">task-add</a>]
  event --> event-slot-add[<a href="#event-slot-add">slot-add</a>]
  event --> event-misc-add[<a href="#event-misc-add">misc-add</a>]

  scalebox --> global[<a href="#global">global</a>]
  global --> global-get[<a href="#global-get">get</a>]
  global --> global-set[<a href="#global-set">set</a>]

  scalebox --> fs[<a href="#fs">fs</a>]
  fs --> fs-ls[<a href="#fs-ls">ls</a>]
  fs --> fs-stat[<a href="#fs-stat">stat</a>]

  scalebox --> cluster
  cluster --> get-parameter

  scalebox --> version

  scalebox --> help

```

## 1.3 <span id="app">app子命令</span>

### 1.3.1 app create

解析应用定义文件，并存到数据库中，完成应用创建。

用法：
```sh
scalebox app create
```

### 1.3.2 app list

列出所有应用的基本信息。

用法：
```sh
scalebox app list
```

### 1.3.3 app set-finished

设置应用已完成，修改其状态为'FINISHED'

用法：
```sh
scalebox app set-finished --job-id ${job_id}
```

### 1.3.4 app add-remote

### 1.3.5 app get-message-router
  
## 1.3 <span id="job">job子命令</span>

### job list

### job info

## 1.4 <span id="task">task子命令</span>

### 1.4.1 task add

#### 参数/环境变量
- app-id/APP_ID
- job-id/JOB_ID: job-id
- sink-job/SINK_JOB: sink job-name
- from-job: job-name
- remote-server: grpc server for remote cluster.
- task-file: multiple messages in file
- ignore-dupkey: add "repeative":"yes" to headers
- headers: headers in json
- header/h: add one header
- to-ip: add "to_ip" to headers
- to-host: add "to_host" to headers
- upsert: overwrite existed task
- async-task-creation
- disable-local-ip

key-text可放在文件 ```${WORK_DIR}/task-body.txt```，该文件为多行文本，每行为一个消息体。

### 1.4.2 task get-header

获取task头信息。

示例：
```sh
scalebox task get-header --task-id 123 from_job
```
### 1.4.3 task set-header

设置新的header（若不存在）或覆盖已有header（若存在）。

示例：
```sh
scalebox task set-header --task-id 123 my_header value
```

### 1.4.4 task remove-header

移除已有header。若不存在，则在stderr上打印"my_header not-exists"

示例：
```sh
scalebox task remove-header --task-id 123 my_header
```


## 1.5 <span id="semaphore">semaphore子命令</span>

- 公共参数：job-id，或app-id
- 环境变量：JOB_ID，或APP_ID

- 信号量命名规则：
  - 字符集：大小写英文字母[A-Za-z]、数字[0-9]、冒号 :、下划线 _ 、中划线 -
  - 首字符为字母、下划线

### 1.5.1 semaphore create

#### 单个信号量的创建
示例：
```sh
scalebox semaphore create ${sema_name} ${int_value}
scalebox semaphore create --app-id ${app_id} ${sema_name} ${int_value}
APP_ID=${app_id} scalebox semaphore create ${sema_name} ${int_value}

scalebox semaphore create --job-id ${job_id} ${sema_name} ${int_value}
JOB_ID=${job_id} scalebox semaphore create ${sema_name} ${int_value}

```

### 信号量组的批量创建
- 命令行方式：受到bash的命令行最大长度2MiB限制。
```sh
scalebox semaphore create '{"semaphores":{"sema1":n1,"sema2":n2}}'
```

- 信号量文件方式：信号量数量通常可以更多
```sh
scalebox semaphore create --sema-file my-sema-file.txt
```

信号量文件为多行文件格式，每行表示一个信号量。

```
"sema1":n1
"sema2":n2
"sema3":n3
```


### 1.5.2 semaphore get

#### 获取单个信号量当前值
```sh
val=$(scalebox semaphore get ${sema_name})
code=$?
```
- ```code```为操作成功与否的标志。
  - 0：OK
  - 1：db error
  - 2： semaphore not-found
- ```val```为新的信号量值（整数）

####  获取信号量组的键值对
信号量组支持以SQL通配符%做通用匹配。

```sh
val=$(scalebox semaphore get ${sema_expr} )
code=$?
```

- sema_expr 为包含SQL通配符（%）的表达式
- ```code```为操作成功与否的标志。0为成功
- ```val```为新的信号量值，如果为多个信号量，返回结果为逗号分隔的信号量名值对。
  ``` ${sema1}:${val1},${sema2}:${val2},${sema3}:${val3} ```

### 1.5.3 semaphore increment

####  单个信号量的增一操作。
```sh
val=$(scalebox semaphore increment ${sema_name})
code=$?
```
- ```code```为操作成功与否的标志。
  - 0：OK
  - 1：db error
  - 2： semaphore not-found
- ```val```为新的信号量值（整数）

#### 信号量组的增一操作。

信号量组支持以SQL通配符%做匹配。

```sh
val=$(scalebox semaphore increment ${sema_expr} )
code=$?
```

- sema_expr 为包含SQL统配符（%）的表达式
- ```code```为操作成功与否的标志。0为成功
- ```val```为新的信号量值，如果为多个信号量，返回结果为逗号分隔的信号量名值对。
  ``` ${sema1}:${val1},${sema2}:${val2},${sema3}:${val3} ```


### 1.5.4 semaphore decrement

#### 单个信号量的减一操作。

```sh
val=$(scalebox semaphore decrement ${sema_expr} )
code=$?
```

用法详见：<a href="#semaphore-increment">semaphore increment</a>

#### 信号量组的减一操作。

用法详见：<a href="#semaphore-increment">semaphore increment</a>

### 1.5.5 semaphore increment-n

#### 单个信号量的加n操作。
```sh
val=$(scalebox semaphore increment-n ${sema_name}) ${delta_value}
code=$?
```

用法详见：<a href="#semaphore-increment">semaphore increment</a>

#### 信号量组的加n操作。

### 1.5.6 semaphore group-dist

- 信号量格式：``` progress-counter_${mod_name}_${host_name} ```，并且对应主机的group_id不为空。

示例：
```sh
APP_ID=3 scalebox semaphore group-dist progress-counter_pull-unpack_r04.main
```

## 1.6 <span id="variable">variable子命令</span>

- 公共参数：job-id，或app-id
- 环境变量：JOB_ID，或APP_ID

- 变量名命名：同信号量命名

### 1.6.1 variable create

```sh
scalebox variable create --app-id ${app_id} var_name ${str_value}
APP_ID=${app_id} scalebox variable create var_name ${str_value}

scalebox variable create --job-id ${job_id} var_name ${str_value}
JOB_ID=${job_id} scalebox variable create  var_name ${str_value}
```

### 1.6.2 variable get

示例：
```sh
export APP_ID=${app_id}
var_val=$(scalebox variable get ${var_name})
code=$?
[[ $code -ne 0 ]] && echo "[ERROR] variable-get ${var_name}, exit_code:$code" >&2

```
### 1.6.3 variable set

示例：
```sh
scalebox variable set ${var_name} ${str_value}
```

## 1.7 <span id="event">event子命令</span>

支持各类event的add操作。

### 基本命令： 

- xxxx为："task"/"slot"/"misc"

```sh
scalebox event xxxx-add "${tag_name}" "${level_name}" ["${code}" ["${txt}" ["${json}"]]]
```
- 若code为空，则取值0
- 若txt为空，则取值""
- 若json为空，则取值"{}"

```sh
scalebox event xxxx-add --txt-file "${txt_file}" --json-file "${json_file}" "${tag_name}" "${level_name}" ["${code}"]
```

则txt、json从文件中读取。


### 1.7.1 event task-add

通过环境变量TASK_ID或参数 --task-id  指定task-id。
```sh
scalebox event task-add --task-id ${task_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

```scalebox event task-add ``` 可简写为 ``` scalebox event add  ```

### 1.7.2 event slot-add

通过环境变量SLOT_ID或参数 --slot-id  指定slot-id。
```sh
scalebox event slot-add --slot-id ${slot_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

### 1.7.3 event misc-add

```sh
scalebox event misc-add ${tag_name} ${level_name} ${code} ${txt} ${json}
```

## 1.8 <span id="global">global子命令</span>

全局变量

### 1.8.1 global get

```sh
scalebox global get ${global_name}
```

### 1.8.2 global set

```sh
scalebox global set ${global_name} ${global_value}
```

## 1.9 <span id="fs">fs子命令</span>

scalebox-fs以文件系统形式，将分布式计算节点上的文件组织在同一个名字空间中。后期可提供mount支持、跨节点迁移等特性。

### 1.8.1 fs ls

- 主要参数：
  - include-removed-file
  - with-hostname
  - with-file-size
  - hostname=${host-name}

```sh
scalebox fs ls ${path_expr}
```

### 1.8.2 fs stat

查看1个或多个文件的元数据。每个节点上的文件名跟全局文件名一致。

元数据主要包括：
- 虚拟文件名
- 文件所属主机号
- 创建时间
- 删除时间
- 
```sh
scalebox fs stat ${path_expr}
```
