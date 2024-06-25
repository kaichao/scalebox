# Status Code



# Task/Task_exec Exit/Return Code

| Code        | Number      | Description |
| ----------- | ----------- | ----------- |
| ExFileNotExists  |   240   |   |
| ExMessageSendException  |   245   |   |


## Status Code Table

- code range(16-bit): [-32768..32767]

| Code        | Number      | Description |
| ----------- | ----------- | ----------- |
| app_code  |   0~255   |   |
| OK  |   0   |   |
| READY  |   -1   |   |
| QUEUED |   -2   | FROM 'READY' to 'RUNNING' |
| RUNNING|   -3   |  |
| INITIAL|   -9   | Initial Status (used by dynamic scheduling) |
| ERROR|   -32~-63  |  |

- task status_code : for task scheduling
- task_exec status_code : task-exec history recording
- ssh / docker status_code (in actuator)

## ret_code vs. exit_code



## Task status_code (32-bit)
- range: [-128..-1]


-102    timeout
-1002   ?
-1003   ?



## Task_exec Status Code Structure
```
 10987654321098765432109876543210
-+------++------++------++-------
 |  sum  |prepare|cleanup| run  |
-+------++------++------++-------
```

```sh
  task_exec_code = combine(task_sum_code,prepare_code, cleanup_code, run_code)
```

### task_sum_code (left-most 8-bit, combined_code)
- range : [1..127]
- [Status codes and their use in gRPC](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)

| Code        | Number      | Description |
| ----------- | ----------- | ----------- |
| RETRIED|   10~14  | 10 ~ 12 : 3; 13 : 10; 14 : 30 |
|    | 16  | Network to Control-server Unavailable，Failed to dial target host |
|    | 17  | Control-Service Unavailable |
|    | 18  | Control-Service Unauthenticated |
|    | 19  | Control-Service Permission Denied |
|    | 20  | Control-Service Timeout |
|    | 21  | Control-Service Method Unimplemented |
|    | 24  | Invalid Input Message Format |
|    | 25  | No Valid Output  |
|    | 127 | Unknown  |


### Sub-Task Status Code (for prepare / run / cleanup)
- Range: [0..255)

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
