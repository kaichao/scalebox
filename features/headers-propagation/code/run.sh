#!/bin/bash
# $1 = task body, $2 = headers JSON
source /usr/local/lib/scalebox/functions.sh

case "${PLAT_MODULE_NAME}" in
    producer)
        # 非路由模块写 sink：_ 前缀 header 由 agent 自动继承
        echo "consumer,$1" > "${WORK_DIR}/sink-tasks.txt"
        ;;
    consumer)
        demo=$(scalebox::task_header "$2" _demo)
        echo "_demo=${demo}"
        ;;
esac
exit 0
