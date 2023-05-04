## CentOS 8上singularity安装

- 参考文档
[Installing SingularityCE](https://docs.sylabs.io/guides/3.11/admin-guide/installation.html)

```sh
yum install -y \
   libseccomp-devel \
   glib2-devel \
   squashfs-tools \
   cryptsetup \
   runc \
   crun

rpm -ivh https://github.com/sylabs/singularity/releases/download/v3.11.2/singularity-ce-3.11.2-1.el8.x86_64.rpm

```
## Singularity使用
- [Introduction to SingularityCE](https://docs.sylabs.io/guides/3.11/user-guide/introduction.html)

容器(container)是一个包含用户软件和依赖的镜像系统，可独立运行某一条或者多条命令。
Singularity没有镜像的概念，用户创建和运行的都是一个一个容器。

Singularity容器有两种存在形式：
- SIF（Singularity Image File）：压缩后的只读（read-only）的Singularity镜像文件，是生产使用的主要形式。
- Sandbox：可写(writable)的容器存在形式，是文件系统中的一个目录，常用于开发或者创建自己的容器，是开发使用的主要形式。

- 常用命令
```sh
singularity pull hello-world.sif docker://hello-world

singularity run hello-world.sif

singularity run docker://hello-world

singularity exec --bind /mnt/nfs:/mnt/nfs jason-tf.sif python /opt/test.py
```

## Singularity选项
```

--env 环境变量
--rocm          enable experimental Rocm support
-S, --scratch strings   include a scratch directory within the container that is linked to a temporary dir (use -W to force location)
-W, --workdir string   working directory to be used for /tmp, /var/tmp and $HOME (if -c/--contain was also used)
-w, --writable    by default all Singularity containers are available as read only. This option makes the file system accessible as read/write.
--writable-tmpfs  makes the file system accessible as read-write with non persistent data (with overlay support only)

```
