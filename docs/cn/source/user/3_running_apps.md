# 3. 运行应用

## 3.1 应用创建与运行

### 3.1.1 应用定义文件

Scalebox应用通过YAML配置文件定义。主要包括：

- **应用名称**：唯一标识符，支持变量替换
- **集群配置**：目标集群名称
- **模块定义**：包含的模块及其配置
- **参数设置**：运行时参数

### 3.1.2 应用参数文件

应用参数通过环境变量文件定义，支持变量替换。

### 3.1.3 创建应用

```bash
# 使用默认配置文件
scalebox run

# 指定配置文件
scalebox run --env-file scalebox.env --app-file app.yaml

# 带参数创建
scalebox run --tag v1 --num-groups 4
```

## 3.2 应用管理

### 3.2.1 查看应用状态

```bash
# 列出所有应用
scalebox app list

# 查看特定应用
scalebox app show <app-name>

# 查看应用详细信息
scalebox app describe <app-name>
```

### 3.2.2 应用操作

```bash
# 停止应用
scalebox app stop <app-name>

# 删除应用
scalebox app delete <app-name>

# 查看应用日志
scalebox app logs <app-name>
```

## 3.3 任务管理

### 3.3.1 查看任务状态

若不指定app-id，则缺省值为最新的app-id
```bash
# 列出应用的任务
scalebox task list --app-id <app-id>

# 按状态过滤任务
scalebox task list --app-id <app-id> --status READY
scalebox task list --app-id <app-id> --status COMPLETED
scalebox task list --app-id <app-id> --status FAILED
```

### 3.3.2 任务操作

```bash
# 查看任务详情
scalebox task show <task-id>

# 查看任务日志
scalebox task log <task-id>

# 重试失败任务
scalebox task retry <task-id>

# 删除任务
scalebox task delete <task-id>
```

## 3.4 集群管理

### 3.4.1 查看集群状态

```bash
# 查看集群状态
scalebox cluster status

# 列出集群节点
scalebox node list

# 查看节点详情
scalebox node show <node-name>
```

### 3.4.2 运行槽管理

```bash
# 列出槽位
scalebox slot list

# 查看槽位详情
scalebox slot show <slot-id>
```

## 3.5 应用监控

### 3.5.1 状态监控

```bash
# 实时监控应用状态
scalebox app status <app-name> --watch

# 监控任务执行
scalebox task monitor --app <app-name>
```

### 3.5.2 性能指标

```bash
# 查看应用性能
scalebox app metrics <app-name>

# 查看节点资源使用
scalebox node metrics <node-name>
```

## 3.6 故障排查

### 3.6.1 常见问题

1. **应用创建失败**
   - 检查配置文件格式
   - 验证参数设置
   - 查看错误日志

2. **任务执行失败**
   - 检查任务日志
   - 验证模块镜像
   - 检查资源可用性

3. **集群连接问题**
   - 验证网络连通性
   - 检查服务状态
   - 验证认证配置

### 3.6.2 调试工具

```bash
# 验证配置
scalebox validate app.yaml

# 测试模块
scalebox test module <module-name>

# 检查系统状态
scalebox system check
```