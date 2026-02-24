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
  scalebox --> run[<a href="#run">run</a>]

  scalebox --> cluster[<a href="#cluster">cluster</a>]
  cluster --> cluster-get-parameter[<a href="#cluster-get-parameter">get-parameter</a>]
  cluster --> cluster-check-status[<a href="#cluster-check-status">check-status</a>]
  cluster --> cluster-dist-image[<a href="#cluster-dist-image">dist-image</a>]
  cluster --> cluster-app-view[<a href="#cluster-app-view">app-view</a>]
  cluster --> cluster-host-view[<a href="#cluster-host-view">host-view</a>]

  scalebox --> host[<a href="#host">host</a>]
  host --> host-check-status[<a href="#host-check-status">check-status</a>]
  host --> host-get-info[<a href="#host-get-info">get-info</a>] 
  host --> host-add-node[<a href="#host-add-node">add-node</a>]
  host --> host-dist-image[<a href="#host-dist-image">dist-image</a>]
  host --> host-recover[<a href="#host-recover">recover</a>]
  host --> host-migrate[<a href="#host-migrate">migrate</a>]
  host --> host-replace[<a href="#host-replace">replace</a>]
  host --> host-asign[<a href="#host-asign">asign</a>]
  
  scalebox --> slot[<a href="#slot">slot</a>]
  slot --> slot-add[<a href="#slot-add">add</a>]
  slot --> slot-add-group[<a href="#slot-add-group">add-group</a>]
  slot --> slot-update[<a href="#slot-update">update</a>]

  scalebox --> app[<a href="#app">app</a>]
  app --> app-create[<a href="#app-create">create</a>]
  app --> app-run[<a href="#app-run">run</a>]
  app --> main-router[<a href="#app-main-router">main-router</a>]
  app --> app-list[<a href="#app-list">list</a>]
  app --> app-add-remote[<a href="#app-add-remote">add-remote</a>]
  app --> app-set-finished[<a href="#app-set-finished">set-finished</a>]

  scalebox --> module[<a href="#module">module</a>]
  module --> module-list[<a href="#module-list">list</a>]
  module --> module-info[<a href="#module-info">info</a>]

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
  semaphore --> increment-n[<a href="#semaphore-increment-n">increment-n</a>]
  semaphore --> semaphore-group[<a href="#semaphore-group">group</a>]

  scalebox --> semagroup[<a href="#semagroup">semagroup</a>]
  semagroup --> semagroup-min[<a href="#semagroup-min">min</a>]
  semagroup --> semagroup-max[<a href="#semagroup-max">max</a>]
  semagroup --> semagroup-increment[<a href="#semagroup-increment">increment</a>]
  semagroup --> semagroup-decrement[<a href="#semagroup-decrement">decrement</a>]
  semagroup --> semagroup-diffmin[<a href="#semagroup-diffmin">diffmin</a>]
  semagroup --> semagroup-diffmax[<a href="#semagroup-diffmax">diffmax</a>]

  scalebox --> variable[<a href="#variable">variable</a>]
  variable --> variable-get[<a href="#variable-get">get</a>]
  variable --> variable-set[<a href="#variable-set">set</a>]

  scalebox --> global[<a href="#global">global</a>]
  global --> global-get[<a href="#global-get">get</a>]
  global --> global-set[<a href="#global-set">set</a>]

  scalebox --> channel[<a href="#channel">channel</a>]
  channel --> channel-pull[<a href="#channel-pull">pull</a>]
  channel --> channel-push[<a href="#channel-push">push</a>]

  scalebox --> event[<a href="#event">event</a>]
  event --> event-task-add[<a href="#event-task-add">task-add</a>]
  event --> event-slot-add[<a href="#event-slot-add">slot-add</a>]
  event --> event-misc-add[<a href="#event-misc-add">misc-add</a>]

  scalebox --> fs[<a href="#fs">fs</a>]
  fs --> fs-ls[<a href="#fs-ls">ls</a>]
  fs --> fs-stat[<a href="#fs-stat">stat</a>]

  scalebox --> status

  scalebox --> help

