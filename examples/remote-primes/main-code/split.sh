#!/bin/bash

if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <total_size> <app_id>"
    exit 1
fi

num_groups=${NUM_GROUPS:-10}

total_size=$1
app_id=$2

# 自动计算每组大小（向上取整）
group_size=$(( (total_size + num_groups - 1) / num_groups ))

echo "total_size=$total_size"
echo "num_groups=$num_groups"
echo "group_size=$group_size"

# 打印num_groups行，但确保起始值不超过total_size
for ((i=0; i<num_groups; i++)); do
    # 计算起始值
    start=$((i * group_size + 1))
    end=$(( (i+1) * group_size ))
    (( end > total_size )) && end=$total_size

    # 格式化起始值为9位数字，并循环输出1或2
    # printf "%09d,%d\n" "$start" "$((i % 2 + 1))"
    
    n="$((i % 2 + 1))"
    body=$(printf "%09d_%09d" "$start" "$end")
    scalebox task add --app-id="$app_id" --remote-cluster="cluster${n}" "$body"

    (( end == total_size )) && break
done
