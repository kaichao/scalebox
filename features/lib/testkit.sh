#!/bin/bash
# feature 集成测试公共库。
#
# 契约：test.sh 依次调用 tk_init → tk_run_app → 断言；成功返回 0，失败返回
# 非 0。tk_init 注册 EXIT 清理：删除本次 app 的记录，并停止同一 APP_ID 的
# slot 容器（SCALEBOX_TEST_KEEP=1 可保留现场）。
#
# 运行位置：控制面所在机器（需要 scalebox CLI 与数据库访问）。本机 Mac 没有
# actuator，不能直接跑；用 `ssh fxcorr` 在测试机运行。
#
# 环境变量：
#   SCALEBOX_CLI                     scalebox CLI（默认 scalebox）
#   SCALEBOX_DB_CONTAINER            数据库容器名（默认 database）
#   SCALEBOX_DB_USER / _NAME         psql 用户 / 库（默认 scalebox / scalebox）
#   SCALEBOX_DB_CMD                  自定义 SQL 执行命令；从 stdin 读 SQL，
#                                    优先于 SCALEBOX_DB_CONTAINER
#   SCALEBOX_TEST_TIMEOUT            等待任务终态的超时秒数（默认 120）
#   SCALEBOX_TEST_POLL               轮询间隔秒数（默认 2）
#   SCALEBOX_TEST_KEEP=1             保留 app 与容器，不清理
#   SCALEBOX_TEST_CLEAN_CONTAINERS=0 不停止 slot 容器（只删 app 记录）

TK_DB_CONTAINER=${SCALEBOX_DB_CONTAINER:-database}
TK_DB_USER=${SCALEBOX_DB_USER:-scalebox}
TK_DB_NAME=${SCALEBOX_DB_NAME:-scalebox}

tk_fail() {
    echo "FAIL [$TK_NAME] $*" >&2
    exit 1
}

# 执行 SQL，psql 以 -tA 输出（无表头、无对齐）。
tk_db() {
    local sql=$1
    if [ -n "${SCALEBOX_DB_CMD:-}" ]; then
        printf '%s\n' "$sql" | sh -c "$SCALEBOX_DB_CMD"
    else
        printf '%s\n' "$sql" | docker exec -i "$TK_DB_CONTAINER" \
            psql -U "$TK_DB_USER" -d "$TK_DB_NAME" -tA
    fi
}

tk_init() {
    TK_NAME=${TK_NAME:-$(basename "$PWD")}
    TK_CLI=${SCALEBOX_CLI:-scalebox}
    TK_TIMEOUT=${SCALEBOX_TEST_TIMEOUT:-120}
    TK_POLL=${SCALEBOX_TEST_POLL:-2}
    TK_KEEP=${SCALEBOX_TEST_KEEP:-}
    TK_APP_ID=
    TK_LOG_DIR="${TMPDIR:-/tmp}/scalebox-feature-tests/${TK_NAME}.$$"
    mkdir -p "$TK_LOG_DIR"

    [ -f feature.yaml ] || tk_fail "当前目录没有 feature.yaml"
    command -v "$TK_CLI" >/dev/null 2>&1 || \
        tk_fail "找不到 $TK_CLI；请在控制面所在机器运行，或用 SCALEBOX_CLI 指定"

    local one
    one=$(tk_db "select 1;")
    [ "$one" = "1" ] || tk_fail "数据库不可访问；请设置 SCALEBOX_DB_CMD（见 lib/testkit.sh 头部）"

    trap tk_cleanup EXIT
    echo "[$TK_NAME] log: $TK_LOG_DIR"
}

# 创建 app。env_file 可选；feature.yaml 中的 ${PWD} 需要当前目录正确。
tk_run_app() {
    local env_file=${1:-} out rc
    if [ -n "$env_file" ]; then
        out=$("$TK_CLI" run --env-file="$env_file" --app-file=feature.yaml 2>&1)
    else
        out=$("$TK_CLI" run --app-file=feature.yaml 2>&1)
    fi
    rc=$?
    printf '%s\n' "$out" > "$TK_LOG_DIR/run.log"
    [ "$rc" -eq 0 ] || tk_fail "scalebox run 失败（rc=$rc，见 $TK_LOG_DIR/run.log）"

    TK_APP_ID=$(printf '%s\n' "$out" | \
        sed -n 's/.*"app_id"[[:space:]]*:[[:space:]]*\([0-9][0-9]*\).*/\1/p' | tail -1)
    [ -n "$TK_APP_ID" ] || tk_fail "未解析到 app_id（见 $TK_LOG_DIR/run.log）"
    echo "[$TK_NAME] app_id=$TK_APP_ID"
}

