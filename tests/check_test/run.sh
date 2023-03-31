#!/bin/bash

echo "check_test done!" >> /work/user-file.txt

result_txt=`cat /work/user-file.txt`
scalebox app set-finished --job-id ${JOB_ID} "$result_txt"

exit 0
