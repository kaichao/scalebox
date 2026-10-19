# check-test

准入控制与空间流控：agent 在每轮领取任务**之前**执行检查，不满足条件就停止领取（任务不失败、不消耗），条件满足后自动恢复。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `free_space_gb` | `module.arguments` | `{"<路径>": <GB>, ...}` | 空 | 路径所在分区**至少保留**的空闲空间；不足则停领 |
| `dir_quota_gb` | `module.arguments` | `{"<路径>": <GB>, ...}` | 空 | 目录**最大占用**空间；超限则停领 |
| `ACTION_CHECK` | `environments` | 脚本路径 | `/app/bin/check.sh`（存在时） | 自定义准入脚本；返回非 0 则停领 |
| `slot_options` | `module.arguments` | `slot_on_head` 等 | 空 | 单 slot 测试常用 `slot_on_head` |

`free_space_gb` / `dir_quota_gb` 的值支持 slot 表达式（如 `((@seq+5))`），用于给不同 slot 不同阈值；规则以 hex 编码的 JSON 传给 agent（`PLAT_FREE_SPACE_GB` / `PLAT_DIR_QUOTA_GB`）。

## 行为与不变式

- 检查顺序：`free_space_gb` → `dir_quota_gb` → `ACTION_CHECK`；任一失败，本轮不领任务。
- 空间检查失败返回码 4（`spaceutil.CheckFreeSpace` / `CheckDirQuota`），**不是任务失败**：任务保持 `status_code=-1`、无 `begin_time`、fail=0。
- 条件恢复后约 5 秒内自动放行（下一轮 poll）。
- `ACTION_CHECK` 失败同样只停领；脚本可用 `timeout` 包裹（受 `task_max_seconds` 限制）。
- 值为 0 表示不限制。

## 坑

- **参数名是 `free_space_gb`**，不是 `space_free_gb`；后者在源码中不存在，写了不会生效。
- **挂载位置决定会不会死锁**：任务级空间闸门 + 组内任务原子性会死锁（盘上数据不再增长、fail=0）。批级流控应挂在批入口模块（如 `vtask-head`，每节点 1 slot，`@seq` 恒 0）；`VTASK_SIZE=1` 可消除首轮盲区。
- **阶梯阈值是背压阀**：`((@seq+5))` 让每个 slot 阈值不同，用于启动后累积阶段；空盘时第一波仍会全部放行。
- `slot_on_head` 必须写在 `arguments.slot_options` 下；写在 `parameters` 下不会被识别。

## 最小例子

见 `feature.yaml`：

```yaml
modules:
  check_test:
    parameters:
      start_task: ANY
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
      dir_quota_gb: ${DIR_QUOTA_GB}
      free_space_gb: ${FREE_SPACE_GB}
    environments:
      - ACTION_CHECK=${ACTION_CHECK}
```

`scalebox.env` 中 `DIR_QUOTA_GB` / `FREE_SPACE_GB` / `ACTION_CHECK` 需先定义（空值也要写，避免 `${VAR}` 字面量透传）。

## 测试

```sh
./test.sh   # 阶段 3 提供
# 当前手工：
DIR_QUOTA_GB='{"/data/ssd/tmp":22}' FREE_SPACE_GB='{"/data/ssd/tmp":1000}' \
  ACTION_CHECK=/app/bin/mycheck.sh scalebox run --app-file=feature.yaml
```

预期：阈值不满足时 slot 不领任务；阈值满足后自动恢复；`t_slot_event` 中可见 `free_space_gb` / `dir_quota_gb` 的检查记录。

## 源码指针

- `go-scalebox/cmd/agent/check.go`（`check`）—— 检查顺序与失败返回码
- `go-scalebox/cmd/agent/spaceutil/`—— `CheckFreeSpace` / `CheckDirQuota`
- `go-scalebox/cmd/controld/grpc/app_env_vol.go`—— 参数到 `PLAT_FREE_SPACE_GB` / `PLAT_DIR_QUOTA_GB` 的 hex 编码
- `scalebox/features/check-test/code/mycheck.sh`—— 自定义 `ACTION_CHECK` 示例

## 相关文档

- Sphinx：`docs/cn/source/developer/7_advanced_features.md` §7.2（准入控制）
- 参数表：`docs/reference/app-yaml.md`（`module.arguments` 的 `free_space_gb` / `dir_quota_gb`）
- 实战：app-fxcorr `docs/dynamic-flow-control.md`（挂载位置与死锁）
