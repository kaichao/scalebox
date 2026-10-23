#!/bin/bash
set -euo pipefail

# agent 调用形式：run.sh <task body...> '<headers-json>'
# 最后一个参数是 headers，不打印。
if [ "$#" -gt 1 ]; then
  body="${*:1:$#-1}"
else
  body="no task body"
fi
printf 'hello from ai-agent-hello: %s\n' "$body"
