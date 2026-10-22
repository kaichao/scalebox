# Scalebox AI Agent 编程支持方案

> **历史记录（2026-10-17）**：本文为当时的方案草稿，文中路径与目录为当时状态；现行中文站结构见 `docs/cn/source/index.rst`。

| 项 | 内容 |
|---|---|
| 状态 | 草案，待确认 |
| 日期 | 2026-10-17 |
| 适用范围 | `workspace/scalebox`（文档/示例/特性）、`go-scalebox`（源码/生成脚本）、应用仓库（app-fxcorr、app-mwa 等） |
| 语言策略 | agent 面向内容**仅中文**，术语中英对照；Sphinx 面向人的发布物维持中英双语 |

## 1. 背景与目标

### 1.1 现状

- 文档站 `docs.scalebox.net` 源在 `docs/cn/source` 与 `docs/en/source`（Sphinx），分类为入门教程 / 使用指南 / 编程指南 / 附录，面向人阅读。
- 平台源码在 `go-scalebox`，`docs/` 下已有一批设计文档（`vtask-design.md`、`dag-design.md`、`env-loading.md`、`error-handling.md`、`task-headers.md`、`grpc-api.md` 等）。
- 应用仓库（`app-fxcorr`、`app-mwa`）在开发中沉淀了大量平台机制实证（`CLAUDE.md`、`docs/` 中的"坑/结案"章节、模块 README）。
- `features/` 有 9 个特性目录，格式不统一，无统一测试入口。
- `examples/` 面向人，部分 README 落后于现行 CLI。
- 无 `AGENTS.md`；平台仓库 `CLAUDE.md` 为维护者视角；应用侧没有统一的 AI 入口。

### 1.2 问题

1. **面向人**：叙述式、教程式，agent 检索与定位成本高。
2. **更新不及时**：文档有，但不全、不同步。
3. **知识分散**：平台机制的关键约束（不变式与坑）散落在应用仓库，平台文档缺失。
4. **依赖人工**：开发应用时需频繁查源码、问平台开发者。

现状实证：

- `features/README.md` 索引 5 项中 3 项指向不存在的目录（`timeout-gen`、`task-perspective`、`hello-scalebox-singularity`）；9 个实际目录中 7 个未登记（含 `timeout` 被写成 `timeout-gen`）。
- `appendix/2_parameter_reference.md` 仍含"拟修改""暂未使用""（与 xxx 的区别？）等标记。
- `examples/README.md` 仍使用 `app create` / `app list` 旧命令写法，与现行 CLI（`scalebox run` / `scalebox app list`）不一致。

### 1.3 目标与非目标

**目标**

- AI agent（Claude Code / Codex / 其他，不绑定工具）主要依靠仓库内资料完成 scalebox 应用开发：建应用、写模块、配置流控、调试、验收。
- 平台知识单源维护、可执行验证、可持续更新。
- 为应用仓库提供可复制的 AI 入口模板。

**非目标**

- 不追求平台参考文档的全自动生成（按信源成熟度分档，见 §5）。
- 不替代 Sphinx 面向人的发布物；agent 内容不进 Sphinx 构建。
- 本期不做 MCP server（列为后续可选，见 §9）。

## 2. 设计原则

1. **一套内容、两个入口**：知识单源；人走 Sphinx，agent 走 `AGENTS.md` + `features/` + `docs/agents/` + `docs/reference/`。不为 AI 另建第二套文档。
2. **三层知识模型**：不变式/意图（人写、稳定）、参考（生成或半生成、完整）、任务 recipe（人写 + 实测）。人写面越小，越跟得上平台演进。
3. **可执行优先**：一个 feature 目录 = 说明 + 集成测试 + 最小样例，三合一。
4. **单语中文 + 术语中英对照**：agent 内容只写中文；代码标识符、日志、CLI help 保留英文；`glossary.md` 做中英对照，作为未来英文翻译的锚点。
5. **显式优于隐式**：文件名、命令、配置显式表达意图；避免"依赖目录内容改变行为"的设计。
6. **实测驱动**：文档验收标准是"agent 只靠仓库资料能否独立完成典型任务"；卡点即缺口。
7. **先建后简**：先建 agent 内容，稳定后再简化 Sphinx 发布物，避免文档空窗。

