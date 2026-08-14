# 5. Node-Local Computing

Node-Local Compute is a computing optimization technique for high-I/O data-intensive applications. By moving computational data to the local storage of compute nodes, it avoids shared storage system bottlenecks, greatly improving computing efficiency.

## 5.1 Technical Background and Requirements

### 5.1.1 Challenges of High-I/O Applications

Traditional HPC (High-Performance Computing) frameworks struggle to meet the requirements of high-I/O applications, which include:

- **Astronomical computing**: large-scale astronomical image processing, radio telescope data processing, with data scales at the PB/EB level
- **Bioinformatics data processing**: genome sequencing analysis, protein structure prediction
- **Digital video processing**: HD video encoding/decoding, video content analysis
- **Medical image analysis**: medical imaging processing (CT, MRI, X-ray)
- **AI training scenarios**: large-scale training data loading for bioinformatics AI and medical AI

A common characteristic of these applications is that **their I/O bandwidth requirements far exceed the data loading capability of shared storage systems under existing computing frameworks**. When compute nodes need to access data through parallel file systems, the physical bandwidth limit of shared storage becomes the bottleneck of the entire computing process.

### 5.1.2 The Storage-Compute Ratio Imbalance Problem

- **Definition of the storage-compute ratio**
  - Storage-compute ratio = access bandwidth / computing power, reflecting the efficiency balance between storage and computation during the computing process.
    - Access bandwidth (storage power): the ability to access data from storage hardware (such as memory, hard disks), usually measured in GB/s or MB/s.
    - Computing power:
      - Usually measured in floating-point operations per second (FLOPs)
      - For data processing applications, can be measured in data processing bandwidth (GB/s or MB/s)
  - **Balance of the storage-compute ratio**
    - An imbalanced storage-compute ratio leads to low overall computing efficiency
      - Storage bottleneck: insufficient storage bandwidth causes idle computing resources.
      - Computing bottleneck: insufficient computing power cannot process data in time.
    - **Cost of storage and compute hardware**
      - Usually, the unit cost of storage power is far lower than that of computing power; computing power cost is the most important factor affecting total computing cost
      - Storage power exceeding computing power increases storage costs but helps little with overall computing efficiency

## 5.2 Storage Hierarchy and Performance Analysis

### 5.2.1 Storage Hierarchy of HPC Systems

The storage hierarchy of modern computing systems, from fast to slow, includes:

1. **Local CPU cache**: nanosecond-level access latency, KB~MB capacity
2. **Local memory cache and the tmpfs memory file system**: hundreds of nanoseconds latency, GB-level capacity
3. **Local external storage**: HDD/SATA SSD/NVMe SSD, microsecond~millisecond latency
4. **Cluster storage (parallel file system)**: millisecond latency, TB~PB capacity

### 5.2.2 Data Access Path Comparison

#### Traditional HPC Application Access Path
- Data access path based on the shared parallel file system
- Total bandwidth shared by all applications: typical HPC storage systems have read/write bandwidth of about 100GB/s (tens to hundreds of GB/s)
- As the number of nodes increases, the effective bandwidth available to each node drops rapidly

#### Access Path Based on Node-Local Storage
- The aggregate bandwidth of node-local storage far exceeds parallel file system bandwidth and is more stable
- Each node independently accesses its own local storage without shared contention

### 5.2.3 Storage Media Performance Comparison

- **SATA SSD**: 500MB/s (single device)
- **NVMe SSD**: 5GB/s (single device)
- **Memory tmpfs**: 6GB/s (memory bandwidth)

### 5.2.4 Aggregate Bandwidth Advantage Analysis

Estimated local bandwidth table for compute cluster nodes (single thread/single component)

|  Nodes |  SATA SSD  |  NVMe SSD   |     tmpfs   |
| ------ | ---------- | ----------- | ----------- |
|     10 |     5 GB/s |     50 GB/s |     60 GB/s |
|     30 |    15 GB/s |    150 GB/s |    180 GB/s |
|    100 |    50 GB/s |    500 GB/s |    600 GB/s |
|    300 |   150 GB/s |  1,500 GB/s |  1,800 GB/s |
|  1,000 |   500 GB/s |  5,000 GB/s |  6,000 GB/s |
|  3,000 | 1,500 GB/s | 15,000 GB/s | 18,000 GB/s |
| 10,000 | 5,000 GB/s | 50,000 GB/s | 60,000 GB/s |

**Key observation**: in large-scale clusters, the aggregate bandwidth of node-local storage can reach tens to hundreds of times that of shared storage. Yet the shared storage bandwidth of a typical computing cluster usually does not exceed 100GB/s

## 5.3 Core Principles of Node-Local Computing

### 5.3.1 The Basic Idea

The core idea of node-local computing is **moving computational data to the local compute node** to avoid shared storage bottlenecks. By leveraging each compute node's own storage resources (memory, SSD, etc.), it builds a distributed, high-bandwidth data access environment.

### 5.3.2 The Three-Stage Processing Model

1. **Data loading stage**: read data from the parallel file system into node-local storage
2. **Local computation stage**: read and write data on node-local storage to complete the algorithm function
3. **Result write-back stage**: write computation results from local storage back to the parallel file system

