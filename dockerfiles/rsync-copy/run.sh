#!/bin/bash

set -e

# Support singularity
[[ ! $WORK_DIR ]] && { echo "[ERROR] WORK_DIR is null, Check the permissions of the directory /tmp/scalebox." >&2; exit 110; }
echo "[DEBUG] WORK_DIR:${WORK_DIR}:" >&2
cd "${WORK_DIR}"

m=$1
if [ "$SOURCE_URL" == "" ] || [ "$TARGET_URL" == "" ]; then
    if [[ ! $m =~ ^([^~]*)~([^~]*)~([^~]*)$ ]]; then
        echo "wrong message format, message:"$1 >&2
        exit 21
    fi
    s1=${BASH_REMATCH[1]}
    s2=${BASH_REMATCH[2]}
    s3=${BASH_REMATCH[3]}
    source_url=${SOURCE_URL:-$s1}
    target_url=${TARGET_URL:-$s3}
    if [ "$s2" == "" ]; then
        echo "null mesage_body, message:"$1 >&2
        exit 24
    fi
    m=$s2
else
    source_url=$SOURCE_URL
    target_url=$TARGET_URL
fi

arr_source=($(/app/share/bin/parse.sh $source_url))
code=$?
[[ $code -ne 0 ]] && echo error while parse_source_url, error_code:$code, source_url:$source_url >&2 && exit $code
source_mode=${arr_source[0]}
source_url=${arr_source[1]}

arr_target=($(/app/share/bin/parse.sh $target_url))
code=$?
[[ $code -ne 0 ]] && echo error while parse_target_url, error_code:$code, source_url:$target_url >&2 && exit $code
target_mode=${arr_target[0]}
target_url=${arr_target[1]}

ssh_args="-T -c aes128-gcm@openssh.com -o Compression=no -x"
if [ $JUMP_SERVERS ]; then
    ssh_args=$ssh_args" -J '${JUMP_SERVERS}'"
fi
if [[ $ZSTD_CLEVEL != "" ]]; then 
    rsync_args="--cc=xxh3 --compress --compress-choice=zstd --compress-level=${ZSTD_CLEVEL}"
fi

echo "[DEBUG]source_url:$source_url,target_url:$target_url,message:$m"

date --iso-8601=ns >> ${WORK_DIR}/timestamps.txt

case $source_mode in
"LOCAL")
    case $target_mode in
    "LOCAL")    exit 31 ;;
    "SSH")
        ssh_port=${arr_target[2]}
        ssh_args="ssh -p ${ssh_port} ${ssh_args}"
        full_file_name=${dest_dir}/$(basename $m)

        # create directory in target side.
        my_arr=($(echo $target_url | tr ":" " "))

        if [[ $JUMP_SERVERS == "" ]]; then 
            cmd="ssh -p ${ssh_port} ${my_arr[0]} \"mkdir -p ${my_arr[1]}\""
        else
            cmd="ssh -p ${ssh_port} -J $JUMP_SERVERS ${my_arr[0]} \"mkdir -p ${my_arr[1]}\""
        fi

        # echo cmd:$cmd
        eval $cmd; code=$?
        # ssh -p ${ssh_port} ${my_arr[0]} "mkdir -p ${my_arr[1]}"
        [[ $code -ne 0 ]] && echo "[ERROR] mkdir in ssh-server,cmd:$cmd" >&2 && exit $code

        cd "/local"$source_url
        rsync -Rut ${rsync_args} -e "${ssh_args}" $m $target_url
        ;;
    "RSYNC")
        cd "/local"$source_url
        rsync -Rut ${rsync_args} $m $target_url
        ;;
    *)      exit 32 ;;
    esac
    ;;
"SSH")
    case $target_mode in
    "LOCAL")
        ssh_port=${arr_source[2]}
        ssh_args="ssh -p ${ssh_port} ${ssh_args}"
        if [[ $target_url == /data* ]]; then
            dest_dir=$(dirname ${target_url}/$m)
        else
            dest_dir=$(dirname "/local"${target_url}/$m)
        fi
        full_file_name=${dest_dir}/$(basename $m)
        echo dest_dir:$dest_dir
        mkdir -p ${dest_dir}
        rsync -ut ${rsync_args} -e "${ssh_args}" $source_url/$m ${dest_dir}
        ;;
    "SSH")      exit 33 ;;
    "RSYNC")    exit 34 ;;
    *)          exit 35;;
    esac
    ;;
"RSYNC")
    case $target_mode in
    "LOCAL")
        dest_dir=$(dirname "/local"${target_url}/$m)
        mkdir -p ${dest_dir}
        full_file_name=${dest_dir}/$(basename $m)
        rsync -ut ${rsync_args} $target_url/$m ${dest_dir}
        ;;
    "SSH")      exit 36 ;;
    "RSYNC")    exit 37 ;;
    *)          exit 38 ;;
    esac
    ;;
*)          exit 39 ;;
esac

code=$?
date --iso-8601=ns >> ${WORK_DIR}/timestamps.txt

if [ $code -ne 0 ]; then
    if [ $code -eq 23 ];then
        # rsync error: some files/attrs were not transferred (see previous errors) (code 23) at main.c(1819) [generator=3.2.3]
        code=23
    elif [ $code -eq 11 ];then
        # Input/output error (5)
        # rsync error: error in file IO (code 11) at receiver.c(871) [receiver=3.2.3]
        code=11
    elif [ $code -eq 255 ];then
        # ssh: connect to host 60.245.209.223 port 22: Connection timed out
        # rsync: connection unexpectedly closed (0 bytes received so far) [sender]
        # rsync error: unexplained error (code 255) at io.c(231) [sender=3.2.6]
        code=255
    else
        echo ret_code=$code
        # code == 10
        # rsync: [Receiver] failed to connect to 10.169.0.68 (10.169.0.68): Connection timed out (110)
        # rsync error: error in socket IO (code 10) at clientserver.c(138) [Receiver=3.2.6]
    fi
fi
[[ $code -ne 0 ]] && exit $code

# cat << EOF > /work/task-exec.json
# {
#     "inputBytes":$total_bytes,
#     "outputBytes":$total_bytes,
#     "timestamps":["${ds0}","${ds1}"]
# }
# EOF
send-message $m; code=$?

echo source_mode:$source_mode, KEEP_SOURCE_FILE:$KEEP_SOURCE_FILE
full_path_file="/local${source_url}/${m}"
echo full_path_file:$full_path_file
[ "$source_mode" = "LOCAL" ] && [ "$KEEP_SOURCE_FILE" = "no" ] && echo $full_path_file >> ${WORK_DIR}/removed-files.txt

if [ "$source_mode" = "LOCAL" ]; then
    echo $full_path_file >> ${WORK_DIR}/input-files.txt
    echo $full_path_file >> ${WORK_DIR}/output-files.txt
fi

exit $code
