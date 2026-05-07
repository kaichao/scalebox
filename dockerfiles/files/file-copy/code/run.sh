#!/usr/bin/env bash

set -e

source /usr/local/bin/functions.sh
source /app/share/bin/functions.sh
source /app/share/bin/common.sh
source /app/share/bin/transfer.sh

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
dispatch_func="transfer_${source_mode}_${target_mode}"
if declare -f "$dispatch_func" > /dev/null 2>&1; then
    "$dispatch_func" "$1" "$2"
else
    echo "[ERROR] Unsupported mode: ${source_mode} -> ${target_mode}" >&2
    exit 40
fi

date --iso-8601=ns | sed 's/,/./' >> ${WORK_DIR}/timestamps.txt

if [ -n "$SINK_MODULE" ]; then
    echo "$1" > ${WORK_DIR}/sink-tasks.txt
fi
exit 0
