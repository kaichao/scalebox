#!/bin/bash
# headers-propagation：`_` 前缀 header 沿 sink 自动传播到下游任务。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app

# producer 尚无 start_task，手动投递入口任务并带上 _demo。
if ! echo "hello" | "$TK_CLI" task add --app-id "$TK_APP_ID" \
        --sink-module producer --header _demo=hello \
        > "$TK_LOG_DIR/task-add.log" 2>&1; then
    tk_fail "投递 producer 任务失败（见 $TK_LOG_DIR/task-add.log）"
fi

# producer + consumer 各一个任务。
tk_wait_tasks 2

consumer_stdout=$(tk_stdout_of_module consumer)
tk_assert_contains "$consumer_stdout" "_demo=hello" "header 自动传播"

echo "PASS [$TK_NAME] _ 前缀 header 传播"
