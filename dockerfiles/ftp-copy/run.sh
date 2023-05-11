#!/bin/bash

filename=$(basename $1)

. /env.sh

if [[ $REMOTE_URL =~ (ftp://([^:]+:[^@]+@)?[^/:]+(:[^/]+)?)(/.*) ]]; then
    ftp_url=${BASH_REMATCH[1]}
    remote_root=${BASH_REMATCH[4]}
else
    echo "REMOTE_URL did not match regex! exit_code:6" >&2
    exit 6
fi

if [[ $remote_root == '/' ]]; then
    remote_dir=$(dirname /$1)
else
    remote_dir=$(dirname ${remote_root}/$1)
fi

if [[ $LOCAL_ROOT == '/' ]]; then
    local_dir=$(dirname "/local"/$1)
else
    local_dir=$(dirname "/local"${LOCAL_ROOT}/$1)
fi

printf "[DEBUG]ACTION=%s, REMOTE_URL=%s, LOCAL_ROOT=%s, file=%s\n"  $ACTION $REMOTE_URL $LOCAL_ROOT $1
printf "[DEBUG]ftp_ur=%s remote_dir=%s, local_dir=%s, filename=%s \n" $ftp_url $remote_dir $local_dir $filename

if [ "$ENABLE_LOCAL_RELAY" == "yes" ] || [ "$ACTION" == "PUSH_RECHECK" ]; then 
    work_dir="/work"
else
    work_dir="$local_dir"
fi

# 多空格分隔
cmd_remote_filesize="lftp -c \"open ${ftp_url}; ls ${remote_dir}|grep ${filename}\"|awk -F'  *' '{print \$5}'"
cmd_local_filesize="stat --printf=\"%s\" ${local_dir}/${filename}"
cmd_mkdir_remote="lftp -c \"open ${ftp_url}; mkdir -p ${remote_dir}\" 2> /dev/null"
cmd_rm_local_file="rm -f ${local_dir}/${filename} /work/${filename}"
cmd_put_file="lftp -c \"open ${ftp_url}; cd ${remote_dir}; lcd ${work_dir}; put ${filename}\""
cmd_get_file="lftp -e \"lcd ${work_dir};pget -n $NUM_PGET_CONN ${ftp_url}${remote_dir}/${filename};exit\""


declare -i input_bytes=0

date --iso-8601=ns >> /work/timestamps.txt

if [[ $ACTION == 'PUSH' ]]; then
    eval ${cmd_mkdir_remote}

    if [[ $ENABLE_LOCAL_RELAY == 'yes' ]]; then
        cmd="cp $local_dir/$filename /work"
        eval $cmd; code=$?
        if [[ $code -ne 0 ]]; then
            echo cmd:$cmd, ret_code:$code, exit_code:11 >&2
            rm -f /work/${filename}
            exit 11
        fi
        # [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:11 >&2 && exit 11
    fi
    date --iso-8601=ns >> /work/timestamps.txt

    # echo remote_dir:$remote_dir, local_dir:$local_dir, filename:$filename
    # push to ftp-server
    eval $cmd_put_file; code=$?

    local_size=$(eval $cmd_local_filesize)
    input_bytes=$local_size
    if [[ $ENABLE_LOCAL_RELAY == 'yes' ]]; then
        rm -f /work/$filename
        let input_bytes*=2
    fi
    [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:12 >&2 && exit 12
    date --iso-8601=ns >> /work/timestamps.txt

    if [[ $ENABLE_RECHECK_SIZE == 'yes' ]]; then
        remote_size=$(eval $cmd_remote_filesize)
        [[ $local_size -ne $remote_size ]] && echo inconsistent file size, local_size:$local_size, remote_size:$remote_size, exit_code:14 >&2 && exit 14
    fi
elif [[ $ACTION == 'PUSH_RECHECK' ]]; then
    eval ${cmd_mkdir_remote}

    # ENABLE_LOCAL_RELAY == 'yes'
    cmd="cp $local_dir/$filename /work"
    eval $cmd; code=$?
    if [[ $code -ne 0 ]]; then
        echo cmd:$cmd, ret_code:$code, exit_code:11 >&2
        rm -f /work/${filename}
        exit 11
    fi
    date --iso-8601=ns >> /work/timestamps.txt

    # push to ftp-server
    cd /work
    # cmd="lftp -c \"open ${ftp_url}; cd ${remote_dir}; lcd /work; put ${filename}\""
    eval $cmd_put_file; code=$?
    if [[ $code -ne 0 ]]; then
        rm -f /work/${filename}
        echo cmd:$cmd, ret_code:$code, exit_code:12 >&2
        exit 12
    fi
    date --iso-8601=ns >> /work/timestamps.txt

    # re-pull from ftp-server
    mv /work/${filename} /work/${filename}.orig
    eval $cmd_get_file; code=$?
    if [[ $code -ne 0 ]]; then
        echo cmd:$cmd, ret_code:$code, exit_code:13 >&2
        rm -f /work/${filename} /work/${filename}.orig
        exit 13
    fi
    date --iso-8601=ns >> /work/timestamps.txt

    local_size=$(eval $cmd_local_filesize)
    # compare original file with re-pulled file.
    cmd="diff /work/${filename} /work/${filename}.orig"
    eval $cmd; code=$?
    if [[ $code -ne 0 ]]; then
        remote_size=$(eval "stat --printf=\"%s\" /work/${filename}")
        rm -f /work/${filename} /work/${filename}.orig
        if [[ $local_size -ne $remote_size ]]; then
            echo diff remote_file local_filed, exit_code:14 >&2
            exit 14
        else
            echo diff remote_file local_filed, exit_code:15 >&2
            exit 15
        fi
    fi
    rm -f /work/${filename} /work/${filename}.orig
    echo "The original file or the returned file is the same."

    input_bytes=$local_size
    let input_bytes*=3
elif [[ $ACTION == 'PULL' ]]; then
    cmd="mkdir -p ${local_dir}"
    eval $cmd; code=$?
    [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:21 >&2 && exit 21

    eval $cmd_rm_local_file
    eval $cmd_get_file; code=$?

    date --iso-8601=ns >> /work/timestamps.txt
    [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:22 >&2 && exit 22

    if [[ $ENABLE_LOCAL_RELAY == 'yes' ]]; then
        mv /work/$filename $local_dir
    fi
    local_size=$(eval $cmd_local_filesize)
    input_bytes=$local_size
    if [[ $ENABLE_LOCAL_RELAY == 'yes' ]]; then
        date --iso-8601=ns >> /work/timestamps.txt
        let input_bytes*=2
    fi

    if [[ $ENABLE_RECHECK_SIZE == 'yes' ]]; then
        remote_size=$(eval $cmd_remote_filesize)
        date --iso-8601=ns >> /work/timestamps.txt

        [[ $local_size -ne $remote_size ]] && echo inconsistent file size, local_size:$local_size, remote_size:$remote_size, exit_code:23 >&2 && exit 23
    fi
fi
date --iso-8601=ns >> /work/timestamps.txt
rm -f /work/${filename} /work/${filename}.orig

cat << EOF > /work/task-exec.json
{
	"inputBytes":${input_bytes},
	"outputBytes":${input_bytes}
}
EOF

# lftp ${ftp_url} << EOF
#     cd ${remote_dir}
#     lcd ${work_dir}
#     put ${filename}
#     by
# EOF

exit $code
