# 任务食谱（recipes）

面向 agent 的"怎么做"清单。每篇固定四段：**前置条件 → 步骤（可复制命令）→ 判据 → 失败下一步**。
规则与红线见 [`../10-invariants.md`](../10-invariants.md)；症状排查看 [`../20-diagnosis.md`](../20-diagnosis.md)。

| recipe | 用于 |
|---|---|
| [`new-app`](./new-app.md) | 从零建一个能跑的最小 app（单模块） |
| [`add-module`](./add-module.md) | 在已有 app 里加模块并接线 |
| [`flow-control`](./flow-control.md) | 加空间/数量闸门，避免死锁 |
| [`vtask`](./vtask.md) | 建 wait-queue → head → core → tail 管道 |
| [`multi-node`](./multi-node.md) | 多节点落点与共享盘 |
| [`debug-stuck-chain`](./debug-stuck-chain.md) | 链卡住时的处置流程 |
| [`acceptance`](./acceptance.md) | 一轮运行结束后的验收清单 |

约定：

- `<...>` 为占位符：`<app_id>`、`<module_id>`、`<task_id>`、`<N>`。
- 命令默认在控制面机器执行；需要容器清理时以控制面为本地集群为例。
- 每条步骤末尾的"判据"不满足时，先按失败下一步处理，不要盲目重跑。
