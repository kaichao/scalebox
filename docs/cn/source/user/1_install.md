# 1. 安装部署

ScaleBox集群有通常有1个头节点（HEAD）、若干个计算节点（NODE）组成。HEAD节点运行着controld、actuator、database等服务；NODE节点执行具体计算的任务。

## 1.1 环境要求

所有节点安装64位Linux（CentOS7/8/9、Debian 12/13、Ubuntu 20/22）等。

- 芯片架构及操作系统

|  芯片架构  | 操作系统  |  说明              |
| -------- | -------- | ----------------- |
| x86_64   | Linux    |                   |
| arm64    | Linux    |                   |
| x86_64   | MacOS    | 用于开发环境        |
| arm64    | MacOS    | 待测试，用于开发环境 |
| x86_64   | Win64    | 待测试，用于开发环境 |

- 容器引擎/容器运行时

|  容器引擎/容器运行时     | 版本号            |  说明     |
| --------------------- | ----------------- |  ------- |
| docker-ce             | 20.10<sup>+</sup> |          |
| podman                | 4.8<sup>+</sup>   | 待测试    |
| containerd + nerdctl  | 1.6<sup>+</sup>   | 待测试    |
| singularity           | 3.8<sup>+</sup>   |          |
| apptainer             |                   | 待测试    |
| Kata Containers       |                   | 待测试    |


## 1.2 安装步骤

本示例中有1个HEAD节点、4个NODE节点。相关配置如下：

| 名称 | 类型 |    IP地址   |
| --- | ---- | ---------- |
| h0  | HEAD | 10.0.6.100 |
| n0  | NODE | 10.0.6.101 |
| n1  | NODE | 10.0.6.102 |
| n2  | NODE | 10.0.6.103 |
| n3  | NODE | 10.0.6.104 |

每个节点配置：
| 类型     | 值      |
| ------- | ------- |
| CPU     |     4核 |
| 内存    |    16GB |
| 本地硬盘 |   200GB |
| 操作系统 | CentOS8 |


### 1.2.1 基础安装
- 软件包获取
- 依赖安装
- 核心组件部署

## 1.2.1 头节点安装


## 1.2.2 计算节点安装


### 1.2.2 集群部署
- 主节点配置
- 工作节点加入
- 网络连通性测试

- 设置免密登录到HEAD节点
```sh
ssh-copy-id root@10.0.6.100
```
- 设置HEAD节点上host文件
```sh
hostname h0

cat >> /etc/hosts << EOF
10.0.6.100  h0
10.0.6.101  n0
10.0.6.102  n1
10.0.6.103  n2
10.0.6.104  n3
EOF
```
- 设置HEAD节点到NODE节点root账号的ssh免密登录
```sh
ssh-keygen
ssh-copy-id n0
ssh-copy-id n1
ssh-copy-id n2
ssh-copy-id n3


scp /etc/hosts n0:/etc
scp /etc/hosts n1:/etc
scp /etc/hosts n2:/etc
scp /etc/hosts n3:/etc

ssh n0 hostname n0
ssh n1 hostname n1
ssh n2 hostname n2
ssh n3 hostname n3
```

- 针对所有节点，更换国内源
```sh
rm -rf /etc/yum.repos.d/*
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.sau.edu.cn/repo/Centos-8.repo
yum makecache
```

- 设置NTP

- 针对所有NODE，更换国内源
```sh
rm -rf /etc/yum.repos.d/*
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.sau.edu.cn/repo/Centos-8.repo
yum makecache
```

- 设置NTP

#### 安装gluster7，构建集群存储
- 每个节点上安装软件
```sh
# CentOS8/CentOS7
yum install -y centos-release-gluster6

yum -y install glusterfs-server 
systemctl enable --now glusterd.service
systemctl start glusterd.service

rpm -qi glusterfs-server 
gluster --version

mkdir /opt/vol-0 /gfs
```
- HEAD节点h0上配置
```sh
gluster peer probe n0
gluster peer probe n1
gluster peer probe n2
gluster peer probe n3

# gluster volume create vol 110.0.6.{100,101,102,103,104}:/opt/vol-0

gluster volume create vol-0 disperse 4 redundancy 1 n{0,1,2,3}:/opt/vol-0
gluster volume start vol-0
gluster volume info
```

- 每个节点上mount glusterfs
```sh
mount -t glusterfs h0:vol-0 /gfs
```

## 三、安装docker

- 所有节点安装docker
```sh
yum remove -y docker-selinux docker-engine-selinux podman buildah
docker docker-common docker-selinux docker-engine

yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

yum install -y docker-ce docker-ce-cli containerd.io

# add user scalebox
useradd scalebox
usermod -aG docker scalebox
# set a unified password for the scalebox user of all nodes in the cluster
passwd scalebox

systemctl enable --now docker
systemctl start docker
```


## 1.3 配置说明
### 1.3.1 基础配置
- 核心参数
- 日志设置
- 资源限制

### 1.3.2 高级配置
- 网络定制
- 存储后端
- 安全策略

## 1.4 验证安装
### 1.4.1 健康检查
- 服务状态
- 组件连通性
- 基本功能测试

### 1.4.2 性能基准
- 单节点测试
- 集群测试
- 压力测试


