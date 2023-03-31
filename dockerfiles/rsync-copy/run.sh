#!/bin/bash

# --cc=xxh3 , xxh3 Hash
# --compress --compress-choice=zstd

. /env.sh

#  LOCAL_ROOT/REMOTE_ROOT have been added prefix '/local'
# env

if [[ $REMOTE_MODE == 'SSH' ]]; then 
    if [[ $ENABLE_ZSTD == 'yes' ]]; then 
        if test -f "/rsync_ver_ge_323"; then
            rsync_args="--cc=xxh3 --compress --compress-choice=zstd --compress-level=${ZSTD_CLEVEL}"
        fi
    fi
    ssh_args="-T -c aes128-gcm@openssh.com -o Compression=no -x"
    #  ssh args in /root/.ssh/config
    rsync_url=${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_ROOT}
else
    if [[ $ENABLE_ZSTD == 'yes' ]]; then 
        rsync_args="--cc=xxh3 --compress --compress-choice=zstd --compress-level=${ZSTD_CLEVEL}"
    fi

    if [ -z ${REMOTE_USER} ]; then
        rsync_url="rsync://"
    else
        rsync_url="rsync://"${REMOTE_USER}@
    fi
    rsync_url=${rsync_url}${REMOTE_HOST}:${REMOTE_PORT}${REMOTE_ROOT}
fi

echo [DEBUG]rsync_url:$rsync_url >&2

declare -i total_bytes=0 bytes
arr=($(echo $1 | tr "," " "))
# multiple files to copy
ds0=$(date --iso-8601=ns)
for file in ${arr[@]}; do
    echo "copying file:"$file
    if [[ $ACTION == 'PUSH' ]]; then
        cd ${LOCAL_ROOT}
        if [[ $REMOTE_MODE == 'SSH' ]]; then
            rsync -Rut ${rsync_args} -e "ssh -p ${REMOTE_PORT} ${ssh_args}" $file $rsync_url
        else
            rsync -Rut ${rsync_args} $file $rsync_url
        fi
        code=$?
    else
        dest_dir=$(dirname ${LOCAL_ROOT}/$file)
        mkdir -p ${dest_dir}
        if [[ $REMOTE_MODE == 'SSH' ]]; then
            rsync -ut ${rsync_args} -e "ssh -p ${REMOTE_PORT} ${ssh_args}" $rsync_url/$file ${dest_dir}
            # rsync -Rut ${rsync_args} -e "ssh -p ${REMOTE_PORT} ${ssh_args}" $rsync_url/$file ${LOCAL_ROOT}
            code=$?
        else
            if [[ $ENABLE_LOCAL_RELAY == 'yes' ]]; then
                rsync -ut ${rsync_args} $rsync_url/$file /tmp
                code=$?
                if [ $code -eq 0 ]; then
                    mv /tmp/$(basename $file) ${dest_dir}
                    code=$?
                fi
            else
                rsync -ut ${rsync_args} $rsync_url/$file ${dest_dir}
                code=$?
            fi
        fi
    fi

    if [ $code -ne 0 ]; then
        if [ $code -eq 23 ];then
            # rsync error: some files/attrs were not transferred (see previous errors) (code 23) at main.c(1819) [generator=3.2.3]
            code=100
        elif [ $code -eq 11 ];then
            # Input/output error (5)
            # rsync error: error in file IO (code 11) at receiver.c(871) [receiver=3.2.3]
            code=101
        elif [ $code -eq 255 ];then
            # ssh: connect to host 60.245.209.223 port 22: Connection timed out
            # rsync: connection unexpectedly closed (0 bytes received so far) [sender]
            # rsync error: unexplained error (code 255) at io.c(231) [sender=3.2.6]
            code=101
        else
            echo ret_code=$code
            # code == 10
            # rsync: [Receiver] failed to connect to 10.169.0.68 (10.169.0.68): Connection timed out (110)
            # rsync error: error in socket IO (code 10) at clientserver.c(138) [Receiver=3.2.6]
        fi

        break
    fi
    bytes=$(stat --printf="%s" ${LOCAL_ROOT}/$file)
    ((total_bytes=total_bytes+bytes))
done
ds1=$(date --iso-8601=ns)

if [ $code -eq 0 ]; then
cat << EOF > /work/task-exec.json
{
	"inputBytes":$total_bytes,
	"outputBytes":$total_bytes,
    "timestamps":["${ds0}","${ds1}"]
}
EOF
    send-message $1
    code=$?
else
cat << EOF > /work/task-exec.json
{
	"statusCode":$code
}
EOF
fi

if [ $code -lt -127 ];then
    echo "actual ret_code:"$code
    ret_code=-127
elif [ $code -gt 127 ];then
    echo "actual ret_code:"$code
    ret_code=127
else
    ret_code=$code
fi

exit $ret_code
