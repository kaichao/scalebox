#!/bin/bash
# 运行 features 下的集成测试。
#
#   ./run-tests.sh                 跑所有 */test.sh
#   ./run-tests.sh -k timeout      只跑名字匹配的
#   ./run-tests.sh --keep          保留 app 与容器
#   ./run-tests.sh --list          只列出测试
#
# 需在控制面所在机器运行（见 lib/testkit.sh）。
set -uo pipefail
cd "$(dirname "$0")"

pattern=""
keep=0
list_only=0
while [ $# -gt 0 ]; do
    case "$1" in
        -k) pattern=${2:-}; shift 2 ;;
        --keep) keep=1; shift ;;
        --list) list_only=1; shift ;;
        -h|--help) sed -n '1,15p' "$0"; exit 0 ;;
        *) echo "未知参数: $1" >&2; exit 2 ;;
    esac
done

if [ "$list_only" = 1 ]; then
    for t in */test.sh; do
        [ -f "$t" ] || continue
        echo "$(dirname "$t")"
    done
    exit 0
fi

total=0
pass=0
fail=0
skip=0
failed=""
start_all=$SECONDS

for t in */test.sh; do
    [ -f "$t" ] || continue
    if [ -n "$pattern" ] && [[ "$t" != *"$pattern"* ]]; then
        continue
    fi
    name=$(dirname "$t")
    total=$((total + 1))
    start=$SECONDS
    if [ "$keep" = 1 ]; then
        SCALEBOX_TEST_KEEP=1 bash "$t"
    else
        bash "$t"
    fi
    rc=$?
    dur=$((SECONDS - start))
    if [ "$rc" -eq 0 ]; then
        pass=$((pass + 1))
        printf 'PASS %s (%ss)\n' "$name" "$dur"
    elif [ "$rc" -eq 77 ]; then
        skip=$((skip + 1))
        printf 'SKIP %s (%ss)\n' "$name" "$dur"
    else
        fail=$((fail + 1))
        failed="$failed $name"
        printf 'FAIL %s (%ss)\n' "$name" "$dur"
    fi
done

echo "----------------------------------------"
if [ "$total" -eq 0 ]; then
    echo "没有找到 test.sh" >&2
    exit 1
fi
printf '%d/%d passed' "$pass" "$total"
[ "$skip" -gt 0 ] && printf '; skipped:%d' "$skip"
[ -n "$failed" ] && printf '; failed:%s' "$failed"
printf ' (%ss)\n' "$((SECONDS - start_all))"

[ "$fail" -eq 0 ]
