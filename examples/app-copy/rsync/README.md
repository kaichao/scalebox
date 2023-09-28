
## rsync over ssh
- PUSH
```sh
SOURCE_URL=/etc TARGET_URL=scalebox@10.255.128.1/tmp/etc3 scalebox app create 
```
- PULL
```sh
SOURCE_URL=scalebox@10.255.128.1/etc TARGET_URL=/tmp/etc5 scalebox app create 
```

## rsync
- PUSH
```sh
SOURCE_URL=/etc TARGET_URL=scalebox@10.255.128.1/tmp/etc3 scalebox app create 
```
- PULL
```sh
SOURCE_URL=scalebox@10.255.128.1/etc TARGET_URL=/tmp/etc5 scalebox app create 
```
