#!/bin/bash

data_dir=$1

if [ "${MODE}" = "SSH" ]; then
    # "ssh version",  sed-match does not work under macos(Linux-only)
    if [ $data_dir = "." ]; then
        rsync_url=${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_ROOT}
    else 
        rsync_url=${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_ROOT}/${data_dir}
    fi
    rsync -avn -e "ssh -p ${REMOTE_PORT}" ${rsync_url} \
        | grep ^\- | awk {'print $5'}  \
        | sed 's/^[^/]\+\//\//' \
        | awk -v p="$data_dir" '$0=p$0' \
        | egrep "${REGEX_FILTER}" 
elif [ "${MODE}" = "RSYNC" ]; then
    if [ -z ${REMOTE_USER} ]; then
        rsync_url="rsync://"
    else
        rsync_url="rsync://"${REMOTE_USER}@
    fi
    if [ $data_dir = "." ]; then
        rsync_url=${rsync_url}${REMOTE_HOST}:${REMOTE_PORT}${REMOTE_ROOT}
    else 
        rsync_url=${rsync_url}${REMOTE_HOST}:${REMOTE_PORT}${REMOTE_ROOT}/${data_dir}
    fi
    echo "[INFO]rsync_url:"$rsync_url >&2

    # "rsync version"
    rsync -avn ${rsync_url} \
        | grep ^\- | awk {'print $5'} \
        | sed 's/^[^/]\+\//\//' \
        | awk -v p="$data_dir" '$0=p$0' \
        | egrep "${REGEX_FILTER}" 
# rsync: [Receiver] read error: Connection reset by peer (104)
# rsync error: error in socket IO (code 10) at io.c(806) [Receiver=3.2.7]
# [INFO]pipe_status:10 0 0 0 0 0

else
    if [ "${MODE}" = "FTP" ]; then

        # if ! mountpoint -q /remote; then
        #     echo "FTP_URL:"$FTP_URL >&2
        #     if [[ $SOURCE_URL =~ (ftp://([^:]+:[^@]+@)?[^/:]+(:[^/]+)?)(/.*) ]]; then
        #         ftp_url=${BASH_REMATCH[1]}
        #         echo ftp_url:$ftp_url >&2
        #         curlftpfs -o ssl ${ftp_url} /remote
        #     else
        #         echo "FTP_URL did not match regex!" >&2
        #     fi
        # fi
        # # curlftpfs -o ssl ${FTP_URL} /remote
        LOCAL_ROOT="/remote"${REMOTE_ROOT}
    fi
    # MODE = 'LOCAL'
    cd ${LOCAL_ROOT} && find ${data_dir} -type f \
        | sed 's/^\.\///' \
        | egrep "${REGEX_FILTER}"
fi

# exit status of egrep
#   0 if a line is selected
#   1 if no lines were selected
#   2 if an error occurred.
status=(${PIPESTATUS[@]})
echo "[INFO]pipe_status:"${status[*]} >&2
n=${#status[*]}
if [ $n == 1 ]; then
    # line 38
    if [ ${status[0]} -ne 0 ]; then
        echo "[ERROR]local mode, dir: "${LOCAL_ROOT}" not found" >&2
        exit ${status[0]}
    fi
fi

declare -i code
for ((i=n-1; i>=0; i--)); do
    code=${status[i]}
    if [ $i == $((n-1)) ]; then
        if [ $code == 1 ]; then
            echo "[WARN]All of data are filtered, empty dataset!" >&2
            code=0
        fi
    fi
    if [ $code -ne 0 ]; then
        break
    fi
done

exit $code
