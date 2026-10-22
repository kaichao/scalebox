# 1. 技术规范与命名规则

> 本章给运行时组件、schema 骨架与命名规则；字段清单见 [reference/app-yaml.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)。

Scalebox 应用是运行在 Scalebox 平台上的应用程序；典型场景包括高通量数据处理、大规模数据传输等。

## 1.1 Scalebox运行时组件

### 1.1.1 核心运行时组件

- **controld**：基于 gRPC 的控制服务，管理 actuator 和计算节点
- **actuator**：启动器服务，通过 SSH 或外部调度器在计算节点上启动 slot
- **database**：PostgreSQL 数据库，存储 app、module、task、slot 等元数据
- **webui**：Scalebox 操作管理的 Web 接口
- **agent**：运行在计算节点上的代理，负责任务执行和状态管理
- **cluster-admin**：与 HPC 调度系统交互的资源申请、管理模块，也是一个系统级应用

### 1.1.2 组件间通信

- **控制平面通信**：gRPC，端口 50051
- **数据平面通信**：agent 从数据库领取任务，模块通过 `sink-tasks.txt` 回流；数据经共享存储交换
- **节点间通信**：SSH 或专用网络协议

## 1.2 应用定义规范

应用定义文件（缺省 `app.yaml`）定义 App、Module 与 Cluster。最小骨架：

```yaml
name: perf-test.my-app
label: perf-test
version: 1.0.0
cluster: ${CLUSTER}
parameters:
  initial_status: RUNNING
  main_router: main-router
  is_cluster_admin: yes
  default_sleep_count: 20
modules:
  module-1:
    arguments: {...}
    parameters: {...}
clusters:
  cluster-1: {...}
```

- 字段清单与类型见 [`reference/app-yaml.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)。
- 若文件仅用于 cluster 定义，主配置属性不生效。
- 缺省文件名为当前目录的 `app.yaml`。

模板变量：除虚拟模板参数 `CLUSTER_DATA_DIR` 外，模板变量必须先定义再使用。优先级从低到高：

- `/etc/scalebox/environments`
- `${HOME}/.scalebox/environments`
- `${PWD}/scalebox.env`
- `${PWD}/${{env-defined}}.env`
- `${PWD}/${{env-defined}}_${{app-defined}}.env`
- 当前命令行中已定义的环境变量

## 1.3 模块定义规范

```yaml
my-module:
  label: My First Module
  base_image: scalebox.net/platform/agent:latest
  cluster: my-cluster
  command: docker run -d --network=host {{ENVS}} {{VOLUMES}} {{IMAGE}}
  arguments:
    code_path: ${PWD}/code
  parameters:
    task_dist_mode: DEFAULT
  environments:
    - MY_VAR=value
  volumes:
    - /host/path:/container/path
  slots:
    - ${nodes}[:n]
  sink_modules:
    - next-module
  comment: This is an algorithm module.
```

- 字段清单与类型见 [`reference/app-yaml.md`](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)。
- `slots` 支持两种写法：

  ```yaml
  slots:
    - ${nodes}[:n]
    - ${nodes}:${n}:${group_prefix}
  ```

  第二行用于 HOST-BOUND 的全局/分组 slot；`nodes` 写主机短名或正则（如 `h0`、`n0-[01]`），不要写 FQDN；`n` 为数量；`group_prefix` 为正则，在分组 slot 的 `parameters` 中设置。
- `sink_modules` 标识模块间的物理关联（跨集群）；`sink_vmodules` 表示模块间的逻辑关系。

## 1.4 集群定义规范

```yaml
mycluster:
  label: My new cluster
  parameters:
    uname: myuser
    port: 10022
    data_root: /global-fs/scalebox/mydata
    local_ip_index: 2
    grpc_server: 192.168.3.123:50051
  total_resources:
    num_cores: cpu cores
    total_mem_gb:
    total_disk_tb:
  status: ON
  comment:
```

`parameters` 的字段语义见 [appendix/2 §2.2.1](2_parameter_reference.md)。`num_of_executors` 与 `channel_size` 仅 Inline cluster 使用。

## 1.5 标识符命名规则

### 1.5.1 文件名命名规则

文件名字符：数字、英文字母大小写、下划线、点。

### 1.5.2 URI 命名规则

App、Module 等资源通过 URI（Uniform Resource Identifier，统一资源标识符）作唯一标识。URI 包括 URL（Uniform Resource Locator）与 URN（Uniform Resource Name）；URL 和 URN 都是 URI 的子集。App、Module 等通过 URN 定义。

### 1.5.3 版本号命名规则

App、Module 等资源可以用版本表示，遵循语义化版本。版本格式为 `主版本号.次版本号.修订号`：

- 主版本号：不兼容的 API 修改
- 次版本号：向下兼容的功能性新增
- 修订号：向下兼容的问题修正

先行版本号及版本编译元数据可以追加在版本号之后。

## 1.6 SSH 加密算法性能

- 单线程性能对比：

| 算法 | 参数值 | 硬件 AES-NI（x64） | 软件实现（无 AES-NI） | 说明 |
|---|---|---|---|---|
| AES-128-GCM | `aes128-gcm@openssh.com` | 5~7 GB/s | 150~200 MB/s | 最优选择，硬件加速极快 |
| AES-256-GCM | `aes256-gcm@openssh.com` | 3~5 GB/s | 120~150 MB/s | 稍慢但安全性更高 |
| ChaCha20-Poly1305 | `chacha20-poly1305@openssh.com` | 1~2 GB/s | 600~700 MB/s | 软件效率高，尤其在无 AES-NI 时比 AES 快 |
