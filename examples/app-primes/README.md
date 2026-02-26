# app-primes

## 一、应用介绍

基于Scalebox，以分布式计算实现方式，计算给定整数区间[1..max]内的质数数量。

通过将区间分解为若干个子区间，每个区间的起止为range_start、range_end

## 二、应用设计

应用分为2个模块：主模块（路由模块）、计算模块。
- 主模块：做区间划分、计算结果汇总；以shell实现。
- 计算模块：按区间起止范围，计算出区间内质数数量，可以用不同程序语言实现。



### 2.1 主模块设计

主模块功能：
1. 计算区间分解：分解为若干消息，并启动若干个计算模块做并行计算
2. 计算结果汇总：

- 格式一：初始task
  - body：总分组数
- 格式二：中间计算结果
  - body：分组起始值
  - headers：
  - ```part-primes```：中间结果
- 环境变量
  - 区间长度：LENGTH

通过信号量```app-primes:sum_value```纪录当前计算的累计中间结果

### 2.2 计算模块设计

计算模块可以用不同程序语言实现。calc目录中为python版本。其他语言版本参见```misc/lang/```，包括以下语言实现：

- C
- FORTRAN
- BASH
- golang
- java
- nodejs
- octave/matlab
- php
- pl/pgsql(postgres) : 基于数据库存储过程实现
- R
- Rust
- scala
- julia


- 计算task格式
  - 输入格式：
    - body：${begin}_${end}，区间起止数值
  - 输出格式
    - body：同输入
    - headers：
- 环境变量

## 三、运行测试

### 3.1 构建镜像


### 3.2 命令行方式运行

- 计算1000以内的质数数量
```sh
export TAG=bash
cd app-primes/
app_id=$( echo 1000 | scalebox run | cut -d':' -f2 | tr -d '}' )
```

### 3.3 查看计算结果
```sh
APP_ID=${app_id} scalebox semaphore get app-primes:sum_value
```

## 四、高级特性

### 4.1 分布式部署
```sh
export TASK_BATCH_SIZE=5
echo 1000 | scalebox run -e inline.env
```

### 4.2 支持多消息同时处理
```sh
export TASK_BATCH_SIZE=5
echo 1000 | scalebox run
```

### 4.3 多集群部署


## 五、多语言版本测试

### 5.1 镜像构建

- bash
- c
- fortran
- golang
- java
- julia
- nodejs
- octave
- php
- python
- r
- rust
- scala

```sh
cd app-primes/misc/lang
TAG=bash make build
```

### 5.2 多语言版本运行

```sh

cd app-primes/
export TAG=bash
app_id=$( echo 1000 | scalebox run | cut -d':' -f2 | tr -d '}' )

```

### 5.3 postgres版本运行

```sh
cd app-primes/
app_id=$( echo 1000 | scalebox run --app-file=postgres.yaml | cut -d':' -f2 | tr -d '}' )
```

### 5.4 查看计算结果
```sh
APP_ID=${app_id} scalebox semaphore get app-primes:sum_value
```
