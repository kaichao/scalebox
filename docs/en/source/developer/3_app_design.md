# 3. App Design

## 3.1 Scalebox Applications

### 3.1.1 App Characteristics
- Containerized computing programs optimized for data-intensive scenarios
- Rely on the Scalebox execution framework, implementing multi-node distributed, parallel computing in an event-driven form
- Mainly applied to astronomical computing scenarios; in the future, applicable to bio-data processing, digital film production, video processing, and other high-data-I/O scenarios

### 3.1.2 Comparison with Traditional MPI Programming

| Feature | Scalebox | Traditional MPI/Charm++ |
| -------   | ------------------  | -------------------------|
| Programming paradigm | Message-driven, over-decomposed | Process communication, single program multiple data (SPMD) |
| Load balancing | Dynamic, automated | Mostly subjective, static |
| Fault tolerance | Built-in support | Must be implemented separately |
| Scalability | Stream-driven, supports extremely large-scale parallelism | Faces challenges at large scale |
| Programming complexity | High-level abstraction, low | Low-level control, high |

## 3.2 App Composition

### 3.2.1 Module Composition
Scalebox applications are composed of multiple containerized software modules:

- **Main router module**: the system main program, responsible for distributing tasks to custom algorithm modules and common base modules, and maintaining the running state of the entire app
- **Custom algorithm modules**: user business modules that implement specific algorithm functions in a message-driven form
- **Common base modules**:
  - File transfer: supports cross-node data transfer based on the computing system's built-in ssh
  - Directory listing: helps generate files to process

### 3.2.2 App Definition Files
- **App template file**: defines the modules included in the app in YAML file format
- **Parameter definition file**: defines parameters related to the app module file in environment variable file form

## 3.3 App Design Principles

### 3.3.1 The Main Router Uniformly Manages Business Logic
- Handles only state management and task distribution, not specific business, for efficient implementation
- Runs serially to avoid concurrency problems

### 3.3.2 Algorithm Modules Are Stateless
- Flexible invocation. Different logic processing defined through environment variables and task headers
- Stateless design; implement task idempotency at the module level as much as possible

### 3.3.3 Divide Modules by Logical Function
- Subdivided modules support pipeline parallelism
- Dividing by function improves reusability

### 3.3.4 Node-Local Computing Mode
- Improves I/O bandwidth
- Reduces intermediate storage; can generate partial results and produce final results as early as possible

### 3.3.5 Reuse Standard Modules
- Use existing standard modules as much as possible
- Reduce duplicated development work

## 3.4 App Definition

For the complete specification of the app template file (app.yaml) and parameter file (scalebox.env), see :doc:`Technical Specification Reference <../appendix/1_tech_specifications>`.

## 3.5 App Design Patterns

### 3.5.1 Simple Pipeline Pattern
```
Data input → Processing module → Output results
```

Characteristics:
- Single-module processing
- Suitable for simple computing tasks
- Easy to understand and debug

### 3.5.2 Multi-Stage Pipeline Pattern
```
Data input → Preprocessing → Computation → Postprocessing → Output
```

Characteristics:
- Multi-module collaboration
- Data processed step by step
- Supports complex computing flows

### 3.5.3 Parallel Pipeline Pattern
```
       → Process A →
Input → Process B → Merge → Output
       → Process C →
```

Characteristics:
- Module-level parallel processing
- Improves processing efficiency
- Suitable for data-parallel tasks

### 3.5.4 Hybrid Parallel Pattern
```
       → Parallel A →
Input → Serial processing → Parallel B → Output
       → Parallel C →
```

Characteristics:
- Combines serial and parallel
- Flexible task scheduling
- Suitable for complex application scenarios

### 3.5.5 Cross-Cluster Pattern
```
Cluster A: data generation → Cluster B: data processing → Cluster C: result analysis
```

Characteristics:
- Cross-cluster computing
- Heterogeneous cluster support
- Suitable for large-scale distributed computing

## 3.6 App Design Best Practices

### 3.6.1 Module Division Principles
1. **Single responsibility**: each module is responsible for only one function
2. **High cohesion, low coupling**: modules are tightly related internally and loosely coupled externally
3. **Reusability**: design generic modules to improve reuse
4. **Extensibility**: support dynamic addition and replacement of modules

### 3.6.2 Data Flow Design
1. **Data locality**: prefer local storage and reduce data transfer
2. **Pipeline parallelism**: arrange module order appropriately to maximize parallelism
3. **Batch processing**: set batch sizes appropriately, balancing latency and throughput
4. **Fault-tolerant design**: configure retry mechanisms and error handling

### 3.6.3 Performance Optimization Strategies
1. **Parallelism settings**: set parallelism appropriately based on data volume and compute resources
2. **Resource allocation**: allocate compute resources appropriately based on task requirements
3. **Caching strategies**: use caching appropriately to reduce redundant computation
4. **Load balancing**: ensure balanced loads across nodes to avoid resource waste

### 3.6.4 Maintainability Design
1. **Modular design**: modularize functions for easier maintenance and upgrades
2. **Configuration-based management**: manage apps through configuration files to reduce maintenance costs
3. **Monitoring and alerting**: integrate monitoring and alerting mechanisms to detect and solve problems promptly
4. **Complete documentation**: provide detailed documentation and usage instructions

