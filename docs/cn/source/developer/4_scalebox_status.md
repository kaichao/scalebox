# 4. Scalebox状态管理

- 状态保持：状态变量
  - 全局变量
  - 共享变量
  - 信号量
- 状态迁移：主路由的任务头信息传递
  - 单次传递
  - vtask头属性：跨vtask生命周期
  - 

信号量是scalebox中用于任务间同步的重要概念，在复杂应用逻辑场景下，集中管理全系统的运行状态，使得计算模块无状态。信号量通常仅在message-router中被读写，以避免并发导致的问题。	而在普通算法模块中可以读取相应值。

## 4.1 信号量（semaphore）

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

## 4.2 共享变量（variable）

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

## 4.3 全局变量

## 4.4 状态迁移

## 4.5 vtask状态

## 4.6 状态变量命名规范
