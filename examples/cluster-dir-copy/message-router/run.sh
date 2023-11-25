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
"cluster-dir-list") # list all files in dir
    if [[ $1 =~ ^{.+}$ ]]; then
        # dataset metadata in json-format
        echo message for dataset metadata 
        # send-message "data-grouping" 水平分组
        # send-message "data-grouping" 垂直分组
    else
        read -r fmt < "/work/.scalebox/my.txt"
        m=$(printf $fmt $1)
        send-message "cluster-file-copy" $m
    fi
    ;;
"data-grouping")
    # if message为水平分组消息 then
    #     发一个垂直分组的消息回 data-grouping
    # else
    #     所有文件拷贝完成，
    #     发送消息给data-grouping，通知其退出运行
    #     按需设置app完成的标志，置最终结果
    # fi
    ;;
"cluster-file-copy")
    # ignore
    ;;
*)
    # 发送给cluster-dir-list模块，启动
    source_cluster=$(echo $1 | cut -d "~" -f 1)
    dir=$(echo $1 | cut -d "~" -f 2)
    target_cluster=$(echo $1 | cut -d "~" -f 3)
    rpath=$(echo $dir | cut -d "#" -f 1)
    echo "$source_cluster~$rpath#%s~$target_cluster" > /work/.scalebox/my.txt

    if [ "$source_cluster" != "" ]; then
        send-message "cluster-dir-list" "${source_cluster}~${dir}"
    else
        send-message "cluster-dir-list" "~${dir}"
    fi
    ;;
esac

exit $?
