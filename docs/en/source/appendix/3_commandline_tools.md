# 3. Command Line Tool Guide

## 3.1 The scalebox Command Line Tool in Detail

## 3.2 Common Operation Commands

## 3.3 Scripting Examples

## 3.4 Automated Operations Scripts



The scalebox command line tool

## 1.1 Command Line Options

| Option             | Default        | Description                              |
| ----------------- | -------------- | --------------------------------- |
| -e / --env-file   | scalebox.env   | environment variable file, sets the environment variables for command execution. |
| --debug           | 'no'           | sets the debug flag, outputs more debugging and troubleshooting information |

Environment variables are an important method of parameter passing in Scalebox applications. Environment variable definitions in an app can come from multiple environment variable definition files and system-level environment variables. If duplicate variable names exist across different definition files and system-level variables, they are loaded in the following order (files that do not exist are ignored):

- System-level environment variables
- User-defined-name env files
- The scalebox.env file in the current directory
- ${HOME}/.scalebox/environments
- /etc/scalebox/environments

Among them, user-defined-name env files can be cascade-loaded by file name.

Example:
If the user-defined env file name is p419_48nodes_1266932744.env, the following files are loaded in order of priority from high to low:
- p419_48nodes_1266932744.env
- p419_48nodes.env
- p419.env


## 1.2 Subcommand Overview

```
scalebox
├── run                    start an app (create App + send initial messages)
├── cluster                cluster management
│   ├── list               list clusters
│   ├── create             create a cluster
│   ├── show               view cluster details
│   ├── set-status         set cluster status
│   ├── get-parameter      get cluster parameters
│   ├── allocate           allocate cluster resources
│   └── release            release cluster resources
├── host                   host management
│   ├── list               list hosts
│   ├── create             add a host
│   ├── show               view host details
│   ├── set-status         set host status
│   ├── delete             delete a host
│   └── renew              renew a host
├── slot                   slot management
│   ├── list               list slots
│   ├── add                add slots
│   ├── remove             delete slots
│   └── set-status         set slot status
├── app                    app management
│   ├── list               list apps
│   ├── create             create an app
│   ├── show               view app details
│   ├── set-status         set app status
│   ├── set-finished       mark an app finished
│   ├── delete             delete an app
│   ├── add-remote         add a remote app link (cross-cluster)
│   └── add-slots          dynamically add slots
├── module                 module management
│   ├── list               list modules
│   └── show               view module details
├── task                   task management
│   ├── list               list tasks
│   ├── show               view task details
│   ├── add                add tasks
│   ├── delete             delete tasks
│   ├── get-header         get a task header
│   ├── set-header         set a task header
│   ├── remove-header      delete a task header
│   └── log                view task logs
├── validate               validate the app definition file (app.yaml)
├── semaphore              semaphore management
│   ├── list               list semaphores
│   ├── create             create semaphores
│   ├── get                get a semaphore value
│   ├── increment          increment by one
│   ├── decrement          decrement by one
│   └── add-value          increment/decrement by N
├── semagroup              semaphore group operations
│   ├── max                maximum value
│   ├── min                minimum value
│   ├── increment          increment the minimum by one
│   ├── decrement          decrement the maximum by one
│   ├── diff-min           difference from the minimum
│   └── diff-max           difference from the maximum
├── variable               shared variable management
│   ├── list               list variables
│   ├── get                get a variable value
│   ├── set                set a variable value
│   └── delete             delete a variable
├── global                 global variable management
│   ├── list               list global variables
│   ├── get                get a value
│   ├── set                set a value
│   └── delete             delete
├── vtask                  virtual task management
│   ├── list               list vtasks
│   ├── list-subtasks      list sub-tasks
│   ├── get                view vtask details
│   ├── fail               mark a vtask failed
│   ├── bind               bind resources
│   ├── unbind             unbind resources
│   ├── add-subtask        add sub-tasks
│   ├── get-variable       get vtask-scoped variables
│   ├── set-variable       set vtask-scoped variables
│   ├── create-semaphore   create vtask-scoped semaphores
│   ├── get-semaphore      get vtask-scoped semaphores
│   ├── add-semaphore-value increment/decrement vtask-scoped semaphores
│   └── delete-semaphore   delete vtask-scoped semaphores
├── channel                priority queues (cross-app communication)
│   ├── pull               dequeue
│   └── push               enqueue
├── fs                     file system operations
│   ├── ls                 list files
│   └── stat               view file metadata
├── event                  event recording
│   ├── task-add           task events
│   ├── slot-add           slot events
│   └── misc-add           miscellaneous events
├── status                 overall system status
└── help                   help information
```