The above steps support pipelined parallel operation, i.e., different data blocks are processed simultaneously at different stages, maximizing system throughput.

### 5.3.3 Key Technical Optimizations

- **Data packing**: pack multiple small files into large files, reducing file system metadata overhead
- **Data compression**: reduce the amount of data transferred and stored, improving I/O efficiency
- **Pipeline parallelism**: hide the latency of data loading and write-back, keeping the computation stage running continuously

## 5.4 Key Technical Challenges

### 5.4.1 Admission Control Mechanisms

- **Admission control between preceding and following modules**
  - Match running speeds to avoid local disk/memory overflow
  - Semaphores are an important technical means to implement admission control
  - Capacity detection for local disk partitions
  - Custom scripts for flexible admission control
- **Synchronized admission control for the same module across multiple nodes**
  - Unify overall progress in multi-node parallel synchronization scenarios
  - Processing of different channels / data downloads on different nodes
  - Differences in the computational volume of pending data and node processing capabilities can lead to different computation progress, in turn affecting intermediate data accumulation and memory overflow

### 5.4.2 Fault Handling

- **Automatic fault handling**
  - Errors affect the processing of part of the data
  - Admission control can stop the entire pipeline
  - The impact of automatic fault tolerance on admission control

### 5.4.3 Resource Management and Scheduling

- **Task processing ordering for key modules**: ensure critical-path tasks execute first
- **Decomposition of non-key modules**
  - Subdivision increases the admission control burden
  - Increasing pipeline parallelism improves performance but increases memory usage for intermediate storage, affecting pipeline stability

## 5.5 Design Approach in Scalebox

### 5.5.1 Task Distribution Modes

Scalebox supports node-local computing through task distribution modes:

- **HOST-BOUND**: tasks bound to a specific host, specified via `headers->>'to_host'`
- **SLOT-BOUND**: tasks bound to an execution slot, specified via `headers->>'to_slot'`; can be used for task distribution of vtask head nodes
- **GROUP-BOUND**: tasks bound to a group; the task's `headers->>'to_host'` is specified via the slot's `parameters->>'group_prefix'`

### 5.5.2 Location-Aware Scheduling

- **pod_id mechanism**: identifies modules managed by pods; if the pod that a task originates from has the same pod_id, all tasks are marked as using local computing (task_dist_mode is HOST_BOUND)
- **Location information in task headers**: supports local cascaded processing between preceding and following modules

### 5.5.3 The vtask Support Framework

- The vtask serves as the basic execution unit of the node-local computing model
- Through parameter design, simplifies the implementation of admission control and fault tolerance for local computing in applications
- Implements idempotency at the cross-module vtask level, ensuring eventual consistency

### 5.5.4 The Semaphore Mechanism

- Centrally manages the running state of the entire system, making compute modules stateless
- Usually read and written only in the main-router to avoid concurrency problems
- Regular algorithm modules can read the corresponding values for admission control decisions

## 5.6 Hardware and Platform Requirements

### 5.6.1 Hardware Infrastructure

- **Genuine local storage support is required**: virtual machines in the cloud are usually not suitable for node-local computing because their storage is often network-attached
- **Self-built server clusters (without virtualization)**: physical servers directly access local storage
- **HPC compute nodes**: compute nodes in traditional HPC environments

### 5.6.2 Applicable Scenario Summary

Node-local computing is particularly suitable for the following scenarios:

1. **Data-intensive scientific computing**: astronomy, bioinformatics, climate simulation, etc.
2. **Large-scale media processing**: video encoding/decoding, image analysis
3. **AI training and inference**: scenarios requiring fast loading of training datasets
4. **Any application whose I/O bandwidth requirements exceed shared storage capabilities**

### 5.6.3 Performance Expectations

- **Small-scale clusters (10-100 nodes)**: I/O performance improved by 5-50x
- **Medium-scale clusters (100-1000 nodes)**: I/O performance improved by 50-500x
- **Large-scale clusters (1000+ nodes)**: I/O performance improved by over 500x

## 5.7 Summary

Node-local computing technology solves the shared storage system bottleneck in traditional HPC frameworks by redesigning the data access path and fully leveraging the high-bandwidth characteristics of node-local storage. This technology applies not only to astronomical computing but also broadly to bioinformatics processing, medical image analysis, digital video processing, AI training, and many other high-I/O application scenarios.

The Scalebox platform provides a complete implementation framework for node-local computing through technical means such as task distribution modes, location-aware scheduling, vtask support, and the semaphore mechanism. In actual deployments, the hardware infrastructure must provide genuine local storage support to fully realize the technical advantages of node-local computing.

## 5.8 Performance Optimization Strategies

### 5.8.1 Pipeline Design Optimization
- Batch processing and pipeline optimization to reduce intermediate waiting time

### 5.8.2 Main Router Design Optimization
- Parallel logic optimization: multiple main router instances execute in parallel
- Coarse-grained task design to reduce the number of tasks and lower scheduling overhead

### 5.8.3 Core Algorithm Optimization
- I/O performance optimization strategies to reduce the number of disk read/write operations
- Fully leverage the node-local computing mode

### 5.8.4 Compute Skew Handling
- Local computing + grouped compute nodes to balance load distribution
