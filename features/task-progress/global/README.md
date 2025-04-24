# task_progress_global_diff

Test for task_progress_global_diff.


## Test 

### 1. Create app

```sh
cd features/task-progress/global/test
app_id=$( scalebox app create | cut -d':' -f2 | tr -d '}' )
```

### 2. Add initial messages

```sh
APP_ID=${app_id} scalebox task add --sink-job task-progress-global --task-file messages.txt
```

### 3. Add new messages

```sh
APP_ID=${app_id} scalebox task add --sink-job task-progress-global --header to_host=n-02.inline 022
APP_ID=${app_id} scalebox task add --sink-job task-progress-global --header to_host=n-02.inline 023

APP_ID=${app_id} scalebox task add --sink-job task-progress-global --header to_host=n-03.inline 030
APP_ID=${app_id} scalebox task add --sink-job task-progress-global --header to_host=n-03.inline 031
```
