# examples for dir-list & ftp-copy

## 1. anonymous ftp-server to local

```sh

SOURCE_URL=ftp://ftp.ncbi.nlm.nih.gov/genbank/docs TARGET_URL=/tmp/genbank DIR_NAME=. REGEX_FILTER=^.+\.txt\$ scalebox run

SOURCE_URL= TARGET_URL=/tmp/genbank DIR_NAME=ftp://ftp.ncbi.nlm.nih.gov/genbank/docs%. REGEX_FILTER=^.+\.txt\$ scalebox run


FILE_NAME=ftp://ftp.ncbi.nlm.nih.gov/genbank~docs/FTv11.1.txt~/tmp scalebox run

```

## 2. non-anonymous ftp-server to local
```sh

```

## 3. local to ftp-server
```sh

```
