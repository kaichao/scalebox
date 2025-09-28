# vtask_size


## 1. default 

### 1.1 Create app

```sh
export CLUSTER=local
export HOSTS=h0
export TASK_DIST_MODE=
app_id=$( cat default-tasks.txt | scalebox app run | cut -d':' -f2 | tr -d '}' )
```

### 1.2 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} vtask_size:vtask
```

## 2. host-bound

### 2.1 Create app

```sh
export CLUSTER=inline
export HOSTS=n-0[01]
export TASK_DIST_MODE=HOST-BOUND
app_id=$( cat host-tasks.txt | scalebox app run | cut -d':' -f2 | tr -d '}' )
```

### 2.2 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} vtask_size:vtask:n-00

scalebox semaphore increment --app-id=${app_id} vtask_size:vtask:n-01
```

## 3. slot-bound

### 3.1 Create app

```sh
export CLUSTER=local
export HOSTS=h0:2
export TASK_DIST_MODE=SLOT-BOUND
app_id=$( scalebox app run | cut -d':' -f2 | tr -d '}' )
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
### 3.3 increment semaphore

```sh
scalebox semaphore increment --app-id=${app_id} vtask_size:vtask:$slot_id

scalebox semaphore increment --app-id=${app_id} vtask_size:vtask:$((slot_id+1))
```