```

## 1.3 <span id="cluster">cluster子命令</span>

### 1.3.1 cluster get-parameter

### 1.3.2 cluster check-status

- 检测cluster配置
- 检测所有host的状态

### 1.3.3 cluster dist-image

### 1.3.4 cluster app-view

### 1.3.5 cluster host-view


## 1.4 <span id="host">host子命令</span>

用于管理系统中的计算节点，包括分组、向流水线应用增添新的节点，替换、释放已有计算节点等功能。

### 1.4.1 host check-status

检查一个指定节点的运行状态，输出其ip_addr、uname、port、parameters、group_id、reg_time、status、last_active、comments等信息

### 1.4.2 host get-info

由hostname获取host信息

### 1.4.2 host add-node

将一个或一组节点分配给一个新的应用。
- 检索目标应用所有非slot-on-head的任务
- 如有需要，调用 ```host dist-image```，将所需镜像分发到指定节点
- 对所有```host_vtask_size```属性的module，创建对应的信号量
- 调用```slot add-group```，在节点上创建对应数量的slot

### 1.4.3 host dist-image

将docker镜像分发到指定节点或节点组

### 1.4.4 host recover

动态申请的计算节点发生故障时，替换用新节点替换。

- 设置原节点所有相关的task状态为-1为-4（确定没有-2、-3状态的task）：down机后，task状态不再变化？
- 将修改后原节点的所有对应host-bound任务的SLOT状态设置为OFF；
- 调用host replace替换节点
- 针对存在HOST-BOUND/SLOT-BOUND任务的节点，启动新节点加载数据模块，恢复新节点的本地数据；
- 将修改后原节点的所有对应host-bound任务的SLOT状态设置为READY；
- 将状态为-4的task设置为-1

### 1.4.5 host migrate

按计划替换即将到期的计算节点。

- 在```*_vtask_size```相关属性的module上，通过减少对应信号量，使得该node上task逐步退出；
- 第一步完成后，调用```host replace```替换节点；(需要定时运行？或者系统检测task全部完成后运行？)
- 增加```*_vtask_size```对应信号量，恢复新节点上的slot/task的运行；

### 1.4.6 host replace

执行节点替换任务

- 原节点中除hostname外的属性（ip_addr、uname、port、parameters、reg_time、status、last_active、comments）替换为新节点对应属性
- 修改新节点的node-agent的host属性，指向原节点的hostname；
- 删除host表的新节点纪录

### 1.4.7 host asign

将已有节点重新分组。
需要参数：src_group_id, dst_group_id, num_groups, group_size
- 从一组或多组源节点中，筛选出满足条件的节点
- 对满足条件的节点，根据输入取出若干个，修改其group_id为dst_group_id，同时将节点重命名为${dst_group_id}-${group_idx}${hostidx}的格式。
- 当满足条件的节点数量不足时，仅分配整组节点，输出分配节点的具体数量
- 在src_group_id包含dst_group_id，或指定了参数时，使用新节点填充dst_group_id中group_idx 0~num_groups，使每组均有group_size个节点。否则，group_idx从现有最大值开始递增，共尝试分配num_groups*group_size个节点。
- group_size为-1时，将所有节点加入dst_group_id,并从0开始重设group_idx。

## 1.5 <span id="slot">slot子命令</span>

### 1.5.1 slot add

### 1.5.2 slot add-group

根据给定的配置参数，在指定节点上为每个module创建对应数量的节点。

### 1.5.3 slot update


## 1.6 <span id="app">app子命令</span>

### 1.6.1 app create

解析应用定义文件，并存到数据库中，完成应用创建。

用法：
```sh
scalebox app create
```

### 1.6.2 app run

以命令行方式，启动scalebox应用。（未来代替app create ?）

- 环境变量文件：```./scalebox.env```
- 主模块代码目录：```./code/```
- 路由模块代码目录：```./mr-code/``` 

#### 单启动消息
```sh
export ENV0=v0
export ENV1=v1

