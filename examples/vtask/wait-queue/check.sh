#!/bin/bash

# wait-queue 的 slot 级准入控制（agent 在 GetTaskList 之前执行 ACTION_CHECK）
# 信号量组 max > 0 才放行领取，额度耗尽时返回非 0，agent 停领并轮询重试

if [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    sema_prefix='^host_vtask_size:vtask-head'
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    sema_prefix='^slot_vtask_size:vtask-head'
else
    exit 0
fi

# 检查max值
ret=$(scalebox semagroup max "$sema_prefix")
# ret格式：'"name":value' 或 '":name:subname":value'
#   "host_vtask_size:vtask-head:n0-0":2
#   "slot_vtask_size:vtask-head:1":3
code=$?

if [ $code -ne 0 ]; then
    echo "semagroup max error, code=$code" >&2
    exit $code
fi

max_value=${ret##*:}
# 确保max_value为整数，进行数值比较
if [ "$max_value" -le 0 ]; then
    # 最大值小于等于0，直接忽略，返回1
    exit 1
fi

exit 0
