<!-- AUTO-GENERATED from `proto/*.proto` by `make docs-gen` (go-scalebox); DO NOT EDIT -->
<!-- source: go-scalebox@3a8dd2e (2026-10-21) -->


# gRPC / Proto 参考

> controld 的 gRPC 接口事实：服务、方法、消息字段。谁调用、何时调用等语义
> 不在本页，见 `docs/agents/` 与 `features/`。

## 概览

| 文件 | package | 服务 | 消息 | 枚举 |
|---|---|---|---|---|
| `app.proto` | `scalebox` | 1 | 82 | 0 |
| `control.proto` | `scalebox` | 1 | 12 | 3 |
| `coordination.proto` | `scalebox` | 1 | 40 | 0 |
| `dashboard.proto` | `scalebox` | 1 | 12 | 0 |
| `error_detail.proto` | `scalebox` | 0 | 1 | 0 |
| `infra.proto` | `scalebox` | 1 | 44 | 0 |
| `user.proto` | `scalebox` | 1 | 14 | 0 |

## app.proto

### 服务

#### `AppService`

| 方法 | 请求 | 响应 | 说明 |
|---|---|---|---|
| `ListApps` | `ListAppsRequest` | `ListAppsResponse` |  |
| `GetApp` | `GetAppRequest` | `GetAppResponse` |  |
| `CreateApp` | `CreateAppRequest` | `CreateAppResponse` |  |
| `DeleteApp` | `DeleteAppRequest` | `google.protobuf.Empty` |  |
| `SetAppStatus` | `SetAppStatusRequest` | `google.protobuf.Empty` |  |
| `SetAppFinished` | `SetAppFinishedRequest` | `google.protobuf.Empty` |  |
| `AppAddSlots` | `AppAddSlotsRequest` | `AppAddSlotsResponse` |  |
| `GetAppRouter` | `GetAppRouterRequest` | `GetAppRouterResponse` |  |
| `GetAppDAG` | `GetAppDAGRequest` | `GetAppDAGResponse` |  |
| `ListModules` | `ListModulesRequest` | `ListModulesResponse` |  |
| `GetModule` | `GetModuleRequest` | `GetModuleResponse` |  |
| `ListTasks` | `ListTasksRequest` | `ListTasksResponse` |  |
| `GetTask` | `GetTaskRequest` | `GetTaskResponse` |  |
| `GetTaskHeader` | `GetTaskHeaderRequest` | `GetTaskHeaderResponse` |  |
| `SetTaskHeader` | `SetTaskHeaderRequest` | `google.protobuf.Empty` |  |
| `RemoveTaskHeader` | `RemoveTaskHeaderRequest` | `google.protobuf.Empty` |  |
| `DeleteTask` | `DeleteTaskRequest` | `google.protobuf.Empty` |  |
| `ListTaskExecs` | `ListTaskExecsRequest` | `ListTaskExecsResponse` |  |
| `GetTaskExec` | `GetTaskExecRequest` | `TaskExecRecord` |  |
| `RegisterRemoteAppLink` | `RegisterRemoteAppLinkRequest` | `RegisterRemoteAppLinkResponse` |  |
| `ValidateApp` | `ValidateAppRequest` | `ValidateAppResponse` |  |
| `GetVtaskVariable` | `GetVtaskVariableRequest` | `GetVtaskVariableResponse` |  |
| `SetVtaskVariable` | `SetVtaskVariableRequest` | `SetVtaskVariableResponse` |  |
| `GetVtaskSemaphore` | `GetVtaskSemaphoreRequest` | `GetVtaskSemaphoreResponse` |  |
| `CreateVtaskSemaphore` | `CreateVtaskSemaphoreRequest` | `CreateVtaskSemaphoreResponse` |  |
| `AddVtaskSemaphore` | `AddVtaskSemaphoreRequest` | `AddVtaskSemaphoreResponse` |  |
| `DeleteVtaskSemaphore` | `DeleteVtaskSemaphoreRequest` | `google.protobuf.Empty` |  |
| `ListVTasks` | `ListVTasksRequest` | `ListVTasksResponse` |  |
| `ListVTaskSubtasks` | `ListVTaskSubtasksRequest` | `ListVTasksResponse` |  |
| `BindVtaskResource` | `BindVtaskResourceRequest` | `BindVtaskResourceResponse` |  |
| `UnbindVtaskResource` | `UnbindVtaskResourceRequest` | `google.protobuf.Empty` |  |
| `AddVtaskSubtask` | `AddVtaskSubtaskRequest` | `AddVtaskSubtaskResponse` |  |
| `GetVtask` | `GetVtaskRequest` | `GetVtaskResponse` |  |
| `FailVtask` | `FailVtaskRequest` | `google.protobuf.Empty` |  |
| `GetVtaskSharedBodyModuleStats` | `GetVtaskSharedBodyModuleStatsRequest` | `GetVtaskSharedBodyModuleStatsResponse` |  |
| `ListVtaskSharedBodyTaskMatrix` | `ListVtaskSharedBodyTaskMatrixRequest` | `ListVtaskSharedBodyTaskMatrixResponse` |  |
| `ListVtaskSharedBodyEntries` | `ListVtaskSharedBodyEntriesRequest` | `ListVtaskSharedBodyEntriesResponse` |  |
| `RetryVtaskSharedBodyTask` | `RetryVtaskSharedBodyTaskRequest` | `RetryVtaskSharedBodyTaskResponse` |  |
| `LookupVtaskIdBySharedBody` | `LookupVtaskIdBySharedBodyRequest` | `LookupVtaskIdBySharedBodyResponse` |  |

### 消息

