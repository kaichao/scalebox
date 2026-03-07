#!/bin/bash

code_dir=$(dirname "$0")

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
    "mod-1")
        echo "mod-2,$1" > ${WORK_DIR}/sink-tasks.txt
        ;;
    "mod-2")
        # do nothing
        ;;
    *)
        n=$1
        for ((i=0; i<n; i++)); do
            echo "$i,$((i % 2))"
            scalebox task add --sink-module=mod-1 --header to_host=n0-$((i % 2)) $i
        done
        ;;
esac

exit 0
