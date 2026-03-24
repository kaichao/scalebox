#!/bin/bash

if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <total_size> <app_id>"
    exit 1
fi

num_groups=${NUM_GROUPS:-10}

total_size=$1
app_id=$2

group_size=$(( (total_size + num_groups - 1) / num_groups ))

# to sub-ranges
for ((i=0; i<num_groups; i++)); do
    start=$((i * group_size + 1))
    end=$(( (i+1) * group_size ))
    (( end > total_size )) && end=$total_size

    n="$((i % 2 + 1))"
    body=$(printf "%09d_%09d" "$start" "$end")
    scalebox task add --app-id="$app_id" --remote-cluster="cluster${n}" "$body"

    (( end == total_size )) && break
done
