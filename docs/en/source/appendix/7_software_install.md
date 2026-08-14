# 7. Base Software Installation and Configuration

- Switch to domestic mirrors
```sh
rm -rf /etc/yum.repos.d/*
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.sau.edu.cn/repo/Centos-8.repo
yum makecache
```

## 7.1 Installing docker

```sh
yum remove -y docker-selinux docker-engine-selinux podman buildah
docker docker-common docker-selinux docker-engine

yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

yum install -y docker-ce docker-ce-cli containerd.io

systemctl enable --now docker
systemctl start docker

```

## 7.2 Installing and Configuring gluster for Cluster Storage

### 7.2.1 Installing gluster
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
### 7.2.2 Cluster Software Configuration

Configure on the HEAD node h0
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

- Mount glusterfs on each node
```sh
mount -t glusterfs h0:vol-0 /gfs
```
