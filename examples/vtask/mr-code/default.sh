#!/bin/bash

echo "In default.sh"
body=$1

if [ "$TASK_DIST_MODE" = "" ]; then
    scalebox task add --sink-module=vtask-head $body
elif [ "$TASK_DIST_MODE" = "HOST-BOUND" ]; then
    scalebox task add --sink-module=wait-queue $body
elif [ "$TASK_DIST_MODE" = "SLOT-BOUND" ]; then
    scalebox task add --sink-module=wait-queue $body
else
    # SLOT-BOUND / HOST-BOUND
    echo "Invalid environment variable TASK_DIST_MODE:${TASK_DIST_MODE}" >&2
    exit 103
fi

exit $?
