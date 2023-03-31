```sh
# local version
CLUSTER=local SOURCE_URL=/ DIR_NAME=etc/postfix REGEX_FILTER= scalebox app create app.yaml
CLUSTER=local SOURCE_URL=/etc/postfix DIR_NAME=. REGEX_FILTER= scalebox app create app.yaml
CLUSTER=local SOURCE_URL=/ DIR_NAME=etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create app.yaml

# rsync-over-ssh version
CLUSTER=local SOURCE_URL=scalebox@10.255.128.1/ DIR_NAME=etc/postfix REGEX_FILTER= scalebox app create app.yaml
CLUSTER=local SOURCE_URL=scalebox@10.255.128.1/etc/postfix DIR_NAME=. REGEX_FILTER= scalebox app create app.yaml
CLUSTER=local SOURCE_URL=scalebox@10.255.128.1/ DIR_NAME=etc/postfix REGEX_FILTER=^.*cf\$ scalebox app create app.yaml


# rsync-native version
CLUSTER=local SOURCE_URL=rsync://fast.cstcloud.cn/doi/10.1038/s41586-021-03878-5 DIR_NAME=20191021 REGEX_FILTER= scalebox app create app.yaml

CLUSTER=local SOURCE_URL=rsync://fast.cstcloud.cn/doi/10.1038/s41586-021-03878-5/20191021 DIR_NAME=. REGEX_FILTER= scalebox app create app.yaml

CLUSTER=local SOURCE_URL=rsync://fast@fast.cstcloud.cn/doi/10.1038/s41586-021-03878-5 DIR_NAME=20191021 REGEX_FILTER= RSYNC_PASSWORD=nao12345 scalebox app create app.yaml

CLUSTER=local SOURCE_URL=rsync://fast@fast.cstcloud.cn/files DIR_NAME=FRB121102/20190830 REGEX_FILTER= RSYNC_PASSWORD=nao12345 scalebox app create app.yaml


# ncbi data
CLUSTER=local SOURCE_URL=rsync://ftp.ncbi.nlm.nih.gov/1000genomes DIR_NAME=. REGEX_FILTER= RSYNC_PASSWORD= scalebox app create app.yaml

```
