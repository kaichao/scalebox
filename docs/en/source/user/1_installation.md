# 1. System Installation and Deployment

## 1.1 Environment Requirements

### 1.1.1 Hardware Requirements

#### Chip Architectures and Operating Systems

|  Architecture | OS       |  Notes                    |
| -------- | --------- | ----------------- |
| x86_64   | Linux     |                   |
| arm64    | Linux     |                   |
| x86_64   | MacOS     | For development environments |
| arm64    | MacOS     | To be tested, for development environments |
| x86_64   | Win64/WSL | To be tested, for development environments |

In production environments, all nodes run 64-bit Linux (CentOS 7/8/9, Debian 12/13, Ubuntu 20/22/24/26), etc.

#### Memory Requirements
- Head node: ≥8GB
- Compute nodes: on demand, ≥8GB recommended

#### Storage Space
- Head node: ≥100GB
- Compute nodes: on demand, ≥100GB recommended

### 1.1.2 Software Dependencies

#### Container Engines / Container Runtimes

- Compute nodes

|  Container engine / runtime | Version           |  Notes   |
| --------------------- | ----------------- |  ------- |
| docker-ce             | 26.1<sup>+</sup> |          |
| singularity           | 3.8<sup>+</sup>   |          |
| podman                | 4.8<sup>+</sup>   | To be tested |
| containerd + nerdctl  | 1.6<sup>+</sup>   | To be tested |
| apptainer             |                   | To be tested |
| Kata Containers       |                   | To be tested |

- Head node: docker-ce 20.10<sup>+</sup>, Docker Compose v2<sup>+</sup>

#### Database
- Head node: postgresql 18<sup>+</sup>, deployed as a container

#### Cluster Storage (optional)
- Install the cluster storage client on all nodes to build unified cluster storage
  - glusterfs
  - NFS

#### Network Configuration

- Head node and compute nodes must be in the same subnet of the internal network
- If cross-cluster computing is required, the head node and related transfer nodes need external addresses for communicating with other clusters

## 1.2 Installation Steps

### 1.2.1 Single-Node Quick Deployment (Docker Compose)

Suitable for development and test environments; all services run on a single machine.

**Get the code**:
```bash
git clone https://github.com/kaichao/scalebox.git
cd scalebox
```

**Generate the key system** (certificates + JWT signing key):
```bash
cd build && bash gen-secrets.sh   # → ./secrets/
```

**Build images**:
```bash
make -C build/                    # build all service images
```

**Start services**:
```bash
docker compose -f build/compose.yaml up -d
```

Services started include:
- **controld**: gRPC control service (:50051)
- **actuator**: container launcher
- **database**: PostgreSQL database (:5432)
- **webui**: Web management interface (:8088)

**Check service status**:
```bash
docker compose -f build/compose.yaml ps
```

**Security layer (optional)**:

Enable gRPC TLS transport encryption:
```bash
docker compose -f build/compose.yaml -f build/compose.tls.yaml up -d
```

Enable database certificate authentication:
```bash
docker compose -f build/compose.yaml -f build/compose.db-cert.yaml up -d
```

Enable JWT authentication (with Token Service):
```bash
docker compose -f build/compose.yaml -f build/compose.security.yaml up -d
```

### 1.2.2 Multi-Node Cluster Deployment

For production multi-node deployment, taking 1 HEAD node + 4 NODE nodes as an example:

| Name | Type | IP Address |
| --- | ---- | ------- |
| h0  | HEAD | 10.0.6.100 |
| n0  | NODE | 10.0.6.101 |
| n1  | NODE | 10.0.6.102 |
| n2  | NODE | 10.0.6.103 |
| n3  | NODE | 10.0.6.104 |

**HEAD node configuration**:
- Install Docker 26.1+
- Clone the code repository and run `make -C build/` to build images
- Start the controld + actuator + database + webui containers
- Configure `/etc/hosts` with IP mappings for all nodes

**Compute node configuration**:
- Install Docker 26.1+
- Set up passwordless SSH from HEAD to NODE (actuator starts containers on compute nodes via SSH)
- Configure hostnames and `/etc/hosts`

**Cluster storage (optional)**:
- Install the glusterfs client and mount the shared storage volume
- See :doc:`Appendix 7 - Base Software Installation <../appendix/7_software_install>` for details

## 1.3 Configuration

### 1.3.1 Environment Variables

Common environment variables:

| Variable | Description | Default |
|------|------|--------|
| `GRPC_SERVER` | controld gRPC address | `localhost:50051` |
| `DATABASE_URL` | Database connection string | — |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `MY_APP_CLUSTER` | Default cluster name | — |

Database connection priority: `DATABASE_URL` > `PGHOST`/`PGPORT`/`PGUSER`/`PGDB`/`PGPASS`

### 1.3.2 Security Configuration

| Variable | Description |
|------|------|
| `SECURITY_ENABLED` | Security master switch (default false) |
| `GRPC_TLS_ENABLED` | gRPC TLS encryption |
| `AUTH_TOKEN` | JWT authentication token |

See :doc:`Security Framework document <../../go-scalebox/docs/security>` for security configuration details (in English).

## 1.4 Verifying the Installation

### 1.4.1 Health Check

```bash
# Check service status
docker compose -f build/compose.yaml ps

# Check controld gRPC connectivity
grpcurl -plaintext localhost:50051 list

# CLI connection test
scalebox cluster list
```

### 1.4.2 Run an Example App

```bash
cd examples/hello-scalebox
echo "Hello Scalebox" | scalebox run
scalebox app list
```

### 1.4.3 Access the WebUI

Open `http://<head-node-ip>:8088` in your browser to view the Dashboard, App list, cluster resources, and more.
