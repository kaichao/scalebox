#!/bin/bash

source functions.sh

if [ "$TASK_DIST_MODE" = "SLOT-BOUND" ] || [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    # vtask处理完成，增加信号量，恢复编程级vtask流控信号量
    sema=$(get_header "$2" "_vtask_size_sema")
    sema=":${sema}"
    echo  "semaphore increment ${sema}"
    scalebox semaphore increment "${sema}"
    # 当前未正确返回错误码
    # code=$?
    # exit $code
fi

exit 0
