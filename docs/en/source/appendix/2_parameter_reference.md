# 2. Parameter Configuration Manual

## 2.1 App Configuration Parameter Tables

### 2.1.1 app-parameters Table
| Parameter name     |   Meaning                          |
| ------------------- | -------------------------------- |
| initial_status      | 'RUNNING' / 'INITIAL' / 'PAUSED' |
| main_router         |  main router module name          |
| default_idle_polls  | the default number of idle polls for the modules it belongs to |
| is_cluster_admin    |                                  |
| default_sleep_count | the default value of each module's sleep_count |
| global_directories  | list of global directories, used to classify read/write data volumes of the modules it belongs to (local/global/tmpfs) |
| slot_group          | multi-node multi-module-slot configuration in map form, used for automatic slot creation for newly added nodes. Example: '{"module0":n0,"module1",n1}' |

### 2.1.2 module-arguments Parameter Table

| Parameter name        | Standard environment variable | Meaning |
| ---------------------- | ---------------------- | ---------------------------------------------------------------------- |
| grpc_server           | GRPC_SERVER            | controld's service endpoint, ```${ip_addr}:${port}```, default port 50051 |
| code_path             |                        | the module code directory, mapped into the container at /app/bin via a data volume |
| local_ip_index        | LOCAL_IP_INDEX         | 'hostname -I' returns an IP address list; this parameter specifies the nth IP address in the list as the local IP address. (Not yet used in agent?) |
|                       | CLUSTER                | the cluster name it belongs to (difference from CLUSTER_NAME, _CLUSTER?) |
|                       | PLAT_MODULE_NAME       | the current module name |
|                       | MODULE_ID              | module-id |
|                       | FROM_MODULE               | module-name |
|                       | FROM_IP                | from-ip |
|                       | SINK_MODULE               | the name of the default sink_module |
|                       | IS_SINGULARITY         | the container engine is singularity or apptainer |
| task_max_seconds      | PLAT_TASK_MAX_SECONDS       | the timeout in seconds for each task run; if the run time exceeds this limit, the task run is interrupted, returning timeout code 124; global timeout returns 123. |
| task_min_seconds      | PLAT_TASK_MIN_SECONDS       | the minimum seconds for each task run; if the run time is continuously below this value multiple times, the slot is judged as GREEDY anomaly and exits |
| poll_interval_seconds | PLAT_POLL_INTERVAL_SECONDS | the polling interval for slots checking new tasks. The slot sleeps and periodically checks task availability; this parameter specifies the interval in seconds, default 6 seconds. |
| max_idle_polls        | PLAT_MAX_IDLE_POLLS    | the maximum number of idle polls a slot can perform before exiting. The maximum number of sleeps before a slot exits. Default 100 (10 minutes) |
| dir_quota_gb          | PLAT_DIR_QUOTA_GB      | standard flow control parameter, the capacity quota limit of the directory itself. Used to specify the maximum space of a directory in GB. Format: ```'{"/dir-1":10,"/dir-2":100}'``` |
| free_space_gb         | PLAT_FREE_SPACE_GB     | standard flow control parameter, the minimum space that must be reserved on the disk where the directory resides. Used to specify the minimum reserved space in GB of the partition where the directory resides. Format: ```'{"/dir-3":10,"/dir-4":100}'``` |
| heartbeat_seconds     | PLAT_HEARTBEAT_SECONDS | heartbeat interval in seconds, default 60; if a non-positive integer, heartbeat is disabled |
| output_text_size      | PLAT_OUTPUT_TEXT_SIZE  | the maximum bytes of large text fields (stdout/stderr/auxout) in the task run record t_task_exec. Default 65535; maximum can be 10MB (for varchar) or 1GB (for text) |
| text_trunc_mode       | PLAT_TEXT_TRUNC_MODE   | 'HEAD'/'TAIL', default value is 'HEAD', head truncation, keeping the tail part |
| timezone_mode         |                        | 'HOST'/'UTC'/'NONE' |
| slot_options          |                        | comma-separated slot options |
|  - always_running     | PLAT_ALWAYS_RUNNING    | the slot keeps running and never exits proactively (generally only for debugging) |
|  - reserved_on_exit   |                        | after the slot exits, keep the container for troubleshooting. (docker-only, remove --rm from the command line) |
|  - tmpfs_workdir      |   TMPFS_WORKDIR        | use the tmpfs file system for the working directory /work (for docker, add --tmpfs /work to the resolved command line; for singularity, add the environment variable TMPFS_WORKDIR=yes after resolution) |
|  - disable_local_mapping |                     | do not generate the mapping from the local root directory to /local_data_root in the container |
|  - disable_data_mapping  |                     | do not generate the mapping from the cluster data directory to /cluster_data_root in the container |
|  - slot_on_head          |                     | generate a single slot instance on the head node |

