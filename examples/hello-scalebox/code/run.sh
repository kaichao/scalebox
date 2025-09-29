#!/bin/bash

echo "Input message:$1"

scalebox app set-finished --module-id=${MODULE_ID} "Hello $1, it is OK!"

exit $?
