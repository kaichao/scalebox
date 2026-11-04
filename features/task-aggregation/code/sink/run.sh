#!/bin/bash
# 下游模块：只记录收到的 body。
#
# 验收点：本模块任务的 t_task.from_mod 应该是 "agg"（聚合模块名），不是主路由——
# 这就是"直连"的库内证据。
echo "sink received body=${1:-}"
exit 0
