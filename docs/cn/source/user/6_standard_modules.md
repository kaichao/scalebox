# 6. 标准模块详解

## 6.1 基础模块（agent）

### 6.1.1 模块功能
- 提供了模块底座功能
- 应用模块的中任务执行的生命周期管理
- 与server端runtime的controld交互，纪录任务执行状态

### 6.1.2 模块特性
- **标准化接口**：提供统一的任务执行接口
- **状态管理**：管理任务的生命周期状态
- **错误处理**：支持任务失败重试和错误报告
- **资源管理**：管理计算资源和执行环境

### 6.1.3 配置参数
```yaml
modules:
  my-module:
    base_image: scalebox.net/platform/agent:latest
    arguments:
      code_path: ${PWD}/code
      task_max_seconds: 600
    parameters:
      task_dist_mode: DEFAULT
```

### 6.1.4 非agent模块的构建
```dockerfile
FROM scalebox.net/platform/agent:latest
COPY --from=scalebox.net/platform/agent:latest /usr/local/ /usr/local/
COPY ./code /app/bin/
```

## 6.2 文件传输模块

### 6.2.1 通用配置
文件传输模块通过rsync-over-ssh实现，支持跨节点数据传输。

#### 任务体格式
- **task-body**: 待传输到远端的本地相对路径

#### 任务头参数
| 任务头       | 缺省值 | 说明                       |
| ---------- | ----- | ------------------------- |
| source_url |       | "/local/dir"              |
| target_url |       | "user@remote-ip:remote-port/remote/dir" |
| keep_source | "yes" | "yes" / "no"，是否保留源端文件/目录 |

### 6.2.2 file-copy 模块
用于文件拷贝操作。

#### file-copy 使用示例
```yaml
modules:
  file-copy:
    base_image: scalebox.net/platform/file-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

#### file-copy 配置参数
- **source_url**: 源文件路径
- **target_url**: 目标文件路径
- **keep_source**: 是否保留源文件

### 6.2.3 dir-copy 模块
用于目录拷贝操作。

#### dir-copy 使用示例
```yaml
modules:
  dir-copy:
    base_image: scalebox.net/platform/dir-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

### 6.2.4 rsync-copy模块
基于rsync的高效数据传输模块。

#### rsync-copy 使用示例
```yaml
modules:
  rsync-copy:
    base_image: scalebox.net/platform/rsync-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

#### rsync-copy 高级参数
- **rsync_options**: rsync命令行选项
- **compress**: 是否启用压缩传输
- **delete**: 是否删除目标端多余文件

### 6.2.5 ftp-copy模块
支持FTP协议的数据传输模块。

#### ftp-copy 使用示例
```yaml
modules:
  ftp-copy:
    base_image: scalebox.net/platform/ftp-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

#### ftp-copy 配置参数
- **ftp_server**: FTP服务器地址
- **ftp_user**: FTP用户名
- **ftp_password**: FTP密码
- **ftp_port**: FTP端口（默认21）

## 6.3 辅助功能模块

### 6.3.1 cron模块
定时任务模块，支持按计划执行任务。

#### cron 使用示例
```yaml
modules:
  cron:
    base_image: scalebox.net/platform/cron
    arguments:
      cron_expression: "0 2 * * *"
      command: "/app/bin/backup.sh"
```

#### cron 配置参数
- **cron_expression**: Cron表达式（分 时 日 月 周）
- **command**: 要执行的命令
- **timezone**: 时区设置（默认UTC）

### 6.3.2 cluster-head模块
集群头节点管理模块。

#### cluster-head 使用示例
```yaml
modules:
  cluster-head:
    base_image: scalebox.net/platform/cluster-head
    parameters:
      slot_options: slot_on_head
```

#### cluster-head 功能特性
- 集群管理协调
- 资源调度优化
- 节点状态监控

### 6.3.3 node-agent模块
节点代理模块，运行在计算节点上。

#### node-agent 使用示例
```yaml
modules:
  node-agent:
    base_image: scalebox.net/platform/node-agent
    parameters:
      cluster: ${CLUSTER}
```

#### node-agent 功能特性
- 节点资源管理
- 任务执行监控
- 系统状态报告

### 6.3.4 dir-list模块
目录列表模块，用于生成待处理文件列表。

#### dir-list 使用示例
```yaml
modules:
  dir-list:
    base_image: scalebox.net/platform/dir-list
    arguments:
      target_dir: /data/input
      pattern: "*.txt"
```

#### dir-list 配置参数
- **target_dir**: 目标目录路径
- **pattern**: 文件匹配模式
- **recursive**: 是否递归扫描子目录
- **output_format**: 输出格式（json/csv）

## 6.4 模块配置最佳实践

### 6.4.1 模块选择原则
1. **优先使用标准模块**：减少重复开发工作
2. **考虑性能需求**：选择合适的传输协议和算法
3. **支持可扩展性**：便于后续功能扩展

### 6.4.2 参数配置建议
1. **合理设置超时**：根据任务复杂度设置task_max_seconds
2. **优化并行度**：根据资源情况设置slot数量
3. **错误处理**：配置合适的重试策略

### 6.4.3 模块组合模式
1. **流水线模式**：多个模块顺序执行
2. **并行模式**：多个模块并行处理
3. **混合模式**：结合流水线和并行处理

## 6.5 模块使用示例

### 6.5.1 完整的数据处理流水线
```yaml
modules:
  # 生成文件列表
  dir-list:
    base_image: scalebox.net/platform/dir-list
    arguments:
      target_dir: /data/input
      pattern: "*.dat"
  
  # 数据传输
  file-copy:
    base_image: scalebox.net/platform/file-copy
    parameters:
      task_dist_mode: SLOT-BOUND
  
  # 数据处理
  data-process:
    base_image: my-registry/data-process:latest
    parameters:
      task_dist_mode: HOST-BOUND
  
  # 结果传输
  rsync-copy:
    base_image: scalebox.net/platform/rsync-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

### 6.5.2 定时备份任务
```yaml
modules:
  # 定时触发
  cron:
    base_image: scalebox.net/platform/cron
    arguments:
      cron_expression: "0 2 * * *"
      command: "/app/bin/backup.sh"
  
  # 数据备份
  rsync-copy:
    base_image: scalebox.net/platform/rsync-copy
    parameters:
      task_dist_mode: SLOT-BOUND
      compress: yes
```

## 6.6 模块开发建议

### 6.6.1 模块设计原则
1. **单一职责**：每个模块只负责一个特定功能
2. **无状态设计**：支持任务的重复执行
3. **标准化接口**：遵循Scalebox模块规范
4. **错误处理**：完善的错误报告和恢复机制

### 6.6.2 性能优化
1. **批量处理**：支持批量任务处理提高效率
2. **缓存机制**：合理使用缓存减少重复计算
3. **并行优化**：充分利用多核和分布式计算
4. **I/O优化**：优化数据读写性能

### 6.6.3 测试建议
1. **单元测试**：测试模块的基本功能
2. **集成测试**：测试模块间的协作
3. **性能测试**：验证模块的性能指标
4. **容错测试**：测试模块的容错能力