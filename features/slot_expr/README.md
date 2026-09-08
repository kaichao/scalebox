# slot_expr

验证 slot 启动命令表达式的三类变量在 slot 启动时动态求值：

| 表达式 | 变量来源 | 类型 | 说明 |
|--------|----------|------|------|
| `((@seq))` | slot 序号 | 数值 | 每个节点内从 0 计数 |
| `((@p:my_param))` | slot 参数 | 字符串 | `t_slot.parameters`（jsonb），每 slot 独立 |
| `((@v:my_var))` | 全局变量 | 字符串 | `t_global`，跨 slot 共享 |

app.yaml 的模块命令将三个表达式作为容器环境变量传入，`code/run.sh` 打印三个变量的值，用于确认求值结果。表达式详细说明见文档《高级特性》§7.7（slot 启动命令表达式）。

## 操作实例

```sh
export APP_ID=$(scalebox run | cut -d':' -f2 | tr -d '}' )
export MODULE_ID=$(scalebox module list --no-head | grep mod-slot-expr | awk '{print $1}')
export SLOT_ID=$(scalebox slot list --no-head | awk '{print $2}')

# 设置全局变量与 slot 参数
scalebox global set my_var VAR_VALUE
scalebox slot set-parameter my_param PARAM_VALUE --slot-id ${SLOT_ID}

# 添加任务，触发 slot 启动
scalebox task add "my-task-body"
```

任务输出中预期看到：

```text
SEQ=0
MY_PARAM=PARAM_VALUE
MY_VAR=VAR_VALUE
```

注意：slot 参数未设置时，`((@p:my_param))` 按空值替换。