| 消息 | 字段 | 类型 | 编号 | 说明 |
|---|---|---|---|---|
| `Pagination` | `page` | `int32` | 1 |  |
| `Pagination` | `page_size` | `int32` | 2 |  |
| `PageInfo` | `page` | `int32` | 1 |  |
| `PageInfo` | `page_size` | `int32` | 2 |  |
| `PageInfo` | `total` | `int32` | 3 |  |
| `ListTasksRequest` | `app_id` | `int32` | 1 |  |
| `ListTasksRequest` | `module_id` | `int32` | 2 |  |
| `ListTasksRequest` | `status` | `string` | 3 |  |
| `ListTasksRequest` | `from_module` | `string` | 4 |  |
| `ListTasksRequest` | `from_time` | `google.protobuf.Timestamp` | 5 |  |
| `ListTasksRequest` | `to_time` | `google.protobuf.Timestamp` | 6 |  |
| `ListTasksRequest` | `status_class` | `string` | 7 |  |
| `ListTasksRequest` | `pagination` | `Pagination` | 10 |  |
| `ListTasksResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListTasksResponse` | `page_info` | `PageInfo` | 2 |  |
| `ListTasksResponse.Item` | `id` | `int64` | 1 |  |
| `ListTasksResponse.Item` | `module_id` | `int32` | 2 |  |
| `ListTasksResponse.Item` | `module_name` | `string` | 3 |  |
| `ListTasksResponse.Item` | `status` | `string` | 4 |  |
| `ListTasksResponse.Item` | `body` | `string` | 5 |  |
| `ListTasksResponse.Item` | `from_module` | `string` | 6 |  |
| `ListTasksResponse.Item` | `from_ip` | `string` | 7 |  |
| `ListTasksResponse.Item` | `created_at` | `google.protobuf.Timestamp` | 8 |  |
| `ListTasksResponse.Item` | `finished_at` | `google.protobuf.Timestamp` | 9 |  |
| `ListTasksResponse.Item` | `status_code` | `string` | 11 |  |
| `ListTasksResponse.Item` | `to_host` | `string` | 12 |  |
| `ListTasksResponse.Item` | `vtask_id` | `int64` | 13 |  |
| `GetTaskRequest` | `task_id` | `int64` | 1 |  |
| `GetTaskRequest` | `omit_exec_logs` | `bool` | 2 |  |
| `GetTaskResponse` | `id` | `int64` | 1 |  |
| `GetTaskResponse` | `module_id` | `int32` | 2 |  |
| `GetTaskResponse` | `module_name` | `string` | 3 |  |
| `GetTaskResponse` | `status` | `string` | 4 |  |
| `GetTaskResponse` | `body` | `string` | 5 |  |
| `GetTaskResponse` | `headers` | `string` | 6 |  |
| `GetTaskResponse` | `stdout` | `string` | 7 |  |
| `GetTaskResponse` | `stderr` | `string` | 8 |  |
| `GetTaskResponse` | `from_module` | `string` | 9 |  |
| `GetTaskResponse` | `from_ip` | `string` | 10 |  |
| `GetTaskResponse` | `input_bytes` | `int64` | 12 |  |
| `GetTaskResponse` | `output_bytes` | `int64` | 13 |  |
| `GetTaskResponse` | `extras` | `string` | 14 |  |
| `GetTaskResponse` | `created_at` | `google.protobuf.Timestamp` | 15 |  |
| `GetTaskResponse` | `started_at` | `google.protobuf.Timestamp` | 16 |  |
| `GetTaskResponse` | `finished_at` | `google.protobuf.Timestamp` | 17 |  |
| `GetTaskResponse` | `vtask_id` | `int64` | 18 |  |
| `GetTaskResponse` | `to_host` | `string` | 19 |  |
| `GetTaskResponse` | `status_code` | `string` | 20 |  |
| `GetTaskResponse` | `exec_records` | `TaskExecRecord (repeated)` | 21 |  |
| `TaskExecRecord` | `id` | `int64` | 1 |  |
| `TaskExecRecord` | `task_id` | `int64` | 2 |  |
| `TaskExecRecord` | `input_bytes` | `int64` | 4 |  |
| `TaskExecRecord` | `output_bytes` | `int64` | 5 |  |
| `TaskExecRecord` | `stdout` | `string` | 6 |  |
| `TaskExecRecord` | `stderr` | `string` | 7 |  |
| `TaskExecRecord` | `extras` | `string` | 8 |  |
| `TaskExecRecord` | `started_at` | `google.protobuf.Timestamp` | 9 |  |
| `TaskExecRecord` | `finished_at` | `google.protobuf.Timestamp` | 10 |  |
| `TaskExecRecord` | `mod_id` | `int32` | 11 |  |
| `TaskExecRecord` | `slot` | `int32` | 12 |  |
| `TaskExecRecord` | `status_code` | `string` | 13 |  |
| `TaskExecRecord` | `body` | `string` | 14 |  |
| `TaskExecRecord` | `auxout` | `string` | 15 |  |
| `TaskExecRecord` | `ctxlog` | `string` | 16 |  |
| `TaskExecRecord` | `t0` | `google.protobuf.Timestamp` | 17 |  |
| `TaskExecRecord` | `t1` | `google.protobuf.Timestamp` | 18 |  |
| `TaskExecRecord` | `t2` | `google.protobuf.Timestamp` | 19 |  |
| `TaskExecRecord` | `t3` | `google.protobuf.Timestamp` | 20 |  |
| `TaskExecRecord` | `ts_prepared` | `google.protobuf.Timestamp` | 21 |  |
| `TaskExecRecord` | `has_stdout` | `bool` | 22 |  |
| `TaskExecRecord` | `has_stderr` | `bool` | 23 |  |
| `TaskExecRecord` | `has_auxout` | `bool` | 24 |  |
| `TaskExecRecord` | `has_ctxlog` | `bool` | 25 |  |
| `TaskExecRecord` | `has_extras` | `bool` | 26 |  |
| `ListTaskExecsRequest` | `app_id` | `int32` | 1 |  |
| `ListTaskExecsRequest` | `module_id` | `int32` | 2 |  |
| `ListTaskExecsRequest` | `slot_id` | `int32` | 3 |  |
| `ListTaskExecsRequest` | `pagination` | `Pagination` | 10 |  |
| `ListTaskExecsResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListTaskExecsResponse` | `page_info` | `PageInfo` | 2 |  |
| `ListTaskExecsResponse.Item` | `id` | `int64` | 1 |  |
| `ListTaskExecsResponse.Item` | `task_id` | `int64` | 2 |  |
| `ListTaskExecsResponse.Item` | `mod_id` | `int32` | 3 |  |
| `ListTaskExecsResponse.Item` | `slot` | `int32` | 4 |  |
| `ListTaskExecsResponse.Item` | `status_code` | `string` | 5 |  |
| `ListTaskExecsResponse.Item` | `body` | `string` | 6 |  |
| `ListTaskExecsResponse.Item` | `input_bytes` | `int64` | 7 |  |
| `ListTaskExecsResponse.Item` | `output_bytes` | `int64` | 8 |  |
| `ListTaskExecsResponse.Item` | `started_at` | `google.protobuf.Timestamp` | 9 |  |
| `ListTaskExecsResponse.Item` | `finished_at` | `google.protobuf.Timestamp` | 10 |  |
| `GetTaskExecRequest` | `exec_id` | `int64` | 1 |  |
| `GetTaskHeaderRequest` | `task_id` | `int64` | 1 |  |
| `GetTaskHeaderRequest` | `name` | `string` | 2 |  |
| `GetTaskHeaderResponse` | `value` | `string` | 1 |  |
| `SetTaskHeaderRequest` | `task_id` | `int64` | 1 |  |
| `SetTaskHeaderRequest` | `name` | `string` | 2 |  |
| `SetTaskHeaderRequest` | `value` | `string` | 3 |  |
| `RemoveTaskHeaderRequest` | `task_id` | `int64` | 1 |  |
| `RemoveTaskHeaderRequest` | `name` | `string` | 2 |  |
| `ListAppsRequest` | `cluster_name` | `string` | 4 |  |
| `ListAppsRequest` | `status` | `string` | 2 |  |
| `ListAppsRequest` | `user_id` | `int32` | 3 |  |
| `ListAppsRequest` | `pagination` | `Pagination` | 10 |  |
| `ListAppsResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListAppsResponse` | `page_info` | `PageInfo` | 2 |  |
| `ListAppsResponse.Item` | `id` | `int32` | 1 |  |
| `ListAppsResponse.Item` | `name` | `string` | 2 |  |
| `ListAppsResponse.Item` | `cluster_name` | `string` | 4 |  |
| `ListAppsResponse.Item` | `status` | `string` | 5 |  |
| `ListAppsResponse.Item` | `module_count` | `int32` | 6 |  |
| `ListAppsResponse.Item` | `user_id` | `int32` | 7 |  |
| `ListAppsResponse.Item` | `created_at` | `google.protobuf.Timestamp` | 8 |  |
| `ListAppsResponse.Item` | `ext` | `string` | 9 |  |
| `GetAppRequest` | `app_id` | `int32` | 1 |  |
| `GetAppResponse` | `id` | `int32` | 1 |  |
| `GetAppResponse` | `name` | `string` | 2 |  |
| `GetAppResponse` | `cluster_name` | `string` | 8 |  |
| `GetAppResponse` | `status` | `string` | 4 |  |
| `GetAppResponse` | `user_id` | `int32` | 5 |  |
| `GetAppResponse` | `parameters` | `string` | 6 |  |
| `GetAppResponse` | `created_at` | `google.protobuf.Timestamp` | 7 |  |
| `GetAppResponse` | `yaml_text` | `string` | 9 |  |
| `CreateAppRequest` | `yaml_content` | `string` | 1 |  |
| `CreateAppRequest` | `env_vars` | `map<string, string>` | 2 |  |
| `CreateAppRequest` | `main_app_id` | `int32` | 3 |  |
| `CreateAppRequest` | `main_grpc_server` | `string` | 4 |  |
| `CreateAppResponse` | `app_id` | `int32` | 1 |  |
| `CreateAppResponse` | `name` | `string` | 2 |  |
| `DeleteAppRequest` | `app_id` | `int32` | 1 |  |
| `SetAppStatusRequest` | `app_id` | `int32` | 1 |  |
| `SetAppStatusRequest` | `status` | `string` | 2 |  |
| `SetAppFinishedRequest` | `app_id` | `int32` | 1 |  |
| `SetAppFinishedRequest` | `result` | `string` | 2 |  |
| `AppAddSlotsRequest` | `app_id` | `int32` | 1 |  |
| `AppAddSlotsRequest` | `module_id` | `int32` | 2 |  |
| `AppAddSlotsRequest` | `count` | `int32` | 3 |  |
| `AppAddSlotsRequest` | `host_regex` | `string` | 4 |  |
| `AppAddSlotsResponse` | `added` | `int32` | 1 |  |
| `GetAppRouterRequest` | `app_id` | `int32` | 1 |  |
| `GetAppRouterResponse` | `module_id` | `int32` | 1 |  |
| `GetAppRouterResponse` | `module_name` | `string` | 2 |  |
| `GetAppDAGRequest` | `app_id` | `int32` | 1 |  |
| `GetAppDAGRequest` | `dag_mode` | `string` | 2 |  |
| `GetAppDAGResponse` | `nodes` | `ModuleNode (repeated)` | 1 |  |
| `GetAppDAGResponse` | `edges` | `Edge (repeated)` | 2 |  |
| `GetAppDAGResponse` | `dag_mode` | `string` | 3 |  |
| `GetAppDAGResponse.ModuleNode` | `name` | `string` | 1 |  |
| `GetAppDAGResponse.ModuleNode` | `slot_count` | `int32` | 2 |  |
| `GetAppDAGResponse.ModuleNode` | `slot_regex` | `string` | 3 |  |
| `GetAppDAGResponse.ModuleNode` | `status` | `string` | 4 |  |
| `GetAppDAGResponse.ModuleNode` | `image` | `string` | 5 |  |
| `GetAppDAGResponse.ModuleNode` | `cluster` | `string` | 6 |  |
| `GetAppDAGResponse.Edge` | `from` | `string` | 1 |  |
| `GetAppDAGResponse.Edge` | `to` | `string` | 2 |  |
| `ListModulesRequest` | `app_id` | `int32` | 1 |  |
| `ListModulesResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListModulesResponse.Item` | `id` | `int32` | 1 |  |
| `ListModulesResponse.Item` | `app_id` | `int32` | 2 |  |
| `ListModulesResponse.Item` | `name` | `string` | 3 |  |
| `ListModulesResponse.Item` | `image` | `string` | 4 |  |
| `ListModulesResponse.Item` | `slot_regex` | `string` | 5 |  |
| `ListModulesResponse.Item` | `status` | `string` | 6 |  |
| `ListModulesResponse.Item` | `slot_count` | `int32` | 7 |  |
| `ListModulesResponse.Item` | `is_main_router` | `bool` | 8 |  |
| `ListModulesResponse.Item` | `cluster` | `string` | 9 |  |
| `GetModuleRequest` | `app_id` | `int32` | 1 |  |
| `GetModuleRequest` | `module_name` | `string` | 2 |  |
| `GetModuleResponse` | `id` | `int32` | 1 |  |
| `GetModuleResponse` | `name` | `string` | 2 |  |
| `GetModuleResponse` | `image` | `string` | 3 |  |
| `GetModuleResponse` | `slot_regex` | `string` | 4 |  |
| `GetModuleResponse` | `status` | `string` | 5 |  |
| `GetModuleResponse` | `slot_count` | `int32` | 6 |  |
| `GetModuleResponse` | `is_main_router` | `bool` | 7 |  |
| `GetModuleResponse` | `parameters` | `map<string, string>` | 8 |  |
| `RegisterRemoteAppLinkRequest` | `main_app_id` | `int32` | 1 |  |
| `RegisterRemoteAppLinkRequest` | `new_app_id` | `int32` | 2 |  |
| `RegisterRemoteAppLinkRequest` | `new_cluster` | `string` | 3 |  |
| `RegisterRemoteAppLinkRequest` | `yaml_text` | `string` | 4 |  |
| `RegisterRemoteAppLinkResponse` | `local_remote_links` | `RemoteLinkEntry (repeated)` | 1 |  |
| `RegisterRemoteAppLinkResponse` | `local_vlinks` | `ModuleVLinkEntry (repeated)` | 2 |  |
| `RemoteLinkEntry` | `app_id` | `int32` | 1 |  |
| `RemoteLinkEntry` | `remote_cluster` | `string` | 2 |  |
| `RemoteLinkEntry` | `remote_app_id` | `int32` | 3 |  |
| `ModuleVLinkEntry` | `remote_cluster` | `string` | 1 |  |
| `ModuleVLinkEntry` | `remote_mod_id` | `int32` | 2 |  |
| `ModuleVLinkEntry` | `remote_mod_name` | `string` | 3 |  |
| `ModuleVLinkEntry` | `mod_id` | `int32` | 4 |  |
| `ModuleVLinkEntry` | `mod_name` | `string` | 5 |  |
| `GetVtaskVariableRequest` | `app_id` | `int32` | 1 |  |
| `GetVtaskVariableRequest` | `vtask_id` | `int64` | 2 |  |
| `GetVtaskVariableRequest` | `name` | `string` | 3 |  |
| `GetVtaskVariableResponse` | `name` | `string` | 1 |  |
| `GetVtaskVariableResponse` | `value` | `string` | 2 |  |
| `SetVtaskVariableRequest` | `app_id` | `int32` | 1 |  |
| `SetVtaskVariableRequest` | `vtask_id` | `int64` | 2 |  |
| `SetVtaskVariableRequest` | `name` | `string` | 3 |  |
| `SetVtaskVariableRequest` | `value` | `string` | 4 |  |
| `SetVtaskVariableResponse` | `name` | `string` | 1 |  |
| `SetVtaskVariableResponse` | `value` | `string` | 2 |  |
| `GetVtaskSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `GetVtaskSemaphoreRequest` | `vtask_id` | `int64` | 2 |  |
| `GetVtaskSemaphoreRequest` | `name` | `string` | 3 |  |
| `GetVtaskSemaphoreResponse` | `name` | `string` | 1 |  |
| `GetVtaskSemaphoreResponse` | `value` | `int32` | 2 |  |
| `CreateVtaskSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `CreateVtaskSemaphoreRequest` | `vtask_id` | `int64` | 2 |  |
| `CreateVtaskSemaphoreRequest` | `name` | `string` | 3 |  |
| `CreateVtaskSemaphoreRequest` | `value` | `int32` | 4 |  |
| `CreateVtaskSemaphoreResponse` | `name` | `string` | 1 |  |
| `CreateVtaskSemaphoreResponse` | `value` | `int32` | 2 |  |
| `AddVtaskSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `AddVtaskSemaphoreRequest` | `vtask_id` | `int64` | 2 |  |
| `AddVtaskSemaphoreRequest` | `name` | `string` | 3 |  |
| `AddVtaskSemaphoreRequest` | `delta` | `int32` | 4 |  |
| `AddVtaskSemaphoreResponse` | `name` | `string` | 1 |  |
| `AddVtaskSemaphoreResponse` | `value` | `int32` | 2 |  |
| `DeleteVtaskSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `DeleteVtaskSemaphoreRequest` | `vtask_id` | `int64` | 2 |  |
| `DeleteVtaskSemaphoreRequest` | `name` | `string` | 3 |  |
| `ListVTasksRequest` | `app_id` | `int32` | 1 |  |
| `ListVTasksResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListVTasksResponse.Item` | `id` | `int64` | 1 |  |
| `ListVTasksResponse.Item` | `module_id` | `int32` | 2 |  |
| `ListVTasksResponse.Item` | `module_name` | `string` | 3 |  |
| `ListVTasksResponse.Item` | `status` | `string` | 4 |  |
| `ListVTasksResponse.Item` | `body` | `string` | 5 |  |
| `ListVTasksResponse.Item` | `created_at` | `google.protobuf.Timestamp` | 6 |  |
| `ListVTasksResponse.Item` | `from_module` | `string` | 7 |  |
| `ListVTasksResponse.Item` | `app_id` | `int32` | 8 |  |
| `ListVTasksResponse.Item` | `status_code` | `string` | 9 | 根 task 原始 status_code（-1/-2/-3/0/1） |
| `ListVTaskSubtasksRequest` | `vtask_id` | `int64` | 1 |  |
| `BindVtaskResourceRequest` | `app_id` | `int32` | 1 |  |
| `BindVtaskResourceRequest` | `module_name` | `string` | 2 | 可选，不指定时自动查找 vtask_role=head 的模块 |
| `BindVtaskResourceResponse` | `resource` | `string` | 1 | 绑定的资源标识：hostname（HOST-BOUND vtask）/ 组号（GROUP-BOUND vtask，即 head slot 的 seq）/ 空（DEFAULT） |
| `BindVtaskResourceResponse` | `mode` | `string` | 2 | 实际使用的 task_dist_mode |
| `UnbindVtaskResourceRequest` | `app_id` | `int32` | 1 |  |
| `UnbindVtaskResourceRequest` | `vtask_id` | `int64` | 2 | 方式一：通过 vtask ID 自动查找 sema |
| `UnbindVtaskResourceRequest` | `sema_name` | `string` | 3 | 方式二：直接指定信号量名（如 host_vtask_size:vtask-head:n0-0） |
| `AddVtaskSubtaskRequest` | `app_id` | `int32` | 1 |  |
| `AddVtaskSubtaskRequest` | `module` | `string` | 2 | 目标模块名 |
| `AddVtaskSubtaskRequest` | `body` | `string` | 3 | 任务体 |
| `AddVtaskSubtaskRequest` | `headers` | `map<string, string>` | 4 | 额外 headers |
| `AddVtaskSubtaskResponse` | `task_id` | `int64` | 1 |  |
| `GetVtaskRequest` | `vtask_id` | `int64` | 1 |  |
| `GetVtaskResponse` | `id` | `int64` | 1 |  |
| `GetVtaskResponse` | `app_id` | `int32` | 2 |  |
| `GetVtaskResponse` | `body` | `string` | 3 |  |
| `GetVtaskResponse` | `status` | `string` | 4 | 派生状态 |
| `GetVtaskResponse` | `module_name` | `string` | 5 | 当前所在模块 |
| `GetVtaskResponse` | `bound_resource` | `string` | 6 | 绑定的 host/slot |
| `GetVtaskResponse` | `sema_name` | `string` | 7 | _vtask_size_sema 值 |
| `GetVtaskResponse` | `subtask_count` | `int32` | 8 |  |
| `GetVtaskResponse` | `finished_count` | `int32` | 9 |  |
| `GetVtaskResponse` | `created_at` | `google.protobuf.Timestamp` | 10 |  |
| `GetVtaskResponse` | `status_code` | `string` | 11 | 根 task 原始 status_code（-1/-2/-3/0/1） |
| `FailVtaskRequest` | `vtask_id` | `int64` | 1 |  |
| `GetVtaskSharedBodyModuleStatsRequest` | `vtask_id` | `int64` | 1 |  |
| `GetVtaskSharedBodyModuleStatsRequest` | `modules` | `string (repeated)` | 2 | 必选，首元素为基准模块 |
| `GetVtaskSharedBodyModuleStatsResponse` | `vtask_id` | `int64` | 1 |  |
| `GetVtaskSharedBodyModuleStatsResponse` | `modules` | `string (repeated)` | 2 | 回显（按输入顺序） |
| `GetVtaskSharedBodyModuleStatsResponse` | `task_count` | `int64` | 3 | 首模块的 task 数 |
| `GetVtaskSharedBodyModuleStatsResponse` | `finished_count` | `int64` | 4 | 全部指定模块 status_code=0 的 body 数 |
| `GetVtaskSharedBodyModuleStatsResponse` | `input_bytes` | `int64` | 5 | 指定模块累加 |
| `GetVtaskSharedBodyModuleStatsResponse` | `output_bytes` | `int64` | 6 | 指定模块累加 |
| `GetVtaskSharedBodyModuleStatsResponse` | `begin_time` | `google.protobuf.Timestamp` | 7 | 指定模块 MIN(begin_time) |
| `GetVtaskSharedBodyModuleStatsResponse` | `end_time` | `google.protobuf.Timestamp` | 8 | 指定模块 MAX(end_time) |
| `VtaskSharedBodyTaskFields` | `task_id` | `int64` | 1 |  |
| `VtaskSharedBodyTaskFields` | `status_code` | `string` | 2 | 与既有消息一致，用 string 承载整数状态码 |
| `VtaskSharedBodyTaskFields` | `input_bytes` | `int64` | 3 |  |
| `VtaskSharedBodyTaskFields` | `output_bytes` | `int64` | 4 |  |
| `VtaskSharedBodyTaskFields` | `begin_time` | `google.protobuf.Timestamp` | 5 |  |
| `VtaskSharedBodyTaskFields` | `end_time` | `google.protobuf.Timestamp` | 6 |  |
| `VtaskSharedBodyTaskFields` | `headers` | `map<string, string>` | 7 | 仅填充请求的 header 名 |
| `ListVtaskSharedBodyTaskMatrixRequest` | `vtask_id` | `int64` | 1 |  |
| `ListVtaskSharedBodyTaskMatrixRequest` | `modules` | `string (repeated)` | 2 | 空 = 全部模块 |
| `ListVtaskSharedBodyTaskMatrixRequest` | `fields` | `string (repeated)` | 3 | 空 = 全部字段 |
| `ListVtaskSharedBodyTaskMatrixRequest` | `headers` | `string (repeated)` | 4 | 请求的 header 名，空 = 不输出 headers |
| `ListVtaskSharedBodyTaskMatrixRequest` | `body` | `string` | 5 | 模糊匹配词，空 = 不过滤 |
| `ListVtaskSharedBodyTaskMatrixRequest` | `status` | `string` | 6 | "abnormal" | ""（= all） |
| `ListVtaskSharedBodyTaskMatrixRequest` | `pagination` | `Pagination` | 10 |  |
| `ListVtaskSharedBodyTaskMatrixResponse` | `items` | `Row (repeated)` | 1 |  |
| `ListVtaskSharedBodyTaskMatrixResponse` | `page_info` | `PageInfo` | 2 |  |
| `ListVtaskSharedBodyTaskMatrixResponse.Row` | `body` | `string` | 1 |  |
| `ListVtaskSharedBodyTaskMatrixResponse.Row` | `modules` | `map<string, VtaskSharedBodyTaskFields>` | 2 | key = 模块名 |
| `ListVtaskSharedBodyEntriesRequest` | `vtask_id` | `int64` | 1 |  |
| `ListVtaskSharedBodyEntriesRequest` | `prefix` | `string` | 2 |  |
| `ListVtaskSharedBodyEntriesRequest` | `modules` | `string (repeated)` | 3 | 空 = 全部模块 |
| `ListVtaskSharedBodyEntriesRequest` | `fields` | `string (repeated)` | 4 | 控制文件条目字段集，空 = 全部字段 |
| `ListVtaskSharedBodyEntriesRequest` | `headers` | `string (repeated)` | 5 | 请求的 header 名，空 = 不输出 headers |
| `ListVtaskSharedBodyEntriesRequest` | `pagination` | `Pagination` | 10 |  |
| `ListVtaskSharedBodyEntriesResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListVtaskSharedBodyEntriesResponse` | `page_info` | `PageInfo` | 2 |  |
| `ListVtaskSharedBodyEntriesResponse.Item` | `name` | `string` | 1 | 相对 prefix 的下一层名称 |
| `ListVtaskSharedBodyEntriesResponse.Item` | `type` | `string` | 2 | "dir" | "file" |
| `ListVtaskSharedBodyEntriesResponse.Item` | `body_count` | `int64` | 3 | dir：目录下（任意深度）去重 body 数 |
| `ListVtaskSharedBodyEntriesResponse.Item` | `begin_time` | `google.protobuf.Timestamp` | 4 | dir：MIN(begin_time) |
| `ListVtaskSharedBodyEntriesResponse.Item` | `end_time` | `google.protobuf.Timestamp` | 5 | dir：MAX(end_time) |
| `ListVtaskSharedBodyEntriesResponse.Item` | `modules` | `map<string, VtaskSharedBodyTaskFields>` | 6 | file：各模块字段集 |
| `ListVtaskSharedBodyEntriesResponse.Item` | `body` | `string` | 7 | file：原始 t.body（管道入口）；dir 为空 |
| `RetryVtaskSharedBodyTaskRequest` | `vtask_id` | `int64` | 1 |  |
| `RetryVtaskSharedBodyTaskRequest` | `body` | `string` | 2 |  |
| `RetryVtaskSharedBodyTaskResponse` | `task_id` | `int64` | 1 |  |
| `RetryVtaskSharedBodyTaskResponse` | `deleted` | `int64` | 2 |  |
| `LookupVtaskIdBySharedBodyRequest` | `app_id` | `int32` | 1 |  |
| `LookupVtaskIdBySharedBodyRequest` | `module` | `string` | 2 |  |
| `LookupVtaskIdBySharedBodyRequest` | `body` | `string` | 3 | 输入 body，可为 JSON 文本；服务端归一化 |
| `LookupVtaskIdBySharedBodyResponse` | `vtask_id` | `int64` | 1 |  |
| `LookupVtaskIdBySharedBodyResponse` | `task_id` | `int64` | 2 |  |
| `LookupVtaskIdBySharedBodyResponse` | `body` | `string` | 3 | 命中的原始 t.body |
| `DeleteTaskRequest` | `task_id` | `int64` | 1 |  |
| `ValidateAppRequest` | `yaml_content` | `string` | 1 |  |
| `ValidateAppResponse` | `valid` | `bool` | 1 |  |
| `ValidateAppResponse` | `app_name` | `string` | 2 |  |
| `ValidateAppResponse` | `module_count` | `int32` | 3 |  |
| `ValidateAppResponse` | `errors` | `string (repeated)` | 4 |  |
| `ValidateAppResponse` | `warnings` | `string (repeated)` | 5 |  |


