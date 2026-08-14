# 5. Operations

## 5.1 System Monitoring

### 5.1.1 Service Status Monitoring

```bash
# Check service status
docker compose ps

# View service logs
docker compose logs controld
docker compose logs actuator
docker compose logs database

# Check service health
scalebox cluster status
```

### 5.1.2 Resource Usage Monitoring

```bash
# View node resources
scalebox host metrics

# View storage usage
docker system df

# View container resource usage
docker stats
```

## 5.2 Troubleshooting

### 5.2.1 Database Problems

**Symptom**: database connection failed

**Solution**:
1. Wait for the database to fully start: `sleep 10`
2. Manually check the database: `docker exec database psql -U scalebox -c "\l"`
3. Recreate the database volume: `docker compose down -v && docker compose up -d`

### 5.2.2 Task Execution Problems

**Symptom**: tasks remain in READY state

**Solution**:
1. Check the actuator logs: `docker logs actuator`
2. Check slot status: `scalebox slot list`
3. Restart the actuator: `docker compose restart actuator`

### 5.2.3 Network Problems

**Symptom**: inter-node communication failed

**Solution**:
1. Check network connectivity: `ping <node-ip>`
2. Check open ports: `netstat -tlnp | grep 50051`
3. Check the firewall: `ufw status`

## 5.3 Performance Optimization

### 5.3.1 System Parameter Tuning

```bash
# Adjust Docker parameters
dockerd --max-concurrent-downloads 10 --max-concurrent-uploads 10

# Adjust system parameters
sysctl -w vm.max_map_count=262144
sysctl -w fs.file-max=100000
```

### 5.3.2 App Performance Optimization

1. **Set parallelism appropriately**
   - Set module parallelism based on data volume
   - Avoid resource contention caused by excessive parallelism

2. **Optimize data locality**
   - Prefer local storage
   - Reduce cross-node data transfer

3. **Batch processing optimization**
   - Set batch sizes appropriately
   - Balance latency and throughput

## 5.4 Backup and Recovery

### 5.4.1 Database Backup

```bash
# Back up the database
docker exec database pg_dump -U scalebox scalebox > backup.sql

# Restore the database
docker exec -i database psql -U scalebox scalebox < backup.sql
```

### 5.4.2 Configuration Backup

```bash
# Back up configuration files
tar -czf scalebox-config.tar.gz examples/ dockerfiles/ runtime/

# Restore configuration files
tar -xzf scalebox-config.tar.gz
```

## 5.5 Security Maintenance

### 5.5.1 Access Control

```bash
# Update keys
cd runtime && make update-pubkey

# Reset the user password
docker exec database psql -U scalebox -c "ALTER USER scalebox WITH PASSWORD 'newpassword';"
```

### 5.5.2 Security Updates

```bash
# Update container images
cd runtime && make pull-all

# Restart services
docker compose down && docker compose up -d
```

## 5.6 Log Management

### 5.6.1 Log Collection

```bash
# Collect logs from all services
docker compose logs > scalebox.log

# Filter logs by time
docker compose logs --since "2024-01-01" --until "2024-01-02"
```

### 5.6.2 Log Analysis

```bash
# Find error logs
grep -i error scalebox.log

# Calculate task execution times
grep "task completed" scalebox.log | awk '{print $NF}' | sort -n
```

## 5.7 System Upgrade

### 5.7.1 Version Upgrade

```bash
# Pull the new version
git pull origin main

# Update container images
make pull-all

# Smoothly restart services
docker compose restart
```

### 5.7.2 Configuration Migration

1. Back up the current configuration
2. Update configuration files
3. Verify configuration validity
4. Restart services

## 5.8 Common Problems Quick Reference

| Symptom | Possible cause | Solution |
| -------- | -------- | -------- |
| Database connection failed | Database not started | Wait for the database to finish starting |
| Tasks cannot execute | Insufficient Slots | Add compute nodes or slots |
| Network connection failed | Firewall blocking | Open the corresponding ports |
| Insufficient disk space | Log or data accumulation | Clean up unnecessary files |
| Insufficient memory | Parallelism too high | Reduce parallelism or add memory |

## 5.9 Cluster Management

### 5.9.1 Cluster Classification

| Cluster name | Primary purpose | Description |
| ------ | -------- | ---- |
| Single-node Cluster | Testing, development | A single node serves as both the management node and the compute node |
| Static Cluster | Production | Management nodes include one or more head nodes; compute nodes are fixed |
| Dynamic Cluster | Elastic computing | Management nodes include one or more head nodes; compute nodes are dynamically allocated by an external scheduler |
| Inline Cluster | Hybrid mode | Includes one or more head nodes and several compute nodes; nodes are usually acquired dynamically |

### 5.9.2 Multi-Cluster Collaborative Computing

- Compute resource management: supports both static definition and dynamic acquisition
- Deploy the cluster-admin app to manage clusters, supporting dynamic expansion of cluster nodes
- System management and dynamic monitoring of compute nodes through node-agent

### 5.9.3 Cluster Architecture

```
                          +-----------------------+
                          |   Control Plane       |
                          |-----------------------|
                          | - Global database     |
                          | - Scheduler controld  |
                          | - Slot launcher actuator |
                          +----------+------------+
                                     |
                   +-------------------------------------+
                   |                 |                   |
            +------------+    +------------+      +------------+
            | node-agent |    | node-agent |      | node-agent |
            +------+-----+    +------+-----+      +------+-----+
                   |                 |                   |
           +-------------+    +-------------+    +--------------+
           |   slot 1    |    |   slot 1    |    |   slot 1     |
           |   slot 2    |    |   slot 2    |    |   slot 2     |
           |   slot N    |    |   slot N    |    |   slot N     |
           +-------------+    +-------------+    +--------------+
                    |
              Module Pipeline
        (moduleA → moduleB → moduleC)
```
