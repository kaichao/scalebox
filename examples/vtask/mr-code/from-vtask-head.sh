#!/bin/bash

source functions.sh

echo "In from-vtask-head.sh"

body=$1

if [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    # 参照default.sh，根据body十位数、个位数的奇偶性，设置to_host(n0-0, n0-1, n1-0, n1-1)
    
    # 提取十位数（从右往左第二位），长度不足时默认为0
    ten_digit=0
    [ ${#body} -ge 2 ] && ten_digit=${body: -2:1}
    
    # 提取个位数（最后一位），长度不足时默认为0
    unit_digit=0
    [ ${#body} -ge 1 ] && unit_digit=${body: -1:1}
    
    # 十位数为偶数则host_part="n0"，否则为"n1"
    # 个位数为偶数则slot_part="0"，否则为"1"
    if [ $((ten_digit % 2)) -eq 0 ]; then
        host_part="n0"
    else
        host_part="n1"
    fi
    
    if [ $((unit_digit % 2)) -eq 0 ]; then
        slot_part="0"
    else
        slot_part="1"
    fi
    
    to_host="${host_part}-${slot_part}"
    echo $to_host
    scalebox task add --sink-module=vtask-core --header to_host=$to_host $body
    code=$?
    if [ $code -eq 0 ]; then
        # 增加信号量，打开wait-queue的开关 
        echo  "semaphore increment vtask_size:wait-queue"
        scalebox semaphore increment "vtask_size:wait-queue"
        code=$?
    fi

elif [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    from_ip=$(get_header "$2" "from_ip")
    scalebox task add --sink-module=vtask-core --header to_ip=$from_ip $body
    code=$?
else
    scalebox task add --sink-module=vtask-core $body
    code=$?
    # echo "vtask-core,$body" > $WORK_DIR/sink-tasks.txt
fi

exit $code
