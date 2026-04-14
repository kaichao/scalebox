# 5. 运维管理

## 5.1 系统监控

### 5.1.1 服务状态监控

```bash
# 检查服务状态
docker-compose ps

# 查看服务日志
docker-compose logs controld
docker-compose logs actuator
docker-compose logs database

# 检查服务健康
scalebox cluster status
```

### 5.1.2 资源使用监控

```bash
# 查看节点资源
scalebox node metrics

# 查看存储使用
docker system df

# 查看容器资源使用
docker stats
```

## 5.2 故障排查

### 5.2.1 数据库问题

**问题现象**：数据库连接失败

**解决方案**：
1. 等待数据库完全启动：`sleep 10`
2. 手动检查数据库：`docker exec database psql -U scalebox -c "\l"`
3. 重新创建数据库卷：`docker-compose down -v && docker-compose up -d`

### 5.2.2 任务执行问题

**问题现象**：任务一直处于READY状态

**解决方案**：
1. 检查actuator日志：`docker logs actuator`
2. 检查slot状态：`scalebox slot list`
3. 重新启动actuator：`docker-compose restart actuator`

### 5.2.3 网络问题

**问题现象**：节点间通信失败

**解决方案**：
1. 检查网络连通性：`ping <node-ip>`
2. 检查端口开放：`netstat -tlnp | grep 50051`
3. 检查防火墙：`ufw status`

## 5.3 性能优化

### 5.3.1 系统参数调优

```bash
# 调整Docker参数
dockerd --max-concurrent-downloads 10 --max-concurrent-uploads 10

# 调整系统参数
sysctl -w vm.max_map_count=262144
sysctl -w fs.file-max=100000
```

### 5.3.2 应用性能优化

1. **合理设置并行度**
   - 根据数据量设置模块并行度
   - 避免过度并行导致资源竞争

2. **优化数据局部性**
   - 优先使用本地存储
   - 减少跨节点数据传输

3. **批量处理优化**
   - 合理设置批量大小
   - 平衡延迟和吞吐量

## 5.4 备份与恢复

### 5.4.1 数据库备份

```bash
# 备份数据库
docker exec database pg_dump -U scalebox scalebox > backup.sql

# 恢复数据库
docker exec -i database psql -U scalebox scalebox < backup.sql
```

### 5.4.2 配置备份

```bash
# 备份配置文件
tar -czf scalebox-config.tar.gz examples/ dockerfiles/ runtime/

# 恢复配置文件
tar -xzf scalebox-config.tar.gz
```

## 5.5 安全维护

### 5.5.1 访问控制

```bash
# 更新密钥
cd runtime && make update-pubkey

# 重置用户密码
docker exec database psql -U scalebox -c "ALTER USER scalebox WITH PASSWORD 'newpassword';"
```

### 5.5.2 安全更新

```bash
# 更新容器镜像
cd runtime && make pull-all

# 重启服务
docker-compose down && docker-compose up -d
```

## 5.6 日志管理

### 5.6.1 日志收集

```bash
# 收集所有服务日志
docker-compose logs > scalebox.log

# 按时间筛选日志
docker-compose logs --since "2024-01-01" --until "2024-01-02"
```

### 5.6.2 日志分析

```bash
# 查找错误日志
grep -i error scalebox.log

# 统计任务执行时间
grep "task completed" scalebox.log | awk '{print $NF}' | sort -n
```

## 5.7 系统升级

### 5.7.1 版本升级

```bash
# 拉取新版本
git pull origin main

# 更新容器镜像
make pull-all

# 平滑重启服务
docker-compose restart
```

### 5.7.2 配置迁移

1. 备份当前配置
2. 更新配置文件
3. 验证配置有效性
4. 重启服务

## 5.8 常见问题速查

| 问题现象 | 可能原因 | 解决方案 |
| -------- | -------- | -------- |
| 数据库连接失败 | 数据库未启动 | 等待数据库启动完成 |
| 任务无法执行 | Slot不足 | 增加计算节点或槽位 |
| 网络连接失败 | 防火墙阻止 | 开放相应端口 |
| 磁盘空间不足 | 日志或数据积累 | 清理不必要的文件 |
| 内存不足 | 并行度过高 | 降低并行度或增加内存 |