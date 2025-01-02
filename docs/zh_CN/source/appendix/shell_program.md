# 4. shell编程


模块内脚本通常以shell实现。要使用内置函数，容器内需安装jq以支持json解析。

以下是debian/ubuntu类镜像的Dockerfile中，安装jq的示例代码：
```Dockerfile
RUN apt update \
    && apt-get install -y jq \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
```

## 4.1 常用scalebox内置函数

### 4.1.1 get_host_dir

参数：容器内dir名
返回：容器可访问的主机目录

### 4.1.2 get_json_value

功能：从json中提取参数值
参数列表：
    json文本
    json字段名
返回：json字段值（字符串）

### 4.1.3 get_parameter

功能：从json中提取参数值，若不存在，则从环境变量中提取值（环境变量名为参数名对应的全大写字母）
参数列表：
    json文本
    json字段名：以消息字母、下划线定义
返回：参数值（字符串）

### 4.1.4 parse_json

功能：将json字符串映射为hash table
参数列表：
    json文本
    bash的hashtable：返回值
返回：无


## 4.2 内置函数的用法示例

```bash
#!/usr/bin/env bash

source functions.sh

my_param=$(get_parameter "$2" "my_parameter")

path_in_container="mypath"
host_dir=$(get_host_dir ${path_in_container})

```

## 4.3 容器内可访问的数据目录

容器内缺省可访问外部目录包括：
- /tmp：本地临时文件目录
- /dev/shm：本地缓存目录（tmpfs）
- /local_data_root：计算节点本地根目录
- /cluster_data_root：集群数据根目录（在集群定义中，用```base_data_dir```定义）

要访问其它目录，需要在模块定义的```paths```中定义映射关系。
