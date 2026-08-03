#!/bin/bash

echo "In default.sh"
body=$1

if [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    scalebox task add --sink-module=wait-queue $body
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    scalebox task add --sink-module=wait-queue $body
else
    scalebox task add --sink-module=vtask-head $body
fi

exit $?
