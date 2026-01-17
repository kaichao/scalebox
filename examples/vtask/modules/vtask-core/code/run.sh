#!/bin/bash

if [ "${RUNNING_SECONDS}" -gt 0 ]; then
    sleep_seconds="${RUNNING_SECONDS}"
else
    # 生成15到120之间的随机数
    sleep_seconds=$(( RANDOM % 106 + 15 ))  # 106 = 120 - 15 + 1
fi
sleep "$sleep_seconds"

scalebox task add $1
code=$?

echo "sleep_seconds:$sleep_seconds, exit_code=$code" >> "$WORK_DIR/auxout.txt"

exit $code
