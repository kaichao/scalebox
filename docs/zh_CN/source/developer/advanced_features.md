# 6. 高级编程特性

## 6.1 跨集群应用


## 6.2 动态集群

部署集群管理应用cluster-admin，支持集群节点的动态扩容，通过node-agent实现计算节点的系统管理、动态监控等。


## 6.3 本地计算（Local Compute）

- task_dist_mode: 
  - HOST-BOUND: headers->>'to_host'
  - SLOT-BOUND
  - GROUP-BOUND

硬件条件：
- 需要真正的本地存储支持，云上的虚机不适合本地计算。
- 自建服务器集群（不带虚拟化）
- HPC计算节点


本地计算的主要技术问题
- 前后job间的流控
  - 运行速度匹配，避免本地硬盘/内存爆满
  - 信号量是实现流控的技术手段
- 同一job在多节点间同步流控
  - 不同通道的处理/不同节点的数据下载
  - 节点处理能力差异，因此中间数据累积，导致内存爆满
- 自动容错处理
  - 出错影响部分数据的处理
  - 流控会导致整个流水线停止
  - 自动容错对流控的影响
- 关键job的task处理排序
- 非关键模块分解
  - 细分会增加流控负担
  - 增加流水线并行提升性能，会增加中间存储的内存使用，影响流水线运行的稳定性



## 6.4 消息路由

开启task头获取功能。
- variables->>slot_options 'with_headers'
- WITH_HEADERS=yes

- 多message-router实例设置
  - 将所有相关信号量模块放到一个messsage-router中。
  - message-router设定为SLOT-BOUND
  - task add过程中，在task-headers中设定to_slot

可用不同程序语言实现。
- bash
- python
- golang

## 6.5 消息排序

- key_group_regex
- key_group_index
- sorted_tag

## 6.6 流控管理
流控是控制每个task运行进度的机制。

### 6.6.1 task级流控
- 标准流控属性（job间流控）
  - dir_limit_gb
  - dir_free_gb
- 代码定制流控（可基于信号量）
  - ACTION_CHECK/check.sh

### 6.6.2 node间并行同步流控
- 标准流控属性
  - 按节点的进度计数器，基于同步信号量实现
  - progress_counter_diff：当前节点进度与最慢进度的差值作为控制变量

- 代码定制流控
  - 基于ACTION_CHECK/check.sh，自定义逻辑

### 6.6.2 业务代码精准流控
精准控制数据输出，若数据空间不满足，在代码中显式调用sleep()，等待其他并行任务完成，清理出空间，待空间容量满足要求后再继续运行。

避免内存空间的峰值需求，以增加等待时间为代价，在低内存配置的机器上运行本地计算的应用。


## 6.7 信号量

信号量是scalebox中用于任务间同步的重要概念，在复杂应用逻辑场景下，集中管理全系统的运行状态，使得计算模块无状态。信号量通常仅在message-router中被调用，以避免并发导致的问题。					

### 信号量名称的命名规范