### 2.1.3 module-parameters Parameter Table

| Parameter name        | Description |
| -------------------- | ---------------------------------------------------------------------------- |
| priority             | priority (not yet used) |
| key_group_regex      | the regular expression for extracting groups from messages (renamed to group_regex?) |
| key_group_index      | the group ordering number (renamed to group_index?) |
| task_dist_mode       | task distribution mode, 'HOST-BOUND'/'SLOT-BOUND'/'DEFAULT' |
| start_task        | given an initial message; if it is 'FILE:{filename}', each line of the file is treated as an initial message |
| initial_task_status  | the initial status of tasks, 'READY'/'INITIAL' |
| initial_slot_status  | the initial status of slots, 'READY'/'OFF' |
| retry_rules          | retry rules based on exit codes<br>```"['exit_code_1:num_retries',...,'exit_code_n:num_retries']"```. The default value of num_retries is 1; the exit code wildcard is '*'. Proposed format change to ```{"exit_code_1":num_retries',...,"exit_code_n":num_retries}``` |
| slot_recoverable     | 'yes', supports setting a slot that errored and exited from 'ERROR' to 'READY' to support slot-level retries (replaced by slot_max_retries? TIMED-OUT/GREEDY handled separately) |
| slot_max_retries     | the number of retries for setting a slot status from 'TIMED-OUT' to 'READY' (?) |
| slot_timeout_minutes | if a slot does not start normally, it stays in the 'STARTING' state. Set a timeout in minutes; upon expiry, the status transitions to 'TIMEOUT'. Default 15 minutes. For slot instances that cannot be restarted repeatedly (using GPUs, etc.), a larger value can be set. |
| task_global_timeout_scale | if external causes (abnormal slot exits, etc.) leave a task running (status code -3). Through the global timeout setting, restore the task status code to 123. This value is a multiple relative to task_max_seconds, default 2.0. The global exit return code is 123. Proposed rename to task_timeout_scale_factor? |
| router_index | in apps with multiple main router instances, specifies that the current module sends to the nth main router. Default 0; usually set >0 to specify a particular main-router instance |
| pod_id               | identifies that this module is managed by a pod; if the pod the message originates from has the same pod_id, all tasks are marked as using local computing (task_dist_mode is HOST_BOUND) |
| task_dedup_cache_ttl_minutes | the time-to-live of the task deduplication cache; needs setting under high load. Sets the cache expiration time (in minutes) for duplicate task-id detection, default 30 minutes, cleared at n+1 minutes. Avoids multiple distributions of the same task. Normally, its time should be greater than the value of ```task_max_seconds```. |
| task_batch_size       | the number of tasks fetched in a batch. Sets the maximum number of messages a slot reads in a single batch, default 1. For tasks with run times within 5 seconds, batch message reads can be configured to avoid server overload and data inconsistency from frequent reads. |
| visible              | whether visible in the pipeline logic diagram. Default 'yes' |
| task_id_in_headers   | 'no'/'yes', the headers of tasks in this module include the task_id value. |
| node_progress_gap    | standard flow control parameter, for synchronizing the execution of nodes within the same group of a specified module, the maximum difference in running tasks between the fastest node and the slowest node; its value is an integer. When the corresponding slot is generated, the corresponding semaphore is automatically created, named ```node_progress:${mod_name}:${hostname}```, initial value 0. Example format of this parameter: ```{"prefix1":4,"node_prefix2":6}```. Proposed rename to node_progress_max_gap? |
| vtask_role   | 'head' / 'core' / 'tail', the current role in vtask processing, valid only for non-router modules. head is the starting module of vtask processing; tail is the ending module; core is the algorithm module. |
| vtask_size   | standard flow control parameter, defines the upper limit of the number of vtasks that can be simultaneously in ready/running states; when the app is parsed, the corresponding semaphore and initial value are created.<br/>The semaphore name for global vtask flow control is: ```vtask_size:${mod_name}```;<br/>The semaphore name for SLOT-BOUND (GROUP-BOUND vtask) flow control is: ```group_vtask_size:${mod_name}:${group_seq}```;<br/>The semaphore name for HOST-BOUND vtask flow control is: ```host_vtask_size:${mod_name}:${hostname}``` |
| vtask_size_sema_copy | 'no' / 'yes'. Generate a copy of the vtask_size semaphore for programmatic control. The copy semaphore name is the original semaphore name prefixed with a colon ```:```, with the same initial value as the original semaphore. |

### 2.1.4 task-headers Standard Parameter Table

