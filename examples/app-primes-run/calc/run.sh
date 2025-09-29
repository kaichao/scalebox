#!/bin/bash

cd /app/bin


# 001_100
m=$1

# 解析变量m，生成两个整数变量start,len
start="001"
len=100


set -o pipefail
num=$(/app/bin/primes $start $len |tail -1)
code=$?
set +o pipefail

if [ $code -eq 0 ]; then
    scalebox task add -h part_primes=$num $1
    code=$?
    echo "exit_code:$code"
fi

exit $code
