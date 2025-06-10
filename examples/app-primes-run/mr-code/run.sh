#!/bin/bash

# set -euo pipefail

code_dir=$(dirname $0)

echo "num-args:$#"
echo "args:$*"

# headers='{"from_ip":"10.11.16.80","from_job":"dir-list","from_host":"n1.dcu"}'
headers=$2
echo "headers:$headers"

pattern='"from_job":"([^"]+)"'
if [[ $headers =~ $pattern ]]; then
    from_job="${BASH_REMATCH[1]}"
else
    # no from_job in json 
    from_job=""
fi

echo "from_job:$from_job"

env

case $from_job in
    "my-job")
        # 未定义或为空
        if [ -z "${BULK_MESSAGE_SIZE}" ]; then
            echo 001
            "${code_dir}/sum.sh" "$1" "$2"
        else
            echo 000
            "${code_dir}/sum-bulk.sh" "$1" "$2"
        fi

        ;;
    *)  # null from_job
        if [ -z "${TASK_DIST_MODE}" ]; then
            echo 101
            "${code_dir}/split.sh" "$1" "$2"
        else
            echo 100
            "${code_dir}/split-host-bound.sh" "$1" "$2"
        fi
        ;;
esac

exit $?