## 3. 总体结构

### 3.1 仓库角色

| 仓库 | 角色 | 本方案中的职责 |
|---|---|---|
| `workspace/scalebox` | 文档站 + 标准模块 + 示例 + 特性测试 | AI 支持主载体：`AGENTS.md`、`features/`、`docs/agents/`、`docs/reference/` |
| `go-scalebox` | 平台源码（controld/agent/actuator/cli/webui） | 设计意图来源；reference 生成脚本的天然位置（贴 proto/CLI/schema） |
| 应用仓库（app-fxcorr、app-mwa 等） | 真实应用 | 验收场 + 知识蒸馏源；复制应用侧 `AGENTS.md` 模板 |

### 3.2 目标目录结构

```text
scalebox/
  AGENTS.md                     # 跨工具入口：心智模型一页 + 导航表 + 铁律摘要
  CLAUDE.md                     # 保留（维护者视角），指向 AGENTS.md
  features/                     # 可执行知识单元（agent 主参考）
    README.md                   # 索引 + 覆盖矩阵（目录 ↔ 索引一致性由脚本校验）
    <feature>/
      README.md                 # 面向 agent 的说明（模板见 §4.3）
      feature.yaml              # 最小可跑应用（与 app.yaml 同 schema）
      test.sh                   # 一键：建 app → 跑 → 断言 → 清理，PASS/FAIL
      test.meta.yaml            # 可选：断言/期望数据
      code/                     # 可选：模块代码，优先复用标准模块
  examples/                     # 面向人的完整应用（也可供 agent 参考）
  docs/
    ai-agent-support-plan.md    # 本方案
    agents/                     # agent 入口内容（仓库直读，不进 Sphinx）
      00-model.md               # 一页纸心智模型：五层概念、消息流、控制面/数据面
      10-invariants.md          # 铁律与坑（一页内，agent 每次全读）
      recipes/                  # 任务食谱：new-app / add-module / flow-control /
                                #   vtask / multi-node / debug-stuck-chain / acceptance
      glossary.md               # 术语表（中英对照）
    reference/                  # 参考层（生成或半生成，不进 Sphinx）
      cli.md                    # 命令参考（A 档，生成）
      app-yaml.md               # 应用模板字段（A 档，生成）
      module.md                 # 模块条目规范（B 档，半自动）
      task-headers.md           # task headers 全表（B 档）
      semaphore-variable.md     # 信号量/变量语义（B 档）
      env-vars.md               # 环境变量全表（B 档）
      exit-codes.md             # 退出码（A 档，生成）
      rpc.md                    # gRPC 参考（A 档，生成）
    cn/  en/                    # Sphinx 发布物（逐步简化，链接到上述内容）
```

### 3.3 三层知识映射

| 层 | 回答的问题 | 位置 | 维护方式 |
|---|---|---|---|
| 不变式 / 意图 | 为什么、红线、踩过的坑 | `docs/agents/10-invariants.md`、`docs/agents/recipes/`、`go-scalebox/docs/` 设计文档 | 人写（平台开发者），ADR 式模板 |
| 参考 | 有什么：字段/命令/表/API | `docs/reference/`、`features/*/README.md` 的配置节 | A 档生成、B 档半自动、C 档手工 |
| 任务 | 怎么做 | `features/*/`、`docs/agents/recipes/`、`examples/` | 人写 + 实测（feature 测试即验证） |

### 3.4 features 与 examples 的边界

| | features/ | examples/ |
|---|---|---|
| 内容 | 单点机制（一个参数、一个接口约定） | 完整应用场景（多模块、端到端） |
| 主要读者 | AI agent（其次是人） | 人（也可供 agent 参考） |
| 判定 | 一个测试能断言、一句话能说清 | 展示一个可用应用及其配置 |
| 例子 | `slot_expr`、`task_batch_size`、`retry_test` | `app-primes`、`hello-scalebox` |

