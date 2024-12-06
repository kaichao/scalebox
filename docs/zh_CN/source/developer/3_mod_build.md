# 3. 模块构建


## 3.1 模块镜像定义


## 3.2 模块脚本脚本

sidecar模式：
run: task的单次运行
check: 检测run的前置条件
setup: 设置环境
teardown: 清除环境

### 3.2.1 算法运行run.sh

### 3.2.2 流控检测check.sh

### 3.2.3 初始设置setup.sh

### 3.2.4 结束退出teardown.sh

## 3.3 模块单元测试

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

## 3.4 迭代优化

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
