#!/bin/bash

dir=$1

if [ "$SOURCE_CLUSTER" == "" ]; then
    #  extract SOURCE_CLUSTER from message-body
    #       ${SOURCE_CLUSTER}~${RELATIVE_PATH}#${local_dir}
    if [[ "$dir" =~ ^([^~]*)~([^#]+)#([^#]+)$ ]]; then
        SOURCE_CLUSTER=${BASH_REMATCH[1]}
        dir=${BASH_REMATCH[2]}#${BASH_REMATCH[3]}
    fi
fi
echo dir-name:$dir

if [ "$REGEX_2D_DATASET" ]; then
    meta=$(list-files.sh $dir | get_2d_meta $REGEX_2D_DATASET $INDEX_2D_DATASET)
    code=$?
    if [[ $code -eq 0 ]]; then
        echo ${meta} > /work/key-text.txt
        echo metadata for 2d-dataset:#${meta}#
        # key text in file /work/key-text.txt
        scalebox task add; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] scalebox task add metadata:${meta}, error_code:$code" >&2 && exit $code
    elif [[ $code -eq 3 ]]; then
        echo "[WARN] get_2d_metadata, error_code:$code" >&2
    fi
fi

list-files.sh $dir | while read line; do 
    send-message ${line}; code=$?
    [[ $code -ne 0 ]] && echo "Error send-message, file:"$line >&2 && exit $code
done

code=${PIPESTATUS[0]}
[[ $code -ne 0 ]] && echo "Error run list-files.sh, dir:"$dir >&2 && exit $code

exit 0