## control.proto

### 服务

#### `ControlService`

| 方法 | 请求 | 响应 | 说明 |
|---|---|---|---|
| `GetSlotList` | `google.protobuf.StringValue` | `SlotList` |  |
| `SetSlotInitialized` | `google.protobuf.Int32Value` | `SetSlotInitializedResponse` |  |
| `SetSlotTerminated` | `google.protobuf.Int32Value` | `google.protobuf.Empty` |  |
| `RegisterSlot` | `google.protobuf.StringValue` | `RegisterSlotResponse` |  |
| `DeregisterSlot` | `google.protobuf.Int32Value` | `google.protobuf.Empty` |  |
| `SetSlotGreedy` | `google.protobuf.Int32Value` | `google.protobuf.Empty` |  |
| `SubmitSlotEvent` | `SlotEvent` | `StringMap` |  |
| `GetTaskList` | `google.protobuf.Int32Value` | `TaskItemList` |  |
| `SetTaskStarted` | `google.protobuf.Int64Value` | `google.protobuf.Empty` |  |
| `SetTaskFinished` | `TaskExecMessage` | `google.protobuf.Empty` |  |
| `UpdateBatchTask` | `BatchTaskStatus` | `google.protobuf.Empty` |  |
| `AddTaskList` | `TaskList` | `google.protobuf.Int64Value` |  |
| `StoreEvent` | `Event` | `google.protobuf.Empty` |  |

