#!/bin/bash

# 入口路由：仅处理无 from_module 的入口 task（由 run.sh 的空 from_module 分支调用）
source /usr/local/lib/scalebox/functions.sh

echo "In default.sh"
body=$1
headers=$2

# VTASK_MODE：vtask 类型（DEFAULT / HOST-BOUND / GROUP-BOUND）
# HOST-BOUND / GROUP-BOUND 需经 wait-queue 分配资源，DEFAULT 直连 vtask-head
if [ "$VTASK_MODE" = "HOST-BOUND" ] || [ "$VTASK_MODE" = "GROUP-BOUND" ]; then
    # 应用属性（_vtask_attr_*）显式转发给 wait-queue：wait-queue 无 vtask_role，
    # agent 对 router→wait-queue 这一跳的 _ 前缀自动继承会被门控拦下
    #（cmd/agent/sink_tasks.go 的 router 分支）；标准模块 run.sh 再转发给 vtask-head。
    declare -A in_headers
    if [ -n "${headers}" ] && [ "${headers}" != "{}" ]; then
        scalebox::json_parse "${headers}" in_headers
    fi
    attr_headers=()
    for k in "${!in_headers[@]}"; do
        if [[ "$k" == _vtask_attr_* ]]; then
            attr_headers+=(--header "$k=${in_headers[$k]}")
        fi
    done
    scalebox task add --sink-module=wait-queue "${attr_headers[@]}" $body
else
    # DEFAULT：sink 是 vtask-head（vtask_role=head），_vtask_attr_* 由 agent 自动继承
    scalebox task add --sink-module=vtask-head $body
fi

exit $?
