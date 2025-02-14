# examples for dir-list & rsync-copy

## 1. ssh-server to local

```sh

SOURCE_URL=root@10.0.6.101/etc TARGET_URL=/tmp FILE_NAME=postfix/master.cf scalebox app create

SOURCE_URL=root@10.0.6.101/etc TARGET_URL=/tmp/etc FILE_NAME=postfix/master.cf scalebox app create

SOURCE_URL=root@10.0.6.101/etc/postfix TARGET_URL=/tmp FILE_NAME=master.cf scalebox app create

SOURCE_URL=root@10.0.6.101/etc/postfix SOURCE_MODE=SSH TARGET_URL=/tmp FILE_NAME=master.cf scalebox app create
```


```sh
SOURCE_URL=root@10.0.6.101/tmp TARGET_URL=/tmp FILE_NAME=master.cf KEEP_SOURCE_FILE=no scalebox app create

SOURCE_URL=root@10.0.6.101/tmp TARGET_URL=/tmp FILE_NAME=postfix/master.cf KEEP_SOURCE_FILE=no scalebox app create

SOURCE_URL=root@10.0.6.101/tmp TARGET_URL=/tmp FILE_NAME=etc/postfix/master.cf KEEP_SOURCE_FILE=no scalebox app create
```

```sh
SOURCE_URL=scalebox@159.226.237.136:10022/raid0/tmp/mwa/tar1257010784 TARGET_URL=/data1/mydata/mwa/tar FILE_NAME=1257010784/1257015406_1257015435_ch119.dat.tar.zst  scalebox app create

SOURCE_URL=scalebox@159.226.237.136:10022/raid0/tmp/mwa/tar1257010784 TARGET_URL=/data1/mydata/mwa/tar FILE_NAME=1257010784/1257015406_1257015435_ch120.dat.tar.zst SOURCE_MODE=RSYNC_OVER_SSH  scalebox app create

SOURCE_URL=scalebox@159.226.237.136:10022/raid0/tmp/mwa/tar1301240224 TARGET_URL=/data1/mydata/mwa/tar FILE_NAME=FILE:filenames.txt SOURCE_MODE=RSYNC_OVER_SSH  scalebox app create

```

- KEEP_SOURCE_FILE test

```sh
ssh n0 'echo hello > /tmp/myfile.txt'
rm -f /tmp/myfile.txt
SOURCE_URL=root@10.0.6.101/tmp TARGET_URL=/tmp FILE_NAME=myfile.txt SOURCE_MODE=SSH KEEP_SOURCE_FILE=no scalebox app create

SOURCE_URL=root@10.0.6.101/tmp TARGET_URL=/tmp FILE_NAME=myfile.txt KEEP_SOURCE_FILE=no scalebox app create


SOURCE_URL=/tmp TARGET_URL=root@10.0.6.101/tmp FILE_NAME=myfile.txt KEEP_SOURCE_FILE=no scalebox app create

```

## 2. local to ssh-server
```sh
SOURCE_URL=/etc TARGET_URL=root@10.0.6.101/tmp FILE_NAME=postfix/master.cf scalebox app create

SOURCE_URL=/ TARGET_URL=root@10.0.6.101/tmp/ FILE_NAME=etc/postfix/master.cf scalebox app create

SOURCE_URL=/etc/postfix TARGET_URL=root@10.0.6.101/tmp FILE_NAME=master.cf scalebox app create

SOURCE_URL=/etc/postfix TARGET_URL=root@10.0.6.101/tmp TARGET_MODE=SSH FILE_NAME=master.cf scalebox app create

```

- KEEP_SOURCE_FILE test

```sh
SOURCE_URL=/tmp TARGET_URL=root@10.0.6.101/tmp FILE_NAME=master.cf KEEP_SOURCE_FILE=no scalebox app create

SOURCE_URL=/tmp TARGET_URL=root@10.0.6.101/tmp FILE_NAME=postfix/master.cf KEEP_SOURCE_FILE=no scalebox app create

SOURCE_URL=/tmp TARGET_URL=root@10.0.6.101/tmp FILE_NAME=etc/postfix/master.cf KEEP_SOURCE_FILE=no scalebox app create
```


```sh
echo "hello" > /tmp/myfile.txt

SOURCE_URL=/tmp TARGET_URL=root@10.0.6.101/tmp FILE_NAME=myfile.txt TARGET_MODE=SSH KEEP_SOURCE_FILE=no scalebox app create


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
SOURCE_URL=/tmp TARGET_URL=rsync://root:cas12345@10.0.6.101:50873/tmp FILE_NAME=myfile.txt scalebox app create

SOURCE_URL=/tmp TARGET_URL=rsync://root:cas12345@10.0.6.101:50873/tmp FILE_NAME=a/myfile.txt scalebox app create

```

## 5. rsync-server to local
```sh
ssh n0 'echo hello > /tmp/myfile.txt'
rm -f /tmp/myfile.txt

SOURCE_URL=rsync://root:cas12345@10.0.6.101:50873/tmp TARGET_URL=/tmp FILE_NAME=myfile.txt scalebox app create

SOURCE_URL=rsync://root:cas12345@10.0.6.101:50873/tmp TARGET_URL=/tmp FILE_NAME=myfile.txt KEEP_SOURCE_FILE=no scalebox app create

SOURCE_URL=rsync://root:cas12345@10.8.1.78:50873/tmp TARGET_URL=/tmp FILE_NAME=myfile.txt scalebox app create

```
