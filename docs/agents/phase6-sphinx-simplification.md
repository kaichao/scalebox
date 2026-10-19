# 阶段 6 — Sphinx 简化执行记录

日期：2026-10-18。范围：`docs/cn/source/`（中文站；英文站留到 6b）。
目标：Sphinx 只保留原理、设计、教程、规范；操作细节指向 `features/`、`docs/reference/`、`docs/agents/`；无重复内容，链接可用。

## 拍板

- 外链基底：GitHub `dev`（`https://github.com/kaichao/scalebox/blob/dev/...`）。
- 中英范围：先中文；英文 6b。
- `appendix/2` 的 task-headers、cluster/host/slot 参数、环境变量保留在 Sphinx。
- 构建验收：本机容器 `cd docs && make build-cn`（镜像 `kaichao/read-the-docs`）；`ssh fxcorr` venv 作为备用。

## 基线环境与命令

本机：Python 3.14.2（仅运行标准库脚本）；Docker 可用。

推荐构建：

```sh
cd docs && make build-cn
```

`docs/Makefile` 使用 `${PWD}` 挂载；`make -C docs build-cn` 会把当前目录挂到 `/docs`，导致 `make html` 找不到 `docs/cn/Makefile`。必须在 `docs/` 目录内执行。构建产物在 `docs/cn/build/`，已被 `docs/.gitignore` 忽略。

备用构建（`ssh fxcorr`，Python 3.9.25 + venv + Sphinx 6.2.1）：

