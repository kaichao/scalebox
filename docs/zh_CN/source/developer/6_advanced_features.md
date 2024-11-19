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

## 6.4 高级消息路由

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

## 6.6 高级流控管理

对于复杂应用程序的本地计算，高级流控可大大简化流控逻辑，避免运行异常。

### 6.6.1 node间并行同步流控
- 标准流控属性```progress_counter_diff```
  - 按节点的进度计数器，基于同步信号量实现
  - progress_counter_diff：当前节点进度与最慢进度的差值作为控制变量

- 消息路由中，按task处理进度，修改对应信号量；
- 业务模块基于信号量计数值，实现同步流控；

### 6.6.2 vtask计数流控

#### vtask定义

消息及其后续消息在多个流水线模块中的完整处理过程。


#### 按计算节点的vtask数量流控
- 属性名：```host_running_vtasks```，计算节点上最多可运行的vtask数量

- yaml应用文件解析时，生成对应信号量：```host-running-vtasks_${mod_name}:${hostname}```，其初值为参数值。
  - mod_name为首模块名，首模块通常为HOST-BOUND，每节点上仅有1个slot
  - hostname为slot对应的hostname

- 处理过程
  - 在首模块入口处，检查信号量计数器，若小于等于0，则流控不通过；
  - 在首模块的消息路由处理中，针对每一条消息，对该计数值自动减一。

  - vtask处理完成，则在消息路由中，通过信号量操作对该计数值加一。

#### 按计算节点组的vtask数量流控
- 属性名：```group_running_vtasks```，计算节点组上最多可运行vtask数量
