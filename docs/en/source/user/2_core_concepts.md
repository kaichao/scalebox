# 2. Core Concepts

## 2.1 Basic Architecture

Scalebox applications are organized in layers. The module layer encapsulates algorithms as standard containers and provides an event-driven container-layer interface; the application layer implements task distribution through the main router module and supports complex parallel logic.

## 2.2 Core Concepts in Detail

| Concept | English | 中文 | Description |
| :------ | :--------- | :------ | :------------------------------------------------------ |
| Application | **App**    | **应用** | An application that performs a specific computing task, containing multiple Modules. An App can dynamically acquire the compute resources it needs by invoking resource scheduling systems (such as Slurm). |
| Software unit | **Module** | **模块** | The basic programmable unit that makes up an App, a collection of program code packaged in a container. Usually an independent algorithm component, or a component for transfer between nodes. Components form data processing pipelines through cascading. |
| Execution instance | **Task**   | **任务** | The basic execution unit — the process of running input data on a specific Module. |
| Compute node | **Host**   | **节点** | The server node that executes computing tasks. |
| Resource unit | **Slot**   | **插槽** | A slice of compute resources on a node corresponding to a Module, the basic unit for fine-grained Task scheduling |
| Virtual task | **VTask**  | **任务组** | A cross-module collection of Tasks, an application-level coarse-grained compute unit. With built-in semaphores and shared variables, it forms a pipeline of wait-queue, vtask-head, vtask-core, and vtask-tail for unified flow control, resource binding, and state tracking. |
| Resource collection | **Cluster** | **集群** | A logical grouping of compute resources containing multiple Hosts. Supports single-node clusters, static clusters, dynamic clusters, inline clusters, and other types. Cross-cluster computing is implemented through full-mesh replication of t_cluster + gRPC proxy. |

## 2.3 Application (App)

### Application Composition
- **Main router module**: the system main program, responsible for distributing tasks to custom algorithm modules and common base modules, and maintaining the running state of the entire app.
- **Custom algorithm modules**: user business modules that implement specific algorithm functions in a message-driven form.
- **Common base modules**: base function modules such as file transfer and directory listing.

### Application Definition
Applications are defined through YAML configuration files, including:
- Application name and identifier
- Module definitions and configurations
- Parameter settings
- Cluster configuration

## 2.4 Module

### Module Types
1. **Algorithm modules**: user-defined computing logic
2. **Transfer modules**: data transfer and file operations
3. **Router modules**: message routing and task distribution
4. **Helper modules**: directory listing, scheduled tasks, etc.

### Module Characteristics
- **Containerized packaging**: independent runtime environment, environment isolation
- **Message-driven**: receives task messages through standard input
- **Stateless design**: supports repeated execution and fault tolerance
- **Multi-language support**: supports Python, Go, C++, Java, Shell, and other languages

## 2.5 Task

### Task Lifecycle
1. **Created**: task generated when the app is created
2. **Ready**: task is ready and waiting to execute
3. **Running**: task is executing
4. **Completed**: task executed successfully
5. **Failed**: task execution failed (retryable)

### Task Format
- **Task body**: task input data, JSON format or a string without null characters
- **Task headers**: task metadata, a JSON format string
- **Task ID**: unique identifier

## 2.6 Slot

- **Slot functions**
  - **Resource allocation**: a slice of compute resources on a node
  - **Task execution**: the execution environment for specific tasks
  - **State management**: manages the execution state of tasks


## 2.7 Cluster

### Cluster Types
1. **Single-node cluster**: for testing and development; a single node serves as both the management node and the compute node
2. **Inline cluster**: includes one or more head nodes and several compute nodes
3. **External cluster**: a compute cluster managed by an external resource scheduling system (such as Slurm)

### Cluster Components
- **Head node**: runs control services such as controld, actuator, and database
- **Compute nodes**: execute specific computing tasks
- **Shared storage**: optional, for data sharing between nodes

## 2.8 Virtual Tasks (VTask / Task Groups)

VTask is an application-level coarse-grained compute unit that Scalebox establishes on top of fine-grained Tasks. A VTask is a cross-module collection of Tasks that solves three core problems:

1. **Cross-module association**: when multiple Modules within the same App collaborate to complete a piece of work, the Tasks scattered across Modules are linked through the `_vtask_id` header
2. **Flow control and resource binding**: semaphores limit the number of concurrent VTasks, supporting binding VTasks to specific compute nodes or node groups
3. **State tracking and fault tolerance**: `vtask get` / `vtask list` show progress in real time, and `vtask fail` terminates the pipeline and cascades cleanup

### Three Modes

| Mode | Resource binding | Pipeline structure | Applicable scenarios |
|------|---------|---------|---------|
| DEFAULT | None | task → vtask-head → vtask-core → vtask-tail | Lightweight batch processing |
| HOST-BOUND | Single node | task → wait-queue → vtask-head → vtask-core → vtask-tail | Exclusive computing on a single node |
| GROUP-BOUND | Node group | task → wait-queue → vtask-head → vtask-core-1 → vtask-core-2 → ... → vtask-tail | Multi-node collaborative computing |

> See :doc:`VTask Virtual Tasks <../developer/9_vtask>` for the detailed design.

## 2.9 Clusters in Detail

### Cluster Classification

| Type | Purpose | Description |
|------|------|------|
| Single-node cluster | Testing, development | A single node serves as both management node and compute node |
| Static cluster | Production | Head node + fixed compute nodes |
| Dynamic cluster | Elastic computing | Head node + external scheduler (Slurm) dynamically allocates nodes |
| Inline cluster | Hybrid mode | Head node + dynamically acquired compute nodes |

### Cross-Cluster Computing

Scalebox supports unified scheduling across heterogeneous compute clusters over WANs:
- The `t_cluster` table is replicated in full mesh across all cluster databases; every controld can query the gRPC address of any cluster
- The CLI always connects to the local controld; cross-cluster operations are forwarded by server-side gRPC proxies
- Cross-cluster App references are established through `RegisterRemoteAppLink`; `t_app_rlink` + `t_module_rvlink` manage module cascade routing

> See :doc:`Cross-Cluster Architecture <../developer/10_cross_cluster>` for the detailed design.

## 2.10 Parallelization Methods

### Intra-Module Parallelism
Code-level parallelism inside the algorithm, such as multithreading and GPU acceleration.

### Module-Level Parallelism
Multiple instances of the same module process different data in parallel, achieving data parallelism.

### Inter-Module Parallelism
Different modules execute in parallel in a pipeline fashion, achieving pipeline parallelism.
