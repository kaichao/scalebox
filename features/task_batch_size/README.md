# task_batch_size


```sh
export TASK_BATCH_SIZE=3
for i in {0..9}; do printf "%03d\n" $i; done | scalebox run
```

- 单个task_exec对应多个task