## 1.3 cluster Subcommands

### 1.3.1 cluster list

List all clusters.

```bash
scalebox cluster list
```

### 1.3.2 cluster show

View the details of a specified cluster.

```bash
scalebox cluster show <cluster-name>
```

### 1.3.3 cluster create

Create a new cluster.

```bash
scalebox cluster create --name my-cluster --grpc-server 10.0.0.1:50051
```

### 1.3.4 cluster set-status

Set cluster status.

```bash
scalebox cluster set-status <cluster-name> ON
```

### 1.3.5 cluster get-parameter

Get cluster parameters.

```bash
scalebox cluster get-parameter <cluster-name> <param-name>
```

### 1.3.6 cluster allocate / release

Allocation and release of dynamic cluster resources.

```bash
scalebox cluster allocate --cluster my-cluster --num-hosts 4
scalebox cluster release --cluster my-cluster
```

## 1.4 host Subcommands

Host (compute node) management.

### 1.4.1 host list

List all hosts.

```bash
scalebox host list
scalebox host list --cluster my-cluster
```

### 1.4.2 host show

View host details.

```bash
scalebox host show <hostname>
```

### 1.4.3 host create

Add a new host.

```bash
scalebox host create --hostname n0 --ip-addr 10.0.6.101 --cluster my-cluster
```

### 1.4.4 host set-status

Set host status.

```bash
scalebox host set-status <hostname> READY
```

### 1.4.5 host delete / renew

Delete a host or renew a dynamic host.

```bash
scalebox host delete <hostname>
scalebox host renew <hostname>
```

## 1.5 slot Subcommands

Slot management.

### 1.5.1 slot list

List slots.

```bash
scalebox slot list
scalebox slot list --host <hostname>
scalebox slot list --module <module-name>
```

### 1.5.2 slot add

Add slots.

```bash
scalebox slot add --module my-module --host n0 --count 4
```

### 1.5.3 slot remove

Delete slots.

```bash
scalebox slot remove --module my-module --host n0
```

### 1.5.4 slot set-status

Set slot status.

```bash
scalebox slot set-status <slot-id> READY
```


## 1.6 app Subcommands

### 1.6.1 app list

List all apps.

```bash
scalebox app list
```

### 1.6.2 app show

View app details.

```bash
scalebox app show <app-id>
```

### 1.6.3 app create

Create an app from a definition file.

```bash
scalebox app create --app-file app.yaml --env-file scalebox.env
```

### 1.6.4 app run

Start an app from the command line. The environment variable file defaults to `./scalebox.env`.

**Single start message**:
```bash
export ENV0=v0
scalebox run --cluster my-cluster --image-name my-image:latest start-item
```

**Parameter table**:

| Parameter name | Corresponding environment variable | Default | Description |
|--------|------------|--------|------|
| app-name | `_APP_NAME` | — | app name |
| cluster | `_CLUSTER` | local | cluster name |
| image-name | `_IMAGE_NAME` | scalebox/agent:latest | main module image |
| code-path | `_CODE_PATH` | ./code (if it exists) | main module code directory |
| slot-regex | `_SLOT_REGEX` | h0 | main module slot configuration |
| mr-image-name | `_MR_IMAGE_NAME` | — | router module image |
| mr-code-path | `_MR_CODE_PATH` | ./mr-code (if it exists) | router module code directory |
| app-file / -f | — | app.yaml (if it exists) | app definition file |
| env-file / -e | — | scalebox.env | environment variable file |

**Piped multi-start tasks**:
```bash
find /data/input -type f | scalebox run --image-name my-image:latest --slot-regex h0:2
```

### 1.6.5 app set-status / set-finished

```bash
scalebox app set-status <app-id> RUNNING
scalebox app set-finished <app-id>
```

### 1.6.6 app delete

```bash
scalebox app delete <app-id>
```

### 1.6.7 app add-remote

Add a cross-cluster app remote link.

```bash
scalebox app add-remote --app-id <local-app-id> --remote-app-id <remote-app-id> --remote-grpc <grpc-addr>
```

