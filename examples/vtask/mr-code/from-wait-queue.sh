#!/bin/bash

# only "$TASK_DIST_MODE" == "SLOT-BOUND" OR "HOST-BOUND"
body=$1

if [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    ret=$(scalebox semagroup decrement ':host_vtask_size:vtask-head')
    code=$?
    if [ $code -ne 0 ] ; then
        echo "semagroup decrement error, code=$code" >&2
        exit $code
    fi

    # ret == ":host_vtask_size:vtask-head:h0-1":3
    # 使用正则表达式提取引号内的hostname
    if [[ $ret =~ \"([^\"]+)\" ]]; then
        sema_name=${BASH_REMATCH[1]}
        # 提取引号内最后一个冒号后的host名称
        to_host=${sema_name##*:}
    else
        echo "Invalid format of semagroup decrement, ret_value=$ret" >&2
        exit 107
    fi

    scalebox task add --sink-module=vtask-head --header to_host=$to_host $body
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    ret=$(scalebox semagroup decrement ':slot_vtask_size:vtask-head')
    code=$?
    echo "ret_value of semagroup decrement:$ret" >> /work/custom-out.txt
    if [ $code -ne 0 ] ; then
        echo "semagroup decrement error, code=$code" >&2
        exit $code
    fi

    # ret == ":slot_vtask_size:vtask-head:1":3
    # 使用正则表达式提取引号内slot_seq
    if [[ $ret =~ \"([^\"]+)\" ]]; then
        sema_name=${BASH_REMATCH[1]}
        # 提取引号内最后一个冒号后的数字
        to_slot_index=${sema_name##*:}
    else
        echo "Invalid format of semagroup decrement, ret_value=$ret" >&2
        exit 107
    fi

    scalebox task add --sink-module=vtask-head --header to_slot_index=$to_slot_index $body
else
    echo "Invalid TASK_DIST_MODE:$TASK_DIST_MODE" >&2
    exit 108
fi

exit $?
