#!/bin/bash

cd /app/bin

set -o pipefail
num=$(/app/bin/primes $1 $LENGTH |tail -1)
code=$?
set +o pipefail

if [ $code -eq 0 ]; then
    scalebox task add -h part_primes=$num $1
    code=$?
    echo "exit_code:$code"
fi

exit $code
