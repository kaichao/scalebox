# 跨集群质数计算
跨物理集群的质数计算。

## 1. 确定两个集群头节点的IP地址
```sh
export SERVER_1=10.0.6.100
export SERVER_2=10.0.6.104
```
## 2. 获取两个集群的最新可用app-id
```sh
export APP_ID_1=`PGHOST=${SERVER_1} scalebox app get-next-id`
export APP_ID_2=`PGHOST=${SERVER_2} scalebox app get-next-id`
```
## 3. 在集群2上，启动calc模块；在集群1上，启动scatter+gather模块
```sh
PGHOST=${SERVER_2} GRPC_SERVER=${SERVER_2} scalebox app create --debug cluster2.yaml
PGHOST=${SERVER_1} GRPC_SERVER=${SERVER_1} scalebox app create --debug cluster1.yaml
```
## 4. 通过cli或gui，检查计算结果的有效性
