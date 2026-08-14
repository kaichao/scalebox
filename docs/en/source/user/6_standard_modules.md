# 6. Standard Modules in Detail

## 6.1 Base Module (agent)

### 6.1.1 Module Functions
- Provides the module base functions
- Manages the lifecycle of task execution within app modules
- Interacts with controld in the server-side runtime to record task execution state

### 6.1.2 Module Features
- **Standardized interface**: provides a unified task execution interface
- **State management**: manages the lifecycle state of tasks
- **Error handling**: supports task failure retries and error reporting
- **Resource management**: manages compute resources and the execution environment

### 6.1.3 Configuration Parameters
```yaml
modules:
  my-module:
    base_image: scalebox.net/platform/agent:latest
    arguments:
      code_path: ${PWD}/code
      task_max_seconds: 600
    parameters:
      task_dist_mode: DEFAULT
```

### 6.1.4 Building Non-Agent Modules
```dockerfile
FROM scalebox.net/platform/agent:latest
COPY --from=scalebox.net/platform/agent:latest /usr/local/ /usr/local/
COPY ./code /app/bin/
```

## 6.2 File Transfer Modules

### 6.2.1 Common Configuration
File transfer modules are implemented through rsync-over-ssh, supporting cross-node data transfer.

#### Task Body Format
- **task-body**: the local relative path to be transferred to the remote end

#### Task Header Parameters
| Task header | Default | Description |
| ---------- | ----- | ------------------------- |
| source_url |       | "/local/dir" |
| target_url |       | "user@remote-ip:remote-port/remote/dir" |
| keep_source | "yes" | "yes" / "no", whether to keep the source file/directory |

### 6.2.2 file-copy Module
Used for file copy operations.

#### file-copy Usage Example
```yaml
modules:
  file-copy:
    base_image: scalebox.net/platform/file-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

#### file-copy Configuration Parameters
- **source_url**: source file path
- **target_url**: target file path
- **keep_source**: whether to keep the source file

### 6.2.3 dir-copy Module
Used for directory copy operations.

#### dir-copy Usage Example
```yaml
modules:
  dir-copy:
    base_image: scalebox.net/platform/dir-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

### 6.2.4 rsync-copy Module
An efficient data transfer module based on rsync.

#### rsync-copy Usage Example
```yaml
modules:
  rsync-copy:
    base_image: scalebox.net/platform/rsync-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

#### rsync-copy Advanced Parameters
- **rsync_options**: rsync command-line options
- **compress**: whether to enable compressed transfer
- **delete**: whether to delete extra files on the target side

### 6.2.5 ftp-copy Module
A data transfer module supporting the FTP protocol.

#### ftp-copy Usage Example
```yaml
modules:
  ftp-copy:
    base_image: scalebox.net/platform/ftp-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

#### ftp-copy Configuration Parameters
- **ftp_server**: FTP server address
- **ftp_user**: FTP username
- **ftp_password**: FTP password
- **ftp_port**: FTP port (default 21)

## 6.3 Helper Function Modules

### 6.3.1 cron Module
A scheduled task module that supports executing tasks on a schedule.

#### cron Usage Example
```yaml
modules:
  cron:
    base_image: scalebox.net/platform/cron
    arguments:
      cron_expression: "0 2 * * *"
      command: "/app/bin/backup.sh"
```

#### cron Configuration Parameters
- **cron_expression**: cron expression (minute hour day month weekday)
- **command**: the command to execute
- **timezone**: timezone setting (default UTC)

### 6.3.2 cluster-head Module
A cluster head node management module.

#### cluster-head Usage Example
```yaml
modules:
  cluster-head:
    base_image: scalebox.net/platform/cluster-head
    parameters:
      slot_options: slot_on_head
```

#### cluster-head Features
- Cluster management coordination
- Resource scheduling optimization
- Node status monitoring

### 6.3.3 node-agent Module
A node agent module that runs on compute nodes.

#### node-agent Usage Example
```yaml
modules:
  node-agent:
    base_image: scalebox.net/platform/node-agent
    parameters:
      cluster: ${CLUSTER}
```

#### node-agent Features
- Node resource management
- Task execution monitoring
- System status reporting

### 6.3.4 dir-list Module
A directory listing module used to generate lists of files to process.

#### dir-list Usage Example
```yaml
modules:
  dir-list:
    base_image: scalebox.net/platform/dir-list
    arguments:
      target_dir: /data/input
      pattern: "*.txt"
```

#### dir-list Configuration Parameters
- **target_dir**: target directory path
- **pattern**: file matching pattern
- **recursive**: whether to scan subdirectories recursively
- **output_format**: output format (json/csv)

## 6.4 Module Configuration Best Practices

### 6.4.1 Module Selection Principles
1. **Prefer standard modules**: reduce duplicated development work
2. **Consider performance requirements**: choose appropriate transfer protocols and algorithms
3. **Support extensibility**: facilitate future feature extensions

### 6.4.2 Parameter Configuration Suggestions
1. **Set timeouts appropriately**: set task_max_seconds based on task complexity
2. **Optimize parallelism**: set the number of slots based on resources
3. **Error handling**: configure an appropriate retry strategy

### 6.4.3 Module Composition Patterns
1. **Pipeline pattern**: multiple modules execute sequentially
2. **Parallel pattern**: multiple modules process in parallel
3. **Hybrid pattern**: combines pipeline and parallel processing

## 6.5 Module Usage Examples

### 6.5.1 A Complete Data Processing Pipeline
```yaml
modules:
  # Generate the file list
  dir-list:
    base_image: scalebox.net/platform/dir-list
    arguments:
      target_dir: /data/input
      pattern: "*.dat"
  
  # Data transfer
  file-copy:
    base_image: scalebox.net/platform/file-copy
    parameters:
      task_dist_mode: SLOT-BOUND
  
  # Data processing
  data-process:
    base_image: my-registry/data-process:latest
    parameters:
      task_dist_mode: HOST-BOUND
  
  # Result transfer
  rsync-copy:
    base_image: scalebox.net/platform/rsync-copy
    parameters:
      task_dist_mode: SLOT-BOUND
```

### 6.5.2 Scheduled Backup Tasks
```yaml
modules:
  # Scheduled trigger
  cron:
    base_image: scalebox.net/platform/cron
    arguments:
      cron_expression: "0 2 * * *"
      command: "/app/bin/backup.sh"
  
  # Data backup
  rsync-copy:
    base_image: scalebox.net/platform/rsync-copy
    parameters:
      task_dist_mode: SLOT-BOUND
      compress: yes
```

## 6.6 Module Development Suggestions

### 6.6.1 Module Design Principles
1. **Single responsibility**: each module is responsible for only one specific function
2. **Stateless design**: supports repeated task execution
3. **Standardized interface**: follow the Scalebox module specification
4. **Error handling**: complete error reporting and recovery mechanisms

### 6.6.2 Performance Optimization
1. **Batch processing**: support batch task processing to improve efficiency
2. **Caching mechanisms**: use caching appropriately to reduce redundant computation
3. **Parallel optimization**: fully leverage multi-core and distributed computing
4. **I/O optimization**: optimize data read/write performance

### 6.6.3 Testing Suggestions
1. **Unit tests**: test the basic functions of the module
2. **Integration tests**: test collaboration between modules
3. **Performance tests**: verify the module's performance metrics
4. **Fault-tolerance tests**: test the module's fault-tolerance capability
