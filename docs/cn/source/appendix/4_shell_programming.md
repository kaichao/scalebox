# 4. shell编程接口

## 4.1 模块脚本编写规范

## 4.2 标准输入输出接口

## 4.3 文件交换接口规范

## 4.4 时间戳与性能统计


模块内脚本通常以shell实现。要使用内置函数，容器内需安装jq以支持json解析。

以下是debian/ubuntu类镜像的Dockerfile中，安装jq的示例代码：
```Dockerfile
RUN apt update \
    && apt-get install -y jq \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
```

## 4.5 常用 scalebox 内置函数

### 4.5.1 scalebox::json_val

功能：从 JSON 中提取参数值

| 参数 | 说明 |
|------|------|
| $1 | JSON 文本 |
| $2 | JSON 字段名 |

返回：字段值（字符串）

### 4.5.2 scalebox::task_header

功能：从 JSON 消息头中提取值，若不存在则从环境变量中提取（环境变量名为消息头对应的大写形式）

| 参数 | 说明 |
|------|------|
| $1 | JSON 文本 |
| $2 | 字段名（字母、下划线） |

返回：参数值（字符串）

### 4.5.3 scalebox::json_parse

功能：将 JSON 字符串解析为 bash 关联数组

| 参数 | 说明 |
|------|------|
| $1 | JSON 文本 |
| $2 | 关联数组变量名（返回值） |

### 4.5.4 scalebox::is_task_body

功能：判断任务体是否为无空字符的简单文本（非 JSON）

返回：0 = 简单文本，1 = JSON

### 4.5.5 scalebox::append_to_file

功能：追加内容到文件（原子写入，带锁）

| 参数 | 说明 |
|------|------|
| $1 | 文件路径 |
| $2 | 内容 |

### 4.5.6 scalebox::get_local_ip

功能：获取本机 IPv4 地址

返回：IP 地址字符串

### 4.5.7 scalebox::get_slot_seq

功能：获取当前 slot 的序号（0-based）

返回：整数

### 4.5.8 path::host_path

功能：将容器内路径映射为主机路径

| 参数 | 说明 |
|------|------|
| $1 | 容器内目录名 |

返回：容器可访问的主机目录

### 4.5.9 path::size

功能：获取目录字节数

| 参数 | 说明 |
|------|------|
| $1 | 容器可访问的主机目录 |

返回：字节数

### 4.5.10 path::is_dir

功能：判断路径是否为目录

返回：0 = 是目录，非 0 = 不是

## 4.6 内置函数的用法示例

```bash
#!/usr/bin/env bash

source /usr/local/lib/scalebox/functions.sh

my_header=$(scalebox::task_header "$2" "my_header")

path_in_container="mypath"
host_dir=$(path::host_path ${path_in_container})

```

## 4.7 容器内可访问的数据目录

容器内缺省可访问外部目录包括：
- /tmp：本地临时文件目录
- /dev/shm：本地缓存目录（tmpfs）
- /cluster_data_root：集群数据根目录（在集群定义中，用```data_root```定义）
- /local_data_root：计算节点本地根目录

要访问其它目录，需要在模块定义的```volumes```中定义映射关系。
