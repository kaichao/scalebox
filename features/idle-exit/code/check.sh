#!/bin/bash
# 自定义 check 让 agent 每个 poll 都提交 SLOT_CHECK，从而快速刷新
# t_slot_runtime.idle_polls_remaining；空检查不会提交该事件。
exit 0
