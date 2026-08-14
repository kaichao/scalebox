# 7. Advanced Programming Features

- §7.1 Fault-tolerance mechanisms (task-level/platform-level/cross-module fault tolerance + tiered strategies)
- §7.2 Admission control (storage capacity/counting/synchronization/custom flow control)
- §7.3 Task execution ordering (sort_tag / group_regex / group_index)
- §7.4 Timeout settings (task_max_seconds + timeout wrapper recommendations)
- §7.5 Timestamp formats (RFC3339Nano / bash / golang examples)
- §7.6 Cross-cluster apps
- §7.7 Slot launch command expressions (`((@seq))` / `((@v:var_name))` dynamic evaluation)
- §7.8 gRPC Metadata parameter passing (semaphore-auto-create)
- §7.9 Slot auto-scaling

## 7.1 Fault-Tolerance Mechanisms

Module execution errors generally fall into two categories:
- Logical errors: incorrect results caused by code logic;
- Non-logical errors: incorrect results caused by anomalies in the external software/hardware environment; usually resolvable by retrying execution multiple times after the external environment recovers;

Non-logical errors occur with a certain probability. In large-scale data processing, automatically handling non-logical errors greatly improves the efficiency and reliability of data processing.

### Common External Runtime Environment Anomalies
- Unstable hardware and base software
  - Unstable network device links
  - Rising ambient temperatures causing compute nodes to temporarily stop serving
  - Unstable media in disk arrays causing temporary RAID failures
  - Temporary, fixable bugs in cluster storage software
- Random anomalies during operation
  - High load causing increased core storage response latency
    - Variable latency (seconds to minutes)
    - Temporarily unwritable
  - High load causing low-probability memory oversubscription and insufficient application memory
  - High system load and inappropriate scheduling policies causing algorithm execution timeouts
- Data anomalies
  - Abnormal input data
  - Random errors in preceding processing software
  - Incomplete data caused by various runtime anomalies

### Fault Handling for Non-Logical Errors
- Judgment and localization of system anomalies
  - Read/execute/write-back times
  - Task return codes
  - Anomaly information from system monitoring
- Fault handling for anomalous tasks
  - Automatic/manual


### 7.1.1 Module-Level Fault Tolerance
- Set the module parameter retry_rules.
Set different error codes in the main script according to the algorithm code's handling of non-logical errors. Based on the returned error code, configure automatic retries in the module definition, thereby implementing automatic retries.

In the Module definition, set ```retry_rules``` to implement the above functionality.

```
    retry_rules: "['91','92:3']"
```
In the above setting, if the return code is 91, retry once; if the return code is 92, retry 3 times.

The return code is the result code of the module execution; the return code of a failed module execution is the error code.

### 7.1.2 Platform-Level Failures and Recovery

- The slot starts and runs, but the task errors quickly when running, causing short task runtimes and allowing multiple tasks to complete per unit time;
- The module's max_tasks_per_minutes parameter (value >=3): when exceeded, the corresponding slot is set to 'ERROR' to avoid interfering with the normal waiting of other tasks;
- Combined with retry_rules, problematic slots can be stopped when a large number of slots are erroring;
- Finally, manual recovery can be performed.

### 7.1.3 Cross-Module Fault Tolerance

- Based on vtasks
- Add management task headers to tasks to implement vtask-based management (deletion?)
- Clean up tasks, semaphores, variables, etc. designed in the vtask


Fault tolerance tiers
- task-level: exit-code-based retry mechanism. Retrying tasks within a vtask requires adjusting the vtask state. Recover at the task level as much as possible.
- slot-level:
- vtask-level: in node-local computing mode, ensure final idempotency. The vtask layer maintains strong consistency.
  - Retrying tasks within a vtask requires adjusting the vtask state
  - Partial failures (node exits) can only be recovered at the vtask level.


## 7.2 Admission Control Mechanisms

Admission Control

Local computing has limited local resources, so admission control must be introduced to avoid overloading key local node resources (memory capacity, disk space).
For local computing in complex applications, advanced admission control can greatly simplify business logic and avoid runtime anomalies.

Admission control for algorithm modules controls the progress of streaming data processing, so that the non-shareable core computing resources required for data processing (memory cache, GPU memory, local disk space) do not exceed the physical limits of local resources, which would cause abnormal computation interruptions, including abnormal exits and program deadlocks.

The implementation mechanism of flow control is to check, before message processing, whether external conditions are met through the module's standard flow control rules and custom flow control rules. If not, explicitly call sleep() in the code, wait for a period, then check again. If multiple checks fail, exit the module execution.

