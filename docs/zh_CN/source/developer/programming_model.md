# 1. 编程模型


- sidecar模式：
- 天然分布式设计

分层：
- 算法模块层：多种程序语言实现
- 模块连接层：shell编程，算法执行结果反馈到系统（后续消息路径、执行时间、I/O数据、执行结果（stdout/stderr））
- 并行程序全局算法间层：用yaml定义的模块间关联。

{画一个示例图}

程序表达能力：
- 结构化程序：顺序、分支、循环
- goto语句

scalebox程序，通过消息发送，可实现类似于goto语句的功能。从而支持完整的应用逻辑表达。

模块间的应用逻辑表达：基于消息发送，实现控制流，类似于goto的功能
    在并行程序中，每个模块（算法模块、传输模块、控制模块等）相当于操作原语    
    

## 标准模块的可编程特性

- 标准模块：其功能脚本可放在/app/share/bin下，子模块的功能脚本在/app/bin下。模块识别规范是优先使用/app/bin，再搜索/app/share/bin目录下；
- message-router消息格式定制：按标准模块的消息格式定制，便于使用标准模块功能；

这样可充分利用标准模块的功能。


