# 3. 命令行参考（摘要）

> 本章只保留命令行的使用约定与语义补充。完整的命令、usage 与 flags 见
> [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)（由 `scalebox --help` 生成）。

## 1.1 命令行选项与环境变量

| 选项 | 缺省值 | 描述 |
|---|---|---|
| `-e, --env-file` | `scalebox.env` | 环境变量文件，设置命令运行的环境变量 |

环境变量是 Scalebox 应用程序中参数传递的重要方法。应用中环境变量定义可来自于多个环境变量定义文件、系统级环境变量；若在不同定义文件、系统级变量中存在重复的变量名，则按照以下顺序加载（若文件不存在，则忽略）：

- 系统级环境变量
- 用户自定义名的 env 文件
- 当前目录下 `scalebox.env` 文件
- `${HOME}/.scalebox/environments`
- `/etc/scalebox/environments`

其中，用户自定义名的 env 文件可按文件名执行级联加载。示例如下：用户自定义 env 文件名为 `p419_48nodes_1266932744.env`，则按优先级从高到低依次加载：

- `p419_48nodes_1266932744.env`
- `p419_48nodes.env`
- `p419.env`

## 1.2 常用任务示例

### 启动应用

```bash
# 用 app.yaml 创建并启动应用
scalebox run --app-file app.yaml --env-file scalebox.env

# 命令行模式：指定镜像与起始 body
scalebox run --app-cluster local --image-name my-image:latest "start-item"

# 只打印生成的 YAML，不创建应用
scalebox run --app-file app.yaml --print-yaml
```

### 查看状态

```bash
scalebox app list
scalebox module list --app-id N
scalebox slot list --module-id M
scalebox task list --app-id N
scalebox task show T
```

### 投递任务

```bash
# 从 stdin 读取 body，指定下游模块与 header
echo "body" | scalebox task add --app-id N --sink-module consumer -H _demo=hello

# 从任务文件读取，每行一条
scalebox task add --app-id N --task-file tasks.txt
```

### 清理

```bash
scalebox app delete --app-id N
```

