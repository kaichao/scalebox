#!/bin/bash
# task-identity：相同任务身份（mod_id/from_mod/to_host/vtask/body）重复投递被拒绝。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app

# feature.yaml 的 start_task:0 已投递第一个任务。
tk_wait_tasks 1
before=$(tk_task_count)

# 第一次投递新 body：成功。
if ! echo "body-1" | "$TK_CLI" task add --app-id "$TK_APP_ID" \
        --sink-module identity > "$TK_LOG_DIR/add-1.log" 2>&1; then
    tk_fail "第一次投递失败（见 $TK_LOG_DIR/add-1.log）"
fi
tk_wait_tasks 2

# 第二次投递相同身份：应报 duplicate key（112），且不新增任务。
echo "body-1" | "$TK_CLI" task add --app-id "$TK_APP_ID" \
    --sink-module identity > "$TK_LOG_DIR/add-2.log" 2>&1
rc=$?
after=$(tk_task_count)
tk_assert_eq "$((before + 1))" "$after" "重复投递不应新增任务"

if [ "$rc" -eq 0 ]; then
    if ! grep -qE 'duplicate key|112' "$TK_LOG_DIR/add-2.log"; then
        tk_fail "重复投递既未返回非 0，也没有 duplicate key 提示（见 $TK_LOG_DIR/add-2.log）"
    fi
fi

echo "PASS [$TK_NAME] uniq_task 重复投递被拒绝（rc=$rc）"
