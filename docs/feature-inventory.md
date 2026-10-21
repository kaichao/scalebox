# Scalebox 特性覆盖矩阵（阶段 0 盘点）

> 状态：阶段 0 产物（2026-10-17），待补充实战来源后定稿。
> 上位方案：`docs/ai-agent-support-plan.md`。本文件用于确定 feature 建设优先级，不替代任何设计文档。

## 0. 方法与输入

**主干**：`docs/cn/source/appendix/2_parameter_reference.md`（7 类参数表）。

**叠加四路输入**：

1. `features/` 现有 9 个目录与 `features/README.md`
2. `docs/cn/source/`（Sphinx 面向人文档）
3. `go-scalebox/docs/`（设计文档）与源码（proto、CLI、app-schema.json、SQL）
4. `app-fxcorr`、`app-mwa` 的实战记录（CLAUDE.md、docs/、模块 README）

**图例**

| 缺口 | 含义 |
|---|---|
| 高 | agent 高频使用 + 文档不足/过期/有坑 + 无可靠测试 |
| 中 | 有文档但缺少测试，或文档需重写，或使用频率一般 |
| 低 | 已有 feature 测试且文档基本准确 |
| 未确认 | 参数表标有"暂未使用/拟修改/？"，或源码中未找到消费点 |

| 去向 | 含义 |
|---|---|
| feature | 新建/统一 feature 目录（README + feature.yaml + test.sh） |
| reference | 只需进 `docs/reference/`（表格式参考） |
| invariant | 进 `docs/agents/10-invariants.md`（规则/坑） |
| docs | 留在 Sphinx（原理性内容） |
| 已有 | 现有覆盖足够，补充索引即可 |

## 1. 初步关键发现

1. **features 索引与目录不一致**：`features/README.md` 索引 5 项中 3 项指向不存在的目录（`timeout-gen`、`task-perspective`、`hello-scalebox-singularity`）；9 个实际目录中 7 个未正确登记（`timeout` 被写成 `timeout-gen`）。
2. **features 参数名错误**：`features/check_test/` 的 README 与 app.yaml 使用 `space_free_gb`，但源码与参数表均为 `free_space_gb`（全仓库搜索 `space_free_gb` 为 0 次）。
3. **examples 命令过时**：`examples/README.md` 使用 `app create` / `app list`；`features/node_progress_gap/README.md` 使用 `scalebox app run`。现行 CLI 为 `scalebox run` / `scalebox app list`（`cmd/cli/run/command.go` 是 `run` 的唯一实现）。
4. **参数生效层不同**：`slot_recoverable`、`slot_max_retries`、`slot_timeout_minutes`、`task_global_timeout_scale` 不在 Go 代码字符串中，由 SQL/pgcron 或 JSON Schema 侧消费（`build/database/initdb.d/*.sql`、`cmd/webui/frontend/src/schema/app-schema.json`），文档需注明"生效层"。
5. **`scalebox run` 已有隐式探测**：`--app-file` 为空时自动使用当前目录的 `app.yaml`（`cmd/cli/run/command.go`）。因此 feature 目录若用 `feature.yaml` 需显式 `--app-file=feature.yaml`；不建议再增加缺省探测。
6. **A 档生成源已具备**：`--help`（CLI）、`cmd/webui/frontend/src/schema/app-schema.json`（app 模板字段）、`proto/*.proto`（RPC）、`appendix/6_exit_code_spec.md`。
7. **features 形态不齐**：8/9 有 README（`retry_test` 无）；均无统一 `test.sh` 与断言；`agent-interface-files` 无 app.yaml。
8. **设计文档质量高但状态未回填**：`go-scalebox/docs/` 26 篇设计文档集中了平台不变式（最佳：`vtask-design`、`cli-grpc-architecture`、`slot-status`、`task-headers`、`env-loading`），但多篇有过时/矛盾——`grpc-api.md` 的 RPC 计数有 97/102/104 三种口径；slot 状态 `CLOSED` 已移除但仍出现在时序图；`architecture-analysis.md` 仍写 WebUI 为 Go+htmx（实为 React SPA）；`webui-vscode-plan.md` 同时有 85/102 与 90 端点两种说法。
9. **app-fxcorr 沉淀 38 条平台机制实证**：覆盖任务流转/身份幂等/vtask 流控/slot 落点/环境生命周期/观测运维 6 类；Top 10 高频坑（NODES 正则、僵尸节点、`uniq_task` 112、未定义变量字面量、app 快照、wait-queue 解锁、WQ_VTS、信号量 API 混用、任务级 `free_space_gb` 死锁、stdin 任务列表）是 `docs/agents/10-invariants.md` 的第一批直接来源。完整清单见附录 A。
10. **Sphinx 文档与 CLI/仓库结构大面积不一致**：命令不存在（`cluster status`、`app status`、`task monitor`、`app metrics`、`node metrics` 等）、路径错误（`build/` 应为 `runtime/`）、编号错乱与空章节（appendix/3、developer/7）、状态码表自相矛盾、`space_free_gb` vs `free_space_gb`、`vtask_size_sema_copy` 已废弃仍列为有效。另有 17 个平台已有但文档缺失的主题（wait-queue、动态集群/cluster-admin、channel、fs、event、pod 生命周期、slot_options 场景、task_batch_size、retry_rules 完整语义、接口文件格式、I/O 计量、环境变量优先级权威表、表结构、REST 端点、DAG、singularity 指南、监控指标）。详见附录 B。
11. **app-mwa 沉淀 29 条平台机制实证**：覆盖路由装配/vtask 管道/信号量流控/任务头契约/存储虚拟化/运维排错 6 类；与附录 A 有重叠也有补充，且存在 1 处**条件性冲突**——`CONFLICT_ACTION=OVERWRITE` 在 app-mwa 用于人工重投，而 app-fxcorr 明确"router 重建下游任务时禁用"（会撞 112）。蒸馏时必须按场景拆分，不能合并成一条规则。详见附录 C。

## 2. 覆盖矩阵

