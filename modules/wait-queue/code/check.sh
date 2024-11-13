#!/usr/bin/env bash

# set -e
source /usr/local/bin/functions.sh

if [[ -z "$SEMA_NAME" ]]; then
    echo "Null Variable SEMA_NAME"
    exit 11
fi
sema_name=$(get_parameter "$2" "sema_name")
val=$(scalebox semaphore get "${sema_name}")
code=$?
[[ $code -ne 0 ]] && echo "[ERROR] semaphore get ${sema_name}, exit_code:$code" >&2 && exit $code

# (( val > 0 )) number compare
if (( val > 0 )); then
    scalebox semaphore decrement "${sema_name}"
    exit $?
else
    exit 10
fi
