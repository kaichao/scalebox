# 1. 常见问题（FAQ）

## 1. 部署与连接

### 1.1 非标准端口设置

Scalebox 涉及的网络端口：
- 数据库端口，标准 5432
- controld gRPC 端口，标准 50051
- WebUI HTTP 端口，标准 8088

**服务端（Docker Compose）**：

```yaml
# runtime/compose.yaml
services:
  controld:
    ports:
      - <custom-grpc-port>:50051
  database:
    ports:
      - <custom-pg-port>:5432
  webui:
    ports:
      - <custom-web-port>:8088
```

**客户端**：设置 `${HOME}/.scalebox/environments`：

```bash
GRPC_SERVER=<internal-ip>:<custom-grpc-port>
PGHOST=<internal-ip>
PGPORT=<custom-pg-port>
```

### 1.2 数据库连接失败

**现象**：`dial tcp: connect: connection refused`

**解决**：
```bash
# 等待数据库完全启动
sleep 10

# 手动检查
docker exec database psql -U scalebox -c "\l"

# 重新创建数据库卷
docker compose -f runtime/compose.yaml down -v
docker compose -f runtime/compose.yaml up -d
```

### 1.3 gRPC 连接失败

**现象**：CLI 命令报 `Unavailable` 或 `connection refused`

**解决**：
```bash
# 检查 GRPC_SERVER 设置
echo $GRPC_SERVER

# 检查 controld 服务状态
docker compose -f runtime/compose.yaml ps controld

# 验证 gRPC 连通性
grpcurl -plaintext localhost:50051 list
```

## 2. 任务与调度

### 2.1 任务一直处于 READY 状态

见 [卡链与失败诊断](https://github.com/kaichao/scalebox/blob/dev/docs/agents/20-diagnosis.md) 的“任务 READY 不动”和“slot 全是 0 行”。

### 2.2 任务执行失败（exit code != 0）

先用 `scalebox task show <task-id>` 查看 stdout/stderr 与 `status_code`；退出码含义见
[退出码规范](appendix/6_exit_code_spec.md)。任务不动或下游不出现，按
[卡链与失败诊断](https://github.com/kaichao/scalebox/blob/dev/docs/agents/20-diagnosis.md) 排查。

### 2.3 批量任务重复创建

相同 `(mod, from, host, vtask, body)` 重复投递会触发 uniq_task/112；用 `--conflict-action IGNORE` 或 header `repeatable: yes`。见
[task-identity](https://github.com/kaichao/scalebox/blob/dev/features/task-identity/README.md)。

## 3. 信号量与流控

### 3.1 信号量不存在错误

设置 `SEMAPHORE_AUTO_CREATE=yes` 自动创建（初值 0），或显式创建信号量。命令与返回码见
[appendix/3 §1.3.3](appendix/3_commandline_tools.md)。

### 3.2 信号量值异常（不收敛）

检查流控版（无前缀，controld 维护）与可编程版（`:` 前缀，脚本操作）是否匹配。见
[VTask §9.4](developer/9_vtask.md)。

## 4. WebUI

### 4.1 页面无法访问

```bash
# 检查 WebUI 服务
docker compose -f runtime/compose.yaml ps webui

# 检查端口
curl http://localhost:8088/api/v1/stats/system
```

### 4.2 创建操作返回错误

常见原因：proto 字段名不匹配。REST API 同时接受 camelCase 和 snake_case。

### 4.3 DAG 页面显示空白

controld 侧 `GetAppDAG` RPC 可能未实现，页面降级为模块拓扑。

## 5. VS Code 插件

### 5.1 TreeView 显示空列表

检查配置：
- `scalebox.serverUrl` 是否正确指向 WebUI
- WebUI 是否可访问：`curl <serverUrl>/api/v1/apps`
- 查看 VS Code Output 面板（Scalebox 通道）

### 5.2 app.yaml 校验不生效

确认：
- 文件名为 `app*.yaml` 或 `scalebox*.yaml` 格式
- 文件语言模式为 YAML
- VS Code 未禁用 Scalebox 扩展

## 6. 日志与调试

### 6.1 调整日志详细度

```bash
# 详细模式（完整错误链 + 位置信息）
LOG_LEVEL=debug scalebox task list

# 静默模式
LOG_LEVEL=error scalebox task list
```

### 6.2 收集服务日志

```bash
# 所有服务
docker compose -f runtime/compose.yaml logs > scalebox.log

# 按服务
docker compose -f runtime/compose.yaml logs controld
docker compose -f runtime/compose.yaml logs actuator

# 按时间
docker compose -f runtime/compose.yaml logs --since "2026-01-01"
```

### 6.3 gRPC 错误信息解读

错误信息格式：`rpc error: code = <Code> desc = <Message>`

常见 gRPC 状态码与对应场景：

| 状态码 | 含义 | 常见原因 |
|--------|------|---------|
| `Unavailable` | 服务不可达 | controld 未启动或网络不通 |
| `NotFound` | 资源不存在 | App/Module/Task ID 无效 |
| `AlreadyExists` | 资源已存在 | 重复创建同名校验失败 |
| `InvalidArgument` | 参数错误 | flag 值不合法 |
| `Internal` | 内部错误 | DB 异常或业务逻辑错误 |

## 7. 概念与架构

### 7.1 主路由模块与平台控制面是什么关系？

主路由是应用内的控制面模块（应用的一部分），执行并行逻辑、读写应用状态、分发任务、组织模块协作；平台控制面是 controld/actuator，负责资源、slot、容器与任务调度。主路由不管理容器/节点/slot，平台不执行应用的计算逻辑。

详见 {doc}`核心概念 §2.5 <user/2_core_concepts>` 与 {doc}`主路由与应用状态：实现与优化 <developer/4_main_router_status>`。

### 7.2 模块都无状态，应用状态存在哪里？

应用状态存在平台存储中：semaphore（计数/流控/进度）、variable（应用内字符串状态）、global（同一物理集群内跨应用共享）。跨物理集群的信息共享使用 queue。task 不携带应用状态；module 不持有应用状态；主路由集中读写，controld 在 vtask 等场景也会读写信号量。

详见 {doc}`核心概念 §2.4 <user/2_core_concepts>`。

### 7.3 计算逻辑、算法逻辑、并行逻辑是什么关系？

计算逻辑 = 算法逻辑（模块内）+ 并行逻辑（模块间）。算法逻辑回答“一个小问题怎么算”，由模块实现；并行逻辑回答“大问题怎么分解、分解后怎么协作”，由主路由执行。应用逻辑与计算逻辑同义。

详见 {doc}`核心概念 §2.3 <user/2_core_concepts>`。

### 7.4 VTask 和主路由是什么关系？

VTask 是对并行逻辑及其状态的粗粒度封装，把跨模块 task 组织为一个计算单元；它不是独立调度器，调度由平台与主路由完成。主路由执行并行逻辑并读写应用状态；VTask 在跨模块层承载状态、流控、资源绑定与容错。

详见 {doc}`核心概念 §2.6 <user/2_core_concepts>` 与 {doc}`VTask 虚拟任务 <developer/9_vtask>`。

### 7.5 global 和 queue 有什么区别？

global 用于同一物理集群内多个应用之间共享信息/配置；queue 用于跨物理集群的管理面消息共享。queue 不是 wait-queue，也不是平台内部任务队列。

详见 {doc}`核心概念 §2.10 <user/2_core_concepts>` 与 {doc}`跨集群架构 §10.6 <developer/10_cross_cluster>`。

> 持续更新中，欢迎贡献问题与答案。
