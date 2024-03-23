#!/bin/bash


if [[ $1 == *~* ]]; then
    m=$1
    if [ "$SOURCE_URL" != "" ] || [ "$TARGET_URL" != "" ]; then
        echo "WARNING: SOURCE_URL/TARGET_URL not null, but skipped!" >&2
    fi
else
    m="${SOURCE_URL}~$1~${TARGET_URL}"
fi

#  extract SOURCE_CLUSTER/TARGET_CLUSTER from message-body
if [[ "$m" =~ ^([^~]+)~([^~]+)~([^~]+)$ ]]; then
    #   pull:   ${SOURCE_CLUSTER}~${FILE_PATH}~
    source_url=${BASH_REMATCH[1]}
    m=${BASH_REMATCH[2]}
    target_url=${BASH_REMATCH[3]}
else
    echo "not valid message:$m " >&2
    exit 11
fi

arr_source=($(/app/share/bin/parse.sh $source_url))
code=$?
[[ $code -ne 0 ]] && echo error while parse_source_url, error_code:$code, source_url:$source_url >&2 && exit $code
source_mode=${arr_source[0]}

arr_target=($(/app/share/bin/parse.sh $target_url))
code=$?
[[ $code -ne 0 ]] && echo error while parse_target_url, error_code:$code, source_url:$target_url >&2 && exit $code
target_mode=${arr_target[0]}

if [ "$source_mode" == "LOCAL" ] || [ "$target_mode" == "FTP" ]; then
    action="PUSH"
    local_root=${arr_source[1]}
    ftp_url=${arr_target[1]}
    remote_root=${arr_target[2]}
elif [ "$source_mode" == "FTP" ] || [ "$target_mode" == "LOCAL" ]; then
    action="PULL"
    local_root=${arr_target[1]}
    ftp_url=${arr_source[1]}
    remote_root=${arr_source[2]}
else
    echo "Only FTP on one end and LOCAL on the other end are allowed, message:"$1 >&2
    exit 21
fi

filename=$(basename $m)

if [[ $remote_root == '/' ]]; then
    remote_dir=$(dirname /$m)
else
    remote_dir=$(dirname ${remote_root}/$m)
fi

if [[ $local_root == '/' ]]; then
    local_dir=$(dirname "/local"/$m)
else
    local_dir=$(dirname "/local"${local_root}/$m)
fi

work_dir="$local_dir"
[[ "$ENABLE_LOCAL_RELAY" == "yes" ]] && work_dir="/work"
[[ "$ENABLE_RECHECK_PUSH" == "yes" ]] && [ "$action" == "PUSH" ] && work_dir="/work"

printf "[DEBUG]ftp_ur=%s remote_dir=%s, local_dir=%s, work_dir=%s, filename=%s \n" \
    $ftp_url $remote_dir $local_dir $work_dir $filename

# 多空格分隔
cmd_remote_filesize="lftp -c \"open ${ftp_url}; ls ${remote_dir}|grep ${filename}\"|awk -F'  *' '{print \$5}'"
cmd_local_filesize="stat --printf=\"%s\" ${work_dir}/${filename}"
cmd_mkdir_remote="lftp -c \"open ${ftp_url}; mkdir -p ${remote_dir}\" 2> /work/mkdir-stderr.txt"
cmd_rm_local_file="rm -f ${local_dir}/${filename} /work/${filename}"
cmd_put_file="lftp -c \"open ${ftp_url}; cd ${remote_dir}; lcd ${work_dir}; put ${filename}\""
cmd_get_file="lftp -e \"lcd ${work_dir};pget -n $NUM_PGET_CONN ${ftp_url}${remote_dir}/${filename};exit\""

