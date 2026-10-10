#!/bin/bash

# vtask 资源额度已由 controld 在 vtask-tail 任务完成时（doVTaskFinished）
# 自动加回（读 tail task 的 _vtask_size_sema），脚本侧无需再调用 vtask unbind

# 示例：应用属性（_vtask_attr_*）随 vtask 管道传至 tail，打印到 stdout，
# 用 `scalebox task show <tail-task-id>` 可见（Headers 里也含原始 header）
source /usr/local/lib/scalebox/functions.sh
headers_json=${2:-}
declare -A task_headers
if [ -n "${headers_json}" ] && [ "${headers_json}" != "{}" ]; then
    scalebox::json_parse "${headers_json}" task_headers
fi
for k in "${!task_headers[@]}"; do
    if [[ "$k" == _vtask_attr_* ]]; then
        echo "[vtask-tail] ${k}=${task_headers[$k]}"
    fi
done

exit 0
