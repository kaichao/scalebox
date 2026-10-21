# reference（参考层）

本目录分两类：

- **生成物**：`cli.md`、`app-yaml.md`、`exit-codes.md`、`rpc.md`、`env-vars.md`——不要手工编辑，
  改动会被下一次生成覆盖，`make docs-check` 也会报差异。
- **手写/半自动**：`task-headers.md`、`module.md`、`semaphore-variable.md`——结构以源码为准，
  由人维护；修改时对照源码并注明依据。

结构事实（命令、字段、状态码、header、环境变量、RPC）以本目录为唯一信源；语义说明
（何时用、副作用、坑）在 [`features/`](../../features/) 与 Sphinx。

| 文件 | 来源 | 方式 |
|---|---|---|
| `cli.md` | `scalebox <cmd> --help`（Cobra 命令树） | 生成：递归遍历 CLI help |
| `app-yaml.md` | go-scalebox `cmd/webui/frontend/src/schema/app-schema.json` | 生成：JSON Schema 字段表 |
| `exit-codes.md` | go-scalebox `internal/task/status_code.go` + 运行时约定 | 生成：解析 Go 常量，补充 123/124 |
| `rpc.md` | go-scalebox `proto/*.proto` | 生成：解析 service/message 定义 |
| `env-vars.md` | go-scalebox 源码中的环境变量读取点 | 生成：扫 `os.Getenv` / 常量 |
| `task-headers.md` | 源码（`cmd/agent`、`cmd/controld`）+ 实测 | 手写 |
| `module.md` | `app-schema.json` modules 段 + agent 接口 | 半自动：字段表生成，语义手写 |
| `semaphore-variable.md` | `pkg/semaphore`、`pkg/variable` + `docs/vtask-design.md` | 手写 |

生成头部的 `source: go-scalebox@<sha> (<date>)` 是生成时的源码版本；源码变化后重新生成，diff 会显示变化。

## 重新生成

在 **go-scalebox 仓库根目录**运行：

```sh
go build -o tmp/scalebox-docs ./cmd/cli
python3 scripts/gen-reference/gen_reference.py \
  --out-dir ../../workspace/scalebox/docs/reference
```

如果本地 Makefile 含 docs 目标（注意该仓库的 `Makefile` 被 `.gitignore` 忽略）：

```sh
make docs-gen     # 更新本目录三个文件
make docs-check   # 生成到临时目录并 diff；源码漂移时失败
```

输出目录可用 `--out-dir` / `DOCS_DIR` 覆盖。生成脚本说明见 go-scalebox 的 `scripts/gen-reference/README.md`。
