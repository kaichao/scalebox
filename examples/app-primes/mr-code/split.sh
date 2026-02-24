#!/bin/bash

if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <total_size> <num_groups>"
    exit 1
fi

total_size=$1
num_groups=$2

# 自动计算每组大小（向上取整）
group_size=$(( (total_size + num_groups - 1) / num_groups ))

echo "total_size=$total_size"
echo "num_groups=$num_groups"
echo "group_size=$group_size"

: "${WORK_DIR:=.}"

for ((g=0; g<num_groups; g++)); do
    start=$((g * group_size + 1))
    end=$(( (g+1) * group_size ))
    (( end > total_size )) && end=$total_size

    printf "my-module,%09d_%09d,{\"len\":\"%s\"}\n" "$start" "$end" "$group_size">> "${WORK_DIR}/sink-tasks.txt"

    (( end == total_size )) && break
done


# declare -i num_groups group_size

# num_groups=${NUM_GROUPS:-10}
# group_size=${GROUP_SIZE:-10000}

# env

# for ((i=num_groups*group_size-group_size+1; i>0; i=i-group_size))
# do
# 	printf "my-module,%09d\n" ${i} >> ${WORK_DIR}/sink-tasks.txt
# done

