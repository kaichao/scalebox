#!/bin/bash

# 入口路由：仅处理无 from_module 的入口 task（由 run.sh 的空 from_module 分支调用）
echo "In default.sh"
body=$1

# VTASK_MODE：vtask 类型（DEFAULT / HOST-BOUND / GROUP-BOUND）
# HOST-BOUND / GROUP-BOUND 需经 wait-queue 分配资源，DEFAULT 直连 vtask-head
if [ "$VTASK_MODE" = "HOST-BOUND" ] || [ "$VTASK_MODE" = "GROUP-BOUND" ]; then
    scalebox task add --sink-module=wait-queue $body
else
    scalebox task add --sink-module=vtask-head $body
fi

exit $?
