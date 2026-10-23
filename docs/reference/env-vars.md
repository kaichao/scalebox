<!-- AUTO-GENERATED from Go 源码中的 `os.Getenv` / `os.LookupEnv` / `getEnvString` 读取点 by `make docs-gen` (go-scalebox); DO NOT EDIT -->
<!-- source: go-scalebox@3a8dd2e (2026-10-21) -->


# 环境变量参考（平台读取点）

> 本页是平台各进程读取的环境变量事实表。模块容器内可见的变量与注入顺序
> （`ACTION_*`、`WORK_DIR`、`PLAT_*`、`SINK_MODULE` 等）见
> [`module.md`](./module.md) §4。

## agent

| 变量 | 默认值 | 读取位置 |
|---|---|---|
| `ACTION_CHECK` |  | `cmd/agent/actions.go:183` |
| `ACTION_RUN` |  | `cmd/agent/actions.go:116` |
| `ACTION_SETUP` |  | `cmd/agent/actions.go:139` |
| `ACTION_TEARDOWN` |  | `cmd/agent/actions.go:161` |
| `APP_ID` |  | `cmd/agent/init.go:134` |
| `AUTH_TOKEN` |  | `cmd/agent/init.go:179` |
| `CONFLICT_ACTION` |  | `cmd/agent/sink_tasks.go:116` |
| `DEBUG_SIGBUS` |  | `cmd/agent/main.go:131` |
| `GRPC_SERVER` |  | `cmd/agent/init.go:161` |
| `IS_SINGULARITY` |  | `cmd/agent/init.go:257` |
| `LOCAL_IP` |  | `cmd/agent/iobytes.go:186`, `cmd/agent/config/config.go:45` |
| `LOCAL_SHMDIR` |  | `cmd/agent/init.go:318` |
| `LOCAL_TMPDIR` |  | `cmd/agent/init.go:292` |
| `LOG_FORMAT` | `json` | `cmd/agent/config/config.go:57` |
| `LOG_LEVEL` | `info` | `cmd/agent/config/config.go:56` |
| `MAX_SINK_TASKS` |  | `cmd/agent/sink_tasks.go:125` |
| `MODULE_ID` |  | `cmd/agent/init.go:126` |
| `NUM_IDLE_POLLS` |  | `cmd/agent/run.go:49`, `cmd/agent/run.go:50`, `cmd/agent/run.go:149` |
| `PLAT_ALWAYS_RUNNING` |  | `cmd/agent/init.go:188` |
| `PLAT_DIR_QUOTA_GB` |  | `cmd/agent/check.go:99`, `cmd/agent/init.go:275` |
| `PLAT_FREE_SPACE_GB` |  | `cmd/agent/check.go:75`, `cmd/agent/init.go:276` |
| `PLAT_HEARTBEAT_SECONDS` |  | `cmd/agent/init.go:267` |
| `PLAT_IS_ROUTER` |  | `cmd/agent/init.go:278` |
| `PLAT_MAX_IDLE_POLLS` |  | `cmd/agent/init.go:195` |
| `PLAT_MODULE_NAME` |  | `cmd/agent/init.go:281`, `cmd/agent/init.go:374` |
| `PLAT_OUTPUT_TEXT_SIZE` |  | `cmd/agent/init.go:214` |
| `PLAT_POLL_INTERVAL_SECONDS` |  | `cmd/agent/init.go:208` |
| `PLAT_SRV_THROTTLE` |  | `cmd/agent/run.go:262` |
| `PLAT_TASK_MAX_SECONDS` |  | `cmd/agent/init.go:191` |
| `PLAT_TASK_MIN_SECONDS` |  | `cmd/agent/init.go:202` |
| `PLAT_TEXT_TRUNC_MODE` |  | `cmd/agent/init.go:221` |
| `SINGULARITY_IMAGE` |  | `cmd/agent/config/config.go:62` |
| `SINK_MODULE` |  | `cmd/agent/sink_tasks.go:66`, `cmd/agent/init.go:142` |
| `SLOT_ID` |  | `cmd/agent/run.go:114`, `cmd/agent/config/config.go:44` |
| `TMPDIR_GROUP` |  | `cmd/agent/init.go:343` |
| `TMPFS_WORKDIR` |  | `cmd/agent/init.go:369` |
| `WORK_DIR` | `/work` | `cmd/agent/config/config.go:60` |

