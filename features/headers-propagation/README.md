# headers-propagation

`_` 前缀 header 的自动传播：平台把任务上下文中以下划线开头的 header 沿 sink 链传给下游，应用不需要手动搬运。

## 配置

无平台参数。模块从 `$2`（headers JSON）读取；bash 用 `scalebox::task_header "$2" <name>`。

| header | 来源 | 说明 |
|---|---|---|
| `_vtask_id` | 平台注入 | vtask 树标识；应用不要手写 |
| `_vtask_size_sema` | 平台注入 | 流控信号量名；应用不要手写 |
| `_vtask_attr_<name>` | 应用 | 自定义属性；平台只运输、不解释 |

## 行为与不变式

- **非路由模块**写 `sink-tasks.txt` 时，当前任务的 `_` 前缀 header 自动继承到下游任务。
- **经 router 分支**时，`_` 前缀的自动继承只对 `vtask_role≠""` 的 sink 模块放行。
- 经 **wait-queue** 一跳时，属性必须由生产方显式转发（wait-queue 没有 `vtask_role`，自动继承被门控拦下）。
- `_vtask_id` / `_vtask_size_sema` 由平台维护；应用只应使用 `_vtask_attr_*` 子命名空间。
- 普通 header（无 `_` 前缀）不会自动传播，需要显式写入 sink 行。

## 坑

- **router 转发后属性丢失**：因为门控只放行 `vtask_role≠""` 的 sink；需要在 router 的 `fromXxx` 中显式转发。
- **wait-queue 跳丢失**：入口任务经 router 进入 wait-queue 时，必须显式带上 `_vtask_attr_*`。
- **手写 `_vtask_id`**：会让 vtask 树指向错误的根；只读不写。
- **`_` 头与同名环境变量**：模块取值时 header 优先（见 `timeout` / side-car 相关条目）。
- **`scalebox::*` 需要先 source**：脚本开头写 `source /usr/local/lib/scalebox/functions.sh`，否则报 `command not found`、取值为空。

## 最小例子

`feature.yaml`：`producer` 把任务转发给 `consumer`，入口任务带 `_demo=hello`。

```yaml
modules:
  producer:
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
  consumer:
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
```

`code/run.sh` 按 `PLAT_MODULE_NAME` 分支：`producer` 写 `consumer,$1` 到 sink 文件；`consumer` 打印 `_demo` 的值。

## 测试

```sh
./test.sh
# 当前手工：
APP_ID=$(scalebox run --app-file=feature.yaml | ...)
echo "hello" | scalebox task add --app-id=$APP_ID --sink-module=producer --header _demo=hello
# 查看 consumer 的 t_task_exec.stdout：应包含 _demo=hello
```

预期：`consumer` 读到 `_demo=hello`；把 `_demo` 改成普通 `demo` 后，consumer 读不到。

## 源码指针

- `go-scalebox/cmd/agent/sink_tasks.go`（`addSinkTasks`）—— `_` 前缀继承与 router 门控
- `go-scalebox/cmd/agent/`—— 任务 header 解析
- `scalebox/examples/vtask/README.md`—— wait-queue 一跳的显式转发示例

## 相关文档

- 参数表：`docs/cn/source/reference/2_parameter_reference.md` §2.1.4（task-headers）
- Sphinx：`docs/cn/source/design/3_vtask.md`（header 传播）
- 设计文档：`go-scalebox/docs/task-headers.md`、`go-scalebox/docs/vtask-design.md`
