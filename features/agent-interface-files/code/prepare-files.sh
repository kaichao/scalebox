#!/bin/bash

# prepare files/dirs in /tmp/agent-interface-files for this test
# 创建目录结构
TMP_DIR="/tmp/agent-interface-files"
mkdir -p "$TMP_DIR/input-dir/sub-dir1"
mkdir -p "$TMP_DIR/input-dir/sub-dir2"
mkdir -p "$TMP_DIR/input-dir/sub-dir3"
mkdir -p "$TMP_DIR/output-dir/sub-dir1"
mkdir -p "$TMP_DIR/output-dir/sub-dir2"

# 创建input-dir下的文件
# sub-dir1: 2个文件，每个1kB
dd if=/dev/zero of="$TMP_DIR/input-dir/sub-dir1/file10.txt" bs=1024 count=1 2>/dev/null
dd if=/dev/zero of="$TMP_DIR/input-dir/sub-dir1/file11.txt" bs=1024 count=1 2>/dev/null

# sub-dir2: 2个文件，每个2kB
dd if=/dev/zero of="$TMP_DIR/input-dir/sub-dir2/file20.txt" bs=1024 count=2 2>/dev/null
dd if=/dev/zero of="$TMP_DIR/input-dir/sub-dir2/file21.txt" bs=1024 count=2 2>/dev/null

# sub-dir3: 2个文件，每个3kB
dd if=/dev/zero of="$TMP_DIR/input-dir/sub-dir3/file30.txt" bs=1024 count=3 2>/dev/null
dd if=/dev/zero of="$TMP_DIR/input-dir/sub-dir3/file31.txt" bs=1024 count=3 2>/dev/null

# 创建output-dir下的文件
# sub-dir1: 2个文件，每个10kB
dd if=/dev/zero of="$TMP_DIR/output-dir/sub-dir1/file10.txt" bs=1024 count=10 2>/dev/null
dd if=/dev/zero of="$TMP_DIR/output-dir/sub-dir1/file11.txt" bs=1024 count=10 2>/dev/null

# sub-dir2: 2个文件，每个20kB
dd if=/dev/zero of="$TMP_DIR/output-dir/sub-dir2/file20.txt" bs=1024 count=20 2>/dev/null
dd if=/dev/zero of="$TMP_DIR/output-dir/sub-dir2/file21.txt" bs=1024 count=20 2>/dev/null
