#!/bin/bash

if [ "$SOURCE_URL" == "" ]; then
    dir=$1
else
    dir=$SOURCE_URL~$1
fi

if [ "$REGEX_2D_DATASET" ]; then
    meta=$(list-files.sh $dir | get_2d_meta $REGEX_2D_DATASET $INDEX_2D_DATASET)
    code=$?
    [[ $code -ne 0 ]] && echo cmd: get_2d_meta, error_code:$code && exit $code
    echo ${meta} > /work/key-text.txt
    echo metadata for 2d-dataset:#${meta}#
    # key text in file /work/key-text.txt
    scalebox task add
    code=$?
    [[ $code -ne 0 ]] && echo cmd: scalebox task add, error_code:$code && exit $code
fi

prefix=$(echo $dir | cut -d "~" -f 1)
ret_code=0
list-files.sh $dir | while read line; do 
    send-message "${prefix}~${line}"
    code=$?
    [[ $code -ne 0 ]] && echo "Error send-message, file:"$line >&2 && exit $code
done

code=${PIPESTATUS[0]}
[[ $code -ne 0 ]] && echo "Error run list-files.sh, dir:"$dir >&2 >&2 && exit $code

exit 0
