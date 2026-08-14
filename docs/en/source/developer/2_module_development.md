# 2. Module Development

## 2.1 Module Design Principles

### 2.1.1 Stateless Design
- Task execution is stateless and repeatable
- Task-level idempotency
- Flexible invocation; different logic processing defined through environment variables and task headers

### 2.1.2 Generality Design
- Supports task-header-driven execution
- Main-router-based, task-header-driven flow control
- Standardized interface design

## 2.2 Module Design Methods

### 2.2.1 task-body Design
The task-body is the task identifier and must be unique within a module:
- Simple tasks: directly use the relative path of the file to process
- Complex tasks: use a canonical identifier (URI) of the input path
- Supports JSON-format compressed representation without null characters

### 2.2.2 task-headers Design
- Configuration parameters
- Configuration parameters that need runtime adjustment are usually passed as environment variables
- Configuration parameters can be modified in external orchestration systems

### 2.2.3 In-Module Directory Design

#### Code Directory / Configuration Directory
- Usually small in data volume; the key consideration is flexible configuration
- Can be placed in shared storage or packaged directly into the container image
- Code execution order:
  1. Environment variable specification: ACTION_RUN, ACTION_CHECK, ACTION_SETUP, ACTION_TEARDOWN
  2. /app/bin/{run.sh,check.sh,setup.sh,teardown.sh}
  3. /app/share/bin/{run.sh,check.sh,setup.sh,teardown.sh}

#### Data Directories
In compute modules with large data volumes, directories should be classified and flexibly configurable according to their usage characteristics during computation:

| Directory type | Description |
| ------------ | -------------------------------------------------- |
| Input data directory | Different data storage forms can be flexibly configured considering read frequency and total data volume |
| Intermediate file directory | Usually on local storage |
| Output result directory | Different data storage forms can be flexibly configured considering write frequency and total data volume |

## 2.3 Module Integration Implementation

### 2.3.1 The Sidecar Pattern
Modules run in the sidecar pattern:
- **run**: a single execution of the task
- **check**: pre-condition check before task execution
- **setup**: set up the initial execution environment
- **teardown**: clean up the compute environment before exit

### 2.3.2 Algorithm Execution run.sh


#### Agent Interface Files

After the user program finishes running, the agent post-processes the execution results. The files are under ```${WORK_DIR}```, and information is exchanged mainly through the following files:

| File name | Description |
| ----------------- | ----------------------------------------- |
| task-exec.yaml    | The main file for task execution results, recording the user program's execution results in YAML |
| extra-attrs.yaml  | The additional attributes file for the run, stored in the extras field as jsonb |
| sink-tasks.txt    | The list of subsequent tasks, one task per line |
| timestamps.txt    | Custom timestamp file, often used for debugging programs and testing program performance; recorded in extras->>'timestamps' |
| input-files.txt   | List of input files (directories) (absolute paths), used to count input file bytes |
| output-files.txt  | List of output files (directories) (absolute paths), used to count output file bytes |
| network-files.txt | List of files (directories) read/written over the network (absolute paths), used to count network read/write and corresponding file read/write bytes |
| removed-files.txt | List of files (directories) to delete (absolute paths), deleted after read/write statistics are complete |
| cleanup-files.txt | Bash cleanup commands before exit (removed-files.txt only deletes local files) (to be implemented) |
| auxout.txt        | Auxiliary output file, recording user-visible output information |

#### sink-tasks.txt File Format

| No. | Name | Example |
| --- | ------------------------------ | ---------------------------------------------- |
| 1   | text-body                      | body1                                          |
| 2   | json-body                      | {"bh0":"a","body":"body2"}                     |
| 3   | text-body + headers            | body3,{"h0":"a","h1":"b"}                      |
| 4   | json-body + headers            | {"bh0":"a","body":"body4"},{"h0":"a","h1":"b"} |
| 5   | sink-mod + text-body           | sink-mod0,body5                                |
| 6   | sink-mod + text-body + headers | sink_mod1,body6,{"h0":"a","h1":"b"}            |
| 7   | sink-mod + json-body + headers | sink-mod2,{"bh0":"a","body":"body7"},{"h0":"a","h1":"b"} |

For formats 1-4, the sink-module is determined by the environment variable SINK_MODULE.

#### timestamps.txt File Format
- One record per line
- Lines whose first character is ```#``` are comment lines
- Line format: ```<timestamp>[,<label>]```, where label is optional

- Timestamp formats

| name        | format                              |
| ----------- | ----------------------------------- |
| RFC3339Nano | 2006-01-02T15:04:05.999999999Z07:00 |
| RFC3339     | 2006-01-02T15:04:05Z07:00           |
|             | 2006-01-02T15:04:05.999999999       |
|             | 2006-01-02T15:04:05                 |
|             | 2006-01-02 15:04:05.999999999       |
|             | 2006-01-02 15:04:05                 |

#### input-files.txt / output-files.txt Format and Processing
- One record per line
- Format output by the application per line: ```<path-item>[,<bytes>]```
- Format after client processing: ```<path-item>,<bytes>```
- After merging, recorded in the database field ```extra->>'iobytes'```; the following items are all numeric byte counts
  - global_input
  - global_output
  - tmpfs_input
  - tmpfs_output
  - local_input
  - local_output
  - input_bytes
  - output_bytes

