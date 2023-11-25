#!/bin/bash

if [ "$SOURCE_URL" == "" ]; then
    dir=$1
else
    dir=$SOURCE_URL~$1
fi

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

prefix=$(echo $dir | cut -d "~" -f 1)
list-files.sh $dir | while read line; do 
    send-message "${prefix}~${line}"
    code=$?
    [[ $code -ne 0 ]] && echo "Error send-message, file:"$line >&2 && exit $code
done

code=${PIPESTATUS[0]}
[[ $code -ne 0 ]] && echo "Error run list-files.sh, dir:"$dir >&2 >&2 && exit $code

exit 0
