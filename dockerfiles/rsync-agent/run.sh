#!/bin/bash

# FILENAME=$(cat /tasks/$1 | head -n 1)
# FILENAME=/data/ZD2020_1_1_2bit/$1
FILENAME=/data/$1

# REMOTE_HOST=10.255.0.10

set -e

if [ $TRANSPORT_TYPE = "rsync" ]; then 
    # RSYNC_PASSWORD="cnic123" rsync -RPut --port $RSYNC_PORT $FILENAME rsync://root@$RECEIVER_HOST:/share
    RSYNC_PASSWORD="cnic123" rsync -RPut $FILENAME rsync://root@$REMOTE_HOST:$RSYNC_PORT/share
else 
    rsync -RPut -e "ssh -p ${SSH_PORT}  -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" \
        ${SSH_USER}@${REMOTE_HOST}:${FILENAME} /output
    # rsync -a -e "ssh -p ${SSH_PORT}  -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" \
    #     ${SSH_USER}@${REMOTE_HOST}:${FILENAME} /tmp
fi 

echo ret_code=$?

exit $ret_code
