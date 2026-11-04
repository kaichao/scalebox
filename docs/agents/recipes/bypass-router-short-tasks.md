# recipe — bypass-router-short-tasks（短任务聚合 + 直连下游）

> 用于：主路由成为瓶颈——某个来源的回流任务数量巨大、单任务极短，router 排队时间随
> 节点数一起涨，加节点不提速。
>
> 本 recipe 讲**怎么改应用**；为什么主路由有这个上限，见
> `docs/cn/source/design/1_programming_model.md` §1.7.5。

## 前置

- 已有 app 能跑通，且**已定位到主路由**，不是下游模块不够（否则先按
  [`flow-control`](./flow-control.md) 加下游 slot）。
- 已按 [`20-diagnosis.md`](../20-diagnosis.md) 看过任务级数据：知道哪个 `from_module`
  的回流最多、排队多久。

## 何时用它

三条**同时**成立才动手，缺一条就先别改结构：

| 判据 | 怎么量 | 典型值 |
|---|---|---|
| 该来源回流量远大于其他来源 | 按 `t_task.from_mod` 分组数回流条数 | 某来源占 router 总量的 > 50% |
| 单任务执行时间与调度开销同量级 | `t_task_exec` 的 `ts_finished - ts_started` | 毫秒~秒级，且远小于下游单任务 |
| router 排队随规模恶化 | `ts_started - t_task.create_time`，对比 1/2/4 节点 | 1 节点几秒、4 节点上百秒 |

若只是"某个来源特别重、但任务不算短"，先试
[`router-multi-instance`](../../../features/router-multi-instance/README.md)（按来源拆实例）；
两条路互不冲突，通常是先拆实例、仍不够再聚合。

## 原理

主路由单实例的吞吐 ≈ `1 / 单任务处理时间`，且**同一个 `from_module` 只能单实例**
（处理逻辑共享状态，属于编程模型约束，不是平台限制）。所以任务越短，上限越低，
而这个上限**不随节点数增长**——加节点只会让上游更快地产出回流，排队反而更长。

出路只有两条：**减少穿过主路由的次数**、或**拆分来源**。本 recipe 是第一条：

```text
改造前：上游 ──N 个短任务──▶ 主路由 ──N 个下游任务──▶ 下游
改造后：上游 ──M 个粗粒度任务──▶ 下游        （M ≪ N，主路由 0 次）
```

聚合后的任务在**算法侧**完成原来 N 个短任务的工作（模块 `run.sh` 里写循环），完成后
通过 `sink_modules` 直接投下游，`sink-tasks.txt` 的第一列就是下游模块名。

## 步骤

1. **量化现状并留档**（改完要对比）：

```text
按 from_module 统计：回流条数、排队 avg/p90（ts_started - t_task.create_time）、
执行 avg（ts_finished - ts_started）；再记一份下游 slot 利用率与全链 wall。
```

2. **选聚合维度**。取"数据上天然成组、下游也按组消费"的那一层，别强行拼：

| 可用的分组依据 | 例子 |
|---|---|
| 同一批 / 同一时间片 | 一个批的所有分片 |
| 同一文件 / 同一目录 | 一个目录下的全部小文件 |
| 同一组（组内元素对下游是一次性输入） | 一个频段组 |

聚合粒度的上下界：**下界**=1 个任务至少覆盖到"下游一次消费"。**上界**=单任务时长不要
超过下游单任务的量级，否则聚合任务自身变成新的长尾，并且失败重跑代价变大。

3. **改算法模块**：一个任务处理全部成员，成员清单由上游用 header 下发（逗号分隔即可），
   不要让它去猜目录结构。

```bash
# code/run.sh（骨架）
#!/bin/bash
set -uo pipefail
body=$1; headers=$2
members=$(printf '%s' "$headers" | sed -n 's/.*"members":"\([^"]*\)".*/\1/p')

rc=0
IFS=',' read -ra list <<< "$members"
for m in "${list[@]}"; do
    do_one_unit "$body" "$m" || { rc=$?; break; }   # 任一成员失败即整组失败
done

# 成功才写；第一列是 sink_modules 里的下游模块名（INV-04）
[ "$rc" -eq 0 ] && printf '%s,%s\n' "${SINK_MODULE}" "$body" > "${WORK_DIR}/sink-tasks.txt"
exit "$rc"
```

