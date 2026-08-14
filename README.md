# Scalebox — A Cloud-Native Stream Computing Engine

Scalebox is a cloud-native stream computing engine that runs containerized user algorithms on distributed, heterogeneous computing clusters, enabling large-scale parallel processing through pipelines of hierarchical modules with task-level fault tolerance. Compared with existing big data processing and parallel computing frameworks, its technical characteristics are particularly suited to scenarios with distributed data, distributed compute resources, and complex algorithms.

## ✨ Key Features

- **Cloud-Native Design**: all algorithm modules are containerized and embedded into data processing pipelines via the sidecar pattern; control messages and data channels are separated, and adjacent modules are linked through a message bus, enabling multi-language, non-intrusive parallel programming.

- **Virtual Tasks (VTask / Task Groups)**: an application-level coarse-grained compute unit built on top of fine-grained Tasks. Cross-module Task collections incorporate semaphores and shared variables, managed by the wait-queue → vtask-head → vtask-core → vtask-tail pipeline for flow control, resource binding, and state tracking. Supports DEFAULT / HOST-BOUND / GROUP-BOUND modes, covering everything from lightweight batch processing to multi-node collaborative computing.

- **Cross-Cluster Computing**: algorithm modules and transfer modules are normalized and processed uniformly through pipelines for intra-cluster and cross-cluster data. Full-mesh replication of t_cluster plus gRPC proxy forwarding; the CLI always connects to the local controld, and cross-cluster operations are transparently routed by the server.

- **Task-Level Fault Tolerance**: automatic retries based on exit-code rules for transient errors caused by hardware failures, software bugs, network issues, or data anomalies. Fine-grained task-level fault tolerance enables trustworthy data analysis on unreliable hardware.

- **Node-Local Computing Optimization**: centered on node-local storage, using "operator-level spatial expansion + local storage residency + explicit data-flow orchestration" to move inter-node communication into the node itself, significantly reducing dependence on network bandwidth and external storage.

- **WebUI Management Interface**: REST gateway (:8088) with 76 endpoints covering all management functions. SPA frontend embedded as a single-file deployment, providing Dashboard statistics, App card grid, DAG topology, Task/VTask management, and semaphore/variable hierarchy trees.

- **VS Code Extension**: sidebar TreeView (three-level lazy-loading tree of App → Modules/Tasks/VTasks), task log Webview, DAG force-directed graph (ECharts), app.yaml LSP (syntax validation/autocompletion/hover hints). Packaged with Docker, zero local Node.js dependencies.

- **Security Framework**: JWT authentication + RBAC role system (admin/viewer/operator/automation) + gRPC TLS encryption + database certificate authentication. Compose-tiered optional enablement, disabled by default with zero overhead.

- **Multi-Level Parallelism**: intra-module algorithm parallelism, module-level data parallelism, and inter-module pipeline parallelism.

- **Multiple Container Engines**: Docker (default), Singularity/Apptainer, Podman.

## 🚀 Quick Start

### 1. Prepare the Environment

```bash
# Docker 20.10+
curl -fsSL https://get.docker.com | sh
sudo systemctl enable --now docker
docker --version
```

### 2. Start Scalebox

```bash
git clone https://github.com/kaichao/scalebox.git
cd scalebox

# Generate secrets (optional)
cd build && bash gen-secrets.sh && cd ..

# Build and start
make -C build/
docker compose -f build/compose.yaml up -d
```

Open `http://localhost:8088` in your browser to access the WebUI.

### 3. Run an Example

```bash
cd examples/hello-scalebox
echo "Hello Scalebox" | scalebox run
scalebox app list
```

### 4. Install the VS Code Extension

```bash
cd vscode-scalebox && make install
```

## 📊 Core Concepts

| Concept | Name    | Description |
|:---|:---|:---|
| App | Application | An application that performs a specific computing task, consisting of multiple Modules |
| Module | Module | A containerized algorithm component that forms pipelines through cascading |
| Task | Task | The basic execution unit — the process of running input data on a Module |
| VTask | Virtual Task | A cross-module Task collection with built-in semaphores and variables for unified flow control and state management |
| Host | Host | The server that executes computing tasks |
| Slot | Slot | A slice of compute resources on a node, the basic unit for Task scheduling |
| Cluster | Cluster | A logical grouping of compute resources, supporting collaboration across wide-area networks |

## 📚 Documentation

### Chinese User Documentation ([→ Full TOC](docs/cn/source/index.rst))

