#!/bin/bash

set -e

source /usr/local/bin/functions.sh
source /app/share/bin/functions.sh

# Support singularity
[[ ! $WORK_DIR ]] && { echo "[ERROR] WORK_DIR is null, Check the permissions of the directory /tmp/scalebox." >&2; exit 110; }
echo "[DEBUG] WORK_DIR:${WORK_DIR}:" >> ${WORK_DIR}/custom-out.txt
cd "${WORK_DIR}"

source_url=$(get_header "$2" "source_url")
target_url=$(get_header "$2" "target_url")

if [ -z "$SOURCE_MODE" ]; then
    source_mode=$(get_mode "$source_url")
else
    source_mode="$SOURCE_MODE"
fi

if [ -z "$TARGET_MODE" ]; then
    target_mode=$(get_mode "$target_url")
else
    target_mode="$TARGET_MODE"
fi

source_dir=$(get_data_root "$source_url")
target_dir=$(get_data_root "$target_url")

echo "[DEBUG]source_url:$source_url,target_url:$target_url,message:$1" >> ${WORK_DIR}/custom-out.txt
echo "[DEBUG]source_dir:$source_dir,target_dir:$target_dir" >> ${WORK_DIR}/custom-out.txt

if [[ $ZSTD_CLEVEL != "" ]]; then
    rsync_args="--cc=xxh3 --compress --compress-choice=zstd --compress-level=${ZSTD_CLEVEL}"
else
    rsync_args=""
fi

date --iso-8601=ns | sed 's/,/./' >> timestamps.txt

