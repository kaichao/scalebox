# app-copy

## ftp
```sh
cd ftp
```
- PUSH
将CentOS8的/etc目录传输到远程ftp服务器
```sh
SOURCE_URL=/etc TARGET_URL=ftp://fast-obs:acddade009@obsftp.cstcloud.cn/etc4 N_RETRIES=0 scalebox app create 

SOURCE_URL=/etc TARGET_URL=ftp://fast-obs:acddade009@obsftp.cstcloud.cn/etc9 N_RETRIES=2 scalebox app create 

SOURCE_URL=/etc TARGET_URL=ftp://fast-obs:acddade009@obsftp.cstcloud.cn/etc11 N_RETRIES=10 ACTION=PUSH_RECHECK scalebox app create 

```
测试结果如下：
|最大重试次数| 总文件数 | 成功文件数  | 失败文件数(返回码14) | 总运行次数|
| ---- | ---- | ---- | ---- | ---- | 
| 0 | 1847 | 1773 | 74 | 1847 |
| 0 | 1847 | 1768 | 79 | 1847 |
| 0 | 1847 | 1777 | 70 | 1847 |
| 1 | 1847 | 1814 | 33 | 1916 |
| 1 | 1847 | 1818 | 29 | 1916 |
| 1 | 1847 | 1818 | 29 | 1918 |
| 2 | 1847 |  |  |  |
| 2 | 1847 |  |  |  |
| 2 | 1847 |  |  |  |

考虑在ftp-copy中增加容错处理。在app.yaml文件的file-copy中，设置以下参数：
```
parameters:
  retry_rules: "['14:${N_RETRIES}']"
```
针对返回码14，再重试N_RETRIES次。再次测试结果：
|测试的总文件数 | 成功传输文件数  | 字节数不一致的文件数(返回码13) | 
|  ----  | ----  | ----  | 
| 1847  |  |  |


```sh
SOURCE_URL=/etc TARGET_URL=ftp://fast-obs:acddade009@obsftp.cstcloud.cn/etc DIR_NAME=sysconfig scalebox app create 


```

- PULL

```sh

SOURCE_URL=ftp://fast-obs:acddade009@obsftp.cstcloud.cn/etc TARGET_URL=/tmp/etc scalebox app create

SOURCE_URL=ftp://fast-obs:acddade009@obsftp.cstcloud.cn/etc TARGET_URL=/tmp/etc REGEX_FILTER=.*allhosts.* scalebox app create

SOURCE_URL=ftp://fast-obs:acddade009@obsftp.cstcloud.cn/etc TARGET_URL=/tmp/etc REGEX_FILTER=.*allhosts.* ENABLE_LOCAL_RELAY=yes scalebox app create

```