| Category | Document | Description |
|------|------|------|
| Getting Started | [Introduction to Scalebox](docs/cn/source/started/1_introduction.md) | Product positioning, key features, core value |
| Getting Started | [Quick Start](docs/cn/source/started/2_quick_start.md) | Environment preparation, deployment, installation verification |
| User Guide | [Installation](docs/cn/source/user/1_installation.md) | Docker Compose single-node / multi-node clusters |
| User Guide | [Core Concepts](docs/cn/source/user/2_core_concepts.md) | App/Module/Task/VTask/Host/Slot/Cluster in detail |
| User Guide | [Running Apps](docs/cn/source/user/3_running_apps.md) | Creating, managing, and monitoring apps and tasks |
| User Guide | [Example Apps](docs/cn/source/user/4_example_apps.md) | From entry-level to advanced examples |
| User Guide | [Operations](docs/cn/source/user/5_operations.md) | Monitoring, troubleshooting, backup, cluster management |
| User Guide | [Standard Modules](docs/cn/source/user/6_standard_modules.md) | File transfer, scheduled tasks, directory listing, etc. |
| User Guide | [WebUI Guide](docs/cn/source/user/7_webui.md) | REST gateway, page navigation, build and deployment |
| User Guide | [VS Code Extension](docs/cn/source/user/8_vscode.md) | Installation and configuration, TreeView, DAG, LSP |
| Developer Guide | [Programming Model](docs/cn/source/developer/1_programming_model.md) | Two-level programming model, event-driven architecture |
| Developer Guide | [Module Development](docs/cn/source/developer/2_module_development.md) | Module design, sidecar pattern, unit testing |
| Developer Guide | [App Design](docs/cn/source/developer/3_app_design.md) | Design principles, pipeline patterns, best practices |
| Developer Guide | [Main Router & Status](docs/cn/source/developer/4_main_router_status.md) | Semaphores/variables/global variables, VTask overview |
| Developer Guide | [Node-Local Computing](docs/cn/source/developer/5_node_local_compute.md) | Computing-storage integration, three-stage model, performance analysis |
| Developer Guide | [Advanced Features](docs/cn/source/developer/7_advanced_features.md) | Fault tolerance, admission control, timeouts, slot auto-scaling |
| Developer Guide | [VTask Design](docs/cn/source/developer/9_vtask.md) | Conceptual model, module structure, semaphores, pipeline flow |
| Developer Guide | [Cross-Cluster Architecture](docs/cn/source/developer/10_cross_cluster.md) | Data replication, address resolution, gRPC proxy |
| Developer Guide | [Security Framework](docs/cn/source/developer/11_security.md) | JWT + RBAC + TLS, certificate management |
| Reference | [Technical Specifications](docs/cn/source/appendix/1_tech_specifications.md) | app.yaml / module / cluster definition specifications |
| Reference | [Parameter Reference](docs/cn/source/appendix/2_parameter_reference.md) | All parameters and environment variables |
| Reference | [CLI Commands](docs/cn/source/appendix/3_commandline_tools.md) | Complete reference of 24 subcommands |
| Reference | [Shell Programming](docs/cn/source/appendix/4_shell_programming.md) | Built-in functions, container directories, file exchange interfaces |
| Reference | [Exit Code Specification](docs/cn/source/appendix/6_exit_code_spec.md) | App exit code conventions and scheduling policies |
| Reference | [FAQ](docs/cn/source/faq.md) | Common questions on deployment, tasks, semaphores, and WebUI |

### English Design Documents

| Document | Description |
|------|------|
| [VTask Design](docs/vtask-design.md) | Complete VTask design document |
| [WebUI / VS Code](docs/webui-vscode-plan.md) | WebUI + VS Code extension architecture plan |
| [REST API](docs/rest-api-reference.md) | Complete reference of 76 endpoints |
| [gRPC API](docs/grpc-api.md) | 95 RPC call examples |
| [Security Framework](docs/security.md) | JWT + RBAC + TLS configuration |

## 🧪 Example Apps

- **[hello-scalebox](examples/hello-scalebox/)** — the first introductory app
- **[app-primes](examples/app-primes/)** — prime number calculation, demonstrating data parallelism
- **[app-copy](examples/app-copy/)** — cross-node data transfer
- **[remote-primes](examples/remote-primes/)** — cross-cluster prime number calculation
- **[vtask](examples/vtask/)** — VTask pipeline examples (DEFAULT / HOST-BOUND)

## 🏗️ Standard Modules

- **dir-list** — directory listing, generates lists of files to process
- **file-copy** — single-file copy based on rsync-over-ssh
- **dir-copy** — directory-level copy
- **rsync-copy** — efficient rsync transfer
- **ftp-copy** — FTP protocol transfer
- **cron** — scheduled message triggering
- **cluster-head** / **node-agent** — cluster node management

## 🔗 Related Software

- [PostgreSQL](https://github.com/postgres/postgres) — metadata storage
- [gRPC](https://github.com/grpc/grpc) — high-efficiency inter-component communication protocol
- [Go](https://github.com/golang/go) — the language for cloud-native application development

## 🤝 Contributing

Issues and Pull Requests are welcome at [GitHub Issues](https://github.com/kaichao/scalebox/issues/new).

## 📄 License

[Apache License 2.0](LICENSE) © Kaichao Wu
