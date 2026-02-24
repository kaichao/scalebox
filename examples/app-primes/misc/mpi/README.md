## Build MPI program

- put source code on /gfsdata/app-primes
```sh
make build
```


## Run MPI program

### Single-node test
- user account : mpich
```sh
mkdir -p /gfsdata/app-primes/result && cd /gfsdata/app-primes/result

for i in {01..20};do mkdir $i;done

for i in 01 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 ; do echo $i;ssh r$i nohup /gfsdata/app-primes/mpi/run-single-node-test.sh 1 1 & ; done

for i in 05 ; do echo $i; ssh r$i nohup /gfsdata/app-primes/mpi/run-single-node-test.sh 1 1 10000000 & ; echo $i; done

```
### 4-node test
```sh
cd /gfsdata/app-primes/mpi
export OMP_NUM_THREADS=1
mpiexec -f hosts/4-nodes_1 -n 4 /gfsdata/app-primes/mpi/primes_mpi 1000000
```

### 16-node test
```sh
cd /gfsdata/app-primes/mpi
export OMP_NUM_THREADS=1
nohup mpiexec -f hosts/16-nodes_1p -n 16 /gfsdata/app-primes/mpi/primes_mpi 100000000 &
nohup mpiexec -f hosts/16-nodes_3p -n 48 /gfsdata/app-primes/mpi/primes_mpi 100000000 &
nohup mpiexec -f hosts/16-nodes_8p -n 128 /gfsdata/app-primes/mpi/primes_mpi 100000000 &
nohup mpiexec -f hosts/16-nodes_24p -n 384 /gfsdata/app-primes/mpi/primes_mpi 100000000 &
```

### Result

|  计算范围   | 串行方法  | OpenMP(2核) | OpenMP(4核) | OpenMP(8核) | OpenMP(16核) |
|  ----  | ----:|----:|----:|----:|----:|