### 1.6.8 app add-slots

Dynamically add slots.

```bash
scalebox app add-slots --app-id <app-id> --module <module-name> --host <hostname> --count 4
```

## 1.7 module Subcommands

### 1.7.1 module list

List the app's modules.

```bash
scalebox module list --app-id <app-id>
```

### 1.7.2 module show

View module details.

```bash
scalebox module show --app-id <app-id> <module-name>
```

## 1.8 task Subcommands

### 1.8.1 task list

List tasks, supporting filtering by module and status.

```bash
scalebox task list --app-id <app-id>
scalebox task list --app-id <app-id> --status FAILED
scalebox task list --app-id <app-id> --module <module-name>
```

### 1.8.2 task show

View task details (including stdout/stderr/body/headers).

```bash
scalebox task show <task-id>
```

### 1.8.3 task delete

Delete tasks.

```bash
scalebox task delete <task-id>
```

### 1.8.4 task add

Add tasks.

**Parameters / environment variables**:

| Parameter name | Environment variable name | Description |
|--------|----------|------|
| app-id | `APP_ID` | app ID |
| module-id | `MODULE_ID` | module ID |
| sink-module | `SINK_MODULE` | downstream module name |
| conflict-action | `CONFLICT_ACTION` | conflict handling: ''/'IGNORE'/'OVERWRITE' |
| from-module | — | source module name |
| headers | — | JSON format headers |
| header / -h | — | add a single header (repeatable) |
| to-ip | — | set the to_ip header |
| to-host | — | set the to_host header |
| batch-size | — | batch size for batch adds, default 100 |

**task file format** (default `${WORK_DIR}/sink-tasks.txt`, one per line):

| Type | Example |
|------|------|
| text body | `body` |
| JSON body | `{"hi0":"a","body":"my_body"}` |
| text body + headers | `body,{"h0":"a","h1":"b"}` |
| JSON body + headers | `{"hi0":"a","body":"my_body"},{"h0":"a","h1":"b"}` |
| module name + text body | `module-name,body` |
| module name + text body + headers | `module-name,body,{"h0":"a","h1":"b"}` |

**Control headers**:

| header | Description |
|--------|------|
| `initial_status_code` | initial status, default -1 (READY) |
| `conflict-action` | ''/'IGNORE'/'OVERWRITE' |
| `slot_broadcast` | broadcast to all slots |
| `host_broadcast` | broadcast to all hosts |

### 1.8.5 task get-header / set-header / remove-header

```bash
scalebox task get-header --task-id 123 from_module
scalebox task set-header --task-id 123 my_header value
scalebox task remove-header --task-id 123 my_header
```

### 1.8.6 task log

View task execution logs.

```bash
scalebox task log <task-id>
```

## 1.9 validate Subcommand

Validate the syntax and completeness of the app definition file (app.yaml).

```bash
scalebox validate --app-file=<app-yaml>
scalebox validate --app-file=app.yaml --env-file=scalebox.env
```

## 1.10 vtask Subcommands

VTask (virtual task) is a cross-module collection of tasks within an App, implementing flow control and state management through the head-core-tail pipeline.

Common parameters:
- `--app-id`: app ID (can also be set via the `APP_ID` environment variable)

### 1.9.1 vtask list

List the app's vtasks.

```bash
scalebox vtask list --app-id <app-id>
```

### 1.9.2 vtask list-subtasks

List the sub-tasks of a specified vtask.

```bash
scalebox vtask list-subtasks --app-id <app-id> <vtask-id>
```

### 1.9.3 vtask get

View vtask details (including sub-task counts, semaphore names, etc.).

```bash
scalebox vtask get <vtask-id>
```

### 1.9.4 vtask fail

Mark a vtask failed. Automatically releases the gate semaphore and resource semaphores, and cascade-marks unfinished sub-tasks.

```bash
scalebox vtask fail <vtask-id>
```

### 1.9.5 vtask bind / unbind

Bind/unbind a vtask's compute resources.

```bash
scalebox vtask bind --app-id <app-id> --sema-name <sema-name>
scalebox vtask unbind --app-id <app-id> --sema-name <sema-name>
```

### 1.9.6 vtask add-subtask

Add a sub-task to a vtask. When running inside the agent, headers such as `_vtask_id` and `_vtask_size_sema` are automatically propagated.

