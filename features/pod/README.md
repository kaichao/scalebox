# pod

基于 `pod_id` 的节点本地计算：同一 pod 内、同属 HOST-BOUND 的模块之间，任务自动留在同一节点，避免跨节点搬运数据。

## 配置

| 字段 | 位置 | 取值 | 默认值 | 说明 |
|---|---|---|---|---|
| `pod_id` | `module.parameters` | 字符串 | 空 | 同一 pod 的模块使用相同值 |
| `task_dist_mode` | `module.parameters` | `HOST-BOUND` | `DEFAULT` | 该特性只对 HOST-BOUND 生效 |
| `slots` | `module.slots` | slot 正则列表 | — | 模块部署的节点；匹配 hostname 第一段（如 `h0`），不要写 FQDN |

## 行为与不变式

- 当来源模块与目标模块都是 HOST-BOUND 且 `pod_id` 相同，服务端在创建任务时自动把 `from_ip` 解析为 `to_host`，任务路由回同一节点。
- 主路由中转时，用 `last_from_module` / `last_from_ip` 定位真正的前序算法模块（见 `sink_tasks.go`），避免把主路由所在节点误当来源。
- `pod_id` 只影响路由，不改变 slot 数量与资源分配；节点本地计算的其他部分（本地根、数据映射）由 `slot_options` / 部署参数决定。

## 坑

- **只对 HOST-BOUND 生效**：SLOT-BOUND / DEFAULT 模块不会因 `pod_id` 改变路由。
- **两侧 `pod_id` 必须相同**：来源模块没设或值不同，不会自动推导 `to_host`。
- **主路由中转**：跨主路由时 `from_module` 是主路由，必须依赖 `last_from_*`；不要把主路由的 `pod_id` 设成与算法模块相同。
- **`slots` 写 FQDN 不建 slot**：平台用 `split_part(hostname,'.',1)` 取短名匹配；`h0.local` 匹配不到，写 `h0`。
- 节点本地计算还依赖本地根与数据映射配置；只配 `pod_id` 不等于完整的本地计算模式。

## 最小例子

见 `feature.yaml`：

```yaml
parameters:
  main_router: router
modules:
  mod-1:
    parameters:
      task_dist_mode: HOST-BOUND
      pod_id: my_pod
    slots:
      - n0-[01]
  mod-2:
    parameters:
      task_dist_mode: HOST-BOUND
      pod_id: my_pod
    slots:
      - n0-[01]
  router:
    arguments:
      code_path: ${PWD}/router-code
      slot_options: slot_on_head
```

## 测试

```sh
./test.sh   # 需要 inline 多节点控制面；不可达时 SKIP（退出码 77）
# 当前手工：
echo 4 | scalebox run --app-file=feature.yaml
```

预期：`mod-1` 产生的下游 `mod-2` 任务 `to_host` 与 `mod-1` 执行节点相同；跨节点数量为 0。

## 源码指针

- `go-scalebox/cmd/controld/task/add_local.go`（`pod_id processing`）—— 自动推导 `to_host`
- `go-scalebox/cmd/agent/sink_tasks.go`（`addSinkTasks`）—— `last_from_module` / `last_from_ip` 继承
- `go-scalebox/cmd/controld/task/cache_attrset.go`—— `pod_id` 参数读取

## 相关文档

- Sphinx：`docs/cn/source/developer/5_node_local_compute.md` §5.4–5.5
- 参数表：`docs/reference/app-yaml.md`（`module.parameters` 的 `pod_id`）
- 实战：app-fxcorr `docs/local-compute-plan.md`
