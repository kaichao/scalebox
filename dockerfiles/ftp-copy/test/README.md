# examples for list-dir & ftp-copy

## 1. anonymous ftp-server to local

```sh

SOURCE_URL=ftp://ftp.ncbi.nlm.nih.gov/genbank/docs TARGET_URL=/tmp/genbank DIR_NAME=. REGEX_FILTER=^.+\.txt\$ scalebox app create

SOURCE_URL= TARGET_URL=/tmp/genbank DIR_NAME=ftp://ftp.ncbi.nlm.nih.gov/genbank/docs%. REGEX_FILTER=^.+\.txt\$ scalebox app create

```

## 2. non-anonymous ftp-server to local
```sh

```

## 3. local to ftp-server
```sh

```
