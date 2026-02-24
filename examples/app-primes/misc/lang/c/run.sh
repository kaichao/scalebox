#!/bin/bash

set -o pipefail
num=$(prime $1|tail -1)
ret_code=$?
set +o pipefail

send-message $1,$num

exit $ret_code
