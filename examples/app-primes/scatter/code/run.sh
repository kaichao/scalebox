#!/bin/bash

declare -i num_groups group_size
num_groups=${NUM_GROUPS} 
group_size=${GROUP_SIZE}

# for ((i=num_groups*group_size-group_size+1; i>0; i=i-group_size))
# do
# 	send-message $(printf "%09d" ${i})
# done

env

for ((i=num_groups*group_size-group_size+1; i>0; i=i-group_size))
do
	printf "%09d\n" ${i} >> ${WORK_DIR}/task-body.txt
done

scalebox task add
