#!/bin/bash

# No directory capacity limit
[[ ! $DIR_LIMIT_GB ]] && exit 0

# DIR_LIMIT_GB=/Users/kaichao/WORK/workspace,4
dir0=$(echo $DIR_LIMIT_GB | cut -d "~" -f 1)
limit_gb=$(echo $DIR_LIMIT_GB | cut -d "~" -f 2)

dir=${dir0}
if [[ $dir =~ ^([^/].*)$ ]]; then
    # relative path, insert cluster_root
    cluster_root=$(scalebox cluster get-parameter --cluster $CLUSTER_NAME base_data_dir)
    code=$?
    [[ $code -ne 0 ]] && echo "[ERROR] get_cluster base_data_dir, cluster:$cluster, error_code:$code" >&2 && exit $code
    dir=$cluster_root/$dir
fi
dir="/local"$dir
echo [DEBUG]check-dir=$dir,limit-gb=${limit_gb}GB

mb=($(du -ms $dir)); code=$?
[[ $code -ne 0 ]] && echo "[ERROR] get dir MB, dir:$dir, error_code:$code" >&2 && exit $code
echo [DEBUG]check-dir=$dir,size=${mb}mb

# The total data size of files in the directory exceeds the limit
(( $mb > 1024 * $limit_gb )) && echo "[INFO]The space occupied by directory ${dir0} has exceeded ${limit_gb}GB."&& exit 1

exit 0
