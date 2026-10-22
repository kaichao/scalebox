# side-car

模块生命周期脚本（边车模式）：把初始设置、准入检查、任务运行、退出清理拆成四个可定制脚本。

## 配置

| 环境变量 | 阶段 | 缺省脚本 | 说明 |
|---|---|---|---|
| `ACTION_SETUP` | slot 启动时执行一次 | `/app/bin/setup.sh` | 非 0 → slot 置 ERROR |
| `ACTION_CHECK` | 每轮领取任务之前 | `/app/bin/check.sh` | 非 0 → 本轮不领任务（不失败） |
| `ACTION_RUN` | 每个 task | `/app/bin/run.sh` | 任务主逻辑；stdin = body，参数 = body/headers |
| `ACTION_TEARDOWN` | slot 退出前 | `/app/bin/teardown.sh` | 非 0 → slot 置 ERROR |

- 环境变量为空时用缺省脚本；相对路径自动加 `/app/bin/` 前缀。
- `/app/bin/` 下不存在可执行文件时，回退查找 `/app/share/bin/`。
- `task_max_seconds > 0` 时，四个脚本都会用 `timeout Ns` 包裹。
- `slot_options: reserved_on_exit` 可在 slot 退出后保留容器，便于排错；`max_idle_polls` 控制空闲轮询次数。

## 行为与不变式

- 四个脚本的失败语义不同：SETUP/TEARDOWN 失败会让 slot 进入 ERROR；CHECK 失败只是停领；RUN 失败是任务失败（`status_code > 0`，不回流）。
- CHECK 在 `GetTaskList` 之前执行，失败时任务从未被拾取，因此不会产生失败任务、也不需要重试。
- RUN 的输入约定：stdin 是 task body；`$1` 是 body，`$2` 是 headers JSON（与平台模块脚本一致）。

## 坑

- **CHECK 不是 RUN**：把业务判断放在 CHECK 只会停领，不会让任务失败；需要失败语义就放在 RUN。
- **SETUP 只执行一次**：slot 生命周期内不会重复；需要每任务初始化请放在 RUN。
- **TEARDOWN 不保证执行**：slot 被强杀（节点掉线、超时）时可能不执行；不要依赖它做关键清理。
- **缺省脚本也会执行**：如果镜像里存在 `/app/bin/run.sh`，即使没配 `ACTION_RUN` 也会运行它；定制时注意覆盖。

## 最小例子

见 `feature.yaml`：

```yaml
modules:
  side_car_test:
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head,reserved_on_exit
      max_idle_polls: 50
    environments:
      - ACTION_RUN=${ACTION_RUN}
      - ACTION_SETUP=${ACTION_SETUP}
      - ACTION_TEARDOWN=${ACTION_TEARDOWN}
      - ACTION_CHECK=${ACTION_CHECK}
```

`scalebox.env` 中四个变量需显式定义（空值也要写，避免 `${VAR}` 字面量透传）。

## 测试

```sh
./test.sh   # 定制 ACTION_RUN，断言 run-custom.sh 被调用
# 当前手工：缺省脚本
echo 0 | scalebox run --app-file=feature.yaml

# 定制脚本
ACTION_RUN=/app/bin/run-custom.sh \
ACTION_SETUP=/app/bin/setup-custom.sh \
ACTION_TEARDOWN=/app/bin/teardown-custom.sh \
ACTION_CHECK=/app/bin/check-custom.sh \
scalebox run --app-file=feature.yaml
```

预期：四个脚本按阶段执行；`t_slot_event` 中可见各阶段 stdout/stderr/返回码；定制脚本的输出可在 `t_task_exec` 查看。

## 源码指针

- `go-scalebox/cmd/agent/actions.go`（`setActions` / `runSetup` / `runTeardown`）—— 脚本解析与执行
- `go-scalebox/cmd/agent/check.go`—— `ACTION_CHECK` 的调用时机
- `go-scalebox/cmd/agent/run.go`—— `ACTION_RUN` 与任务主循环

## 相关文档

- Sphinx：`docs/cn/source/application/4_module_development.md` §2.3.1
- 参数表：`docs/cn/source/reference/2_parameter_reference.md` §2.3（环境变量）
- 标准模块：`docs/cn/source/application/6_standard_modules.md`
