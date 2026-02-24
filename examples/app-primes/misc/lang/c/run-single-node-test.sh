#!/bin/bash

# h=$(hostname)
# cd /gfsdata/primes/result/${h:1}
export OMP_NUM_THREADS=1
export GROUP_SIZE=10000
$(dirname $0)/primes-grouped 1 $1 $GROUP_SIZE
