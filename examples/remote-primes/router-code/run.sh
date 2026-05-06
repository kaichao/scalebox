#!/bin/bash

source /usr/local/bin/functions.sh

# set -euo pipefail

from_module=$(get_header "$2" "from_module")
case $from_module in
    "calc")
        echo "from local module calc" >> auxout.txt
        part_primes=$(get_header "$2" "part_primes")
        scalebox task add --remote-cluster="cluster0" --header part_primes="$part_primes" "$1"
        ;;
    *)  
        # from remote module main(cluster cluster0)
        echo "from remote cluster" >> auxout.txt
        echo "calc,$1" > "${WORK_DIR}/sink-tasks.txt"
        ;;
esac

exit $?