scalebox app run --param-name=param-value start-item
```

参数表
| 参数名         |    参数说明     |  对应环境变量    | 缺省值                                       |
| ------------- | -------------- | -------------- | ------------------------------------------- |
| app-name       | app名称       | _APP_NAME       | 
| cluster       | cluster名      | _CLUSTER       | local                                       |
| image-name    | 主模块镜像名     | _IMAGE_NAME    | hub.cstcloud.cn/scalebox/agent:latest       |
| code-path     | 主模块代码目录   | _CODE_PATH     | 若当前目录下有./code/，则为./code;否则为空       |
| slot-regex    | 主模块的slot配置 | _SLOT_REGEX    | 缺省为：h0，在头节点上1个slot                  |
| mr-image-name | 路由模块镜像名   | _MR_IMAGE_NAME | 若mr_code_path已设置，则设置为agent            |
| mr-code-path  | 路由模块代码目录 | _MR_CODE_PATH  | 若当前目录下有./mr-code/，则为./mr-code;否则为空 |
| app-file/f    | 应用定义文件     |                | 若当前目录下有app.yaml，则为app.yaml，否则缺省为空 |
| env-file/e    | 环境变量文件     |                | 若当前目录下有app.yaml，则为app.yaml，否则缺省为空 |

- 若有应用定义文件，则以此创建应用
- 启动消息start-task
  - 若有消息路由，则启动消息发给消息路由
  - 若无消息路由，则启动消息发给首模块

- 启动项start-item
若非json串，则为启动消息start-task；否则start-item中包括前述参数及start-task。json格式定义如下：
```json
{
  "cluster": "my-cluster",
  "image_name": "my-image",
  "code_path": "/path/to/code",
  "slot_regex": "node[0-9]+:2",
  "mr_image_name": "my-mr-image",
  "mr_code_path": "/path/to/mr-code",
  "start_task": "starting task"
}
```
实际应用中，去除json字符串中的无空格、换行等空字符

#### 基于管道的多启动消息

针对多启动消息，可通过管道将多消息按行传递给启动命令。每行的消息体不按前述json格式解析。
```sh
echo "start-item\nstart-task1" | scalebox app run --param-name=param-value
```

示例：
```sh
# 设定源端、目标端URL
export SOURCE_URL=/data2/mydata/mwa/tar
export TARGET_URL=cstu0036@10.100.1.104:65010/work2/cstu0036/mydata/mwa/tar

# 单文件传输
scalebox app run --image-name=hub.cstcloud.cn/scalebox/file-copy:latest 1267459328/1267464090_1267464129_ch127.dat.tar.zst

