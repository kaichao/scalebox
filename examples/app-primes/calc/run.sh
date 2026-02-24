#!/bin/bash

# cd /app/bin

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 <start_end>"
    exit 1
fi

# 输入如：001_100
range="$1"

# 用 IFS 拆分
IFS='_' read -r start_str end_str <<< "$range"

# 去掉前导零，转为整数
start=$((10#$start_str))
end=$((10#$end_str))
len=$(( end - start + 1 ))

set -o pipefail
num=$(/app/bin/primes $start $len |tail -1)
code=$?
set +o pipefail

if [ $code -eq 0 ]; then
    scalebox task add -H part_primes=$num $1
    code=$?
    echo "exit_code:$code"
fi

exit $code
