# 3. 第一个应用

本章从零创建一个最小应用：一个模块、一个任务，看到 stdout 后清理。全程不需要构建镜像，模块代码通过 `code_path` 挂载。

## 3.1 准备

- 控制面已启动，见 [快速上手](2_quick_start.md)。
- `scalebox` CLI 可用。
- 本章命令在应用目录执行；示例目录为 `hello-app/`。

```text
hello-app/
├── app.yaml
└── code/
    └── run.sh
```

## 3.2 写 app.yaml

```yaml
name: hello-app
cluster: local

modules:
  hello:
    base_image: scalebox.net/platform/agent:latest
    arguments:
      code_path: ${PWD}/code
      slot_options: slot_on_head
```

- `base_image` 在测试部署使用 `scalebox.net/platform/agent:latest`；正式部署按实际镜像修改。
- `code_path` 把本地 `code/` 挂到模块容器；修改脚本后无需重建镜像，但要重建 app。
- `slot_options: slot_on_head` 在 local 集群头节点创建一个 slot。若改用 `slots`，写主机短名或正则（如 `h0`、`n0-[01]`），不要写完整 FQDN。

## 3.3 写 code/run.sh

```bash
#!/bin/bash
set -euo pipefail

printf 'hello from task: %s\n' "$1"
```

```bash
chmod +x code/run.sh
```

`$1` 是 task body；stdout 会写入 `t_task_exec.stdout`，`status_code` 为 0 表示成功。
本章没有下游模块，不要创建 `sink-tasks.txt`。

## 3.4 创建应用

```bash
scalebox run --app-file app.yaml
```

输出形如：

```text
{"app_id":1}
```

记下 `app_id`。应用创建时会快照镜像、环境变量与参数；修改 `app.yaml` 后必须删除并重建。

## 3.5 投递任务并看结果

```bash
echo "hello scalebox" | scalebox task add --app-id 1 --sink-module hello
```

```bash
scalebox task list --app-id 1
scalebox task show <task-id>
```

预期 `task show` 的 stdout：

```text
hello from task: hello scalebox
```

且 `status_code` 为 0。

## 3.6 修改后重建

修改 `run.sh` 后，重新投递任务即可（`code_path` 是挂载的）。
修改 `app.yaml`（镜像、参数、环境变量）后，`scalebox` 不会热更新，需要：

```bash
scalebox app delete --app-id 1
scalebox run --app-file app.yaml
```

## 3.7 清理

```bash
scalebox app delete --app-id 1
```

`app delete` 只删除数据库记录，不停止本地 slot 容器；容器清理见
[卡链与失败诊断](https://github.com/kaichao/scalebox/blob/dev/docs/agents/20-diagnosis.md)。

## 3.8 下一步

本章是最小单模块应用，没有主路由；多模块协作时，由主路由执行并行逻辑、读写应用状态。下一步先读核心概念，再按需读主路由实现。

- [核心概念](../user/2_core_concepts.md)：基本概念、编程原语、计算逻辑 / 应用状态 / 主路由 / VTask。
- [运行与管理应用](../user/3_running_apps.md)：投递、查看、清理的完整流程。
- [主路由与应用状态：实现与优化](../developer/4_main_router_status.md)：并行逻辑如何落地。
- [app.yaml 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)
- [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)
- [features/](https://github.com/kaichao/scalebox/tree/dev/features)：可执行的最小特性与测试。
