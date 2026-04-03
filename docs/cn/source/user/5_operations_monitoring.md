# 5. 运维与监控

## 5.1 系统状态监控

## 5.2 日志管理与分析

## 5.3 性能调优指南

## 5.4 故障排查与恢复


本章介绍Scalebox平台的日常运维管理和监控方法，帮助管理员确保系统稳定运行。

## 5.1 系统监控

### 5.1.1 监控指标

Scalebox提供丰富的监控指标，用于评估系统运行状态：

#### 系统级指标
- **集群状态**：集群在线率、节点健康度
- **资源使用**：CPU、内存、磁盘、网络使用率
- **任务统计**：任务总数、运行中任务、已完成任务
- **队列深度**：各模块待处理任务数量

#### 应用级指标
- **应用吞吐量**：单位时间处理的任务数
- **任务成功率**：成功完成的任务比例
- **平均处理时间**：任务从创建到完成的平均耗时
- **错误率**：任务失败的比例及错误类型分布

### 5.1.2 监控工具

#### 命令行工具
```bash
# 查看集群整体状态
scalebox cluster status

# 查看节点状态
scalebox host list

# 查看slot状态
scalebox slot list --all

# 查看任务统计
scalebox task stats --period 1h

# 查看系统资源使用
scalebox cluster stats
```

#### 监控命令详解

**集群状态监控**：
```bash
# 查看集群详细状态
scalebox cluster status --detail

# 输出格式
scalebox cluster status --format json
scalebox cluster status --format table
```

**节点监控**：
```bash
# 查看所有节点状态
scalebox host list --status ALL

# 查看异常节点
scalebox host list --status ERROR

# 查看节点详细信息
scalebox host show <host_id>
```

**任务监控**：
```bash
# 查看任务统计
scalebox task stats --app <app_name> --period 24h

# 查看各状态任务数量
scalebox task count --status READY
scalebox task count --status RUNNING
scalebox task count --status COMPLETED
scalebox task count --status ERROR

# 查看任务执行时间分布
scalebox task stats --histogram --period 1d
```

## 5.2 日常运维

### 5.2.1 系统维护操作

#### 集群维护
```bash
# 添加新节点
scalebox cluster add-host <cluster_name> <host_ip> --uname <username> --port <port>

# 移除故障节点
scalebox cluster remove-host <host_id>

# 暂停集群
scalebox cluster pause <cluster_name>

# 恢复集群
scalebox cluster resume <cluster_name>

# 重启集群服务
scalebox cluster restart <cluster_name>
```

#### 应用维护
```bash
# 暂停应用
scalebox app pause <app_name>

# 恢复应用
scalebox app resume <app_name>

# 重新加载应用配置
scalebox app reload <app_name>

# 导出应用状态
scalebox app export <app_name> --output app_state.json
```

#### 模块维护
```bash
# 暂停模块
scalebox module pause <module_name>

# 恢复模块
scalebox module resume <module_name>

# 重新创建slot
scalebox slot delete --module <module_name> --all
scalebox app reload <app_name>
```

### 5.2.2 备份与恢复

#### 数据库备份
```bash
# 备份Scalebox数据库
docker exec scalebox_database_1 pg_dump -U scalebox scalebox > scalebox_backup_$(date +%Y%m%d).sql

# 自动备份脚本
#!/bin/bash
BACKUP_DIR=/backup/scalebox
DATE=$(date +%Y%m%d_%H%M%S)
docker exec scalebox_database_1 pg_dump -U scalebox scalebox > ${BACKUP_DIR}/scalebox_${DATE}.sql
gzip ${BACKUP_DIR}/scalebox_${DATE}.sql
# 保留最近7天的备份
find ${BACKUP_DIR} -name "scalebox_*.sql.gz" -mtime +7 -delete
```

#### 配置备份
```bash
# 备份应用配置
tar -czf scalebox_config_backup_$(date +%Y%m%d).tar.gz \
  /etc/scalebox/ \
  ${HOME}/.scalebox/ \
  apps/ \
  examples/
```

#### 数据恢复
```bash
# 恢复数据库
gunzip -c scalebox_backup_20240101.sql.gz | \
  docker exec -i scalebox_database_1 psql -U scalebox scalebox

# 恢复后验证
scalebox cluster status
scalebox app list
```

### 5.2.3 日志管理

#### 日志级别设置
```bash
# 设置日志级别（从低到高）
export LOG_LEVEL=error     # 仅错误
export LOG_LEVEL=warn      # 警告和错误
export LOG_LEVEL=info      # 信息、警告、错误（默认）
export LOG_LEVEL=debug     # 调试信息
export LOG_LEVEL=trace     # 最详细的跟踪信息
```

