IMAGE_NAME:=app-primes/builder:python

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:python app-primes/calc:python
