#!/bin/bash
# task-aggregation：N 个短任务在算法侧聚合成 1 个粗粒度任务，完成后直连下游，
# 不经过主路由。
#
# 断言四件事：
#   1. 聚合模块任务数 = 组数（不是 组数 × 成员数）
#   2. 下游任务数 = 组数
#   3. 下游任务的 from_mod = "agg"（直连证据，不是经 router 中转）
#   4. 单个聚合任务内部确实跑了 N 步（聚合真的发生，不是空转）
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app

agg_id=$(tk_module_id agg)
sink_id=$(tk_module_id sink)
[ -n "$agg_id" ] || tk_fail "未找到模块 agg"
[ -n "$sink_id" ] || tk_fail "未找到模块 sink"

GROUPS=3
if ! for i in $(seq 1 "$GROUPS"); do echo "batch-$i"; done | \
        "$TK_CLI" task add --app-id "$TK_APP_ID" --sink-module agg \
        > "$TK_LOG_DIR/add.log" 2>&1; then
    tk_fail "投递任务失败（见 $TK_LOG_DIR/add.log）"
fi

# 3 组 × (1 个 agg + 1 个 sink) = 6 条任务
tk_wait_tasks $((GROUPS * 2)) 180

n_agg=$(tk_db "select count(*) from t_task where mod_id=$agg_id;")
n_sink=$(tk_db "select count(*) from t_task where mod_id=$sink_id;")
tk_assert_eq "$GROUPS" "$n_agg" "聚合模块任务数（应为组数，不是组数×成员数）"
tk_assert_eq "$GROUPS" "$n_sink" "下游任务数"

# 直连证据：下游任务的来源模块是聚合模块，不是主路由
froms=$(tk_db "select string_agg(distinct coalesce(from_mod,'<null>'), ',') from t_task where mod_id=$sink_id;")
tk_assert_eq "agg" "$froms" "下游任务的来源模块（直连时应为聚合模块名）"

# 聚合真的发生：单条聚合任务内部跑了成员数 × 步数（每个单元输出一行 UNIT）
units=$(tk_db "select coalesce(sum(length(stdout) - length(replace(stdout,'UNIT',''))),0) / 4 from t_task_exec where mod_id=$agg_id;")
expected=$((GROUPS * 3 * 8))   # 组数 × 默认成员数(3) × AGG_STEPS(8)
tk_assert_eq "$expected" "$units" "聚合任务内部执行的单元步数总数"

echo "PASS [$TK_NAME] $GROUPS 组 → $n_agg 个聚合任务 → $n_sink 个下游任务，来源均为 agg，共 $units 步"
