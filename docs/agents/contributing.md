# 知识回流（contributing）

> **何时用**：应用开发中发现平台级机制（新坑、新行为、新参数、新边界）之后；以及**每轮应用开发结束时**。
> **目标**：把"一次踩坑"变成"下次不踩"，并保证每个事实只有一个信源。
> 应用侧入口页（如 `app-mwa/AGENTS.md`）会指到本页。

## 1. 先判定落点

判据只有一句：**换一个应用还会不会遇到？** 会 → 平台；不会 → 应用。

| 发现类型 | 落点 | 维护方式 |
|---|---|---|
| 平台规则、红线、坑（为什么 + 违反后果） | [`10-invariants.md`](./10-invariants.md) | 按 [`_template-invariant.md`](./_template-invariant.md) 手写 |
| 可跑的最小特性（配置 + 行为 + 测试） | `features/<kebab-name>/` | 按 [`_template/README.md`](../../features/_template/README.md) + `feature.yaml` + `test.sh` |
| 结构性事实（命令、字段、状态码、header、环境变量、RPC） | [`docs/reference/`](../reference/) | 生成物改生成器；手写参考对照源码 |
| 设计原理、架构、跨模块机制 | go-scalebox `docs/` 或 Sphinx `design/` | 平台开发者 |
| 应用业务约定（目录布局、数据格式、业务命名） | 应用仓的 `AGENTS.md` / `docs/` | 应用作者 |
| 单点、未验证、只在一个应用出现 | 应用仓，标注"未交叉验证" | 先不进 invariant |

**去项目化**是第一步：删掉项目私有参数（`FXCORR_*`、`DIR_*` 之类业务变量）、节点名、数据布局。
例：把"`FXCORR_BATCH_RANGE` 必须同时设 `X_MERGE_FORCE`"抽象成"部分批运行时，下游归并的期望集合
必须与投递集合同源"。

## 2. 各落点的动作

### 2.1 invariant

1. 查重：先看 [`feature-inventory.md`](../feature-inventory.md) 与现有 INV，避免同规则两条。
2. 按 [`_template-invariant.md`](./_template-invariant.md) 填：现象 / 规则 / 为什么 / 反例 / 源码指针 / 置信度 / 适用范围。
3. 置信度：≥2 个应用或应用+官方文档 → **已交叉验证**；否则标注**单项目来源**。
4. 同步文件顶部的 TL;DR 列表；同一主题放在同一小节。

### 2.2 feature

1. 建目录 `features/<kebab-name>/`，README 按模板（配置 / 行为与不变式 / 坑 / 最小例子 / 测试 / 源码指针）。
2. `feature.yaml` 做最小可跑形态；`base_image` 用 `scalebox.net/platform/agent:latest`。
3. `test.sh` 按 [`lib/testkit.sh`](../../features/lib/testkit.sh) 契约：`tk_init` → `tk_run_app` →
   `tk_wait_tasks` → 断言；不能自动判定的前置条件用 `tk_skip`（退出码 77）。
4. 更新 [`features/README.md`](../../features/README.md) 索引（`check-index.sh` 会校验）。
5. 若 Sphinx 章节讲同一特性：Sphinx 只留指针，正文放 feature README。

### 2.3 reference

- **生成物**（`cli.md`、`app-yaml.md`、`exit-codes.md`、`rpc.md`、`env-vars.md`）：改 go-scalebox 的
  生成器或数据源（schema/proto/源码），再 `make docs-gen`；**不要手改生成物**。
- **手写参考**（`task-headers.md`、`module.md`、`semaphore-variable.md`）：对照源码修改，注明依据。
- 新增事实类别时更新 [`reference/README.md`](../reference/README.md) 的信源表。

## 3. 每轮应用开发结束的检查清单

```text
[ ] 扫应用仓本轮提交、docs/、CLAUDE.md/AGENTS.md：列出新坑、新参数、新行为
[ ] 对每条判定落点（§1），去项目化（§2 开头）
[ ] 平台级条目：按模板写入，更新索引 / TL;DR / 信源表
[ ] 生成物：make docs-gen && make docs-check
[ ] 校验：check-links.py --lang cn --scope all；check-index.sh；check-content.py；sphinx-build -W
[ ] 应用侧 AGENTS.md：本地约定是否需要更新（不回写平台规则）
[ ] 产出合并提案：文件清单 + 每条一句话理由 + 置信度，人工确认后提交
```

**合并提案**示例（给平台开发者看的最小形式）：

```text
1. docs/agents/10-invariants.md  INV-21 部分批运行需同步归并期望集合（已交叉验证：app-fxcorr）
2. features/partial-batch/       新增 feature（单项目来源）
3. docs/reference/task-headers.md  补 _vtask_attr_ 的转发边界（源码核对）
```

## 4. 反模式

- **复制平台规则到应用仓**：两份必然腐化，应用仓只放指针与本地约定。
- **把单次实验写成规则**：样本量 1 先标"单项目来源"，等第二个应用复现再升级。
- **手改生成物**：下一次 `make docs-gen` 会覆盖，CI 漂移检查会红。
- **落点错位**：在 reference 里写"为什么"（应进 invariant）；在 invariant 里抄字段表（应进 reference）。
- **新增 feature 不更新索引**：`check-index.sh` 会失败；README 与目录必须同步。
- **提交拆分错误**：新增 invariant 与它的 TL;DR 更新同 commit；feature 的 README/yaml/test.sh 同 commit；
  源码与生成物同 commit。

## 5. 相关

- 应用侧入口：[`app-mwa/AGENTS.md`](https://github.com/kaichao/app-mwa/blob/main/AGENTS.md)（模板参考）
- 平台入口：[`AGENTS.md`](../../AGENTS.md)
- 盘点矩阵：[`feature-inventory.md`](../feature-inventory.md)
