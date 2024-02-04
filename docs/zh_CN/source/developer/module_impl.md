# 4. 模块设计实现

## 4.1 task-key定义

## 4.2 Dockerfile定义

## 4.3 脚本编程

### 4.3.1 算法运行run.sh

### 4.3.2 初始设置setup.sh

### 4.3.3 结束退出teardown.sh

### 4.3.4 流控检查check.sh


## 4.4 集成测试

- 多个初始化消息

## 4.5 迭代优化

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
