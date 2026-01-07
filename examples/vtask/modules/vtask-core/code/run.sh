#!/bin/bash

if [ "${RUNNING_SECONDS}" -gt 0 ]; then
    # 生成15到120之间的随机数
    sleep_seconds=$(( RANDOM % 106 + 15 ))  # 106 = 120 - 15 + 1
    sleep $sleep_seconds
else
    sleep "${RUNNING_SECONDS}"
fi

scalebox task add $1

exit $?
