#!/bin/bash

echo "check_test done!" >> /work/custom-out.txt

result_txt=`cat /work/custom-out.txt`
scalebox app set-finished --job-id ${JOB_ID} "$result_txt"

exit 0
