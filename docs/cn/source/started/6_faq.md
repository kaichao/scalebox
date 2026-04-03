# 6. 常见问题（FAQ）

- 安装部署常见问题

- 使用过程中的疑难解答

- 最佳实践建议

## 6.1 安装与部署问题

### 6.1.1 本地集群的hello-scalebox应用不能正常运行，如何检查？

**问题现象**：Hello Scalebox示例应用无法启动或运行失败。

**排查步骤**：

1. **检查SSH端口**：
   ```bash
   # 检查本地节点的sshd端口号是否为缺省的22
   sudo netstat -tlnp | grep :22
   ```
   
   如果不是默认端口22，请在 `server/local/mycluster.yaml` 文件中设置端口号：
   ```yaml
   parameters:
     port: <port_number>
   ```

2. **检查本地IP地址**：
   ```bash
   # 检查本地IP地址是否正确设置
   hostname -i
   hostname -I
   ```
   
   如果 `hostname -i` 返回不正确，使用 `hostname -I` 列出所有本地IP地址，然后在 `defs.mk` 中设置变量：
   ```
   LOCAL_ADDR=<your_ip_address>
   # 或
   LOCAL_IP_INDEX=<index_number>
   ```

3. **检查Docker容器状态**：
   ```bash
   cd examples/hello-scalebox
   docker-compose ps
   docker-compose logs
   ```

### 6.1.2 Docker容器无法启动

**问题现象**：`docker-compose up` 命令执行失败。

**解决方案**：
1. 检查Docker服务状态：
   ```bash
   sudo systemctl status docker
   sudo systemctl restart docker
   ```

2. 检查端口占用：
   ```bash
   netstat -tlnp | grep -E '(50051|5432)'
   ```

3. 清理旧容器和数据：
   ```bash
   docker-compose down -v
   docker system prune -a
   ```

4. 检查磁盘空间：
   ```bash
   df -h
   ```

### 6.1.3 数据库连接失败

**问题现象**：应用创建失败，显示数据库连接错误。

**解决方案**：
1. 等待数据库完全启动：
   ```bash
   sleep 15
   docker-compose logs database
   ```

2. 手动检查数据库连接：
   ```bash
   docker exec -it hello-scalebox_database_1 psql -U scalebox -c "\l"
   ```

3. 重新创建数据库卷：
   ```bash
   docker-compose down -v
   docker-compose up -d database
   sleep 30
   docker-compose up -d
   ```

## 6.2 程序执行与排错

### 6.2.1 程序执行中如何排错（Debug）？

**方法1：启用调试日志**
```bash
# scalebox命令启动前，设定环境变量
export LOG_LEVEL=debug  # 或 trace

# 在agent脚本中增加环境变量
export PLAT_LOG_LEVEL=debug
```

**方法2：检查组件日志**
```bash
# 查看actuator日志
docker logs hello-scalebox_actuator_1

# 查看controld日志
docker logs hello-scalebox_controld_1

# 查看数据库日志
docker logs hello-scalebox_database_1
```

**方法3：使用scalebox命令行工具**
```bash
# 查看集群状态
scalebox cluster status

# 查看slot状态
scalebox slot list

# 查看任务状态
scalebox task list --all --limit 20

# 查看任务详情
scalebox task show <task_id>
```

**方法4：启用任务详细输出**
```bash
# 在模块定义中添加详细输出参数
parameters:
  output_text_size: 1048576  # 1MB
  text_trunc_mode: "TAIL"    # 保留尾部输出
```

### 6.2.2 任务无法执行或一直处于READY状态

**问题现象**：任务创建后一直处于READY状态，不进入RUNNING状态。

**排查步骤**：

1. **检查slot状态**：
   ```bash
   scalebox slot list
   ```
   
   确认有可用的slot处于READY状态。

2. **检查actuator状态**：
   ```bash
   docker logs hello-scalebox_actuator_1 --tail 100
   ```

3. **检查网络连接**：
   ```bash
   # 检查grpc连接
   grpcurl -plaintext localhost:50051 list
   
   # 检查数据库连接
   docker exec hello-scalebox_database_1 pg_isready -U scalebox
   ```

4. **检查资源配置**：
   ```bash
   # 检查是否有足够的系统资源
   free -h
   df -h /
   
   # 检查Docker资源限制
   docker info | grep -A5 "Resources"
   ```

**解决方案**：
1. 重启actuator：
   ```bash
   docker-compose restart actuator
   ```

2. 重新创建slot：
   ```bash
   scalebox slot delete --module <module_name> --all
   scalebox app reload
   ```

3. 检查任务分发模式：
   ```yaml
   # 确认模块的task_dist_mode设置正确
   parameters:
     task_dist_mode: "SLOT-BOUND"  # 或 "HOST-BOUND"
   ```

### 6.2.3 任务执行超时

**问题现象**：任务执行时间过长，最终超时失败。

**解决方案**：
1. 调整任务超时设置：
   ```yaml
   parameters:
     task_max_seconds: 3600  # 1小时
     task_global_timeout_scale: 3.0  # 全局超时为task_max_seconds的3倍
   ```

2. 检查用户程序性能：
   ```bash
   # 在容器内直接运行程序测试
   docker run -it <module_image> /app/bin/run.sh test_input
   ```

3. 优化资源分配：
   ```yaml
   # 增加slot资源限制
   slots:
     - ${nodes}:2  # 每个节点2个slot
   ```

## 6.3 性能优化问题

### 6.3.1 如何提高任务处理速度？

**优化建议**：

1. **批量处理任务**：
   ```yaml
   parameters:
     task_batch_size: 10  # 单批次处理10个任务
   ```