### 2.1 App 级参数

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| `initial_status`（RUNNING/INITIAL/PAUSED） | 无 | appendix/2 §2.1.1 | 待补 | 中 | reference |
| `main_router` | examples/*、各应用 | appendix/2、developer/3 | app-fxcorr（主路由形态） | 低 | 已有 + reference |
| `default_idle_polls` / `default_sleep_count` | 无 | appendix/2 | 待补 | 低 | reference |
| `is_cluster_admin` | 无 | appendix/2（说明空白） | 待补 | 中 | reference |
| `global_directories` | 无 | appendix/2 | 待补 | 低 | reference |
| `slot_group`（多节点 slot 自动创建） | 无 | appendix/2 | app-fxcorr（NODES/slot 配置） | 中 | feature |

### 2.2 module-arguments（运行/资源）

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| `code_path` | 全部 features/examples | appendix/2 | app-fxcorr modules | 低 | reference |
| `grpc_server` | 无 | appendix/2 | 测试环境 | 低 | reference |
| `local_ip_index` | 无 | appendix/2（"agent 中暂未使用？"） | 待补 | 未确认 | 源码核实后定 |
| `task_max_seconds`（超时 124 / 全局 123） | features/timeout | appendix/2、appendix/6 | 待补 | 中（缺统一模板与断言） | feature |
| `task_min_seconds`（GREEDY 判定） | 无 | appendix/2 | 待补 | 中 | feature |
| `poll_interval_seconds` | 无 | appendix/2 | 待补 | 低 | reference |
| `max_idle_polls` | features/side-car | appendix/2 | 待补 | 低 | feature 统一 |
| `dir_quota_gb` | features/check_test | appendix/2 | 待补 | 中（示例参数名需修） | feature |
| `free_space_gb` | features/check_test（写成 `space_free_gb`） | appendix/2 | app-fxcorr（阶梯流控 `((@seq+5))`；挂错模块死锁） | 高 | invariant + feature |
| `heartbeat_seconds` | 无 | appendix/2 | 待补 | 低 | reference |
| `output_text_size` / `text_trunc_mode` | 无 | appendix/2 | app-fxcorr（查 stderr） | 低 | reference |
| `timezone_mode` | 无 | appendix/2（说明缺失） | 待补 | 中 | reference |
| `slot_options`（always_running / reserved_on_exit / tmpfs_workdir / disable_local_mapping / disable_data_mapping / slot_on_head） | side-car、pod、多处 | appendix/2 | app-fxcorr（tmpfs A/B） | 中 | feature + reference |
| `environments`（环境注入） | 多处 | developer/2（待确认）、go-scalebox `docs/env-loading.md` | app-fxcorr（创建时快照；同名变量优先级） | 高 | invariant + reference |

### 2.3 module-parameters（行为/路由/流控）

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| `key_group_regex` / `key_group_index` | 无 | appendix/2（"改名为 group_regex?"） | app-fxcorr（分组） | 中 | reference（源码核实用名） |
| `task_dist_mode`（HOST-BOUND/SLOT-BOUND/DEFAULT） | features/pod、examples/vtask、app-mwa | appendix/2、developer/9 | 多应用 | 高 | feature + invariant |
| `start_task` | retry_test、timeout、examples | appendix/2 | 待补 | 低 | feature 统一 |
| `initial_task_status` / `initial_slot_status` | 无 | appendix/2 | 待补 | 低 | reference |
| `retry_rules` | features/retry_test（无 README） | appendix/2 | app-fxcorr（v1 不配 retry_rules） | 中 | feature 补齐 |
| `slot_recoverable` / `slot_max_retries` / `slot_timeout_minutes` | 无 | appendix/2（SQL/schema 侧） | 待补 | 中 | reference（注明生效层） |
| `slot_max_tasks`（slot 生命周期） | 无 | appendix/2 | 待补 | 中 | feature |
| `task_global_timeout_scale`（拟改名） | 无 | appendix/2 | 待补 | 中 | reference |
| `router_index`（多主路由） | 无 | appendix/2 | app-fxcorr（多主路由?） | 中 | feature |
| `pod_id`（节点本地计算） | features/pod | appendix/2 | app-fxcorr（route=1/2） | 中（feature 太简） | feature 统一 |
| `task_dedup_cache_ttl_minutes` / `repeatable` | 无 | appendix/2 | 待补 | 中 | feature + invariant |
| `task_batch_size` | features/task_batch_size（README 三行） | appendix/2 | app-fxcorr（v9 定案 T=8） | 中 | feature 统一 + invariant |
| `visible` | 无 | appendix/2 | 待补 | 低 | reference |
| `task_id_in_headers` | 无 | appendix/2 | 待补 | 低 | reference |
| `node_progress_gap` | features/node_progress_gap（无断言） | appendix/2 | 待补 | 中 | feature 补齐 |
| `vtask_role` / `vtask_size` / `vtask_size_sema_copy` | examples/vtask、app-fxcorr | appendix/2、developer/9、examples/vtask | 大量 | 低（sema_copy 已废弃需标注） | feature + reference |

### 2.4 task-headers 与传播

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| `to_ip` / `to_host` / `from_ip` / `from_host` / `from_ip_last` / `from_host` / `from_module_last` / `from_cluster` | 多应用间接覆盖 | appendix/2 §2.1.4、go-scalebox `docs/task-headers.md` | app-fxcorr（路由） | 中 | reference |
| `to_slot` / `to_slot_index` | examples/vtask | appendix/2 | 待补 | 中 | reference |
| `repeatable` / `conflict_action` | 无 | appendix/2 | 待补 | 中 | feature |
| `slot_last_task` | 无 | appendix/2 | 待补 | 中 | feature |
| `_vtask_attr_*`（应用自定义属性） | examples/vtask | appendix/2、examples/vtask | 待补 | 中 | feature + reference |
| `_` 前缀 header 自动传播规则 | 无 | appendix/2（一句） | app-fxcorr（多处实证） | 高 | invariant |
| headers 类型约定（全 string、数组逗号分隔） | 无 | go-scalebox `docs/task-headers.md` | app-fxcorr（`_vtask_id` 等） | 高 | invariant + reference |

### 2.5 接口文件（agent-interface-files）

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| `task-exec.yaml` / `extra-attrs.yaml` | features/agent-interface-files（无 app.yaml） | features README（空白小节） | 待补 | 高 | feature 补齐 |
| `sink-tasks.txt` | agent-interface-files | features README | 待补 | 高 | feature 补齐 |
| `timestamps.txt` | agent-interface-files | features README | app-fxcorr（三档计时） | 高 | feature + reference |
| `input-files.txt` / `output-files.txt` / `removed-files.txt` | agent-interface-files | features README | 待补 | 高 | feature 补齐 |
| `auxout.txt` | agent-interface-files | features README | 待补 | 中 | feature |

### 2.6 VTask、信号量与流控

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| vtask 三种模式（DEFAULT/HOST-BOUND/GROUP-BOUND） | examples/vtask | developer/9、examples/vtask、go-scalebox `docs/vtask-design.md` | app-fxcorr | 低 | 已有 |
| wait-queue 标准模块（check + bind + add-subtask） | examples/vtask | developer/9、examples/vtask | app-fxcorr | 低 | 已有 + reference |
| 信号量/变量/global/semagroup 语义 | app-primes、examples/vtask | developer/9（部分）、appendix/3 | app-fxcorr（计数器） | 高 | reference + invariant |
| vtask add-subtask 双路径（agent 内异步 / `--direct`） | examples/vtask | developer/9、examples/vtask | app-fxcorr | 中 | reference |
| vtask fail / unbind / 资源归还 | examples/vtask | developer/9 | 待补 | 中 | feature |
| 批级/任务级流控组合 | app-fxcorr 实战 | 应用侧 `docs/dynamic-flow-control.md` | 高价值（挂错模块死锁） | 高 | invariant + recipe |
| `free_space_gb` 阶梯兜底 `((@seq+5))` | app-fxcorr | 应用侧 node-tuning | 待迁移 | 高 | invariant |

### 2.7 节点、容器与数据路径

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| 容器引擎 docker / singularity / podman | hello-scalebox、p419 实际 | user/1、appendix/7 | 大量 | 中 | reference + feature |
| slot 启动命令表达式 `((@seq))` / `((@p:))` / `((@v:))` | features/slot_expr | developer/7 §7.7 | app-fxcorr（`((@seq+5))`） | 低 | 已有 |
| 节点本地计算（pod_id + local mapping） | features/pod | developer/5 | app-fxcorr route=1/2 | 中 | feature + invariant |
| 数据卷映射（`/cluster_data_root`、`/local_data_root`、`WORK_DIR`） | 多应用 | developer/5、modules/README | 大量 | 高 | reference + invariant |
| `slot_group` / 多节点 slot | 无 | appendix/2 | app-fxcorr（NODES） | 中 | feature |
| `always_running` / `reserved_on_exit` / `tmpfs_workdir` | side-car、多处 | appendix/2 | app-fxcorr（tmpfs A/B） | 中 | feature + reference |

### 2.8 失败、重试、超时与可观测

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| 退出码语义（0/1/123/124/…） | timeout | appendix/6 | app-fxcorr（fail=0 判据） | 中 | reference（A 档生成） |
| `retry_rules` | features/retry_test | appendix/2 | app-fxcorr（未启用） | 中 | feature |
| 心跳 / GREEDY / slot 状态机 / `slot_timeout_minutes` | 无 | appendix/2、user/5、go-scalebox `docs/slot-status.md`、`docs/seq-diagram.md` | app-fxcorr（节点掉线） | 中 | reference + feature |
| `t_task_exec`（stdout/stderr/extras）诊断 | 多应用 | user/5 | 大量（平台机制实证） | 高 | reference + recipe |
| 时间戳三档 / `timestamps.txt` | agent-interface-files | 应用侧 `docs/timing/` | app-fxcorr | 高 | feature + reference |

### 2.9 集群、跨集群与安全

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| cluster-parameters（data_root / uname / port / grpc_server / remote_check_mode / initial_status） | remote-primes | appendix/2、developer/10 | 待补 | 中 | reference |
| 跨集群 app（main_cluster / 远端子应用） | examples/remote-primes | developer/10、go-scalebox docs | 待补 | 中 | feature |
| host-parameters（slurm_node / default_runtime / use_home_tmp 等） | 无 | appendix/2 | p419 实战 | 中 | reference |
| slot-parameters（group_prefix 等） | 无 | appendix/2 | app-fxcorr（group slot） | 中 | reference |
| 安全（SECURITY_ENABLED / TLS / JWT / 角色） | 无 | developer/11、go-scalebox docs/security | 待补 | 中 | reference |
| WebUI / VS Code | 无 | user/7、user/8 | 待补 | 低 | docs |

### 2.10 CLI 与运维

| 特性/参数 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|
| `scalebox run`（--app-file / -e / --print-yaml） | 全部示例 | appendix/3 | 大量 | 中 | reference（A 档生成） |
| `task add`（--app-id / --sink-module / --header / --conflict-action） | 多应用 | appendix/3 | 大量 | 中 | reference（A 档生成） |
| `app / module / slot / semaphore / variable / global / vtask / cluster / host / queue` 子命令 | 分散 | appendix/3（997 行） | 大量 | 中 | reference（A 档生成） |
| app 创建时环境变量快照 | 无 | 无 | app-fxcorr（多次踩坑） | 高 | invariant |
| 应用级编排模式（批号区间、计数器、两级归并） | app-fxcorr/app-mwa | 应用侧文档 | 高价值 | 中 | recipe |

### 2.11 平台内部机制（设计文档覆盖，平台维护/高级场景）

| 特性/机制 | 现有文档 | 可提炼条目 | 缺口 | 去向 |
|---|---|---|---|---|
| CLI gRPC 架构与跨集群地址解析 | go-scalebox `docs/cli-grpc-architecture.md` | CLI 永远连本地 controld；本地/远端用 `grpc_server` 字符串比较，不用 IP/CIDR | 中 | reference |
| Cluster 状态机与动态集群 allocate/release/renew | `docs/cluster-status.md`、`docs/cluster-admin.md`、`docs/cluster-renew-design.md` | ON→OFF 级联 PAUSED/ERROR；allocations 以组为单位、释放写 `released_at` 不删记录 | 中 | reference + feature |
| Host 状态机与 pg_cron | `docs/host-status.md` | EXPIRED 才能转 PENDING；动态节点必须 node-agent | 中 | reference |
| Slot 八状态机与恢复 | `docs/slot-status.md`、`docs/seq-diagram.md` | 状态变更以 SQL WHERE 为最终防线；SSHFAILED/TIMEOUT 仅 STARTING | 中 | reference + feature |
| Task headers string-only 契约与转义 | `docs/task-headers.md` | 转义顺序反斜杠→引号→逗号，任意值 round-trip 无损 | 高 | invariant + feature |
| 环境变量加载与 `${VAR}` 展开 | `docs/env-loading.md` | 进程环境（含空值）永不被文件覆盖；循环引用保留字面量 | 高 | invariant + feature |
| DAG（物理/逻辑双视角） | `docs/dag-design.md` | 有 `main_router` 的应用缺省 logical；逻辑 DAG 剔除主路由 | 低 | reference |
| 错误处理与 `user_msg` | `docs/error-handling.md` | 业务层只返回 TracedError；高频写路径必须带 `user_msg` | 低 | reference |
| 安全（RBAC/DAC/TLS/JWT） | `docs/security.md`、`docs/security-verification.md` | methodMap 只读映射、写操作仅 admin；agent 不直连 DB | 中 | reference |
| REST API（WebUI） | `docs/rest-api-reference.md`、`docs/webui-vscode-plan.md` | 路径段需 encodeURIComponent；protojson 兼容 camelCase | 中 | reference（A 档生成） |
| shared-body v2（5 接口） | `docs/vtask-shared-body-rest-api.md` | 唯一约束 `(mod_id,from_mod,to_host,vtask,body)` | 中 | reference |
| 手工测试用例集 | `docs/manual-test-plan.md` | 可转 feature 测试的错误路径清单 | 中 | feature |
| 调试（pprof） | `docs/debug.md` | profile 启停命令与端点 | 低 | reference |

## 3. features 现状（2026-10-21 更新）

12 个 feature 目录全部具备 README（按模板）、`feature.yaml`、`test.sh`；测试机 `ssh fxcorr`
连续两轮 10/12 PASS、2 SKIP（多节点环境）、0 FAIL。

| 目录 | README | feature 配置 | 一键测试 | 覆盖特性 | 备注 |
|---|---|---|---|---|---|
| `agent-interface-files` | 有 | `feature.yaml` | `test.sh` ✓ | 8 类接口文件 + sink 7 种格式 | 9 个任务的端到端断言 |
| `check-test` | 有 | `feature.yaml` + `scalebox.env` | `test.sh` ✓ | `free_space_gb` / `dir_quota_gb` / `ACTION_CHECK` | 两次建 app：停领 → 放行 |
| `env-lifecycle` | 有 | `feature.yaml` + `scalebox.env` | `test.sh` ✓ | 快照 / 未定义变量 / APP_ID | — |
| `headers-propagation` | 有 | `feature.yaml` | `test.sh` ✓ | `_` 前缀 header | — |
| `node-progress-gap` | 有 | `feature.yaml` + `tasks.txt` | `test.sh` ✓ | `node_progress_gap` 信号量 | 需 inline 多节点，否则 SKIP |
| `pod` | 有 | `feature.yaml` | `test.sh` ✓ | `pod_id` 同节点路由 | 需 inline 多节点，否则 SKIP |
| `retry-test` | 有 | `feature.yaml` | `test.sh` ✓ | `retry_rules` | 慢测：等 pgcron（5 分钟周期） |
| `side-car` | 有 | `feature.yaml` + `scalebox.env` | `test.sh` ✓ | `ACTION_RUN/SETUP/TEARDOWN/CHECK` | 断言定制 ACTION_RUN |
| `slot-expr` | 有 | `feature.yaml` + `feature-test.yaml` | `test.sh` ✓ | `@seq` / `@v:` / `@p:` | `@p:` 在 slot 无参数时的平台缺陷已修（go-scalebox），镜像重建后可去掉变体文件 |
| `task-batch-size` | 有 | `feature.yaml` + `scalebox.env` | `test.sh` ✓ | `task_batch_size` | 从容器日志验证单次领取 3 条 |
| `task-identity` | 有 | `feature.yaml` | `test.sh` ✓ | `uniq_task` / 112 / `conflict_action` | — |
| `timeout` | 有 | `feature.yaml` + `scalebox.env` | `test.sh` ✓ | `task_max_seconds`（124） | — |

## 4. 缺口分级汇总（初步）

> 阶段 0 的分级结论，保留作优先级依据；各项的当前落地状态见 §3 与 `docs/agents/`、`docs/reference/`、`features/`。

**高（优先做）**
- `_` 前缀 header 传播规则（invariant）
- `free_space_gb` / 流控挂载位置（invariant + feature）
- app 创建时环境快照 + 未定义变量保留字面量 + 非 tty stdin 任务列表（invariant）
- 接口文件族（feature + reference）
- 数据路径规则 / 容器可见路径白名单（reference + invariant）
- 信号量/变量/semagroup 语义（reference + invariant）
- 两套信号量 API 不可混用 + wait-queue gate 解锁（invariant）
- 任务身份与幂等（`uniq_task` / 112 / 禁用 OVERWRITE）（invariant + feature）
- sink-tasks 回流规则与 `last_from_module`（invariant + feature）
- `NODES` 主机正则与节点池可达性（invariant）
- `t_task_exec` 诊断（recipe + reference）
- wait-queue / vtask-tail 契约（不二次 bind、不二次归还、head 必须释放闸门）（invariant）
- header 优先于环境变量（invariant）
- 嵌套 `sink-tasks.txt`（外层先 `rm -f`）（invariant）
- 重投/去重的场景区分（OVERWRITE 与 router 重建的冲突）（invariant，待核实）

**中**
- `task_dist_mode` / pod_id / node_progress_gap / retry_rules / slot_max_tasks / slot_options / 超时与恢复
- 集群/主机/slot 参数、安全、跨集群

**低**
- 大部分 app/module 参数、CLI 命令、WebUI/VS Code（有文档，补 index 即可）

**未确认**
- `local_ip_index`（"暂未使用？"）
- `key_group_regex` / `key_group_index`（"改名?"）
- `task_global_timeout_scale`（"拟改名?"）
- `vtask_size_sema_copy`（已废弃）

## 5. 建议的 feature 优先级（第一批）

1. `check-flow-control`：`free_space_gb` / `dir_quota_gb` + check 挂载位置（修复现有 check_test 参数名）
2. `headers-propagation`：`_` 前缀传播（含 `_vtask_id` / `_vtask_size_sema` / `_vtask_attr_*`）
3. `agent-interface-files`：补齐 app.yaml + test.sh（8 种接口文件）
4. `task-batch-size`：补行为断言（单 task_exec 对多 task）
5. `vtask-flow-control`：从 examples/vtask 抽出最小可断言测试
6. `retry-timeout`：合并 retry_test/timeout，补 README 与断言
7. `pod-local-compute`：补说明与断言
8. `slot-expr`：补 test.sh（现有 README 质量好）
9. `task-identity`：`uniq_task` / 112 / `conflict-action` / `repeatable` 重跑契约
10. `env-lifecycle`：app 快照 / 未定义变量字面量 / non-tty stdin
11. `sink-tasks`：回流格式、`last_from_module`、vs `vtask add-subtask` 的区别

## 6. 待确认

- 哪些参数实际已废弃（`slot_recoverable` 等）——由平台开发者确认；`vtask_size_sema_copy` 已确认由 app-mwa 于 2026-09 退役（见附录 C #33）。
- `features/` 是否允许引用需要真实节点的资源（多节点 feature 的 CI 前提）。
- 第一批 feature 的优先级是否按 §5 执行。
- 设计文档的过时项（RPC 计数、状态名、WebUI 技术栈等）是否在本轮一并修正，还是留到 reference 生成时以源码为准重写。
- `go-scalebox/docs/manual-test-plan.md` 是否作为 feature 测试用例的直接来源。
- `CONFLICT_ACTION=OVERWRITE` 的适用条件：app-mwa（人工重投）与 app-fxcorr（router 重建下游任务，禁用）结论冲突，需平台开发者给出权威边界（按 `from_module`？按是否重建下游？）。
- `task_min_seconds` 的语义易被误解为"补足时长"，官方文档需明确为 slot 健康阈值。
- I/O 计量文件名以连字符为准（`input-files.txt`），需确认 `modules/README.md` 中的下划线写法是否为旧文档。

## 附录 A app-fxcorr 实战实证清单（38 条，阶段 2 蒸馏输入）

来源：`/Users/kaichao/WORK/app-fxcorr`（`CLAUDE.md`、`router/README.md`、`modules/README.md`、`docs/*`）。

### A.1 任务流转与回流

| # | 机制 | 规则/实证 | 来源 | 去向 |
|---|---|---|---|---|
| 1 | 入口默认路由 | `task add <body>` 不带 `--sink-module` 时投给 `main_router` 指定的模块，由它按 `from_module` 分派 | router/README §3.1；main-router-golang-guide §2 | invariant |
| 2 | 回流靠 sink 文件 | 算法模块成功后写 `sink-tasks.txt` 才产生下游；不写 sink = 无回流 | router/README E1/E6、§3.1；modules/README §2.1 | invariant |
| 3 | sink 行格式 | 纯 body / `<mod>,<body>` / `<mod>,<body>,{json headers}`；行内 headers 可直连 vtask-tail 并携带 `batch_status` | modules/README §2.1；router/README §5/E7、§8 | feature |
| 4 | 失败门控与重试 | 非 0 退出不回流；`retry_rules` 按退出码重试，重试执行带 `repeatable:yes`，成功后只回流一次 | router/README E4/§9；main-router-golang-guide §7 | invariant |
| 5 | `last_from_module` | 经主路由转发的任务自动携带上一个非路由模块标识，同一 body 的多阶段模块靠它区分 | router/README E5/§2 | feature |
| 6 | vtask-tail 不自动创建 | 平台不会在 vtask 结束时自动投 tail，必须由 core 的 sink 直连或 router 显式投递 | route3-mpi-plan §3；router/README §5/E7 | invariant |

### A.2 任务身份、幂等与重跑

| # | 机制 | 规则/实证 | 来源 | 去向 |
|---|---|---|---|---|
| 7 | `uniq_task` 唯一索引 | 任务身份 = `(mod_id, from_mod, to_host, vtask, body)`；重复投递报 duplicate key、退出码 112 | router/README E3/§7；modules/README §5 坑1 | invariant |
| 8 | 不要用 OVERWRITE 去重 | `--conflict-action=OVERWRITE` 与 router 重建下游任务撞 duplicate key（112）；重跑应清 `t_task_exec`/`t_task` 并加 app 过滤 | modules/README §5 坑1/坑12；router/README §9 | invariant |
| 9 | 人工重跑契约 | 重投 core 任务须抄回原 `_vtask_id` 与业务头；成功后计数器恰好减一次 | router/README §9 | reference |
| 10 | `task add` 必填项 | 不带 `--app-id` 报 `No valid ref-module-id`；手工测试应显式指定 `--module-id` | CLAUDE.md；modules/fxcorr-x-merge/README | reference |
| 11 | 缺省 router 无差别转发 | 未自研 router 时回流任务被转回算法模块，单模块测试常见 112 噪音；自研 router 按 `from_module` 分派则无回环 | router/README E3；modules/README §2.1/§5 坑9 | reference |
| 12 | 计数器原子归零 | 初值必须精确（多一卡死、少一提前触发）；`AddSemaphoreValue` 返回原子新值、唯一归零者发下游；重复创建 = duplicate key | router/README §7；main-router-golang-guide §6.2 | invariant |

### A.3 vtask、信号量与流控

| # | 机制 | 规则/实证 | 来源 | 去向 |
|---|---|---|---|---|
| 13 | host-bound vtask 配额 | HOST-BOUND 下 `vtask_size` 展开为每节点 `host_vtask_size:<mod>:<host>`；bind 原子预扣，tail 完成/失败由 controld 自动归还 | router/README §5/E8 | invariant |
| 14 | wait-queue gate | `vtask_size:wait-queue` 是批进入 gate（缺省 1 串行）；from_head 末尾必须 `AddValue(...,1)` 解锁，漏做则 wait-queue 停领 | router/README §5；main-router-golang-guide §6.4 | invariant |
| 15 | WQ_VTS 与多节点零扩展 | `WQ_VTS=1` 时 4 节点 wall 与单节点相同（443.7 vs 440.1 s）；`WQ_VTS=8` 后 242.6 s、R 24.9→45.6；建议 `WQ_VTS≈2N` | p419-scale-plan §14.5 | feature |
| 16 | 两套信号量 API 不可混用 | `vtask.CreateSemaphore` 建 vtask 作用域，须配 `vtask.AddSemaphoreValue`；`semaphore.Create` 须配 `semaphore.AddValue`；混用报 `semaphore not found: sql: no rows` | router/CLAUDE.md；main-router-golang-guide §6.2/6.3 | invariant |
| 17 | 实验级任务与 `_` 头传播 | 无 `vtask_role` 的实验级任务不占 batch 配额；`_vtask_*` 只传播给 `vtask_role≠""` 的模块 | router/README §5 | invariant |
| 18 | `free_space_gb` 准入 | 渲染为 hex 编码的 `PLAT_FREE_SPACE_GB`；不满足则 `status_code=-1`、无 `begin_time`、fail=0、留在原 slot 每轮重查，释放后约 5 秒放行 | dynamic-flow-control §2/§7 | feature |
| 19 | 任务级空间闸门会死锁 | 任务级 `free_space_gb` × 组原子性 = 死锁（fail=0、数据不增长）；闸门应挂批入口 vtask-head；`VTASK_SIZE=1` 可消除首轮盲区 | dynamic-flow-control §7.1；local.env | invariant |
| 20 | 阶梯阈值是背压阀 | `((@seq+5))` 给每个 slot 不同阈值；空盘时第一波仍全部放行；真正封顶 = 在产任务数 × 单任务大小 | dynamic-flow-control §5；archive/node-tuning §8.3 | reference |

### A.4 Slot、节点与落点

| # | 机制 | 规则/实证 | 来源 | 去向 |
|---|---|---|---|---|
| 21 | `NODES` 是主机正则 | 不是逗号列表；逗号串整串匹配不到 host → 计算侧 slot 全不建、`host_vtask_size` 缺失、wait-queue 卡 READY、controld 反复报 `GetSemagroupMax: no rows` | archive/node-tuning §10.1；p419.env | invariant |
| 22 | `${NODES}:N` 语义 | 每节点 N 个 slot，序号 0..N-1；服务型模块用 `slots: - ${IO_NODES}` 落 I/O 节点 | archive/node-tuning §8；app.yaml | feature |
| 23 | `slot_on_head` / `reserved_on_exit` | `slot_on_head` 让实验级模块固定落头节点且优先于 `slots`；`reserved_on_exit` 让 vtask-head 的 slot 常驻 | app.yaml；modules/fxcorr-fits/README；local-compute-plan §2.2 | feature |
| 24 | HOST-BOUND 与数据面 | HOST-BOUND 把 vtask 整体钉在一台节点；需要数据面的模块必须落 I/O 节点，控制面节点可能看不到共享 workdir | archive/node-tuning §10.1；router/CLAUDE.md；app.yaml | invariant |
| 25 | 节点池必须 ON 且可达 | Slurm 重分配会换节点名；旧名仍标 ON 但物理已死 → 投递静默定格/链卡死；每轮核对 `t_host` 并 ssh 抽查 | p419.env；p419-scale-plan §13/§15 | docs |
| 26 | slot 启动与派发节奏 | slot 拉起约 5 秒/个；常驻 slot 派发到 run.sh 约 19 ms；短任务看似串行，加 slot 会早饱和 | archive/node-tuning §8.1/§10.7；route3-mpi-plan §4 | reference |
| 27 | 镜像分发形态 | bio 每节点独立 docker daemon，镜像须逐节点分发；p419 共享盘 `.sif` 一份即可；镜像不匹配表现为 STARTING→SSHFAILED | p419-scale-plan §13；modules/README §5 坑11 | docs |

### A.5 环境变量与 app 生命周期

| # | 机制 | 规则/实证 | 来源 | 去向 |
|---|---|---|---|---|
| 28 | env 文件合并顺序 | 先加载基线 `scalebox.env` 再加载 `--env-file`，同名后者覆盖；shell 已导出变量优先级最高 | bio-distributed-steps §2.2；local.env | feature |
| 29 | 未定义变量保留字面量 | 未定义变量展开为字面 `${VAR}` 而非空串，会把字面量传进 app.yaml/模块环境；每个新变量必须显式定义（空值也要写） | p419-distributed-steps §4.3；bio-distributed-steps §7 | invariant |
| 30 | 建 app 即快照 | env、镜像名、slots、`code_path` 在建 app 时快照；run.sh 权限由 slot 容器启动时缓存；改动后必须重建 app | CLAUDE.md；p419-scale-plan §9.3；modules/README §5 坑7 | invariant |
| 31 | 非 tty stdin 是任务列表 | `scalebox run` 把非 tty stdin 内容当起始任务列表；脚本/管道包装必须 `< /dev/null` | bio-distributed-steps 附录 A；p419-scale-plan §9.3 | invariant |
| 32 | 容器不净化宿主环境 | 平台无 `--cleanenv` 等价行为；宿主 `OPAL_*`/`OMPI_MCA_*` 原样透传；模块内须 `unset` 并固定 `LD_LIBRARY_PATH` | route3-mpi-plan §3；modules/mpifxcorr/README | feature |
| 33 | 容器可见路径 | 只有 `/tmp`、`/dev/shm`、`/cluster_data_root`、`/local_data_root`（及 `code_path` 的 `/app/bin`）；workdir 须放 data_root 之下 | modules/README §3 | invariant |

### A.6 观测、数据库与运维

| # | 机制 | 规则/实证 | 来源 | 去向 |
|---|---|---|---|---|
| 34 | 任务状态码语义 | `t_task_exec.status_code`：-1 = 等待/未开始（含被闸门拦）、-3 = 执行中、0 = 完成、>0 = 失败 | route3-mpi-plan §3；dynamic-flow-control §7.4 | invariant |
| 35 | 计时用 t0–t3 | `t_task_exec.t0/t1/t2/t3` 与容器时钟一致；`t_task.create_time/begin_time/end_time` 与 agent 时钟不一致，不可用于计时 | route3-mpi-plan §3 | invariant |
| 36 | 接口文件与日志落库 | `timestamps.txt`、`auxout.txt` 被采集进 `t_task_exec.extras`；`sink-tasks.txt` 触发下游；单任务 stderr 上限 1 GB | modules/README §7；timing/design §7 | reference |
| 37 | 数据库查询口径 | `t_task` 无 app 列（需 join `t_module.m.app`）；`t_task_exec` 经 task 关联；`t_app` 无 `created_at`；`t_host` 用 `uname/ip_addr/port` | p419-distributed-steps §5 | docs |
| 38 | 平台重启后端口同步 | 改 grpc 端口要同步 `t_cluster.parameters->grpc_server`；CLI 需 `PGHOST/PGPORT/GRPC_SERVER`；否则 slot 卡 STARTING | modules/README §5 坑2 | docs |

## 附录 B Sphinx 文档盘点（过时点与缺失主题）

### B.1 过时/矛盾要点（按文件）

| 文件 | 主要问题 | 迁移去向 |
|---|---|---|
| started/1_introduction.md | CLI 命令数"18 域 / 24 叶子"与实际（16 域 / 80+ 叶子）不符 | reference/cli.md |
| started/2_quick_start.md | `cluster status` 不存在；推荐不存在的 `app-copy` | recipes/quickstart |
| user/1_installation.md | `build/` 路径不存在（实为 `runtime/`）；多个 compose overlay 不存在；多处"待测试" | recipes/deploy |
| user/2_core_concepts.md | 集群分类三处自相矛盾（单节点/内联/外部 vs 静态/动态） | glossary + reference |
| user/3_running_apps.md | 8 个命令不存在；`--app-file` 与位置参数写法矛盾；环境变量优先级与 appendix/3 不一致 | recipes/run-app |
| user/4_example_apps.md | `examples/app-copy` 不存在；Makefile 目标与 CLI 用法不符 | recipes/examples |
| user/5_operations.md | `docker compose` 缺 `-f runtime/compose.yaml`；metrics 命令不存在 | recipes/ops、debug |
| user/6_standard_modules.md | `slot_options`/`cluster` 放错层（parameters vs arguments）；未收录 wait-queue | reference/module.md |
| user/7_webui.md | `build/compose.yaml` 不存在；端点数 76/77 不一致；"开发者"角色未定义 | reference/webui.md |
| user/8_vscode.md | `vscode-scalebox/` 不在本仓库；GetAppDAG 状态过时 | reference/vscode.md |
| developer/1_programming_model.md | 标题编号错乱；重复"图1" | docs（保留，修编号） |
| developer/2_module_development.md | `cleanup-files.txt (待实现)`；若干字段空说明 | reference/module.md + feature |
| developer/3_app_design.md | 笔误；无实质过时 | docs（保留） |
| developer/4_main_router_status.md | `semaphore increment-n`、`variable create` 与 CLI 不一致；vtask 表空白 | reference |
| developer/5_node_local_compute.md | 估算值无实测口径 | docs（保留） |
| developer/7_advanced_features.md | `space_free_gb` 与 `free_space_gb` 冲突；编号重号；空章节（7.5/7.6）；空代码块 | reference + features |
| developer/9_vtask.md | `vtask_size_sema_copy` 已废弃仍列有效；tail unbind 描述自相矛盾；命令参数与 appendix/3 不一致；子命令数 13/14 不一致 | reference + invariants |
| developer/10_cross_cluster.md | `app get-router`、`task add --remote-cluster` 未在 CLI 参考 | reference |
| developer/11_security.md | `build/`、`gen-secrets.sh`、`make copy-secrets` 均不存在 | reference/security.md |
| appendix/1_tech_specifications.md | `environmens` 拼写；占位值；schema 大量 `...`；孤立性能表 | reference/app-yaml.md |
| appendix/2_parameter_reference.md | 多处"（？）""拟修改""暂未使用"；`port` 描述错误；`ssh_cepher` 拼写 | reference（生成/半自动） |
| appendix/3_commandline_tools.md | 标题无正文；编号错乱；`app create`/`app run` 旧命令；占位符 `xxxx` | reference/cli.md（A 档生成） |
| appendix/4_shell_programming.md | 标题无正文；内置函数无版本说明 | reference/module.md |
| appendix/6_exit_code_spec.md | 16 位状态码与 0–255 退出码关系未说明；225 重复且描述互换 | reference/exit-codes.md（A 档生成） |
| appendix/7_software_install.md | CentOS 8 已 EOL；非法命令；建议归档 | 归档 |
| faq.md | `build/compose.yaml` 不存在；仍讲已废弃的可编程信号量副本；GetAppDAG 状态过时 | recipes/debug |

### B.2 缺失主题（平台已有、文档没有）

1. wait-queue 标准模块（仅 vtask 章节提及）
2. Slurm/动态集群与 cluster-admin（资源申请、节点续期）
3. channel 跨应用优先队列（仅命令参考）
4. scalebox-fs（仅命令参考）
5. event 事件系统（仅命令参考）
6. pod_id / pod 部署生命周期
7. slot_options 运行选项的使用场景与排错
8. task_batch_size 的批量领取与 task_exec 聚合关系
9. retry_rules 完整语义（退出码匹配、次数、vtask 联动）
10. agent 接口文件完整格式（task-exec.yaml、extra-attrs.yaml 等 8 种）
11. I/O 计量与数据库视图（v_task_iobytes、v_task_network_io）
12. 环境变量优先级权威表（三处说法不一）
13. 平台表结构（t_app/t_task/t_slot/t_vtask 等）
14. REST API 端点表（Sphinx 只外链）
15. DAG / GetAppDAG 数据模型
16. singularity/apptainer 运行指南
17. 监控指标（host/app/node metrics 的命令与口径）

### B.3 与 features 重复的文档

| 文档 | 重复 feature |
|---|---|
| developer/2_module_development.md | agent-interface-files、side-car、retry_test |
| developer/7_advanced_features.md | retry_test、check_test、node_progress_gap、timeout、slot_expr、task_batch_size |
| appendix/2_parameter_reference.md | task_batch_size、node_progress_gap、check_test |
| appendix/3_commandline_tools.md | slot_expr |
| user/2_core_concepts.md、developer/5_node_local_compute.md | pod |
| user/6_standard_modules.md | side-car、timeout |
| faq.md | check_test |

## 附录 C app-mwa 实战实证清单（35 条：阶段 2 蒸馏输入 + 2026-10 增量）

来源：`/Users/kaichao/WORK/workspace/apps/app-mwa/pipelines/beam-form`（`CLAUDE.md`、`Design.md`、`app/`、`modules/*/README.md`、`docs/upgrade-impact-analysis.md`）；C.7 增量来自提交 `7e55008`、`f89a621`、`91ba83e`、`ae71311` 与 `pkg/module` 源码核对。

### C.1 路由与应用装配

| # | 机制 | 规则/实证 | 来源 | 去向 | 备注 |
|---|---|---|---|---|---|
| 1 | 主路由分发约定 | `module.Run(fromFuncs, "from_module")` 按 header 分派，`""` 是入口；参数/header 解析失败时会跳过 handler | app/router/main.go、app-base/router/main.go | invariant | 补充 A1 |
| 2 | init 时序 | `module.InitLogging()` 必须早于 `initApp()`，且不能放在 `init()` 里，否则 `module.LogEntry` 为 nil | app/router/init.go | invariant | 新增 |
| 3 | 应用启动与状态 | `start_task` 是创建时注入的首条消息；`initial_status: RUNNING` 免手工放行；运行中可 `app set-status` | app-base/app.yaml、app/README.md | reference | 补充 |
| 4 | 路由模块参数 | 路由用 `task_dist_mode: SLOT-BOUND` + `task_batch_size` 批量取消息；`router_index` 指定第几个主路由实例 | app/app.yaml | reference | 补充 |
| 5 | `slot_group` | `{模块: 数量}` 声明新增节点时自动创建的 slot 数，扩容不必逐台维护 `slots` | app/app.yaml | reference | 补充 |
| 6 | `pod_id` / `key_group_regex` | `pod_id` 让同源 task 走 HOST-BOUND 本地计算；`key_group_regex` 从 body 提取分组键，共同决定落点 | app/app.yaml | reference | 补充 |

### C.2 vtask 流水线

| # | 机制 | 规则/实证 | 来源 | 去向 | 备注 |
|---|---|---|---|---|---|
| 7 | vtask 四段结构 | `wait-queue → vtask-head → core → vtask-tail`，用 `sink_vmodules` 串接；wait-queue 配 `VTASK_MODE=GROUP-BOUND`、`vtask_size: 1` 串行准入 | app/app.yaml、Design.md §3 | feature | 补充 |
| 8 | wait-queue 不重复 bind | 标准 run.sh 已做准入 + `vtask bind` + `add-subtask --direct`；路由端 `fromWaitQueue` 必须空实现，不能再 bind 或写 sink-tasks | app/router/wait_queue.go | invariant | 补充 A14 |
| 9 | vtask-head 释放闸门 | head 创建 `cube-vtask-done:<body>` 计数器并对 `vtask_size:wait-queue` 加 1；漏做后续 vtask 永远进不来 | app/router/vtask_head.go | invariant | 补充 A14 |
| 10 | vtask-tail 不重复归还 | tail 只递减计数器；资源由 controld `doVTaskFinished` 自动归还，不能再调 `UnbindResource` | app/router/vtask_tail.go、upgrade-impact-analysis §3 | invariant | 补充 A13 |
| 11 | `_` 前缀 header 传播 | 所有 `_` 前缀 header 自动向下游传播；`_vtask_id`/`_vtask_size_sema` 平台注入，应用自定义 `_vtask_cube_name` 等靠此传递 | Design.md §3、app/router/beam_make.go | invariant | 补充 A17 |

### C.3 信号量与流控

| # | 机制 | 规则/实证 | 来源 | 去向 | 备注 |
|---|---|---|---|---|---|
| 12 | 倒计数信号量 | 计数器按本批预期条数初始化，每完成一条 decrement，返回值 ≤0 触发下游；多并发安全 | app/router/pull_unpack.go | invariant | 补充 A12 |
| 13 | `node_progress_gap` | 按主机名正则配置快慢节点允许的进度差（如 `{"[a-z][0-9][0-9]":3}`），平台自动建 `node_progress:<module>:<host>` | app/app.yaml | feature | 补充 |
| 14 | 空间双闸门 | `free_space_gb` 是分区保留最小空间，`dir_quota_gb` 是目录容量上限，均在 `module.arguments` 按路径配置 | app/app.yaml | feature | 补充 A18/A19 |
| 15 | 重试与重投 | `retry_rules` 按退出码配次数（`*:3` 兜底）；重复投递设 `CONFLICT_ACTION=OVERWRITE`，人工重跑加 `repeatable=yes`；去重窗口要覆盖任务最长运行时间 | app/app.yaml、p419/README.md | feature | **与 A8 条件冲突，蒸馏时按场景拆分** |
| 16 | 任务时长参数 | `task_min_seconds` 是 slot 健康阈值（连续低于判 GREEDY），不是补足任务时长；`task_max_seconds` 才是超时 | app/app.yaml | reference | 新增（易误解） |

### C.4 任务头与模块契约

| # | 机制 | 规则/实证 | 来源 | 去向 | 备注 |
|---|---|---|---|---|---|
| 17 | header 优先于环境变量 | `scalebox::task_header` 取值时 task header 优先于同名环境变量；测试可用 `task add -h` 覆盖而不改镜像 | modules/pull-unpack/README.md | invariant | 新增 |
| 18 | 每任务下游与超时 | `task.Add*` 的 env 中 `SINK_MODULE` 指定下游、`TIMEOUT_SECONDS` 覆盖模块级超时 | app/router/beam_make.go | reference | 新增 |
| 19 | `sink-tasks.txt` 契约 | 成功时写 `${WORK_DIR}/sink-tasks.txt`；嵌套调用 file-copy 也会生成同名文件，外层必须先 `rm -f` 再写，失败时不要写 | modules/fits-redist/code/run.sh | invariant | 补充 A2/A3 |
| 20 | I/O 计量文件 | 连字符文件名：`input-files.txt`、`output-files.txt`、`removed-files.txt`、`network-files.txt`、`extras.yaml`；平台汇总到 `t_task_exec`；调试信息追加 `auxout.txt` | modules/README §一 | feature | 补充 §2.5 |
| 21 | `tmpfs_workdir` | `slot_options: tmpfs_workdir` 把工作目录放 tmpfs，避开 `/tmp` 清理；大文件场景要核对容量 | app/app.yaml | reference | 补充 |

### C.5 存储虚拟化与 gRPC（项目私有为主）

| # | 机制 | 规则/实证 | 来源 | 去向 | 备注 |
|---|---|---|---|---|---|
| 22 | vpath 池成员与容量 | 池成员由信号量组 `vpath:free-gb:<pool>` 决定；成员路径存共享变量；容量即信号量值 | app/vpath.yaml、internal/vpath/README.md | reference | **项目私有，不进平台文档** |
| 23 | `Locate` vs `Allocate` | `Locate` 只查不写；`Allocate` 在未分配时选池并扣容量且对同 key 幂等；查询误用会凭空占容量 | internal/vpath/README.md | invariant | 项目私有 |
| 24 | 成员路径含冒号 | 信号量名解析必须 `SplitN(name, ":", 4)[3]`，不能用 `Split` | internal/vpath/README.md | invariant | 项目私有 |
| 25 | gRPC 变量不存在 | 返回 `codes.Internal` + `sql: no rows in result set`，无 `sql.ErrNoRows`；客户端判空要覆盖三种形态 | internal/varerr/varerr.go | reference | 疑似平台通用 |

### C.6 运维与排错

| # | 机制 | 规则/实证 | 来源 | 去向 | 备注 |
|---|---|---|---|---|---|
| 26 | 平台库查进度 | 直接查 `t_task`/`t_slot`/`t_semaphore` 做进度横表；`status_code` -1 READY / -3 RUNNING / 0 成功 / >0 失败；人工重跑改回 -1 并加 `repeatable=yes` | app/sql-htable.md、p419/README.md | reference | 补充 A34 |
| 27 | 节点被 drain 的假 RUNNING | Slurm drain/down 后任务仍停在 RUNNING(-3) 不自动失败；>60 节点出错率上升，120 节点时 >1/4 任务超时，需人工重置 | app-base/README.md §三 | reference | 新增 |
| 28 | 单模块测试 | `scalebox app run --image-name --code-path --slot-regex` 起测试应用，再 `task add --sink-module next -h k=v`；`reserved_on_exit` 保留容器排查 | modules/pull-unpack/README.md | docs | 补充 |
| 29 | 应用内重启清单 | 节点 renew 清临时空间 → 重置 slot → `node_progress` 归零 → 扣减已完成部分 → 重置首次加载信号量 → 清共享临时目录 | app/README.md | docs | 新增（recipe 候选） |

### C.7 平台标准化改造增量（2026-09 / 2026-10）

本轮从 app-mwa 的标准化改造提交与 `pkg/module` 源码核对得到、C.1–C.6 未覆盖的条目。

| # | 机制 | 规则/实证 | 来源 | 去向 | 备注 |
|---|---|---|---|---|---|
| 30 | `pkg/module` 统一模块入口 | `module.Run(handlers, "from_module")` 接管参数校验、header 解析、按 header 分派与退出码；模块 `main` 只剩 handlers 表 + `initApp`。`module.InitLogging()` 需早于 `initApp()`，且幂等（`Run` 内部再调一次无害） | app/router/main.go；`scalebox/pkg/module/module.go`（`Run`） | reference | C1 补充；`module.md` 待补「Go 模块入口」一节 |
| 31 | 模块入口退出码契约 | `0` OK；`1` 参数不足 **或 handler 返回未包装的普通 error**；`2` headers JSON 解析失败（`ExitBadHeaders`）；`4` 无匹配 handler。handler 返回 `*errors.TracedError` 时用其 Code | `scalebox/pkg/module/module.go:20-23,58-96` | reference | 新增；注意与任务状态码 `2=StatusCodeInvalidArgs` 同名不同源；普通 error 会被压成 1，掩盖原因，应 `errors.WrapE(err, code, ...)` |
| 32 | 标准 wait-queue 接管 vtask 生命周期 | 直接用镜像 `scalebox.net/platform/wait-queue`（`VTASK_MODE` 驱动准入/bind 模式，`vtask_size: 1` 串行）；`run.sh` 内做 `vtask bind` + `add-subtask --direct`。路由端 `fromWaitQueue` / `fromVtaskTail` 必须空实现，本地 `check.sh` 不再需要 | app/app.yaml；app/router/wait_queue.go | invariant | 补充 C8 |
| 33 | vtask 应用属性与退役的双信号量 | 应用自定义数据用 `_vtask_cube_name` 等 `_` 头随 sink 链传播；`slot_vtask_size` / `:slot_vtask_size` 双信号量与 `vtask_size_sema_copy` 已退役 | app/router/beam_make.go；upgrade-impact-analysis §3 | reference | 新增；关闭 §6 的 `vtask_size_sema_copy` 待确认 |
| 34 | 路由端不释放 vtask 资源 | vtask-tail 只递减业务计数器；资源由 controld `doVTaskFinished` 按 `_vtask_size_sema` 自动归还，脚本再调 `UnbindResource` 会多归还（额度虚增） | app/router/vtask_tail.go | invariant | 补充 C10 |
| 35 | env 文件分层（应用约定） | `scalebox.env` 只放空基线（每个被引用变量都定义，防 `${VAR}` 字面量透传），站点 env（`p419*.env`、`dcu*.env`）覆盖；不靠 shell 变量跨站点 | pipelines/main/scalebox.env、p419.env | 应用仓 | 项目约定，不进平台；与 INV-02 互为表里 |

