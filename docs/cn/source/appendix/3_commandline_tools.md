# 3. 命令行工具指南

## 3.1 scalebox命令行工具详解

## 3.2 常用操作命令

## 3.3 脚本编写示例

## 3.4 自动化运维脚本



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


## 1.2 子命令概览

```
scalebox
├── run                    启动应用（创建 App + 发送初始消息）
├── cluster                集群管理
│   ├── list               列出集群
│   ├── create             创建集群
│   ├── show               查看集群详情
│   ├── set-status         设置集群状态
│   ├── get-parameter      获取集群参数
│   ├── allocate           分配集群资源
│   └── release            释放集群资源
├── host                   主机管理
│   ├── list               列出主机
│   ├── create             添加主机
│   ├── show               查看主机详情
│   ├── set-status         设置主机状态
│   ├── delete             删除主机
│   └── renew              续期主机
├── slot                   插槽管理
│   ├── list               列出插槽
│   ├── add                添加插槽
│   ├── remove             删除插槽
│   └── set-status         设置插槽状态
├── app                    应用管理
│   ├── list               列出应用
│   ├── create             创建应用
│   ├── show               查看应用详情
│   ├── set-status         设置应用状态
│   ├── set-finished       标记应用完成
│   ├── delete             删除应用
│   ├── add-remote         添加远程应用链接（跨集群）
│   └── add-slots          动态添加插槽
├── module                 模块管理
│   ├── list               列出模块
│   └── show               查看模块详情
├── task                   任务管理
│   ├── list               列出任务
│   ├── show               查看任务详情
│   ├── add                添加任务
│   ├── delete             删除任务
│   ├── get-header         获取任务头
│   ├── set-header         设置任务头
│   ├── remove-header      删除任务头
│   └── log                查看任务日志
├── validate               校验应用定义文件（app.yaml）
├── semaphore              信号量管理
│   ├── list               列出信号量
│   ├── create             创建信号量
│   ├── get                获取信号量值
│   ├── increment          增一
│   ├── decrement          减一
│   └── add-value          增减 N
├── semagroup              信号量组操作
│   ├── max                最大值
│   ├── min                最小值
│   ├── increment          最小值加一
│   ├── decrement          最大值减一
│   ├── diff-min           与最小值的差
│   └── diff-max           与最大值的差
├── variable               共享变量管理
│   ├── list               列出变量
│   ├── get                获取变量值
│   ├── set                设置变量值
│   └── delete             删除变量
├── global                 全局变量管理
│   ├── list               列出全局变量
│   ├── get                获取值
│   ├── set                设置值
│   └── delete             删除
├── vtask                  虚拟任务管理
│   ├── list               列出 vtask
│   ├── list-subtasks      列出子任务
│   ├── get                查看 vtask 详情
│   ├── fail               标记 vtask 失败
│   ├── bind               绑定资源
│   ├── unbind             解绑资源
│   ├── add-subtask        添加子任务
│   ├── get-variable       获取 vtask 作用域变量
│   ├── set-variable       设置 vtask 作用域变量
│   ├── create-semaphore   创建 vtask 作用域信号量
│   ├── get-semaphore      获取 vtask 作用域信号量
│   ├── add-semaphore-value 增减 vtask 作用域信号量
│   └── delete-semaphore   删除 vtask 作用域信号量
├── channel                优先级队列（跨应用通信）
│   ├── pull               出队
│   └── push               入队
├── fs                     文件系统操作
│   ├── ls                 列出文件
│   └── stat               查看文件元数据
├── event                  事件记录
│   ├── task-add           任务事件
│   ├── slot-add           插槽事件
│   └── misc-add           杂项事件
├── status                 系统整体状态
└── help                   帮助信息
```

## 1.3 cluster 子命令

### 1.3.1 cluster list

列出所有集群。

```bash
scalebox cluster list
```

### 1.3.2 cluster show

查看指定集群详情。

```bash
scalebox cluster show <cluster-name>
```

### 1.3.3 cluster create

创建新集群。

```bash
scalebox cluster create --name my-cluster --grpc-server 10.0.0.1:50051
```

