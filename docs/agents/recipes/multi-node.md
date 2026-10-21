# recipe — multi-node（多节点落点与数据）

> 用于：把 app 从单节点扩到多节点，或排查"节点在、任务不来"。

## 前置

- 控制面可管理目标节点（`scalebox host list`）。
- 知道数据在共享盘还是节点本地盘——决定容器的数据路径。

## 步骤

1. 核对节点池（**ON ≠ 可达**，INV-16）：

```bash
scalebox host list
# 每个节点：status=ON；ssh 可达（端口/密钥正确）
```

节点名变化（Slurm 重新分配）后，同步更新 app 的 `slots` / `NODES`。

2. 写 `slots`：**短名或正则，不用逗号、不写 FQDN**（INV-15）：

```yaml
slots:
  - n0-[01]        # 匹配 n0-0、n0-1
  - n1-0
```

3. 决定数据方案：

| 方式 | 适用 | 要点 |
|---|---|---|
| 共享盘 | 所有节点读同一份数据 | 工作目录放在 `data_root` 之下（容器内 `/cluster_data_root`） |
| 节点本地盘 | 每节点独立处理本地数据 | 放到 `/tmp` 或 `/dev/shm`；节点间不共享，需要搬运模块 |
| 本地产数据 | 数据由节点自己生成 | `FXCORR_ROUTE=2` 之类模式，工作目录用直接映射路径 |

容器内只有 `/tmp`、`/dev/shm`、`/cluster_data_root`、`/local_data_root`、`/app/bin` 可见（INV-17）。

4. 建 app，验证 slot 展开：

```bash
scalebox run --env-file=<env> --app-file=app.yaml < /dev/null
scalebox module list --app-id <app_id>
scalebox slot list --module-id <module_id>     # 每个节点 × 每模块 slot 数是否符合预期
```

5. 跑一轮，按节点核对任务分布：

```bash
scalebox slot list --module-id <module_id>     # host 列
scalebox task list --app-id <app_id>
```

## 判据

- 每个模块的 slot 数 = 目标节点数 × 该模块的 slot 数；状态 ON/READY。
- 任务落到预期节点（HOST-BOUND 的 `to_host` 与业务分组一致）。
- 无 `STARTING → SSHFAILED`；数据在共享盘时各节点读到同一份。

## 失败下一步

| 现象 | 下一步 |
|---|---|
| 某模块 0 个 slot | `slots` 写错：逗号串、FQDN、正则不匹配 → `20-diagnosis`「slot 全是 0 行」 |
| 节点 ON 但连接失败 | 僵尸节点 / Slurm 重分配 / ssh 端口（如非 22 需显式 `SSH_PORT`）→ INV-16 |
| 任务全挤在一个节点 | `to_host` 缺失或 pod_id 绑定过强；HOST-BOUND 要求目标 host 有 slot |
| 节点读不到数据 | 路径不在 `data_root` 下，或用了本地盘却没搬运 |
| 大 M 跑前几批正常、后面变慢 | 冷缓存：先顺序预热共享盘数据（应用侧经验） |
