# 逗号分隔字符串转数组
```bash
# comma-seperated multi-message
arr=($(echo ${MESSAGE_BODY} | tr "," " "))

arr=($(/app/bin/url_parser ${SOURCE_URL}))

export MODE="${arr[0]}"

```

# 数组迭代处理
```bash
# comma-seperated multi-message
arr=($(echo ${MESSAGE_BODY} | tr "," " "))

for m in ${arr[@]}; do
    send-message $m
done
```

# bash逻辑条件判断

```bash
if [ "${MODE}" = "ERROR_FORMAT" ]; then
    echo "url format error, "${arr[1]} >&2
    exit 1
fi

```

# yaml解析
```bash
#!/bin/bash

IMAGE=hub.cstcloud.cn/scalebox/parser
PGHOST=localhost
PGPORT=5432

# Pass all environment variables except PATH to the container
docker run --network host --rm -v $(cd $(dirname $1);pwd):/data:ro \
    --env-file <(env|grep -v ^PATH=) -e PGHOST=${PGHOST} -e PGPORT=${PGPORT} \
    ${IMAGE} $(basename $1)
```

# timestamp in bash
```sh
#  rfc-3339 in seconds
$ date --rfc-3339=seconds | sed 's/ /T/'
2014-03-19T18:35:03-04:00

# rfc-3339 in milliseconds
$ date --rfc-3339=ns | sed 's/ /T/; s/\(\....\).*\([+-]\)/\1\2/g'
2014-03-19T18:42:52.362-04:00

# rfc-3339 in microseconds
$ date --rfc-3339=ns | sed 's/ /T/; s/\(\.......\).*\([+-]\)/\1\2/g'
2022-11-06T22:35:04.030320+08:00

# iso-8601 on Linux
$ date --iso-8601=ns
2022-12-02T08:54:04,170460829+0800

```
