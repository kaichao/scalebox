#!/usr/bin/env bash

set -e

source /usr/local/bin/functions.sh
source /app/share/bin/functions.sh

# Support singularity
[[ ! $WORK_DIR ]] && { echo "[ERROR] WORK_DIR is null, Check the permissions of the directory /tmp/scalebox." >&2; exit 110; }
echo "[DEBUG] WORK_DIR:${WORK_DIR}:" >> ${WORK_DIR}/auxout.txt
cd "${WORK_DIR}"

source_url=$(get_header "$2" "source_url")
target_url=$(get_header "$2" "target_url")

keep_source=$(get_header "$2" "keep_source")

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


# 如果source_url/target_url/source_mode/target_mode任意一个为空，则显示输入错误，退出，退出码120
if [ -z "$source_url" ] || [ -z "$target_url" ] || [ -z "$source_mode" ] || [ -z "$target_mode" ]; then
    echo "[ERROR] Invalid input: source_url, target_url, source_mode, or target_mode is empty." >&2
    exit 120
fi

echo "[DEBUG]source_url:$source_url,target_url:$target_url,task-body:$1" >> ${WORK_DIR}/auxout.txt
echo "[DEBUG]source_dir:$source_dir,target_dir:$target_dir" >> ${WORK_DIR}/auxout.txt

if [[ $ZSTD_CLEVEL != "" ]]; then
    rsync_args="--cc=xxh3 --compress --compress-choice=zstd --compress-level=${ZSTD_CLEVEL}"
else
    rsync_args=""
fi

if [ "$source_mode" = "RSYNC_OVER_SSH" ] && [ "$target_mode" = "RSYNC_OVER_SSH" ]; then
    source_mode="SSH"
    target_mode="SSH"
fi

date --iso-8601=ns | sed 's/,/./' >> timestamps.txt