| Parameter name  | Meaning |
| ---------------- | ----------------------------------------------- |
| to_ip            | the IP of the host to process the current task |
| to_host          | the hostname to process the current task (t_host primary key) |
| from_ip          | the IP of the host that generated the task |
| from_ip_last     | if from_module is the main router, the host IP |
| from_host        | the hostname that generated the task (t_host primary key) |
| from_module      | the module name that generated the task |
| from_module_last | if from_module is the main router, the from_module before the main router |
| from_cluster     | the previous cluster name in cross-cluster apps |
| to_slot          | the slot_id to process the current task in SLOT-BOUND modules |
| to_slot_index    | the slot_id expressed as the local seq |
| repeatable       | by default, tasks cannot be distributed repeatedly within a specified time; the default value can be customized through the module's task_cache_expired_minutes parameter; in retry operations and specific scenarios, if repeated message distribution is needed, set this parameter to 'yes' |

Among them, from_ip, from_module, from_module_last, etc. are automatically generated by the system.

- For HOST-BOUND modules, to_host must be set in the code; or set through the module's pod_id equality.
- For SLOT-BOUND modules, to_slot must be set in the code

task-header parameters can be customized, where parameters starting with ```_``` are standard passing parameters; parameters starting with ```_vtask_``` are vtask parameters maintained by the system.

## 2.2 Resource Configuration Parameter Tables

### 2.2.1 cluster-parameters Table
| Parameter name      |   Meaning                                              |
| ------------------ | ------------------------------------------------------ |
| data_root          |  the /cluster_data_root directory in the container (rename to data_root?) |
| code_dir           |  unused                                                |
| uname              | ssh login username |
| port               | ssh login username |
| grpc_server        | controld connection info, ${controld_ip}[:${port}], default port 50051 |
| remote_grpc_server | controld connection info for remote routing connections in cross-cluster apps |
| remote_check_mode  | 'ip': exact IP comparison, for test clusters, multi-cluster same-segment; 'subnet': subnet matching, for production clusters, cross-segment deployments. |
| initial_status     | 'ON'/'OFF' |

### 2.2.2 host-parameters Table

| Parameter name   |   Meaning                       |
| --------------- | ----------------------------- |
| uname           | ssh login username |
| port            | ssh port number |
| node_slot       | the slot number of this node's node-agent |
| group_id        | the number of the node group this node belongs to |
| slurm_node      | the corresponding node number in the slurm scheduling system |
| reg_time        | registration time in scalebox |
| slot_module_id  | the slurm module id of the node-agent in the slurm scheduling system |
| default_runtime | default container engine, 'docker'/'singularity'/'podman' |
| use_home_tmp    | for nodes that do not allow running programs under /tmp (I/O nodes), use $HOME/tmp instead of /tmp as the temporary directory for remotely starting slots. |
| ssh_cepher      | 'aes128-gcm@openssh.com' / 'aes256-gcm@openssh.com' / 'chacha20-poly1305@openssh.com' |

### 2.2.3 slot-parameters Table

| Parameter name | Meaning |
| ------------ | --------------------------------------- |
| reg_time     | registration time |
| last_access  | last access time |
| group_prefix | the slot's group prefix, for grouped slots in HOST-BOUND modules |

## 2.3 Environment Variable Reference

- Standard environment variable table for module programming

| Environment variable name | Description |
| --------------- | --------------------------------------------------------------------------------- |
| ACTION_CHECK    |  the path of the custom admission control script (check). If the return value is 0, OK; non-zero, admission restricted. |
| ACTION_RUN      |  the path of the custom main algorithm script (run) |
| ACTION_SETUP    |  the path of the custom initial setup script (setup). If the return value is 0, OK; non-zero, initialization failed, set the current slot status to ERROR. |
| ACTION_TEARDOWN |  the path of the custom exit script (teardown). If the return value is 0, OK; non-zero, exit failed, set the current slot to ERROR. |
| APP_ID          | |
| MODULE_ID       | |
| SLOT_ID         | |
| CLUSTER         | CLUSTER |
| WORK_DIR        | local working directory, default /work |
| LOG_LEVEL       | 'info'/'debug'/'trace'; trace in principle is only for single-module debugging and troubleshooting. |
| LOCAL_IP        | local IP address |
| FROM_IP         | already included in headers? |
| FROM_MODULE     | already included in headers? |
| SINK_MODULE     | in routed apps, non-router modules point to the router module; in router-less apps, point to the next module. |
| GRPC_SERVER     | the server-side grpc address, format: ```server_name[:port]```, default port 50051. |
| PGHOST          | PGHOST, query job-attr-set from db |
| IS_SINGULARITY  | IS_SINGULARITY |
| LOCAL_SHMDIR    | the working directory under local tmpfs (/dev/shm) |
| LOCAL_TMPDIR    | the working directory under local /tmp |
| TMPDIR_GROUP    | the group number (numeric) in TMPDIR |
| TMPFS_WORKDIR   | use tmpfs for the local working directory to improve large-file loading performance |
| SLOT_ROLE       | ''/'group', the group slot in group computing mode, which can fetch all tasks in the group across nodes |

All environment variables starting with PLAT_ do not need to be accessed by user programs.

