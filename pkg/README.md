# scalebox/pkg

Go SDK for scalebox controld gRPC services.

## Package overview

| Package | Implementation | gRPC Service | Notes |
|---------|:---:|------|------|
| `client` | gRPC client setup | — | Connection factory, TLS + JWT config |
| `common` | Pure Go utilities | — | File I/O, JSON, net helpers |
| `pb` | Generated protobuf | — | `protoc` output from `proto/` |
| `postgres` | Direct DB | — | `*sql.DB` / `*pgxpool.Pool` connections |
| `task` | **CLI 子进程** | — | `exec.RunReturnAll("scalebox task add ...")`，见下方说明 |
| `vtask` | gRPC | AppService | bind / unbind / add-subtask / get / fail / list / list-subtasks / semaphore / variable |
| `semaphore` | gRPC | CoordinationService | get / add / create |
| `semagroup` | gRPC | CoordinationService | min / max / diffmin / diffmax / increment / decrement |
| `variable` | gRPC | CoordinationService | get / set |
| `global` | gRPC | CoordinationService | get / set / ls |

## Why `task.Add` uses CLI subprocess instead of gRPC

`task.Add()` 是 pkg 下唯一不通过 gRPC 直连的接口。它通过 `exec.RunReturnAll` 启动 `scalebox task add` CLI 子进程完成 task 创建，而非直接调用 gRPC。三个原因：

### 1. Agent 内异步路径不可替代

CLI 的 `task add` 在 agent 容器内运行时（`PLAT_MODULE_NAME` 环境变量存在），走 sink-tasks.txt 异步路径：写入文件 → agent `addSinkTasks` 统一处理。这一步包含了 SDK 无法复刻的关键逻辑：

- `_` 前缀 header 自动传播（`_vtask_id`、`_vtask_size_sema`）
- `_vtask_*` 门控：main-router 只对 `vtask_role != ""` 的 sink 模块传播 vtask header
- `from_module` / `from_ip` 注入（由 agent 运行时上下文填充，比 CLI 推断更准确）
- 批量发送（100 条/批），共享 gRPC 连接

如果 `task.Add` 绕开 CLI 直连 gRPC，所有这些能力需要 Go SDK 重新实现，造成代码重复且行为不一致。

### 2. Task 生命周期门控

Agent 在 task 执行完毕后才调用 `addSinkTasks` 处理 sink-tasks.txt。若 task 失败（`StatusCodeTaskPostExecError`），agent 不发送下行 task。这意味着 **sink task 的创建受当前 task 的执行结果门控**。gRPC 直连会绕过这个门控——task 执行中就能创建下游 task，失败后无法撤销。

### 3. 没有对应的 gRPC 接口

AppService 没有通用的 `AddTaskList` RPC。`AddVtaskSubtask` 是 vtask 专用的，需要 `_vtask_id` header 来设置 `vtask` 列。`task.Add` 的语义是"在当前 agent 运行上下文中创建下行 task"，底层操作是写入 sink-tasks.txt，不是插入数据库。

### 设计原则

`task.Add` 委托给 CLI 子进程，确保**所有调用方（shell 脚本、Go SDK、外部应用）走同一条 task 创建路径**。无论调用来源是什么，header 传播、门控、生命周期管理的行为完全一致。这是有意为之的架构选择，不是实现限制。
