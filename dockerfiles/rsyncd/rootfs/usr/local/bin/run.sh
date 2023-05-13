#!/bin/bash

set -e

if [ $MODE = "rsyncd" ]; then
    echo "starting rsyncd on port $RSYNC_PORT"
    rsync --daemon --port $RSYNC_PORT --no-detach --log-file /dev/stdout 
else 
    echo "starting sshd on port $SSH_PORT"
    /usr/sbin/sshd -p $SSH_PORT -D
fi 
