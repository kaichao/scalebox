IMAGE_NAME:=app-primes/builder:golang

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:golang app-primes/calc:golang
