#!/bin/bash

m=$1
input_files="\"/etc/passwd\",\"/etc/group\""
output_files="\"/etc/fstab\",\"/etc/hosts\""
input_bytes=9999
output_bytes=999

user_text="user-defined text\nHello,scalebox!"

cat << EOF > /work/timestamps.txt
2008-03-19T18:35:03-08:00
2009-11-05T17:50:20.154+08:00
2010-11-05T17:50:20.154918+08:00
2011-11-05T17:50:20.154918780+08:00
2012-11-17T08:52:21,963572856+08:00
EOF

cat << EOF > /work/user-file.txt
This is user-defined data in a file.
Multi-line is supported.
EOF


if [ "$m" = "0" ]; then
    rm -f /work/timestamps.txt /work/user-file.txt
    echo "stdout in message-${m}."
    echo "stderr in message-${m}." >&2
cat << EOF > /work/task-exec.json
{
    "statusCode":0,
	"inputBytes":${input_bytes},
	"outputBytes":${output_bytes},
    "userText":"user-defined text\nHello scalebox in message-${m}",
    "timestamps":["2018-03-19T18:35:03-08:00","2019-11-05T17:50:20.154+08:00","2020-11-05T17:50:20.154918+08:00","2021-11-05T17:50:20.154918780+08:00","2022-11-17T08:52:21,963572856+08:00"],
    "sinkJob":"task-exec-files",
    "messageBody":"1"
}
EOF

elif [ "$m" = "1" ]; then
    echo "stdout in message-${m}."
    echo "stderr in message-${m}." >&2

cat << EOF > /work/task-exec.json
{
    "statusCode":0,
	"inputFiles":[${input_files}],
	"outputFiles":[${output_files}],
    "sinkJob":"task-exec-files",
    "messageBody":"2"
}
EOF

elif [ "$m" = "2" ]; then
    rm -f /work/task-exec.json
    echo $m
else
    echo $m
fi

exit 0
