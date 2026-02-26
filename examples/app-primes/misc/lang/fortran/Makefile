IMAGE_NAME:=app-primes/builder:fortran

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:fortran app-primes/calc:fortran
