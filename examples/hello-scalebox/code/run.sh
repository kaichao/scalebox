#!/bin/bash

echo "task-body:$1"

scalebox app set-finished "Hello $1, it is OK!"

exit $?
