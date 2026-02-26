IMAGE_NAME:=app-primes/builder:c

build:
	docker build --network=host -t $(IMAGE_NAME) .

run:
	docker run -it --rm -e OMP_NUM_THREADS=4 --entrypoint bash $(IMAGE_NAME)

local-run: local-build
	./primes 1000001 100000
	./primes-grouped 1000001 100000 1000

local-build:
	docker run -it --rm -v `pwd`:/src -w /src gcc:11-bullseye gcc -std=c99 -D CONTAINER_APP -Iinclude -fopenmp -O3 primes.c -o primes
	docker run -it --rm -v `pwd`:/src -w /src gcc:11-bullseye gcc -std=c99 -Iinclude -fopenmp -O3 primes_grouped.c -o primes-grouped
	# docker run -it --rm -v `cd ..;pwd`:/src -w /src/openmp gcc:11-bullseye gcc -std=c99 -Iinclude -fopenmp -O3 primes.c -o primes

clean:
	docker rmi -f app-primes/builder:c app-primes/calc:c
	docker rmi -f gcc:15-trixie
