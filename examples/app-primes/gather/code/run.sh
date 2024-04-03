#!/bin/bash

# second column is num
num=$(echo $1 | awk -F "," '{print $2}')

# define integer variables
declare -i count sum

line=$(head -1 /tmp/result.txt)
count=$(echo $line | awk -F " " '{print $1}')
sum=$(echo $line | awk -F " " '{print $2}')
((count++))
sum=$(($sum + $num)) 
echo -n $count $sum > /tmp/result.txt

# save the result to t_app, and set the status of the application to FINISHED
if [[ "$count" = "${NUM_GROUPS}" ]]; then
    scalebox app set-finished --job-id=${JOB_ID} "Result is "${sum}
fi