# 多文件传输
cd /data2/mydata/mwa/tar
find 1267459328 -type f | scalebox app run --image-name=hub.cstcloud.cn/scalebox/file-copy:latest --slot-regex=h0:2
```

### 1.6.3 app main-router
  

### 1.6.4 app list

列出所有应用的基本信息。

用法：
```sh
scalebox app list
```

### 1.6.5 app set-finished

设置应用已完成，修改其状态为'FINISHED'

用法：
```sh
scalebox app set-finished --module-id ${module_id}
```

### 1.6.6 app add-remote

## 1.7 <span id="module">module子命令</span>

### 1.7.1 module list

### 1.7.2 module info

## 1.8 <span id="task">task子命令</span>

### 1.8.1 task add

#### 参数/环境变量

| 参数名         |  环境变量名         |    说明                                 |
| --------------- | --------------- | --------------------------------------- |
| app-id          | APP_ID          |                                         |
| module-id       | MODULE_ID       |                                         |
| sink-module     | SINK_MODULE     | sink-module name                        |
| conflict-action | CONFLICT_ACTION | 数据库插入时发生冲突的缺省动作，''/'IGNORE'/'OVERWRITE' |
| from-module     |                 | module-name                             |
| remote-server   |                 | grpc server for remote cluster, 格式为{ip_addr}:{port} |
| task-file       |                 | multiple tasks in file               |
| ignore-dupkey   |                 | add "repeative":"yes" to headers        |
| headers         |                 | headers in json                         |
| header/h        |                 | add one header                          |
| to-ip           |                 | add "to_ip" to headers (以-h to_ip=$ip_addr代替)   |
| to-host         |                 | add "to_host" to headers (以-h to_host=$host_name代替) |
| disable-local-ip |                |                                         |
| batch-size      |                 | 批量task添加中，指定批次大小。缺省值100。     |


task文件缺省为 ```${WORK_DIR}/sink-tasks.txt```，该文件为多行文本，每行为 消息体+消息头。

task文件每行格式如下：
| 类型                   |  示例                                        |
| --------------------- | ----------------------------------------------- |
| 文本body               | body                                            |
| json body             | {"hi0":"a","body":"my_body"}                     |
| 文本body+headers       | body,{"h0":"a","h1":"b"}                        |
| json-body+headers     | {"hi0":"a","body":"my_body"},{"h0":"a","h1":"b"} |
| 模块名+文本body         | module-name,body                                            |
| 模块名+文本body+headers | module-name,body,{"h0":"a","h1":"b"}                        |

用于控制的task头（header）：
| header              |  说明                                   |
| ------------------- | -------------------------------------- |
| initial_status_code | 缺省为-1,'READY'                        |
| upsert              | overwrite existed task    （删除）       |
| conflict-action     | ''/'IGNORE'/'OVERWRITE'                |
| async-task-creation |                                        |
| slot_broadcast      |                                        |
| host_broadcast      |                                        |

### 1.8.2 task get-header

获取task头信息。

示例：
```sh
scalebox task get-header --task-id 123 from_module
```
### 1.8.3 task set-header

设置新的header（若不存在）或覆盖已有header（若存在）。

示例：
```sh
scalebox task set-header --task-id 123 my_header value
```

### 1.8.4 task remove-header

移除已有header。若不存在，则在stderr上打印"my_header not-exists"

示例：
```sh
scalebox task remove-header --task-id 123 my_header
```

## 1.9 <span id="semaphore">semaphore子命令</span>

- 公共参数：module-id，或app-id
- 环境变量：MODULE_ID，或APP_ID、SEMAPHORE_AUTO_CREATE

- 信号量命名规则：
  - 字符集：大小写英文字母[A-Za-z]、数字[0-9]、冒号 :、下划线 _ 、中划线 -
  - 首字符为字母、下划线

- 信号量表达式：表示一组信号量的正则表达式
  - 字符集：信号量字符集，加上 '.*+?^$[]{}()|\'

### 1.9.1 semaphore create

- 参数：batch-size：用于批量信号量创建中，指定批次大小，缺省值为100。

#### 单个信号量的创建
示例：
```sh
scalebox semaphore create ${sema_name} ${int_value}
scalebox semaphore create --app-id ${app_id} ${sema_name} ${int_value}
APP_ID=${app_id} scalebox semaphore create ${sema_name} ${int_value}

scalebox semaphore create --module-id=${module_id} ${sema_name} ${int_value}
MODULE_ID=${module_id} scalebox semaphore create ${sema_name} ${int_value}

