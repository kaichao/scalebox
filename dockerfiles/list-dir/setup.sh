#!/bin/bash

# curlftpfs -f -v -o debug,ftpfs_debug=3 -o allow_other -o ssl ${FTP_URL} /remote

if [[ $SOURCE_URL =~ (ftp://([^:]+:[^@]+@)?[^/:]+(:[^/]+)?)(/.*) ]]; then
    ftp_url=${BASH_REMATCH[1]}
    echo ftp_url:$ftp_url >&2
    curlftpfs -o ssl ${ftp_url} /remote
else
    echo "SOURCE_URL did not match regex!" >&2
fi