4. **改接线**：把该模块的 `sink_modules` 直接指向下游，**不要**回主路由。

```yaml
modules:
  aggregate-stage:
    sink_modules:
      - downstream-stage     # 直连：不再经 main_router
```

5. **迁移主路由原本承担的副作用**——漏掉任何一个都会表现为"下游起不来"：

| 副作用 | 原来谁做 | 改造后 |
|---|---|---|
| 批级计数器 / 信号量 | router 在 head 回流时创建 | head 分支里创建，或在聚合模块投下游前创建 |
| 中间产物清理（purge 路径） | router 在投下游时下发 | 聚合模块写 `sink-tasks.txt` 时把路径写进第三个字段 |
| 下游需要的 header（成员清单、状态） | router 拼装 | 聚合模块随 `sink-tasks.txt` 的 JSON 字段下发 |

6. **重建 app**（配置在创建时快照，INV-01），跑同一批数据，按判据验收。

## 判据

| 指标 | 期望 |
|---|---|
| 主路由中该来源的回流条数 | **0**（按 `t_task.from_mod` 统计） |
| 聚合模块任务数 | = 批数/组数，不再是 批数 × 分片数 |
| 下游任务的 `from_mod` | = 聚合模块名（证明是直连，不是经 router 中转） |
| 主路由排队时间 | 回落到与任务量无关的量级（毫秒~秒） |
| 下游 slot 利用率 / 全链 wall | 利用率升、wall 降 |

任一条不满足 → 见下表。

## 失败下一步

| 现象 | 原因 | 下一步 |
|---|---|---|
| 主路由仍有该来源回流 | `sink_modules` 没改，或聚合模块没写 `sink-tasks.txt` | 查 `task show` 的 `from_mod`；确认模块的 `SINK_MODULE` 环境变量 |
| 下游一条任务都没有 | `sink-tasks.txt` 空文件 / 第一列模块名拼错 | 空文件会被判后处理错误；第一列必须与 `modules:` 键完全一致 |
| 下游 READY 不动 | 原来由 router 建的计数器没迁过来 | 回到步骤 5，检查计数器初值（INV-08） |
| 中间产物没释放，磁盘涨 | purge 路径没随任务下发 | 把组目录写进 `sink-tasks.txt` 的第三个字段 |
| 聚合任务成新长尾 | 聚合粒度太粗 | 缩小每组成员数，或按批长拆两层 |
| 失败重跑代价大 | 一个成员失败整组重跑 | 组内成员数控制在下游单任务量级；需要细粒度重试时改为两层聚合 |

## 实证（app-fxcorr，2026-11-04）

下游单任务 57 s，上游每个批产出 32 个"毫秒级"短任务再回流主路由。4 节点、120 批：

| 指标 | 逐任务回流 | 聚合 + 直连 |
|---|---|---|
| 主路由回流条数 | 3840 | **0** |
| 主路由排队 avg | 96.67 s | **1.14 s**（max 6 s） |
| 下游 slot 利用率 | 62.5% | **70%**（后续由 vtask 名额决定） |
| 全链 wall | 609.8 s | **407.8 s（−33%，1.50×）** |

同轮 2 节点对照：主路由排队只有 3.8 s——**排队量随节点数增长**正是本模式的判据。

## 相关

- 原理与约束：`docs/cn/source/design/1_programming_model.md` §1.7.5（主路由多实例的编程约束）
- 同一问题的另一条出路：[`router-multi-instance`](../../../features/router-multi-instance/README.md)
- 可跑最小形态：[`features/task-aggregation`](../../../features/task-aggregation/README.md)
- 回流格式与 header：[`../../reference/task-headers.md`](../../reference/task-headers.md)、
  [`../../reference/module.md`](../../reference/module.md)（`sink_modules` / `SINK_MODULE`）
