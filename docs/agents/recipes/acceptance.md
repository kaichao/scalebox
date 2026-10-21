# recipe — acceptance（一轮运行的验收）

> 用于：实验/生产跑完一轮后判定"这轮算不算成功"。**不要只看 app 状态**。

## 前置

- 一轮任务已结束（或有明确的观察窗口）。
- 已知本轮的预期：任务条数、产物清单、判据（文件数、行数、md5 等）。

## 步骤

1. 任务全绿：

```bash
scalebox task list --app-id <app_id>
# 期望：没有 status_code > 0；终态任务数 = 预期条数
```

2. 计数核对（防止"少跑一批"被当成成功）：

```sh
docker exec database psql -U scalebox -d scalebox -c \
  "select status_code, count(*) from t_task t join t_module m on t.mod_id=m.id
   where m.app=<app_id> group by status_code order by 1;"
```

3. 产物核对：按本轮预期逐项验证（文件数、总大小、行数、校验和）。有基准就逐字节比对，
   没有基准就与上一轮一致（同配置复跑的稳定性）。

4. 失败与超时扫描：

- `status_code`：`0` 完成、`>0` 失败（112/123/124 等见 [`../20-diagnosis.md`](../20-diagnosis.md)）。
- 扫 `t_task_exec.stderr` 与 `extras` 中的 ERROR/异常（拓扑上全绿也可能有静默错误）。

5. 清理：

```bash
scalebox app delete --app-id <app_id>
# app delete 不停止本地 slot 容器；按 APP_ID 环境变量停容器（见 20-diagnosis）
```

6. 复跑一轮（重要实验建议）：同配置重跑，确认结果一致与幂等（无残留、无 112）。

## 判据

| 项 | 通过条件 |
|---|---|
| 任务 | 无 `>0`；终态数 = 预期条数 |
| 产物 | 清单齐全；大小/行数/校验和符合基准 |
| 日志 | 无未处理的 ERROR |
| 资源 | 信号量回到初值；无在制 vtask 残留 |
| 现场 | app 删除后无容器/目录残留 |

## 失败下一步

| 现象 | 处理 |
|---|---|
| 任务数少了 | 查是否部分批未投递（`sink-tasks` 回流链断）或失败即停 |
| 有少量失败码 | 失败任务不回流；按码定位（112 重投、124 超时、1–10 参数/输入） |
| 全绿但产物不对 | 数据面问题：路径映射、共享盘、上游读到的数据不完整 |
| 复跑结果不一致 | 先排除并发造数/写数据；再查非确定性算法或时间戳 |