Where,
- global is determined by the cluster's data_root or the app's global_directories attribute
- tmpfs is determined by the /dev/shm prefix
- local is everything besides global/tmpfs
- input_bytes/output_bytes are the totals

The above records can be viewed through the view ```v_task_iobytes```.

#### network-files.txt Format and Processing
- One record per line
- Format output by the application per line: ```<from_to>,<remote_host>,(<local_path>|<local_bytes>),<remote_path>[,<remote_bytes>]```
  - from_to: input/output direction; from means reading from outside, to means writing to an external node
  - remote_host: hostname or IP address; uniformly converted to hostname before being stored server-side
  - local_path: file name or directory name; converted to byte counts on the client side
  - local_bytes: byte count of the local path item
  - remote_path:
  - remote_bytes: if not set, equals local_bytes
- Format after client processing: ```<from_host>,<to_host>,<net_bytes>,<io_host>,<io_path>,<io_bytes>```
  - from_host: source host
  - to_host: target host
  - net_bytes: bytes transferred over the network; if local_path, computed locally as bytes
  - io_host: corresponds to remote_host in the original record
  - io_path: corresponds to remote_path in the original record
  - io_bytes: corresponds to remote_bytes in the original record
- After merging, recorded in the database field ```extra->>'network_io'```
  - Record format is an array; each line format:
    - ```<from_host>,<to_host>,<net_bytes>,<io_host>,<io_type>,<io_bytes>```
  - from_host:
  - to_host:
  - net_bytes:
  - io_host: hostname of the other end of the network transfer
  - io_type: one of ```global_input/global_output/local_input/local_output/tmpfs_input/tmpfs_output```; determined from io_path before storage

The above records can be viewed through the view ```v_task_network_io```.


## 2.4 Module Image Packaging

### 2.4.1 Building an Algorithm Module Based on agent
```dockerfile
FROM scalebox.net/platform/agent:latest
```

### 2.4.2 Copying agent Components into the Algorithm Module
```dockerfile
COPY --from=scalebox.net/platform/agent:latest /usr/local/ /usr/local/
```

## 2.5 Module Unit Testing

### 2.5.1 Testing a Standalone Module
The default module code directory is: `./code`

```bash
# Start module testing from the command line alone
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path}

# Create the module test first, then add test tasks
app_id=$(scalebox run --image-name ${module_image_name} --code-path ${code_path}| cut -d':' -f2 | tr -d '}')
scalebox task add --app-id=${app_id} --header header1=${header1} ${task_body}
```

### 2.5.2 Unit Testing with a Router
The default router code directory is: `./mr-code`

```bash
echo ${task-body} | scalebox run --image-name ${module_image_name} --code-path ${code_path} --mr-image-name ${mr_image_name} --mr-code-path ${mr_code_path}
```

### 2.5.3 Unit Testing Based on Pipeline Apps
- For more complex unit tests, write a standalone app definition file app.yaml and an environment variable definition file scalebox.env

## 2.6 Module Types in Detail

### 2.6.1 Base Module (agent)
- Provides the module base functions
- Manages the lifecycle of task execution within app modules
- Interacts with controld in the server-side runtime to record task execution state

### 2.6.2 File Transfer Modules
- **file-copy**: file transfer based on rsync-over-ssh
- **dir-copy**: directory copy module
- **rsync-copy**: rsync copy module
- **ftp-copy**: FTP protocol transfer module

### 2.6.3 Helper Function Modules
- **dir-list**: directory listing module
- **cron**: scheduled task module
- **cluster-head**: cluster head node module
- **node-agent**: node agent module

## 2.7 Error Handling and Fault-Tolerance Mechanisms

### 2.7.1 Task-Level Fault Tolerance
- Automatic retries scheduled based on individual task return codes
- Supports the task retry mechanism

### 2.7.2 Combined Fault Tolerance
- Solves complex non-logical errors
- Built on task-level fault tolerance
- Slot-level error handling

## 2.8 Best Practices

### 2.8.1 Module Development Standards
1. Keep each module to a single responsibility
2. Support environment variable configuration
3. Implement task idempotency
4. Set reasonable timeout values

### 2.8.2 Performance Optimization
1. Optimize data locality
2. Reduce unnecessary file operations
3. Use caching appropriately
4. Batch processing optimization

### 2.8.3 Testing Suggestions
1. Unit tests cover core logic
2. Integration tests verify module collaboration
3. Performance tests ensure requirements are met
4. Fault-tolerance tests verify reliability

## 2.9 App-Level Debugging and Optimization

### 2.9.1 Debugging Methods

1. **Unit tests**: test the functional correctness of each module
2. **Integration tests**: test collaboration and communication between modules
3. **Performance tests**: test the app's performance and scalability
4. **Fault-tolerance tests**: test the app's fault-tolerance and recovery capabilities

### 2.9.2 Performance Optimization

1. **Bottleneck analysis**: analyze performance bottlenecks and locate optimization points
2. **Algorithm optimization**: optimize algorithm implementations to improve computing efficiency
3. **Parallel optimization**: optimize parallel strategies to improve resource utilization
4. **I/O optimization**: optimize data read/write to reduce I/O overhead

### 2.9.3 Monitoring and Tuning

1. **Real-time monitoring**: monitor the app's running state in real time
2. **Performance analysis**: analyze performance data to find optimization opportunities
3. **Dynamic adjustment**: dynamically adjust configuration parameters based on load
4. **Capacity planning**: plan resource capacity based on business requirements
