# vtask

主要功能：
- 粗粒度容错；
- 计数流控，以配合HPC动态资源分配；
- 节点本地计算的组织；

包括：
- wait-queue：仅用于组模式的定制流控
- vtask-head：vtask起始标识模块。组模式下，多个slot需部署在单节点上。
- vtask-core：vtask核心处理模块，可为多个；
- vtask-tail：vtask结束标识模块。task的主体（body）与vtask-head一致。

task_dist_mode设置

|            |   全局模式   |   节点模式   |   组模式    |
| ---------- | ----------- | ---------- | ---------- |
| wait-queue |             |            |            |
| vtask-head |             | HOST-BOUND | SLOT-BOUND |
| vtask-core |             | HOST-BOUND | HOST-BOUND |
| vtask-tail |             |            |            |


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
export CLUSTER=inline
export TASK_DIST_MODE=HOST-BOUND
export HEAD_SLOTS=n0-[01]
export CORE_MODE=HOST-BOUND
export CORE_SLOTS=n0-[01]
app_id=$( cat host-tasks.txt | scalebox app run | cut -d':' -f2 | tr -d '}' )
```

### 2.2 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} host_vtask_size:vtask-head:n-00

scalebox semaphore increment --app-id=${app_id} host_vtask_size:vtask-head:n-01
```

## 3. slot-bound

group-bound

### 3.1 Create app

```sh

export CLUSTER=inline
export TASK_DIST_MODE=SLOT-BOUND
export HEAD_SLOTS=h0:2
export CORE_MODE=HOST-BOUND
export CORE_SLOTS=n[01]-[01]
export VTASK_AUTO_COUNT=no
app_id=$( cat slot-tasks.txt | scalebox app run | cut -d':' -f2 | tr -d '}' )

```
### 3.2 add tasks

```sh
export slot_id=15
for i in {0..3}; do
  echo "$i"
  scalebox task add --app-id=$app_id --header to_slot=$slot_id 0${i}0
  scalebox task add --app-id=$app_id --header to_slot=$((slot_id+1)) 0${i}1
done

```
或
```sh
for i in {0..3}; do
  echo "$i"
  scalebox task add --app-id=$app_id --header to_slot_index=$i 0${i}0
  scalebox task add --app-id=$app_id --header to_slot_index=$((i+1)) 0${i}1
done

```

### 3.3 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} vtask_size:vtask-head:$slot_id

scalebox semaphore increment --app-id=${app_id} vtask_size:vtask-head:$((slot_id+1))
```
