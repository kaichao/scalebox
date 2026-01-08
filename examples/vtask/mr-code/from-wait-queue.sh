#!/bin/bash

# only "$TASK_DIST_MODE" == "SLOT-BOUND" OR "HOST-BOUND"
body=$1

# 函数：从semagroup获取子名称
# 输入：sema_prefix
# 输出：提取的最后一个冒号后的值（通过标准输出返回）
# 退出码：0-成功，1-最大值小于等于0（忽略），其他-失败
get_semagroup_subname() {
    local sema_prefix="$1"
    
    # 检查max值
    local ret=$(scalebox semagroup max "$sema_prefix")
    # ret格式：'"name":value' 或 '":name:subname":value'
    #   ":host_vtask_size:vtask-head:n0-0":2
    #   ":slot_vtask_size:vtask-head:1":3
    local code=$?
    if [ $code -ne 0 ]; then
        echo "semagroup max error, code=$code" >&2
        return $code
    fi

    local max_value=${ret##*:}
    # max_value=${max_value//\"/}
    # 确保max_value为整数，进行数值比较
    if [ "$max_value" -le 0 ]; then
        # 最大值小于等于0，直接忽略，返回1
        return 1
    fi

    # 执行decrement
    ret=$(scalebox semagroup decrement "$sema_prefix")
    code=$?
    if [ $code -ne 0 ]; then
        echo "semagroup decrement error, code=$code" >&2
        return $code
    fi

    # 提取引号内的内容
    if [[ $ret =~ \"([^\"]+)\" ]]; then
        local sema_name=${BASH_REMATCH[1]}
        # 提取最后一个冒号后的值
        echo "${sema_name##*:}"
        return 0
    else
        echo "Invalid format of semagroup decrement, ret_value=$ret" >&2
        return 107
    fi
}

if [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    sema_prefix=':host_vtask_size:vtask-head'
    to_host=$(get_semagroup_subname "$sema_prefix")
    if [ $? -eq 0 ]; then
        scalebox task add --sink-module=vtask-head --header to_host=$to_host $body
    fi
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    sema_prefix=':slot_vtask_size:vtask-head'
    to_slot_index=$(get_semagroup_subname "$sema_prefix")
    if [ $? -eq 0 ]; then
        scalebox task add --sink-module=vtask-head --header to_slot_index=$to_slot_index $body
    fi
else
    echo "Invalid TASK_DIST_MODE:$TASK_DIST_MODE" >&2
    exit 108
fi

exit 0
