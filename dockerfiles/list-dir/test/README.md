# examples for list-dir

## 1. local path
```sh
SOURCE_URL=/ DIR_NAME=etc/postfix scalebox app create

SOURCE_URL=/etc/postfix DIR_NAME=. scalebox app create

SOURCE_URL=/ DIR_NAME=etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create
```

## 2. rsync-over-ssh
```sh
SOURCE_URL=scalebox@10.255.128.1/ DIR_NAME=etc/postfix REGEX_FILTER= scalebox app create

SOURCE_URL=scalebox@10.255.128.1/etc/postfix DIR_NAME=. REGEX_FILTER= scalebox app create

SOURCE_URL=scalebox@10.255.128.1/ DIR_NAME=etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create
```

## 3. rsync-native
- anonymous
```sh
# ncbi data
SOURCE_URL=rsync://ftp.ncbi.nlm.nih.gov/1000genomes DIR_NAME=. scalebox app create

SOURCE_URL=rsync://ftp.ncbi.nlm.nih.gov/1000genomes DIR_NAME=. REGEX_FILTER=.*gz scalebox app create
```

- non-anonymous
```sh
SOURCE_URL=rsync://fast@fast.cstcloud.cn/doi/10.1038/s41586-021-03878-5 DIR_NAME=20191021 RSYNC_PASSWORD=<rsync-password> scalebox app create
```


## 4. ftp

- anonymous ftp
```sh
SOURCE_URL=ftp://ftp.ncbi.nlm.nih.gov/1000genomes/ftp/release/2008_12 DIR_NAME=. scalebox app create

SOURCE_URL=ftp://ftp.ncbi.nlm.nih.gov/1000genomes/ftp/release/ DIR_NAME=2008_12 scalebox app create

SOURCE_URL=ftp://ftp.ncbi.nlm.nih.gov/1000genomes/ftp/release/ DIR_NAME=2008_12 REGEX_FILTER=^.*gz\$ scalebox app create
```

- non-anonymous ftp
```sh
SOURCE_URL=ftp://<ftp-user>:<ftp-pass>@<ftp-host>/<ftp-path> DIR_NAME=. scalebox app create
```
