#!/bin/bash

# 确认命令行参数3个，支持消息的headers
if [ $# -ne 2 ]; then
    echo "Usage: $0 message-body headers" >&2
    exit 11
fi

# extract from_job from headers
headers=$2
if [[ $headers =~ \"from_job\":\"(.+)\" ]]; then
    from_job=${BASH_REMATCH[1]}
fi

case $from_job in
"dir-list") # list all files in dir
    read -r target_url < "/work/.scalebox/my.txt"
    m=$1~$target_url
    if [[ $1 =~ ^{.+}$ ]]; then
        # ignore dataset metadata in json-format
        echo message for dataset metadata 
    elif [[ $m =~ ftp:// ]]; then
        send-message "ftp-copy" $m
    else
        send-message "rsync-copy" $m
    fi
    ;;
"rsync-copy")
    # ignore
    ;;
"ftp-copy")
    # ignore
    ;;
*)
    #    发送给模块，启动
    source_url=$(echo $1 | cut -d "~" -f 1)
    target_url=$(echo $1 | cut -d "~" -f 2)
    echo $target_url > /work/.scalebox/my.txt
    send-message "dir-list" $source_url~.
    ;;
esac

exit $?
