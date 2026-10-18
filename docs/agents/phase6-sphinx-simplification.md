# 阶段 6 — Sphinx 简化执行记录

日期：2026-10-18。范围：`docs/cn/source/`（中文站；英文站留到 6b）。
目标：Sphinx 只保留原理、设计、教程、规范；操作细节指向 `features/`、`docs/reference/`、`docs/agents/`；无重复内容，链接可用。

## 拍板

- 外链基底：GitHub `dev`（`https://github.com/kaichao/scalebox/blob/dev/...`）。
- 中英范围：先中文；英文 6b。
- `appendix/2` 的 task-headers、cluster/host/slot 参数、环境变量保留在 Sphinx。
- 构建验收：`ssh fxcorr`（Python 3.9.25 + venv + Sphinx 6.2.1）。

## 基线环境与命令

本机：Python 3.14.2，仅运行标准库脚本。

测试机 `ssh fxcorr`：Python 3.9.25；`python3 -m venv` 可用；PyPI 可达。
临时目录：`/tmp/scalebox-docs-phase6/`（rsync 自本机 `docs/`，不写仓库）。

```sh
rsync -a --delete --exclude build docs/ fxcorr:/tmp/scalebox-docs-phase6/
ssh fxcorr 'cd /tmp/scalebox-docs-phase6 && python3 -m venv .venv && \
  .venv/bin/python -m pip install -q -r requirements.txt'
ssh fxcorr 'cd /tmp/scalebox-docs-phase6 && \
  .venv/bin/sphinx-build -b html cn/source build/html'
ssh fxcorr 'cd /tmp/scalebox-docs-phase6 && \
  .venv/bin/sphinx-build -W -b html cn/source build/html-w'
```

静态检查：

```sh
python3 docs/check-links.py --lang cn
python3 docs/check-content.py
```

## 基线结果

| 检查 | 结果 |
|---|---|
| 链接检查 | 29 个文件；6 个目标缺失；19 处 RST doc 角色写在 Markdown 中 |
| CLI 覆盖 | `appendix/3` 提取 78 条；`reference/cli.md` 未覆盖 23 条 |
| 参数覆盖 | `appendix/2` §2.1.1–2.1.3 提取 50 个；`app-yaml.md` 未覆盖 7 个 |
| 退出码覆盖 | 256 个；`exit-codes.md` 未覆盖 0 个 |
| Sphinx 构建 | exit 0；1 warning |
| `sphinx-build -W` | exit 2（warning 视为错误） |

构建 warning：

```text
developer/7_advanced_features.md:323: duplicate label ***应用定义***,
other instance in developer/7_advanced_features.md
```

## 发现与处理决定

### F-01 RST doc 角色在 Markdown 中不生效（高）

MyST 的角色语法是 `{doc}`，不是 RST 的冒号 doc 语法。构建产物中 19 处该语法全部以字面文本出现，没有链接；Sphinx 也不报 warning。

处理：

- 6.5/6.6 将 `.md` 中的 doc 角色全部改为 `{doc}`（已验证可生成链接）。
- 其中 6 个指向 `../../go-scalebox/docs/*` 的目标在 Sphinx 源树外，go-scalebox 无 GitHub 远端；改为文字指针（“见 go-scalebox 仓库 `docs/xxx.md`”），不保留断链。
- `check-links.py` 现在对 `.md` 中的 doc 角色输出 `SYNTAX` 行，作为后续验收项。

缺失目标（6 个，全部为 go-scalebox 文档）：

```text
developer/10_cross_cluster.md  ../../go-scalebox/docs/cli-grpc-architecture
developer/11_security.md       ../../go-scalebox/docs/security
developer/9_vtask.md           ../../go-scalebox/docs/vtask-review
user/1_installation.md         ../../go-scalebox/docs/security
user/7_webui.md                ../../go-scalebox/docs/rest-api-reference
user/8_vscode.md               ../../go-scalebox/docs/webui-vscode-plan
```

### F-02 重复的 autosectionlabel（中）

`developer/7_advanced_features.md` 中 `***应用定义***` 出现两次（约 278 行与 323 行），触发 duplicate label。

处理：6.4 把第二处改名为“应用定义（slot 参数示例）”，保留第一处。

### F-03 内联三反引号（低）

`appendix/3` 有 4 处、`developer/7` 有 1 处形如 ```` ```value``` ```` 的内联代码。它们不是未闭合围栏，但会干扰工具判断；6.2/6.4 顺带改为单反引号。

### F-04 CLI 覆盖缺口（高）

`appendix/3` 有 23 条命令未出现在生成参考中，需在 6.2 逐条判定“改名 / 移到别章 / 已废弃”：

