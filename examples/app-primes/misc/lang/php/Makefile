IMAGE_NAME:=app-primes/builder:php

build:
	docker build --network=host -t $(IMAGE_NAME) .

clean:
	docker rmi -f app-primes/builder:php app-primes/calc:php
	docker rmi -f php:8-cli