```bash
scalebox vtask add-subtask --app-id <app-id> --module <module-name> --body <task-body>
```

### 1.9.7 vtask-Scoped Variables

```bash
# Get a vtask variable
scalebox vtask get-variable --vtask-id <vtask-id> <var-name>

# Set a vtask variable
scalebox vtask set-variable --vtask-id <vtask-id> <var-name> <value>
```

### 1.9.8 vtask-Scoped Semaphores

```bash
# Create
scalebox vtask create-semaphore --vtask-id <vtask-id> <sema-name> <initial-value>

# Get
scalebox vtask get-semaphore --vtask-id <vtask-id> <sema-name>

# Increment/decrement
scalebox vtask add-semaphore-value --vtask-id <vtask-id> <sema-name> <delta>

# Delete
scalebox vtask delete-semaphore --vtask-id <vtask-id> <sema-name>
```

## 1.11 semaphore Subcommands

- Common parameters: `--app-id` or `--module-id` (can also be set via the `APP_ID` / `MODULE_ID` environment variables)
- Environment variable: when `SEMAPHORE_AUTO_CREATE=yes`, semaphores that do not exist are automatically created (initial value 0). The CLI passes this to controld via the gRPC metadata `semaphore-auto-create: yes`

- Semaphore naming rules:
  - Character set: `[A-Za-z0-9:_-]`
  - First character is a letter or underscore

- Semaphore expressions: regular expressions representing a group of semaphores; the character set additionally includes `.*+?^$[]{}()|\`

### 1.11.1 semaphore list

List semaphores, supporting prefix and leaf-node filtering.

```bash
scalebox semaphore list --app-id <app-id>
scalebox semaphore list --app-id <app-id> --prefix vtask_size
scalebox semaphore list --app-id <app-id> --leaf-only
```

### 1.11.2 semaphore create

- Parameter: batch-size: used in batch semaphore creation to specify the batch size; default 100.

#### Creating a single semaphore
Example:
```sh
scalebox semaphore create ${sema_name} ${int_value}
scalebox semaphore create --app-id ${app_id} ${sema_name} ${int_value}
APP_ID=${app_id} scalebox semaphore create ${sema_name} ${int_value}

scalebox semaphore create --module-id=${module_id} ${sema_name} ${int_value}
MODULE_ID=${module_id} scalebox semaphore create ${sema_name} ${int_value}

