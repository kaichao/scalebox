# host_vtask

Wait queue test for host vtask.


## Test 

### 1. Create app

```sh
cd features/vtask/host/test
app_id=$( scalebox app create | cut -d':' -f2 | tr -d '}' )
```

### 2. Add messages

```sh
APP_ID=${app_id} scalebox task add --sink-job host-vtask --task-file messages.txt
```

### 3. increment semaphore

```sh
APP_ID=${app_id} scalebox semaphore increment host_vtask_size:host-vtask:n-00

APP_ID=${app_id} scalebox semaphore increment host_vtask_size:host-vtask:n-01
```
