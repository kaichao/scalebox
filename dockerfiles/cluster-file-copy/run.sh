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

if [[ $1 == *~* ]]; then
    m=$1
else
    m="${SOURCE_CLUSTER}~$1~${TARGET_CLUSTER}"
fi

#  extract SOURCE_CLUSTER/TARGET_CLUSTER from message-body
if [[ "$m" =~ ^([^~]+)~([^~]+)~$ ]]; then
    #   pull:   ${SOURCE_CLUSTER}~${FILE_PATH}~
    SOURCE_CLUSTER=${BASH_REMATCH[1]}
    m=${BASH_REMATCH[2]}
elif [[ "$m" =~ ^~([^~]+)~([^~]+)$ ]]; then
    #   push:   ~${FILE_PATH}~${TARGET_CLUSTER}
    m=${BASH_REMATCH[1]}
    TARGET_CLUSTER=${BASH_REMATCH[2]}
else
    echo "message $m not valid, only one of SOURCE_CLUSTER and TARGET_CLUSTER is allowed in the message body"
    exit 11
fi

file=$(echo $m | tr "#" "/")

echo source_cluster:$SOURCE_CLUSTER
echo target_cluster:$TARGET_CLUSTER
echo filename:$file

ssh_args="-T -c aes128-gcm@openssh.com -o Compression=no -x"
if [[ $ZSTD_CLEVEL != "" ]]; then 
    rsync_args="--cc=xxh3 --compress --compress-choice=zstd --compress-level=${ZSTD_CLEVEL}"
else
    rsync_args=""
fi

echo "ssh_args:"$ssh_args

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
ssh_args="ssh -p ${ssh_port} ${ssh_args}"

if [ "$SOURCE_CLUSTER" != "" ]; then
    dest_dir=$(dirname /data/$file)
    mkdir -p ${dest_dir}
    cmd="rsync -ut -L ${rsync_args} -e \"${ssh_args}\" ${rsync_prefix}/${file} ${dest_dir}"
else
    cd /data
    cmd="rsync -Rut -L ${rsync_args} -e \"${ssh_args}\" $file $rsync_prefix"
fi

eval $cmd
code=$?

if [[ $code -eq 0 ]]; then
    echo "rsync-over-ssh runs successfully."
    echo cmd:$cmd
    send-message $(echo $m | cut -d "#" -f 2)
    code=$?
fi

exit $code
