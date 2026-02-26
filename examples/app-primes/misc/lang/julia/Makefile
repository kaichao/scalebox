IMAGE_NAME:=app-primes/builder:julia

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:julia app-primes/calc:julia
	docker rmi -f julia:1.12
