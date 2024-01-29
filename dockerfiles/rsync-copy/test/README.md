# examples for dir-list & rsync-copy

## 1. ssh-server to local

```sh

FILE_NAME=scalebox@10.255.128.1/etc~postfix/master.cf~/tmp scalebox app create

FILE_NAME=scalebox@10.255.128.1/etc/postfix~master.cf~/tmp scalebox app create

FILE_NAME=scalebox@10.255.128.1/~etc/postfix/master.cf~/tmp scalebox app create

```

## 2. local to ssh-server
```sh
FILE_NAME=/etc~postfix/master.cf~scalebox@10.255.128.1/tmp scalebox app create

FILE_NAME=/~etc/postfix/master.cf~scalebox@10.255.128.1/tmp scalebox app create

FILE_NAME=/etc/postfix~master.cf~scalebox@10.255.128.1/tmp scalebox app create
```

- KEEP_SOURCE_FILE test

```sh
touch /tmp/myfile.txt

FILE_NAME=/tmp~myfile.txt~scalebox@10.255.128.1/tmp scalebox app create

KEEP_SOURCE_FILE=no FILE_NAME=/tmp~myfile.txt~scalebox@10.255.128.1/tmp scalebox app create

FILE_NAME=/tmp~mydir/myfile.txt~scalebox@10.255.128.1/tmp scalebox app create

KEEP_SOURCE_FILE=no FILE_NAME=/tmp~mydir/myfile.txt~scalebox@10.255.128.1/tmp scalebox app create

```

## 3. rsync-server to local
```sh

```
## 4. local to rsync-server
```sh

```
