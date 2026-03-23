#!/bin/bash

source /usr/local/bin/functions.sh

# set -euo pipefail

code_dir=$(dirname $0)

echo "num-args:$#"
echo "args:$*"

from_module=$(get_header "$2" "from_module")
app_id=$(get_header "$2" "app_id")

case $from_module in
    "router")
        # remote-side
        part_primes=$(get_header "$2" "part_primes")
        "${code_dir}/sum.sh" "$part_primes"
        ;;
    *)  
        "${code_dir}/split.sh" "$1" "$app_id"
        ;;
esac

exit $?
