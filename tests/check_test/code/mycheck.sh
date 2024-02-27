#!/bin/bash

# code in [0..2]
code=$(($RANDOM%3))

echo "check_code:$code" >> /work/custom-out.txt

exit $code
