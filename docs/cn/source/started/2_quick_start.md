# 2. 快速上手

本章将指导您快速部署和运行第一个Scalebox应用，让您在几分钟内体验Scalebox的核心功能。

## 2.1 环境准备

### 2.1.1 系统要求

- **操作系统**：Linux (CentOS 7+/Ubuntu 18.04+/Debian 10+)、macOS 10.15+
- **Docker**：20.10+（推荐使用最新稳定版）
- **内存**：≥ 4GB
- **磁盘空间**：≥ 30GB

### 2.1.2 Docker安装

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
```

- 验证Docker安装：
```bash
docker --version
docker run hello-world
```

## 2.2 获取Scalebox代码

```bash
git clone https://github.com/kaichao/scalebox.git
cd scalebox
```

## 2.3 运行Hello Scalebox示例

### 2.3.1 启动单节点集群

```bash
make -C runtime all
```

查看服务状态：
```bash
docker ps
```

### 2.3.2 创建并运行应用

```bash
# 创建应用
cd examples/hello-scalebox

echo "Docker-based_Scalebox" | scalebox run 

# 查看应用状态
scalebox app list

# 查看任务执行情况
scalebox task list
```

### 2.3.3 查看运行结果

```bash
# 查看任务执行日志
scalebox task log <task_id>
```

## 2.4 验证安装

### 2.4.1 健康检查

```bash
# 检查各服务状态
scalebox cluster status
```

## 2.5 下一步

成功运行Hello Scalebox后，您可以：

1. **学习核心概念**：阅读用户指南了解App、Module、Task等核心概念
2. **尝试更多示例**：探索其他示例应用（app-primes、app-copy等）
3. **开发自定义应用**：参考开发指南，开发自己的Scalebox应用