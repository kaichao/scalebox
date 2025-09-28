# node_progress_gap

Test for node_progress_gap.


## Test 

### 1. Create app and add initial tasks

```sh
cd features/node_progress_gap/
app_id=$(cat tasks.txt | scalebox app run| cut -d':' -f2 | tr -d '}' )
```

### 2. Add new tasks

```sh

scalebox task add --app-id=${app_id} --header to_host=n-00.inline 00e

APP_ID=${app_id} scalebox task add --header to_host=n-10.inline 10c
APP_ID=${app_id} scalebox task add --header to_host=n-10.inline 10d

APP_ID=${app_id} scalebox task add --header to_host=n-11.inline 11b
```
