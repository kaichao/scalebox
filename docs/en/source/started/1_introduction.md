# 1. Introduction to Scalebox

## 1.1 What is Scalebox?

Scalebox is a programmable containerized computing framework, comprising a parallel programming model for containerized computing and a runtime environment that supports efficient data-intensive computing.

Through "operator-level spatial expansion + local storage residency + explicit data-flow orchestration", Scalebox moves a large portion of the communication that would otherwise occur between nodes into the node itself, significantly reducing dependence on inter-node communication bandwidth and external storage bandwidth.

### 1.1.1 From Traditional Containerization to Programmable Containerization

- Traditional applications: application dependencies and library files must be installed on the compute node;
- Containerized applications: dependencies, library files, and the application itself are packaged together into a container, simplifying software deployment; the container's unified resource-layer interface (southbound interface) connects to the underlying compute platform.
- Programmable containerization: building on containerized applications, standard interfaces are defined to form standard modules. Modules run when driven by external messages. A unified application interface (northbound interface) specification is defined on top of the algorithm container.

![programmable_container](../diagrams/programmable_container.drawio.svg)

Figure 1: Traditional Containerization vs. Programmable Containerization

- Traditional containerization: encapsulation connects to underlying compute resources through the unified container-layer downward interface, without standardizing the upper-layer application interface; the interfaces called externally are microservice interfaces provided by the algorithm code inside the container, in different forms such as gRPC and REST.
- Programmable containerization: building on containerization and the algorithm code's service interfaces, a message-driven unified container-level programming interface (northbound interface) specification is defined for upper-layer applications, forming programmable, standard algorithm containers.


### 1.1.2 Node-Level Computing-Storage Integration Architecture

| Traditional HPC Architecture | Scalebox          |
| ------------- | ----------------- |
| Centered on HPC parallel storage | Local SSD/memory cache |
| Storage bandwidth drives performance | Results drive performance |
| Runtime scheduling | Pipeline topology priority scheduling |
| Time-multiplexing model centered on global storage | Data-flow model centered on local storage/memory cache + spatial expansion |

Scalebox encapsulates algorithms in containers and provides programmable control, supporting stream computing on distributed, cross-WAN heterogeneous compute clusters. Computing tasks are driven by data flow under a non-von Neumann architecture, achieving ultra-large-scale efficient parallelism.


## 1.2 Key Features of Scalebox

### 1.2.1 Cloud-Native Design
Algorithm implementations support multiple programming languages; cross-platform deployment via containerized encapsulation; modular design based on the sidecar pattern.

### 1.2.2 Non-Intrusive Programming
User algorithms can be containerized without modification, interacting with the platform through standard interfaces, separating parallel logic from algorithm logic.

### 1.2.3 Cross-Cluster Deployment
Supports unified scheduling across heterogeneous compute clusters over WANs, enabling efficient computing in scenarios with distributed data and distributed compute resources.

### 1.2.4 Node-Local Computing Optimization
Centered on node-local storage, reducing dependence on global storage and improving I/O efficiency through data locality.

### 1.2.5 High Parallel Efficiency
Supports multi-level parallelism: intra-module algorithm parallelism, module-level data parallelism, and inter-module pipeline parallelism.

### 1.2.6 Virtual Tasks (VTask / Task Groups)
Cross-module Task collections with built-in semaphores and shared variables, managed by the head→core→tail pipeline for unified flow control, resource binding, and state tracking. Supports DEFAULT / HOST-BOUND / GROUP-BOUND modes.

### 1.2.7 Multiple Interfaces
- **WebUI**: REST gateway (:8088), SPA frontend, full-function management of Dashboard / App / DAG / Task / VTask
- **VS Code extension**: sidebar TreeView, task log Webview, DAG topology, app.yaml syntax support
- **CLI**: 18 command groups (including validate) under the scalebox main command, 24 leaf commands in total covering all operations, supporting scripting and CI/CD integration

## 1.3 Core Value

- **Simplified parallel programming**: separates complex parallel logic from algorithm logic through a two-level programming model, lowering the difficulty of parallel program development.
- **Improved computing efficiency**: significantly improves the execution efficiency of data-intensive applications through node-local computing optimization and multi-level parallelism.
- **Enhanced system reliability**: built-in task-level fault tolerance supports trustworthy data analysis on unreliable hardware.
- **Heterogeneous environment support**: cross-cluster, cross-platform deployment capabilities fully leverage distributed heterogeneous compute resources.
- **Lower operational barriers**: WebUI visualization + VS Code extension cover the full workflow of development, debugging, and monitoring.
