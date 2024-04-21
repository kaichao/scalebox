# examples for dir-list & rsync-copy

## 1. ssh-server to local

```sh

SOURCE_URL=root@10.0.6.101/etc TARGET_URL=/tmp FILE_NAME=postfix/master.cf scalebox app create

SOURCE_URL=root@10.0.6.101/etc TARGET_URL=/tmp/etc FILE_NAME=postfix/master.cf scalebox app create

SOURCE_URL=root@10.0.6.101/etc/postfix TARGET_URL=/tmp FILE_NAME=master.cf scalebox app create


```

## 2. local to ssh-server
```sh
SOURCE_URL=/etc TARGET_URL=root@10.0.6.101/tmp/etc FILE_NAME=postfix/master.cf scalebox app create

SOURCE_URL=/ TARGET_URL=root@10.0.6.101/tmp/ FILE_NAME=etc/postfix/master.cf scalebox app create

SOURCE_URL=/etc/postfix TARGET_URL=root@10.0.6.101/tmp FILE_NAME=master.cf scalebox app create

```

- KEEP_SOURCE_FILE test

```sh
echo "hello" > /tmp/myfile.txt

FILE_NAME=/tmp~myfile.txt~root@10.0.6.101/tmp scalebox app create

KEEP_SOURCE_FILE=no FILE_NAME=/tmp~myfile.txt~root@10.0.6.101/tmp scalebox app create

FILE_NAME=root@10.0.6.101/tmp~myfile.txt~/tmp scalebox app create
KEEP_SOURCE_FILE=no FILE_NAME=root@10.0.6.101/tmp~myfile.txt~/tmp scalebox app create

```

## 3. ssh-server to ssh-server
```sh
SOURCE_URL=root@10.0.6.101/etc TARGET_URL=root@10.0.6.102/tmp/etc FILE_NAME=postfix/master.cf scalebox app create

```
## 4. local to rsync-server
```sh

```
