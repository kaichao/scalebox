# reference（参考层）

本目录的 `cli.md`、`app-yaml.md`、`exit-codes.md` 是**生成物，不要手工编辑**；手工改动会被下一次生成覆盖，`make docs-check` 也会报差异。结构事实（命令、字段、状态码）由源码生成；语义说明（何时用、副作用、坑）在 [`features/`](../../features/) 与 Sphinx。

| 文件 | 来源 | 生成方式 |
|---|---|---|
| `cli.md` | `scalebox <cmd> --help`（Cobra 命令树） | 递归遍历 CLI help |
| `app-yaml.md` | go-scalebox `cmd/webui/frontend/src/schema/app-schema.json` | JSON Schema 字段表 |
| `exit-codes.md` | go-scalebox `internal/task/status_code.go` + 运行时约定 | 解析 Go 常量，补充 123/124 |

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
