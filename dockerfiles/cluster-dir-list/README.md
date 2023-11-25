# cluster-dir-list

## 一、模块介绍
集群内目录的文件列表，远端目录基于rsync-over-ssh实现

## 二、环境变量
  - SOURCE_CLUSTER: 若为空，则再由消息体确定
  - JUMP_SERVER_OPTION: "source"，通过source端的头节点作为跳板
  - REGEX_FILTER: 
  - REGEX_2D_DATASET: 
  - INDEX_2D_DATASET: 
  
## 三、输入消息格式

[<SOURCE_CLUSTER>]~<RELATIVE_PATH>#<local_dir>

- 若环境变量为空，则以输入消息中的<SOURCE_CLUSTER>来替代
- RELATIVE_PATH，相对CLUSTER_DATA_ROOT的相对根目录，不包含在产生的消息体中
- local_dir，本地路径，包含在产生的消息体中

## 四、用户应用的退出码
- 0 : OK 

## 五、输出消息格式
### 5.1 数据集元数据消息

示例如下：
```json
{
		"datasetID":"datasetID",
		"horizontalWidth": 24,
		"verticalStart": 1,
		"verticalHeight": 855
}
```

- datasetID为<RELATIVE_PATH>#<local_dir>。
- horizontalWidth为二维表的水平宽度
- verticalStart、verticalHeight则为垂直方向的起始值、高度。


### 5.2 文件实体消息

针对目录下文件，生成一条消息，消息格式为：
从<local-dir>起的相对路径。

### 5.3 目录扫描完成的控制消息
- 若退出码为0，则输出与输入消息相同的消息
- 退出码非0，则不输出该消息
