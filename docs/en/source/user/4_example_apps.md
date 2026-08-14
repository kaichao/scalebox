# 4. Example Apps

## 4.1 hello-scalebox: An Entry-Level App

The first scalebox app, used to verify installation and basic functionality.

### hello-scalebox Run Example

```bash
cd examples/hello-scalebox
echo "Scalebox" | scalebox run 
```

### hello-scalebox App Features
- Simple message processing
- Single-module app
- Suitable for beginners to understand the basic concepts

## 4.2 app-primes: A Compute-Intensive App

Finds the number of primes in [1..max_value]. Mainly used to demonstrate the key features of scalebox.

### app-primes Run Example

```bash
cd examples/app-primes
make run NUM_GROUPS=4 CALC_NODE=local NUM_PARALLEL=2
```

### app-primes App Features
- Data-parallel processing
- Multi-module collaboration
- Multi-language implementations (Python, Go, C++, etc.)

## 4.3 app-copy: A Data Transfer App

File copy operations are common during computation.

### app-copy Run Example

```bash
cd examples/app-copy
scalebox run --source /data/source --target /data/target
```

### app-copy App Features
- Cross-node data transfer
- Supports multiple transfer protocols
- Pipeline-parallel optimization

## 4.4 remote-primes: A Cross-Cluster App

The cross-cluster version of app-primes, using the compute power of multiple clusters to find the number of primes in a given integer range.

### remote-primes Run Example

```bash
cd examples/remote-primes
make run CLUSTER0=cluster0 CLUSTER1=cluster1
```

### remote-primes App Features
- Cross-cluster computing
- Heterogeneous cluster support
- Dynamic resource allocation

## 4.5 vtask: A Virtual Task App

### vtask Run Example

```bash
cd examples/vtask
scalebox app run --vtasks 100
```

### vtask App Features
- Virtual task management
- Batch task processing
- Resource-optimized scheduling

## 4.6 Design Patterns for Complex Pipeline Apps

### 4.6.1 Multi-Stage Pipeline

```
Data input → Preprocessing → Computation → Postprocessing → Output
```

### 4.6.2 Parallel Pipeline

```
       → Process A →
Input → Process B → Merge → Output
       → Process C →
```

### 4.6.3 Hybrid Mode

```
       → Parallel A →
Input → Serial processing → Parallel B → Output
       → Parallel C →
```

### 4.6.4 Cross-Cluster Pipeline

```
Cluster A: data generation → Cluster B: data processing → Cluster C: result analysis
```

## 4.7 App Best Practices

### 4.7.1 Module Design
- Keep modules stateless
- Support task idempotency
- Set reasonable task timeouts
- Optimize data locality

### 4.7.2 Performance Optimization
- Set parallelism appropriately
- Optimize data distribution
- Reduce communication overhead
- Leverage local storage

### 4.7.3 Fault-Tolerant Design
- Configure retry mechanisms
- Set reasonable timeout values
- Monitor task status
- Handle failed tasks promptly
