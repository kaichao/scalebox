#!/bin/bash

arr=($(/app/bin/url_parser ${SOURCE_URL}))

export MODE="${arr[0]}"
if [ "${MODE}" = "ERROR_FORMAT" ]; then
    echo "url format error, "${arr[1]} >&2
    exit 1
fi

if [ "${MODE}" = "LOCAL" ]; then
	#  MODE | LOCAL_ROOT
    export LOCAL_ROOT="${arr[1]}"
elif [ "${MODE}" = "FTP" ]; then
    export FTP_URL="${arr[1]}"
    export REMOTE_ROOT="${arr[2]}"
else
    # MODE | REMOTE_HOST | REMOTE_PORT |  REMOTE_ROOT | REMOTE_USER 
    export REMOTE_HOST="${arr[1]}"
    export REMOTE_PORT="${arr[2]}"
    export REMOTE_ROOT="${arr[3]}"
    if [ ${#arr[*]} -eq 5 ]; then
        export REMOTE_USER="${arr[4]}"
    else
        export REMOTE_USER=
    fi
fi

env

ret_code=0
/app/bin/list-files.sh $1 | while read line; do 
    if [[ $line == ./* ]]; then
        # remove prefix './'
        line=${line:2}
    fi
    send-message $line
    code=$?
    if [ $code -ne 0 ]; then
        ret_code=$code
        echo "Error send-message, message:"$line 2>&
    fi
done

code=${PIPESTATUS[0]}
if [ $code -ne 0 ]; then
    ret_code=$code
    echo "Error run list-files.sh "$1 2>&
fi

exit $ret_code
