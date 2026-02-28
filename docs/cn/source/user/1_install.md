# 1. 安装部署

ScaleBox集群有通常有1个头节点（HEAD）、若干个计算节点（NODE）组成。HEAD节点运行着controld、actuator、database等服务；NODE节点执行具体计算的任务。

## 1.1 环境要求


### 1.1.1 硬环境要求

- 芯片架构及操作系统

|  芯片架构  | 操作系统  |  说明              |
| -------- | -------- | ----------------- |
| x86_64   | Linux    |                   |
| arm64    | Linux    |                   |
| x86_64   | MacOS    | 用于开发环境        |
| arm64    | MacOS    | 待测试，用于开发环境 |
| x86_64   | Win64    | 待测试，用于开发环境 |

生产环境中，所有节点安装64位Linux（CentOS7/8/9、Debian 12/13、Ubuntu 20/22/24）等。

- CPU架构
  - x86_64
  - arm64（已测试）
  - riscv（待测试）
- 内存需求
  - 头节点：≥8GB
  - 计算节点：按需，推荐 ≥8GB

- 存储空间
  - 头节点：≥100GB
  - 计算节点：按需，推荐 ≥100GB

### 1.1.2 软件依赖

#### 容器引擎/容器运行时

  - 计算节点
|  容器引擎/容器运行时     | 版本号            |  说明     |
| --------------------- | ----------------- |  ------- |
| docker-ce             | 20.10<sup>+</sup> |          |
| singularity           | 3.8<sup>+</sup>   |          |
| podman                | 4.8<sup>+</sup>   | 待测试    |
| containerd + nerdctl  | 1.6<sup>+</sup>   | 待测试    |
| apptainer             |                   | 待测试    |
| Kata Containers       |                   | 待测试    |

  - 头节点：docker-ce 20.10<sup>+</sup>

#### 数据库
- 头节点：postgresql 17<sup>+</sup>，以容器化部署

#### 集群存储（可选）
- 所有节点上安装glusterfs，构建统一的集群存储

#### 网络配置

- 头节点、计算节点属于内网的同一子网内
- 若需支持跨集群计算，头节点及相关传输节点需与其他集群通信的外网地址

## 1.2 安装步骤
### 1.2.1 基础安装
#### 头节点安装docker
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

#### 软件包获取
在头节点上安装git，获取软件安装包：
```sh
git clone github.com/kaichao/scalebox
```
#### 容器镜像获取
在头节点上下载容器镜像
```sh

docker pull hub.cstcloud.cn/scalebox/controld:latest
docker pull hub.cstcloud.cn/scalebox/actuator:latest
docker pull hub.cstcloud.cn/scalebox/database:latest

```

### 1.2.2 集群部署

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

#### 所有节点配置

- 针对所有节点，更换国内源
```sh
rm -rf /etc/yum.repos.d/*
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.sau.edu.cn/repo/Centos-8.repo
yum makecache
```

- 设置时间同步NTP

- 安装集群存储软件glusterfs（可选）
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

#### 主节点配置
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

#### 计算节点配置

- 设置HEAD节点到NODE节点root账号的ssh免密登录
```sh
ssh-keygen
ssh-copy-id n0
ssh-copy-id n1
ssh-copy-id n2
ssh-copy-id n3
```

- 设置主机名及hostname
  
```sh

scp /etc/hosts n0:/etc
scp /etc/hosts n1:/etc
scp /etc/hosts n2:/etc
scp /etc/hosts n3:/etc

ssh n0 hostname n0
ssh n1 hostname n1
ssh n2 hostname n2
ssh n3 hostname n3

```

#### 网络连通性测试

#### 配置glusterfs

- HEAD节点h0上配置
```sh
gluster peer probe n0
gluster peer probe n1
gluster peer probe n2
gluster peer probe n3

# gluster volume create vol 10.0.6.{100,101,102,103,104}:/opt/vol-0

gluster volume create vol-0 disperse 4 redundancy 1 n{0,1,2,3}:/opt/vol-0
gluster volume start vol-0
gluster volume info
```

- 每个节点上mount glusterfs
```sh
mount -t glusterfs h0:vol-0 /gfs
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
