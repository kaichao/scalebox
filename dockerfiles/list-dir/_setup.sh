#!/bin/bash

# curlftpfs -f -v -o debug,ftpfs_debug=3 -o allow_other -o ssl ${FTP_URL} /remote
if [[ $SOURCE_URL =~ (ftp://([^:]+:[^@]+@)?[^/:]+(:[^/]+)?)(/.*) ]]; then
    ftp_url=${BASH_REMATCH[1]}
    echo ftp_url:$ftp_url >&2
    if [[ $SOURCE_URL =~ (ftp://([^:]+:[^@]+@)[^/:]+(:[^/]+)?)(/.*) ]]; then
        # non-anonymous ftp
        curlftpfs -o ssl ${ftp_url} /remote
    else
        # anonymous ftp
        curlftpfs ${ftp_url} /remote
    fi
# else
#     echo "[INFO] SOURCE_URL is not valid ftp url."
fi
