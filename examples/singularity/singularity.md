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

### 常用命令
```sh
singularity pull hello-world.sif docker://hello-world

singularity pull singularity/scalebox/hello-scalebox.sif docker://hub.cstcloud.cn/scalebox/hello-scalebox

singularity build ~/singularity/scalebox/hello-scalebox.sif docker://hub.cstcloud.cn/scalebox/hello-scalebox

singularity build ~/singularity/app-primes/calc.sif docker-daemon://app-primes/calc:latest

singularity run hello-world.sif

singularity run docker://hello-world

singularity exec --bind /mnt/nfs:/mnt/nfs jason-tf.sif python /opt/test.py
```

### 常见运行命令
- 交互式运行
```sh
singularity shell docker://ubuntu
```

- 执行一个命令并退出
```sh
singularity exec docker://ubuntu bash -c  "pwd && id"
```

- 运行一个容器
```sh
singularity run docker://ubuntu
```

### 后台运行容器实例
- 启动实例
```sh
singularity instance start docker://ubuntu test-instance
```

- 查看实例
```sh
singularity instance list
```

- 操作实例
```sh
singularity shell instance://test-instance
```
- 使用exec执行命令
```sh
singularity exec instance://test-instance ps -ef
```
- 停止实例
```sh
singularity instance stop test-instance
```

### 绑定目录
在 shell, run, instance start 等命令中通过 "-B" 选项来实现Docker中“-v”选项提供挂载卷的功能
```sh
singularity shell -B /apps:/apps docker://ubuntu
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
