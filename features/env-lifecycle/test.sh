#!/bin/bash
# env-lifecycle：app 创建时快照环境变量；未定义变量保留字面量；注入 APP_ID。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app scalebox.env
tk_wait_tasks 1

stdout=$(tk_stdout_of_module env-lifecycle)
tk_assert_contains "$stdout" "DEMO_VALUE=hello" "environments 注入"
tk_assert_contains "$stdout" 'LITERAL_DEMO=${LITERAL_DEMO}' "未定义变量保留字面量"
tk_assert_contains "$stdout" "APP_ID=$TK_APP_ID" "APP_ID 注入"

echo "PASS [$TK_NAME] environments 注入 / 未定义变量 / APP_ID"
