#!/bin/bash

source /usr/local/lib/scalebox/functions.sh

echo "In from-vtask-head.sh"

body=$1

# _vtask_id / _vtask_size_sema 由 agent addSinkTasks 自动传播，不再需要显式传递
sema_name=$(scalebox::task_header "$2" "_vtask_size_sema")

if [ "$VTASK_MODE" = "HOST-BOUND" ]; then
    # to_host 由 pod_id 自动推导：vtask-head 与 vtask-core 配置相同的 pod_id，
    # 且未显式指定 to_host/to_ip 时，controld 按来源节点（from_ip）解析 to_host。
    scalebox vtask add-subtask --module=vtask-core $body
elif [ "$VTASK_MODE" = "GROUP-BOUND" ]; then
    # 组号取自 _vtask_size_sema（group_vtask_size:vtask-head:<组号>）
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

# 2026-11-02（L1）：串行门改由 wait-queue/run.sh 在 bind + 投递成功后立即归还，
# 这里不再 +1——否则双重归还，vtask_size:wait-queue 只涨不落。

exit $code
