# 1. Frequently Asked Questions (FAQ)

## 1. Deployment and Connection

### 1.1 Non-Standard Port Settings

Network ports involved in Scalebox:
- database port, standard 5432
- controld gRPC port, standard 50051
- WebUI HTTP port, standard 8088

**Server side (Docker Compose)**:

```yaml
# build/compose.yaml
services:
  controld:
    ports:
      - <custom-grpc-port>:50051
  database:
    ports:
      - <custom-pg-port>:5432
  webui:
    ports:
      - <custom-web-port>:8088
```

**Client side**: set `${HOME}/.scalebox/environments`:

```bash
GRPC_SERVER=<internal-ip>:<custom-grpc-port>
PGHOST=<internal-ip>
PGPORT=<custom-pg-port>
```

### 1.2 Database Connection Failed

**Symptom**: `dial tcp: connect: connection refused`

**Solution**:
```bash
# Wait for the database to fully start
sleep 10

# Check manually
docker exec database psql -U scalebox -c "\l"

# Recreate the database volume
docker compose -f build/compose.yaml down -v
docker compose -f build/compose.yaml up -d
```

### 1.3 gRPC Connection Failed

**Symptom**: CLI commands report `Unavailable` or `connection refused`

**Solution**:
```bash
# Check the GRPC_SERVER setting
echo $GRPC_SERVER

# Check the controld service status
docker compose -f build/compose.yaml ps controld

# Verify gRPC connectivity
grpcurl -plaintext localhost:50051 list
```

## 2. Tasks and Scheduling

### 2.1 Tasks Stuck in READY State

**Cause**: insufficient Slots; no available execution slots.

**Solution**:
```bash
# Check slot status
scalebox slot list
scalebox slot list --host <hostname>

# Add slots
scalebox slot add --module <module-name> --host <hostname> --count 4

# Check the actuator logs
docker logs actuator
```

### 2.2 Task Execution Failed (exit code != 0)

**Troubleshooting steps**:
1. View task logs: `scalebox task log <task-id>`
2. Check whether the module image exists: `docker images | grep <image-name>`
3. Check node resources: `scalebox host show <hostname>`
4. Check the meaning of the exit code (see :doc:`Appendix 6 - Exit Code Specification <appendix/6_exit_code_spec>`)

### 2.3 Duplicate Creation of Batch Tasks

Always use `--conflict-action IGNORE` to avoid duplicates. Or set the header `repeatable: yes` to allow repeated distribution.

## 3. Semaphores and Flow Control

### 3.1 Semaphore Not Found Error

```bash
# Method 1: enable automatic creation
SEMAPHORE_AUTO_CREATE=yes scalebox semaphore get my_sema

# Method 2: create manually
scalebox semaphore create --app-id <app-id> my_sema 0
```

### 3.2 Abnormal Semaphore Values (not converging)

Check whether the flow control version and the programmable version match:
- Flow control version (no prefix) is maintained automatically by controld
- Programmable version (`:` prefix) is operated by scripts
- See :doc:`VTask Semaphore Mechanism <developer/9_vtask>` §9.4

## 4. WebUI

### 4.1 Pages Inaccessible

```bash
# Check the WebUI service
docker compose -f build/compose.yaml ps webui

# Check the port
curl http://localhost:8088/api/v1/stats/system
```

### 4.2 Create Operations Return Errors

Common cause: proto field name mismatch. The REST API accepts both camelCase and snake_case.

### 4.3 DAG Page Shows Blank

The `GetAppDAG` RPC on the controld side may not be implemented; the page degrades to module topology.

## 5. VS Code Extension

### 5.1 TreeView Shows an Empty List

Check the configuration:
- Does `scalebox.serverUrl` point to the WebUI correctly
- Is the WebUI accessible: `curl <serverUrl>/api/v1/apps`
- Check the VS Code Output panel (Scalebox channel)

### 5.2 app.yaml Validation Not Working

Confirm:
- The file name matches the `app*.yaml` or `scalebox*.yaml` pattern
- The file's language mode is YAML
- VS Code has not disabled the Scalebox extension

## 6. Logs and Debugging

### 6.1 Adjusting Log Verbosity

```bash
# Verbose mode (complete error chain + location information)
LOG_LEVEL=debug scalebox task list

# Quiet mode
LOG_LEVEL=error scalebox task list
```

### 6.2 Collecting Service Logs

```bash
# All services
docker compose -f build/compose.yaml logs > scalebox.log

# By service
docker compose -f build/compose.yaml logs controld
docker compose -f build/compose.yaml logs actuator

# By time
docker compose -f build/compose.yaml logs --since "2026-01-01"
```

### 6.3 Interpreting gRPC Error Messages

Error message format: `rpc error: code = <Code> desc = <Message>`

Common gRPC status codes and corresponding scenarios:

| Status code | Meaning | Common causes |
|--------|------|---------|
| `Unavailable` | service unreachable | controld not started or network unreachable |
| `NotFound` | resource does not exist | invalid App/Module/Task ID |
| `AlreadyExists` | resource already exists | duplicate name creation validation failed |
| `InvalidArgument` | parameter error | invalid flag value |
| `Internal` | internal error | DB anomaly or business logic error |

> Continuously updated; contributions of questions and answers are welcome.
