#!/bin/bash

docker rm -f $(head -n +1 /tmp/container_id)
code=$?

rm -f /tmp/container_id

exit $code