```

#### Batch creation of semaphore groups
- Command line method: subject to the bash command line maximum length limit of 2MiB.
```sh
scalebox semaphore create '{"semaphores":{"sema1":n1,"sema2":n2}}'
```

- Semaphore file method: usually allows more semaphores
```sh
scalebox semaphore create --sema-file my-sema-file.txt
```

The semaphore file is a multi-line file format, with each line representing one semaphore.

```
"sema1":n1
"sema2":n2
"sema3":n3
```

### 1.11.3 semaphore get

#### Get the current value of a single semaphore
```sh
val=$(scalebox semaphore get ${sema_name})
code=$?
```
- ```code``` is the flag of operation success or failure.
  - 0: OK
  - 1: db error
  - 2: semaphore not-found
- ```val``` is the new semaphore value (integer)

If the environment variable SEMAPHORE_AUTO_CREATE=yes is set, the semaphore is automatically created with initial value 0

```sh
val=$(SEMAPHORE_AUTO_CREATE=yes scalebox semaphore get ${sema_name})
code=$?
```
- ```code``` is the flag of operation success or failure.
  - 0: OK
  - 1: db error
- ```val``` is the new semaphore value (integer)

####  Get the JSON key-value pairs of a semaphore group
Semaphore groups support generic matching of variable names with regular expressions.

```sh
val=$(scalebox semaphore get ${sema_expr} )
code=$?
```

- sema_expr is a regular expression
- ```code``` is the flag of operation success or failure. 0 means success
- ```val``` is the new semaphore value; if there are multiple semaphores, the returned result is a JSON map of semaphore name-value pairs.
  ```{"sema1":n1,"sema2":n2,"sema3":n3}```

### 1.11.4 semaphore increment

####  Incrementing a single semaphore by one
```sh
val=$(scalebox semaphore increment ${sema_name})
code=$?
```

- ```code``` is the flag of operation success or failure.
  - 0: OK
  - 1: db error
  - 2: semaphore not-found
- ```val``` is the new semaphore value (integer)

If the environment variable SEMAPHORE_AUTO_CREATE=yes is set, the semaphore is automatically created with initial value 0 and incremented by 1.

```sh
val=$(SEMAPHORE_AUTO_CREATE=yes scalebox semaphore increment ${sema_name})
code=$?
```
- ```code``` is the flag of operation success or failure.
  - 0: OK
- ```val``` is the new semaphore value (integer)

#### Incrementing a semaphore group by one

Semaphore groups support generic matching of variable names with regular expressions.

```sh
val=$(scalebox semaphore increment ${sema_expr} )
code=$?
```

- sema_expr is a regular expression
- ```code``` is the flag of operation success or failure. 0 means success
- ```val``` is the new semaphore value; if there are multiple semaphores, the returned result is a JSON map of semaphore name-value pairs.
  ```{"sema1":n1,"sema2":n2,"sema3":n3}```

### 1.11.5 semaphore decrement

#### Decrementing a single semaphore by one.

```sh
val=$(scalebox semaphore decrement ${sema_expr} )
code=$?
```

For usage details, see: <a href="#semaphore-increment">semaphore increment</a>

#### Decrementing a semaphore group by one.

For usage details, see: <a href="#semaphore-increment">semaphore increment</a>

### 1.11.6 semaphore add-value

#### Adding n to a single semaphore.
```sh
val=$(scalebox semaphore add-value ${sema_name} ${delta_value})
code=$?
```

For usage details, see: <a href="#semaphore-increment">semaphore increment</a>

#### Adding n to a semaphore group.

For usage details, see: <a href="#semaphore-increment">semaphore increment</a>


### 1.11.7 semaphore delete

Delete semaphores.

```sh
scalebox semaphore delete --app-id ${app_id} ${sema_name}
```

## 1.12 semagroup Subcommands

- Multiple semaphores form a semaphore group, identified by semaphore name prefixes and regular expressions


### 1.12.1 semagroup max
- The maximum value in the semaphore group
```sh
val=$(scalebox semagroup max ${sema_expr})
code=$?
```

### 1.12.2 semagroup min
- The minimum value in the semaphore group
```sh
val=$(scalebox semagroup min ${sema_expr})
code=$?
```
- sema_expr is a regular expression or prefix of semaphore names
- The return value val is an integer string

### 1.12.3 semagroup increment
- Select the minimum value in the semaphore group and increment it by one

### 1.12.4 semagroup decrement
- Select the maximum value in the semaphore group and decrement it by one

### 1.12.5 semagroup diff-max
- The difference between the maximum value of the semaphore group and the current value of the semaphore
```sh
val=$(scalebox semagroup diff-max ${sema_expr})
code=$?
```
- sema_expr is a semaphore with a group definition, e.g. ```(group-prefix):sema-suffix```
- The return value val is an integer string

### 1.12.6 semagroup diff-min
- The difference between the current value of the semaphore and the minimum value of the semaphore group
```sh
val=$(scalebox semagroup diff-min ${sema_expr})
code=$?
```
- sema_expr is a semaphore with a group definition, e.g. ```(group-part:)sema-suffix```
- The return value val is an integer string


## 1.13 variable Subcommands

- Common parameters: `--app-id` or `--module-id` (can also be set via the `APP_ID` / `MODULE_ID` environment variables)
- Variable naming: same as semaphore naming (`[A-Za-z0-9:_-]`, first character is a letter or underscore)
- Supports `list`, `get`, `set`, `delete` operations

### 1.16.1 variable list

List variables, supporting prefix and leaf-node filtering.

```bash
scalebox variable list --app-id <app-id>
scalebox variable list --app-id <app-id> --prefix my_prefix
scalebox variable list --app-id <app-id> --leaf-only
```

### 1.16.2 variable get

#### Get the current value of a single variable
- Example:
```sh
val=$(scalebox variable get ${var_name})
code=$?
[[ $code -ne 0 ]] && echo "[ERROR] variable-get ${var_name}, exit_code:$code" >&2
```
- ```code``` is the flag of operation success or failure.
  - 0: OK
  - 1: db error
  - 2: variable not-found
- ```val``` is the new variable value

####  Get the JSON key-value pairs of a variable group
Variable groups support generic matching of variable names with regular expressions.

```sh
val=$(scalebox variable get ${var_expr} )
code=$?
```

- var_expr is a regular expression
- ```code``` is the flag of operation success or failure. 0 means success
- ```val``` is the new variable value; the returned result is a JSON map of variable name-value pairs.
  ```{"var1":"val1","var2":"val2","var3":"val3"}```

### 1.15.3 variable set

```sh
scalebox variable set --app-id ${app_id} ${var_name} ${str_value}
APP_ID=${app_id} scalebox variable set ${var_name} ${str_value}
```

### 1.14.4 variable delete

```sh
scalebox variable delete --app-id ${app_id} ${var_name}
```

## 1.14 global Subcommands

Global variables, shared across apps.

### 1.16.1 global list

List global variables, supporting prefix and leaf-node filtering.

```bash
scalebox global list
scalebox global list --prefix my_prefix
scalebox global list --leaf-only
```

### 1.16.2 global get

```sh
scalebox global get ${global_name}
```

### 1.15.3 global set

```sh
scalebox global set ${global_name} ${global_value}
```

### 1.14.4 global delete

```sh
scalebox global delete ${global_name}
```


## 1.15 channel Subcommands

channel is used for cross-app communication and is a priority queue.

- Common parameters: module-id, or app-id
- Environment variables: MODULE_ID, or APP_ID

- Priority queue naming: same as semaphore naming

### 1.16.1 channel create

- head-app : app-id
- tail-app

If not specified, the current app is used

### 1.16.2 channel pull

- Get the current value of the queue
- Example:
```sh
val=$(scalebox channel pull ${pp_name})
code=$?
[[ $code -ne 0 ]] && echo "[ERROR] channel-pull ${pp_name}, exit_code:$code" >&2
```
- ```code``` is the flag of operation success or failure.
  - 0: OK
  - 1: db error
  - 2: channel not-found
- ```val``` is the new variable value

### 1.15.3 channel push

- priority is the priority, a floating-point number. Smaller values have higher priority.
  
```sh
scalebox channel push --app-id ${app_id} ${pp_name} ${str_value} [${priority}]
APP_ID=${app_id} scalebox channel push ${pp_name} ${str_value}

