# app-primes-run



## 传统方式运行
```sh
scalebox app create
```

## 分布式部署
```sh
BULK_MESSAGE_SIZE=1000 scalebox app create -e inline.env
```


## 自定义参数
```sh
NUM_GROUPS=10 \
GROUP_SIZE=100 \
scalebox app create
```

## 支持多消息同时处理
```sh
BULK_MESSAGE_SIZE=5 scalebox app create
```

## 命令行方式运行

```sh
scalebox app run --image-name=app-primes/calc ANY

echo ANY | scalebox app run --image-name=app-primes/calc
```

## 查看计算结果
```sh
APP_ID=60 scalebox semaphore get app-primes:sum_value
```
