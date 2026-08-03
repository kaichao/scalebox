#!/bin/bash

source /usr/local/lib/scalebox/functions.sh

sema=$(scalebox::task_header "$2" "_vtask_size_sema")
# 仅 HOST-BOUND / SLOT-BOUND 需要 unbind（有 bind 才有 unbind）
# DEFAULT 的 _vtask_size_sema=vtask_size:vtask-head 不触发 unbind
case "$sema" in
    host_vtask_size:*|slot_vtask_size:*)
        scalebox vtask unbind --sema-name=":${sema}"
        code=$?
        echo "exit-code of vtask-unbind is $code"
        ;;
esac

exit 0
