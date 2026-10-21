#!/bin/bash
# node-progress-gap：HOST-BOUND 多节点，node_progress 信号量按节点建立。
# 领取时的 gap 限制是动态行为（依赖任务持续供给），本测试只做基线验证：
# 任务全绿 + 每个节点一个 node_progress 信号量。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init

# 前置：inline 多节点集群的控制面可达
grpc=$(tk_db "select coalesce(parameters->>'grpc_server', grpc_server) from t_cluster where name='inline';")
[ -n "$grpc" ] || tk_skip "没有 inline 集群定义"
if ! timeout 3 bash -c "echo > /dev/tcp/${grpc%%:*}/${grpc##*:}" 2>/dev/null; then
    tk_skip "inline 集群控制面 $grpc 不可达（需要多节点部署）"
fi

tk_run_app scalebox.env
if ! "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module node-progress-gap \
        < tasks.txt > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递任务失败（见 $TK_LOG_DIR/add.log）"
fi
tk_wait_tasks 10 180

sema=$(tk_db "select count(*) from t_semaphore where app=$TK_APP_ID and name ~ '^node_progress:';")
tk_assert_eq "4" "$sema" "4 个节点各一个 node_progress 信号量"

echo "PASS [$TK_NAME] node_progress 信号量按节点建立"
