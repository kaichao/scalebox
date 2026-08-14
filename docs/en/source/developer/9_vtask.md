# 9. VTask Virtual Tasks

**VTask** (Virtual Task) is an application-level coarse-grained compute unit that Scalebox establishes on top of fine-grained tasks. A vtask is a cross-module collection of tasks with built-in semaphores and shared variables (variables), composed of wait-queue (optional), vtask-head, one or more cascaded vtask-core stages, and vtask-tail forming a pipeline that uniformly manages flow control, resource binding, state tracking, and full-chain context propagation.

## 9.1 Why VTask Is Needed

### 9.1.1 Limitations of the task Model

Scalebox's base model is module → slot → task. Each task is an independent execution unit that works well at the fine-grained scheduling level, but has three gaps at the application level:

1. **Lack of cross-module state management**: when multiple modules within the same App collaborate, the tasks in each module are isolated from each other, with no concept of "this batch of tasks belongs to the same computational job"
2. **Cannot answer application-level questions**: "how far has a computational job progressed", "how many sub-tasks in total, how many completed", "which resources should be released on failure"
3. **Lack of resource isolation**: no resource reservation mechanism when different computational jobs execute on the same node

### 9.1.2 VTask's Solution

VTask overlays cross-module state management on top of tasks:

```
Fine-grained layer (task):   single execution unit, slot scheduling, status codes -1/-2/-3/0/>0
Coarse-grained layer (vtask):   cross-module task collection, flow control + state + fault tolerance + resource binding
```

| Capability | Implementation | Effect |
|------|------|------|
| Cross-module association | `vtask` column + `_vtask_id` header linking head→core→tail | Scripts don't need to manually maintain task references |
| Flow control | `vtask_size` semaphore limits concurrency | No extra modules needed for queuing logic |
| State tracking | `vtask get` / `vtask list` show progress in real time | One command replaces manual task table queries |
| Coarse-grained fault tolerance | `vtask fail` terminates the pipeline + releases resources | Retry the whole vtask; no per-task cleanup |
| Resource binding | bind/unbind reserves/returns compute nodes | Declarative management, reuse instead of hardcoding |
| Context propagation | `_`-prefixed headers auto-propagate across the whole chain | Intermediate modules need zero code to sense context origins |

## 9.2 Conceptual Model

### 9.2.1 Three Types

| | DEFAULT | HOST-BOUND | GROUP-BOUND |
|----|---------|------------|-------------|
| Resource binding | None | Single node | Node group |
| wait-queue | No | Yes | Yes |
| Flow control semaphore | `vtask_size:<mod>` | `host_vtask_size:<mod>:<host>` | `slot_vtask_size:<mod>:<seq>` |
| vtask-head mode | no task_dist_mode | `HOST-BOUND` | `SLOT-BOUND` |
| vtask-core mode | no task_dist_mode | `HOST-BOUND` | `HOST-BOUND` |
| Applicable scenarios | Lightweight batch processing | Exclusive computing on a single node | Multi-node collaborative computing |

**Pipeline structures** (`vtask-core` can cascade multiple stages):

```
DEFAULT:       task → vtask-head → vtask-core-1 → vtask-core-2 → ... → vtask-tail
HOST-BOUND:    task → wait-queue → vtask-head → vtask-core → vtask-tail
GROUP-BOUND:   task → wait-queue → vtask-head → vtask-core-1 → vtask-core-2 → ... → vtask-tail
```

### 9.2.2 Capability Matrix

| Capability | DEFAULT | HOST-BOUND | GROUP-BOUND |
|------|:--:|:--:|:--:|
| Flow control + tracking + fault tolerance | ✅ | ✅ | ✅ |
| header DAG propagation | ✅ | ✅ | ✅ |
| Deterministic performance (node exclusive) | ❌ | ✅ | ✅ |
| Local I/O (NVMe) | ❌ | ✅ | ✅ |
| Fault isolation | ❌ | ✅ | ✅ |
| vtask multi-node parallelism | ❌ | ❌ | ✅ |

| Type | Positioning | Typical scenarios |
|------|------|---------|
| DEFAULT | Agile mode | Log aggregation, config distribution, batch thumbnails, data format conversion |
| HOST-BOUND | Performance mode | Single-node GPU inference, large-scale data transformation |
| GROUP-BOUND | Scale mode | Multi-channel beamforming, multi-node collaborative data analysis |

### 9.2.3 Two-Level Abstraction

| Level | Concept | Concern |
|------|------|--------|
| **vtask** (application layer) | Resource binding form: HOST-BOUND / GROUP-BOUND / DEFAULT | User view: runs on a single node or a node group |
| **module** (underlying) | Task routing method: HOST-BOUND / SLOT-BOUND / DEFAULT | Framework view: how tasks route to specific slots |

