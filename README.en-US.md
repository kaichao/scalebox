# Scalebox - A cloud-native streaming computing engine

Scalebox is a cloud-native streaming computing engine that can run containerized stand-alone user algorithms on distributed and heterogeneous computing clusters, organize large-scale parallel processing at the module level with pipelines, and support task-level fault tolerance. Compared with existing frameworks such as big data processing and parallel computing, its technical characteristics are especially suitable for application scenarios such as data distribution, computing power resource distribution, and complex algorithms.

Scalebox has the following important properties:

- **Cloud-native design**: Encapsulate all algorithm modules and transmission modules in containers, and embed them into the data processing pipeline designed for the cloud environment through the sidecar model; the software basic platform is completely based on the cloud-native design; the platform will control messages, The data channel is separated, and the front and back modules are connected by a message bus to realize non-intrusive parallel programming, which greatly simplifies the difficulty of parallel computing.
- [**Cross-cluster computing**](./tests/cross-cluster-primes/): Normalized processing algorithm module, transmission module, unified processing of intra-cluster/cross-cluster data through the pipeline, shielding data and computing Cross-cluster differences; cross-cluster message-driven stream processing supports the deployment of a single pipeline application on multiple heterogeneous computing power clusters.
- [**Task-level fault tolerance**](./tests/retry_test/): For sporadic errors caused by hardware failures, software bugs, network problems, data anomalies, etc., automatic fault-tolerant processing is implemented based on rules; fine-grained task-level Fault tolerance for trusted data analytics on unreliable hardware.
- **Location-aware scheduling**: The IP address of the sender can be configured in the message body, which supports local cascading processing between front and rear modules; separates the processing layer and storage layer to reduce coupling, and messages in the horizontal direction (processing layer) drive the vertical direction Data reading and writing (from the processing layer to the storage layer) reduces the east-west network traffic in the computing cluster and eliminates the I/O bottleneck of the cluster storage; then realizes local computing without shared storage, pure local loading of large files, and effectively supports horizontal expansion.
- [**Task Perspective**](./tests/task-perspective/): The computing task is a message-driven process, and the task perspective records the detailed running status of each task in detail; including user program return code, standard output, Standard error, program custom text, number of bytes read and written by user programs, etc., also includes various system-level and user-defined time stamps on the computing container side and control side during the task execution cycle (message generation/distribution/processing, result recording) . Task perspective provides basic support for precise positioning, application optimization, and data statistics in troubleshooting.

- **Multiple parallelization methods**
  - Algorithm parallelism within the module
  - Module-level data parallelism
  - Cross-module pipeline parallelism

- **multi-computing backend**
  - Multiple types of computing clusters (self-managed clusters, HPC clusters, k8s container clusters, etc.)
  - Various container engines
    - docker: the default container engine
    - [singularity](./tests/hello-scalebox-singularity/)
    - k8s: TODO

This repository contains:

1. Scalebox server environment based on docker-compose ([Service Environment](./server/README.md))
2. Dockerfile definition for scalebox standard modules([standard module](./dockerfiles/README.md))
3. Application example of scalebox ([Application Example](./examples/README.md))
4. Test of the main features of scalebox ([feature test](./tests/README.md))


## Table of Contents

- [Scalebox - A cloud-native streaming computing engine](#scalebox---a-cloud-native-streaming-computing-engine)
  - [Table of Contents](#table-of-contents)
  - [Background](#background)
  - [Install](#install)
  - [Usage](#usage)
  - [Examples](#examples)
  - [Feature tests](#feature-tests)
  - [Related Softwares](#related-softwares)
  - [Maintainers](#maintainers)
  - [Contributing](#contributing)
  - [License](#license)

## Background

## Install

## Usage

## Examples

## [Feature tests](./tests)

- [Fault tolerance Setup](tests/retry_test/)
- [Timeout-setup](tests/timeout-gen/)
- [Flow control management](tests/check_test/)
- [Task-perspective](tests/task-perspective/)
- [Cross-cluster-computing](tests/cross-cluster-primes/)
- [singularity](tests/hello-scalebox-singularity/)

## Related Softwares

- [PostgreSQL Database Management System](https://github.com/postgres/postgres) — Scalebox backend database
- [gRPC – An RPC library and framework](https://github.com/grpc/grpc) — Efficient communication protocol between different software modules
- [The Go Programming Language](https://github.com/golang/go) — Programming language for cloud-native applications
- [Pony ORM ER Diagram Editor](https://editor.ponyorm.com/) - Magical ER Diagram Tool


## Maintainers

[@Kaichao](https://github.com/kaichao).

## Contributing

Feel free to dive in! [Open an issue](https://github.com/kaichao/scalebox/issues/new) or submit Pull Requests.

## License

[Apache](LICENSE) © Kaichao Wu
