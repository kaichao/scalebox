
# dir-list

## 一、模块介绍

dir-list is a common module in scalebox. Its function is to traverse the file list of the local directory or the rsync remote directory, generate messages and send it to the subsequent module.

dir-list supports four types of directories:
- LOCAL: Local Directory
- ssh: 基于SSH的Server端，未安装rsync
- rsync-over-ssh: 基于SSH的Server端，已安装rsync
- native rsync: 基于rsyc的Server端

## 二、环境变量

| variable name   | Description  |
|  ----  | ----  |
| PREFIX_URL  | See table below. |
| REGEX_FILTER | File filtering rules represented by regular expressions |
| JUMP_SERVERS | 基于ssh的Server端的跳板机 |
| FILE_CHECKSUM_ALGO | 基于local、ssh的PREFIX_URL，在生产的文件消息头中，添加该文件的文件校验和。校验和算法包括：MD-5 / SHA-1 / SHA-256 / SHA-512 / RIPEMD-160 |
| RSYNC_PASSWORD | Non-anonymous rsync user password |


### 三、PREFIX_URL

| type | description |
| --- | ---- |
| local | represented by an absolute path ```</absolute-path> ```|
| rsync | anonymous access: ```rsync://<rsync-host><rsync-base-dir>```<br/> non-anonymous access: ```rsync://<rsync-user>@<rsync-host><rsync-base-dir>```|
| rsync-over-ssh | The ssh public key is stored in the ssh-server account to support password-free access <br/> ``` [<ssh-user>@]<ssh-host>[:<ssh-port>][#<jump-servers>#]<ssh-base-dir>```, default ssh-user is root, default ssh-port is 22. The format of jump-servers is ```[<user1>@]<host1>[:<port1>],[<user>@]<host2>[:<port2>] ``` |

## 四、输入消息

DIR_NAME: 
- relative_path: Relative to the subdirectory of SOURCE_URL, use "." to represent the current directory of SOURCE_URL

## 五、输出消息
mei

## 六、模块退出码（Exit Code）



## 九、TO-DO

- 针对所有文件，可在消息头中增加file-size、modify-time；
- 针对本地文件，可在消息中增加MD5、SHA1摘要。

