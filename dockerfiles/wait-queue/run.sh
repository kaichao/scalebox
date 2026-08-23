#!/bin/bash

# wait-queue 标准模块：vtask 资源分配 + 转发（slot 主逻辑）
# 执行位置：wait-queue slot（拾取来自 main-router 的 task）
# 完成后不写 sink-tasks.txt，task 在本 slot 终结（vtask-head 任务由 add-subtask --direct 直投）
#
# 环境变量：
#   VTASK_MODE       vtask 类型：HOST-BOUND / GROUP-BOUND（缺省 HOST-BOUND）
#   HEAD_MODULE      vtask-head 模块名（缺省 vtask-head）
#   PLAT_MODULE_NAME 当前模块名（agent 注入，用于 gate 信号量）

body=$1
head_module=${HEAD_MODULE:-vtask-head}
module_name=${PLAT_MODULE_NAME:-wait-queue}

if [ "$VTASK_MODE" = "GROUP-BOUND" ]; then
    # GROUP-BOUND vtask：bind 返回组号（head slot 的 seq）
    group_seq=$(scalebox vtask bind) || exit 1
    sema_name="group_vtask_size:${head_module}:${group_seq}"
    scalebox vtask add-subtask --direct --module=$head_module \
        --header to_slot_index=$group_seq --header _vtask_size_sema=$sema_name $body || {
        scalebox vtask unbind --sema-name="${sema_name}"
        scalebox semaphore increment vtask_size:${module_name}
        exit 1
    }
else
    # HOST-BOUND vtask（DEFAULT 不经 wait-queue）
    to_host=$(scalebox vtask bind) || exit 1
    sema_name="host_vtask_size:${head_module}:${to_host}"
    scalebox vtask add-subtask --direct --module=$head_module \
        --header to_host=$to_host --header _vtask_size_sema=$sema_name $body || {
        # 回滚已扣减的资源，并放行下一个 vtask
        scalebox vtask unbind --sema-name="${sema_name}"
        scalebox semaphore increment vtask_size:${module_name}
        exit 1
    }
fi

exit 0