Flow control rules introduce a waiting mechanism, which lowers the utilization of local computing resources but ensures the smooth execution of the application as a whole.
Avoiding peak memory space requirements, at the cost of increased waiting time, allows running local computing applications on machines with low memory configurations.

Flow control is divided into slot-level flow control and inter-node parallel synchronization flow control.
Flow control supports the synchronization of module execution to avoid runtime anomalies caused by limited system resources.

Precisely control data output; if data space is insufficient, wait for other parallel tasks to complete and release resources, clear out space, and continue running after the space capacity meets requirements.

### 7.2.1 Admission Control Classification
#### 7.2.1.1 task addmission
- Rule-based Admission Control
  - Storage Capacity Admission Control
  - Counter-based Admission Control
- Custom Admission Control

#### 7.2.1.2 slot addmission
How many slots can be started?

### 7.2.2 Storage Capacity Admission Control

(Storage Capacity Admission Control)
#### 7.2.2.1 Directory Maximum Space Flow Control

- Standard attribute name: dir_quota_gb; flow control implemented by limiting the maximum storage space occupied by a directory.

#### 7.2.2.2 Minimum Free Space on the Directory Partition

- Standard flow control attribute: space_free_gb; flow control implemented by limiting the minimum free space of the storage partition where the directory resides

### 7.2.3 vtask Counting Admission Control

- Standard flow control attribute ```vtask_size```
  - Understood as the maximum number of tasks that can run
  - In the vtask framework, the counter is decremented in the vtask-head module and incremented in the vtask-tail module
  - In standalone modules, the counter is decremented. The increment operation must be done manually in the main router.

Modules with vtasks enabled automatically decrement the counter when a task starts; if the counter value <= 0, counting flow control takes effect and new tasks cannot run. From an implementation perspective, the vtask mechanism sets the waiting queue length.

After subsequent module task processing completes, the vtask counter is automatically restored in the main router. The counter can also be operated manually to implement start/stop.

- Processing flow
  - At the module entry, check the semaphore counter value; if <= 0, there are no available slots and admission control fails;
  - The main router automatically decrements the counter, occupying one running slot. This operation is automatic and requires no programming control;
  - When vtask processing completes, the counter is incremented through semaphore operations in the corresponding message routing, releasing an empty slot.


#### SLOT-BOUND Modules
- Generates the corresponding semaphore: ```slot_vtask_size:${mod_name}:${slot_id}```, with the initial value being the parameter value.
  - mod_name is the first module name
  - slot_id is the ID of the SLOT
- The maximum number of vtasks that can run on a compute slot, usually representing the length of the running queue on a compute node group.

#### HOST-BOUND Modules
- Generates the corresponding semaphore: ```host_vtask_size:${mod_name}:${hostname}```, with the initial value being the parameter value.
  - mod_name is the first module name
  - hostname is the hostname in the t_host table
- The maximum number of vtasks that can run on a compute node, representing the length of the running queue on the node.

#### Default-Type Modules
- Generates the corresponding semaphore: ```vtask_size:${mod_name}```, with the initial value being the parameter value.
  - mod_name is the first module name
- The maximum number of vtasks that can run globally, representing the length of the running queue on the node.



### 7.2.4 Node Group Parallel Synchronization Admission Control
- Standard flow control attribute ```node_progress_gap```
  - Based on per-node progress counters, implemented through synchronization semaphores
  - node_progress_diff: the difference between the current node's progress and the slowest progress serves as the control variable

- In the main router, modify the corresponding semaphores according to task processing progress;
- Business modules implement synchronized flow control based on semaphore counter values;



### 7.2.5 Task-Level Custom Admission Control

- Code-customized flow control (can be based on semaphores)
  - ACTION_CHECK/check.sh: custom flow control logic


## 7.3 Task Execution Ordering for Key Modules

Task execution ordering is an important foundation for running scalebox applications.

For modules involving multi-node collaborative processing, ordering allows all data from the same time period to run in parallel on different compute resources, enabling collaboration (grouping, etc.) in subsequent modules.

Ordering is implemented by setting the following parameters.

### 7.3.1  Sort Tags
  The sort_tag in the message header is a dedicated tag for message ordering, with the highest priority, usually set by the main-router

### 7.3.2 Task Group Numbers
- group_regex: a regular expression that extracts the relevant grouping string from the task body.
- group_index: the group number corresponding to the regular expression.

### 7.3.3 Task Processing Order

If the aforementioned ordering methods are not set, tasks are processed in the order they were generated by default

