#!/bin/bash

echo "check_test done!" >> /work/auxout.txt

result_txt=$(cat /work/auxout.txt)
scalebox app set-finished  "$result_txt"

exit 0
