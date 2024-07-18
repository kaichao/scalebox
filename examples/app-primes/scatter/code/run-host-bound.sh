#!/bin/bash

env

declare -i num_groups group_size
num_groups=${NUM_GROUPS} 
group_size=${GROUP_SIZE}
num_hosts=4

# 生成host-bound消息
for ((i=num_groups*group_size-group_size+1; i>0; i=i-group_size))
do
	n=$(((i -1) / group_size % num_hosts))
	echo "i=$i, n=$n"
	printf "%09d\n" ${i} >> "${WORK_DIR}/task-body-${n}.txt"
done

ls -l ${WORK_DIR}/*

# 导入到数据库
for ((n=0; n<num_hosts; n=n+1))
do
	scalebox task add --header to_host=n-0${n}.inline-cluster --task-file "${WORK_DIR}/task-body-${n}.txt"
done
