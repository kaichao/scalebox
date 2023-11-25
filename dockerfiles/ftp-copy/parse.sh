#!/bin/bash

s=$1
if [[ $s =~ ^(/.*)$ ]]; then
    # /my-root
    echo "LOCAL $1"
elif [[ $s =~ ^(ftp://([^@:]+(:[^@]+)?@)??[^:/]+(:[0-9]+)?)(/.*)$ ]]; then
# ^ftp://([^@:]+(:[^@]+)?@)?[^:/]+(:[0-9]+)?(/.*)$
    # ftp://user:pass@myhost:21/my-root
    # ftp://user@myhost/my-root
    ftp_url=${BASH_REMATCH[1]}
    ftp_root=${BASH_REMATCH[5]}
    echo "FTP $ftp_url $ftp_root"
else
    echo "wrong message format, message:"$1 >&2
    exit 26
fi

exit 0

# if [[ $REMOTE_URL =~ (ftp://([^:]+:[^@]+@)?[^/:]+(:[^/]+)?)(/.*) ]]; then
#     ftp_url=${BASH_REMATCH[1]}
#     remote_root=${BASH_REMATCH[4]}
#     local_root=
# else
#     echo "REMOTE_URL did not match regex! exit_code:6" >&2
#     exit 6
# fi

