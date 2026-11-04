#!/bin/bash
# 聚合模块：一个 task 覆盖一组"短任务"成员。
#
#   $1 = task body（组标识）
#   $2 = headers JSON（成员清单由上游下发，例如 {"members":"m1,m2,m3"}）
#
# 组内循环完成后，把下游任务直接写进 sink-tasks.txt——第一列是 sink_modules 的模块名
# （这里用 agent 注入的 SINK_MODULE），**不经主路由**。
#
# 失败语义：任一成员失败即整组失败，不写 sink-tasks.txt（下游不启动）。
set -uo pipefail

body=${1:-}
headers=${2:-}
steps=${AGG_STEPS:-8}

members=$(printf '%s' "$headers" | sed -n 's/.*"members":"\([^"]*\)".*/\1/p')
[ -n "$members" ] || members="m1,m2,m3"

echo "agg body=${body} members=${members} per_member=${steps}"

rc=0
n=0
IFS=',' read -ra list <<< "$members"
for m in "${list[@]}"; do
    for i in $(seq 1 "$steps"); do
        # 真实应用在这里调用一次"短任务"要做的工作；UNIT 行供测试计数用
        echo "  UNIT ${m} ${i}"
        n=$((n + 1))
    done
done

echo "agg body=${body} units_done=${n}"

[ "$rc" -eq 0 ] || exit "$rc"

printf '%s,%s\n' "${SINK_MODULE:-sink}" "$body" > "${WORK_DIR}/sink-tasks.txt"
exit 0
