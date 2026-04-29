# 6. 应用模块退出码规范

## 6.1 设计目标

本规范定义应用模块在通用计算框架中运行时的退出码约定，用于：

- 向框架报告运行结果（成功 / 失败 / 分类）
- 支持调度系统进行自动决策（如重试、告警、终止）
- 提供跨应用的一致“最小语义层”

⚠️ 说明：
退出码仅用于粗粒度分类，详细错误信息应通过标准输出（stdout）、标准错误（stderr）或结构化日志提供。

## 6.2 基本规则
- 退出码范围：0–255（仅低 8 位有效）
- 0 表示成功，必须遵守
- 非 0 表示失败或非成功状态
- 126–127 为系统保留，不可使用
- 128–255 表示进程被信号终止，不可自定义

## 6.3 总体分区

| 范围     | 类型        | 说明                    |
| ------- | ----------- | ---------------------- |
| 0       | 成功        | 任务成功完成              |
| 1–10    | 全局统一语义 | 框架可理解，所有应用必须遵守 |
| 11–99   | 应用自定义   | 应用自由定义，框架不解析细节 |
| 100–119 | 推荐共享语义 | 可选标准，跨应用共享       |
| 120–125 | 保留        | 预留未来扩展             |
| 126–127 | 系统保留     | shell 语义             |
| 128–255 | 信号退出     | 被信号终止              |

## 6.4 全局统一语义（1–10）

所有应用模块必须遵守以下定义：
| 退出码 | 名称               |       含义           |
| ----- | ----------------- | -------------------- |
| 1     | UNKNOWN           | 未分类错误（兜底）      |
| 2     | INVALID_ARGS      | 参数错误              |
| 3     | BAD_CONFIG        | 配置错误              |
| 4     | BAD_INPUT         | 输入数据错误           |
| 5     | OUTPUT_ERROR      | 输出失败              |
| 6     | RESOURCE_LACK     | 资源不足（内存/磁盘等） |
| 7     | PERMISSION_DENIED | 权限不足             |
| 8     | TIMEOUT           | 执行超时             |
| 9     | DEPENDENCY_FAIL   | 外部依赖失败          |
| 10    | INTERNAL_ERROR    | 应用内部错误          |

## 6.5 应用自定义区（11–99）

该区间完全由应用模块定义，框架不解析具体含义。

| 范围   | 建议用途        |
| ----- | ----------- |
| 11–29 | 参数 / 输入细分错误 |
| 30–49 | 业务逻辑错误      |
| 50–69 | 数据处理错误      |
| 70–89 | 外部系统 / 依赖错误 |
| 90–99 | 扩展 / 预留     |

✅ 建议：应用内部维护自己的错误码文档

## 6.6 推荐共享语义（100–119）

该区间为可选标准，用于跨应用统一行为语义（特别适用于调度策略）。

| 退出码 | 名称                   | 含义     |
| --- | -------------------- | ------ |
| 100 | RETRYABLE_ERROR      | 可重试错误  |
| 101 | NON_RETRYABLE        | 不可重试   |
| 102 | PARTIAL_SUCCESS      | 部分成功   |
| 103 | SKIPPED              | 主动跳过   |
| 104 | IDEMPOTENCY_CONFLICT | 幂等冲突   |
| 105 | DATA_CORRUPTION      | 数据损坏   |
| 106 | RATE_LIMITED         | 被限流    |
| 107 | DEPENDENCY_TIMEOUT   | 下游依赖超时 |
| 108 | ENVIRONMENT_ERROR    | 运行环境错误 |
| 109 | CANCELLED            | 被上游取消  |

⚠️ 注意：
- 不强制使用
- 建议跨团队统一采用以提升系统协同能力

## 6.7 保留区（120–125）

120–125 RESERVED

用于未来扩展，不得使用。


## 6.8 系统保留区

| 范围     | 含义                   |
| ------- | ---------------------- |
| 126     | 命令存在但不可执行        |
| 127     | 命令不存在               |
| 128–255 | 被信号终止（128 + signal）|

- 常见信号示例

| 退出码 | 信号     | 含义               |
| ----- | ------- | ------------------ |
| 137   | SIGKILL | 通常为 OOM 或强制终止 |
| 143   | SIGTERM | 正常终止（如调度器停止）|


## 6.9 框架处理建议
### 6.9.1 基本分类逻辑

| 条件       | 处理建议         |
| -------- | ---------------- |
| exit = 0 | 成功              |
| 1–10     | 按统一语义处理      |
| 11–99    | 应用错误           |
| 100–119  | 可用于调度策略      |
| ≥128     | 信号终止（特殊处理） |

