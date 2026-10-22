# 阶段 7 — 中文站文档重建归档

> 状态：已完成（2026-10-22）。
> 范围：`docs/cn/source` 中文站；英文站未检查、未构建。
> 用途：替代 `tmp/scalebox-cn-stage-a/` 中的一次性审计制品，保留最终结构、关键决策与验收结果。

## 1. 最终结构

- 30 个源文件：29 个页面 + `index.rst`。
- 六大板块与目录一一对应：

| 板块 | 目录 | 页面数 |
|---|---|---:|
| 一、入门教程 | `started/` | 3 |
| 二、应用开发 | `application/` | 8 |
| 三、部署与运维 | `deployment/` | 2 |
| 四、设计与原理 | `design/` | 7 |
| 五、参考 | `reference/` | 4 |
| 六、附录与其他 | `appendix/` | 5 |

- 29 个页面全部位于对应板块目录；文件名前缀与页面 H1 编号一致。
- 目录重组涉及 26 个文件重命名；旧 URL 未加重定向（明确选择）。
- 当前入口：`docs/cn/source/index.rst`。

## 2. 编号规则

- 板块标题（toctree caption）：汉字序号，一、至六、。
- 页面 H1：板块内阿拉伯序号。
- 小节 H2/H3：`x.y` / `x.y.z`。
- 不使用 `.0` 小节。原 `3.3.0`、`4.5.0` 已收口为 `3.3.1`、`4.5.1`，后续小节顺延。

## 3. 核心概念定稿

- 计算逻辑 = 算法逻辑（模块内）+ 并行逻辑（模块间）。
- 核心概念采用“3+1”：计算逻辑、应用状态、主路由 + VTask。
- 主路由是应用控制面模块，执行并行逻辑、读写应用状态、分发任务、组织模块协作；不是平台控制面，不执行算法逻辑。
- 应用状态存放在平台存储：semaphore / variable / global；跨物理集群的信息共享使用 queue。
- 编程原语单独成层：semaphore、semagroup、variable、global、queue。
  - global：同一物理集群内跨应用。
  - queue：跨物理集群的管理面消息队列，不是 wait-queue，也不是平台内部任务队列。
- 主路由原则上单实例；高负载可按模块分设 SLOT-BOUND 多实例，使用 `router_index` 分流。同一应用状态不得被多个实例并发写，并发正确性由程序员保证。
- VTask 是并行逻辑及其状态的粗粒度封装，不是独立调度器，不封装算法逻辑。
- 算法模块无状态；task 不携带应用状态；模块内并行属于算法逻辑，模块级/模块间并行属于并行逻辑。

## 4. 读者定位

- 中文站正文面向人：用户、应用开发者、模块开发者、系统管理员。
- AI agent 与自动化内容留在仓库直读：
  - `AGENTS.md`
  - `docs/agents/`
  - `features/`
  - `docs/reference/`
- 中文站只链接上述内容，不复制。VTask 实现索引与信号量生命周期表保留在站点，但标注“面向平台贡献者/高级开发者”。

## 5. 施工批次

| 批次 | 内容 | 结果 |
|---|---|---|
| B1 | 核心概念页重写 | 完成 |
| B2 | 主路由章重写、导航归位 | 完成 |
| B3 | 编程模型、应用设计接入 | 完成 |
| B4 | VTask、节点本地计算、高级特性、跨集群接入 | 完成 |
| B5 | 入门教程、示例、模块开发接入与去重 | 完成 |
| B6 | FAQ、首页、链接与术语终审 | 完成 |
| B7 | 命令行参考、FAQ、about 页编号修复 | 完成 |
| B8 | 概念入口、编程原语补充、审计收口 | 完成 |
| A | 板块汉字序号、页面与小节重编号 | 完成 |
| A2 | 附录补号、主路由标题精简、通信机制修正 | 完成 |
| C | 六板块目录重组与文件重命名 | 完成 |
| D | 去掉 `.0` 小节并顺延编号 | 完成 |

主要提交：`30d204a`、`c025c7f`、`e39f746`、`b5a750a`。

## 6. 验收

- `make -C docs clean && make -C docs build-cn`：Sphinx 9.1.0 + myst 5.0.0，30 个源文件，构建成功，无告警。
- `python3 docs/check-links.py --lang cn --scope all`：0 缺失（归档前 71 个文件；归档后 72 个）。
- 六板块 H1 编号连续；无 `.0` 标题；H2/H3 与所属页面编号一致。
- `README.md` 与 `README.zh-CN.md` 的文档表均覆盖 29 个页面，六板块分类一致。
- 活跃文档中无旧 `user/`、`developer/`、旧 `appendix/` 路径引用。

## 7. 旧审计制品处置

`tmp/scalebox-cn-stage-a/` 为本次施工的临时审计目录，已被 `.gitignore` 忽略。本文提交后可按以下方式处置：

- 可直接删除：`heading-index.md`、`content-inventory.md`、`audience-screening.md`、`heading-denumber-map.md`。
- 如需保留底稿：`concept-spec.md`、`target-ia-and-mapping.md`、`heading-renumber-map-a.md`；长期以本文为准。
- 删除前确认本文已提交；`tmp/` 不在版本控制中，删除后无法从 git 恢复。

## 8. 相关记录

- `docs/agents/phase6-sphinx-simplification.md`：阶段 6 Sphinx 简化记录（历史）。
- `docs/ai-agent-support-plan.md`：AI agent 支持方案（历史）。
- `docs/feature-inventory.md`：特性覆盖矩阵（历史盘点）。
- `docs/cn/source/index.rst`：当前中文站入口。
