#!/bin/bash

headers=$2
pattern='"from_module":"([^"]+)"'
if [[ $headers =~ $pattern ]]; then
    from_module="${BASH_REMATCH[1]}"
else
    # no from_module in json 
    from_module=""
fi

if [ "$from_module" == "" ]; then
    echo "my-module,$1" > "${WORK_DIR}/sink-tasks.txt"
fi

exit 0