#### 日志文件位置
```
# 主要组件日志
/var/log/scalebox/controld.log     # controld服务日志
/var/log/scalebox/actuator.log     # actuator服务日志
/var/log/scalebox/database.log     # 数据库日志

# 任务执行日志
/var/log/scalebox/tasks/<app_name>/<module_name>/

# Docker容器日志
docker logs <container_name>
```

#### 日志轮转配置
```bash
# 使用logrotate管理日志
cat > /etc/logrotate.d/scalebox << EOF
/var/log/scalebox/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
    postrotate
        docker exec scalebox_controld_1 kill -USR1 1
        docker exec scalebox_actuator_1 kill -USR1 1
    endscript
}
EOF
```

## 5.3 性能调优

### 5.3.1 系统参数调优

#### 数据库参数优化
```sql
-- 调整PostgreSQL参数
ALTER SYSTEM SET shared_buffers = '4GB';
ALTER SYSTEM SET work_mem = '64MB';
ALTER SYSTEM SET maintenance_work_mem = '1GB';
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET effective_cache_size = '12GB';

-- 重启数据库使配置生效
docker restart scalebox_database_1
```

#### Scalebox服务参数
```yaml
# controld配置优化
controld:
  max_workers: 100           # 最大工作线程数
  grpc_max_recv_msg_size: 104857600  # 100MB消息限制
  task_batch_size: 100       # 批量任务处理大小
  
# actuator配置优化
actuator:
  max_concurrent_slots: 50   # 最大并发slot数
  heartbeat_interval: 30     # 心跳间隔（秒）
  slot_check_interval: 10    # slot检查间隔（秒）
```

### 5.3.2 应用性能调优

#### 模块配置优化
```yaml
modules:
  fast-processing:
    parameters:
      # 批量处理优化
      task_batch_size: 50            # 增大批量大小
      poll_interval_seconds: 2       # 减少轮询间隔
      
      # 资源优化
      task_max_seconds: 300          # 合理设置超时
      slot_timeout_minutes: 10       # 减少slot启动超时
      
      # 并发优化
      vtask_size: 100                # 增加虚拟任务并发数
```

#### Slot配置优化
```yaml
slots:
  # 根据节点性能调整slot数量
  - high_perf_nodes:8    # 高性能节点：8个slot
  - medium_perf_nodes:4  # 中等性能节点：4个slot
  - low_perf_nodes:2     # 低性能节点：2个slot
```

### 5.3.3 网络优化

#### 跨集群网络优化
```yaml
clusters:
  remote-cluster:
    parameters:
      # 使用专用网络通道
      network_interface: eth1
      
      # 优化传输参数
      tcp_keepalive: yes
      compression: gzip
      chunk_size: 65536      # 64KB分块传输
```

#### 本地网络优化
```bash
# 调整网络参数
sudo sysctl -w net.core.rmem_max=134217728
sudo sysctl -w net.core.wmem_max=134217728
sudo sysctl -w net.ipv4.tcp_rmem="4096 87380 134217728"
sudo sysctl -w net.ipv4.tcp_wmem="4096 65536 134217728"
```

## 5.4 故障处理

### 5.4.1 常见故障处理

#### 数据库连接失败
```bash
# 检查数据库状态
docker exec scalebox_database_1 pg_isready -U scalebox

# 重启数据库服务
docker restart scalebox_database_1

# 检查数据库日志
docker logs scalebox_database_1 --tail 100

# 修复数据库连接
scalebox database repair
```

#### Slot异常退出
```bash
# 查看异常slot
scalebox slot list --status ERROR

# 分析slot日志
scalebox slot log <slot_id>

# 重新创建slot
scalebox slot delete --id <slot_id>
scalebox app reload <app_name>

# 如果问题持续，检查模块配置
scalebox module show <module_name>
```

#### 任务堆积
```bash
# 查看任务队列
scalebox task list --status READY --limit 100

# 分析瓶颈模块
scalebox module stats --period 1h

# 临时增加slot数量
scalebox slot add --module <module_name> --count 5

# 调整任务分发策略
scalebox module update <module_name> --param task_batch_size=20
```

### 5.4.2 灾难恢复

#### 完整系统恢复流程
```bash
# 1. 停止所有服务
docker-compose down

# 2. 恢复数据库
gunzip -c latest_backup.sql.gz | docker exec -i scalebox_database_1 psql -U scalebox scalebox

# 3. 恢复配置文件
tar -xzf config_backup.tar.gz -C /

# 4. 启动服务
docker-compose up -d

# 5. 验证恢复
scalebox cluster status
scalebox app list
scalebox task count --all
```

#### 数据一致性检查
```bash
# 检查任务状态一致性
scalebox database check-consistency

# 修复不一致数据
scalebox database repair --fix-inconsistent

# 验证修复结果
scalebox database verify
```

