#!/bin/bash
# group-node-drain：DRAINING host 上 skip 模块停领并退出，其他模块继续消费。
# 需要 inline 集群的 n0-0 / n0-1（见 scalebox.env 的 CLUSTER=inline）。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init

# 目标主机短名。默认要求 inline 多节点集群；设置 SCALEBOX_DRAIN_HOST 可改用
# 单节点变体（如 local 集群的 h0，配合 SCALEBOX_ENV_FILE=scalebox-local.env）。
drain_host_short=${SCALEBOX_DRAIN_HOST:-n0-0}

if [ -z "${SCALEBOX_DRAIN_HOST:-}" ]; then
    grpc=$(tk_db "select coalesce(parameters->>'grpc_server', grpc_server) from t_cluster where name='inline';")
    [ -n "$grpc" ] || tk_skip "没有 inline 集群定义"
    if ! timeout 3 bash -c "echo > /dev/tcp/${grpc%%:*}/${grpc##*:}" 2>/dev/null; then
        tk_skip "inline 集群控制面 $grpc 不可达（需要多节点部署）"
    fi
fi

tk_run_app "${SCALEBOX_ENV_FILE:-scalebox.env}"

cluster=$(tk_db "select cluster from t_module where app=$TK_APP_ID order by id limit 1;")
host=$(tk_db "select hostname from t_host where split_part(hostname,'.',1)='$drain_host_short' and cluster='$cluster' and status='ON' order by id desc limit 1;")
[ -n "$host" ] || tk_skip "集群 $cluster 没有可用的 $drain_host_short 节点"

restore_host() {
    "$TK_CLI" host undrain "$host" >/dev/null 2>&1 || true
}

producer_mod=$(tk_module_id producer)
consumer_mod=$(tk_module_id consumer)
[ -n "$producer_mod" ] || tk_fail "找不到 producer 模块"
[ -n "$consumer_mod" ] || tk_fail "找不到 consumer 模块"

prod_slot=$(tk_db "select s.id from t_slot s where s.mod_id=$producer_mod and s.host='$host' order by s.id limit 1;")
cons_slot=$(tk_db "select s.id from t_slot s where s.mod_id=$consumer_mod and s.host='$host' order by s.id limit 1;")
[ -n "$prod_slot" ] || tk_fail "producer 在 $host 上没有 slot"
[ -n "$cons_slot" ] || tk_fail "consumer 在 $host 上没有 slot"

# 预热：确认两个模块都能在组节点执行
echo "drain-producer-1" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module producer \
    > "$TK_LOG_DIR/p1.log" 2>&1 || tk_fail "投递 producer 任务失败"
echo "drain-consumer-1" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module consumer \
    > "$TK_LOG_DIR/c1.log" 2>&1 || tk_fail "投递 consumer 任务失败"
tk_wait_tasks 2 120

# host 进入 DRAINING
"$TK_CLI" host drain "$host" > "$TK_LOG_DIR/drain.log" 2>&1 || \
    tk_fail "host drain 失败（见 $TK_LOG_DIR/drain.log）"

# 生产者 agent 应在 poll 周期内退出，slot 回 READY
start=$SECONDS
while :; do
    st=$(tk_db "select status from t_slot where id=$prod_slot;")
    [ "$st" = "READY" ] && break
    if [ $((SECONDS - start)) -ge 60 ]; then
        restore_host
        tk_fail "producer slot 未在 60 秒内退出（status=$st）"
    fi
    sleep "$TK_POLL"
done

# 新任务：producer 保持 READY 不被领取；consumer 继续执行
echo "drain-producer-2" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module producer \
    > "$TK_LOG_DIR/p2.log" 2>&1 || tk_fail "投递 producer 任务 2 失败"
echo "drain-consumer-2" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module consumer \
    > "$TK_LOG_DIR/c2.log" 2>&1 || tk_fail "投递 consumer 任务 2 失败"

start=$SECONDS
while :; do
    cs=$(tk_db "select coalesce(status_code,-1) from t_task where mod_id=$consumer_mod and body='drain-consumer-2' order by id desc limit 1;")
    [ "$cs" = "0" ] && break
    if [ $((SECONDS - start)) -ge 60 ]; then
        restore_host
        tk_fail "consumer 任务未在 60 秒内完成（status=$cs）"
    fi
    sleep "$TK_POLL"
done

# 给 producer 留一个完整的 poll 窗口，确认它确实没领
sleep 15
ps=$(tk_db "select coalesce(status_code,-1) from t_task where mod_id=$producer_mod and body='drain-producer-2' order by id desc limit 1;")
restore_host
tk_assert_eq "-1" "$ps" "DRAINING host 上 producer 不应领取新 task"

echo "PASS [$TK_NAME] DRAINING 停生产者、消费者继续"
