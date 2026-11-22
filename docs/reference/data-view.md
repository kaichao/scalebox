# Data View 参考

> 唯一信源：Data View 的定义、App 绑定、SQL 模板与 renderer 配置。
> 设计与安全边界见 go-scalebox `docs/data-view-design.md` 与
> `docs/app-custom-view-design.md`；使用入口见
> [`docs/cn/source/application/7_webui.md`](../cn/source/application/7_webui.md)。
> 最后校验：2026-11-22，对照 go-scalebox@14360f2。

## 1. 定位

Data View 是 WebUI 的**只读 SQL 查看组件**：管理员预定义 SQL 模板，绑定到 App 后以
表格或应用定制可视化展示查询结果。

- 不是通用 SQL 控制台：只允许单条 `SELECT` / `WITH`。
- 不是 BI 平台：不提供拖拽式图表编辑器；新可视化类型需要前端新增 renderer 组件。
- 只在**当前集群 PostgreSQL** 执行；跨集群应用需在对应集群的 WebUI 查看。
- 定义与绑定是本地表，不参与跨集群复制。

## 2. 数据模型

| 表 | 作用 | 关键字段 |
|---|---|---|
| `t_data_view` | 组件定义 | `id`、`name`（唯一 slug）、`title`、`description`、`sql_template`、`params`（JSONB）、`layout`（JSONB）、`max_rows`、`status` |
| `t_app_data_view` | App 绑定 | `id`（binding id）、`app_id`、`data_view_id`、`title`、`bindings`（JSONB）、`tab_order`、`visible` |

- 同一 Data View 可以被同一 App 绑定多次（不同参数绑定），REST 用 **binding id** 寻址。
- `max_rows` 缺省 500，上限 5000；结果行数、字节数、单元格大小都有限制。
- `layout` 同时承载表格列配置和 Overview 定制视图配置；未知字段忽略，未知 renderer 回退表格。

## 3. 参数与绑定

### 3.1 `params`

```json
[
  {"name":"module","label":"Module","type":"module","required":false,"default":""},
  {"name":"since", "label":"Since", "type":"string","default":"1 hour"},
  {"name":"limit", "label":"Limit", "type":"int",   "default":"100"}
]
```

`type` 可为 `string` / `int` / `bool` / `select` / `module` / `slot` / `task` / `vtask`。
`module` 等语义类型只用于标识参数含义；查看页不渲染参数表单，参数值来自绑定或默认值。

### 3.2 `bindings`

```json
{
  "app_id": {"source":"context",  "value":"app.id"},
  "user":   {"source":"constant", "value":"kaichao"},
  "since":  {"source":"runtime",  "default":"1 hour"}
}
```

| `source` | 值来源 | 客户端可覆盖 |
|---|---|---|
| `context` | 服务端解析 `app.id` / `app.name` / `app.cluster` / `user.id` | 否 |
| `constant` | 固定值 | 否 |
| `runtime` | 客户端传入 → binding default → params default | 是 |

`context` / `constant` 参数由服务端解析并忽略客户端传值，避免 viewer 改 `app_id` 读其他 App 数据。

## 4. SQL 模板与安全

- 占位符 `{{name}}` 编译为 `$n` 参数绑定，不做文本拼接。
- 占位符必须独立作为值出现；`'{{x}}'`、字符串/注释/dollar-quote 内的占位符直接拒绝。
- 首关键字必须是 `SELECT` 或 `WITH`；中间不得有多余分号；禁止 `$数字`。
- 表名黑名单：`t_user`、`t_api_key`、`t_role_binding`、`t_token_blacklist`
  （`DATA_VIEW_TABLE_DENYLIST` 可覆盖）。这是 best-effort，不是完整安全边界。
- 执行使用独立连接池 + `BEGIN READ ONLY` + `statement_timeout`（缺省 30 秒）；
  结果按行数 / 总字节 / 单元格大小截断，超时返回已取到的部分结果。
- `DATA_VIEW_STATEMENT_TIMEOUT_MS`、`DATA_VIEW_MAX_RESULT_BYTES`、
  `DATA_VIEW_CELL_MAX_BYTES` 等环境变量见 [`env-vars.md`](./env-vars.md)。

## 5. REST API

前缀 `/api/v1`。