## 5.5 安全运维

### 5.5.1 访问控制

#### 用户权限管理
```bash
# 创建管理员用户
scalebox user create admin --role administrator

# 创建普通用户
scalebox user create user1 --role operator

# 修改用户权限
scalebox user update user1 --role administrator

# 查看用户列表
scalebox user list

# 删除用户
scalebox user delete user1
```

#### API访问控制
```yaml
# 配置API访问策略
security:
  api_keys:
    - name: "monitoring_key"
      key: "secret_monitoring_key"
      permissions: ["read:cluster", "read:app", "read:task"]
      
    - name: "admin_key"
      key: "secret_admin_key"
      permissions: ["*"]
```

### 5.5.2 网络安全

#### 防火墙配置
```bash
# 只允许必要端口
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 50051/tcp   # gRPC
sudo ufw allow 5432/tcp    # PostgreSQL
sudo ufw default deny
sudo ufw enable
```

#### TLS加密配置
```yaml
# 启用TLS加密
tls:
  enabled: true
  cert_file: /etc/scalebox/cert.pem
  key_file: /etc/scalebox/key.pem
  ca_file: /etc/scalebox/ca.pem
```

## 5.6 监控仪表板

### 5.6.1 内置监控界面

Scalebox提供内置的Web监控界面：

```bash
# 启动监控界面（如果已安装）
scalebox monitor start

# 访问监控界面
# 默认地址：http://localhost:3000
```

### 5.6.2 第三方监控集成

#### Prometheus集成
```yaml
# Scalebox Prometheus exporter配置
prometheus:
  enabled: true
  port: 9091
  metrics_path: /metrics
  
  # 监控指标
  collect:
    - cluster_metrics
    - task_metrics
    - slot_metrics
    - system_metrics
```

#### Grafana仪表板
```bash
# 导入Scalebox Grafana仪表板
# 仪表板ID：12345（示例）
# 可以从Scalebox GitHub仓库获取预配置仪表板
```

### 5.6.3 自定义监控脚本

#### 健康检查脚本
```bash
#!/bin/bash
# scalebox_health_check.sh

# 检查服务状态
check_service() {
    service=$1
    if docker ps | grep -q $service; then
        echo "✓ $service is running"
        return 0
    else
        echo "✗ $service is not running"
        return 1
    fi
}

# 检查数据库连接
check_database() {
    if docker exec scalebox_database_1 pg_isready -U scalebox > /dev/null 2>&1; then
        echo "✓ Database is reachable"
        return 0
    else
        echo "✗ Database is not reachable"
        return 1
    fi
}

# 检查任务处理
check_tasks() {
    stuck_tasks=$(scalebox task list --status RUNNING --older-than 1h | wc -l)
    if [ $stuck_tasks -gt 10 ]; then
        echo "⚠ $stuck_tasks tasks stuck for more than 1 hour"
        return 1
    else
        echo "✓ Task processing is normal"
        return 0
    fi
}

# 执行所有检查
check_service scalebox_controld_1
check_service scalebox_actuator_1
check_database
check_tasks
```

## 5.7 最佳实践

### 5.7.1 运维最佳实践

1. **定期备份**：每天备份数据库和关键配置
2. **监控告警**：设置关键指标告警阈值
3. **容量规划**：定期评估系统容量，提前扩容
4. **变更管理**：所有配置变更先测试后上线
5. **文档更新**：系统变更时同步更新文档

### 5.7.2 性能最佳实践

1. **合理分配资源**：根据任务类型调整slot数量和配置
2. **批量处理**：合理设置task_batch_size减少数据库压力
3. **本地化计算**：尽量使用HOST-BOUND模式减少网络开销
4. **定期优化**：定期分析慢查询和性能瓶颈

### 5.7.3 安全最佳实践

1. **最小权限原则**：用户只授予必要权限
2. **定期审计**：定期检查访问日志和操作记录
3. **及时更新**：保持系统和依赖组件的最新版本
4. **加密传输**：生产环境启用TLS加密

## 5.8 故障演练

### 5.8.1 演练场景

定期进行以下故障演练：
- 数据库故障恢复
- 节点宕机处理
- 网络分区恢复
- 应用配置错误恢复

### 5.8.2 演练记录

记录每次演练的：
- 故障现象
- 处理步骤
- 恢复时间
- 经验教训
- 改进措施

---

通过本章的学习，您应该能够：
1. 监控Scalebox系统运行状态
2. 进行日常维护和故障处理
3. 优化系统性能
4. 确保系统安全稳定运行
5. 建立完善的运维流程

**提示**：运维工作需要持续学习和实践，建议定期回顾本章内容，并根据实际运行情况调整运维策略。