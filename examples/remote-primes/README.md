# remote-primes

## 一、跨集群应用介绍

- 主应用、远端子应用都需基于路由实现。


## 二、跨集群质数计算

## 三、应用创建

## 3.1 创建集群

- 在所有的scalebox runtime实例上，创建所有的集群cluster。

```sh
scalebox cluster create cluster0.yaml 
scalebox cluster create cluster1.yaml 
scalebox cluster create cluster2.yaml 
```

在scalebox runtime实例上，创建当前应用（app）。
### 3.2 先创建本地应用
```sh
app_id_0=$(scalebox run --app-file main.yaml | cut -d':' -f2 | tr -d '}')
```

### 3.3 依次创建远端应用
```sh
app_id_1=$(CLUSTER=cluster1 scalebox run --app-file calc.yaml --remote-app-id=$app_id_0| cut -d':' -f2 | tr -d '}' )
app_id_2=$(CLUSTER=cluster2 scalebox run --app-file calc.yaml --remote-app-id=$app_id_0| cut -d':' -f2 | tr -d '}' )

```

### 3.4 设置各应用为运行状态

```sh
scalebox app set-status --app-id=$app_id_0 RUNNING
CLUSTER=cluster1 scalebox app set-status --app-id=$app_id_1 RUNNING
CLUSTER=cluster2 scalebox app set-status --app-id=$app_id_2 RUNNING
```

### 3.5 给主应用发送起始任务
```sh
echo '1000' | scalebox task add --app-id=$app_id_0
```


## 四、问题与讨论

### 问题：如何设计参数，使得可创建远端应用（app）

### 创建远端任务（task）

- 指定环境变量GRPC_SERVER
```sh
GRPC_SERVER=10.0.6.100 scalebox task add 
```

- 指定环境变量CLUSTER（本地查表转换）

- module-id的标识
  - 直接module-id
  - module-ref：module-id + sink-module
  - app-ref：app-id + sink-module
- remote module-id的标识
  - 当前：```remote_server + [app-id] + sink-module```
  - 调整为：
    - cluster_name + app-id + sink-module
    - remote_server + app-id
    - ```cluster_name + [app-id]```
  
grpc_remote_server
