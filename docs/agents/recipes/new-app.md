# recipe — new-app（从零建最小 app）

> 用于：验证控制面可用、给新想法搭骨架、或给 feature 写最小例子。单模块、无下游、不构建镜像。

## 前置

- 控制面已启动，`scalebox` CLI 可用（快速上手见 Sphinx `started/2_quick_start`）。
- 有可用的 agent 镜像：测试部署用 `scalebox.net/platform/agent:latest`。

## 步骤

1. 建目录与文件：

```text
hello-app/
├── app.yaml
└── code/run.sh
```

```yaml
# app.yaml
name: hello-app
cluster: local

modules:
  hello:
    base_image: scalebox.net/platform/agent:latest
    arguments:
      code_path: ${PWD}/code        # 绝对路径；挂到容器 /app/bin
      slot_options: slot_on_head    # 单节点调试；多节点改用 slots
```

```bash
# code/run.sh
#!/bin/bash
set -euo pipefail
echo "hello from task: $1"          # $1 = task body
```

```bash
chmod +x code/run.sh
```

2. 建 app（脚本里务必 `< /dev/null`，否则 stdin 会被当成初始任务，INV-03）：

```bash
scalebox run --app-file=app.yaml < /dev/null
# {"app_id":<app_id>}
```

输出里可能有一条 `No standard input detected`（info 级）：它只表示 stdin 为空、没有
额外初始任务，不是错误。

3. 投一个任务并查结果：

```bash
echo "hello scalebox" | scalebox task add --app-id <app_id> --sink-module hello
scalebox task list --app-id <app_id>
scalebox task show <task_id>
```

4. 改 `run.sh` 后直接重投（`code_path` 是挂载的）；改 `app.yaml` 必须重建：

```bash
scalebox app delete --app-id <app_id>
scalebox run --app-file=app.yaml < /dev/null
```

## 可选：主路由多实例（回流分流）

主路由只有一个 slot 时，所有来源的回流串行处理。某个来源产出特别密集（例如每批几十条任务）
会把它后面的轻量来源一起堵住。这时可以给主路由开多个 slot、按来源模块分流：

```yaml
modules:
  router:
    parameters:
      task_dist_mode: SLOT-BOUND   # 必须写在主路由模块上
      slot_start_batch_size: 2
    slots:
      - h0:2                       # seq 0 / 1，同一节点
  heavy-source:
    parameters:
      router_index: 1              # 该来源的回流固定进 seq=1
```

要点：`slots` 必须落在**单节点**（跨节点时 seq 会重复）；`router_index` 与实例数不匹配时
任务不会卡死，但会静默退化为负载均衡。完整说明、坑与实测见
[`features/router-multi-instance`](../../../features/router-multi-instance/README.md)。

## 判据

- `task show` 的 `status_code == 0`，stdout 为预期文本。
- `scalebox slot status --module-id <module_id>` 至少 1 个 READY/忙碌 slot（`slot_on_head` 时在头节点）。
- 没写 `sink-tasks.txt`（单模块无下游）。

## 失败下一步

| 现象 | 下一步 |
|---|---|
| app 建了但 `t_slot` 0 行 | `slots`/`slot_on_head` 配置问题 → [`20-diagnosis.md`](../20-diagnosis.md)「slot 全是 0 行」 |
| slot 一直 STARTING / SSHFAILED | 镜像拉取或节点可达性 → 同上「slot 有记录但一直 STARTING」 |
| 任务 READY 不动 | slot 未建或归属错误 → 同上「任务 READY（-1）不动」 |
| 重复投递报 112 | 任务身份冲突 → `features/task-identity`、INV-06 |
| 脚本没跑起来 | `run.sh` 可执行位、shebang、`ACTION_RUN` 是否指向它 |
