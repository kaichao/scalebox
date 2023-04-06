# cron module

## Introduction
The cron module is a scalebox public module similar to the cron function under UNIX, which periodically sends messages to each sub-module to start related timing operations.

## Usage

Define timing operations for multiple modules through the file cron.txt.

In the cron.txt file, lines beginning with '#' are comment lines.

Each line defines the timing operation for a module, including timing operation interval definition and module name, separated by commas.

A sample cron.txt file is as follows:
```
# comments for cron.txt
@every 1m,mod0
@every 1m30s,mod1
```

The definition part of the timing operation interval refers to:
[cron-doc](https://pkg.go.dev/github.com/robfig/cron)

