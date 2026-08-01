#!/bin/bash

echo "task-body:$1"

scalebox app set-status 'FINISHED' "Hello $1, it is OK!"

exit $?
