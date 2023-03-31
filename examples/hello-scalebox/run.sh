#!/bin/bash

echo "Input message:"$1
echo "Hello, $1!"

scalebox app set-finished -job-id=${JOB_ID} "Hello, Scalebox is OK!"

exit 0
