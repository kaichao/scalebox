# task-batch-size

批量领取任务：一个 slot 一次从服务端取多条任务，形成一个 `task_exec` 记录对应多个 task，降低短任务的领取开销。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `task_batch_size` | `module.parameters` | 整数 | 1 | 单次领取的最大 task 数 |

## 行为与不变式

- 配置后，controld 一次返回多条待执行任务；agent 在一次执行周期内处理这一批。
- 单个 `t_task_exec` 记录对应多个 task（批量任务聚合），不是每个 task 一条 exec 记录。
- 适合作业时长在 5 秒以内的短任务；长任务批量领取会拉长单 slot 的占用。
- 与 `node_progress_gap` 同时配置时，实际领取数取两者的较小值。

## 坑

- **不是越多越好**：批量过大时单个 slot 长时间占用，其他 slot 空转；短任务从 2–4 起步测。
- **task_exec 语义变化**：排查时按 task 计数与按 exec 计数会不同；统计任务量用 `t_task`，统计执行批次用 `t_task_exec`。
- **空队列时不会等待凑满**：有多少取多少，不保证每批都是 N 条。
- 与 `task_queue_length` / 去重窗口无直接关系，不要混用。

## 最小例子

见 `feature.yaml`：

```yaml
modules:
  task-batch-size:
    arguments:
      slot_options: slot_on_head
    parameters:
      task_batch_size: '${TASK_BATCH_SIZE}'
```

`scalebox.env` 中定义 `TASK_BATCH_SIZE=3`。

## 测试

```sh
./test.sh   # 阶段 3 提供
# 当前手工：
export TASK_BATCH_SIZE=3
for i in {0..9}; do printf "%03d\n" $i; done | scalebox run --app-file=feature.yaml
```

预期：10 个任务形成约 4 个 `t_task_exec` 记录（3+3+3+1）；任务本身全部完成。

## 源码指针

- `go-scalebox/cmd/controld/task/cache_attrset.go`—— `task_batch_size` 参数读取
- `go-scalebox/cmd/controld/task/get.go`—— 批量取任务的限制
- `go-scalebox/cmd/agent/run.go`—— 单次领取与执行循环

## 相关文档

- 参数表：`docs/cn/source/appendix/2_parameter_reference.md` §2.1.3
- Sphinx：`docs/cn/source/developer/7_advanced_features.md` §7.3（任务排序/批量相关）
- 实战：app-fxcorr v9 流水线（`TASK_BATCH_SIZE=8` 定案，见其 `CLAUDE.md`）