### 1.3.4 cluster set-status

设置集群状态。

```bash
scalebox cluster set-status <cluster-name> ON
```

### 1.3.5 cluster get-parameter

获取集群参数。

```bash
scalebox cluster get-parameter <cluster-name> <param-name>
```

### 1.3.6 cluster allocate / release

动态集群资源的分配和释放。

```bash
scalebox cluster allocate --cluster my-cluster --num-hosts 4
scalebox cluster release --cluster my-cluster
```

## 1.4 host 子命令

主机（计算节点）管理。

### 1.4.1 host list

列出所有主机。

```bash
scalebox host list
scalebox host list --cluster my-cluster
```

### 1.4.2 host show

查看主机详情。

```bash
scalebox host show <hostname>
```

### 1.4.3 host create

添加新主机。

```bash
scalebox host create --hostname n0 --ip-addr 10.0.6.101 --cluster my-cluster
```

### 1.4.4 host set-status

设置主机状态。

```bash
scalebox host set-status <hostname> READY
```

### 1.4.5 host delete / renew

删除主机或续期动态主机。

```bash
scalebox host delete <hostname>
scalebox host renew <hostname>
```

## 1.5 slot 子命令

执行槽（Slot）管理。

### 1.5.1 slot list

列出插槽。

```bash
scalebox slot list
scalebox slot list --host <hostname>
scalebox slot list --module <module-name>
```

### 1.5.2 slot add

添加插槽。

```bash
scalebox slot add --module my-module --host n0 --count 4
```

### 1.5.3 slot remove

删除插槽。

```bash
scalebox slot remove --module my-module --host n0
```

### 1.5.4 slot set-status

设置插槽状态。

```bash
scalebox slot set-status <slot-id> READY
```

### 1.5.5 slot get-parameter

读取插槽参数值（```t_slot.parameters```，jsonb）。

```bash
scalebox slot get-parameter <parameter> --slot-id <slot-id>
```

参数未设置时返回空字符串。

### 1.5.6 slot set-parameter

设置插槽参数值（```t_slot.parameters```，jsonb）。值以字符串类型存储。

```bash
scalebox slot set-parameter <parameter> <value> --slot-id <slot-id>
```

插槽参数可用于slot启动命令表达式（```((@p:param_name))```），详见《高级特性》§7.7。


## 1.6 app 子命令

### 1.6.1 app list

列出所有应用。

```bash
scalebox app list
```

### 1.6.2 app show

查看应用详情。

```bash
scalebox app show <app-id>
```

### 1.6.3 app create

从定义文件创建应用。

```bash
scalebox app create --app-file app.yaml --env-file scalebox.env
```

### 1.6.4 app run

以命令行方式启动应用。环境变量文件缺省为 `./scalebox.env`。

**单启动消息**：
```bash
export ENV0=v0
scalebox run --cluster my-cluster --image-name my-image:latest start-item
```

**参数表**：

| 参数名 | 对应环境变量 | 缺省值 | 说明 |
|--------|------------|--------|------|
| app-name | `_APP_NAME` | — | 应用名称 |
| cluster | `_CLUSTER` | local | 集群名 |
| image-name | `_IMAGE_NAME` | scalebox/agent:latest | 主模块镜像 |
| code-path | `_CODE_PATH` | ./code（若存在） | 主模块代码目录 |
| slot-regex | `_SLOT_REGEX` | h0 | 主模块 slot 配置 |
| mr-image-name | `_MR_IMAGE_NAME` | — | 路由模块镜像 |
| mr-code-path | `_MR_CODE_PATH` | ./mr-code（若存在） | 路由模块代码目录 |
| app-file / -f | — | app.yaml（若存在） | 应用定义文件 |
| env-file / -e | — | scalebox.env | 环境变量文件 |

**管道式多启动任务**：
```bash
find /data/input -type f | scalebox run --image-name my-image:latest --slot-regex h0:2
```

### 1.6.5 app set-status / set-finished

```bash
scalebox app set-status <app-id> RUNNING
scalebox app set-finished <app-id>
```

