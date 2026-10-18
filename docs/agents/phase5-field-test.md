# 阶段 5 — AI 实测记录

日期：2026-10-18。执行环境：测试机 `ssh fxcorr`（控制面 + local 集群，host `h0.local`）。
目标：只用仓库资料完成两个典型任务，记录卡点并修复文档；不询问平台开发者。

## 任务 A：从零写一个两模块应用并跑通

**任务**：写 `producer → consumer` 的 app，producer 在 `${WORK_DIR}/sink-tasks.txt`
回流，consumer 读取 `_demo` header；入口任务带 `_demo=hello`。

**使用资料**：`features/headers-propagation/README.md`、`docs/reference/app-yaml.md`、
`docs/reference/cli.md`、`features/README.md`、`docs/agents/10-invariants.md`。

**结果**：通过。app 104 的 consumer stdout 为 `_demo=hello`；清理后无残留。

**卡点**：

| # | 卡点 | 级别 | 修复 |
|---|---|---|---|
| A1 | 仓库没有 agent 入口页；README 是英文平台介绍，不知道从 features 还是 reference 开始 | 高 | 新增根 `AGENTS.md`（心智模型 + 导航表 + 铁律 + 常用命令） |
| A2 | `base_image` 必填，但测试部署没有默认 `platform/base:latest`；reference 只列字段，不说明该用哪个镜像 | 中 | `AGENTS.md` 增补“测试部署用 `platform/agent:latest`”；`features/README.md` 已有约定 |
| A3 | 入口任务如何投递（`task add` 从 stdin 读 body、`--sink-module`、`-H`）只在 feature README 的例子里；CLI reference 只有 flags | 中 | `AGENTS.md` 常用命令补完整示例 |
| A4 | 查任务结果需要知道用 `task list` / `task show`（含 stdout）；没有“如何确认成功”的入口 | 中 | `AGENTS.md` 常用命令列出；诊断文档给出 60 秒检查顺序 |
| A5 | `scalebox app delete` 只删数据库记录，不停止本地 slot 容器（实测容器仍在） | 高 | `docs/agents/20-diagnosis.md` 增加清理节；`AGENTS.md` 命令旁括注 |

## 任务 B：诊断一个卡住的链

**任务**：创建一个 app，模块 `slots` 写成测试机不存在的节点；观察 app RUNNING、任务 READY，
只用仓库资料定位并修复。

**症状**：`scalebox app show` 显示 RUNNING；`scalebox task list` 有 1 个任务 `status=-1`；
`t_slot` 0 行；`scalebox host list` 显示 `h0.local`。

**使用资料**：`docs/agents/10-invariants.md` INV-15、`docs/reference/app-yaml.md` 的 `slots`、
`features/pod/README.md`、go-scalebox `cmd/controld/grpc/app_slots.go` 的
`expandRegexHostList`。

**过程**：先按 `host list` 的完整 hostname 把 `slots` 改成 `h0.local`，仍不建 slot；
读源码发现匹配对象是 `split_part(hostname,'.',1)`（短名），改成 `h0` 后 slot 建出
（`h0.local`，ON），任务成功，stdout `worker ok: 0`。

**卡点**：

| # | 卡点 | 级别 | 修复 |
|---|---|---|---|
| B1 | 没有“卡链症状 → 检查顺序”的诊断入口 | 高 | 新增 `docs/agents/20-diagnosis.md` |
| B2 | `slots` 匹配短名（hostname 第一段）未写清；`host list` 显示 FQDN 会诱导写成 `h0.local` | 高 | INV-15 补充规则与反例；`features/pod/README.md` 补坑；schema `slots` description 更新并重新生成 `app-yaml.md` |
| B3 | `scalebox app show` 不显示 module/slot 摘要，需要 `module list` + `slot list` 组合 | 中 | 诊断文档给出命令组合；CLI 改进留待后续 |
| B4 | `slots` 写错时没有报错，只有 slot 0 行；需要知道查 `t_slot` | 中 | 诊断文档“slot 全是 0 行”一节给出 CLI 与数据库兜底 |

## 修复清单

- 新增 `AGENTS.md`（平台 agent 入口）。
- 新增 `docs/agents/20-diagnosis.md`（症状驱动诊断）。
- 更新 `docs/agents/10-invariants.md` INV-15（短名规则、反例、置信度）。
- 更新 `features/pod/README.md`（slots 短名说明与坑）。
- 更新 `features/README.md`（AGENTS/诊断入口）。
- 更新 `scalebox/CLAUDE.md`（指向 AGENTS/诊断）。
- 更新 go-scalebox `cmd/webui/frontend/src/schema/app-schema.json` 的 `slots` description；
  重新生成 `docs/reference/app-yaml.md`（`make docs-check` OK）。

## 复测与验证

- 任务 A：app 104 跑通（`_demo=hello`），按诊断文档清理，无残留。
- 任务 B：`slots: h0.local` 卡住 → 改 `slots: h0` 后 app 107 任务成功（`worker ok: 0`）。
- `make docs-check` OK；`features/check-index.sh` OK。

## 遗留（不阻塞阶段 5）

- `scalebox app show` 不显示 slot/task 摘要；agent 需组合 `module list`/`slot list`/`task list`。
- 根 `README.md` Quick Start 仍用 `build/` 路径（阶段 0 已知，阶段 6 统一修）。
- 未提供 `test.sh` 的 feature README 仍保留“阶段 3 提供”的占位注释（本次只修 pod；其余随阶段 6 或后续测试补齐）。
- 测试部署没有默认 `platform/base:latest`，属于部署镜像问题，已在文档中说明替代镜像。
