# Scalebox集群

Scalebox集群是运行scalebox应用的软硬件资源集合。分为
- 内联集群（Inline Cluster）：资源由Scalebox直接管理的计算集群。
  - 按其资源集合可否扩展，分为固定集群、可扩展集群（比如公有云）；
  - 按作业调度分配方式，分为资源预分配、资源动态分配。
- 外部集群（External Cluster）：资源由外部系统管理的计算集群，可支持用slurm、k8s等外部程序。

Scalebox的主要对象：
- 集群Cluster：
- 主机Host：
- 应用App：
- Job：
- 任务Task：
- Slot：


## 一、Scalebox集群介绍
### 1.1 头节点
头节点上安装了单个scalebox集群的管理服务，主要包括：
- actuator：启动端，负责在计算节点上启动slot；针对内联集群，直接通过免密ssh启动；针对外部集群，通过调用集群调度系统来启动。
- controld：面向actuator、计算节点，提供基于grpc的控制端应用服务；
- database：基于postgresql的核心数据库，存放app、job、task、slot等相关数据，面向controld等提供数据存储、检索等服务；也可面向计算节点、命令行工具提供直接数据库访问。

头节点以及头节点服务（controld/actuator/database），可为多个集群所共享。

### 1.2 计算节点
计算节点分为两类：
- 内部计算节点：scalebox内部调度管理的计算节点，通过免密ssh启动slot；
- 外部计算节点：通过外部调度程序（slurm/k8s等）启动的计算节点，在节点上启动slot。

## 二、Scalebox单节点集群安装及配置
- 操作系统：
  - CentOS 7以上（其他版本的Linux待测试）
  - macOS 10.15(amd64)以上（ARM版的macos待测试）
- 容器化引擎：DockerCE / rootless docker，版本20.10+
- docker-compose: 1.29.2+
- 安装dstat、htop、rsync、zstd、gmake、git、rsync、glances、dstat、pdsh等工具软件，用于性能监控、开发运行等

### 2.1 安装CentOS 7/8下基本软件
- 以root用户安装
```bash
yum install -y epel-release
yum install -y rsync pv pdsh wget make dstat nmon glances htop
```
- CentOS 7安装git v2

```sh
yum -y remove git

yum install -y https://packages.endpointdev.com/rhel/7/os/x86_64/endpoint-repo.x86_64.rpm

yum install -y git
```

- sshd
  - turn off UseDNS/GSSAPIAuthentication

- Setup Linux Time Calibration

### 2.2 安装macos下基本软件

启动本机sshd

基于Homebrew，安装所需基本软件

### 2.3 安装docker及docker-compose
安装docker可以有两种方式：
- 以root用户安装docker-ce，docker可供本节点的所有用户共享，或
- 以当前用户安装rootless docker

#### 2.3.1 安装docker-ce

```sh
yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
systemctl start docker
systemctl enable docker
```

#### 2.3.2 安装rootless docker

- 以Rocky Linux 9 x64为例
以root安装必须软件
```sh
dnf install -y fuse-overlayfs
dnf install -y iptables
modprobe ip_tables
```

- 参考文档
[Run the Docker daemon as a non-root user (Rootless mode)](https://docs.docker.com/engine/security/rootless/)

- 新增用户scalebox，并以scalebox用户进行后续操作

- 设置环境变量
```bash

echo "export DOCKER_HOST=unix:///run/user/$(id -u)/docker.sock" >> ~/.bashrc
echo "export DOCKER_HOST=unix:///run/user/$(id -u)/docker.sock" >> ~/.bash_profile
```

- 安装rootless docker
```bash
curl -fsSL https://get.docker.com/rootless | sh
```

- 启动rootless docker
```bash
systemctl --user enable docker
systemctl --user start docker

```

#### 2.3.3 验证docker有效
```bash
docker run --rm hello-world
```

### 2.3.4 以scalebox用户，下载安装docker-compose
- 安装docker-compose
```bash
mkdir -p ~/bin
wget -O ~/bin/docker-compose https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64
chmod +x ~/bin/docker-compose
```

### 2.4 scalebox的安装配置

#### 2.4.1 下载scalebox

在${HOME}下，下载scalebox

```bash
cd && git clone https://github.com/kaichao/scalebox && cd scalebox && git checkout dev
```

#### 2.4.2 获取最新容器镜像及命令行工具
```bash
cd ~/scalebox/server
make pull-all
make get-cli
```

#### 2.4.3 更新系统缺省的公钥（可选）
```bash
cd ~/scalebox/server
make update-pubkey
```

#### 2.4.4 准备运行环境
- 创建运行相关目录
- 设置环境变量
- 设置内联集群内的免密ssh

```bash
make prepare
```

#### 2.4.5 配置头节点服务的本地IP地址
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
```
则可在defs.mak中设置为：
```LOCAL_IP_INDEX=2```
或
```LOCAL_ADDR=192.168.56.21```

#### 2.4.6 退出当前shell，并重新登录，使得配置生效


#### 2.4.7 启动scalebox控制端
```bash
cd ~/scalebox/server && make all
```
至此，scalebox控制端启动完成。下一步，可以到[examples](../examples/)目录下，运行```hello-scalebox、app-primes```两个应用案例。

## 三、Scalebox多节点集群安装及配置

### 3.1 头节点安装

在单集群基础上，头节点安装time-server、pdsh、pv、rsync
```sh
yum install -y pdsh pdsh-rcmd-ssh
export PDSH_RCMD_TYPE=ssh
```

### 3.2 设置集群的共享存储
设置可用于计算节点间的外部共享存储，共享存储可以是NFS或glusterfs等。

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

- 参照目录 ```inline-cluster```下的```mycluster.yaml```，定义计算节点配置文件，并更新```Makefile```中的```clusters```变量

### 3.4 外部计算节点
- 操作系统：CentOS 7以上
- 单机容器化引擎：
  - Docker 20.10以上
  - podman ？版本；
  - singularity 3.8以上
- k8s集群

### 3.5 集群定义文件

- 参照```inline-cluster```目录，完成集群定义文件```mycluster.yaml```
- 将新集群名称，加入到Makefile文件的clusters变量中
- 确认计算容器中的本地IP获取正确，可按需定制及集群定义中parameters的local_ip_index，保证每个计算节点能正确获取本地IP地址。

### 3.5 启动scalebox集群服务端
```bash
cd ~/scalebox/server && make all
```

至此，scalebox服务端安装完成，可以通过[scalebox/examples](../examples)中的```hello-scalebox```、```app-primes```两个应用来测试平台安装是否正确。

## 四、Scalebox多集群安装及跨集群应用配置

