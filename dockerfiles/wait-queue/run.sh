#!/bin/bash

# wait-queue 标准模块：vtask 资源分配 + 转发（slot 主逻辑）
# 执行位置：wait-queue slot（拾取来自 main-router 的 task）
# 完成后不写 sink-tasks.txt，task 在本 slot 终结（vtask-head 任务由 add-subtask --direct 直投）
#
# 环境变量：
#   VTASK_MODE       vtask 类型：HOST-BOUND / GROUP-BOUND（缺省 HOST-BOUND）
#   HEAD_MODULE      vtask-head 模块名（缺省 vtask-head）
#   PLAT_MODULE_NAME 当前模块名（agent 注入，用于 gate 信号量）
#   BIND_RETRY_SECONDS bind 失败时的重试上限（**秒**，缺省 120；内部按 0.2 s
#                      粒度轮询，即最多 600 次。2026-11-03 前实现把该值当
#                      循环次数用，实际窗口只有 24 秒，单节点高在制批数时
#                      会误判为"无空额"而丢整批——见文件尾注）
#
# 入参：$1 = task body，$2 = task headers（JSON，agent 注入）
# 转发规则：`_vtask_attr_*` 前缀的 header 为应用属性，原样转发给 vtask-head
#（平台只运输、不解释）；其余 `_vtask_*` 不继承——_vtask_id 由服务端自指，
# _vtask_size_sema 用本次 bind 的新值。
#
# 2026-11-02（vtask 规模化 L1）三处改动：
#   1) bind + add-subtask 成功后**立即**归还串行门。原来这一步在 main-router 的
#      from-vtask-head.sh 里（要等 head 执行完、回流），一个 vtask 要穿过两个轮询
#      边界；挪到这里后只剩一个，节拍 2×poll → 1×poll。与 from-vtask-head.sh 的
#      +1 必须二选一，否则双重归还、信号量只涨不落。
#   2) bind 失败时**等待重试**，而不是直接失败任务。多 slot 并发时 check.sh 的
#      `semagroup max > 0` 挡不住全部（updateVTask 的扣减走 asyncbatch，有 10 ms
#      窗口），会出现「读到位子、bind 时没了」。原来 `|| exit 1` 让这些任务直接
#      status=1，120 个里丢了 34 个（2026-11-02 实测）。
#   3) 重试耗尽或 add-subtask 失败时归还 gate。原实现的 `|| exit 1` 不还，几次之后
#      vtask_size 卡在 0、整条链停摆（app 218 实测）。

source /usr/local/lib/scalebox/functions.sh

body=$1
head_module=${HEAD_MODULE:-vtask-head}
module_name=${PLAT_MODULE_NAME:-wait-queue}
bind_retry=${BIND_RETRY_SECONDS:-120}

# 收集应用属性（_vtask_attr_*），随 add-subtask 转发进 vtask 范围
headers_json=${2:-}
declare -A task_headers
if [ -n "${headers_json}" ] && [ "${headers_json}" != "{}" ]; then
    scalebox::json_parse "${headers_json}" task_headers
fi
attr_headers=()
for k in "${!task_headers[@]}"; do
    if [[ "$k" == _vtask_attr_* ]]; then
        attr_headers+=(--header "$k=${task_headers[$k]}")
    fi
done

# bind 重试：无空额时等待（其他 vtask 还在跑），不丢任务。
# 轮询粒度 0.2 s，循环次数 = 秒数 × 5（2026-11-03 前误把秒数直接当循环次数，
# 实际窗口只有标称值的 1/5）。
bind_with_retry() {
    local i out
    for ((i = 0; i < bind_retry * 5; i++)); do
        if out=$(scalebox vtask bind 2>/dev/null); then
            echo "$out"
            return 0
        fi
        sleep 0.2
    done
    return 1
}

if [ "$VTASK_MODE" = "GROUP-BOUND" ]; then
    # GROUP-BOUND vtask：bind 返回组号（head slot 的 seq）
    group_seq=$(bind_with_retry) || {
        scalebox semaphore increment vtask_size:${module_name}
        exit 1
    }
    sema_name="group_vtask_size:${head_module}:${group_seq}"
    scalebox vtask add-subtask --direct --module=$head_module \
        "${attr_headers[@]}" \
        --header to_slot_index=$group_seq --header _vtask_size_sema=$sema_name $body || {
        scalebox vtask unbind --sema-name="${sema_name}"
        scalebox semaphore increment vtask_size:${module_name}
        exit 1
    }
    scalebox semaphore increment vtask_size:${module_name}
else
    # HOST-BOUND vtask（DEFAULT 不经 wait-queue）
    to_host=$(bind_with_retry) || {
        scalebox semaphore increment vtask_size:${module_name}
        exit 1
    }
    sema_name="host_vtask_size:${head_module}:${to_host}"
    scalebox vtask add-subtask --direct --module=$head_module \
        "${attr_headers[@]}" \
        --header to_host=$to_host --header _vtask_size_sema=$sema_name $body || {
        # 回滚已扣减的资源，并放行下一个 vtask
        scalebox vtask unbind --sema-name="${sema_name}"
        scalebox semaphore increment vtask_size:${module_name}
        exit 1
    }
    scalebox semaphore increment vtask_size:${module_name}
fi

exit 0

# 2026-11-03：BIND_RETRY_SECONDS 的单位修正。原实现是
#     for ((i = 0; i < bind_retry; i++)); do ...; sleep 0.2; done
# 把"秒"当成"循环次数"用，标称 120 秒的实际窗口只有 24 秒。app 328
# （单节点、VTS=12、120 批）实测：批 00000024 于 12:46:26 开始、12:46:53
# 失败，恰好 26 秒。单节点上一个 vtask 从 bind 到完成约 127 秒，24 秒窗口
# 等不到空额，整批 80 个任务丢失（wait-queue status=1，下游 vtask-head
# 只得到 119/120）。现改为 bind_retry * 5 次，标称值即真实秒数。
#
# 应用侧兼容：app.yaml 的 BIND_RETRY_SECONDS 若按旧语义放大过（如 1200
# 是为了凑 240 秒），在新镜像下会变成 1200 秒——功能上无害（有空额立即
# 返回），但故障发现变慢，建议同步回调到 300 左右。
