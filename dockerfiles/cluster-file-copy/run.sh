#!/bin/bash

# create hashtable cluster_map
declare -A cluster_map
if [ -e /work/.scalebox/cluster_data.txt ]; then
    while read line; do
        # separate key and value
        key=$(echo $line | cut -d " " -f 1)
        value=$(echo $line | cut -d " " -f 2)
        # put key value to hashtable
        cluster_map[$key]=$value
    done < /work/.scalebox/cluster_data.txt
fi

file=$(echo $1 | tr "#" "/")
if [ "$SOURCE_CLUSTER" == "" ] && [ "$TARGET_CLUSTER" == "" ]; then
    # Both variables are empty
    #  extract SOURCE_CLUSTER/TARGET_CLUSTER from message-body
    #   pull:   ${SOURCE_CLUSTER}~${FILE_PATH}
    #   push:   ~${FILE_PATH}~${TARGET_CLUSTER}
    if [[ "$file" =~ ^([^~]+)~([^~]+)$ ]]; then
        SOURCE_CLUSTER=${BASH_REMATCH[1]}
        file=${BASH_REMATCH[2]}
    elif [[ "$file" =~ ^~([^~]+)~([^~]+)$ ]]; then
        file=${BASH_REMATCH[1]}
        TARGET_CLUSTER=${BASH_REMATCH[2]}
    else
        echo no valid SOURCE_CLUSTER or TARGET_CLUSTER
        exit 11
    fi
fi
echo filename:$file

ssh_args="-T -c aes128-gcm@openssh.com -o Compression=no -x"
if [[ $ZSTD_CLEVEL != "" ]]; then 
    rsync_args="--cc=xxh3 --compress --compress-choice=zstd --compress-level=${ZSTD_CLEVEL}"
else
    rsync_args=""
fi

echo "ssh_args:"$ssh_args


# scalebox cluster get-parameter --cluster $SOURCE_CLUSTER rsync_prefix
# return value:  <ssh-user>@<host>:/<data-dir>#<port>#<jump-servers>
if [ "$SOURCE_CLUSTER" != "" ] && [ "$TARGET_CLUSTER" != "" ]; then
    # dual-remote copy
    exit 0
fi

if [ "$SOURCE_CLUSTER" != "" ]; then
    cluster=$SOURCE_CLUSTER
else
    cluster=$TARGET_CLUSTER
fi

v=cluster_map[$cluster]
if [ "$v" != "" ]; then
    v=$(scalebox cluster get-parameter --cluster $cluster rsync)
    code=$?
    [[ $code -ne 0 ]] && echo cmd: get_cluster_rsync, error_code:$code && exit $code
    cluster_map[$cluster]=$v
    echo $cluster $v >> /work/.scalebox/cluster_data.txt
fi
rsync_prefix=$(echo $v | cut -d "#" -f 1)
ssh_port=$(echo $v | cut -d "#" -f 2)
# jump_servers=$(echo $v | cut -d "#" -f 3)

if [ "$SOURCE_CLUSTER" != "" ]; then
    dest_dir=$(dirname /data/$file)
    mkdir -p ${dest_dir}
    cmd="rsync -ut -L ${rsync_args} -e \"ssh -p ${ssh_port} ${ssh_args}\" ${rsync_prefix}/${file} ${dest_dir}"
else
    cd /data
    cmd="rsync -Rut -L ${rsync_args} -e \"ssh -p ${ssh_port} ${ssh_args}\" $file $rsync_prefix"
fi

echo cmd:$cmd
eval $cmd
code=$?

if [[ $code -eq 0 ]]; then
    send-message $(echo $1 | cut -d "#" -f 2)
    code=$?
fi

exit $code
