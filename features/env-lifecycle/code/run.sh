#!/bin/bash
# 一次性写出，避免极短任务下 agent 逐行读取丢失尾部。
printf 'DEMO_VALUE=%s\nLITERAL_DEMO=%s\nAPP_ID=%s\n' \
    "$DEMO_VALUE" "$LITERAL_DEMO" "$APP_ID"
# 给 agent 的 stdout 读取线程留出时间（超短任务的已知竞态）。
sleep 0.3
exit 0
