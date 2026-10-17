# env-lifecycle

应用环境变量的加载、展开与生命周期：从 CLI 到 app 快照再到容器环境，哪一步发生什么。

## 配置

| 来源 | 优先级 | 说明 |
|---|---|---|
| 进程环境（shell 已导出） | 最高 | 覆盖 env 文件中的同名键（含空值） |
| `--env-file` 指定文件 | 次之 | 按出现顺序，后者覆盖前者 |
| 基线 `scalebox.env` | 再次 | CLI 先加载基线，再加载 `--env-file` |
| `~/.scalebox` / `/etc/scalebox` | 最低 | 全局配置 |

模块收到的变量由 `module.environments` 注入；`arguments` 中的部分参数会转换为 `PLAT_*` 环境变量。

## 行为与不变式

- 未定义的 `${VAR}` 在模板展开时**保留字面量**，不会变成空串；每个变量必须显式定义（空值也写）。
- app 创建时对 env、镜像名、slots、`code_path` 做快照；之后修改这些配置不会影响已创建的 app。
- `scalebox run` 在 stdin 非 tty 时把 stdin 内容当初始任务列表；脚本包装必须 `< /dev/null`。
- 进程环境中值为空的 key 不会被 env 文件覆盖（`getModuleEnvMap` 的删除规则）。

## 坑

- **漏定义变量**：app.yaml 里出现字面量 `${VAR}`，模块收到错误路径/阈值。
- **改了 env 不重建 app**：旧快照继续生效；必须重建。
- **脚本里忘记 `< /dev/null`**：脚本剩余文本被创建成任务。
- **同名变量多来源**：shell 导出 > `--env-file` > 基线；排错时先确认最终值来源。

## 最小例子

`feature.yaml` 注入 `DEMO_VALUE` 与一个字面量占位 `LITERAL_DEMO`；`scalebox.env` 只定义 `DEMO_VALUE`。

```yaml
modules:
  env-lifecycle:
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
    parameters:
      start_task: 0
    environments:
      - DEMO_VALUE=${DEMO_VALUE}
      - LITERAL_DEMO=${LITERAL_DEMO}
```

`code/run.sh` 打印 `DEMO_VALUE`、`LITERAL_DEMO`、`APP_ID` 等。

## 测试

```sh
./test.sh   # 阶段 3 提供
# 当前手工：
DEMO_VALUE=hello scalebox run --env-file=scalebox.env --app-file=feature.yaml
# 查看模块 stdout：DEMO_VALUE=hello；LITERAL_DEMO 应为字面量 ${LITERAL_DEMO}
```

预期：未定义变量以字面量传入；改 `scalebox.env` 后不重建 app，旧值仍生效。

## 源码指针

- `go-scalebox/internal/common/yaml_template.go`（`RenderYAMLTemplate`）—— 变量展开与未定义保留
- `go-scalebox/cmd/cli/root.go`（`--env-file`）—— env 文件加载
- `go-scalebox/cmd/cli/run/stdin_start_tasks.go`—— 非 tty stdin 任务列表
- `go-scalebox/cmd/controld/grpc/app_create.go`—— app 快照
- `go-scalebox/docs/env-loading.md`—— 优先级与展开规则

## 相关文档

- 参数表：`docs/cn/source/appendix/2_parameter_reference.md` §2.3（环境变量）
- 设计文档：`go-scalebox/docs/env-loading.md`
- 实战：`docs/agents/10-invariants.md` INV-01 / INV-02 / INV-03
