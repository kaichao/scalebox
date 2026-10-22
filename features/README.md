# features（特性测试）

每个目录一个平台特性，三合一：面向 AI agent 的说明（`README.md`）、最小可跑应用（`feature.yaml`）、可判定的集成测试（`test.sh`，首批已提供，其余待补）。

> agent 入口见仓库根 [`AGENTS.md`](../AGENTS.md)；卡链排查见 [`docs/agents/20-diagnosis.md`](../docs/agents/20-diagnosis.md)。

## 约定

- 目录名：小写连字符（kebab-case），如 `task-batch-size`；单词目录直接用小写（如 `pod`）。
- `README.md`：按 [`_template/README.md`](./_template/) 的模板编写，中文。
- `feature.yaml`：与 `app.yaml` 同 schema；用 `scalebox run --app-file=feature.yaml` **显式**加载（`scalebox run` 的缺省探测只找 `app.yaml`）。
- `test.sh`：一键建 app → 跑 → 断言 → 清理，返回 0/非 0；公共函数见 [`lib/testkit.sh`](./lib/testkit.sh)，用 [`run-tests.sh`](./run-tests.sh) 批量运行。
- `feature.yaml` 的模块显式写 `base_image: scalebox.net/platform/agent:latest`；测试部署通常没有默认基础镜像。
- 索引一致性由 `./check-index.sh` 校验；新增/删除 feature 后运行。

## 索引

| feature | 覆盖特性 | 状态 |
|---|---|---|
| [agent-interface-files](./agent-interface-files/) | 模块与 agent 的 8 类接口文件 | `test.sh` ✓ |
| [check-test](./check-test/) | `free_space_gb` / `dir_quota_gb` / `ACTION_CHECK` | `test.sh` ✓ |
| [env-lifecycle](./env-lifecycle/) | app 快照 / 未定义变量 / non-tty stdin | `test.sh` ✓ |
| [headers-propagation](./headers-propagation/) | `_` 前缀传播 / `_vtask_attr_*` | `test.sh` ✓ |
| [node-progress-gap](./node-progress-gap/) | HOST-BOUND 节点进度差 | `test.sh` ✓（需 inline 多节点；bio 实跑 PASS） |
| [pod](./pod/) | `pod_id` 节点本地计算 | `test.sh` ✓（需 inline 多节点；bio 实跑 PASS） |
| [retry-test](./retry-test/) | `retry_rules` 退出码重试 | `test.sh` ✓（慢测：pgcron 5 分钟周期） |
| [side-car](./side-car/) | `ACTION_RUN/SETUP/TEARDOWN/CHECK` | `test.sh` ✓ |
| [slot-expr](./slot-expr/) | slot 启动表达式 | `test.sh` ✓（`@p:` 受限，见 README 坑） |
| [task-batch-size](./task-batch-size/) | `task_batch_size` 批量领取 | `test.sh` ✓ |
| [task-identity](./task-identity/) | `uniq_task` / 112 / `conflict-action` / `repeatable` | `test.sh` ✓ |
| [timeout](./timeout/) | `task_max_seconds` 超时 | `test.sh` ✓ |

## 规划中（未建目录）

来自 `docs/feature-inventory.md` §5：

| feature | 覆盖特性 |
|---|---|
| sink-tasks | 回流格式 / `last_from_module` / 与 vtask 的区别（部分已并入 agent-interface-files） |
| vtask-flow-control | wait-queue / bind / gate / fail |
| check-flow-control | 空间闸门挂载位置与死锁（可并入 check-test） |

## 测试运行

```sh
./run-tests.sh                 # 跑所有 test.sh
./run-tests.sh -k timeout      # 只跑名字匹配的
./run-tests.sh --keep          # 保留 app 与容器，便于排查
./run-tests.sh --list          # 只列出测试
```

**SKIP**：前置条件不满足时（如需要多节点控制面而当前只有 local 集群），`test.sh` 调用
`tk_skip` 以 77 退出，汇总显示 `SKIP`，不计失败。`pod` 与 `node-progress-gap` 要求
inline 集群控制面可达：单机 local 环境会 SKIP；在有 inline 多节点的测试机（如 bio / h0，
节点 `n0-0`、`n0-1`、`n1-0`、`n1-1`）上已实跑 PASS（2026-10-22）。

补跑命令：

```sh
ssh bio 'export PATH=/root/bin:$PATH; bash /root/scalebox/features/pod/test.sh'
ssh bio 'export PATH=/root/bin:$PATH; bash /root/scalebox/features/node-progress-gap/test.sh'
```

**自定义 app 文件**：`SCALEBOX_FEATURE_FILE=feature-test.yaml ./test.sh` 可加载同目录下的
变体文件（`slot-expr` 用它规避 `@p:` 的平台限制）。

**运行位置**：控制面所在机器（如测试机 `ssh fxcorr`）。需要 `scalebox` CLI 和数据库访问；默认通过 `docker exec database psql -U scalebox -d scalebox` 查询，可用 `SCALEBOX_DB_CMD` 覆盖。本机 Mac 没有 actuator，不能直接跑。

**`test.sh` 契约**：

1. `cd "$(dirname "$0")"` 后 `source ../lib/testkit.sh`。
2. `tk_init`：检查 CLI/数据库，注册 EXIT 清理。
3. `tk_run_app [env_file]`：运行 `scalebox run --app-file=feature.yaml`，解析 `app_id`。
4. `tk_wait_tasks <期望任务数> [超时]`：轮询到任务全部进入终态（`status_code >= 0`；`-1/-2/-3` 是 READY/QUEUED/RUNNING）。
5. `tk_assert_eq` / `tk_assert_contains` 断言；失败用 `tk_fail`。
6. 清理：`tk_init` 自动停 slot 容器并 `scalebox app delete`（CLI 不支持时回退数据库删除）；`SCALEBOX_TEST_KEEP=1` 保留现场。

## 校验

```sh
./check-index.sh
```
