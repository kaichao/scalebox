#!/bin/bash

# wait-queue 标准模块：slot 级准入控制（agent 在 GetTaskList 之前执行 ACTION_CHECK）
# 信号量组 max > 0 才放行领取，额度耗尽时返回非 0，agent 停领并轮询重试

head_module=${HEAD_MODULE:-vtask-head}

# VTASK_MODE：vtask 类型（HOST-BOUND / GROUP-BOUND）
# DEFAULT 不经 wait-queue；缺省按 HOST-BOUND 处理
if [ "$VTASK_MODE" = "GROUP-BOUND" ]; then
    sema_prefix="^group_vtask_size:${head_module}"
else
    sema_prefix="^host_vtask_size:${head_module}"
fi

# 检查 max 值
ret=$(scalebox semagroup max "$sema_prefix")
# ret 格式：'"name":value' 或 '":name:subname":value'
#   "host_vtask_size:vtask-head:n0-0":2
#   "group_vtask_size:vtask-head:1":3
code=$?

if [ $code -ne 0 ]; then
    echo "semagroup max error, code=$code" >&2
    exit $code
fi

max_value=${ret##*:}
# 确保 max_value 为整数，进行数值比较
if [ "$max_value" -le 0 ]; then
    # 最大值小于等于 0，直接忽略，返回 1
    exit 1
fi

exit 0
