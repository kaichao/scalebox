# agent-interface-files

模块与 agent 之间的接口文件：任务执行元数据、I/O 计量、时间戳、下游任务。模块只需在 `${WORK_DIR}`（缺省 `/work`）下写文件，agent 在任务结束后采集并落库。

## 配置

无平台参数；模块在运行目录（`WORK_DIR`）写文件即可。文件名与格式固定。

| 文件 | 格式 | 落库位置 |
|---|---|---|
| `task-exec.yaml` | `statusCode` / `inputBytes` / `outputBytes` | `t_task_exec` 对应字段；退出码 > 0 时以退出码为准 |
| `extra-attrs.yaml` | 任意 YAML map（示例顶层键 `extra-attr-set`） | `t_task_exec.extras`（JSON） |
| `timestamps.txt` | 每行 `<timestamp>[,<label>]`；label 缺省 `no-label` | `extras.timestamps` |
| `input-files.txt` | 每行 `路径` 或 `路径,字节数`；目录递归统计 | `t_task_exec.input_bytes` / 明细 |
| `output-files.txt` | 同上 | `t_task_exec.output_bytes` / 明细 |
| `network-files.txt` | `<from\|to>,<remote_host>,<本地路径或字节数>,<远端路径>[,<远端字节数>]` | 网络 I/O 明细与汇总 |
| `removed-files.txt` | 每行一个路径 | 删除文件统计 |
| `cleanup-files.txt` | 退出前 bash 清理命令 | 尚未实现，不要依赖 |
| `auxout.txt` | 任意文本 | `t_task_exec.auxout` |
| `sink-tasks.txt` | 下游任务，7 种格式（见下） | 由 agent 创建后续 task |

`sink-tasks.txt` 的 7 种格式（`<module>` 可省略，省略时用 `SINK_MODULE`）：

| # | 格式 | 说明 |
|---|---|---|
| 1 | `<body>` | 纯文本 body |
| 2 | `{"body":"...","k":"v"}` | JSON body + headers |
| 3 | `<body>,{"k":"v"}` | 文本 body + headers |
| 4 | `{"body":"..."},{"k":"v"}` | JSON body + headers（两段） |
| 5 | `<module>,<body>` | 指定下游模块 |
| 6 | `<module>,<body>,{"k":"v"}` | 指定模块 + headers |
| 7 | `<module>,{"body":"..."},{"k":"v"}` | 指定模块 + JSON body + headers |

## 行为与不变式

- `timestamps.txt` 每行 `<timestamp>[,<label>]`，`#` 开头为注释，label 缺省 `no-label`。支持 RFC3339Nano（`2006-01-02T15:04:05.999999999Z07:00`）、RFC3339（`2006-01-02T15:04:05Z07:00`）、`2006-01-02T15:04:05[.fffffffff]`、`2006-01-02 15:04:05[.fffffffff]`；格式错误会导致任务后处理错误。
- `input-files.txt` / `output-files.txt` 的路径若以 `/local_data_root` 开头，落库前会剥掉该前缀。
- `network-files.txt` 第 1 段 `from` 表示从外部读取、`to` 表示写到外部节点；第 3 段是本地路径（自动统计大小）或直接写字节数；第 5 段缺省等于第 3 段；`remote_host` 入库前统一转主机名。
- `task-exec.yaml` 的 `statusCode` 在进程退出码 > 0 时被退出码覆盖。
- `sink-tasks.txt` 为空文件会判为任务后处理错误（`null sink-tasks.txt`），不是"无下游"；无下游时**不要创建该文件**。

## 落库字段与视图

- `input-files.txt` / `output-files.txt` 合并到 `t_task_exec.extras->>'iobytes'`，字段均为字节数：`global_input`、`global_output`、`tmpfs_input`、`tmpfs_output`、`local_input`、`local_output`、`input_bytes`、`output_bytes`。分类：global 按 cluster 的 `data_root` 或 app 的 `global_directories`；tmpfs 以 `/dev/shm` 前缀；其余为 local；`input_bytes` / `output_bytes` 为总计量。视图 `v_task_iobytes`。
- `network-files.txt` 合并到 `extras->>'network_io'`，每行 `<from_host>,<to_host>,<net_bytes>,<io_host>,<io_type>,<io_bytes>`；`io_type` 取 `global_input` / `global_output` / `local_input` / `local_output` / `tmpfs_input` / `tmpfs_output`。agent 处理后的中间格式为 `<from_host>,<to_host>,<net_bytes>,<io_host>,<io_path>,<io_bytes>`，入库时用 `io_path` 判定 `io_type`。视图 `v_task_network_io`。
- `removed-files.txt` 的路径在完成读写量统计后删除。
- `cleanup-files.txt` 尚未实现。

## 坑

- **嵌套脚本覆盖 sink 文件**：模块内调用 file-copy 等子脚本也会生成 `sink-tasks.txt`。外层必须先 `rm -f` 再写自己的，否则会重复投递或投错。
- **失败时不要写 sink 文件**：只有成功才写下游；失败任务不应触发下游。
- **文件名用连字符**：`input-files.txt`，不是下划线；`extras.yaml` 不是接口文件（接口用 `extra-attrs.yaml`）。
- **sink 行的 headers 会直传下游**：可用于带 `batch_status` 直连 `vtask-tail` 等场景。

## 最小例子

`feature.yaml`：主路由模块 `mr-module` + 业务模块 `my-module`。业务模块的 `code/run.sh` 依次生成全部接口文件；主路由把入口任务转给业务模块。

```yaml
parameters:
  main_router: mr-module
modules:
  mr-module:
    arguments:
      code_path: ${PWD}/mr-code
      slot_options: slot_on_head
  my-module:
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
```

## 测试

```sh
./test.sh   # 9 个任务：入口 → my-module → 7 条 sink 格式
# 当前手工：
echo start-task | scalebox run --app-file=feature.yaml
# 观察 my-module 的 t_task_exec：extras.timestamps、input_bytes/output_bytes、auxout、sink 任务
```

预期：`t_task_exec.extras` 含 `timestamps`（`before-sleep` / `after-sleep` 两个 label）；`sink-tasks.txt` 的 7 种格式全部创建成功。

## 源码指针

- `go-scalebox/cmd/agent/postexec.go`（`normlizeTaskExec` / `readTaskExec` / `getExtras` / `getTimestamps`）—— `task-exec.yaml`、`extra-attrs.yaml`、`timestamps.txt`、`auxout.txt`
- `go-scalebox/cmd/agent/iobytes.go`（`sumIOBytes`）—— `input-files.txt` / `output-files.txt` / `network-files.txt`
- `go-scalebox/cmd/agent/sink_tasks.go`（`addSinkTasks`）—— `sink-tasks.txt` 与 `last_from_*` 继承
- `go-scalebox/cmd/agent/run.go`—— 任务后处理调用顺序

## 相关文档

- Sphinx：`docs/cn/source/developer/2_module_development.md` §2.3.2（摘要与指针）
- 参数表：`docs/cn/source/appendix/2_parameter_reference.md` §2.3（环境变量）
- 实战：app-fxcorr `modules/README.md` §7（接口文件与日志落库）
