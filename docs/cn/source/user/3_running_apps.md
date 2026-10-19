# 3. 运行与管理应用

本章按“建应用 → 投任务 → 查状态 → 清理”走一遍。完整命令、usage 与 flags 见
[CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)。

## 3.1 创建应用

```bash
# 用 app.yaml 创建并启动
scalebox run --app-file app.yaml --env-file scalebox.env

# 只打印生成的 YAML，不创建应用
scalebox run --app-file app.yaml --print-yaml
```

`run` 的主要参数（`--app-cluster`、`--image-name`、`--code-path`、`--slot-regex` 等）见
[appendix/3 §1.3.1](../appendix/3_commandline_tools.md)。字段定义见
[app-yaml](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)。
环境变量文件与模板变量的优先级见 [appendix/1 §1.2](../appendix/1_tech_specifications.md) 和
[appendix/3 §1.1](../appendix/3_commandline_tools.md)。

## 3.2 投递任务

```bash
# 从 stdin 读取 body，指定下游模块与 header
echo "body" | scalebox task add --app-id N --sink-module consumer -H _demo=hello

# 从任务文件读取，每行一条
scalebox task add --app-id N --task-file tasks.txt
```

任务文件格式、控制 header 与 `--conflict-action` 见
[appendix/3 §1.3.2](../appendix/3_commandline_tools.md)。

## 3.3 查看状态

```bash
scalebox app list
scalebox module list --app-id N
scalebox slot list --module-id M
scalebox task list --app-id N
scalebox task show T
```

`task show` 输出 stdout/stderr、body、headers 与 `status_code`。任务卡住或失败时，按
[卡链与失败诊断](https://github.com/kaichao/scalebox/blob/dev/docs/agents/20-diagnosis.md) 排查。

## 3.4 停止与清理

```bash
scalebox app delete --app-id N
```

`app delete` 只删除数据库记录，不停止本地 slot 容器；容器清理方式见诊断文档的“清理现场”一节。

## 3.5 完整命令参考

- [CLI 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/cli.md)
- [app.yaml 参考](https://github.com/kaichao/scalebox/blob/dev/docs/reference/app-yaml.md)
- [参数语义](../appendix/2_parameter_reference.md)
