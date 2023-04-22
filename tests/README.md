# tests(特性测试)

scalebox平台的关键特性测试。

## 容错支持：[retry_test](./retry_test/)
在模板定义文件的job定义中，设置
```
parameters:
  retry_rules: "['1','2:3']"
```
表示：若task的返回错误码为1，自动重做1次；若task的返回错误码为2，自动重做3次。
## 超时设置：[timeout-gen](./timeout-gen/)
在模板定义文件的job定义中，设置
```
variables:
  task_timeout_seconds:	10
```
表示，task的最大运行时间为10秒，若超过10秒，则该task退出，返回timeout错误码。
## 流控管理：[check_test](./check_test/)

流控管理通过设置容器的环境变量ACTION_CHECK（缺省值为:/app/bin/check.sh）来实现。若ACTION_CHECK返回值非0，则该容器将跳过server端消息获取，从而不能进行后续处理流程。

## 应用交互：[task-exec-files](./task-exec-files/)

用户应用程序与平台交互，通过以下3个文件实现
- /work/timestamps.txt：纪录用户程序分段的时间戳，可以纪录task运行数据库中。格式如下：
```
2008-03-19T18:35:03-08:00
2009-11-05T17:50:20.154+08:00
2010-11-05T17:50:20.154918+08:00
2011-11-05T17:50:20.154918780+08:00
2012-11-17T08:52:21,963572856+08:00
```

- /work/user-file.txt：用户自身产生的关于运行过程的数据文件，无格式限制，直接记录在task运行数据库中
- /work/task-exec.json：用户运行过程的各类统计数据、状态数据，需纪录到task运行数据库中。其示例格式如下：

```json
{
    "statusCode":"<status_code>",
	  "inputBytes":"<input_bytes>",
	  "outputBytes":"<output_bytes>",
    "userText":"user-defined text\nHello scalebox in message-${m}",
    "timestamps":["2018-03-19T18:35:03-08:00","2019-11-05T17:50:20.154+08:00","2020-11-05T17:50:20.154918+08:00","2021-11-05T17:50:20.154918780+08:00","2022-11-17T08:52:21,963572856+08:00"],
    "sinkJob":"task-exec-files",
    "messageBody":"1"
}
```
