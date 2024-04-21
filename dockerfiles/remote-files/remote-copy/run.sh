#!/bin/bash

# set -e
source /usr/local/bin/functions.sh
source /app/share/bin/functions.sh

declare -F

# Support singularity
[[ ! $WORK_DIR ]] && { echo "[ERROR] WORK_DIR is null, Check the permissions of the directory /tmp/scalebox." >&2; exit 110; }
echo "[DEBUG] WORK_DIR:${WORK_DIR}:" >&2
cd "${WORK_DIR}"

source_url=$(get_parameter "$2" "source_url")
target_url=$(get_parameter "$2" "target_url")

source_mode=$(get_mode $source_url)
target_mode=$(get_mode $target_url)

source_dir=$(get_data_root $source_url)
target_dir=$(get_data_root $target_url)

echo "[DEBUG]source_url:$source_url,target_url:$target_url,message:$1" >> custom-out.txt
echo "[DEBUG]source_dir:$source_dir,target_dir:$target_dir" >> custom-out.txt

echo "[DEBUG]source_mode:$source_mode,target_mode:$target_mode" >> custom-out.txt

date --iso-8601=ns >> timestamps.txt


case $source_mode in
"LOCAL")
    case $target_mode in
    "LOCAL")    exit 31 ;;
    "SSH")
        echo "LOCAL-SSH" >> custom-out.txt
        ssh_cmd=$(get_ssh_cmd "$2" "target_url" "target_jump_servers")
        echo "[DEBUG] ssh_cmd:$ssh_cmd" >> custom-out.txt
        local_file="/local$source_dir/$1"
        remote_file="$target_dir/$1"

        # create directory in remote side.
        remote_dir=$(dirname $remote_file)
        eval "$ssh_cmd mkdir -p $remote_dir"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in remote dir,dir_name:$remote_dir, error_code:$code" >&2 && exit $code

        cmd="cat ${local_file} | pv | $ssh_cmd \"cat > $remote_file\""
        echo "[DEBUG] cmd:$cmd" >> ${WORK_DIR}/custom-out.txt
        eval "$cmd"; code=$?
        [[ $code -ne 0 ]] && echo "[ERROR] cp file from local to remote, cmd=$cmd, error_code:$code" >&2 && exit $code

        if [ "$KEEP_SOURCE_FILE" = "no" ]; then
            echo $local_file >> ${WORK_DIR}/removed-files.txt
        fi
        echo $local_file >> ${WORK_DIR}/input-files.txt
        echo $local_file >> ${WORK_DIR}/output-files.txt

        echo "[DEBUG] local_file:$local_file" >> ${WORK_DIR}/custom-out.txt
        ;;
    "RSYNC-OVER-SSH") exit 33;;
    "RSYNC") exit 34;;
    *)      exit 35 ;;
    esac
    ;;
"SSH")
    case $target_mode in
    "LOCAL")
        echo "SSH-LOCAL" >> custom-out.txt
        ssh_cmd=$(get_ssh_cmd "$2" "source_url" "source_jump_servers")
        echo "[DEBUG] ssh_cmd:$ssh_cmd" >> ${WORK_DIR}/custom-out.txt
        local_file="/local$target_dir/$1"
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
        echo "SSH-SSH" >> custom-out.txt
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
    "RSYNC")    exit 34 ;;
    *)          exit 35;;
    esac
    ;;
"RSYNC-OVER-SSH")
    case $target_mode in
    "LOCAL") exit 36 ;;
    "SSH")      exit 36 ;;
    "RSYNC")    exit 37 ;;
    *)          exit 38 ;;
    esac
    ;;
"RSYNC")
    case $target_mode in
    "LOCAL") exit 36 ;;
    "SSH")      exit 36 ;;
    "RSYNC")    exit 37 ;;
    *)          exit 38 ;;
    esac
    ;;
*)          exit 39 ;;
esac

date --iso-8601=ns >> ${WORK_DIR}/timestamps.txt

# Delay writing messages
echo $1 > ${WORK_DIR}/messages.txt

exit $code
