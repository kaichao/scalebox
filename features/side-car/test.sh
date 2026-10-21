#!/bin/bash
# side-car：ACTION 脚本按阶段执行；定制 ACTION_RUN 的 stdout 可在执行记录中看到。
# SETUP/CHECK/TEARDOWN 的输出在 t_slot_event，不在此断言。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

export ACTION_RUN=/app/bin/run-custom.sh
export ACTION_SETUP=/app/bin/setup-custom.sh
export ACTION_TEARDOWN=/app/bin/teardown-custom.sh
export ACTION_CHECK=/app/bin/check-custom.sh

tk_init
tk_run_app scalebox.env
if ! echo "body-1" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module side_car_test \
        > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递任务失败（见 $TK_LOG_DIR/add.log）"
fi
tk_wait_tasks 1 60

stdout=$(tk_stdout_of_module side_car_test)
tk_assert_contains "$stdout" "run-custom.sh" "ACTION_RUN 使用定制脚本"
tk_assert_contains "$stdout" "task-body:body-1" "任务 body 传入"

echo "PASS [$TK_NAME] 定制 ACTION_RUN 生效"
