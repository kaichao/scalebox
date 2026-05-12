#!/bin/bash

date +%Y-%m-%dT%H:%M:%S.%6N > "${WORK_DIR}/timestamps.txt"

code_dir=$(dirname $0)
# prepare files/dirs in /tmp/agent-interface-files for this test
$code_dir/prepare-files.sh

# 生成接口文件
# task-exec.yaml
cat > task-exec.yaml << 'EOF'
statusCode: 0
inputBytes: 102400
outputBytes: 204800
EOF

# sink-tasks.txt，支持7种格式
cat > sink-tasks.txt << 'EOF'
body1
{"bh0":"a","body":"body2"}
body3,{"h0":"a","h1":"b"}
{"bh0":"a","body":"body4"},{"h0":"a","h1":"b"}
mr-module,body5
mr-module,body6,{"h0":"a","h1":"b"}
mr-module,{"bh0":"a","body":"body7"},{"h0":"a","h1":"b"}
EOF


# extra-attrs.yaml
cat > extra-attrs.yaml << 'EOF'
  extra-attr-set:
    myattr1: 1024
    myattr2: 3.14
    myattr3: str
    myattr4: false
    myattr5: 2026-05-02T10:30:00Z
    arr_attr:
      - v0
      - v1
    obj_attr:
      id: 1
      key: mykey
      value: myval
EOF

# timestamps.txt
echo "$(date +%Y-%m-%dT%H:%M:%S.%6N),before-sleep" >> timestamps.txt
sleep 3
echo "$(date +%Y-%m-%dT%H:%M:%S.%6N),after-sleep" >> timestamps.txt

# input-files.txt
cat > input-files.txt << 'EOF'
/tmp/agent-interface-files/input-dir/sub-dir1
/tmp/agent-interface-files/input-dir/sub-dir2
/tmp/agent-interface-files/input-dir/sub-dir3/file30.txt
/tmp/agent-interface-files/input-dir/sub-dir3/file31.txt
/dev/shm/input-dir,1024
/shared/mydata/input-file,2048
EOF

# output-files.txt
cat > output-files.txt << 'EOF'
/tmp/agent-interface-files/output-dir/sub-dir1
/tmp/agent-interface-files/output-dir/sub-dir2/file20.txt
/tmp/agent-interface-files/output-dir/sub-dir2/file21.txt
/dev/shm/output-dir,1024
/shared/mydata/output-file,2048
EOF

# network-files.txt
cat > network-files.txt << 'EOF'
from,10.0.6.100,/tmp/agent-interface-files/input-dir/sub-dir1
from,10.0.6.100,1024
to,10.0.6.100,/tmp/agent-interface-files/output-dir/sub-dir1
to,10.0.6.101,1024
EOF

# removed-files.txt
cat > removed-files.txt << 'EOF'
/tmp/agent-interface-files/input-dir/sub-dir2
/tmp/agent-interface-files/input-dir/sub-dir3/file30.txt
/tmp/agent-interface-files/input-dir/sub-dir3/file31.txt
EOF

# auxout.txt
v="my output"
cat > auxout.txt << EOF
user-defined output
v=$v
EOF

exit 0