2. **优化slot配置**：
   ```yaml
   # 增加并发slot数量
   slots:
     - ${nodes}:4  # 每个节点4个slot
   ```

3. **本地计算优化**：
   ```yaml
   parameters:
     task_dist_mode: "HOST-BOUND"  # 主机绑定模式
   ```

4. **内存和存储优化**：
   ```yaml
   parameters:
     slot_options:
       - tmpfs_workdir  # 使用tmpfs文件系统
   ```

### 6.3.2 如何减少网络传输开销？

**优化建议**：

1. **数据本地化**：
   ```yaml
   # 使用本地存储路径
   volumes:
     - /local/data/path:/data
   ```

2. **压缩传输数据**：
   ```bash
   # 在模块命令中添加数据压缩
   command: docker run ... --compress
   ```

3. **减少消息大小**：
   ```yaml
   # 优化消息格式
   parameters:
     output_text_size: 65536  # 限制输出大小
   ```

## 6.4 跨集群部署问题

### 6.4.1 跨集群通信失败

**问题现象**：跨集群任务无法正常执行。

**排查步骤**：

1. **检查网络连通性**：
   ```bash
   # 测试集群间网络连接
   ping <remote_cluster_ip>
   nc -zv <remote_cluster_ip> 50051
   ```

2. **检查集群配置**：
   ```yaml
   clusters:
     remote-cluster:
       parameters:
         grpc_server: <remote_ip>:50051
         remote_grpc_server: <local_ip>:50051
   ```

3. **检查防火墙设置**：
   ```bash
   # 确认端口开放
   sudo ufw status
   sudo iptables -L -n
   ```

### 6.4.2 跨集群数据同步问题

**解决方案**：
1. 使用共享存储：
   ```yaml
   volumes:
     - /nfs/shared/data:/shared_data
   ```

2. 配置数据复制：
   ```yaml
   # 使用数据传输模块
   modules:
     data-transfer:
       base_image: scalebox/data-transfer
   ```

## 6.5 高级功能问题

### 6.5.1 如何使用信号量和共享变量？

**示例代码**：
```bash
# 创建信号量
scalebox semaphore create my_semaphore

# 信号量操作
scalebox semaphore increment my_semaphore
scalebox semaphore decrement my_semaphore
scalebox semaphore get my_semaphore

# 创建共享变量
scalebox variable create my_variable

# 变量操作
scalebox variable set my_variable "value"
scalebox variable get my_variable
```

### 6.5.2 如何实现任务级容错？

**配置示例**：
```yaml
parameters:
  retry_rules: '[{"0":0,"*":3}]'  # 退出码0不重试，其他退出码重试3次
  slot_recoverable: "yes"
  slot_max_retries: 2
```

## 6.6 监控与日志问题

### 6.6.1 如何查看详细的运行日志？

**方法**：
```bash
# 启用详细日志
export LOG_LEVEL=trace

# 查看组件日志
docker-compose logs -f --tail 100

# 查看任务执行日志
scalebox task log <task_id> --full

# 导出任务执行记录
scalebox task export --status COMPLETED --limit 1000 > tasks.csv
```

### 6.6.2 如何监控系统性能？

**监控工具**：
```bash
# 查看系统资源使用
scalebox cluster stats

# 查看任务统计
scalebox task stats --period 1h

# 导出性能数据
scalebox metrics export --start "2024-01-01" --end "2024-01-02"
```

## 6.7 其他常见问题

### 6.7.1 Scalebox支持哪些容器引擎？

**支持的引擎**：
- Docker（推荐）
- Singularity
- Apptainer
- Podman（实验性支持）

**配置方法**：
```yaml
parameters:
  default_runtime: "docker"  # 或 "singularity"
```

### 6.7.2 如何备份和恢复Scalebox数据？

**备份**：
```bash
# 备份数据库
docker exec hello-scalebox_database_1 pg_dump -U scalebox scalebox > backup.sql

# 备份配置
tar -czf scalebox-backup.tar.gz examples/ config/
```

**恢复**：
```bash
# 恢复数据库
cat backup.sql | docker exec -i hello-scalebox_database_1 psql -U scalebox scalebox

# 恢复配置
tar -xzf scalebox-backup.tar.gz
```

---

## 6.8 获取更多帮助

如果上述解决方案无法解决您的问题，请：

1. **查看官方文档**：访问 [Scalebox文档网站](https://scalebox.readthedocs.io/)
2. **检查GitHub Issues**：查看 [GitHub Issues](https://github.com/kaichao/scalebox/issues) 是否有类似问题
3. **提交新问题**：提供完整的错误日志、环境信息和复现步骤
4. **社区讨论**：加入Scalebox用户社区讨论

**报告问题时请提供**：
- Scalebox版本号
- 操作系统和Docker版本
- 完整的错误日志
- 应用配置文件（脱敏后）
- 复现步骤

---


## 1. local集群的hello-scalebox应用不能正常运行，如何检查？
- 检查本地节点的sshd的端口号是否为缺省的22，若不是，则在server/local/mycluster.yaml文件中设置该端口号，具体在parameters中加上`port=<port_number>`
- 用```hostname -i```检查本地ip地址是否正确设置？若不正确，则用```hostname -I```列出所有本地IP地址，再在defs.mk中设置变量```LOCAL_ADDR```或```LOCAL_IP_INDEX```

## 2. 程序执行中如何排错（Debug）？
- scalebox命令启动前，设定环境变量为LOG_LEVEL=debug或trace
- agent脚本中，增加环境变量，设定值为debug或trace
- actuator/controld的启动文件中， 设定环境变量为LOG_LEVEL=debug或trace
