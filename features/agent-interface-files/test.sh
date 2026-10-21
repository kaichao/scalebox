#!/bin/bash
# agent-interface-files：模块写接口文件，agent 采集落库；sink-tasks.txt 7 种格式。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app

# 入口任务 → mr-module → my-module
if ! echo "start" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module mr-module \
        > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递入口任务失败（见 $TK_LOG_DIR/add.log）"
fi

# 1 入口 + 1 my-module + 7 条 sink（全部回到 mr-module）
tk_wait_tasks 9

extras=$(tk_db "select coalesce(extras::text,'') from t_task_exec e join t_module m on e.mod_id=m.id where m.app=$TK_APP_ID and m.name='my-module' order by e.id desc limit 1;")
tk_assert_contains "$extras" "before-sleep" "timestamps.txt：before-sleep"
tk_assert_contains "$extras" "after-sleep" "timestamps.txt：after-sleep"
tk_assert_contains "$extras" "myattr1" "extra-attrs.yaml 落库"
tk_assert_contains "$extras" "iobytes" "input/output 计量落库"

auxout=$(tk_db "select coalesce(auxout,'') from t_task_exec e join t_module m on e.mod_id=m.id where m.app=$TK_APP_ID and m.name='my-module' order by e.id desc limit 1;")
tk_assert_contains "$auxout" "user-defined output" "auxout.txt 落库"

bodies=$(tk_db "select string_agg(body, '|' order by id) from t_task where mod_id=(select id from t_module where app=$TK_APP_ID and name='mr-module');")
for b in body1 body2 body3 body4 body5 body6 body7; do
    tk_assert_contains "$bodies" "$b" "sink 格式：$b"
done

echo "PASS [$TK_NAME] 接口文件落库 + sink 7 种格式"
