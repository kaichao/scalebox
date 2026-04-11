# 3. 模块设计与实现

## 3.1 模块设计原则

- 无状态设计
  - 任务运行的无状态，可重复运行
- 任务级幂等
- 通用性设计，支持任务头驱动
  - 基于主路由的任务头驱动的流程控制

## 3.2 模块设计方法

### 3.2.1 task-body设计

task-body为任务标识，在模块中具有唯一性

- 简单任务，直接用待处理文件名相对路径
- 复杂任务，可用输入路径的规范标识（URI）
- 支持用json格式的非空字符压缩来表示

### 3.2.2 task-headers设计

- 配置参数

在运行时需调整的配置参数通常以环境变量形式传递。在外部编排系统中可以修改配置参数。(运行线程数、)


### 3.2.3 模块内目录设计

#### 代码目录/配置目录
- 通常数据量较小，重点考虑可灵活配置。可放置与共享存储中，或直接打包进容器镜像。

代码的运行顺序：
- 环境变量指定：ACTION_RUN、ACTION_CHECK、ACTION_SETUP、ACTION_TEARDOWN
- /app/bin/{run.sh,check.sh,setup.sh,teardown.sh}
- /app/share/bin/{run.sh,check.sh,setup.sh,teardown.sh}

在模块运行中，通过argument参数code_path，将代码目录映射为容器中的/app/bin。


#### 数据目录

在数据处理量较大的计算模块中，通常需按照计算过程中的使用特性，将目录分类并支持对其灵活配置，以适配不同的计算环境中不同类型的存储形式（内存缓存、本机SSD、本机HDD、网络存储等）,进而提升数据I/O效率。

| 目录类型      | 目录说明                                             |
| ------------ | -------------------------------------------------- |
| 输入数据目录 | 综合考虑读取频度、数据总量，可灵活配置不同数据存储形式。如配置为本机存储，则输入数据由上游模块生成，或通过上游模块从共享存储传输过来 |
| 中间文件目录 | 通常位于本机存储 |
| 输出结果目录 | 综合考虑写出频度、数据总量，可灵活配置不同数据存储形式。如配置为本机存储，则输出数据可供下游游模块读取，或通过下游模块传输至共享存储 |

- 数据目录映射
  - 临时目录 ```/dev/shm```、```/tmp```，自动映射到算法容器中
  - 集群数据目录映射到算法容器中的```/cluster_data_root```
  - 计算节点的```/```映射到容器中的```/local_data_root```
  - 如需额外目录映射，通过```volumes```再做定制映射。


## 3.3 模块集成的编码实现

![slot_arch](../diagrams/slot_arch.drawio.svg)

支持用多种语言实现，推荐使用shell。

sidecar模式：
- run: 任务的单次运行
- check: 任务运行的前置条件检测
- setup: 设置环境
- teardown: 清除环境


### 3.3.1 算法运行run.sh

用户程序（run.sh）运行结束后，agent对用户程序的运行做后处理，采集输入、输出字节数，纪录时间戳及调试信息，并纪录到核心数据库中。两者之间主要通过用户程序的标准输出（stdout）、标准错误（stderr）以及以下文件来交换信息：

|  文件名                           |  文件说明                                   |
| -------------------------------- |  ----------------------------------------- |
| ${WORK_DIR}/task-exec.json       | 任务运行结果主文件，以json形式纪录用户程序运行结果   |
| ${WORK_DIR}/sink-tasks.txt       | 后续任务列表文件，每行一个任务                  |
| ${WORK_DIR}/extra-attributes.txt | 运行附加属性文件，存放于t_task_exec表中extras的extra_attributes中 |
| ${WORK_DIR}/timestamps.txt       | 自定义时间戳文件，常用于调试程序、测试程序性能等     |
| ${WORK_DIR}/input-files.txt      | 输入文件列表，用于统计输入文件字节数             |
| ${WORK_DIR}/output-files.txt     | 输出文件列表，用于统计输出文件字节数             |
| ${WORK_DIR}/removed-files.txt    | 待删除文件列表，一般用于统计文件字节数后再删除     |
| ${WORK_DIR}/auxout.txt           | 辅助输出文件，纪录在最终的任务执行(task_exec)表中 |


### 3.3.2 准入控制check.sh

- 返回值为0，表示准入控制通过，可进行后续任务处理

### 3.3.3 初始设置setup.sh

设置初始运行环境。

### 3.3.4 结束退出teardown.sh

退出前清理计算环境。

## 3.4 模块的镜像封装

### 3.4.1 基于agent构建算法模块
```dockerfile
FROM hub.cstcloud.cn/scalebox/agent:latest
```
### 3.4.2 拷贝agent组件至算法模块
```dockerfile
COPY --from=hub.cstcloud.cn/scalebox/agent:latest /usr/local/ /usr/local/
```


## 3.5 模块单元测试

各模块集成前，对独立算法模块做测试、调试、排错。

### 3.5.1 独立模块的单元测试

缺省的模块代码目录为： ```./code```

- 单独命令行，启动模块测试

```sh
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path}
```

- 先创建模块测试，再添加测试任务
```sh
app_id=$(scalebox run --image-name ${module_image_name} --code-path ${code_path}| cut -d':' -f2 | tr -d '}')

scalebox task add --app-id=${app_id} --header header1=${header1} ${task_body}
```

### 3.5.2 带路由的单元测试

复杂的单元测试，可按需写用于单元测试的路由模块脚本。

用户程序：用任意语言写；

集成程序：一般用bash写。将用户程序的结果写回。

缺省的路由代码目录为： ```./mr-code```

```sh
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path} --mr-image-name ${mr_image_name} --mr-code-path ${mr_code_path}
```

### 3.5.3 基于流水线应用的单元测试

- 写独立的应用定义文件app.yaml、环境变量定义文件scalebox.env


## 3.6 错误处理与容错机制




