#!/bin/bash

# local_ip_index : 1..n
# local_ip_index=4

if [ -z "$local_ip_index" ]; then
  local_ip=$(hostname -i)
else
  local_ip=$(hostname -I| cut -d ' ' -f ${local_ip_index})
fi
echo "local-ip:$local_ip"

docker exec -it database psql -Uscalebox -P pager=off  -c "
DELETE
FROM t_host
WHERE cluster<>'local' AND cluster NOT IN(
  SELECT cluster
  FROM t_host
  WHERE ip_addr='$local_ip'
);

"