## controld

| 变量 | 默认值 | 读取位置 |
|---|---|---|
| `CLEANUP_INTERVAL_MINUTES` |  | `cmd/controld/task/cache_attrset.go:67` |
| `CROSS_SERVER_MAX_RETRY_ATTEMPTS` |  | `cmd/controld/task/add_remote.go:18` |
| `DATABASE_URL` |  | `cmd/controld/config/config.go:46` |
| `DEFAULT_EXPIRATION_MINUTES` |  | `cmd/controld/task/cache_attrset.go:62` |
| `DEFAULT_SLOT_STATUS` | `READY` | `cmd/controld/config/config.go:57` |
| `EVENT_MODE` |  | `cmd/controld/grpc/event.go:19` |
| `EVENT_POOL_SIZE` |  | `cmd/controld/grpc/event.go:103` |
| `GRPC_TLS_ENABLED` |  | `cmd/controld/main.go:84` |
| `INITIAL_SLOT_STATUS` | `READY` | `cmd/controld/config/config.go:58` |
| `LOG_LEVEL` | `info` | `cmd/controld/grpc/app_env_vol.go:35`, `cmd/controld/config/config.go:41`, `cmd/controld/task/init.go:14` |
| `MAX_GRPC_STREAMS` |  | `cmd/controld/main.go:75` |
| `MY_APP_CLUSTER` |  | `cmd/controld/config/config.go:51` |
| `MY_CODE_PATH` |  | `cmd/controld/config/config.go:53` |
| `MY_IMAGE_NAME` |  | `cmd/controld/config/config.go:52` |
| `MY_SLOT_REGEX` |  | `cmd/controld/config/config.go:54` |
| `PG_BATCH_SIZE` |  | `cmd/controld/grpc/event.go:121` |
| `PORT` |  | `cmd/controld/main.go:145` |
| `SEMAPHORE_AUTO_CREATE` |  | `cmd/controld/grpc/coordination_service.go:388` |
| `TASK_POOL_SIZE` |  | `cmd/controld/task/task.go:75` |

## actuator

| 变量 | 默认值 | 读取位置 |
|---|---|---|
| `AUTH_TOKEN` |  | `cmd/actuator/cluster.go:139`, `cmd/actuator/init.go:61` |
| `GRPC_SERVER` |  | `cmd/actuator/init.go:56` |
| `GRPC_TLS_CA_FILE` |  | `cmd/actuator/cluster.go:144` |
| `GRPC_TLS_ENABLED` |  | `cmd/actuator/cluster.go:142` |
| `LOG_LEVEL` | `info` | `cmd/actuator/init.go:30`, `cmd/actuator/init.go:39` |
| `SLOT_POLL_INTERVAL_SECONDS` |  | `cmd/actuator/init.go:97` |
| `TIMEOUT_RUN_IMAGE_SECONDS` |  | `cmd/actuator/init.go:102` |

## webui

| 变量 | 默认值 | 读取位置 |
|---|---|---|
| `AUTH_TOKEN` |  | `cmd/webui/client/client.go:126` |
| `WEBUI_PORT` |  | `cmd/webui/main.go:44` |

## cli

