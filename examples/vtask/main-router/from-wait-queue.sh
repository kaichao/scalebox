#!/bin/bash

# vtask 资源分配 + 转发
# APP_ID 由 agent 容器环境变量设定，无需 --app-id 传参

body=$1

if [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    to_host=$(scalebox vtask bind) || exit 1
    sema_name="host_vtask_size:vtask-head:${to_host}"
    scalebox vtask add-subtask --direct --module=vtask-head \
        --header to_host=$to_host --header _vtask_size_sema=$sema_name $body || {
        scalebox vtask unbind --sema-name=":${sema_name}"
        scalebox semaphore increment vtask_size:wait-queue
        exit 1
    }
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    to_slot_index=$(scalebox vtask bind) || exit 1
    sema_name="slot_vtask_size:vtask-head:${to_slot_index}"
    scalebox vtask add-subtask --direct --module=vtask-head \
        --header to_slot_index=$to_slot_index --header _vtask_size_sema=$sema_name $body || {
        scalebox vtask unbind --sema-name=":${sema_name}"
        scalebox semaphore increment vtask_size:wait-queue
        exit 1
    }
else
    echo "Invalid TASK_DIST_MODE:$TASK_DIST_MODE" >&2
    exit 108
fi

exit 0
