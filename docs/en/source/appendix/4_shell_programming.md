# 4. Shell Programming Interface

## 4.1 Module Script Writing Standards

## 4.2 Standard Input/Output Interface

## 4.3 File Exchange Interface Specification

## 4.4 Timestamps and Performance Statistics


Scripts inside modules are usually implemented in shell. To use the built-in functions, jq must be installed in the container to support JSON parsing.

The following is example code for installing jq in a Dockerfile of a debian/ubuntu-based image:
```Dockerfile
RUN apt update \
    && apt-get install -y jq \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
```

## 4.5 Common scalebox Built-In Functions

### 4.5.1 scalebox::json_val

Function: extract a parameter value from JSON

| Parameter | Description |
|------|------|
| $1 | JSON text |
| $2 | JSON field name |

Returns: the field value (string)

### 4.5.2 scalebox::task_header

Function: extract a value from the JSON message header; if absent, extract from the environment variable (the environment variable name is the uppercase form of the message header)

| Parameter | Description |
|------|------|
| $1 | JSON text |
| $2 | field name (letters, underscores) |

Returns: the parameter value (string)

### 4.5.3 scalebox::json_parse

Function: parse a JSON string into a bash associative array

| Parameter | Description |
|------|------|
| $1 | JSON text |
| $2 | associative array variable name (return value) |

### 4.5.4 scalebox::is_task_body

Function: determine whether the task body is simple text without null characters (non-JSON)

Returns: 0 = simple text, 1 = JSON

### 4.5.5 scalebox::append_to_file

Function: append content to a file (atomic write, with lock)

| Parameter | Description |
|------|------|
| $1 | file path |
| $2 | content |

### 4.5.6 scalebox::get_local_ip

Function: get the local IPv4 address

Returns: IP address string

### 4.5.7 scalebox::get_slot_seq

Function: get the current slot's sequence number (0-based)

Returns: integer

### 4.5.8 path::host_path

Function: map a path inside the container to a host path

| Parameter | Description |
|------|------|
| $1 | directory name inside the container |

Returns: the host directory accessible from the container

### 4.5.9 path::size

Function: get the byte size of a directory

| Parameter | Description |
|------|------|
| $1 | host directory accessible from the container |

Returns: byte count

### 4.5.10 path::is_dir

Function: determine whether a path is a directory

Returns: 0 = is a directory, non-zero = is not

## 4.6 Usage Examples of Built-In Functions

```bash
#!/usr/bin/env bash

source /usr/local/lib/scalebox/functions.sh

my_header=$(scalebox::task_header "$2" "my_header")

path_in_container="mypath"
host_dir=$(path::host_path ${path_in_container})

```

## 4.7 Data Directories Accessible Inside the Container

The external directories accessible inside the container by default include:
- /tmp: local temporary file directory
- /dev/shm: local cache directory (tmpfs)
- /cluster_data_root: cluster data root directory (defined by ```data_root``` in the cluster definition)
- /local_data_root: compute node local root directory

To access other directories, define the mapping relationships in the module definition's ```volumes```.
