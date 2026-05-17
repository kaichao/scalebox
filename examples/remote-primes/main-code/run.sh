#!/bin/bash

source /usr/local/lib/scalebox/functions.sh

# set -euo pipefail

code_dir=$(dirname $0)

echo "num-args:$#"
echo "args:$*"

from_module=$(scalebox::task_header "$2" "from_module")

case $from_module in
    "router")
        # remote-side
        part_primes=$(scalebox::task_header "$2" "part_primes")
        "${code_dir}/sum.sh" "$part_primes"
        ;;
    *)  
        "${code_dir}/split.sh" "$1"
        ;;
esac

exit $?
