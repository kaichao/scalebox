#!/bin/bash
# task-batch-size：一次领取多条任务。批量效应体现在 agent 的取任务调用次数，
# 不在 t_task_exec 条数（每个 task 仍一条 exec 记录）；本测试从 slot 容器日志
# 的 "GetTaskList(), len of tasks=N" 验证单次领取 > 1。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app scalebox.env

mod_id=$(tk_module_id task-batch-size)
[ -n "$mod_id" ] || tk_fail "未找到模块 task-batch-size"

# 先投递 10 条：平台只启动有待执行任务的 slot（GetSlotList 要求 t_task 有 -1/-2），
# 且这样第一次领取时队列里已有全部任务，才能观察到批量领取。
if ! for i in $(seq 0 9); do printf '%03d\n' "$i"; done | \
        "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module task-batch-size \
        > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递任务失败（见 $TK_LOG_DIR/add.log）"
fi

start=$SECONDS
while [ $((SECONDS - start)) -lt 60 ]; do
    on=$(tk_db "select count(*) from t_slot where mod_id=$mod_id and status='ON';")
    [ "$on" -ge 1 ] && break
    sleep 2
done
[ "${on:-0}" -ge 1 ] || tk_fail "slot 60 秒内未进入 ON"

tk_wait_tasks 10 120

# 从 slot 容器日志取单次领取的最大条数
cid=""
for c in $(docker ps -q 2>/dev/null); do
    app=$(docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$c" 2>/dev/null | \
        sed -n 's/^APP_ID=//p' | head -1)
    if [ "$app" = "$TK_APP_ID" ]; then cid=$c; break; fi
done
[ -n "$cid" ] || tk_fail "未找到本 app 的 slot 容器"

maxbatch=$(docker logs "$cid" 2>&1 | sed -n 's/.*len of tasks=\([0-9][0-9]*\).*/\1/p' | sort -n | tail -1)
[ -n "$maxbatch" ] || tk_fail "容器日志中没有 len of tasks= 记录（LOG_LEVEL=trace 是否生效？）"
if [ "$maxbatch" -lt 2 ]; then
    tk_fail "未观察到批量领取：单次最大 $maxbatch 条（期望 >=2）"
fi

echo "PASS [$TK_NAME] 批量领取：单次最多 $maxbatch 条"