`SLOT-BOUND` is a module-level implementation detail, used only on GROUP-BOUND head modules.

### 9.2.4 Root Task Identification

The task of a `vtask_role=head` module is the vtask root. When controld creates a head module task, it automatically establishes a self-reference:

```text
UPDATE t_task SET vtask = $1,
    headers = jsonb_set(headers, '{_vtask_id}', to_jsonb($1::text), true)
WHERE id = $1
```

The `ListVTasks` query condition is `WHERE vtask = id`.

### 9.2.5 State Derivation

The vtask state is derived from the root task's `status_code` + the completion ratio of sub-tasks:

| Root task status | Sub-tasks | Derived state |
|:---:|------|------|
| `-1` | 0 | READY |
| `-2` | 0 | QUEUED |
| `-3` | 0 | RUNNING |
| `-3` | >0, partially complete | RUNNING (X/Y finished) |
| `0` | all complete | FINISHED (Y tasks) |
| `1` | — | FAILED (exit=N) |

## 9.3 Module Structure

### 9.3.1 DEFAULT (4 modules)

```yaml
modules:
  main-router:   # entry router: task add --sink-module=vtask-head $body
  vtask-head:    # vtask_role: head, vtask_size: 2
  vtask-core:    # vtask_role: core
  vtask-tail:    # vtask_role: tail (no unbind)
```

### 9.3.2 HOST-BOUND (5 modules)

```yaml
modules:
  main-router:   # task add --sink-module=wait-queue $body
  wait-queue:    # vtask_size: 1 (serialization gate)
  vtask-head:    # vtask_role: head, task_dist_mode: HOST-BOUND
  vtask-core:    # vtask_role: core, task_dist_mode: HOST-BOUND
  vtask-tail:    # vtask_role: tail
```

**Key script operations**:

| Script | Operation |
|------|------|
| `from-wait-queue.sh` | `vtask bind` → `vtask add-subtask --direct` (on failure, unbind + gate increment rollback) |
| `from-vtask-head.sh` | `vtask add-subtask --module=vtask-core` (async), finally `semaphore increment vtask_size:wait-queue` |
| `from-vtask-core.sh` | `task add --sink-module=vtask-tail $1` |
| `from-vtask-tail.sh` | case match `host_vtask_size:*` → `vtask unbind --sema-name=":${sema}"` |

### 9.3.3 GROUP-BOUND

Same structure as HOST-BOUND, with differences:
- `vtask-head.task_dist_mode: SLOT-BOUND`, bind returns slot_seq
- Flow control semaphore format `slot_vtask_size`, head slot serves as the node group anchor
- `from-vtask-head.sh` computes `to_host` from the body (business logic)

### 9.3.4 vtask-core Cascading

`vtask-core` can contain multiple computation stages:

```yaml
modules:
  vtask-core-beamforming:     # data distribution / preprocessing
    vtask_role: core
    task_dist_mode: HOST-BOUND
  vtask-core-cross-correlate: # core computation
    vtask_role: core
    task_dist_mode: HOST-BOUND
  vtask-core-merge:           # result merging
    vtask_role: core
    task_dist_mode: DEFAULT
```

Cascading rules:
- All core modules share `vtask_role: core`; controld handles them with unified logic
- `_`-prefixed headers (`_vtask_id`, `_vtask_size_sema`) auto-propagate along the core chain
- Each core module independently configures slot counts, `task_dist_mode`, and images

## 9.4 The Semaphore Mechanism

### 9.4.1 Naming Conventions

| Distribution mode | Format | Example |
|---------|------|------|
| DEFAULT | `vtask_size:<module>` | `vtask_size:vtask-head` |
| HOST-BOUND | `host_vtask_size:<module>:<host>` | `host_vtask_size:vtask-head:n0-0` |
| SLOT-BOUND | `slot_vtask_size:<module>:<seq>` | `slot_vtask_size:vtask-head:2` |

### 9.4.2 Dual Semaphores

Each flow control semaphore has two versions:

| Version | Prefix | Managed by | Purpose |
|------|------|--------|------|
| Flow control version | no prefix | controld automatic | `getVtaskBatchSize` / `updateVTask` / `doVTaskFinished` |
| Programmable version | `:` prefix | scripts manual | `vtask bind` / `unbind` operations |

Set `vtask_size_sema_copy: yes` when creating the app; both versions are established simultaneously with the same initial value.

### 9.4.3 Flow Control Process

