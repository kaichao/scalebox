#!/bin/bash

source /usr/local/lib/scalebox/functions.sh

echo "In from-vtask-head.sh"

body=$1

# _vtask_id / _vtask_size_sema 由 agent addSinkTasks 自动传播，不再需要显式传递
sema_name=$(scalebox::task_header "$2" "_vtask_size_sema")

if [ "$VTASK_MODE" = "HOST-BOUND" ]; then
    # 不传路由 header：controld 优先用 wait-queue 写入的 _vtask_host 兜底；
    # 老路径/非 vtask 任务可用 vtask-head 与 vtask-core 相同的 pod_id 兜底。
    scalebox vtask add-subtask --module=vtask-core $body
elif [ "$VTASK_MODE" = "GROUP-BOUND" ]; then
    # 组号优先取 _vtask_group；老 wait-queue 回退 _vtask_size_sema
    group_part=$(scalebox::task_header "$2" "_vtask_group")
    if [ -z "$group_part" ]; then
        group_part=${sema_name##*:}
    fi
    # 组内跨 host 分发是业务逻辑，必须显式给 to_host（平台不猜组内节点）；
    # async 路径由 agent 传播 _vtask_*；若改用 --direct，须显式转发
    # _vtask_id/_vtask_group/_vtask_head_module。
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
