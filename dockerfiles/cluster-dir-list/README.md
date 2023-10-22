# cluster-dir-list

## 模块介绍
集群内目录的文件列表，远端目录基于rsync-over-ssh实现

## 环境变量
  - SOURCE_CLUSTER: 若为空，则再由消息体确定
  - JUMP_SERVER_OPTION: "source"
  - REGEX_FILTER: 
  - REGEX_2D_DATASET: 
  - INDEX_2D_DATASET: 
  
## 输入消息格式

[<SOURCE_CLUSTER>~]<RELATIVE_PATH>#<local_dir>

- 若环境变量为空，则以输入消息中的<SOURCE_CLUSTER>来替代
- RELATIVE_PATH，相对CLUSTER_DATA_ROOT的相对根目录，不包含在产生的消息体中
- local_dir，本地路径，包含在产生的消息体中

## 用户应用的退出码
- 0 : OK 

## 输出消息格式
- 若退出码为0，则输出与输入消息相同的消息。
- 退出码非0，则不输出消息