```

#### 信号量组的批量创建
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

### 1.9.2 semaphore get

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

若设置环境变量SEMAPHORE_AUTO_CREATE=yes，则自动创建初值为0的信号量

```sh
val=$(SEMAPHORE_AUTO_CREATE=yes scalebox semaphore get ${sema_name})
code=$?
```
- ```code```为操作成功与否的标志。
  - 0：OK
  - 1：db error
- ```val```为新的信号量值（整数）

####  获取信号量组的json键值对
信号量组支持变量名以正则表达式做通用匹配。

```sh
val=$(scalebox semaphore get ${sema_expr} )
code=$?
```

- sema_expr 为正则表达式
- ```code```为操作成功与否的标志。0为成功
- ```val```为新的信号量值，如果为多个信号量，返回结果为json map表示的信号量名值对。
  ```{"sema1":n1,"sema2":n2,"sema3":n3}```

### 1.9.3 semaphore increment

####  单个信号量的增一操作
```sh
val=$(scalebox semaphore increment ${sema_name})
code=$?
```

- ```code```为操作成功与否的标志。
  - 0：OK
  - 1：db error
  - 2： semaphore not-found
- ```val```为新的信号量值（整数）

若设置环境变量SEMAPHORE_AUTO_CREATE=yes，则自动创建初值为0的信号量，并增1.

```sh
val=$(SEMAPHORE_AUTO_CREATE=yes scalebox semaphore increment ${sema_name})
code=$?
```
- ```code```为操作成功与否的标志。
  - 0：OK
- ```val```为新的信号量值（整数）

#### 信号量组的增一操作

信号量组支持变量名以正则表达式做通用匹配。

```sh
val=$(scalebox semaphore increment ${sema_expr} )
code=$?
```

- sema_expr 为正则表达式
- ```code```为操作成功与否的标志。0为成功
- ```val```为新的信号量值，如果为多个信号量，返回结果为json map表示的信号量名值对。
  ```{"sema1":n1,"sema2":n2,"sema3":n3}```

### 1.9.4 semaphore decrement

#### 单个信号量的减一操作。

```sh
val=$(scalebox semaphore decrement ${sema_expr} )
code=$?
```

用法详见：<a href="#semaphore-increment">semaphore increment</a>

#### 信号量组的减一操作。

用法详见：<a href="#semaphore-increment">semaphore increment</a>

### 1.9.5 semaphore increment-n

#### 单个信号量的加n操作。
```sh
val=$(scalebox semaphore increment-n ${sema_name} ${delta_value})
code=$?
```

用法详见：<a href="#semaphore-increment">semaphore increment</a>

#### 信号量组的加n操作。

用法详见：<a href="#semaphore-increment">semaphore increment</a>


### 1.9.7 semaphore global-dist

(改为global-offset?)

- 作用范围：t_host表中group_id不为NULL的所有host

- 信号量格式：``` task_progress:${mod_name}:${host_name} ```，并且对应主机的group_id不为空。

示例：
```sh
APP_ID=3 scalebox semaphore global-dist task_progress:beam-make:r04.main
```

### 1.9.8 semaphore group-dist

(改为group-offset?)

- 作用范围：t_host表中group_id相同的host分为一组（为NULL的也是一组）

- 信号量格式：``` task_progress:${mod_name}:${host_name} ```，并且对应主机的group_id不为空。

示例：
```sh
APP_ID=3 scalebox semaphore group-dist task_progress:beam-make:r04.main
```

## 1.10 <span id="semagroup">semagroup子命令</span>

- 多个信号量组成信号量组，用信号量名前缀、正则表达式标识信号量组


#### 1.10.1 semagroup max
- 信号量组中最大值
```sh
val=$(scalebox semagroup max ${sema_expr})
code=$?
```

#### 1.10.2 semagroup min
- 信号量组中最小值
```sh
val=$(scalebox semagroup min ${sema_expr})
code=$?
```
- sema_expr为信号量名的正则表达式或前缀
- 返回值val为整数字符串

#### 1.10.3 semagroup increment
- 选取信号量组中最小值，并加一

#### 1.10.4 semagroup decrement
- 选取信号量组中最大值，并减一

#### 1.10.5 semagroup diffmax
- 信号量组最大值与信号量当前值的差值
```sh
val=$(scalebox semagroup diffmax ${sema_expr})
code=$?
```
- sema_expr为含分组定义的信号量，示例为```(group-prefix):sema-suffix```
- 返回值val为整数字符串

#### 1.10.6 semagroup diffmin
- 信号量当前值与信号量组最小值的差值
```sh
val=$(scalebox semagroup diffmin ${sema_expr})
code=$?
```
- sema_expr为含分组定义的信号量，示例为```(group-part:)sema-suffix```
- 返回值val为整数字符串


## 1.11 <span id="variable">variable子命令</span>

- 公共参数：module-id，或app-id
- 环境变量：MODULE_ID，或APP_ID

- 变量名命名：同信号量命名

### 1.11.1 variable get

#### 获取单个变量当前值
- 示例：
```sh
val=$(scalebox variable get ${var_name})
code=$?
[[ $code -ne 0 ]] && echo "[ERROR] variable-get ${var_name}, exit_code:$code" >&2
```
- ```code```为操作成功与否的标志。
  - 0：OK
  - 1：db error
  - 2： variable not-found
- ```val```为新的变量值

####  获取变量组的json键值对
变量组支持变量名以正则表达式做通用匹配。

```sh
val=$(scalebox variable get ${var_expr} )
code=$?
```

- var_expr 为正则表达式
- ```code```为操作成功与否的标志。0为成功
- ```val```为新的变量量值，返回结果为json map表示的信号量名值对。
  ```{"var1":"val1","var2":"val2","var3":"val3"}```

### 1.11.2 variable set

```sh
scalebox variable set --app-id ${app_id} ${var_name} ${str_value}
APP_ID=${app_id} scalebox variable set ${var_name} ${str_value}