declare -i input_bytes=0
date --iso-8601=ns >> /work/timestamps.txt
if [[ $action == 'PULL' ]]; then
    cmd="mkdir -p ${local_dir}"
    eval $cmd; code=$?
    [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code >&2 && exit 21
    eval $cmd_rm_local_file
    echo 011, ok, mkdir_local_dir and rm_local_file

    eval $cmd_get_file; code=$?
    date --iso-8601=ns >> /work/timestamps.txt
    [[ $code -ne 0 ]] && echo cmd:$cmd_get_file, ret_code:$code >&2 && exit 22
    echo 012, ok, get_file from ftp-server

    local_size=$(eval $cmd_local_filesize);code=$?
    [[ $code -ne 0 ]] && echo cmd:$cmd_local_filesize, ret_code:$code >&2 && exit 23

    input_bytes=$local_size
    echo 013, ok, get_local_file_size

    if [[ $ENABLE_LOCAL_RELAY == 'yes' ]]; then
        mv /work/$filename $local_dir; code=$?
        [[ $code -ne 0 ]] && echo mv file from work_dir to local_dir, ret_code:$code >&2 && exit 24

        date --iso-8601=ns >> /work/timestamps.txt
        let input_bytes*=2
        echo 014, ok, mv local_file from work_dir
    fi

    if [[ $ENABLE_RECHECK_SIZE == 'yes' ]]; then
        remote_size=$(eval $cmd_remote_filesize); code=$?
        date --iso-8601=ns >> /work/timestamps.txt
        [[ $local_size -ne $remote_size ]] && \
            echo inconsistent file size, local_size:$local_size, remote_size:$remote_size, ret_code:$code >&2 && exit 25
        echo 015, ok, recheck local_file_size and remote_file_size
    fi
else
    date --iso-8601=ns >> /work/timestamps.txt
    eval ${cmd_mkdir_remote}; code=$?
    date --iso-8601=ns >> /work/timestamps.txt
    if [[ $code -ne 0 ]]; then
        #   550  目录已存在
        if ! grep -q "550" /work/mkdir-stderr.txt; then
            echo err while creating remote-dir >&2
            exit 31
        fi
    fi 
    echo 021, ok, mkdir_remote_dir

    if [[ $ENABLE_LOCAL_RELAY == 'yes' ]] || [[ $ENABLE_RECHECK_PUSH == 'yes' ]]; then
        cmd="cp $local_dir/$filename /work"
        eval $cmd; code=$?
        [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:32 >&2 && exit 32
        echo 022, ok, cp file from local_dir to work_dir
        let input_bytes*=2
    fi

    # push to ftp-server
    eval $cmd_put_file; code=$?
    [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:33 >&2 && exit 33
    echo 023, ok, put file to ftp-server

    local_size=$(eval $cmd_local_filesize)
    [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:34 >&2 && exit 34
    date --iso-8601=ns >> /work/timestamps.txt
    echo 024, ok, get local_file_size in work_dir

    if [[ $ENABLE_RECHECK_SIZE == 'yes' ]]; then
        remote_size=$(eval $cmd_remote_filesize)
        [[ $local_size -ne $remote_size ]] && echo inconsistent file size, local_size:$local_size, remote_size:$remote_size, exit_code:35 >&2 && exit 35
        echo 025, ok, get remote_file_size in ftp-server
    fi

    if [[ $ENABLE_RECHECK_PUSH == 'yes' ]]; then
        sha1sum /work/$filename > /work/sha1sum.txt
        rm -f /work/$filename

        eval $cmd_get_file; code=$?
        [[ $code -ne 0 ]] && echo cmd:$cmd, ret_code:$code, exit_code:36 >&2 && exit 36
        date --iso-8601=ns >> /work/timestamps.txt
        echo 026, ok, get remote_file from ftp-server for re-check

        sha1sum -c /work/sha1sum.txt; code=$?
        [[ $code -ne 0 ]] && echo sha1sum recheck error, ret_code:$code, exit_code:37 >&2 && exit 37
        echo 027, ok, re-check local_file in work-dir

        let input_bytes+=local_size
    fi
fi

date --iso-8601=ns >> /work/timestamps.txt

cat << EOF > /work/task-exec.json
{
	"inputBytes":${input_bytes},
	"outputBytes":${input_bytes}
}
EOF

exit $code
