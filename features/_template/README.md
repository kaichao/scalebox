# <feature 名称>

<一句话说明该特性是什么、解决什么问题。>

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `<field>` | app.parameters / module.parameters / module.arguments / environments | ... | ... | ... |

## 行为与不变式

- <必须成立的行为；违反会怎样。>
- <边界条件：只对某模式生效、只在某阶段生效、依赖某前置条件。>

## 坑

- **<现象>**：<原因>；<规避方式>。

## 最小例子

见 `feature.yaml`。关键片段：

```yaml
...
```

## 测试

```sh
./test.sh
```

预期：PASS；断言点：<断言什么>。

## 源码指针

- `go-scalebox/<path>`（<符号/函数>）—— <说明>
- `scalebox/<path>` —— <说明>

## 相关文档

- Sphinx：`docs/cn/source/...`
- 设计文档：`go-scalebox/docs/...`
- 其他：...
