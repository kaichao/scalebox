# 1. Technical Specification Reference

A Scalebox app is an application running on the Scalebox platform. Typical Scalebox applications include high-throughput data processing and large-scale data transfer.

## 1.1 Scalebox Runtime Components

### 1.1.1 Core Runtime Components

Scalebox's runtime environment includes the following core components:

- **controld**: a gRPC-based control service that manages actuators and compute nodes
- **actuator**: a launcher service that starts slots on compute nodes via SSH or external schedulers
- **database**: a PostgreSQL database storing metadata for apps, modules, tasks, slots, etc.
- **webui**: the web interface for operating and managing scalebox
- **agent**: an agent running on compute nodes, responsible for task execution and state management
- **cluster-admin**: a software module for resource acquisition and management that interacts with HPC scheduling systems; also a system-level app

### 1.1.2 Inter-Component Communication

- **Control plane communication**: uses the gRPC protocol, port 50051
- **Data plane communication**: through message queues and the database
- **Inter-node communication**: through SSH or dedicated network protocols

## 1.2 App Definition Specification (complete app.yaml schema)

The app definition file is a YAML text file used to define a Scalebox application (App) and its Modules, and also supports Clusters.

An example app definition file format:

```yaml
name: perf-test.my-app
label: perf-test
version: 1.0.0
cluster: ${CLUSTER}
parameters:
  initial_status: RUNNING
  main_router: main-router
  is_cluster_admin: yes
  default_sleep_count: 20
  comment: This is a sample app.

modules:
  module-1:
    arguments:
      ...
    parameters:
      ...
    ...
  module-2:
    ...

clusters:
  cluster-1:
    ...
  cluster-2:
    ...


```

The app definition file is divided into the main configuration attributes, module definitions, cluster definitions, and other parts.

If the app definition file is used only for cluster definitions, all attributes under the main configuration do not take effect.

The default name of the app definition file is app.yaml in the current directory.

Field descriptions under the main configuration:
- *name*: the app name, usually composed of dot-separated identifiers
- *label*: the app name displayed in the system; can usually be identified in Chinese.
- *version*: the app version number
- *cluster*: the default cluster name for all modules in the app
- *comment*: comment information
- *parameters*: the app's parameter list
  - *initial_status*: the app's initial status, values: 'RUNNING'/'INITIAL'. If set to 'RUNNING', the App enters the running state directly after creation;
  - *main_router*: specifies the default router for modules in the app. If a module's subsequent module (sink-module) is empty, main-router is designated as the subsequent module;
  - *default_sleep_count*: the default max_sleep_count parameter for all modules; the default value is 100 (6 seconds per unit, 10 minutes total)

For detailed definitions of module and cluster, see the subsequent sections.


- Template parameter files in the app definition file

Except for the virtual template parameter CLUSTER_DATA_DIR, template variables must be defined before use. The following are template variable definitions and their priority ordering (from low to high)
- ```/etc/scalebox/environments```
- ```${HOME}/.scalebox/environments```
- ```${PWD}/scalebox.env```
- ```${PWD}/${{env-defined}}.env```
- ```${PWD}/${{env-defined}}_${{app-defined}}.env```
- Environment variables already defined on the current command line


## 1.2 Module Definition Specification

```yaml
  my-module:
    label: My First Module
    base_image: scalebox/agent
    cluster: my-cluster
    command: docker run -d --network=host {{ENVS}} {{VOLUMES}} {{IMAGE}}
    arguments:
      ...
    parameters:
      ...
    environmens:
      ...
    volumes:
      ...
    slots:
      ...
    sink_modules:
      ...
    sink_vmodules:
      ...
    comment: This is new algorithm module.

```

Field descriptions for Module:
- *label*: the Module displayed in the app interface
- *base_image*: the container image name
- *cluster*: the cluster name
- *command*: the command template for running the container
- *arguments*: standard variables on the container side, usually mapped to environment variables
        ...
- *parameters*: server-side parameters of the module
        ...
- *environmens*: environment variables
        ...
- *volumes*: physical paths mapped through volumes
        ...
- *slots*: the module's slot definitions
        ...
```yaml
  slots:
    - ${nodes}[:n]
    - ${nodes}:${n}:${group_prefix}
```
The second line is used for global/grouped slot definitions in host-bound scenarios.
- nodes is the node where the slot resides
- n is the number of slots
- the regular expression of the group prefix
For grouped slots, set group_prefix in the slot's parameters to the group prefix

- *sink_modules*: identifies physical associations between Modules; in cross-cluster apps, used to identify associations between Modules across clusters;
        ...
- *sink_vmodules* (string array): used for logical relationships between Modules.

## 1.3 Cluster Definition Specification

An example cluster definition:

```yaml
  mycluster:
    label: My new clster
    parameters:
      uname: myuser
      port: 10022
      data_root: /global-fs/scalebox/mydata
      local_ip_index: 2
      num_of_executors: Inline cluster only
      channel_size: channel size fo executor, Inline cluster only
      grpc_server: 192.168.3.123:50051
    total_resources:
      num_cores: cpu cores
      total_mem_gb:
      total_disk_tb:
    status: ON
    comment:

```

- *label*:
- *parameters*:
  - *port*: the host's default port number
  - *uname*: the host's default username
  - *data_root*: the cluster's data directory
  - *local_ip_index*: the index number used to extract the local IP address (hostname -I)
  - *grpc_server*:	there is an external field with the same name.
  - *remote_grpc_server*:

## 1.4 Identifier Naming Rules
### 1.4.1	File Name Naming Rules
File name characters: digits, uppercase and lowercase English letters, underscores, dots.
### 1.4.2	URI Naming Rules
Resources such as Apps and Modules are uniquely identified by URIs (Uniform Resource Identifiers). Common URIs mainly fall into two categories: URLs (Uniform Resource Locators) and URNs (Uniform Resource Names). A URI refers to a resource, a URL locates a resource by address, and a URN locates a resource by name. That is, URL and URN are subsets of URI.
Apps, Modules, etc. are defined through URNs.

### 1.4.3	Version Number Naming Rules
Resource types such as Apps and Modules can be represented by versions; version definitions follow semantic versioning.
The version format is as follows: major.minor.patch, with the following increment rules:
·	Major version: when you make incompatible API changes,
·	Minor version: when you add functionality in a backward-compatible manner,
·	Patch version: when you make backward-compatible bug fixes.

Pre-release versions and build metadata can be appended after "major.minor.patch" as extensions.



- Single-thread performance comparison

| Algorithm	        | Parameter value              |Hardware AES-NI(x64)| Software implementation (no AES-NI)| Notes              |
| ---------------- | ---------------------------- | ------------- | ----------------| -------------------|
| AES-128-GCM      | aes128-gcm@openssh.com       | 5~7 GB/s      |	150~200 MB/s    | Best choice, extremely fast with hardware acceleration |
| AES-256-GCM	     | aes256-gcm@openssh.com       | 3~5 GB/s	    | 120~150 MB/s    | Slightly slower but more secure      |
| ChaCha20-Poly1305| chacha20-poly1305@openssh.com| 1~2 GB/s      | 600~700 MB/s    | Highly efficient in software, especially faster than AES without AES-NI |

