<!-- AUTO-GENERATED from `internal/task/status_code.go` (Go 常量) + 运行时约定 by `make docs-gen` (go-scalebox); DO NOT EDIT -->
<!-- source: go-scalebox@67c994e (2026-10-21) -->


# 退出码与任务状态码

> 应用模块的退出码写入 `t_task.status_code`（0–255）；调度阶段用负值表示中间态。
> 源码常量之外的 123/124 是平台运行时约定，单列在最后。

## 调度状态（负值）

| 值 | 常量 | 含义 |
|---|---|---|
| `-9` | `StatusCodeInitial` | 初始状态（动态调度使用） |
| `-3` | `StatusCodeRunning` | 任务运行中 |
| `-2` | `StatusCodeQueued` | 任务已排队，从 READY 到 RUNNING 的中间状态 |
| `-1` | `StatusCodeReady` | 任务就绪，等待调度 |

## 成功与全局统一语义（0–10）

| 值 | 常量 | 含义 |
|---|---|---|
| `0` | `StatusCodeOK` | 任务成功完成 |
| `1` | `StatusCodeUnknown` | 未分类错误（兜底） |
| `2` | `StatusCodeInvalidArgs` | 参数错误 |
| `3` | `StatusCodeBadConfig` | 配置错误 |
| `4` | `StatusCodeBadInput` | 输入数据错误 |
| `5` | `StatusCodeOutputError` | 输出失败 |
| `6` | `StatusCodeResourceLack` | 资源不足（内存/磁盘等） |
| `7` | `StatusCodePermissionDenied` | 权限不足 |
| `8` | `StatusCodeTimeout` | 执行超时 |
| `9` | `StatusCodeDependencyFail` | 外部依赖失败 |
| `10` | `StatusCodeInternalError` | 应用内部错误 |

## 推荐共享语义（100–109）

| 值 | 常量 | 含义 |
|---|---|---|
| `100` | `StatusCodeRetryableError` | 可重试错误 |
| `101` | `StatusCodeNonRetryable` | 不可重试 |
| `102` | `StatusCodePartialSuccess` | 部分成功 |
| `103` | `StatusCodeSkipped` | 主动跳过 |
| `104` | `StatusCodeIdempotencyConflict` | 幂等冲突 |
| `105` | `StatusCodeDataCorruption` | 数据损坏 |
| `106` | `StatusCodeRateLimited` | 被限流 |
| `107` | `StatusCodeDependencyTimeout` | 下游依赖超时 |
| `108` | `StatusCodeEnvironmentError` | 运行环境错误 |
| `109` | `StatusCodeCancelled` | 被上游取消 |

## 系统内部（110–125）

| 值 | 常量 | 含义 |
|---|---|---|
| `111` | `StatusCodeSetupError` | setup 错误 |
| `112` | `StatusCodeTaskPostExecError` | task-post-exec 错误 |
| `125` | `StatusCodeExecError` | 执行错误（Docker daemon error） |

## 运行时补充（非源码常量）

| 值 | 含义 | 来源 |
|---|---|---|
| `123` | 全局超时：任务长期处于运行态，由 pgcron 按 `task_max_seconds × task_global_timeout_scale` 恢复 | `build/database/initdb.d/62-pgcron_task.sql` |
| `124` | 单任务超时：agent 用 `timeout <task_max_seconds>` 包裹 ACTION_RUN/SETUP/TEARDOWN/CHECK，GNU timeout 的退出码 | `cmd/agent/run.go` |

## 退出码分区（0–255）

| 范围 | 类型 | 说明 |
|---|---|---|
| `0` | 成功 | 任务成功完成 |
| `1–10` | 全局统一语义 | 框架可理解，所有应用必须遵守 |
| `11–99` | 应用自定义 | 应用自由定义，框架不解析细节 |
| `100–119` | 推荐共享语义 | 可选标准，跨应用共享 |
| `120–125` | 保留 | 预留扩展 |
| `126–127` | 系统保留 | shell 语义 |
| `128–255` | 信号退出 | 进程被信号终止 |

> 分区为本项目约定；源码常量以本页上方表格为准。