scalebox variable set --module-id=${module_id} ${var_name} ${str_value}
MODULE_ID=${module_id} scalebox variable set ${var_name} ${str_value}
```

## 1.12 <span id="global">global子命令</span>

全局变量

### 1.12.1 global get

```sh
scalebox global get ${global_name}
```

### 1.12.2 global set

```sh
scalebox global set ${global_name} ${global_value}
```


## 1.13 <span id="channel">channel子命令</span>

channel用于跨应用间的通信，是一个有优先级队列。

- 公共参数：module-id，或app-id
- 环境变量：MODULE_ID，或APP_ID

- 优先队列命名：同信号量命名

### 1.13.1 channel create

- head-app : app-id
- tail-app

若不指定，则为当前app

### 1.13.2 channel pull

- 获取队列当前值
- 示例：
```sh
val=$(scalebox channel pull ${pp_name})
code=$?
[[ $code -ne 0 ]] && echo "[ERROR] channel-pull ${pp_name}, exit_code:$code" >&2
```
- ```code```为操作成功与否的标志。
  - 0：OK
  - 1：db error
  - 2： channel not-found
- ```val```为新的变量值

### 1.13.3 channel push

- priority为优先级，浮点数。数值小，优先级高。
  
```sh
scalebox channel push --app-id ${app_id} ${pp_name} ${str_value} [${priority}]
APP_ID=${app_id} scalebox channel push ${pp_name} ${str_value}

scalebox channel push --module-id ${module_id} ${pp_name} ${str_value} [${priority}]
MODULE_ID=${module_id} scalebox channel push ${pp_name} ${str_value}
```


## 1.14 <span id="fs">fs子命令</span>

scalebox-fs以文件系统形式，将分布式计算节点上的文件组织在同一个名字空间中。后期可提供mount支持、跨节点迁移等特性。

### 1.14.1 fs ls

- 主要参数：
  - include-removed-file
  - with-hostname
  - with-file-size
  - hostname=${host-name}

```sh
scalebox fs ls ${path_expr}
```

### 1.14.2 fs stat

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

## 1.15 <span id="status">status</span>

- 系统整体状态：local集群头节点（actuator到local头节点有效）
- cluster列表：不同状态的host数量
- app列表：不同状态

## 1.16 <span id="event">event子命令</span>

支持各类event的add操作。

- 基本命令： 

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


### 1.16.1 event task-add

通过环境变量TASK_ID或参数 --task-id  指定task-id。
```sh
scalebox event task-add --task-id ${task_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

```scalebox event task-add ``` 可简写为 ``` scalebox event add  ```

### 1.16.2 event slot-add

通过环境变量SLOT_ID或参数 --slot-id  指定slot-id。
```sh
scalebox event slot-add --slot-id ${slot_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

### 1.16.3 event misc-add

```sh
scalebox event misc-add ${tag_name} ${level_name} ${code} ${txt} ${json}
```
