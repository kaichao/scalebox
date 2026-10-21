#!/bin/bash
# slot-expr：((@seq)) / ((@v:name)) 在 slot 启动时求值。
# ((@p:name)) 在 slot 无参数时不会被替换（平台缺陷，见 README「坑」），
# 因此用 feature-test.yaml（去掉 @p 项）执行本测试。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

export SCALEBOX_FEATURE_FILE=feature-test.yaml

tk_init
"$TK_CLI" global set my_var VAR_VALUE >/dev/null 2>&1 || tk_fail "global set 失败"
tk_run_app

if ! echo "body" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module mod-slot-expr \
        > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递任务失败（见 $TK_LOG_DIR/add.log）"
fi
tk_wait_tasks 1 60

stdout=$(tk_stdout_of_module mod-slot-expr)
tk_assert_contains "$stdout" "SEQ=0" "@seq 求值"
tk_assert_contains "$stdout" "MY_VAR=VAR_VALUE" "@v: 全局变量求值"

echo "PASS [$TK_NAME] slot 启动表达式求值（@seq / @v:）"
