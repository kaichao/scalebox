#!/bin/bash

set -e
source "/usr/local/bin/functions.sh"

n=$(get_header "$2" "part_primes")
echo "part_primes:$n"
export SEMAPHORE_AUTO_CREATE=yes
val=$(scalebox semaphore increment-n app-primes:sum_value ${n})
code=$?

echo "part_sum=${val}"
exit $code