`app delete` 只删除数据库记录，不停止本地 slot 容器；清理方式见
[卡链与失败诊断](https://github.com/kaichao/scalebox/blob/dev/docs/agents/20-diagnosis.md)。

### 跨集群队列

`scalebox queue` 提供跨集群应用通信：`ls`、`peek`、`pull`、`push`、`purge`。
参数与优先级规则见 CLI 参考的 `scalebox queue` 章节。

## 1.3 语义补充

### 1.3.1 `scalebox run` 参数

| 参数 | 说明 |
|---|---|
| `--app-name` | 应用名称；缺省为 `<目录名>_<时间戳>` |
| `--app-cluster` | 主应用所在集群；缺省 `local` |
| `--app-file` | 应用定义文件路径 |
| `--image-name` | 主模块镜像 |
| `--code-path` | 主模块代码目录；存在 `./code` 时缺省使用 |
| `--slot-regex` | 主模块 slot 匹配（主机短名或正则） |
| `--mr-image-name` | 主路由模块镜像 |
| `--mr-code-path` | 主路由代码目录；存在 `./mr-code` 时缺省使用 |
| `--main-app-id` | 跨集群链接的主应用 ID |
| `--main-cluster` | 跨集群链接的主应用集群 |
| `--print-yaml` | 只打印生成的 YAML，不创建应用 |
| `-e, --env-file` | 环境变量文件 |

CLI 内部使用 `MY_APP_NAME`、`MY_APP_CLUSTER`、`MY_IMAGE_NAME`、`MY_CODE_PATH`、
`MY_SLOT_REGEX`、`MY_MR_IMAGE_NAME`、`MY_MR_CODE_PATH` 等环境变量；一般直接使用上述 flags。

### 1.3.2 `task add` 输入格式与控制 header

主要 flags：

| flag | 说明 |
|---|---|
| `--app-id` / `--module-id` | 应用 ID / 模块 ID |
| `--sink-module` | 下游模块名 |
| `-H, --header` | 单个 header，`key=value`，可重复 |
| `--headers` | JSON 格式的多个 header |
| `--conflict-action` | 冲突处理：`''` / `IGNORE` / `OVERWRITE` |
| `--task-file` | 任务文件路径；每行一条 |
| `--from-module` | 来源模块名 |
| `--remote-cluster` | 投递到远端集群 |
| `--direct-write-in-agent` | 跳过 agent 内逻辑 |

不带 `--task-file` 时从 stdin 读取。任务文件每行格式：

| # | 格式 | 示例 |
|---|---|---|
| 1 | 文本 body | `body` |
| 2 | JSON body | `{"hi0":"a","body":"my_body"}` |
| 3 | 文本 body + headers | `body,{"h0":"a","h1":"b"}` |
| 4 | JSON body + headers | `{"hi0":"a","body":"my_body"},{"h0":"a","h1":"b"}` |
| 5 | 模块名 + 文本 body | `module-name,body` |
| 6 | 模块名 + 文本 body + headers | `module-name,body,{"h0":"a","h1":"b"}` |
| 7 | 模块名 + JSON body + headers | `module-name,{"hi0":"a","body":"my_body"},{"h0":"a","h1":"b"}` |

控制 header：

| header | 说明 |
|---|---|
| `slot_broadcast=yes` | 广播到该模块的所有 slot |
| `host_broadcast=yes` | 广播到该模块的所有 host |
| `conflict-action` | 冲突处理策略，与同名 flag 等价 |

`_` 前缀的 header 会沿下游模块自动传播；`_vtask_attr_*` 是 vtask 范围内的应用自定义属性。
详见 [headers-propagation](https://github.com/kaichao/scalebox/blob/dev/features/headers-propagation/README.md)。

### 1.3.3 信号量

- 公共参数：`--app-id` 或 `--module-id`（也可用 `APP_ID` / `MODULE_ID` 环境变量）。
- 自动创建：设置 `SEMAPHORE_AUTO_CREATE=yes` 时，信号量不存在则以初值 0 创建；CLI 通过 gRPC metadata `semaphore-auto-create: yes` 传给 controld。
- 命名：字符集 `[A-Za-z0-9:_-]`，首字符为字母或下划线。
- 信号量表达式：在命名规则基础上增加 `.*+?^$[]{}()|\`，用于按正则匹配一组信号量。

常用操作：

```bash
scalebox semaphore ls [prefix] --leaf-only
scalebox semaphore create <name> <value>
scalebox semaphore create '{"semaphores":{"sema1":n1,"sema2":n2}}'
scalebox semaphore create --sema-file my-sema-file.txt --batch-size 100
scalebox semaphore get <name-or-expr>
scalebox semaphore increment <name-or-expr>
scalebox semaphore decrement <name-or-expr>
scalebox semaphore add-value <name> <delta>
scalebox semaphore delete <name>
```

`get` / `increment` / `decrement` / `add-value` 的退出码：`0` 成功，`1` 数据库错误，
`2` 信号量不存在。表达式匹配多个信号量时，输出为 JSON map：
`{"sema1":n1,"sema2":n2}`。

### 1.3.4 信号量组

`semagroup` 以信号量名前缀或正则标识一组信号量：

| 命令 | 语义 |
|---|---|
| `scalebox semagroup max <expr>` | 组内最大值 |
| `scalebox semagroup min <expr>` | 组内最小值 |
| `scalebox semagroup increment <expr>` | 取组内最小值并加一 |
| `scalebox semagroup decrement <expr>` | 取组内最大值并减一 |
| `scalebox semagroup diff-max <expr>` | 组内最大值与当前值的差值；`expr` 形如 `(group-prefix):sema-suffix` |
| `scalebox semagroup diff-min <expr>` | 当前值与组内最小值的差值；`expr` 形如 `(group-prefix):sema-suffix` |

返回值与退出码约定同 `semaphore`。

### 1.3.5 变量

- 公共参数：`--app-id` 或 `--module-id`。
- 命名：同信号量（`[A-Za-z0-9:_-]`，首字符为字母或下划线）。
- 支持 `ls`、`get`、`set`、`delete`。

```bash
scalebox variable ls [prefix] --leaf-only
scalebox variable get <name-or-expr>
scalebox variable set <name> <value>
scalebox variable delete <name>
```

`get` 的退出码：`0` 成功，`1` 数据库错误，`2` 变量不存在。表达式匹配多个变量时，
输出为 JSON map：`{"var1":"val1","var2":"val2"}`。

### 1.3.6 其他约定

- `global` 变量跨应用共享；`variable` 属于应用或模块。两者都用 `ls` / `get` / `set` / `delete`。
- `cluster init-local` 必须在 local 集群头节点的物理机上执行，不能在容器内执行。
- 旧命令名已更新：`global list`→`global ls`、`semaphore list`→`semaphore ls`、`slot remove`→`slot delete`、`variable list`→`variable ls`、`host create`→`host new`。
- 旧 `channel` 已由 `scalebox queue` 取代；系统整体状态请用 `cluster status`、`host status`、`slot status`。

## 1.4 完整命令参考

- [docs/reference/cli.md](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)
- 由 `scalebox <command> --help` 递归生成；当前覆盖 97 个命令。
- 命令、flags、usage 以该页为准；本章不再逐条列出子命令。
