#!/bin/bash

# vtask 资源额度已由 controld 在 vtask-tail 任务完成时（doVTaskFinished）
# 自动加回（读 tail task 的 _vtask_size_sema），脚本侧无需再调用 vtask unbind

exit 0
