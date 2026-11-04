# task-aggregation

把**大量极短的算法任务**在算法侧聚合成**一个粗粒度任务**（任务内部循环处理所有成员），
完成后用 `sink_modules` **直连下游**，从而把穿过主路由的回流次数从 N 降到 0。

用来解决"主路由单实例吞吐有限 + 同一来源只能单实例"造成的架构性上限——加节点无效，
因为排队量随节点数一起涨。原理见 `docs/cn/source/design/1_programming_model.md` §1.7.5；
操作步骤见 [`recipes/bypass-router-short-tasks`](../../docs/agents/recipes/bypass-router-short-tasks.md)。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `sink_modules` | 聚合模块的 `module` 条目 | 模块名列表 | — | 直连的下游；容器内 `SINK_MODULE` 取**第一个** |
| `code_path` | 聚合模块的 `arguments` | 路径 | — | 模块代码；`run.sh` 里写"组内循环" |
| `task_batch_size` | 聚合模块的 `parameters` | 整数 | 1 | 建议 1（逐条领取），让"任务数 = 组数"在库里可见 |
| — | 聚合模块的 `run.sh` | — | — | 用 `$2` 解析上游下发的成员清单 header；成功才写 `sink-tasks.txt` |

**不需要**新增平台参数：本模式完全由"模块代码粒度 + `sink_modules`"两件事构成。

## 行为与不变式

- 聚合任务的 `t_task` 条数 = 组数；下游任务的 `t_task` 条数 = 组数（一对一直连）。
- 下游任务的 `t_task.from_mod` 是**聚合模块名**，不是主路由名——这是"没有经 router 中转"
  的库内证据。
- `sink_modules` 由 agent 体现在环境变量 `SINK_MODULE` 上（取第一个），`run.sh` 用它写
  `sink-tasks.txt` 的第一列。
- 失败语义**必须**是"任一成员失败 → 整组失败"：不写 `sink-tasks.txt`，下游不启动。这与
  原来"N 个任务各自成功才凑齐下游输入"的语义一致；想让失败粒度更细，就拆小聚合粒度。
- 主路由不再是这批任务的必经之路，但**批级入口任务仍走 router**（`from_head` 之类）。
  本模式只去掉"回流"，不去掉"入口"。
- 主路由原本承担的副作用要一起搬走：批级计数器/信号量、中间产物 purge 路径、下游需要的
  header。漏一个就表现为"下游起不来"或"磁盘涨"。

## 坑

- **只改了 `sink_modules` 没改任务粒度**：回流量不变，只是换了投递目标，没有收益。本模式的
  收益全部来自"N 个短任务 → 1 个长任务"。
- **聚合任务变成新长尾**：一个任务跑太久，下游一直等它，且失败重跑代价大。聚合粒度控制在
  "单任务时长 ≈ 下游单任务时长"量级。
- **成员清单让它自己去猜**：不要让聚合模块扫目录推断成员，上游必须通过 header 显式下发；
  否则并发时容易把别的组的产物算进来。
- **purge 路径没下发**：原来的清理是 router 投下游时带的，改直连后没人带 → 中间产物堆满盘。
- **以为拆 router 多实例也能解决**：同一 `from_module` 只能单实例（编程约束），拆实例对这条
  来源无效；见 [`router-multi-instance`](../router-multi-instance/README.md) 的适用范围。
- **在 `run.sh` 里用 `set -e`**：组内循环任一成员失败会让脚本静默退出、`sink-tasks.txt` 写一半；
  改成显式收集返回码再统一 `exit`（本 feature 的 `code/agg/run.sh` 是范式）。

## 最小例子

见 `feature.yaml`。关键片段：

```yaml
modules:
  agg:
    arguments:
      code_path: ${PWD}/code/agg
    parameters:
      task_batch_size: 1
    sink_modules:
      - sink                 # 直连下游，不经 main_router
    environments:
      - AGG_STEPS=8

  sink:
    arguments:
      code_path: ${PWD}/code/sink
```

`code/agg/run.sh` 的核心三行：

```bash
for m in "${list[@]}"; do do_one_unit "$body" "$m" || { rc=$?; break; }; done
[ "$rc" -eq 0 ] || exit "$rc"
printf '%s,%s\n' "${SINK_MODULE}" "$body" > "${WORK_DIR}/sink-tasks.txt"
```

## 测试

```sh
./test.sh
```

按 [`lib/testkit.sh`](../lib/testkit.sh) 的契约：`tk_init` → `tk_run_app` → 投 3 个组 →
`tk_wait_tasks 6` → 断言。预期 PASS。断言点：

- 聚合模块任务数 = 3（不是 3 × 成员数）
- 下游任务数 = 3，且 `from_mod` 全为 `agg`（直连证据）
- 聚合任务内部执行步数 = 3 组 × 3 成员 × 8 步 = 72（聚合真的发生）

## 实证（app-fxcorr，2026-11-04）

4 节点、0.512 s × 120 批。上游每批产出 32 个毫秒级短任务并回流主路由：

| 指标 | 逐任务回流 | 聚合 + 直连 |
|---|---|---|
| 主路由回流条数 | 3840 | **0** |
| 主路由排队 avg / max | 96.67 s / — | **1.14 s / 6 s** |
| 下游 slot 利用率 | 62.5% | 70% |
| 全链 wall | 609.8 s | **407.8 s（−33%）** |

同轮 2 节点对照主路由排队只有 3.8 s——排队量随节点数增长，是本模式的判据。

## 源码指针

- `go-scalebox/cmd/agent/sink_tasks.go`（`addSinkTasks`）—— 读 `sink-tasks.txt` 投下游、
  填 `from_module`。
- `go-scalebox/cmd/agent/run.go` —— 注入 `SINK_MODULE` / `FROM_MODULE` 等环境变量。
- `go-scalebox/cmd/controld/task/add_local.go`（`cacheMainRouterSlotID`）—— 主路由多实例
  分流（另一条出路）。
- `scalebox/docs/cn/source/design/1_programming_model.md` §1.7.5 —— 主路由单实例的
  吞吐上限与两条出路。

## 相关

- 操作手册：[`recipes/bypass-router-short-tasks`](../../docs/agents/recipes/bypass-router-short-tasks.md)
- 另一条出路：[`router-multi-instance`](../router-multi-instance/README.md)
- 回流格式：[`docs/reference/task-headers.md`](../../docs/reference/task-headers.md)、
  [`docs/reference/module.md`](../../docs/reference/module.md)
- 规则：`docs/agents/10-invariants.md` 的 INV-24
