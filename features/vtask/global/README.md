# global_vtask

Wait queue test for global vtask.


## Test 

### 1. Create app

```sh
cd features/vtask/global/test
app_id=$( scalebox app create | cut -d':' -f2 | tr -d '}' )
```

### 2. Add messages

```sh
APP_ID=${app_id} scalebox task add --sink-job global-vtask --task-file messages.txt
```

### 3. increment semaphore

```sh
APP_ID=${app_id} scalebox semaphore increment global_vtask_size:global-vtask
```
