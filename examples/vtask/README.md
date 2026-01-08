# vtask

vtask作为节点本地计算编程模型的基本单元，通过参数设计，简化本地计算在应用程序中的流控、容错等的实现；

主要功能：
- 计数流控：以配合HPC计算资源的动态调度，实现超长任务计算；
- 粗粒度容错；实现细粒度检查点功能；

示例中包括以下模块：
- wait-queue：vtask的全局等待队列，尚未绑定计算资源。全局容错的基本单元。
- vtask-head：vtask起始标识模块。为vtask分配计算资源（单节点/资源组），单节点模式使用HOST-BOUND，直接运行在计算节点上；资源组模式使用SLOT-BOUND，多个slot需部署在单节点（通常为头节点）上。
- vtask-core：vtask核心处理模块，针对资源组模式，通常用hostname前缀来标识；可以为多个模块，通过pod实现模块间本地分发。
- vtask-tail：vtask结束标识模块，通常部署在头节点上。task的主体（body）与vtask-head一致。

task_dist_mode设置

|            |   全局模式   |  单节点模式  | 资源组模式  |
| ---------- | ----------- | ---------- | ---------- |
| wait-queue |             |            |            |
| vtask-head |             | HOST-BOUND | SLOT-BOUND |
| vtask-core |             | HOST-BOUND | HOST-BOUND |
| vtask-tail |             |            |            |

## 信号量设计

流控信号量：```slot_vtask_size:vtask_head:${slot_seq}```

编程信号量：```:slot_vtask_size:vtask_head:${slot_seq}```，与流控信号量初值一致，用于编程控制。

流控信号量的值范围在[0..n]；可编程信号量的值短期可出现-1.

## 1. default 

### 1.1 Create app

```sh
export CLUSTER=local
export TASK_DIST_MODE=
export HEAD_SLOTS=h0
export CORE_MODE=
export CORE_SLOTS=h0
app_id=$( cat default-tasks.txt | scalebox app run | cut -d':' -f2 | tr -d '}' )
```

### 1.2 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} vtask_size:vtask-head
```

## 2. host-bound

### 2.1 Create app

```sh
cd /shared/scalebox/examples/vtask

export CLUSTER=inline
export TASK_DIST_MODE=HOST-BOUND
export HEAD_SLOTS=n0-[01]
export CORE_MODE=HOST-BOUND
export CORE_SLOTS=n0-[01]
app_id=$( cat host-tasks.txt | scalebox app run | cut -d':' -f2 | tr -d '}' )
```

### 2.2 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} host_vtask_size:wait-queue

```

## 3. slot-bound

group-bound

### 3.1 Create app

```sh
cd /shared/scalebox/examples/vtask

export CLUSTER=inline
export TASK_DIST_MODE=SLOT-BOUND
export HEAD_SLOTS=h0:2
export CORE_MODE=HOST-BOUND
export CORE_SLOTS=n[01]-[01]
app_id=$( cat slot-tasks.txt | scalebox app run | cut -d':' -f2 | tr -d '}' )

```
### 3.2 add tasks

```sh
for i in {0..3}; do
  echo "$i"
  scalebox task add --app-id=$app_id --header to_slot_index=$i 0${i}0
  scalebox task add --app-id=$app_id --header to_slot_index=$((i+1)) 0${i}1
done

```

### 3.3 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} host_vtask_size:wait-queue
```