### 1.6.6 app delete

```bash
scalebox app delete <app-id>
```

### 1.6.7 app add-remote

添加跨集群应用远程链接。

```bash
scalebox app add-remote --app-id <local-app-id> --remote-app-id <remote-app-id> --remote-grpc <grpc-addr>
```

### 1.6.8 app add-slots

动态添加插槽。

```bash
scalebox app add-slots --app-id <app-id> --module <module-name> --host <hostname> --count 4
```

## 1.7 module 子命令

### 1.7.1 module list

列出应用的模块。

```bash
scalebox module list --app-id <app-id>
```

### 1.7.2 module show

查看模块详情。

```bash
scalebox module show --app-id <app-id> <module-name>
```

## 1.8 task 子命令

### 1.8.1 task list

列出任务，支持按模块、状态过滤。

```bash
scalebox task list --app-id <app-id>
scalebox task list --app-id <app-id> --status FAILED
scalebox task list --app-id <app-id> --module <module-name>
```

### 1.8.2 task show

查看任务详情（含 stdout/stderr/body/headers）。

```bash
scalebox task show <task-id>
```

### 1.8.3 task delete

删除任务。

```bash
scalebox task delete <task-id>
```

### 1.8.4 task add

添加任务。

**参数/环境变量**：

| 参数名 | 环境变量名 | 说明 |
|--------|----------|------|
| app-id | `APP_ID` | 应用 ID |
| module-id | `MODULE_ID` | 模块 ID |
| sink-module | `SINK_MODULE` | 下游模块名 |
| conflict-action | `CONFLICT_ACTION` | 冲突处理：''/'IGNORE'/'OVERWRITE' |
| from-module | — | 来源模块名 |
| headers | — | JSON 格式 headers |
| header / -h | — | 添加单个 header（可重复） |
| to-ip | — | 设置 to_ip header |
| to-host | — | 设置 to_host header |
| batch-size | — | 批量添加批次大小，缺省 100 |

**task 文件格式**（缺省为 `${WORK_DIR}/sink-tasks.txt`，每行一条）：

| 类型 | 示例 |
|------|------|
| 文本 body | `body` |
| JSON body | `{"hi0":"a","body":"my_body"}` |
| 文本 body + headers | `body,{"h0":"a","h1":"b"}` |
| JSON body + headers | `{"hi0":"a","body":"my_body"},{"h0":"a","h1":"b"}` |
| 模块名 + 文本 body | `module-name,body` |
| 模块名 + 文本 body + headers | `module-name,body,{"h0":"a","h1":"b"}` |

**控制 headers**：

| header | 说明 |
|--------|------|
| `initial_status_code` | 初始状态，缺省 -1（READY） |
| `conflict-action` | ''/'IGNORE'/'OVERWRITE' |
| `slot_broadcast` | 广播到所有 slot |
| `host_broadcast` | 广播到所有 host |

### 1.8.5 task get-header / set-header / remove-header

```bash
scalebox task get-header --task-id 123 from_module
scalebox task set-header --task-id 123 my_header value
scalebox task remove-header --task-id 123 my_header
```

### 1.8.6 task log

查看任务执行日志。

```bash
scalebox task log <task-id>
```

## 1.9 validate 子命令

校验应用定义文件（app.yaml）的语法和完整性。

```bash
scalebox validate --app-file=<app-yaml>
scalebox validate --app-file=app.yaml --env-file=scalebox.env
```

## 1.10 vtask 子命令

VTask（虚拟任务）是 App 内跨模块的 task 集合，通过头-核心-尾管道实现流控和状态管理。

公共参数：
- `--app-id`：应用 ID（也可通过 `APP_ID` 环境变量设置）

### 1.9.1 vtask list

列出应用的 vtask 列表。

```bash
scalebox vtask list --app-id <app-id>
```

### 1.9.2 vtask list-subtasks

列出指定 vtask 的子任务。

```bash
scalebox vtask list-subtasks --app-id <app-id> <vtask-id>
```

### 1.9.3 vtask get

查看 vtask 详情（含子任务计数、信号量名等）。