feature 与 reference 互补：reference 给全表（有什么），feature 给深挖与验证（确实如此）。参数表每条应能指向一个 feature，或明确标注"未测试"。

## 4. 内容规范

### 4.1 feature 目录骨架

见 §3.2。要求：

- `test.sh` 幂等：可重复运行、失败可清理、返回 0/非 0。
- 优先复用标准模块；必须自带代码时放 `code/`。
- `feature.yaml` 只覆盖该特性，不引入无关模块。

### 4.2 feature.yaml 决策

- **命名**：用 `feature.yaml`，不用 `test.yaml`（避免与单元测试配置混淆）；与 `app.yaml` **同 schema**，复用同一校验。
- **加载方式**：`scalebox run --app-file=feature.yaml`（现有能力，零平台改动）。
- **不做缺省探测**：不在 `scalebox run` 中增加"自动查找 feature.yaml/test.yaml"的行为。理由：
  - 隐式行为依赖目录内容，agent 难以排错；
  - feature 的真正入口是 `test.sh`（建 → 跑 → 断言 → 清理），换文件名不解决自动化；
  - 若需"一键跑全部 feature"，应新增显式命令（如 `scalebox feature test`，扫描目录 + 索引），属于后续增强，本期不做。
- **feature 专属元数据**（断言、期望输出）放 `test.sh` 或 `test.meta.yaml`，不污染应用模板 schema。

### 4.3 feature README 模板

```markdown
# <feature 名称>

一句话说明该特性。

## 配置

| 字段 | 位置（app/module/argument） | 取值 | 默认值 | 说明 |
|---|---|---|---|---|

## 行为与不变式

- 规则 1；规则 2……

## 坑

- 坑 1（现象 → 原因 → 规避）。

## 最小例子

见 `feature.yaml`；关键片段：

```yaml
...
```

## 测试

```sh
./test.sh
```

预期：PASS，断言点……

## 源码指针

- `go-scalebox/<path>:<line>` —— 该特性的实现位置

## 相关文档

- 链接到 Sphinx 章节 / reference / 设计文档
```

### 4.4 不变式条目模板

```markdown
### INV-<编号> <一句话规则>

- **现象**：什么场景下遇到。
- **规则**：必须遵守什么 / 绝对不能做什么。
- **为什么**：机制层面的原因。
- **反例/后果**：违反时的具体表现（越具体越好）。
- **源码指针**：`go-scalebox/<path>:<line>`。
- **置信度**：已交叉验证（≥2 个应用）/ 单项目来源。
```

### 4.5 AGENTS.md

**平台仓库（`scalebox/`）入口**：一页内包含

- 心智模型：App / Module / Task / Slot / VTask 五层概念，消息流与控制面/数据面。
- 导航表："要做 X → 读 Y"。
- 铁律摘要：最关键的 5–10 条，指向 `docs/agents/10-invariants.md`。

**应用仓库模板**：可复制到 app-fxcorr/app-mwa 等仓库，内容为

- 本应用的流水线拓扑与入口命令；
- 平台资料入口（指向 scalebox 仓库的 features/reference/invariants）；
- 本应用的私有约定与坑（仅本应用有效的知识留在应用仓库）。

`CLAUDE.md` 保留并指向 `AGENTS.md`，避免同一内容两处维护。

### 4.6 reference 规范

- 生成物头部标注来源与命令：

  ```markdown
  <!-- AUTO-GENERATED from proto/*.proto by make docs-gen; DO NOT EDIT -->
  ```

- 每篇带"最后校验版本/日期"字段。
- 关键条目带源码指针（包路径或 `file:line`），文档过期时 agent 可顺指针自查源码。
- 语义说明与结构事实分开：结构（名字/类型/默认值）可生成，语义（何时用、副作用）手写。

## 5. 生成管线

### 5.1 分档

