#!/bin/bash

s=$1
if [[ $s =~ ^(/.*)$ ]]; then
    # /my-root
    echo "LOCAL $1"
elif [[ $s =~ ^(rsync://([^@:]+(:[^@]+)@)?[^:/]+(:[0-9]+)?/.*)$ ]]; then
    # rsync://user:pass@myhost:873/my-root
    # rsync://user@myhost/my-root
    rsync_url=${BASH_REMATCH[1]}
    rsync_pass=${BASH_REMATCH[3]}
    rsync_port=${BASH_REMATCH[4]}
    # remove rsync_pass from rsync_url
    rsync_url=${rsync_url//$rsync_pass/}
    echo "RSYNC $rsync_url ${rsync_pass:1}"
elif [[ $s =~ ^([^@]+@[^:/]+)(:[0-9]+)?(/.*)$ ]]; then
    # user@myhost:22/my-root#my-dir
    # user@myhost/my-root#my-dir
    ssh_host=${BASH_REMATCH[1]}
    ssh_port=${BASH_REMATCH[2]}
    data_root=${BASH_REMATCH[3]}
    rsync_url="${ssh_host}:${data_root}"
    if [ "$ssh_port" == "" ]; then
        ssh_port="22"
    else
        ssh_port=${ssh_port:1}
    fi
    echo "SSH $rsync_url $ssh_port"
else
    echo "wrong message format, message:"$1 >&2
    exit 26
fi

exit 0
