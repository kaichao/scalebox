# examples for list-dir

## 1. local path
```sh
SOURCE_URL=/ DIR_NAME=etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create
DIR_NAME=/#etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create
DIR_NAME=/etc#postfix REGEX_FILTER=^.*cf\$ scalebox app create
DIR_NAME=/etc/postfix#. REGEX_FILTER=^.*cf\$ scalebox app create
```

## 2. rsync-over-ssh
```sh
SOURCE_URL=scalebox@10.255.128.1/etc DIR_NAME=postfix REGEX_FILTER=^.*cf\$ scalebox app create

DIR_NAME=scalebox@10.255.128.1/etc#postfix REGEX_FILTER=^.*cf\$ scalebox app create
DIR_NAME=scalebox@10.255.128.1:22/etc#postfix REGEX_FILTER=^.*cf\$ scalebox app create
DIR_NAME=scalebox@10.255.128.1/etc/postfix#. REGEX_FILTER=^.*cf\$ scalebox app create
DIR_NAME=scalebox@10.255.128.1/#etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create

```

## 3. rsync-native
- anonymous
```sh
# ncbi data
SOURCE_URL=rsync://ftp.ncbi.nlm.nih.gov/1000genomes/ftp DIR_NAME=. scalebox app create

SOURCE_URL=rsync://ftp.ncbi.nlm.nih.gov/1000genomes/ftp DIR_NAME=release/2008_12 scalebox app create

SOURCE_URL=rsync://ftp.ncbi.nlm.nih.gov/1000genomes DIR_NAME=. REGEX_FILTER=.*gz scalebox app create
```

- non-anonymous
```sh
SOURCE_URL=rsync://fast@fast.cstcloud.cn/doi/10.1038/s41586-021-03878-5 DIR_NAME=20191021 RSYNC_PASSWORD=<rsync-password> scalebox app create

DIR_NAME=rsync://fast:<rsync-password>@fast.cstcloud.cn/doi/10.1038/s41586-021-03878-5#20191021 scalebox app create
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

## 5. combined-message
```sh
SOURCE_URL= DIR_NAME=/etc%postfix scalebox app create

SOURCE_URL= DIR_NAME=ftp://ftp.ncbi.nlm.nih.gov/genbank/docs%. REGEX_FILTER=^.+\.txt\$ scalebox app create

```

## 6. 2D-dataset
```sh
DIR_NAME=/raidz/fast-fz/ZD2022_1_1%Dec+4120_03_03/20221102 \
REGEX_2D_DATASET='^(.+%)?([^/]+/[^/]+)/.+M([0-9]+)_([0-9]+).+$' \
INDEX_2D_DATASET='2,3,4' \
scalebox app create

```
