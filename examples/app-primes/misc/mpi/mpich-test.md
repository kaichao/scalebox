## Single-node test

- 6-core / 4-parallel
```sh
/gfsdata/primes/mpi/run-single-node-test.sh 6 4 10000000
/gfsdata/primes/mpi/run-single-node-test.sh 6 4 30000000
/gfsdata/primes/mpi/run-single-node-test.sh 6 4 100000000
```

- 1核24并运行
```sh
/gfsdata/primes/mpi/run-single-node-test.sh 1 24 10000000
/gfsdata/primes/mpi/run-single-node-test.sh 1 24 30000000
/gfsdata/primes/mpi/run-single-node-test.sh 1 24 100000000
```

## 16-node parallel

```sh
export OMP_NUM_THREADS=1
echo \n1c384p 10M
for i in {1..4}; do
  mpiexec -np 384 -f nodes-16.txt ./primes_mpi 10000000
  mpiexec -np 64 -f nodes.txt /gfsdata/app-primes/mpi/primes_mpi 10000000
done

echo \n1c384p 30M
for i in {1..4}; do
  mpiexec -np 384 -f nodes-16.txt ./primes_mpi 30000000 
done


export OMP_NUM_THREADS=6
echo \n6c64p 10M
for i in {1..4}; do
  mpiexec -np 64 -f nodes-16.txt ./primes_mpi 10000000 
done

echo \n6c64p 30M
for i in {1..4}; do
  mpiexec -np 64 -f nodes-16.txt ./primes_mpi 30000000 
done

```
## 12-node parallel
```sh
export OMP_NUM_THREADS=1
echo \n1c288p 100M
for i in {1..4}; do
  mpiexec -np 288 -f nodes-12.txt ./primes_mpi 10000000
done
for i in {1..4}; do
  mpiexec -np 288 -f nodes-12.txt ./primes_mpi 30000000
done
for i in {1..4}; do
  mpiexec -np 288 -f nodes-12.txt .primes_mpi 100000000
done

export OMP_NUM_THREADS=6
echo \n6c48p 100M
for i in {1..4}; do
  mpiexec -np 48 -f nodes-12.txt ./primes_mpi 10000000
done
for i in {1..4}; do
  mpiexec -np 48 -f nodes-12.txt ./primes_mpi 30000000
done
for i in {1..4}; do
  mpiexec -np 48 -f nodes-12.txt ./primes_mpi 100000000
done
```

## 8-node parallel
```sh
export OMP_NUM_THREADS=1
echo \n1c192p 10M
for i in {1..2}; do
  mpiexec -np 192 -f nodes-8-1.txt ./primes_mpi 10000000
done
echo \n1c192p 30M
for i in {1..2}; do
  mpiexec -np 192 -f nodes-8-1.txt ./primes_mpi 30000000
done
echo \n1c192p 100M
for i in {1..2}; do
  mpiexec -np 192 -f nodes-8-1.txt ./primes_mpi 100000000
done

export OMP_NUM_THREADS=6
echo \n6c32p 10M
for i in {1..2}; do
  mpiexec -np 32 -f nodes-8-1.txt ./primes_mpi 10000000
done
echo \n6c32p 30M
for i in {1..2}; do
  mpiexec -np 32 -f nodes-8-1.txt .primes_mpi 30000000
done
echo \n6c32p 100M
for i in {1..2}; do
  mpiexec -np 32 -f nodes-8-1.txt ./primes_mpi 100000000
done
```

## 4-node parallel
```sh
export OMP_NUM_THREADS=1
echo \n1c96p 10M
mpiexec -np 96 -f nodes-4-1.txt ./primes_mpi 10000000
echo \n1c96p 30M
mpiexec -np 96 -f nodes-4-1.txt ./primes_mpi 30000000
echo \n1c96p 100M
mpiexec -np 96 -f nodes-4-1.txt ./primes_mpi 100000000

export OMP_NUM_THREADS=6
echo \n6c16p 10M
mpiexec -np 16 -f nodes-4-1.txt ./primes_mpi 10000000
echo \n6c16p 30M
mpiexec -np 16 -f nodes-4-1.txt ./primes_mpi 30000000
echo \n6c16p 100M
mpiexec -np 16 -f nodes-4-1.txt ./primes_mpi 100000000
```
