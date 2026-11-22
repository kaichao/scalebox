#!/bin/bash
# idle-exit feature: the task body is only used to trigger a non-empty batch,
# so the agent resets its idle budget. Keep it trivial and always successful.
echo "idle-exit task: ${1:-}"
exit 0
