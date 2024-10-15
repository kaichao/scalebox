# 3. Scalebox模块定义规范

## 3.1 模块介绍

###  模块结构

附图 模块结构示意图



- 模块内主要脚本

| 环境变量 | 脚本名 | 脚本说明 |
| --------------- | ----------- | ------------ |
| ACTION_RUN      | run.sh      | 主算法脚本      |
| ACTION_CHECK    | check.sh    | 流控脚本。返回值为0，OK；非0，流控限制 |
| ACTION_SETUP    | setup.sh    | 初始设置脚本。返回值为0，OK；非0，初始化失败，设置对应slot为错误。 |
| ACTION_TEARDOWN | teardown.sh | 结束退出脚本。返回值为0，OK；非0，退出失败，设置对应slot为错误。 |


- 脚本返回码



## 3.2 消息体body格式

## 3.3 消息头格式

### 标准消息头

### 自定义消息头


## 3.4 环境变量


环境变量可为消息头设定初始值。

### 标准环境变量

- WORK_DIR: 工作目录，缺省为/work

## 3.5 封装脚本

支持用多种语言实现，推荐使用shell。

- run.sh
- setup.sh
- teardown.sh
- check.sh

### 数据目录路径

- Cluster数据目录（base_data_dir），可在容器中用/data引用。
- 本机目录：在目录前加上/local，可访问


## 3.6 返回码

- 通用返回码
  - 0：OK
  - >0：错误

- 主算法返回码

