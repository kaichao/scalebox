#!/bin/bash

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <part_primes>"
    exit 1
fi

set -e
export SEMAPHORE_AUTO_CREATE=yes
val=$(scalebox semaphore increment-n app-primes:sum_value $1)
code=$?

echo "part_sum=${val}"
exit $code
