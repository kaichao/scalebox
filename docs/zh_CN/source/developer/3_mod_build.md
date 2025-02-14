# 3. 模块构建

## 3.1 模块及代码规范


### 代码结构

（模块结构图）

支持消息驱动式结构。

模块连接代码：

算法代码；

- 减少中间存储量，可以生成部分结果，尽可能早生成。

### 数据目录
通常情况下，模块处理的数据量较大，需支持灵活配置，将输入数据、中间结果、输出数据、配置文件等不同类数据放置在不同种类的存储单元（本地内存、本地SSD、本地HDD、网络存储等），进而提升I/O效率。

- 代码目录/配置目录：位于共享存储或容器镜像内；
- 输入目录：位于共享存储。本地计算模式中，由上游模块生成的数据，位于本地存储。
- 中间数据：通常位于本地存储
- 输出结果：位于共享存储。本地计算模式中，若需要在游模块使用输出数据，位于本地存储。

### 配置参数

在运行时需调整的配置参数通常以环境变量形式传递。在外部编排系统中可以修改配置参数。(运行线程数、)


## 3.2 模块镜像定义


## 3.3 模块脚本脚本

sidecar模式：
run: task的单次运行
check: 检测run的前置条件
setup: 设置环境
teardown: 清除环境

### 3.3.1 算法运行run.sh

### 3.3.2 流控检测check.sh

### 3.3.3 初始设置setup.sh

### 3.3.4 结束退出teardown.sh

## 3.4 模块单元测试

用户程序：用任意语言写；

集成程序：一般用bash写。将用户程序的结果写回。

用户程序与集成程序间接口：
- 输入消息文件：${WORK_DIR}/input-messages.txt，单次处理多个消息，可用于消息路由中批量消息的高效处理。设定环境变量BULK_MESSGES=n，缺省值为1.
- 运行结果文件：${WORK_DIR}/task-exec.json
- 用户自定义时间戳：${WORK_DIR}/timestamps.txt
- 运行附加属性文件：${WORK_DIR}/extra-attributes.txt，存放于t_task_exec表中extras的extra_attributes中。
- 用户数据文件：${WORK_DIR}/custom-out.txt
- 输入文件列表：${WORK_DIR}/input-files.txt
- 输出文件列表：${WORK_DIR}/output-files.txt
- 待删除文件列表：${WORK_DIR}/removed-files.txt
- 输出消息文件：${WORK_DIR}/output-messages.txt。原始文件为：${WORK_DIR}/messages.txt

- 多个初始化消息

## 3.5 迭代优化

用户程序（run.sh）运行结束后，agent对用户程序的运行进行统计，并纪录到核心数据库中。两者之间主要通过用户程序的标准输出（stdout）、标准错误（stderr）以及以下文件来交换信息：

|  文件名  |  文件说明    |
| --------- |  ------- |
| /work/task-exec.json |  主控制文件，以json形式，用户程序运行结果 |
| /work/messages.txt | 消息列表文件，产生后续模块所需的消息；每行一个消息 |
| /work/timestamps.txt |  时间戳文件，用于调试程序、测试程序性能时使用 |
| /work/input-files.txt |  输入文件列表，用于统计输入文件字节数 |
| /work/output-files.txt |  输出文件列表，用于统计输出文件字节数 |
| /work/removed-files.txt |  待删除文件列表，一般用于统计文件字节数后再删除 |

- 采集输入、输出字节数
- 增加时间戳
