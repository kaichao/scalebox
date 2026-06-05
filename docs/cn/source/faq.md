# 常见问题（FAQ）

## 1. Scalebox Runtime非标准端口的设置

scalebox涉及的网络端口包括：
- 数据库端口`<custom-pg-port>`，其标准端口号5432
- controld的grpc端口`<custom-grpc-port>`，其标准端口号50051

### 1.1 服务端设置
- 修改`runtime/compose.yaml`

```yaml
  controld:
    ports:
      - <custom-grpc-port>:50051
  database:
    ports:
      - <custom-pg-port>:5432
```
- 重启runtime

### 1.2 客户端设置
- 设置`${HOME}/.scalebox/environments`，使得当前用户下所有操作都可以用非标准端口操作
`
PGHOST=<internal-ip>:<custom-grpc-port>
GRPC_SERVER=<internal-ip>:<custom-pg-port>
`

### 1.3 导入cluster配置
- 修改cluster配置文件，重新导入数据库
```yaml
cluster-name:
  parameters:
    data_root: /raid0/scalebox/mydata
    grpc_server: <internal-ip>:<custom-grpc-port>
    remote_grpc_server: <external-ip>:<custom-grpc-port>
    pghost: <internal-ip>:<custom-pg-port>
    remote_pghost: <external-ip>:<custom-pg-port>
    uname: <user-name>
    port: <ssh-port>
```


> 持续更新中，欢迎贡献问题与答案。
