# Scalebox集群

Scalebox集群是运行scalebox应用的硬件资源集合。分为
- 内联集群（Inline Cluster）：由Scalebox直接管理的计算资源组成的集群。支持静态资源分配、动态资源分配两种方式。
- 外部集群（External Cluster）：有外部调度系统管理的计算资源组成的集群，可支持用slurm、k8s等。

## Scalebox集群介绍
### 头节点
头节点上安装了单个scalebox集群的管理服务，主要包括：
- controld：面向actuator、计算节点，提供基于grpc的控制端应用服务；
- database：基于postgresql的数据库，存放app、job、task、slot等相关数据，面向controld等提供数据存储、检索等服务。
- actuator：启动端，负责在计算节点上启动slot。

头节点以及头节点服务（controld/actuator/database）可为多个集群所共享。

### 计算节点
计算节点分为两类：
- 内部计算节点：scalebox内部调度管理的计算节点，通过免密ssh启动slot；
- 外部计算节点：通过外部调度程序（slurm/k8s等）启动的计算节点，由外部启动slot。

## Scalebox单节点集群安装
- 操作系统：CentOS 7以上；
- 容器化引擎：DockerCE 20.10+
- docker-compose: 1.29.2+
- 安装dstat、htop、zstd、gmake、git、rsync等工具软件，用于性能监控等

### 安装CentOS 7/8下基本软件
- 以root用户安装
```bash
yum install -y epel-release
yum install -y htop dstat rsync pv pdsh wget make

```
- 安装git v2


- sshd
  - turn off UseDNS/GSSAPIAuthentication

- Setup Linux Time Calibration
### 安装scalebox命令工具
- 以当前普通用户身份

```bash
mkdir -p ~/bin
wget -O ~/bin/docker-compose https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64
chmod +x ~/bin/docker-compose

mkdir -p ~/.ssh ~/.scalebox/log
chmod 700 ~/.ssh

docker pull hub.cstcloud.cn/scalebox/cli
id=$(docker create hub.cstcloud.cn/scalebox/cli) 
docker cp $id:/usr/local/bin/scalebox ~/bin/ 
docker rm -v $id

chmod +x ~/bin/*

echo "alias app='scalebox app'" >> ${HOME}/.bash_profile
echo "alias job='scalebox job'" >> ${HOME}/.bash_profile
echo "alias task='scalebox task'" >> ${HOME}/.bash_profile

```
### 下载安装docker-scalebox
```bash
cd && git clone https://github.com/kaichao/docker-scalebox
```
- setup passwordless from actuator to node(设置actuator可通过免密ssh启动slot)
```bash
cat ~/docker-scalebox/cluster/id_rsa.pub >> ${HOME}/.ssh/authorized_keys
```

### 启动scalebox控制端
- 通过以下命令，检查本地IP地址是否需要设定
```bash
hostname -i
```
- 前一步的本地IP地址不正确，则需要设置local/defs.mk文件中LOCAL_IP_INDEX或LOCAL_ADDR变量

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
- 启动运行scalebox控制端
```bash
cd ~/docker-scalebox/cluster && make all
```

## Scalebox多节点集群安装

### 头节点安装

在单集群基础上，安装time-server、pdsh、pv
```sh
yum install -y pdsh pdsh-rcmd-ssh
export PDSH_RCMD_TYPE=ssh
```


### 内部计算节点安装
- 操作系统：CentOS 7以上
- 容器化引擎：Docker 20.10以上
- 安装dstat、htop、zstd等工具软件，用于性能监控等
- 
可选：
- Linux Time Calibration
- 增加头节点管理用户、actuator的公钥
- docker-ce
- dstat
- htop
- zstd
- sshd 
  - turn off UseDNS GSSAPIAuthentication

```sh
yum install -y epel-release
yum install -y htop dstat rsync pv
```


### 外部计算节点
- 操作系统：CentOS 7以上
- 单机容器化引擎：
  - Docker 20.10以上
  - podman ？版本；
  - singularity 3.8；
- k8s集群

## 计算节点上gluster存储安装