| 方法 | 路径 | 说明 | 权限 |
|---|---|---|---|
| `GET` | `/data-views` | 列出定义，可按 `status` 过滤 | read |
| `POST` | `/data-views` | 创建定义 | admin |
| `GET` | `/data-views/{id}` | 读取定义 | read |
| `PUT` | `/data-views/{id}` | 更新定义 | admin |
| `DELETE` | `/data-views/{id}` | 删除定义 | admin |
| `POST` | `/data-views/preview` | 编辑器试运行（不落绑定） | admin |
| `GET` | `/apps/{id}/data-views` | 列出 App 绑定 | read |
| `POST` | `/apps/{id}/data-views` | 绑定 Data View | admin |
| `PUT` | `/apps/{id}/data-views/{bindingId}` | 更新绑定 | admin |
| `DELETE` | `/apps/{id}/data-views/{bindingId}` | 解绑 | admin |
| `POST` | `/apps/{id}/data-views/{bindingId}/query` | 执行绑定查询 | read |

结果形状：

```json
{
  "columns": [{"name":"task_id","type":"int8"}],
  "rowsJson": ["[1,\"ok\"]"],
  "rowCount": 1,
  "truncated": false,
  "elapsedMs": 12
}
```

`ResourceExhausted` 映射 HTTP 429；参数/模板错误映射 400。

## 6. App 页面与 Overview 定制视图

### 6.1 App → Data

App 页的 `Data` tab 列出该 App 的 Data View 绑定。查看页不渲染参数表单；需要不同参数时，
修改绑定默认值，或为同一视图创建多个绑定（如“最近 1 小时”/“最近 24 小时”）。
`Refresh` 只重新执行当前 tab 的查询。

### 6.2 `layout` 与 renderer

`t_data_view.layout` 是 JSON object；缺省/无效时按表格展示。

```json
{
  "version": 1,
  "overview": true,
  "renderer": {
    "type": "scatter3d",
    "x": "channel",
    "y": "time_sec",
    "z": "pointing",
    "color": "value",
    "axisTypes": {"x": "value", "y": "value", "z": "category"},
    "labels": {"x": "Channel", "y": "Time", "z": "Pointing", "color": "SNR"}
  }
}
```

| 字段 | 说明 |
|---|---|
| `overview` | `true` 时挂到 App Overview；不是 bool 按 false |
| `renderer.type` | `table`（缺省）/ `scatter3d`；未知值回退表格并提示 |
| `renderer.x/y/z` | `scatter3d` 必填，引用 SQL 结果列名 |
| `renderer.color` / `colorMap` / `colorLabels` | 数值着色或离散色块 |
| `renderer.axisTypes` / `axisRange` / `axisLabelIntervals` | 轴类型、固定范围、标签间隔；value 轴自动整数格式化 |
| `renderer.companion` | 3D 下方二维 heatmap；`aggregate=avg|count|min`，支持 `filter` / `views[]` |
| `renderer.height` / `pointSize` / `opacity` | 200–800 / 1–20 / 0.1–1 |

- SQL 列名是数据契约；renderer 只引用结果列名。
- 时间轴建议在 SQL 中转为相对整数秒或 epoch 数值。
- 缺列、WebGL 不可用、未知 renderer 都降级为错误卡片或表格，不拖垮整个 Overview。
- `layout.overview` 是定义级开关；同一视图绑定到多个 App 时会在所有绑定 App 的 Overview 出现。

## 7. 限制与排障

- Data View 只读；不能写库、不能调用外部服务。
- 表名黑名单是 best-effort；不要把它当唯一安全边界，敏感表应通过数据库权限另行限制。
- 全量大数据集必须在 SQL 内聚合/分桶；Data View 有行数和字节上限。
- 查询失败先看：模板是否单条 SELECT/WITH、占位符是否都在 `params` 中、绑定参数是否缺失、
  是否触发敏感表名单、是否超时或超过行数上限。
- MCP server 当前不暴露 Data View；需要定义/执行时用 WebUI 或 REST。

## 8. 源码指针

- `build/database/initdb.d/84-data-view.sql`、`87-data-view-link-aliases.sql`
- `internal/dataview/`：模板编译、参数解析、只读执行、SQLStore
- `cmd/controld/grpc/data_view_service.go`：DataViewService
- `cmd/webui/handler/data_view.go`、`cmd/webui/router/router.go`：REST
- `cmd/webui/frontend/src/pages/dataViews/`、`pages/apps/AppDataViews.tsx`、
  `components/dataViewRenderers/`：前端
- go-scalebox `docs/data-view-design.md`、`docs/app-custom-view-design.md`
