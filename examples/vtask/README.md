# vtask

vtask是同一应用中跨模块的task集合，还包含信号量、共享变量的集合等，vtask是节点本地计算编程模型的基本单元，通过参数设计，简化本地计算在应用程序中的准入控制、容错等的实现；

主要功能：
- 任务管理：跨模块的任务（task）组成vtask，简化任务管理、状态管理
- 计数准入控制：以配合HPC计算资源的动态调度，实现超长任务计算；
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
cd scalebox/examples/vtask

cat default-tasks.txt | scalebox run
```

### 1.2 semaphore ls

- 以当前最新app-id
```sh
scalebox semaphore ls vtask_size
```

## 2. host-bound

### 2.1 Create app

```sh
cd scalebox/examples/vtask

app_id=$( cat host-tasks.txt | scalebox run -e host-bound.env| cut -d':' -f2 | tr -d '}' )
```

### 2.2 semaphore ls

```sh
scalebox semaphore ls host_vtask_size:vtask-head
scalebox semaphore ls vtask_size
```

## 3. slot-bound

group-bound

### 3.1 Create app

```sh
cd scalebox/examples/vtask

app_id=$( cat slot-tasks.txt | scalebox run -e slot-bound.env| cut -d':' -f2 | tr -d '}' )

```
### 3.2 add tasks

```sh
for i in {0..3}; do
  echo "$i"
  scalebox task add --app-id=$app_id --header to_slot_index=$i 0${i}0
  scalebox task add --app-id=$app_id --header to_slot_index=$((i+1)) 0${i}1
done

```

### 3.3 semaphore ls

```sh
scalebox semaphore increment --app-id=${app_id} host_vtask_size:wait-queue
```