| 变量 | 默认值 | 读取位置 |
|---|---|---|
| `APP_ID` |  | `cmd/cli/app/add_slots.go:32`, `cmd/cli/app/app.go:62` |
| `CLUSTER` |  | `cmd/cli/task/task.go:134` |
| `DATABASE_URL` |  | `cmd/cli/init.go:75` |
| `GRPC_SERVER` |  | `cmd/cli/init.go:66`, `cmd/cli/root.go:90`, `cmd/cli/task/task.go:244` … (+2) |
| `LOG_LEVEL` | `info` | `cmd/cli/init.go:26`, `cmd/cli/run/command.go:423` |
| `MODULE_NAME` |  | `cmd/cli/vtask/vtask.go:444`, `cmd/cli/task/task.go:86` |
| `MY_APP_CLUSTER` |  | `cmd/cli/run/command.go:415` |
| `MY_APP_NAME` |  | `cmd/cli/run/command.go:195` |
| `MY_MR_CODE_PATH` |  | `cmd/cli/run/command.go:181` |
| `MY_MR_IMAGE_NAME` |  | `cmd/cli/run/command.go:116`, `cmd/cli/run/command.go:186` |
| `PGHOST` |  | `cmd/cli/init.go:75`, `cmd/cli/init.go:79` |
| `PGURL` |  | `cmd/cli/init.go:75` |
| `PLAT_MODULE_NAME` |  | `cmd/cli/vtask/vtask.go:442`, `cmd/cli/vtask/vtask.go:462`, `cmd/cli/task/task.go:84` … (+1) |
| `SEMAPHORE_AUTO_CREATE` |  | `cmd/cli/semaphore/semaphore.go:378` |
| `SLOT_ID` |  | `cmd/cli/event.go:34` |
| `TASK_ID` |  | `cmd/cli/event.go:29` |
| `USER` |  | `cmd/cli/cluster/init_local.go:73` |
| `WORK_DIR` | `/work` | `cmd/cli/vtask/vtask.go:463`, `cmd/cli/task/task.go:175` |

## internal

| 变量 | 默认值 | 读取位置 |
|---|---|---|
| `AUTH_TOKEN` |  | `internal/grpc/grpc.go:91` |
| `CONFLICT_ACTION` |  | `internal/semaphore/create.go:24`, `internal/semaphore/create.go:194` |
| `DATABASE_URL` |  | `internal/postgres/init.go:23` |
| `GRPC_MAX_RETRY_ATTEMPTS` |  | `internal/grpc/init.go:13` |
| `GRPC_SERVER` |  | `internal/grpc/grpc.go:270`, `internal/postgres/init.go:88` |
| `GRPC_TLS_CA_FILE` |  | `internal/grpc/grpc.go:90` |
| `GRPC_TLS_ENABLED` |  | `internal/grpc/grpc.go:89` |
| `LOCAL_ADDR` |  | `internal/postgres/init.go:91` |
| `LOCAL_IP` |  | `internal/common/net.go:19`, `internal/postgres/init.go:127` |
| `LOCAL_IP_INDEX` |  | `internal/common/net.go:23`, `internal/postgres/init.go:131` |
| `LOGRUS_ASYNC_SECONDS` |  | `internal/common/log.go:70` |
| `LOG_LEVEL` | `info` | `internal/common/log.go:21`, `internal/common/log.go:76` |
| `PGDB` |  | `internal/postgres/init.go:113` |
| `PGHOST` |  | `internal/postgres/init.go:75` |
| `PGPASS` |  | `internal/postgres/init.go:62` |
| `PGPORT` |  | `internal/postgres/init.go:76` |
| `PGURL` |  | `internal/postgres/init.go:29` |
| `PGUSER` |  | `internal/postgres/init.go:106` |
| `PG_MAX_CONNS` |  | `internal/postgres/pgx.go:44` |
| `PG_MAX_CONN_IDLE_TIME_MIN` |  | `internal/postgres/pgx.go:61` |
| `PG_MAX_CONN_LIFETIME_MIN` |  | `internal/postgres/pgx.go:56` |
| `PG_MAX_IDLE_CONNS` |  | `internal/postgres/sql.go:49` |
| `PG_MAX_OPEN_CONNS` |  | `internal/postgres/sql.go:54` |
| `PG_MIN_CONNS` |  | `internal/postgres/pgx.go:49` |
| `PLAT_MODULE_NAME` |  | `internal/task/add.go:159` |
| `SCALEBOX_CERTS_DIR` |  | `internal/postgres/init.go:40` |
| `TIMEOUT_SECONDS` |  | `internal/task/add.go:101` |
