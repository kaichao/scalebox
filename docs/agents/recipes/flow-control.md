# recipe — flow-control（加闸门，不死锁）

> 用于：限制在制批数、限制磁盘占用，或在资源紧张时让 slot "停领"而不是失败。
> 细则见 `features/check-test`、INV-12/INV-14。

## 前置

- 链是"批/组"结构（一批任务有共同的入口模块，如 vtask-head）。
- 知道要控制的维度：**数量**（在制批数）还是**空间**（磁盘剩余/目录配额）。

## 步骤

1. 选闸门类型：

| 维度 | 手段 | 挂载位置 |
|---|---|---|
| 在制批数 | 批级数量闸：`vtask_size`（vtask 管道）或 slot 数限制 + `VTASK_SIZE` | vtask 层 |
| 磁盘剩余 | `arguments.free_space_gb`（保留空间） | **批入口模块**（如 vtask-head，每节点 1 slot） |
| 目录配额 | `arguments.dir_quota_gb` | 同上 |
| 自定义条件 | `ACTION_CHECK` 脚本 | slot 领取前执行 |

2. 把空间闸门挂在批入口，**不要挂任务级模块**：

```yaml
modules:
  vtask-head:
    base_image: scalebox.net/platform/agent:latest
    arguments:
      free_space_gb: '{"/data":50}'    # 该目录至少保留 50 GB 才领任务
      slot_options: slot_on_head
    parameters:
      vtask_role: head
      vtask_size: 1                    # 先串行，消除首轮盲区
```

3. 重建 app → 跑一批 → 人为制造"空间不足"验证拦截：

```bash
# 让占用超过阈值（测试机上用 loop 文件或已有大数据），然后投递
scalebox task list --app-id <app_id>       # 任务停在 -1，fail=0
```

4. 释放空间，确认自动放行（无需重启 app/slot）：

```bash
scalebox task list --app-id <app_id>       # 数秒内转为 -3/0
```

## 判据

- 空间不足时：新任务保持 READY（`-1`），`fail=0`，没有任务失败。
- 空间恢复后：任务在轮询间隔内自动开始（默认 6 秒）。
- 在制批数 ≤ `vtask_size`；信号量值始终在 `[0, vtask_size]`。

## 失败下一步

| 现象 | 下一步 |
|---|---|
| 任务一直 `-1`、fail=0、数据不增长 | 闸门在工作；查 `free_space_gb` → `dir_quota_gb` → `ACTION_CHECK` 顺序 |
| 死锁：组内部分任务占资源、部分等闸门 | 闸门挂到了任务级 → 移到批入口（INV-14） |
| 在制批数超过 `vtask_size` | 检查 head 是否对 `vtask_size:wait-queue` 加 1（INV-10）、是否重复 unbind（INV-11） |
| 一启动就被拦 | 首轮盲区：`VTASK_SIZE=1` 串行化；或阈值高于当前可用空间 |
