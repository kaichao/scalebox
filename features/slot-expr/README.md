# slot-expr

Slot 启动命令表达式：在模块的 `command` 中写 `((@...))`，平台在 slot 启动时求值，把结果作为容器参数/环境变量。

## 配置

| 表达式 | 变量来源 | 类型 | 说明 |
|---|---|---|---|
| `((@seq))` | slot 序号 | 数值 | 每个节点内从 0 计数 |
| `((@p:<name>))` | slot 参数 | 字符串 | `t_slot.parameters`（jsonb），每 slot 独立 |
| `((@v:<name>))` | 全局变量 | 字符串 | `t_global`，跨 slot 共享 |
| `((@seq + @v:<name>))` 等 | 组合 | 数值/字符串 | 支持算术；非数字参与算术会报错并保留原表达式 |

slot 参数用 `scalebox slot set-parameter` 设置；全局变量用 `scalebox global set` 设置。

## 行为与不变式

- 表达式在 **slot 启动时**求值一次，之后不随变量变化。
- `@seq` 是节点内序号，不是全局序号；多节点时各节点都从 0 开始。
- `@v:` 查询 `t_global`；`@p:` 查询当前 slot 的 `parameters`。
- 算术表达式按数值求值；变量值不是数字时白名单校验失败，报错并原样保留表达式。
- slot 参数未设置时按空值替换（不会报错）。

## 坑

- **`@seq` 是节点内序号**：需要全局唯一标识时不要直接用它。
- **表达式只求值一次**：运行中改变全局变量不会影响已启动的 slot。
- **命令里要加引号**：表达式可能展开出空格或特殊字符时才需要；数值表达式按求值结果替换。
- **调试用 `--print-yaml`**：`scalebox run --print-yaml` 可先看模板展开结果（注意表达式是 slot 启动时才求值，print 阶段可能仍是原样）。

## 最小例子

见 `feature.yaml`：

```yaml
modules:
  mod-slot-expr:
    command: 'docker run -d --network=host -e SEQ=((@seq)) -e MY_PARAM=((@p:my_param)) -e MY_VAR=((@v:my_var)) --rm {{ENVS}} {{VOLUMES}} {{IMAGE}}'
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
```

`code/run.sh` 打印三个变量的值，用于确认求值结果。

## 测试

```sh
./test.sh   # 阶段 3 提供
# 当前手工：
app_id=$(scalebox run --app-file=feature.yaml | cut -d':' -f2 | tr -d '}')
scalebox global set my_var VAR_VALUE
scalebox slot set-parameter my_param PARAM_VALUE --slot-id <slot_id>
echo "my-task-body" | scalebox task add --app-id=$app_id
```

预期：任务输出中看到

```text
SEQ=0
MY_PARAM=PARAM_VALUE
MY_VAR=VAR_VALUE
```

## 源码指针

- `go-scalebox/cmd/controld/slot/expr/expr.go`（`Evaluate`）—— 表达式解析
- `go-scalebox/cmd/controld/slot/expr/README.md`—— 求值规则与示例
- `go-scalebox/cmd/controld/grpc/app_slots.go`—— slot 启动命令中的表达式替换

## 相关文档

- Sphinx：`docs/cn/source/developer/7_advanced_features.md` §7.7
- 参数表：`docs/cn/source/appendix/2_parameter_reference.md` §2.2.3（slot-parameters）
