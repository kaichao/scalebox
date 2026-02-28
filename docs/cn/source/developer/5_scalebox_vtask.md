# 5. vtask管理与编程

vtask作为节点本地计算编程模型的基本单元，通过参数设计，简化本地计算在应用程序中的流控、容错等的实现；


## 5.1 vtask主要功能
- 任务管理：跨模块的任务（task）组成vtask，可简化任务管理，并支持状态管理，进而支撑节点本地计算模型；
- 计数准入控制：以配合HPC计算资源的动态调度，实现超长任务计算；
- 粗粒度容错；实现细粒度检查点功能，自动容错。

## 5.2 vtask模块结构

- 前置任务队列（wait-queue）：尚未绑定计算资源。全局容错的基本单元。
- vtask头模块（vtask-head）：vtask起始标识模块。为vtask分配计算资源（单节点/资源组），单节点模式使用HOST-BOUND，直接运行在计算节点上；资源组模式使用SLOT-BOUND，多个slot需部署在单节点（通常为头节点）上。
- vtask算法模块（vtask-core）：vtask核心处理模块，针对资源组模式，通常用hostname前缀来标识；可以为多个模块，通过pod实现模块间本地分发。
- vtask尾模块（vtask-tail）：vtask结束标识模块，通常部署在头节点上。任务体（task-body）与vtask-head保持一致。


## 5.3 vtask模式分类

task_dist_mode设置

|            |   全局模式   |  单节点模式  | 资源组模式  |
| ---------- | ----------- | ---------- | ---------- |
| wait-queue |             |            |            |
| vtask-head |             | HOST-BOUND | SLOT-BOUND |
| vtask-core |             | HOST-BOUND | HOST-BOUND |
| vtask-tail |             |            |            |

### 5.3.1 全局模式vtask

### 5.3.2 单节点模式vtask

### 5.3.3 分组模式vtask

vtask跨越输入加载、计算、结果回写的整个过程，基于vtask的容错简单、直接。

vtask标识一种跨越多个模块的虚拟task，是基于计数的流控机制，分为全局计数（global_vtask）、分组计数（group_vtask）、节点计数（host_vtask）等三类。
