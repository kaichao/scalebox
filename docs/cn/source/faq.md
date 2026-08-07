# 常见问题（FAQ）

## 1. 部署与连接

### 1.1 非标准端口设置

Scalebox 涉及的网络端口：
- 数据库端口，标准 5432
- controld gRPC 端口，标准 50051
- WebUI HTTP 端口，标准 8088

**服务端（Docker Compose）**：

```yaml
# build/compose.yaml
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
docker compose -f build/compose.yaml down -v
docker compose -f build/compose.yaml up -d
```

### 1.3 gRPC 连接失败

**现象**：CLI 命令报 `Unavailable` 或 `connection refused`

**解决**：
```bash
# 检查 GRPC_SERVER 设置
echo $GRPC_SERVER

# 检查 controld 服务状态
docker compose -f build/compose.yaml ps controld

# 验证 gRPC 连通性
grpcurl -plaintext localhost:50051 list
```

## 2. 任务与调度

### 2.1 任务一直处于 READY 状态

**原因**：Slot 不足，无可用执行槽。

**解决**：
```bash
# 检查 slot 状态
scalebox slot list
scalebox slot list --host <hostname>

# 增加 slot
scalebox slot add --module <module-name> --host <hostname> --count 4

# 检查 actuator 日志
docker logs actuator
```

### 2.2 任务执行失败（exit code != 0）

**排查步骤**：
1. 查看任务日志：`scalebox task log <task-id>`
2. 检查模块镜像是否存在：`docker images | grep <image-name>`
3. 检查节点资源：`scalebox host show <hostname>`
4. 查看退出码含义（参见 :doc:`附录 6 - 退出码规范 <appendix/6_exit_code_spec>`）

### 2.3 批量任务重复创建

始终使用 `--conflict-action IGNORE` 避免重复。或设置 header `repeatable: yes` 允许重复分发。

## 3. 信号量与流控

### 3.1 信号量不存在错误

```bash
# 方式一：启用自动创建
SEMAPHORE_AUTO_CREATE=yes scalebox semaphore get my_sema

# 方式二：手动创建
scalebox semaphore create --app-id <app-id> my_sema 0
```

### 3.2 信号量值异常（不收敛）

检查流控版和可编程版是否匹配：
- 流控版（无前缀）由 controld 自动维护
- 可编程版（`:` 前缀）由脚本操作
- 参见 :doc:`VTask 信号量机制 <developer/9_vtask>` §9.4

## 4. WebUI

### 4.1 页面无法访问

```bash
# 检查 WebUI 服务
docker compose -f build/compose.yaml ps webui

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
docker compose -f build/compose.yaml logs > scalebox.log

# 按服务
docker compose -f build/compose.yaml logs controld
docker compose -f build/compose.yaml logs actuator

# 按时间
docker compose -f build/compose.yaml logs --since "2026-01-01"
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

> 持续更新中，欢迎贡献问题与答案。
