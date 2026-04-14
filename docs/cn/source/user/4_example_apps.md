# 4. 示例应用

## 4.1 hello-scalebox：入门级应用

第一个scalebox应用，用于验证安装和基本功能。

### 运行示例

```bash
cd examples/hello-scalebox
echo "Docker-based_Scalebox" | scalebox run 
```

### 应用特点
- 简单的消息处理
- 单模块应用
- 适合初学者理解基本概念

## 4.2 app-primes：计算密集型应用

求解[1..max_value]之间的质数数量。主要用于展示scalebox的主要特性。

### 运行示例

```bash
cd examples/app-primes
make run NUM_GROUPS=4 CALC_NODE=local NUM_PARALLEL=2
```

### 应用特点
- 数据并行处理
- 多模块协作
- 支持多语言实现（Python、Go、C++等）

## 4.3 app-copy：数据传输应用

文件拷贝操作是计算过程中常见操作。

### 运行示例

```bash
cd examples/app-copy
scalebox run --source /data/source --target /data/target
```

### 应用特点
- 跨节点数据传输
- 支持多种传输协议
- 流水线并行优化

## 4.4 remote-primes：跨集群应用

app-primes的跨集群版本，用多个集群算力求解给定整数区间内质数数量。

### 运行示例

```bash
cd examples/remote-primes
make run CLUSTER0=cluster0 CLUSTER1=cluster1
```

### 应用特点
- 跨集群计算
- 异构集群支持
- 动态资源分配

## 4.5 vtask：虚拟任务应用

### 运行示例

```bash
cd examples/vtask
scalebox app run --vtasks 100
```

### 应用特点
- 虚拟任务管理
- 批量任务处理
- 资源优化调度

## 4.6 复杂流水线应用设计模式

### 4.6.1 多级流水线

```
数据输入 → 预处理 → 计算 → 后处理 → 输出
```

### 4.6.2 并行流水线

```
       → 处理A →
输入 → 处理B → 合并 → 输出
       → 处理C →
```

### 4.6.3 混合模式

```
       → 并行A →
输入 → 串行处理 → 并行B → 输出
       → 并行C →
```

### 4.6.4 跨集群流水线

```
集群A: 数据生成 → 集群B: 数据处理 → 集群C: 结果分析
```

## 4.7 应用最佳实践

### 4.7.1 模块设计
- 保持模块无状态
- 支持任务幂等性
- 合理设置任务超时
- 优化数据局部性

### 4.7.2 性能优化
- 合理设置并行度
- 优化数据分布
- 减少通信开销
- 利用本地存储

### 4.7.3 容错设计
- 设置重试机制
- 合理设置超时时间
- 监控任务状态
- 及时处理失败任务