# 00 — 平台心智模型（一页纸）

> 读法：本页只建立全局图，不展开规则。红线与坑见 [`10-invariants.md`](./10-invariants.md)；
> 卡链排查见 [`20-diagnosis.md`](./20-diagnosis.md)；入口与命令见仓库根 [`AGENTS.md`](../../AGENTS.md)。

## 1. 五层对象

从大到小：**Cluster → Host → Slot → Module → Task**；**App** 把 Module 组织成一条流水线，
**VTask** 把跨 Module 的 Task 串成一个应用级计算单元。

| 层 | 是什么 | 关键事实 |
|---|---|---|
| Cluster | 计算资源的逻辑分组 | 单节点 / 静态 / 动态（Slurm）/ 内联；跨集群靠 `t_cluster` 全网格复制 + gRPC 代理 |
| Host | 物理节点 | 必须 `status=ON` **且 ssh 可达**；`slots` 按 hostname 短名匹配 |
| Slot | 节点上与某 Module 对应的资源切片 | 由 actuator 按 `slots` 正则展开后创建；容器在 slot 里长期存在、循环领任务 |
| Module | 容器化算法单元 | app 创建时快照镜像/参数/slot 展开结果；容器内跑 `run.sh` |
| Task | 一次执行 | `body`（stdin）+ headers（环境变量）+ `status_code`；无状态、可重跑 |

**App** = `app.yaml` 里的一组 Module + 参数 + 集群。**VTask** = 带 `_vtask_id` 的 Task 集合，
管道为 `task → [wait-queue] → vtask-head → vtask-core → … → vtask-tail`，用于批/组级流控与资源绑定。

## 2. 一条任务的旅程

```text
task add（入口，body 走 stdin，可带 -H header）
  → controld 落库 t_task / t_task_exec（身份 = mod_id, from_mod, to_host, vtask, body）
  → 目标 Module 的空闲 Slot 领取（agent 在 slot 内轮询）
  → agent 启动容器：body 进 stdin，headers/平台变量进环境，WORK_DIR 挂载
  → run.sh 执行；成功退出码 0
  → agent 读 ${WORK_DIR}/sink-tasks.txt，为每条下游投递新任务
  → 下游重复以上过程，直到没有回流
```

三个要点：

1. **回流只靠 `sink-tasks.txt`**——成功才写；空文件是错误；没写就没有下游（INV-04）。
2. **`_` 前缀 header 自动传播**，普通 header 不自动传播；router 与 wait-queue 有门控（INV-13）。
3. **失败即停**：任务 `status_code > 0` 时不回流，链停在该批/该分支。

## 3. 控制面与数据面

| 面 | 组件 | 职责 |
|---|---|---|
| 控制面 | `controld` | gRPC 服务：app/module/slot/task/vtask 的增删查与调度决策 |
| 控制面 | `actuator` | 在节点上启动 slot（SSH / 外部调度器） |
| 控制面 | `database` | PostgreSQL：`t_app`、`t_module`、`t_slot`、`t_task`、`t_task_exec` 等 |
| 控制面 | `webui` / `cluster-admin` | Web 接口；与 HPC 调度器交互的资源申请应用 |
| 数据面 | `agent` | 节点代理：领任务、起容器、收集 stdout/stderr 与退出码、写回流 |
| 数据面 | Slot / 容器 | 实际执行单元；只有 `/tmp`、`/dev/shm`、`/cluster_data_root`、`/local_data_root`、`/app/bin` 可见 |
| 数据面 | 共享存储 | 可选；放 `data_root` 之下的工作目录，容器间通过挂载可见 |

**判断归属**：改配置/查状态 → 控制面与 DB；数据读写/性能 → 数据面与挂载盘。

## 4. 三层知识怎么查

| 问题类型 | 去哪 |
|---|---|
| 为什么、红线、踩过的坑 | [`10-invariants.md`](./10-invariants.md)（规则）、`features/*/README.md`（特性细节） |
| 有什么：字段/命令/表/API | [`docs/reference/`](../reference/)（唯一信源） |
| 怎么做：步骤与判据 | [`recipes/`](./recipes/README.md)（任务食谱）、`features/*/`（可跑测试）、`examples/` |

## 5. 常见误读

- **Module ≠ Slot ≠ Task**：Module 是定义，Slot 是节点上的执行实例，Task 是一次调用；
  同一 Module 可以有多个 Slot，一个 Slot 顺序执行多个 Task。
- **App 状态 ≠ 任务进度**：`app show` 显示 RUNNING 只表示 app 未被终止；任务可能停在 READY。
- **建 app 是快照**：改 env / 镜像 / `slots` / `code_path` / `run.sh` 后必须重建 app（INV-01）。
- **`slots` 是主机正则**：匹配 hostname 第一段（短名），不是逗号列表，不要写 FQDN（INV-15）。
- **vtask 不是 "多节点"**：它是跨模块的任务集合；是否多节点由 HOST-BOUND / GROUP-BOUND 绑定决定。
