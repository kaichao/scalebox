# node-progress-gap

HOST-BOUND 模式下，同一模块在不同节点之间的进度差控制：快节点不能领先慢节点太多，防止节点间负载严重不均。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `node_progress_gap` | `module.parameters` | `{"<主机名前缀正则>": <gap>, ...}` | 空 | 同一组内快节点比慢节点最多多跑的 task 数 |
| `task_dist_mode` | `module.parameters` | `HOST-BOUND` | `DEFAULT` | 该特性只对 HOST-BOUND 生效 |
| `task_batch_size` | `module.parameters` | 整数 | 1 | 每次领取上限；与进度差取小 |

主机名匹配规则：用配置的键作为前缀正则匹配短主机名（去掉 cluster 后缀），后缀保留，信号量名为
`(node_progress:<module>:<前缀>)<后缀>`。

## 行为与不变式

- 每次领取前，agent/controld 计算 `gap - DiffMin(node_progress 信号量)`，实际领取数取 `task_batch_size` 与该值的最小值。
- 任务完成后 `node_progress:<module>:<host>` 加 1；任务失败时恢复（不加）。
- 信号量不存在时按 0 处理（不阻塞）。
- 只对 `HOST-BOUND` 且配置了 `node_progress_gap` 的模块生效；其他情况直接返回 `task_batch_size`。
- 进度按 `to_host` 归属，而不是按执行 task 的 slot 归属（组外计算也算到组内节点上）。

## 坑

- **正则前缀不是完整主机名**：`{"n-[01]":2}` 匹配 `n-0`/`n-1` 开头的短主机名；写错会让规则不生效。
- **信号量初值 0**：平台在 slot 生成时自动创建；手工建 slot 或改配置后需确认信号量存在。
- 进度差是在**领取时**限制，不改变已排队任务；节点掉线时仍可能观察到较大的已领取差。

## 最小例子

见 `feature.yaml`：

```yaml
modules:
  node-progress-gap:
    parameters:
      task_dist_mode: HOST-BOUND
      node_progress_gap: '{"n-[01]":2,"n":3}'
    slots:
      - n-[01][01]
```

## 测试

```sh
./test.sh   # 需要 inline 多节点控制面；不可达时 SKIP（退出码 77）
# 基线验证：任务全绿 + 每个节点一个 node_progress 信号量。
# 领取时的 gap 限制是动态行为（依赖持续任务供给），暂不自动断言。
# 当前手工（tasks.txt 提供初始任务）：
cat tasks.txt | scalebox run --app-file=feature.yaml
scalebox task add --app-id=$APP_ID --header to_host=n0-0.inline 00e
scalebox task add --app-id=$APP_ID --header to_host=n1-0.inline 10c
```

预期：快节点 `n-0x` 最多领先慢节点 2 个 task；`semagroup list` 中可见 `node_progress:<module>:<host>` 的值差不超过 gap。

## 源码指针

- `go-scalebox/cmd/controld/task/node_progress.go`（`getNodeProgressGapSize` / `updateNodeProgress` / `restoreNodeProgressSemaphore`）—— 领取限制与计数
- `go-scalebox/cmd/controld/task/cache_attrset.go`—— `node_progress_gap` 参数读取
- `go-scalebox/cmd/controld/grpc/app_env_vol.go`—— 信号量创建

## 相关文档

- 参数表：`docs/reference/app-yaml.md`（`module.parameters` 的 `node_progress_gap`）
- Sphinx：`docs/cn/source/developer/7_advanced_features.md` §7.2.4
