# examples for list-dir & rsync-copy

## 1. ssh-server to local

```sh

SOURCE_URL=root@$(hostname -i)/etc TARGET_URL=/tmp/etc DIR_NAME=. REGEX_FILTER=^.+\.cf\$ scalebox app create

SOURCE_URL= TARGET_URL=/tmp/etc DIR_NAME=root@$(hostname -i)/etc%. REGEX_FILTER=^.+\.cf\$ scalebox app create

```

## 2. local to ssh-server
```sh

```

## 3. rsync-server to local
```sh

```
## 4. local to rsync-server
```sh

```
