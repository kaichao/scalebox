# singularity-based apps

## hello-scalebox
- build app
In directory hello-scalebox app, run command:
```sh
scalebox app create
```
- run slot-init command (from table t_job)
```sh
singularity run  --env JOB_NAME=hello-scalebox --env JOB_ID=1 --env GRPC_SERVER=10.100.1.30  --bind /tmp:/var/log/scalebox --bind /etc/localtime:/etc/localtime:ro singularity/scalebox/hello-scalebox.sif
```

## app-primes
- build app
In directory app-primes app, run command:
```sh
scalebox app create
```
