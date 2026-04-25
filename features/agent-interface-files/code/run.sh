#!/bin/bash

echo "start,$(date +%Y-%m-%dT%H:%M:%S.%6N)" > "${WORK_DIR}/timestamps.txt"

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
  iobytes:
    global_ibytes: 1024
    global_obytes: 2048
    ssd_ibytes: 3072
    ssd_obytes: 4096
    tmpfs_ibytes: 5120
    tmpfs_obytes: 6144
    input_bytes: 9216
    output_bytes: 12288
EOF

# timestamps.txt
echo "before-sleep,$(date +%Y-%m-%dT%H:%M:%S.%6N)" >> timestamps.txt
sleep 3
echo "after-sleep,$(date +%Y-%m-%dT%H:%M:%S.%6N)" >> timestamps.txt

# input-files.txt
cat > input-files.txt << 'EOF'
/tmp/agent-interface-files/input-dir/sub-dir1
/tmp/agent-interface-files/input-dir/sub-dir2
/tmp/agent-interface-files/input-dir/sub-dir3/file30.txt
/tmp/agent-interface-files/input-dir/sub-dir3/file31.txt
EOF

# output-files.txt
cat > output-files.txt << 'EOF'
/tmp/agent-interface-files/output-dir/sub-dir1
/tmp/agent-interface-files/output-dir/sub-dir2/file20.txt
/tmp/agent-interface-files/output-dir/sub-dir2/file21.txt
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
