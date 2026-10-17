# features（特性测试）

每个目录一个平台特性，三合一：面向 AI agent 的说明（`README.md`）、最小可跑应用（`feature.yaml`）、可判定的集成测试（`test.sh`，阶段 3 提供）。

## 约定

- 目录名：小写连字符（kebab-case），如 `task-batch-size`；单词目录直接用小写（如 `pod`）。
- `README.md`：按 [`_template/README.md`](./_template/) 的模板编写，中文。
- `feature.yaml`：与 `app.yaml` 同 schema；用 `scalebox run --app-file=feature.yaml` **显式**加载（`scalebox run` 的缺省探测只找 `app.yaml`）。
- `test.sh`：一键建 app → 跑 → 断言 → 清理，返回 0/非 0（阶段 3 补齐）。
- 索引一致性由 `./check-index.sh` 校验；新增/删除 feature 后运行。

## 索引

| feature | 覆盖特性 | 状态 |
|---|---|---|
| [agent-interface-files](./agent-interface-files/) | 模块与 agent 的 8 类接口文件 | README 已套模板；测试待阶段 3 |
| [check-test](./check-test/) | `free_space_gb` / `dir_quota_gb` / `ACTION_CHECK` | 已修正参数名与 slot 选项写法 |
| [env-lifecycle](./env-lifecycle/) | app 快照 / 未定义变量 / non-tty stdin | 阶段 2 新建 |
| [headers-propagation](./headers-propagation/) | `_` 前缀传播 / `_vtask_attr_*` | 阶段 2 新建 |
| [node-progress-gap](./node-progress-gap/) | HOST-BOUND 节点进度差 | README 已套模板 |
| [pod](./pod/) | `pod_id` 节点本地计算 | README 已套模板 |
| [retry-test](./retry-test/) | `retry_rules` 退出码重试 | 新增 README |
| [side-car](./side-car/) | `ACTION_RUN/SETUP/TEARDOWN/CHECK` | README 已套模板 |
| [slot-expr](./slot-expr/) | slot 启动表达式 | README 已套模板 |
| [task-batch-size](./task-batch-size/) | `task_batch_size` 批量领取 | README 已套模板 |
| [task-identity](./task-identity/) | `uniq_task` / 112 / `conflict-action` / `repeatable` | 阶段 2 新建 |
| [timeout](./timeout/) | `task_max_seconds` 超时 | README 已套模板 |

## 规划中（未建目录）

来自 `docs/feature-inventory.md` §5：

| feature | 覆盖特性 |
|---|---|
| sink-tasks | 回流格式 / `last_from_module` / 与 vtask 的区别（部分已并入 agent-interface-files） |
| vtask-flow-control | wait-queue / bind / gate / fail |
| check-flow-control | 空间闸门挂载位置与死锁（可并入 check-test） |

## 校验

```sh
./check-index.sh
```
