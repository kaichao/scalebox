# rsyncd

rsync module name is '/local'
## rsync mode
```sh
make run-rsync-mode
```
## rsync over ssh
```sh
make run-rsync-over-ssh
```

## setup

uid in rsyncd.conf should have read and write permissions for the local directory
- set uid = root
- set local directory = 777
- map uid in rsyncd.conf to an outside user
