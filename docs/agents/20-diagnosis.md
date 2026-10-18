# 20 — 卡链与失败诊断

症状驱动的排查顺序。先看 app/module/slot/task 四层，再查单个任务的执行记录；规则细节见
[`10-invariants.md`](./10-invariants.md)。

## 60 秒检查

```sh
scalebox app show --app-id N              # app 状态
scalebox module list --app-id N           # 模块与 module-id
scalebox slot status --module-id M        # 该模块 slot 汇总
scalebox slot list --module-id M          # slot 明细（host/seq/status）
scalebox task list --app-id N             # 任务与 status_code
scalebox task show --task-id T            # 单任务详情（含 stdout/stderr）
```

数据库兜底（CLI 版本过旧或需要精确查询时）：

```sh
docker exec database psql -U scalebox -d scalebox -c \
  "select id,mod_id,body,status_code from t_task where mod_id in \
   (select id from t_module where app=N) order by id;"
```

## 症状 → 检查 → 修复

### app RUNNING，但任务没有进展

1. `scalebox task list --app-id N` 看是否有任务；没有任务 → 检查入口 `start_task` 或手动投递。
2. 任务停在 `-1`（READY）→ 看 slot：`scalebox slot status --module-id M`。
3. slot 为 0 行 → 见下一节“slot 全是 0 行”。

### slot 全是 0 行

最常见原因是 `slots` 没匹配到任何 host。

```sh
scalebox host list                         # 看 Hostname 与 Status
```

- `slots` 匹配的是 **hostname 第一段（短名）**，例如 hostname `h0.local` 要写 `h0`；
  写完整 FQDN `h0.local` 匹配不到（`expandRegexHostList` 用 `split_part(hostname,'.',1)`）。
- 多条 host 用正则/字符类（`n0-[01]`），**不要用逗号**。
- host 必须存在且 `status=ON`；僵尸节点见 INV-16。
- 单节点调试可改用 `arguments.slot_options: slot_on_head`（在头节点建一个 slot）。

修复后在 `app.yaml` 改 `slots`，**重建 app**（配置在创建时生效）；旧 app 的 slot 不会自动补。

### slot 有记录但一直 STARTING

- 用 `scalebox slot list --module-id M` 看 slot 明细（host/seq/status）。
- 检查执行节点：SSH/容器引擎是否可用；镜像是否存在（`base_image` 必须能在节点拉到）。
- 节点状态：`scalebox host list`，`ON` 不代表物理机可达（INV-16）。

### 任务 READY（-1）不动

- slot 未建 / 未 READY → 按上一节查 slot。
- slot READY 但任务不动 → 看 slot 是否属于正确的 module；`HOST-BOUND` 任务还要求目标 host 有该模块的 slot（见 `features/pod`）。
- 检查是否有空闲 slot 被其他 app 占用（`scalebox slot status --module-id M`）。

### 任务 RUNNING（-3）很久

- 单任务超时由 `task_max_seconds` 控制，超时返回 124；全局超时由 pgcron 按倍数恢复为 123。
- `-3` 是运行中，不是失败；先区分“任务真的在跑”与“进程已死但状态未更新”。
- 自定义 `ACTION_CHECK`/`ACTION_SETUP` 也被 `timeout` 包裹（见 `features/check-test`、`features/timeout`）。

### 任务失败码

| 码 | 含义 | 常见原因 |
|---|---|---|
| `112` | uniq_task 冲突 | 重复投递相同身份；重跑先清记录或显式处理冲突（INV-06/07） |
| `123` | 全局超时 | 任务长期运行态；检查 `task_global_timeout_scale` |
| `124` | 单任务超时 | `task_max_seconds` 到期；见 `features/timeout` |
| `1–10` | 全局语义 | 参数/输入/资源/依赖等，见 `docs/reference/exit-codes.md` |
| `-1/-2/-3` | 调度态 | READY/QUEUED/RUNNING，不是失败 |

### 下游任务不出现

- 非路由模块必须写 `${WORK_DIR}/sink-tasks.txt`；格式见 `features/headers-propagation` 与 `features/agent-interface-files`。
- 第一列是下游模块名；名字必须与 `modules:` 中的键完全一致。
- 上游任务失败（`status_code > 0`）不会回流下游（失败即停）。
- 经 router/wait-queue 时，`_` header 的自动继承有门控（INV-05/08）；需要显式转发。

### `_` header 丢失

- 非路由模块：`_` 前缀自动继承；普通 header 不继承。
- router 分支只放行 `vtask_role≠""` 的 sink；wait-queue 一跳需要显式转发。
- 脚本取值前先 `source /usr/local/lib/scalebox/functions.sh`，否则 `scalebox::task_header` 不存在。

### 空间闸门/流控卡住

- `free_space_gb` / `dir_quota_gb` 是“停领任务”，不是任务失败：任务保持 READY、fail=0。
- 任务级空间闸门 + 组内原子性会死锁；应挂批入口（vtask-head），见 INV-12。
- 检查顺序：`free_space_gb` → `dir_quota_gb` → `ACTION_CHECK`。

## 清理现场

```sh
scalebox app delete --app-id N
```

**注意**：`app delete` 只删数据库记录，本地集群的 slot 容器不会自动停止。需要按容器环境变量
`APP_ID=N` 找到并停止（测试机/本地集群示例）：

```sh
for c in $(docker ps -q); do
  app=$(docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$c" \
    | sed -n 's/^APP_ID=//p' | head -1)
  [ "$app" = "N" ] && docker stop "$c"
done
```

远端节点上的容器由该节点的 actuator/agent 管理，不能在控制面机器用本地 `docker` 清理。
