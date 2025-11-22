#!/bin/bash

declare -i num_groups group_size

num_groups=${NUM_GROUPS:-10}
group_size=${GROUP_SIZE:-10000}

env

for ((i=num_groups*group_size-group_size+1; i>0; i=i-group_size))
do
	printf "%09d\n" ${i} >> ${WORK_DIR}/sink-tasks.txt
done

scalebox task add --sink-job=my-job
