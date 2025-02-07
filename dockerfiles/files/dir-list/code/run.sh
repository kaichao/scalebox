#!/bin/bash

# set -e
source /usr/local/bin/functions.sh
source /app/share/bin/functions.sh

export REGEX_FILTER=$(get_header "$2" "regex_filter")
export JUMP_SERVERS=$(get_header "$2" "jump_servers")

prefix_url=$(get_header "$2" "prefix_url")

dir_name=$1

echo prefix_url:$prefix_url

/app/share/bin/list-files.sh "$prefix_url~$dir_name" | while read line; do 
    scalebox task add --header source_url="$prefix_url" "${line}"
    code=$?
    [[ $code -ne 0 ]] && echo "Error send-message, file:"$line >&2 && exit $code
done

code=${PIPESTATUS[0]}
[[ $code -ne 0 ]] && echo "Error run list-files.sh, dir:"$prefix_url >&2 && exit $code

exit 0
