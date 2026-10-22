# 4. 示例应用

## 4.1 hello-scalebox：入门级应用

第一个scalebox应用，用于验证安装和基本功能。

### hello-scalebox 运行示例

```bash
cd examples/hello-scalebox
echo "Scalebox" | scalebox run 
```

### hello-scalebox 应用特点
- 简单的消息处理
- 单模块应用，没有主路由；主路由在多模块协作时出现
- 适合初学者理解基本概念

## 4.2 app-primes：计算密集型应用

求解[1..max_value]之间的质数数量。主要用于展示scalebox的主要特性。

### app-primes 运行示例

```bash
cd examples/app-primes
export NUM_GROUPS=4 CALC_NODE=local NUM_PARALLEL=2
echo 1000 | scalebox run
```

### app-primes 应用特点
- 数据并行处理
- 多模块协作
- 支持多语言实现（Python、Go、C++等）

## 4.3 app-copy：数据传输应用

文件拷贝操作是计算过程中常见操作。

### app-copy 运行示例

仓库中没有独立的 `examples/app-copy` 目录；文件拷贝由标准模块 `file-copy` /
`dir-copy` 提供，可运行示例见 `dockerfiles/files/app-dir-copy/`（`app.yaml` +
`scalebox run`），模块参数见 `dockerfiles/files/file-copy/README.md`。

### app-copy 应用特点
- 跨节点数据传输
- 支持多种传输协议
- 流水线并行优化

## 4.4 remote-primes：跨集群应用

app-primes的跨集群版本，用多个集群算力求解给定整数区间内质数数量。

### remote-primes 运行示例

```bash
cd examples/remote-primes
scalebox cluster create cluster0.yaml    # cluster1/cluster2 同理
scalebox run --app-file main.yaml
```

远端子应用（`calc.yaml --main-app-id`）与投递步骤见 `examples/remote-primes/README.md` §3。

### remote-primes 应用特点
- 跨集群计算
- 异构集群支持
- 动态资源分配

## 4.5 vtask：虚拟任务应用

### vtask 运行示例

```bash
cd examples/vtask
cat default-tasks.txt | scalebox run -e scalebox.env
```

### vtask 应用特点
- 虚拟任务管理
- 批量任务处理
- 资源优化调度

## 4.6 设计模式与最佳实践

示例只展示“怎么跑”；应用设计模式（简单流水线、多级流水线、并行流水线、混合并行、跨集群）与最佳实践（模块划分、数据流设计、性能优化、可维护性、容错设计）见 {doc}`应用设计 <../developer/3_app_design>` §3.5/§3.6。