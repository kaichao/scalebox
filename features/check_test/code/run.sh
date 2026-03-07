#!/bin/bash

echo "check_test done!" >> /work/custom-out.txt

result_txt=$(cat /work/auxout.txt)
scalebox app set-finished --module-id ${PLAT_MODULE_ID} "$result_txt"

exit 0