echo "[DEBUG]source_mode:$source_mode,target_mode:$target_mode" >> ${WORK_DIR}/auxout.txt
case $source_mode in
"LOCAL")
    case $target_mode in
    "LOCAL")    
        echo "[ERROR] both source_mode and target_mode are local" >> ${WORK_DIR}/auxout.txt
        exit 31 ;;
    "SSH")
        echo "LOCAL-SSH" >> ${WORK_DIR}/auxout.txt
        ssh_cmd=$(get_ssh_cmd "$2" "target_url" "target_jump")
        echo "[DEBUG] ssh_cmd:$ssh_cmd" >> auxout.txt
        local_file=$(get_host_path "$source_dir/$1")

        remote_file="$target_dir/$1"

        # source file not exists ?
        [ ! -f "$local_file" ] && echo "file $local_file not exists, exit " >&2 && exit 101

        # create directory in remote side.
        remote_dir=$(dirname $remote_file)
        eval "$ssh_cmd mkdir -p $remote_dir"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in remote dir,dir_name:$remote_dir, error_code:$code" >&2 && exit $code

        cmd="cat ${local_file} | pv -q | $ssh_cmd \"cat > $remote_file\""
        echo "[DEBUG] cmd:$cmd" >> ${WORK_DIR}/auxout.txt
        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from local to remote, cmd=$cmd, error_code:$code" >&2 && exit $code

        if [ "$keep_source" = "no" ]; then
            echo "$local_file" >> ${WORK_DIR}/removed-files.txt
        fi
        echo "$local_file" >> ${WORK_DIR}/input-files.txt
        echo "$local_file" >> ${WORK_DIR}/output-files.txt

        echo "[DEBUG] local_file:$local_file" >> ${WORK_DIR}/auxout.txt
        ;;
    "RSYNC_OVER_SSH") 
        echo "LOCAL to RSYNC_OVER_SSH"
        target_ssh_option=$(get_ssh_option "$2" "target_url" "target_jump")

        if [[ $source_url == /data/* ]]; then
            source_dir="${source_url}/"
        else
            source_dir=$(get_host_path "${source_url}/")
        fi
        echo "source_dir:$source_dir" >> ${WORK_DIR}/auxout.txt

        target_ssh_url=$(to_ssh_url $target_url)
        target_dir=$(dirname ${target_ssh_url#*:}/$1)
        echo target_dir:$target_dir
        cmd="rsync -Rut ${rsync_args} -e \"ssh ${target_ssh_option}\" --rsync-path=\"mkdir -p ${target_dir} && rsync\" $1 $target_ssh_url/"

        echo "cmd=$cmd" >> ${WORK_DIR}/auxout.txt
        cd $source_dir && eval $cmd
        code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to remote, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code
        if [ "$keep_source" = "no" ]; then
            echo "$source_dir/$1" >> ${WORK_DIR}/removed-files.txt
        fi
    ;;
    "RSYNC")
        export RSYNC_PASSWORD=cas12345

        rsync_url=$(to_rsync_url "$target_url")

        if [ "$keep_source" = "no" ]; then
            rsync_args=" --remove-source-files ${rsync_args} "
        fi
        local_dir=$(get_host_path "${source_dir}")
        cmd="cd $local_dir && rsync -ut ${rsync_args} $1 $rsync_url"
        echo "cmd:$cmd"

        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from local to remote, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code
    ;;
    *)      echo "[ERROR] Unsupported mode: LOCAL -> $target_mode" >&2; exit 40 ;;
    esac
    ;;
"SSH")
    case $target_mode in
    "LOCAL")
        echo "SSH to LOCAL"
        ssh_cmd=$(get_ssh_cmd "$2" "source_url" "source_jump")
        echo "[DEBUG] ssh_cmd:$ssh_cmd" >> ${WORK_DIR}/auxout.txt
        local_file=$(get_host_path "$target_dir/$1")
        remote_file="$source_dir/$1"

        # create directory in local side.
        local_dir=$(dirname $local_file)
        mkdir -p $local_dir; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in local dir,dir_name:$local_dir, error_code:$code" >&2 && exit $code

        cmd="$ssh_cmd \"cat < $remote_file\" - | pv -q > ${local_file}"
        echo "[DEBUG] cmd:$cmd" >> ${WORK_DIR}/auxout.txt
        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to local, cmd=$cmd, error_code:$code" >&2 && exit $code

        if [ "$keep_source" = "no" ]; then
            eval "$ssh_cmd rm -f $remote_file"
        fi
        echo $local_file >> ${WORK_DIR}/input-files.txt
        echo $local_file >> ${WORK_DIR}/output-files.txt
        ;;
    "SSH")
        echo "SSH to SSH"
        source_ssh_cmd=$(get_ssh_cmd "$2" "source_url" "source_jump")
        target_ssh_cmd=$(get_ssh_cmd "$2" "target_url" "target_jump")

        source_file="$source_dir/$1"
        target_file="$target_dir/$1"

        # create directory in remote side.
        target_dir=$(dirname $target_file)
        eval "$target_ssh_cmd mkdir -p $target_dir"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in remote dir,dir_name:$remote_dir, error_code:$code" >&2 && exit $code

        cmd="$source_ssh_cmd \"cat < $source_file\" - | pv -q -n | $target_ssh_cmd \"cat > $target_file\""
        bytes_transferred=$(eval "$cmd"); code=$?
        echo "bytes_transferred: $bytes_transferred" >> auxout.txt

        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to remote, cmd=$cmd, error_code:$code" >&2 && exit $code

        if [ "$keep_source" = "no" ]; then
            eval "$source_ssh_cmd rm -f $source_file"
        fi
        # ssh user1@node1 "cat /path/to/file" | ssh user2@node2 "cat > /path/to/destination"
        ;;
    "RSYNC")    echo "[ERROR] Unsupported mode: SSH -> RSYNC" >&2; exit 40 ;;
    *)          echo "[ERROR] Unsupported mode: SSH -> $target_mode" >&2; exit 40 ;;
    esac
    ;;
"RSYNC_OVER_SSH")
    case $target_mode in
    "LOCAL")
        echo "RSYNC_OVER_SSH to LOCAL"
        source_ssh_option=$(get_ssh_option "$2" "source_url" "source_jump")
        local_file=$(get_host_path "$target_dir/$1")
        
        remote_file="$source_dir/$1"

        if [[ $target_url == /data/* ]]; then
            dest_dir=$(dirname ${target_url}/$1)
        else
            dest_dir=$(dirname "/local_data_root"${target_url}/$1)
        fi
        echo "dest_dir:$dest_dir" >> ${WORK_DIR}/auxout.txt
        mkdir -p ${dest_dir}
        cd ${dest_dir}
        source_ssh_url=$(to_ssh_url $source_url)
        cmd="rsync -ut ${rsync_args} -e \"ssh ${source_ssh_option}\" $source_ssh_url/$1 ${dest_dir}"
        echo "cmd=$cmd" >> ${WORK_DIR}/auxout.txt
        if [ -n "$TASK_TIMEOUT_SECONDS" ]; then
            # 若timeout超时，返回124编码。否则实际完成后，返回0。导致后续脚本错误，误删除出错的源文件
            cmd="timeout ${TASK_TIMEOUT_SECONDS}s $cmd"
        fi
        eval $cmd
        # rsync -ut ${rsync_args} -e "ssh ${ssh_option}" $source_url/$1 ${dest_dir}
        code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from rsync-over-ssh to local, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code
        if [ "$keep_source" = "no" ]; then
            cmd="ssh ${source_ssh_option} $(get_ssh_host $source_url) rm -f $remote_file"
            echo cmd_remove_source_file: $cmd >> ${WORK_DIR}/auxout.txt
            eval $cmd
            [[ $? -ne 0 ]] && echo "[WARN] error while remove remote source file :$source_file" >> ${WORK_DIR}/auxout.txt
        fi
    ;;
    "SSH")      echo "[ERROR] Unsupported mode: RSYNC_OVER_SSH -> SSH" >&2; exit 40 ;;
    "RSYNC")    echo "[ERROR] Unsupported mode: RSYNC_OVER_SSH -> RSYNC" >&2; exit 40 ;;
    *)          echo "[ERROR] Unsupported mode: RSYNC_OVER_SSH -> $target_mode" >&2; exit 40 ;;
    esac
    ;;
"RSYNC")
    case $target_mode in
    "LOCAL") 
        dest_dir=$(dirname "/local_data_root"${target_dir}/$1)
        mkdir -p ${dest_dir}
        # full_file_name=${dest_dir}/$(basename $1)
        rsync_url=$(to_rsync_url "$source_url")
        export RSYNC_PASSWORD=cas12345
        if [ "$keep_source" = "no" ]; then
            rsync_args=" --remove-source-files ${rsync_args} "
        fi
        cmd="rsync -ut ${rsync_args} $rsync_url/$1 ${dest_dir}"
        echo "cmd:$cmd"
        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to local, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code
    ;;
    "SSH")      echo "[ERROR] Unsupported mode: RSYNC -> SSH" >&2; exit 40 ;;
    "RSYNC")    echo "[ERROR] Unsupported mode: RSYNC -> RSYNC" >&2; exit 40 ;;
    *)          echo "[ERROR] Unsupported mode: RSYNC -> $target_mode" >&2; exit 40 ;;
    esac
    ;;
*)          echo "[ERROR] Unsupported source_mode: $source_mode" >&2; exit 40 ;;
esac

date --iso-8601=ns | sed 's/,/./' >> ${WORK_DIR}/timestamps.txt

if [ -n "$SINK_MODULE" ]; then
    echo "$1" > ${WORK_DIR}/sink-tasks.txt
fi
exit 0
