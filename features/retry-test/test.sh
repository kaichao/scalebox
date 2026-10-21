#!/bin/bash
# retry-test：退出码命中 retry_rules 后，pgcron（每 5 分钟）把任务重置为 READY 并重跑。
# 慢测：最长等待 420 秒（一个 pgcron 周期 + 重跑）。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app

# body=1 → run.sh 退出码 1，命中 '1:3'
if ! echo "1" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module retry_test \
        > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递任务失败（见 $TK_LOG_DIR/add.log）"
fi

task_id=$(tk_db "select id from t_task where mod_id=(select id from t_module where app=$TK_APP_ID and name='retry_test') order by id desc limit 1;")
[ -n "$task_id" ] || tk_fail "未找到投递的任务"

# 第一次执行以失败结束
start=$SECONDS
while [ $((SECONDS - start)) -lt 60 ]; do
    code=$(tk_db "select coalesce(status_code,-1) from t_task where id=$task_id;")
    [ "$code" != "-1" ] && [ "$code" != "-3" ] && break
    sleep 2
done
tk_assert_eq "1" "$code" "第一次执行失败（退出码 1）"

# 等 pgcron 重置并重跑：同一 task 出现第 2 条 exec 记录
start=$SECONDS
n=0
while [ $((SECONDS - start)) -lt 420 ]; do
    n=$(tk_db "select count(*) from t_task_exec where task=$task_id;")
    [ "$n" -ge 2 ] && break
    sleep 10
done
[ "$n" -ge 2 ] || tk_fail "等待 420 秒仍未观察到重试（exec 记录数 $n）"

rep=$(tk_db "select coalesce(headers->>'repeatable','') from t_task where id=$task_id;")
tk_assert_eq "yes" "$rep" "重试任务带 repeatable=yes"

echo "PASS [$TK_NAME] 退出码 1 触发重试（exec 记录数 $n）"
