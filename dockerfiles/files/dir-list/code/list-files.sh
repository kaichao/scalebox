#!/bin/bash

if [[ "$ENTRY_TYPE" == "DIR" ]]; then
    entry_type="d"
    rsync_args="--include '*/' --exclude '*'"
else
    entry_type="f"
fi

s=$1

echo "message:"$s >&2
if [[ $s =~ ^(/[^~]*)~(.+)$ ]]; then
    # /my-root~my-dir
    echo "local-dir" >&2
    data_root=${BASH_REMATCH[1]}
    dir=${BASH_REMATCH[2]}
    data_dir="/local${data_root}/$dir"
    echo dir:$dir >&2
    echo data-dir:$data_dir >&2
    echo "ENTRY_TYPE:$ENTRY_TYPE, entry_type:$entry_type" >&2
    
    # Use percent sign % as separator
    cd ${data_dir} && find -L . -type $entry_type \
        | sed "s%^\.%${dir}%" \
        | sed 's/^\.\///' \
        | egrep "${REGEX_FILTER}"
elif [[ $s =~ ^(rsync://([^@:]+(:[^@]+)@)?[^:/]+(:[0-9]+)?/[^~]*)~(.+)$ ]]; then
    # rsync://user:pass@myhost:873/my-root~my-dir
    # rsync://user@myhost/my-root~my-dir
    rsync_url=${BASH_REMATCH[1]}
    rsync_pass=${BASH_REMATCH[3]}
    rsync_port=${BASH_REMATCH[4]}
    dir=${BASH_REMATCH[5]}
    rsync_url=${rsync_url//$rsync_pass/}
    # rsync_url=${rsync_url//$port/}
    export RSYNC_PASSWORD=${rsync_pass:1}

    # echo "rsync_url:${rsync_url}" >&2
    # echo dir:$dir >&2

    
    rsync -avn "${rsync_args}" "${rsync_url}/${dir}" \
        | grep ^\- | awk {'print $5'} \
        | awk '{ gsub(/^[^\/]+?\//,""); print $0 }' \
        | sed "s%^%${dir}\/%" \
        | sed 's/^\.\///' \
        | egrep "${REGEX_FILTER}"
elif [[ $s =~ ^([^@]+@[^:/]+)(:[0-9]+)?(/[^~]*)~(.+)$ ]]; then
    # user@myhost:22/my-root~my-dir
    # user@myhost/my-root~my-dir
    echo "rsync-over-ssh" >&2
    ssh_host=${BASH_REMATCH[1]}
    ssh_port=${BASH_REMATCH[2]}
    data_root=${BASH_REMATCH[3]}
    dir=${BASH_REMATCH[4]}
    rsync_url="${ssh_host}:${data_root}"
    if [ "$ssh_port" == "" ]; then
        ssh_port="22"
    else
        ssh_port=${ssh_port:1}
    fi

    echo "rsync_url:${rsync_url}" >&2
    echo port:$ssh_port >&2
    echo dir:$dir  >&2

    ssh_args="-T -c aes128-gcm@openssh.com -o Compression=no -x"
    echo rsync -avn -L ${rsync_args} -e \"ssh -p ${ssh_port}\" ${rsync_url}/${dir} >&2

    if [[ "$ENTRY_TYPE" == "FILE" ]]; then
        rsync -avn -L -e "ssh -p ${ssh_port} ${ssh_args}" ${rsync_url}/${dir} \
            | grep ^\- | awk {'print $5'}  \
            | awk '{ gsub(/^[^\/]+?\//,""); print $0 }' \
            | sed "s%^%${dir}\/%" \
            | sed 's/^\.\///' \
            | egrep "${REGEX_FILTER}"
    elif [[ "$ENTRY_TYPE" == "DIR" ]]; then
        rsync -avn -L ${rsync_args} -e "ssh -p ${ssh_port}" ${rsync_url}/${dir} \
            | grep ^dr[w\-]x | awk {'print $5'}  \
            | awk '{ gsub(/^[^\/]+?\//,""); print $0 }' \
            | sed "s%^%${dir}\/%" \
            | sed 's/^\.\///' \
            | egrep "${REGEX_FILTER}"
    fi
else
    echo "wrong message format, message:"$1 >&2
    exit 21
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
        echo "[ERROR]local mode, dir: ${LOCAL_ROOT} not found" >&2
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
