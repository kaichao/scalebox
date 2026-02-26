IMAGE_NAME:=app-primes/builder:r

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:r app-primes/calc:r
	docker rmi -f r-base:4.5.2