```sh
rsync -a --delete --exclude build docs/ fxcorr:/tmp/scalebox-docs-phase6/
ssh fxcorr 'cd /tmp/scalebox-docs-phase6 && python3 -m venv .venv && \
  .venv/bin/python -m pip install -q -r requirements.txt'
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
| CLI 覆盖 | `appendix/3` 提取 78 条；`reference/cli.md` 未覆盖 23 条（P-1 修复后 18 条） |
| 参数覆盖 | `appendix/2` §2.1.1–2.1.3 提取 50 个；`app-yaml.md` 未覆盖 7 个 |
| 退出码覆盖 | 256 个；`exit-codes.md` 未覆盖 0 个 |
| Sphinx 构建 | exit 0；1 warning |
| `sphinx-build -W` | exit 2（warning 视为错误） |

构建 warning：

```text
developer/7_advanced_features.md:323: duplicate label ***应用定义***,
other instance in developer/7_advanced_features.md
```

## 6.2 进度

- `appendix/3`（2026-10-19）：997 → 203 行；CLI 提取 26 条，`reference/cli.md` 0 缺口；容器干净构建仅剩 `developer/7` 的已知重复标签 warning；6 个 go-scalebox 断链与 19 个 doc 角色按计划留到 6.5。
- `developer/2`（2026-10-19）：262 → 180 行；§2.3.2 改为分类表 + 指针；`iobytes` / `network_io` 字段与视图、时间戳精确格式、`cleanup-files` 未实现迁移到 `features/agent-interface-files`（88 → 96 行）。
- `developer/9`（2026-10-19）：488 → 460 行；§9.6 从 106 行压到约 79 行；保留三层操作体系、fail/bind/unbind 语义与 add-subtask 三条路径，命令细节改指 `reference/cli.md`，并更新 `list`→`ls`、`diffmin/diffmax`→`diff-min/diff-max`。

## 6.3 进度

- `appendix/2`（2026-10-19）：170 → 116 行；2.1.1–2.1.3 的字段表改为 `app-yaml.md` 指针；保留 2.1.4 task-headers、2.2 cluster/host/slot、2.3 环境变量；`slot_max_tasks` 作为 schema 未覆盖的唯一补充。
- `appendix/1`（2026-10-19）：206 → 141 行；app/module schema 字段说明改指针；保留 slots 语法、命名规则、SSH 性能表；修正重复的 1.2 编号。
- `appendix/6`（2026-10-19）：233 → 209 行；6.11 改为 `reference/exit-codes.md` 指针；删除过期的 224/225 子任务码表；保留应用退出码规范与外部参考。
- `gen_exit_codes.py` 去掉对 Sphinx 的反向引用；`exit-codes.md` 重新生成；`make docs-check` OK。
- `check-content.py` 参数段改为只检查 §2.1.3 补充表；当前报告 `slot_max_tasks` 未出现在 `app-yaml.md`，这是预期缺口。

## 6.4 进度

- `developer/7`（2026-10-19）：467 → 370 行；retry/timestamps/slot 示例改指 features；vtask 计数改指 `developer/9 §9.4`；修正重复编号与重复 section label。干净 `sphinx-build -W` 现在通过（0 warning）。
- `user/3`（2026-10-19）：166 → 58 行；重写为“建应用 → 投任务 → 查状态 → 清理”，删除 7 个不存在/过时命令。
- `user/5`（2026-10-19）：228 → 217 行；删除 `host metrics`，任务 READY 排障改指诊断文档，5.8 指 FAQ。
- `faq`（2026-10-19）：195 → 174 行；任务/信号量排障改指针，`task log` 改 `task show`，`build/compose.yaml` 改 `runtime/compose.yaml`。
- `developer/1 §1.3` 已是指向 `user/2` 的指针，本轮不改；`:doc:` 统一转换留到 6.5。
- 当前链接检查：6 个 go-scalebox 缺失 + 17 个 Markdown doc 角色（FAQ 已减少 2 个）。

## 6.5 进度

- 17 个 Markdown doc 角色全部转换：11 个内部目标改为 MyST `{doc}`，6 个 go-scalebox 目标改为纯文字指针；`docs/cn/source` 中已无该角色语法。
- `check-links.py --scope all`：54 个文件，0 缺失、0 个 `SYNTAX`。
- 干净 `sphinx-build -W`：0 warning，exit 0；`make -C docs build-cn` 也通过。
- features 参数指针更新 7 个文件（timeout、node-progress-gap、task-batch-size、retry-test、pod、task-identity、check-test），改指 `reference/app-yaml.md` 或保留 §2.1.4。
- 根 README：`build/` 改 `runtime/`，删除不存在的 `gen-secrets.sh` 与 `examples/app-copy`，英文设计文档改为 go-scalebox 文字指针。
- `index.rst` 增加 AGENTS/features/reference 外链入口。
- `docs/Makefile` 的 `${PWD}` 改为 `$(CURDIR)`；`make -C docs build-cn` 可用。
- 修正 `developer/7` 的 `space_free_gb` → `free_space_gb`。

## 6.6 进度

- `index.rst` 重排为 6 个分区：入门教程 / 应用开发 / 部署与运维 / 设计与原理 / 参考 / 附录与其他；文件路径不变。
- 新增 `started/3_first_app.md`（110 行）：app.yaml + run.sh + 投递任务 + 查看 stdout + 清理，全程不需要构建镜像。
- 7 个页面增加标题/边界声明：`appendix/1`、`appendix/2`、`appendix/6`、`developer/2`、`developer/7`、`developer/9`、`developer/4`。
- 中文站最终验证：`check-links.py --scope all` 55 文件、0 缺失、0 doc 角色语法；干净 `sphinx-build -W` 0 warning；`make -C docs build-cn` 通过；`check-content.py` CLI 26/0、参数补充 1 条预期、退出码 256/0；`features/check-index.sh` OK。
- 中文 Sphinx 正文从 5838 行降到约 4677 行，重复的操作细节已迁到 `features/` 与 `docs/reference/`。
- 英文站 6b 待做：同样的 doc 角色转换、断链修复与内容简化。

## 整体审计（2026-10-19）

- 无文件删除、无重命名；暂存区 35 个文件，另有 4 个审计补丁未暂存。
- 关键内容去向：
  - CLI 命令/flags → `docs/reference/cli.md`；app/module 字段 → `app-yaml.md`；平台状态码 → `exit-codes.md`。
  - 接口文件格式、I/O 落库字段、sink 7 种格式 → `features/agent-interface-files`。
  - retry / slot 表达式 / 准入控制示例 → `features/retry-test`、`slot-expr`、`check-test` 等。
  - 卡链与排障 → `docs/agents/20-diagnosis.md`；vtask 语义 → `developer/9 §9.4/§9.6`。
- 审计补丁：`appendix/2 §2.2.3` 补 slot 参数空值/字符串语义；`developer/7 §7.5` 补 bash/Go 时间戳示例；`user/3 §3.1` 补环境变量/模板变量优先级指针；`features/slot-expr` 补多 GPU command 模板。
- 最终验证：`check-links.py --scope all` 55 文件 0 缺失、0 doc 角色；干净 `sphinx-build -W` 0 warning；`make -C docs build-cn` 通过；`check-content.py` CLI 26/0、参数补充 1 条预期、退出码 256/0；`features/check-index.sh` OK。

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

### F-08 CLI 生成器漏掉长命令名（高）

`gen_cli.py` 解析 `Available Commands:` 时用 `^\s{2}(\S+)\s{2,}(.*)$`；Cobra 对长命令名（`get-parameter`、`set-parameter`、`remove-header`、`add-semaphore-value`）只留一个空格，这些命令因此未进入 `cli.md`。`make docs-check` 仍然 OK，因为生成器与生成物同源，无法发现自身漏项。

处理（已完成 2026-10-19）：`scripts/gen-reference/gen_cli.py` 的分隔正则改为 `\s+`；重新生成 `cli.md`，新增 115 行、命令总数 92→97；`make docs-check` OK。这是 go-scalebox 的独立 commit。

### F-09 appendix/3 命令分类（高）

- 生成器漏项（F-08 已修复，reference 已覆盖）：`cluster get-parameter`、`slot get-parameter`、`slot set-parameter`、`task remove-header`、`vtask add-semaphore-value`。
- 改名：`global list`→`global ls`、`semaphore list`→`semaphore ls`、`slot remove`→`slot delete`、`variable list`→`variable ls`、`host create`→`host new`。
- 当前源码已无：`app create`、`app add-remote`、`app set-finished`、`cluster show`、`channel pull/push`、`event *`、`fs *`、`task log`、`scalebox status`；从 appendix/3 删除，不迁移。
- `channel` 的跨集群通信语义由 `scalebox queue` 取代（reference 已有 queue 章节）。

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
- 容器构建直接写 `docs/cn/build/`（已被 `docs/.gitignore` 忽略）；fxcorr 备用构建使用 `/tmp` 副本。
- P-1 已完成：`gen_cli.py` 分隔正则修复 + `cli.md` 重新生成；剩余 18 条 CLI 缺口均为改名或已移除，见 F-09。
- `cli.md` 头部 source 仍为 `go-scalebox@d343cc1`；提交 `gen_cli.py` 后需再跑一次 `make docs-gen` 更新 stamp，再提交 `cli.md`。
- `docs/Makefile` 的 `${PWD}` 使 `make -C docs` 失效；6.6 可改为 `$(CURDIR)`。
- `check-content.py` 默认报告不阻塞；`--strict` 可用于验收。