- Idempotency: node-local computing carries the possibility of node failure, causing local storage to fail and making task-level idempotency unguaranteed; idempotency is implemented at the cross-module vtask level.

## 7.4 Timeout Settings (timeout)

Due to code bugs, data errors, and other causes, algorithm modules may enter infinite loops, occupying slots indefinitely and failing to exit normally.

Timeout settings solve the above problem. By setting the module's ```task_max_seconds``` parameter, the task exits after the time limit is exceeded, returning error code 124.

When a wrapper script calls the algorithm code, it is recommended to also add timeout-related settings in the wrapper script. Without this setting, if the outer script times out, the inner script will still run to completion, causing uncertain results.
If the algorithm code is named run.py, the corresponding part of the wrapper code can be written as:
```sh
timeout ${TASK_TIMEOUT_SECONDS}s run.py $*
```


## 7.4 task-perspective

- scalebox standard time format
```2023-10-24T18:00:00.123456+0800```

- bash generates the current time 
```sh
date +"%Y-%m-%dT%H:%M:%S.%6N%z"
```

- golang generates the current time 
```go
formattedTime := time.Now().Format("2006-01-02T15:04:05.999999")
fmt.Println(formattedTime)
```

- python generates the current time 
```python

```


## 7.5 Git App Code Repositories

## 7.6 Cross-Cluster Apps

## 7.7 Slot Launch Command Expressions

Expressions delimited by `(( ))` are supported; at slot launch time, the expression is dynamically evaluated according to the external environment to obtain the final launch command. Applicable to the following scenarios:
- A single node with multiple GPU cards, where different commands are needed to launch the slots corresponding to different GPU cards
- In admission control rules, the slot sequence number is needed as a variable to precisely define storage capacity limits, ensuring full resource utilization.

### 7.7.1 Variable Types

In expressions, variable names start with ```@``` and are divided into numeric variables and string variables.

#### Numeric Variables (slot sequence numbers)
Each module, for each node, has a sequence number counted from 0, which can be used in slot launch command expressions to use different launch commands for different slot sequence numbers.

Use ```@seq``` to identify the count sequence number

#### String Variables (global variables)
The final slot command is resolved based on global variable values; image name version numbers can be specified

Use ```@v:var_name``` to identify string variables

### 7.7.2 Slot Sequence Number Variable Examples

A single compute node can be configured with multiple GPU accelerator cards. Usually each GPU is configured with an independent slot, and module algorithms do not need parallel optimization for multiple GPU cards; such configurations are usually also more efficient.
In this scenario, different GPUs need different slot launch commands.
Slot launch commands support parameterized configuration.


#### ***App Definition***
```yaml
  beam-make:
    command: ${ROCM_COMMAND}
    arguments:
      free_space_gb: '{"${LOCAL_SHMDIR}":${BEAM_MAKE_FREE_GB}}'
```
#### ***Admission Control Parameter Definition***
```bash
LOCAL_SHMDIR=/dev/shm/scalebox/mydata
BEAM_MAKE_FREE_GB='((@seq*5+11))'

```
#### singularity multi-GPU configuration command
```bash
ROCM_COMMAND='singularity exec --rocm --env ROCR_VISIBLE_DEVICES=((@seq)) {{ENVS}} {{VOLUMES}} {{IMAGE}} goagent'
```

#### docker multi-GPU configuration command 1
```bash
ROCM_COMMAND='docker run -d --rm --network host --tmpfs=/work --device=/dev/kfd --device=/dev/dri --security-opt seccomp=unconfined --group-add video -e ROCR_VISIBLE_DEVICES=((@seq)) {{ENVS}} {{VOLUMES}} {{IMAGE}}'
```

#### docker multi-GPU configuration command 2
```bash
docker run -d --rm --network=host --tmpfs=/work --device=/dev/kfd --device=/dev/dri/card((@seq%2)) --device=/dev/dri/renderD((n%2+128)) --security-opt seccomp=unconfined --group-add video --cap-add=SYS_PTRACE {{ENVS}} {{VOLUMES}} {{IMAGE}}
```

### 7.7.3 global Variable Examples

```yaml
  my-module:
    base_image: my-image:(($v:pipeline_version))
```

Define the global variable ```pipeline_version``` as ```20260411``` in global; then at launch time, the actual image name is ```my-image:20260411``` 

```sh
scalebox global set pipeline_version 20260411
```

## 7.8 gRPC Metadata Parameter Passing

Some flags are passed through gRPC metadata to avoid environment variable pollution across processes.

### 7.8.1 semaphore-auto-create

