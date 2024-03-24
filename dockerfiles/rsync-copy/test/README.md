# examples for dir-list & rsync-copy

## 1. ssh-server to local

```sh

FILE_NAME=root@10.0.6.101/etc~postfix/master.cf~/tmp scalebox app create

FILE_NAME=root@10.0.6.101/etc~postfix/master.cf~/tmp/etc scalebox app create

FILE_NAME=root@10.0.6.101/etc/postfix~master.cf~/tmp scalebox app create

FILE_NAME=root@10.0.6.101/etc/postfix~master.cf~/tmp scalebox app create



```

## 2. local to ssh-server
```sh
FILE_NAME=/etc~postfix/master.cf~root@10.0.6.101/tmp scalebox app create

FILE_NAME=/etc~postfix/master.cf~root@10.0.6.101/tmp/etc scalebox app create

FILE_NAME=/~etc/postfix/master.cf~root@10.0.6.101/tmp scalebox app create

FILE_NAME=/etc/postfix~master.cf~root@10.0.6.101/tmp scalebox app create
```

- KEEP_SOURCE_FILE test

```sh
touch /tmp/myfile.txt

FILE_NAME=/tmp~myfile.txt~root@10.0.6.101/tmp scalebox app create

KEEP_SOURCE_FILE=no FILE_NAME=/tmp~myfile.txt~root@10.0.6.101/tmp scalebox app create

FILE_NAME=/tmp~mydir/myfile.txt~root@10.0.6.101/tmp scalebox app create

KEEP_SOURCE_FILE=no FILE_NAME=/tmp~mydir/myfile.txt~root@10.0.6.101/tmp scalebox app create

```

## 3. rsync-server to local
```sh

```
## 4. local to rsync-server
```sh

```
