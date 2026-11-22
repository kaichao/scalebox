#!/bin/bash
# idle-exit：验证 agent 空闲预算的四条行为：
#   1. 非空 task 批次后重置 idle_polls_remaining；
#   2. 预算耗尽时 ON → READY，SLOT_COMPLETE 记录 slot exit=idle timeout；
#   3. always_running 豁免空闲退出；
#   4. max_idle_polls=0 回落 100，不立即退出。
set -uo pipefail
cd "$(dirname "$0")"
source ../lib/testkit.sh

tk_init
tk_run_app

# 等待模块的 slot 行出现，返回 slot id。
wait_slot_exists() {
    local mod=$1 timeout=${2:-60} start=$SECONDS slot=""
    while :; do
        slot=$(tk_db "select id from t_slot where mod_id=$mod order by id limit 1;" 2>/dev/null || true)
        if [ -n "$slot" ]; then
            printf '%s\n' "$slot"
            return 0
        fi
        if [ $((SECONDS - start)) -ge "$timeout" ]; then
            tk_fail "module $mod 的 slot 未在 ${timeout}s 内创建"
        fi
        sleep "$TK_POLL"
    done
}

# 等待指定 slot 到达目标生命周期状态。
wait_slot_status() {
    local slot=$1 want=$2 timeout=${3:-60} start=$SECONDS st=""
    while :; do
        st=$(tk_db "select status from t_slot where id=$slot;" 2>/dev/null || true)
        [ "$st" = "$want" ] && return 0
        if [ $((SECONDS - start)) -ge "$timeout" ]; then
            tk_fail "slot $slot 未在 ${timeout}s 内变为 $want（当前 $st）"
        fi
        sleep "$TK_POLL"
    done
}

# idle-reset 有 start_task；等任务完成，证明发生了非空 task 批次。
tk_wait_tasks 1 60

reset_mod=$(tk_module_id idle-reset)
[ -n "$reset_mod" ] || tk_fail "找不到 idle-reset 模块"
reset_slot=$(wait_slot_exists "$reset_mod")
wait_slot_status "$reset_slot" ON 60

# 自定义 ACTION_CHECK 每个 poll 都上报 idle_polls_remaining；任务后应回到 20 附近。
start=$SECONDS
remaining=-1
while :; do
    remaining=$(tk_db "select coalesce(idle_polls_remaining,-1) from t_slot_runtime where slot_id=$reset_slot;" 2>/dev/null || echo -1)
    [ "${remaining:--1}" -ge 15 ] && break
    if [ $((SECONDS - start)) -ge 10 ]; then
        tk_fail "task 完成后 idle_polls_remaining 未重置（当前 $remaining，期望 >=15）"
    fi
    sleep 1
done
echo "[$TK_NAME] reset: idle_polls_remaining=$remaining"

exit_mod=$(tk_module_id idle-exit)
[ -n "$exit_mod" ] || tk_fail "找不到 idle-exit 模块"
exit_slot=$(wait_slot_exists "$exit_mod")
wait_slot_status "$exit_slot" READY 60

exit_reason=$(tk_db "select coalesce(details->>'slot exit','') from t_slot_log where slot=$exit_slot and event_type='SLOT_COMPLETE' order by id desc limit 1;")
tk_assert_eq "idle timeout" "$exit_reason" "空闲退出原因"

exit_max=$(tk_db "select coalesce(details->>'slot exit:max_idle_polls','') from t_slot_log where slot=$exit_slot and event_type='SLOT_COMPLETE' order by id desc limit 1;")
exit_interval=$(tk_db "select coalesce(details->>'slot exit:poll_interval_seconds','') from t_slot_log where slot=$exit_slot and event_type='SLOT_COMPLETE' order by id desc limit 1;")
tk_assert_eq "5" "$exit_max" "退出明细中的 max_idle_polls"
tk_assert_eq "1" "$exit_interval" "退出明细中的 poll_interval_seconds"

always_mod=$(tk_module_id idle-always)
[ -n "$always_mod" ] || tk_fail "找不到 idle-always 模块"
always_slot=$(wait_slot_exists "$always_mod")
wait_slot_status "$always_slot" ON 60

clamp_mod=$(tk_module_id idle-clamp)
[ -n "$clamp_mod" ] || tk_fail "找不到 idle-clamp 模块"
clamp_slot=$(wait_slot_exists "$clamp_mod")
wait_slot_status "$clamp_slot" ON 60

# idle-exit 已退出；再等一个完整的 always_running 预算窗口。
sleep 3
always_status=$(tk_db "select status from t_slot where id=$always_slot;")
clamp_status=$(tk_db "select status from t_slot where id=$clamp_slot;")
tk_assert_eq "ON" "$always_status" "always_running 不应空闲退出"
tk_assert_eq "ON" "$clamp_status" "max_idle_polls=0 应回落 100，不应立即退出"

echo "PASS [$TK_NAME] task 后重置 / 空闲退出明细 / always_running / 0 值回落"
