# wait-queue

vtask 管道的标准 wait-queue 模块：slot 级准入闸门（check.sh）+ 资源分配（run.sh 中 `vtask bind`）。

应用作者无需编写任何脚本或接触信号量名称，只需在 app.yaml 中引用镜像并配置参数。

## 参数

| 环境变量 | 说明 | 缺省值 |
|----------|------|--------|
| `VTASK_MODE` | vtask 类型：`HOST-BOUND` / `GROUP-BOUND`（head slot 锚定节点组） | `HOST-BOUND`（DEFAULT 模式不经 wait-queue） |
| `HEAD_MODULE` | vtask-head 模块名（信号量名与 add-subtask 目标模块） | `vtask-head` |

## 用法（app.yaml 片段）

```yaml
modules:
  wait-queue:
    base_image: scalebox.net/platform/wait-queue
    arguments:
      slot_options: slot_on_head
    parameters:
      vtask_size: 1          # gate 串行化：同一时刻只有一个 vtask 在执行 bind
    environments:
      - VTASK_MODE=${VTASK_MODE}
      # - HEAD_MODULE=my-vtask-head
```

## 行为

1. **准入（check.sh，ACTION_CHECK）**：agent 每轮 `GetTaskList` 之前执行；对应 vtask 类型的信号量组 max ≤ 0 时返回非 0，agent 停领并轮询重试
2. **分配（run.sh）**：`vtask bind` 原子扣减额度并返回资源（hostname / 组号）→ 构造 `_vtask_size_sema` → `vtask add-subtask --direct` 直投 vtask-head（同步感知错误，失败时 `vtask unbind` 回滚 + gate 放行）
3. **终结**：run.sh 不写 sink-tasks.txt，task 在 wait-queue slot 直接终结（无需回流 main-router）

额度由 controld `doVTaskFinished`（tail 完成）/ `FailVtask`（失败）自动加回，脚本侧无需参与。
