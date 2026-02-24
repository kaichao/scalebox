local-build:
	# docker run -it --rm -v `cd ..;pwd`:/src -w /src/ref-result gcc:11-bullseye gcc -std=c99 -I../include -fopenmp -lm primes-gen.c -o primes-gen
	docker run -it --rm -v `cd ..;pwd`:/src -w /src/ref-result gcc:11-bullseye gcc -std=c99 -I../include -fopenmp -lm qprimes.c -o qprimes
