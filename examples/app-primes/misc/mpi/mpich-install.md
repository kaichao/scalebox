## mpich install on CentOS7
```sh
yum install -y mpich-3.2*
echo "export PATH=/usr/lib64/mpich-3.2/bin:$PATH" >> /etc/profile
```

## mpich install on CentOS7 cluster
```sh
for i in {01..20}; do echo $i; ssh root@r$i 'yum install -y mpich-3.2*;echo "export PATH=/usr/lib64/mpich-3.2/bin:$PATH" >> /etc/profile';done
```
## create mpich user/group on each node
```sh
groupadd -g 1500 mpich
useradd mpich -u 1500 -g mpich
```

## set p2p passwordless ssh on mpich@host
```sh

```

## open mpich tcp port 45817 on each host
```sh
systemctl stop firewalld
```

## The mpi program runs on the cluster
```sh
mpirun -np 6 -f nodes-list mpi-program
```
