# 7. 基础软件安装配置

- 更换国内源
```sh
rm -rf /etc/yum.repos.d/*
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.sau.edu.cn/repo/Centos-8.repo
yum makecache
```

## 7.1 docker安装

```sh
yum remove -y docker-selinux docker-engine-selinux podman buildah
docker docker-common docker-selinux docker-engine

yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

yum install -y docker-ce docker-ce-cli containerd.io

systemctl enable --now docker
systemctl start docker

```

## 7.2 集群存储gluster安装配置

### 7.2.1 gluster软件安装
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
### 7.2.2 cluster软件配置

在HEAD节点h0上配置
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