### 消息

| 消息 | 字段 | 类型 | 编号 | 说明 |
|---|---|---|---|---|
| `StringMap` | `values` | `map<string, string>` | 1 |  |
| `RegisterSlotResponse` | `slot_id` | `int32` | 1 |  |
| `RegisterSlotResponse` | `group_prefix` | `string` | 2 |  |
| `SetSlotInitializedResponse` | `group_prefix` | `string` | 1 |  |
| `TaskList` | `module_id` | `int32` | 1 |  |
| `TaskList` | `app_id` | `int32` | 2 |  |
| `TaskList` | `sink_module` | `string` | 3 |  |
| `TaskList` | `remote_cluster` | `string` | 4 |  |
| `TaskList` | `bodies` | `string (repeated)` | 10 |  |
| `TaskList` | `headers` | `map<string, string>` | 21 |  |
| `BatchTaskStatus` | `ids` | `int64 (repeated)` | 1 |  |
| `BatchTaskStatus` | `code` | `int32` | 2 |  |
| `TaskItem` | `id` | `int64` | 1 |  |
| `TaskItem` | `body` | `string` | 10 |  |
| `TaskItem` | `headers` | `map<string, string>` | 11 |  |
| `TaskItemList` | `tasks` | `TaskItem (repeated)` | 1 |  |
| `TaskExecMessage` | `slot` | `int32` | 1 |  |
| `TaskExecMessage` | `status_code` | `int32` | 2 |  |
| `TaskExecMessage` | `task_id` | `int64` | 3 |  |
| `TaskExecMessage` | `input_bytes` | `int64` | 4 |  |
| `TaskExecMessage` | `output_bytes` | `int64` | 5 |  |
| `TaskExecMessage` | `t0` | `google.protobuf.Timestamp` | 6 |  |
| `TaskExecMessage` | `t1` | `google.protobuf.Timestamp` | 7 |  |
| `TaskExecMessage` | `t2` | `google.protobuf.Timestamp` | 8 |  |
| `TaskExecMessage` | `t3` | `google.protobuf.Timestamp` | 9 |  |
| `TaskExecMessage` | `extras` | `string` | 10 |  |
| `TaskExecMessage` | `input_files` | `string (repeated)` | 11 |  |
| `TaskExecMessage` | `output_files` | `string (repeated)` | 12 |  |
| `TaskExecMessage` | `network_files` | `string (repeated)` | 13 |  |
| `TaskExecMessage` | `stdout` | `string` | 20 |  |
| `TaskExecMessage` | `stderr` | `string` | 21 |  |
| `TaskExecMessage` | `auxout` | `string` | 22 |  |
| `TaskExecMessage` | `ctxlog` | `string` | 23 |  |
| `SlotSpec` | `id` | `int32` | 1 |  |
| `SlotSpec` | `host` | `string` | 2 |  |
| `SlotSpec` | `ip_addr` | `string` | 3 |  |
| `SlotSpec` | `port` | `int32` | 4 |  |
| `SlotSpec` | `uname` | `string` | 5 |  |
| `SlotSpec` | `use_home_tmp` | `bool` | 6 |  |
| `SlotSpec` | `command_text` | `string` | 7 |  |
| `SlotList` | `slots` | `SlotSpec (repeated)` | 1 |  |
| `SlotEvent` | `slot_id` | `int32` | 1 |  |
| `SlotEvent` | `event_type` | `SlotEventType` | 2 |  |
| `SlotEvent` | `task_id` | `int64` | 11 |  |
| `SlotEvent` | `task_tm` | `google.protobuf.Timestamp` | 12 |  |
| `SlotEvent` | `details` | `map<string, string>` | 21 |  |
| `SlotEvent` | `extras` | `map<string, string>` | 22 |  |
| `Event` | `type` | `EventType` | 1 |  |
| `Event` | `obj_id` | `int32` | 2 |  |
| `Event` | `t` | `google.protobuf.Timestamp` | 3 |  |
| `Event` | `tag` | `string` | 4 |  |
| `Event` | `level` | `LogLevel` | 5 |  |
| `Event` | `code` | `int32` | 6 |  |
| `Event` | `txt` | `string` | 7 |  |
| `Event` | `json` | `google.protobuf.Struct` | 8 |  |

