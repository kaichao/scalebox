# agent-interface-files

agent的用户程序与基础平台的接口文件。

## 接口文件的格式及示例

### task-exec.yaml

### extra-attrs.yaml

### sink-tasks.txt

### timestamps.txt
- 格式为:label,时间，```#```开头为注释行

### input-files.txt
文件目录的绝对路径，```#```开头为注释行

### output-files.txt
文件目录的绝对路径，```#```开头为注释行

### removed-files.txt
文件目录的绝对路径，```#```开头为注释行

### auxout.txt

用户自定义文件，格式不限。存放在t_task_exec的auxout字段中


## 示例运行

### 示例说明

共有2个模块。

- 主路由模块：启动业务模块，接收返回任务
- 业务模块：展示各个接口文件的格式及用法

### 
```sh
echo start-task | scalebox run
```
