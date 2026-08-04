#!/bin/bash

source /usr/local/lib/scalebox/functions.sh

echo "In from-vtask-head.sh"

body=$1

# _vtask_id / _vtask_size_sema 由 agent addSinkTasks 自动传播，不再需要显式传递
sema_name=$(scalebox::task_header "$2" "_vtask_size_sema")

if [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    from_ip=$(scalebox::task_header "$2" "from_ip")
    scalebox vtask add-subtask --module=vtask-core \
        --header to_ip=$from_ip $body
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    # to_host 按 body 哈希分配给计算节点（业务逻辑，保留在脚本中）
    group_part=${sema_name##*:}
    body_int=$((10#$body))
    host_part=$(( (body_int / 2) % 2 ))
    to_host="n${group_part}-${host_part}"
    echo "$to_host"
    scalebox vtask add-subtask --module=vtask-core \
        --header to_host=$to_host $body
else
    # DEFAULT 模式
    scalebox vtask add-subtask --module=vtask-core \
        $body
    exit $?
fi
code=$?

if [ $code -eq 0 ]; then
    # wait-queue vtask_size=1 保持串行化，放行下一个
    scalebox semaphore increment "vtask_size:wait-queue"
    code=$?
fi

exit $code
