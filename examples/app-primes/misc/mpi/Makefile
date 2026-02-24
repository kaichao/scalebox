build:
	mpicc -std=c99 -I../src/c/include -fopenmp -O3 -lm primes_mpi.c -o primes_mpi

run: build
	mpiexec -np 20 ./primes_mpi 1000000
	# mpiexec -np 20 -f nodes.txt ./primes_mpi 1000000

hyprid-run:
	export OMP_NUM_THREADS=4; nohup mpiexec -np 6 ./primes_mpi 10000000 &

batch-run:
	for i in {01..02}; do echo $$i; ssh root@r$$i 'cd /gfsdata/primes/result/$$i; export OMP_NUM_THREADS=12; nohup mpiexec -np 2 ../../mpi/primes_mpi 10000000 &; sleep 10'; done