```
vtask-head task picked up by a slot
  → updateVTask: semaphore decrement (flow control version)
  → group aggregation by _vtask_size_sema (supports multi-slot same-group scenarios)
  → vtask pipeline executes
  → vtask-tail completes
  → doVTaskFinished: semaphore increment (flow control version, quota restored)
  → from-vtask-tail.sh: vtask unbind (programmable version, HOST-BOUND / SLOT-BOUND only)
```

> `updateVTask` aggregates decrements grouped by the task's `_vtask_size_sema` header, avoiding inconsistent decrement/increment names in multi-slot scenarios.

### 9.4.4 Wait Queue Gating

The wait-queue module's `vtask_size` stays at 1 (serialization). `from-vtask-head.sh` manually increments at the end to let the next vtask through.

## 9.5 Pipeline Flow

### 9.5.1 Header Propagation

```
main-router task add
  → agent addSinkTasks: _vtask_* gating (only sink modules with vtask_role≠"" pass)
  → wait-queue task: does not receive _vtask_id / _vtask_size_sema

wait-queue task executes from-wait-queue.sh
  → --direct gRPC creates the vtask-head task
  → vtask_role=head → self-ref vtask=id
  → _vtask_size_sema explicitly set by the script

vtask-head task executes from-vtask-head.sh
  → sink-tasks.txt async creates the vtask-core task
  → agent addSinkTasks: _-prefixed headers auto-propagate

vtask-core → vtask-tail (same path as above)

vtask-tail completes
  → doVTaskFinished: reads _vtask_size_sema → increments the flow control version
  → from-vtask-tail.sh: reads _vtask_size_sema → unbinds the programmable version
```

### 9.5.2 Lifecycle

```
task creation → ON CONFLICT dedup
  ↓
task scheduling → FOR UPDATE SKIP LOCKED (single slot exclusive)
  ↓
task execution → agent → status_code: RUNNING → OK/ERROR
  ↓  _vtask_id / _vtask_size_sema auto-propagate via the _ prefix
  ↓
tail completes → doVTaskFinished → semaphore increment
  ↓  headers.vtask_tail_id = tail task ID
  ↓
from-vtask-tail.sh → unbind (HOST-BOUND / SLOT-BOUND only)
```

### 9.5.3 add-subtask Dual Paths

When `vtask add-subtask` runs inside the agent container, it takes the sink-tasks.txt async path by default; adding `--direct` switches to the gRPC sync path:

| | Async (sink-tasks.txt) | Sync (gRPC `--direct`) |
|--|:---:|:---:|
| Execution timing | agent post-processing (batch-sent after task execution completes) | synchronous call within the script |
| Error awareness | script unaware (write success → exit 0) | script learns immediately (exit code ≠ 0) |
| `_`-prefixed headers | agent auto-propagates | caller explicitly passes `--header` |
| `from_module`/`from_ip` | agent injects (more accurate) | CLI injects from env |

**Selection principles**:

| Script | Path | Reason |
|------|------|------|
| `from-vtask-head.sh` | async (default) | no resource-binding prerequisite, enjoys header auto-propagation |
| `from-wait-queue.sh` | sync (`--direct`) | bind has already deducted resources; must sense success/failure synchronously, rollback on failure |
| Ops manual CLI | sync (gRPC) | no sink-tasks.txt mechanism outside the agent |

## 9.6 Command Reference

### 9.6.1 Three-Layer Operation System

```
Layer 1: basic resource operations (global scope)
  semaphore create / get / increment / decrement / list / delete
  semagroup min / max / diffmin / diffmax / increment / decrement
  variable get / set / delete / list
  global get / set / delete / list

Layer 2: vtask-scoped resource operations
  vtask get-variable / set-variable
  vtask create-semaphore / get-semaphore / add-semaphore-value / delete-semaphore

Layer 3: vtask advanced operations
  vtask bind / unbind / add-subtask / get / fail / list / list-subtasks
```

### 9.6.2 vtask list / list-subtasks

```bash
# List all vtasks of the app (with derived state and sub-task counts)
scalebox vtask list --app-id=77

# List the sub-tasks of a vtask
scalebox vtask list-subtasks --vtask-id=1964
```

### 9.6.3 vtask get

```bash
scalebox vtask get --vtask-id=42
# ID:        42
# App ID:    5
# Status:    RUNNING (2/3 finished)
# Sema:      host_vtask_size:vtask-head:n0-0
# Subtasks:  3 total, 2 finished
```

### 9.6.4 vtask fail

```bash
scalebox vtask fail --vtask-id=42
```

