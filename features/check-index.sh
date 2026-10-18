#!/usr/bin/env bash
# 校验 features/README.md 索引与实际目录的一致性，以及每个 feature 的必备文件。
# 用法：./check-index.sh
set -u
cd "$(dirname "$0")"

errors=0

# 实际 feature 目录（排除 _ 开头的模板/辅助目录与 lib）
actual=$(for d in */; do
    d=${d%/}
    case "$d" in _*|lib) continue ;; esac
    echo "$d"
done | sort)

# 索引中引用的目录（格式：](./name/)）
indexed=$(grep -oE '\]\(\./[^/()]+/\)' README.md 2>/dev/null \
    | sed -E 's#^\]\(\./([^/]+)/\)$#\1#' | sort -u)

# 1. 索引引用了不存在的目录
while IFS= read -r name; do
    [ -z "$name" ] && continue
    if [ ! -d "./$name" ]; then
        echo "ERROR: index references missing directory: $name"
        errors=$((errors + 1))
    fi
done <<< "$indexed"

# 2. 目录存在但未登记
while IFS= read -r name; do
    [ -z "$name" ] && continue
    if ! grep -qF "](./$name/)" README.md 2>/dev/null; then
        echo "ERROR: directory not in index: $name"
        errors=$((errors + 1))
    fi
done <<< "$actual"

# 3. 每个 feature 的必备文件
for name in $actual; do
    [ -f "./$name/README.md" ] || {
        echo "ERROR: $name/README.md missing"
        errors=$((errors + 1))
    }
    if [ -f "./$name/feature.yaml" ]; then
        :
    elif [ -f "./$name/app.yaml" ]; then
        echo "WARN: $name uses legacy app.yaml (rename to feature.yaml)"
    else
        echo "ERROR: $name/feature.yaml missing"
        errors=$((errors + 1))
    fi
    if [ -f "./$name/test.sh" ] && [ ! -x "./$name/test.sh" ]; then
        echo "WARN: $name/test.sh is not executable"
    fi
done

count=$(printf '%s\n' "$actual" | grep -c . || true)
if [ "$errors" -gt 0 ]; then
    echo "FAILED: $errors error(s)"
    exit 1
fi
echo "OK: $count features, index consistent"
