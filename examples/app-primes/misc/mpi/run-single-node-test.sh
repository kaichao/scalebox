#!/bin/bash
h=$(hostname)
cd /gfsdata/app-primes/result/${h:1}
export OMP_NUM_THREADS=$1
# nohup mpiexec -np $2 /gfsdata/app-primes/mpi/primes_mpi $3 &

mpiexec -np $2 /gfsdata/app-primes/mpi/primes_mpi $3
# mpiexec -np $2 /gfsdata/app-primes/mpi/primes_mpi 10000000
# mpiexec -np $2 /gfsdata/app-primes/mpi/primes_mpi 30000000
# mpiexec -np $2 /gfsdata/app-primes/mpi/primes_mpi 100000000