echo "[DEBUG]source_mode:$source_mode,target_mode:$target_mode" >> ${WORK_DIR}/custom-out.txt
case $source_mode in
"LOCAL")
    case $target_mode in
    "LOCAL")    exit 31 ;;
    "SSH")
        echo "LOCAL-SSH" >> custom-out.txt
        ssh_cmd=$(get_ssh_cmd "$2" "target_url" "target_jump_servers")
        echo "[DEBUG] ssh_cmd:$ssh_cmd" >> custom-out.txt
        local_file=$(get_host_path "$source_dir/$1")

        remote_file="$target_dir/$1"

        # source file not exists ?
        [ ! -f "$local_file" ] && echo "file $local_file not exists, exit " && exit 101

        # create directory in remote side.
        remote_dir=$(dirname $remote_file)
        eval "$ssh_cmd mkdir -p $remote_dir"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in remote dir,dir_name:$remote_dir, error_code:$code" >&2 && exit $code

        cmd="cat ${local_file} | pv | $ssh_cmd \"cat > $remote_file\""
        echo "[DEBUG] cmd:$cmd" >> ${WORK_DIR}/custom-out.txt
        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from local to remote, cmd=$cmd, error_code:$code" >&2 && exit $code

        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            echo "$local_file" >> ${WORK_DIR}/removed-files.txt
        fi
        echo "$local_file" >> ${WORK_DIR}/input-files.txt
        echo "$local_file" >> ${WORK_DIR}/output-files.txt

        echo "[DEBUG] local_file:$local_file" >> ${WORK_DIR}/custom-out.txt
        ;;
    "RSYNC_OVER_SSH") 
        echo "LOCAL to RSYNC_OVER_SSH" >> ${WORK_DIR}/custom-out.txt
        target_ssh_option=$(get_ssh_option "$2" "target_url" "target_jump_servers")

        if [[ $source_url == /data/* ]]; then
            source_dir="${source_url}/"
        else
            source_dir=$(get_host_path "${source_url}/")
        fi
        echo "source_dir:$source_dir" >> ${WORK_DIR}/custom-out.txt
        cd $source_dir

        target_ssh_url=$(to_ssh_url $target_url)
        # target_ssh_url=$(dirname "$target_ssh_url/$1")
        cmd="rsync -Rut ${rsync_args} -e \"ssh ${target_ssh_option}\" $1 $target_ssh_url/ "
        echo "cmd=$cmd" >> ${WORK_DIR}/custom-out.txt
        eval $cmd
        # rsync -ut ${rsync_args} -e "ssh ${ssh_option}" $source_url/$1 ${dest_dir}
        code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to remote, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/custom-out.txt && exit $code
        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            echo "$source_dir/$1" > ${WORK_DIR}/removed-files.txt
            # cmd="rm -f $1"
            # echo "cmd_remove_source_file: $cmd" >> ${WORK_DIR}/custom-out.txt
            # eval "$cmd"
            # [[ $? -ne 0 ]] && echo "[WARN] error while remove remote source file :$source_file" >> ${WORK_DIR}/custom-out.txt
        fi
    ;;
    "RSYNC")
        export RSYNC_PASSWORD=cas12345

        rsync_url=$(to_rsync_url "$target_url")

        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            rsync_args=" --remove-source-files ${rsync_args} "
        fi
        local_dir=$(get_host_path "${source_dir}")
        cmd="cd $local_dir && rsync -ut ${rsync_args} $1 $rsync_url"
        echo "cmd:$cmd"

        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from local to remote, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/custom-out.txt && exit $code
    ;;
    *)      exit 45 ;;
    esac
    ;;
"SSH")
    case $target_mode in
    "LOCAL")
        echo "SSH to LOCAL" >> custom-out.txt
        ssh_cmd=$(get_ssh_cmd "$2" "source_url" "source_jump_servers")
        echo "[DEBUG] ssh_cmd:$ssh_cmd" >> ${WORK_DIR}/custom-out.txt
        local_file=$(get_host_path "$target_dir/$1")
        remote_file="$source_dir/$1"

        # create directory in local side.
        local_dir=$(dirname $local_file)
        mkdir -p $local_dir; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in local dir,dir_name:$local_dir, error_code:$code" >&2 && exit $code

        cmd="$ssh_cmd \"cat < $remote_file\" - | pv > ${local_file}"
        echo "[DEBUG] cmd:$cmd" >> ${WORK_DIR}/custom-out.txt
        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to local, cmd=$cmd, error_code:$code" >&2 && exit $code

        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            eval "$ssh_cmd rm -f $remote_file"
        fi
        echo $local_file >> ${WORK_DIR}/input-files.txt
        echo $local_file >> ${WORK_DIR}/output-files.txt
        ;;
    "SSH")
        echo "SSH to SSH" >> ${WORK_DIR}/custom-out.txt
        source_ssh_cmd=$(get_ssh_cmd "$2" "source_url" "source_jump_servers")
        target_ssh_cmd=$(get_ssh_cmd "$2" "target_url" "target_jump_servers")

        source_file="$source_dir/$1"
        target_file="$target_dir/$1"

        # create directory in remote side.
        target_dir=$(dirname $target_file)
        eval "$target_ssh_cmd mkdir -p $target_dir"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in remote dir,dir_name:$remote_dir, error_code:$code" >&2 && exit $code

        cmd="$source_ssh_cmd \"cat < $source_file\" - | pv -n | $target_ssh_cmd \"cat > $target_file\""
        bytes_transferred=$(eval "$cmd"); code=$?
        echo "bytes_transferred: $bytes_transferred" >> custom-out.txt

        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to remote, cmd=$cmd, error_code:$code" >&2 && exit $code

        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            eval "$source_ssh_cmd rm -f $source_file"
        fi
        # echo $local_file >> ${WORK_DIR}/input-files.txt
        # echo $local_file >> ${WORK_DIR}/output-files.txt

    # ssh user1@node1 "cat /path/to/file" | ssh user2@node2 "cat > /path/to/destination"
        ;;
    "RSYNC")    exit 54;;
    *)          exit 55;;
    esac
    ;;
"RSYNC_OVER_SSH")
    case $target_mode in
    "LOCAL")
        echo "RSYNC_OVER_SSH to LOCAL" >> ${WORK_DIR}/custom-out.txt
        source_ssh_option=$(get_ssh_option "$2" "source_url" "source_jump_servers")
        local_file=$(get_host_path "$target_dir/$1")
        
        remote_file="$source_dir/$1"

        if [[ $target_url == /data/* ]]; then
            dest_dir=$(dirname ${target_url}/$1)
        else
            dest_dir=$(dirname "/local"${target_url}/$1)
        fi
        echo "dest_dir:$dest_dir" >> ${WORK_DIR}/custom-out.txt
        mkdir -p ${dest_dir}
        cd ${dest_dir}
        source_ssh_url=$(to_ssh_url $source_url)
        cmd="rsync -ut ${rsync_args} -e \"ssh ${source_ssh_option}\" $source_ssh_url/$1 ${dest_dir}"
        echo "cmd=$cmd" >> ${WORK_DIR}/custom-out.txt
        if [ -n "$TASK_TIMEOUT_SECONDS" ]; then
            # 若timeout超时，返回124编码。否则实际完成后，返回0。导致后续脚本错误，误删除出错的源文件
            cmd="timeout ${TASK_TIMEOUT_SECONDS}s $cmd"
        fi
        eval $cmd
        # rsync -ut ${rsync_args} -e "ssh ${ssh_option}" $source_url/$1 ${dest_dir}
        code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from rsync-over-ssh to local, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/custom-out.txt && exit $code
        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            cmd="ssh ${source_ssh_option} $(get_ssh_host $source_url) rm -f $remote_file"
            echo cmd_remove_source_file: $cmd >> ${WORK_DIR}/custom-out.txt
            eval $cmd
            [[ $? -ne 0 ]] && echo "[WARN] error while remove remote source file :$source_file" >> ${WORK_DIR}/custom-out.txt
        fi
    ;;
    "SSH")      exit 66 ;;
    "RSYNC")    exit 67 ;;
    *)          exit 68 ;;
    esac
    ;;
"RSYNC")
    case $target_mode in
    "LOCAL") 
        dest_dir=$(dirname "/local"${target_dir}/$1)
        mkdir -p ${dest_dir}
        # full_file_name=${dest_dir}/$(basename $1)
        rsync_url=$(to_rsync_url "$source_url")
        export RSYNC_PASSWORD=cas12345
        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            rsync_args=" --remove-source-files ${rsync_args} "
        fi
        cmd="rsync -ut ${rsync_args} $rsync_url/$1 ${dest_dir}"
        echo "cmd:$cmd"
        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to local, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/custom-out.txt && exit $code
    ;;
    "SSH")      exit 77 ;;
    "RSYNC")    exit 78 ;;
    *)          exit 79 ;;
    esac
    ;;
*)          exit 80 ;;
esac

date --iso-8601=ns | sed 's/,/./' >> ${WORK_DIR}/timestamps.txt

if [ -n "$SINK_JOB" ]; then
    echo "$1" > ${WORK_DIR}/messages.txt
fi
exit 0
