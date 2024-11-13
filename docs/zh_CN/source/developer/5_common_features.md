# 5. 常见编程特性

## 5.1 超时设置（timeout）

由于代码的bug、数据错误等原因，可能导致算法模块进入死循环，一直占用slot运行，无法正常退出。

超时设置可以解决以上问题。通过设置job的```task_timeout_seconds```参数，在超出该时限后，该任务（task）会退出运行，返回错误码 124.

在封装脚本调用算法代码时，建议也在封装脚本中加入timeout的相关设置。若无此设置，外层脚本timeout，内层脚本还会运行结束，导致不确定的结果。
若算法代码名为run.py，则封装代码中对应部分可写成：
```sh
timeout ${TASK_TIMEOUT_SECONDS}s run.py $*
```

## 5.2 自动重试 retry

返回码是模块运行的结果代码，模块运行出错的返回码就是错误码。

模块运行错误通常分为两类：
- 逻辑错误：代码逻辑导致的结果错误；
- 非逻辑错误：外部软硬件环境异常导致的结果错误；通常在外部环境恢复后，重复多次运行，可解决；

非逻辑错发生有一定的概率，大规模数据处理中，非逻辑错的自动处理，将大大提升数据处理的效率及可靠性。

### 常见的外部运行环境异常
- 不稳定的硬件和基础软件
  - 网络设备链路的不稳定
  - 环境温度升高，导致计算节点暂时停止服务
  - 磁盘盘阵中介质不稳定，导致RAID暂时失效
  - 集群存储软件的临时性、可修复的bug
- 运行中的随机异常
  - 高负载导致核心存储响应延时增加
    - 可变延时（数秒~几分钟）
    - 暂时性不可写入
  - 高负载导致低概率的内存超分，应用内存不足
  - 系统高负载，不合适的调度策略，导致算法运行超时
- 数据异常
  - 异常输入数据
  - 前导处理软件的随机错误
  - 各类运行异常导致的数据不完整

### 非逻辑错的容错处理
- 系统异常的判断与定位
  - 读取/执行/回写的时间
  - 任务返回码
  - 系统监控的异常信息
- 异常任务的容错处理
  - 自动/人工

### 运行出错的自动重试（retry）

按照算法代码对于非逻辑错的错误，设置主脚本不同的错误码。按返回错误码，在模块定义中设置自动重试，从而实现自动重试。

在job定义中，设置```retry_rules```实现以上功能。

```
    retry_rules: "['91','92:3']"
```
在以上设置中，若返回码为91，则重试1次；若返回码为92，则重试3次。

## 5.3 流控设置

算法模块的流量控制，就是控制流式数据处理的进度，以免数据处理所需不可共享的核心计算资源（内存缓存、GPU显存、本地磁盘空间）超过本地资源的物理限额，导致计算过程中出现异常，包括异常退出、程序死锁等。

流量控制的实现机制是在消息处理前，通过模块的标准流控规则、自定义流控规则检查是否符合外部条件，若不符合，则在代码中显式调用sleep()，等待一段时间后再做检查。经过多次检查未果，则退出模块运行。

流控规则通过引入等待机制，降低了局部计算资源的利用率，但从整体上保证应用程序的顺利执行。
避免内存空间的峰值需求，以增加等待时间为代价，在低内存配置的机器上运行本地计算的应用。

流控分为slot级流控、node间并行同步流控。
流控支持模块运行的同步，以免系统资源受限而导致运行异常。

精准控制数据输出，若数据空间不满足，等待其他并行任务完成后释放资源，清理出空间，待空间容量满足要求后再继续运行。



### 5.3.1 slot级级流控
- 标准流控属性（job间流控）
  - dir_limit_gb：通过限制目录的最大占用存储空间来实现流控。
  - dir_free_gb：通过限制目录所在存储分区的最小空余空间实现流控
- 代码定制流控（可基于信号量）
  - ACTION_CHECK/check.sh：自定义流控逻辑

### 5.3.2 node间并行同步流控
- 标准流控属性```progress_counter_diff```
  - 按节点的进度计数器，基于同步信号量实现
  - progress_counter_diff：当前节点进度与最慢进度的差值作为控制变量

- 消息路由中，按task处理进度，修改对应信号量；
- 业务模块基于信号量，实现同步流控；


## 5.4 关键模块的task运行排序

task运行排序是scalebox应用运行的重要基础。

通过设置以下参数，实现排序。

- key_group_regex：正则表达式，从消息体中提取相关分组字符串。
- key_group_index：正则表达式对应的分组编号。
- sorted_tag：高优先级排序编号，通常有message-router设置


## 5.5 多GPU配置

单个计算节点可配置多个GPU加速卡。通常为每个GPU配置独立的slot，模块算法就不需要针对多GPU卡做并行优化，这样配置通常也更具高运行效率。
在这种场景下，不同GPU需要用不同的slot启动命令。
scalebox启动命令支持参数化配置。

参数化命令配置如下：

```
docker run -d --rm --network host --tmpfs=/work --device=/dev/kfd --device=/dev/dri --security-opt seccomp=unconfined --group-add video -e ROCR_VISIBLE_DEVICES={~n~} {{ENVS}} {{VOLUMES}} {{IMAGE}}

docker run -d --rm --network=host --tmpfs=/work --device=/dev/kfd --device=/dev/dri/card{~n%2~} --device=/dev/dri/renderD{~n%2+128~} --security-opt seccomp=unconfined --group-add video --cap-add=SYS_PTRACE {{ENVS}} {{VOLUMES}} {{IMAGE}}
```

其中， ```{~ ~}```中间的为运行时表达式，其中n为节点上slot编号（从0开始）。在启动slot前，动态解析为对应值。

## 5.6 共享变量（信号量及普通变量）

信号量是scalebox中用于任务间同步的重要概念，在复杂应用逻辑场景下，集中管理全系统的运行状态，使得计算模块无状态。信号量通常仅在message-router中被读写，以避免并发导致的问题。	而在普通算法模块中可以读取相应值。

### 信号量

- 信号量创建
```sh
scalebox semaphore create ${sema_name}
```
- 信号量读取
```sh
scalebox semaphore get ${sema_name}
```

- 信号量增一

```sh
scalebox semaphore increment ${sema_name}
```

- 信号量减一
```sh
scalebox semaphore decrement ${sema_name}
```

- 信号量增减
```sh
scalebox semaphore increment-n {sema_name} ${n}
```

- 信号量组距离
```sh
scalebox semaphore group-dist ${sema_name} ${n}
```


### 普通变量

普通变量通常为字符串类型。

- 变量创建
```sh
scalebox variable create ${var_name}
```
- 变量读取
```sh
scalebox variable get ${var_name}
```
- 变量写入
```sh
scalebox variable set ${var_name} ${value}
```


### 信号量名称的命名规范



## 5.7 task-perspective

- scalebox标准时间格式
```2023-10-24T18:00:00.123456+0800```

- bash生成当前时间 
```sh
date +"%Y-%m-%dT%H:%M:%S.%6N%z"
```

- golang生成当前时间 
```go
formattedTime := time.Now().Format("2006-01-02T15:04:05.999999")
fmt.Println(formattedTime)
```

- python生成当前时间 
```python

```


## 5.8 git应用代码库