| 档 | 内容 | 信源状态 | 方式 |
|---|---|---|---|
| A | CLI 命令、app.yaml 字段、exit code、RPC 列表/签名 | 已结构化：Cobra help、webui 的 app.yaml JSON Schema、规范文件、proto | **现在生成**，先做 2–3 个样板 |
| B | app/module parameters 语义表、环境变量表 | 半结构化：Go 代码中的定义/`Getenv`，注释不足 | **半自动**：脚本抽取草稿 + 人工/AI 对照源码润色；推动源码侧集中定义/加注释后转全自动 |
| C | 不变式、recipe、设计意图、原理性文档 | 非结构化 | 手工维护，不生成 |

### 5.2 第一批样板

1. `reference/cli.md`：从 `scalebox --help` / 各子命令 `--help` 生成。
2. `reference/app-yaml.md`：从 webui 现有 app.yaml JSON Schema 生成字段表。
3. `reference/exit-codes.md`：从现有规范文件（`appendix/6_exit_code_spec.md` 与源码常量）生成。

跑通后再扩 A 档其余项。

### 5.3 校验机制

- `make docs-gen` 可复现生成结果。
- 检查脚本比对"重新生成结果与提交内容"，不一致即失败（CI 或本地钩子）。
- 目的：让文档落后在提交阶段暴露，而不是靠人记得更新。

### 5.4 脚本放置

- 生成脚本放 `go-scalebox`（紧贴源码与 proto），产物提交到 `scalebox/docs/reference/`。
- 若跨仓库不便，最低要求：脚本读源码的路径可配置，且生成命令写入产物头部。

## 6. 知识回流

### 6.1 蒸馏来源（按优先级）

1. `app-fxcorr`、`app-mwa` 的 `CLAUDE.md`、`docs/`（尤其"坑/结案"章节）、模块 README。
2. `examples/` 下的复杂应用（补充场景与配置样例）。
3. 新应用开发过程中的新发现（边做边沉淀）。

### 6.2 蒸馏流程

1. 按 §4.4 模板记录（现象/规则/为什么/反例/源码指针）。
2. 分流：跨应用机制 → `docs/agents/10-invariants.md`；单点特性 → `features/<name>/README.md`（无目录则新建并登记索引）。
3. **去项目化**：只保留机制，不出现项目私有参数（如 `FXCORR_*`）与数据布局。
4. 标注置信度：两个及以上应用出现过的规则标"已交叉验证"；单项目来源单独标注。

### 6.3 持续机制

- 应用仓库开发流程中加一条轻规则：发现平台机制的新发现 → 提交到 scalebox 仓库对应 feature/invariant。
- 可定期（如每轮应用开发结束）由 AI 扫描应用仓库增量，产出合并提案，人工确认后入库。

## 7. 执行计划

| 阶段 | 内容 | 产物 | 依赖 | 验收 |
|---|---|---|---|---|
| 0 | 盘点矩阵（只读） | 特性覆盖矩阵：特性 | 参数 | 测试 | 文档 | 坑记录 | 来源 | 缺口等级 | 无 | 覆盖 parameters/arguments 主要特性；缺口分级 |
| 1 | 模板与索引 | feature README 模板、invariant 模板、`features/README.md` 索引、索引校验脚本 | 0 | 现有 9 个目录套模板；索引与目录一致 |
| 2 | 蒸馏第一批 | `docs/agents/10-invariants.md` 初版；从 app-fxcorr/app-mwa 补齐的 features | 1 | 每条有源码指针；去项目化；置信度标注 |
| 3 | feature 测试骨架 | `test.sh` 约定 + 首批可判定测试 | 1 | 一键运行、PASS/FAIL、幂等清理 |
| 4 | reference 生成样板 | `cli.md`、`app-yaml.md` 生成脚本 + `make docs-gen` + 校验 | 1 | 源码变 → 重新生成 → diff 可见 |
| 5 | AI 实测 | 卡点清单 + 文档修复 | 2、3、4 | agent 只靠仓库资料完成典型任务 |
| 6 | Sphinx 简化 | 内容迁移映射；操作细节改指针 | 5 | 无重复内容；链接可用 |

