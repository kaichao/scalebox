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
        "${code_dir}/from-vtask-head.sh" "$1"
        ;;
    "vtask-core")
        "${code_dir}/from-vtask-core.sh" "$1"
        ;;
    "vtask-tail")
        "${code_dir}/from-vtask-tail.sh" "$1"
        ;;
    *)
        "${code_dir}/default.sh" "$1"
        ;;
esac

exit $?
