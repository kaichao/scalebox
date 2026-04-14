# 2. 模块开发

## 2.1 模块设计原则

### 2.1.1 无状态设计
- 任务运行的无状态，可重复运行
- 任务级幂等性
- 灵活调用，通过环境变量、任务头定义不同逻辑处理

### 2.1.2 通用性设计
- 支持任务头驱动
- 基于主路由的任务头驱动的流程控制
- 标准化接口设计

## 2.2 模块设计方法

### 2.2.1 task-body设计
task-body为任务标识，在模块中具有唯一性：
- 简单任务：直接用待处理文件名相对路径
- 复杂任务：用输入路径的规范标识（URI）
- 支持JSON格式的非空字符压缩表示

### 2.2.2 task-headers设计
- 配置参数
- 运行时需调整的配置参数通常以环境变量形式传递
- 在外部编排系统中可以修改配置参数

### 2.2.3 模块内目录设计

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

## 2.3 模块集成实现

### 2.3.1 Sidecar模式
模块采用sidecar模式运行：
- **run**：任务的单次运行
- **check**：任务运行的前置条件检测
- **setup**：设置初始运行环境
- **teardown**：退出前清理计算环境

### 2.3.2 算法运行run.sh
用户程序运行结束后，agent对用户程序的运行做后处理，主要通过以下文件交换信息：

| 文件名                           | 文件说明                                   |
| -------------------------------- | ----------------------------------------- |
| ${WORK_DIR}/task-exec.json       | 任务运行结果主文件，以JSON形式纪录用户程序运行结果 |
| ${WORK_DIR}/sink-tasks.txt       | 后续任务列表文件，每行一个任务 |
| ${WORK_DIR}/extra-attributes.txt | 运行附加属性文件 |
| ${WORK_DIR}/timestamps.txt       | 自定义时间戳文件，常用于调试程序、测试程序性能 |
| ${WORK_DIR}/input-files.txt      | 输入文件列表，用于统计输入文件字节数 |
| ${WORK_DIR}/output-files.txt     | 输出文件列表，用于统计输出文件字节数 |
| ${WORK_DIR}/removed-files.txt    | 待删除文件列表 |
| ${WORK_DIR}/auxout.txt           | 辅助输出文件 |

## 2.4 模块的镜像封装

### 2.4.1 基于agent构建算法模块
```dockerfile
FROM hub.cstcloud.cn/scalebox/agent:latest
```

### 2.4.2 拷贝agent组件至算法模块
```dockerfile
COPY --from=hub.cstcloud.cn/scalebox/agent:latest /usr/local/ /usr/local/
```

## 2.5 模块单元测试

### 2.5.1 独立模块测试
缺省的模块代码目录为：`./code`

```bash
# 单独命令行启动模块测试
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path}

# 先创建模块测试，再添加测试任务
app_id=$(scalebox run --image-name ${module_image_name} --code-path ${code_path}| cut -d':' -f2 | tr -d '}')
scalebox task add --app-id=${app_id} --header header1=${header1} ${task_body}
```

### 2.5.2 带路由的单元测试
缺省的路由代码目录为：`./mr-code`

```bash
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path} --mr-image-name ${mr_image_name} --mr-code-path ${mr_code_path}
```

### 2.5.3 基于流水线应用的单元测试
- 写独立的应用定义文件app.yaml、环境变量定义文件scalebox.env

## 2.6 模块类型详解

### 2.6.1 基础模块（agent）
- 提供了模块底座功能
- 应用模块的中任务执行的生命周期管理
- 与server端runtime的controld交互，纪录任务执行状态

### 2.6.2 文件传输模块
- **file-copy**：基于rsync-over-ssh的文件传输
- **dir-copy**：目录拷贝模块
- **rsync-copy**：rsync拷贝模块
- **ftp-copy**：FTP协议传输模块

### 2.6.3 辅助功能模块
- **dir-list**：目录列表模块
- **cron**：定时任务模块
- **cluster-head**：集群头节点模块
- **node-agent**：节点代理模块

## 2.7 错误处理与容错机制

### 2.7.1 任务级容错
- 按单个task返回码，安排自动重试
- 支持任务重试机制

### 2.7.2 组合式容错
- 解决复杂的非逻辑错
- 以任务级容错为基础
- slot级出错处理

## 2.8 最佳实践

### 2.8.1 模块开发规范
1. 保持模块单一职责
2. 支持环境变量配置
3. 实现任务幂等性
4. 合理设置超时时间

### 2.8.2 性能优化
1. 优化数据局部性
2. 减少不必要的文件操作
3. 合理使用缓存
4. 批量处理优化

### 2.8.3 测试建议
1. 单元测试覆盖核心逻辑
2. 集成测试验证模块协作
3. 性能测试确保满足需求
4. 容错测试验证可靠性