# 返回 "总数|未完成数"。终态 = status_code >= 0；-1/-2/-3 分别是
# READY/QUEUED/RUNNING。
tk_task_stats() {
    tk_db "select count(*), coalesce(sum(case when coalesce(status_code,-1) < 0 then 1 else 0 end),0) from t_task where mod_id in (select id from t_module where app=$TK_APP_ID);"
}

tk_task_count() {
    tk_db "select count(*) from t_task where mod_id in (select id from t_module where app=$TK_APP_ID);"
}

# 等待至少 expected 个任务，且全部进入终态（status_code >= 0）。
tk_wait_tasks() {
    local expected=$1 timeout=${2:-$TK_TIMEOUT} start=$SECONDS stats total pending
    while true; do
        stats=$(tk_task_stats 2>/dev/null || true)
        if [[ "$stats" =~ ^([0-9]+)\|([0-9]+)$ ]]; then
            total=${BASH_REMATCH[1]}
            pending=${BASH_REMATCH[2]}
            if [ "$total" -ge "$expected" ] && [ "$pending" -eq 0 ]; then
                echo "[$TK_NAME] tasks: total=$total pending=0"
                return 0
            fi
        fi
        if [ $((SECONDS - start)) -ge "$timeout" ]; then
            tk_fail "等待任务超时（期望 >=$expected，实际 ${stats:-无结果}）；任务：$(tk_db "select id,mod_id,coalesce(status_code,-1),left(body,40) from t_task where mod_id in (select id from t_module where app=$TK_APP_ID) order by id;")"
        fi
        sleep "$TK_POLL"
    done
}

tk_module_id() {
    tk_db "select id from t_module where app=$TK_APP_ID and name='$1' order by id limit 1;"
}

# 某个模块最近一次执行的 stdout。
tk_stdout_of_module() {
    tk_db "select coalesce(e.stdout,'') from t_task_exec e join t_module m on e.mod_id=m.id where m.app=$TK_APP_ID and m.name='$1' order by e.id desc limit 1;"
}

# 某个模块最近一次执行的退出码。
tk_status_of_module() {
    tk_db "select coalesce(e.status_code,-1) from t_task_exec e join t_module m on e.mod_id=m.id where m.app=$TK_APP_ID and m.name='$1' order by e.id desc limit 1;"
}

tk_assert_eq() {
    local expected=$1 actual=$2 msg=$3
    if [ "$expected" != "$actual" ]; then
        tk_fail "$msg：期望 [$expected]，实际 [$actual]"
    fi
}

tk_assert_contains() {
    local haystack=$1 needle=$2 msg=$3
    case "$haystack" in
        *"$needle"*) ;;
        *) tk_fail "$msg：输出中未找到 [$needle]；实际：[$haystack]" ;;
    esac
}

tk_assert_not_contains() {
    local haystack=$1 needle=$2 msg=$3
    case "$haystack" in
        *"$needle"*) tk_fail "$msg：输出中不应出现 [$needle]；实际：[$haystack]" ;;
        *) ;;
    esac
}

# 停止同一 APP_ID 的 slot 容器（仅 local 集群；远端节点跳过）。
tk_stop_containers() {
    command -v docker >/dev/null 2>&1 || return 0
    local c app
    for c in $(docker ps -q 2>/dev/null); do
        app=$(docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$c" 2>/dev/null | \
            sed -n 's/^APP_ID=//p' | head -1)
        if [ "$app" = "$TK_APP_ID" ]; then
            docker stop "$c" >/dev/null 2>&1 || true
        fi
    done
}

tk_cleanup() {
    local rc=$?
    trap - EXIT

    if [ -n "${TK_APP_ID:-}" ]; then
        if [ -n "${TK_KEEP:-}" ]; then
            echo "[$TK_NAME] KEEP app_id=$TK_APP_ID (log: $TK_LOG_DIR)"
        else
            if [ "${SCALEBOX_TEST_CLEAN_CONTAINERS:-1}" != "0" ]; then
                tk_stop_containers
            fi
            if "$TK_CLI" app delete --app-id "$TK_APP_ID" >/dev/null 2>&1; then
                echo "[$TK_NAME] cleanup: app $TK_APP_ID deleted"
            elif tk_db "delete from t_app where id=$TK_APP_ID;" >/dev/null 2>&1; then
                echo "[$TK_NAME] cleanup: app $TK_APP_ID deleted (db fallback)"
            else
                echo "[$TK_NAME] WARN: 清理 app $TK_APP_ID 失败" >&2
            fi
        fi
    fi

    if [ "$rc" -ne 0 ]; then
        echo "[$TK_NAME] log: $TK_LOG_DIR" >&2
    fi
    exit "$rc"
}