阶段 2、3、4 可并行；阶段 5 是总验收，必须用真实任务（如"从零写一个 N 模块应用并跑通"、"诊断一个卡住的链"）。

### 阶段 0 盘点输入

- `docs/cn/source/appendix/2_parameter_reference.md`（参数表）
- `features/` 现有目录与 `features/README.md`
- `docs/cn/source/developer/`（编程指南各章）
- `go-scalebox/docs/*.md`（设计文档）
- `app-fxcorr`：`CLAUDE.md`、`docs/*.md`（坑/结案）、`router/README.md`、`modules/README.md`
- `app-mwa`：`CLAUDE.md`、`docs/`、各模块 README
- `examples/` 各 README

## 8. 验收标准

1. **任务完成**：agent 在只给仓库资料的条件下，完成"新建应用 / 加模块 / 配流控 / 调试卡链"中的典型任务。
2. **可执行**：每个已登记 feature 有可判定的 `test.sh`，且索引与目录一致。
3. **可追溯**：每条不变式有源码指针；参考文档标明生成来源或最后校验日期。
4. **无重复**：Sphinx 发布物与 `features/`、`reference/` 不重复维护同一事实。
5. **可持续**：源码变更能被生成/校验机制捕获；应用侧新发现能按流程回流。

## 9. 风险与待决事项

| 风险 | 应对 |
|---|---|
| 模板未定就蒸馏，内容走样 | 阶段 1 先行；阶段 2 严格按模板 |
| 信源质量不足，生成物只有骨架 | A 档先做；B 档半自动；推动源码注释/集中定义 |
| 追求全覆盖，工期失控 | 盘点矩阵定优先级；先最小闭环（invariants + 2–3 篇 reference + 2–3 个 recipe） |
| 双写导致再次不同步 | Sphinx 只留原理 + 指针；agent 内容单源 |
| 隐式行为增加排错成本 | `feature.yaml` 显式加载；不新增 `scalebox run` 缺省探测 |

**待决**

1. `docs/agents/` 与 `docs/reference/` 放 `docs/` 根（仓库直读、不进 Sphinx）已在 §3.2 按此设计，需确认。
2. 生成脚本放 `go-scalebox`（建议）还是 `scalebox`，需按团队仓库权限确认。
3. `scalebox feature test` 等自动发现命令是否纳入后续版本（本期不做）。
4. MCP server（包装 webui REST/CLI 查询 app/task/slot/日志）列为后续可选，不在本期范围。

## 附录 A 现状 → 目标映射

| 现状 | 目标 | 处理 |
|---|---|---|
| `docs/cn/source/appendix/2_parameter_reference.md` | `docs/reference/` | A/B 档重构；Sphinx 保留摘要 + 链接 |
| `docs/cn/source/appendix/3_commandline_tools.md` | `docs/reference/cli.md` | A 档生成；Sphinx 指向 |
| `docs/cn/source/appendix/6_exit_code_spec.md` | `docs/reference/exit-codes.md` | A 档生成 |
| `docs/cn/source/developer/*` | 保留原理内容；操作细节迁移 | 阶段 6 逐步执行 |
| `docs/cn/source/user/*` | 保留（安装/运行/示例导读） | 基本不动 |
| `features/*/README.md` | 按模板重写 + `test.sh` | 阶段 1–3 |
| `examples/*` | 保持面向人；修正过时命令 | 随阶段 6 顺手修 |

## 附录 B feature 盘点表模板

| 特性/参数 | 位置 | 现有测试 | 现有文档 | 实战记录 | 缺口 | 去向 |
|---|---|---|---|---|---|---|
| 例：`task_dist_mode` | module-parameters | 无 | appendix/2 | 无 | 测试+文档 | `features/task-dist-mode/` |
| 例：`retry_rules` | module-parameters | `features/retry_test` | appendix/2 | 无 | 模板统一 | `features/retry_test/` |
