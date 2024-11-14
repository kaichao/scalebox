#!/bin/bash

m=$1
# random number: [0..2]
# exit_code=$((${RANDOM}%3))

exit_code=$(($1%4))
if [ "$m" = "0" ]; then
    for ((i = 1; i < 4; i++));do
        send-message $i
    done
fi

exit $m