```text
scalebox app add-remote        scalebox app create
scalebox app set-finished      scalebox channel pull
scalebox channel push          scalebox cluster get-parameter
scalebox cluster show          scalebox event add
scalebox event misc-add        scalebox event slot-add
scalebox event task-add        scalebox fs ls
scalebox fs stat               scalebox global list
scalebox host create           scalebox semaphore list
scalebox slot get-parameter    scalebox slot remove
scalebox slot set-parameter    scalebox task log
scalebox task remove-header    scalebox variable list
scalebox vtask add-semaphore-value
```

已知改名候选：`global list`→`global ls`、`semaphore list`→`semaphore ls`、`slot remove`→`slot delete`、`variable list`→`variable ls`、`host create`→`host new`。其余需对照 CLI 源码确认；确认不存在的不迁移。

### F-05 参数覆盖缺口（中）

`appendix/2` §2.1.1–2.1.3 有 7 个标识符未出现在 `app-yaml.md`：

```text
always_running  reserved_on_exit  tmpfs_workdir
disable_local_mapping  disable_data_mapping  slot_on_head
slot_max_tasks
```

前 6 个是 `slot_options` 的子选项，`app-yaml.md` 的 `slot_options` 行已描述，属于检查器粒度问题；6.3 决定是否在检查器中跳过。
`slot_max_tasks` 未出现在 go-scalebox `app-schema.json`，6.3 决定补 schema（推荐）或标注为 Sphinx 唯一信源。

### F-06 退出码无缺口（低）

`appendix/6` 的 256 个码全部被 `exit-codes.md` 覆盖；6.3 只删除重复的 §6.11 状态码表，保留应用退出码规范与外部参考表。

### F-07 根 README 的英文设计文档链接（中）

`README.md` 的 5 个链接指向 `docs/vtask-design.md`、`docs/webui-vscode-plan.md`、`docs/rest-api-reference.md`、`docs/grpc-api.md`、`docs/security.md`，这些文件不在本仓库（属于 go-scalebox）；`examples/app-copy/` 目录也不存在。6.5 改为文字指针或删除。

## 事实-信源表

| 事实 | 唯一信源 |
|---|---|
| CLI 命令、usage、flags | `docs/reference/cli.md`（生成） |
| app/module schema 字段 | `docs/reference/app-yaml.md`（生成） |
| 平台状态码、退出码常量 | `docs/reference/exit-codes.md`（生成） |
| 应用退出码语义与分区 | Sphinx `appendix/6` |
| 参数语义细节（`retry_rules`、`slot_options` 等） | Sphinx `appendix/2` |
| task-headers、cluster/host/slot 参数、环境变量 | Sphinx `appendix/2`（本阶段保留） |
| 接口文件格式、sink 7 种格式 | `features/agent-interface-files` |
| 特性操作步骤、配置、坑 | `features/*/README.md` |
| 卡链诊断 | `docs/agents/20-diagnosis.md` |
| 平台不变式 | `docs/agents/10-invariants.md` |
| 设计原理、架构、教程 | Sphinx `developer/*`、`user/*` |

## 批次与提交拆分

| 批次 | 内容 | 建议 commit |
|---|---|---|
| 6.1 | `check-links.py`、`check-content.py`、本记录 | `docs: add the Sphinx link/content checks and the phase-6 baseline` |
| 6.2 | `appendix/3`、`developer/2`、`developer/9 §9.6` | 三个 commit（每文件一个） |
| 6.3 | `appendix/2`、`appendix/6`、`appendix/1`；go-scalebox 生成器 | 附录各一个；生成器 + 生成物单独一个 |
| 6.4 | `developer/7`、`user/3`、`user/5`、`faq`、`user/2` vs `developer/1` | 每个主题一个 |
| 6.5 | `features/*/README.md` 反向指针；doc 角色转换；根 `README.md`；`index.rst` 外链 | 交叉指针、断链、README、index 分开 |
| 6.6 | `started/3_first_app.md`；`index.rst` 目录重排 | 教程一个；导航一个 |

每批先列“大块删除确认清单”（文件、节、行数、目标位置、保留摘要），确认后执行。

## 待办与风险

- `origin/dev` 当前停在 `f8008b5`，不包含阶段 3–5；GitHub `dev` 外链需你把本地 `dev` 推送后才是最新内容。
- 英文站仍保留同样的 doc 角色问题与 6 个断链，留到 6b。
- 构建使用 `/tmp` 副本；临时 venv 与构建产物不入库。
- `check-content.py` 默认报告不阻塞；`--strict` 可用于验收。
