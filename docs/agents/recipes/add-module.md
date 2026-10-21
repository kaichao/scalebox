# recipe — add-module（给已有 app 加模块）

> 用于：在跑通的 app 上增加一个处理环节，并让上游把任务正确交给它。

## 前置

- 已有 app 能跑通（先按 [`new-app`](./new-app.md) 验证过）。
- 明确**接线方式**：上游是路由模块（改 `sink_modules`）还是算法模块（写 `sink-tasks.txt`）。

## 步骤

1. 画接线，二选一：

| 上游类型 | 下游如何被投递 |
|---|---|
| 路由模块（`main_router`） | 在 router 的 `sink_modules` 里加新模块名；router 脚本按 `from_module`/header 选择 |
| 普通算法模块 | 上游成功时写 `${WORK_DIR}/sink-tasks.txt`，第一列为新模块名 |

2. 写模块代码：

```bash
# code/run.sh
#!/bin/bash
set -euo pipefail
echo "process: $1"
# 有下游时：把下游任务写进 sink-tasks.txt（格式：模块名,body,headers-json）
# echo "next,${1},{}" >> "${WORK_DIR}/sink-tasks.txt"
```

注意：成功才写；文件为空会被判后处理错误；嵌套子脚本先 `rm -f`（INV-04）。

3. `app.yaml` 增加模块条目：

```yaml
modules:
  stage2:
    base_image: scalebox.net/platform/agent:latest
    arguments:
      code_path: ${PWD}/stage2-code
      slot_options: slot_on_head     # 或 slots: [h0]
    parameters:
      task_dist_mode: DEFAULT        # HOST-BOUND 时任务需 to_host
    environments:
      - MY_PARAM=value
```

并在上游条目里加 `sink_modules: [stage2]`（仅路由模块需要）。

4. 重建 app（配置在创建时快照，INV-01）：

```bash
scalebox app delete --app-id <old_app_id>
scalebox run --app-file=app.yaml < /dev/null
```

5. 投一个能流过上游的测试任务，确认下游任务出现：

```bash
echo "test-body" | scalebox task add --app-id <new_app_id> --sink-module <upstream>
scalebox task list --app-id <new_app_id>
```

## 判据

- 上游任务 `status_code == 0`；下游出现新任务，body 与预期一致。
- 下游模块的 slot 已建（`scalebox slot status --module-id <stage2_mod_id>`）。
- `_` 前缀 header 按预期到达下游（`task show` 的 `headers`）。

## 失败下一步

| 现象 | 下一步 |
|---|---|
| 上游成功但下游无任务 | 检查 `sink-tasks.txt` 是否写、第一列模块名是否与 `modules:` 键完全一致、空文件 |
| 下游任务 112 | 身份冲突 → `features/task-identity` |
| 下游 READY 不动 | `task_dist_mode: HOST-BOUND` 但没给 `to_host`；或 slot 未建 |
| `_` header 丢失 | 非路由模块只自动继承 `_`；普通 header 要显式写进 sink-tasks、INV-13 |