```bash
scalebox vtask get <vtask-id>
```

### 1.9.4 vtask fail

标记 vtask 失败。将自动释放门控信号量和资源信号量，并级联标记未完成子任务。

```bash
scalebox vtask fail <vtask-id>
```

### 1.9.5 vtask bind / unbind

绑定/解绑 vtask 的计算资源。

```bash
scalebox vtask bind --app-id <app-id> --sema-name <sema-name>
scalebox vtask unbind --app-id <app-id> --sema-name <sema-name>
```

### 1.9.6 vtask add-subtask

向 vtask 添加子任务。在 agent 内部运行时自动传播 `_vtask_id`、`_vtask_size_sema` 等 header。

```bash
scalebox vtask add-subtask --app-id <app-id> --module <module-name> --body <task-body>
```

### 1.9.7 vtask 作用域变量

```bash
# 获取 vtask 变量
scalebox vtask get-variable --vtask-id <vtask-id> <var-name>

# 设置 vtask 变量
scalebox vtask set-variable --vtask-id <vtask-id> <var-name> <value>
```

### 1.9.8 vtask 作用域信号量

```bash
# 创建
scalebox vtask create-semaphore --vtask-id <vtask-id> <sema-name> <initial-value>

# 获取
scalebox vtask get-semaphore --vtask-id <vtask-id> <sema-name>

# 增减
scalebox vtask add-semaphore-value --vtask-id <vtask-id> <sema-name> <delta>

# 删除
scalebox vtask delete-semaphore --vtask-id <vtask-id> <sema-name>
```

## 1.11 semaphore 子命令

- 公共参数：`--app-id` 或 `--module-id`（也可通过 `APP_ID` / `MODULE_ID` 环境变量）
- 环境变量：`SEMAPHORE_AUTO_CREATE=yes` 时，信号量不存在则自动创建（初值为 0）。CLI 通过 gRPC metadata `semaphore-auto-create: yes` 传递给 controld

- 信号量命名规则：
  - 字符集：`[A-Za-z0-9:_-]`
  - 首字符为字母或下划线

