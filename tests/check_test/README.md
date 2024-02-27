# check

## dir_limit_gb

Standard flow control parameters, used to limit the maximum number of Gigabytes in the directory

```sh
# FALSE
DIR_LIMIT_GB=/data/ssd/tmp~18 scalebox app create

# TRUE
DIR_LIMIT_GB=/data/ssd/tmp~22 scalebox app create
```

目录最大占用空间19GB

## dir_free_gb

Standard flow control parameters, used to limit the minimum number of Gigabytes of the partition where the directory is located

```sh
# TRUE
DIR_FREE_GB=/data/ssd/tmp~1000 scalebox app create

# FALSE
DIR_FREE_GB=/data/ssd/tmp~1500 scalebox app create
```


DIR_LIMIT_GB=/data/ssd/tmp~22 DIR_FREE_GB=/data/ssd/tmp~1000 ACTION_CHECK=/app/bin/mycheck.sh scalebox app create

目录所在分区需至少要保留1000GB的空余空间

## user-defined flow control
```sh
ACTION_CHECK=/app/bin/mycheck.sh scalebox app create
```
