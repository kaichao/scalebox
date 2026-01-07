#!/bin/bash

source functions.sh

echo "In from-vtask-head.sh"

body=$1

if [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    from_ip=$(get_header "$2" "from_ip")
    scalebox task add --sink-module=vtask-core --header to_ip=$from_ip $body
    code=$?
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    # 根据body的值、头_vtask_size_sema，设置to_host(n0-0, n0-1, n1-0, n1-1)
    sema_name=$(get_header "$2" "_vtask_size_sema")
    group_part=${sema_name##*:}

    # body为3位整数，host_part 为body的整数除以2，再用2取模（结果为0或1）
    # 将body转换为整数，计算 (body / 2) % 2
    body_int=$((10#$body))  # 确保按十进制解析
    host_part=$(( (body_int / 2) % 2 ))

    to_host="n${group_part}-${host_part}"
    echo $to_host
    scalebox task add --sink-module=vtask-core --header to_host=$to_host $body
    code=$?
else
    # TASK_DIST_MODE == ""
    # echo "vtask-core,$body" > $WORK_DIR/sink-tasks.txt
    scalebox task add --sink-module=vtask-core $body
    exit $?
fi

if [ $code -eq 0 ]; then
    # 增加信号量，打开wait-queue的开关
    echo  "semaphore increment vtask_size:wait-queue"
    scalebox semaphore increment "vtask_size:wait-queue"
    code=$?
fi

exit $code
