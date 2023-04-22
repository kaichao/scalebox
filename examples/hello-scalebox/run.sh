#!/bin/bash

echo "Input message:"$1

scalebox app set-finished --job-id=${JOB_ID} "Hello $1, it is OK!"

exit $?
