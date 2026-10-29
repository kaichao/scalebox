# 术语表（中英对照）

> 面向 agent。每条一句话；需要细节时按“详见”列。术语同时是源码标识符的，保留原文拼写。

## 核心对象

| 中文 | 英文 / 标识符 | 含义 | 详见 |
|---|---|---|---|
| 应用 | App | 一次应用实例，由 `app.yaml` 定义模块与参数；创建时快照配置 | `reference/app-yaml.md` |
| 模块 | Module | 容器化的算法单元；app 内的一个处理环节 | `reference/app-yaml.md`、`00-model.md` |
| 任务 | Task | 模块的一次执行：`body` + headers + `status_code` | `10-invariants.md` INV-04/06 |
| 节点 | Host | 执行计算的物理节点；须 `status=ON`（同 cluster）且 ssh 可达；`DRAINING` 为排空状态：skip 模块停领并退出，其他模块继续 | `10-invariants.md` INV-16 / INV-21 |
| 运行槽 | Slot | 节点上与某模块对应的资源切片；本 cluster `ON`/`DRAINING` host 上的 READY slot 会被启动（DRAINING 上按 `skip_draining_host` 过滤），容器在其内循环领任务，空闲退出回 READY | `00-model.md` §1 |
| 集群 | Cluster | Host 的逻辑分组（单节点/静态/动态/内联） | `application/1_core_concepts`（Sphinx） |
| 虚拟任务 | VTask | 跨模块的 Task 集合，用于批/组级流控与资源绑定 | `10-invariants.md` INV-09…14 |

## 运行时组件

| 中文 | 英文 / 标识符 | 含义 |
|---|---|---|
| 控制服务 | `controld` | gRPC 控制面：调度决策与元数据操作 |
| 启动器 | `actuator` | 在节点上启动 slot（SSH 或外部调度器） |
| 数据库 | `database` | PostgreSQL，存 app/module/slot/task/vtask 元数据 |
| 节点代理 | `agent` | 领任务、起容器、收集结果、写回流 |
| Web 界面 | `webui` | 浏览器操作入口；生成 `app-schema.json` 的宿主 |
| 集群管理应用 | `cluster-admin` | 与 HPC 调度系统交互的资源申请应用（系统级 app） |
| 主路由 | `main_router` | 应用内的控制面模块：执行并行逻辑、读写应用状态、分发任务、组织模块协作；不是平台控制面（详见 `application/1_core_concepts`） |
| 等待队列 | wait-queue | vtask 标准模块：准入、`bind`、预扣额度 |

## 任务与接口

| 中文 | 英文 / 标识符 | 含义 | 详见 |
|---|---|---|---|
| 任务体 | `body` | 任务输入，从 stdin 传入；JSON 或无空白字符串 | `features/agent-interface-files` |
| 任务头 | task header | 随任务传播的键值元数据；经环境变量注入 | `reference/task-headers.md` |
| 回流 | `sink-tasks.txt` | 模块成功后写的下游任务清单；无文件即无下游 | `features/headers-propagation` |
| 任务身份 | `uniq_task` | `(mod_id, from_mod, to_host, vtask, body)`；重复投递给 112 | `features/task-identity` |
| 冲突处理 | `conflict_action` | 重复投递策略：`OVERWRITE` / `IGNORE`；一次投递单事务，任意一行冲突整批回滚 | `reference/task-headers.md`、INV-07 |
| 状态码 | `status_code` | `-1` 等待、`-3` 执行中、`0` 完成、`>0` 失败 | `reference/exit-codes.md` |
| 退出码 | exit code | 任务进程返回值；平台语义见参考页 | `reference/exit-codes.md` |

## 流控与信号量

| 中文 | 英文 / 标识符 | 含义 | 详见 |
|---|---|---|---|
| 信号量 | semaphore | 计数同步原语；两套 API 不可混用 | `reference/semaphore-variable.md`、INV-12 |
| 任务组头 | vtask-head | 管道入口；`add-subtask` 后必须归还 wait-queue gate | INV-10 |
| 任务组核心 | vtask-core | 实际计算段；HOST-BOUND/GROUP-BOUND 的绑定落点 | `features/pod` |
| 任务组尾 | vtask-tail | 管道收尾；资源由 `doVTaskFinished` 自动归还 | INV-11 |
| 空间闸门 | `free_space_gb` / `dir_quota_gb` | slot 停领任务的准入条件；挂批入口，不挂任务级 | `features/check-test`、INV-14 |
| 批大小 | `task_batch_size` | 一次领取的任务条数 | `features/task-batch-size` |

## 绑定与落点

| 中文 | 英文 / 标识符 | 含义 | 详见 |
|---|---|---|---|
| 主机绑定 | HOST-BOUND | vtask 绑定到单节点；需要 `to_host` 或 `pod_id` | `features/pod` |
| 组绑定 | GROUP-BOUND | vtask 绑定到节点组；组内 slot 协同 | `application/1_core_concepts`（Sphinx） |
| 槽绑定 | SLOT-BOUND | 任务指定 `to_slot` 落点 | `reference/2_parameter_reference §2.1.4` |
| 节点本地计算 | `pod_id` | 同一 pod 的模块落在同一节点 | `features/pod` |
| 独立组节点 | group slot | `slots` 条目的 `<host>:<count>:<group_prefix>` 形式；组锚点节点上领取匹配前缀的 task | `features/group-node-drain`、Sphinx `design/3_vtask.md` §3.3.3 |
| 排空 | `DRAINING` / `skip_draining_host` | host 排空状态；skip 模块停领并退出 slot，其他模块继续消费本地链 | `features/group-node-drain`、INV-21 |
| 头节点槽 | `slot_on_head` | 在头节点建一个 slot，便于单节点调试 | `reference/app-yaml.md` |

## 配置字段

| 中文 | 英文 / 标识符 | 含义 | 详见 |
|---|---|---|---|
| 基础镜像 | `base_image` | 模块容器的默认镜像；测试部署常用 `scalebox.net/platform/agent:latest` | `reference/app-yaml.md` |
| 代码挂载 | `code_path` | 把模块代码挂到容器 `/app/bin` | `reference/app-yaml.md` |
| 工作目录 | `WORK_DIR` | 模块的本地工作目录，缺省 `/work` | `reference/2_parameter_reference §2.3`（Sphinx） |
| 数据根 | `data_root` | 集群数据根，容器内 `/cluster_data_root` | `reference/app-yaml.md` |
| 变量文件 | `--env-file` / `scalebox.env` | 建 app 时的变量来源；未定义变量保留字面量 | INV-02 |
| 脚本动作 | `ACTION_RUN` / `SETUP` / `TEARDOWN` / `CHECK` | sidecar 启动/收尾/准入脚本路径 | `features/side-car` |
| 平台前缀 | `PLAT_*` | 平台注入的环境变量；用户程序无需访问 | `reference/env-vars.md` |
