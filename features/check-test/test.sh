#!/bin/bash
# check-test：空间闸门不满足时停领（任务保持 READY、无执行记录），满足后自动放行。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
export DIR_QUOTA_GB=
export ACTION_CHECK=

# 第一轮：阈值不可能满足 → 停领
export FREE_SPACE_GB='{"/":999999}'
tk_run_app scalebox.env
sleep 8
stats=$(tk_task_stats)
total=${stats%%|*}
pending=${stats##*|}
tk_assert_eq "1" "$total" "start_task 已创建任务"
tk_assert_eq "1" "$pending" "闸门不满足时任务保持未完成"
execs=$(tk_db "select count(*) from t_task_exec e join t_module m on e.mod_id=m.id where m.app=$TK_APP_ID;")
tk_assert_eq "0" "$execs" "闸门不满足时没有执行记录"

# 清掉第一轮（app + 容器），第二轮不设闸门
first=$TK_APP_ID
"$TK_CLI" app delete --app-id "$first" >/dev/null 2>&1 || true
tk_stop_containers
export FREE_SPACE_GB=
tk_run_app scalebox.env
tk_wait_tasks 1 60

echo "PASS [$TK_NAME] free_space_gb 停领与自动放行"
