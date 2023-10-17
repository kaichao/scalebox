#!/bin/bash

dir=$1
if [ -z "$SOURCE_URL" ]; then
    arr=($(echo $1 | tr "%" " ")) 
    if [ ${#arr[@]} -ge 2 ]; then
        export SOURCE_URL="${arr[0]}"
        dir="${arr[1]}"
        percent_seperated="yes"
    elif [ ${#arr[@]} -eq 1 ]; then
        echo "SOURCE_URL is null, and not included in the input message" >&2
        exit 11
    fi
fi

if [[ $SOURCE_URL =~ (ftp://([^:]+:[^@]+@)?[^/:]+(:[^/]+)?)(/.*) ]]; then
# curlftpfs -f -v -o debug,ftpfs_debug=3 -o allow_other -o ssl ${FTP_URL} /remote
    ftp_url=${BASH_REMATCH[1]}
    echo ftp_url:$ftp_url >&2
    source_ftp="yes"
    if [[ $SOURCE_URL =~ (ftp://([^:]+:[^@]+@)[^/:]+(:[^/]+)?)(/.*) ]]; then
        # non-anonymous ftp
        curlftpfs -o ssl ${ftp_url} /remote
    else
        # anonymous ftp
        curlftpfs ${ftp_url} /remote
    fi
fi

arr=($(/usr/local/bin/url_parser ${SOURCE_URL}))
code=$?
if [ $code -ne 0 ]; then
    # url format error
    exit $code
fi

export MODE="${arr[0]}"
if [ "${MODE}" = "ERROR_FORMAT" ]; then
    echo "url format error, "${arr[1]} >&2
    exit 1
fi

if [ "${MODE}" = "LOCAL" ]; then
	#  MODE | LOCAL_ROOT
    export LOCAL_ROOT="${arr[1]}"
elif [ "${MODE}" = "FTP" ]; then
    export FTP_URL="${arr[1]}"
    export REMOTE_ROOT="${arr[2]}"
else
    # MODE | REMOTE_HOST | REMOTE_PORT |  REMOTE_ROOT | REMOTE_USER 
    export REMOTE_HOST="${arr[1]}"
    export REMOTE_PORT="${arr[2]}"
    export REMOTE_ROOT="${arr[3]}"
    if [ ${#arr[*]} -eq 5 ]; then
        export REMOTE_USER="${arr[4]}"
    else
        export REMOTE_USER=
    fi
fi

env

if [ "$REGEX_2D_DATASET" ]; then
    meta=$(/usr/local/bin/list-files.sh $dir | get_2d_meta $REGEX_2D_DATASET $INDEX_2D_DATASET)
    code=$?
    [[ $code -ne 0 ]] && echo cmd: get_2d_meta, error_code:$code && exit $code
    echo ${meta} > /work/key-text.txt
    echo metadata for 2d-dataset:#${meta}#
    # key text in file /work/key-text.txt
    scalebox task add
    code=$?
    [[ $code -ne 0 ]] && echo cmd: scalebox task add, error_code:$code && exit $code
fi

ret_code=0
/usr/local/bin/list-files.sh $dir | while read line; do 
    if [[ $line == ./* ]]; then
        # remove prefix './'
        line=${line:2}
    fi
    if [ "${percent_seperated}" = "yes" ]; then
        line=${SOURCE_URL}%${line}
    fi

    send-message $line
    code=$?
    if [ $code -ne 0 ]; then
        ret_code=$code
        echo "Error send-message, message:"$line >&2
    fi
done

code=${PIPESTATUS[0]}
if [ $code -ne 0 ]; then
    ret_code=$code
    echo "Error run list-files.sh "$dir >&2
fi

if [ "${source_ftp}" = "yes" ]; then
    umount /remote
fi

exit $ret_code
