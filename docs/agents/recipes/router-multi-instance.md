# recipe — router-multi-instance（按来源拆主路由实例）

> 用于：主路由**单实例的处理时间**成为整轮 wall 的下限——回流总条数 × 单条处理时间
> 已经超过下游并行执行的下限，加节点不提速。
>
> 本 recipe 讲**怎么决定拆几个、每个放什么**。机制与配置字段的完整语义见
> [`features/router-multi-instance/README.md`](../../../features/router-multi-instance/README.md)；
> 若某个来源的任务本身极短且数量巨大，另一条路是
> [`bypass-router-short-tasks`](./bypass-router-short-tasks.md)（聚合并直连下游）。

## 前置

- 已按 [`20-diagnosis.md`](../20-diagnosis.md) 确认瓶颈在主路由，不是下游 slot 不够
  （下游 slot 不够时该来源的排队在**它自己**身上，不在 router）。
- 已按 `t_task.from_mod` 分组量过各来源的条数与执行时间（SQL 见下）。

## 何时用它

两条**同时**成立：

| 判据 | 怎么量 | 典型值 |
|---|---|---|
| router 纯处理时间 ≳ 整轮 wall | `SUM(end_time - begin_time)` 与 wall 比 | > 60% |
| 单来源贡献占比可分辨 | 按 `from_mod` 分组 | 有 2 个以上来源各占 > 10% |

**只做数量优化的情形**：如果 90% 的条数来自**同一个**来源，拆实例只是把它从一个进程
挪到另一个进程，总量不变——先看这篇；若该来源的单条执行时间又是毫秒级，
改走 [`bypass-router-short-tasks`](./bypass-router-short-tasks.md)。

## 步骤

### 1. 量化各来源

```sh
docker exec database psql -U scalebox -d scalebox -c "
SELECT t.from_mod AS src,
       count(*) n,
       round(avg(extract(epoch from (t.end_time - t.begin_time)))::numeric,4) avg_exec_s,
       round(sum(extract(epoch from (t.end_time - t.begin_time)))::numeric,1) total_exec_s,
       round(avg(extract(epoch from (t.begin_time - t.create_time)))::numeric,1) avg_queue_s
FROM t_task t JOIN t_module m ON m.id = t.mod_id
WHERE m.app = <app_id> AND m.name = '<main_router>'
  AND t.begin_time IS NOT NULL
GROUP BY 1 ORDER BY total_exec_s DESC;
"
```

**判据**：`total_exec_s` 之和与整轮 wall 同量级；`avg_queue_s` 随规模线性增长。

### 2. 定实例数

规则：**每个高负载来源独占一个实例，缺省实例只收零散流量**。

- 实例数 = 1（缺省）+ 需要隔离的来源个数
- 单个来源模块**只能钉一个实例**（`router_index` 是模块级参数），所以这是上限
- 零散来源（每实验几次的那种）留在 seq 0，不要为它们单独开实例

### 3. 改配置

主路由模块（`task_dist_mode: SLOT-BOUND`、`slots` 单节点、`slot_start_batch_size` 跟实例数）：

```yaml
modules:
  router:
    parameters:
      task_dist_mode: SLOT-BOUND
      slot_start_batch_size: 4
    slots:
      - n0:4          # seq 0..3
```

来源模块各加一行（**必须 ≤ 实例数 - 1**）：

```yaml
  module-heavy-a:
    parameters:
      router_index: 1
  module-heavy-b:
    parameters:
      router_index: 2
  module-tail:
    parameters:
      router_index: 3
```

**改完必须重建 app**——`parameters` 在建 app 时快照，已有 app 不受影响。

### 4. 验证

```sh
# 重跑后：各来源的排队应降到秒级
docker exec database psql -U scalebox -d scalebox -c "
SELECT t.from_mod, count(*),
       round(avg(extract(epoch from (t.begin_time - t.create_time)))::numeric,1) avg_queue_s
FROM t_task t JOIN t_module m ON m.id = t.mod_id
WHERE m.app = <new_app_id> AND m.name = '<main_router>' AND t.begin_time IS NOT NULL
GROUP BY 1;"
```

**判据**：最大 `avg_queue_s` 从数十秒降到个位数，且整轮 wall 不再由 router 的
`SUM(end_time - begin_time)` 决定。

**失败下一步**：若某个来源的排队仍然很高，说明它自己就是那个"单来源占 90% 条数"
的情形——回到 [`bypass-router-short-tasks`](./bypass-router-short-tasks.md)。

## 实证（app-fxcorr，2026-11-05）

8 节点 × 0.128 s × 3840 批的 B 方案（fpm → xp 链）。**单实例**实测：

| 来源 | 次数 | 平均执行 | 总执行 | 占 router |
|---|---|---|---|---|
| `vtask-head` | 3840 | 0.2722 s | 1045.3 s | 50% |
| `fxcorr-xp` | 15360 | 0.0552 s | 847.9 s | 40% |
| `vtask-tail` | 3840 | 0.0497 s | 190.7 s | 9% |
| `fxcorr-prep` | 1 | 16.85 s | 16.8 s | — |
| 合计 | 23043 | | **2101 s** | |

整轮 wall 3366 s，router 占 62%；该轮下游 xp 的并行执行下限是
272333 s / (24×8 槽) = **1418 s**——router 比下游还慢，是瓶颈。

特征：head 条数少但单条最重（读 `batch.json` + 解析 `.input` + 建计数器 + 展开任务），
xp 条数最多但单条最轻。**不能只看条数排序，要按 `total_exec_s` 排。**

按"每个来源一个实例"分配 4 实例后，最大分支 1045 s < 下游下限 1418 s，router 不再是瓶颈。

## 坑

（配置类坑见 feature 文档，这里只记拆实例特有的）

- **按条数而不是按处理时间分配**：上例中 head 条数只有 xp 的 1/4，处理时间却更多。
  用 `total_exec_s` 排序。
- **给零散来源开实例**：prep / 入口 / 实验级归并每轮只有几次，塞进缺省实例即可；
  单独开一个只会让 `slot_start_batch_size` 变大、启动更慢。
- **忘了删旧的 `router_index`**：历史上给某个模块钉过实例，模块换用途后那个值还在，
  会与新分配撞车。改配置时全局搜一遍 `router_index`。
- **以为拆实例能突破同一来源的上限**：一个来源模块只能钉一个实例，规模再涨也还是
  单进程。真正的上限突破要靠减少回流条数。
- **和"下游 slot 不够"混淆**：拆实例只提高回流处理并发，不改变下游消费能力。
  下游排队高时拆 router 无效。
