#!/usr/bin/env bash

# 公共工具函数，供 transfer.sh 中的传输函数调用
# 依赖库: /usr/local/lib/scalebox/

# 创建本地目录，失败时退出
ensure_local_dir() {
    local dir="$1"
    mkdir -p "$dir"
    local code=$?
    [[ $code -ne 0 ]] && echo "[ERROR] mkdir in local dir, dir_name:${dir}, error_code:${code}" >&2 && exit $code
}

# 通过 SSH 创建远程目录，失败时退出
ensure_remote_dir() {
    local ssh_cmd="$1"
    local remote_dir="$2"
    eval "$ssh_cmd mkdir -p $remote_dir"
    local code=$?
    [[ $code -ne 0 ]] && echo "[ERROR] mkdir in remote dir, dir_name:${remote_dir}, error_code:${code}" >&2 && exit $code
}

# keep_source=no 时清理源端文件
# 参数: mode path keep_source [ssh_cmd]
#   mode: LOCAL / SSH / RSYNC
#   RSYNC 模式的清理由 rsync --remove-source-files 完成，此函数无操作
cleanup_source() {
    local mode="$1"
    local path="$2"
    local keep_source="$3"
    local ssh_cmd="${4:-}"

    [[ "$keep_source" != "no" ]] && return

    case "$mode" in
        "LOCAL")
            echo "$path" >> "${WORK_DIR}/removed-files.txt"
            ;;
        "SSH")
            eval "$ssh_cmd rm -f $path"
            ;;
        # RSYNC: cleanup already done by --remove-source-files in rsync command
    esac
}
