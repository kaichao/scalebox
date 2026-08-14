# 2. Quick Start

This chapter guides you through deploying and running your first Scalebox application, letting you experience Scalebox's core functionality in minutes.

## 2.1 Preparing the Environment

### 2.1.1 System Requirements

- **Operating system**: Linux (CentOS 7+/Ubuntu 18.04+/Debian 10+), macOS 10.15+
- **Docker**: 20.10+ (latest stable version recommended)
- **Memory**: ≥ 8GB
- **Disk space**: ≥ 100GB

### 2.1.2 Installing Docker

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y docker.io
sudo systemctl enable --now docker

# CentOS/RHEL
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl enable --now docker
```

- Verify the Docker installation:
```bash
docker --version
docker run hello-world
```

## 2.2 Getting the Scalebox Code

```bash
git clone https://github.com/kaichao/scalebox.git
cd scalebox
```

## 2.3 Running the Hello Scalebox Example

### 2.3.1 Starting a Single-Node Cluster

```bash
make -C runtime all
```

Check the service status:
```bash
docker ps
```

### 2.3.2 Creating and Running an App

```bash
# Create an app
cd examples/hello-scalebox

echo "Docker-based_Scalebox" | scalebox run 

# Check app status
scalebox app list

# Check task execution
scalebox task list
```

## 2.4 Verifying the Installation

### 2.4.1 Health Check

```bash
# Check the status of each service
scalebox cluster status
```

## 2.5 Next Steps

After successfully running Hello Scalebox, you can:

1. **Learn the core concepts**: read the User Guide to understand core concepts such as App, Module, Task, VTask (task group), and Cluster
2. **Try more examples**: explore other example apps (app-primes, app-copy, vtask, etc.)
3. **Use the WebUI**: open `http://localhost:8088` in your browser to enter the visual management interface
4. **Install the VS Code extension**: `cd vscode-scalebox && make install` for TreeView + DAG + app.yaml syntax support
5. **Develop custom apps**: refer to the Developer Guide to develop your own Scalebox applications
