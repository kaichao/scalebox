# 1. 命令行scalebox用法

命令行工具scalebox

## 1.1 命令行选项

| 选项               | 缺省值          | 描述                              |
| ----------------- | -------------- | --------------------------------- |
| -e / --env-file   | scalebox.env   | 环境变量文件，设置命令运行的环境变量。 |
| --debug           | 'no'           | 设置调试标志位，输出更多调试、排错的信息 |

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

  scalebox --> semaphore[<a href="#semaphore">semaphore</a>]
  semaphore --> sema-create[<a href="#semaphore-create">create</a>]
  semaphore --> count-down[<a href="#semaphore-count-down">count-down</a>]
  semaphore --> semaphore-get[<a href="#semaphore-get">get</a>]
  semaphore --> semaphore-group-dist[<a href="#semaphore-group-dist">group-dist</a>]

  scalebox --> cluster
  cluster --> get-parameter

  scalebox --> config
  config --> config-get[get]
  config --> config-set[set]  

  scalebox --> version

  scalebox --> help

```


## 1.3 scalebox app 子命令 {#app}

### 1.3.1 app create{#app-create}

### 1.3.2 app list {#app-list}

### 1.3.3 app add-remote {#app-add-remote}

### 1.3.4 app set-finished {#app-set-finished}

### 1.3.5 app get-message-router {#app-get-message-router}
  
## 1.3 scalebox job 子命令 {#job}

### job list {#job-list}

### job info

## 1.4 scalebox task 子命令{#task}

### task add {#task-add}

key-text可放在文件 ```${WORK_DIR}/key-text.txt```，该文件为多行文本，每行为一个消息体。


## 1.5 scalebox semaphore 子命令{#semaphore}

### 1.5.1 semaphore create{#semaphore-create}


### 1.5.2 semaphore count-down{#semaphore-count-down}


### 1.5.3 semaphore get{#semaphore-get}


### 1.5.4 semaphore group-dist{#semaphore-group-dist}

