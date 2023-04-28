# Scalebox集群

Scalebox集群是运行scalebox应用的硬件资源集合。分为
- 内联集群（Inline Cluster）：由Scalebox直接管理的计算资源组成的集群。支持静态资源分配、动态资源分配两种方式。
- 外部集群（External Cluster）：有外部调度系统管理的计算资源组成的集群，可支持用slurm、k8s等。

## 一、Scalebox集群介绍
### 1.1 头节点
头节点上安装了单个scalebox集群的管理服务，主要包括：
- actuator：启动端，负责在计算节点上启动slot；针对内联集群，直接通过免密ssh启动；针对外部集群，通过调用集群调度系统来启动。
- controld：面向actuator、计算节点，提供基于grpc的控制端应用服务；
- database：基于postgresql的数据库，存放app、job、task、slot等相关数据，面向controld等提供数据存储、检索等服务；也可面向计算节点、命令行工具提供直接数据库访问。

头节点以及头节点服务（controld/actuator/database），可为多个集群所共享。

### 1.2 计算节点
计算节点分为两类：
- 内部计算节点：scalebox内部调度管理的计算节点，通过免密ssh启动slot；
- 外部计算节点：通过外部调度程序（slurm/k8s等）启动的计算节点，在节点上启动slot。

## 二、Scalebox单节点集群安装及配置
- 操作系统：
  - CentOS 7以上（其他版本的Linux待测试）
  - macos 10.15(amd64)以上（ARM版的macos待测试）
- 容器化引擎：DockerCE 20.10+(rootless)
- docker-compose: 1.29.2+
- 安装dstat、htop、rsync、zstd、gmake、git、rsync等工具软件，用于性能监控等

### 2.1 安装CentOS 7/8下基本软件
- 以root用户安装
```bash
yum install -y epel-release
yum install -y htop dstat rsync pv pdsh wget make

```
- 安装git v2


- sshd
  - turn off UseDNS/GSSAPIAuthentication

- Setup Linux Time Calibration

### 2.2 安装macos下基本软件

启动本机sshd

基于Homebrew，安装所需基本软件

### 2.3 以当前用户身份，下载安装docker (rootless)
- 参考文档

[Run the Docker daemon as a non-root user (Rootless mode)](https://docs.docker.com/engine/security/rootless/)

- 设置环境变量
```bash
cat >> ${HOME}/.bashrc << EOF
export XDG_RUNTIME_DIR=${HOME}/.docker/run
export PATH=${HOME}/bin:$PATH
export DOCKER_HOST=unix://${HOME}/.docker/run/docker.sock
EOF

source ~/.bashrc
```
- 安装rootless docker
```bash
curl -fsSL https://get.docker.com/rootless | sh
```

- 启动rootless docker
```bash
nohup dockerd-rootless.sh &
```

- 验证docker有效
```bash
docker run --rm hello-world
```
### 2.4 以当前用户身份，下载安装docker-scalebox
- 安装docker-compose
```bash
mkdir -p ~/bin
wget -O ~/bin/docker-compose https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64
chmod +x ~/bin/docker-compose
```

- 下载docker-scalebox

在${HOME}下，下载docker-scalebox

```bash
cd && git clone https://github.com/kaichao/docker-scalebox
```

### 2.5 配置docker-scalebox

#### 准备运行环境
- 创建运行相关目录
- 设置环境变量
- 设置内联集群内的免密ssh

```bash
make prepare
```

#### 获取最新的命令行工具
```bash
make get-cli
```
#### 配置本地IP地址
- 通过以下命令，检查本地IP地址是否需要设定
```bash
hostname -i
```
- 前一步的本地IP地址不正确，则需要设置defs.mk文件中LOCAL_IP_INDEX或LOCAL_ADDR变量

可通过以下命令```hostname -I```找到正确的本地IP地址。
- 示例
```
[user@node-1 local]$ hostname -I
10.0.2.21 192.168.56.21 172.17.0.1 172.20.0.1 172.19.0.1 172.22.0.1 

则可以设置为：
LOCAL_IP_INDEX=2
或
LOCAL_ADDR=192.168.56.21
```

#### 重新启动shell，使得配置生效


### 2.6 启动scalebox控制端
```bash
cd ~/docker-scalebox/cluster && make all
```

## 三、Scalebox多节点集群安装及配置

### 3.1 头节点安装

在单集群基础上，安装time-server、pdsh、pv、rsync
```sh
yum install -y pdsh pdsh-rcmd-ssh
export PDSH_RCMD_TYPE=ssh
```

### 3.2 设置集群的共享存储
设置可用于计算节点间共享的存储，可以是基于NFS或glusterfs等。

并设置defs.mk中的```SHARED_DIR```变量

### 3.3 内部计算节点安装
- 操作系统：CentOS 7以上
- 容器化引擎：Docker 20.10以上
- 安装dstat、htop、zstd等工具软件，用于性能监控、数据压缩等
- 
可选：
- Linux Time Calibration
- 增加头节点管理用户、actuator的公钥
- dstat
- htop
- zstd
- sshd 
  - turn off UseDNS GSSAPIAuthentication

```sh
yum install -y epel-release
yum install -y htop dstat pv
```

- 参照目录 ```bio-down```下的```mycluster.yaml```，定义计算节点配置文件，并更新```Makefile```中的```clusters```变量

### 3.4 外部计算节点
- 操作系统：CentOS 7以上
- 单机容器化引擎：
  - Docker 20.10以上
  - podman ？版本；
  - singularity 3.8；
- k8s集群

### 3.5 集群定义文件

- 参照bio-down目录，完成集群定义文件mycluster.yaml
- 将新集群名称，加入到Makefile文件的clusters变量中
- 确认计算容器中的本地IP获取正确，可按需定制及集群定义中parameters的local_ip_index，保证每个计算节点能正确获取本地IP地址。


### 3.5 启动scalebox集群控制端
```bash
cd ~/docker-scalebox/cluster && make all
```

