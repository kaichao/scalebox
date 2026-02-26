IMAGE_NAME:=app-primes/builder:rust

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:rust app-primes/calc:rust
	docker rmi -f rust:1.93
