# 1. The Scalebox Programming Model


## 1.1 System Architecture



### 1.1.1 System Architecture Diagram
![system architecture diagram](../diagrams/scalebox_arch.drawio.svg)

Figure 1: The scalebox system architecture

### 1.1.2 Layered Architecture

Scalebox uses a layered architecture design, divided into three main layers:

1. **Algorithm Module Layer**
   - Supports multiple programming languages (Python, Go, C++, Java, etc.)
   - Each module is an independent algorithm component or functional unit
   - Containerized packaging for environment isolation and deployment flexibility

2. **Module Connection Layer**
   - Module-to-module connections implemented through shell programming
   - Algorithm execution results are fed back to the system through standard output
   - Information passed includes: subsequent message paths, execution times, I/O data volumes, execution results

3. **Parallel Program Global Layer**
   - YAML configuration files define associations between modules
   - Implements cross-node distributed computing flow control
   - Supports pipeline parallelism and data parallelism modes

### 1.1.3 Core Components

Scalebox's core components include:

- **controld**: a gRPC-based control service that manages actuators and compute nodes
- **actuator**: a launcher service that starts slots on compute nodes via SSH or external schedulers
- **database**: a PostgreSQL database storing metadata for apps, modules, tasks, slots, etc.
- **webui**: the web interface for operating Scalebox
- **agent**: an agent running on compute nodes, responsible for task execution and state management
- **cluster-admin**: a software module for resource acquisition and management that interacts with HPC scheduling systems; also a system-level app

### 1.1.4 Architectural Characteristics

- **Naturally distributed design**: modules are loosely coupled, natively supporting distributed deployment
- **Event-driven execution**: a message-passing communication mechanism
- **Cloud-native design**: containerized packaging, supporting cross-platform deployment
- **Cross-cluster support**: unified scheduling across heterogeneous compute clusters over WANs
- **Cascading support**: agents at the module layer are cascaded, supporting pipeline parallelism

## 1.2 The Programming Model

### 1.2.1 Two-Level Programming Model

Scalebox adopts a two-level programming model, layering complex computing logic:

1. **Algorithm Logic Layer**
   - Focuses on single-node algorithm implementation
   - No need to deliberately consider multi-node parallelism
   - Containerized packaging with standard interfaces

2. **Parallel Logic Layer**
   - Handles business logic related to parallel execution of algorithms
   - Includes resource scheduling, data distribution, task assignment
   - Supports global state management, communication primitives, flow control

### 1.2.2 Event-Driven Architecture

Scalebox uses an event-driven programming model:

```
Event trigger → Task distribution → Module execution → Result feedback → Subsequent task trigger
```

Each module is equivalent to an operation primitive, implementing complex distributed computing logic through the event-driven mechanism. The system distributes input data to the corresponding modules for processing through task queues and schedulers, and the processing results trigger the execution of subsequent modules, forming a complete data processing pipeline.

### 1.1.2 Parallel Program Structure

![parallel_app](../diagrams/parallel_app.drawio.svg)

Figure 2: Traditional Parallel Programs vs. Containerized Parallel Programs


![scalebox_app](../diagrams/scalebox_app.drawio.svg)

Figure 3: Programmable Containerized Programs


## 1.3 Core Concepts

For detailed definitions of core concepts such as App, Module, Task, Slot, Host, and Cluster, see :doc:`User Guide - Core Concepts <../user/2_core_concepts>`.

## 1.4 Design Principles

### 1.4.1 Programmable Containerization

The evolution from traditional containerization to programmable containerization:

- **Traditional containerization**: only packages the application and its dependencies, providing a unified resource-layer interface
- **Programmable containerization**: defines a unified northbound interface specification on top of the algorithm container, driven by external messages

### 1.4.2 The Hourglass Model

Scalebox adopts the hourglass model, concentrating standardized protocols or interfaces in the "narrow waist" to balance system flexibility and complexity:

- Upper layer: complex application-layer protocols
- Middle: the standardized container interface
- Lower layer: flexible underlying implementations

The hourglass model is a commonly used model in the architectural design of network communication and distributed systems. Its design goal is to achieve flexibility, interoperability, and standardization. This design philosophy stems from the need for layered system structures: placing the core protocol or interface as the middle "narrow waist" ensures the independence of the upper and lower layers while enabling unified coordination across layers.

All container technologies follow the OCI (Open Container Initiative) model, making it possible to design independently of application types and of the communication protocols/data representations of distributed systems. However, current container technology serves only as an application packaging technique and does not solve the core protocols and interfaces of distributed software.

![hourglass](../diagrams/hourglass.drawio.svg)

Figure 1: The Hourglass Model of Distributed Computing (Programmable Containerization)

Designing the programmable containerized computing model according to the hourglass model above lets you leverage the model's advantages. Combined with the deployment advantages of container technology, this greatly improves support for distributed software development.

### 1.4.3 Stateless Design

Modules follow the stateless design principle:

- Task execution is stateless and repeatable
- Supports task-level idempotency
- Different logic processing defined through environment variables and task headers

## 1.5 Module Architecture

![module architecture diagram](../diagrams/module_arch.drawio.svg)

- The control plane manages task queues, using an event-driven programming model
- The core of the system is the event-driven task distribution and execution mechanism
- Slots in the data plane obtain tasks through queues, implementing distributed parallel computing
- Flexible and efficient

## 1.6 Slot Architecture

![slot architecture diagram](../diagrams/slot_arch.drawio.svg)

### 1.6.1 The Sidecar Pattern

Each algorithm module runs as a sidecar, decoupled from the system's core services:
- Interacts with the system through standard input and output
- Supports hot swapping and dynamic replacement

### 1.6.2 Slot Management

A Slot is a slice of compute resources on a node corresponding to a Module, the basic unit for fine-grained Task scheduling:
- Each slot is responsible for executing specific task types
- Supports dynamic scaling
- Supports fault recovery and task retries

## 1.7 Parallelization Strategies

### 1.7.1 Intra-Module Parallelism
Code-level parallelism inside the algorithm, such as multithreading and GPU acceleration.

### 1.7.2 Module-Level Parallelism
Multiple instances of the same module process different data in parallel, achieving data parallelism.

### 1.7.3 Inter-Module Parallelism
Different modules execute in parallel in a pipeline fashion, achieving pipeline parallelism.

### 1.7.4 Hybrid Parallelism
Combines the above parallelism methods to achieve optimal parallel efficiency.

## 1.8 The Runtime Environment

### 1.8.1 Control Plane
- controld: the core control service
- actuator: slot startup and management
- database: metadata storage

### 1.8.2 Data Plane
- agent: the task execution agent
- slot: the unit of compute resources
- message bus: inter-module communication

### 1.8.3 Communication Mechanisms
- gRPC: control plane communication
- message queues: data plane communication
- database: state persistence
