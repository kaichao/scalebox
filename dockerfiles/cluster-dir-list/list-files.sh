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

if [ "$SOURCE_CLUSTER" == "" ]; then
    cluster=$CLUSTER_NAME
else
    cluster=$SOURCE_CLUSTER
fi
v=${cluster_map[$cluster]}
if [ "$v" == "" ]; then
    v=$(scalebox cluster get-parameter --cluster $cluster rsync)
    code=$?
    [[ $code -ne 0 ]] && echo cmd: get_cluster_rsync, error_code:$code && exit $code
    cluster_map[$cluster]=$v
    echo $cluster $v >> /work/.scalebox/cluster_data.txt
fi
rsync_prefix=$(echo $v | cut -d "#" -f 1)
ssh_port=$(echo $v | cut -d "#" -f 2)
dir1=$(echo $1 | cut -d "#" -f 1)
dir2=$(echo $1 | cut -d "#" -f 2)

if [ "$SOURCE_CLUSTER" == "" ]; then
    # local dir
    cluster_root=$(echo $rsync_prefix | cut -d ":" -f 2)
    # set /local to support symlink
    data_dir="/local${cluster_root}/${dir1}/${dir2}"
    echo local data-dir:${data_dir} >&2
    cd ${data_dir} && find -L . -type f \
        | sed "s/^\./${dir2}/" \
        | egrep "${REGEX_FILTER}"
        # | sed 's/^\.\///' \
else
    # rsync-over-ssh
    ssh_args="-T -c aes128-gcm@openssh.com -o Compression=no -x"
    data_dir=$(echo $1 | cut -d "#" -f 1)
    echo remote data-dir:${rsync_prefix}/${dir1}/${dir2} >&2
    rsync -avn -L -e "ssh -p ${ssh_port} ${ssh_args}" ${rsync_prefix}/${dir1}/${dir2} \
        | grep ^\- | awk {'print $5'}  \
        | awk '{ gsub(/^[^\/]+?\//,""); print $0 }' \
        | sed "s/^/${dir2}\//" \
        | egrep "${REGEX_FILTER}" 
fi

# exit status of egrep
#   0 if a line is selected
#   1 if no lines were selected
#   2 if an error occurred.
status=(${PIPESTATUS[@]})
echo "[INFO]pipe_status:"${status[*]} >&2
n=${#status[*]}
if [ $n == 1 ]; then
    if [ ${status[0]} -ne 0 ]; then
        echo "[ERROR]local mode, dir: "${LOCAL_ROOT}" not found" >&2
        exit ${status[0]}
    fi
fi

declare -i code
for ((i=n-1; i>=0; i--)); do
    code=${status[i]}
    if [ $i == $((n-1)) ]; then
        if [ $code == 1 ]; then
            echo "[WARN]All of data are filtered, empty dataset!" >&2
            code=0
        fi
    fi
    if [ $code -ne 0 ]; then
        break
    fi
done

exit $code
