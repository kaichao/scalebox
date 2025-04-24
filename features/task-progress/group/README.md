# task_progress_group_diff

Test for task_progress_group_diff.


## Test 

### 1. Create app

```sh
cd features/task-progress/group/test
app_id=$( scalebox app create | cut -d':' -f2 | tr -d '}' )
```

### 2. Add initial messages

```sh
APP_ID=${app_id} scalebox task add --sink-job task-progress-group --task-file messages.txt
```

### 3. Add new messages

```sh
APP_ID=${app_id} scalebox task add --sink-job task-progress-group --header to_host=n-01.inline 006

APP_ID=${app_id} scalebox task add --sink-job task-progress-group --header to_host=n-00.inline 007
```
