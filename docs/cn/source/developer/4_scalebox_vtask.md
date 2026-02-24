# 4. vtask管理

vtask跨越输入加载、计算、结果回写的整个过程，基于vtask的容错简单、直接。

vtask标识一种跨越多个模块的虚拟task，是基于计数的流控机制，分为全局计数（global_vtask）、分组计数（group_vtask）、节点计数（host_vtask）等三类。

- 前置任务队列
- vtask头模块
- vtask算法模块（可为多个）
- vtask尾模块

