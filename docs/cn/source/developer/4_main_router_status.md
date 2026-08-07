# 4. 主路由及应用状态

Scalebox应用是由主路由模块驱动的分布式应用程序，其状态管理是其中的关键技术。状态管理通过状态保持、状态迁移来实现。

主路由模块，又称为路由模块，应用的主控程序。


## 4.1 主路由模块

- 多main-router实例设置
  - 将所有相关信号量模块放到一个main-router中。
  - task-dist-mode设定为SLOT-BOUND
  - task add过程中，在task-headers中设定to_slot

可用不同程序语言实现。
- bash
- golang
- python


### 4.1.1 优化设计
- 主路由模块从设计上看，是单个串行模块。但可按不同来源消息部署为多个独立模块，实现并行处理。
- 减少路由相关的输入任务数量；（粗粒度的控制消息）
- 减少全系统的信号量数量；（粗粒度信号量）

### 4.1.2 高效实现
- 在主路由中，单次处理多个任务（一次操作读取多个任务，并批量处理后续任务生成、信号量修改等）
- 后续任务生成的批量操作（减少数据库操作次数）
- 信号量处理的批量生成（减少数据库操作次数）

### 4.1.3 状态管理

- 状态保持
- 状态迁移


信号量是scalebox中用于任务间同步、状态管理的重要概念，在复杂应用逻辑场景下，集中管理全系统的运行状态，使得计算模块无状态。信号量通常仅在main-router中被读写，以避免并发导致的问题。	而在普通算法模块中可以读取相应值。

## 4.2 状态保持
- ：以状态变量形式，存储应用模块的当前状态。
  - 全局变量
  - 共享变量
  - 信号量

  - 状态保存（semaphore/variable/global/......）

### 4.2.1 信号量（semaphore）

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

- 信号量组操作（semagroup）

### 4.2.2 共享变量（variable）

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

## 4.2.3 全局变量

- 跨应用的共享变量
- 全局配置参数

## 4.3 状态转移

### 4.3.1 任务头状态传递
- 主路由的任务头信息传递
- 任务头以下划线```_```开头，标识为待传递任务头
- 传递路径：主路由 -> 算法模块 -> 主路由
- 单次传递，需人工维护状态

### 4.3.2 vtask状态保持
- 通过任务头传递实现vtask状态保持
- 任务头以下划线```_vtask_```开头，标识为vtask任务头属性
- vtask属性的保持：横跨vtask的各个模块
- vtask头属性：跨vtask生命周期，自动维护


## 4.4 vtask管理

> 完整的 VTask 设计文档（概念模型、模块结构、信号量机制、管道流程、命令参考）见 :doc:`VTask 虚拟任务 <9_vtask>`。以下为概览。

vtask作为节点本地计算编程模型的基本单元，通过参数设计，简化本地计算在应用程序中的流控、容错等的实现；

vtask跨越输入加载、计算、结果回写的整个过程，基于vtask的容错简单、直接。

vtask标识一种跨越多个模块的虚拟task，是基于计数的流控机制，分为全局计数（global_vtask）、分组计数（group_vtask）、节点计数（host_vtask）等三类。

## 4.4.1 vtask主要功能
- 任务管理：跨模块的任务（task）组成vtask，可简化任务管理，并支持状态管理，进而支撑节点本地计算模型；
- 计数准入控制：以配合HPC计算资源的动态调度，实现超长任务计算；
- 粗粒度容错；实现细粒度检查点功能，自动容错。

## 4.4.2 vtask模块结构

- 前置任务队列（wait-queue）：尚未绑定计算资源。全局容错的基本单元。
- vtask头模块（vtask-head）：vtask起始标识模块。为vtask分配计算资源（单节点/资源组），单节点模式使用HOST-BOUND，直接运行在计算节点上；资源组模式使用SLOT-BOUND，多个slot需部署在单节点（通常为头节点）上。
- vtask算法模块（vtask-core）：vtask核心处理模块，针对资源组模式，通常用hostname前缀来标识；可以为多个模块，通过pod实现模块间本地分发。
- vtask尾模块（vtask-tail）：vtask结束标识模块，通常部署在头节点上。任务体（task-body）与vtask-head保持一致。


## 4.4.3 vtask模式分类

task_dist_mode设置

|            |   全局模式   |  单节点模式  | 资源组模式  |
| ---------- | ----------- | ---------- | ---------- |
| wait-queue |             |            |            |
| vtask-head |             | HOST-BOUND | SLOT-BOUND |
| vtask-core |             | HOST-BOUND | HOST-BOUND |
| vtask-tail |             |            |            |


## 4.5 状态变量命名规范

- 信号量（semaphore）：字符集 `[A-Za-z0-9:_-]`，首字符为字母或下划线
- 共享变量（variable）：同信号量命名规则
- 全局变量（global）：同信号量命名规则
- vtask 信号量命名模式：
  - 全局：`vtask_size:${mod_name}`
  - 节点级：`host_vtask_size:${mod_name}:${hostname}`
  - 槽位级：`slot_vtask_size:${mod_name}:${slot_id}`

## 4.6 最佳实践

- **应用分解**：路由模块为全系统集中模块，设计中减少消息数量，有利于在消息路由中批量处理算法模块的信号量、消息发送等，提高系统效率
- **信号量设计**：尽可能采用粗粒度，减小管理数据的规模
- **主路由高效实现**：避免因为与流控联动影响并行效率
- **批量消息处理**：一次操作读取多个任务，批量生成后续任务，减少数据库操作次数
- **信号量批量操作**：批量增减信号量，减少数据库往返