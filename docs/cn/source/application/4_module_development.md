# 4. 模块开发

> 本章讲模块设计；接口文件格式见 [features/agent-interface-files](https://github.com/kaichao/scalebox/blob/dev/features/agent-interface-files/README.md)。

## 4.1 模块设计原则

### 4.1.1 无状态设计
- 任务运行的无状态，可重复运行
- 任务级幂等性
- 灵活调用，通过环境变量、任务头定义不同逻辑处理
- 应用状态由主路由与平台存储承担，模块不持有应用状态；详见 {doc}`核心概念 §1.8 <1_core_concepts>`

### 4.1.2 通用性设计
- 支持任务头驱动
- 任务头驱动的流程控制（控制信息由主路由设置）
- 标准化接口设计
- 主路由执行并行逻辑，模块通过任务头接收控制信息；详见 {doc}`核心概念 §1.3 <1_core_concepts>` 与 {doc}`主路由与应用状态 <../design/2_main_router_status>`

## 4.2 模块设计方法

### 4.2.1 task-body设计
task-body为任务标识，在模块中具有唯一性：
- 简单任务：直接用待处理文件名相对路径
- 复杂任务：用输入路径的规范标识（URI）
- 支持JSON格式的非空字符压缩表示

### 4.2.2 task-headers设计
- 配置参数
- 运行时需调整的配置参数通常以环境变量形式传递
- 在外部编排系统中可以修改配置参数

### 4.2.3 模块内目录设计

#### 代码目录/配置目录
- 通常数据量较小，重点考虑可灵活配置
- 可放置于共享存储中，或直接打包进容器镜像
- 代码运行顺序：
  1. 环境变量指定：ACTION_RUN、ACTION_CHECK、ACTION_SETUP、ACTION_TEARDOWN
  2. /app/bin/{run.sh,check.sh,setup.sh,teardown.sh}
  3. /app/share/bin/{run.sh,check.sh,setup.sh,teardown.sh}

#### 数据目录
在数据处理量较大的计算模块中，需按照计算过程中的使用特性，将目录分类并支持灵活配置：

| 目录类型      | 目录说明                                             |
| ------------ | -------------------------------------------------- |
| 输入数据目录 | 综合考虑读取频度、数据总量，可灵活配置不同数据存储形式 |
| 中间文件目录 | 通常位于本机存储 |
| 输出结果目录 | 综合考虑写出频度、数据总量，可灵活配置不同数据存储形式 |

## 4.3 模块集成实现

### 4.3.1 Sidecar模式
模块采用sidecar模式运行：
- **run**：任务的单次运行
- **check**：任务运行的前置条件检测
- **setup**：设置初始运行环境
- **teardown**：退出前清理计算环境

### 4.3.2 算法运行run.sh

#### agent接口文件

`run.sh` 是模块的 `ACTION_RUN`：读取任务输入，调用算法，并把执行结果写到 `${WORK_DIR}` 下的接口文件。agent 在任务结束后采集这些文件并落库。

| 类别 | 文件 |
|---|---|
| 执行结果与元数据 | `task-exec.yaml`、`extra-attrs.yaml`、`auxout.txt` |
| 时间戳 | `timestamps.txt` |
| I/O 计量 | `input-files.txt`、`output-files.txt`、`network-files.txt`、`removed-files.txt` |
| 下游任务 | `sink-tasks.txt` |
| 未实现 | `cleanup-files.txt` |

完整字段、7 种 `sink-tasks.txt` 格式、时间戳格式、I/O 落库字段与视图、坑和最小例子，统一见
[agent-interface-files](https://github.com/kaichao/scalebox/blob/dev/features/agent-interface-files/README.md)。

## 4.4 模块的镜像封装

### 4.4.1 基于agent构建算法模块
```dockerfile
FROM scalebox.net/platform/agent:latest
```

### 4.4.2 拷贝agent组件至算法模块
```dockerfile
COPY --from=scalebox.net/platform/agent:latest /usr/local/ /usr/local/
```

## 4.5 模块单元测试

### 4.5.1 独立模块测试
缺省的模块代码目录为：`./code`

```bash
# 单独命令行启动模块测试
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path}

# 先创建模块测试，再添加测试任务
app_id=$(scalebox run --image-name ${module_image_name} --code-path ${code_path}| cut -d':' -f2 | tr -d '}')
scalebox task add --app-id=${app_id} --header header1=${header1} ${task_body}
```

### 4.5.2 带路由的单元测试
缺省的路由代码目录为：`./mr-code`

```bash
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path} --mr-image-name ${mr_image_name} --mr-code-path ${mr_code_path}
```

### 4.5.3 基于流水线应用的单元测试
- 针对较复杂的单元测试，可写独立的应用定义文件app.yaml、环境变量定义文件scalebox.env

## 4.6 模块类型详解

### 4.6.1 基础模块（agent）
- 提供了模块底座功能
- 应用模块的中任务执行的生命周期管理
- 与server端runtime的controld交互，纪录任务执行状态

### 4.6.2 文件传输模块
- **file-copy**：基于rsync-over-ssh的文件传输
- **dir-copy**：目录拷贝模块
- **rsync-copy**：rsync拷贝模块
- **ftp-copy**：FTP协议传输模块

### 4.6.3 辅助功能模块
- **dir-list**：目录列表模块
- **cron**：定时任务模块
- **cluster-head**：集群头节点模块
- **node-agent**：节点代理模块

## 4.7 错误处理与容错机制

### 4.7.1 任务级容错
- 按单个task返回码，安排自动重试
- 支持任务重试机制

### 4.7.2 组合式容错
- 解决复杂的非逻辑错
- 以任务级容错为基础
- slot级出错处理

## 4.8 最佳实践

### 4.8.1 模块开发规范
1. 保持模块单一职责
2. 支持环境变量配置
3. 实现任务幂等性
4. 合理设置超时时间

### 4.8.2 性能优化
1. 优化数据局部性
2. 减少不必要的文件操作
3. 合理使用缓存
4. 批量处理优化

### 4.8.3 测试建议
1. 单元测试覆盖核心逻辑
2. 集成测试验证模块协作
3. 性能测试确保满足需求
4. 容错测试验证可靠性

## 4.9 应用级调试与优化

### 4.9.1 调试方法

1. **单元测试**：测试每个模块的功能正确性
2. **集成测试**：测试模块间的协作和通信
3. **性能测试**：测试应用的性能和扩展性
4. **容错测试**：测试应用的容错和恢复能力

### 4.9.2 性能优化

1. **瓶颈分析**：分析性能瓶颈，定位优化点
2. **算法优化**：优化算法实现，提高计算效率
3. **并行优化**：优化并行策略，提高资源利用率
4. **I/O优化**：优化数据读写，减少I/O开销

### 4.9.3 监控调优

1. **实时监控**：实时监控应用运行状态
2. **性能分析**：分析性能数据，找出优化空间
3. **动态调整**：根据负载动态调整配置参数
4. **容量规划**：根据业务需求规划资源容量