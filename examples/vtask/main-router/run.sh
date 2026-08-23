#!/bin/bash

code_dir=$(dirname "$0")

# headers='{"_vtask_id":"36","from_ip":"10.0.6.100","from_module":"vtask-head","to_slot":"27"}'
headers=$2
echo "headers:$headers"

pattern='"from_module":"([^"]+)"'
if [[ $headers =~ $pattern ]]; then
    from_module="${BASH_REMATCH[1]}"
else
    # no from_module in json 
    from_module=""
fi
echo "from_module:$from_module"

case $from_module in
    "vtask-head")
        "${code_dir}/from-vtask-head.sh" "$1" "$2"
        ;;
    "vtask-core")
        "${code_dir}/from-vtask-core.sh" "$1" "$2"
        ;;
    "vtask-tail")
        "${code_dir}/from-vtask-tail.sh" "$1" "$2"
        ;;
    "wait-queue")
        # 标准模块 run.sh 已完成 bind + add-subtask 且不回流；
        # 若因自定义 run.sh 等原因回流，空操作终结，避免落入 default.sh 重新入队造成死循环
        exit 0
        ;;
    "")
        # 仅入口 task（无 from_module）走入口路由
        "${code_dir}/default.sh" "$1" "$2"
        ;;
    *)
        echo "unexpected from_module:${from_module}, drop task" >&2
        exit 0
        ;;
esac

exit $?
