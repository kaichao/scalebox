#!/bin/bash

# echo "vtask-tail,$1" > $WORK_DIR/sink-tasks.txt
scalebox task add --sink-module=vtask-tail "$1"

exit $?
