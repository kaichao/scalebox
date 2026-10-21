#!/bin/bash
# pod：HOST-BOUND 且 pod_id 相同的模块之间，下游任务自动留在同一节点。
# 需要 inline 集群的 n0-0 / n0-1（见 scalebox.env 的 CLUSTER=inline）。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init

# 前置：inline 多节点集群的控制面可达（单机测试环境只跑 local 集群）
grpc=$(tk_db "select coalesce(parameters->>'grpc_server', grpc_server) from t_cluster where name='inline';")
[ -n "$grpc" ] || tk_skip "没有 inline 集群定义"
if ! timeout 3 bash -c "echo > /dev/tcp/${grpc%%:*}/${grpc##*:}" 2>/dev/null; then
    tk_skip "inline 集群控制面 $grpc 不可达（需要多节点部署）"
fi

tk_run_app scalebox.env
if ! echo "4" | "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module router \
        > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递入口任务失败（见 $TK_LOG_DIR/add.log）"
fi

# 1 入口 + 4 个 mod-1 + 4 个 mod-2
tk_wait_tasks 9 180

mismatch=$(tk_db "
select count(*) from (
  select t1.body, t1.headers->>'to_host' as h1, t2.headers->>'to_host' as h2
  from t_task t1
  join t_module m1 on t1.mod_id=m1.id and m1.app=$TK_APP_ID and m1.name='mod-1'
  join t_task t2 on t2.body=t1.body
  join t_module m2 on t2.mod_id=m2.id and m2.app=$TK_APP_ID and m2.name='mod-2'
) s
where coalesce(h1,'')='' or coalesce(h2,'')=''
   or split_part(h1,'.',1) <> split_part(h2,'.',1);")
tk_assert_eq "0" "$mismatch" "mod-1 → mod-2 的 to_host 必须一致"

echo "PASS [$TK_NAME] pod_id 路由保持同节点"
