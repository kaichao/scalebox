#!/bin/bash
# timeout：task_max_seconds 到期后任务以 124 结束。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app scalebox.env
tk_wait_tasks 1 60

task_status=$(tk_db "select coalesce(status_code,-1) from t_task where mod_id in (select id from t_module where app=$TK_APP_ID) order by id desc limit 1;")
exec_status=$(tk_status_of_module timeout)
stdout=$(tk_stdout_of_module timeout)

tk_assert_eq 124 "$task_status" "t_task 超时退出码"
tk_assert_eq 124 "$exec_status" "t_task_exec 超时退出码"
tk_assert_contains "$stdout" "SLEEP_SECONDS:20" "任务已启动后被中断"

echo "PASS [$TK_NAME] task_max_seconds 超时（124）"
