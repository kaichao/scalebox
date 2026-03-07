# 4. 用户工具

## 4.1 命令行工具

```{mermaid}

graph LR
  scalebox --> run[<a href="#run">run</a>]

  scalebox --> cluster[<a href="#cluster">cluster</a>]
  cluster --> cluster-get-parameter[<a href="#cluster-get-parameter">get-parameter</a>]
  cluster --> cluster-check-status[<a href="#cluster-check-status">check-status</a>]
  cluster --> cluster-dist-image[<a href="#cluster-dist-image">dist-image</a>]
  cluster --> cluster-app-view[<a href="#cluster-app-view">app-view</a>]
  cluster --> cluster-host-view[<a href="#cluster-host-view">host-view</a>]

  scalebox --> host[<a href="#host">host</a>]
  host --> host-check-status[<a href="#host-check-status">check-status</a>]
  host --> host-get-info[<a href="#host-get-info">get-info</a>] 
  host --> host-add-node[<a href="#host-add-node">add-node</a>]
  host --> host-dist-image[<a href="#host-dist-image">dist-image</a>]
  host --> host-recover[<a href="#host-recover">recover</a>]
  host --> host-migrate[<a href="#host-migrate">migrate</a>]
  host --> host-replace[<a href="#host-replace">replace</a>]
  host --> host-asign[<a href="#host-asign">asign</a>]
  
  scalebox --> slot[<a href="#slot">slot</a>]
  slot --> slot-add[<a href="#slot-add">add</a>]
  slot --> slot-add-group[<a href="#slot-add-group">add-group</a>]
  slot --> slot-update[<a href="#slot-update">update</a>]

  scalebox --> app[<a href="#app">app</a>]
  app --> main-router[<a href="#app-main-router">main-router</a>]
  app --> app-list[<a href="#app-list">list</a>]
  app --> app-add-remote[<a href="#app-add-remote">add-remote</a>]
  app --> app-set-finished[<a href="#app-set-finished">set-finished</a>]

  scalebox --> module[<a href="#module">module</a>]
  module --> module-list[<a href="#module-list">list</a>]
  module --> module-info[<a href="#module-info">info</a>]

  scalebox --> task[<a href="#task">task</a>]
  task --> task-add[<a href="#task-add">add</a>]
  task --> task-get-header[<a href="#task-get-header">get-header</a>]
  task --> task-set-header[<a href="#task-set-header">set-header</a>]
  task --> task-remove-header[<a href="#task-remove-header">remove-header</a>]

  scalebox --> semaphore[<a href="#semaphore">semaphore</a>]
  semaphore --> sema-create[<a href="#semaphore-create">create</a>]
  semaphore --> semaphore-get[<a href="#semaphore-get">get</a>]
  semaphore --> increment[<a href="#semaphore-increment">increment</a>]
  semaphore --> decrement[<a href="#semaphore-decrement">decrement</a>]
  semaphore --> increment-n[<a href="#semaphore-increment-n">increment-n</a>]
  semaphore --> semaphore-group[<a href="#semaphore-group">group</a>]

  scalebox --> semagroup[<a href="#semagroup">semagroup</a>]
  semagroup --> semagroup-min[<a href="#semagroup-min">min</a>]
  semagroup --> semagroup-max[<a href="#semagroup-max">max</a>]
  semagroup --> semagroup-increment[<a href="#semagroup-increment">increment</a>]
  semagroup --> semagroup-decrement[<a href="#semagroup-decrement">decrement</a>]
  semagroup --> semagroup-diffmin[<a href="#semagroup-diffmin">diffmin</a>]
  semagroup --> semagroup-diffmax[<a href="#semagroup-diffmax">diffmax</a>]

  scalebox --> variable[<a href="#variable">variable</a>]
  variable --> variable-get[<a href="#variable-get">get</a>]
  variable --> variable-set[<a href="#variable-set">set</a>]

  scalebox --> global[<a href="#global">global</a>]
  global --> global-get[<a href="#global-get">get</a>]
  global --> global-set[<a href="#global-set">set</a>]


  scalebox --> fs[<a href="#fs">fs</a>]
  fs --> fs-ls[<a href="#fs-ls">ls</a>]
  fs --> fs-stat[<a href="#fs-stat">stat</a>]

  scalebox --> status

  scalebox --> help

```

## 4.2 WebUI工具

## 4.3 数据库工具(adminer)
