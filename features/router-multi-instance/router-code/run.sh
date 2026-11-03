#!/bin/bash
# 最小主路由：按来源模块把回流任务投给对应的下游。
# 两个实例跑同一份代码——分派由 from_module 决定，实例身份只影响
# 任务被哪个 slot 领走（to_slot 由 controld 按 router_index 写入）。

headers=$2
from_module=""
if [[ $headers =~ \"from_module\":\"([^\"]+)\" ]]; then
    from_module="${BASH_REMATCH[1]}"
fi

# 记录本实例实际收到的来源，便于 test.sh 断言
echo "router instance received: from_module=${from_module} body=$1" >> /tmp/router-multi-instance.log

case "$from_module" in
    "src-a")
        echo "sink-a,$1" > "${WORK_DIR}/sink-tasks.txt"
        ;;
    "src-b")
        echo "sink-b,$1" > "${WORK_DIR}/sink-tasks.txt"
        ;;
    *)
        # 其他来源（如入口）不投下游
        ;;
esac