- 信号量表达式：表示一组信号量的正则表达式，字符集加上 `.*+?^$[]{}()|\`

### 1.11.1 semaphore list

列出信号量，支持前缀和叶子节点过滤。

```bash
scalebox semaphore list --app-id <app-id>
scalebox semaphore list --app-id <app-id> --prefix vtask_size
scalebox semaphore list --app-id <app-id> --leaf-only
```

### 1.11.2 semaphore create

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

### 1.11.3 semaphore get

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

### 1.11.4 semaphore increment

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

### 1.11.5 semaphore decrement

#### 单个信号量的减一操作。

```sh
val=$(scalebox semaphore decrement ${sema_expr} )
code=$?
```

用法详见：<a href="#semaphore-increment">semaphore increment</a>

#### 信号量组的减一操作。

用法详见：<a href="#semaphore-increment">semaphore increment</a>

### 1.11.6 semaphore add-value

#### 单个信号量的加n操作。
```sh
val=$(scalebox semaphore add-value ${sema_name} ${delta_value})
code=$?
```

用法详见：<a href="#semaphore-increment">semaphore increment</a>

#### 信号量组的加n操作。

用法详见：<a href="#semaphore-increment">semaphore increment</a>


### 1.11.7 semaphore delete

删除信号量。

```sh
scalebox semaphore delete --app-id ${app_id} ${sema_name}
```

## 1.12 semagroup 子命令

- 多个信号量组成信号量组，用信号量名前缀、正则表达式标识信号量组


### 1.12.1 semagroup max
- 信号量组中最大值
```sh
val=$(scalebox semagroup max ${sema_expr})
code=$?
```

### 1.12.2 semagroup min
- 信号量组中最小值
```sh
val=$(scalebox semagroup min ${sema_expr})
code=$?
```
- sema_expr为信号量名的正则表达式或前缀
- 返回值val为整数字符串

### 1.12.3 semagroup increment
- 选取信号量组中最小值，并加一

### 1.12.4 semagroup decrement
- 选取信号量组中最大值，并减一

### 1.12.5 semagroup diff-max
- 信号量组最大值与信号量当前值的差值
```sh
val=$(scalebox semagroup diff-max ${sema_expr})
code=$?
```
- sema_expr为含分组定义的信号量，示例为```(group-prefix):sema-suffix```
- 返回值val为整数字符串

### 1.12.6 semagroup diff-min
- 信号量当前值与信号量组最小值的差值
```sh
val=$(scalebox semagroup diff-min ${sema_expr})
code=$?
```
- sema_expr为含分组定义的信号量，示例为```(group-part:)sema-suffix```
- 返回值val为整数字符串


## 1.13 variable 子命令

- 公共参数：`--app-id` 或 `--module-id`（也可通过 `APP_ID` / `MODULE_ID` 环境变量）
- 变量名命名：同信号量命名（`[A-Za-z0-9:_-]`，首字符为字母或下划线）
- 支持 `list`、`get`、`set`、`delete` 操作

### 1.16.1 variable list

列出变量，支持前缀和叶子节点过滤。

```bash
scalebox variable list --app-id <app-id>
scalebox variable list --app-id <app-id> --prefix my_prefix
scalebox variable list --app-id <app-id> --leaf-only
```

### 1.16.2 variable get

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

### 1.15.3 variable set

```sh
scalebox variable set --app-id ${app_id} ${var_name} ${str_value}
APP_ID=${app_id} scalebox variable set ${var_name} ${str_value}
```

### 1.14.4 variable delete

```sh
scalebox variable delete --app-id ${app_id} ${var_name}
```

## 1.14 global 子命令

全局变量，跨应用共享。

### 1.16.1 global list

列出全局变量，支持前缀和叶子节点过滤。

```bash
scalebox global list
scalebox global list --prefix my_prefix
scalebox global list --leaf-only
```

### 1.16.2 global get

```sh
scalebox global get ${global_name}
```

### 1.15.3 global set

```sh
scalebox global set ${global_name} ${global_value}
```

### 1.14.4 global delete

```sh
scalebox global delete ${global_name}
```


## 1.15 channel 子命令

channel用于跨应用间的通信，是一个有优先级队列。

- 公共参数：module-id，或app-id
- 环境变量：MODULE_ID，或APP_ID

- 优先队列命名：同信号量命名

### 1.16.1 channel create

- head-app : app-id
- tail-app

若不指定，则为当前app

### 1.16.2 channel pull

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

### 1.15.3 channel push

- priority为优先级，浮点数。数值小，优先级高。
  
```sh
scalebox channel push --app-id ${app_id} ${pp_name} ${str_value} [${priority}]
APP_ID=${app_id} scalebox channel push ${pp_name} ${str_value}

scalebox channel push --module-id ${module_id} ${pp_name} ${str_value} [${priority}]
MODULE_ID=${module_id} scalebox channel push ${pp_name} ${str_value}
```


## 1.16 fs 子命令

scalebox-fs以文件系统形式，将分布式计算节点上的文件组织在同一个名字空间中。后期可提供mount支持、跨节点迁移等特性。

### 1.16.1 fs ls

- 主要参数：
  - include-removed-file
  - with-hostname
  - with-file-size
  - hostname=${host-name}

```sh
scalebox fs ls ${path_expr}
```

### 1.16.2 fs stat

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

## 1.17 status

- 系统整体状态：local集群头节点（actuator到local头节点有效）
- cluster列表：不同状态的host数量
- app列表：不同状态

## 1.18 event 子命令

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


### 1.18.1 event task-add

通过环境变量TASK_ID或参数 --task-id  指定task-id。
```sh
scalebox event task-add --task-id ${task_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

```scalebox event task-add ``` 可简写为 ``` scalebox event add  ```

### 1.18.2 event slot-add

通过环境变量SLOT_ID或参数 --slot-id  指定slot-id。
```sh
scalebox event slot-add --slot-id ${slot_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

### 1.18.3 event misc-add

```sh
scalebox event misc-add ${tag_name} ${level_name} ${code} ${txt} ${json}
```
