# hello-scalebox

第一个入门的scalebox应用。

## 一、应用介绍

介绍scalebox应用的基本概念，实现第一个应用。

- 应用：
- 模块：
- 消息：

## 二、功能设计及实现

### 2.1 编写模块脚本

```sh

```

### 2.2 设计运行参数

| 参数名 | 参数值 | 参数说明 |
| ------ | ------- | ------- |
| cluster | local   | 应用所在集群名，缺省为'local' |
| slot-regex | h0  | 计算节点的正则表达式，缺省为头节点'h0' |
| code-path | /path/to/hello-scalebox/code | 脚本代码所在目录的绝对路径，缺省为当前路径下的相对目录./code |
| image-name | /path/to/agent.sif | 镜像名，针对singularity容器，为sif文件路径。缺省为'hub.cstcloud.cn/scalebox/agent:latest' |

### 2.3 运行应用程序

- 所有运行参数全部采用缺省值
```sh
echo "I'm scalebox" | scalebox app run 
```
将消息通过管道传递给scalebox应用创建程序，启动该应用，并将消息传递给该应用的初始模块。

- 非缺省参数运行

```sh
echo "I'm scalebox" | scalebox app run --app-name=hello --cluster=p419 --slot-regex=n0 --image-name=~/singularity/scalebox/agent.env
```

### 2.4 查看应用状态

```sh
scalebox app list
```
