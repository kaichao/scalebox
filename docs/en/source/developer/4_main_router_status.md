# 4. Main Router and App State

A Scalebox app is a distributed application driven by the main router module, and its state management is a key technology. State management is implemented through state retention and state transition.

The main router module, also called the router module, is the app's main control program.


## 4.1 The Main Router Module

- Multiple main-router instance settings
  - Put all related semaphore modules into one main-router.
  - Set task-dist-mode to SLOT-BOUND
  - During task add, set to_slot in the task-headers

Can be implemented in different programming languages.
- bash
- golang
- python


### 4.1.1 Optimized Design
- By design, the main router module is a single serial module. However, it can be deployed as multiple independent modules based on message sources, achieving parallel processing.
- Reduce the number of router-related input tasks (coarse-grained control messages)
- Reduce the number of system-wide semaphores (coarse-grained semaphores)

### 4.1.2 Efficient Implementation
- In the main router, process multiple tasks in a single pass (one operation reads multiple tasks and batch-processes subsequent task generation, semaphore modifications, etc.)
- Batch operations for subsequent task generation (reducing the number of database operations)
- Batch generation for semaphore processing (reducing the number of database operations)

### 4.1.3 State Management

- State retention
- State transition


Semaphores are an important concept in scalebox for inter-task synchronization and state management. In complex application logic scenarios, they centrally manage the running state of the entire system, making compute modules stateless. Semaphores are usually read and written only in the main-router to avoid concurrency problems. Regular algorithm modules can read the corresponding values.

## 4.2 State Retention
- : Store the current state of app modules in the form of state variables.
  - Global variables
  - Shared variables
  - Semaphores

  - State persistence (semaphore/variable/global/......)

### 4.2.1 Semaphores

- Create a semaphore
```sh
scalebox semaphore create ${sema_name}
```

- Read a semaphore
```sh
scalebox semaphore get ${sema_name}
```

- Increment a semaphore

```sh
scalebox semaphore increment ${sema_name}
```

- Decrement a semaphore
```sh
scalebox semaphore decrement ${sema_name}
```

- Increment/decrement a semaphore
```sh
scalebox semaphore increment-n {sema_name} ${n}
```

- Semaphore group operations (semagroup)

### 4.2.2 Shared Variables (variable)

Regular variables are usually of string type.

- Create a variable
```sh
scalebox variable create ${var_name}
```
- Read a variable
```sh
scalebox variable get ${var_name}
```
- Write a variable
```sh
scalebox variable set ${var_name} ${value}
```

## 4.2.3 Global Variables

- Shared variables across apps
- Global configuration parameters

## 4.3 State Transition

### 4.3.1 Task Header State Passing
- Task header information passing in the main router
- Task headers starting with an underscore ```_``` are marked as headers to pass along
- Passing path: main router -> algorithm module -> main router
- Passed only once; state must be maintained manually

### 4.3.2 vtask State Retention
- vtask state retention implemented through task header passing
- Task headers starting with ```_vtask_``` are marked as vtask task header attributes
- vtask attribute retention: spans all modules of the vtask
- vtask header attributes: maintained automatically across the vtask lifecycle


## 4.4 vtask Management

> The complete VTask design document (conceptual model, module structure, semaphore mechanism, pipeline flow, command reference) is at :doc:`VTask Virtual Tasks <9_vtask>`. The following is an overview.

The vtask, as the basic unit of the node-local computing programming model, simplifies the implementation of flow control and fault tolerance for local computing in applications through parameter design;

A vtask spans the entire process of input loading, computation, and result write-back; fault tolerance based on vtasks is simple and direct.

A vtask identifies a virtual task spanning multiple modules. It is a count-based flow control mechanism, divided into three types: global count (global_vtask), group count (group_vtask), and node count (host_vtask).

## 4.4.1 Main vtask Functions
- Task management: tasks spanning modules form a vtask, simplifying task management and supporting state management, in turn supporting the node-local computing model;
- Count-based admission control: coordinates with the dynamic scheduling of HPC compute resources to implement very long task computations;
- Coarse-grained fault tolerance: implements fine-grained checkpointing for automatic fault tolerance.

## 4.4.2 vtask Module Structure

- Pre-task queue (wait-queue): not yet bound to compute resources. The basic unit of global fault tolerance.
- vtask head module (vtask-head): the module marking the start of a vtask. Allocates compute resources (single node/resource group) for the vtask; the single-node mode uses HOST-BOUND, running directly on the compute node; the resource group mode uses SLOT-BOUND, with multiple slots deployed on a single node (usually the head node).
- vtask algorithm module (vtask-core): the core processing module of the vtask; in resource group mode, usually identified by a hostname prefix; can be multiple modules, with local distribution between modules implemented through pods.
- vtask tail module (vtask-tail): the module marking the end of a vtask, usually deployed on the head node. The task-body is consistent with vtask-head.


## 4.4.3 vtask Mode Classification

task_dist_mode settings

|            |   Global mode  |  Single-node mode  | Resource group mode |
| ---------- | ----------- | ---------- | ---------- |
| wait-queue |             |            |            |
| vtask-head |             | HOST-BOUND | SLOT-BOUND |
| vtask-core |             | HOST-BOUND | HOST-BOUND |
| vtask-tail |             |            |            |


## 4.5 State Variable Naming Conventions

- Semaphores: character set `[A-Za-z0-9:_-]`, first character is a letter or underscore
- Shared variables: same naming rules as semaphores
- Global variables: same naming rules as semaphores
- vtask semaphore naming patterns:
  - Global: `vtask_size:${mod_name}`
  - Node-level: `host_vtask_size:${mod_name}:${hostname}`
  - Group-level (GROUP-BOUND vtask): `group_vtask_size:${mod_name}:${group_seq}`

## 4.6 Best Practices

- **App decomposition**: the router module is a system-wide centralized module; reducing the number of messages in the design helps batch-process algorithm module semaphores and message sending in message routing, improving system efficiency
- **Semaphore design**: use coarse granularity as much as possible to reduce the scale of management data
- **Efficient main router implementation**: avoid affecting parallel efficiency due to linkage with flow control
- **Batch message processing**: read multiple tasks in one operation and generate subsequent tasks in batches, reducing the number of database operations
- **Batch semaphore operations**: increment/decrement semaphores in batches, reducing database round trips
