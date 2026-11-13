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
| Flow control semaphore | `vtask_size:<mod>` | `host_vtask_size:<mod>:<host>` | `group_vtask_size:<mod>:<group>` |
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

`SLOT-BOUND` is a module-level implementation detail, used only on GROUP-BOUND head modules; there is no such thing as a "pure SLOT-BOUND vtask".

**Environment variable layering** (one-to-one with the two-level abstraction):

| Variable | Level | Values | Consumers |
|------|------|------|--------|
| `VTASK_MODE` | vtask application layer | `DEFAULT` / `HOST-BOUND` / `GROUP-BOUND` | Application scripts (main-router routing, wait-queue standard module) |
| `TASK_DIST_MODE` | Platform implementation layer | `HOST-BOUND` / `SLOT-BOUND` / empty | Only drives the vtask-head module's `task_dist_mode` platform parameter |

The mapping converges at a single place in the app env file (e.g. `group-bound.env`):

```bash
VTASK_MODE=GROUP-BOUND
TASK_DIST_MODE=SLOT-BOUND   # only drives vtask-head's platform parameter, implementation-layer enum
```

Application authors only touch `VTASK_MODE` at the vtask layer; `SLOT-BOUND` appears only at the implementation layer.

### 9.2.4 Compute Resource Identification

The resource allocated by `vtask bind` is essentially a **placement** — the vtask anchors to a host (HOST-BOUND) or a node group (GROUP-BOUND), after which its tasks are routed via `to_host` / `to_slot` to slots at that location. The platform has no notion of physical capacity (CPU cores, memory); **capacity is configured by the application author via the `vtask_size` parameter** (the number of in-flight vtasks the location can carry simultaneously).

**Resource pool source**: available resources = the expansion of the head module (`vtask_role=head`) slots:

| vtask type | Resource pool | Identifier | Example |
|-----------|--------|------|------|
| HOST-BOUND | each host expanded from the head module's `slots` regex (with `status='ON'` in `t_host`) | short hostname (cluster suffix removed) | `n0-0` |
| GROUP-BOUND | each slot of the head module (slot list length = group count) | group number = the head slot's `seq` | `0` |

Semaphores are created at app creation and re-created on dynamic slot expansion, with initial value = `vtask_size`.

**Identifier system** (semaphore naming is the identifier):

| Type | Semaphore name | Identifier part |
|------|---------|---------|
| DEFAULT | `vtask_size:<module>` | none (global) |
| HOST-BOUND | `host_vtask_size:<module>:<hostname>` | short hostname |
| GROUP-BOUND | `group_vtask_size:<module>:<group number>` | head slot's seq |

The resource returned by `vtask bind` is the semaphore name suffix; `vtask get` shows `host=n0-0` / `group=2`.

