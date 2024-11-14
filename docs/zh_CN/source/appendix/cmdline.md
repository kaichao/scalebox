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


## 1.2 子命令

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
  task --> task-info[<a href="#task-info">info</a>]

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

  scalebox --> cluster
  cluster --> get-parameter

  scalebox --> config
  config --> config-get[get]
  config --> config-set[set]  

  scalebox --> version

  scalebox --> help

```


## 1.3 <span id="app">app 子命令</span>

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
  
## 1.3 job 子命令 {#job}

### job list

### job info

## 1.4 scalebox task 子命令{#task}

### 1.4.1 task add

#### 环境变量/参数
- APP_ID/app-id
- JOB_ID/job-id
- SINK_JOB/sink-job
- 
- 
key-text可放在文件 ```${WORK_DIR}/task-body.txt```，该文件为多行文本，每行为一个消息体。

### 1.4.2 task info


## 1.5 <span id="semaphore">semaphore子命令</span>

- 公共参数：job-id，或app-id
- 环境变量：JOB_ID，或APP_ID

### 1.5.1 semaphore create

示例：
```sh
scalebox semaphore create sema_name ${int_value}
scalebox semaphore create --app-id ${app_id} ${sema_name} ${int_value}
APP_ID=${app_id} scalebox semaphore create ${sema_name} ${int_value}

scalebox semaphore create --job-id ${job_id} ${sema_name} ${int_value}
JOB_ID=${job_id} scalebox semaphore create ${sema_name} ${int_value}

```

### 1.5.2 semaphore get


### 1.5.3 semaphore increment

单个/单组（前缀匹配）信号量增一的操作。

### 1.5.4 semaphore decrement

单个/单组（前缀匹配）信号量减一的操作。

### 1.5.5 semaphore increment-n

单个/单组（前缀匹配）信号量加n的操作。

### 1.5.6 semaphore group-dist

- 信号量格式：信号量名:host全名，并且对应主机的group_id不为空。

示例：
```sh
APP_ID=3 scalebox semaphore group-dist progress-counter_pull-unpack:r04.main
```

## 1.6 variable子命令{#variable}

- 公共参数：job-id，或app-id
- 环境变量：JOB_ID，或APP_ID

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
