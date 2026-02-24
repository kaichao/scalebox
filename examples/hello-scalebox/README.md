# hello-scalebox

第一个入门的scalebox应用。

## 一、应用介绍

介绍scalebox应用的基本概念，实现第一个应用。

- 应用(App)：
- 模块(Module)：
- 任务(Task)：

## 二、功能设计及实现

### 2.1 系统设计

### 2.2 编写模块脚本

```sh

```

### 2.3 设计运行参数

| 参数名      | 参数值                        | 参数说明                             |
| ---------- | ---------------------------- | ----------------------------------- |
| cluster    | local                        | 应用所在集群名，缺省为'local'          |
| slot-regex | h0                           | 计算节点的正则表达式，缺省为头节点'h0'   |
| code-path  | /path/to/hello-scalebox/code | 脚本代码所在目录的绝对路径，缺省为当前路径下的./code |
| image-name | /path/to/agent.sif           | 镜像名，针对singularity容器，为sif文件路径。缺省为'hub.cstcloud.cn/scalebox/agent:latest' |

### 2.4 运行应用程序

#### 2.4.1 所有运行参数全部采用缺省值
```sh
echo "Docker-based_Scalebox" | scalebox app run 
```
将任务列表通过管道传递给scalebox命令，创建应用程序，启动该应用，并将任务列表发送给该应用的初始模块。

#### 2.4.2 非缺省参数运行

- 构建singularity镜像文件
```sh
	sif_file=/shared/singularity/scalebox/agent.sif
	docker save hub.cstcloud.cn/scalebox/agent:latest -o agent.tar
	mkdir -p $(dirname $sif_file)
	singularity build -F $sif_file docker-archive://agent.tar
	rm -f agent.tar
```

- 运行应用
```sh
echo "Singularity-based_Scalebox" | scalebox run --image-name=$sif_file
```

### 2.5 查看应用状态

```sh
scalebox app list
```