### 枚举

| 枚举 | 值 | 编号 | 说明 |
|---|---|---|---|
| `SlotEventType` | `SLOT_EVENT_UNKNOWN` | 0 |  |
| `SlotEventType` | `SLOT_CREATE` | 1 |  |
| `SlotEventType` | `SLOT_START` | 2 |  |
| `SlotEventType` | `SLOT_SETUP` | 3 |  |
| `SlotEventType` | `SLOT_TEARDOWN` | 4 |  |
| `SlotEventType` | `SLOT_CHECK` | 5 |  |
| `SlotEventType` | `SLOT_COMPLETE` | 6 |  |
| `SlotEventType` | `SLOT_ONGOING` | 10 |  |
| `SlotEventType` | `SLOT_IDLE` | 11 |  |
| `SlotEventType` | `SLOT_DEBUG` | 20 |  |
| `SlotEventType` | `SLOT_ERROR` | 21 |  |
| `Event.EventType` | `TASK` | 0 |  |
| `Event.EventType` | `SLOT` | 1 |  |
| `Event.EventType` | `MISC` | 2 |  |
| `Event.LogLevel` | `DEBUG` | 0 |  |
| `Event.LogLevel` | `INFO` | 1 |  |
| `Event.LogLevel` | `WARN` | 2 |  |
| `Event.LogLevel` | `ERROR` | 3 |  |
| `Event.LogLevel` | `FATAL` | 4 |  |


## coordination.proto

### 服务

#### `CoordinationService`

| 方法 | 请求 | 响应 | 说明 |
|---|---|---|---|
| `GetSemaphore` | `GetSemaphoreRequest` | `GetSemaphoreResponse` |  |
| `ListSemaphores` | `ListSemaphoresRequest` | `ListSemaphoresResponse` |  |
| `CreateSemaphore` | `CreateSemaphoreRequest` | `CreateSemaphoreResponse` |  |
| `UpdateSemaphore` | `UpdateSemaphoreRequest` | `UpdateSemaphoreResponse` |  |
| `DeleteSemaphore` | `DeleteSemaphoreRequest` | `google.protobuf.Empty` |  |
| `GetVariable` | `GetVariableRequest` | `GetVariableResponse` |  |
| `ListVariables` | `ListVariablesRequest` | `ListVariablesResponse` |  |
| `SetVariable` | `SetVariableRequest` | `SetVariableResponse` |  |
| `DeleteVariable` | `DeleteVariableRequest` | `google.protobuf.Empty` |  |
| `GetGlobal` | `GetGlobalRequest` | `GetGlobalResponse` |  |
| `ListGlobals` | `ListGlobalsRequest` | `ListGlobalsResponse` |  |
| `SetGlobal` | `SetGlobalRequest` | `google.protobuf.Empty` |  |
| `DeleteGlobal` | `DeleteGlobalRequest` | `google.protobuf.Empty` |  |
| `GetSemagroupMin` | `SemagroupRequest` | `SemagroupKeyValueResponse` |  |
| `GetSemagroupMax` | `SemagroupRequest` | `SemagroupKeyValueResponse` |  |
| `DiffSemagroupMin` | `SemagroupRequest` | `SemagroupDiffResponse` |  |
| `DiffSemagroupMax` | `SemagroupRequest` | `SemagroupDiffResponse` |  |
| `IncrementSemagroup` | `SemagroupRequest` | `SemagroupKeyValueResponse` |  |
| `DecrementSemagroup` | `SemagroupRequest` | `SemagroupKeyValueResponse` |  |
| `PushQueue` | `PushQueueRequest` | `PushQueueResponse` |  |
| `PullQueue` | `PullQueueRequest` | `PullQueueResponse` |  |
| `PeekQueue` | `PeekQueueRequest` | `PeekQueueResponse` |  |
| `ListQueues` | `ListQueuesRequest` | `ListQueuesResponse` |  |
| `PurgeQueue` | `PurgeQueueRequest` | `PurgeQueueResponse` |  |

### 消息

