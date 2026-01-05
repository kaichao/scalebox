#!/bin/bash

# only "$TASK_DIST_MODE" = "SLOT-BOUND" allowed

body=$1

ret=$(scalebox semagroup decrement 'slot_vtask_size:vtask-head')
code=$?

if [ $code -ne 0 ] ; then
    echo "semagroup decrement error, code=$code" >&2
    exit $code
fi
# ret == "slot_vtask_size:vtask-head:1":3
# 使用正则表达式提取引号内的内容
if [[ $ret =~ \"([^\"]+)\" ]]; then
    sema_name=${BASH_REMATCH[1]}
    # 提取引号内最后一个冒号后的数字
    to_slot_index=${sema_name##*:}
else
    echo "invalid format of semagroup decrement, ret_value=$ret"
    to_slot_index=0  # 默认值
fi

scalebox task add --sink-module=vtask-head --header to_slot_index=$to_slot_index $body

exit $?
