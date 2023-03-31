#!/bin/bash

echo "SLEEP_SECONDS:"$SLEEP_SECONDS
# sleep $SLEEP_SECONDS

for ((i = 0; i < $SLEEP_SECONDS; i++))
do
    sleep 1
    date
done

exit 0
