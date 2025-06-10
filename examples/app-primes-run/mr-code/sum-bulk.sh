#!/bin/bash

set -e
source "/usr/local/bin/functions.sh"

# 2025-06-08T21:24:14.369583794+08:00
ds=$(date --iso-8601=ns | sed 's/,/./')
echo "$ds,before-sum" >> ${WORK_DIR}/timestamps.txt

sum=0
pattern='"part_primes":"([^"]+)"'
# 767896,000090001,{"from_ip":"10.0.6.100","from_job":"my-job","part_primes":"879","to_slot":"255"}
while read -r line; do
    # task_id=$(echo "$line" | cut -d',' -f1)
    # task_body=$(echo "$line" | cut -d',' -f2)
    headers=$(echo "$line" | cut -d',' -f3-)
    # very slow
    # n=$(get_header "$headers" "part_primes")
    if [[ $headers =~ $pattern ]]; then
        n="${BASH_REMATCH[1]}"
    else
        # no part_primes in json 
        n="0"
    fi

    sum=$((sum + n))
done < ${WORK_DIR}/input-messages.txt

ds=$(date --iso-8601=ns | sed 's/,/./')
echo "$ds,after-sum" >> ${WORK_DIR}/timestamps.txt

export SEMAPHORE_AUTO_CREATE=yes
val=$(scalebox semaphore increment-n app-primes:sum_value ${sum})
code=$?

date --iso-8601=ns | sed 's/,/./' >> ${WORK_DIR}/timestamps.txt

echo "part_sum=${val}"

exit $code
