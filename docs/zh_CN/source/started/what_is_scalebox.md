# 1. Scalebox是什么？

Scalebox是一种云原生的流式并行计算引擎，用户算法以容器化封装，通过串行的脚本编程，将容器化算法构建为在分布式、异构、多集群环境下、以流式方式运行的并行计算程序。Scalebox应用支持广域网跨集群、任务级容错、本机数据读写等特性，并支持容器间流水线并行、容器级数据级并行、容器内代码并行等多级并行化。

与已有大数据处理、并行计算等框架相比，其技术特点特别适用于数据分布、算力分布、复杂算法封装、极高强度I/O读写等应用场景。简化超大规模并行数据处理的实现。

主要特性包括：
- 云原生设计
- 非侵入式设计：
- 跨集群计算
- 本地计算
- 任务透视及优化支持

## 1.1 Scalebox能做什么？

Scalebox非常适用于超大大规模数据处理，尤其适用于中间数据规模超大的场景。

- 天文计算：天文学家使用大型望远镜和探测器来收集天空中的观测数据，其数据量巨大，包含了大量天体观测数据。天文数据处理通常需要巨量计算资源和复杂的数据分析技术。
- 基因组学和生物信息学：包括基因组测序、基因组组装、基因组比对、蛋白质结构预测等领域的应用。这些应用需要处理大量的生物数据，并利用计算来分析和解释这些数据。
- 高能物理：包括粒子对撞机实验、核物理实验等领域。这些实验产生了大量的高能粒子碰撞数据，需要大量计算资源和复杂的数据分析算法来研究和解释这些数据。
- 大规模数据传输：利用横向扩展支持数据传输，通过数据压缩提升网络带宽利用，通过数据校验提升网络传输的可靠性，利用流水线并行提升运行效率。

Scalebox还可应用于复杂算法处理的应用场景，通过容器化封装，将用户算法、大平台集成分解为相对独立部分。


## 1.2 Scalebox不适合做什么？

Scalebox基于server端做计算模块间通信，针对小数据块的高频度通信，其效率相对较低。因而不适用于通信密集型并行计算。

包括以下类型：
- 科学仿真和建模：包括天气预报、气候建模、地震模拟、流体动力学、空气动力学等领域的科学计算。这些应用通常涉及大量的数学计算和复杂的物理模型，计算过程中涉及到高频度通信。
- 分子动力学模拟：用于研究分子结构和化学反应动力学的过程模拟，该计算过程同样需要高频度通信。
