#!/bin/bash

echo "In default.sh"
body=$1

if [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    # echo "wait-queue,$1" > $WORK_DIR/sink-tasks.txt
    scalebox task add --sink-module=wait-queue $body
elif [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    # 提取个位数（最后一位），长度不足时默认为0
    unit_digit=0
    [ ${#body} -ge 1 ] && unit_digit=${body: -1:1}

    # 个位数为奇数则to_host="n0-1"，否则为"n0-0"
    if [ $((unit_digit % 2)) -eq 1 ]; then
        to_host="n0-1.inline"
    else
        to_host="n0-0.inline"
    fi
    scalebox task add --sink-module=vtask-head --header to_host=$to_host $body
else
    scalebox task add --sink-module=vtask-head $body
fi

exit $?
