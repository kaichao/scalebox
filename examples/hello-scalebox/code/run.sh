#!/bin/bash

echo "task-body:$1"

scalebox app set-status --status=FINISHED --result="Hello $1, it is OK!"

exit $?