scalebox channel push --module-id ${module_id} ${pp_name} ${str_value} [${priority}]
MODULE_ID=${module_id} scalebox channel push ${pp_name} ${str_value}
```


## 1.16 fs Subcommands

scalebox-fs organizes files on distributed compute nodes into the same namespace in file system form. Later, features such as mount support and cross-node migration can be provided.

### 1.16.1 fs ls

- Main parameters:
  - include-removed-file
  - with-hostname
  - with-file-size
  - hostname=${host-name}

```sh
scalebox fs ls ${path_expr}
```

### 1.16.2 fs stat

View the metadata of one or more files. The file name on each node is consistent with the global file name.

The metadata mainly includes:
- virtual file name
- the host number the file belongs to
- creation time
- deletion time
- 
```sh
scalebox fs stat ${path_expr}
```

## 1.17 status

- Overall system status: local cluster head nodes (actuator to local head nodes is valid)
- cluster list: the number of hosts in different states
- app list: different states

## 1.18 event Subcommands

Supports add operations for various events.

- Basic commands: 

- xxxx is: "task"/"slot"/"misc"

```sh
scalebox event xxxx-add "${tag_name}" "${level_name}" ["${code}" ["${txt}" ["${json}"]]]
```
- If code is empty, the value 0 is used
- If txt is empty, the value "" is used
- If json is empty, the value "{}" is used

```sh
scalebox event xxxx-add --txt-file "${txt_file}" --json-file "${json_file}" "${tag_name}" "${level_name}" ["${code}"]
```

txt and json are read from files.


### 1.18.1 event task-add

Specify the task-id through the environment variable TASK_ID or the parameter --task-id.
```sh
scalebox event task-add --task-id ${task_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

```scalebox event task-add ``` can be abbreviated as ``` scalebox event add  ```

### 1.18.2 event slot-add

Specify the slot-id through the environment variable SLOT_ID or the parameter --slot-id.
```sh
scalebox event slot-add --slot-id ${slot_id} ${tag_name} ${level_name} ${code} ${txt} ${json}
```

### 1.18.3 event misc-add

```sh
scalebox event misc-add ${tag_name} ${level_name} ${code} ${txt} ${json}
```
