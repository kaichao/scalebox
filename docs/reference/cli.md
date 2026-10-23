<!-- AUTO-GENERATED from `scalebox <command> --help` (Cobra) by `make docs-gen` (go-scalebox); DO NOT EDIT -->
<!-- source: go-scalebox@3a8dd2e (2026-10-21) -->


# Scalebox CLI 参考

> 由 `scalebox <命令> --help` 递归生成；命令按 help 输出顺序。
> `completion`、`help` 是 Cobra 自动命令，未收录。

## `scalebox`

Scalebox command line tool

```text
Usage:
  scalebox [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `app` | Manage apps |
| `cluster` | Manage clusters |
| `global` | Manage global variables |
| `host` | Manage hosts |
| `mcp` | Model Context Protocol server for AI clients |
| `module` | Manage modules |
| `queue` | Cross-cluster message queue for management-plane data exchange |
| `run` | Run an application |
| `semagroup` | Manage semagroups |
| `semaphore` | Manage semaphores |
| `slot` | Manage slots |
| `ssh` | SSH to hosts matched by hostname (exact or regex) |
| `task` | Manage tasks |
| `validate` | Validate an app YAML file |
| `variable` | Manage application or module variables |
| `version` | Display version information |
| `vtask` | Manage vtasks |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |
| `-h, --help` | `—` | help for scalebox |

## `scalebox app`

Manage apps

```text
Usage:
  scalebox app [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `add-slots` | Create slots on hosts for a specific application |
| `delete` | Delete an app |
| `get-router` | Get app router |
| `list` | List apps |
| `set-status` | Set app status. Use FINISHED with optional result text (written to t_app.ext). |
| `show` | Show app details with module/slot/task summary |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-h, --help` | `—` | help for app |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox app add-slots`

Create slots on hosts for a specific application

```text
Usage:
  scalebox app add-slots <host_regex> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for add-slots |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox app delete`

Delete an app

```text
Usage:
  scalebox app delete [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox app get-router`

Get app router

```text
Usage:
  scalebox app get-router [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get-router |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox app list`

List apps

```text
Usage:
  scalebox app list [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for list |
| `--no-head` | `—` | skip header lines |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox app set-status`

Set app status. Use FINISHED with optional result text (written to t_app.ext).

```text
Usage:
  scalebox app set-status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-status |
| `--result` | `string` | Result text (written to t_app.ext when FINISHED) |
| `--status` | `string` | Target status (RUNNING/OFF/FINISHED) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox app show`

Show app details with module/slot/task summary

```text
Usage:
  scalebox app show [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for show |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster`

Manage clusters

```text
Usage:
  scalebox cluster [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `allocate` | Allocate dynamic hosts for a cluster |
| `create` | Create a cluster from YAML file |
| `get-parameter` | Get a cluster parameter |
| `init-local` | Initialize the local cluster (run on the head node's physical machine, not in a container) |
| `list` | List all clusters |
| `release` | Release dynamic hosts back to the scheduler |
| `renew` | Renew dynamic hosts for a cluster (rolling update) |
| `set-status` | Set cluster status |
| `status` | Show cluster status |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-h, --help` | `—` | help for cluster |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster allocate`

Allocate dynamic hosts for a cluster

```text
Usage:
  scalebox cluster allocate [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--account` | `string` | SLURM account (default from backend config) |
| `--cpu-only` | `string` | CPU-only nodes: yes/true (default from backend config) |
| `--group-spec` | `string` | Group spec: prefix:group_size:num_groups, e.g. d:24:2 |
| `-h, --help` | `—` | help for allocate |
| `--min-tmp-gb` | `int` | Min /tmp space in GB (default from backend config) |
| `--partition` | `string` | SLURM partition (default from backend config) |
| `--time-limit` | `string` | SLURM walltime (default from backend config) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster create`

Create a cluster from YAML file

```text
Usage:
  scalebox cluster create <yaml-file> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--admin-app` | `int` | Admin app ID (for dynamic clusters) |
| `-h, --help` | `—` | help for create |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster get-parameter`

Get a cluster parameter

```text
Usage:
  scalebox cluster get-parameter <parameter> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get-parameter |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster init-local`

Initialize the local cluster (run on the head node's physical machine, not in a container)

```text
Usage:
  scalebox cluster init-local [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--data-root` | `string` | Cluster data root directory (default "/tmp/scalebox/mydata") |
| `--grpc-server` | `string` | gRPC server address of the local cluster (default: <head-ip>:50051) |
| `--head-ip` | `string` | Head node IP address (default: auto-detected) |
| `-h, --help` | `—` | help for init-local |
| `--ssh-port` | `int` | SSH port of the head node (default 22) |
| `--username` | `string` | SSH username on the head node (default: current user) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster list`

List all clusters

```text
Usage:
  scalebox cluster list [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for list |
| `--no-head` | `—` | skip header lines |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster release`

Release dynamic hosts back to the scheduler

```text
Usage:
  scalebox cluster release [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--all` | `—` | Release all dynamic hosts |
| `--group-name` | `string` | Group name, e.g. d00 (format: prefix+index) |
| `-h, --help` | `—` | help for release |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster renew`

Renew dynamic hosts for a cluster (rolling update)

```text
Usage:
  scalebox cluster renew [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--all` | `—` | Renew all active allocation groups |
| `--group-name` | `string` | Group name, e.g. d00 |
| `-h, --help` | `—` | help for renew |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster set-status`

Set cluster status

```text
Usage:
  scalebox cluster set-status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-status |
| `--status` | `string` | Target status (ON/OFF/ERROR) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox cluster status`

Show cluster status

```text
Usage:
  scalebox cluster status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for status |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox global`

Manage global variables that are accessible across the entire system. Global variables are key-value pairs that can be stored and retrieved from anywhere in the system. They are useful for storing system-wide configuration values, shared state information, or any other data that needs to be globally accessible. Examples: # Get a global variable scalebox global get system_timeout # Set a global variable scalebox global set system_timeout "300"

```text
Usage:
  scalebox global [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `delete` | Delete a global variable |
| `get` | Get the value of a global variable |
| `ls` | List global variables with smart display |
| `set` | Set the value of a global variable |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for global |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox global delete`

Delete a global variable

```text
Usage:
  scalebox global delete <name> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox global get`

Get the value of a global variable

```text
Usage:
  scalebox global get <name> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox global ls`

List global variables with smart display

```text
Usage:
  scalebox global ls [prefix] [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for ls |
| `--leaf-only` | `—` | Show all entries flat, skipping tree-group display |
| `--no-header` | `—` | Do not print table header |
| `--no-truncate` | `—` | Do not truncate long name and value strings |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox global set`

Set the value of a global variable

```text
Usage:
  scalebox global set <name> <value> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host`

Manage hosts

```text
Usage:
  scalebox host [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `delete` | Delete a static host |
| `list` | List hosts in a cluster or by prefix |
| `new` | Add a new static host |
| `renew` | Renew a single host |
| `set-status` | Manually set a host's status |
| `show` | Show detailed information for a host |
| `status` | Show host status |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for host |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host delete`

Delete a static host

```text
Usage:
  scalebox host delete [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete |
| `--id` | `int` | Host ID |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host list`

List hosts in a cluster or by prefix

```text
Usage:
  scalebox host list [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `--groupID` | `string` | Group ID filter |
| `-h, --help` | `—` | help for list |
| `--no-head` | `—` | skip header lines |
| `--prefix` | `string` | Hostname prefix filter |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host new`

Add a new static host

```text
Usage:
  scalebox host new <hostname> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Cluster name |
| `-h, --help` | `—` | help for new |
| `--ip-addr` | `string` | IP address of the host |
| `--port` | `int` | SSH port (default 22) |
| `--uname` | `string` | SSH username (default root) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host renew`

Renew a single host

```text
Usage:
  scalebox host renew <hostname> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for renew |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host set-status`

Manually set a host's status

```text
Usage:
  scalebox host set-status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-status |
| `--id` | `int` | Host ID |
| `--status` | `string` | Target status |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host show`

Show detailed information for a host

```text
Usage:
  scalebox host show <hostname> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for show |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox host status`

Show host status

```text
Usage:
  scalebox host status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for status |
| `--prefix` | `string` | Hostname prefix filter |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox mcp`

Model Context Protocol server for AI clients

```text
Usage:
  scalebox mcp [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `serve` | Serve Scalebox read-only MCP tools over stdio |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for mcp |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox mcp serve`

Serve Scalebox read-only MCP tools over stdio

```text
Usage:
  scalebox mcp serve [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for serve |
| `--read-only` | `—` | expose read-only tools only (default true) |
| `--transport` | `string` | transport: stdio (default "stdio") |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox module`

Manage modules

```text
Usage:
  scalebox module [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `list` | List modules |
| `show` | Show module details |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-h, --help` | `—` | help for module |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox module list`

List modules

```text
Usage:
  scalebox module list [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID |
| `-h, --help` | `—` | help for list |
| `--no-head` | `—` | skip header lines |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox module show`

Show module details

```text
Usage:
  scalebox module show <module-name> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for show |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | App ID (or set APP_ID env) |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox queue`

Message queue for cross-cluster app communication. Queues live in the target cluster's database; --cluster selects a remote cluster (resolved from t_cluster, forwarded by the local controld). Dequeue order: priority DESC, id ASC (higher priority first, FIFO within the same priority). Pull removes messages; peek does not.

```text
Usage:
  scalebox queue [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `ls` | List queues with backlog stats |
| `peek` | Peek messages in a queue (keeps them) |
| `pull` | Pull messages from a queue (removes them) |
| `purge` | Remove all messages in a queue |
| `push` | Push a message into a queue |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for queue |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox queue ls`

List queues with backlog stats

```text
Usage:
  scalebox queue ls [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | remote cluster name (default: local) |
| `-h, --help` | `—` | help for ls |
| `--no-head` | `—` | skip header lines |
| `--prefix` | `string` | queue name prefix filter |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox queue peek`

Peek messages in a queue (keeps them)

```text
Usage:
  scalebox queue peek [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | remote cluster name (default: local) |
| `--count` | `int` | max messages to fetch (up to 1000) (default 1) |
| `-h, --help` | `—` | help for peek |
| `--json` | `—` | output full message metadata as JSON |
| `--queue` | `string` | queue name (required) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox queue pull`

Pull messages from a queue (removes them)

```text
Usage:
  scalebox queue pull [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | remote cluster name (default: local) |
| `--count` | `int` | max messages to fetch (up to 1000) (default 1) |
| `-h, --help` | `—` | help for pull |
| `--json` | `—` | output full message metadata as JSON |
| `--queue` | `string` | queue name (required) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox queue purge`

Remove all messages in a queue

```text
Usage:
  scalebox queue purge [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | remote cluster name (default: local) |
| `-h, --help` | `—` | help for purge |
| `--queue` | `string` | queue name (required) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox queue push`

Push a message into a queue. The body comes from --body, or from stdin when --body is omitted and stdin is a pipe/redirect (handy for multi-line content): scalebox queue push --queue=evt --body="done" cat result.json | scalebox queue push --queue=evt

```text
Usage:
  scalebox queue push [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--body` | `string` | message body (or pipe content via stdin) |
| `--cluster` | `string` | remote cluster name (default: local) |
| `-H, --header` | `strings` | custom headers in key=value format (can be specified multiple times) |
| `-h, --help` | `—` | help for push |
| `--priority` | `int` | message priority (higher dequeues first) |
| `--queue` | `string` | queue name (required) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox run`

Run an application. Examples: # Run with app.yaml file scalebox run --app-file app.yaml # Run with command line parameters scalebox run --app-cluster my-cluster --image-name my-image "start_task" # Print generated YAML without running scalebox run --print-yaml --app-cluster my-cluster "start_task"

```text
Usage:
  scalebox run [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-cluster` | `string` |  |
| `--app-file` | `string` |  |
| `--app-name` | `string` | application name (default: <dirname>_<timestamp>) |
| `--code-path` | `string` |  |
| `-h, --help` | `—` | help for run |
| `--image-name` | `string` |  |
| `--main-app-id` | `int` |  |
| `--main-cluster` | `string` | cluster name of the main app (for cross-cluster linking) |
| `--mr-code-path` | `string` |  |
| `--mr-image-name` | `string` |  |
| `--print-yaml` | `—` |  |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semagroup`

Manage semagroups

```text
Usage:
  scalebox semagroup [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `decrement` | Decrement semagroup |
| `diff-max` | Get semagroup diff-max value |
| `diff-min` | Get semagroup diff-min value |
| `increment` | Increment semagroup |
| `max` | Get semagroup maximum value |
| `min` | Get semagroup minimum value |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-h, --help` | `—` | help for semagroup |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semagroup decrement`

Decrement semagroup

```text
Usage:
  scalebox semagroup decrement [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for decrement |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semagroup diff-max`

Get semagroup diff-max value

```text
Usage:
  scalebox semagroup diff-max [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for diff-max |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semagroup diff-min`

Get semagroup diff-min value

```text
Usage:
  scalebox semagroup diff-min [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for diff-min |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semagroup increment`

Increment semagroup

```text
Usage:
  scalebox semagroup increment [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for increment |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semagroup max`

Get semagroup maximum value

```text
Usage:
  scalebox semagroup max [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for max |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semagroup min`

Get semagroup minimum value

```text
Usage:
  scalebox semagroup min [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for min |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore`

Manage semaphores

```text
Usage:
  scalebox semaphore [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `add-value` | Add a delta value to a semaphore |
| `create` | Create a semaphore |
| `decrement` | Decrement a semaphore by 1 |
| `delete` | Delete a semaphore |
| `get` | Get semaphore value |
| `increment` | Increment a semaphore by 1 |
| `ls` | List semaphores with smart display |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-h, --help` | `—` | help for semaphore |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore add-value`

Add a delta value to a semaphore

```text
Usage:
  scalebox semaphore add-value <name> <value> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for add-value |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore create`

Create a semaphore

```text
Usage:
  scalebox semaphore create [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--batch-size` | `int` | (default 100) |
| `-h, --help` | `—` | help for create |
| `--sema-file` | `string` | sema file path |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore decrement`

Decrement a semaphore by 1

```text
Usage:
  scalebox semaphore decrement [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for decrement |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore delete`

Delete a semaphore

```text
Usage:
  scalebox semaphore delete <name> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore get`

Get semaphore value

```text
Usage:
  scalebox semaphore get [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore increment`

Increment a semaphore by 1

```text
Usage:
  scalebox semaphore increment [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for increment |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox semaphore ls`

List semaphores with smart display

```text
Usage:
  scalebox semaphore ls [prefix] [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for ls |
| `--leaf-only` | `—` | Show all entries flat, skipping tree-group display |
| `--no-header` | `—` | Do not print table header |
| `--no-truncate` | `—` | Do not truncate long name and value strings |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot`

Manage slots

```text
Usage:
  scalebox slot [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `add` | Add READY slots for hosts matching a regex pattern |
| `delete` | Remove slots for hosts matching a regex pattern |
| `get-parameter` | Get a single slot parameter |
| `list` | List all slots for a specific module |
| `set-parameter` | Set a single slot parameter |
| `set-status` | Manually set a slot's status |
| `show` | Show detailed information for a single slot |
| `status` | Show slot status summary for a module |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for slot |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot add`

Add READY slots for hosts matching a regex pattern

```text
Usage:
  scalebox slot add [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for add |
| `--host-regex` | `string` | host-regex for matching hosts |
| `--module-id` | `int` | module-id |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot delete`

Remove slots for hosts matching a regex pattern

```text
Usage:
  scalebox slot delete [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete |
| `--host-regex` | `string` | host-regex for matching hosts |
| `--module-id` | `int` | module-id |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot get-parameter`

Get a single slot parameter

```text
Usage:
  scalebox slot get-parameter <parameter> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get-parameter |
| `--slot-id` | `int` | slot-id |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot list`

List all slots for a specific module

```text
Usage:
  scalebox slot list [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for list |
| `--module-id` | `int` | module-id |
| `--no-head` | `—` | skip header lines |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot set-parameter`

Set a single slot parameter

```text
Usage:
  scalebox slot set-parameter <parameter> <value> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-parameter |
| `--slot-id` | `int` | slot-id |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot set-status`

Manually set a slot's status

```text
Usage:
  scalebox slot set-status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-status |
| `--slot-id` | `int` | slot-id |
| `--status` | `string` | target status (READY, OFF) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot show`

Show detailed information for a single slot

```text
Usage:
  scalebox slot show [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for show |
| `--slot-id` | `int` | slot-id |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox slot status`

Show slot status summary for a module

```text
Usage:
  scalebox slot status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for status |
| `--module-id` | `int` | module-id |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox ssh`

SSH to hosts matched by hostname. <host-expr> is matched against t_host.hostname (status must be ON): exact match first, Go regexp partial match as fallback. Without a command an interactive session is opened (requires exactly one matching host); with a command it runs on every matched host: sequentially by default (each host's output stays grouped), or concurrently with --parallel. Port: --port wins; otherwise t_host.port and 22 are TCP-probed in order. Unknown host keys are accepted automatically (StrictHostKeyChecking=accept-new); keys that changed after a first connection are still rejected.

```text
Usage:
  scalebox ssh <host-expr> [command...] [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--cluster` | `string` | Limit host match to a cluster |
| `--concurrency` | `int` | Max parallel SSH sessions with --parallel (0 = unlimited) (default 5) |
| `-h, --help` | `—` | help for ssh |
| `--ip` | `—` | Connect to ip_addr instead of hostname |
| `--parallel` | `—` | Run on matched hosts concurrently (output may interleave) |
| `--port` | `int` | SSH port (skip t_host.port/22 auto-detect) |
| `--user` | `string` | SSH user (default t_host.uname, fallback root) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task`

Manage tasks

```text
Usage:
  scalebox task [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `add` | Add a task |
| `delete` | Delete a task |
| `get` | Get task details |
| `get-header` | Get a task header |
| `list` | List tasks |
| `remove-header` | Remove a task header |
| `set-header` | Set a task header |
| `show` | Show task details |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for task |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task add`

Add a task

```text
Usage:
  scalebox task add [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` | app ID |
| `--conflict-action` | `string` | conflict-action |
| `--direct-write-in-agent` | `—` | skip in-agent logic |
| `--from-module` | `string` | from module |
| `-H, --header` | `strings` | Custom headers in key=value format (can be specified multiple times) |
| `--headers` | `string` | headers in JSON format |
| `-h, --help` | `—` | help for add |
| `--module-id` | `int` | module ID |
| `--remote-cluster` | `string` | remote cluster |
| `--sink-module` | `string` | sink module name |
| `--task-file` | `string` | task file path |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task delete`

Delete a task

```text
Usage:
  scalebox task delete <task-id> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task get`

Get task details

```text
Usage:
  scalebox task get <task-id> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task get-header`

Get a task header

```text
Usage:
  scalebox task get-header [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get-header |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task list`

List tasks

```text
Usage:
  scalebox task list [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-h, --help` | `—` | help for list |
| `--limit` | `int` | max items (default 10) |
| `--module-id` | `int` |  |
| `--no-head` | `—` | skip header lines |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task remove-header`

Remove a task header

```text
Usage:
  scalebox task remove-header [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for remove-header |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task set-header`

Set a task header

```text
Usage:
  scalebox task set-header [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-header |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox task show`

Show task details

```text
Usage:
  scalebox task show <task-id> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for show |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox validate`

Validate an app YAML file

```text
Usage:
  scalebox validate [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-file` | `string` | Path to app YAML file |
| `-h, --help` | `—` | help for validate |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox variable`

Manage variables at the application or module level. Variables are key-value pairs that can be stored and retrieved for specific applications or modules. They are useful for storing configuration values, state information, or any other data that needs to be persisted. Examples: # Get a variable for application with ID 1 scalebox variable get my_var --app-id=1

```text
Usage:
  scalebox variable [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `delete` | Delete a variable |
| `get` | Get the value of a variable |
| `ls` | List variables with smart display |
| `set` | Set the value of a variable |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-h, --help` | `—` | help for variable |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox variable delete`

Delete a variable

```text
Usage:
  scalebox variable delete <name> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox variable get`

Get the value of a variable

```text
Usage:
  scalebox variable get <name> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox variable ls`

List variables with smart display

```text
Usage:
  scalebox variable ls [prefix] [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for ls |
| `--leaf-only` | `—` | Show all entries flat, skipping tree-group display |
| `--no-header` | `—` | Do not print table header |
| `--no-truncate` | `—` | Do not truncate long name and value strings |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox variable set`

Set the value of a variable

```text
Usage:
  scalebox variable set <name> <value> [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox version`

Display version information

```text
Usage:
  scalebox version [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for version |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox vtask`

Manage vtasks

```text
Usage:
  scalebox vtask [command]
```

**子命令**

| 命令 | 说明 |
|---|---|
| `add-semaphore-value` | Add value to vtask semaphore |
| `add-subtask` | Add a subtask to an existing vtask |
| `bind` | Bind a vtask to a compute resource group |
| `create-semaphore` | Create a vtask semaphore |
| `delete-semaphore` | Delete a vtask semaphore |
| `fail` | Force-terminate a vtask and release its resources |
| `get` | Show detailed information for a vtask |
| `get-semaphore` | Get vtask semaphore value |
| `get-variable` | Get vtask variable value |
| `list` | List vtasks for an app |
| `list-subtasks` | List subtasks for a given vtask |
| `set-status` | Set vtask status |
| `set-variable` | Set vtask variable value |
| `unbind` | Unbind a vtask's bound compute resource |

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-h, --help` | `—` | help for vtask |
| `--vtask-id` | `int` | vtask ID |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-e, --env-file` | `string` | Environment variables file |

## `scalebox vtask add-semaphore-value`

Add value to vtask semaphore

```text
Usage:
  scalebox vtask add-semaphore-value [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for add-semaphore-value |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask add-subtask`

Add a subtask to an existing vtask

```text
Usage:
  scalebox vtask add-subtask [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--direct` | `—` | force gRPC direct mode (bypass sink-tasks.txt even in agent) |
| `--header` | `stringToString` | headers (repeatable, key=value) (default []) |
| `-h, --help` | `—` | help for add-subtask |
| `--module` | `string` | target module name (required) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask bind`

Bind a vtask to a compute resource group

```text
Usage:
  scalebox vtask bind [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for bind |
| `--module` | `string` | target head module (optional, auto-detected) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask create-semaphore`

Create a vtask semaphore

```text
Usage:
  scalebox vtask create-semaphore [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for create-semaphore |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask delete-semaphore`

Delete a vtask semaphore

```text
Usage:
  scalebox vtask delete-semaphore [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for delete-semaphore |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask fail`

Force-terminate a vtask and release its resources

```text
Usage:
  scalebox vtask fail [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for fail |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask get`

Show detailed information for a vtask

```text
Usage:
  scalebox vtask get [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask get-semaphore`

Get vtask semaphore value

```text
Usage:
  scalebox vtask get-semaphore [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get-semaphore |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask get-variable`

Get vtask variable value

```text
Usage:
  scalebox vtask get-variable [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for get-variable |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask list`

List vtasks for an app

```text
Usage:
  scalebox vtask list [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for list |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask list-subtasks`

List subtasks for a given vtask

```text
Usage:
  scalebox vtask list-subtasks [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for list-subtasks |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask set-status`

Set vtask status

```text
Usage:
  scalebox vtask set-status [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-status |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask set-variable`

Set vtask variable value

```text
Usage:
  scalebox vtask set-variable [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for set-variable |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
| `--vtask-id` | `int` | vtask ID |

## `scalebox vtask unbind`

Unbind a vtask's bound compute resource

```text
Usage:
  scalebox vtask unbind [flags]
```

**Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `-h, --help` | `—` | help for unbind |
| `--sema-name` | `string` | full semaphore name |
| `--vtask-id` | `int` | vtask ID (auto-look up _vtask_size_sema) |

**Global Flags**

| Flag | 类型 | 说明 |
|---|---|---|
| `--app-id` | `int` |  |
| `-e, --env-file` | `string` | Environment variables file |