**Node group membership has no platform-level definition**: which nodes belong to a group and how tasks are distributed within it are application business logic (vtask-head script routing + vtask-core deployed HOST-BOUND on the group's nodes). The platform only manages the "group" anchor dimension.

**Resource acquisition and release**: acquisition = `vtask bind` atomic decrement + `_vtask_head_module` + `_vtask_host` (HOST-BOUND) or `_vtask_group` (GROUP-BOUND) propagated with the task chain; release = `doVTaskFinished` / `FailVtask` increments by `_vtask_size_sema` or the rebuilt value (`vtask unbind` for error rollback). The number of in-flight vtasks per host/group is ≤ `vtask_size`.

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
  wait-queue:    # standard module scalebox.net/platform/wait-queue (vtask_size >= slot count: concurrent bind admission + check.sh admission + run.sh bind)
  vtask-head:    # vtask_role: head, task_dist_mode: HOST-BOUND, pod_id optional (non-vtask fallback)
  vtask-core:    # vtask_role: core, task_dist_mode: HOST-BOUND, prefer _vtask_host
  vtask-tail:    # vtask_role: tail
```

**Key script operations**:

| Script | Location | Operation |
|------|------|------|
| `check.sh` / `run.sh` | wait-queue standard module image | Admission gate (ACTION_CHECK) + `vtask bind` → `vtask add-subtask --direct` (on failure, unbind + gate increment rollback); task terminates at this slot, no author scripting needed |
| `from-vtask-head.sh` | main-router/ | `vtask add-subtask --module=vtask-core` (async; HOST-BOUND omits `to_ip`/`to_host` and uses `_vtask_host`, with `pod_id` as the legacy/non-vtask fallback; GROUP-BOUND must set `to_host` explicitly); **no +1 gate** — wait-queue `run.sh` releases it after bind+delivery (2026-11-02 L1) |
| `from-vtask-core.sh` | main-router/ | `task add --sink-module=vtask-tail $1` |
| `from-vtask-tail.sh` | main-router/ | no-op (quota restored automatically by controld `doVTaskFinished`) |

> **2026-11-10 anchor change (updated 2026-11-12)**: wait-queue writes `_vtask_head_module` + `_vtask_host` (HOST-BOUND) or `_vtask_group` (GROUP-BOUND); HOST-BOUND head→core no longer depends on `pod_id`; the GROUP-BOUND first hop must set `to_host` explicitly, core→core inherits `_vtask_host`, and `pod_id` remains the fallback. See §9.3.6.

> **Tail delivery contract**: `vtask-tail` must be delivered once after all core tasks complete.
> The platform does not implicitly ignore duplicate deliveries; a duplicate tail identity rolls back
> the whole batch (112) and may finish the vtask early if the first tail runs before all shards finish.

### 9.3.3 GROUP-BOUND

Same structure as HOST-BOUND, with differences:
- `vtask-head.task_dist_mode: SLOT-BOUND` (platform-level value), bind returns the group number (i.e., the head slot's seq)
- Flow control semaphore format `group_vtask_size`, head slot serves as the node group anchor
- every core task belongs to both the **group** and a **host**: `_vtask_group` is the group number; `_vtask_host` is the host where that core task actually runs
- `from-vtask-head.sh` computes `to_host` from `_vtask_group` (legacy `_vtask_size_sema`) under the `VTASK_MODE=GROUP-BOUND` branch (business logic); `vtask add-subtask` must set `to_host` explicitly — the platform never guesses a host inside the group
- head→core uses the script's explicit `to_host`; the head is SLOT-BOUND, so `pod_id` derivation does not trigger
- core→core: explicit `to_host`/`to_ip` > `_vtask_host` (the actual host of the previous core) > `pod_id` (fallback); core modules no longer need to share the same `pod_id`
- controld updates `_vtask_host` to the explicit `to_host` when the core task is inserted; a GROUP-BOUND core task with no resolvable `to_host` is rejected at insertion
- the group number comes from `_vtask_group` first (legacy wait-queue falls back to `_vtask_size_sema`); the quota semaphore uses only the group number, and `_vtask_host` never participates in sema reconstruction

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
- `_`-prefixed headers (`_vtask_id`, `_vtask_head_module`, `_vtask_host`/`_vtask_group`; `_vtask_size_sema` only for old tasks/custom) auto-propagate along the core chain
- Each core module independently configures slot counts, `task_dist_mode`, and images
- HOST-BOUND vtask core chains prefer `_vtask_host` to stay on the same node; GROUP-BOUND core→core also inherits `_vtask_host` (the previous hop's actual host); `pod_id` remains the fallback

### 9.3.5 Placement Routing: to_host Sources and Priority

For HOST-BOUND tasks, the target host is resolved in this order (`cmd/controld/task/add_local.go`):

```
explicit to_host  >  explicit to_ip  >  _vtask_host (core)  >  pod_id derivation
```

- **Explicit `to_host` / `to_ip` wins**: caller-supplied routing is explicit intent; `_vtask_host`/`pod_id` never overrides it.
- **`_vtask_host` wins over `pod_id`**: applies to `HOST-BOUND` + `vtask_role=core` tasks. A HOST-BOUND vtask uses its immutable bind anchor; a GROUP-BOUND vtask uses the actual host of the previous core. The GROUP-BOUND head→core hop has no `_vtask_host` and must set `to_host` explicitly (see §9.3.6).
- **`pod_id` derivation conditions** (all must hold):
  1. target module has `task_dist_mode=HOST-BOUND`;
  2. source module has `task_dist_mode=HOST-BOUND` (SLOT-BOUND / DEFAULT sources do not trigger it);
  3. both `pod_id` values are non-empty and equal;
  4. the caller did not set `to_host` or `to_ip`.
- When relaying through the main-router, the agent writes `last_from_module` / `last_from_ip`; controld unwraps the router and evaluates the real predecessor module.
- The derived `to_host` gets a `.<cluster>` suffix (the `t_slot.host` storage format).
- If the source IP cannot be resolved from `t_host`, no `to_host` is written and controld logs a warning; the task then stays READY until the host registration is fixed or routing is made explicit.

Configuration by vtask mode:

| Mode | head→core | core→core | Configuration |
|------|-----------|-----------|---------------|
| HOST-BOUND | `_vtask_host` (wait-queue; legacy `pod_id`) | `_vtask_host`/`pod_id` | head/core share the anchor; `pod_id` is core→core/non-vtask fallback |
| GROUP-BOUND | explicit `to_host` from the script (computed from `_vtask_group`) | explicit > `_vtask_host` (previous hop's actual host) > `pod_id` | head is SLOT-BOUND; the first hop must be explicit |
| DEFAULT | no binding | no binding | `pod_id` not needed |

### 9.3.6 vTask Anchor Attributes (2026-11-10, updated 2026-11-12)

A HOST-BOUND vtask brings its core stages back to the head anchor; a GROUP-BOUND vtask
distributes core stages across hosts in the group and keeps each task local. Since
2026-11-10 the platform expresses the anchor with three reserved attributes; since
2026-11-12 the GROUP-BOUND `_vtask_host` is defined as the host where that core task
actually runs:

| Attribute | Scope | Meaning |
|---|---|---|
| `_vtask_head_module` | all vtasks | vtask-head module name (used to rebuild the semaphore name) |
| `_vtask_host` | HOST-BOUND / GROUP-BOUND | HOST-BOUND: bind anchor, full name (e.g. `n0-0.inline`), immutable for the vtask lifetime. GROUP-BOUND: the host where that core task actually runs; updated from an explicit `to_host` and propagated along the core chain |
| `_vtask_group` | GROUP-BOUND | group number (head slot seq); **not a host** |
| `_vtask_size_sema` | internal/compat | old tasks or custom `--sema-name`; new tasks no longer write it; rebuilt from the three attributes above |

**Writer**: wait-queue `run.sh` writes `_vtask_head_module` + `_vtask_host` (HOST-BOUND) or
`_vtask_group` (GROUP-BOUND) on the head task after `vtask bind`; it no longer writes
`_vtask_size_sema`. GROUP-BOUND does not write `_vtask_host`; it is written by controld when
the first core task is inserted. The `_vtask_*` headers propagate through the agent `_` rule;
the router is `SLOT-BOUND`, does not consume the anchor, only passes it through, and is never
moved to a compute node.

**Consumer**: at task insertion, controld applies:

- target `task_dist_mode=HOST-BOUND`, `vtask_role=core`, no `_vtask_group`, and no explicit
  `to_host`/`to_ip`: `to_host = _vtask_host`;
- GROUP-BOUND core: an explicit `to_host`/`to_ip` wins and updates `_vtask_host` to that host;
  without explicit routing, `to_host = _vtask_host` (the previous core's actual host); finally
  `pod_id` is used as a fallback; if `to_host` still cannot be resolved, insertion is rejected.
  `_vtask_group` is only for sema rebuild/observability;
- explicit `to_host`/`to_ip` still wins; `pod_id` remains the fallback for HOST-BOUND
  core→core / non-vtask paths.

**Priority**:

```
explicit to_host > explicit to_ip > _vtask_host (core) > pod_id
```

Both HOST-BOUND and GROUP-BOUND core→core chains include the `_vtask_host` step. On the
GROUP-BOUND head→core hop `_vtask_host` is empty, so `to_host` must be explicit; the group
number is never treated as a host.

**Sema reconstruction** (`VTaskSizeSema`):

- HOST-BOUND: `host_vtask_size:<_vtask_head_module>:<short(_vtask_host)>`
- GROUP-BOUND: `group_vtask_size:<_vtask_head_module>:<_vtask_group>`; `_vtask_host` does not participate
- `updateVTask` pre-deduct check, `doVTaskFinished`, `FailVtask`, `vtask unbind`, and
  `vtask get` all prefer `_vtask_size_sema` (old tasks/custom) and otherwise rebuild it;
- missing fields are logged loudly; the tail completion flag is not blocked.

**Constraints**:

- `_vtask_host` is the full host name (same format as `t_slot.host`); the sema uses the short form.
- GROUP-BOUND's `_vtask_host` is the host where that core task actually runs, not the
  head/entry node; the first hop must be explicit.
- A GROUP-BOUND core explicitly overriding `to_host` is supported: controld updates
  `_vtask_host` and downstream tasks follow by default.
- HOST-BOUND data-plane core tasks should not override `to_host`; if they do, downstream must
  follow the actual source explicitly or local data (`/dev/shm`) will not be readable.
- A `SLOT-BOUND` router never becomes a compute-node task because of `_vtask_host`.
- On the GROUP-BOUND head→core hop, `vtask add-subtask` must set `to_host`; with `--direct`
  it must also forward `_vtask_id`/`_vtask_group`/`_vtask_head_module` explicitly (the CLI
  does not read the current task headers).

## 9.4 The Semaphore Mechanism

### 9.4.1 Naming Conventions

| Distribution mode | Format | Example |
|---------|------|------|
| DEFAULT | `vtask_size:<module>` | `vtask_size:vtask-head` |
| HOST-BOUND | `host_vtask_size:<module>:<host>` | `host_vtask_size:vtask-head:n0-0` |
| SLOT-BOUND (GROUP-BOUND vtask) | `group_vtask_size:<module>:<group>` | `group_vtask_size:vtask-head:2` |

### 9.4.2 Unified Semaphore

Each flow control semaphore has a single instance (no colon prefix), pre-deducted by `vtask bind` and restored automatically by controld; scripts no longer participate in resource return:

| Manager | Operation | Timing |
|--------|------|------|
| wait-queue `check.sh` (ACTION_CHECK) | reads semagroup max; when ≤ 0 returns non-zero → agent stops picking tasks and polls/retries | every agent poll, before `GetTaskList` |
| `vtask bind` (`BindVtaskResource`) | atomic decrement (semagroup picks the member with max value, returns group number/hostname) | wait-queue serialized admission |
| controld `updateVTask` | tasks with `_vtask_size_sema` or `_vtask_head_module`+host/group are pre-deducted and skipped at pickup; only bare tasks are deducted as fallback | head slot picks up tasks |
| controld `doVTaskFinished` / `FailVtask` | increment (by `_vtask_size_sema` or rebuilt from `_vtask_head_module`+host/group) | vtask-tail completes / vtask fails |

> The historical dual-semaphore design (flow control version + `:`-prefixed programmable version) is deprecated; the corresponding parameter no longer takes effect.

### 9.4.3 Flow Control Process

```
wait-queue slot polls each round
  → check.sh: semagroup max ≤ 0 → exit 1, agent stops picking (task never picked up)
  → max > 0 → pick up task → run.sh (inside the standard module image)
  → vtask bind: semagroup.Decrement atomic pre-deduction (returns hostname/group number)
  → vtask-head task picked up: updateVTask skips pre-deducted tasks (bare DEFAULT tasks deducted as fallback)
  → vtask pipeline executes
  → vtask-tail completes → doVTaskFinished: increment (quota restored)
  → check.sh next polling round lets the next task through
```

### 9.4.4 Wait Queue Gating

The wait-queue module's `vtask_size` stays at 1 (serialization). The gate is incremented by wait-queue `run.sh` immediately after bind + delivery (2026-11-02 L1); `from-vtask-head.sh` no longer adds +1.

## 9.5 Pipeline Flow

### 9.5.1 Header Propagation

```
main-router task add
  → agent addSinkTasks: _vtask_* gating (only sink modules with vtask_role≠"" pass)
  → wait-queue task: does not receive _vtask_id / _vtask_size_sema; run.sh writes _vtask_head_module + host/group

wait-queue task executes the standard module run.sh
  → --direct gRPC creates the vtask-head task
  → vtask_role=head → self-ref vtask=id
  → _vtask_head_module + _vtask_host (HOST-BOUND) or _vtask_group (GROUP-BOUND) written by run.sh;
    no _vtask_size_sema; GROUP-BOUND's _vtask_host is written when the first core task is inserted
  → task terminates at the wait-queue slot (no sink-tasks.txt written, no flow-back)

vtask-head task executes from-vtask-head.sh
  → sink-tasks.txt async creates the vtask-core task
  → agent addSinkTasks: _-prefixed headers auto-propagate

vtask-core → vtask-tail (same path as above)

vtask-tail completes
  → doVTaskFinished: reads _vtask_size_sema or rebuilds from the new attributes → increments
  → from-vtask-tail.sh: same (legacy sema / rebuilt attributes)
```

### 9.5.2 Lifecycle

```
task creation → ON CONFLICT dedup
  ↓
task scheduling → FOR UPDATE SKIP LOCKED (single slot exclusive)
  ↓
task execution → agent → status_code: RUNNING → OK/ERROR
  ↓  _vtask_id / _vtask_head_module / _vtask_host / _vtask_group auto-propagate via the _ prefix (_vtask_size_sema compat)
  ↓
tail completes → doVTaskFinished → sema rebuilt/compat → increment
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
| `_`-prefixed headers | agent auto-propagates (`_vtask_id`, `_vtask_head_module`, `_vtask_host`, `_vtask_group`) | caller explicitly passes `--header`; on the GROUP-BOUND first hop forward at least `_vtask_id`/`_vtask_group`/`_vtask_head_module` |
| `from_module`/`from_ip` | agent injects (more accurate) | CLI injects from env |

**Selection principles**:

| Script | Path | Reason |
|------|------|------|
| `from-vtask-head.sh` | async (default) | no resource-binding prerequisite, enjoys header auto-propagation; GROUP-BOUND must set `to_host` explicitly |
| wait-queue standard module run.sh | sync (`--direct`) | bind has already deducted resources; must sense success/failure synchronously, rollback on failure |
| Ops manual CLI | sync (gRPC) | no sink-tasks.txt mechanism outside the agent |
| GROUP-BOUND first hop needing early failure | optional `--direct` | learns the missing-`to_host` rejection synchronously; must forward `_vtask_*` explicitly and unbind on failure |

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
# Bind (auto-finds the head module, returns hostname or group number)
to_host=$(scalebox vtask bind) || exit 1

# Unbind (semaphore name passed for error rollback; normal completion is restored by controld)
scalebox vtask unbind --sema-name="host_vtask_size:vtask-head:n0-0"

# Unbind (ops passes vtask-id; server reads _vtask_size_sema or rebuilds it from the new attributes)
scalebox vtask unbind --vtask-id=42
```

| vtask type | bind returns |
|-----------|----------|
| HOST-BOUND | hostname (e.g. `n0-0`) |
| GROUP-BOUND | group number (head slot's seq, e.g. `0`) |
| DEFAULT | empty (no-op) |

The GROUP-BOUND return value is a group number, not a host. The first head→core hop is placed
by the application's explicit `to_host`; core→core may inherit `_vtask_host`. Quota remains
per group.

### 9.6.6 vtask add-subtask

```bash
# HOST-BOUND: async path, _-prefixed headers auto-propagate; no routing header,
# controld uses _vtask_host (legacy/non-vtask paths may still use the same pod_id)
scalebox vtask add-subtask --module=vtask-core $body

# GROUP-BOUND first hop (head→core): to_host is mandatory (head is SLOT-BOUND;
# the platform cannot derive a host from the group number). Async path:
# the agent propagates _vtask_* automatically.
scalebox vtask add-subtask --module=vtask-core \
    --header to_host=$to_host $body

# GROUP-BOUND first hop needing synchronous error awareness: add --direct.
# The CLI does not read the current task headers, so forward _vtask_* explicitly.
scalebox vtask add-subtask --direct --module=vtask-core \
    --header to_host=$to_host \
    --header _vtask_id=$(scalebox::task_header "$2" _vtask_id) \
    --header _vtask_group=$(scalebox::task_header "$2" _vtask_group) \
    --header _vtask_head_module=$(scalebox::task_header "$2" _vtask_head_module) $body || {
    scalebox vtask unbind --sema-name="group_vtask_size:${head_module}:${group}"
    exit 1
}

# GROUP-BOUND core→core: sink-generated tasks inherit the upstream _vtask_host;
# no repeated to_host/pod_id is needed. Set to_host explicitly only to deviate.

# Ops manual, outside the agent: --app-id is mandatory; a GROUP-BOUND core target
# also needs explicit _vtask_group/to_host
scalebox vtask add-subtask --app-id=10 --module=vtask-core \
    --header _vtask_id=42 --header _vtask_group=0 --header to_host=n0-0 $body
```

Auto-injection: `from_module` (`PLAT_MODULE_NAME` env), `from_ip` (local IPv4). Not overwritten when the caller explicitly passes them.

Both paths finally go through controld's `DoAddLocalTaskList`: a GROUP-BOUND core task with no
resolvable `to_host` (explicit, `_vtask_host`, or `pod_id`) is rejected at insertion and never
lands at a random host. wait-queue's delivery to `vtask-head` is SLOT-BOUND (`to_slot_index`)
and needs no `to_host`; the gate is released by run.sh right after bind + delivery.

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
| `vtask_size >= slot count` (wait-queue) | concurrent bind admission (use total free capacity; too low causes oversell, see INV-23) |

## 9.8 Ops Scenario Quick Reference

| Scenario | Command |
|------|------|
| View global vtask progress | `vtask list --app-id=N` |
| View a single vtask's details | `vtask get --vtask-id=N` |
| View vtask sub-task distribution | `vtask list-subtasks --vtask-id=N` |
| Terminate a single vtask | `vtask fail --vtask-id=N` |
| Batch termination | `vtask list \| awk ... \| xargs vtask fail` |
| Pause global flow control | `semaphore decrement vtask_size:<mod> N` |
| Manually release leaked resources | `vtask unbind --sema-name="host_vtask_size:..."` |
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