### 6.9.2 调度策略示例 
| 退出码         | 建议行为        |
| ------------- | -------------- |
| 2 (参数错误)    | 直接失败，不重试 |
| 6 (资源不足)    | 可延迟重试      |
| 8 (超时)       | 可重试         |
| 100 (可重试)   | 重试           |
| 101 (不可重试)  | 失败          |
| 137 (SIGKILL) | 标记为系统异常   |

## 6.10 设计原则总结
- 最小语义原则：退出码仅表达“分类”，不表达细节
- 分层设计：统一语义 + 应用自由 + 可选共享
- 向前兼容：预留扩展空间
- 系统兼容：遵循 Unix exit code 约定


## 6.11 Status Code Table

- code range(16-bit): [-32768..32767]

| Code      | Number  | Description |
| --------- | ------- | ----------- |
| app_code  | 0~255   |             |
| OK        | 0       |             |
| READY     | -1      |             |
| QUEUED    | -2      | FROM 'READY' to 'RUNNING' |
| RUNNING   | -3      |             |
| INITIAL   | -9      | Initial Status (used by dynamic scheduling) |
| ERROR     | -32~-63 |             |

- task status_code : for task scheduling
- task_exec status_code : task-exec history recording
- ssh / docker status_code (in actuator)


## REF

### REF[1] : /usr/include/sysexits.h IN Linux

| Code        | Number      | Description |
| ----------- | ----------- | ----------- |
| ExOK |    0    |  successful termination |
| ExUsage |    64     | command line usage error |
| ExDataErr |  65    | data format error |
| ExNoInput |   66      | cannot open input |
| ExNoUser |    67   | addressee unknown (for email) |
| ExNoHost |    68   | host name unknown (for email) |
| ExUnavailable |    69   | service unavailable (for email) |
| ExSoftware |     70  | internal software error |
| ExOSErr |      71 | system error (e.g., can't fork) |
| ExOSFile |     72  | critical OS file missing |
| ExCantCreat |  73     | can't create (user) output file |
| ExIOErr |     74  | input/output error |
| ExTempFail |   75    | temp failure; user is invited to retry |
| ExProtocol |   76    | remote error in protocol |
| ExNoPerm |   77    | permission denied |
| ExConfig |   78    | configuration error |

### REF[2] : [Advanced Bash Scripting Guide, Appendix E. Exit Codes With Special Meanings](https://tldp.org/LDP/abs/html/exitcodes.html)

| Code        | Number      | Description |
| ----------- | ----------- | ----------- |
| ExGeneral |     1   |   |
| ExMisuse  |     2  |  |
| ExCantExec |       126  |  |
| ExCmdNotFound |    127  |  |
| ExInvalidExit |    128  |  |
| ExSignals |     129~165  |  |
| ExOutOfRange |      255  |  |

### Additional Scalebox Sub-Task Status Code
| Code        | Number      | Description |
| ----------- | ----------- | ----------- |
| ExUserDef |   32~63  |  |
| ExUserDef |   192~223  |  |
| ExTimeOut   |   224   |  |
| ExCoreDump   |  225    |  |
| ExNotRunnable   |  225    | run program not exists |
| ExExecNotExists   |  225    | run program not runnable |

## Exit Status Code in actuator
### Docker run Exit Status
- [Docker run Exit Status](https://docs.docker.com/engine/reference/run/#exit-status)

| Code         | Description |
| -----------  | ----------- |
|  125    | The error is with Docker daemon itself |
|  126    | The contained command cannot be invoked |
|  127    | The contained command cannot be found |

### SSH Exit Status
- [SSH and SCP Return Codes](https://support.microfocus.com/kb/doc.php?id=7021696)

| Code     | Description |
| -------- | ----------- |
| 0 | Operation was successful |
| 1 | Generic error, usually because invalid command line options or malformed configuration |
| 2 | Connection failed |
| 65 | Host not allowed to connect |
| 66 | General error in ssh protocol |
| 67 | Key exchange failed |
| 68 | Reserved |
| 69 | MAC error |
| 70 | Compression error |
| 71 | Service not available |
| 72 | Protocol version not supported |
| 73 | Host key not verifiable |
| 74 | Connection failed |
| 75 | Disconnected by application |
| 76 | Too many connections |
| 77 | Authentication cancelled by user |
| 78 | No more authentication methods available |
| 79 | Invalid user name |

### GRPC status code
- [Status codes and their use in gRPC](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)