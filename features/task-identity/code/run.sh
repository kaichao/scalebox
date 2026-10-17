#!/bin/bash
# 成功退出；用于重复投递同一 body 观察 uniq_task 约束（退出码 112）
echo "identity task: $1"
exit 0
