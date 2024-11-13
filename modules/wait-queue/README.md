# 等待队列模块

基于信号量实现流控的等待队列

## 环境变量 

SEMA_NAME：用于等待队列判断的信号量名称。若信号量 > 0；则放行1个消息，并将信号量减1；


## 测试过程

### 1. 创建应用

```sh
cd test
scalebox app create

```

### 2. 创建semaphore

```sh
APP_ID=1 scalebox semaphore create my_sema 2
```

### 3. 添加消息

```sh
APP_ID=1 scalebox task add --sink-job wait-queue --task-file messages.txt
```

### 4. 修改semaphore

```sh
APP_ID=1 scalebox semaphore increment my_sema
```
