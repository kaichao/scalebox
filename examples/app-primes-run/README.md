# app-primes-run

## 一、应用介绍

以分布式计算实现全部整数区间内的质数数量。通过将区间分解为若干个子区间，每个区间的起止为range_start、range_end

## 二、应用设计

应用分为2个模块：主模块（消息路由）、计算模块。
- 主模块负责区间划分、计算结果汇总；以shell实现。
- 计算模块负责按区间起止范围，计算出区间内质数数量，可以用不同程序语言实现。

### 2.1 主模块设计

- 消息体格式: 
- 消息头设计
- 环境变量
  - 区间长度：LEN_RANGE
  - 
### 2.2 计算模块设计

- 消息体格式
- 消息头设计
- 环境变量

## 传统方式运行
```sh
scalebox app create
```

## 分布式部署
```sh
BULK_MESSAGE_SIZE=1000 scalebox app create -e inline.env
```


## 自定义参数
```sh
NUM_GROUPS=10 \
GROUP_SIZE=100 \
scalebox app create
```

## 支持多消息同时处理
```sh
BULK_MESSAGE_SIZE=5 scalebox app create
```

## 命令行方式运行

```sh
scalebox app run --image-name=app-primes/calc ANY

echo ANY | scalebox app run --image-name=app-primes/calc
```

## 查看计算结果
```sh
APP_ID=60 scalebox semaphore get app-primes:sum_value
```
