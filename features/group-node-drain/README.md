# group-node-drain

独立组节点（GROUP_SLOTS）的排空与替换：host 进入 `DRAINING` 后，生产者模块停止领取并退出 slot，消费者模块继续运行，把该节点本地链上的数据消费完后，再安全地换 IP / 换节点。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `slots` 的组 slot 条目 | `module.slots` | `<host>:<count>:<group_prefix>`，如 `n0-0:2:n0` | — | 在 `host` 上建 `count` 个组 slot，slot 领取 `to_host` 匹配 `group_prefix.+` 的 task；第三段为单个字母（如 `x`）表示全局前缀 |
| `skip_draining_host` | `module.parameters` | `'yes'` / 缺省 | 缺省（继续领取） | `yes` 时，该模块在 `DRAINING` host 上停止领取 task，并让 slot agent 退出 |
| host 状态 | `t_host.status` | `ON` / `DRAINING` / `OFF` | `OFF`（新建） | 通过 `scalebox host drain <hostname>` / `undrain` / `set-status` 修改 |

独立组节点的典型应用配置（应用 env 变量只是字符串拼接，平台只认上面的 slot 条目语法）：

```bash
# y00-00 是独立组锚点节点，4 个组 slot，group_prefix=x43
GROUP_SLOTS='y00-00:4:x43'
```

## 行为与不变式

- `DRAINING` 的 host **仍参与 `to_host` 位置解析**（`ON` 与 `DRAINING` 同等待遇），否则本地链下游 task 会退化成空 `to_host`。
- `skip_draining_host=yes` 的模块在 `DRAINING` host 上：
  - 领取判断拒绝该 slot（`isSlotDispatchable`）；
  - `GetSlotList` 不再启动该模块的 READY slot；
  - 已运行的 agent 在下一次 poll 收到 `slot_action=stop` 后优雅退出，slot 回到 `READY`。
- 未配置 `skip_draining_host` 的模块继续领取和运行，用于消费该节点本地已有的数据（如 1ch/1chx）。
- 排空完成的口径：
  1. 该 host 上组 slot 无执行中的 task（`t_task_exec.ts_finished IS NULL` 为 0）；
  2. `to_host` 指向该 host 的 pending/running task 为 0；
  3. 本地数据目录不再增长（业务侧确认）。
- 排空后置 `OFF`、换 IP / 换主机名，再置 `ON`（或 `undrain`）；READY slot 会重新启动。
- 恢复 `ON` 后，之前被跳过的 slot 重新可领取。

## 坑

- **`DRAINING` ≠ `OFF`**：直接置 `OFF` 会切断位置解析和领取资格，正在排空的本地链会断掉（表现为空 `to_host` 或下游 task 无人领取）。先 `drain`，排空后再 `OFF`。
- **只给生产者配 `skip_draining_host`**：把消费者也配成 `yes`，本地链没有模块消费，排空永远不结束。生产者是产生本地数据的模块（如 beam-make），消费者是下游（如 down-sample、fits-redist）。
- **agent 退出有延迟**：退出由 agent 的 poll 周期驱动（默认 10 秒），不是设置 `DRAINING` 立即退出；判断排空时要留出这段窗口。
- **复用 host 名换 IP**：`to_host` 存的是逻辑主机名，本地数据绑在物理机上。排空完成后再改 IP；不要在还有 pending/running task 时换 IP，否则任务会被调度到没有数据的新机器。
- **组 slot 的 `group_prefix` 是正则前缀**：`n0-0:2:n0` 匹配 `to_host ~ 'n0.+'`；写成单字母 `x` 表示匹配所有 `x` 开头的 host。写错会导致组 slot 领不到 task。
- **历史空串 `to_host` / `to_slot`**：平台已用 `NULLIF` 把它们当作"未指定"处理，不会再永久卡在 `-1`；新 task 也不会再写入空串。

## 最小例子

见 `feature.yaml`。关键片段：

```yaml
modules:
  producer:
    parameters:
      task_dist_mode: HOST-BOUND
      skip_draining_host: 'yes'
    slots:
      - ${GROUP_SLOTS}        # 例：n0-0:2:n0

  consumer:
    parameters:
      task_dist_mode: HOST-BOUND
    slots:
      - ${GROUP_SLOTS}
```

排空操作：

```bash
scalebox host drain n0-0.inline     # 停生产者，消费者继续
# 观察排空：见上文的两个查询
scalebox host set-status --id <id> --status OFF   # 或直接改 IP 后 undrain
scalebox host undrain n0-0.inline   # DRAINING -> ON
```

## 测试

```sh
./test.sh
```

默认需要 inline 多节点控制面（`n0-0` 等）；不可达时 `tk_skip`。

单节点环境（如测试机 bio 的 local 集群，只有 `h0`）可用变体跑同一套断言：

```sh
SCALEBOX_ENV_FILE=scalebox-local.env SCALEBOX_DRAIN_HOST=h0 ./test.sh
```

断言点：

1. `producer` 在 `n0-0` 的 slot 在 `drain` 后 60 秒内变为 `READY`（agent 退出）；
2. 新投递给 `producer` 的 task 保持 `-1`（skip 模块不再领取）；
3. 新投递给 `consumer` 的 task 正常完成（非 skip 模块继续运行）；
4. 测试结束恢复 host 为 `ON`。

## 源码指针

- `go-scalebox/cmd/controld/task/get.go`（`isSlotDispatchable` / `getGroupBoundRows`）—— DRAINING 下按模块停领、组 slot 领取条件
- `go-scalebox/cmd/controld/task/cache.go`（`cacheIPAddr2HostName`）—— DRAINING 仍参与 `to_host` 解析
- `go-scalebox/cmd/controld/slot/slot.go`（`GetSlotList`）—— DRAINING host 上按模块决定是否启动 slot
- `go-scalebox/cmd/controld/grpc/slot.go`（`SubmitSlotEvent`）—— 返回 `slot_action=stop`
- `go-scalebox/cmd/agent/run.go`（`runAgent`）—— 收到 `slot_action=stop` 后退出
- `go-scalebox/build/database/initdb.d/83-host-draining.sql`—— host 状态约束加入 `DRAINING`

## 相关文档

- Sphinx：`docs/cn/source/design/3_vtask.md`（GROUP-BOUND 与独立组节点）
- 不变式：`docs/agents/10-invariants.md`、`docs/agents/glossary.md`
- 示例：`examples/vtask/README.md`（GROUP-BOUND 基础管道）