| 消息 | 字段 | 类型 | 编号 | 说明 |
|---|---|---|---|---|
| `GetSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `GetSemaphoreRequest` | `name` | `string` | 2 |  |
| `GetSemaphoreResponse` | `name` | `string` | 1 |  |
| `GetSemaphoreResponse` | `value` | `int32` | 2 |  |
| `GetSemaphoreResponse` | `value0` | `int32` | 3 |  |
| `ListSemaphoresRequest` | `app_id` | `int32` | 1 |  |
| `ListSemaphoresRequest` | `prefix` | `string` | 2 |  |
| `ListSemaphoresRequest` | `leaf_only` | `bool` | 3 |  |
| `ListSemaphoresRequest` | `vtask_id` | `int64` | 4 |  |
| `ListSemaphoresResponse` | `nodes` | `Node (repeated)` | 1 |  |
| `ListSemaphoresResponse.Node` | `name` | `string` | 1 |  |
| `ListSemaphoresResponse.Node` | `value` | `int32` | 2 |  |
| `ListSemaphoresResponse.Node` | `value0` | `int32` | 3 |  |
| `ListSemaphoresResponse.Node` | `is_leaf` | `bool` | 5 |  |
| `ListSemaphoresResponse.Node` | `child_count` | `int32` | 6 |  |
| `ListSemaphoresResponse.Node` | `vtask` | `int64` | 7 |  |
| `ListSemaphoresResponse.Node` | `t` | `google.protobuf.Timestamp` | 8 |  |
| `ListSemaphoresResponse.Node` | `id` | `int32` | 9 |  |
| `ListSemaphoresResponse.Node` | `app` | `int32` | 10 |  |
| `CreateSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `CreateSemaphoreRequest` | `name` | `string` | 2 |  |
| `CreateSemaphoreRequest` | `value` | `int32` | 3 |  |
| `CreateSemaphoreRequest` | `conflict_action` | `string` | 4 |  |
| `CreateSemaphoreResponse` | `name` | `string` | 1 |  |
| `CreateSemaphoreResponse` | `value` | `int32` | 2 |  |
| `UpdateSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `UpdateSemaphoreRequest` | `name` | `string` | 2 |  |
| `UpdateSemaphoreRequest` | `delta` | `int32` | 3 |  |
| `UpdateSemaphoreResponse` | `name` | `string` | 1 |  |
| `UpdateSemaphoreResponse` | `value` | `int32` | 2 |  |
| `DeleteSemaphoreRequest` | `app_id` | `int32` | 1 |  |
| `DeleteSemaphoreRequest` | `name` | `string` | 2 |  |
| `GetVariableRequest` | `app_id` | `int32` | 1 |  |
| `GetVariableRequest` | `name` | `string` | 2 |  |
| `GetVariableResponse` | `name` | `string` | 1 |  |
| `GetVariableResponse` | `value` | `string` | 2 |  |
| `ListVariablesRequest` | `app_id` | `int32` | 1 |  |
| `ListVariablesRequest` | `prefix` | `string` | 2 |  |
| `ListVariablesRequest` | `leaf_only` | `bool` | 3 |  |
| `ListVariablesRequest` | `vtask_id` | `int64` | 4 |  |
| `ListVariablesResponse` | `nodes` | `Node (repeated)` | 1 |  |
| `ListVariablesResponse.Node` | `name` | `string` | 1 |  |
| `ListVariablesResponse.Node` | `value` | `string` | 2 |  |
| `ListVariablesResponse.Node` | `is_leaf` | `bool` | 4 |  |
| `ListVariablesResponse.Node` | `child_count` | `int32` | 5 |  |
| `ListVariablesResponse.Node` | `vtask` | `int64` | 6 |  |
| `ListVariablesResponse.Node` | `t` | `google.protobuf.Timestamp` | 7 |  |
| `ListVariablesResponse.Node` | `id` | `int32` | 8 |  |
| `ListVariablesResponse.Node` | `app` | `int32` | 9 |  |
| `SetVariableRequest` | `app_id` | `int32` | 1 |  |
| `SetVariableRequest` | `name` | `string` | 2 |  |
| `SetVariableRequest` | `value` | `string` | 3 |  |
| `SetVariableResponse` | `name` | `string` | 1 |  |
| `SetVariableResponse` | `value` | `string` | 2 |  |
| `DeleteVariableRequest` | `app_id` | `int32` | 1 |  |
| `DeleteVariableRequest` | `name` | `string` | 2 |  |
| `GetGlobalRequest` | `name` | `string` | 1 |  |
| `GetGlobalResponse` | `name` | `string` | 1 |  |
| `GetGlobalResponse` | `value` | `string` | 2 |  |
| `KeyValue` | `key` | `string` | 1 |  |
| `KeyValue` | `value` | `string` | 2 |  |
| `KeyValue` | `t` | `google.protobuf.Timestamp` | 3 |  |
| `ListGlobalsRequest` | `prefix` | `string` | 1 |  |
| `ListGlobalsRequest` | `leaf_only` | `bool` | 2 |  |
| `ListGlobalsResponse` | `globals` | `KeyValue (repeated)` | 1 |  |
| `SetGlobalRequest` | `name` | `string` | 1 |  |
| `SetGlobalRequest` | `value` | `string` | 2 |  |
| `DeleteGlobalRequest` | `name` | `string` | 1 |  |
| `SemagroupRequest` | `app_id` | `int32` | 1 |  |
| `SemagroupRequest` | `sema_expr` | `string` | 2 |  |
| `SemagroupKeyValueResponse` | `name` | `string` | 1 |  |
| `SemagroupKeyValueResponse` | `value` | `int32` | 2 |  |
| `SemagroupDiffResponse` | `value` | `int32` | 1 |  |
| `QueueMessage` | `id` | `int64` | 1 |  |
| `QueueMessage` | `name` | `string` | 2 |  |
| `QueueMessage` | `priority` | `int32` | 3 |  |
| `QueueMessage` | `body` | `string` | 4 |  |
| `QueueMessage` | `headers` | `map<string, string>` | 5 |  |
| `QueueMessage` | `created_at` | `google.protobuf.Timestamp` | 6 |  |
| `PushQueueRequest` | `name` | `string` | 1 |  |
| `PushQueueRequest` | `body` | `string` | 2 |  |
| `PushQueueRequest` | `priority` | `int32` | 3 |  |
| `PushQueueRequest` | `headers` | `map<string, string>` | 4 |  |
| `PushQueueRequest` | `cluster` | `string` | 5 |  |
| `PushQueueResponse` | `id` | `int64` | 1 |  |
| `PullQueueRequest` | `name` | `string` | 1 |  |
| `PullQueueRequest` | `count` | `int32` | 2 |  |
| `PullQueueRequest` | `cluster` | `string` | 3 |  |
| `PullQueueResponse` | `messages` | `QueueMessage (repeated)` | 1 |  |
| `PeekQueueRequest` | `name` | `string` | 1 |  |
| `PeekQueueRequest` | `count` | `int32` | 2 |  |
| `PeekQueueRequest` | `cluster` | `string` | 3 |  |
| `PeekQueueResponse` | `messages` | `QueueMessage (repeated)` | 1 |  |
| `ListQueuesRequest` | `prefix` | `string` | 1 |  |
| `ListQueuesRequest` | `cluster` | `string` | 2 |  |
| `ListQueuesResponse` | `queues` | `QueueInfo (repeated)` | 1 |  |
| `ListQueuesResponse.QueueInfo` | `name` | `string` | 1 |  |
| `ListQueuesResponse.QueueInfo` | `count` | `int64` | 2 |  |
| `ListQueuesResponse.QueueInfo` | `oldest_at` | `google.protobuf.Timestamp` | 3 |  |
| `PurgeQueueRequest` | `name` | `string` | 1 |  |
| `PurgeQueueRequest` | `cluster` | `string` | 2 |  |
| `PurgeQueueResponse` | `deleted` | `int64` | 1 |  |


## dashboard.proto

### 服务

#### `DashboardService`

| 方法 | 请求 | 响应 | 说明 |
|---|---|---|---|
| `GetSystemStats` | `google.protobuf.Empty` | `SystemStats` |  |
| `GetAppStats` | `GetAppStatsRequest` | `AppStats` |  |
| `GetModuleStats` | `GetModuleStatsRequest` | `ModuleStats` |  |
| `ListModuleStats` | `ListModuleStatsRequest` | `ListModuleStatsResponse` |  |

### 消息

| 消息 | 字段 | 类型 | 编号 | 说明 |
|---|---|---|---|---|
| `SystemStats` | `total_clusters` | `int32` | 1 |  |
| `SystemStats` | `clusters_on` | `int32` | 2 |  |
| `SystemStats` | `clusters_error` | `int32` | 3 |  |
| `SystemStats` | `total_hosts` | `int32` | 4 |  |
| `SystemStats` | `hosts_on` | `int32` | 5 |  |
| `SystemStats` | `total_slots` | `int32` | 6 |  |
| `SystemStats` | `slots_on` | `int32` | 7 |  |
| `SystemStats` | `running_apps` | `int32` | 8 |  |
| `SystemStats` | `active_tasks` | `int32` | 9 |  |
| `GetAppStatsRequest` | `app_id` | `int32` | 1 |  |
| `AppStats` | `total_tasks` | `int32` | 1 |  |
| `AppStats` | `finished_tasks_24h` | `int32` | 2 |  |
| `AppStats` | `error_tasks_24h` | `int32` | 3 |  |
| `AppStats` | `active_slots` | `int32` | 4 |  |
| `AppStats` | `total_slots` | `int32` | 5 |  |
| `AppStats` | `total_io_bytes` | `int64` | 6 |  |
| `AppStats` | `avg_latency_ms` | `double` | 7 |  |
| `AppStats` | `task_trend` | `TaskTrendPoint (repeated)` | 8 |  |
| `AppStats` | `module_latency` | `ModuleLatency (repeated)` | 9 |  |
| `AppStats` | `slot_util` | `SlotUtilization (repeated)` | 10 |  |
| `AppStats.TaskTrendPoint` | `ts` | `google.protobuf.Timestamp` | 1 |  |
| `AppStats.TaskTrendPoint` | `finished` | `int32` | 2 |  |
| `AppStats.TaskTrendPoint` | `error` | `int32` | 3 |  |
| `AppStats.ModuleLatency` | `module_name` | `string` | 1 |  |
| `AppStats.ModuleLatency` | `avg_ms` | `double` | 2 |  |
| `AppStats.ModuleLatency` | `p50_ms` | `double` | 3 |  |
| `AppStats.ModuleLatency` | `p95_ms` | `double` | 4 |  |
| `AppStats.ModuleLatency` | `p99_ms` | `double` | 5 |  |
| `AppStats.SlotUtilization` | `slot_name` | `string` | 1 |  |
| `AppStats.SlotUtilization` | `hostname` | `string` | 2 |  |
| `AppStats.SlotUtilization` | `util_pct` | `double` | 3 |  |
| `GetModuleStatsRequest` | `app_id` | `int32` | 1 |  |
| `GetModuleStatsRequest` | `module_name` | `string` | 2 |  |
| `ModuleStats` | `module_name` | `string` | 1 |  |
| `ModuleStats` | `total_tasks` | `int32` | 2 |  |
| `ModuleStats` | `error_count` | `int32` | 3 |  |
| `ModuleStats` | `task_trend` | `AppStats.TaskTrendPoint (repeated)` | 4 |  |
| `ModuleStats` | `error_distribution` | `ExitCodeCount (repeated)` | 5 |  |
| `ModuleStats` | `slot_util` | `AppStats.SlotUtilization (repeated)` | 6 |  |
| `ModuleStats` | `input_bytes` | `int64` | 7 |  |
| `ModuleStats` | `output_bytes` | `int64` | 8 |  |
| `ModuleStats` | `network_bytes` | `int64` | 9 |  |
| `ModuleStats.ExitCodeCount` | `exit_code` | `int32` | 1 |  |
| `ModuleStats.ExitCodeCount` | `count` | `int32` | 2 |  |
| `ListModuleStatsRequest` | `app_id` | `int32` | 1 |  |
| `ListModuleStatsRequest` | `module_name` | `string` | 2 |  |
| `ListModuleStatsRequest` | `include_main_router` | `bool` | 3 |  |
| `ListModuleStatsResponse` | `items` | `ModuleStatItem (repeated)` | 1 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `module_id` | `int32` | 1 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `module_name` | `string` | 2 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `num_tasks` | `int64` | 3 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `num_execs` | `int64` | 4 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `ntasks_ok` | `int64` | 5 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `ntasks_ready` | `int64` | 6 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `ntasks_queued` | `int64` | 7 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `ntasks_running` | `int64` | 8 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `ntasks_err` | `int64` | 9 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `sum_ibytes` | `int64` | 10 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `sum_obytes` | `int64` | 11 |  |
| `ListModuleStatsResponse.ModuleStatItem` | `total_seconds` | `double` | 12 |  |


## error_detail.proto

### 消息

| 消息 | 字段 | 类型 | 编号 | 说明 |
|---|---|---|---|---|
| `ErrorDetail` | `message` | `string` | 1 |  |
| `ErrorDetail` | `code` | `int32` | 2 |  |
| `ErrorDetail` | `location` | `string` | 3 |  |
| `ErrorDetail` | `timestamp` | `google.protobuf.Timestamp` | 4 |  |
| `ErrorDetail` | `context` | `map<string, string>` | 5 |  |
| `ErrorDetail` | `cause` | `ErrorDetail` | 6 |  |


## infra.proto

### 服务

#### `InfraService`

| 方法 | 请求 | 响应 | 说明 |
|---|---|---|---|
| `ListClusters` | `ListClustersRequest` | `ListClustersResponse` |  |
| `GetCluster` | `GetClusterRequest` | `GetClusterResponse` |  |
| `CreateCluster` | `CreateClusterRequest` | `CreateClusterResponse` |  |
| `DeleteCluster` | `DeleteClusterRequest` | `google.protobuf.Empty` |  |
| `SetClusterStatus` | `SetClusterStatusRequest` | `google.protobuf.Empty` |  |
| `ListHosts` | `ListHostsRequest` | `ListHostsResponse` |  |
| `GetHost` | `GetHostRequest` | `GetHostResponse` |  |
| `CreateHost` | `CreateHostRequest` | `CreateHostResponse` |  |
| `DeleteHost` | `DeleteHostRequest` | `google.protobuf.Empty` |  |
| `SetHostStatus` | `SetHostStatusRequest` | `google.protobuf.Empty` |  |
| `RenewHost` | `RenewHostRequest` | `google.protobuf.Empty` |  |
| `ListSlots` | `ListSlotsRequest` | `ListSlotsResponse` |  |
| `GetSlot` | `GetSlotRequest` | `GetSlotResponse` |  |
| `AddSlots` | `AddSlotsRequest` | `AddSlotsResponse` |  |
| `RemoveSlots` | `RemoveSlotsRequest` | `RemoveSlotsResponse` |  |
| `SetSlotStatus` | `SetSlotStatusRequest` | `google.protobuf.Empty` |  |
| `GetSlotParameter` | `GetSlotParameterRequest` | `GetSlotParameterResponse` |  |
| `SetSlotParameter` | `SetSlotParameterRequest` | `google.protobuf.Empty` |  |
| `GetClusterGrpcServer` | `GetClusterGrpcServerRequest` | `GetClusterGrpcServerResponse` |  |
| `GetClusterByName` | `GetClusterByNameRequest` | `GetClusterResponse` |  |
| `GetClusterParameter` | `GetClusterParameterRequest` | `GetClusterParameterResponse` |  |
| `GetHostByHostname` | `GetHostByHostnameRequest` | `GetHostResponse` |  |
| `AllocateCluster` | `AllocateClusterRequest` | `AllocateClusterResponse` |  |
| `ReleaseCluster` | `ReleaseClusterRequest` | `ReleaseClusterResponse` |  |
| `RenewCluster` | `RenewClusterRequest` | `RenewClusterResponse` |  |

### 消息

| 消息 | 字段 | 类型 | 编号 | 说明 |
|---|---|---|---|---|
| `ListClustersRequest` | `status` | `string` | 1 |  |
| `ListClustersResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListClustersResponse.Item` | `id` | `int32` | 1 |  |
| `ListClustersResponse.Item` | `name` | `string` | 2 |  |
| `ListClustersResponse.Item` | `type` | `string` | 3 |  |
| `ListClustersResponse.Item` | `status` | `string` | 4 |  |
| `ListClustersResponse.Item` | `host_count` | `int32` | 5 |  |
| `ListClustersResponse.Item` | `slot_total` | `int32` | 6 |  |
| `ListClustersResponse.Item` | `slot_used` | `int32` | 7 |  |
| `ListClustersResponse.Item` | `grpc_server` | `string` | 8 |  |
| `GetClusterRequest` | `name` | `string` | 2 |  |
| `GetClusterResponse` | `id` | `int32` | 1 |  |
| `GetClusterResponse` | `name` | `string` | 2 |  |
| `GetClusterResponse` | `type` | `string` | 3 |  |
| `GetClusterResponse` | `status` | `string` | 4 |  |
| `GetClusterResponse` | `grpc_server` | `string` | 5 |  |
| `GetClusterResponse` | `parameters` | `string` | 6 |  |
| `GetClusterResponse` | `created_at` | `google.protobuf.Timestamp` | 7 |  |
| `GetClusterResponse` | `hosts` | `string (repeated)` | 10 |  |
| `GetClusterResponse` | `apps` | `string (repeated)` | 11 |  |
| `CreateClusterRequest` | `yaml_content` | `string` | 1 |  |
| `CreateClusterResponse` | `name` | `string` | 2 |  |
| `DeleteClusterRequest` | `name` | `string` | 2 |  |
| `SetClusterStatusRequest` | `status` | `string` | 2 |  |
| `SetClusterStatusRequest` | `name` | `string` | 3 |  |
| `ListHostsRequest` | `status` | `string` | 2 |  |
| `ListHostsRequest` | `cluster_name` | `string` | 3 |  |
| `ListHostsResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListHostsResponse.Item` | `id` | `int32` | 1 |  |
| `ListHostsResponse.Item` | `hostname` | `string` | 2 |  |
| `ListHostsResponse.Item` | `ip_addr` | `string` | 3 |  |
| `ListHostsResponse.Item` | `cluster_name` | `string` | 5 |  |
| `ListHostsResponse.Item` | `status` | `string` | 6 |  |
| `ListHostsResponse.Item` | `slot_total` | `int32` | 7 |  |
| `ListHostsResponse.Item` | `slot_used` | `int32` | 8 |  |
| `ListHostsResponse.Item` | `created_at` | `google.protobuf.Timestamp` | 9 |  |
| `ListHostsResponse.Item` | `port` | `int32` | 10 |  |
| `ListHostsResponse.Item` | `uname` | `string` | 11 |  |
| `GetHostRequest` | `host_id` | `int32` | 1 |  |
| `GetHostResponse` | `id` | `int32` | 1 |  |
| `GetHostResponse` | `hostname` | `string` | 2 |  |
| `GetHostResponse` | `ip_addr` | `string` | 3 |  |
| `GetHostResponse` | `cluster_name` | `string` | 5 |  |
| `GetHostResponse` | `status` | `string` | 6 |  |
| `GetHostResponse` | `labels` | `string` | 7 |  |
| `GetHostResponse` | `created_at` | `google.protobuf.Timestamp` | 8 |  |
| `CreateHostRequest` | `hostname` | `string` | 1 |  |
| `CreateHostRequest` | `ip_addr` | `string` | 2 |  |
| `CreateHostRequest` | `uname` | `string` | 4 |  |
| `CreateHostRequest` | `port` | `int32` | 5 |  |
| `CreateHostRequest` | `cluster_name` | `string` | 6 |  |
| `CreateHostResponse` | `host_id` | `int32` | 1 |  |
| `DeleteHostRequest` | `host_id` | `int32` | 1 |  |
| `SetHostStatusRequest` | `host_id` | `int32` | 1 |  |
| `SetHostStatusRequest` | `status` | `string` | 2 |  |
| `RenewHostRequest` | `host_id` | `int32` | 1 |  |
| `ListSlotsRequest` | `app_id` | `int32` | 1 |  |
| `ListSlotsRequest` | `host_id` | `int32` | 2 |  |
| `ListSlotsRequest` | `module_id` | `int32` | 3 |  |
| `ListSlotsRequest` | `status` | `string` | 4 |  |
| `ListSlotsRequest` | `hostname` | `string` | 5 |  |
| `ListSlotsResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListSlotsResponse.Item` | `id` | `int32` | 1 |  |
| `ListSlotsResponse.Item` | `slot_name` | `string` | 2 |  |
| `ListSlotsResponse.Item` | `host_id` | `int32` | 3 |  |
| `ListSlotsResponse.Item` | `hostname` | `string` | 4 |  |
| `ListSlotsResponse.Item` | `status` | `string` | 6 |  |
| `ListSlotsResponse.Item` | `module_id` | `int32` | 7 |  |
| `ListSlotsResponse.Item` | `module_name` | `string` | 8 |  |
| `ListSlotsResponse.Item` | `started_at` | `google.protobuf.Timestamp` | 9 |  |
| `ListSlotsResponse.Item` | `seq` | `string` | 10 |  |
| `ListSlotsResponse.Item` | `ip_addr` | `string` | 11 |  |
| `ListSlotsResponse.Item` | `cluster_name` | `string` | 12 |  |
| `ListSlotsResponse.Item` | `heartbeat_at` | `google.protobuf.Timestamp` | 13 |  |
| `GetSlotRequest` | `slot_id` | `int32` | 1 |  |
| `GetSlotResponse` | `id` | `int32` | 1 |  |
| `GetSlotResponse` | `host_id` | `int32` | 3 |  |
| `GetSlotResponse` | `hostname` | `string` | 4 |  |
| `GetSlotResponse` | `ip_addr` | `string` | 5 |  |
| `GetSlotResponse` | `status` | `string` | 7 |  |
| `GetSlotResponse` | `module_id` | `int32` | 8 |  |
| `GetSlotResponse` | `module_name` | `string` | 9 |  |
| `GetSlotResponse` | `command_text` | `string` | 10 |  |
| `GetSlotResponse` | `started_at` | `google.protobuf.Timestamp` | 11 |  |
| `GetSlotResponse` | `heartbeat_at` | `google.protobuf.Timestamp` | 12 |  |
| `GetSlotResponse` | `seq` | `string` | 13 |  |
| `GetSlotResponse` | `parameters` | `string` | 14 |  |
| `GetSlotResponse` | `cluster_name` | `string` | 15 |  |
| `AddSlotsRequest` | `host_id` | `int32` | 1 |  |
| `AddSlotsRequest` | `module_id` | `int32` | 2 |  |
| `AddSlotsRequest` | `count` | `int32` | 3 |  |
| `AddSlotsRequest` | `slot_regex` | `string` | 4 |  |
| `AddSlotsResponse` | `added` | `int32` | 1 |  |
| `AddSlotsResponse` | `slot_names` | `string (repeated)` | 2 |  |
| `RemoveSlotsRequest` | `slot_names` | `string (repeated)` | 1 |  |
| `RemoveSlotsResponse` | `removed` | `int32` | 1 |  |
| `SetSlotStatusRequest` | `slot_id` | `int32` | 1 |  |
| `SetSlotStatusRequest` | `status` | `string` | 2 |  |
| `GetSlotParameterRequest` | `slot_id` | `int32` | 1 |  |
| `GetSlotParameterRequest` | `parameter` | `string` | 2 |  |
| `GetSlotParameterResponse` | `value` | `string` | 1 |  |
| `SetSlotParameterRequest` | `slot_id` | `int32` | 1 |  |
| `SetSlotParameterRequest` | `parameter` | `string` | 2 |  |
| `SetSlotParameterRequest` | `value` | `string` | 3 |  |
| `GetClusterGrpcServerRequest` | `cluster_name` | `string` | 1 |  |
| `GetClusterGrpcServerResponse` | `grpc_server` | `string` | 1 |  |
| `GetClusterByNameRequest` | `name` | `string` | 1 |  |
| `GetClusterParameterRequest` | `cluster_name` | `string` | 1 |  |
| `GetClusterParameterRequest` | `parameter` | `string` | 2 |  |
| `GetClusterParameterResponse` | `value` | `string` | 1 |  |
| `GetHostByHostnameRequest` | `hostname` | `string` | 1 |  |
| `AllocateClusterRequest` | `cluster_name` | `string` | 1 |  |
| `AllocateClusterRequest` | `group_spec` | `string` | 2 |  |
| `AllocateClusterRequest` | `partition` | `string` | 3 |  |
| `AllocateClusterRequest` | `account` | `string` | 4 |  |
| `AllocateClusterRequest` | `time_limit` | `string` | 5 |  |
| `AllocateClusterRequest` | `cpu_only` | `string` | 6 |  |
| `AllocateClusterRequest` | `min_tmp_gb` | `int32` | 7 |  |
| `AllocateClusterResponse` | `num_groups` | `int32` | 1 |  |
| `AllocateClusterResponse` | `group_names` | `string (repeated)` | 2 |  |
| `ReleaseClusterRequest` | `cluster_name` | `string` | 1 |  |
| `ReleaseClusterRequest` | `group_name` | `string` | 2 |  |
| `ReleaseClusterRequest` | `all` | `bool` | 3 |  |
| `ReleaseClusterResponse` | `released_hosts` | `int32` | 1 |  |
| `RenewClusterRequest` | `cluster_name` | `string` | 1 |  |
| `RenewClusterRequest` | `group_name` | `string` | 2 |  |
| `RenewClusterRequest` | `all` | `bool` | 3 |  |
| `RenewClusterResponse` | `renewed_groups` | `int32` | 1 |  |


## user.proto

### 服务

#### `UserService`

| 方法 | 请求 | 响应 | 说明 |
|---|---|---|---|
| `ListUsers` | `ListUsersRequest` | `ListUsersResponse` |  |
| `GetUser` | `GetUserRequest` | `GetUserResponse` |  |
| `CreateUser` | `CreateUserRequest` | `CreateUserResponse` |  |
| `UpdateUser` | `UpdateUserRequest` | `google.protobuf.Empty` |  |
| `DisableUser` | `DisableUserRequest` | `google.protobuf.Empty` |  |
| `BindRole` | `BindRoleRequest` | `google.protobuf.Empty` |  |
| `UnbindRole` | `UnbindRoleRequest` | `google.protobuf.Empty` |  |
| `Login` | `LoginRequest` | `LoginResponse` |  |
| `Logout` | `LogoutRequest` | `google.protobuf.Empty` |  |

### 消息

| 消息 | 字段 | 类型 | 编号 | 说明 |
|---|---|---|---|---|
| `ListUsersRequest` | — | — | — | 空消息 |
| `ListUsersResponse` | `items` | `Item (repeated)` | 1 |  |
| `ListUsersResponse.Item` | `id` | `int32` | 1 |  |
| `ListUsersResponse.Item` | `username` | `string` | 2 |  |
| `ListUsersResponse.Item` | `roles` | `string (repeated)` | 3 |  |
| `ListUsersResponse.Item` | `scope` | `string` | 4 |  |
| `ListUsersResponse.Item` | `active` | `bool` | 5 |  |
| `ListUsersResponse.Item` | `created_at` | `google.protobuf.Timestamp` | 6 |  |
| `GetUserRequest` | `user_id` | `int32` | 1 |  |
| `GetUserResponse` | `id` | `int32` | 1 |  |
| `GetUserResponse` | `username` | `string` | 2 |  |
| `GetUserResponse` | `roles` | `string (repeated)` | 3 |  |
| `GetUserResponse` | `scope` | `string` | 4 |  |
| `GetUserResponse` | `active` | `bool` | 5 |  |
| `GetUserResponse` | `created_at` | `google.protobuf.Timestamp` | 6 |  |
| `CreateUserRequest` | `username` | `string` | 1 |  |
| `CreateUserRequest` | `password` | `string` | 2 |  |
| `CreateUserRequest` | `roles` | `string (repeated)` | 3 |  |
| `CreateUserRequest` | `scope` | `string` | 4 |  |
| `CreateUserResponse` | `user_id` | `int32` | 1 |  |
| `UpdateUserRequest` | `user_id` | `int32` | 1 |  |
| `UpdateUserRequest` | `username` | `string` | 2 |  |
| `UpdateUserRequest` | `password` | `string` | 3 |  |
| `UpdateUserRequest` | `roles` | `string (repeated)` | 4 |  |
| `UpdateUserRequest` | `scope` | `string` | 5 |  |
| `DisableUserRequest` | `user_id` | `int32` | 1 |  |
| `BindRoleRequest` | `user_id` | `int32` | 1 |  |
| `BindRoleRequest` | `role` | `string` | 2 |  |
| `BindRoleRequest` | `scope` | `string` | 3 |  |
| `UnbindRoleRequest` | `user_id` | `int32` | 1 |  |
| `UnbindRoleRequest` | `role` | `string` | 2 |  |
| `UnbindRoleRequest` | `scope` | `string` | 3 |  |
| `LoginRequest` | `username` | `string` | 1 |  |
| `LoginRequest` | `password` | `string` | 2 |  |
| `LoginResponse` | `token` | `string` | 1 |  |
| `LoginResponse` | `user_id` | `int32` | 2 |  |
| `LoginResponse` | `username` | `string` | 3 |  |
| `LoginResponse` | `roles` | `string (repeated)` | 4 |  |
| `LoginResponse` | `expires_at` | `google.protobuf.Timestamp` | 5 |  |
| `LogoutRequest` | `token` | `string` | 1 |  |
