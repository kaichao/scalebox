# 3. Running Apps

## 3.1 Creating and Running Apps

### 3.1.1 App Definition File

Scalebox apps are defined through YAML configuration files. They mainly include:

- **App name**: unique identifier, supports variable substitution
- **Cluster configuration**: target cluster name
- **Module definitions**: the modules included and their configurations
- **Parameter settings**: runtime parameters

### 3.1.2 App Parameter File

App parameters are defined through environment variable files, supporting variable substitution.

### 3.1.3 Creating an App

```bash
# Use the default configuration file and parameter file
scalebox run

# Specify configuration files
scalebox run --env-file scalebox.env --app-file app.yaml
```

## 3.2 App Management

### 3.2.1 Viewing App Status

```bash
# List all apps
scalebox app list

# View app details
scalebox app show --app-id=<app-id>
```

### 3.2.2 App Operations

```bash
# Stop an app
scalebox app set-status --app-id=<app-id> --status=<app-status>

# Delete an app
scalebox app delete <app-name>
```

## 3.3 Task Management

### 3.3.1 Viewing Task Status

If app-id is not specified, the latest app-id is used by default
```bash
# List the app's tasks
scalebox task list --app-id <app-id>

# Filter tasks by status
scalebox task list --app-id <app-id> --status READY
scalebox task list --app-id <app-id> --status COMPLETED
scalebox task list --app-id <app-id> --status FAILED
```

### 3.3.2 Task Operations

```bash
# View task details
scalebox task show <task-id>

# Delete a task
scalebox task delete <task-id>
```

## 3.4 Cluster Management

### 3.4.1 Viewing Cluster Status

```bash
# Check cluster status
scalebox cluster status

# List cluster nodes
scalebox host list

# View node details
scalebox host show --host-id=<host-id>
```

### 3.4.2 Slot Management

```bash
# List slots
scalebox slot list

# View slot details
scalebox slot show --slot-id=<slot-id>
```

## 3.5 App Monitoring

### 3.5.1 Status Monitoring

```bash
# Monitor app status in real time
scalebox app status <app-name> --watch

# Monitor task execution
scalebox task monitor --app <app-name>
```

### 3.5.2 Performance Metrics

```bash
# View app performance
scalebox app metrics <app-name>

# View node resource usage
scalebox node metrics <node-name>
```

## 3.6 Troubleshooting

### 3.6.1 Common Problems

1. **App creation failed**
   - Check the configuration file format
   - Verify the parameter settings
   - Check the error logs

2. **Task execution failed**
   - Check the task logs
   - Verify the module image
   - Check resource availability

3. **Cluster connection problems**
   - Verify network connectivity
   - Check service status
   - Verify authentication configuration

### 3.6.2 Debugging Tools

```bash
# Validate the configuration
scalebox validate --app-file=<app-yaml>

# Test a module
scalebox test module <module-name>

# Check system status
scalebox system check
```

## 3.7 App Parsing and Integration

Apps are parsed and created through the `scalebox run` command:

```bash
scalebox run --env-file scalebox.env app.yaml
```

The priority order of template parameter variables during parsing (from high to low):

1. Environment variables from the command line execution
2. Environment variable file specified on the command line (default scalebox.env)
3. User-level environment variable configuration file: `${HOME}/.scalebox/environments`
4. System-level environment variable configuration file: `/etc/scalebox/environments`