Execution logic:
1. Atomic `UPDATE status_code='1' WHERE id=$1 AND status_code IN ('-1','-2','-3')`
2. `RowsAffected==0` → already finished, return user_msg
3. Release the wait-queue gate semaphore
4. Release resource semaphores (`semaphore increment :<sema_name>`)
5. Cascade-mark unfinished sub-tasks

### 9.6.5 vtask bind / unbind

```bash
# Bind (auto-finds the head module, returns to_host or slot_seq)
to_host=$(scalebox vtask bind) || exit 1

# Unbind (semaphore name passed by the tail script)
scalebox vtask unbind --sema-name=":host_vtask_size:vtask-head:n0-0"

# Unbind (ops passes vtask-id; server auto-queries _vtask_size_sema)
scalebox vtask unbind --vtask-id=42
```

| mode | bind returns |
|------|----------|
| HOST-BOUND | hostname (e.g. `n0-0`) |
| SLOT-BOUND | slot_seq (e.g. `0`) |
| DEFAULT | empty (no-op) |

### 9.6.6 vtask add-subtask

```bash
# Inside agent (common): async path, _-prefixed headers auto-propagate
scalebox vtask add-subtask --module=vtask-core --header to_ip=$from_ip $body

# Inside agent, needing synchronous error awareness: add --direct
scalebox vtask add-subtask --direct --module=vtask-head \
    --header to_host=$to_host --header _vtask_size_sema=$sema_name $body || {
    scalebox vtask unbind --sema-name=":${sema_name}"
    scalebox semaphore increment vtask_size:wait-queue
    exit 1
}

# Ops manual, outside the agent: must explicitly pass --app-id and _vtask_id
scalebox vtask add-subtask --app-id=10 --module=vtask-core \
    --header _vtask_id=42 $body
```

Auto-injection: `from_module` (`PLAT_MODULE_NAME` env), `from_ip` (local IPv4). Not overwritten when the caller explicitly passes them.

### 9.6.7 vtask-Scoped Variables/Semaphores

```bash
# vtask-scoped variables
scalebox vtask get-variable --vtask-id <vtask-id> <var-name>
scalebox vtask set-variable --vtask-id <vtask-id> <var-name> <value>

# vtask-scoped semaphores
scalebox vtask create-semaphore --vtask-id <vtask-id> <sema-name> <initial-value>
scalebox vtask get-semaphore --vtask-id <vtask-id> <sema-name>
scalebox vtask add-semaphore-value --vtask-id <vtask-id> <sema-name> <delta>
scalebox vtask delete-semaphore --vtask-id <vtask-id> <sema-name>
```

## 9.7 Concurrency Safety

| Mechanism | Protection scenario |
|------|---------|
| `FOR UPDATE SKIP LOCKED` | slots never take the same task |
| `ON CONFLICT` | tasks with the same body are not created twice |
| Deadlock retry ×3 | `40P01` deadlock self-healing |
| asyncbatch (10ms/50ms) | semaphore batch writes don't compete per record |
| `vtask_size=1` (wait-queue) | bind operations serialized |

## 9.8 Ops Scenario Quick Reference

| Scenario | Command |
|------|------|
| View global vtask progress | `vtask list --app-id=N` |
| View a single vtask's details | `vtask get --vtask-id=N` |
| View vtask sub-task distribution | `vtask list-subtasks --vtask-id=N` |
| Terminate a single vtask | `vtask fail --vtask-id=N` |
| Batch termination | `vtask list \| awk ... \| xargs vtask fail` |
| Pause global flow control | `semaphore decrement vtask_size:<mod> N` |
| Manually release leaked resources | `vtask unbind --sema-name=":host_vtask_size:..."` |
| Manually return (semaphore name unknown) | `vtask unbind --vtask-id=N` |

## 9.9 Implementation Index

| Layer | File | Responsibility |
|----|------|------|
| Proto | `proto/app.proto` | 5 vtask RPCs + 8 messages |
| Server | `cmd/controld/grpc/vtask_ops.go` | bind / unbind / add-subtask / get / fail handlers |
| Server | `cmd/controld/grpc/app_service.go` | ListVTasks / ListVTaskSubtasks handlers |
| Server | `cmd/controld/task/vtask.go` | `getVtaskBatchSize` / `updateVTask` / `doVTaskFinished` |
| CLI | `cmd/cli/vtask/vtask.go` | 14 subcommands + add-subtask async dual paths |
| Agent | `cmd/agent/sink_tasks.go` | `addSinkTasks`: `_`-prefixed header propagation + vtask_role gating |
| SDK | `pkg/vtask/` | Go wrappers |
| Example | `examples/vtask/main-router/` | 5 scripts |

> For test details and known limitations, see :doc:`VTask Design Review <../../go-scalebox/docs/vtask-review>` (in English).
