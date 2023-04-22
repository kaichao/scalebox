# examples

## 1. hello-scalebox
scalebox的入门应用，展示基于消息触发的基本应用。具体操作步骤如下：

- 构建应用模块的容器镜像
```bash
cd examples/hello-scalebox
make build
```
- 创建应用
```sh
app create
```
- 查看应用状态：
```bash
$ app list

```


## 2. app-primes
计算区间内质数数量的分布式应用。首先将区间分解为若干个子区间，再分别计算各个子区间的质数数量，最后加和得到区间内结果。

操作步骤如下：
```bash
cd examples/app-primes
make build
```
- 创建应用
```sh
app create
```
- 查看应用状态：
```bash
$ app list

```


## 3. app-copy
文件拷贝是scalebox的基本功能，实现集群内、跨集群的数据传输。

scalebox的基本模块支持基于rsync、rsync-over-ssh、ftp等的文件拷贝。

app-copy应用则展示利用基本模块实现数据传输的例子。


