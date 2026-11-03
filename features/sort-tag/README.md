# sort-tag

任务排序标签：主路由在投递时给任务写 `sort_tag`，controld 领任务时按它做第一优先级排序，
让"早批/早分片"的任务先被领走——多批同时在制时保证逐批推进，而不是各批交错。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `sort_tag` | task header（投递时设置） | 字符串 | 空 | controld 领任务 SQL 的第一排序键；通常由主路由在投递下游任务时写入 |

## 行为与不变式

- 领任务的 SQL（`cmd/controld/task/get.go`）：

  ```sql
  SELECT id FROM t_task
  WHERE status_code = -1 AND mod_id = $1 <附加条件>
  ORDER BY headers->>'sort_tag', gid_0, gid_1, gid_2, id
  LIMIT $2
  FOR UPDATE SKIP LOCKED
  ```

- **排序 = 字符串比较**。`"10"` 排在 `"2"` 前面——用数字语义的标签必须零填充
  （`00000010` / `00000002`）。
- **只影响"领取顺序"，不影响"完成顺序"**。多个 slot 并行执行，先领到的任务不保证先完成；
  需要完成顺序保证时靠下游的计数器（vtask / semaphore），不是 `sort_tag`。
- **只在 READY（`status_code = -1`）阶段生效**。已被领走（`-2`/执行中）的任务不受后续
  排序影响。
- 全部任务都不设 `sort_tag` 时全部为 NULL，退化为按 `gid_0..gid_2, id` 排——即"投递顺序"。
- 多 slot 并发领取有 `FOR UPDATE SKIP LOCKED`，不会重复领同一条。

## 坑

- **字符串序不是数值序**：批号 `1..120` 直接当标签，第 10 批会排到第 2 批前面。批号本身是
  零填充格式（如 `00000001`）时可直接用。
- **未设置的任务排在最后**：PostgreSQL 的 `ORDER BY` 对 ASC 缺省是 `NULLS LAST`，键不存在
  时 `->>` 返回 NULL。同一模块里混着"设了 sort_tag"和"没设"的两类任务时，没设的会一直排在
  后面——如果带标签的任务持续投递，未设的可能长时间轮不到（饿死）。
- **空串排最前**：`sort_tag: ""` 与"键不存在"不同，空串比任何非空字符串都小，会抢到最前面。
  不要用空串表示"不排序"，直接不写这个 header。
- **不自动传播到下游**：`sort_tag` 不以 `_` 开头，agent 写 sink 时只传播 `_` 前缀的头
  （见 `features/headers-propagation`）。每一跳的下游任务都要由投递方重新设置。
- **不是节流闸**：`sort_tag` 只改变顺序，不限制在制数量；控制并发用它旁边的 vtask /
  信号量机制。

## 实证（app-fxcorr，2026-11-03）

v9 链（fp → xp）在 2 节点、0.512 s × 120 批下跑，主路由在投 core 任务与 xp 分片时统一写
`sort_tag = <batch>`（批号本身是 8 位零填充）：

| 模块 | 任务数 | 带 sort_tag | 不同批号 | 范围 |
|---|---|---|---|---|
| `fxcorr-fp` | 3840 | 3840（100%） | 120 | 00000001–00000120 |
| `fxcorr-xp` | 480 | 480（100%） | 120 | 00000001–00000120 |

配合 `VTASK_SIZE` 从 8 提到 12：xp 的 slot 利用率 91.2% → 93.9%，全链 `wall_core`
754.9 s → 720.4 s（−4.6%）。

## 最小例子

见 `feature.yaml`。关键片段——主路由投递时写入：

```bash
# 同一个批的任务共用一个标签，不同批之间按标签先后被领走
scalebox task add --sink-module worker --headers='{"sort_tag":"00000001"}' "<batch>-<n>"
```

## 相关

- 任务头总表：`docs/reference/task-headers.md`
- 分组键 `group_regex` / `group_index`：排序的第二、三优先级
- 设计说明：`docs/cn/source/design/7_advanced_features.md` §7.3
