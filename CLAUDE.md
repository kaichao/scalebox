# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Scalebox is a cloud-native stream computing engine designed for distributed, heterogeneous computing clusters. It enables containerized single-machine user algorithms to run on pipeline-organized modular hierarchical large-scale parallel processing with task-level fault tolerance.

## Development Workflow

### Quick Start
```bash
# Clone and setup
git clone https://github.com/kaichao/scalebox.git
cd scalebox

# Setup development environment
cd runtime && make prepare && make pull-all && make get-cli

# Start control plane services
cd runtime && make all

# Build and run hello example
cd examples/hello-scalebox && make build
echo "hello" | scalebox run

# Run distributed primes calculation
cd examples/app-primes && make run
```

### Development Cycle
1. **Module Development**: Create/modify modules in `dockerfiles/`
2. **Build & Test**: Use `make build` in module directories
3. **Integration**: Create app YAML configurations in `examples/`
4. **Deploy**: Run applications with `scalebox run`
5. **Monitor**: Use `scalebox app list` and `scalebox app logs`

## Project Structure

### Core Directories
- `runtime/`: Control plane services (controld, actuator, database)
- `dockerfiles/`: Standard reusable modules with Docker definitions
- `examples/`: Application examples with multi-language implementations
- `pkg/`: Shared Go packages (auth, postgres, semaphore, task, etc.)
- `docs/`: Documentation in English and Chinese
- `features/`: Feature-specific tests and demonstrations

### Key Concepts
- **App**: Complete application workflow defined in YAML
- **Module**: Containerized processing component
- **Task**: Unit of work execution
- **Slot**: Compute resource allocation
- **Cluster**: Collection of compute resources (inline or external)

## Key Commands

### Environment Setup
```bash
# Install dependencies and prepare environment
cd ~/scalebox/runtime && make prepare

# Pull container images and get CLI tools
cd ~/scalebox/runtime && make pull-all
make get-cli

# Start scalebox control services
cd ~/scalebox/runtime && make all
```

### Application Development

#### Building Applications
```bash
# Build hello-scalebox example
cd examples/hello-scalebox && make build

# Build app-primes example  
cd examples/app-primes && make build

# Build with custom tag
cd examples/app-primes && make build TAG=v2
```

#### Running Applications
```bash
# Run hello-scalebox (message-triggered)
cd examples/hello-scalebox && echo "scalebox" | scalebox run

# Run distributed primes calculation
cd examples/app-primes && make run NUM_GROUPS=4 CALC_NODE=local NUM_PARALLEL=2

# Run with custom cluster
cd examples/app-primes && make run CLUSTER=inline NUM_GROUPS=8

# List running applications
scalebox app list

# Stop an application
scalebox app stop <app-name>

# View application logs
scalebox app logs <app-name>
```

#### Application Configuration
Applications are defined using YAML files (e.g., `app.yaml`):
- `name`: Application identifier with variables (e.g., `app-primes-g${NUM_GROUPS}-${TAG}`)
- `cluster`: Target cluster name
- `modules`: Processing components with container images and slot allocations
- `parameters`: Runtime configuration including task distribution mode

### Module Development

#### Building Modules
```bash
# Build all standard modules
cd dockerfiles && make build

# Build specific module
cd dockerfiles/dir-list && make build

# Build with custom tag
cd dockerfiles/dir-list && make build TAG=v2
```

#### Module Structure
Standard modules follow this pattern:
- `Dockerfile`: Container definition with entrypoint
- `main.go`: Go implementation with message handling
- `go.mod`: Go module dependencies
- Tests in `*_test.go` files

#### Module Types
1. **File Operations**: `dir-list`, `file-copy`, `rsync-copy`, `ftp-copy`, `rsyncd`
2. **Data Processing**: `data-grouping-2d`
3. **Utilities**: `cron` (scheduled messaging), `actuator` (key generation)

#### Go Module Development
```bash
# Run module tests
cd dockerfiles/dir-list && go test ./...

# Run with coverage
cd dockerfiles/data-grouping-2d && go test -coverprofile=coverage.out ./...

# Benchmark tests
cd pkg/semagroup && go test -bench=. ./...
```

#### Shared Packages
- `pkg/auth/`: Authentication and authorization
- `pkg/postgres/`: PostgreSQL database operations
- `pkg/semaphore/`: Semaphore-based flow control
- `pkg/task/`: Task management utilities
- `pkg/common/`: Common networking and JSON utilities

### Documentation
```bash
# Build documentation
cd docs && make html
cd docs/cn && make html
```

## Architecture Overview

### Core Components
1. **Control Plane** (`runtime/`):
   - `controld`: gRPC-based control service managing actuators and compute nodes
   - `actuator`: Launcher service that starts slots on compute nodes via SSH or external schedulers
   - `database`: PostgreSQL database storing app/module/task/slot metadata
   - Services communicate via gRPC and database connections

2. **Standard Modules** (`dockerfiles/`):
   - **File Operations**: `dir-list`, `file-copy`, `rsync-copy`, `ftp-copy`, `rsyncd`
   - **Data Processing**: `data-grouping-2d` for 2D dataset grouping operations
   - **Utilities**: `cron` (scheduled messaging), `actuator` (key generation)
   - Modules are containerized with Go-based message handling

3. **Application Examples** (`examples/`):
   - `hello-scalebox`: Basic message-triggered application
   - `app-primes`: Distributed prime number calculation (multi-language)
   - `remote-primes`: Cross-cluster primes calculation demonstration
   - `app-copy`, `cluster-dir-copy`: Cross-cluster data transfer

4. **Shared Packages** (`pkg/`):
   - `pkg/auth/`: JWT and basic authentication
   - `pkg/postgres/`: Database connection and query utilities
   - `pkg/semaphore/`: Task synchronization primitives
   - `pkg/task/`: Task management and status tracking
   - `pkg/variable/`: Variable substitution and expansion

### Key Concepts
- **App**: Complete workflow defined in YAML with modules and parameters
- **Module**: Containerized processing component with message input/output
- **Task**: Unit of work execution triggered by messages
- **Slot**: Compute resource allocation on cluster nodes
- **Cluster**: Resource collection (inline managed or external via Slurm/K8s)

### Development Patterns
1. **Module Development**: Implement `main.go` with message handlers, containerize with Docker
2. **Application Composition**: Define modules, slots, and routing in YAML configs
3. **Multi-language Support**: Python, Go, C, Java, Julia, R via containerized execution
4. **Cross-cluster Operations**: Built-in support for distributed computing across heterogeneous clusters

### Communication Model
- **Control Plane**: gRPC between controld, actuator, and compute nodes
- **Data Flow**: Message-driven between modules via message bus
- **Storage**: PostgreSQL for metadata, local/remote storage for data
- **Authentication**: JWT-based auth for API access

This architecture supports cloud-native design with containerized algorithms, message-driven processing, and task-level fault tolerance for large-scale distributed computing applications.

This architecture supports cloud-native design with containerized algorithms, message-driven processing, and task-level fault tolerance for large-scale distributed computing applications.