| metadata key | value | purpose |
|-------------|------|------|
| `semaphore-auto-create` | `yes` | Automatically create the semaphore if it does not exist (initial value 0) |

**CLI side** (`metadata.AppendToOutgoingContext`):

```go
func semaContext() context.Context {
    if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
        return metadata.AppendToOutgoingContext(
            context.Background(), "semaphore-auto-create", "yes")
    }
    return context.Background()
}
```

**Server side** (`metadata.FromIncomingContext`):

```go
func getSemaphoreAutoCreate(ctx context.Context) bool {
    if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
        return true
    }
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok { return false }
    for _, v := range md.Get("semaphore-auto-create") {
        if v == "yes" { return true }
    }
    return false
}
```

### 7.8.2 Usage

```bash
# After setting the environment variable, the CLI automatically passes it to controld via metadata
SEMAPHORE_AUTO_CREATE=yes scalebox semaphore get my_sema
```

### 7.8.3 Semaphore Function Signatures

```go
// The autoCreate parameter replaces os.Getenv("SEMAPHORE_AUTO_CREATE")
// The server side obtains it through gRPC metadata; the controld task package passes true directly
func GetValue(name string, appID int, autoCreate bool) (int, error)
func AddValue(name string, delta int, appID int, autoCreate bool) (int, error)
```

## 7.9 Slot Auto-Scaling

### 7.9.1 Overview

Slot auto-scaling automatically adjusts the number of task execution slots according to resource usage, intelligently deciding to scale up or down by monitoring host resource utilization, optimizing resource utilization and task execution efficiency.

### 7.9.2 Features

- **Automatic monitoring**: periodically checks host resource usage
- **Intelligent decisions**: automatically decides to scale up/down based on utilization thresholds
- **Dual strategies**: supports both resource-tight scaling up and capacity-insufficient scaling up
- **Safety protection**: prevents excessive scaling; built-in cooldown mechanisms ensure system stability

### 7.9.3 Configuration Parameters

Set the following key configuration items in module parameters:

```json
{
  "autoscaling.enabled": "true",
  "autoscaling.min_slots": "1",
  "autoscaling.max_slots": "10",
  "autoscaling.scale_up_threshold": "0.8",
  "autoscaling.scale_down_threshold": "0.3",
  "autoscaling.scale_up_step": "2",
  "autoscaling.scale_down_step": "1"
}
```

### 7.9.4 Scaling Strategies

#### Scale-Up Trigger Conditions
1. **Resource-tight scaling up**: when resource utilization exceeds the scale-up threshold (default 80%)
2. **Capacity-insufficient scaling up**: when the current slot count is below the threshold ratio of the maximum available capacity

#### Scale-Down Trigger Conditions
1. **Resource-idle scaling down**: when resource utilization falls below the scale-down threshold (default 30%)
2. **Over-provisioned scaling down**: forced scaling down when the slot count exceeds physical capacity

### 7.9.5 Resource Metric Requirements

The system requires three types of resource metrics to work together:
- **Total host resources** (total.*): total resource definitions in the host configuration
- **Used host resources** (used.*): usage data periodically updated by external monitoring systems
- **Slot resource requirements** (slot_req.*): per-slot resource requirements in module configuration

The three types of metrics must maintain consistent key names and units.

### 7.9.6 Usage Examples

```
# Normal scale-up scenario
Initial state: slot=3, max available=10, utilization=90%
Result: triggers scaling up by 2 slots (utilization>80% and current < max capacity × 80%)

# Normal scale-down scenario  
Initial state: slot=8, max available=10, utilization=20%
Result: triggers scaling down by 1 slot (utilization<30%)

# Over-provisioned scale-down
Initial state: slot=5, max available=1
Result: forced scale-down to within physical capacity
```

## 7.10 Related Documents

The following advanced features have separate detailed descriptions:

| Topic | Document | Description |
|------|------|------|
| VTask task groups | :doc:`9_vtask` | Complete conceptual model, module structure, semaphore mechanism, pipeline flow, CLI command reference |
| Cross-cluster architecture | :doc:`10_cross_cluster` | Data replication model, address resolution, gRPC proxy forwarding, cross-cluster Task routing |
| Security framework | :doc:`11_security` | JWT + RBAC + TLS, certificate management, Automation Token rotation |
| WebUI usage | :doc:`WebUI Guide <../user/7_webui>` | REST gateway, navigation structure, build and deployment |
| VS Code extension | :doc:`VS Code Extension <../user/8_vscode>` | TreeView, DAG, LSP, installation and configuration |
