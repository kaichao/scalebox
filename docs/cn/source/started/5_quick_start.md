# 5. 快速上手

- 最小化环境部署

- Hello Scalebox示例

- 核心概念速览

本章将指导您快速部署和运行第一个Scalebox应用，让您在几分钟内体验Scalebox的核心功能。

## 5.1 环境准备

### 5.1.1 系统要求

- **操作系统**：Linux (CentOS 7+/Ubuntu 18.04+/Debian 10+)、macOS 10.15+
- **Docker**：20.10+（推荐使用最新稳定版）
- **内存**：≥ 4GB
- **磁盘空间**：≥ 10GB

### 5.1.2 Docker安装

如果您尚未安装Docker，请按以下步骤安装：

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y docker.io
sudo systemctl enable --now docker

# CentOS/RHEL
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl enable --now docker

# macOS
# 从 https://docs.docker.com/desktop/install/mac-install/ 下载Docker Desktop
```

验证Docker安装：
```bash
docker --version
docker run hello-world
```

## 5.2 获取Scalebox代码

```bash
# 克隆Scalebox仓库
git clone https://github.com/kaichao/scalebox.git
cd scalebox

# 或使用已存在的仓库
cd /path/to/your/scalebox
```

## 5.3 运行Hello Scalebox示例

### 5.3.1 启动单节点集群

```bash
# 进入hello-scalebox示例目录
cd examples/hello-scalebox

# 启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

预期输出：
```
      Name                    Command               State    Ports
------------------------------------------------------------------
hello-scalebox_controld_1   /controld --config-fi ...   Up      0.0.0.0:50051->50051/tcp
hello-scalebox_actuator_1   /actuator --config-fi ...   Up      
hello-scalebox_database_1   docker-entrypoint.sh ...    Up      0.0.0.0:5432->5432/tcp
```

### 5.3.2 创建并运行应用

```bash
# 创建应用
scalebox run

# 查看应用状态
scalebox app list

# 查看任务执行情况
scalebox task list --app hello-scalebox
```

### 5.3.3 查看运行结果

```bash
# 查看任务执行日志
scalebox task log <task_id>

# 或直接查看数据库中的执行记录
docker exec -it hello-scalebox_database_1 psql -U scalebox -d scalebox -c "SELECT * FROM t_task_exec LIMIT 5;"
```

## 5.4 验证安装

### 5.4.1 健康检查

```bash
# 检查各服务状态
scalebox cluster status

# 检查数据库连接
docker exec hello-scalebox_database_1 pg_isready -U scalebox

# 检查gRPC服务
grpcurl -plaintext localhost:50051 list
```

### 5.4.2 基本功能测试

1. **任务创建测试**：
   ```bash
   # 创建测试任务
   scalebox task create --module hello-module --body "test message"
   
   # 查看任务状态
   scalebox task list --status READY
   ```

2. **任务执行测试**：
   ```bash
   # 等待任务执行完成（约30秒）
   sleep 30
   
   # 查看执行结果
   scalebox task list --status COMPLETED
   ```

## 5.5 常见问题排查

### 5.5.1 Docker容器无法启动

**问题现象**：`docker-compose up` 失败

**解决方案**：
1. 检查Docker服务状态：`sudo systemctl status docker`
2. 检查端口占用：`netstat -tlnp | grep -E '(50051|5432)'`
3. 清理旧容器：`docker-compose down -v`

### 5.5.2 数据库连接失败

**问题现象**：应用创建失败，数据库连接错误

**解决方案**：
1. 等待数据库完全启动：`sleep 10`
2. 手动检查数据库：`docker exec hello-scalebox_database_1 psql -U scalebox -c "\l"`
3. 重新创建数据库卷：`docker-compose down -v && docker-compose up -d`

### 5.5.3 任务无法执行

**问题现象**：任务一直处于READY状态

**解决方案**：
1. 检查actuator日志：`docker logs hello-scalebox_actuator_1`
2. 检查slot状态：`scalebox slot list`
3. 重新启动actuator：`docker-compose restart actuator`

## 5.6 下一步

成功运行Hello Scalebox后，您可以：

1. **学习核心概念**：阅读第1-4章，了解Scalebox的架构和设计原理
2. **尝试更多示例**：探索其他示例应用（app-primes、app-copy等）
3. **开发自定义应用**：参考使用指南和编程指南，开发自己的Scalebox应用
4. **部署生产环境**：按照第2章指南部署多节点集群

## 5.7 故障反馈

如果遇到无法解决的问题：

1. 查看详细日志：设置`LOG_LEVEL=debug`环境变量后重新运行
2. 检查Issues：访问 [GitHub Issues](https://github.com/kaichao/scalebox/issues)
3. 提交问题：提供完整的错误日志和环境信息

---

**提示**：本章仅提供快速入门指南，详细配置和高级功能请参考后续章节。