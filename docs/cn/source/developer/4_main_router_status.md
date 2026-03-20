# 4. 主路由及状态管理

Scalebox应用是由主路由模块驱动的分布式应用程序，其状态管理是其中的关键技术。状态管理通过状态保持、状态迁移来实现。

- 状态保持：以状态变量形式，存储应用模块的当前状态。
  - 全局变量
  - 共享变量
  - 信号量
- 状态迁移：主路由的任务头信息传递
  - 单次传递：主路由 -> 算法模块 -> 主路由，需人工维护状态
  - vtask头属性：跨vtask生命周期，自动维护
  - 


## 4.1 主路由模块

- 多main-router实例设置
  - 将所有相关信号量模块放到一个messsage-router中。
  - task-dist-mode设定为SLOT-BOUND
  - task add过程中，在task-headers中设定to_slot

可用不同程序语言实现。
- bash
- python
- golang


### 4.1.1 优化设计
- 消息路由从设计上看，是单个串行模块。但可按不同来源消息部署为多个独立模块，实现并行处理。
- 减少消息路由相关的输入消息数量；（粗粒度的控制消息）
- 减少全系统的信号量数量；（粗粒度信号量）

### 4.1.2 高效实现
- 在消息路由中，单次处理多个消息（一次操作读取多个消息，并批量处理后续消息生成、信号量修改等）
- 后续消息生成的批量操作（减少数据库操作次数）
- 信号量处理的批量生成（减少数据库操作次数）



信号量是scalebox中用于任务间同步的重要概念，在复杂应用逻辑场景下，集中管理全系统的运行状态，使得计算模块无状态。信号量通常仅在main-router中被读写，以避免并发导致的问题。	而在普通算法模块中可以读取相应值。

## 4.2 信号量（semaphore）

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

## 4.3 共享变量（variable）

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

## 4.4 全局变量

## 4.5 状态迁移

## 4.6 vtask状态

## 4.7 状态变量命名规范
