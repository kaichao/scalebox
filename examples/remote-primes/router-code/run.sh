#!/bin/bash

source /usr/local/bin/functions.sh

# set -euo pipefail

from_module=$(get_header "$2" "from_module")
case $from_module in
    "calc")
        # 来自本地计算模块
        part_primes=$(get_header "$2" "part_primes")
        app_id=$(get_header "$2" "app_id")
        scalebox task add --app-id="$app_id" --remote-cluster="cluster0" --header part_primes="$part_primes" "$1"
        ;;
    *)  
        # null from_module（远端主集群）
        echo "calc,$1" > "${WORK_DIR}/sink-tasks.txt"
        ;;
esac

exit $?
