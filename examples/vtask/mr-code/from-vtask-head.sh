#!/bin/bash

echo "In from-vtask-head.sh"

# echo "vtask-core,$1" > $WORK_DIR/sink-tasks.txt
scalebox task add --sink-module=vtask-core $1

exit 0
