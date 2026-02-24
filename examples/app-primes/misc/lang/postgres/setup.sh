#!/bin/bash

docker run  -d ${SLAVE_IMAGE}|tail -1 > /tmp/container_id
code=$?
# wait for container starting  
sleep 5s

exit $code
