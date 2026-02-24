#!/bin/bash

# set -euo pipefail

code_dir=$(dirname $0)

echo "num-args:$#"
echo "args:$*"

# headers='{"from_ip":"10.11.16.80","from_module":"dir-list","from_host":"n1.dcu"}'
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

num_groups=${NUM_GROUPS:-10}

case $from_module in
    "my-module")
        # sum part result
        if [ -z "${TASK_BATCH_SIZE}" ]; then
            "${code_dir}/sum.sh" "$1" "$2"
        else
            "${code_dir}/sum-batch.sh" "$1" "$2"
        fi

        ;;
    *)  
        # null from_module
        if [ -z "${TASK_DIST_MODE}" ]; then
            "${code_dir}/split.sh" "$1" "${num_groups}"
        else
            "${code_dir}/split-host-bound.sh" "$1" "$2"
        fi
        ;;
esac

exit $?
