# 6. App Module Exit Code Specification

## 6.1 Design Goals

This specification defines the exit code conventions for app modules running in the general computing framework, used to:

- Report execution results to the framework (success / failure / classification)
- Support automatic decision-making by the scheduling system (such as retries, alerts, termination)
- Provide a consistent "minimal semantic layer" across apps

⚠️ Note:
Exit codes are only for coarse-grained classification; detailed error information should be provided through standard output (stdout), standard error (stderr), or structured logs.

## 6.2 Basic Rules
- Exit code range: 0–255 (only the lower 8 bits are valid)
- 0 means success and must be followed
- Non-zero means failure or non-success status
- 126–127 are reserved by the system and cannot be used
- 128–255 indicate the process was terminated by a signal and cannot be customized

## 6.3 Overall Partitioning

| Range   | Type         | Description                  |
| ------- | ----------- | ---------------------- |
| 0       | success     | task completed successfully |
| 1–10    | global unified semantics | understood by the framework, must be followed by all apps |
| 11–99   | app-defined | freely defined by apps; the framework does not parse details |
| 100–119 | recommended shared semantics | optional standards shared across apps |
| 120–125 | reserved    | reserved for future extensions |
| 126–127 | system reserved | shell semantics |
| 128–255 | signal exit | terminated by a signal |

## 6.4 Global Unified Semantics (1–10)

All app modules must follow the following definitions:
| Exit code | Name             |       Meaning     |
| ----- | ----------------- | -------------------- |
| 1     | UNKNOWN           | unclassified error (fallback) |
| 2     | INVALID_ARGS      | parameter error |
| 3     | BAD_CONFIG        | configuration error |
| 4     | BAD_INPUT         | input data error |
| 5     | OUTPUT_ERROR      | output failure |
| 6     | RESOURCE_LACK     | insufficient resources (memory/disk, etc.) |
| 7     | PERMISSION_DENIED | insufficient permissions |
| 8     | TIMEOUT           | execution timeout |
| 9     | DEPENDENCY_FAIL   | external dependency failure |
| 10    | INTERNAL_ERROR    | internal app error |

## 6.5 App-Defined Range (11–99)

This range is entirely defined by app modules; the framework does not parse specific meanings.

| Range | Suggested use |
| ----- | ----------- |
| 11–29 | parameter / input subdivision errors |
| 30–49 | business logic errors |
| 50–69 | data processing errors |
| 70–89 | external system / dependency errors |
| 90–99 | extensions / reserved |

✅ Recommendation: maintain your own error code documentation within the app

## 6.6 Recommended Shared Semantics (100–119)

This range is an optional standard for unified behavioral semantics across apps (particularly suitable for scheduling policies).

| Exit code | Name                 | Meaning |
| --- | -------------------- | ------ |
| 100 | RETRYABLE_ERROR      | retryable error |
| 101 | NON_RETRYABLE        | non-retryable |
| 102 | PARTIAL_SUCCESS      | partial success |
| 103 | SKIPPED              | actively skipped |
| 104 | IDEMPOTENCY_CONFLICT | idempotency conflict |
| 105 | DATA_CORRUPTION      | data corruption |
| 106 | RATE_LIMITED         | rate limited |
| 107 | DEPENDENCY_TIMEOUT   | downstream dependency timeout |
| 108 | ENVIRONMENT_ERROR    | runtime environment error |
| 109 | CANCELLED            | cancelled by upstream |

⚠️ Note:
- Not mandatory
- Recommended for cross-team adoption to improve system coordination

## 6.7 Reserved Range (120–125)

120–125 RESERVED

For future extensions; must not be used.


## 6.8 System Reserved Range

| Range   | Meaning               |
| ------- | ---------------------- |
| 126     | command exists but is not executable |
| 127     | command not found |
| 128–255 | terminated by a signal (128 + signal) |

- Common signal examples

| Exit code | Signal   | Meaning               |
| ----- | ------- | ------------------ |
| 137   | SIGKILL | usually OOM or forced termination |
| 143   | SIGTERM | normal termination (e.g. stopped by the scheduler) |


## 6.9 Framework Handling Recommendations
### 6.9.1 Basic Classification Logic

| Condition | Handling recommendation |
| -------- | ---------------- |
| exit = 0 | success |
| 1–10     | handle by unified semantics |
| 11–99    | app error |
| 100–119  | can be used for scheduling policies |
| ≥128     | signal termination (special handling) |

### 6.9.2 Scheduling Policy Examples
| Exit code     | Recommended behavior |
| ------------- | -------------- |
| 2 (parameter error) | fail directly, no retry |
| 6 (insufficient resources) | delayed retry possible |
| 8 (timeout)   | retryable |
| 100 (retryable) | retry |
| 101 (non-retryable) | fail |
| 137 (SIGKILL) | mark as system anomaly |

## 6.10 Design Principles Summary
- Minimal semantics principle: exit codes express only "classification", not details
- Layered design: unified semantics + app freedom + optional sharing
- Forward compatibility: reserved extension space
- System compatibility: follow Unix exit code conventions


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
