#!/usr/bin/env bash

# 传输函数，每个函数处理一种 source_mode -> target_mode 组合
# 参数: $1 = 文件名 (task body), $2 = JSON headers
# 依赖全局变量: WORK_DIR, source_dir, target_dir, keep_source, rsync_args, source_url, target_url
# 依赖库: /usr/local/lib/scalebox/

# ===== LOCAL -> LOCAL =====
transfer_LOCAL_LOCAL() {
    echo "[ERROR] both source_mode and target_mode are local" >> ${WORK_DIR}/auxout.txt
    exit 31
}

# ===== LOCAL -> SSH =====
transfer_LOCAL_SSH() {
    local target_ssh_option=$(url::ssh_option "$2" "target_url" "target_jump")

    if [[ $source_url == /data/* ]]; then
        source_dir="${source_url}/"
    else
        source_dir=$(path::host_path "${source_url}/")
    fi
    echo "source_dir:$source_dir" >> ${WORK_DIR}/auxout.txt

    local target_ssh_url=$(url::to_ssh_conn $target_url)
    local target_dir_rsync=$(dirname ${target_ssh_url#*:}/$1)
    local cmd="rsync -Rut ${rsync_args} -e \"ssh ${target_ssh_option}\" --rsync-path=\"mkdir -p ${target_dir_rsync} && rsync\" $1 $target_ssh_url/"

    echo "cmd=$cmd" >> ${WORK_DIR}/auxout.txt
    cd $source_dir && eval $cmd
    local code=$?
    [[ $code -ne 0 ]] && echo "[ERROR] cp file from local to remote, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code

    cleanup_source "LOCAL" "$source_dir/$1" "$keep_source"
    echo "$source_dir/$1" >> "${WORK_DIR}/input-files.txt"
    echo "to,$(url::host "$target_url"),$source_dir/$1,${target_ssh_url#*:}/$1" >> "${WORK_DIR}/network-files.txt"
}

# ===== LOCAL -> RSYNC =====
transfer_LOCAL_RSYNC() {
    export RSYNC_PASSWORD=cas12345

    local rsync_url=$(url::to_rsync_conn "$target_url")

    if [ "$keep_source" = "no" ]; then
        rsync_args=" --remove-source-files ${rsync_args} "
    fi
    local local_dir=$(path::host_path "${source_dir}")
    local cmd="cd $local_dir && rsync -ut ${rsync_args} $1 $rsync_url"
    echo "cmd:$cmd"

    eval "$cmd"; local code=$?
    [[ $code -ne 0 ]] && echo "[ERROR] cp file from local to remote, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code
    local local_file=$(path::host_path "${source_dir}/$1")
    echo "$local_file" >> "${WORK_DIR}/input-files.txt"
    echo "to,$(url::host "$target_url"),$local_file,$target_dir" >> "${WORK_DIR}/network-files.txt"
}

# ===== SSH -> LOCAL =====
transfer_SSH_LOCAL() {
    local source_ssh_option=$(url::ssh_option "$2" "source_url" "source_jump")
    local local_file=$(path::host_path "$target_dir/$1")
    local remote_file="$source_dir/$1"

    if [[ $target_url == /data/* ]]; then
        local dest_dir=$(dirname ${target_url}/$1)
    else
        local dest_dir=$(dirname "/local_data_root"${target_url}/$1)
    fi
    echo "dest_dir:$dest_dir" >> ${WORK_DIR}/auxout.txt
    mkdir -p ${dest_dir}
    cd ${dest_dir}
    local source_ssh_url=$(url::to_ssh_conn $source_url)
    local cmd="rsync -ut ${rsync_args} -e \"ssh ${source_ssh_option}\" $source_ssh_url/$1 ${dest_dir}"
    echo "cmd=$cmd" >> ${WORK_DIR}/auxout.txt
    if [ -n "$TASK_TIMEOUT_SECONDS" ]; then
        cmd="timeout ${TASK_TIMEOUT_SECONDS}s $cmd"
    fi
    eval $cmd
    local code=$?
    [[ $code -ne 0 ]] && echo "[ERROR] cp file from rsync-over-ssh to local, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code
    if [ "$keep_source" = "no" ]; then
        local cmd_rm="ssh ${source_ssh_option} $(url::ssh_host $source_url) rm -f $remote_file"
        echo cmd_remove_source_file: $cmd_rm >> ${WORK_DIR}/auxout.txt
        eval $cmd_rm
        [[ $? -ne 0 ]] && echo "[WARN] error while remove remote source file :$source_file" >> ${WORK_DIR}/auxout.txt
    fi
    if [[ $target_url == /data/* ]]; then
        local local_file="${target_url}/$1"
    else
        local local_file="/local_data_root${target_url}/$1"
    fi
    echo "$local_file" >> "${WORK_DIR}/output-files.txt"
    echo "from,$(url::host "$source_url"),$local_file,$remote_file" >> "${WORK_DIR}/network-files.txt"
}

# ===== SSH -> SSH =====
transfer_SSH_SSH() {
    local source_ssh_cmd=$(url::ssh_cmd "$2" "source_url" "source_jump")
    local target_ssh_cmd=$(url::ssh_cmd "$2" "target_url" "target_jump")

    local source_file="$source_dir/$1"
    local target_file="$target_dir/$1"

    local target_dir_ssh=$(dirname $target_file)
    ensure_remote_dir "$target_ssh_cmd" "$target_dir_ssh"

    local cmd="$source_ssh_cmd \"cat < $source_file\" - | pv -q -n | $target_ssh_cmd \"cat > $target_file\""
    local bytes_transferred=$(eval "$cmd"); local code=$?
    echo "bytes_transferred: $bytes_transferred" >> auxout.txt

    [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to remote, cmd=$cmd, error_code:$code" >&2 && exit $code

    cleanup_source "SSH" "$source_file" "$keep_source" "$source_ssh_cmd"
    echo "from,$(url::host "$source_url"),$bytes_transferred,$source_file" >> "${WORK_DIR}/network-files.txt"
    echo "to,$(url::host "$target_url"),$bytes_transferred,$target_file" >> "${WORK_DIR}/network-files.txt"
}

# ===== RSYNC -> LOCAL =====
transfer_RSYNC_LOCAL() {
    local dest_dir=$(dirname "/local_data_root"${target_dir}/$1)
    mkdir -p ${dest_dir}
    local rsync_url=$(url::to_rsync_conn "$source_url")
    export RSYNC_PASSWORD=cas12345
    if [ "$keep_source" = "no" ]; then
        rsync_args=" --remove-source-files ${rsync_args} "
    fi
    local cmd="rsync -ut ${rsync_args} $rsync_url/$1 ${dest_dir}"
    echo "cmd:$cmd"
    eval "$cmd"; local code=$?
    [[ $code -ne 0 ]] && echo "[ERROR] cp file from remote to local, cmd=$cmd, error_code:$code" >> ${WORK_DIR}/auxout.txt && exit $code
    local local_file="/local_data_root${target_dir}/$1"
    echo "$local_file" >> "${WORK_DIR}/output-files.txt"
    echo "from,$(url::host "$source_url"),$local_file,$target_dir" >> "${WORK_DIR}/network-files.txt"
}
